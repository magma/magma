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

// FromRows scans the sql response data into AuditLog.
func (al *AuditLog) FromRows(rows *sql.Rows) error {
	var scanal struct {
		ID                int
		CreatedAt         sql.NullTime
		UpdatedAt         sql.NullTime
		ActingUserID      sql.NullInt64
		Organization      sql.NullString
		MutationType      sql.NullString
		ObjectID          sql.NullString
		ObjectType        sql.NullString
		ObjectDisplayName sql.NullString
		MutationData      []byte
		URL               sql.NullString
		IPAddress         sql.NullString
		Status            sql.NullString
		StatusCode        sql.NullString
	}
	// the order here should be the same as in the `auditlog.Columns`.
	if err := rows.Scan(
		&scanal.ID,
		&scanal.CreatedAt,
		&scanal.UpdatedAt,
		&scanal.ActingUserID,
		&scanal.Organization,
		&scanal.MutationType,
		&scanal.ObjectID,
		&scanal.ObjectType,
		&scanal.ObjectDisplayName,
		&scanal.MutationData,
		&scanal.URL,
		&scanal.IPAddress,
		&scanal.Status,
		&scanal.StatusCode,
	); err != nil {
		return err
	}
	al.ID = scanal.ID
	al.CreatedAt = scanal.CreatedAt.Time
	al.UpdatedAt = scanal.UpdatedAt.Time
	al.ActingUserID = int(scanal.ActingUserID.Int64)
	al.Organization = scanal.Organization.String
	al.MutationType = scanal.MutationType.String
	al.ObjectID = scanal.ObjectID.String
	al.ObjectType = scanal.ObjectType.String
	al.ObjectDisplayName = scanal.ObjectDisplayName.String
	if value := scanal.MutationData; len(value) > 0 {
		if err := json.Unmarshal(value, &al.MutationData); err != nil {
			return fmt.Errorf("unmarshal field mutation_data: %v", err)
		}
	}
	al.URL = scanal.URL.String
	al.IPAddress = scanal.IPAddress.String
	al.Status = scanal.Status.String
	al.StatusCode = scanal.StatusCode.String
	return nil
}

// Update returns a builder for updating this AuditLog.
// Note that, you need to call AuditLog.Unwrap() before calling this method, if this AuditLog
// was returned from a transaction, and the transaction was committed or rolled back.
func (al *AuditLog) Update() *AuditLogUpdateOne {
	return (&AuditLogClient{al.config}).UpdateOne(al)
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

// FromRows scans the sql response data into AuditLogs.
func (al *AuditLogs) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanal := &AuditLog{}
		if err := scanal.FromRows(rows); err != nil {
			return err
		}
		*al = append(*al, scanal)
	}
	return nil
}

func (al AuditLogs) config(cfg config) {
	for _i := range al {
		al[_i].config = cfg
	}
}
