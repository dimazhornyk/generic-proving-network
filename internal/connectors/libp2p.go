package connectors

import (
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"multi-proving-client/internal/common"
	"os"
)

func NewHost(cfg *common.Config) (host.Host, error) {
	privBytes, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	priv, err := crypto.UnmarshalEd25519PrivateKey(privBytes)
	if err != nil {
		return nil, err
	}

	h, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", cfg.Port)),
		libp2p.Identity(priv),
	)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("host ID %s\n", host.ID().String())
	//fmt.Printf("following are the assigned addresses\n")
	//for _, addr := range host.Addrs() {
	//	fmt.Printf("%s\n", addr.String())
	//}
	//fmt.Printf("\n")

	return h, err
}
