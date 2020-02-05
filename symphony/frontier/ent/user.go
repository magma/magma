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
	"github.com/facebookincubator/symphony/frontier/ent/user"
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
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the UserQuery when eager-loading is set.
	Edges UserEdges `json:"edges"`
}

// UserEdges holds the relations/edges for other nodes in the graph.
type UserEdges struct {
	// Tokens holds the value of the tokens edge.
	Tokens []*Token
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// TokensOrErr returns the Tokens value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) TokensOrErr() ([]*Token, error) {
	if e.loadedTypes[0] {
		return e.Tokens, nil
	}
	return nil, &NotLoadedError{edge: "tokens"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*User) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // created_at
		&sql.NullTime{},   // updated_at
		&sql.NullString{}, // email
		&sql.NullString{}, // password
		&sql.NullInt64{},  // role
		&sql.NullString{}, // tenant
		&[]byte{},         // networks
		&[]byte{},         // tabs
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the User fields.
func (u *User) assignValues(values ...interface{}) error {
	if m, n := len(values), len(user.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	u.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field created_at", values[0])
	} else if value.Valid {
		u.CreatedAt = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field updated_at", values[1])
	} else if value.Valid {
		u.UpdatedAt = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field email", values[2])
	} else if value.Valid {
		u.Email = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field password", values[3])
	} else if value.Valid {
		u.Password = value.String
	}
	if value, ok := values[4].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field role", values[4])
	} else if value.Valid {
		u.Role = int(value.Int64)
	}
	if value, ok := values[5].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field tenant", values[5])
	} else if value.Valid {
		u.Tenant = value.String
	}

	if value, ok := values[6].(*[]byte); !ok {
		return fmt.Errorf("unexpected type %T for field networks", values[6])
	} else if value != nil && len(*value) > 0 {
		if err := json.Unmarshal(*value, &u.Networks); err != nil {
			return fmt.Errorf("unmarshal field networks: %v", err)
		}
	}

	if value, ok := values[7].(*[]byte); !ok {
		return fmt.Errorf("unexpected type %T for field tabs", values[7])
	} else if value != nil && len(*value) > 0 {
		if err := json.Unmarshal(*value, &u.Tabs); err != nil {
			return fmt.Errorf("unmarshal field tabs: %v", err)
		}
	}
	return nil
}

// QueryTokens queries the tokens edge of the User.
func (u *User) QueryTokens() *TokenQuery {
	return (&UserClient{u.config}).QueryTokens(u)
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

func (u Users) config(cfg config) {
	for _i := range u {
		u[_i].config = cfg
	}
}
