// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"bytes"
	"math"
	"testing"
)

func TestInteger64(t *testing.T) {
	n := Integer64(math.MaxInt64)
	b := []byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	if x := n.Serialize(); !bytes.Equal(b, x) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, x)
	}
	if n.Len() != 8 {
		t.Fatalf("Unexpected len. Want 8, have %d", n.Len())
	}
	if n.Padding() != 0 {
		t.Fatalf("Unexpected padding. Want 0, have %d", n.Padding())
	}
	if n.Type() != Integer64Type {
		t.Fatalf("Unexpected type. Want %d, have %d",
			Integer64Type, n.Type())
	}
	if len(n.String()) == 0 {
		t.Fatalf("Unexpected empty string")
	}
}

func TestNegativeInteger64(t *testing.T) {
	n := Integer64(math.MinInt64)
	b := []byte{0x80, 0, 0, 0, 0, 0, 0, 0}
	if x := n.Serialize(); !bytes.Equal(b, x) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, x)
	}
}

func TestDecodeInteger64(t *testing.T) {
	b := []byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	n, err := DecodeInteger64(b)
	if err != nil {
		t.Fatal(err)
	}
	z := int64(math.MaxInt64)
	if int64(n.(Integer64)) != z {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", z, n)
	}
}

func TestDecodeNegativeInteger64(t *testing.T) {
	b := []byte{0x80, 0, 0, 0, 0, 0, 0, 0}
	n, err := DecodeInteger64(b)
	if err != nil {
		t.Fatal(err)
	}
	z := int64(math.MinInt64)
	if int64(n.(Integer64)) != z {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", z, n)
	}
}

func BenchmarkInteger64(b *testing.B) {
	v := Integer64(math.MaxInt64)
	for n := 0; n < b.N; n++ {
		v.Serialize()
	}
}

func BenchmarkDecodeInteger64(b *testing.B) {
	v := []byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	for n := 0; n < b.N; n++ {
		DecodeInteger64(v)
	}
}
