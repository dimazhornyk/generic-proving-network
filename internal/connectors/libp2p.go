package connectors

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/pkg/errors"
	"log/slog"
	"os"
	"strings"
)

//nolint:ireturn
func NewPrivateKey(cfg *common.Config) (*ecdsa.PrivateKey, error) {
	privBytes, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading private key")
	}

	trimmed := strings.TrimSpace(string(privBytes))
	priv, err := ethCrypto.HexToECDSA(trimmed)

	return priv, nil
}

//nolint:ireturn
func NewHost(cfg *common.Config, privECDSA *ecdsa.PrivateKey) (host.Host, error) {
	b := ethCrypto.FromECDSA(privECDSA)
	priv, err := crypto.UnmarshalSecp256k1PrivateKey(b)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling private key")
	}

	h, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", cfg.Port)),
		libp2p.Identity(priv),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating a new host")
	}

	slog.Info("libp2p host created", slog.String("hostID", h.ID().String()))

	return h, nil
}
