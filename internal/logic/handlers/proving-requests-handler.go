package handlers

import (
	"context"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/connectors"
	"github.com/dimazhornyk/generic-proving-network/internal/logic"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/pkg/errors"
	"log/slog"
	"time"
)

type ProvingRequestsHandler struct {
	host    host.Host
	storage *logic.Storage
	service *logic.Service
	pubsub  *connectors.PubSub
}

func NewProvingRequestsHandler(host host.Host, storage *logic.Storage, service *logic.Service, pubsub *connectors.PubSub) *ProvingRequestsHandler {
	return &ProvingRequestsHandler{
		host:    host,
		storage: storage,
		service: service,
		pubsub:  pubsub,
	}
}

func (h *ProvingRequestsHandler) Handle(ctx context.Context, msg common.ProvingRequestMessage) {
	if err := validateProvingRequest(msg); err != nil {
		// TODO: if it is invalid - take punishing actions
		slog.Error("invalid proving request", slog.String("err", err.Error()))

		return
	}

	if err := h.storage.SaveRequest(msg); err != nil {
		slog.Error("error saving proving request", slog.String("err", err.Error()))

		return
	}

	if err := h.service.HandleProverSelection(ctx, msg); err != nil {
		slog.Error("error handling prover selection", slog.String("err", err.Error()))

		return
	}
}

func validateProvingRequest(msg common.ProvingRequestMessage) error {
	if msg.ID == "" {
		return errors.New("requestID is empty")
	}

	if msg.ConsumerImage == "" {
		return errors.New("consumerImage is empty")
	}

	if msg.ConsumerAddress == "" {
		return errors.New("consumerAddress is empty")
	}

	if len(msg.Signature) == 0 {
		return errors.New("signature is empty")
	}

	if len(msg.Data) == 0 {
		return errors.New("data is empty")
	}

	t := time.Unix(0, msg.Timestamp)
	if t.After(time.Now()) {
		return errors.New("timestamp is in the future")
	}

	if time.Now().Sub(t) > time.Hour {
		return errors.New("timestamp is too old")
	}

	return verifySignature(msg.ConsumerImage, msg.Signature)
}

func verifySignature(consumerImage string, signature []byte) error {
	// TODO: implement
	return nil
}
