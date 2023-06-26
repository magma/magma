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
	"magma/dp/cloud/go/services/dp/storage"
)

func TestCalcAvailableFrequencies(t *testing.T) {
	testData := []struct {
		name     string
		channels []storage.Channel
		eirps    []float64
		expected []uint32
	}{{
		name:     "Should handle no channels",
		channels: nil,
		eirps:    []float64{0, 0, 0, 0},
		expected: []uint32{0, 0, 0, 0},
	}, {
		name: "Should handle single channel",
		channels: []storage.Channel{{
			LowFrequencyHz:  3590e6,
			HighFrequencyHz: 3610e6,
			MaxEirp:         37,
		}},
		eirps: []float64{0, 0, 0, 0},
		expected: []uint32{
			1<<9 | 1<<10 | 1<<11,
			1<<9 | 1<<10 | 1<<11,
			1 << 10,
			1 << 10,
		},
	}, {
		name: "Should handle joined channels",
		channels: []storage.Channel{{
			LowFrequencyHz:  3590e6,
			HighFrequencyHz: 3600e6,
			MaxEirp:         37,
		}, {
			LowFrequencyHz:  3600e6,
			HighFrequencyHz: 3610e6,
			MaxEirp:         37,
		}},
		eirps: []float64{0, 0, 0, 0},
		expected: []uint32{
			1<<9 | 1<<10 | 1<<11,
			1<<9 | 1<<10 | 1<<11,
			1 << 10,
			1 << 10,
		},
	}, {
		name: "Should handle disjoint channels",
		channels: []storage.Channel{{
			LowFrequencyHz:  3590e6,
			HighFrequencyHz: 3600e6,
			MaxEirp:         37,
		}, {
			LowFrequencyHz:  3610e6,
			HighFrequencyHz: 3620e6,
			MaxEirp:         37,
		}},
		eirps: []float64{0, 0, 0, 0},
		expected: []uint32{
			1<<9 | 1<<13,
			1<<9 | 1<<13,
			0,
			0,
		},
	}, {
		name: "Should handle nested channels",
		channels: []storage.Channel{{
			LowFrequencyHz:  3595e6,
			HighFrequencyHz: 3605e6,
			MaxEirp:         37,
		}, {
			LowFrequencyHz:  3590e6,
			HighFrequencyHz: 3610e6,
			MaxEirp:         37,
		}},
		eirps: []float64{0, 0, 0, 0},
		expected: []uint32{
			1<<9 | 1<<10 | 1<<11,
			1<<9 | 1<<10 | 1<<11,
			1 << 10,
			1 << 10,
		},
	}, {
		name: "Should handle borders",
		channels: []storage.Channel{{
			LowFrequencyHz:  3550e6,
			HighFrequencyHz: 3700e6,
			MaxEirp:         37,
		}},
		eirps: []float64{0, 0, 0, 0},
		expected: []uint32{
			1<<30 - 1<<1,
			1<<30 - 1<<1,
			1<<29 - 1<<2,
			1<<29 - 1<<2,
		},
	},
		{
			name: "Should calculate channels not aligned to multiple of 5MHz",
			channels: []storage.Channel{{
				LowFrequencyHz:  3591e6,
				HighFrequencyHz: 3600e6,
				MaxEirp:         37,
			}, {
				LowFrequencyHz:  3610e6,
				HighFrequencyHz: 3629e6,
				MaxEirp:         37,
			}},
			eirps: []float64{0, 0, 0, 0},
			expected: []uint32{
				1<<9 | 1<<13 | 1<<14 | 1<<15,
				1<<13 | 1<<14,
				1 << 14,
				0,
			},
		}, {
			name: "Should skip channels with too low eirp",
			channels: []storage.Channel{{
				LowFrequencyHz:  3590e6,
				HighFrequencyHz: 3610e6,
				MaxEirp:         -1,
			}},
			eirps:    []float64{0, 0, 0, 0},
			expected: []uint32{0, 0, 0, 0},
		}, {
			name: "Should use correct eirp per bandwidth",
			channels: []storage.Channel{{
				LowFrequencyHz:  3590e6,
				HighFrequencyHz: 3600e6,
				MaxEirp:         5,
			}, {
				LowFrequencyHz:  3600e6,
				HighFrequencyHz: 3610e6,
				MaxEirp:         10,
			}},
			eirps: []float64{11, 10, 9, 5},
			expected: []uint32{
				0,
				1 << 11,
				0,
				1 << 10,
			},
		}}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			calc := &stubEirpCalculator{eirps: tt.eirps}
			actual := grant.CalcAvailableFrequencies(tt.channels, calc)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

type stubEirpCalculator struct {
	eirps []float64
}

func (s *stubEirpCalculator) CalcLowerBound(bandwidthHz int) float64 {
	return s.eirps[(bandwidthHz/5e6)-1]
}
