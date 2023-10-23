package handlers

import (
	"context"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/connectors"
	"github.com/dimazhornyk/generic-proving-network/internal/logic"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"log/slog"
)

type ProofsHandler struct {
	host    host.Host
	state   *logic.State
	service *logic.ServiceStruct
	pubsub  *connectors.PubSub
}

func NewProofsHandler(host host.Host, state *logic.State, service *logic.ServiceStruct, pubsub *connectors.PubSub) *ProofsHandler {
	return &ProofsHandler{
		host:    host,
		state:   state,
		service: service,
		pubsub:  pubsub,
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

	// TODO: what to do with validation result? send to the network?
}
