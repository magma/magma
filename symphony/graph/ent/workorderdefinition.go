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
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
)

// WorkOrderDefinition is the model entity for the WorkOrderDefinition schema.
type WorkOrderDefinition struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Index holds the value of the "index" field.
	Index int `json:"index,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*WorkOrderDefinition) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},
		&sql.NullTime{},
		&sql.NullTime{},
		&sql.NullInt64{},
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the WorkOrderDefinition fields.
func (wod *WorkOrderDefinition) assignValues(values ...interface{}) error {
	if m, n := len(values), len(workorderdefinition.Columns); m != n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	wod.ID = strconv.FormatInt(value.Int64, 10)
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
	return nil
}

// QueryType queries the type edge of the WorkOrderDefinition.
func (wod *WorkOrderDefinition) QueryType() *WorkOrderTypeQuery {
	return (&WorkOrderDefinitionClient{wod.config}).QueryType(wod)
}

// QueryProjectType queries the project_type edge of the WorkOrderDefinition.
func (wod *WorkOrderDefinition) QueryProjectType() *ProjectTypeQuery {
	return (&WorkOrderDefinitionClient{wod.config}).QueryProjectType(wod)
}

// Update returns a builder for updating this WorkOrderDefinition.
// Note that, you need to call WorkOrderDefinition.Unwrap() before calling this method, if this WorkOrderDefinition
// was returned from a transaction, and the transaction was committed or rolled back.
func (wod *WorkOrderDefinition) Update() *WorkOrderDefinitionUpdateOne {
	return (&WorkOrderDefinitionClient{wod.config}).UpdateOne(wod)
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

// id returns the int representation of the ID field.
func (wod *WorkOrderDefinition) id() int {
	id, _ := strconv.Atoi(wod.ID)
	return id
}

// WorkOrderDefinitions is a parsable slice of WorkOrderDefinition.
type WorkOrderDefinitions []*WorkOrderDefinition

func (wod WorkOrderDefinitions) config(cfg config) {
	for _i := range wod {
		wod[_i].config = cfg
	}
}
