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
	dbGrant := &storage.DBGrant{}
	dbActiveModeConfig := &storage.DBActiveModeConfig{}
	dbCbsd := &storage.DBCbsd{}
	dbCbsdState := &storage.DBCbsdState{}
	dbGrantState := &storage.DBGrantState{}
	testCases := []struct {
		name     string
		model    db.Model
		expected []db.BaseType
	}{
		{
			name:  "check field names for DBGrantState",
			model: dbGrantState,
			expected: []db.BaseType{
				db.IntType{X: &dbGrantState.Id},
				db.StringType{X: &dbGrantState.Name},
			},
		},
		{
			name:  "check field names for DBGrant",
			model: dbGrant,
			expected: []db.BaseType{
				db.IntType{X: &dbGrant.Id},
				db.IntType{X: &dbGrant.StateId},
				db.IntType{X: &dbGrant.CbsdId},
				db.StringType{X: &dbGrant.GrantId},
				db.TimeType{X: &dbGrant.GrantExpireTime},
				db.TimeType{X: &dbGrant.TransmitExpireTime},
				db.IntType{X: &dbGrant.HeartbeatInterval},
				db.StringType{X: &dbGrant.ChannelType},
				db.IntType{X: &dbGrant.LowFrequency},
				db.IntType{X: &dbGrant.HighFrequency},
				db.FloatType{X: &dbGrant.MaxEirp},
			},
		},
		{
			name:  "check field names for DBCbsdState",
			model: dbCbsdState,
			expected: []db.BaseType{
				db.IntType{X: &dbCbsdState.Id},
				db.StringType{X: &dbCbsdState.Name},
			},
		},
		{
			name:  "check field names for DBCbsd",
			model: dbCbsd,
			expected: []db.BaseType{
				db.IntType{X: &dbCbsd.Id},
				db.StringType{X: &dbCbsd.NetworkId},
				db.IntType{X: &dbCbsd.StateId},
				db.StringType{X: &dbCbsd.CbsdId},
				db.StringType{X: &dbCbsd.UserId},
				db.StringType{X: &dbCbsd.FccId},
				db.StringType{X: &dbCbsd.CbsdSerialNumber},
				db.TimeType{X: &dbCbsd.LastSeen},
				db.IntType{X: &dbCbsd.GrantAttempts},
				db.IntType{X: &dbCbsd.PreferredBandwidthMHz},
				db.StringType{X: &dbCbsd.PreferredFrequenciesMHz},
				db.FloatType{X: &dbCbsd.MinPower},
				db.FloatType{X: &dbCbsd.MaxPower},
				db.FloatType{X: &dbCbsd.AntennaGain},
				db.IntType{X: &dbCbsd.NumberOfPorts},
				db.BoolType{X: &dbCbsd.IsDeleted},
				db.BoolType{X: &dbCbsd.ShouldDeregister},
			},
		},
		{
			name:  "check field names for DBActiveModeConfig",
			model: dbActiveModeConfig,
			expected: []db.BaseType{
				db.IntType{X: &dbActiveModeConfig.Id},
				db.IntType{X: &dbActiveModeConfig.CbsdId},
				db.IntType{X: &dbActiveModeConfig.DesiredStateId},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.model.Fields())
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
				Properties: []*db.Field{
					{
						Name:    "id",
						SqlType: sqorc.ColumnTypeInt,
					},
					{
						Name:    "name",
						SqlType: sqorc.ColumnTypeText,
					},
				},
			},
		},
		{
			name:  "check ModelMetadata structure for DBGrant",
			model: &storage.DBGrant{},
			expected: db.ModelMetadata{
				Table: storage.GrantTable,
				Properties: []*db.Field{
					{
						Name:    "id",
						SqlType: sqorc.ColumnTypeInt,
					},
					{
						Name:     "state_id",
						SqlType:  sqorc.ColumnTypeInt,
						Relation: storage.GrantStateTable,
					},
					{
						Name:     "cbsd_id",
						SqlType:  sqorc.ColumnTypeInt,
						Nullable: true,
						Relation: storage.CbsdTable,
					},
					{
						Name:    "grant_id",
						SqlType: sqorc.ColumnTypeText,
					},
					{
						Name:     "grant_expire_time",
						SqlType:  sqorc.ColumnTypeDatetime,
						Nullable: true,
					},
					{
						Name:     "transmit_expire_time",
						SqlType:  sqorc.ColumnTypeDatetime,
						Nullable: true,
					},
					{
						Name:     "heartbeat_interval",
						SqlType:  sqorc.ColumnTypeInt,
						Nullable: true,
					},
					{
						Name:     "channel_type",
						SqlType:  sqorc.ColumnTypeText,
						Nullable: true,
					},
					{
						Name:    "low_frequency",
						SqlType: sqorc.ColumnTypeInt,
					},
					{
						Name:    "high_frequency",
						SqlType: sqorc.ColumnTypeInt,
					},
					{
						Name:    "max_eirp",
						SqlType: sqorc.ColumnTypeReal,
					},
				},
			},
		},
		{
			name:  "check ModelMetadata structure for DBCbsdState",
			model: &storage.DBCbsdState{},
			expected: db.ModelMetadata{
				Table: storage.CbsdStateTable,
				Properties: []*db.Field{
					{
						Name:    "id",
						SqlType: sqorc.ColumnTypeInt,
					},
					{
						Name:    "name",
						SqlType: sqorc.ColumnTypeText,
					},
				},
			},
		},
		{
			name:  "check ModelMetadata structure for DBCbsd",
			model: &storage.DBCbsd{},
			expected: db.ModelMetadata{
				Table: storage.CbsdTable,
				Properties: []*db.Field{
					{
						Name:    "id",
						SqlType: sqorc.ColumnTypeInt,
					},
					{
						Name:    "network_id",
						SqlType: sqorc.ColumnTypeText,
					},
					{
						Name:     "state_id",
						SqlType:  sqorc.ColumnTypeInt,
						Relation: storage.CbsdStateTable,
					},
					{
						Name:     "cbsd_id",
						SqlType:  sqorc.ColumnTypeText,
						Nullable: true,
					},
					{
						Name:     "user_id",
						SqlType:  sqorc.ColumnTypeText,
						Nullable: true,
					},
					{
						Name:     "fcc_id",
						SqlType:  sqorc.ColumnTypeText,
						Nullable: true,
					},
					{
						Name:     "cbsd_serial_number",
						SqlType:  sqorc.ColumnTypeText,
						Nullable: true,
						Unique:   true,
					},
					{
						Name:     "last_seen",
						SqlType:  sqorc.ColumnTypeDatetime,
						Nullable: true,
					},
					{
						Name:         "grant_attempts",
						SqlType:      sqorc.ColumnTypeInt,
						HasDefault:   true,
						DefaultValue: 0,
					},
					{
						Name:    "preferred_bandwidth_mhz",
						SqlType: sqorc.ColumnTypeInt,
					},
					{
						Name:    "preferred_frequencies_mhz",
						SqlType: sqorc.ColumnTypeText,
					},
					{
						Name:     "min_power",
						SqlType:  sqorc.ColumnTypeReal,
						Nullable: true,
					},
					{
						Name:     "max_power",
						SqlType:  sqorc.ColumnTypeReal,
						Nullable: true,
					},
					{
						Name:     "antenna_gain",
						SqlType:  sqorc.ColumnTypeReal,
						Nullable: true,
					},
					{
						Name:     "number_of_ports",
						SqlType:  sqorc.ColumnTypeInt,
						Nullable: true,
					},
					{
						Name:         "is_deleted",
						SqlType:      sqorc.ColumnTypeBool,
						HasDefault:   true,
						DefaultValue: false,
					},
					{
						Name:         "should_deregister",
						SqlType:      sqorc.ColumnTypeBool,
						HasDefault:   true,
						DefaultValue: false,
					},
				},
			},
		},
		{
			name:  "check ModelMetadata structure for DBActiveModeConfig",
			model: &storage.DBActiveModeConfig{},
			expected: db.ModelMetadata{
				Table: storage.ActiveModeConfigTable,
				Properties: []*db.Field{
					{
						Name:    "id",
						SqlType: sqorc.ColumnTypeInt,
					},
					{
						Name:     "cbsd_id",
						SqlType:  sqorc.ColumnTypeInt,
						Relation: storage.CbsdTable,
					},
					{
						Name:     "desired_state_id",
						SqlType:  sqorc.ColumnTypeInt,
						Relation: storage.CbsdStateTable,
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.model.GetMetadata()
			obj := actual.CreateObject()

			assert.Equal(t, tc.expected.Properties, actual.Properties)
			assert.Equal(t, tc.expected.Table, actual.Table)
			assert.Equal(t, tc.model, obj)
		})
	}
}
