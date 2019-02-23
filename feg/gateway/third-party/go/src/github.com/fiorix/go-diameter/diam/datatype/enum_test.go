// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"bytes"
	"testing"
)

func TestEnumerated(t *testing.T) {
	s := Enumerated(7)
	b := []byte{0x00, 0x00, 0x00, 0x07}
	if v := s.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	if s.Len() != 4 {
		t.Fatalf("Unexpected len. Want 4, have %d", s.Len())
	}
	if s.Padding() != 0 {
		t.Fatalf("Unexpected padding. Want 0, have %d", s.Padding())
	}
	if s.Type() != EnumeratedType {
		t.Fatalf("Unexpected type. Want %d, have %d",
			EnumeratedType, s.Type())
	}
	if len(s.String()) == 0 {
		t.Fatalf("Unexpected empty string")
	}
}

func TestDecodeEnumerated(t *testing.T) {
	b := []byte{0x00, 0x00, 0x00, 0xFF}
	s, err := DecodeEnumerated(b)
	if err != nil {
		t.Fatal(err)
	}
	if v := s.(Enumerated); v != 255 {
		t.Fatalf("Unexpected value. Want 255, have %d", v)
	}
	if s.Len() != 4 {
		t.Fatalf("Unexpected len. Want 4, have %d", s.Len())
	}
	if s.Padding() != 0 {
		t.Fatalf("Unexpected padding. Want 0, have %d", s.Padding())
	}
}
