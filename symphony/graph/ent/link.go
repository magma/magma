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

// Link is the model entity for the Link schema.
type Link struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// FutureState holds the value of the "future_state" field.
	FutureState string `json:"future_state,omitempty"`
}

// FromRows scans the sql response data into Link.
func (l *Link) FromRows(rows *sql.Rows) error {
	var scanl struct {
		ID          int
		CreateTime  sql.NullTime
		UpdateTime  sql.NullTime
		FutureState sql.NullString
	}
	// the order here should be the same as in the `link.Columns`.
	if err := rows.Scan(
		&scanl.ID,
		&scanl.CreateTime,
		&scanl.UpdateTime,
		&scanl.FutureState,
	); err != nil {
		return err
	}
	l.ID = strconv.Itoa(scanl.ID)
	l.CreateTime = scanl.CreateTime.Time
	l.UpdateTime = scanl.UpdateTime.Time
	l.FutureState = scanl.FutureState.String
	return nil
}

// QueryPorts queries the ports edge of the Link.
func (l *Link) QueryPorts() *EquipmentPortQuery {
	return (&LinkClient{l.config}).QueryPorts(l)
}

// QueryWorkOrder queries the work_order edge of the Link.
func (l *Link) QueryWorkOrder() *WorkOrderQuery {
	return (&LinkClient{l.config}).QueryWorkOrder(l)
}

// QueryProperties queries the properties edge of the Link.
func (l *Link) QueryProperties() *PropertyQuery {
	return (&LinkClient{l.config}).QueryProperties(l)
}

// QueryService queries the service edge of the Link.
func (l *Link) QueryService() *ServiceQuery {
	return (&LinkClient{l.config}).QueryService(l)
}

// Update returns a builder for updating this Link.
// Note that, you need to call Link.Unwrap() before calling this method, if this Link
// was returned from a transaction, and the transaction was committed or rolled back.
func (l *Link) Update() *LinkUpdateOne {
	return (&LinkClient{l.config}).UpdateOne(l)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (l *Link) Unwrap() *Link {
	tx, ok := l.config.driver.(*txDriver)
	if !ok {
		panic("ent: Link is not a transactional entity")
	}
	l.config.driver = tx.drv
	return l
}

// String implements the fmt.Stringer.
func (l *Link) String() string {
	var builder strings.Builder
	builder.WriteString("Link(")
	builder.WriteString(fmt.Sprintf("id=%v", l.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(l.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(l.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", future_state=")
	builder.WriteString(l.FutureState)
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (l *Link) id() int {
	id, _ := strconv.Atoi(l.ID)
	return id
}

// Links is a parsable slice of Link.
type Links []*Link

// FromRows scans the sql response data into Links.
func (l *Links) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanl := &Link{}
		if err := scanl.FromRows(rows); err != nil {
			return err
		}
		*l = append(*l, scanl)
	}
	return nil
}

func (l Links) config(cfg config) {
	for _i := range l {
		l[_i].config = cfg
	}
}
