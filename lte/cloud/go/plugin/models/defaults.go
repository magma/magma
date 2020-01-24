/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import "github.com/go-openapi/swag"

func NewDefaultTDDNetworkConfig() *NetworkCellularConfigs {
	return &NetworkCellularConfigs{
		Ran: &NetworkRanConfigs{
			BandwidthMhz: 20,
			TddConfig: &NetworkRanConfigsTddConfig{
				Earfcndl:               44590,
				SubframeAssignment:     2,
				SpecialSubframePattern: 7,
			},
		},
		Epc: &NetworkEpcConfigs{
			Mcc: "001",
			Mnc: "01",
			Tac: 1,
			// 16 bytes of \x11
			LteAuthOp:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf: []byte("\x80\x00"),

			RelayEnabled:             swag.Bool(false),
			CloudSubscriberdbEnabled: false,
			DefaultRuleID:            "",
		},
	}
}

func NewDefaultFDDNetworkConfig() *NetworkCellularConfigs {
	return &NetworkCellularConfigs{
		Ran: &NetworkRanConfigs{
			BandwidthMhz: 20,
			FddConfig: &NetworkRanConfigsFddConfig{
				Earfcndl: 1,
				Earfcnul: 18001,
			},
		},
		Epc: &NetworkEpcConfigs{
			Mcc: "001",
			Mnc: "01",
			Tac: 1,
			// 16 bytes of \x11
			LteAuthOp:                []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:               []byte("\x80\x00"),
			RelayEnabled:             swag.Bool(false),
			CloudSubscriberdbEnabled: false,
			DefaultRuleID:            "",
		},
	}
}

func NewDefaultSubscriberConfig() *NetworkSubscriberConfig {
	return &NetworkSubscriberConfig{
		NetworkWideBaseNames: []BaseName{"base1"},
		NetworkWideRuleNames: []string{"rule1"},
	}
}

func NewDefaultEnodebStatus() *EnodebState {
	return &EnodebState{
		EnodebConfigured: swag.Bool(true),
		EnodebConnected:  swag.Bool(true),
		GpsConnected:     swag.Bool(true),
		GpsLatitude:      swag.String("1.1"),
		GpsLongitude:     swag.String("2.2"),
		OpstateEnabled:   swag.Bool(true),
		RfTxOn:           swag.Bool(true),
		RfTxDesired:      swag.Bool(false),
		PtpConnected:     swag.Bool(false),
		MmeConnected:     swag.Bool(true),
		FsmState:         swag.String("TEST"),
	}
}
