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

package grant

import (
	"sort"

	"magma/dp/cloud/go/services/dp/storage"
)

func CalcAvailableFrequencies(channels []storage.Channel, calc eirpCalculator) []uint32 {
	masks := make([]uint32, 4)
	sort.Slice(channels, func(i, j int) bool {
		return channels[i].LowFrequencyHz < channels[j].LowFrequencyHz
	})
	for i := 0; i < 4; i++ {
		bw := (i + 1) * unitToHz
		minEirp := calc.CalcLowerBound(bw)
		masks[i] = calcAvailableFrequenciesForBandwidth(channels, minEirp, int64(bw))
	}
	return masks
}

type eirpCalculator interface {
	CalcLowerBound(bandwidthHz int) float64
}

func calcAvailableFrequenciesForBandwidth(channels []storage.Channel, minEirp float64, band int64) uint32 {
	mask, begin, end := uint32(0), int64(0), int64(0)
	for _, c := range channels {
		if c.MaxEirp < minEirp {
			continue
		}
		if c.LowFrequencyHz > end {
			mask |= makeMaskForRange(begin, end, band)
			begin = c.LowFrequencyHz
		}
		if c.HighFrequencyHz > end {
			end = c.HighFrequencyHz
		}
	}
	return mask | makeMaskForRange(begin, end, band)
}

func makeMaskForRange(begin int64, end int64, band int64) uint32 {
	if end == 0 {
		return 0
	}
	l := hzToMask(begin + band/2 + unitToHz - 1)
	r := hzToMask(end - band/2)
	if r < l {
		return 0
	}
	return r<<1 - l
}

func UnsetGrantFrequency(cbsd *storage.DBCbsd, grant *storage.DBGrant) []uint32 {
	low := grant.LowFrequencyHz.Int64
	high := grant.HighFrequencyHz.Int64
	if cbsd.AvailableFrequencies == nil || low == 0 || high == 0 {
		return cbsd.AvailableFrequencies
	}

	frequencies := make([]uint32, len(cbsd.AvailableFrequencies))
	copy(frequencies, cbsd.AvailableFrequencies)

	bwHz := high - low
	mid := (low + high) / 2
	bitToUnset := (mid - int64(3550*1e6)) / int64(5*1e6)
	bwIndex := bwHz/int64(5*1e6) - 1

	frequencies[bwIndex] = frequencies[bwIndex] & ^(1 << bitToUnset)
	return frequencies
}
