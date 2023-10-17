package logic

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"log/slog"
	"multi-proving-client/internal/common"
	"time"
)

type ProvingRequestsHandler struct {
	state *State
}

func NewProvingRequestsHandler(state *State) *ProvingRequestsHandler {
	return &ProvingRequestsHandler{
		state: state,
	}
}

func (h *ProvingRequestsHandler) Handle(peerID peer.ID, msg common.ProvingRequestMessage) {
	if err := validateProvingRequest(msg); err != nil {
		// TODO: if it is invalid - take punishing actions
		slog.Error("invalid proving request", err)

		return
	}

	if err := h.state.SaveRequest(msg); err != nil {
		slog.Error("error saving proving request", err)
	}
}

func validateProvingRequest(msg common.ProvingRequestMessage) error {
	if msg.ID == "" {
		return errors.New("requestID is empty")
	}

	if msg.ConsumerName == "" {
		return errors.New("consumerName is empty")
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

	return verifySignature(msg.ConsumerName, msg.Signature)
}

func verifySignature(consumerName string, signature []byte) error {
	// TODO: implement
	return nil
}
