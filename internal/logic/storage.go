package logic

import (
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"time"
)

var errUnknownRequest = errors.New("unknown request")

type RequestExtension struct {
	common.ProvingRequestMessage
	ProvingPeers []peer.ID
	proofs       map[peer.ID]common.ZKProof
}

type Storage struct {
	latestProofs    map[string]common.ZKProof
	provingRequests map[common.RequestID]RequestExtension
}

func NewStorage() *Storage {
	return &Storage{
		latestProofs:    make(map[string]common.ZKProof),
		provingRequests: make(map[string]RequestExtension),
	}
}

func (s *Storage) GetLatestProof(consumerName string) (common.ZKProof, error) {
	if proof, ok := s.latestProofs[consumerName]; ok {
		return proof, nil
	}

	return common.ZKProof{}, errors.New("unknown consumer")
}

func (s *Storage) GetProvingRequestByID(requestID common.RequestID) (RequestExtension, error) {
	data, ok := s.provingRequests[requestID]
	if !ok {
		return RequestExtension{}, errUnknownRequest
	}

	return data, nil
}

func (s *Storage) HasRequest(requestID common.RequestID) bool {
	_, ok := s.provingRequests[requestID]

	return ok
}

func (s *Storage) SaveRequest(data common.ProvingRequestMessage) error {
	s.provingRequests[data.ID] = RequestExtension{
		ProvingRequestMessage: data,
		ProvingPeers:          make([]peer.ID, 0),
		proofs:                make(map[peer.ID]common.ZKProof),
	}

	return nil
}

func (s *Storage) AddProvingPeer(requestID common.RequestID, peerID peer.ID) error {
	if !s.HasRequest(requestID) {
		return errUnknownRequest
	}

	req := s.provingRequests[requestID]
	req.ProvingPeers = append(req.ProvingPeers, peerID)
	s.provingRequests[requestID] = req

	return nil
}

func (s *Storage) AddProof(requestID common.RequestID, peerID peer.ID, proofID common.ProofID, proof []byte) error {
	if !s.HasRequest(requestID) {
		return errUnknownRequest
	}

	req := s.provingRequests[requestID]
	req.proofs[peerID] = common.ZKProof{
		ProofID:   proofID,
		Proof:     proof,
		Timestamp: time.Now().UnixNano(),
	}
	s.provingRequests[requestID] = req

	return nil
}

func (s *Storage) FinishProving(requestID common.RequestID) error {
	if !s.HasRequest(requestID) {
		return errUnknownRequest
	}

	provers := s.provingRequests[requestID].ProvingPeers
	if len(provers) == 0 {
		return errors.New("no provers")
	}

	proverID := provers[len(provers)-1]
	proof, ok := s.provingRequests[requestID].proofs[proverID]
	if !ok {
		return errors.New("no proof for the latest prover")
	}

	consumerName := s.provingRequests[requestID].ConsumerName
	s.latestProofs[consumerName] = proof
	delete(s.provingRequests, requestID)

	return nil
}
