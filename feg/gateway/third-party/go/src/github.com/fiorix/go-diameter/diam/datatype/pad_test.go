// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import "testing"

func TestPad4(t *testing.T) {
	if n := pad4(2); n != 4 {
		t.Fatalf("Unexpected result. Want 4, have %d", n)
	}
}
