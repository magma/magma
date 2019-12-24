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
}

// FromRows scans the sql response data into Project.
func (pr *Project) FromRows(rows *sql.Rows) error {
	var scanpr struct {
		ID          int
		CreateTime  sql.NullTime
		UpdateTime  sql.NullTime
		Name        sql.NullString
		Description sql.NullString
		Creator     sql.NullString
	}
	// the order here should be the same as in the `project.Columns`.
	if err := rows.Scan(
		&scanpr.ID,
		&scanpr.CreateTime,
		&scanpr.UpdateTime,
		&scanpr.Name,
		&scanpr.Description,
		&scanpr.Creator,
	); err != nil {
		return err
	}
	pr.ID = strconv.Itoa(scanpr.ID)
	pr.CreateTime = scanpr.CreateTime.Time
	pr.UpdateTime = scanpr.UpdateTime.Time
	pr.Name = scanpr.Name.String
	if scanpr.Description.Valid {
		pr.Description = new(string)
		*pr.Description = scanpr.Description.String
	}
	if scanpr.Creator.Valid {
		pr.Creator = new(string)
		*pr.Creator = scanpr.Creator.String
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

// FromRows scans the sql response data into Projects.
func (pr *Projects) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanpr := &Project{}
		if err := scanpr.FromRows(rows); err != nil {
			return err
		}
		*pr = append(*pr, scanpr)
	}
	return nil
}

func (pr Projects) config(cfg config) {
	for _i := range pr {
		pr[_i].config = cfg
	}
}
