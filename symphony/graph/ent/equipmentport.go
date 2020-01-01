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

// EquipmentPort is the model entity for the EquipmentPort schema.
type EquipmentPort struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
}

// FromRows scans the sql response data into EquipmentPort.
func (ep *EquipmentPort) FromRows(rows *sql.Rows) error {
	var scanep struct {
		ID         int
		CreateTime sql.NullTime
		UpdateTime sql.NullTime
	}
	// the order here should be the same as in the `equipmentport.Columns`.
	if err := rows.Scan(
		&scanep.ID,
		&scanep.CreateTime,
		&scanep.UpdateTime,
	); err != nil {
		return err
	}
	ep.ID = strconv.Itoa(scanep.ID)
	ep.CreateTime = scanep.CreateTime.Time
	ep.UpdateTime = scanep.UpdateTime.Time
	return nil
}

// QueryDefinition queries the definition edge of the EquipmentPort.
func (ep *EquipmentPort) QueryDefinition() *EquipmentPortDefinitionQuery {
	return (&EquipmentPortClient{ep.config}).QueryDefinition(ep)
}

// QueryParent queries the parent edge of the EquipmentPort.
func (ep *EquipmentPort) QueryParent() *EquipmentQuery {
	return (&EquipmentPortClient{ep.config}).QueryParent(ep)
}

// QueryLink queries the link edge of the EquipmentPort.
func (ep *EquipmentPort) QueryLink() *LinkQuery {
	return (&EquipmentPortClient{ep.config}).QueryLink(ep)
}

// QueryProperties queries the properties edge of the EquipmentPort.
func (ep *EquipmentPort) QueryProperties() *PropertyQuery {
	return (&EquipmentPortClient{ep.config}).QueryProperties(ep)
}

// QueryEndpoints queries the endpoints edge of the EquipmentPort.
func (ep *EquipmentPort) QueryEndpoints() *ServiceEndpointQuery {
	return (&EquipmentPortClient{ep.config}).QueryEndpoints(ep)
}

// Update returns a builder for updating this EquipmentPort.
// Note that, you need to call EquipmentPort.Unwrap() before calling this method, if this EquipmentPort
// was returned from a transaction, and the transaction was committed or rolled back.
func (ep *EquipmentPort) Update() *EquipmentPortUpdateOne {
	return (&EquipmentPortClient{ep.config}).UpdateOne(ep)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (ep *EquipmentPort) Unwrap() *EquipmentPort {
	tx, ok := ep.config.driver.(*txDriver)
	if !ok {
		panic("ent: EquipmentPort is not a transactional entity")
	}
	ep.config.driver = tx.drv
	return ep
}

// String implements the fmt.Stringer.
func (ep *EquipmentPort) String() string {
	var builder strings.Builder
	builder.WriteString("EquipmentPort(")
	builder.WriteString(fmt.Sprintf("id=%v", ep.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(ep.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(ep.UpdateTime.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (ep *EquipmentPort) id() int {
	id, _ := strconv.Atoi(ep.ID)
	return id
}

// EquipmentPorts is a parsable slice of EquipmentPort.
type EquipmentPorts []*EquipmentPort

// FromRows scans the sql response data into EquipmentPorts.
func (ep *EquipmentPorts) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanep := &EquipmentPort{}
		if err := scanep.FromRows(rows); err != nil {
			return err
		}
		*ep = append(*ep, scanep)
	}
	return nil
}

func (ep EquipmentPorts) config(cfg config) {
	for _i := range ep {
		ep[_i].config = cfg
	}
}
