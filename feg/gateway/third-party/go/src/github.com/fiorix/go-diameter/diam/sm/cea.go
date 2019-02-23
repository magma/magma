// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/sm/smparser"
	"github.com/fiorix/go-diameter/diam/sm/smpeer"
)

// handleCEA handles Capabilities-Exchange-Answer messages.
func handleCEA(sm *StateMachine, errc chan error) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		cea := new(smparser.CEA)
		if err := cea.Parse(m, smparser.Client); err != nil {
			errc <- err
			return
		}
		meta := smpeer.FromCEA(cea)
		c.SetContext(smpeer.NewContext(c.Context(), meta))
		// Notify about peer passing the handshake.
		select {
		case sm.hsNotifyc <- c:
		default:
		}
		// Done receiving and validating this CEA.
		close(errc)
	}
}
