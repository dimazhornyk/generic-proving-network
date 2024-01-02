package logic

import (
	"context"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/connectors"
	"log/slog"
	"sync"
	"time"
)

type StatusSharing struct {
	pubsub *connectors.PubSub
	status common.Status
	mu     sync.Mutex
}

func NewGlobalMessaging(pubsub *connectors.PubSub) (*StatusSharing, error) {
	return &StatusSharing{
		pubsub: pubsub,
		status: common.StatusIdle,
	}, nil
}

func (s *StatusSharing) Init(ctx context.Context, consumers []string) error {
	payload := common.StatusMessage{
		Status:  common.StatusInit,
		Payload: consumers,
	}

	if err := s.pubsub.SendStatusMessage(ctx, payload); err != nil {
		return err
	}

	go s.worker(ctx)

	return nil
}

func (s *StatusSharing) SetStatus(ctx context.Context, status common.Status) {
	s.mu.Lock()
	s.status = status
	s.mu.Unlock()

	s.shareStatus(ctx)
}

func (s *StatusSharing) worker(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-ctx.Done():
			s.SetStatus(context.Background(), common.StatusShuttingDown)
			return
		case <-ticker.C:
			s.shareStatus(ctx)
		}
	}
}

func (s *StatusSharing) shareStatus(ctx context.Context) {
	s.mu.Lock()
	status := s.status
	s.mu.Unlock()

	payload := common.StatusMessage{
		Status: status,
	}

	if err := s.pubsub.SendStatusMessage(ctx, payload); err != nil {
		slog.Error("error on publishing status message", err)
	}
}
