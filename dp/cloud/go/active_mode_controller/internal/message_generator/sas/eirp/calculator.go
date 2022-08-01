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
	"math"

	"golang.org/x/exp/constraints"

	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

type calculator struct {
	minPower    float64
	maxPower    float64
	antennaGain float64
	noPorts     float64
}

func NewCalculator(antennaGain float32, capabilities *active_mode.EirpCapabilities) *calculator {
	return &calculator{
		minPower:    float64(capabilities.GetMinPower()),
		maxPower:    float64(capabilities.GetMaxPower()),
		antennaGain: float64(antennaGain),
		noPorts:     float64(capabilities.GetNumberOfPorts()),
	}
}

func (c *calculator) CalcLowerBound(bandwidthHz int) float64 {
	return math.Ceil(c.calcEirp(c.minPower, bandwidthHz))
}

func (c *calculator) CalcUpperBoundForRange(channels []*active_mode.Channel, low int64, high int64) float64 {
	eirp := c.calcUpperBound(int(high - low))
	return calculateMaxEirp(channels, low, high, eirp)
}

func (c *calculator) calcEirp(power float64, bandwidthHz int) float64 {
	bwMHz := float64(bandwidthHz / 1e6)
	return power + c.antennaGain - 10*math.Log10(bwMHz/c.noPorts)
}

func (c *calculator) calcUpperBound(bandwidthHz int) float64 {
	return math.Floor(c.calcEirp(c.maxPower, bandwidthHz))
}

const (
	minSASEirp = -137
	maxSASEirp = 37
)

func calculateMaxEirp(channels []*active_mode.Channel, lowFrequencyHz int64, highFrequencyHz int64, maxEirp float64) float64 {
	bw := int((highFrequencyHz - lowFrequencyHz) / 1e6)
	eirps := make([]float64, bw+1)
	for i := range eirps {
		eirps[i] = minSASEirp
	}
	for _, c := range channels {
		if c.HighFrequencyHz >= lowFrequencyHz && c.LowFrequencyHz <= highFrequencyHz {
			updateMaxEirpsForChannel(c, eirps, lowFrequencyHz, highFrequencyHz)
		}
	}
	eirp := min(maxEirp, maxSASEirp)
	for _, e := range eirps {
		eirp = min(eirp, e)
	}
	return eirp
}

func updateMaxEirpsForChannel(c *active_mode.Channel, eirps []float64, lowFrequencyHz int64, highFrequencyHz int64) {
	low := max(lowFrequencyHz, c.LowFrequencyHz)
	high := min(highFrequencyHz, c.HighFrequencyHz)
	l := int((low - lowFrequencyHz + 1e6 - 1) / 1e6)
	r := int((high - lowFrequencyHz) / 1e6)
	eirp := float64(maxSASEirp)
	if c.MaxEirp != nil {
		eirp = float64(c.MaxEirp.Value)
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
