package presenters

import (
	"context"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/logic"
	"github.com/dimazhornyk/generic-proving-network/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type API struct {
	proto.UnimplementedProvingNetworkServiceServer
	service *logic.Service
}

func NewAPI(service *logic.Service) *API {
	return &API{
		service: service,
	}
}

func (a *API) ComputeProof(ctx context.Context, req *proto.ComputeProofRequest) (*emptypb.Empty, error) {
	r := toCommonRequest(req)

	return &emptypb.Empty{}, a.service.InitiateProofCalculation(ctx, r)
}

func (a *API) GetProof(_ context.Context, req *proto.GetProofRequest) (*proto.GetProofResponse, error) {
	proof, err := a.service.GetProof(req.GetRequestId())
	if err != nil {
		if errors.Is(err, logic.ErrNoProof) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.GetProofResponse{
		ProofId:   proof.ProofID,
		Proof:     proof.Proof,
		Timestamp: proof.Timestamp,
	}, nil
}

func toCommonRequest(req *proto.ComputeProofRequest) common.ComputeProofRequest {
	return common.ComputeProofRequest{
		ID:              req.GetRequestId(),
		ConsumerName:    req.GetConsumerName(),
		ConsumerAddress: req.GetConsumerAddress(),
		Data:            req.GetData(),
		Signature:       req.GetSignature(),
	}
}
