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
	RequestTypeTable      = "request_types"
	RequestStateTable     = "request_states"
	RequestTable          = "requests"
	ResponseTable         = "responses"
	GrantStateTable       = "grant_states"
	GrantTable            = "grants"
	CbsdStateTable        = "cbsd_states"
	CbsdTable             = "cbsds"
	ChannelTable          = "channels"
	ActiveModeConfigTable = "active_mode_configs"
	LogTable              = "domain_proxy_logs"
)

type DBRequestType struct {
	Id   sql.NullInt64
	Name sql.NullString
}

func (rt *DBRequestType) Fields() db.FieldMap {
	return db.FieldMap{
		"id": &db.Field{
			BaseType: db.IntType{X: &rt.Id},
		},
		"name": &db.Field{
			BaseType: db.StringType{X: &rt.Name},
		},
	}
}

func (rt *DBRequestType) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table:     RequestTypeTable,
		Relations: map[string]string{},
		CreateObject: func() db.Model {
			return &DBRequestType{}
		},
	}
}

type DBRequestState struct {
	Id   sql.NullInt64
	Name sql.NullString
}

func (rs *DBRequestState) Fields() db.FieldMap {
	return db.FieldMap{
		"id": &db.Field{
			BaseType: db.IntType{X: &rs.Id},
		},
		"name": &db.Field{
			BaseType: db.StringType{X: &rs.Name},
		},
	}
}

func (rs *DBRequestState) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table:     RequestStateTable,
		Relations: map[string]string{},
		CreateObject: func() db.Model {
			return &DBRequestState{}
		},
	}
}

type DBRequest struct {
	Id      sql.NullInt64
	TypeId  sql.NullInt64
	StateId sql.NullInt64
	CbsdId  sql.NullInt64
	Payload sql.NullString
}

func (r *DBRequest) Fields() db.FieldMap {
	return db.FieldMap{
		"id": &db.Field{
			BaseType: db.IntType{X: &r.Id},
		},
		"type_id": &db.Field{
			BaseType: db.IntType{X: &r.TypeId},
			Nullable: true,
		},
		"state_id": &db.Field{
			BaseType: db.IntType{X: &r.StateId},
			Nullable: true,
		},
		"cbsd_id": &db.Field{
			BaseType: db.IntType{X: &r.CbsdId},
			Nullable: true,
		},
		"payload": &db.Field{
			BaseType: db.StringType{X: &r.Payload},
			Nullable: true,
		},
	}
}

func (r *DBRequest) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: RequestTable,
		Relations: map[string]string{
			RequestTypeTable:  "type_id",
			RequestStateTable: "state_id",
			CbsdTable:         "cbsd_id",
		},
		CreateObject: func() db.Model {
			return &DBRequest{}
		},
	}
}

type DBResponse struct {
	Id           sql.NullInt64
	RequestId    sql.NullInt64
	GrantId      sql.NullInt64
	ResponseCode sql.NullInt64
	Payload      sql.NullString
}

func (r *DBResponse) Fields() db.FieldMap {
	return db.FieldMap{
		"id": &db.Field{
			BaseType: db.IntType{X: &r.Id},
		},
		"request_id": &db.Field{
			BaseType: db.IntType{X: &r.RequestId},
			Nullable: true,
		},
		"grant_id": &db.Field{
			BaseType: db.IntType{X: &r.GrantId},
			Nullable: true,
		},
		"response_code": &db.Field{
			BaseType: db.IntType{X: &r.Id},
		},
		"payload": &db.Field{
			BaseType: db.StringType{X: &r.Payload},
			Nullable: true,
		},
	}
}

func (r *DBResponse) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: ResponseTable,
		Relations: map[string]string{
			RequestTable: "request_id",
			GrantTable:   "grant_id",
		},
		CreateObject: func() db.Model {
			return &DBResponse{}
		},
	}
}

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
	Id               sql.NullInt64
	NetworkId        sql.NullString
	StateId          sql.NullInt64
	CbsdId           sql.NullString
	UserId           sql.NullString
	FccId            sql.NullString
	CbsdSerialNumber sql.NullString
	LastSeen         sql.NullTime
	MinPower         sql.NullFloat64
	MaxPower         sql.NullFloat64
	AntennaGain      sql.NullFloat64
	NumberOfPorts    sql.NullInt64
	IsDeleted        sql.NullBool
	IsUpdated        sql.NullBool
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
		},
		"last_seen": &db.Field{
			BaseType: db.TimeType{X: &c.LastSeen},
			Nullable: true,
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

type DBChannel struct {
	Id              sql.NullInt64
	CbsdId          sql.NullInt64
	LowFrequency    sql.NullInt64
	HighFrequency   sql.NullInt64
	ChannelType     sql.NullString
	RuleApplied     sql.NullString
	MaxEirp         sql.NullFloat64
	LastUsedMaxEirp sql.NullFloat64
}

func (c *DBChannel) Fields() db.FieldMap {
	return db.FieldMap{
		"id": &db.Field{
			BaseType: db.IntType{X: &c.Id},
		},
		"cbsd_id": &db.Field{
			BaseType: db.IntType{X: &c.CbsdId},
			Nullable: true,
		},
		"low_frequency": &db.Field{
			BaseType: db.IntType{X: &c.LowFrequency},
		},
		"high_frequency": &db.Field{
			BaseType: db.IntType{X: &c.HighFrequency},
		},
		"channel_type": &db.Field{
			BaseType: db.StringType{X: &c.ChannelType},
		},
		"rule_applied": &db.Field{
			BaseType: db.StringType{X: &c.RuleApplied},
		},
		"max_eirp": &db.Field{
			BaseType: db.FloatType{X: &c.MaxEirp},
			Nullable: true,
		},
	}
}

func (c *DBChannel) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: ChannelTable,
		Relations: map[string]string{
			CbsdTable: "cbsd_id",
		},
		CreateObject: func() db.Model {
			return &DBChannel{}
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

type DBLog struct {
	Id           sql.NullInt64
	NetworkId    sql.NullString
	From         sql.NullString
	To           sql.NullString
	Name         sql.NullString
	Message      sql.NullString
	SerialNumber sql.NullString
	FccId        sql.NullString
	ResponseCode sql.NullInt64
	CreatedDate  sql.NullTime
}

func (l *DBLog) Fields() db.FieldMap {
	return db.FieldMap{
		"id": &db.Field{
			BaseType: db.IntType{X: &l.Id},
		},
		"network_id": &db.Field{
			BaseType: db.StringType{X: &l.NetworkId},
		},
		"log_from": &db.Field{
			BaseType: db.StringType{X: &l.From},
		},
		"log_to": &db.Field{
			BaseType: db.StringType{X: &l.To},
		},
		"log_name": &db.Field{
			BaseType: db.StringType{X: &l.Name},
		},
		"log_message": &db.Field{
			BaseType: db.StringType{X: &l.Message},
		},
		"cbsd_serial_number": &db.Field{
			BaseType: db.StringType{X: &l.SerialNumber},
		},
		"fcc_id": &db.Field{
			BaseType: db.StringType{X: &l.FccId},
		},
		"response_code": &db.Field{
			BaseType: db.IntType{X: &l.ResponseCode},
			Nullable: true,
		},
		"created_date": &db.Field{
			BaseType: db.TimeType{X: &l.CreatedDate},
		},
	}
}

func (l *DBLog) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: LogTable,
		CreateObject: func() db.Model {
			return &DBLog{}
		},
	}
}
