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

package gy

import (
	"flag"
	"log"

	"github.com/fiorix/go-diameter/v4/diam"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/session_proxy/credit_control"
	managed_configs "magma/gateway/mconfig"
)

// OCS Environment Variables
const (
	OCSAddrEnv              = "OCS_ADDR"
	GyNetworkEnv            = "GY_NETWORK"
	GyDiamHostEnv           = "GY_DIAM_HOST"
	GyDiamRealmEnv          = "GY_DIAM_REALM"
	GyDiamProductEnv        = "GY_DIAM_PRODUCT"
	GyInitMethodEnv         = "GY_INIT_METHOD"
	GyLocalAddr             = "GY_LOCAL_ADDR"
	OCSHostEnv              = "OCS_HOST"
	OCSRealmEnv             = "OCS_REALM"
	OCSApnOverwriteEnv      = "OCS_APN_OVERWRITE"
	OCSServiceIdentifierEnv = "OCS_SERVICE_IDENTIFIER_OVERWRITE"
	DisableDestHostEnv      = "DISABLE_DEST_HOST"
	OverwriteDestHostEnv    = "GY_OVERWRITE_DEST_HOST"
	UseGyForAuthOnlyEnv     = "USE_GY_FOR_AUTH_ONLY"
	GySupportedVendorIDsEnv = "GY_SUPPORTED_VENDOR_IDS"
	GyServiceContextIdEnv   = "GY_SERVICE_CONTEXT_ID"

	GyInitMethodFlag         = "gy_init_method"
	OCSApnOverwriteFlag      = "ocs_apn_overwrite"
	OCSServiceIdentifierFlag = "ocs_service_identifier_overwrite"
)

var (
	_ = flag.String(GyInitMethodFlag, "", "Gy init method (per_key|per_session)")
	_ = flag.String(OCSApnOverwriteFlag, "", "OCS APN to use instead of request's APN")
	_ = flag.String(OCSServiceIdentifierFlag, "", "OCS ServiceIdentifier to use in Gy requests")
)

// InitMethod describes the type of ways sessions can be initialized through the
// Gy interface
type InitMethod uint8

// InitMethod enum values
const (
	// 1 CCR-Init per session, multiple CCR-Updates to get initial credit
	PerSessionInit InitMethod = 1
	// CCR-Init per charging key
	PerKeyInit InitMethod = 2
)

// GetInitMethod returns the init method for this gy client based on the flags
// or environment variables
func GetInitMethod() InitMethod {
	initMethod := PerKeyInit
	configsPtr := &mconfig.SessionProxyConfig{}
	err := managed_configs.GetServiceConfigs(credit_control.SessionProxyServiceName, configsPtr)
	if err != nil || configsPtr.Gy == nil {
		log.Printf("%s Managed Gy Configs Load Error: %v", credit_control.SessionProxyServiceName, err)
	} else {
		initMethod = InitMethod(configsPtr.Gy.InitMethod)
	}
	method := diameter.GetValueOrEnv(GyInitMethodFlag, GyInitMethodEnv, "")
	switch method {
	case "per_session":
		return PerSessionInit
	case "per_key":
		return PerKeyInit
	case "":
	default:
		log.Printf("Invalid Gy Init Method specified: %s, will use %d", method, initMethod)
	}
	return initMethod
}

