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

package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/lib/go/merrors"
)

func TestMultiError(t *testing.T) {
	// Correct Multi error implementation tests
	if returnNil() != nil {
		t.Error("'if err == nil' check should succeed for returned error")
	}
	if returnNilMulti().AsError() != nil {
		t.Error("'if err == nil' check should succeed for returned MultiError")
	}
	me := merrors.NewMulti()
	assert.Nil(t, me)
	var e error = nil
	me = me.Add(e)
	assert.Nil(t, me)
	if me != nil {
		t.Error("'if err == nil' expected to succeed for multi error")
	}
	if me.AsError() != nil {
		t.Error("'if err == nil' check should succeed for converted error")
	}
	me = me.Add(merrors.ErrNotFound)
	me = me.Add(merrors.ErrNotFound)
	assert.NotNil(t, me)
	assert.NotNil(t, me.AsError())
	assert.GreaterOrEqual(t, len(me.AsError().Error()), 20)
	assert.Equal(t, me.AsError().Error(), me.Error())
	assert.Equal(t, 2, len(me.Get()))

	multi := returnMulti(merrors.ErrNotFound, merrors.ErrAlreadyExists, nil)
	assert.NotNil(t, multi)
	assert.Equal(t, 2, len(multi.Get()))
	assert.NotEmpty(t, multi.Error())

	multiErr := returnMultiError(merrors.ErrNotFound, nil, merrors.ErrAlreadyExists)
	assert.NotNil(t, multiErr)
	assert.NotEmpty(t, multiErr.Error())

	// test AddFmt
	me = merrors.NewMulti().AddFmt(fmt.Errorf("foo bar"), "Multi error has %s (%d)", "one", 1)
	assert.Len(t, me.Get(), 1)
	assert.Equal(t, "Multi error has one (1) foo bar", me.Get()[0].Error())
	assert.Equal(t, "Multi error has one (1) foo bar", me.Error())
	me = me.AddFmt(fmt.Errorf("foo bars"), "Multi error has %s (%d)", "two", 2)
	assert.Len(t, me.Get(), 2)
	assert.Equal(t, "Multi error has two (2) foo bars", me.Get()[1].Error())
	assert.Equal(t, "errors: [0: Multi error has one (1) foo bar; 1: Multi error has two (2) foo bars]", me.Error())

	var nilMulti *merrors.Multi
	nilMulti = nilMulti.AddFmt(fmt.Errorf("just an error"), "")
	assert.NotNil(t, nilMulti)
	assert.Error(t, nilMulti.AsError())
}

func returnNil() error {
	return merrors.NewMulti().AsError()
}

func returnNilMulti() *merrors.Multi {
	return merrors.NewMulti()
}

func returnMulti(e1, e2, e3 error) *merrors.Multi {
	return merrors.NewMulti(e1, e2, e3)
}

func returnMultiError(e1, e2, e3 error) error {
	return merrors.NewMulti(e1, e2, e3)
}
