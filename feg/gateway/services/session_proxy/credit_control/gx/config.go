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

package gx

import (
	"flag"
	"log"

	"github.com/fiorix/go-diameter/v4/diam"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/session_proxy/credit_control"
	managed_configs "magma/gateway/mconfig"
)

// PCRF Environment Variables
const (
	PCRFAddrEnv               = "PCRF_ADDR"
	GxNetworkEnv              = "GX_NETWORK"
	GxDiamHostEnv             = "GX_DIAM_HOST"
	GxDiamRealmEnv            = "GX_DIAM_REALM"
	GxDiamProductEnv          = "GX_DIAM_PRODUCT"
	GxLocalAddr               = "GX_LOCAL_ADDR"
	PCRFHostEnv               = "PCRF_HOST"
	PCRFRealmEnv              = "PCRF_REALM"
	PCRF91CompliantEnv        = "PCRF_91_COMPLIANT"
	DisableDestHostEnv        = "DISABLE_DEST_HOST"
	OverwriteDestHostEnv      = "GX_OVERWRITE_DEST_HOST"
	DisableEUIIPv6IfNoIPEnv   = "DISABLE_EUI64_IPV6_IF_NO_IP"
	FramedIPv4AddrRequiredEnv = "FRAMED_IPV4_ADDR_REQUIRED"
	DefaultFramedIPv4AddrEnv  = "DEFAULT_FRAMED_IPV4_ADDR"
	GxSupportedVendorIDsEnv   = "GX_SUPPORTED_VENDOR_IDS"

	PCRF91CompliantFlag      = "pcrf_91_compliant"
	DisableEUIIPv6IfNoIPFlag = "disable_eui64_ipv6_prefix"
)

var (
	pcrf91Compliant = flag.Bool(
		PCRF91CompliantFlag, false, "Set to support 29.212 release 9.1 compliant PCRF")
	disableEUIIpIfEmpty = flag.Bool(
		DisableEUIIPv6IfNoIPFlag, false, "Don't use MAC based EUI-64 IPv6 address for Gx CCR if IP is not provided")
)

// TODO: refactor those functions to make it more simple
// GetPCRFConfiguration returns a slice containing all configuration for all known PCRF
func GetPCRFConfiguration() []*diameter.DiameterServerConfig {
	configsPtr := &mconfig.SessionProxyConfig{}
	err := managed_configs.GetServiceConfigs(credit_control.SessionProxyServiceName, configsPtr)
	if err != nil || !validGxConfig(configsPtr) {
		log.Printf("%s Managed Gx PCRF Configs Load Error: %v", credit_control.SessionProxyServiceName, err)
		return []*diameter.DiameterServerConfig{
			{
				DiameterServerConnConfig: diameter.DiameterServerConnConfig{
					Addr:      diameter.GetValueOrEnv(diameter.AddrFlag, PCRFAddrEnv, "127.0.0.1:3870"),
					Protocol:  diameter.GetValueOrEnv(diameter.NetworkFlag, GxNetworkEnv, "tcp"),
					LocalAddr: diameter.GetValueOrEnv(diameter.LocalAddrFlag, GxLocalAddr, ""),
				},
				DestHost:          diameter.GetValueOrEnv(diameter.DestHostFlag, PCRFHostEnv, ""),
				DestRealm:         diameter.GetValueOrEnv(diameter.DestRealmFlag, PCRFRealmEnv, ""),
				DisableDestHost:   diameter.GetBoolValueOrEnv(diameter.DisableDestHostFlag, DisableDestHostEnv, false),
				OverwriteDestHost: diameter.GetBoolValueOrEnv(diameter.DisableDestHostFlag, OverwriteDestHostEnv, false),
			},
		}
	}

	gxConfigs := configsPtr.GetGx().GetServers()
	//TODO: remove this once backwards compatibility is not needed for the field server
	if len(gxConfigs) == 0 {
		server := configsPtr.GetGx().GetServer()
		if server == nil {
			log.Print("Server configuration for Gx servers not found!!")
		} else {
			gxConfigs = append(gxConfigs, server)
			log.Print("Gx Server configuration using legacy swagger attribute Server (not Servers)")
		}
	}

	// Iterate over the slice of servers. VarEnv will apply only to index 0
	diamServerConfigs := []*diameter.DiameterServerConfig{}
	for i, gxCfg := range gxConfigs {
		diamSrvCfg := &diameter.DiameterServerConfig{
			DiameterServerConnConfig: diameter.DiameterServerConnConfig{
				Addr:      diameter.GetValueOrEnv(diameter.AddrFlag, PCRFAddrEnv, gxCfg.GetAddress(), i),
				Protocol:  diameter.GetValueOrEnv(diameter.NetworkFlag, GxNetworkEnv, gxCfg.GetProtocol(), i),
				LocalAddr: diameter.GetValueOrEnv(diameter.LocalAddrFlag, GxLocalAddr, gxCfg.GetLocalAddress(), i),
			},
			DestHost:          diameter.GetValueOrEnv(diameter.DestHostFlag, PCRFHostEnv, gxCfg.GetDestHost(), i),
			DestRealm:         diameter.GetValueOrEnv(diameter.DestRealmFlag, PCRFRealmEnv, gxCfg.GetDestRealm(), i),
			DisableDestHost:   diameter.GetBoolValueOrEnv(diameter.DisableDestHostFlag, DisableDestHostEnv, gxCfg.GetDisableDestHost(), i),
			OverwriteDestHost: diameter.GetBoolValueOrEnv(diameter.OverwriteDestHostFlag, OverwriteDestHostEnv, gxCfg.GetOverwriteDestHost(), i),
		}
		diamServerConfigs = append(diamServerConfigs, diamSrvCfg)
	}

	return diamServerConfigs

}

