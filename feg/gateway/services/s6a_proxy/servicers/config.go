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

package servicers

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/plmn_filter"
	managed_configs "magma/gateway/mconfig"
)

// S6a Environment Variables to overwrite default configs
const (
	HSSAddrEnv           = "HSS_ADDR"
	S6aNetworkEnv        = "S6A_NETWORK"
	S6aDiamHostEnv       = "S6A_DIAM_HOST"
	S6aDiamRealmEnv      = "S6A_DIAM_REALM"
	S6aDiamProductEnv    = "S6A_DIAM_PRODUCT"
	S6aLocalAddrEnv      = "S6A_LOCAL_ADDR"
	HSSHostEnv           = "HSS_HOST"
	HSSRealmEnv          = "HSS_REALM"
	DisableDestHostEnv   = "DISABLE_DEST_HOST"
	OverwriteDestHostEnv = "OVERWRITE_DEST_HOST"

	S6aProxyServiceName = "s6a_proxy"
	DefaultS6aDiamRealm = "epc.mnc070.mcc722.3gppnetwork.org"
	DefaultS6aDiamHost  = "feg-s6a.epc.mnc070.mcc722.3gppnetwork.org"
)

// Get GetS6aProxyConfigs returns the server config for an HSS based on the
// input flags and environment variables
func GetS6aProxyConfigs() *S6aProxyConfig {
	serviceBaseName := filepath.Base(os.Args[0])
	serviceBaseName = strings.TrimSuffix(serviceBaseName, filepath.Ext(serviceBaseName))
	if S6aProxyServiceName != serviceBaseName {
		log.Printf(
			"NOTE: S6a Proxy Base Service name: %s does not match its managed configs key: %s",
			serviceBaseName, S6aProxyServiceName)
	}
	configsPtr := &mconfig.S6AConfig{}
	err := managed_configs.GetServiceConfigs(S6aProxyServiceName, configsPtr)
	if err != nil || configsPtr.Server == nil {
		log.Printf("%s Managed Configs Load Error: %v", S6aProxyServiceName, err)
		return &S6aProxyConfig{
			ClientCfg: &diameter.DiameterClientConfig{
				Host:        diameter.GetValueOrEnv(diameter.HostFlag, S6aDiamHostEnv, DefaultS6aDiamHost),
				Realm:       diameter.GetValueOrEnv(diameter.RealmFlag, S6aDiamRealmEnv, DefaultS6aDiamRealm),
				ProductName: diameter.GetValueOrEnv(diameter.ProductFlag, S6aDiamProductEnv, diameter.DiamProductName),
			},
			ServerCfg: &diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
				Addr:      diameter.GetValueOrEnv(diameter.AddrFlag, HSSAddrEnv, ""),
				Protocol:  diameter.GetValueOrEnv(diameter.NetworkFlag, S6aNetworkEnv, "sctp"),
				LocalAddr: diameter.GetValueOrEnv(diameter.LocalAddrFlag, S6aLocalAddrEnv, "")},
				DestHost:          diameter.GetValueOrEnv(diameter.DestHostFlag, HSSHostEnv, ""),
				DestRealm:         diameter.GetValueOrEnv(diameter.DestRealmFlag, HSSRealmEnv, ""),
				DisableDestHost:   diameter.GetBoolValueOrEnv(diameter.DisableDestHostFlag, DisableDestHostEnv, false),
				OverwriteDestHost: diameter.GetBoolValueOrEnv(diameter.OverwriteDestHostFlag, OverwriteDestHostEnv, false),
			},
			PlmnIds: plmn_filter.PlmnIdVals{},
		}
	}

	log.Printf("Loaded %s configs: %+v", S6aProxyServiceName, *configsPtr)

	return &S6aProxyConfig{
		ClientCfg: &diameter.DiameterClientConfig{
			Host:             diameter.GetValueOrEnv(diameter.HostFlag, S6aDiamHostEnv, configsPtr.Server.Host),
			Realm:            diameter.GetValueOrEnv(diameter.RealmFlag, S6aDiamRealmEnv, configsPtr.Server.Realm),
			ProductName:      diameter.GetValueOrEnv(diameter.ProductFlag, S6aDiamProductEnv, configsPtr.Server.ProductName),
			Retransmits:      uint(configsPtr.Server.Retransmits),
			WatchdogInterval: uint(configsPtr.Server.WatchdogInterval),
			RetryCount:       uint(configsPtr.Server.RetryCount),
		},
		ServerCfg: &diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
			Addr:      diameter.GetValueOrEnv(diameter.AddrFlag, HSSAddrEnv, configsPtr.Server.Address),
			Protocol:  diameter.GetValueOrEnv(diameter.NetworkFlag, S6aNetworkEnv, configsPtr.Server.Protocol),
			LocalAddr: diameter.GetValueOrEnv(diameter.LocalAddrFlag, S6aLocalAddrEnv, configsPtr.Server.LocalAddress)},
			DestHost:          diameter.GetValueOrEnv(diameter.DestHostFlag, HSSHostEnv, configsPtr.Server.DestHost),
			DestRealm:         diameter.GetValueOrEnv(diameter.DestRealmFlag, HSSRealmEnv, configsPtr.Server.DestRealm),
			DisableDestHost:   diameter.GetBoolValueOrEnv(diameter.DisableDestHostFlag, DisableDestHostEnv, configsPtr.GetServer().GetDisableDestHost()),
			OverwriteDestHost: diameter.GetBoolValueOrEnv(diameter.OverwriteDestHostFlag, OverwriteDestHostEnv, configsPtr.GetServer().GetOverwriteDestHost()),
		},
		PlmnIds: plmn_filter.GetPlmnVals(configsPtr.PlmnIds),
	}
}
