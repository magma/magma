// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package equipmentpositiondefinition

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the equipmentpositiondefinition type in the database.
	Label = "equipment_position_definition"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldName holds the string denoting the name vertex property in the database.
	FieldName = "name"
	// FieldIndex holds the string denoting the index vertex property in the database.
	FieldIndex = "index"
	// FieldVisibilityLabel holds the string denoting the visibility_label vertex property in the database.
	FieldVisibilityLabel = "visibility_label"

	// Table holds the table name of the equipmentpositiondefinition in the database.
	Table = "equipment_position_definitions"
	// PositionsTable is the table the holds the positions relation/edge.
	PositionsTable = "equipment_positions"
	// PositionsInverseTable is the table name for the EquipmentPosition entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentposition" package.
	PositionsInverseTable = "equipment_positions"
	// PositionsColumn is the table column denoting the positions relation/edge.
	PositionsColumn = "equipment_position_definition"
	// EquipmentTypeTable is the table the holds the equipment_type relation/edge.
	EquipmentTypeTable = "equipment_position_definitions"
	// EquipmentTypeInverseTable is the table name for the EquipmentType entity.
	// It exists in this package in order to avoid circular dependency with the "equipmenttype" package.
	EquipmentTypeInverseTable = "equipment_types"
	// EquipmentTypeColumn is the table column denoting the equipment_type relation/edge.
	EquipmentTypeColumn = "equipment_type_position_definitions"
)

// Columns holds all SQL columns for equipmentpositiondefinition fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldIndex,
	FieldVisibilityLabel,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the EquipmentPositionDefinition type.
var ForeignKeys = []string{
	"equipment_type_position_definitions",
}

var (
	mixin       = schema.EquipmentPositionDefinition{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.EquipmentPositionDefinition{}.Fields()

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
