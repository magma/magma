/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Shared diameter settings across magma cloud
package diameter

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	DiamHost        = "magma-fedgw.magma.com"
	DiamRealm       = "magma.com"
	DiamProductName = "magma"
	Vendor3GPP      = uint32(10415) // diameter code for a 3GPP application

	AddrFlag            = "addr"
	NetworkFlag         = "network"
	HostFlag            = "host"
	RealmFlag           = "realm"
	ProductFlag         = "product"
	LocalAddrFlag       = "laddr"
	DestHostFlag        = "dest_host"
	DestRealmFlag       = "dest_realm"
	DisableDestHostFlag = "disable_dest_host"

	DefaultWatchdogIntervalSeconds = 3
)

// Diameter flags
var (
	_ = flag.String(AddrFlag, "", "Server address (host:port)")
	_ = flag.String(NetworkFlag, "", "protocol (sctp/tcp)")
	_ = flag.String(HostFlag, "", "Diameter host")
	_ = flag.String(RealmFlag, "", "Diameter realm")
	_ = flag.String(ProductFlag, "", "Diameter product name")
	_ = flag.String(LocalAddrFlag, "", "Client local address to bind to (IP:port OR :port)")
	_ = flag.String(DestHostFlag, "", "Diameter server host name")
	_ = flag.String(DestRealmFlag, "", "Diameter server realm")
	_ = flag.String(DisableDestHostFlag, "", "Disable sending dest-host AVP in requests")
)

type DiameterServerConnConfig struct {
	Addr      string // host:port
	Protocol  string // tcp/sctp
	LocalAddr string // IP:port or :port
}

type DiameterServerConfig struct {
	DiameterServerConnConfig
	DestHost        string
	DestRealm       string
	DisableDestHost bool
}

// DiameterClientConfig holds information for connecting with a diameter server
type DiameterClientConfig struct {
	Host               string // diameter host
	Realm              string // diameter realm
	ProductName        string
	AppID              uint32
	AuthAppID          uint32
	Retransmits        uint
	WatchdogInterval   uint
	RetryCount         uint // number of times to reconnect after connection lost
	SupportedVendorIDs string
	ServiceContextId   string
}

func (cfg *DiameterServerConfig) Validate() error {
	if cfg == nil {
		return fmt.Errorf("Nil server config")
	}
	// validate network address, replace 'sctp' with 'tcp' to check resolving
	network := cfg.Protocol
	if len(network) == 0 {
		return fmt.Errorf("Empty network protocol")
	} else if strings.Index(network, "sctp") == 0 {
		network = "tcp" + network[4:]
	}
	_, err := net.ResolveTCPAddr(network, cfg.Addr)
	if err != nil {
		return fmt.Errorf("Invalid Diameter Address (%s://%s): %v", cfg.Protocol, cfg.Addr, err)
	}
	return nil
}

func (cfg *DiameterClientConfig) Validate() error {
	if cfg == nil {
		return fmt.Errorf("Nil client config")
	}
	if len(cfg.Host) == 0 {
		return fmt.Errorf("Invalid Diameter Host")
	}
	if len(cfg.Realm) == 0 {
		return fmt.Errorf("Invalid Diameter Realm")
	}
	return nil
}

func (srcCfg *DiameterClientConfig) FillInDefaults() *DiameterClientConfig {
	if srcCfg == nil {
		return nil
	}
	cfg := *srcCfg
	if len(cfg.Host) == 0 {
		cfg.Host = DiamHost
	}
	if len(cfg.Realm) == 0 {
		cfg.Realm = DiamRealm
	}
	if len(cfg.ProductName) == 0 {
		cfg.ProductName = DiamProductName
	}
	if cfg.Retransmits == 0 {
		cfg.Retransmits = 3
	}
	return &cfg
}

// getUint64FlagValue looks up the flag and either returns its uint64 value
// or an error.
func getUint64FlagValue(flagName string) (uint64, error) {
	f := flag.Lookup(flagName)
	if f == nil {
		return 0, fmt.Errorf("Flag not found: %s", flagName)
	}
	if f.Value == nil {
		return 0, fmt.Errorf("Flag value is nil: %s", flagName)
	}
	getter, ok := f.Value.(flag.Getter)
	if !ok {
		return 0, fmt.Errorf("Flag value has no Getter: %s", flagName)
	}
	value, ok := getter.Get().(uint64)
	if !ok {
		return 0, fmt.Errorf("Flag value is not of type uint64: %s", flagName)
	}
	return value, nil
}

// GetValueUint64 returns value of the flagValue if it exists, or defaultValue if not
func GetValueUint64(flagName string, defaultValue uint64) uint64 {
	if len(flagName) > 0 {
		value, err := getUint64FlagValue(flagName)
		if err != nil {
			return value
		}

	}
	log.Printf("Using default value: %v for flag: %s", defaultValue, flagName)
	return defaultValue
}

// GetValue returns value of the flagValue if it exists, or defaultValue if not
func GetValue(flagName, defaultValue string) string {
	flagValue := getFlagValue(flagName)
	if len(flagValue) != 0 {
		return flagValue
	}
	log.Printf("Using default value: %s for flag: %s", defaultValue, flagName)
	return defaultValue
}

// GetValueOrEnv returns value of the flagValue if it exists, then the environment
// variable if it exists, or defaultValue if not
func GetValueOrEnv(flagName, envVariable, defaultValue string) string {
	flagValue := getFlagValue(flagName)
	if len(flagValue) != 0 {
		return flagValue
	}
	if len(envVariable) > 0 {
		envValue := os.Getenv(envVariable)
		if len(envValue) > 0 {
			log.Printf("Using Environment Parameter: %s => %s (default: '%s')", envVariable, envValue, defaultValue)
			return envValue
		}
	}
	log.Printf("Using default value: %s for flag: %s", defaultValue, flagName)
	return defaultValue
}

// GetBoolValueOrEnv returns value of the flagValue if it exists, then the environment
// variable if it exists, or defaultValue if not
func GetBoolValueOrEnv(flagName string, envVariable string, defaultValue bool) bool {
	flagValue := getFlagValue(flagName)
	flagValueBool, err := strconv.ParseBool(flagValue)
	if len(flagValue) != 0 && err == nil {
		return flagValueBool
	}
	if len(envVariable) > 0 {
		envValue := os.Getenv(envVariable)
		envValueBool, err := strconv.ParseBool(envValue)
		if len(envValue) > 0 && err == nil {
			log.Printf("Using Environment Parameter: %s => %t (default: '%t')", envVariable, envValueBool, defaultValue)
			return envValueBool
		}
	}
	log.Printf("Using default value: %t for flag: %s", defaultValue, flagName)
	return defaultValue
}

// getFlagValue returns the value of the flagValue if it exists, or an empty string if not
func getFlagValue(flagName string) string {
	if len(flagName) > 0 {
		f := flag.Lookup(flagName)
		if f != nil && len(f.Value.String()) != 0 {
			res := f.Value.String()
			log.Printf("Using runtime flag: %s => %s", flagName, res)
			return res
		}
	}
	return ""
}
