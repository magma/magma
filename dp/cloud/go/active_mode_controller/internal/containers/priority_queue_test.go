/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package containers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/active_mode_controller/internal/containers"
)

func TestPriorityQueue_Insert(t *testing.T) {
	p := &containers.PriorityQueue{}
	p.Push(2)
	p.Push(3)
	p.Push(1)

	assert.Equal(t, 3, p.Top())
}

func TestPriorityQueue_DeleteMax(t *testing.T) {
	p := &containers.PriorityQueue{}
	key := p.Push(2)
	p.Push(1)
	p.Delete(key)

	assert.Equal(t, 1, p.Top())
}

func TestPriorityQueue_Delete(t *testing.T) {
	p := containers.PriorityQueue{}
	key1 := p.Push(3)
	key2 := p.Push(2)
	p.Push(1)
	p.Delete(key2)
	p.Delete(key1)

	assert.Equal(t, 1, p.Top())
}
