// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/pkg/ent/usersgroup"
)

// UsersGroup is the model entity for the UsersGroup schema.
type UsersGroup struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// Status holds the value of the "status" field.
	Status usersgroup.Status `json:"status,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the UsersGroupQuery when eager-loading is set.
	Edges UsersGroupEdges `json:"edges"`
}

// UsersGroupEdges holds the relations/edges for other nodes in the graph.
type UsersGroupEdges struct {
	// Members holds the value of the members edge.
	Members []*User
	// Policies holds the value of the policies edge.
	Policies []*PermissionsPolicy
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// MembersOrErr returns the Members value or an error if the edge
// was not loaded in eager-loading.
func (e UsersGroupEdges) MembersOrErr() ([]*User, error) {
	if e.loadedTypes[0] {
		return e.Members, nil
	}
	return nil, &NotLoadedError{edge: "members"}
}

// PoliciesOrErr returns the Policies value or an error if the edge
// was not loaded in eager-loading.
func (e UsersGroupEdges) PoliciesOrErr() ([]*PermissionsPolicy, error) {
	if e.loadedTypes[1] {
		return e.Policies, nil
	}
	return nil, &NotLoadedError{edge: "policies"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*UsersGroup) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // name
		&sql.NullString{}, // description
		&sql.NullString{}, // status
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the UsersGroup fields.
func (ug *UsersGroup) assignValues(values ...interface{}) error {
	if m, n := len(values), len(usersgroup.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	ug.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		ug.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		ug.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[2])
	} else if value.Valid {
		ug.Name = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field description", values[3])
	} else if value.Valid {
		ug.Description = value.String
	}
	if value, ok := values[4].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field status", values[4])
	} else if value.Valid {
		ug.Status = usersgroup.Status(value.String)
	}
	return nil
}

// QueryMembers queries the members edge of the UsersGroup.
func (ug *UsersGroup) QueryMembers() *UserQuery {
	return (&UsersGroupClient{config: ug.config}).QueryMembers(ug)
}

// QueryPolicies queries the policies edge of the UsersGroup.
func (ug *UsersGroup) QueryPolicies() *PermissionsPolicyQuery {
	return (&UsersGroupClient{config: ug.config}).QueryPolicies(ug)
}

// Update returns a builder for updating this UsersGroup.
// Note that, you need to call UsersGroup.Unwrap() before calling this method, if this UsersGroup
// was returned from a transaction, and the transaction was committed or rolled back.
func (ug *UsersGroup) Update() *UsersGroupUpdateOne {
	return (&UsersGroupClient{config: ug.config}).UpdateOne(ug)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (ug *UsersGroup) Unwrap() *UsersGroup {
	tx, ok := ug.config.driver.(*txDriver)
	if !ok {
		panic("ent: UsersGroup is not a transactional entity")
	}
	ug.config.driver = tx.drv
	return ug
}

// String implements the fmt.Stringer.
func (ug *UsersGroup) String() string {
	var builder strings.Builder
	builder.WriteString("UsersGroup(")
	builder.WriteString(fmt.Sprintf("id=%v", ug.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(ug.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(ug.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(ug.Name)
	builder.WriteString(", description=")
	builder.WriteString(ug.Description)
	builder.WriteString(", status=")
	builder.WriteString(fmt.Sprintf("%v", ug.Status))
	builder.WriteByte(')')
	return builder.String()
}

// UsersGroups is a parsable slice of UsersGroup.
type UsersGroups []*UsersGroup

func (ug UsersGroups) config(cfg config) {
	for _i := range ug {
		ug[_i].config = cfg
	}
}
