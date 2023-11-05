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
	host     host.Host
	nodesMap logic.NodesMap
	storage  *logic.Storage
	service  *logic.Service
	pubsub   *connectors.PubSub
	key      crypto.PrivKey
}

func NewProofsHandler(key crypto.PrivKey, host host.Host, storage *logic.Storage, service *logic.Service, pubsub *connectors.PubSub, nodesMap logic.NodesMap) *ProofsHandler {
	return &ProofsHandler{
		host:     host,
		storage:  storage,
		service:  service,
		pubsub:   pubsub,
		key:      key,
		nodesMap: nodesMap,
	}
}

func (h *ProofsHandler) Handle(ctx context.Context, peerID peer.ID, msg common.ProofSubmissionMessage) {
	if peerID == h.host.ID() {
		return // no need to verify the proof that we just generated
	}

	reqData, err := h.storage.GetProvingRequestByID(msg.RequestID)
	if err != nil {
		slog.Error("error getting proving request data", slog.String("err", err.Error()))

		return
	}

	if err := h.storage.AddProof(msg.RequestID, peerID, msg.ProofID, msg.Proof); err != nil {
		slog.Error("error adding proof to storage", slog.String("err", err.Error()))

		return
	}

	// TODO: check if the node's status was Proving

	valid, err := h.service.ValidateProof(msg.RequestID, reqData.ConsumerName, reqData.Data, msg.Proof)
	if err != nil {
		slog.Error("error validating proof", slog.String("err", err.Error()))

		return
	}

	votingPayload := common.ValidationPayload{
		RequestID: msg.RequestID,
		ProverID:  peerID,
		IsValid:   valid,
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
	votingPayload.PoolSize = h.nodesMap.CountNodesForConsumer(reqData.ConsumerName)
	votingPayload.ValidationTimestamp = time.Now().UnixNano()

	votingMsg := common.VotingMessage{
		Type:    common.VoteValidation,
		Payload: votingPayload,
	}

	if err := h.pubsub.Publish(ctx, common.VotingTopic, votingMsg); err != nil {
		slog.Error("error publishing validation message", slog.String("err", err.Error()))

		return
	}
}
