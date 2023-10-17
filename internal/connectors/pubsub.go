package connectors

import (
	"context"
	"encoding/json"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/pkg/errors"
	"multi-proving-client/internal/common"
)

type PubSub struct {
	ps            *pubsub.PubSub
	globalTopic   *pubsub.Topic
	requestsTopic *pubsub.Topic
	votingTopic   *pubsub.Topic
}

func NewPubSub(ctx context.Context, host host.Host) (*PubSub, error) {
	gossipSub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, err
	}

	globalTopic, err := gossipSub.Join(common.GlobalTopic.String())
	if err != nil {
		return nil, err
	}

	requestsTopic, err := gossipSub.Join(common.RequestsTopic.String())
	if err != nil {
		return nil, err
	}

	votingTopic, err := gossipSub.Join(common.VotingTopic.String())
	if err != nil {
		return nil, err
	}

	return &PubSub{
		ps:            gossipSub,
		globalTopic:   globalTopic,
		requestsTopic: requestsTopic,
		votingTopic:   votingTopic,
	}, nil
}

func (p *PubSub) SendStatusMessage(ctx context.Context, payload common.StatusMessage) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.globalTopic.Publish(ctx, b)
}

func (p *PubSub) Publish(ctx context.Context, topic common.Topic, message any) error {
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	t, err := p.selectTopic(topic)
	if err != nil {
		return err
	}

	return t.Publish(ctx, b)
}

func (p *PubSub) Subscribe(topic common.Topic) (*pubsub.Subscription, error) {
	t, err := p.selectTopic(topic)
	if err != nil {
		return nil, err
	}

	return t.Subscribe()
}

func (p *PubSub) selectTopic(topic common.Topic) (*pubsub.Topic, error) {
	switch topic {
	case common.GlobalTopic:
		return p.globalTopic, nil
	case common.RequestsTopic:
		return p.requestsTopic, nil
	case common.VotingTopic:
		return p.votingTopic, nil
	default:
		return nil, errors.New("unknown topic")
	}
}
