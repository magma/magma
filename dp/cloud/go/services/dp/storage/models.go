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

func (rt *DBRequestType) Fields() map[string]db.Field {
	return map[string]db.Field{
		"id":   db.IntField{X: &rt.Id},
		"name": db.StringField{X: &rt.Name},
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

func (rs *DBRequestState) Fields() map[string]db.Field {
	return map[string]db.Field{
		"id":   db.IntField{X: &rs.Id},
		"name": db.StringField{X: &rs.Name},
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

func (r *DBRequest) Fields() map[string]db.Field {
	return map[string]db.Field{
		"id":       db.IntField{X: &r.Id},
		"type_id":  db.NullIntField{X: &r.TypeId},
		"state_id": db.NullIntField{X: &r.StateId},
		"cbsd_id":  db.NullIntField{X: &r.CbsdId},
		"payload":  db.NullStringField{X: &r.Payload},
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

func (r *DBResponse) Fields() map[string]db.Field {
	return map[string]db.Field{
		"id":            db.IntField{X: &r.Id},
		"request_id":    db.NullIntField{X: &r.RequestId},
		"grant_id":      db.NullIntField{X: &r.GrantId},
		"response_code": db.IntField{X: &r.Id},
		"payload":       db.NullStringField{X: &r.Payload},
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

func (gs *DBGrantState) Fields() map[string]db.Field {
	return map[string]db.Field{
		"id":   db.IntField{X: &gs.Id},
		"name": db.StringField{X: &gs.Name},
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
	ChannelId          sql.NullInt64
	GrantId            sql.NullInt64
	GrantExpireTime    sql.NullTime
	TransmitExpireTime sql.NullTime
	HeartbeatInterval  sql.NullInt64
	ChannelType        sql.NullString
}

func (g *DBGrant) Fields() map[string]db.Field {
	return map[string]db.Field{
		"id":                   db.IntField{X: &g.Id},
		"state_id":             db.IntField{X: &g.StateId},
		"cbsd_id":              db.NullIntField{X: &g.CbsdId},
		"channel_id":           db.NullIntField{X: &g.ChannelId},
		"grant_id":             db.IntField{X: &g.GrantId},
		"grant_expire_time":    db.NullTimeField{X: &g.GrantExpireTime},
		"transmit_expire_time": db.NullTimeField{X: &g.TransmitExpireTime},
		"heartbeat_interval":   db.NullIntField{X: &g.HeartbeatInterval},
		"channel_type":         db.NullStringField{X: &g.ChannelType},
	}
}

func (g *DBGrant) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: GrantTable,
		Relations: map[string]string{
			GrantStateTable: "state_id",
			CbsdTable:       "cbsd_id",
			ChannelTable:    "channel_id",
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

func (cs *DBCbsdState) Fields() map[string]db.Field {
	return map[string]db.Field{
		"id":   db.IntField{X: &cs.Id},
		"name": db.StringField{X: &cs.Name},
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
}

func (c *DBCbsd) Fields() map[string]db.Field {
	return map[string]db.Field{
		"id":                 db.IntField{X: &c.Id},
		"network_id":         db.StringField{X: &c.NetworkId},
		"state_id":           db.IntField{X: &c.StateId},
		"cbsd_id":            db.NullStringField{X: &c.CbsdId},
		"user_id":            db.NullStringField{X: &c.UserId},
		"fcc_id":             db.NullStringField{X: &c.FccId},
		"cbsd_serial_number": db.NullStringField{X: &c.CbsdSerialNumber},
		"last_seen":          db.NullTimeField{X: &c.LastSeen},
		"min_power":          db.NullFloatField{X: &c.MinPower},
		"max_power":          db.NullFloatField{X: &c.MaxPower},
		"antenna_gain":       db.NullFloatField{X: &c.AntennaGain},
		"number_of_ports":    db.NullIntField{X: &c.NumberOfPorts},
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

func (c *DBChannel) Fields() map[string]db.Field {
	return map[string]db.Field{
		"id":                 db.IntField{X: &c.Id},
		"cbsd_id":            db.NullIntField{X: &c.CbsdId},
		"low_frequency":      db.IntField{X: &c.LowFrequency},
		"high_frequency":     db.IntField{X: &c.HighFrequency},
		"channel_type":       db.StringField{X: &c.ChannelType},
		"rule_applied":       db.StringField{X: &c.RuleApplied},
		"max_eirp":           db.NullFloatField{X: &c.MaxEirp},
		"last_used_max_eirp": db.NullFloatField{X: &c.LastUsedMaxEirp},
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

func (amc *DBActiveModeConfig) Fields() map[string]db.Field {
	return map[string]db.Field{
		"id":               db.IntField{X: &amc.Id},
		"cbsd_id":          db.IntField{X: &amc.CbsdId},
		"desired_state_id": db.IntField{X: &amc.DesiredStateId},
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

func (l *DBLog) Fields() map[string]db.Field {
	return map[string]db.Field{
		"id":                 db.IntField{X: &l.Id},
		"network_id":         db.StringField{X: &l.NetworkId},
		"log_from":           db.StringField{X: &l.From},
		"log_to":             db.StringField{X: &l.To},
		"log_name":           db.StringField{X: &l.Name},
		"log_message":        db.StringField{X: &l.Message},
		"cbsd_serial_number": db.StringField{X: &l.SerialNumber},
		"fcc_id":             db.StringField{X: &l.FccId},
		"response_code":      db.NullIntField{X: &l.ResponseCode},
		"created_date":       db.TimeField{X: &l.CreatedDate},
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
