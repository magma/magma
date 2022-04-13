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

package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
)

func TestFields(t *testing.T) {
	testCases := []struct {
		name     string
		model    db.Model
		expected []string
	}{
		{
			name:     "check field names for DBGrantState",
			model:    &storage.DBGrantState{},
			expected: []string{"id", "name"},
		},
		{
			name:  "check field names for DBGrant",
			model: &storage.DBGrant{},
			expected: []string{
				"id", "state_id", "cbsd_id", "grant_id",
				"grant_expire_time", "transmit_expire_time",
				"heartbeat_interval", "channel_type",
				"low_frequency", "high_frequency", "max_eirp",
			},
		},
		{
			name:     "check field names for DBCbsdState",
			model:    &storage.DBCbsdState{},
			expected: []string{"id", "name"},
		},
		{
			name:  "check field names for DBCbsd",
			model: &storage.DBCbsd{},
			expected: []string{
				"id", "network_id", "state_id", "cbsd_id", "user_id",
				"fcc_id", "cbsd_serial_number", "last_seen", "grant_attempts",
				"preferred_bandwidth_mhz", "preferred_frequencies_mhz",
				"min_power", "max_power", "antenna_gain", "number_of_ports",
				"is_deleted", "is_updated",
			},
		},
		{
			name:     "check field names for DBActiveModeConfig",
			model:    &storage.DBActiveModeConfig{},
			expected: []string{"id", "cbsd_id", "desired_state_id"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.model.Fields()

			keys := make([]string, 0, len(actual))
			for k := range actual {
				keys = append(keys, k)
			}
			assert.ElementsMatch(t, tc.expected, keys)
		})
	}
}

func TestGetMetadata(t *testing.T) {
	testCases := []struct {
		name     string
		model    db.Model
		expected db.ModelMetadata
	}{
		{
			name:  "check ModelMetadata structure for DBGrantState",
			model: &storage.DBGrantState{},
			expected: db.ModelMetadata{
				Table:     storage.GrantStateTable,
				Relations: map[string]string{},
			},
		},
		{
			name:  "check ModelMetadata structure for DBGrant",
			model: &storage.DBGrant{},
			expected: db.ModelMetadata{
				Table: storage.GrantTable,
				Relations: map[string]string{
					storage.CbsdTable:       "cbsd_id",
					storage.GrantStateTable: "state_id",
				},
			},
		},
		{
			name:  "check ModelMetadata structure for DBCbsdState",
			model: &storage.DBCbsdState{},
			expected: db.ModelMetadata{
				Table:     storage.CbsdStateTable,
				Relations: map[string]string{},
			},
		},
		{
			name:  "check ModelMetadata structure for DBCbsd",
			model: &storage.DBCbsd{},
			expected: db.ModelMetadata{
				Table: storage.CbsdTable,
				Relations: map[string]string{
					storage.CbsdStateTable: "state_id",
				},
			},
		},
		{
			name:  "check ModelMetadata structure for DBActiveModeConfig",
			model: &storage.DBActiveModeConfig{},
			expected: db.ModelMetadata{
				Table: storage.ActiveModeConfigTable,
				Relations: map[string]string{
					storage.CbsdTable:      "cbsd_id",
					storage.CbsdStateTable: "desired_state_id",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.model.GetMetadata()
			obj := actual.CreateObject()

			assert.Equal(t, tc.expected.Relations, actual.Relations)
			assert.Equal(t, tc.expected.Table, actual.Table)
			assert.Equal(t, tc.model, obj)
		})
	}
}
