package eap

import (
	"fbc/cwf/radius/counters"
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
