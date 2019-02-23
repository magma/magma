// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"bytes"
	"testing"
)

// TODO(fiorix): Use some real AVP bytes here.

func TestGrouped(t *testing.T) {
	s := Grouped("hello")
	b := []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f}
	if v := s.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	if s.Len() != 5 {
		t.Fatalf("Unexpected len. Want 5, have %d", s.Len())
	}
	if s.Padding() != 0 {
		t.Fatalf("Unexpected padding. Want 0, have %d", s.Padding())
	}
	if s.Type() != GroupedType {
		t.Fatalf("Unexpected type. Want %d, have %d",
			GroupedType, s.Type())
	}
	if len(s.String()) == 0 {
		t.Fatalf("Unexpected empty string")
	}
}

func TestDecodeGrouped(t *testing.T) {
	b := []byte{
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c,
		0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64,
	}
	s, err := DecodeGrouped(b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal([]byte(s.(Grouped)), b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, s)
	}
	if s.Len() != 12 {
		t.Fatalf("Unexpected len. Want 12, have %d", s.Len())
	}
	if s.Padding() != 0 {
		t.Fatalf("Unexpected padding. Want 0, have %d", s.Padding())
	}
	if v := string(s.(Grouped)); v != "hello, world" {
		t.Fatalf("Unexpected string. Want 'hello, world', have %q", v)
	}
}
