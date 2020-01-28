// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/project"
)

// Project is the model entity for the Project schema.
type Project struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Description holds the value of the "description" field.
	Description *string `json:"description,omitempty"`
	// Creator holds the value of the "creator" field.
	Creator *string `json:"creator,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ProjectQuery when eager-loading is set.
	Edges               ProjectEdges `json:"edges"`
	project_location_id *string
	type_id             *string
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
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Project) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // name
		&sql.NullString{}, // description
		&sql.NullString{}, // creator
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*Project) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // project_location_id
		&sql.NullInt64{}, // type_id
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
	pr.ID = strconv.FormatInt(value.Int64, 10)
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
	if value, ok := values[4].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field creator", values[4])
	} else if value.Valid {
		pr.Creator = new(string)
		*pr.Creator = value.String
	}
	values = values[5:]
	if len(values) == len(project.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field project_location_id", value)
		} else if value.Valid {
			pr.project_location_id = new(string)
			*pr.project_location_id = strconv.FormatInt(value.Int64, 10)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field type_id", value)
		} else if value.Valid {
			pr.type_id = new(string)
			*pr.type_id = strconv.FormatInt(value.Int64, 10)
		}
	}
	return nil
}

// QueryType queries the type edge of the Project.
func (pr *Project) QueryType() *ProjectTypeQuery {
	return (&ProjectClient{pr.config}).QueryType(pr)
}

// QueryLocation queries the location edge of the Project.
func (pr *Project) QueryLocation() *LocationQuery {
	return (&ProjectClient{pr.config}).QueryLocation(pr)
}

// QueryComments queries the comments edge of the Project.
func (pr *Project) QueryComments() *CommentQuery {
	return (&ProjectClient{pr.config}).QueryComments(pr)
}

// QueryWorkOrders queries the work_orders edge of the Project.
func (pr *Project) QueryWorkOrders() *WorkOrderQuery {
	return (&ProjectClient{pr.config}).QueryWorkOrders(pr)
}

// QueryProperties queries the properties edge of the Project.
func (pr *Project) QueryProperties() *PropertyQuery {
	return (&ProjectClient{pr.config}).QueryProperties(pr)
}

// Update returns a builder for updating this Project.
// Note that, you need to call Project.Unwrap() before calling this method, if this Project
// was returned from a transaction, and the transaction was committed or rolled back.
func (pr *Project) Update() *ProjectUpdateOne {
	return (&ProjectClient{pr.config}).UpdateOne(pr)
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
	if v := pr.Creator; v != nil {
		builder.WriteString(", creator=")
		builder.WriteString(*v)
	}
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (pr *Project) id() int {
	id, _ := strconv.Atoi(pr.ID)
	return id
}

// Projects is a parsable slice of Project.
type Projects []*Project

func (pr Projects) config(cfg config) {
	for _i := range pr {
		pr[_i].config = cfg
	}
}
