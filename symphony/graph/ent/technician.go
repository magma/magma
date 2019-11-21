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

// Technician is the model entity for the Technician schema.
type Technician struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Email holds the value of the "email" field.
	Email string `json:"email,omitempty"`
}

// FromRows scans the sql response data into Technician.
func (t *Technician) FromRows(rows *sql.Rows) error {
	var scant struct {
		ID         int
		CreateTime sql.NullTime
		UpdateTime sql.NullTime
		Name       sql.NullString
		Email      sql.NullString
	}
	// the order here should be the same as in the `technician.Columns`.
	if err := rows.Scan(
		&scant.ID,
		&scant.CreateTime,
		&scant.UpdateTime,
		&scant.Name,
		&scant.Email,
	); err != nil {
		return err
	}
	t.ID = strconv.Itoa(scant.ID)
	t.CreateTime = scant.CreateTime.Time
	t.UpdateTime = scant.UpdateTime.Time
	t.Name = scant.Name.String
	t.Email = scant.Email.String
	return nil
}

// QueryWorkOrders queries the work_orders edge of the Technician.
func (t *Technician) QueryWorkOrders() *WorkOrderQuery {
	return (&TechnicianClient{t.config}).QueryWorkOrders(t)
}

// Update returns a builder for updating this Technician.
// Note that, you need to call Technician.Unwrap() before calling this method, if this Technician
// was returned from a transaction, and the transaction was committed or rolled back.
func (t *Technician) Update() *TechnicianUpdateOne {
	return (&TechnicianClient{t.config}).UpdateOne(t)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (t *Technician) Unwrap() *Technician {
	tx, ok := t.config.driver.(*txDriver)
	if !ok {
		panic("ent: Technician is not a transactional entity")
	}
	t.config.driver = tx.drv
	return t
}

// String implements the fmt.Stringer.
func (t *Technician) String() string {
	var builder strings.Builder
	builder.WriteString("Technician(")
	builder.WriteString(fmt.Sprintf("id=%v", t.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(t.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(t.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(t.Name)
	builder.WriteString(", email=")
	builder.WriteString(t.Email)
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (t *Technician) id() int {
	id, _ := strconv.Atoi(t.ID)
	return id
}

// Technicians is a parsable slice of Technician.
type Technicians []*Technician

// FromRows scans the sql response data into Technicians.
func (t *Technicians) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scant := &Technician{}
		if err := scant.FromRows(rows); err != nil {
			return err
		}
		*t = append(*t, scant)
	}
	return nil
}

func (t Technicians) config(cfg config) {
	for _i := range t {
		t[_i].config = cfg
	}
}
