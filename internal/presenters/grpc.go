package presenters

import (
	"context"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/logic"
)

type API struct {
	service *logic.ServiceStruct
}

func NewAPI(service *logic.ServiceStruct) *API {
	return &API{
		service: service,
	}
}

func (a API) CalculateProof(ctx context.Context, req common.CalculateProofRequest) ([]byte, error) {
	//return a.service.CalculateProof(req)
	return nil, nil
}
