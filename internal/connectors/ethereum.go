package connectors

import (
	"context"
	"crypto/ecdsa"
	gpn "github.com/dimazhornyk/generic-proving-network/internal/abi"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"log/slog"
	"strings"
)

type Ethereum struct {
	address   ethcommon.Address
	client    *gpn.ProvingNetwork
	ethClient *ethclient.Client
}

func NewEthereum(cfg *common.Config, privateKey *ecdsa.PrivateKey) (*Ethereum, error) {
	ethClient, err := ethclient.Dial("https://sepolia.infura.io/v3/8c2d91ab6e0d476f9ea87f5e19ea6fb7")
	if err != nil {
		return nil, err
	}

	contractAddr := ethcommon.HexToAddress(cfg.ContractAddress)
	client, err := gpn.NewProvingNetwork(contractAddr, ethClient)
	if err != nil {
		return nil, err
	}

	return &Ethereum{
		address:   ethCrypto.PubkeyToAddress(privateKey.PublicKey),
		client:    client,
		ethClient: ethClient,
	}, nil
}

func (e *Ethereum) GetAllConsumers(ctx context.Context) ([]common.Consumer, error) {
	opts := &bind.CallOpts{
		Context: ctx,
		From:    e.address,
	}

	consumers, err := e.client.GetConsumers(opts)
	if err != nil {
		return nil, err
	}

	var result []common.Consumer
	for _, consumer := range consumers {
		result = append(result, common.Consumer{
			Image:   consumer.ContainerName,
			Address: consumer.Addr,
			Balance: consumer.Balance,
			Name:    strings.Split(consumer.ContainerName, ":")[0], // remove image tag
		})
	}

	return result, nil
}

func (e *Ethereum) GetAllProvers(ctx context.Context) ([]gpn.NetworkProverView, error) {
	opts := &bind.CallOpts{
		Context: ctx,
		From:    e.address,
	}

	return e.client.GetProvers(opts)
}

func (e *Ethereum) SubmitValidationSignatures(ctx context.Context, request common.ProvingRequestMessage, signatures [][]byte) error {
	if len(signatures) == 0 {
		return errors.New("no signatures provided")
	}

	opts := &bind.TransactOpts{
		Context: ctx,
		From:    e.address,
	}

	var err error
	rs := make([][32]byte, len(signatures)+1)
	ss := make([][32]byte, len(signatures)+1)
	vs := make([]uint8, len(signatures)+1)

	// consumer's signature is a first element in the signatures array
	rs[0], ss[0], vs[0], err = common.GetRSV(request.Signature)
	if err != nil {
		return errors.Wrap(err, "error getting RSV of the consumer's signature")
	}

	for i, signature := range signatures {
		rs[i+1], ss[i+1], vs[i+1], err = common.GetRSV(signature)
		if err != nil {
			return errors.Wrap(err, "error getting RSV of the validator's signature")
		}
	}

	tx, err := e.client.SubmitSignedProof(opts, request.ID, &request.Reward, rs, ss, vs)
	if err != nil {
		return errors.Wrap(err, "error submitting signed proof")
	}

	receipt, err := bind.WaitMined(ctx, e.ethClient, tx)
	if err != nil {
		return errors.Wrap(err, "error waiting for the transaction to be mined")
	}

	slog.Info("Transaction mined", "tx", receipt.TxHash.String(), "status", receipt.Status)

	return nil
}
