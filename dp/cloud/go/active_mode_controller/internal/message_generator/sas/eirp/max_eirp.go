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

package eirp

import (
	"golang.org/x/exp/constraints"

	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

const (
	minSASEirp = -137
	maxSASEirp = 37
)

func CalculateMaxEirp(channels []*active_mode.Channel, lowFrequencyHz int64, highFrequencyHz int64, maxEirp float32) float32 {
	bw := int((highFrequencyHz - lowFrequencyHz) / 1e6)
	eirps := make([]float32, bw+1)
	for i := range eirps {
		eirps[i] = minSASEirp
	}
	for _, c := range channels {
		if c.HighFrequencyHz >= lowFrequencyHz && c.LowFrequencyHz <= highFrequencyHz {
			updateMaxEirpsForChannel(c, eirps, lowFrequencyHz, highFrequencyHz)
		}
	}
	eirp := min(maxEirp, float32(maxSASEirp))
	for _, e := range eirps {
		eirp = min(eirp, e)
	}
	return eirp
}

func updateMaxEirpsForChannel(c *active_mode.Channel, eirps []float32, lowFrequencyHz int64, highFrequencyHz int64) {
	low := max(lowFrequencyHz, c.LowFrequencyHz)
	high := min(highFrequencyHz, c.HighFrequencyHz)
	l := int((low - lowFrequencyHz + 1e6 - 1) / 1e6)
	r := int((high - lowFrequencyHz) / 1e6)
	eirp := float32(maxSASEirp)
	if c.MaxEirp != nil {
		eirp = c.MaxEirp.Value
	}
	for ; l <= r; l++ {
		eirps[l] = max(eirps[l], eirp)
	}

}

func min[T constraints.Ordered](a T, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T constraints.Ordered](a T, b T) T {
	if a > b {
		return a
	}
	return b
}
