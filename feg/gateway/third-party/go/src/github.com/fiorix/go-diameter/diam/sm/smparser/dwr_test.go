// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smparser

import (
	"testing"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

func TestDWR_MissingOriginHost(t *testing.T) {
	m := diam.NewRequest(diam.DeviceWatchdog, 0, dict.Default)
	dwr := new(DWR)
	err := dwr.Parse(m)
	if err != nil && err != ErrMissingOriginHost {
		t.Fatal("Unexpected error:", err)
	}
}

func TestDWR_MissingOriginRealm(t *testing.T) {
	m := diam.NewRequest(diam.CapabilitiesExchange, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("foobar"))
	dwr := new(DWR)
	err := dwr.Parse(m)
	if err != nil && err != ErrMissingOriginRealm {
		t.Fatal("Unexpected error:", err)
	}
}

func TestDWR_OK(t *testing.T) {
	m := diam.NewRequest(diam.CapabilitiesExchange, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("foobar"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test"))
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	dwr := new(DWR)
	err := dwr.Parse(m)
	if err != nil {
		t.Fatal(err)
	}
}
