package presenters

import (
	"context"
	"encoding/json"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/connectors"
	"github.com/dimazhornyk/generic-proving-network/internal/logic/handlers"
	"github.com/pkg/errors"
	"log/slog"
)

type Listener struct {
	pubsub               *connectors.PubSub
	votingHandler        *handlers.VotingHandler
	requestsHandler      *handlers.ProvingRequestsHandler
	statusUpdatesHandler *handlers.StatusUpdatesHandler
	proofsHandler        *handlers.ProofsHandler
}

func NewListener(pubsub *connectors.PubSub, vh *handlers.VotingHandler, rh *handlers.ProvingRequestsHandler, sh *handlers.StatusUpdatesHandler) *Listener {
	return &Listener{
		pubsub:               pubsub,
		votingHandler:        vh,
		requestsHandler:      rh,
		statusUpdatesHandler: sh,
	}
}

func (l *Listener) Listen(ctx context.Context) error {
	funcs := []func(context.Context) error{
		l.ListenStateUpdates,
		l.ListenProvingRequests,
		l.ListenVoting,
		l.ListenProofs,
	}

	errs := make(chan error, len(funcs))
	for _, f := range funcs {
		function := f
		go func(function func(context.Context) error) {
			errs <- function(ctx)
		}(function)
	}

	cnt := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errs:
			if err != nil {
				return err
			} else {
				cnt++
				if cnt == len(funcs) {
					return nil
				}
			}
		}
	}
}

func (l *Listener) ListenStateUpdates(ctx context.Context) error {
	subscription, err := l.pubsub.Subscribe(common.GlobalTopic)
	if err != nil {
		return errors.Wrap(err, "error subscribing to state updates topic")
	}

	for {
		pubsubMsg, err := subscription.Next(ctx)
		if err != nil {
			slog.Error("error getting next message from subscription", err)

			continue
		}

		var msg common.StatusMessage
		if err := json.Unmarshal(pubsubMsg.Data, &msg); err != nil {
			slog.Error("error unmarshalling state update message", err)

			continue
		}

		go l.statusUpdatesHandler.Handle(pubsubMsg.ReceivedFrom, msg)
	}
}

func (l *Listener) ListenProvingRequests(ctx context.Context) error {
	subscription, err := l.pubsub.Subscribe(common.RequestsTopic)
	if err != nil {
		return errors.Wrap(err, "error subscribing to requests topic")
	}

	for {
		pubsubMsg, err := subscription.Next(ctx)
		if err != nil {
			slog.Error("error getting next message from subscription", err)

			continue
		}

		var msg common.ProvingRequestMessage
		if err := json.Unmarshal(pubsubMsg.Data, &msg); err != nil {
			slog.Error("error unmarshalling proving request message", err)

			continue
		}

		go l.requestsHandler.Handle(ctx, msg)
	}
}

func (l *Listener) ListenProofs(ctx context.Context) error {
	subscription, err := l.pubsub.Subscribe(common.ProofsTopic)
	if err != nil {
		return errors.Wrap(err, "error subscribing to voting topic")
	}

	for {
		pubsubMsg, err := subscription.Next(ctx)
		if err != nil {
			slog.Error("error getting next message from subscription", err)

			continue
		}

		var msg common.ProofSubmissionMessage
		if err := json.Unmarshal(pubsubMsg.Data, &msg); err != nil {
			slog.Error("error unmarshalling voting message", err)

			continue
		}

		go l.proofsHandler.Handle(ctx, pubsubMsg.ReceivedFrom, msg)
	}
}

func (l *Listener) ListenVoting(ctx context.Context) error {
	subscription, err := l.pubsub.Subscribe(common.VotingTopic)
	if err != nil {
		return errors.Wrap(err, "error subscribing to voting topic")
	}

	for {
		pubsubMsg, err := subscription.Next(ctx)
		if err != nil {
			slog.Error("error getting next message from subscription", err)

			continue
		}

		var msg common.VotingMessage
		if err := json.Unmarshal(pubsubMsg.Data, &msg); err != nil {
			slog.Error("error unmarshalling voting message", err)

			continue
		}

		go l.votingHandler.Handle(ctx, pubsubMsg.ReceivedFrom, msg)
	}
}
