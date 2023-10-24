package connectors

import (
	"fmt"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/pkg/errors"
	"os"
)

func NewPrivateKey(cfg *common.Config) (crypto.PrivKey, error) {
	privBytes, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading private key")
	}

	priv, err := crypto.UnmarshalEd25519PrivateKey(privBytes)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling private key")
	}

	return priv, nil
}

//nolint:ireturn
func NewHost(cfg *common.Config, priv crypto.PrivKey) (host.Host, error) {
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", cfg.Port)),
		libp2p.Identity(priv),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating a new host")
	}

	// fmt.Printf("host ID %s\n", host.ID().String())
	// fmt.Printf("following are the assigned addresses\n")
	// for _, addr := range host.Addrs() {
	// 	fmt.Printf("%s\n", addr.String())
	// }
	// fmt.Printf("\n")

	return h, nil
}
