package service

import (
	"fmt"
	"net"
	"strings"
)

type S6aProxyConfig struct {
	HssAddr, // host:port
	Protocol, // tcp/sctp
	Host, // diameter host
	Realm string // diameter realm
	Retransmits,
	WatchdogInterval uint
}

func (cfg *S6aProxyConfig) Validate() error {
	if cfg == nil {
		return fmt.Errorf("Nil S6aProxy config")
	}
	network := cfg.Protocol
	if len(network) == 0 {
		cfg.Protocol = "sctp"
	} else if strings.Index(network, "sctp") == 0 {
		network = "tcp" + network[4:]
	}
	if len(cfg.Host) == 0 {
		return fmt.Errorf("Invalid Diameter Host")
	}
	if len(cfg.Realm) == 0 {
		return fmt.Errorf("Invalid Diameter Realm")
	}
	_, err := net.ResolveTCPAddr(network, cfg.HssAddr)
	if err != nil {
		return fmt.Errorf("Invalid HSS Address (%s://%s): %v", cfg.Protocol, cfg.HssAddr, err)
	}
	return nil
}

func (srcCfg *S6aProxyConfig) CloneWithDefaults() *S6aProxyConfig {
	if srcCfg == nil {
		return nil
	}
	cfg := *srcCfg
	if len(cfg.Protocol) == 0 {
		cfg.Protocol = "sctp"
	}
	if len(cfg.Host) == 0 {
		cfg.Protocol = "protocol.s6a.proxy"
	}
	if len(cfg.Realm) == 0 {
		cfg.Protocol = "realm.s6a.proxy"
	}
	if cfg.Retransmits == 0 {
		cfg.Retransmits = 3
	}
	if cfg.WatchdogInterval == 0 {
		cfg.WatchdogInterval = 7
	}
	return &cfg
}
