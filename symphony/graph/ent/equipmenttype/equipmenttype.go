// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package equipmenttype

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the equipmenttype type in the database.
	Label = "equipment_type"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldName holds the string denoting the name vertex property in the database.
	FieldName = "name"

	// Table holds the table name of the equipmenttype in the database.
	Table = "equipment_types"
	// PortDefinitionsTable is the table the holds the port_definitions relation/edge.
	PortDefinitionsTable = "equipment_port_definitions"
	// PortDefinitionsInverseTable is the table name for the EquipmentPortDefinition entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentportdefinition" package.
	PortDefinitionsInverseTable = "equipment_port_definitions"
	// PortDefinitionsColumn is the table column denoting the port_definitions relation/edge.
	PortDefinitionsColumn = "equipment_type_port_definitions"
	// PositionDefinitionsTable is the table the holds the position_definitions relation/edge.
	PositionDefinitionsTable = "equipment_position_definitions"
	// PositionDefinitionsInverseTable is the table name for the EquipmentPositionDefinition entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentpositiondefinition" package.
	PositionDefinitionsInverseTable = "equipment_position_definitions"
	// PositionDefinitionsColumn is the table column denoting the position_definitions relation/edge.
	PositionDefinitionsColumn = "equipment_type_position_definitions"
	// PropertyTypesTable is the table the holds the property_types relation/edge.
	PropertyTypesTable = "property_types"
	// PropertyTypesInverseTable is the table name for the PropertyType entity.
	// It exists in this package in order to avoid circular dependency with the "propertytype" package.
	PropertyTypesInverseTable = "property_types"
	// PropertyTypesColumn is the table column denoting the property_types relation/edge.
	PropertyTypesColumn = "equipment_type_property_types"
	// EquipmentTable is the table the holds the equipment relation/edge.
	EquipmentTable = "equipment"
	// EquipmentInverseTable is the table name for the Equipment entity.
	// It exists in this package in order to avoid circular dependency with the "equipment" package.
	EquipmentInverseTable = "equipment"
	// EquipmentColumn is the table column denoting the equipment relation/edge.
	EquipmentColumn = "equipment_type"
	// CategoryTable is the table the holds the category relation/edge.
	CategoryTable = "equipment_types"
	// CategoryInverseTable is the table name for the EquipmentCategory entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentcategory" package.
	CategoryInverseTable = "equipment_categories"
	// CategoryColumn is the table column denoting the category relation/edge.
	CategoryColumn = "equipment_type_category"
)

// Columns holds all SQL columns for equipmenttype fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the EquipmentType type.
var ForeignKeys = []string{
	"equipment_type_category",
}

var (
	mixin       = schema.EquipmentType{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.EquipmentType{}.Fields()

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
