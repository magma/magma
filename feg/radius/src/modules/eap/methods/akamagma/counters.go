package akamagma

import (
	"fbc/cwf/radius/counters"
)

var (
	// MarshalProtocolState marshal eap state
	MarshalProtocolState = counters.NewOperation("eap_marshal_state")

	// UnmarshalProtocolState marshal eap state
	UnmarshalProtocolState = counters.NewOperation("eap_marshal_state")
)