// TODO: refactor those functions to make it more simple
// GetOCSConfiguration returns the server configuration for the set OCS
func GetOCSConfiguration() []*diameter.DiameterServerConfig {
	configsPtr := &mconfig.SessionProxyConfig{}
	err := managed_configs.GetServiceConfigs(credit_control.SessionProxyServiceName, configsPtr)
	if err != nil || !validGyConfig(configsPtr) {
		log.Printf("%s Managed Gy Server Configs Load Error: %v", credit_control.SessionProxyServiceName, err)
		return []*diameter.DiameterServerConfig{
			{
				DiameterServerConnConfig: diameter.DiameterServerConnConfig{
					Addr:      diameter.GetValueOrEnv(diameter.AddrFlag, OCSAddrEnv, "127.0.0.1:3869"),
					Protocol:  diameter.GetValueOrEnv(diameter.NetworkFlag, GyNetworkEnv, "tcp"),
					LocalAddr: diameter.GetValueOrEnv(diameter.LocalAddrFlag, GyLocalAddr, ""),
				},
				DestHost:          diameter.GetValueOrEnv(diameter.DestHostFlag, OCSHostEnv, ""),
				DestRealm:         diameter.GetValueOrEnv(diameter.DestRealmFlag, OCSRealmEnv, ""),
				DisableDestHost:   diameter.GetBoolValueOrEnv(diameter.DisableDestHostFlag, DisableDestHostEnv, false),
				OverwriteDestHost: diameter.GetBoolValueOrEnv(diameter.OverwriteDestHostFlag, OverwriteDestHostEnv, false),
			},
		}
	}

	gyConfigs := configsPtr.GetGy().GetServers()
	//TODO: remove this once backwards compatibility is not needed for the field server
	if len(gyConfigs) == 0 {
		server := configsPtr.GetGy().GetServer()
		if server == nil {
			log.Print("Server configuration for Gy servers not found!!")
		} else {
			gyConfigs = append(gyConfigs, server)
			log.Print("Gy Server configuration using legacy swagger attribute Server (not Servers)")
		}
	}

	// Iterate over the slice of servers. VarEnv will apply only to index 0
	diamServerConfigs := []*diameter.DiameterServerConfig{}
	for i, gyCfg := range gyConfigs {
		diamSrvCfg := &diameter.DiameterServerConfig{
			DiameterServerConnConfig: diameter.DiameterServerConnConfig{
				Addr:      diameter.GetValueOrEnv(diameter.AddrFlag, OCSAddrEnv, gyCfg.GetAddress(), i),
				Protocol:  diameter.GetValueOrEnv(diameter.NetworkFlag, GyNetworkEnv, gyCfg.GetProtocol(), i),
				LocalAddr: diameter.GetValueOrEnv(diameter.LocalAddrFlag, GyLocalAddr, gyCfg.GetLocalAddress(), i),
			},
			DestHost:          diameter.GetValueOrEnv(diameter.DestHostFlag, OCSHostEnv, gyCfg.GetDestHost(), i),
			DestRealm:         diameter.GetValueOrEnv(diameter.DestRealmFlag, OCSRealmEnv, gyCfg.GetDestRealm(), i),
			DisableDestHost:   diameter.GetBoolValueOrEnv(diameter.DisableDestHostFlag, DisableDestHostEnv, gyCfg.GetDisableDestHost(), i),
			OverwriteDestHost: diameter.GetBoolValueOrEnv(diameter.OverwriteDestHostFlag, OverwriteDestHostEnv, gyCfg.GetOverwriteDestHost(), i),
		}
		diamServerConfigs = append(diamServerConfigs, diamSrvCfg)
	}
	return diamServerConfigs
}

