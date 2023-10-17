package logic

import (
	"cmp"
	"context"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"multi-proving-client/internal/common"
	"multi-proving-client/internal/connectors"
	"slices"
	"time"
)

type Service struct {
	docker    *connectors.Docker
	pubsub    *connectors.PubSub
	host      host.Host
	nodes     NodesMap
	state     State
	consumers []common.Consumer
}

func NewService(cfg common.Config, docker *connectors.Docker, pubsub *connectors.PubSub, nodes NodesMap, state State, host host.Host) *Service {
	return &Service{
		docker:    docker,
		pubsub:    pubsub,
		nodes:     nodes,
		state:     state,
		host:      host,
		consumers: common.GetConsumers(cfg.Consumers),
	}
}

func (s Service) Start(ctx context.Context) error {
	images := common.Map(s.consumers, func(c common.Consumer) string {
		return c.Image
	})

	if err := s.preloadImages(images); err != nil {
		return err
	}

	return nil
}

func (s Service) InitiateProofCalculation(req common.CalculateProofRequest) ([]byte, error) {
	msg := common.ProvingRequestMessage{
		ID:              uuid.New().String(),
		ConsumerName:    req.ConsumerName,
		ConsumerAddress: req.ConsumerAddress,
		Signature:       req.Signature,
		Data:            req.Data,
		Timestamp:       time.Now().UnixNano(),
	}

	if err := s.pubsub.Publish(context.Background(), common.RequestsTopic, msg); err != nil {
		return nil, errors.Wrap(err, "error publishing the proving request")
	}

	// TODO think about delaying this so everyone has saved the request data
	if err := s.StartSelectionVote(msg.ID); err != nil {
		return nil, err
	}

	// TODO: wait till the proving end and respond (or make API async so client would poll for response)
	return nil, nil
}

func (s Service) StartSelectionVote(requestID common.RequestID) error {
	msg := common.VotingMessage{
		Type:    common.InitProverSelectionVoting,
		Payload: requestID,
	}

	if err := s.pubsub.Publish(context.Background(), common.VotingTopic, msg); err != nil {
		return errors.Wrap(err, "error publishing the message to start voting")
	}

	return nil
}

func (s Service) SelectProvingNode(consumerName string, requestTimestamp int64) (peer.ID, error) {
	nodes := make([]common.NodeData, 0)
	for _, node := range s.nodes {
		if slices.Contains(node.Commitments, consumerName) && isNodeAppropriate(node, requestTimestamp) {
			nodes = append(nodes, node)
		}
	}
	slices.SortStableFunc(nodes, func(a, b common.NodeData) int {
		return cmp.Compare(a.PeerID.String(), b.PeerID.String())
	})

	lastProof, err := s.state.GetLatestProof(consumerName)
	if err != nil {
		return "", errors.Wrap(err, "error getting latest proof")
	}

	random, err := common.ZKPToRandom(lastProof)
	if err != nil {
		return "", errors.Wrap(err, "error getting random from last proof")
	}

	idx := random.Intn(len(nodes))

	return nodes[idx].PeerID, nil
}

func isNodeAppropriate(node common.NodeData, maxTimestamp int64) bool {
	return node.Status == common.StatusIdle && node.AvailableSince < maxTimestamp
}

func (s Service) preloadImages(images []string) error {
	for _, img := range images {
		if err := s.docker.Pull(img); err != nil {
			return errors.Wrap(err, "error preloading docker images")
		}
	}

	return nil
}

func (s Service) GenerateProof(request common.ProvingRequestMessage) error {

	return nil
}
