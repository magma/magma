// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"encoding/hex"
	"testing"
)

var testHeader = []byte{ // CER
	0x01, 0x00, 0x00, 0x74,
	0x80, 0x00, 0x01, 0x01,
	0x00, 0x00, 0x00, 0x01,
	0x2c, 0x0b, 0x61, 0x49,
	0xdb, 0xbf, 0xd3, 0x85,
}

func TestDecodeHeader(t *testing.T) {
	hdr, err := DecodeHeader(testHeader)
	if err != nil {
		t.Fatal(err)
	}
	switch {
	case hdr.Version != 1:
		t.Fatalf("Unexpected Version. Want 1, have %d", hdr.Version)
	case hdr.MessageLength != 116:
		t.Fatalf("Unexpected MessageLength. Want 116, have %d", hdr.MessageLength)
	case hdr.CommandFlags != RequestFlag:
		t.Fatalf("Unexpected CommandFlags. Want %#x, have %#x", RequestFlag, hdr.CommandFlags)
	case hdr.CommandCode != 257:
		t.Fatalf("Unexpected CommandCode. Want 257, have %d", hdr.CommandCode)
	case hdr.ApplicationID != 1:
		t.Fatalf("Unexpected ApplicationId. Want 1, have %d", hdr.ApplicationID)
	case hdr.HopByHopID != 0x2c0b6149:
		t.Fatalf("Unexpected HopByHopId. Want 0x2c0b6149, have 0x%x", hdr.HopByHopID)
	case hdr.EndToEndID != 0xdbbfd385:
		t.Fatalf("Unexpected EndToEndId. Want 0xdbbf0385, have 0x%x", hdr.EndToEndID)
	}
	t.Log(hdr)
}

func TestDecodeHeaderMalformed(t *testing.T) {
	_, err := DecodeHeader(testHeader[:10])
	if err == nil {
		t.Fatal("Malformed header decoded with no errors")
	}
}

func TestEncodeHeader(t *testing.T) {
	hdr := &Header{
		Version:       1,
		MessageLength: 116,
		CommandFlags:  RequestFlag,
		CommandCode:   CapabilitiesExchange,
		ApplicationID: 1,
		HopByHopID:    0x2c0b6149,
		EndToEndID:    0xdbbfd385,
	}
	b := hdr.Serialize()
	if !bytes.Equal(testHeader, b) {
		t.Fatalf("Unexpected packet.\nWant:\n%s\nHave:\n%s",
			hex.Dump(testHeader), hex.Dump(b))
	}
}

func BenchmarkDecodeHeader(b *testing.B) {
	for n := 0; n < b.N; n++ {
		DecodeHeader(testHeader)
	}
}

func BenchmarkEncodeHeader(b *testing.B) {
	hdr := &Header{
		Version:       1,
		MessageLength: 116,
		CommandFlags:  RequestFlag,
		CommandCode:   CapabilitiesExchange,
		ApplicationID: 1,
		HopByHopID:    0x2c0b6149,
		EndToEndID:    0xdbbfd385,
	}
	for n := 0; n < b.N; n++ {
		hdr.Serialize()
	}
}
