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

func (gs *DBGrantState) Fields() map[string]db.BaseType {
	return map[string]db.BaseType{
		"id":   db.IntType{X: &gs.Id},
		"name": db.StringType{X: &gs.Name},
	}
}

func (gs *DBGrantState) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: GrantStateTable,
		Properties: map[string]*db.Field{
			"id": {
				SqlType: sqorc.ColumnTypeInt,
			},
			"name": {
				SqlType: sqorc.ColumnTypeText,
			},
		},
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

func (g *DBGrant) Fields() map[string]db.BaseType {
	return map[string]db.BaseType{
		"id":                   db.IntType{X: &g.Id},
		"state_id":             db.IntType{X: &g.StateId},
		"cbsd_id":              db.IntType{X: &g.CbsdId},
		"grant_id":             db.StringType{X: &g.GrantId},
		"grant_expire_time":    db.TimeType{X: &g.GrantExpireTime},
		"transmit_expire_time": db.TimeType{X: &g.TransmitExpireTime},
		"heartbeat_interval":   db.IntType{X: &g.HeartbeatInterval},
		"channel_type":         db.StringType{X: &g.ChannelType},
		"low_frequency":        db.IntType{X: &g.LowFrequency},
		"high_frequency":       db.IntType{X: &g.HighFrequency},
		"max_eirp":             db.FloatType{X: &g.MaxEirp},
	}
}

func (g *DBGrant) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: GrantTable,
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

func (cs *DBCbsdState) Fields() map[string]db.BaseType {
	return map[string]db.BaseType{
		"id":   db.IntType{X: &cs.Id},
		"name": db.StringType{X: &cs.Name},
	}
}

func (cs *DBCbsdState) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: "cbsd_states",
		Properties: map[string]*db.Field{
			"id": {
				SqlType: sqorc.ColumnTypeInt,
			},
			"name": {
				SqlType: sqorc.ColumnTypeText,
			},
		},
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

func (c *DBCbsd) Fields() map[string]db.BaseType {
	return map[string]db.BaseType{
		"id":                        db.IntType{X: &c.Id},
		"network_id":                db.StringType{X: &c.NetworkId},
		"state_id":                  db.IntType{X: &c.StateId},
		"cbsd_id":                   db.StringType{X: &c.CbsdId},
		"user_id":                   db.StringType{X: &c.UserId},
		"fcc_id":                    db.StringType{X: &c.FccId},
		"cbsd_serial_number":        db.StringType{X: &c.CbsdSerialNumber},
		"last_seen":                 db.TimeType{X: &c.LastSeen},
		"grant_attempts":            db.IntType{X: &c.GrantAttempts},
		"preferred_bandwidth_mhz":   db.IntType{X: &c.PreferredBandwidthMHz},
		"preferred_frequencies_mhz": db.StringType{X: &c.PreferredFrequenciesMHz},
		"min_power":                 db.FloatType{X: &c.MinPower},
		"max_power":                 db.FloatType{X: &c.MaxPower},
		"antenna_gain":              db.FloatType{X: &c.AntennaGain},
		"number_of_ports":           db.IntType{X: &c.NumberOfPorts},
		"is_deleted":                db.BoolType{X: &c.IsDeleted},
		"should_deregister":         db.BoolType{X: &c.ShouldDeregister},
	}
}

func (c *DBCbsd) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: CbsdTable,
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

func (amc *DBActiveModeConfig) Fields() map[string]db.BaseType {
	return map[string]db.BaseType{
		"id":               db.IntType{X: &amc.Id},
		"cbsd_id":          db.IntType{X: &amc.CbsdId},
		"desired_state_id": db.IntType{X: &amc.DesiredStateId},
	}
}

func (amc *DBActiveModeConfig) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: ActiveModeConfigTable,
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
			CbsdTable:      "cbsd_id",
			CbsdStateTable: "desired_state_id",
		},
		CreateObject: func() db.Model {
			return &DBActiveModeConfig{}
		},
	}
}
