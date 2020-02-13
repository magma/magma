/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

func NewDefaultNetworkFederationConfigs() *NetworkFederationConfigs {
	return &NetworkFederationConfigs{
		AaaServer:        newDefaultAaaServer(),
		EapAka:           newDefaultEapAka(),
		Gx:               newDefaultGx(),
		Gy:               newDefaultGy(),
		Health:           newDefaultHealth(),
		Hss:              newDefaultHss(),
		S6a:              newDefaultS6a(),
		ServedNetworkIds: newDefaultServedNetworkIds(),
		Swx:              newDefaultSwx(),
		Csfb:             newDefaultCsfb(),
	}
}

func NewDefaultModifiedNetworkFederationConfigs() *NetworkFederationConfigs {
	configs := NewDefaultNetworkFederationConfigs()
	configs.AaaServer = &AaaServer{
		AccountingEnabled:    true,
		CreateSessionOnAuth:  true,
		IDLESessionTimeoutMs: 11600000,
	}
	return configs
}

func NewDefaultGatewayFederationConfig() *GatewayFederationConfigs {
	return &GatewayFederationConfigs{
		AaaServer:        newDefaultAaaServer(),
		EapAka:           newDefaultEapAka(),
		Gx:               newDefaultGx(),
		Gy:               newDefaultGy(),
		Health:           newDefaultHealth(),
		Hss:              newDefaultHss(),
		S6a:              newDefaultS6a(),
		ServedNetworkIds: newDefaultServedNetworkIds(),
		Swx:              newDefaultSwx(),
		Csfb:             newDefaultCsfb(),
	}
}

func NewDefaultFederatedNetworkConfigs() *FederatedNetworkConfigs {
	fegNetworkID := "n1"
	return &FederatedNetworkConfigs{
		FegNetworkID: &fegNetworkID,
	}
}

func newDefaultAaaServer() *AaaServer {
	return &AaaServer{
		AccountingEnabled:    false,
		CreateSessionOnAuth:  false,
		IDLESessionTimeoutMs: 21600000,
	}
}

func newDefaultEapAka() *EapAka {
	return &EapAka{
		PlmnIds: []string{"123456"},
		Timeout: &EapAkaTimeouts{
			ChallengeMs:            20000,
			ErrorNotificationMs:    10000,
			SessionAuthenticatedMs: 5000,
			SessionMs:              43200000,
		},
	}
}

func newDefaultGx() *Gx {
	return &Gx{
		Server: newDefaultDiameterClientConfigs(),
	}
}

func newDefaultGy() *Gy {
	initMethod := uint32(float32(1))
	return &Gy{
		InitMethod: &initMethod,
		Server:     newDefaultDiameterClientConfigs(),
	}
}

func newDefaultDiameterClientConfigs() *DiameterClientConfigs {
	return &DiameterClientConfigs{
		Address:          "foo.bar.com:5555",
		DestHost:         "magma-fedgw.magma.com",
		DestRealm:        "magma.com",
		DisableDestHost:  false,
		Host:             "string",
		LocalAddress:     ":56789",
		ProductName:      "string",
		Protocol:         "tcp",
		Realm:            "string",
		Retransmits:      0,
		RetryCount:       0,
		WatchdogInterval: 0,
	}
}

func newDefaultHealth() *Health {
	return &Health{
		CloudDisablePeriodSecs:   10,
		CPUUtilizationThreshold:  0.9,
		HealthServices:           []string{"S6A_PROXY", "SESSION_PROXY", "SWX_PROXY"},
		LocalDisablePeriodSecs:   1,
		MemoryAvailableThreshold: 0.9,
		MinimumRequestThreshold:  1,
		RequestFailureThreshold:  0.5,
		UpdateFailureThreshold:   3,
		UpdateIntervalSecs:       10,
	}
}

func newDefaultHss() *Hss {
	return &Hss{
		DefaultSubProfile: &SubscriptionProfile{
			MaxDlBitRate: 200000000,
			MaxUlBitRate: 100000000,
		},
		LteAuthAmf: []byte("gAA="),
		LteAuthOp:  []byte("EREREREREREREREREREREQ=="),
		Server: &DiameterServerConfigs{
			Address:      "foo.bar.com:5555",
			DestHost:     "magma-fedgw.magma.com",
			DestRealm:    "magma.com",
			LocalAddress: ":56789",
			Protocol:     "tcp",
		},
		StreamSubscribers: false,
		SubProfiles: map[string]SubscriptionProfile{
			"additionalProp1": {
				MaxDlBitRate: 200000000,
				MaxUlBitRate: 100000000,
			},
			"additionalProp2": {
				MaxDlBitRate: 200000000,
				MaxUlBitRate: 100000000,
			},
			"additionalProp3": {
				MaxDlBitRate: 200000000,
				MaxUlBitRate: 100000000,
			},
		},
	}
}

func newDefaultS6a() *S6a {
	return &S6a{
		Server: newDefaultDiameterClientConfigs(),
	}
}

func newDefaultSwx() *Swx {
	return &Swx{
		CacheTTLSeconds:       10800,
		DeriveUnregisterRealm: false,
		RegisterOnAuth:        false,
		Server:                newDefaultDiameterClientConfigs(),
		VerifyAuthorization:   false,
	}
}

func newDefaultCsfb() *Csfb {
	return &Csfb{
		Client: &SctpClientConfigs{
			LocalAddress:  "foo.bar.com:5555",
			ServerAddress: ":56789"},
	}
}

func newDefaultServedNetworkIds() []string {
	return []string{"string"}
}
