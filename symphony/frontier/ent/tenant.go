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
	"github.com/facebookincubator/symphony/frontier/ent/tenant"
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

// scanValues returns the types for scanning values from sql.Rows.
func (*Tenant) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // created_at
		&sql.NullTime{},   // updated_at
		&sql.NullString{}, // name
		&[]byte{},         // domains
		&[]byte{},         // networks
		&[]byte{},         // tabs
		&sql.NullString{}, // SSOCert
		&sql.NullString{}, // SSOEntryPoint
		&sql.NullString{}, // SSOIssuer
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Tenant fields.
func (t *Tenant) assignValues(values ...interface{}) error {
	if m, n := len(values), len(tenant.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	t.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field created_at", values[0])
	} else if value.Valid {
		t.CreatedAt = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field updated_at", values[1])
	} else if value.Valid {
		t.UpdatedAt = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[2])
	} else if value.Valid {
		t.Name = value.String
	}

	if value, ok := values[3].(*[]byte); !ok {
		return fmt.Errorf("unexpected type %T for field domains", values[3])
	} else if value != nil && len(*value) > 0 {
		if err := json.Unmarshal(*value, &t.Domains); err != nil {
			return fmt.Errorf("unmarshal field domains: %v", err)
		}
	}

	if value, ok := values[4].(*[]byte); !ok {
		return fmt.Errorf("unexpected type %T for field networks", values[4])
	} else if value != nil && len(*value) > 0 {
		if err := json.Unmarshal(*value, &t.Networks); err != nil {
			return fmt.Errorf("unmarshal field networks: %v", err)
		}
	}

	if value, ok := values[5].(*[]byte); !ok {
		return fmt.Errorf("unexpected type %T for field tabs", values[5])
	} else if value != nil && len(*value) > 0 {
		if err := json.Unmarshal(*value, &t.Tabs); err != nil {
			return fmt.Errorf("unmarshal field tabs: %v", err)
		}
	}
	if value, ok := values[6].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field SSOCert", values[6])
	} else if value.Valid {
		t.SSOCert = value.String
	}
	if value, ok := values[7].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field SSOEntryPoint", values[7])
	} else if value.Valid {
		t.SSOEntryPoint = value.String
	}
	if value, ok := values[8].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field SSOIssuer", values[8])
	} else if value.Valid {
		t.SSOIssuer = value.String
	}
	return nil
}

// Update returns a builder for updating this Tenant.
// Note that, you need to call Tenant.Unwrap() before calling this method, if this Tenant
// was returned from a transaction, and the transaction was committed or rolled back.
func (t *Tenant) Update() *TenantUpdateOne {
	return (&TenantClient{config: t.config}).UpdateOne(t)
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

func (t Tenants) config(cfg config) {
	for _i := range t {
		t[_i].config = cfg
	}
}
