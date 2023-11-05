package handlers

import (
	"context"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/connectors"
	"github.com/dimazhornyk/generic-proving-network/internal/logic"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"log/slog"
	"time"
)

const maxProvingAttempts = 3
const DoubleCheckInterval = time.Second * 2
const SelectionVotingDuration = time.Second * 4
const ValidationVotingDuration = time.Second * 6

type VotingHandler struct {
	host              host.Host
	storage           *logic.Storage
	service           *logic.Service
	pubsub            *connectors.PubSub
	ethereum          *connectors.Ethereum
	selectionVotings  logic.VotingMap[common.RequestID, peer.ID]
	validationVotings logic.VotingMap[common.RequestID, bool]
}

func NewVotingHandler(host host.Host, service *logic.Service, storage *logic.Storage, pubsub *connectors.PubSub, eth *connectors.Ethereum) *VotingHandler {
	return &VotingHandler{
		host:              host,
		service:           service,
		storage:           storage,
		pubsub:            pubsub,
		ethereum:          eth,
		selectionVotings:  make(logic.VotingMap[common.RequestID, peer.ID]),
		validationVotings: make(logic.VotingMap[common.RequestID, bool]),
	}
}

func (h *VotingHandler) Handle(ctx context.Context, peerID peer.ID, msg common.VotingMessage) {
	var err error
	switch msg.Type {
	case common.VoteProverSelection:
		err = h.handleSelectionVoting(peerID, msg)
	case common.VoteValidation:
		err = h.handleValidationVoting(ctx, peerID, msg)
	}

	if err != nil {
		slog.Error("error handling voting message", err)
	}
}

func (h *VotingHandler) handleSelectionVoting(voterID peer.ID, message common.VotingMessage) error {
	payload, ok := message.Payload.(common.ProverSelectionPayload)
	if !ok {
		return errors.New("invalid payload type for VoteProverSelection")
	}

	if !h.storage.HasRequest(payload.RequestID) {
		slog.Info("unknown requestID, double checking after timeout", slog.String("requestID", payload.RequestID))
		time.Sleep(DoubleCheckInterval)
		if !h.storage.HasRequest(payload.RequestID) {
			return errors.New("unknown requestID")
		}
	}

	if !h.selectionVotings.Add(payload.RequestID, voterID, payload.PeerID) {
		time.Sleep(SelectionVotingDuration)
		winner, err := h.selectionVotings.GetWinner(payload.RequestID) // TODO: handle draw and empty voting
		if err != nil {
			return errors.Wrap(err, "error getting winner")
		}

		h.selectionVotings.Delete(payload.RequestID)
		if winner == nil {
			return errors.New("no selection votes")
		}

		if err := h.storage.AddProvingPeer(payload.RequestID, *winner); err != nil {
			return errors.Wrap(err, "error adding proving peer")
		}
	}

	return nil
}

func (h *VotingHandler) handleValidationVoting(ctx context.Context, voterID peer.ID, message common.VotingMessage) error {
	payload, ok := message.Payload.(common.ValidationPayload)
	if !ok {
		return errors.New("invalid payload type for VoteValidation")
	}

	if !h.storage.HasRequest(payload.RequestID) {
		return errors.New("unknown requestID")
	}

	if err := h.storage.AddValidationSignature(payload.RequestID, voterID, payload.Signature); err != nil {
		return errors.Wrap(err, "error adding validation signature")
	}

	if !h.validationVotings.Add(payload.RequestID, voterID, payload.IsValid) {
		time.Sleep(ValidationVotingDuration)
		isProofValid, err := h.validationVotings.GetWinner(payload.RequestID)
		if err != nil {
			return errors.Wrap(err, "error getting winner")
		}

		h.validationVotings.Delete(payload.RequestID)
		if isProofValid == nil {
			return errors.New("no validation votes")
		}

		if !*isProofValid {
			return h.handleInvalidProof(ctx, payload.RequestID)
		}

		signatures, err := h.storage.GetValidationSignatures(payload.RequestID)
		if err != nil {
			return errors.Wrap(err, "error getting validation signatures")
		}

		// TODO: optimize by batching the signatures
		if err := h.ethereum.SubmitValidationSignatures(ctx, payload.RequestID, signatures, true); err != nil {
			return errors.Wrap(err, "error submitting validation signatures")
		}

		if err := h.storage.DeleteProvingRequest(payload.RequestID); err != nil {
			return errors.Wrap(err, "error finishing proving")
		}
	}

	return nil
}

func (h *VotingHandler) handleInvalidProof(ctx context.Context, requestID common.RequestID) error {
	req, err := h.storage.GetProvingRequestByID(requestID)
	if err != nil {
		return errors.Wrap(err, "error getting proving request")
	}

	// TODO: punish the node that has submitted the proof

	if len(req.ProvingPeers) < maxProvingAttempts {
		return h.service.HandleProverSelection(ctx, req.ProvingRequestMessage, req.ProvingPeers...)
	}

	// TODO: collect signatures for all (3) invalid proofs and send them to the contract
	// TODO: delete request from storage

	return nil
}
