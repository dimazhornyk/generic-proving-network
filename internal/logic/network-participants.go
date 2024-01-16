package logic

import (
	"context"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/connectors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"sync"
)

type NetworkParticipants struct {
	sync.Mutex

	eth       *connectors.Ethereum
	provers   map[ethcommon.Address]struct{}
	consumers map[ethcommon.Address]common.Consumer
}

func NewNetworkParticipants(ctx context.Context, eth *connectors.Ethereum) (*NetworkParticipants, error) {
	np := &NetworkParticipants{
		eth:       eth,
		provers:   make(map[ethcommon.Address]struct{}),
		consumers: make(map[ethcommon.Address]common.Consumer),
	}

	eg := errgroup.Group{}
	eg.Go(func() error {
		return np.GetConsumers(ctx)
	})
	eg.Go(func() error {
		return np.GetProvers(ctx)
	})

	return np, eg.Wait()
}

func (np *NetworkParticipants) GetConsumers(ctx context.Context) error {
	consumers, err := np.eth.GetAllConsumers(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting consumers from ethereum")
	}

	np.Lock()
	for _, consumer := range consumers {
		np.consumers[consumer.Address] = consumer
	}
	np.Unlock()

	ch, err := np.eth.ListenForNewConsumers(ctx)
	if err != nil {
		return errors.Wrap(err, "error listening for new consumers")
	}

	go func() {
		for {
			select {
			case msg := <-ch:
				np.Lock()
				if msg.IsAdded {
					np.consumers[msg.Addr] = common.Consumer{
						Address: msg.Addr,
						Balance: msg.Balance,
						Image:   msg.ContainerName,
					}
				} else {
					delete(np.consumers, msg.Addr)
				}
				np.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (np *NetworkParticipants) GetProvers(ctx context.Context) error {
	addrs, err := np.eth.GetAllProvers(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting provers from ethereum")
	}

	np.Lock()
	for _, addr := range addrs {
		np.provers[addr] = struct{}{}
	}
	np.Unlock()

	ch, err := np.eth.ListenForNewProvers(ctx)
	if err != nil {
		return errors.Wrap(err, "error listening for new provers")
	}

	go func() {
		for {
			select {
			case msg := <-ch:
				np.Lock()
				if msg.IsAdded {
					np.provers[msg.Addr] = struct{}{}
				} else {
					delete(np.provers, msg.Addr)
				}
				np.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (np *NetworkParticipants) IsKnownProver(addr ethcommon.Address) bool {
	np.Lock()
	defer np.Unlock()

	_, ok := np.provers[addr]

	return ok
}

func (np *NetworkParticipants) IsKnownConsumer(addr ethcommon.Address) bool {
	np.Lock()
	defer np.Unlock()

	_, ok := np.consumers[addr]

	return ok
}

func (np *NetworkParticipants) GetAllConsumers() []common.Consumer {
	np.Lock()
	defer np.Unlock()

	var result []common.Consumer
	for _, consumer := range np.consumers {
		result = append(result, consumer)
	}

	return result
}
