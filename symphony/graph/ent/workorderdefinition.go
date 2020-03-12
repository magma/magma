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
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// WorkOrderDefinition is the model entity for the WorkOrderDefinition schema.
type WorkOrderDefinition struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Index holds the value of the "index" field.
	Index int `json:"index,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the WorkOrderDefinitionQuery when eager-loading is set.
	Edges                      WorkOrderDefinitionEdges `json:"edges"`
	project_type_work_orders   *int
	work_order_definition_type *int
}

// WorkOrderDefinitionEdges holds the relations/edges for other nodes in the graph.
type WorkOrderDefinitionEdges struct {
	// Type holds the value of the type edge.
	Type *WorkOrderType
	// ProjectType holds the value of the project_type edge.
	ProjectType *ProjectType
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// TypeOrErr returns the Type value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e WorkOrderDefinitionEdges) TypeOrErr() (*WorkOrderType, error) {
	if e.loadedTypes[0] {
		if e.Type == nil {
			// The edge type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: workordertype.Label}
		}
		return e.Type, nil
	}
	return nil, &NotLoadedError{edge: "type"}
}

// ProjectTypeOrErr returns the ProjectType value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e WorkOrderDefinitionEdges) ProjectTypeOrErr() (*ProjectType, error) {
	if e.loadedTypes[1] {
		if e.ProjectType == nil {
			// The edge project_type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: projecttype.Label}
		}
		return e.ProjectType, nil
	}
	return nil, &NotLoadedError{edge: "project_type"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*WorkOrderDefinition) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // id
		&sql.NullTime{},  // create_time
		&sql.NullTime{},  // update_time
		&sql.NullInt64{}, // index
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*WorkOrderDefinition) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // project_type_work_orders
		&sql.NullInt64{}, // work_order_definition_type
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the WorkOrderDefinition fields.
func (wod *WorkOrderDefinition) assignValues(values ...interface{}) error {
	if m, n := len(values), len(workorderdefinition.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	wod.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		wod.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		wod.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field index", values[2])
	} else if value.Valid {
		wod.Index = int(value.Int64)
	}
	values = values[3:]
	if len(values) == len(workorderdefinition.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field project_type_work_orders", value)
		} else if value.Valid {
			wod.project_type_work_orders = new(int)
			*wod.project_type_work_orders = int(value.Int64)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field work_order_definition_type", value)
		} else if value.Valid {
			wod.work_order_definition_type = new(int)
			*wod.work_order_definition_type = int(value.Int64)
		}
	}
	return nil
}

// QueryType queries the type edge of the WorkOrderDefinition.
func (wod *WorkOrderDefinition) QueryType() *WorkOrderTypeQuery {
	return (&WorkOrderDefinitionClient{config: wod.config}).QueryType(wod)
}

// QueryProjectType queries the project_type edge of the WorkOrderDefinition.
func (wod *WorkOrderDefinition) QueryProjectType() *ProjectTypeQuery {
	return (&WorkOrderDefinitionClient{config: wod.config}).QueryProjectType(wod)
}

// Update returns a builder for updating this WorkOrderDefinition.
// Note that, you need to call WorkOrderDefinition.Unwrap() before calling this method, if this WorkOrderDefinition
// was returned from a transaction, and the transaction was committed or rolled back.
func (wod *WorkOrderDefinition) Update() *WorkOrderDefinitionUpdateOne {
	return (&WorkOrderDefinitionClient{config: wod.config}).UpdateOne(wod)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (wod *WorkOrderDefinition) Unwrap() *WorkOrderDefinition {
	tx, ok := wod.config.driver.(*txDriver)
	if !ok {
		panic("ent: WorkOrderDefinition is not a transactional entity")
	}
	wod.config.driver = tx.drv
	return wod
}

// String implements the fmt.Stringer.
func (wod *WorkOrderDefinition) String() string {
	var builder strings.Builder
	builder.WriteString("WorkOrderDefinition(")
	builder.WriteString(fmt.Sprintf("id=%v", wod.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(wod.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(wod.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", index=")
	builder.WriteString(fmt.Sprintf("%v", wod.Index))
	builder.WriteByte(')')
	return builder.String()
}

// WorkOrderDefinitions is a parsable slice of WorkOrderDefinition.
type WorkOrderDefinitions []*WorkOrderDefinition

func (wod WorkOrderDefinitions) config(cfg config) {
	for _i := range wod {
		wod[_i].config = cfg
	}
}
