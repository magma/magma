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

// FromRows scans the sql response data into WorkOrderDefinition.
func (wod *WorkOrderDefinition) FromRows(rows *sql.Rows) error {
	var scanwod struct {
		ID         int
		CreateTime sql.NullTime
		UpdateTime sql.NullTime
		Index      sql.NullInt64
	}
	// the order here should be the same as in the `workorderdefinition.Columns`.
	if err := rows.Scan(
		&scanwod.ID,
		&scanwod.CreateTime,
		&scanwod.UpdateTime,
		&scanwod.Index,
	); err != nil {
		return err
	}
	wod.ID = strconv.Itoa(scanwod.ID)
	wod.CreateTime = scanwod.CreateTime.Time
	wod.UpdateTime = scanwod.UpdateTime.Time
	wod.Index = int(scanwod.Index.Int64)
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

// FromRows scans the sql response data into WorkOrderDefinitions.
func (wod *WorkOrderDefinitions) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanwod := &WorkOrderDefinition{}
		if err := scanwod.FromRows(rows); err != nil {
			return err
		}
		*wod = append(*wod, scanwod)
	}
	return nil
}

func (wod WorkOrderDefinitions) config(cfg config) {
	for _i := range wod {
		wod[_i].config = cfg
	}
}
