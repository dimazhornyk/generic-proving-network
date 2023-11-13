package common

import (
	"github.com/caarlos0/env"
	"github.com/libp2p/go-libp2p/core"
	"github.com/pkg/errors"
)

type Config struct {
	ProtocolID      core.ProtocolID `env:"PROTOCOL_ID" envDefault:"/p2p/gpc-node/1.0.0"`
	SyncProtocolID  core.ProtocolID `env:"SYNC_PROTOCOL_ID" envDefault:"/p2p/gpc-sync/1.0.0"`
	Namespace       string          `env:"NAMESPACE" envDefault:"mpc-pubsub"`
	PrivateKeyPath  string          `env:"PRIVATE_KEY_PATH" envDefault:"/app/keys/priv.key"`
	ContractAddress string          `env:"CONTRACT_ADDRESS" envDefault:"0x1E0447b19BB6EcFdAe1e4AE1694b0C3659614e4e"`
	Port            string          `env:"PORT" envDefault:"0"`
	Consumers       []string        `env:"CONSUMERS" envDefault:"zksync-prover,scroll-prover"`
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
