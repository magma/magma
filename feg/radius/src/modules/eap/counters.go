/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package eap

import (
	"fbc/cwf/radius/monitoring"
)

var (
	// ExtractEapPacket extraction EAP-Message from RADIUS message
	ExtractEapPacket = monitoring.NewOperation("eap_extract_packet_from_radius")

	// RestoreProtocolState restore eap state from storage
	RestoreProtocolState = monitoring.NewOperation("eap_restore_state")

	// HandleEapPacket handling EAP packet
	HandleEapPacket = monitoring.NewOperation("eap_handle")

	// PersistProtocolState writing new state, after handling, to storage
	PersistProtocolState = monitoring.NewOperation("eap_persist_state")
)
