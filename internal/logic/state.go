package logic

import (
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
)

var errUnknownRequest = errors.New("unknown request")

type RequestExtension struct {
	common.ProvingRequestMessage
	ProvingPeers []peer.ID
}

type State struct {
	latestProofs        map[string][]byte
	provingRequestsByID map[common.RequestID]RequestExtension
}

func NewState() State {
	return State{
		latestProofs:        make(map[string][]byte),
		provingRequestsByID: make(map[string]RequestExtension),
	}
}

func (s *State) GetLatestProof(consumerName string) ([]byte, error) {
	if proof, ok := s.latestProofs[consumerName]; ok {
		return proof, nil
	}

	return nil, errors.New("unknown consumer")
}

func (s *State) GetProvingRequestByID(requestID common.RequestID) (RequestExtension, error) {
	data, ok := s.provingRequestsByID[requestID]
	if !ok {
		return RequestExtension{}, errUnknownRequest
	}

	return data, nil
}

func (s *State) HasRequest(requestID common.RequestID) bool {
	_, ok := s.provingRequestsByID[requestID]

	return ok
}

func (s *State) SaveRequest(data common.ProvingRequestMessage) error {
	s.provingRequestsByID[data.ID] = RequestExtension{
		ProvingRequestMessage: data,
		ProvingPeers:          make([]peer.ID, 0),
	}

	return nil
}

func (s *State) AddProvingPeer(requestID common.RequestID, peerID peer.ID) error {
	if !s.HasRequest(requestID) {
		return errUnknownRequest
	}

	req := s.provingRequestsByID[requestID]
	req.ProvingPeers = append(req.ProvingPeers, peerID)
	s.provingRequestsByID[requestID] = req

	return nil
}
