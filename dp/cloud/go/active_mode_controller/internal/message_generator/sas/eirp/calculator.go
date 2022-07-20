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

func (c *calculator) CalcUpperBound(bandwidthHz int) float64 {
	return math.Floor(c.calcEirp(c.maxPower, bandwidthHz))
}

func (c *calculator) calcEirp(power float64, bandwidthHz int) float64 {
	bwMHz := float64(bandwidthHz / 1e6)
	return power + c.antennaGain - 10*math.Log10(bwMHz/c.noPorts)
}
