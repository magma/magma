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
	"github.com/facebookincubator/symphony/pkg/ent/projecttype"
)

// ProjectType is the model entity for the ProjectType schema.
type ProjectType struct {
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
	Description *string `json:"description,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ProjectTypeQuery when eager-loading is set.
	Edges ProjectTypeEdges `json:"edges"`
}

// ProjectTypeEdges holds the relations/edges for other nodes in the graph.
type ProjectTypeEdges struct {
	// Projects holds the value of the projects edge.
	Projects []*Project
	// Properties holds the value of the properties edge.
	Properties []*PropertyType
	// WorkOrders holds the value of the work_orders edge.
	WorkOrders []*WorkOrderDefinition
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// ProjectsOrErr returns the Projects value or an error if the edge
// was not loaded in eager-loading.
func (e ProjectTypeEdges) ProjectsOrErr() ([]*Project, error) {
	if e.loadedTypes[0] {
		return e.Projects, nil
	}
	return nil, &NotLoadedError{edge: "projects"}
}

// PropertiesOrErr returns the Properties value or an error if the edge
// was not loaded in eager-loading.
func (e ProjectTypeEdges) PropertiesOrErr() ([]*PropertyType, error) {
	if e.loadedTypes[1] {
		return e.Properties, nil
	}
	return nil, &NotLoadedError{edge: "properties"}
}

// WorkOrdersOrErr returns the WorkOrders value or an error if the edge
// was not loaded in eager-loading.
func (e ProjectTypeEdges) WorkOrdersOrErr() ([]*WorkOrderDefinition, error) {
	if e.loadedTypes[2] {
		return e.WorkOrders, nil
	}
	return nil, &NotLoadedError{edge: "work_orders"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ProjectType) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // name
		&sql.NullString{}, // description
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ProjectType fields.
func (pt *ProjectType) assignValues(values ...interface{}) error {
	if m, n := len(values), len(projecttype.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	pt.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		pt.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		pt.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[2])
	} else if value.Valid {
		pt.Name = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field description", values[3])
	} else if value.Valid {
		pt.Description = new(string)
		*pt.Description = value.String
	}
	return nil
}

// QueryProjects queries the projects edge of the ProjectType.
func (pt *ProjectType) QueryProjects() *ProjectQuery {
	return (&ProjectTypeClient{config: pt.config}).QueryProjects(pt)
}

// QueryProperties queries the properties edge of the ProjectType.
func (pt *ProjectType) QueryProperties() *PropertyTypeQuery {
	return (&ProjectTypeClient{config: pt.config}).QueryProperties(pt)
}

// QueryWorkOrders queries the work_orders edge of the ProjectType.
func (pt *ProjectType) QueryWorkOrders() *WorkOrderDefinitionQuery {
	return (&ProjectTypeClient{config: pt.config}).QueryWorkOrders(pt)
}

// Update returns a builder for updating this ProjectType.
// Note that, you need to call ProjectType.Unwrap() before calling this method, if this ProjectType
// was returned from a transaction, and the transaction was committed or rolled back.
func (pt *ProjectType) Update() *ProjectTypeUpdateOne {
	return (&ProjectTypeClient{config: pt.config}).UpdateOne(pt)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (pt *ProjectType) Unwrap() *ProjectType {
	tx, ok := pt.config.driver.(*txDriver)
	if !ok {
		panic("ent: ProjectType is not a transactional entity")
	}
	pt.config.driver = tx.drv
	return pt
}

// String implements the fmt.Stringer.
func (pt *ProjectType) String() string {
	var builder strings.Builder
	builder.WriteString("ProjectType(")
	builder.WriteString(fmt.Sprintf("id=%v", pt.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(pt.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(pt.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(pt.Name)
	if v := pt.Description; v != nil {
		builder.WriteString(", description=")
		builder.WriteString(*v)
	}
	builder.WriteByte(')')
	return builder.String()
}

// ProjectTypes is a parsable slice of ProjectType.
type ProjectTypes []*ProjectType

func (pt ProjectTypes) config(cfg config) {
	for _i := range pt {
		pt[_i].config = cfg
	}
}
