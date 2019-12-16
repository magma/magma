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

// Customer is the model entity for the Customer schema.
type Customer struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// ExternalID holds the value of the "external_id" field.
	ExternalID *string `json:"external_id,omitempty"`
}

// FromRows scans the sql response data into Customer.
func (c *Customer) FromRows(rows *sql.Rows) error {
	var scanc struct {
		ID         int
		CreateTime sql.NullTime
		UpdateTime sql.NullTime
		Name       sql.NullString
		ExternalID sql.NullString
	}
	// the order here should be the same as in the `customer.Columns`.
	if err := rows.Scan(
		&scanc.ID,
		&scanc.CreateTime,
		&scanc.UpdateTime,
		&scanc.Name,
		&scanc.ExternalID,
	); err != nil {
		return err
	}
	c.ID = strconv.Itoa(scanc.ID)
	c.CreateTime = scanc.CreateTime.Time
	c.UpdateTime = scanc.UpdateTime.Time
	c.Name = scanc.Name.String
	if scanc.ExternalID.Valid {
		c.ExternalID = new(string)
		*c.ExternalID = scanc.ExternalID.String
	}
	return nil
}

// QueryServices queries the services edge of the Customer.
func (c *Customer) QueryServices() *ServiceQuery {
	return (&CustomerClient{c.config}).QueryServices(c)
}

// Update returns a builder for updating this Customer.
// Note that, you need to call Customer.Unwrap() before calling this method, if this Customer
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Customer) Update() *CustomerUpdateOne {
	return (&CustomerClient{c.config}).UpdateOne(c)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (c *Customer) Unwrap() *Customer {
	tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Customer is not a transactional entity")
	}
	c.config.driver = tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Customer) String() string {
	var builder strings.Builder
	builder.WriteString("Customer(")
	builder.WriteString(fmt.Sprintf("id=%v", c.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(c.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(c.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(c.Name)
	if v := c.ExternalID; v != nil {
		builder.WriteString(", external_id=")
		builder.WriteString(*v)
	}
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (c *Customer) id() int {
	id, _ := strconv.Atoi(c.ID)
	return id
}

// Customers is a parsable slice of Customer.
type Customers []*Customer

// FromRows scans the sql response data into Customers.
func (c *Customers) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanc := &Customer{}
		if err := scanc.FromRows(rows); err != nil {
			return err
		}
		*c = append(*c, scanc)
	}
	return nil
}

func (c Customers) config(cfg config) {
	for _i := range c {
		c[_i].config = cfg
	}
}
