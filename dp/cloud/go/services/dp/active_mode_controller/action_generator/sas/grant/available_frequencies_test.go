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
	"magma/dp/cloud/go/services/dp/storage/db"
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

func TestUnsetGrantFrequency(t *testing.T) {
	testData := []struct {
		origAvailFreq     []uint32
		lowFreq           int64
		highFreq          int64
		expectedAvailFreq []uint32
	}{{
		nil,
		3560e6,
		3580e6,
		nil,
	}, {
		[]uint32{0b1111, 0b110, 0b1100, 0b1010},
		35625e5,
		35675e5,
		[]uint32{0b0111, 0b110, 0b1100, 0b1010},
	}, {
		[]uint32{0b0, 0b110, 0b1100, 0b1010},
		3550e6,
		3560e6,
		[]uint32{0b0, 0b100, 0b1100, 0b1010},
	}, {
		[]uint32{0b0, 0b110, 0b1111100, 0b1010},
		35725e5,
		35875e5,
		[]uint32{0b0, 0b110, 0b0111100, 0b1010},
	}, {
		[]uint32{0b0, 0b110, 0b1111100, 0b10101},
		3560e6,
		3580e6,
		[]uint32{0b0, 0b110, 0b1111100, 0b00101},
	},
	}
	for _, tt := range testData {
		t.Run("", func(t *testing.T) {
			cbsd := &storage.DBCbsd{
				Id:                   db.MakeInt(1),
				AvailableFrequencies: tt.origAvailFreq,
			}
			gt := &storage.DBGrant{
				CbsdId:          cbsd.Id,
				LowFrequencyHz:  db.MakeInt(tt.lowFreq),
				HighFrequencyHz: db.MakeInt(tt.highFreq),
			}
			newFrequencies := grant.UnsetGrantFrequency(cbsd, gt)

			assert.Equal(t, tt.expectedAvailFreq, newFrequencies)
		})
	}
}

type stubEirpCalculator struct {
	eirps []float64
}

func (s *stubEirpCalculator) CalcLowerBound(bandwidthHz int) float64 {
	return s.eirps[(bandwidthHz/5e6)-1]
}
