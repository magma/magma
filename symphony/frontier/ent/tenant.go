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

// Tenant is the model entity for the Tenant schema.
type Tenant struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Domains holds the value of the "domains" field.
	Domains []string `json:"domains,omitempty"`
	// Networks holds the value of the "networks" field.
	Networks []string `json:"networks,omitempty"`
	// Tabs holds the value of the "tabs" field.
	Tabs []string `json:"tabs,omitempty"`
	// SSOCert holds the value of the "SSOCert" field.
	SSOCert string `json:"SSOCert,omitempty"`
	// SSOEntryPoint holds the value of the "SSOEntryPoint" field.
	SSOEntryPoint string `json:"SSOEntryPoint,omitempty"`
	// SSOIssuer holds the value of the "SSOIssuer" field.
	SSOIssuer string `json:"SSOIssuer,omitempty"`
}

// FromRows scans the sql response data into Tenant.
func (t *Tenant) FromRows(rows *sql.Rows) error {
	var scant struct {
		ID            int
		CreatedAt     sql.NullTime
		UpdatedAt     sql.NullTime
		Name          sql.NullString
		Domains       []byte
		Networks      []byte
		Tabs          []byte
		SSOCert       sql.NullString
		SSOEntryPoint sql.NullString
		SSOIssuer     sql.NullString
	}
	// the order here should be the same as in the `tenant.Columns`.
	if err := rows.Scan(
		&scant.ID,
		&scant.CreatedAt,
		&scant.UpdatedAt,
		&scant.Name,
		&scant.Domains,
		&scant.Networks,
		&scant.Tabs,
		&scant.SSOCert,
		&scant.SSOEntryPoint,
		&scant.SSOIssuer,
	); err != nil {
		return err
	}
	t.ID = scant.ID
	t.CreatedAt = scant.CreatedAt.Time
	t.UpdatedAt = scant.UpdatedAt.Time
	t.Name = scant.Name.String
	if value := scant.Domains; len(value) > 0 {
		if err := json.Unmarshal(value, &t.Domains); err != nil {
			return fmt.Errorf("unmarshal field domains: %v", err)
		}
	}
	if value := scant.Networks; len(value) > 0 {
		if err := json.Unmarshal(value, &t.Networks); err != nil {
			return fmt.Errorf("unmarshal field networks: %v", err)
		}
	}
	if value := scant.Tabs; len(value) > 0 {
		if err := json.Unmarshal(value, &t.Tabs); err != nil {
			return fmt.Errorf("unmarshal field tabs: %v", err)
		}
	}
	t.SSOCert = scant.SSOCert.String
	t.SSOEntryPoint = scant.SSOEntryPoint.String
	t.SSOIssuer = scant.SSOIssuer.String
	return nil
}

// Update returns a builder for updating this Tenant.
// Note that, you need to call Tenant.Unwrap() before calling this method, if this Tenant
// was returned from a transaction, and the transaction was committed or rolled back.
func (t *Tenant) Update() *TenantUpdateOne {
	return (&TenantClient{t.config}).UpdateOne(t)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (t *Tenant) Unwrap() *Tenant {
	tx, ok := t.config.driver.(*txDriver)
	if !ok {
		panic("ent: Tenant is not a transactional entity")
	}
	t.config.driver = tx.drv
	return t
}

// String implements the fmt.Stringer.
func (t *Tenant) String() string {
	var builder strings.Builder
	builder.WriteString("Tenant(")
	builder.WriteString(fmt.Sprintf("id=%v", t.ID))
	builder.WriteString(", created_at=")
	builder.WriteString(t.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(t.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(t.Name)
	builder.WriteString(", domains=")
	builder.WriteString(fmt.Sprintf("%v", t.Domains))
	builder.WriteString(", networks=")
	builder.WriteString(fmt.Sprintf("%v", t.Networks))
	builder.WriteString(", tabs=")
	builder.WriteString(fmt.Sprintf("%v", t.Tabs))
	builder.WriteString(", SSOCert=")
	builder.WriteString(t.SSOCert)
	builder.WriteString(", SSOEntryPoint=")
	builder.WriteString(t.SSOEntryPoint)
	builder.WriteString(", SSOIssuer=")
	builder.WriteString(t.SSOIssuer)
	builder.WriteByte(')')
	return builder.String()
}

// Tenants is a parsable slice of Tenant.
type Tenants []*Tenant

// FromRows scans the sql response data into Tenants.
func (t *Tenants) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scant := &Tenant{}
		if err := scant.FromRows(rows); err != nil {
			return err
		}
		*t = append(*t, scant)
	}
	return nil
}

func (t Tenants) config(cfg config) {
	for _i := range t {
		t[_i].config = cfg
	}
}
