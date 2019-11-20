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

// User is the model entity for the User schema.
type User struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Email holds the value of the "email" field.
	Email string `json:"email,omitempty"`
	// Password holds the value of the "password" field.
	Password string `json:"-"`
	// Role holds the value of the "role" field.
	Role int `json:"role,omitempty"`
	// Tenant holds the value of the "tenant" field.
	Tenant string `json:"tenant,omitempty"`
	// Networks holds the value of the "networks" field.
	Networks []string `json:"networks,omitempty"`
	// Tabs holds the value of the "tabs" field.
	Tabs []string `json:"tabs,omitempty"`
}

// FromRows scans the sql response data into User.
func (u *User) FromRows(rows *sql.Rows) error {
	var scanu struct {
		ID        int
		CreatedAt sql.NullTime
		UpdatedAt sql.NullTime
		Email     sql.NullString
		Password  sql.NullString
		Role      sql.NullInt64
		Tenant    sql.NullString
		Networks  []byte
		Tabs      []byte
	}
	// the order here should be the same as in the `user.Columns`.
	if err := rows.Scan(
		&scanu.ID,
		&scanu.CreatedAt,
		&scanu.UpdatedAt,
		&scanu.Email,
		&scanu.Password,
		&scanu.Role,
		&scanu.Tenant,
		&scanu.Networks,
		&scanu.Tabs,
	); err != nil {
		return err
	}
	u.ID = scanu.ID
	u.CreatedAt = scanu.CreatedAt.Time
	u.UpdatedAt = scanu.UpdatedAt.Time
	u.Email = scanu.Email.String
	u.Password = scanu.Password.String
	u.Role = int(scanu.Role.Int64)
	u.Tenant = scanu.Tenant.String
	if value := scanu.Networks; len(value) > 0 {
		if err := json.Unmarshal(value, &u.Networks); err != nil {
			return fmt.Errorf("unmarshal field networks: %v", err)
		}
	}
	if value := scanu.Tabs; len(value) > 0 {
		if err := json.Unmarshal(value, &u.Tabs); err != nil {
			return fmt.Errorf("unmarshal field tabs: %v", err)
		}
	}
	return nil
}

// Update returns a builder for updating this User.
// Note that, you need to call User.Unwrap() before calling this method, if this User
// was returned from a transaction, and the transaction was committed or rolled back.
func (u *User) Update() *UserUpdateOne {
	return (&UserClient{u.config}).UpdateOne(u)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (u *User) Unwrap() *User {
	tx, ok := u.config.driver.(*txDriver)
	if !ok {
		panic("ent: User is not a transactional entity")
	}
	u.config.driver = tx.drv
	return u
}

// String implements the fmt.Stringer.
func (u *User) String() string {
	var builder strings.Builder
	builder.WriteString("User(")
	builder.WriteString(fmt.Sprintf("id=%v", u.ID))
	builder.WriteString(", created_at=")
	builder.WriteString(u.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(u.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", email=")
	builder.WriteString(u.Email)
	builder.WriteString(", password=<sensitive>")
	builder.WriteString(", role=")
	builder.WriteString(fmt.Sprintf("%v", u.Role))
	builder.WriteString(", tenant=")
	builder.WriteString(u.Tenant)
	builder.WriteString(", networks=")
	builder.WriteString(fmt.Sprintf("%v", u.Networks))
	builder.WriteString(", tabs=")
	builder.WriteString(fmt.Sprintf("%v", u.Tabs))
	builder.WriteByte(')')
	return builder.String()
}

// Users is a parsable slice of User.
type Users []*User

// FromRows scans the sql response data into Users.
func (u *Users) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanu := &User{}
		if err := scanu.FromRows(rows); err != nil {
			return err
		}
		*u = append(*u, scanu)
	}
	return nil
}

func (u Users) config(cfg config) {
	for _i := range u {
		u[_i].config = cfg
	}
}
