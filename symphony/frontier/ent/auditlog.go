// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/frontier/ent/auditlog"
)

// AuditLog is the model entity for the AuditLog schema.
type AuditLog struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// ActingUserID holds the value of the "acting_user_id" field.
	ActingUserID int `json:"acting_user_id,omitempty"`
	// Organization holds the value of the "organization" field.
	Organization string `json:"organization,omitempty"`
	// MutationType holds the value of the "mutation_type" field.
	MutationType string `json:"mutation_type,omitempty"`
	// ObjectID holds the value of the "object_id" field.
	ObjectID string `json:"object_id,omitempty"`
	// ObjectType holds the value of the "object_type" field.
	ObjectType string `json:"object_type,omitempty"`
	// ObjectDisplayName holds the value of the "object_display_name" field.
	ObjectDisplayName string `json:"object_display_name,omitempty"`
	// MutationData holds the value of the "mutation_data" field.
	MutationData map[string]string `json:"mutation_data,omitempty"`
	// URL holds the value of the "url" field.
	URL string `json:"url,omitempty"`
	// IPAddress holds the value of the "ip_address" field.
	IPAddress string `json:"ip_address,omitempty"`
	// Status holds the value of the "status" field.
	Status string `json:"status,omitempty"`
	// StatusCode holds the value of the "status_code" field.
	StatusCode string `json:"status_code,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*AuditLog) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // created_at
		&sql.NullTime{},   // updated_at
		&sql.NullInt64{},  // acting_user_id
		&sql.NullString{}, // organization
		&sql.NullString{}, // mutation_type
		&sql.NullString{}, // object_id
		&sql.NullString{}, // object_type
		&sql.NullString{}, // object_display_name
		&[]byte{},         // mutation_data
		&sql.NullString{}, // url
		&sql.NullString{}, // ip_address
		&sql.NullString{}, // status
		&sql.NullString{}, // status_code
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the AuditLog fields.
func (al *AuditLog) assignValues(values ...interface{}) error {
	if m, n := len(values), len(auditlog.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	al.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field created_at", values[0])
	} else if value.Valid {
		al.CreatedAt = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field updated_at", values[1])
	} else if value.Valid {
		al.UpdatedAt = value.Time
	}
	if value, ok := values[2].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field acting_user_id", values[2])
	} else if value.Valid {
		al.ActingUserID = int(value.Int64)
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field organization", values[3])
	} else if value.Valid {
		al.Organization = value.String
	}
	if value, ok := values[4].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field mutation_type", values[4])
	} else if value.Valid {
		al.MutationType = value.String
	}
	if value, ok := values[5].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field object_id", values[5])
	} else if value.Valid {
		al.ObjectID = value.String
	}
	if value, ok := values[6].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field object_type", values[6])
	} else if value.Valid {
		al.ObjectType = value.String
	}
	if value, ok := values[7].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field object_display_name", values[7])
	} else if value.Valid {
		al.ObjectDisplayName = value.String
	}

	if value, ok := values[8].(*[]byte); !ok {
		return fmt.Errorf("unexpected type %T for field mutation_data", values[8])
	} else if value != nil && len(*value) > 0 {
		if err := json.Unmarshal(*value, &al.MutationData); err != nil {
			return fmt.Errorf("unmarshal field mutation_data: %v", err)
		}
	}
	if value, ok := values[9].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field url", values[9])
	} else if value.Valid {
		al.URL = value.String
	}
	if value, ok := values[10].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field ip_address", values[10])
	} else if value.Valid {
		al.IPAddress = value.String
	}
	if value, ok := values[11].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field status", values[11])
	} else if value.Valid {
		al.Status = value.String
	}
	if value, ok := values[12].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field status_code", values[12])
	} else if value.Valid {
		al.StatusCode = value.String
	}
	return nil
}

// Update returns a builder for updating this AuditLog.
// Note that, you need to call AuditLog.Unwrap() before calling this method, if this AuditLog
// was returned from a transaction, and the transaction was committed or rolled back.
func (al *AuditLog) Update() *AuditLogUpdateOne {
	return (&AuditLogClient{config: al.config}).UpdateOne(al)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (al *AuditLog) Unwrap() *AuditLog {
	tx, ok := al.config.driver.(*txDriver)
	if !ok {
		panic("ent: AuditLog is not a transactional entity")
	}
	al.config.driver = tx.drv
	return al
}

// String implements the fmt.Stringer.
func (al *AuditLog) String() string {
	var builder strings.Builder
	builder.WriteString("AuditLog(")
	builder.WriteString(fmt.Sprintf("id=%v", al.ID))
	builder.WriteString(", created_at=")
	builder.WriteString(al.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(al.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", acting_user_id=")
	builder.WriteString(fmt.Sprintf("%v", al.ActingUserID))
	builder.WriteString(", organization=")
	builder.WriteString(al.Organization)
	builder.WriteString(", mutation_type=")
	builder.WriteString(al.MutationType)
	builder.WriteString(", object_id=")
	builder.WriteString(al.ObjectID)
	builder.WriteString(", object_type=")
	builder.WriteString(al.ObjectType)
	builder.WriteString(", object_display_name=")
	builder.WriteString(al.ObjectDisplayName)
	builder.WriteString(", mutation_data=")
	builder.WriteString(fmt.Sprintf("%v", al.MutationData))
	builder.WriteString(", url=")
	builder.WriteString(al.URL)
	builder.WriteString(", ip_address=")
	builder.WriteString(al.IPAddress)
	builder.WriteString(", status=")
	builder.WriteString(al.Status)
	builder.WriteString(", status_code=")
	builder.WriteString(al.StatusCode)
	builder.WriteByte(')')
	return builder.String()
}

// AuditLogs is a parsable slice of AuditLog.
type AuditLogs []*AuditLog

func (al AuditLogs) config(cfg config) {
	for _i := range al {
		al[_i].config = cfg
	}
}
