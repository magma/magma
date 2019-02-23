// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smparser

import (
	"testing"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
)

func TestCEA_MissingResultCode(t *testing.T) {
	m := diam.NewMessage(diam.CapabilitiesExchange, 0, 0, 0, 0, nil)
	cea := new(CEA)
	err := cea.Parse(m, Client)
	if err == nil {
		t.Fatal("Broken CEA was parsed with no errors")
	}
	if err != nil && err != ErrMissingResultCode {
		t.Fatal("Unexpected error:", err)
	}
}

func TestCEA_MissingOriginHost(t *testing.T) {
	m := diam.NewMessage(diam.CapabilitiesExchange, 0, 0, 0, 0, nil)
	m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(diam.Success))
	cea := new(CEA)
	err := cea.Parse(m, Client)
	if err == nil {
		t.Fatal("Broken CEA was parsed with no errors")
	}
	if err != nil && err != ErrMissingOriginHost {
		t.Fatal("Unexpected error:", err)
	}
}

func TestCEA_MissingOriginRealm(t *testing.T) {
	m := diam.NewMessage(diam.CapabilitiesExchange, 0, 0, 0, 0, nil)
	m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(diam.Success))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("foobar"))
	cea := new(CEA)
	err := cea.Parse(m, Client)
	if err == nil {
		t.Fatal("Broken CEA was parsed with no errors")
	}
	if err != nil && err != ErrMissingOriginRealm {
		t.Fatal("Unexpected error:", err)
	}
}

func TestCEA_MissingApplication(t *testing.T) {
	m := diam.NewMessage(diam.CapabilitiesExchange, 0, 0, 0, 0, dict.Default)
	m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(diam.Success))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("foobar"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test"))
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	cea := new(CEA)
	err := cea.Parse(m, Client)
	if err == nil {
		t.Fatal("Broken CEA was parsed with no errors")
	}
	if err != ErrMissingApplication {
		t.Fatal("Unexpected error:", err)
	}
}

func TestCEA_MissingApplicationWithError(t *testing.T) {
	m := diam.NewMessage(diam.CapabilitiesExchange, 0, 0, 0, 0, dict.Default)
	m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(diam.ResourcesExceeded))
	m.NewAVP(avp.FailedAVP, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1000)),
		},
	})

	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("foobar"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test"))
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	cea := new(CEA)
	err := cea.Parse(m, Client)
	if err == nil {
		t.Fatal("Broken CEA was parsed with no errors")
	}
	e, ok := err.(*ErrFailedResultCode)
	if !ok {
		t.Fatalf("Unexpected error type. Want *ErrFailedResultCode, have %T", err)
	}
	if e.ResultCode != diam.ResourcesExceeded {
		t.Fatalf("Unexpected ResultCode. Want %d, have %d", diam.ResourcesExceeded, e.ResultCode)
	}
	g, ok := e.FailedAVP[0].Data.(*diam.GroupedAVP)
	if !ok {
		t.Fatalf("Unexpected type. Want *diam.GroupedAVP, have %T", e.FailedAVP[0].Data)
	}
	d, ok := g.AVP[0].Data.(datatype.Unsigned32)
	if !ok {
		t.Fatalf("Unexpected type. Want *datatype.Unsigned32, have %T", e.FailedAVP[0].Data)
	}
	if d != 1000 {
		t.Fatalf("Wrong value for FailedAVP. Want 1000, have %d", d)
	}
}

func TestCEA_NoCommonApplication(t *testing.T) {
	m := diam.NewMessage(diam.CapabilitiesExchange, 0, 0, 0, 0, dict.Default)
	m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(diam.Success))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("foobar"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test"))
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(2))
	cea := new(CEA)
	err := cea.Parse(m, Server)
	if err == nil {
		t.Fatal("Broken CEA was parsed with no errors")
	}
	if err != ErrNoCommonApplication {
		t.Fatal("Unexpected error:", err)
	}
}

func TestCEA_FailedAcctAppID(t *testing.T) {
	m := diam.NewMessage(diam.CapabilitiesExchange, 0, 0, 0, 0, nil)
	m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(diam.Success))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("foobar"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test"))
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1000))
	cea := new(CEA)
	err := cea.Parse(m, Server)
	if err == nil {
		t.Fatal("Broken CEA was parsed with no errors")
	}
	if err != ErrNoCommonApplication {
		t.Fatal("Unexpected error:", err)
	}
}

func TestCEA(t *testing.T) {
	m := diam.NewMessage(diam.CapabilitiesExchange, 0, 0, 0, 0, nil)
	m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(diam.Success))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("foobar"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test"))
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4))
	cea := new(CEA)
	if err := cea.Parse(m, Server); err != nil {
		t.Fatal(err)
	}
	if cea.ResultCode != diam.Success {
		t.Fatalf("Unexpected Result-Code. Want %d, have %d",
			diam.Success, cea.ResultCode)
	}
	if cea.OriginStateID != 1 {
		t.Fatalf("Unexpected Origin-State-ID. Want 1, have %d", cea.OriginStateID)
	}
	if app := cea.Applications(); len(app) != 1 {
		if app[0] != 4 {
			t.Fatalf("Unexpected app ID. Want 4, have %d", app[0])
		}
	}
}
