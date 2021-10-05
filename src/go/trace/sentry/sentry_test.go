// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sentry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrace_NewTrace(t *testing.T) {
	operation := "testSpan"
	tr := NewTrace(context.Background(), operation)
	trace, ok := tr.(*Trace)
	assert.True(
		t, ok, "NewCrash must return *sentry.Trace, have=%+v", tr)
	assert.NotNil(t, trace.Span)
	assert.Equal(t, operation, trace.Op)
}

func TestTrace_SetTag(t *testing.T) {
	operation := "testSpan"
	testName := "testName"
	testVal := "testVal"
	tr := NewTrace(context.Background(), operation)
	tr.SetTag(testName, testVal)

	trace, ok := tr.(*Trace)
	assert.True(
		t, ok, "NewCrash must return *sentry.Trace, have=%+v", tr)
	assert.NotNil(t, trace.Span)
	assert.Equal(t, testVal, trace.Tags[testName])
}

func TestTrace_Finish(t *testing.T) {
	operation := "testSpan"
	tr := NewTrace(context.Background(), operation)
	tr.Finish()
}
