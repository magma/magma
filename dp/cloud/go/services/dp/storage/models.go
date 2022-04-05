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

package storage

import (
	"database/sql"

	"magma/dp/cloud/go/services/dp/storage/db"
)

const (
	GrantStateTable       = "grant_states"
	GrantTable            = "grants"
	CbsdStateTable        = "cbsd_states"
	CbsdTable             = "cbsds"
	ActiveModeConfigTable = "active_mode_configs"
)

type DBGrantState struct {
	Id   sql.NullInt64
	Name sql.NullString
}

func (gs *DBGrantState) Fields() db.FieldMap {
	return db.FieldMap{
		"id": &db.Field{
			BaseType: db.IntType{X: &gs.Id},
		},
		"name": &db.Field{
			BaseType: db.StringType{X: &gs.Name},
		},
	}
}

func (gs *DBGrantState) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table:     GrantStateTable,
		Relations: map[string]string{},
		CreateObject: func() db.Model {
			return &DBGrantState{}
		},
	}
}

type DBGrant struct {
	Id                 sql.NullInt64
	StateId            sql.NullInt64
	CbsdId             sql.NullInt64
	GrantId            sql.NullInt64
	GrantExpireTime    sql.NullTime
	TransmitExpireTime sql.NullTime
	HeartbeatInterval  sql.NullInt64
	ChannelType        sql.NullString
	LowFrequency       sql.NullInt64
	HighFrequency      sql.NullInt64
	MaxEirp            sql.NullFloat64
}

func (g *DBGrant) Fields() db.FieldMap {
	return db.FieldMap{
		"id": &db.Field{
			BaseType: db.IntType{X: &g.Id},
		},
		"state_id": &db.Field{
			BaseType: db.IntType{X: &g.StateId},
		},
		"cbsd_id": &db.Field{
			BaseType: db.IntType{X: &g.CbsdId},
			Nullable: true,
		},
		"grant_id": &db.Field{
			BaseType: db.IntType{X: &g.GrantId},
		},
		"grant_expire_time": &db.Field{
			BaseType: db.TimeType{X: &g.GrantExpireTime},
			Nullable: true,
		},
		"transmit_expire_time": &db.Field{
			BaseType: db.TimeType{X: &g.TransmitExpireTime},
			Nullable: true,
		},
		"heartbeat_interval": &db.Field{
			BaseType: db.IntType{X: &g.HeartbeatInterval},
			Nullable: true,
		},
		"channel_type": &db.Field{
			BaseType: db.StringType{X: &g.ChannelType},
			Nullable: true,
		},
		"low_frequency": &db.Field{
			BaseType: db.IntType{X: &g.LowFrequency},
		},
		"high_frequency": &db.Field{
			BaseType: db.IntType{X: &g.HighFrequency},
		},
		"max_eirp": &db.Field{
			BaseType: db.FloatType{X: &g.MaxEirp},
		},
	}
}

func (g *DBGrant) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: GrantTable,
		Relations: map[string]string{
			GrantStateTable: "state_id",
			CbsdTable:       "cbsd_id",
		},
		CreateObject: func() db.Model {
			return &DBGrant{}
		},
	}
}

type DBCbsdState struct {
	Id   sql.NullInt64
	Name sql.NullString
}

func (cs *DBCbsdState) Fields() db.FieldMap {
	return db.FieldMap{
		"id": &db.Field{
			BaseType: db.IntType{X: &cs.Id},
		},
		"name": &db.Field{
			BaseType: db.StringType{X: &cs.Name},
		},
	}
}

func (cs *DBCbsdState) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table:     "cbsd_states",
		Relations: map[string]string{},
		CreateObject: func() db.Model {
			return &DBCbsdState{}
		},
	}
}

