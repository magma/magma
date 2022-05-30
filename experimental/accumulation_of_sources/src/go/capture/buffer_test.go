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

package capture

import (
	"testing"

	"github.com/magma/magma/src/go/protos/magma/capture"
	"github.com/stretchr/testify/assert"
)

func TestNewBuffer(t *testing.T) {
	buf := NewBuffer()
	assert.Equal(t, []*capture.UnaryCall(nil), buf.unaryCalls)
}

func TestBuffer_Write(t *testing.T) {
	buf := NewBuffer()
	call := &capture.UnaryCall{
		Method: "test",
	}
	call2 := &capture.UnaryCall{
		Method: "test2",
	}

	buf.Write(call)
	assert.Equal(t, call, buf.unaryCalls[0])
	buf.Write(call2)
	assert.Equal(t, call2, buf.unaryCalls[1])
}

func TestBuffer_Flush(t *testing.T) {
	buf := NewBuffer()
	call := &capture.UnaryCall{
		Method: "test",
	}
	call2 := &capture.UnaryCall{
		Method: "test2",
	}
	call3 := &capture.UnaryCall{
		Method: "test3",
	}

	buf.Write(call)
	buf.Write(call2)
	calls := buf.Flush()
	assert.Equal(t, call, calls[0])
	assert.Equal(t, call2, calls[1])

	buf.Write(call3)
	assert.Equal(t, call3, buf.unaryCalls[0])
}
