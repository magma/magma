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
	"fmt"

	mcfgprotos "magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/plmn_filter"
	"magma/feg/gateway/services/swx_proxy/cache"
	"magma/gateway/mconfig"

	"github.com/golang/glog"
)

const (
	SwxProxyServiceName = "swx_proxy"

	HSSAddrEnv           = "HSS_ADDR"
	SwxNetworkEnv        = "SWX_NETWORK"
	SwxDiamHostEnv       = "SWX_DIAM_HOST"
	SwxDiamRealmEnv      = "SWX_DIAM_REALM"
	SwxDiamProductEnv    = "SWX_DIAM_PRODUCT"
	SwxLocalAddrEnv      = "SWX_LOCAL_ADDR"
	HSSHostEnv           = "HSS_HOST"
	HSSRealmEnv          = "HSS_REALM"
	DisableDestHostEnv   = "DISABLE_DEST_HOST"
	OverwriteDestHostEnv = "OVERWRITE_DEST_HOST"

	DefaultSwxDiamRealm          = "epc.mnc070.mcc722.3gppnetwork.org"
	DefaultSwxDiamHost           = "feg-swx.epc.mnc070.mcc722.3gppnetwork.org"
	DefaultVerifyAuthorization   = false
	DefaultRegisterOnAuth        = false
	DefaultDeriveUnregisterRealm = false
)

