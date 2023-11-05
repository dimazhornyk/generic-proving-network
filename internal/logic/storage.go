package logic

import (
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var errUnknownRequest = errors.New("unknown request")

type RequestExtension struct {
	common.ProvingRequestMessage
	ProvingPeers         []peer.ID
	proofs               map[peer.ID]common.ZKProof
	validationSignatures []common.ValidationSignature
}

type Storage struct {
	latestProofs    map[string]common.ZKProof
	provingRequests map[common.RequestID]RequestExtension
	mu              sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		latestProofs:    make(map[string]common.ZKProof),
		provingRequests: make(map[string]RequestExtension),
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

func (s *Storage) GetProvingRequestByID(requestID common.RequestID) (RequestExtension, error) {
	s.mu.RLock()
	data, ok := s.provingRequests[requestID]
	s.mu.RUnlock()

	if !ok {
		return RequestExtension{}, errUnknownRequest
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
	s.provingRequests[data.ID] = RequestExtension{
		ProvingRequestMessage: data,
		ProvingPeers:          make([]peer.ID, 0),
		proofs:                make(map[peer.ID]common.ZKProof),
		validationSignatures:  make([]common.ValidationSignature, 0),
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

	req.proofs[peerID] = common.ZKProof{
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
	proof, ok := req.proofs[proverID]
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

	req.validationSignatures = append(req.validationSignatures, common.ValidationSignature{
		PeerID:    peerID,
		Signature: signature,
	})
	s.mu.Lock()
	s.provingRequests[requestID] = req
	s.mu.Unlock()

	return nil
}

func (s *Storage) GetValidationSignatures(requestID common.RequestID) ([]common.ValidationSignature, error) {
	if !s.HasRequest(requestID) {
		return nil, errUnknownRequest
	}

	s.mu.RLock()
	req := s.provingRequests[requestID]
	s.mu.RUnlock()

	return req.validationSignatures, nil
}
