package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DialTimeout     time.Duration `envconfig:"DIAL_TIMEOUT"`
	RequestTimeout  time.Duration `envconfig:"REQUEST_TIMEOUT"`
	PollingInterval time.Duration `envconfig:"POLLING_INTERVAL"`
	GrpcService     string        `envconfig:"GRPC_SERVICE"`
	GrpcPort        int           `envconfig:"GRPC_PORT"`
}

func Read() (*Config, error) {
	cfg := &Config{
		DialTimeout:     time.Second * 60,
		RequestTimeout:  time.Second * 5,
		PollingInterval: time.Second * 30,
		GrpcService:     "domain-proxy-radio-controller",
		GrpcPort:        50053,
	}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
