/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import "github.com/go-openapi/strfmt"

func NewDefaultTDDNetworkConfig() *NetworkCellularConfigs {
	return &NetworkCellularConfigs{
		Ran: &NetworkRanConfigs{
			BandwidthMhz: uint32Ptr(20),
			TddConfig: &NetworkRanConfigsTddConfig{
				Earfcndl:               uint32Ptr(44590),
				SubframeAssignment:     uint32Ptr(2),
				SpecialSubframePattern: uint32Ptr(7),
			},
		},
		Epc: &NetworkEpcConfigs{
			Mcc: sPtr("001"),
			Mnc: sPtr("01"),
			Tac: uint32Ptr(1),
			// 16 bytes of \x11
			LteAuthOp:  bytesPtr([]byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11")),
			LteAuthAmf: bytesPtr([]byte("\x80\x00")),
		},
	}
}

func uint32Ptr(i uint32) *uint32 {
	return &i
}

func sPtr(s string) *string {
	return &s
}

func bytesPtr(b []byte) *strfmt.Base64 {
	s := strfmt.Base64(b)
	return &s
}
