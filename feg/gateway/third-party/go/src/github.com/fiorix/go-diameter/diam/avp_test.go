// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

var testAVP = [][]byte{ // Body of a CER message
	{ // Origin-Host
		0x00, 0x00, 0x01, 0x08,
		0x40, 0x00, 0x00, 0x0e,
		0x63, 0x6c, 0x69, 0x65,
		0x6e, 0x74, 0x00, 0x00,
	},
	{ // Origin-Realm
		0x00, 0x00, 0x01, 0x28,
		0x40, 0x00, 0x00, 0x11,
		0x6c, 0x6f, 0x63, 0x61,
		0x6c, 0x68, 0x6f, 0x73,
		0x74, 0x00, 0x00, 0x00,
	},
	{ // Host-IP-Address
		0x00, 0x00, 0x01, 0x01,
		0x40, 0x00, 0x00, 0x0e,
		0x00, 0x01, 0xc0, 0xa8,
		0xf2, 0x7a, 0x00, 0x00,
	},
	{ // Vendor-Id
		0x00, 0x00, 0x01, 0x0a,
		0x40, 0x00, 0x00, 0x0c,
		0x00, 0x00, 0x00, 0x0d,
	},
	{ // Product-Name
		0x00, 0x00, 0x01, 0x0d,
		0x40, 0x00, 0x00, 0x13,
		0x67, 0x6f, 0x2d, 0x64,
		0x69, 0x61, 0x6d, 0x65,
		0x74, 0x65, 0x72, 0x00,
	},
	{ // Origin-State-Id
		0x00, 0x00, 0x01, 0x16,
		0x40, 0x00, 0x00, 0x0c,
		0xe8, 0x3e, 0x3b, 0x84,
	},
}

func TestNewAVP(t *testing.T) {
	a := NewAVP(
		avp.OriginHost,                      // Code
		avp.Mbit,                            // Flags
		0,                                   // Vendor
		datatype.DiameterIdentity("foobar"), // Data
	)
	if a.Length != 14 { // Length in the AVP header
		t.Fatalf("Unexpected length. Want 14, have %d", a.Length)
	}
	if a.Len() != 16 { // With padding
		t.Fatalf("Unexpected length (with padding). Want 16, have %d\n", a.Len())
	}
	t.Log(a)
}

func TestDecodeAVP(t *testing.T) {
	a, err := DecodeAVP(testAVP[0], 1, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	switch {
	case a.Code != avp.OriginHost:
		t.Fatalf("Unexpected Code. Want %d, have %d", avp.OriginHost, a.Code)
	case a.Flags != avp.Mbit:
		t.Fatalf("Unexpected Code. Want %#x, have %#x", avp.Mbit, a.Flags)
	case a.Length != 14:
		t.Fatalf("Unexpected Length. Want 14, have %d", a.Length)
	case a.Data.Padding() != 2:
		t.Fatalf("Unexpected Padding. Want 2, have %d", a.Data.Padding())
	}
	t.Log(a)
}

func TestDecodeAVPMalformed(t *testing.T) {
	_, err := DecodeAVP(testAVP[0][:1], 1, dict.Default)
	if err == nil {
		t.Fatal("Malformed AVP decoded with no error")
	}
}

func TestDecodeAVPWithVendorID(t *testing.T) {
	var userNameVendorXML = `<?xml version="1.0" encoding="UTF-8"?>
<diameter>
  <application id="1">
    <avp name="Session-Start-Indicator" code="1" vendor-id="999">
      <data type="UTF8String" />
    </avp>
  </application>
</diameter>`
	dict.Default.Load(bytes.NewReader([]byte(userNameVendorXML)))
	a := NewAVP(avp.UserName, avp.Mbit|avp.Vbit, 999, datatype.UTF8String("foobar"))
	b, err := a.Serialize()
	if err != nil {
		t.Fatal("Failed to serialize AVP:", err)
	}
	a, err = DecodeAVP(b, 1, dict.Default)
	if err != nil {
		t.Fatal("Failed to decode AVP:", err)
	}
	if a.VendorID != 999 {
		t.Fatalf("Unexpected VendorID. Want 999, have %d", a.VendorID)
	}
}

func TestEncodeAVP(t *testing.T) {
	a := &AVP{
		Code:  avp.OriginHost,
		Flags: avp.Mbit,
		Data:  datatype.DiameterIdentity("client"),
	}
	b, err := a.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, testAVP[0]) {
		t.Fatalf("AVPs do not match.\nWant:\n%s\nHave:\n%s",
			hex.Dump(testAVP[0]), hex.Dump(b))
	}
	t.Log(hex.Dump(b))
}

func TestEncodeAVPWithoutData(t *testing.T) {
	a := &AVP{
		Code:  avp.OriginHost,
		Flags: avp.Mbit,
	}
	_, err := a.Serialize()
	if err != nil {
		t.Log("Expected:", err)
	} else {
		t.Fatal("Unexpected serialization succeeded")
	}
}

func BenchmarkDecodeAVP(b *testing.B) {
	for n := 0; n < b.N; n++ {
		DecodeAVP(testAVP[0], 1, dict.Default)
	}
}

func BenchmarkEncodeAVP(b *testing.B) {
	a := NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("client"))
	for n := 0; n < b.N; n++ {
		a.Serialize()
	}
}
