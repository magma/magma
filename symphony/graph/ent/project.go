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
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/user"
)

// Project is the model entity for the Project schema.
type Project struct {
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
	// The values are being populated by the ProjectQuery when eager-loading is set.
	Edges                 ProjectEdges `json:"edges"`
	project_location      *int
	project_creator       *int
	project_type_projects *int
}

// ProjectEdges holds the relations/edges for other nodes in the graph.
type ProjectEdges struct {
	// Type holds the value of the type edge.
	Type *ProjectType
	// Location holds the value of the location edge.
	Location *Location
	// Comments holds the value of the comments edge.
	Comments []*Comment
	// WorkOrders holds the value of the work_orders edge.
	WorkOrders []*WorkOrder
	// Properties holds the value of the properties edge.
	Properties []*Property
	// Creator holds the value of the creator edge.
	Creator *User
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [6]bool
}

// TypeOrErr returns the Type value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ProjectEdges) TypeOrErr() (*ProjectType, error) {
	if e.loadedTypes[0] {
		if e.Type == nil {
			// The edge type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: projecttype.Label}
		}
		return e.Type, nil
	}
	return nil, &NotLoadedError{edge: "type"}
}

// LocationOrErr returns the Location value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ProjectEdges) LocationOrErr() (*Location, error) {
	if e.loadedTypes[1] {
		if e.Location == nil {
			// The edge location was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: location.Label}
		}
		return e.Location, nil
	}
	return nil, &NotLoadedError{edge: "location"}
}

// CommentsOrErr returns the Comments value or an error if the edge
// was not loaded in eager-loading.
func (e ProjectEdges) CommentsOrErr() ([]*Comment, error) {
	if e.loadedTypes[2] {
		return e.Comments, nil
	}
	return nil, &NotLoadedError{edge: "comments"}
}

// WorkOrdersOrErr returns the WorkOrders value or an error if the edge
// was not loaded in eager-loading.
func (e ProjectEdges) WorkOrdersOrErr() ([]*WorkOrder, error) {
	if e.loadedTypes[3] {
		return e.WorkOrders, nil
	}
	return nil, &NotLoadedError{edge: "work_orders"}
}

// PropertiesOrErr returns the Properties value or an error if the edge
// was not loaded in eager-loading.
func (e ProjectEdges) PropertiesOrErr() ([]*Property, error) {
	if e.loadedTypes[4] {
		return e.Properties, nil
	}
	return nil, &NotLoadedError{edge: "properties"}
}

// CreatorOrErr returns the Creator value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ProjectEdges) CreatorOrErr() (*User, error) {
	if e.loadedTypes[5] {
		if e.Creator == nil {
			// The edge creator was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: user.Label}
		}
		return e.Creator, nil
	}
	return nil, &NotLoadedError{edge: "creator"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Project) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // name
		&sql.NullString{}, // description
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*Project) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // project_location
		&sql.NullInt64{}, // project_creator
		&sql.NullInt64{}, // project_type_projects
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Project fields.
func (pr *Project) assignValues(values ...interface{}) error {
	if m, n := len(values), len(project.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	pr.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		pr.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		pr.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[2])
	} else if value.Valid {
		pr.Name = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field description", values[3])
	} else if value.Valid {
		pr.Description = new(string)
		*pr.Description = value.String
	}
	values = values[4:]
	if len(values) == len(project.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field project_location", value)
		} else if value.Valid {
			pr.project_location = new(int)
			*pr.project_location = int(value.Int64)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field project_creator", value)
		} else if value.Valid {
			pr.project_creator = new(int)
			*pr.project_creator = int(value.Int64)
		}
		if value, ok := values[2].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field project_type_projects", value)
		} else if value.Valid {
			pr.project_type_projects = new(int)
			*pr.project_type_projects = int(value.Int64)
		}
	}
	return nil
}

// QueryType queries the type edge of the Project.
func (pr *Project) QueryType() *ProjectTypeQuery {
	return (&ProjectClient{config: pr.config}).QueryType(pr)
}

// QueryLocation queries the location edge of the Project.
func (pr *Project) QueryLocation() *LocationQuery {
	return (&ProjectClient{config: pr.config}).QueryLocation(pr)
}

// QueryComments queries the comments edge of the Project.
func (pr *Project) QueryComments() *CommentQuery {
	return (&ProjectClient{config: pr.config}).QueryComments(pr)
}

// QueryWorkOrders queries the work_orders edge of the Project.
func (pr *Project) QueryWorkOrders() *WorkOrderQuery {
	return (&ProjectClient{config: pr.config}).QueryWorkOrders(pr)
}

// QueryProperties queries the properties edge of the Project.
func (pr *Project) QueryProperties() *PropertyQuery {
	return (&ProjectClient{config: pr.config}).QueryProperties(pr)
}

// QueryCreator queries the creator edge of the Project.
func (pr *Project) QueryCreator() *UserQuery {
	return (&ProjectClient{config: pr.config}).QueryCreator(pr)
}

// Update returns a builder for updating this Project.
// Note that, you need to call Project.Unwrap() before calling this method, if this Project
// was returned from a transaction, and the transaction was committed or rolled back.
func (pr *Project) Update() *ProjectUpdateOne {
	return (&ProjectClient{config: pr.config}).UpdateOne(pr)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (pr *Project) Unwrap() *Project {
	tx, ok := pr.config.driver.(*txDriver)
	if !ok {
		panic("ent: Project is not a transactional entity")
	}
	pr.config.driver = tx.drv
	return pr
}

// String implements the fmt.Stringer.
func (pr *Project) String() string {
	var builder strings.Builder
	builder.WriteString("Project(")
	builder.WriteString(fmt.Sprintf("id=%v", pr.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(pr.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(pr.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(pr.Name)
	if v := pr.Description; v != nil {
		builder.WriteString(", description=")
		builder.WriteString(*v)
	}
	builder.WriteByte(')')
	return builder.String()
}

// Projects is a parsable slice of Project.
type Projects []*Project

func (pr Projects) config(cfg config) {
	for _i := range pr {
		pr[_i].config = cfg
	}
}
