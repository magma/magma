/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package protos_test

import (
	"testing"

	"magma/orc8r/cloud/go/protos"
)

func TestOverlapFillIn(t *testing.T) {
	type Embeded1 struct {
		a, B, C int
	}
	type Embeded2 struct {
		C, d, E int
	}

	type S_A struct {
		S      string
		I      int
		E      Embeded1
		hidden bool
		S1     string
		M      map[string]Embeded1
	}

	type S_B struct {
		I      int
		S1     string
		hidden bool
		E      Embeded2
		M      map[string]*Embeded1
	}

	a := S_A{"str1", 1, Embeded1{1, 2, 3}, true, "str2",
		map[string]Embeded1{}}
	b := S_B{11, "str3", false, Embeded2{111, 222, 333},
		map[string]*Embeded1{"key": &Embeded1{1, 2, 3}}}

	count := protos.FillIn(&b, &a)
	if count <= 0 || a.I != b.I || a.S1 != b.S1 || a.S != "str1" || b.I != 11 ||
		a.hidden == b.hidden || a.E.C != b.E.C || a.E.a != 1 || a.E.B != 2 ||
		len(a.M) != 1 || a.M["key"].C != b.M["key"].C {

		t.Fatalf("Invalid assignment:\n\tcount: %d\n\ta: %+#v\n\tb: %+#v\n",
			count, a, b)
	}
}
