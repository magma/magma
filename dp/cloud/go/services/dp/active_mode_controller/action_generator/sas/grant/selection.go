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
	"math/bits"

	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator/sas/frequency"
	"magma/dp/cloud/go/services/dp/storage"
)

type Processor[T any] interface {
	ProcessGrant(frequency int64, bandwidth int64) T
}

type Processors[T any] struct {
	Keep Processor[T]
	Del  Processor[T]
	Add  Processor[T]
}

func ProcessGrants[T any](cbsd *storage.DBCbsd, grants []*storage.DBGrant, processors Processors[T], index int) []T {
	oldBw, oldGrants := calculateOldGrants(grants)
	bw, newGrants := selectGrants(cbsd, oldGrants, oldBw, index)
	return processGrants(processors, oldGrants, newGrants, bw)
}

func calculateOldGrants(grants []*storage.DBGrant) (int64, uint32) {
	mask := uint32(0)
	bw := int64(0)
	for i, g := range grants {
		mask |= hzToMask((g.HighFrequencyHz.Int64 + g.LowFrequencyHz.Int64) / 2)
		if i == 0 {
			bw = g.HighFrequencyHz.Int64 - g.LowFrequencyHz.Int64
		}
	}
	return bw, mask
}

func hzToMask(hz int64) uint32 {
	return 1 << ((hz - frequency.LowestHz) / unitToHz)
}

const unitToHz = 5e6

func selectGrants(cbsd *storage.DBCbsd, oldGrants uint32, oldBandwidthHz int64, index int) (int64, uint32) {
	prefMask := preferencesToMask(cbsd.PreferredFrequenciesMHz)
	order := PickBandwidthSelectionOrder(cbsd, cbsd.PreferredBandwidthMHz.Int64*1e6, oldBandwidthHz)
	for _, o := range order {
		newGrants := selectGrantsForBandwidth(o, cbsd, oldGrants, prefMask, index)
		if newGrants != 0 {
			return o.BandwidthHz, newGrants
		}
	}
	return 0, 0
}

func preferencesToMask(frequenciesMHz []int64) []uint32 {
	masks := make([]uint32, len(frequenciesMHz))
	for i, f := range frequenciesMHz {
		masks[i] = hzToMask(f * 1e6)
	}
	return masks
}

func selectGrantsForBandwidth(data *SelectionData, cbsd *storage.DBCbsd, grants uint32, pref []uint32, index int) uint32 {
	minWidth := int(data.BandwidthHz/unitToHz - 1)
	maxWidth := int((cbsd.MaxIbwMhx.Int64*1e6 - data.BandwidthHz) / unitToHz)
	available := cbsd.AvailableFrequencies[minWidth]
	if minWidth > maxWidth {
		maxWidth = minWidth
	}
	if data.UseRedundancy == NoRedundancy {
		return SelectGrantsWithoutRedundancy(available, grants, pref, index)
	}
	newGrants := SelectGrantsWithRedundancy(available, grants, pref, minWidth, maxWidth, index)
	if newGrants != 0 || data.UseRedundancy == MustHaveTwo {
		return newGrants
	}
	return SelectGrantsWithoutRedundancy(available, grants, pref, index)
}

func processGrants[T any](processors Processors[T], oldGrants uint32, newGrants uint32, bandwidthHz int64) []T {
	toKeep := oldGrants & newGrants
	toDel := oldGrants &^ newGrants
	toAdd := newGrants &^ oldGrants
	r1 := processTypedGrants(processors.Del, toDel, bandwidthHz)
	r2 := processTypedGrants(processors.Add, toAdd, bandwidthHz)
	r3 := processTypedGrants(processors.Keep, toKeep, bandwidthHz)
	return append(append(r1, r2...), r3...)
}

func processTypedGrants[T any](processor Processor[T], grants uint32, bandwidthHz int64) []T {
	reqs := make([]T, 0, bits.OnesCount32(grants))
	for grants > 0 {
		x := grants & -grants
		grants -= x
		reqs = append(reqs, processor.ProcessGrant(maskToHz(x), bandwidthHz))
	}
	return reqs
}

func maskToHz(mask uint32) int64 {
	return int64(bits.TrailingZeros32(mask))*unitToHz + frequency.LowestHz
}
