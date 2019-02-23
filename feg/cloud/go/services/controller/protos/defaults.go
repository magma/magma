/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package protos

import (
	"github.com/golang/protobuf/proto"
)

var defaultConfig = Config{
	S6A: &S6AConfig{
		Server: &DiamClientConfig{
			Protocol:         "sctp",
			Retransmits:      3,
			WatchdogInterval: 1,
			RetryCount:       5,
			ProductName:      "magma",
			Host:             "magma-fedgw.magma.com",
			Realm:            "magma.com",
		},
	},
	Gx: &GxConfig{
		Server: &DiamClientConfig{
			Protocol:         "tcp",
			Retransmits:      3,
			WatchdogInterval: 1,
			RetryCount:       5,
			ProductName:      "magma",
			Host:             "magma-fedgw.magma.com",
			Realm:            "magma.com",
		},
	},
	Gy: &GyConfig{
		Server: &DiamClientConfig{
			Protocol:         "tcp",
			Retransmits:      3,
			WatchdogInterval: 1,
			RetryCount:       5,
			ProductName:      "magma",
			Host:             "magma-fedgw.magma.com",
			Realm:            "magma.com",
		},
		InitMethod: GyInitMethod_PER_SESSION,
	},
	Hss: &HSSConfig{
		Server: &DiamServerConfig{
			Protocol:  "tcp",
			DestHost:  "magma.com",
			DestRealm: "magma.com",
		},
		LteAuthOp:  []byte("EREREREREREREREREREREQ=="),
		LteAuthAmf: []byte("gA"),
		DefaultSubProfile: &HSSConfig_SubscriptionProfile{
			MaxUlBitRate: 100000000, // 100 Mbps
			MaxDlBitRate: 200000000, // 200 Mbps
		},
		SubProfiles: make(map[string]*HSSConfig_SubscriptionProfile),
	},
	Swx: &SwxConfig{
		Server: &DiamClientConfig{
			Protocol:         "sctp",
			Retransmits:      3,
			WatchdogInterval: 1,
			RetryCount:       5,
			ProductName:      "magma",
			Host:             "magma-fedgw.magma.com",
			Realm:            "magma.com",
		},
	},
	ServedNetworkIds: []string{},
}

func NewDefaultNetworkConfig() *Config {
	return proto.Clone(&defaultConfig).(*Config)
}

func NewDefaultGatewayConfig() *Config {
	return proto.Clone(&defaultConfig).(*Config)
}
