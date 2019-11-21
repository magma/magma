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

// EquipmentPortType is the model entity for the EquipmentPortType schema.
type EquipmentPortType struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
}

// FromRows scans the sql response data into EquipmentPortType.
func (ept *EquipmentPortType) FromRows(rows *sql.Rows) error {
	var scanept struct {
		ID         int
		CreateTime sql.NullTime
		UpdateTime sql.NullTime
		Name       sql.NullString
	}
	// the order here should be the same as in the `equipmentporttype.Columns`.
	if err := rows.Scan(
		&scanept.ID,
		&scanept.CreateTime,
		&scanept.UpdateTime,
		&scanept.Name,
	); err != nil {
		return err
	}
	ept.ID = strconv.Itoa(scanept.ID)
	ept.CreateTime = scanept.CreateTime.Time
	ept.UpdateTime = scanept.UpdateTime.Time
	ept.Name = scanept.Name.String
	return nil
}

// QueryPropertyTypes queries the property_types edge of the EquipmentPortType.
func (ept *EquipmentPortType) QueryPropertyTypes() *PropertyTypeQuery {
	return (&EquipmentPortTypeClient{ept.config}).QueryPropertyTypes(ept)
}

// QueryLinkPropertyTypes queries the link_property_types edge of the EquipmentPortType.
func (ept *EquipmentPortType) QueryLinkPropertyTypes() *PropertyTypeQuery {
	return (&EquipmentPortTypeClient{ept.config}).QueryLinkPropertyTypes(ept)
}

// QueryPortDefinitions queries the port_definitions edge of the EquipmentPortType.
func (ept *EquipmentPortType) QueryPortDefinitions() *EquipmentPortDefinitionQuery {
	return (&EquipmentPortTypeClient{ept.config}).QueryPortDefinitions(ept)
}

// Update returns a builder for updating this EquipmentPortType.
// Note that, you need to call EquipmentPortType.Unwrap() before calling this method, if this EquipmentPortType
// was returned from a transaction, and the transaction was committed or rolled back.
func (ept *EquipmentPortType) Update() *EquipmentPortTypeUpdateOne {
	return (&EquipmentPortTypeClient{ept.config}).UpdateOne(ept)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (ept *EquipmentPortType) Unwrap() *EquipmentPortType {
	tx, ok := ept.config.driver.(*txDriver)
	if !ok {
		panic("ent: EquipmentPortType is not a transactional entity")
	}
	ept.config.driver = tx.drv
	return ept
}

// String implements the fmt.Stringer.
func (ept *EquipmentPortType) String() string {
	var builder strings.Builder
	builder.WriteString("EquipmentPortType(")
	builder.WriteString(fmt.Sprintf("id=%v", ept.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(ept.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(ept.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(ept.Name)
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (ept *EquipmentPortType) id() int {
	id, _ := strconv.Atoi(ept.ID)
	return id
}

// EquipmentPortTypes is a parsable slice of EquipmentPortType.
type EquipmentPortTypes []*EquipmentPortType

// FromRows scans the sql response data into EquipmentPortTypes.
func (ept *EquipmentPortTypes) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanept := &EquipmentPortType{}
		if err := scanept.FromRows(rows); err != nil {
			return err
		}
		*ept = append(*ept, scanept)
	}
	return nil
}

func (ept EquipmentPortTypes) config(cfg config) {
	for _i := range ept {
		ept[_i].config = cfg
	}
}
