package logic

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/connectors"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"time"
)

const proveURL = "http://localhost:%s/prove"
const validateURL = "http://localhost:%s/validate"

var ErrNoProof = errors.New("no proof found")

type Service struct {
	docker    *connectors.Docker
	pubsub    *connectors.PubSub
	host      host.Host
	nodes     StatusMap
	storage   *Storage
	status    *StatusSharing
	consumers []common.Consumer
}

func NewService(ctx context.Context, cfg *common.Config, d *connectors.Docker, pubsub *connectors.PubSub, nodes StatusMap, storage *Storage, status *StatusSharing, host host.Host, eth *connectors.Ethereum) (*Service, error) {
	var consumers []common.Consumer

	if cfg.Mode == common.TestingMode {
		consumers = []common.Consumer{{
			Name:  cfg.Consumers[0],
			Image: cfg.Consumers[0],
		}}
	} else {
		allConsumers, err := eth.GetAllConsumers(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "error getting consumers from ethereum")
		}

		consumers = common.Filter(allConsumers, func(consumer common.Consumer) bool {
			return slices.Contains(cfg.Consumers, consumer.Name)
		})
	}

	if len(consumers) == 0 {
		return nil, errors.New("no consumers found")
	}

	return &Service{
		docker:    d,
		pubsub:    pubsub,
		nodes:     nodes,
		storage:   storage,
		status:    status,
		host:      host,
		consumers: consumers,
	}, nil
}

func (s *Service) Start() error {
	images := common.Map(s.consumers, func(c common.Consumer) string {
		return c.Image
	})

	if err := s.docker.StartContainers(images); err != nil {
		return errors.Wrap(err, "error starting containers")
	}

	return nil
}

func (s *Service) InitiateProofCalculation(ctx context.Context, req common.ComputeProofRequest) error {
	msg := common.ProvingRequestMessage{
		ID:              req.ID,
		ConsumerName:    req.ConsumerName,
		ConsumerAddress: req.ConsumerAddress,
		Signature:       req.Signature,
		Data:            req.Data,
		Timestamp:       time.Now().UnixNano(),
	}

	// todo: check that consumer is in a list in contract, verify signature
	slog.Info("new request", slog.String("requestID", req.ID), slog.String("consumer", req.ConsumerName))
	if err := s.pubsub.Publish(ctx, common.RequestsTopic, msg); err != nil {
		return errors.Wrap(err, "error publishing the proving request")
	}

	return nil
}

func (s *Service) GetProof(requestID common.RequestID) (common.ZKProof, error) {
	proof, err := s.storage.GetFromResultsStorage(requestID)
	if err != nil {
		slog.Warn("no proof in storage", slog.String("requestID", requestID))

		return common.ZKProof{}, ErrNoProof
	}

	return proof, nil
}

func (s *Service) HandleProverSelection(ctx context.Context, msg common.ProvingRequestMessage, excludedPeers ...peer.ID) error {
	proverID, err := s.selectProvingNode(msg.ConsumerName, msg.Timestamp, excludedPeers...)
	if err != nil {
		return errors.Wrap(err, "error selecting prover")
	}

	slog.Info("selected prover",
		slog.String("nodeID", proverID.String()),
		slog.String("requestID", msg.ID),
	)

	if err := s.voteProverSelection(ctx, msg.ID, proverID); err != nil {
		return errors.Wrap(err, "error voting for prover selection")
	}

	if proverID == s.host.ID() {
		slog.Info("I am the selected node, starting proving...")

		proof, err := s.computeProof(ctx, msg)
		if err != nil {
			return errors.Wrap(err, "error computing the proof")
		}

		if err := s.submitProof(msg.ID, proof); err != nil {
			return errors.Wrap(err, "error submitting the proof")
		}
	}

	return nil
}

func (s *Service) submitProof(requestID common.RequestID, proof []byte) error {
	msg := common.ProofSubmissionMessage{
		RequestID: requestID,
		ProofID:   uuid.New().String(),
		Proof:     proof,
	}

	if err := s.pubsub.Publish(context.Background(), common.ProofsTopic, msg); err != nil {
		return errors.Wrap(err, "error publishing the proof")
	}

	return nil
}

