/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Shared diameter settings across magma cloud
package diameter

import (
	"flag"
	"fmt"
	"net"
	"strings"

	"magma/feg/gateway/utils"
)

const (
	DiamHost        = "magma-fedgw.magma.com"
	DiamRealm       = "magma.com"
	DiamProductName = "magma"
	Vendor3GPP      = uint32(10415) // diameter code for a 3GPP application

	AddrFlag              = "addr"
	NetworkFlag           = "network"
	HostFlag              = "host"
	RealmFlag             = "realm"
	ProductFlag           = "product"
	LocalAddrFlag         = "laddr"
	DestHostFlag          = "dest_host"
	DestRealmFlag         = "dest_realm"
	DisableDestHostFlag   = "disable_dest_host"
	OverwriteDestHostFlag = "overwrite_dest_host"

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
	_ = flag.String(OverwriteDestHostFlag, "", "Overwrite dest-host AVP in requests even if message includes it")
)

type DiameterServerConnConfig struct {
	Addr      string // host:port
	Protocol  string // tcp/sctp
	LocalAddr string // IP:port or :port
}

type DiameterServerConfig struct {
	DiameterServerConnConfig
	DestHost          string
	DestRealm         string
	DisableDestHost   bool
	OverwriteDestHost bool
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

// GetValueUint64 returns value of the flagValue if it exists, or defaultValue if not
func GetValueUint64(flagName string, defaultValue uint64) uint64 {
	return utils.GetValueUint64(flagName, defaultValue)
}

// GetValue returns value of the flagValue if it exists, or defaultValue if not
func GetValue(flagName, defaultValue string) string {
	return utils.GetValue(flagName, defaultValue)
}

// GetValueOrEnv returns value of the flagValue if it exists, then the environment
// variable if it exists, or defaultValue if not.
// If idx parameter is passed, then if that idx > 1 defaultValue will be returned.
// Note in case of many idx are passed, only the first idx will be checked.
func GetValueOrEnv(flagName, envVariable, defaultValue string, idx ...int) string {
	if len(idx) > 0 && idx[0] > 0 {
		return defaultValue
	}
	return utils.GetValueOrEnv(flagName, envVariable, defaultValue)
}

// GetBoolValueOrEnv returns value of the flagValue if it exists, then the environment
// variable if it exists, or defaultValue if not.
// If idx parameter is passed, then if that idx > 1 defaultValue will be returned.
// Note in case of many idx are passed, only the first idx will be checked.
func GetBoolValueOrEnv(flagName string, envVariable string, defaultValue bool, idx ...int) bool {
	if len(idx) > 0 && idx[0] > 0 {
		return defaultValue
	}
	return utils.GetBoolValueOrEnv(flagName, envVariable, defaultValue)
}
