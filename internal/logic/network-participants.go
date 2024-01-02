package logic

import (
	"context"
	"github.com/dimazhornyk/generic-proving-network/internal/connectors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"sync"
)

type NetworkParticipants struct {
	sync.Mutex

	Addresses map[ethcommon.Address]struct{}
}

func NewNetworkParticipants(ctx context.Context, eth *connectors.Ethereum) (*NetworkParticipants, error) {
	addrs, err := eth.GetAllProvers(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error getting provers from ethereum")
	}

	np := &NetworkParticipants{
		Addresses: make(map[ethcommon.Address]struct{}),
	}

	for _, addr := range addrs {
		np.Addresses[addr] = struct{}{}
	}

	ch, err := eth.ListenForNewProvers(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error listening for new provers")
	}

	go func() {
		for {
			select {
			case msg := <-ch:
				np.Lock()
				if msg.IsAdded {
					np.Addresses[msg.Addr] = struct{}{}
				} else {
					delete(np.Addresses, msg.Addr)
				}
				np.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()

	return np, nil
}

func (np *NetworkParticipants) IsNetworkParticipant(addr ethcommon.Address) bool {
	np.Lock()
	defer np.Unlock()

	_, ok := np.Addresses[addr]

	return ok
}
