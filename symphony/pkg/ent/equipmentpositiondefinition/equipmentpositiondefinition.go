// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package equipmentpositiondefinition

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the equipmentpositiondefinition type in the database.
	Label = "equipment_position_definition"
	// FieldID holds the string denoting the id field in the database.
	FieldID              = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime      = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime      = "update_time" // FieldName holds the string denoting the name vertex property in the database.
	FieldName            = "name"        // FieldIndex holds the string denoting the index vertex property in the database.
	FieldIndex           = "index"       // FieldVisibilityLabel holds the string denoting the visibility_label vertex property in the database.
	FieldVisibilityLabel = "visibility_label"

	// EdgePositions holds the string denoting the positions edge name in mutations.
	EdgePositions = "positions"
	// EdgeEquipmentType holds the string denoting the equipment_type edge name in mutations.
	EdgeEquipmentType = "equipment_type"

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

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/facebookincubator/symphony/pkg/ent/runtime"
//
var (
	Hooks  [1]ent.Hook
	Policy ent.Policy
	// DefaultCreateTime holds the default value on creation for the create_time field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the update_time field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	UpdateDefaultUpdateTime func() time.Time
)
