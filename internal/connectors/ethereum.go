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
)

type Ethereum struct {
	address   ethcommon.Address
	client    *gpn.ProvingNetwork
	ethClient *ethclient.Client
}

func NewEthereum(cfg *common.Config, privateKey *ecdsa.PrivateKey) (*Ethereum, error) {
	ethClient, err := ethclient.Dial(cfg.EthereumAPI)
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
		})
	}

	return result, nil
}

func (e *Ethereum) GetAllProvers(ctx context.Context) ([]ethcommon.Address, error) {
	opts := &bind.CallOpts{
		Context: ctx,
		From:    e.address,
	}

	return e.client.GetProvers(opts)
}

func (e *Ethereum) ListenForNewProvers(ctx context.Context) (<-chan *gpn.ProvingNetworkProverUpdate, error) {
	opts := &bind.WatchOpts{
		Context: ctx,
	}

	ch := make(chan *gpn.ProvingNetworkProverUpdate)
	sub, err := e.client.WatchProverUpdate(opts, ch)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		sub.Unsubscribe()
	}()

	return ch, nil
}

func (e *Ethereum) ListenForNewConsumers(ctx context.Context) (<-chan *gpn.ProvingNetworkConsumerUpdate, error) {
	opts := &bind.WatchOpts{
		Context: ctx,
	}

	ch := make(chan *gpn.ProvingNetworkConsumerUpdate)
	sub, err := e.client.WatchConsumerUpdate(opts, ch)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		sub.Unsubscribe()
	}()

	return ch, nil
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

	tx, err := e.client.SubmitSignedProof(opts, request.ID, request.Reward, rs, ss, vs)
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
