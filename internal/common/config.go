package common

import (
	"github.com/caarlos0/env"
	"github.com/libp2p/go-libp2p/core"
	"github.com/pkg/errors"
)

type Config struct {
	EthereumAPI     string          `env:"ETHEREUM_API,required"`
	ProtocolID      core.ProtocolID `env:"PROTOCOL_ID" envDefault:"/p2p/gpn-node-te/1.0.0"`
	SyncProtocolID  core.ProtocolID `env:"SYNC_PROTOCOL_ID" envDefault:"/p2p/gpn-sync/1.0.0"`
	Namespace       string          `env:"NAMESPACE" envDefault:"mpc-pubsub"`
	PrivateKeyPath  string          `env:"PRIVATE_KEY_PATH" envDefault:"priv.key"`
	ContractAddress string          `env:"CONTRACT_ADDRESS" envDefault:"0x5510E82f2A7f0B1397Ef60FE1751DCB722C66ED9"`
	Port            string          `env:"PORT" envDefault:"0"`
	Consumers       []string        `env:"CONSUMERS" envDefault:"matterlabs/prover,scroll-tech/scroll-prover"`
	Mode            string          `env:"MODE" envDefault:"production"`
}

func NewConfig() (*Config, error) {
	conf := new(Config)
	if err := env.Parse(conf); err != nil {
		return nil, errors.Wrap(err, "error on parsing config")
	}

	if err := validateConfig(*conf); err != nil {
		return nil, errors.Wrap(err, "error on validating config")
	}

	return conf, nil
}

func validateConfig(cfg Config) error {
	if cfg.ProtocolID == "" {
		return errors.New("protocol ID is required")
	}

	if cfg.Namespace == "" {
		return errors.New("namespace is required")
	}

	if cfg.PrivateKeyPath == "" {
		return errors.New("private key path is required")
	}

	if cfg.Port == "" {
		return errors.New("port is required")
	}

	return nil
}
