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

package eirp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator/sas/eirp"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
)

func TestCalcLowerBound(t *testing.T) {
	cbsd := &storage.DBCbsd{
		MinPower:       db.MakeFloat(0),
		MaxPower:       db.MakeFloat(20),
		AntennaGainDbi: db.MakeFloat(15),
		NumberOfPorts:  db.MakeInt(2),
	}
	c := eirp.NewCalculator(cbsd)

	actual := c.CalcLowerBound(10 * 1e6)
	assert.Equal(t, 9.0, actual)
}

func TestCalculateUpperBoundForRange(t *testing.T) {
	data := []struct {
		name            string
		channels        []storage.Channel
		lowFrequencyHz  int64
		highFrequencyHz int64
		maxEirp         float64
		expected        float64
	}{{
		name: "Should calculate eirp for channel matching exactly",
		channels: []storage.Channel{{
			LowFrequencyHz:  3595e6,
			HighFrequencyHz: 3605e6,
			MaxEirp:         20,
		}},
		lowFrequencyHz:  3595e6,
		highFrequencyHz: 3605e6,
		maxEirp:         30,
		expected:        20,
	}, {
		name: "Should calculate eirp for non overlapping channels",
		channels: []storage.Channel{{
			LowFrequencyHz:  3590e6,
			HighFrequencyHz: 3600e6,
			MaxEirp:         25,
		}, {
			LowFrequencyHz:  3600e6,
			HighFrequencyHz: 3610e6,
			MaxEirp:         20,
		}},
		lowFrequencyHz:  3590e6,
		highFrequencyHz: 3610e6,
		maxEirp:         30,
		expected:        20,
	}, {
		name: "Should calculate eirp for overlapping channels",
		channels: []storage.Channel{{
			LowFrequencyHz:  3585e6,
			HighFrequencyHz: 3595e6,
			MaxEirp:         25,
		}, {
			LowFrequencyHz:  3590e6,
			HighFrequencyHz: 3600e6,
			MaxEirp:         15,
		}, {
			LowFrequencyHz:  3595e6,
			HighFrequencyHz: 3615e6,
			MaxEirp:         20,
		}},
		lowFrequencyHz:  3590e6,
		highFrequencyHz: 3610e6,
		maxEirp:         30,
		expected:        20,
	}, {
		name: "Should use given max eirp is it is smallest",
		channels: []storage.Channel{{
			LowFrequencyHz:  3590e6,
			HighFrequencyHz: 3610e6,
			MaxEirp:         25,
		}},
		lowFrequencyHz:  3590e6,
		highFrequencyHz: 3610e6,
		maxEirp:         20,
		expected:        20,
	}, {
		name: "Should use max sas eirp by default",
		channels: []storage.Channel{{
			LowFrequencyHz:  3590e6,
			HighFrequencyHz: 3610e6,
			MaxEirp:         37, // TODO can this be null?
		}},
		lowFrequencyHz:  3590e6,
		highFrequencyHz: 3610e6,
		maxEirp:         40,
		expected:        37,
	}, {
		name: "Should skip outside channels",
		channels: []storage.Channel{{
			LowFrequencyHz:  3590e6,
			HighFrequencyHz: 3610e6,
			MaxEirp:         37,
		}, {
			LowFrequencyHz:  3550e6,
			HighFrequencyHz: 3570e6,
			MaxEirp:         20,
		}},
		lowFrequencyHz:  3590e6,
		highFrequencyHz: 3610e6,
		maxEirp:         40,
		expected:        37,
	}}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			cbsd := &storage.DBCbsd{
				MaxPower:      db.MakeFloat(tt.maxEirp),
				NumberOfPorts: db.MakeInt((tt.highFrequencyHz - tt.lowFrequencyHz) / 1e6),
			}
			calc := eirp.NewCalculator(cbsd)
			actual := calc.CalcUpperBoundForRange(tt.channels, tt.lowFrequencyHz, tt.highFrequencyHz)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
