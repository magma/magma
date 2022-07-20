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

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas/grant"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func TestPickBandwidthSelectionOrder(t *testing.T) {
	testData := []struct {
		name           string
		settings       *active_mode.GrantSettings
		maxBandwidthHz int64
		oldBandwidthHz int64
		expected       []*grant.SelectionData
	}{{
		name:           "Should pick no redundancy order when redundancy is disabled",
		settings:       &active_mode.GrantSettings{},
		maxBandwidthHz: 20e6,
		expected: []*grant.SelectionData{
			{BandwidthHz: 20e6, UseRedundancy: grant.NoRedundancy},
			{BandwidthHz: 15e6, UseRedundancy: grant.NoRedundancy},
			{BandwidthHz: 10e6, UseRedundancy: grant.NoRedundancy},
			{BandwidthHz: 5e6, UseRedundancy: grant.NoRedundancy},
		},
	}, {
		name: "Should pick best effort order for redundancy without carrier aggregation",
		settings: &active_mode.GrantSettings{
			GrantRedundancyEnabled: true,
		},
		maxBandwidthHz: 20e6,
		expected: []*grant.SelectionData{
			{BandwidthHz: 20e6, UseRedundancy: grant.BestEffort},
			{BandwidthHz: 15e6, UseRedundancy: grant.BestEffort},
			{BandwidthHz: 10e6, UseRedundancy: grant.BestEffort},
			{BandwidthHz: 5e6, UseRedundancy: grant.BestEffort},
		},
	}, {
		name: "Should pick custom order for carrier aggregation",
		settings: &active_mode.GrantSettings{
			GrantRedundancyEnabled:    true,
			CarrierAggregationEnabled: true,
		},
		maxBandwidthHz: 20e6,
		expected: []*grant.SelectionData{
			{BandwidthHz: 20e6, UseRedundancy: grant.BestEffort},
			{BandwidthHz: 10e6, UseRedundancy: grant.MustHaveTwo},
			{BandwidthHz: 15e6, UseRedundancy: grant.BestEffort},
			{BandwidthHz: 10e6, UseRedundancy: grant.NoRedundancy},
			{BandwidthHz: 5e6, UseRedundancy: grant.BestEffort},
		},
	}, {
		name: "Should filter out too large bandwidths",
		settings: &active_mode.GrantSettings{
			GrantRedundancyEnabled:    true,
			CarrierAggregationEnabled: true,
		},
		maxBandwidthHz: 10e6,
		expected: []*grant.SelectionData{
			{BandwidthHz: 10e6, UseRedundancy: grant.MustHaveTwo},
			{BandwidthHz: 10e6, UseRedundancy: grant.NoRedundancy},
			{BandwidthHz: 5e6, UseRedundancy: grant.BestEffort},
		},
	}, {
		name:           "Should pick no redundancy for existing bandwidth without redundancy",
		settings:       &active_mode.GrantSettings{},
		maxBandwidthHz: 20e6,
		oldBandwidthHz: 15e6,
		expected: []*grant.SelectionData{
			{BandwidthHz: 15e6, UseRedundancy: grant.NoRedundancy},
		},
	}, {
		name: "Should pick best effort for existing bandwidth with redundancy",
		settings: &active_mode.GrantSettings{
			GrantRedundancyEnabled: true,
		},
		maxBandwidthHz: 20e6,
		oldBandwidthHz: 15e6,
		expected: []*grant.SelectionData{
			{BandwidthHz: 15e6, UseRedundancy: grant.BestEffort},
		},
	}}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			actual := grant.PickBandwidthSelectionOrder(tt.settings, tt.maxBandwidthHz, tt.oldBandwidthHz)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
