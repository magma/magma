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
	RequestTypeTable = "request_types"
	RequestTable     = "requests"
	GrantStateTable  = "grant_states"
	GrantTable       = "grants"
	CbsdStateTable   = "cbsd_states"
	CbsdTable        = "cbsds"
)

type EnumModel interface {
	GetId() int64
	GetName() string
}

type DBRequestType struct {
	Id   sql.NullInt64
	Name sql.NullString
}

func (rt *DBRequestType) Fields() []db.BaseType {
	return []db.BaseType{
		db.IntType{X: &rt.Id},
		db.StringType{X: &rt.Name},
	}
}

func (rt *DBRequestType) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: RequestTypeTable,
		Properties: []*db.Field{{
			Name:    "id",
			SqlType: sqorc.ColumnTypeInt,
		}, {
			Name:    "name",
			SqlType: sqorc.ColumnTypeText,
		}},
		CreateObject: func() db.Model {
			return &DBRequestType{}
		},
	}
}

func (rt *DBRequestType) GetId() int64 {
	return rt.Id.Int64
}

func (rt *DBRequestType) GetName() string {
	return rt.Name.String
}

type DBRequest struct {
	Id      sql.NullInt64
	TypeId  sql.NullInt64
	CbsdId  sql.NullInt64
	Payload any
}

func (r *DBRequest) Fields() []db.BaseType {
	return []db.BaseType{
		db.IntType{X: &r.Id},
		db.IntType{X: &r.TypeId},
		db.IntType{X: &r.CbsdId},
		db.JsonType{X: &r.Payload},
	}
}

func (r *DBRequest) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: RequestTable,
		Properties: []*db.Field{{
			Name:    "id",
			SqlType: sqorc.ColumnTypeInt,
		}, {
			Name:     "type_id",
			SqlType:  sqorc.ColumnTypeInt,
			Relation: RequestTypeTable,
		}, {
			Name:     "cbsd_id",
			SqlType:  sqorc.ColumnTypeInt,
			Relation: CbsdTable,
		}, {
			Name:    "payload",
			SqlType: sqorc.ColumnTypeText,
		}},
		CreateObject: func() db.Model {
			return &DBRequest{}
		},
	}
}

type DBGrantState struct {
	Id   sql.NullInt64
	Name sql.NullString
}

func (gs *DBGrantState) Fields() []db.BaseType {
	return []db.BaseType{
		db.IntType{X: &gs.Id},
		db.StringType{X: &gs.Name},
	}
}

func (gs *DBGrantState) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: GrantStateTable,
		Properties: []*db.Field{{
			Name:    "id",
			SqlType: sqorc.ColumnTypeInt,
		}, {
			Name:    "name",
			SqlType: sqorc.ColumnTypeText,
		}},
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
	Id                       sql.NullInt64
	StateId                  sql.NullInt64
	CbsdId                   sql.NullInt64
	GrantId                  sql.NullString
	GrantExpireTime          sql.NullTime
	TransmitExpireTime       sql.NullTime
	HeartbeatIntervalSec     sql.NullInt64
	LastHeartbeatRequestTime sql.NullTime
	ChannelType              sql.NullString
	LowFrequencyHz           sql.NullInt64
	HighFrequencyHz          sql.NullInt64
	MaxEirp                  sql.NullFloat64
}

func (g *DBGrant) Fields() []db.BaseType {
	return []db.BaseType{
		db.IntType{X: &g.Id},
		db.IntType{X: &g.StateId},
		db.IntType{X: &g.CbsdId},
		db.StringType{X: &g.GrantId},
		db.TimeType{X: &g.GrantExpireTime},
		db.TimeType{X: &g.TransmitExpireTime},
		db.IntType{X: &g.HeartbeatIntervalSec},
		db.TimeType{X: &g.LastHeartbeatRequestTime},
		db.StringType{X: &g.ChannelType},
		db.IntType{X: &g.LowFrequencyHz},
		db.IntType{X: &g.HighFrequencyHz},
		db.FloatType{X: &g.MaxEirp},
	}
}

