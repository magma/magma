/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"fmt"

	mcfgprotos "magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/diameter"
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
func GetSwxProxyConfig() *SwxProxyConfig {
	configsPtr := &mcfgprotos.SwxConfig{}
	hlrPlmnIds := map[string]PlmnIdVal{}
	err := mconfig.GetServiceConfigs(SwxProxyServiceName, configsPtr)

	if err != nil || configsPtr.Server == nil {
		glog.V(2).Infof("%s Managed Configs Load Error: %v", SwxProxyServiceName, err)

		return &SwxProxyConfig{
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
			HlrPlmnIds:            hlrPlmnIds,
		}
	}

	glog.V(2).Infof("Loaded %s configs: %+v", SwxProxyServiceName, *configsPtr)

	for _, plmnid := range configsPtr.HlrPlmnIds {
		glog.Infof("Adding HLR PLMN ID: %s", plmnid)
		l := len(plmnid)
		switch l {
		case 5:
			hlrPlmnIds[plmnid] = PlmnIdVal{l5: true}
		case 6:
			plmnid5 := plmnid[:5]
			val, _ := hlrPlmnIds[plmnid5]
			val.b6 = plmnid[5]
			hlrPlmnIds[plmnid5] = val
		default:
			glog.Warningf("Invalid HLR PLMN ID: %s", plmnid)
		}
	}

	ttl := configsPtr.CacheTTLSeconds
	if ttl < uint32(cache.DefaultGcInterval.Seconds()) {
		ttl = uint32(cache.DefaultTtl.Seconds())
	}
	return &SwxProxyConfig{
		ClientCfg: &diameter.DiameterClientConfig{
			Host:             diameter.GetValueOrEnv(diameter.HostFlag, SwxDiamHostEnv, configsPtr.GetServer().GetHost()),
			Realm:            diameter.GetValueOrEnv(diameter.RealmFlag, SwxDiamRealmEnv, configsPtr.GetServer().GetRealm()),
			ProductName:      diameter.GetValueOrEnv(diameter.ProductFlag, SwxDiamProductEnv, configsPtr.GetServer().GetProductName()),
			Retransmits:      uint(configsPtr.GetServer().GetRetransmits()),
			WatchdogInterval: uint(configsPtr.GetServer().GetWatchdogInterval()),
			RetryCount:       uint(configsPtr.GetServer().GetRetryCount()),
		},
		ServerCfg: &diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
			Addr:      diameter.GetValueOrEnv(diameter.AddrFlag, HSSAddrEnv, configsPtr.GetServer().GetAddress()),
			Protocol:  diameter.GetValueOrEnv(diameter.NetworkFlag, SwxNetworkEnv, configsPtr.GetServer().GetProtocol()),
			LocalAddr: diameter.GetValueOrEnv(diameter.LocalAddrFlag, SwxLocalAddrEnv, configsPtr.GetServer().GetLocalAddress())},
			DestHost:          diameter.GetValueOrEnv(diameter.DestHostFlag, HSSHostEnv, configsPtr.GetServer().GetDestHost()),
			DestRealm:         diameter.GetValueOrEnv(diameter.DestRealmFlag, HSSRealmEnv, configsPtr.GetServer().GetDestRealm()),
			DisableDestHost:   diameter.GetBoolValueOrEnv(diameter.DisableDestHostFlag, DisableDestHostEnv, configsPtr.GetServer().GetDisableDestHost()),
			OverwriteDestHost: diameter.GetBoolValueOrEnv(diameter.OverwriteDestHostFlag, OverwriteDestHostEnv, configsPtr.GetServer().GetOverwriteDestHost()),
		},
		VerifyAuthorization:   configsPtr.GetVerifyAuthorization(),
		RegisterOnAuth:        configsPtr.GetRegisterOnAuth(),
		DeriveUnregisterRealm: configsPtr.GetDeriveUnregisterRealm(),
		CacheTTLSeconds:       ttl,
		HlrPlmnIds:            hlrPlmnIds,
	}
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

func (config *SwxProxyConfig) IsHlrClient(imsi string) bool {
	if config != nil && len(config.HlrPlmnIds) > 0 {
		if val, ok := config.HlrPlmnIds[string(imsi)[:5]]; ok && (val.l5 || (len(imsi) > 5 && val.b6 == imsi[6])) {
			return true
		}
	}
	return false
}
