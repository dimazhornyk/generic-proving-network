package handlers

import (
	"context"
	"encoding/json"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/connectors"
	"github.com/dimazhornyk/generic-proving-network/internal/logic"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"log/slog"
	"time"
)

type ProofsHandler struct {
	host    host.Host
	state   *logic.State
	service *logic.ServiceStruct
	pubsub  *connectors.PubSub
	key     crypto.PrivKey
}

func NewProofsHandler(key crypto.PrivKey, host host.Host, state *logic.State, service *logic.ServiceStruct, pubsub *connectors.PubSub) *ProofsHandler {
	return &ProofsHandler{
		host:    host,
		state:   state,
		service: service,
		pubsub:  pubsub,
		key:     key,
	}
}

func (h *ProofsHandler) Handle(ctx context.Context, peerID peer.ID, msg common.ProofSubmissionMessage) {
	if peerID == h.host.ID() {
		return // no need to verify the proof that we just generated
	}

	reqData, err := h.state.GetDataByProvingRequestID(msg.RequestID)
	if err != nil {
		slog.Error("error getting proving request data", slog.String("err", err.Error()))

		return
	}

	valid, err := h.service.ValidateProof(msg.RequestID, reqData.ConsumerName, reqData.Data, msg.Proof)
	if err != nil {
		slog.Error("error validating proof", slog.String("err", err.Error()))

		return
	}

	votingPayload := common.ValidationPayload{
		RequestID:           msg.RequestID,
		ProofID:             msg.ProofID,
		IsValid:             valid,
		ValidationTimestamp: time.Now().UnixNano(),
	}

	b, err := json.Marshal(votingPayload)
	if err != nil {
		slog.Error("error marshalling validation message", slog.String("err", err.Error()))

		return
	}

	signature, err := h.key.Sign(b)
	if err != nil {
		slog.Error("error signing validation message", slog.String("err", err.Error()))

		return
	}
	votingPayload.Signature = signature

	votingMsg := common.VotingMessage{
		Type:    common.VoteValidation,
		Payload: votingPayload,
	}

	if err := h.pubsub.Publish(ctx, common.VotingTopic, votingMsg); err != nil {
		slog.Error("error publishing validation message", slog.String("err", err.Error()))

		return
	}
}
