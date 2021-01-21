/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package protos_test

import (
	"testing"

	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
)

func TestOverlapFillIn(t *testing.T) {
	type Embeded1 struct {
		a, B, C int
	}
	type Embeded2 struct {
		C, d, E int
	}

	type S_A struct {
		S          string
		I          int
		E          Embeded1
		hidden     bool
		S1         string
		MixEdCaSe  string
		mixEdCaSe1 string
		MixEdCaSe2 string
		M          map[string]Embeded1
		M1         map[string]*Embeded1
		M2         map[string]Embeded1
	}

	type S_B struct {
		I          int
		S1         string
		hidden     bool
		E          Embeded2
		MiXedCase  string
		MiXedCase1 string
		miXedCase2 string
		M          map[string]*Embeded2
		M1         map[string]*Embeded1
		M2         map[string]Embeded1
	}

	type TestInitStr1 struct {
		I1      int
		F1      float32
		S11     string
		hidden1 bool
		E1      Embeded1
		EPtr1   *Embeded2
		ISlice1 []int
		IArr1   [5]int
	}

	type TestInitStr2 struct {
		I2      int
		S2      string
		hidden2 bool
		E2      Embeded1
		EPtr2   *Embeded2
		TSPtr2  *TestInitStr1
		ISlice2 []int
		IArr2   [5]int
		M2      map[string]int
	}

	a := S_A{"str1", 1, Embeded1{1, 2, 3}, true, "str2",
		"bla Bla bla", "bla1", "bla2",
		map[string]Embeded1{}, map[string]*Embeded1{}, nil}
	b := S_B{11, "str3", false, Embeded2{111, 222, 333},
		"Foo bar", "foo1", "foo2",
		map[string]*Embeded2{"key": {1, 2, 3}},
		map[string]*Embeded1{"key": {1, 2, 3}},
		map[string]Embeded1{"key": {1, 2, 3}}}

	count := protos.FillIn(&b, &a)
	if count <= 0 || a.I != b.I || a.S1 != b.S1 || a.S != "str1" || b.I != 11 ||
		a.hidden == b.hidden || a.E.C != b.E.C || a.E.a != 1 || a.E.B != 2 ||
		a.MixEdCaSe != b.MiXedCase || a.mixEdCaSe1 == b.MiXedCase1 || a.MixEdCaSe2 == b.miXedCase2 ||
		len(a.M) != 1 || a.M["key"].C != 1 ||
		len(a.M1) != 1 || a.M1["key"].a != 1 || a.M1["key"].B != 2 || a.M1["key"].C != 3 ||
		len(a.M2) != 1 || a.M2["key"].a != 1 || a.M2["key"].B != 2 || a.M2["key"].C != 3 {

		t.Fatalf("Invalid assignment:\n\tcount: %d\n\ta: %+#v\n\tb: %+#v\n",
			count, a, b)
	}

	tp := &TestInitStr2{}
	assert.Nil(t, tp.EPtr2)
	if tp.M2 != nil {
		t.Fatal("M2 != nil")
	}
	if tp.ISlice2 != nil {
		t.Fatal("ISlice2 != nil")
	}
	assert.Nil(t, tp.TSPtr2)
	assert.Nil(t, tp.EPtr2)

	tp = protos.SafeInit(tp).(*TestInitStr2)
	assert.NotNil(t, tp.EPtr2)
	if tp.M2 == nil {
		t.Fatal("M2 == nil")
	}
	if tp.ISlice2 == nil {
		t.Fatal("ISlice2 == nil")
	}
	assert.NotNil(t, tp.TSPtr2)
	assert.NotNil(t, tp.EPtr2)
	assert.NotNil(t, tp.TSPtr2.EPtr1)
	if tp.TSPtr2.ISlice1 == nil {
		t.Fatal("TSPtr2.ISlice1 == nil")
	}

	tp = nil
	tp = protos.SafeInit(tp).(*TestInitStr2)
	assert.NotNil(t, tp.EPtr2)
	if tp.M2 == nil {
		t.Fatal("M2 == nil")
	}
	if tp.ISlice2 == nil {
		t.Fatal("ISlice2 == nil")
	}
	assert.NotNil(t, tp.TSPtr2)
	assert.NotNil(t, tp.EPtr2)
	assert.NotNil(t, tp.TSPtr2.EPtr1)
	if tp.TSPtr2.ISlice1 == nil {
		t.Fatal("TSPtr2.ISlice1 == nil")
	}
}
