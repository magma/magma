/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/diameter"
	configs "magma/feg/gateway/mconfig"

	"github.com/golang/glog"
)

// HSS Flag Variables to overwrite default configs
const (
	hssServiceName       = "hss"
	hssDefaultProtocol   = "tcp"
	hssDefaultHost       = "magma.com"
	hssDefaultRealm      = "magma.com"
	hssDefaultLteAuthOp  = "EREREREREREREREREREREQ=="
	lteAuthOpFlag        = "lte_auth_op"
	hssDefaultLteAuthAmf = "gA"
	lteAuthAmfFlag       = "lte_auth_amf"
	maxUlBitRateFlag     = "max_ul_bit_rate"
	maxDlBitRateFlag     = "max_dl_bit_rate"
	defaultMaxUlBitRate  = uint64(100000000)
	defaultMaxDlBitRate  = uint64(200000000)
)

func init() {
	flag.String(lteAuthAmfFlag, hssDefaultLteAuthAmf, "Authentication management field for LTE")
	flag.String(lteAuthOpFlag, hssDefaultLteAuthOp, "Operator configuration field for LTE")
	flag.Uint64(maxUlBitRateFlag, defaultMaxUlBitRate, "Maximum uplink bit rate (AMBR-UL)")
	flag.Uint64(maxDlBitRateFlag, defaultMaxDlBitRate, "Maximum downlink bit rate (AMBR-DL)")
}

// GetHSSConfig returns the server config for an HSS based on the input flags
func GetHSSConfig() (*mconfig.HSSConfig, error) {
	serviceBaseName := filepath.Base(os.Args[0])
	serviceBaseName = strings.TrimSuffix(serviceBaseName, filepath.Ext(serviceBaseName))
	if hssServiceName != serviceBaseName {
		glog.Errorf(
			"NOTE: HSS Service name: %s does not match its managed configs key: %s\n",
			serviceBaseName, hssServiceName)
	}

	configsPtr := &mconfig.HSSConfig{}
	err := configs.GetServiceConfigs(hssServiceName, configsPtr)
	if err != nil || configsPtr.Server == nil || configsPtr.DefaultSubProfile == nil {
		glog.Errorf("%s Managed Configs Load Error: %v\n", hssServiceName, err)
		return &mconfig.HSSConfig{
			Server: &mconfig.DiamServerConfig{
				Address:      diameter.GetValue(diameter.AddrFlag, ""),
				Protocol:     diameter.GetValue(diameter.NetworkFlag, hssDefaultProtocol),
				LocalAddress: diameter.GetValue(diameter.LocalAddrFlag, ""),
				DestHost:     diameter.GetValue(diameter.DestHostFlag, hssDefaultHost),
				DestRealm:    diameter.GetValue(diameter.DestRealmFlag, hssDefaultRealm),
			},
			LteAuthOp:  []byte(diameter.GetValue(lteAuthOpFlag, hssDefaultLteAuthOp)),
			LteAuthAmf: []byte(diameter.GetValue(lteAuthAmfFlag, hssDefaultLteAuthAmf)),
			DefaultSubProfile: &mconfig.HSSConfig_SubscriptionProfile{
				MaxUlBitRate: diameter.GetValueUint64(maxUlBitRateFlag, defaultMaxUlBitRate),
				MaxDlBitRate: diameter.GetValueUint64(maxDlBitRateFlag, defaultMaxDlBitRate),
			},
			SubProfiles: make(map[string]*mconfig.HSSConfig_SubscriptionProfile),
		}, err
	}

	glog.V(2).Infof("Loaded %s configs: %+v\n", hssServiceName, *configsPtr)

	return &mconfig.HSSConfig{
		Server: &mconfig.DiamServerConfig{
			Address:      diameter.GetValue(diameter.AddrFlag, configsPtr.Server.Address),
			Protocol:     diameter.GetValue(diameter.NetworkFlag, configsPtr.Server.Protocol),
			LocalAddress: diameter.GetValue(diameter.LocalAddrFlag, configsPtr.Server.LocalAddress),
			DestHost:     diameter.GetValue(diameter.DestHostFlag, configsPtr.Server.DestHost),
			DestRealm:    diameter.GetValue(diameter.DestRealmFlag, configsPtr.Server.DestRealm),
		},
		LteAuthOp:  []byte(diameter.GetValue(lteAuthOpFlag, string(configsPtr.LteAuthOp))),
		LteAuthAmf: []byte(diameter.GetValue(lteAuthAmfFlag, string(configsPtr.LteAuthAmf))),
		DefaultSubProfile: &mconfig.HSSConfig_SubscriptionProfile{
			MaxUlBitRate: diameter.GetValueUint64(maxUlBitRateFlag, configsPtr.DefaultSubProfile.MaxUlBitRate),
			MaxDlBitRate: diameter.GetValueUint64(maxDlBitRateFlag, configsPtr.DefaultSubProfile.MaxDlBitRate),
		},
		SubProfiles: configsPtr.SubProfiles,
	}, nil
}
