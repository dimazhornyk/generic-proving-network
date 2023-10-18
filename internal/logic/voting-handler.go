package logic

import (
	"context"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/connectors"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"log/slog"
	"time"
)

const VoteDuration = time.Second * 4

type VotingHandler struct {
	host              host.Host
	state             *State
	service           *ServiceStruct
	pubsub            *connectors.PubSub
	selectionVotings  VotingMap[common.RequestID, peer.ID]
	validationVotings VotingMap[common.RequestID, bool]
}

func NewVotingHandler(host host.Host, service *ServiceStruct, state *State, pubsub *connectors.PubSub) *VotingHandler {
	return &VotingHandler{
		host:              host,
		service:           service,
		state:             state,
		pubsub:            pubsub,
		selectionVotings:  make(VotingMap[common.RequestID, peer.ID]),
		validationVotings: make(VotingMap[common.RequestID, bool]),
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
	//h.selectionVotings.Create(requestID, creatorID)

	msg, ok := message.Payload.(common.ProverSelectionMessage)
	if !ok {
		return errors.New("invalid payload type for VoteProverSelection")
	}

	return h.selectionVotings.Add(msg.RequestID, voterID, msg.PeerID)
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
