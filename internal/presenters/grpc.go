package presenters

import (
	"context"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/logic"
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
