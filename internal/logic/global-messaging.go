package logic

import (
	"context"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/connectors"
	"log/slog"
	"time"
)

type GlobalMessaging struct {
	pubsub *connectors.PubSub
}

func NewGlobalMessaging(pubsub *connectors.PubSub) (*GlobalMessaging, error) {
	return &GlobalMessaging{
		pubsub: pubsub,
	}, nil
}

func (gm *GlobalMessaging) Init(ctx context.Context, consumers []string) error {
	payload := common.StatusMessage{
		Status:  common.StatusInit,
		Payload: consumers,
	}

	if err := gm.pubsub.SendStatusMessage(ctx, payload); err != nil {
		return err
	}

	go gm.worker(ctx)

	return nil
}

func (gm *GlobalMessaging) worker(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 5)

	select {
	case <-ctx.Done():
		payload := common.StatusMessage{
			Status: common.StatusShuttingDown,
		}

		if err := gm.pubsub.SendStatusMessage(ctx, payload); err != nil {
			slog.Error("error on publishing status message", err)
		}

		return
	case <-ticker.C:
		payload := common.StatusMessage{
			Status: common.StatusIdle,
		}

		if err := gm.pubsub.SendStatusMessage(ctx, payload); err != nil {
			slog.Error("error on publishing status message", err)
		}
	}
}
