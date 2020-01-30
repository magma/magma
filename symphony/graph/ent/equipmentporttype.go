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
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
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
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the EquipmentPortTypeQuery when eager-loading is set.
	Edges EquipmentPortTypeEdges `json:"edges"`
}

// EquipmentPortTypeEdges holds the relations/edges for other nodes in the graph.
type EquipmentPortTypeEdges struct {
	// PropertyTypes holds the value of the property_types edge.
	PropertyTypes []*PropertyType
	// LinkPropertyTypes holds the value of the link_property_types edge.
	LinkPropertyTypes []*PropertyType
	// PortDefinitions holds the value of the port_definitions edge.
	PortDefinitions []*EquipmentPortDefinition
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// PropertyTypesErr returns the PropertyTypes value or an error if the edge
// was not loaded in eager-loading.
func (e EquipmentPortTypeEdges) PropertyTypesErr() ([]*PropertyType, error) {
	if e.loadedTypes[0] {
		return e.PropertyTypes, nil
	}
	return nil, &NotLoadedError{edge: "property_types"}
}

// LinkPropertyTypesErr returns the LinkPropertyTypes value or an error if the edge
// was not loaded in eager-loading.
func (e EquipmentPortTypeEdges) LinkPropertyTypesErr() ([]*PropertyType, error) {
	if e.loadedTypes[1] {
		return e.LinkPropertyTypes, nil
	}
	return nil, &NotLoadedError{edge: "link_property_types"}
}

// PortDefinitionsErr returns the PortDefinitions value or an error if the edge
// was not loaded in eager-loading.
func (e EquipmentPortTypeEdges) PortDefinitionsErr() ([]*EquipmentPortDefinition, error) {
	if e.loadedTypes[2] {
		return e.PortDefinitions, nil
	}
	return nil, &NotLoadedError{edge: "port_definitions"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*EquipmentPortType) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // name
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the EquipmentPortType fields.
func (ept *EquipmentPortType) assignValues(values ...interface{}) error {
	if m, n := len(values), len(equipmentporttype.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	ept.ID = strconv.FormatInt(value.Int64, 10)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		ept.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		ept.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[2])
	} else if value.Valid {
		ept.Name = value.String
	}
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

func (ept EquipmentPortTypes) config(cfg config) {
	for _i := range ept {
		ept[_i].config = cfg
	}
}
