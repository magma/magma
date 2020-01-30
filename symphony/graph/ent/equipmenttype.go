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
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
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
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the EquipmentTypeQuery when eager-loading is set.
	Edges       EquipmentTypeEdges `json:"edges"`
	category_id *string
}

// EquipmentTypeEdges holds the relations/edges for other nodes in the graph.
type EquipmentTypeEdges struct {
	// PortDefinitions holds the value of the port_definitions edge.
	PortDefinitions []*EquipmentPortDefinition
	// PositionDefinitions holds the value of the position_definitions edge.
	PositionDefinitions []*EquipmentPositionDefinition
	// PropertyTypes holds the value of the property_types edge.
	PropertyTypes []*PropertyType
	// Equipment holds the value of the equipment edge.
	Equipment []*Equipment
	// Category holds the value of the category edge.
	Category *EquipmentCategory
}

// scanValues returns the types for scanning values from sql.Rows.
func (*EquipmentType) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // name
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*EquipmentType) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // category_id
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the EquipmentType fields.
func (et *EquipmentType) assignValues(values ...interface{}) error {
	if m, n := len(values), len(equipmenttype.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	et.ID = strconv.FormatInt(value.Int64, 10)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		et.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		et.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[2])
	} else if value.Valid {
		et.Name = value.String
	}
	values = values[3:]
	if len(values) == len(equipmenttype.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field category_id", value)
		} else if value.Valid {
			et.category_id = new(string)
			*et.category_id = strconv.FormatInt(value.Int64, 10)
		}
	}
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

func (et EquipmentTypes) config(cfg config) {
	for _i := range et {
		et[_i].config = cfg
	}
}
