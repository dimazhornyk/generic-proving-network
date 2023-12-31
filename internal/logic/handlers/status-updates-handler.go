package handlers

import (
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/logic"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"log/slog"
)

type StatusUpdatesHandler struct {
	nodes logic.StatusMap
}

func NewStatusUpdatesHandler(nodes logic.StatusMap) *StatusUpdatesHandler {
	return &StatusUpdatesHandler{
		nodes: nodes,
	}
}

func (h *StatusUpdatesHandler) Handle(peerID peer.ID, msg common.StatusMessage) {
	slog.Info("handling status message", slog.String("status", msg.Status.String()), slog.String("peerID", peerID.String()))

	var err error
	switch msg.Status {
	case common.StatusInit:
		err = h.handleInit(peerID, msg)
	case common.StatusIdle:
		err = h.handleIdle(peerID)
	case common.StatusShuttingDown:
		err = h.handleShuttingDown(peerID)
	case common.StatusProving:
		err = h.handleProving(peerID)
	}

	if err != nil {
		slog.Error("error handling status message",
			slog.String("err", err.Error()),
			slog.String("status", msg.Status.String()),
		)
	}
}

func (h *StatusUpdatesHandler) handleInit(peerID peer.ID, msg common.StatusMessage) error {
	consumers, ok := msg.Payload.([]string)
	if !ok {
		return errors.New("invalid payload type for StatusInit")
	}

	h.nodes.Add(peerID, common.StatusInit, consumers)

	return nil
}

func (h *StatusUpdatesHandler) handleIdle(peerID peer.ID) error {
	if err := h.nodes.UpdateStatus(peerID, common.StatusIdle); err != nil {
		return err
	}

	return nil
}

func (h *StatusUpdatesHandler) handleShuttingDown(peerID peer.ID) error {
	if err := h.nodes.UpdateStatus(peerID, common.StatusShuttingDown); err != nil {
		return err
	}

	return nil
}

func (h *StatusUpdatesHandler) handleProving(peerID peer.ID) error {
	if err := h.nodes.UpdateStatus(peerID, common.StatusProving); err != nil {
		return err
	}

	return nil
}
