// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"bytes"
	"net"
	"testing"
)

func TestIPv4(t *testing.T) {
	ip4 := IPv4(net.ParseIP("10.0.0.1"))
	b := []byte{0x0a, 0x00, 0x00, 0x01}
	if v := ip4.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	if ip4.Len() != 4 {
		t.Fatalf("Unexpected leip4. Want 4, have %d", ip4.Len())
	}
	if ip4.Padding() != 0 {
		t.Fatalf("Unexpected padding. Want 0, have %d", ip4.Padding())
	}
	if ip4.Type() != IPv4Type {
		t.Fatalf("Unexpected type. Want %d, have %d",
			IPv4Type, ip4.Type())
	}
	if len(ip4.String()) == 0 {
		t.Fatalf("Unexpected empty string")
	}
}

func TestIPv4Malformed(t *testing.T) {
	ip4 := IPv4(net.ParseIP("2001:0db8::ff00:0042:8329"))
	b := []byte{0x00, 0x02,
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0x00, 0x00, 0x42, 0x83, 0x29,
	}
	if v := ip4.Serialize(); bytes.Equal(v, b) {
		t.Fatalf("IPv6 match, that's unexpected")
	}
}

func TestDecodeIPv4(t *testing.T) {
	b := []byte{0x0a, 0x00, 0x00, 0x01}
	ip4, err := DecodeIPv4(b)
	if err != nil {
		t.Fatal(err)
	}
	if ip := net.IP(ip4.(IPv4)).String(); ip != "10.0.0.1" {
		t.Fatalf("Unexpected value. Want 10.0.0.1, have %s", ip)
	}
}
