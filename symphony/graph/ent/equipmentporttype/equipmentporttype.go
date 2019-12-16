// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package equipmentporttype

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the equipmentporttype type in the database.
	Label = "equipment_port_type"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldName holds the string denoting the name vertex property in the database.
	FieldName = "name"

	// Table holds the table name of the equipmentporttype in the database.
	Table = "equipment_port_types"
	// PropertyTypesTable is the table the holds the property_types relation/edge.
	PropertyTypesTable = "property_types"
	// PropertyTypesInverseTable is the table name for the PropertyType entity.
	// It exists in this package in order to avoid circular dependency with the "propertytype" package.
	PropertyTypesInverseTable = "property_types"
	// PropertyTypesColumn is the table column denoting the property_types relation/edge.
	PropertyTypesColumn = "equipment_port_type_id"
	// LinkPropertyTypesTable is the table the holds the link_property_types relation/edge.
	LinkPropertyTypesTable = "property_types"
	// LinkPropertyTypesInverseTable is the table name for the PropertyType entity.
	// It exists in this package in order to avoid circular dependency with the "propertytype" package.
	LinkPropertyTypesInverseTable = "property_types"
	// LinkPropertyTypesColumn is the table column denoting the link_property_types relation/edge.
	LinkPropertyTypesColumn = "link_equipment_port_type_id"
	// PortDefinitionsTable is the table the holds the port_definitions relation/edge.
	PortDefinitionsTable = "equipment_port_definitions"
	// PortDefinitionsInverseTable is the table name for the EquipmentPortDefinition entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentportdefinition" package.
	PortDefinitionsInverseTable = "equipment_port_definitions"
	// PortDefinitionsColumn is the table column denoting the port_definitions relation/edge.
	PortDefinitionsColumn = "equipment_port_type_id"
)

// Columns holds all SQL columns are equipmentporttype fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
}

var (
	mixin       = schema.EquipmentPortType{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.EquipmentPortType{}.Fields()

	// descCreateTime is the schema descriptor for create_time field.
	descCreateTime = mixinFields[0][0].Descriptor()
	// DefaultCreateTime holds the default value on creation for the create_time field.
	DefaultCreateTime = descCreateTime.Default.(func() time.Time)

	// descUpdateTime is the schema descriptor for update_time field.
	descUpdateTime = mixinFields[0][1].Descriptor()
	// DefaultUpdateTime holds the default value on creation for the update_time field.
	DefaultUpdateTime = descUpdateTime.Default.(func() time.Time)
	// UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	UpdateDefaultUpdateTime = descUpdateTime.UpdateDefault.(func() time.Time)
)