// GetGyClientConfiguration returns the client diameter configuration
func GetGyClientConfiguration() []*diameter.DiameterClientConfig {
	var retries uint32 = 1
	configsPtr := &mconfig.SessionProxyConfig{}
	err := managed_configs.GetServiceConfigs(credit_control.SessionProxyServiceName, configsPtr)
	if err != nil {
		log.Printf("%s Managed Gy Client Configs Load Error: %v", credit_control.SessionProxyServiceName, err)
		return []*diameter.DiameterClientConfig{
			{
				Host:               diameter.GetValueOrEnv(diameter.HostFlag, GyDiamHostEnv, diameter.DiamHost),
				Realm:              diameter.GetValueOrEnv(diameter.RealmFlag, GyDiamRealmEnv, diameter.DiamRealm),
				ProductName:        diameter.GetValueOrEnv(diameter.ProductFlag, GyDiamProductEnv, diameter.DiamProductName),
				AppID:              diam.CHARGING_CONTROL_APP_ID,
				WatchdogInterval:   diameter.DefaultWatchdogIntervalSeconds,
				RetryCount:         uint(retries),
				SupportedVendorIDs: diameter.GetValueOrEnv("", GySupportedVendorIDsEnv, ""),
				ServiceContextId:   diameter.GetValueOrEnv("", GyServiceContextIdEnv, ""),
			},
		}
	}

	diamClientsConfigs := []*diameter.DiameterClientConfig{}
	gyConfigs := configsPtr.GetGy().GetServers()
	//TODO: remove this once backwards compatibility is not needed for the field server
	if len(gyConfigs) == 0 {
		server := configsPtr.GetGy().GetServer()
		if server == nil {
			log.Print("Client configuration for Gy servers not found!!")
		} else {
			gyConfigs = append(gyConfigs, server)
			log.Print("Gy Client configuration using legacy swagger attribute Server (not Servers)")
		}
	}
	for i, gyCfg := range gyConfigs {
		retries = gyCfg.GetRetryCount()
		if retries < 1 {
			log.Printf("Invalid Gy Server Retry Count for server (%s): %d, must be >0. Will be set to 1", gyCfg.GetAddress(), retries)
			retries = 1
		}

		wdInterval := gyCfg.GetWatchdogInterval()
		if wdInterval == 0 {
			wdInterval = diameter.DefaultWatchdogIntervalSeconds
		}
		diamCliCfg := &diameter.DiameterClientConfig{
			Host:               diameter.GetValueOrEnv(diameter.HostFlag, GyDiamHostEnv, gyCfg.GetHost(), i),
			Realm:              diameter.GetValueOrEnv(diameter.RealmFlag, GyDiamRealmEnv, gyCfg.GetRealm(), i),
			ProductName:        diameter.GetValueOrEnv(diameter.ProductFlag, GyDiamProductEnv, gyCfg.GetProductName(), i),
			AppID:              diam.CHARGING_CONTROL_APP_ID,
			WatchdogInterval:   uint(wdInterval),
			RetryCount:         uint(retries),
			SupportedVendorIDs: diameter.GetValueOrEnv("", GySupportedVendorIDsEnv, "", i),
			ServiceContextId:   diameter.GetValueOrEnv("", GyServiceContextIdEnv, "", i),
		}
		diamClientsConfigs = append(diamClientsConfigs, diamCliCfg)
	}
	return diamClientsConfigs
}

func GetGyGlobalConfig() *GyGlobalConfig {
	configsPtr := &mconfig.SessionProxyConfig{}
	err := managed_configs.GetServiceConfigs(credit_control.SessionProxyServiceName, configsPtr)
	siStr := diameter.GetValueOrEnv(OCSServiceIdentifierFlag, OCSServiceIdentifierEnv, "")
	if err != nil || !validGyConfig(configsPtr) {
		log.Printf("%s Managed Gy Server Configs Load Error: %v", credit_control.SessionProxyServiceName, err)
		return &GyGlobalConfig{
			OCSOverwriteApn:      diameter.GetValueOrEnv(OCSApnOverwriteFlag, OCSApnOverwriteEnv, ""),
			OCSServiceIdentifier: siStr,
			DisableGy:            false,
		}
	}

	return &GyGlobalConfig{
		OCSOverwriteApn:      diameter.GetValueOrEnv(OCSApnOverwriteFlag, OCSApnOverwriteEnv, configsPtr.GetGy().GetOverwriteApn()),
		OCSServiceIdentifier: siStr,
		DisableGy:            configsPtr.GetGy().GetDisableGy(),
		VirtualApnRules:      credit_control.GenerateVirtualApnRules(configsPtr.GetGy().GetVirtualApnRules()),
	}
}

// check if required fields related to Gy are valid in the config
func validGyConfig(config *mconfig.SessionProxyConfig) bool {
	if config == nil || config.Gy == nil ||
		(config.Gy.Server == nil && len(config.Gy.Servers) == 0) ||
		(config.Gy.Server != nil && config.Gy.Server.Address == "") {
		return false
	}
	for _, server := range config.Gy.Servers {
		if server.Address == "" {
			return false
		}
	}
	return true
}