func (g *DBGrant) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: GrantTable,
		Properties: []*db.Field{{
			Name:    "id",
			SqlType: sqorc.ColumnTypeInt,
		}, {
			Name:     "state_id",
			SqlType:  sqorc.ColumnTypeInt,
			Relation: GrantStateTable,
		}, {
			Name:     "cbsd_id",
			SqlType:  sqorc.ColumnTypeInt,
			Nullable: true,
			Relation: CbsdTable,
		}, {
			Name:    "grant_id",
			SqlType: sqorc.ColumnTypeText,
		}, {
			Name:     "grant_expire_time",
			SqlType:  sqorc.ColumnTypeDatetime,
			Nullable: true,
		}, {
			Name:     "transmit_expire_time",
			SqlType:  sqorc.ColumnTypeDatetime,
			Nullable: true,
		}, {
			Name:     "heartbeat_interval",
			SqlType:  sqorc.ColumnTypeInt,
			Nullable: true,
		}, {
			Name:     "last_heartbeat_request_time",
			SqlType:  sqorc.ColumnTypeDatetime,
			Nullable: true,
		}, {
			Name:     "channel_type",
			SqlType:  sqorc.ColumnTypeText,
			Nullable: true,
		}, {
			Name:    "low_frequency",
			SqlType: sqorc.ColumnTypeInt,
		}, {
			Name:    "high_frequency",
			SqlType: sqorc.ColumnTypeInt,
		}, {
			Name:    "max_eirp",
			SqlType: sqorc.ColumnTypeReal,
		}},
		CreateObject: func() db.Model {
			return &DBGrant{}
		},
	}
}

type DBCbsdState struct {
	Id   sql.NullInt64
	Name sql.NullString
}

func (cs *DBCbsdState) Fields() []db.BaseType {
	return []db.BaseType{
		db.IntType{X: &cs.Id},
		db.StringType{X: &cs.Name},
	}
}

func (cs *DBCbsdState) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: "cbsd_states",
		Properties: []*db.Field{{
			Name:    "id",
			SqlType: sqorc.ColumnTypeInt,
		}, {
			Name:    "name",
			SqlType: sqorc.ColumnTypeText,
		}},
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
	Id                        sql.NullInt64
	NetworkId                 sql.NullString
	StateId                   sql.NullInt64
	DesiredStateId            sql.NullInt64
	CbsdId                    sql.NullString
	UserId                    sql.NullString
	FccId                     sql.NullString
	CbsdSerialNumber          sql.NullString
	LastSeen                  sql.NullTime
	PreferredBandwidthMHz     sql.NullInt64
	PreferredFrequenciesMHz   []int64
	MinPower                  sql.NullFloat64
	MaxPower                  sql.NullFloat64
	AntennaGainDbi            sql.NullFloat64
	NumberOfPorts             sql.NullInt64
	IsDeleted                 sql.NullBool
	ShouldDeregister          sql.NullBool
	ShouldRelinquish          sql.NullBool
	SingleStepEnabled         sql.NullBool
	CbsdCategory              sql.NullString
	LatitudeDeg               sql.NullFloat64
	LongitudeDeg              sql.NullFloat64
	HeightM                   sql.NullFloat64
	HeightType                sql.NullString
	IndoorDeployment          sql.NullBool
	CarrierAggregationEnabled sql.NullBool
	GrantRedundancy           sql.NullBool
	MaxIbwMhx                 sql.NullInt64
	AvailableFrequencies      []uint32
	Channels                  []Channel
}

type Channel struct {
	// TODO some of the fields may not be required
	LowFrequencyHz  int64   `json:"low_frequency"`
	HighFrequencyHz int64   `json:"high_frequency"`
	MaxEirp         float64 `json:"max_eirp"`
}

func (c *DBCbsd) Fields() []db.BaseType {
	return []db.BaseType{
		db.IntType{X: &c.Id},
		db.StringType{X: &c.NetworkId},
		db.IntType{X: &c.StateId},
		db.IntType{X: &c.DesiredStateId},
		db.StringType{X: &c.CbsdId},
		db.StringType{X: &c.UserId},
		db.StringType{X: &c.FccId},
		db.StringType{X: &c.CbsdSerialNumber},
		db.TimeType{X: &c.LastSeen},
		db.IntType{X: &c.PreferredBandwidthMHz},
		db.JsonType{X: &c.PreferredFrequenciesMHz},
		db.FloatType{X: &c.MinPower},
		db.FloatType{X: &c.MaxPower},
		db.FloatType{X: &c.AntennaGainDbi},
		db.IntType{X: &c.NumberOfPorts},
		db.BoolType{X: &c.IsDeleted},
		db.BoolType{X: &c.ShouldDeregister},
		db.BoolType{X: &c.ShouldRelinquish},
		db.BoolType{X: &c.SingleStepEnabled},
		db.StringType{X: &c.CbsdCategory},
		db.FloatType{X: &c.LatitudeDeg},
		db.FloatType{X: &c.LongitudeDeg},
		db.FloatType{X: &c.HeightM},
		db.StringType{X: &c.HeightType},
		db.BoolType{X: &c.IndoorDeployment},
		db.BoolType{X: &c.CarrierAggregationEnabled},
		db.BoolType{X: &c.GrantRedundancy},
		db.IntType{X: &c.MaxIbwMhx},
		db.JsonType{X: &c.AvailableFrequencies},
		db.JsonType{X: &c.Channels},
	}
}