// GetSwxProxyConfig returns the service config based on the
// the values in mconfig or default values provided
func GetSwxProxyConfig() []*SwxProxyConfig {
	configsPtr := &mcfgprotos.SwxConfig{}
	err := mconfig.GetServiceConfigs(SwxProxyServiceName, configsPtr)

	if err != nil || !isSWxMConfiValid(configsPtr) {
		glog.V(2).Infof("%s Managed Configs Load Error: %v", SwxProxyServiceName, err)

		return []*SwxProxyConfig{
			{
				ClientCfg: &diameter.DiameterClientConfig{
					Host:        diameter.GetValueOrEnv(diameter.HostFlag, SwxDiamHostEnv, DefaultSwxDiamHost),
					Realm:       diameter.GetValueOrEnv(diameter.RealmFlag, SwxDiamRealmEnv, DefaultSwxDiamRealm),
					ProductName: diameter.GetValueOrEnv(diameter.ProductFlag, SwxDiamProductEnv, diameter.DiamProductName),
				},
				ServerCfg: &diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
					Addr:      diameter.GetValueOrEnv(diameter.AddrFlag, HSSAddrEnv, ""),
					Protocol:  diameter.GetValueOrEnv(diameter.NetworkFlag, SwxNetworkEnv, "sctp"),
					LocalAddr: diameter.GetValueOrEnv(diameter.LocalAddrFlag, SwxLocalAddrEnv, "")},
					DestHost:          diameter.GetValueOrEnv(diameter.DestHostFlag, HSSHostEnv, ""),
					DestRealm:         diameter.GetValueOrEnv(diameter.DestRealmFlag, HSSRealmEnv, ""),
					DisableDestHost:   diameter.GetBoolValueOrEnv(diameter.DisableDestHostFlag, DisableDestHostEnv, false),
					OverwriteDestHost: diameter.GetBoolValueOrEnv(diameter.OverwriteDestHostFlag, OverwriteDestHostEnv, false),
				},
				VerifyAuthorization:   DefaultVerifyAuthorization,
				RegisterOnAuth:        DefaultRegisterOnAuth,
				DeriveUnregisterRealm: DefaultDeriveUnregisterRealm,
				CacheTTLSeconds:       uint32(cache.DefaultTtl.Seconds()),
				HlrPlmnIds:            plmn_filter.PlmnIdVals{},
			},
		}
	}

	glog.V(2).Infof("Loaded %s configs: %+v", SwxProxyServiceName, *configsPtr)

	hlrPlmnIds := plmn_filter.GetPlmnVals(configsPtr.HlrPlmnIds, "HLR")

	ttl := configsPtr.CacheTTLSeconds
	if ttl < uint32(cache.DefaultGcInterval.Seconds()) {
		ttl = uint32(cache.DefaultTtl.Seconds())
	}
	swxConfigs := configsPtr.GetServers()

	//TODO: remove this once backwards compatibility is not needed for the field server
	if len(swxConfigs) == 0 {
		server := configsPtr.GetServer()
		if server == nil {
			glog.V(2).Infof("Server configuration for Swx servers not found!!")
		} else {
			swxConfigs = append(swxConfigs, server)
			glog.V(2).Infof("Swx Server configuration using legacy swagger attribute Server (not Servers)")
		}
	}

	// Iterate over the slice of servers. VarEnv will apply only to index 0
	diamServerConfigs := []*SwxProxyConfig{}
	for i, swxConfig := range swxConfigs {
		diamSrvCfg := &SwxProxyConfig{
			ClientCfg: &diameter.DiameterClientConfig{
				Host:             diameter.GetValueOrEnv(diameter.HostFlag, SwxDiamHostEnv, swxConfig.GetHost(), i),
				Realm:            diameter.GetValueOrEnv(diameter.RealmFlag, SwxDiamRealmEnv, swxConfig.GetRealm(), i),
				ProductName:      diameter.GetValueOrEnv(diameter.ProductFlag, SwxDiamProductEnv, swxConfig.GetProductName(), i),
				Retransmits:      uint(swxConfig.GetRetransmits()),
				WatchdogInterval: uint(swxConfig.GetWatchdogInterval()),
				RetryCount:       uint(swxConfig.GetRetryCount()),
			},
			ServerCfg: &diameter.DiameterServerConfig{
				DiameterServerConnConfig: diameter.DiameterServerConnConfig{
					Addr:      diameter.GetValueOrEnv(diameter.AddrFlag, HSSAddrEnv, swxConfig.GetAddress(), i),
					Protocol:  diameter.GetValueOrEnv(diameter.NetworkFlag, SwxNetworkEnv, swxConfig.GetProtocol(), i),
					LocalAddr: diameter.GetValueOrEnv(diameter.LocalAddrFlag, SwxLocalAddrEnv, swxConfig.GetLocalAddress(), i)},
				DestHost:          diameter.GetValueOrEnv(diameter.DestHostFlag, HSSHostEnv, swxConfig.GetDestHost(), i),
				DestRealm:         diameter.GetValueOrEnv(diameter.DestRealmFlag, HSSRealmEnv, swxConfig.GetDestRealm(), i),
				DisableDestHost:   diameter.GetBoolValueOrEnv(diameter.DisableDestHostFlag, DisableDestHostEnv, swxConfig.GetDisableDestHost(), i),
				OverwriteDestHost: diameter.GetBoolValueOrEnv(diameter.OverwriteDestHostFlag, OverwriteDestHostEnv, swxConfig.GetOverwriteDestHost(), i),
			},
			VerifyAuthorization:   configsPtr.GetVerifyAuthorization(),
			RegisterOnAuth:        configsPtr.GetRegisterOnAuth(),
			DeriveUnregisterRealm: configsPtr.GetDeriveUnregisterRealm(),
			CacheTTLSeconds:       ttl,
			HlrPlmnIds:            hlrPlmnIds,
		}
		diamServerConfigs = append(diamServerConfigs, diamSrvCfg)
	}
	return diamServerConfigs
}

// ValidateSwxProxyConfig ensures that the swx proxy config specified has valid
// diameter client and server configs
func ValidateSwxProxyConfig(config *SwxProxyConfig) error {
	if config == nil {
		return fmt.Errorf("Nil SwxProxyConfig provided")
	}
	if config.ClientCfg == nil {
		return fmt.Errorf("Nil client config provided")
	}
	err := config.ClientCfg.Validate()
	if err != nil {
		return err
	}
	if config.ServerCfg == nil {
		return fmt.Errorf("Nil server config provided")
	}
	return config.ServerCfg.Validate()
}

// isSWxMConfiValid check if required fields are present on SwxConfig proto
func isSWxMConfiValid(config *mcfgprotos.SwxConfig) bool {
	if config == nil ||
		(config.Server == nil && len(config.Servers) == 0) ||
		(config.Server != nil && config.Server.Address == "") {
		return false
	}
	for _, server := range config.Servers {
		if server.Address == "" {
			return false
		}
	}
	return true
}
