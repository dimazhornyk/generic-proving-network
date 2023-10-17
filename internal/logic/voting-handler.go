package logic

import (
	"context"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"log/slog"
	"multi-proving-client/internal/common"
	"multi-proving-client/internal/connectors"
	"time"
)

const VoteDuration = time.Second * 4

type VotingHandler struct {
	host              host.Host
	state             *State
	service           *Service
	pubsub            *connectors.PubSub
	selectionVotings  VotingMap[common.RequestID, peer.ID]
	validationVotings VotingMap[common.RequestID, bool]
}

func NewVotingHandler(host host.Host, service *Service, state *State, pubsub *connectors.PubSub) *VotingHandler {
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
	case common.InitProverSelectionVoting:
		err = h.selectProver(ctx, peerID, msg)
	case common.VoteProverSelection:
		err = h.handleSelectionVoting(peerID, msg)
	case common.InitValidationVoting:
	case common.VoteValidation:
	}

	if err != nil {
		slog.Error("error handling voting message", err)
	}
}

func (h *VotingHandler) handleSelectionVoting(voterID peer.ID, message common.VotingMessage) error {
	msg, ok := message.Payload.(common.ProverSelectionMessage)
	if !ok {
		return errors.New("invalid payload type for VoteProverSelection")
	}

	return h.selectionVotings.Add(msg.RequestID, voterID, msg.PeerID)
}

func (h *VotingHandler) selectProver(ctx context.Context, creatorID peer.ID, msg common.VotingMessage) error {
	requestID, ok := msg.Payload.(common.RequestID)
	if !ok {
		return errors.New("invalid payload type for InitProverSelectionVoting")
	}

	h.selectionVotings.Create(requestID, creatorID)
	data, err := h.state.GetDataByProvingRequestID(requestID)
	if err != nil {
		return errors.Wrap(err, "error when getting data from storage")
	}

	peerID, err := h.service.SelectProvingNode(data.ConsumerName, data.Timestamp)
	if err != nil {
		return errors.Wrap(err, "error selecting proving node")
	}

	msg = common.VotingMessage{
		Type: common.VoteProverSelection,
		Payload: common.ProverSelectionMessage{
			RequestID: requestID,
			PeerID:    peerID,
		},
	}

	if err := h.pubsub.Publish(ctx, common.VotingTopic, msg); err != nil {
		return errors.Wrap(err, "error when publishing to voting topic")
	}

	go func() {

	}()

	return nil
}

func (h *VotingHandler) finalizeVoting(request common.ProvingRequestMessage) {
	time.Sleep(VoteDuration)
	winner, err := h.selectionVotings.GetWinner(request.ID)
	if err != nil {
		slog.Error("error getting winner", slog.String("err", err.Error()))

		return
	}

	slog.Info("proving selection winner",
		slog.String("winner", winner.String()),
		slog.String("requestID", request.ID),
	)

	if h.host.ID() == winner {
		slog.Info("I am the winner, starting proving")

		if err := h.service.GenerateProof(request); err != nil {
			slog.Error("error starting proving", slog.String("err", err.Error()))
		}
	}
}
