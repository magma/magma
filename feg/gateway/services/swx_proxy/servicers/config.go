/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	mcfgprotos "magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/mconfig"

	"github.com/golang/glog"
)

const (
	SwxProxyServiceName = "swx_proxy"

	HSSAddrEnv        = "HSS_ADDR"
	SwxNetworkEnv     = "SWX_NETWORK"
	SwxDiamHostEnv    = "SWX_DIAM_HOST"
	SwxDiamRealmEnv   = "SWX_DIAM_REALM"
	SwxDiamProductEnv = "SWX_DIAM_PRODUCT"
	SwxLocalAddrEnv   = "SWX_LOCAL_ADDR"
	HSSHostEnv        = "HSS_HOST"
	HSSRealmEnv       = "HSS_REALM"

	DefaultSwxDiamRealm = "epc.mnc070.mcc722.3gppnetwork.org"
	DefaultSwxDiamHost  = "feg-swx.epc.mnc070.mcc722.3gppnetwork.org"
)

// Get GetSwxProxyConfigs returns the server config for an HSS based on the
// the values in mconfig or default values provided
func GetSwxProxyConfigs() (*diameter.DiameterClientConfig, *diameter.DiameterServerConfig) {
	configsPtr := &mcfgprotos.SwxConfig{}
	err := mconfig.GetServiceConfigs(SwxProxyServiceName, configsPtr)
	if err != nil || configsPtr.Server == nil {
		glog.V(2).Infof("%s Managed Configs Load Error: %v", SwxProxyServiceName, err)
		return &diameter.DiameterClientConfig{
				Host:        diameter.GetValueOrEnv(diameter.HostFlag, SwxDiamHostEnv, DefaultSwxDiamHost),
				Realm:       diameter.GetValueOrEnv(diameter.RealmFlag, SwxDiamRealmEnv, DefaultSwxDiamRealm),
				ProductName: diameter.GetValueOrEnv(diameter.ProductFlag, SwxDiamProductEnv, diameter.DiamProductName),
			},
			&diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
				Addr:      diameter.GetValueOrEnv(diameter.AddrFlag, HSSAddrEnv, ""),
				Protocol:  diameter.GetValueOrEnv(diameter.NetworkFlag, SwxNetworkEnv, "sctp"),
				LocalAddr: diameter.GetValueOrEnv(diameter.LocalAddrFlag, SwxLocalAddrEnv, "")},
				DestHost:  diameter.GetValueOrEnv(diameter.DestHostFlag, HSSHostEnv, ""),
				DestRealm: diameter.GetValueOrEnv(diameter.DestRealmFlag, HSSRealmEnv, ""),
			}
	}

	glog.V(2).Infof("Loaded %s configs: %+v", SwxProxyServiceName, *configsPtr)

	return &diameter.DiameterClientConfig{
			Host:             diameter.GetValueOrEnv(diameter.HostFlag, SwxDiamHostEnv, configsPtr.Server.Host),
			Realm:            diameter.GetValueOrEnv(diameter.RealmFlag, SwxDiamRealmEnv, configsPtr.Server.Realm),
			ProductName:      diameter.GetValueOrEnv(diameter.ProductFlag, SwxDiamProductEnv, configsPtr.Server.ProductName),
			Retransmits:      uint(configsPtr.Server.Retransmits),
			WatchdogInterval: uint(configsPtr.Server.WatchdogInterval),
			RetryCount:       uint(configsPtr.Server.RetryCount),
		},
		&diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
			Addr:      diameter.GetValueOrEnv(diameter.AddrFlag, HSSAddrEnv, configsPtr.Server.Address),
			Protocol:  diameter.GetValueOrEnv(diameter.NetworkFlag, SwxNetworkEnv, configsPtr.Server.Protocol),
			LocalAddr: diameter.GetValueOrEnv(diameter.LocalAddrFlag, SwxLocalAddrEnv, configsPtr.Server.LocalAddress)},
			DestHost:  diameter.GetValueOrEnv(diameter.DestHostFlag, HSSHostEnv, configsPtr.Server.DestHost),
			DestRealm: diameter.GetValueOrEnv(diameter.DestRealmFlag, HSSRealmEnv, configsPtr.Server.DestRealm),
		}
}