type DBCbsd struct {
	Id                      sql.NullInt64
	NetworkId               sql.NullString
	StateId                 sql.NullInt64
	CbsdId                  sql.NullString
	UserId                  sql.NullString
	FccId                   sql.NullString
	CbsdSerialNumber        sql.NullString
	LastSeen                sql.NullTime
	GrantAttempts           sql.NullInt64
	PreferredBandwidthMHz   sql.NullInt64
	PreferredFrequenciesMHz sql.NullString
	MinPower                sql.NullFloat64
	MaxPower                sql.NullFloat64
	AntennaGain             sql.NullFloat64
	NumberOfPorts           sql.NullInt64
	IsDeleted               sql.NullBool
	IsUpdated               sql.NullBool
}

func (c *DBCbsd) Fields() db.FieldMap {
	return db.FieldMap{
		"id": &db.Field{
			BaseType: db.IntType{X: &c.Id},
		},
		"network_id": &db.Field{
			BaseType: db.StringType{X: &c.NetworkId},
		},
		"state_id": &db.Field{
			BaseType: db.IntType{X: &c.StateId},
		},
		"cbsd_id": &db.Field{
			BaseType: db.StringType{X: &c.CbsdId},
			Nullable: true,
		},
		"user_id": &db.Field{
			BaseType: db.StringType{X: &c.UserId},
			Nullable: true,
		},
		"fcc_id": &db.Field{
			BaseType: db.StringType{X: &c.FccId},
			Nullable: true,
		},
		"cbsd_serial_number": &db.Field{
			BaseType: db.StringType{X: &c.CbsdSerialNumber},
			Nullable: true,
			Unique:   true,
		},
		"last_seen": &db.Field{
			BaseType: db.TimeType{X: &c.LastSeen},
			Nullable: true,
		},
		"grant_attempts": &db.Field{
			BaseType:     db.IntType{X: &c.GrantAttempts},
			HasDefault:   true,
			DefaultValue: 0,
		},
		"preferred_bandwidth_mhz": &db.Field{
			BaseType: db.IntType{X: &c.PreferredBandwidthMHz},
		},
		"preferred_frequencies_mhz": &db.Field{
			BaseType: db.StringType{X: &c.PreferredFrequenciesMHz},
		},
		"min_power": &db.Field{
			BaseType: db.FloatType{X: &c.MinPower},
			Nullable: true,
		},
		"max_power": &db.Field{
			BaseType: db.FloatType{X: &c.MaxPower},
			Nullable: true,
		},
		"antenna_gain": &db.Field{
			BaseType: db.FloatType{X: &c.AntennaGain},
			Nullable: true,
		},
		"number_of_ports": &db.Field{
			BaseType: db.IntType{X: &c.NumberOfPorts},
			Nullable: true,
		},
		"is_deleted": &db.Field{
			BaseType:     db.BoolType{X: &c.IsDeleted},
			HasDefault:   true,
			DefaultValue: false,
		},
		"is_updated": &db.Field{
			BaseType:     db.BoolType{X: &c.IsUpdated},
			HasDefault:   true,
			DefaultValue: false,
		},
	}
}

func (c *DBCbsd) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: CbsdTable,
		Relations: map[string]string{
			CbsdStateTable: "state_id",
		},
		CreateObject: func() db.Model {
			return &DBCbsd{}
		},
	}
}

type DBActiveModeConfig struct {
	Id             sql.NullInt64
	CbsdId         sql.NullInt64
	DesiredStateId sql.NullInt64
}

func (amc *DBActiveModeConfig) Fields() db.FieldMap {
	return db.FieldMap{
		"id": &db.Field{
			BaseType: db.IntType{X: &amc.Id},
		},
		"cbsd_id": &db.Field{
			BaseType: db.IntType{X: &amc.CbsdId},
		},
		"desired_state_id": &db.Field{
			BaseType: db.IntType{X: &amc.DesiredStateId},
		},
	}
}

func (amc *DBActiveModeConfig) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: ActiveModeConfigTable,
		Relations: map[string]string{
			CbsdTable:      "cbsd_id",
			CbsdStateTable: "desired_state_id",
		},
		CreateObject: func() db.Model {
			return &DBActiveModeConfig{}
		},
	}
}
