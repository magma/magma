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
	"magma/orc8r/cloud/go/sqorc"
)

const (
	GrantStateTable       = "grant_states"
	GrantTable            = "grants"
	CbsdStateTable        = "cbsd_states"
	CbsdTable             = "cbsds"
	ActiveModeConfigTable = "active_mode_configs"
)

type EnumModel interface {
	GetId() int64
	GetName() string
}

type DBGrantState struct {
	Id   sql.NullInt64
	Name sql.NullString
}

func (gs *DBGrantState) Fields() db.FieldMap {
	return db.FieldMap{
		"id": &db.Field{
			Item:    db.IntType{X: &gs.Id},
			SqlType: sqorc.ColumnTypeInt,
		},
		"name": &db.Field{
			Item:    db.StringType{X: &gs.Name},
			SqlType: sqorc.ColumnTypeText,
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

func (gs *DBGrantState) GetId() int64 {
	return gs.Id.Int64
}

func (gs *DBGrantState) GetName() string {
	return gs.Name.String
}

type DBGrant struct {
	Id                 sql.NullInt64
	StateId            sql.NullInt64
	CbsdId             sql.NullInt64
	GrantId            sql.NullString
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
			Item:    db.IntType{X: &g.Id},
			SqlType: sqorc.ColumnTypeInt,
		},
		"state_id": &db.Field{
			Item:    db.IntType{X: &g.StateId},
			SqlType: sqorc.ColumnTypeInt,
		},
		"cbsd_id": &db.Field{
			Item:     db.IntType{X: &g.CbsdId},
			SqlType:  sqorc.ColumnTypeInt,
			Nullable: true,
		},
		"grant_id": &db.Field{
			Item:    db.StringType{X: &g.GrantId},
			SqlType: sqorc.ColumnTypeText,
		},
		"grant_expire_time": &db.Field{
			Item:     db.TimeType{X: &g.GrantExpireTime},
			SqlType:  sqorc.ColumnTypeDatetime,
			Nullable: true,
		},
		"transmit_expire_time": &db.Field{
			Item:     db.TimeType{X: &g.TransmitExpireTime},
			SqlType:  sqorc.ColumnTypeDatetime,
			Nullable: true,
		},
		"heartbeat_interval": &db.Field{
			Item:     db.IntType{X: &g.HeartbeatInterval},
			SqlType:  sqorc.ColumnTypeInt,
			Nullable: true,
		},
		"channel_type": &db.Field{
			Item:     db.StringType{X: &g.ChannelType},
			SqlType:  sqorc.ColumnTypeText,
			Nullable: true,
		},
		"low_frequency": &db.Field{
			Item:    db.IntType{X: &g.LowFrequency},
			SqlType: sqorc.ColumnTypeInt,
		},
		"high_frequency": &db.Field{
			Item:    db.IntType{X: &g.HighFrequency},
			SqlType: sqorc.ColumnTypeInt,
		},
		"max_eirp": &db.Field{
			Item:    db.FloatType{X: &g.MaxEirp},
			SqlType: sqorc.ColumnTypeReal,
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
			Item:    db.IntType{X: &cs.Id},
			SqlType: sqorc.ColumnTypeInt,
		},
		"name": &db.Field{
			Item:    db.StringType{X: &cs.Name},
			SqlType: sqorc.ColumnTypeText,
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

func (cs *DBCbsdState) GetId() int64 {
	return cs.Id.Int64
}

func (cs *DBCbsdState) GetName() string {
	return cs.Name.String
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
	ShouldDeregister        sql.NullBool
}

func (c *DBCbsd) Fields() db.FieldMap {
	return db.FieldMap{
		"id": &db.Field{
			Item:    db.IntType{X: &c.Id},
			SqlType: sqorc.ColumnTypeInt,
		},
		"network_id": &db.Field{
			Item:    db.StringType{X: &c.NetworkId},
			SqlType: sqorc.ColumnTypeText,
		},
		"state_id": &db.Field{
			Item:    db.IntType{X: &c.StateId},
			SqlType: sqorc.ColumnTypeInt,
		},
		"cbsd_id": &db.Field{
			Item:     db.StringType{X: &c.CbsdId},
			SqlType:  sqorc.ColumnTypeText,
			Nullable: true,
		},
		"user_id": &db.Field{
			Item:     db.StringType{X: &c.UserId},
			SqlType:  sqorc.ColumnTypeText,
			Nullable: true,
		},
		"fcc_id": &db.Field{
			Item:     db.StringType{X: &c.FccId},
			SqlType:  sqorc.ColumnTypeText,
			Nullable: true,
		},
		"cbsd_serial_number": &db.Field{
			Item:     db.StringType{X: &c.CbsdSerialNumber},
			SqlType:  sqorc.ColumnTypeText,
			Nullable: true,
			Unique:   true,
		},
		"last_seen": &db.Field{
			Item:     db.TimeType{X: &c.LastSeen},
			SqlType:  sqorc.ColumnTypeDatetime,
			Nullable: true,
		},
		"grant_attempts": &db.Field{
			Item:         db.IntType{X: &c.GrantAttempts},
			SqlType:      sqorc.ColumnTypeInt,
			HasDefault:   true,
			DefaultValue: 0,
		},
		"preferred_bandwidth_mhz": &db.Field{
			Item:    db.IntType{X: &c.PreferredBandwidthMHz},
			SqlType: sqorc.ColumnTypeInt,
		},
		"preferred_frequencies_mhz": &db.Field{
			Item:    db.StringType{X: &c.PreferredFrequenciesMHz},
			SqlType: sqorc.ColumnTypeText,
		},
		"min_power": &db.Field{
			Item:     db.FloatType{X: &c.MinPower},
			SqlType:  sqorc.ColumnTypeReal,
			Nullable: true,
		},
		"max_power": &db.Field{
			Item:     db.FloatType{X: &c.MaxPower},
			SqlType:  sqorc.ColumnTypeReal,
			Nullable: true,
		},
		"antenna_gain": &db.Field{
			Item:     db.FloatType{X: &c.AntennaGain},
			SqlType:  sqorc.ColumnTypeReal,
			Nullable: true,
		},
		"number_of_ports": &db.Field{
			Item:     db.IntType{X: &c.NumberOfPorts},
			SqlType:  sqorc.ColumnTypeInt,
			Nullable: true,
		},
		"is_deleted": &db.Field{
			Item:         db.BoolType{X: &c.IsDeleted},
			SqlType:      sqorc.ColumnTypeBool,
			HasDefault:   true,
			DefaultValue: false,
		},
		"should_deregister": &db.Field{
			Item:         db.BoolType{X: &c.ShouldDeregister},
			SqlType:      sqorc.ColumnTypeBool,
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
			Item:    db.IntType{X: &amc.Id},
			SqlType: sqorc.ColumnTypeInt,
		},
		"cbsd_id": &db.Field{
			Item:    db.IntType{X: &amc.CbsdId},
			SqlType: sqorc.ColumnTypeInt,
		},
		"desired_state_id": &db.Field{
			Item:    db.IntType{X: &amc.DesiredStateId},
			SqlType: sqorc.ColumnTypeInt,
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
