// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"bytes"
	"net"
	"testing"
)

func TestAddressIPv4(t *testing.T) {
	address := Address(net.ParseIP("10.0.0.1"))
	b := []byte{0x00, 0x01, 0x0a, 0x00, 0x00, 0x01}
	if v := address.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	if address.Padding() != 2 {
		t.Fatalf("Unexpected padding. Want 2, have %d",
			address.Padding())
	}
	if address.Type() != AddressType {
		t.Fatalf("Unexpected type. Want %d, have %d",
			AddressType, address.Type())
	}
	if address.Len() != 6 {
		t.Fatalf("Unexpected len. Want 6, have %d", address.Len())
	}
	if len(address.String()) == 0 {
		t.Fatalf("Unexpected empty string")
	}
	//t.Log(address)
}

func TestDecodeAddressEmpty(t *testing.T) {
	_, err := DecodeAddress([]byte{})
	if err == nil {
		t.Fatal("Empty Address was decoded with no error.")
	}
}

func TestDecodeAddressInvalid(t *testing.T) {
	_, err := DecodeAddress([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err == nil {
		t.Fatal("Invalid Address was decoded with no error.")
	}
}

func TestDecodeAddressBadIPv4(t *testing.T) {
	_, err := DecodeAddress([]byte{
		0x00, 0x01, 0x0a, 0x00, 0x00, 0x00, 0x00})
	if err == nil {
		t.Fatal("Bad IPv4 was decoded with no error.")
	}
}

func TestDecodeAddressIPv4(t *testing.T) {
	b := []byte{0x00, 0x01, 0x0a, 0x00, 0x00, 0x01}
	address, err := DecodeAddress(b)
	if err != nil {
		t.Fatal(err)
	}
	if ip := net.IP(address.(Address)).String(); ip != "10.0.0.1" {
		t.Fatalf("Unexpected value. Want 10.0.0.1, have %s", ip)
	}
	if address.Padding() != 2 {
		t.Fatalf("Unexpected padding. Want 2, have %d",
			address.Padding())
	}
	//t.Log(address)
}

func TestAddressIPv6(t *testing.T) {
	address := Address(net.ParseIP("2001:0db8::ff00:0042:8329"))
	b := []byte{0x00, 0x02,
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0x00, 0x00, 0x42, 0x83, 0x29,
	}
	if v := address.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	if address.Padding() != 2 {
		t.Fatalf("Unexpected padding. Want 2, have %d",
			address.Padding())
	}
	if address.Type() != AddressType {
		t.Fatalf("Unexpected type. Want %d, have %d",
			AddressType, address.Type())
	}
	if address.Len() != 18 {
		t.Fatalf("Unexpected len. Want 18, have %d", address.Len())
	}
	//t.Log(address)
}

func TestDecodeAddressBadIPv6(t *testing.T) {
	_, err := DecodeAddress([]byte{
		0x00, 0x02, 0x0a, 0x00, 0x00, 0x00, 0x00})
	if err == nil {
		t.Fatal("Unexpected bad IPv6 was decoded with no error.")
	}
}

func TestDecodeAddressIPv6(t *testing.T) {
	b := []byte{0x00, 0x02,
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0x00, 0x00, 0x42, 0x83, 0x29,
	}
	address, err := DecodeAddress(b)
	if err != nil {
		t.Fatal(err)
	}
	want := "2001:db8::ff00:42:8329"
	if ip := net.IP(address.(Address)).String(); ip != want {
		t.Fatalf("Unexpected value. Want %s, have %s", want, ip)
	}
	if address.Padding() != 2 {
		t.Fatalf("Unexpected padding. Want 2, have %d", address.Padding())
	}
}

func TestAddressIPv4_Generic(t *testing.T) {
	address := Address([]byte{0x0a, 0x00, 0x00, 0x01})
	b := []byte{0x00, 0x01, 0x0a, 0x00, 0x00, 0x01}
	if v := address.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	if address.Padding() != 2 {
		t.Fatalf("Unexpected padding. Want 2, have %d",
			address.Padding())
	}
	if address.Type() != AddressType {
		t.Fatalf("Unexpected type. Want %d, have %d",
			AddressType, address.Type())
	}
	if address.Len() != 6 {
		t.Fatalf("Unexpected len. Want 6, have %d", address.Len())
	}
	if len(address.String()) == 0 {
		t.Fatalf("Unexpected empty string")
	}
}

func TestAddressIPv6_Generic(t *testing.T) {
	address := Address([]byte{0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0x00, 0x00, 0x42, 0x83, 0x29,
	})
	b := []byte{0x00, 0x02,
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0x00, 0x00, 0x42, 0x83, 0x29,
	}
	if v := address.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	if address.Padding() != 2 {
		t.Fatalf("Unexpected padding. Want 2, have %d",
			address.Padding())
	}
	if address.Type() != AddressType {
		t.Fatalf("Unexpected type. Want %d, have %d",
			AddressType, address.Type())
	}
	if address.Len() != 18 {
		t.Fatalf("Unexpected len. Want 18, have %d", address.Len())
	}
}

func TestAddressE164_Generic(t *testing.T) {
	var addressBytes []byte
	addressType := []byte{0x0, 0x8}
	addressValue := []byte("48602007060")
	addressBytes = make([]byte, len(addressType)+len(addressValue))
	copy(addressBytes[:2], addressType)
	copy(addressBytes[2:], addressValue)
	address := Address(addressBytes)

	b := []byte{0x00, 0x08, 0x34, 0x38, 0x36, 0x30, 0x32, 0x30, 0x30, 0x37, 0x30, 0x36, 0x30}
	if v := address.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	if address.Padding() != 3 {
		t.Fatalf("Unexpected padding. Want 3, have %d",
			address.Padding())
	}
	if address.Type() != AddressType {
		t.Fatalf("Unexpected type. Want %d, have %d",
			AddressType, address.Type())
	}
	if address.Len() != 13 {
		t.Fatalf("Unexpected len. Want 13, have %d", address.Len())
	}
	if len(address.String()) == 0 {
		t.Fatalf("Unexpected empty string")
	}
}

func TestDecodeAddressE164(t *testing.T) {
	b := []byte{0x00, 0x08, 0x34, 0x38, 0x36, 0x30, 0x32, 0x30, 0x30, 0x37, 0x30, 0x36, 0x30}
	address, err := DecodeAddress(b)
	if err != nil {
		t.Fatal(err)
	}
	if address.Padding() != 3 {
		t.Fatalf("Unexpected padding. Want 3, have %d",
			address.Padding())
	}
}

func BenchmarkAddressIPv4(b *testing.B) {
	address := Address(net.ParseIP("10.0.0.1"))
	for n := 0; n < b.N; n++ {
		address.Serialize()
	}
}

func BenchmarkDecodeAddressIPv4(b *testing.B) {
	v := []byte{0x00, 0x01, 0x0a, 0x00, 0x00, 0x01}
	for n := 0; n < b.N; n++ {
		DecodeAddress(v)
	}
}

func BenchmarkAddressIPv6(b *testing.B) {
	address := Address(net.ParseIP("2001:db8::ff00:42:8329"))
	for n := 0; n < b.N; n++ {
		address.Serialize()
	}
}

func BenchmarkDecodeAddressIPv6(b *testing.B) {
	v := []byte{0x00, 0x02,
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0x00, 0x00, 0x42, 0x83, 0x29,
	}
	for n := 0; n < b.N; n++ {
		DecodeAddress(v)
	}
}
