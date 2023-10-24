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

const VoteDuration = time.Second * 4

type VotingHandler struct {
	host              host.Host
	state             *logic.State
	service           *logic.ServiceStruct
	pubsub            *connectors.PubSub
	selectionVotings  logic.VotingMap[common.RequestID, peer.ID]
	validationVotings logic.VotingMap[common.RequestID, bool]
}

func NewVotingHandler(host host.Host, service *logic.ServiceStruct, state *logic.State, pubsub *connectors.PubSub) *VotingHandler {
	return &VotingHandler{
		host:              host,
		service:           service,
		state:             state,
		pubsub:            pubsub,
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

	if !h.state.HasRequest(payload.RequestID) {
		return errors.New("unknown requestID")
	}

	h.selectionVotings.Add(payload.RequestID, voterID, payload.PeerID)

	return nil
}

func (h *VotingHandler) handleValidationVoting(voterID peer.ID, message common.VotingMessage) error {
	payload, ok := message.Payload.(common.ValidationPayload)
	if !ok {
		return errors.New("invalid payload type for VoteValidation")
	}

	if !h.state.HasRequest(payload.RequestID) {
		return errors.New("unknown requestID")
	}

	h.validationVotings.Add(payload.RequestID, voterID, payload.IsValid)

	return nil
}

//func (h *VotingHandler) finalizeVoting(request common.ProvingRequestMessage) {
//	time.Sleep(VoteDuration)
//	winner, err := h.selectionVotings.GetWinner(request.ID)
//	if err != nil {
//		slog.Error("error getting winner", slog.String("err", err.Error()))
//
//		return
//	}
//}
