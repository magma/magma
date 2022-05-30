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
	"sync"

	"github.com/magma/magma/src/go/protos/magma/capture"
)

// Buffer is a concurrency safe buffer used by the capture service.
type Buffer struct {
	unaryCalls []*capture.UnaryCall

	sync.RWMutex
}

// NewBuffer returns a new empty buffer.
func NewBuffer() *Buffer {
	return &Buffer{}
}

// Write locks the buffer while writing a UnaryCall to the buffer.
func (b *Buffer) Write(call *capture.UnaryCall) {
	b.Lock()
	defer b.Unlock()
	b.unaryCalls = append(b.unaryCalls, call)
}

// Flush locks the buffer, returns all of its contents and replaces the internal call list with nil.
func (b *Buffer) Flush() []*capture.UnaryCall {
	b.RLock()
	defer b.RUnlock()
	calls := b.unaryCalls
	b.unaryCalls = nil
	return calls
}
