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

// ProjectType is the model entity for the ProjectType schema.
type ProjectType struct {
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
}

// FromRows scans the sql response data into ProjectType.
func (pt *ProjectType) FromRows(rows *sql.Rows) error {
	var scanpt struct {
		ID          int
		CreateTime  sql.NullTime
		UpdateTime  sql.NullTime
		Name        sql.NullString
		Description sql.NullString
	}
	// the order here should be the same as in the `projecttype.Columns`.
	if err := rows.Scan(
		&scanpt.ID,
		&scanpt.CreateTime,
		&scanpt.UpdateTime,
		&scanpt.Name,
		&scanpt.Description,
	); err != nil {
		return err
	}
	pt.ID = strconv.Itoa(scanpt.ID)
	pt.CreateTime = scanpt.CreateTime.Time
	pt.UpdateTime = scanpt.UpdateTime.Time
	pt.Name = scanpt.Name.String
	if scanpt.Description.Valid {
		pt.Description = new(string)
		*pt.Description = scanpt.Description.String
	}
	return nil
}

// QueryProjects queries the projects edge of the ProjectType.
func (pt *ProjectType) QueryProjects() *ProjectQuery {
	return (&ProjectTypeClient{pt.config}).QueryProjects(pt)
}

// QueryProperties queries the properties edge of the ProjectType.
func (pt *ProjectType) QueryProperties() *PropertyTypeQuery {
	return (&ProjectTypeClient{pt.config}).QueryProperties(pt)
}

// QueryWorkOrders queries the work_orders edge of the ProjectType.
func (pt *ProjectType) QueryWorkOrders() *WorkOrderDefinitionQuery {
	return (&ProjectTypeClient{pt.config}).QueryWorkOrders(pt)
}

// Update returns a builder for updating this ProjectType.
// Note that, you need to call ProjectType.Unwrap() before calling this method, if this ProjectType
// was returned from a transaction, and the transaction was committed or rolled back.
func (pt *ProjectType) Update() *ProjectTypeUpdateOne {
	return (&ProjectTypeClient{pt.config}).UpdateOne(pt)
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

// id returns the int representation of the ID field.
func (pt *ProjectType) id() int {
	id, _ := strconv.Atoi(pt.ID)
	return id
}

// ProjectTypes is a parsable slice of ProjectType.
type ProjectTypes []*ProjectType

// FromRows scans the sql response data into ProjectTypes.
func (pt *ProjectTypes) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanpt := &ProjectType{}
		if err := scanpt.FromRows(rows); err != nil {
			return err
		}
		*pt = append(*pt, scanpt)
	}
	return nil
}

func (pt ProjectTypes) config(cfg config) {
	for _i := range pt {
		pt[_i].config = cfg
	}
}
