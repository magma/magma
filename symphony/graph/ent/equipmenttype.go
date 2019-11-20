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

// EquipmentType is the model entity for the EquipmentType schema.
type EquipmentType struct {
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

// FromRows scans the sql response data into EquipmentType.
func (et *EquipmentType) FromRows(rows *sql.Rows) error {
	var scanet struct {
		ID         int
		CreateTime sql.NullTime
		UpdateTime sql.NullTime
		Name       sql.NullString
	}
	// the order here should be the same as in the `equipmenttype.Columns`.
	if err := rows.Scan(
		&scanet.ID,
		&scanet.CreateTime,
		&scanet.UpdateTime,
		&scanet.Name,
	); err != nil {
		return err
	}
	et.ID = strconv.Itoa(scanet.ID)
	et.CreateTime = scanet.CreateTime.Time
	et.UpdateTime = scanet.UpdateTime.Time
	et.Name = scanet.Name.String
	return nil
}

// QueryPortDefinitions queries the port_definitions edge of the EquipmentType.
func (et *EquipmentType) QueryPortDefinitions() *EquipmentPortDefinitionQuery {
	return (&EquipmentTypeClient{et.config}).QueryPortDefinitions(et)
}

// QueryPositionDefinitions queries the position_definitions edge of the EquipmentType.
func (et *EquipmentType) QueryPositionDefinitions() *EquipmentPositionDefinitionQuery {
	return (&EquipmentTypeClient{et.config}).QueryPositionDefinitions(et)
}

// QueryPropertyTypes queries the property_types edge of the EquipmentType.
func (et *EquipmentType) QueryPropertyTypes() *PropertyTypeQuery {
	return (&EquipmentTypeClient{et.config}).QueryPropertyTypes(et)
}

// QueryEquipment queries the equipment edge of the EquipmentType.
func (et *EquipmentType) QueryEquipment() *EquipmentQuery {
	return (&EquipmentTypeClient{et.config}).QueryEquipment(et)
}

// QueryCategory queries the category edge of the EquipmentType.
func (et *EquipmentType) QueryCategory() *EquipmentCategoryQuery {
	return (&EquipmentTypeClient{et.config}).QueryCategory(et)
}

// Update returns a builder for updating this EquipmentType.
// Note that, you need to call EquipmentType.Unwrap() before calling this method, if this EquipmentType
// was returned from a transaction, and the transaction was committed or rolled back.
func (et *EquipmentType) Update() *EquipmentTypeUpdateOne {
	return (&EquipmentTypeClient{et.config}).UpdateOne(et)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (et *EquipmentType) Unwrap() *EquipmentType {
	tx, ok := et.config.driver.(*txDriver)
	if !ok {
		panic("ent: EquipmentType is not a transactional entity")
	}
	et.config.driver = tx.drv
	return et
}

// String implements the fmt.Stringer.
func (et *EquipmentType) String() string {
	var builder strings.Builder
	builder.WriteString("EquipmentType(")
	builder.WriteString(fmt.Sprintf("id=%v", et.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(et.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(et.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(et.Name)
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (et *EquipmentType) id() int {
	id, _ := strconv.Atoi(et.ID)
	return id
}

// EquipmentTypes is a parsable slice of EquipmentType.
type EquipmentTypes []*EquipmentType

// FromRows scans the sql response data into EquipmentTypes.
func (et *EquipmentTypes) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanet := &EquipmentType{}
		if err := scanet.FromRows(rows); err != nil {
			return err
		}
		*et = append(*et, scanet)
	}
	return nil
}

func (et EquipmentTypes) config(cfg config) {
	for _i := range et {
		et[_i].config = cfg
	}
}
