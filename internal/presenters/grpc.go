package presenters

import (
	"context"
	"multi-proving-client/internal/common"
	"multi-proving-client/internal/logic"
)

type API struct {
	service *logic.Service
}

func NewAPI(service *logic.Service) *API {
	return &API{
		service: service,
	}
}

func (a API) CalculateProof(ctx context.Context, req common.CalculateProofRequest) ([]byte, error) {
	//return a.service.CalculateProof(req)
	return nil, nil
}
