/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
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

// GetPCRFConfiguration returns a slice containing all configuration for all known PCRF
func GetPCRFConfiguration() []*diameter.DiameterServerConfig {
	configsPtr := &mconfig.SessionProxyConfig{}
	err := managed_configs.GetServiceConfigs(credit_control.SessionProxyServiceName, configsPtr)
	if err != nil || !validGxConfig(configsPtr) {
		log.Printf("%s Managed Gx PCRF Configs Load Error: %v", credit_control.SessionProxyServiceName, err)
		return []*diameter.DiameterServerConfig{
			&diameter.DiameterServerConfig{
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
				Addr:      getValueOrEnvForIndexZero(i, diameter.AddrFlag, PCRFAddrEnv, gxCfg.GetAddress()),
				Protocol:  getValueOrEnvForIndexZero(i, diameter.NetworkFlag, GxNetworkEnv, gxCfg.GetProtocol()),
				LocalAddr: getValueOrEnvForIndexZero(i, diameter.LocalAddrFlag, GxLocalAddr, gxCfg.GetLocalAddress()),
			},
			DestHost:          getValueOrEnvForIndexZero(i, diameter.DestHostFlag, PCRFHostEnv, gxCfg.GetDestHost()),
			DestRealm:         getValueOrEnvForIndexZero(i, diameter.DestRealmFlag, PCRFRealmEnv, gxCfg.GetDestHost()),
			DisableDestHost:   getBoolValueOrEnvForIndexZero(i, diameter.DisableDestHostFlag, DisableDestHostEnv, gxCfg.GetDisableDestHost()),
			OverwriteDestHost: getBoolValueOrEnvForIndexZero(i, diameter.OverwriteDestHostFlag, OverwriteDestHostEnv, gxCfg.GetOverwriteDestHost()),
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
			&diameter.DiameterClientConfig{
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
		diamCliCfg := &diameter.DiameterClientConfig{
			Host:               getValueOrEnvForIndexZero(i, diameter.HostFlag, GxDiamHostEnv, gxCfg.GetHost()),
			Realm:              getValueOrEnvForIndexZero(i, diameter.RealmFlag, GxDiamRealmEnv, gxCfg.GetRealm()),
			ProductName:        getValueOrEnvForIndexZero(i, diameter.ProductFlag, GxDiamProductEnv, gxCfg.GetProductName()),
			AppID:              diam.GX_CHARGING_CONTROL_APP_ID,
			WatchdogInterval:   diameter.DefaultWatchdogIntervalSeconds,
			RetryCount:         uint(retries),
			SupportedVendorIDs: getValueOrEnvForIndexZero(i, "", GxSupportedVendorIDsEnv, ""),
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
		PCFROverwriteApn: configsPtr.GetGx().OverwriteApn,
	}

}

// validGxConfig check if required fields related to Gx are valid in the config
func validGxConfig(config *mconfig.SessionProxyConfig) bool {
	if config == nil || config.Gx == nil || config.Gx.Server == nil || config.Gx.Server.Address == "" {
		return false
	}
	return true
}

// getValueOrEnvForIndexZero gets the Value or Env for the index 0, otherwise it returns Value(string)
func getValueOrEnvForIndexZero(index int, flagName, envVariable, defaultValue string) string {
	if index == 0 {
		return diameter.GetValueOrEnv(flagName, envVariable, defaultValue)
	}
	return defaultValue
}

// getBoolValueOrEnvForIndexZero gets the Value or Env for the index 0, otherwise it returns Value(bool)
func getBoolValueOrEnvForIndexZero(index int, flagName string, envVariable string, defaultValue bool) bool {
	if index == 0 {
		return diameter.GetBoolValueOrEnv(flagName, envVariable, defaultValue)
	}
	return defaultValue
}
