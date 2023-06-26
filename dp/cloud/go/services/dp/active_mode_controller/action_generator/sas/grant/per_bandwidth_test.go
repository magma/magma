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

package grant_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator/sas/grant"
)

const msg = "\nexpected: %b\nactual:   %b"

func TestSelectGrantsWithRedundancy(t *testing.T) {
	data := []struct {
		name      string
		available uint32
		grants    uint32
		pref      []uint32
		minWidth  int
		maxWidth  int
		index     int
		expected  uint32
	}{{
		name:      "Should pick pair matching preference",
		available: 0,
		grants:    ^uint32(0),
		pref:      []uint32{1 << 10},
		minWidth:  0,
		maxWidth:  30,
		expected:  1<<10 | 1<<9,
	}, {
		name:      "Should pick second pair matching preference if first unavailable",
		available: 0,
		grants:    1<<10 | 1<<11 | (1<<5 - 1),
		pref:      []uint32{1 << 20, 1 << 10},
		minWidth:  0,
		maxWidth:  9,
		expected:  1<<10 | 1<<11,
	}, {
		name:      "Should pick nth pair if no matching preferences",
		available: 0,
		grants:    1<<5 - 1,
		pref:      nil,
		minWidth:  1,
		maxWidth:  3,
		expected:  1<<2 | 1<<0,
	}, {
		name:      "Should pick new pair according to preferences if current grants are do not form a valid pair",
		available: 1<<10 | 1<<5 | 1<<0,
		grants:    1<<30 | 1<<20,
		pref:      []uint32{1 << 10},
		minWidth:  4,
		maxWidth:  9,
		expected:  1<<30 | 1<<20 | 1<<10 | 1<<5,
	}, {
		name:      "Should pick nth new pair if no matching preferences",
		available: 1<<10 | 1<<5 | 1<<0,
		grants:    1 << 20,
		pref:      []uint32{1 << 20},
		minWidth:  0,
		maxWidth:  9,
		expected:  1<<20 | 1<<5 | 1<<0,
	}, {
		name:      "Should leave one grant according to preferences if can not find pair",
		available: 0,
		grants:    1<<20 | 1<<10 | 1<<0,
		pref:      []uint32{1 << 20},
		minWidth:  0,
		maxWidth:  5,
		expected:  1 << 20,
	}, {
		name:      "Should leave nth old grant if can not find pair and no matching preferences",
		available: 0,
		grants:    1<<20 | 1<<10 | 1<<0,
		pref:      nil,
		minWidth:  0,
		maxWidth:  5,
		expected:  1 << 0,
	}, {
		name:      "Should not pick pair that is too close to current grants",
		available: 1<<12 | 1<<8,
		grants:    1 << 10,
		pref:      nil,
		minWidth:  3,
		maxWidth:  30,
		expected:  1 << 10,
	}, {
		name:      "Should pick from existing and available grants",
		available: 1 << 10,
		grants:    1<<8 | 1<<1,
		pref:      nil,
		minWidth:  1,
		maxWidth:  5,
		expected:  1<<10 | 1<<8 | 1<<1,
	}, {
		name:      "Should pick nth grant",
		available: 1<<20 | 1<<15 | 1<<10 | 1<<5,
		grants:    0,
		pref:      nil,
		minWidth:  0,
		maxWidth:  30,
		index:     2,
		expected:  1<<15 | 1<<10,
	}, {
		name:      "Should not pick anything if can not find pair",
		available: 1 << 10,
		grants:    0,
		pref:      nil,
		minWidth:  0,
		maxWidth:  30,
		expected:  0,
	}}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			actual := grant.SelectGrantsWithRedundancy(tt.available, tt.grants, tt.pref, tt.minWidth, tt.maxWidth, tt.index)
			assert.Equal(t, tt.expected, actual, msg, tt.expected, actual)
		})
	}
}

func TestSelectGrantsWithoutRedundancy(t *testing.T) {
	data := []struct {
		name      string
		available uint32
		grants    uint32
		pref      []uint32
		index     int
		expected  uint32
	}{{
		name:      "Should pick existing grant according to preferences",
		available: 1 << 15,
		grants:    1<<10 | 1<<5,
		pref:      []uint32{1 << 15, 1 << 10},
		expected:  1 << 10,
	}, {
		name:      "Should pick a new grant if no existing",
		available: 1 << 15,
		grants:    0,
		pref:      []uint32{1 << 15, 1 << 10},
		expected:  1 << 15,
	}}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			actual := grant.SelectGrantsWithoutRedundancy(tt.available, tt.grants, tt.pref, tt.index)
			assert.Equal(t, tt.expected, actual, msg, tt.expected, actual)
		})
	}
}
