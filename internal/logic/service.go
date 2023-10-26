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

// TODO: rename to service after GoLand is fixed
type ServiceStruct struct {
	docker    *connectors.Docker
	pubsub    *connectors.PubSub
	host      host.Host
	nodes     NodesMap
	state     State
	consumers []common.Consumer
}

func NewService(cfg common.Config, d *connectors.Docker, pubsub *connectors.PubSub, nodes NodesMap, state State, host host.Host) *ServiceStruct {
	return &ServiceStruct{
		docker:    d,
		pubsub:    pubsub,
		nodes:     nodes,
		state:     state,
		host:      host,
		consumers: common.GetConsumers(cfg.Consumers),
	}
}

func (s *ServiceStruct) Start() error {
	images := common.Map(s.consumers, func(c common.Consumer) string {
		return c.Image
	})

	if err := s.docker.StartContainers(images); err != nil {
		return errors.Wrap(err, "error starting containers")
	}

	return nil
}

func (s *ServiceStruct) InitiateProofCalculation(req common.CalculateProofRequest) ([]byte, error) {
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

	// TODO: wait till the proving finalization and respond (or make API async so client would poll for response)
	return nil, nil
}

func (s *ServiceStruct) HandleProverSelection(ctx context.Context, msg common.ProvingRequestMessage, excludedPeers ...peer.ID) error {
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

		proof, err := s.computeProof(msg)
		if err != nil {
			return errors.Wrap(err, "error computing the proof")
		}

		if err := s.submitProof(msg.ID, proof); err != nil {
			return errors.Wrap(err, "error submitting the proof")
		}
	}

	return nil
}

func (s *ServiceStruct) submitProof(requestID common.RequestID, proof []byte) error {
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

func (s *ServiceStruct) voteProverSelection(ctx context.Context, requestID common.RequestID, provingNodeID peer.ID) error {
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

func (s *ServiceStruct) selectProvingNode(consumerName string, requestTimestamp int64, excludeList ...peer.ID) (peer.ID, error) {
	nodes := make([]common.NodeData, 0)
	for _, node := range s.nodes {
		// is committed to the consumer, is idle, went up earlier than request was sent, is not in the exclude list
		if slices.Contains(node.Commitments, consumerName) && isNodeAppropriate(node, requestTimestamp) && !slices.Contains(excludeList, node.PeerID) {
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

func (s *ServiceStruct) computeProof(req common.ProvingRequestMessage) ([]byte, error) {
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

func (s *ServiceStruct) ValidateProof(requestID common.RequestID, consumer string, data, proof []byte) (bool, error) {
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

func (s *ServiceStruct) getConsumerImage(consumerName string) string {
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