func (s *Service) voteProverSelection(ctx context.Context, requestID common.RequestID, provingNodeID peer.ID) error {
	msg := common.VotingMessage{
		Type: common.VoteProverSelection,
		Payload: common.ProverSelectionPayload{
			RequestID: requestID,
			PeerID:    provingNodeID,
		},
	}

	if err := s.pubsub.Publish(ctx, common.VotingTopic, msg); err != nil {
		return errors.Wrap(err, "error when publishing to voting topic")
	}

	return nil
}

func (s *Service) selectProvingNode(consumerName string, requestTimestamp int64, excludeList ...peer.ID) (peer.ID, error) {
	nodes := make([]common.NodeData, 0)
	for _, node := range s.nodes {
		fmt.Println("NODE ID", node.PeerID, "COMMITMENTS", node.Commitments, "STATUS", node.Status, "AVAILABLE SINCE", node.AvailableSince, "REQUEST TIMESTAMP", requestTimestamp, "CONSUMER_NAME", consumerName)
		// is committed to the consumer, is idle, went up earlier than request was sent, is not in the exclude list
		if slices.Contains(node.Commitments, consumerName) && isNodeAppropriate(node, requestTimestamp) && !slices.Contains(excludeList, node.PeerID) {
			nodes = append(nodes, node)
		}
	}
	slices.SortStableFunc(nodes, func(a, b common.NodeData) int {
		return cmp.Compare(a.PeerID.String(), b.PeerID.String())
	})

	var seed []byte
	latestProof := s.storage.GetLatestProof(consumerName)
	if latestProof == nil {
		seed = []byte(consumerName)
	} else {
		seed = latestProof.Proof
	}

	random, err := common.BytesToRandom(seed)
	if err != nil {
		return "", errors.Wrap(err, "error getting random from last proof")
	}

	idx := random.Intn(len(nodes))

	return nodes[idx].PeerID, nil
}

func (s *Service) computeProof(ctx context.Context, req common.ProvingRequestMessage) ([]byte, error) {
	s.status.SetStatus(ctx, common.StatusProving)
	defer s.status.SetStatus(ctx, common.StatusIdle)

	image := s.getConsumerImage(req.ConsumerName)
	if image == "" {
		return nil, errors.New("unknown consumer")
	}

	port, err := s.docker.GetContainerPort(image)
	if err != nil {
		return nil, errors.Wrap(err, "error getting container ID")
	}

	msg := common.ProvingMessage{
		RequestID: req.ID,
		Data:      req.Data,
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling prover message")
	}

	resp, err := http.Post(fmt.Sprintf(proveURL, port), "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, errors.Wrap(err, "error requesting the prover's container")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading the response body")
	}

	var response common.ProvingResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling response")
	}

	return response.Proof, nil
}

func (s *Service) ValidateProof(requestID common.RequestID, consumer string, data, proof []byte) (bool, error) {
	image := s.getConsumerImage(consumer)
	if image == "" {
		return false, errors.New("unknown consumer")
	}

	port, err := s.docker.GetContainerPort(image)
	if err != nil {
		return false, errors.Wrap(err, "error getting container ID")
	}

	msg := common.ValidationProverMessage{
		RequestID: requestID,
		Proof:     proof,
		Data:      data,
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return false, errors.Wrap(err, "error marshalling validation message")
	}

	resp, err := http.Post(fmt.Sprintf(validateURL, port), "application/json", bytes.NewReader(b))
	if err != nil {
		return false, errors.Wrap(err, "error requesting the prover's container")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, errors.Wrap(err, "error reading the response body")
	}

	var response common.ValidationResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return false, errors.Wrap(err, "error unmarshalling response")
	}

	return response.Valid, nil
}

func (s *Service) getConsumerImage(consumerName string) string {
	for _, consumer := range s.consumers {
		if consumer.Name == consumerName {
			return consumer.Image
		}
	}

	return ""
}

func isNodeAppropriate(node common.NodeData, maxTimestamp int64) bool {
	return node.Status == common.StatusIdle && node.AvailableSince < maxTimestamp
}
