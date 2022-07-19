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

import "magma/dp/cloud/go/active_mode_controller/protos/active_mode"

type SelectionData struct {
	BandwidthHz   int64
	UseRedundancy RedundancyType
}
type RedundancyType uint8

const (
	NoRedundancy RedundancyType = iota
	MustHaveTwo
	BestEffort
)

func PickBandwidthSelectionOrder(settings *active_mode.GrantSettings, maxBandwidthHz int64, oldBandwidthHz int64) []*SelectionData {
	if oldBandwidthHz != 0 {
		redundancy := NoRedundancy
		if settings.GrantRedundancyEnabled {
			redundancy = BestEffort
		}
		return []*SelectionData{{
			BandwidthHz:   oldBandwidthHz,
			UseRedundancy: redundancy,
		}}
	}
	order := bandwidthSelectionOrder[2]
	if settings.GrantRedundancyEnabled {
		order = bandwidthSelectionOrder[1]
	}
	if settings.CarrierAggregationEnabled {
		order = bandwidthSelectionOrder[0]
	}
	return filterBandwidth(order, maxBandwidthHz)
}

var bandwidthSelectionOrder = [][]*SelectionData{{
	{BandwidthHz: 20e6, UseRedundancy: BestEffort},
	{BandwidthHz: 10e6, UseRedundancy: MustHaveTwo},
	{BandwidthHz: 15e6, UseRedundancy: BestEffort},
	{BandwidthHz: 10e6, UseRedundancy: NoRedundancy},
	{BandwidthHz: 5e6, UseRedundancy: BestEffort},
}, {
	{BandwidthHz: 20e6, UseRedundancy: BestEffort},
	{BandwidthHz: 15e6, UseRedundancy: BestEffort},
	{BandwidthHz: 10e6, UseRedundancy: BestEffort},
	{BandwidthHz: 5e6, UseRedundancy: BestEffort},
}, {
	{BandwidthHz: 20e6, UseRedundancy: NoRedundancy},
	{BandwidthHz: 15e6, UseRedundancy: NoRedundancy},
	{BandwidthHz: 10e6, UseRedundancy: NoRedundancy},
	{BandwidthHz: 5e6, UseRedundancy: NoRedundancy},
}}

func filterBandwidth(data []*SelectionData, bandwidthHz int64) []*SelectionData {
	filtered := make([]*SelectionData, 0, len(data))
	for _, d := range data {
		if d.BandwidthHz <= bandwidthHz {
			filtered = append(filtered, d)
		}
	}
	return filtered
}
