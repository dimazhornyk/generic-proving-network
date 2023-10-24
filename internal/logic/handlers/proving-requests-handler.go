package handlers

import (
	"context"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/connectors"
	"github.com/dimazhornyk/generic-proving-network/internal/logic"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"log/slog"
	"time"
)

type ProvingRequestsHandler struct {
	host    host.Host
	state   *logic.State
	service *logic.ServiceStruct
	pubsub  *connectors.PubSub
}

func NewProvingRequestsHandler(host host.Host, state *logic.State, service *logic.ServiceStruct, pubsub *connectors.PubSub) *ProvingRequestsHandler {
	return &ProvingRequestsHandler{
		host:    host,
		state:   state,
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

	proverID, err := h.selectProver(msg)
	if err != nil {
		slog.Error("error selecting prover", slog.String("err", err.Error()))

		return
	}
	slog.Info("proving selection winner",
		slog.String("winner", proverID.String()),
		slog.String("requestID", msg.ID),
	)

	if err := h.state.SaveRequest(proverID, msg); err != nil {
		slog.Error("error saving proving request", slog.String("err", err.Error()))

		return
	}

	if err := h.voteProverSelection(ctx, msg.ID, proverID); err != nil {
		slog.Error("error voting for prover selection", slog.String("err", err.Error()))

		return
	}

	if proverID == h.host.ID() {
		slog.Info("I am the selected node, starting proving...")

		proof, err := h.service.ComputeProof(msg)
		if err != nil {
			slog.Error("error computing the proof", slog.String("err", err.Error()))

			return
		}

		if err := h.submitProof(msg.ID, proof); err != nil {
			slog.Error("error submitting proof", slog.String("err", err.Error()))
		}
	}
}

func (h *ProvingRequestsHandler) submitProof(requestID common.RequestID, proof []byte) error {
	msg := common.ProofSubmissionMessage{
		RequestID: requestID,
		ProofID:   uuid.New().String(),
		Proof:     proof,
	}

	if err := h.pubsub.Publish(context.Background(), common.ProofsTopic, msg); err != nil {
		return errors.Wrap(err, "error publishing the proof")
	}

	return nil
}

func (h *ProvingRequestsHandler) voteProverSelection(ctx context.Context, requestID common.RequestID, provingNodeID peer.ID) error {
	msg := common.VotingMessage{
		Type: common.VoteProverSelection,
		Payload: common.ProverSelectionPayload{
			RequestID: requestID,
			PeerID:    provingNodeID,
		},
	}

	if err := h.pubsub.Publish(ctx, common.VotingTopic, msg); err != nil {
		return errors.Wrap(err, "error when publishing to voting topic")
	}

	return nil
}

func (h *ProvingRequestsHandler) selectProver(req common.ProvingRequestMessage) (peer.ID, error) {
	peerID, err := h.service.SelectProvingNode(req.ConsumerName, req.Timestamp)
	if err != nil {
		return "", errors.Wrap(err, "error selecting proving node")
	}

	return peerID, nil
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
