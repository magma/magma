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
	"google.golang.org/protobuf/types/known/wrapperspb"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas/eirp"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func TestCalcLowerBound(t *testing.T) {
	ec := &active_mode.EirpCapabilities{
		MinPower:      0,
		MaxPower:      20,
		NumberOfPorts: 2,
	}
	c := eirp.NewCalculator(15, ec)

	actual := c.CalcLowerBound(10 * 1e6)
	assert.Equal(t, 9.0, actual)
}

func TestCalculateUpperBoundForRange(t *testing.T) {
	data := []struct {
		name            string
		channels        []*active_mode.Channel
		lowFrequencyHz  int64
		highFrequencyHz int64
		maxEirp         float64
		expected        float64
	}{{
		name: "Should calculate eirp for channel matching exactly",
		channels: []*active_mode.Channel{{
			LowFrequencyHz:  3595e6,
			HighFrequencyHz: 3605e6,
			MaxEirp:         wrapperspb.Float(20),
		}},
		lowFrequencyHz:  3595e6,
		highFrequencyHz: 3605e6,
		maxEirp:         30,
		expected:        20,
	}, {
		name: "Should calculate eirp for non overlapping channels",
		channels: []*active_mode.Channel{{
			LowFrequencyHz:  3590e6,
			HighFrequencyHz: 3600e6,
			MaxEirp:         wrapperspb.Float(25),
		}, {
			LowFrequencyHz:  3600e6,
			HighFrequencyHz: 3610e6,
			MaxEirp:         wrapperspb.Float(20),
		}},
		lowFrequencyHz:  3590e6,
		highFrequencyHz: 3610e6,
		maxEirp:         30,
		expected:        20,
	}, {
		name: "Should calculate eirp for overlapping channels",
		channels: []*active_mode.Channel{{
			LowFrequencyHz:  3585e6,
			HighFrequencyHz: 3595e6,
			MaxEirp:         wrapperspb.Float(25),
		}, {
			LowFrequencyHz:  3590e6,
			HighFrequencyHz: 3600e6,
			MaxEirp:         wrapperspb.Float(15),
		}, {
			LowFrequencyHz:  3595e6,
			HighFrequencyHz: 3615e6,
			MaxEirp:         wrapperspb.Float(20),
		}},
		lowFrequencyHz:  3590e6,
		highFrequencyHz: 3610e6,
		maxEirp:         30,
		expected:        20,
	}, {
		name: "Should use given max eirp is it is smallest",
		channels: []*active_mode.Channel{{
			LowFrequencyHz:  3590e6,
			HighFrequencyHz: 3610e6,
			MaxEirp:         wrapperspb.Float(25),
		}},
		lowFrequencyHz:  3590e6,
		highFrequencyHz: 3610e6,
		maxEirp:         20,
		expected:        20,
	}, {
		name: "Should use max sas eirp by default",
		channels: []*active_mode.Channel{{
			LowFrequencyHz:  3590e6,
			HighFrequencyHz: 3610e6,
		}},
		lowFrequencyHz:  3590e6,
		highFrequencyHz: 3610e6,
		maxEirp:         40,
		expected:        37,
	}, {
		name: "Should skip outside channels",
		channels: []*active_mode.Channel{{
			LowFrequencyHz:  3590e6,
			HighFrequencyHz: 3610e6,
		}, {
			LowFrequencyHz:  3550e6,
			HighFrequencyHz: 3570e6,
			MaxEirp:         wrapperspb.Float(20),
		}},
		lowFrequencyHz:  3590e6,
		highFrequencyHz: 3610e6,
		maxEirp:         40,
		expected:        37,
	}}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			c := &active_mode.EirpCapabilities{
				MaxPower:      float32(tt.maxEirp),
				NumberOfPorts: int32(tt.highFrequencyHz-tt.lowFrequencyHz) / 1e6,
			}
			calc := eirp.NewCalculator(0, c)
			actual := calc.CalcUpperBoundForRange(tt.channels, tt.lowFrequencyHz, tt.highFrequencyHz)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
