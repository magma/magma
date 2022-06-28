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

func TestMinQueue_PushLess(t *testing.T) {
	m := &containers.MinQueue{}
	m.Push(2)
	m.Push(1)
	assert.Equal(t, 1, m.Top())
}

func TestMinQueue_PushGreater(t *testing.T) {
	m := &containers.MinQueue{}
	m.Push(1)
	m.Push(2)
	assert.Equal(t, 1, m.Top())
}

func TestMinQueue_PopLess(t *testing.T) {
	m := &containers.MinQueue{}
	m.Push(1)
	m.Push(2)
	m.Pop()
	assert.Equal(t, 2, m.Top())
}

func TestMinQueue_PopGreater(t *testing.T) {
	m := &containers.MinQueue{}
	m.Push(2)
	m.Push(1)
	m.Pop()
	assert.Equal(t, 1, m.Top())
}
