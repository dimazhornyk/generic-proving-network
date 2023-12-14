package logic

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var errUnknownRequest = errors.New("unknown request")

type Storage struct {
	latestProofs    map[string]common.ZKProof
	provingRequests map[common.RequestID]common.RequestExtension
	mu              sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		latestProofs:    make(map[string]common.ZKProof),
		provingRequests: make(map[string]common.RequestExtension),
	}
}

func (s *Storage) GetLatestProof(consumerName string) (common.ZKProof, error) {
	s.mu.RLock()
	proof, ok := s.latestProofs[consumerName]
	s.mu.RUnlock()

	if ok {
		return proof, nil
	}

	return common.ZKProof{}, errors.New("unknown consumer")
}

func (s *Storage) GetProvingRequestByID(requestID common.RequestID) (common.RequestExtension, error) {
	s.mu.RLock()
	data, ok := s.provingRequests[requestID]
	s.mu.RUnlock()

	if !ok {
		return common.RequestExtension{}, errUnknownRequest
	}

	return data, nil
}

func (s *Storage) HasRequest(requestID common.RequestID) bool {
	s.mu.RLock()
	_, ok := s.provingRequests[requestID]
	s.mu.RUnlock()

	return ok
}

func (s *Storage) SaveRequest(data common.ProvingRequestMessage) error {
	if s.HasRequest(data.ID) {
		return errors.New("request already exists")
	}

	s.mu.Lock()
	s.provingRequests[data.ID] = common.RequestExtension{
		ProvingRequestMessage: data,
		ProvingPeers:          make([]peer.ID, 0),
		Proofs:                make(map[peer.ID]common.ZKProof),
		ValidationSignatures:  make([][]byte, 0),
	}
	s.mu.Unlock()

	return nil
}

func (s *Storage) AddProvingPeer(requestID common.RequestID, peerID peer.ID) error {
	if !s.HasRequest(requestID) {
		return errUnknownRequest
	}

	s.mu.RLock()
	req := s.provingRequests[requestID]
	s.mu.RUnlock()

	req.ProvingPeers = append(req.ProvingPeers, peerID)
	s.mu.Lock()
	s.provingRequests[requestID] = req
	s.mu.Unlock()

	return nil
}

func (s *Storage) AddProof(requestID common.RequestID, peerID peer.ID, proofID common.ProofID, proof []byte) error {
	if !s.HasRequest(requestID) {
		return errUnknownRequest
	}

	s.mu.RLock()
	req := s.provingRequests[requestID]
	s.mu.RUnlock()

	req.Proofs[peerID] = common.ZKProof{
		ProofID:   proofID,
		Proof:     proof,
		Timestamp: time.Now().UnixNano(),
	}
	s.mu.Lock()
	s.provingRequests[requestID] = req
	s.mu.Unlock()

	return nil
}

func (s *Storage) DeleteProvingRequest(requestID common.RequestID) error {
	if !s.HasRequest(requestID) {
		return errUnknownRequest
	}

	s.mu.RLock()
	req := s.provingRequests[requestID]
	s.mu.RUnlock()
	if len(req.ProvingPeers) == 0 {
		return errors.New("no provers")
	}

	proverID := req.ProvingPeers[len(req.ProvingPeers)-1]
	s.mu.RLock()
	proof, ok := req.Proofs[proverID]
	s.mu.RUnlock()
	if !ok {
		return errors.New("no proof for the latest prover")
	}

	s.mu.Lock()
	s.latestProofs[req.ConsumerName] = proof
	delete(s.provingRequests, requestID)
	s.mu.Unlock()

	return nil
}

func (s *Storage) AddValidationSignature(requestID common.RequestID, peerID peer.ID, signature []byte) error {
	if !s.HasRequest(requestID) {
		return errUnknownRequest
	}

	s.mu.RLock()
	req := s.provingRequests[requestID]
	s.mu.RUnlock()

	req.ValidationSignatures = append(req.ValidationSignatures, signature)
	s.mu.Lock()
	s.provingRequests[requestID] = req
	s.mu.Unlock()

	return nil
}

func (s *Storage) GetValidationSignatures(requestID common.RequestID) ([][]byte, error) {
	if !s.HasRequest(requestID) {
		return nil, errUnknownRequest
	}

	s.mu.RLock()
	req := s.provingRequests[requestID]
	s.mu.RUnlock()

	return req.ValidationSignatures, nil
}

func (s *Storage) GetRequests() map[common.RequestID]common.RequestExtension {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.provingRequests
}

func (s *Storage) GetLatestProofs() map[string]common.ZKProof {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.latestProofs
}

func (s *Storage) SetRequests(requests map[common.RequestID]common.RequestExtension) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.provingRequests = requests
}

func (s *Storage) SetLatestProofs(proofs map[string]common.ZKProof) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.latestProofs = proofs
}

func (s *Storage) GetStorageHash() (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hasher := sha256.New()
	b, err := common.GobEncodeMessage(s.provingRequests)
	if err != nil {
		return "", errors.Wrap(err, "error encoding proving requests")
	}

	hasher.Write(b)
	b, err = common.GobEncodeMessage(s.latestProofs)
	if err != nil {
		return "", errors.Wrap(err, "error encoding latest proofs")
	}

	hasher.Write(b)

	return base64.URLEncoding.EncodeToString(hasher.Sum(nil)), nil
}
