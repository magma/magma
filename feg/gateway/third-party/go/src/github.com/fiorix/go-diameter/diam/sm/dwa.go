// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/sm/smparser"
)

var dwaACK = struct{}{}

// handleDWA handles Device-Watchdog-Answer messages.
func handleDWA(sm *StateMachine, dwac chan struct{}) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		dwa := new(smparser.DWA)
		if err := dwa.Parse(m); err != nil {
			sm.Error(&diam.ErrorReport{
				Conn:    c,
				Message: m,
				Error:   err,
			})
			return
		}
		if dwa.ResultCode != diam.Success {
			return
		}
		select {
		case dwac <- dwaACK:
		default:
		}
	}
}
