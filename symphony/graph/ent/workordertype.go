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

// WorkOrderType is the model entity for the WorkOrderType schema.
type WorkOrderType struct {
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
	Description string `json:"description,omitempty"`
}

// FromRows scans the sql response data into WorkOrderType.
func (wot *WorkOrderType) FromRows(rows *sql.Rows) error {
	var scanwot struct {
		ID          int
		CreateTime  sql.NullTime
		UpdateTime  sql.NullTime
		Name        sql.NullString
		Description sql.NullString
	}
	// the order here should be the same as in the `workordertype.Columns`.
	if err := rows.Scan(
		&scanwot.ID,
		&scanwot.CreateTime,
		&scanwot.UpdateTime,
		&scanwot.Name,
		&scanwot.Description,
	); err != nil {
		return err
	}
	wot.ID = strconv.Itoa(scanwot.ID)
	wot.CreateTime = scanwot.CreateTime.Time
	wot.UpdateTime = scanwot.UpdateTime.Time
	wot.Name = scanwot.Name.String
	wot.Description = scanwot.Description.String
	return nil
}

// QueryWorkOrders queries the work_orders edge of the WorkOrderType.
func (wot *WorkOrderType) QueryWorkOrders() *WorkOrderQuery {
	return (&WorkOrderTypeClient{wot.config}).QueryWorkOrders(wot)
}

// QueryPropertyTypes queries the property_types edge of the WorkOrderType.
func (wot *WorkOrderType) QueryPropertyTypes() *PropertyTypeQuery {
	return (&WorkOrderTypeClient{wot.config}).QueryPropertyTypes(wot)
}

// QueryDefinitions queries the definitions edge of the WorkOrderType.
func (wot *WorkOrderType) QueryDefinitions() *WorkOrderDefinitionQuery {
	return (&WorkOrderTypeClient{wot.config}).QueryDefinitions(wot)
}

// QueryCheckListDefinitions queries the check_list_definitions edge of the WorkOrderType.
func (wot *WorkOrderType) QueryCheckListDefinitions() *CheckListItemDefinitionQuery {
	return (&WorkOrderTypeClient{wot.config}).QueryCheckListDefinitions(wot)
}

// Update returns a builder for updating this WorkOrderType.
// Note that, you need to call WorkOrderType.Unwrap() before calling this method, if this WorkOrderType
// was returned from a transaction, and the transaction was committed or rolled back.
func (wot *WorkOrderType) Update() *WorkOrderTypeUpdateOne {
	return (&WorkOrderTypeClient{wot.config}).UpdateOne(wot)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (wot *WorkOrderType) Unwrap() *WorkOrderType {
	tx, ok := wot.config.driver.(*txDriver)
	if !ok {
		panic("ent: WorkOrderType is not a transactional entity")
	}
	wot.config.driver = tx.drv
	return wot
}

// String implements the fmt.Stringer.
func (wot *WorkOrderType) String() string {
	var builder strings.Builder
	builder.WriteString("WorkOrderType(")
	builder.WriteString(fmt.Sprintf("id=%v", wot.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(wot.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(wot.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(wot.Name)
	builder.WriteString(", description=")
	builder.WriteString(wot.Description)
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (wot *WorkOrderType) id() int {
	id, _ := strconv.Atoi(wot.ID)
	return id
}

// WorkOrderTypes is a parsable slice of WorkOrderType.
type WorkOrderTypes []*WorkOrderType

// FromRows scans the sql response data into WorkOrderTypes.
func (wot *WorkOrderTypes) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanwot := &WorkOrderType{}
		if err := scanwot.FromRows(rows); err != nil {
			return err
		}
		*wot = append(*wot, scanwot)
	}
	return nil
}

func (wot WorkOrderTypes) config(cfg config) {
	for _i := range wot {
		wot[_i].config = cfg
	}
}
