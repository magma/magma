// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"bytes"
	"testing"
)

func TestOctetString(t *testing.T) {
	s := OctetString("hello")
	b := []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f}
	if v := s.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	if s.Len() != 5 {
		t.Fatalf("Unexpected len. Want 5, have %d", s.Len())
	}
	if s.Padding() != 3 {
		t.Fatalf("Unexpected padding. Want 3, have %d", s.Padding())
	}
	if s.Type() != OctetStringType {
		t.Fatalf("Unexpected type. Want %d, have %d",
			OctetStringType, s.Type())
	}
	if len(s.String()) == 0 {
		t.Fatalf("Unexpected empty string")
	}
}

func TestDecodeOctetString(t *testing.T) {
	b := []byte{
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c,
		0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64,
	}
	s, err := DecodeOctetString(b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal([]byte(s.(OctetString)), b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, s)
	}
	if s.Len() != 12 {
		t.Fatalf("Unexpected len. Want 12, have %d", s.Len())
	}
	if s.Padding() != 0 {
		t.Fatalf("Unexpected padding. Want 0, have %d", s.Padding())
	}
	if v := string(s.(OctetString)); v != "hello, world" {
		t.Fatalf("Unexpected string. Want 'hello, world', have %q", v)
	}
}

func BenchmarkOctetString(b *testing.B) {
	v := OctetString("hello, world")
	for n := 0; n < b.N; n++ {
		v.Serialize()
	}
}

func BenchmarkDecodeOctetString(b *testing.B) {
	v := []byte{
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c,
		0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64,
	}
	for n := 0; n < b.N; n++ {
		DecodeOctetString(v)
	}
}
