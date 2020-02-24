/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package test_utils

import (
	"magma/lte/cloud/go/protos"
	orc8rprotos "magma/orc8r/lib/go/protos"
)

const (
	defaultServerHostSubscribers = "magma.com"
	defaultMaxUlBitRate          = uint64(100000000)
	defaultMaxDlBitRate          = uint64(200000000)
)

// GetTestSubscribers returns SubscriberData protos with different settings
// to be used for testing authentication. More users can be added.
func GetTestSubscribers() []*protos.SubscriberData {
	subs := make([]*protos.SubscriberData, 0)

	// Default subscriber
	sub := generateDefaultSub("sub1")
	subs = append(subs, sub)

	// Default Subs with a blank AAA server
	sub = generateDefaultSub("sub1_noAAAsrv")
	sub.State.TgppAaaServerName = ""
	subs = append(subs, sub)

	// Empty sub
	sub = &protos.SubscriberData{
		Sid:       &protos.SubscriberID{Id: "empty_sub"},
		NetworkId: &orc8rprotos.NetworkID{Id: "test"},
	}
	subs = append(subs, sub)

	// Subscriber without auth key
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
			TgppAaaServerName: defaultServerHostSubscribers,
		},
	}
	subs = append(subs, sub)

	return subs
}

func generateDefaultSub(subscriberID string) *protos.SubscriberData {
	// Default user
	sub := &protos.SubscriberData{
		Sid:       &protos.SubscriberID{Id: subscriberID},
		NetworkId: &orc8rprotos.NetworkID{Id: "test"},
		Lte: &protos.LTESubscription{
			State:    protos.LTESubscription_ACTIVE,
			AuthAlgo: protos.LTESubscription_MILENAGE,
			AuthKey:  []byte("\x8b\xafG?/\x8f\xd0\x94\x87\xcc\xcb\xd7\t|hb"),
			AuthOpc:  []byte("\x8e'\xb6\xaf\x0ei.u\x0f2fz;\x14`]"),
		},
		State: &protos.SubscriberState{
			LteAuthNextSeq:          7350,
			TgppAaaServerName:       defaultServerHostSubscribers,
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
			ApnConfig: []*protos.APNConfiguration{&protos.APNConfiguration{
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
			}},
		},
		SubProfile: "test_profile",
	}

	return sub
}
