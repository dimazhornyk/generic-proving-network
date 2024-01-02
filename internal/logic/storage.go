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

// inmemory storage is a temporary solution, it should be replaced with a more persistent storage
type Storage struct {
	resultsStorage  map[common.RequestID]common.ZKProof // TODO: implement properly, use disk
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

func (s *Storage) GetLatestProof(consumerName string) *common.ZKProof {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if proof, ok := s.latestProofs[consumerName]; ok {
		return &proof
	}

	return nil
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
	defer s.mu.RUnlock()

	_, ok := s.provingRequests[requestID]

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
		ValidationSignatures:  make(map[peer.ID]map[peer.ID][]byte),
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
	s.resultsStorage[requestID] = proof
	s.latestProofs[req.ConsumerName] = proof
	delete(s.provingRequests, requestID)
	s.mu.Unlock()

	return nil
}

func (s *Storage) GetFromResultsStorage(request common.RequestID) (common.ZKProof, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	proof, ok := s.resultsStorage[request]
	if !ok {
		return common.ZKProof{}, errors.New("no proof in the results storage")
	}

	return proof, nil
}

func (s *Storage) AddValidationSignature(requestID common.RequestID, voterID, proverID peer.ID, signature []byte) error {
	if !s.HasRequest(requestID) {
		return errUnknownRequest
	}

	s.mu.Lock()
	req := s.provingRequests[requestID]
	if _, ok := req.ValidationSignatures[proverID]; !ok {
		req.ValidationSignatures[proverID] = make(map[peer.ID][]byte)
	}

	req.ValidationSignatures[proverID][voterID] = signature
	s.mu.Unlock()

	return nil
}

func (s *Storage) GetValidationSignatures(requestID common.RequestID, proverID peer.ID) ([][]byte, error) {
	if !s.HasRequest(requestID) {
		return nil, errUnknownRequest
	}

	res := make([][]byte, 0)
	s.mu.RLock()
	defer s.mu.RUnlock()

	req := s.provingRequests[requestID]
	for _, signature := range req.ValidationSignatures[proverID] {
		res = append(res, signature)
	}

	return res, nil
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
