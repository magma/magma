/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_utils

import (
	"magma/lte/cloud/go/protos"
	orc8rprotos "magma/orc8r/cloud/go/protos"
)

const (
	defaultServerHost   = "magma.com"
	defaultMaxUlBitRate = uint64(100000000)
	defaultMaxDlBitRate = uint64(200000000)
)

// GetTestSubscribers returns a slice of SubscriberData protos to be used
// for testing authentication.
func GetTestSubscribers() []*protos.SubscriberData {
	subs := make([]*protos.SubscriberData, 0)

	sub := &protos.SubscriberData{
		Sid:       &protos.SubscriberID{Id: "sub1"},
		NetworkId: &orc8rprotos.NetworkID{Id: "test"},
		Lte: &protos.LTESubscription{
			State:    protos.LTESubscription_ACTIVE,
			AuthAlgo: protos.LTESubscription_MILENAGE,
			AuthKey:  []byte("\x8b\xafG?/\x8f\xd0\x94\x87\xcc\xcb\xd7\t|hb"),
			AuthOpc:  []byte("\x8e'\xb6\xaf\x0ei.u\x0f2fz;\x14`]"),
		},
		State: &protos.SubscriberState{
			LteAuthNextSeq:          7350,
			TgppAaaServerName:       defaultServerHost,
			TgppAaaServerRegistered: false,
		},
		Non_3Gpp: &protos.Non3GPPUserProfile{
			Msisdn:              "12345",
			Non_3GppIpAccess:    protos.Non3GPPUserProfile_NON_3GPP_SUBSCRIPTION_ALLOWED,
			Non_3GppIpAccessApn: protos.Non3GPPUserProfile_NON_3GPP_APNS_ENABLE,
			Ambr: &protos.AggregatedMaximumBitrate{
				MaxBandwidthUl: uint32(defaultMaxUlBitRate),
				MaxBandwidthDl: uint32(defaultMaxDlBitRate),
			},
			ApnConfig: &protos.APNConfiguration{
				ContextId:        10,
				ServiceSelection: "*",
				QosProfile: &protos.APNConfiguration_QoSProfile{
					ClassId:                 7,
					PriorityLevel:           3,
					PreemptionCapability:    true,
					PreemptionVulnerability: true,
				},
				Ambr: &protos.AggregatedMaximumBitrate{
					MaxBandwidthUl: uint32(defaultMaxUlBitRate),
					MaxBandwidthDl: uint32(defaultMaxDlBitRate),
				},
				Pdn: protos.APNConfiguration_IPV6,
			},
		},
		SubProfile: "test_profile",
	}
	subs = append(subs, sub)

	sub = &protos.SubscriberData{
		Sid:       &protos.SubscriberID{Id: "empty_sub"},
		NetworkId: &orc8rprotos.NetworkID{Id: "test"},
	}
	subs = append(subs, sub)

	sub = &protos.SubscriberData{
		Sid:       &protos.SubscriberID{Id: "missing_auth_key"},
		NetworkId: &orc8rprotos.NetworkID{Id: "test"},
		Lte: &protos.LTESubscription{
			State:    protos.LTESubscription_ACTIVE,
			AuthAlgo: protos.LTESubscription_MILENAGE,
			AuthOpc:  []byte("\x8e'\xb6\xaf\x0ei.u\x0f2fz;\x14`]"),
		},
		State: &protos.SubscriberState{
			LteAuthNextSeq:    7350,
			TgppAaaServerName: defaultServerHost,
		},
	}
	subs = append(subs, sub)

	return subs
}
