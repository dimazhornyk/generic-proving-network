package connectors

import (
	"context"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"math/big"
	"strings"
)

type Ethereum struct {
	contractAddress ethcommon.Address
	client          *ethclient.Client
}

func NewEthereum(cfg *common.Config) (*Ethereum, error) {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/8c2d91ab6e0d476f9ea87f5e19ea6fb7")
	if err != nil {
		return nil, err
	}

	return &Ethereum{
		client:          client,
		contractAddress: ethcommon.HexToAddress(cfg.ContractAddress),
	}, nil
}

func (e *Ethereum) GetAllConsumers(ctx context.Context) ([]common.Consumer, error) {
	// TODO: finish and test
	position := 0
	arraySlot := e.getArraySlot(position)
	vf := crypto.Keccak256Hash(arraySlot[:])

	offset := new(big.Int).SetInt64(int64(position - 1))
	v := new(big.Int).SetBytes(vf[:])
	v.Add(v, offset)

	slot := ethcommon.BytesToHash(v.Bytes())
	resp, err := e.client.StorageAt(ctx, e.contractAddress, slot, nil)
	if err != nil {
		return nil, err
	}

	hexutil.Encode(resp)
	res := make([]common.Consumer, 0)
	addresses := make([]string, 0) // TODO: get the addresses from a contract
	for _, addr := range addresses {
		// TODO: get image from map inside the contract
		var image = ""
		fragments := strings.Split(image, "/")
		if len(fragments) < 2 {
			return nil, errors.New("invalid image format")
		}

		res = append(res, common.Consumer{
			Image:   image,
			Address: addr,
			Name:    strings.Split(fragments[1], ":")[0],
		})
	}

	return res, nil
}

func (e *Ethereum) SubmitValidationSignatures(ctx context.Context, requestID common.RequestID, signatures []common.ValidationSignature, success bool) error {
	// TODO: implement
	return nil
}

func (e *Ethereum) getArraySlot(position int) [32]byte {
	return crypto.Keccak256Hash(
		ethcommon.LeftPadBytes(e.contractAddress[:], 32),
		ethcommon.LeftPadBytes(big.NewInt(int64(position)).Bytes(), 32),
	)
}
