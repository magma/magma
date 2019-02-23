// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"testing"
)

func TestUint24to32Malformed(t *testing.T) {
	if v := uint24to32([]byte{1}); v != 0 {
		t.Fatalf("Unexpected result. Want 0, have 0x%x", v)
	}
}

func TestUint24to32(t *testing.T) {
	if v := uint24to32([]byte{218, 190, 239}); v != 0xdabeef {
		t.Fatalf("Unexpected result. Want 0xdeadbeef, have 0x%x", v)
	}
}

func TestUint32to24(t *testing.T) {
	if v := uint32to24(0xdabeef); !bytes.Equal(v, []byte{218, 190, 239}) {
		t.Fatalf("Unexpected result. Want [218 190 239], have %v", v)
	}
}

func BenchmarkUint24to32(b *testing.B) {
	for n := 0; n < b.N; n++ {
		uint24to32([]byte{218, 190, 239})
	}
}

func BenchmarkUint32to24(b *testing.B) {
	for n := 0; n < b.N; n++ {
		uint32to24(0xdeadbeef)
	}
}
