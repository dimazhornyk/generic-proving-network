package logic

import (
	"fmt"
	"github.com/pkg/errors"
	"multi-proving-client/internal/common"
)

type State struct {
	latestProofs        map[string][]byte
	provingRequestsByID map[common.RequestID]common.ProvingRequestMessage
}

func NewState() State {
	return State{
		latestProofs:        make(map[string][]byte),
		provingRequestsByID: make(map[string]common.ProvingRequestMessage),
	}
}

func (s *State) GetLatestProof(consumerName string) ([]byte, error) {
	if proof, ok := s.latestProofs[consumerName]; ok {
		return proof, nil
	}

	return nil, errors.New("unknown consumer")
}

func (s *State) GetDataByProvingRequestID(requestID common.RequestID) (common.ProvingRequestMessage, error) {
	data, ok := s.provingRequestsByID[requestID]
	if !ok {
		return common.ProvingRequestMessage{}, fmt.Errorf("unknown requestID: %s", requestID)
	}

	return data, nil
}

func (s *State) SaveRequest(data common.ProvingRequestMessage) error {
	s.provingRequestsByID[data.ID] = data

	return nil
}
