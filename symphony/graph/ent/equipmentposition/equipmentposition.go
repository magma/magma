// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package equipmentposition

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the equipmentposition type in the database.
	Label = "equipment_position"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"

	// Table holds the table name of the equipmentposition in the database.
	Table = "equipment_positions"
	// DefinitionTable is the table the holds the definition relation/edge.
	DefinitionTable = "equipment_positions"
	// DefinitionInverseTable is the table name for the EquipmentPositionDefinition entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentpositiondefinition" package.
	DefinitionInverseTable = "equipment_position_definitions"
	// DefinitionColumn is the table column denoting the definition relation/edge.
	DefinitionColumn = "definition_id"
	// ParentTable is the table the holds the parent relation/edge.
	ParentTable = "equipment_positions"
	// ParentInverseTable is the table name for the Equipment entity.
	// It exists in this package in order to avoid circular dependency with the "equipment" package.
	ParentInverseTable = "equipment"
	// ParentColumn is the table column denoting the parent relation/edge.
	ParentColumn = "parent_id"
	// AttachmentTable is the table the holds the attachment relation/edge.
	AttachmentTable = "equipment"
	// AttachmentInverseTable is the table name for the Equipment entity.
	// It exists in this package in order to avoid circular dependency with the "equipment" package.
	AttachmentInverseTable = "equipment"
	// AttachmentColumn is the table column denoting the attachment relation/edge.
	AttachmentColumn = "parent_position_id"
)

// Columns holds all SQL columns for equipmentposition fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the EquipmentPosition type.
var ForeignKeys = []string{
	"parent_id",
	"definition_id",
}

var (
	mixin       = schema.EquipmentPosition{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.EquipmentPosition{}.Fields()

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
