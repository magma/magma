package session

import (
	"fbc/cwf/radius/counters"
)

var (
	// ReadSessionState counts reading session state from storage
	ReadSessionState = counters.NewOperation("read_session_state")

	// WriteSessionState counts writing session state from storage
	WriteSessionState = counters.NewOperation("write_session_state")

	// ResetSessionState counts reseting session state from storage
	ResetSessionState = counters.NewOperation("reset_session_state")
)
