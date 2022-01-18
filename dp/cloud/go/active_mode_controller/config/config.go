package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DialTimeout               time.Duration
	HeartbeatSendTimeout      time.Duration
	RequestTimeout            time.Duration
	RequestProcessingInterval time.Duration
	PollingInterval           time.Duration
	GrpcService               string
	GrpcPort                  int
	CbsdInactivityTimeout     time.Duration
}

// Unfortunately it was decided that timeouts should not have units
// and this is why this proxy config is needed
type appConfig struct {
	DialTimeoutSec               int    `envconfig:"DIAL_TIMEOUT_SEC"`
	HeartbeatSendTimeoutSec      int    `envconfig:"HEARTBEAT_SEND_TIMEOUT_SEC"`
	RequestTimeoutSec            int    `envconfig:"REQUEST_TIMEOUT_SEC"`
	RequestProcessingIntervalSec int    `envconfig:"REQUEST_PROCESSING_INTERVAL_SEC"`
	PollingIntervalSec           int    `envconfig:"POLLING_INTERVAL_SEC"`
	GrpcService                  string `envconfig:"GRPC_SERVICE"`
	GrpcPort                     int    `envconfig:"GRPC_PORT"`
	CbsdInactivityTimeoutSec     int    `envconfig:"CBSD_INACTIVITY_TIMEOUT_SEC"`
}

func Read() (*Config, error) {
	cfg := &appConfig{
		DialTimeoutSec:               60,
		HeartbeatSendTimeoutSec:      10,
		RequestTimeoutSec:            5,
		RequestProcessingIntervalSec: 10,
		PollingIntervalSec:           10,
		GrpcService:                  "domain-proxy-radio-controller",
		GrpcPort:                     50053,
		CbsdInactivityTimeoutSec:     4 * 60 * 60,
	}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return toAppConfig(cfg), nil
}

func toAppConfig(c *appConfig) *Config {
	return &Config{
		DialTimeout:               secToDuration(c.DialTimeoutSec),
		HeartbeatSendTimeout:      secToDuration(c.HeartbeatSendTimeoutSec),
		RequestTimeout:            secToDuration(c.RequestTimeoutSec),
		RequestProcessingInterval: secToDuration(c.RequestProcessingIntervalSec),
		PollingInterval:           secToDuration(c.PollingIntervalSec),
		GrpcService:               c.GrpcService,
		GrpcPort:                  c.GrpcPort,
		CbsdInactivityTimeout:     secToDuration(c.CbsdInactivityTimeoutSec),
	}
}

func secToDuration(s int) time.Duration {
	return time.Duration(s) * time.Second
}
