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
	"magma/orc8r/cloud/go/sqorc"
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
				"is_deleted", "should_deregister",
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
				Table: storage.GrantStateTable,
				Properties: map[string]*db.Field{
					"id": {
						SqlType: sqorc.ColumnTypeInt,
					},
					"name": {
						SqlType: sqorc.ColumnTypeText,
					},
				},
				Relations: map[string]string{},
			},
		},
		{
			name:  "check ModelMetadata structure for DBGrant",
			model: &storage.DBGrant{},
			expected: db.ModelMetadata{
				Table: storage.GrantTable,
				Properties: map[string]*db.Field{
					"id": {
						SqlType: sqorc.ColumnTypeInt,
					},
					"state_id": {
						SqlType: sqorc.ColumnTypeInt,
					},
					"cbsd_id": {
						SqlType:  sqorc.ColumnTypeInt,
						Nullable: true,
					},
					"grant_id": {
						SqlType: sqorc.ColumnTypeText,
					},
					"grant_expire_time": {
						SqlType:  sqorc.ColumnTypeDatetime,
						Nullable: true,
					},
					"transmit_expire_time": {
						SqlType:  sqorc.ColumnTypeDatetime,
						Nullable: true,
					},
					"heartbeat_interval": {
						SqlType:  sqorc.ColumnTypeInt,
						Nullable: true,
					},
					"channel_type": {
						SqlType:  sqorc.ColumnTypeText,
						Nullable: true,
					},
					"low_frequency": {
						SqlType: sqorc.ColumnTypeInt,
					},
					"high_frequency": {
						SqlType: sqorc.ColumnTypeInt,
					},
					"max_eirp": {
						SqlType: sqorc.ColumnTypeReal,
					},
				},
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
				Table: storage.CbsdStateTable,
				Properties: map[string]*db.Field{
					"id": {
						SqlType: sqorc.ColumnTypeInt,
					},
					"name": {
						SqlType: sqorc.ColumnTypeText,
					},
				},
				Relations: map[string]string{},
			},
		},
		{
			name:  "check ModelMetadata structure for DBCbsd",
			model: &storage.DBCbsd{},
			expected: db.ModelMetadata{
				Table: storage.CbsdTable,
				Properties: map[string]*db.Field{
					"id": {
						SqlType: sqorc.ColumnTypeInt,
					},
					"network_id": {
						SqlType: sqorc.ColumnTypeText,
					},
					"state_id": {
						SqlType: sqorc.ColumnTypeInt,
					},
					"cbsd_id": {
						SqlType:  sqorc.ColumnTypeText,
						Nullable: true,
					},
					"user_id": {
						SqlType:  sqorc.ColumnTypeText,
						Nullable: true,
					},
					"fcc_id": {
						SqlType:  sqorc.ColumnTypeText,
						Nullable: true,
					},
					"cbsd_serial_number": {
						SqlType:  sqorc.ColumnTypeText,
						Nullable: true,
						Unique:   true,
					},
					"last_seen": {
						SqlType:  sqorc.ColumnTypeDatetime,
						Nullable: true,
					},
					"grant_attempts": {
						SqlType:      sqorc.ColumnTypeInt,
						HasDefault:   true,
						DefaultValue: 0,
					},
					"preferred_bandwidth_mhz": {
						SqlType: sqorc.ColumnTypeInt,
					},
					"preferred_frequencies_mhz": {
						SqlType: sqorc.ColumnTypeText,
					},
					"min_power": {
						SqlType:  sqorc.ColumnTypeReal,
						Nullable: true,
					},
					"max_power": {
						SqlType:  sqorc.ColumnTypeReal,
						Nullable: true,
					},
					"antenna_gain": {
						SqlType:  sqorc.ColumnTypeReal,
						Nullable: true,
					},
					"number_of_ports": {
						SqlType:  sqorc.ColumnTypeInt,
						Nullable: true,
					},
					"is_deleted": {
						SqlType:      sqorc.ColumnTypeBool,
						HasDefault:   true,
						DefaultValue: false,
					},
					"should_deregister": {
						SqlType:      sqorc.ColumnTypeBool,
						HasDefault:   true,
						DefaultValue: false,
					},
				},
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
				Properties: map[string]*db.Field{
					"id": {
						SqlType: sqorc.ColumnTypeInt,
					},
					"cbsd_id": {
						SqlType: sqorc.ColumnTypeInt,
					},
					"desired_state_id": {
						SqlType: sqorc.ColumnTypeInt,
					},
				},
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
			assert.Equal(t, tc.expected.Properties, actual.Properties)
			assert.Equal(t, tc.expected.Table, actual.Table)
			assert.Equal(t, tc.model, obj)
		})
	}
}