// GetGxClientConfiguration returns a slice containing all client diameter configuration
func GetGxClientConfiguration() []*diameter.DiameterClientConfig {
	var retries uint32 = 1
	configsPtr := &mconfig.SessionProxyConfig{}
	err := managed_configs.GetServiceConfigs(credit_control.SessionProxyServiceName, configsPtr)
	if err != nil {
		log.Printf("%s Managed Gx Client Configs Load Error: %v", credit_control.SessionProxyServiceName, err)
		return []*diameter.DiameterClientConfig{
			{
				Host:               diameter.GetValueOrEnv(diameter.HostFlag, GxDiamHostEnv, diameter.DiamHost),
				Realm:              diameter.GetValueOrEnv(diameter.RealmFlag, GxDiamRealmEnv, diameter.DiamRealm),
				ProductName:        diameter.GetValueOrEnv(diameter.ProductFlag, GxDiamProductEnv, diameter.DiamProductName),
				AppID:              diam.GX_CHARGING_CONTROL_APP_ID,
				WatchdogInterval:   diameter.DefaultWatchdogIntervalSeconds,
				RetryCount:         uint(retries),
				SupportedVendorIDs: diameter.GetValueOrEnv("", GxSupportedVendorIDsEnv, ""),
			},
		}
	}

	diamClientsConfigs := []*diameter.DiameterClientConfig{}
	gxConfigs := configsPtr.GetGx().GetServers()
	//TODO: remove this once backwards compatibility is not needed for the field server
	if len(gxConfigs) == 0 {
		server := configsPtr.GetGx().GetServer()
		if server == nil {
			log.Print("Client configuration for Gx servers not found!!")
		} else {
			gxConfigs = append(gxConfigs, server)
			log.Print("Gx Client configuration using legacy swagger attribute Server (not Servers)")
		}
	}

	for i, gxCfg := range gxConfigs {
		retries = gxCfg.GetRetryCount()
		if retries < 1 {
			log.Printf("Invalid Gx Server Retry Count for server (%s): %d, must be >0. Will be set to 1", gxCfg.GetAddress(), retries)
			retries = 1
		}

		wdInterval := gxCfg.GetWatchdogInterval()
		if wdInterval == 0 {
			wdInterval = diameter.DefaultWatchdogIntervalSeconds
		}
		diamCliCfg := &diameter.DiameterClientConfig{
			Host:               diameter.GetValueOrEnv(diameter.HostFlag, GxDiamHostEnv, gxCfg.GetHost(), i),
			Realm:              diameter.GetValueOrEnv(diameter.RealmFlag, GxDiamRealmEnv, gxCfg.GetRealm(), i),
			ProductName:        diameter.GetValueOrEnv(diameter.ProductFlag, GxDiamProductEnv, gxCfg.GetProductName(), i),
			AppID:              diam.GX_CHARGING_CONTROL_APP_ID,
			WatchdogInterval:   uint(wdInterval),
			RetryCount:         uint(retries),
			SupportedVendorIDs: diameter.GetValueOrEnv("", GxSupportedVendorIDsEnv, "", i),
		}
		diamClientsConfigs = append(diamClientsConfigs, diamCliCfg)
	}
	return diamClientsConfigs

}

func GetGxGlobalConfig() *GxGlobalConfig {
	configsPtr := &mconfig.SessionProxyConfig{}
	err := managed_configs.GetServiceConfigs(credit_control.SessionProxyServiceName, configsPtr)
	if err != nil || !validGxConfig(configsPtr) {
		log.Printf("%s Managed Gx Server Configs Load Error: %v", credit_control.SessionProxyServiceName, err)
		return &GxGlobalConfig{}
	}
	return &GxGlobalConfig{
		PCFROverwriteApn: configsPtr.GetGx().GetOverwriteApn(),
		DisableGx:        configsPtr.GetGx().GetDisableGx(),
		VirtualApnRules:  credit_control.GenerateVirtualApnRules(configsPtr.GetGx().GetVirtualApnRules()),
	}
}

// validGxConfig check if required fields related to Gx are valid in the config
func validGxConfig(config *mconfig.SessionProxyConfig) bool {
	if config == nil || config.Gx == nil ||
		(config.Gx.Server == nil && len(config.Gx.Servers) == 0) ||
		(config.Gx.Server != nil && config.Gx.Server.Address == "") {
		return false
	}
	for _, server := range config.Gx.Servers {
		if server.Address == "" {
			return false
		}
	}
	return true
}
