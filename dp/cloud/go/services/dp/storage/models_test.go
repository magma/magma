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
			name:     "check field names for DBRequestType",
			model:    &storage.DBRequestType{},
			expected: []string{"id", "name"},
		},
		{
			name:     "check field names for DBRequestState",
			model:    &storage.DBRequestState{},
			expected: []string{"id", "name"},
		},
		{
			name:  "check field names for DBRequest",
			model: &storage.DBRequest{},
			expected: []string{
				"id", "type_id", "state_id", "cbsd_id", "payload",
			},
		},
		{
			name:  "check field names for DBResponse",
			model: &storage.DBResponse{},
			expected: []string{
				"id", "request_id", "grant_id", "response_code", "payload",
			},
		},
		{
			name:     "check field names for DBGrantState",
			model:    &storage.DBGrantState{},
			expected: []string{"id", "name"},
		},
		{
			name:  "check field names for DBGrant",
			model: &storage.DBGrant{},
			expected: []string{
				"id", "state_id", "cbsd_id", "channel_id", "grant_id",
				"grant_expire_time", "transmit_expire_time",
				"heartbeat_interval", "channel_type",
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
				"fcc_id", "cbsd_serial_number", "last_seen", "min_power",
				"max_power", "antenna_gain", "number_of_ports",
				"is_deleted", "is_updated",
			},
		},
		{
			name:  "check field names for DBChannel",
			model: &storage.DBChannel{},
			expected: []string{
				"id", "cbsd_id", "low_frequency", "high_frequency",
				"channel_type", "rule_applied", "max_eirp",
				"last_used_max_eirp",
			},
		},
		{
			name:     "check field names for DBActiveModeConfig",
			model:    &storage.DBActiveModeConfig{},
			expected: []string{"id", "cbsd_id", "desired_state_id"},
		},
		{
			name:  "check field names for DBLog",
			model: &storage.DBLog{},
			expected: []string{
				"id", "network_id", "log_from", "log_to", "log_name",
				"log_message", "cbsd_serial_number", "fcc_id",
				"response_code", "created_date",
			},
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
			name:  "check ModelMetadata structure for DBRequestType",
			model: &storage.DBRequestType{},
			expected: db.ModelMetadata{
				Table:     storage.RequestTypeTable,
				Relations: map[string]string{},
			},
		},
		{
			name:  "check ModelMetadata structure for DBRequestState",
			model: &storage.DBRequestState{},
			expected: db.ModelMetadata{
				Table:     storage.RequestStateTable,
				Relations: map[string]string{},
			},
		},
		{
			name:  "check ModelMetadata structure for DBRequest",
			model: &storage.DBRequest{},
			expected: db.ModelMetadata{
				Table: storage.RequestTable,
				Relations: map[string]string{
					storage.RequestStateTable: "state_id",
					storage.RequestTypeTable:  "type_id",
					storage.CbsdTable:         "cbsd_id",
				},
			},
		},
		{
			name:  "check ModelMetadata structure for DBResponse",
			model: &storage.DBResponse{},
			expected: db.ModelMetadata{
				Table: storage.ResponseTable,
				Relations: map[string]string{
					storage.RequestTable: "request_id",
					storage.GrantTable:   "grant_id",
				},
			},
		},
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
					storage.ChannelTable:    "channel_id",
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
			name:  "check ModelMetadata structure for DBChannel",
			model: &storage.DBChannel{},
			expected: db.ModelMetadata{
				Table: storage.ChannelTable,
				Relations: map[string]string{
					storage.CbsdTable: "cbsd_id",
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
		{
			name:  "check ModelMetadata structure for DBLog",
			model: &storage.DBLog{},
			expected: db.ModelMetadata{
				Table: storage.LogTable,
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
