/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package session

import (
	"fbc/cwf/radius/monitoring/counters"
)

var (
	// ReadSessionState counts reading session state from storage
	ReadSessionState = counters.NewOperation("read_session_state")

	// WriteSessionState counts writing session state from storage
	WriteSessionState = counters.NewOperation("write_session_state")

	// ResetSessionState counts reseting session state from storage
	ResetSessionState = counters.NewOperation("reset_session_state")
)
