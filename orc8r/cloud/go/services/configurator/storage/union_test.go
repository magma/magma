/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnionFind(t *testing.T) {
	uf := newUnionFind([]string{"1", "2", "3", "4", "5"})

	// 1 | 2 | 3 | 4 | 5
	assert.Equal(
		t,
		[][]string{
			{"1"},
			{"2"},
			{"3"},
			{"4"},
			{"5"},
		},
		uf.getComponents(),
	)

	uf.union("1", "2")
	uf.union("3", "4")
	// 1 -> 2 | 3 -> 4 | 5
	assert.Equal(
		t,
		[][]string{
			{"5"},
			{"1", "2"},
			{"3", "4"},
		},
		uf.getComponents(),
	)

	uf.union("4", "2")
	// 1 -> (2, 3 -> 4) | 5
	assert.Equal(
		t,
		[][]string{
			{"5"},
			{"1", "2", "3", "4"},
		},
		uf.getComponents(),
	)

	// paths were compressed from last call to getComponents
	uf.union("4", "5")
	// 1 -> (2, 3, 4, 5)
	assert.Equal(
		t,
		[][]string{
			{"1", "2", "3", "4", "5"},
		},
		uf.getComponents(),
	)

	uf = newUnionFind([]string{"1", "2", "3", "4", "5"})
	uf.union("1", "2")
	uf.union("2", "3")
	uf.union("4", "5")
	// 1 -> (2, 3) | 4 -> 5
	uf.union("5", "3")
	// 1 -> (2, 3, 4 -> 5)
	// 4 -> (1 -> [2, 3], 5)
	assert.Equal(
		t,
		[][]string{
			{"1", "2", "3", "4", "5"},
		},
		uf.getComponents(),
	)

	// Trigger xRank < yRank
	uf = newUnionFind([]string{"1", "2", "3"})
	uf.union("1", "2")
	// 1 -> 2 | 3
	uf.union("3", "1")
	// 1 -> (2, 3)
	assert.Equal(
		t,
		[][]string{
			{"1", "2", "3"},
		},
		uf.getComponents(),
	)

	// Trigger xRoot == yRoot
	uf = newUnionFind([]string{"1", "2"})
	uf.union("1", "2")
	uf.union("1", "2")
	assert.Equal(
		t,
		[][]string{
			{"1", "2"},
		},
		uf.getComponents(),
	)
}