func (c *DBCbsd) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: CbsdTable,
		Properties: []*db.Field{{
			Name:    "id",
			SqlType: sqorc.ColumnTypeInt,
		}, {
			Name:    "network_id",
			SqlType: sqorc.ColumnTypeText,
		}, {
			Name:     "state_id",
			SqlType:  sqorc.ColumnTypeInt,
			Relation: CbsdStateTable,
		}, {
			Name:     "desired_state_id",
			SqlType:  sqorc.ColumnTypeInt,
			Relation: CbsdStateTable,
		}, {
			Name:     "cbsd_id",
			SqlType:  sqorc.ColumnTypeText,
			Nullable: true,
		}, {
			Name:     "user_id",
			SqlType:  sqorc.ColumnTypeText,
			Nullable: true,
		}, {
			Name:     "fcc_id",
			SqlType:  sqorc.ColumnTypeText,
			Nullable: true,
		}, {
			Name:     "cbsd_serial_number",
			SqlType:  sqorc.ColumnTypeText,
			Nullable: true,
			Unique:   true,
		}, {
			Name:     "last_seen",
			SqlType:  sqorc.ColumnTypeDatetime,
			Nullable: true,
		}, {
			Name:    "preferred_bandwidth_mhz",
			SqlType: sqorc.ColumnTypeInt,
		}, {
			Name:    "preferred_frequencies_mhz",
			SqlType: sqorc.ColumnTypeText,
		}, {
			Name:     "min_power",
			SqlType:  sqorc.ColumnTypeReal,
			Nullable: true,
		}, {
			Name:     "max_power",
			SqlType:  sqorc.ColumnTypeReal,
			Nullable: true,
		}, {
			Name:     "antenna_gain",
			SqlType:  sqorc.ColumnTypeReal,
			Nullable: true,
		}, {
			Name:     "number_of_ports",
			SqlType:  sqorc.ColumnTypeInt,
			Nullable: true,
		}, {
			Name:         "is_deleted",
			SqlType:      sqorc.ColumnTypeBool,
			HasDefault:   true,
			DefaultValue: false,
		}, {
			Name:         "should_deregister",
			SqlType:      sqorc.ColumnTypeBool,
			HasDefault:   true,
			DefaultValue: false,
		}, {
			Name:         "should_relinquish",
			SqlType:      sqorc.ColumnTypeBool,
			HasDefault:   true,
			DefaultValue: false,
		}, {
			Name:         "single_step_enabled",
			SqlType:      sqorc.ColumnTypeBool,
			HasDefault:   true,
			DefaultValue: false,
		}, {
			Name:         "cbsd_category",
			SqlType:      sqorc.ColumnTypeText,
			HasDefault:   true,
			DefaultValue: "b",
		}, {
			Name:     "latitude_deg",
			SqlType:  sqorc.ColumnTypeReal,
			Nullable: true,
		}, {
			Name:     "longitude_deg",
			SqlType:  sqorc.ColumnTypeReal,
			Nullable: true,
		}, {
			Name:     "height_m",
			SqlType:  sqorc.ColumnTypeInt,
			Nullable: true,
		}, {
			Name:     "height_type",
			SqlType:  sqorc.ColumnTypeText,
			Nullable: true,
		}, {
			Name:         "indoor_deployment",
			SqlType:      sqorc.ColumnTypeBool,
			HasDefault:   true,
			DefaultValue: false,
		}, {
			Name:         "carrier_aggregation_enabled",
			SqlType:      sqorc.ColumnTypeBool,
			HasDefault:   true,
			DefaultValue: false,
		}, {
			Name:         "grant_redundancy",
			SqlType:      sqorc.ColumnTypeBool,
			HasDefault:   true,
			DefaultValue: true,
		}, {
			Name:         "max_ibw_mhz",
			SqlType:      sqorc.ColumnTypeInt,
			HasDefault:   true,
			DefaultValue: 150,
		}, {
			Name:     "available_frequencies",
			SqlType:  sqorc.ColumnTypeText,
			Nullable: true,
		}, {
			Name:         "channels",
			SqlType:      sqorc.ColumnTypeText,
			Nullable:     false,
			HasDefault:   true,
			DefaultValue: "'[]'",
		}},
		CreateObject: func() db.Model {
			return &DBCbsd{}
		},
	}
}
