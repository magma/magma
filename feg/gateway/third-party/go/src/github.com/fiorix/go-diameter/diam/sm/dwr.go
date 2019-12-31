// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/sm/smparser"
)

// handleDWR handles Device-Watchdog-Request messages.
//
// If mandatory AVPs such as Origin-Host or Origin-Realm are missing,
// we ignore the message.
//
// See RFC 6733 section 5.5 for details.
func handleDWR(sm *StateMachine) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		dwr := new(smparser.DWR)
		err := dwr.Parse(m)
		if err != nil {
			sm.Error(&diam.ErrorReport{
				Conn:    c,
				Message: m,
				Error:   err,
			})
			return
		}
		a := m.Answer(diam.Success)
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, sm.cfg.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, sm.cfg.OriginRealm)
		if sm.cfg.OriginStateID != 0 {
			stateid := datatype.Unsigned32(sm.cfg.OriginStateID)
			m.NewAVP(avp.OriginStateID, avp.Mbit, 0, stateid)
		}
		_, err = a.WriteTo(c)
		if err != nil {
			sm.Error(&diam.ErrorReport{
				Conn:    c,
				Message: m,
				Error:   err,
			})
		}
	}
}
