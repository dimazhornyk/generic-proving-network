package logic

import (
	"fmt"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
)

type RequestExtension struct {
	common.ProvingRequestMessage
	ProvingPeerID peer.ID
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

func (s *State) GetDataByProvingRequestID(requestID common.RequestID) (RequestExtension, error) {
	data, ok := s.provingRequestsByID[requestID]
	if !ok {
		return RequestExtension{}, fmt.Errorf("unknown requestID: %s", requestID)
	}

	return data, nil
}

func (s *State) SaveRequest(provingNode peer.ID, data common.ProvingRequestMessage) error {
	s.provingRequestsByID[data.ID] = RequestExtension{ProvingRequestMessage: data, ProvingPeerID: provingNode}

	return nil
}
