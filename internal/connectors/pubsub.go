package connectors

import (
	"context"
	"encoding/json"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/pkg/errors"
)

type PubSub struct {
	ps            *pubsub.PubSub
	globalTopic   *pubsub.Topic
	requestsTopic *pubsub.Topic
	votingTopic   *pubsub.Topic
	proofsTopic   *pubsub.Topic
}

func NewPubSub(ctx context.Context, host host.Host) (*PubSub, error) {
	gossipSub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, errors.Wrap(err, "error creating a new gossip sub")
	}

	globalTopic, err := gossipSub.Join(common.GlobalTopic.String())
	if err != nil {
		return nil, errors.Wrap(err, "error joining a global topic")
	}

	requestsTopic, err := gossipSub.Join(common.RequestsTopic.String())
	if err != nil {
		return nil, errors.Wrap(err, "error joining a requests topic")
	}

	votingTopic, err := gossipSub.Join(common.VotingTopic.String())
	if err != nil {
		return nil, errors.Wrap(err, "error joining a voting topic")
	}

	proofsTopic, err := gossipSub.Join(common.ProofsTopic.String())
	if err != nil {
		return nil, errors.Wrap(err, "error joining a proofs topic")
	}

	return &PubSub{
		ps:            gossipSub,
		globalTopic:   globalTopic,
		requestsTopic: requestsTopic,
		votingTopic:   votingTopic,
		proofsTopic:   proofsTopic,
	}, nil
}

func (p *PubSub) SendStatusMessage(ctx context.Context, payload common.StatusMessage) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "error marshalling a status message")
	}

	if err := p.globalTopic.Publish(ctx, b); err != nil {
		return errors.Wrap(err, "error publishing a status message")
	}

	return nil
}

func (p *PubSub) Publish(ctx context.Context, topic common.Topic, message any) error {
	b, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "error marshalling a message")
	}

	t, err := p.selectTopic(topic)
	if err != nil {
		return errors.Wrap(err, "error selecting a topic")
	}

	return errors.Wrap(t.Publish(ctx, b), "error publishing a message")
}

func (p *PubSub) Subscribe(topic common.Topic) (*pubsub.Subscription, error) {
	t, err := p.selectTopic(topic)
	if err != nil {
		return nil, errors.Wrap(err, "error selecting a topic")
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
	case common.ProofsTopic:
		return p.proofsTopic, nil
	default:
		return nil, errors.New("unknown topic")
	}
}
