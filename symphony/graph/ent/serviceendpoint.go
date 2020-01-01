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

// ServiceEndpoint is the model entity for the ServiceEndpoint schema.
type ServiceEndpoint struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Role holds the value of the "role" field.
	Role string `json:"role,omitempty"`
}

// FromRows scans the sql response data into ServiceEndpoint.
func (se *ServiceEndpoint) FromRows(rows *sql.Rows) error {
	var scanse struct {
		ID         int
		CreateTime sql.NullTime
		UpdateTime sql.NullTime
		Role       sql.NullString
	}
	// the order here should be the same as in the `serviceendpoint.Columns`.
	if err := rows.Scan(
		&scanse.ID,
		&scanse.CreateTime,
		&scanse.UpdateTime,
		&scanse.Role,
	); err != nil {
		return err
	}
	se.ID = strconv.Itoa(scanse.ID)
	se.CreateTime = scanse.CreateTime.Time
	se.UpdateTime = scanse.UpdateTime.Time
	se.Role = scanse.Role.String
	return nil
}

// QueryPort queries the port edge of the ServiceEndpoint.
func (se *ServiceEndpoint) QueryPort() *EquipmentPortQuery {
	return (&ServiceEndpointClient{se.config}).QueryPort(se)
}

// QueryService queries the service edge of the ServiceEndpoint.
func (se *ServiceEndpoint) QueryService() *ServiceQuery {
	return (&ServiceEndpointClient{se.config}).QueryService(se)
}

// Update returns a builder for updating this ServiceEndpoint.
// Note that, you need to call ServiceEndpoint.Unwrap() before calling this method, if this ServiceEndpoint
// was returned from a transaction, and the transaction was committed or rolled back.
func (se *ServiceEndpoint) Update() *ServiceEndpointUpdateOne {
	return (&ServiceEndpointClient{se.config}).UpdateOne(se)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (se *ServiceEndpoint) Unwrap() *ServiceEndpoint {
	tx, ok := se.config.driver.(*txDriver)
	if !ok {
		panic("ent: ServiceEndpoint is not a transactional entity")
	}
	se.config.driver = tx.drv
	return se
}

// String implements the fmt.Stringer.
func (se *ServiceEndpoint) String() string {
	var builder strings.Builder
	builder.WriteString("ServiceEndpoint(")
	builder.WriteString(fmt.Sprintf("id=%v", se.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(se.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(se.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", role=")
	builder.WriteString(se.Role)
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (se *ServiceEndpoint) id() int {
	id, _ := strconv.Atoi(se.ID)
	return id
}

// ServiceEndpoints is a parsable slice of ServiceEndpoint.
type ServiceEndpoints []*ServiceEndpoint

// FromRows scans the sql response data into ServiceEndpoints.
func (se *ServiceEndpoints) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanse := &ServiceEndpoint{}
		if err := scanse.FromRows(rows); err != nil {
			return err
		}
		*se = append(*se, scanse)
	}
	return nil
}

func (se ServiceEndpoints) config(cfg config) {
	for _i := range se {
		se[_i].config = cfg
	}
}
