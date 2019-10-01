/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package eap

import (
	"fbc/cwf/radius/monitoring/counters"
)

var (
	// ExtractEapPacket extraction EAP-Message from RADIUS message
	ExtractEapPacket = counters.NewOperation("eap_extract_packet_from_radius")

	// RestoreProtocolState restore eap state from storage
	RestoreProtocolState = counters.NewOperation("eap_restore_state")

	// HandleEapPacket handling EAP packet
	HandleEapPacket = counters.NewOperation("eap_handle")

	// PersistProtocolState writing new state, after handling, to storage
	PersistProtocolState = counters.NewOperation("eap_persist_state")
)
