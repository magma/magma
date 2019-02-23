// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smparser

import (
	"testing"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
)

func TestDWA(t *testing.T) {
	m := diam.NewMessage(diam.CapabilitiesExchange, 0, 0, 0, 0, nil)
	m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(diam.Success))
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	dwa := new(DWA)
	if err := dwa.Parse(m); err != nil {
		t.Fatal(err)
	}
	if dwa.ResultCode != diam.Success {
		t.Fatalf("Unexpected Result-Code. Want %d, have %d",
			diam.Success, dwa.ResultCode)
	}
	if dwa.OriginStateID != 1 {
		t.Fatalf("Unexpected Result-Code. Want 1, have %d", dwa.OriginStateID)
	}
}
