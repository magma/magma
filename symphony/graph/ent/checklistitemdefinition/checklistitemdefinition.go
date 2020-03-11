// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package checklistitemdefinition

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the checklistitemdefinition type in the database.
	Label = "check_list_item_definition"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldTitle holds the string denoting the title vertex property in the database.
	FieldTitle = "title"
	// FieldType holds the string denoting the type vertex property in the database.
	FieldType = "type"
	// FieldIndex holds the string denoting the index vertex property in the database.
	FieldIndex = "index"
	// FieldEnumValues holds the string denoting the enum_values vertex property in the database.
	FieldEnumValues = "enum_values"
	// FieldHelpText holds the string denoting the help_text vertex property in the database.
	FieldHelpText = "help_text"

	// Table holds the table name of the checklistitemdefinition in the database.
	Table = "check_list_item_definitions"
	// WorkOrderTypeTable is the table the holds the work_order_type relation/edge.
	WorkOrderTypeTable = "check_list_item_definitions"
	// WorkOrderTypeInverseTable is the table name for the WorkOrderType entity.
	// It exists in this package in order to avoid circular dependency with the "workordertype" package.
	WorkOrderTypeInverseTable = "work_order_types"
	// WorkOrderTypeColumn is the table column denoting the work_order_type relation/edge.
	WorkOrderTypeColumn = "work_order_type_check_list_definitions"
)

// Columns holds all SQL columns for checklistitemdefinition fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldTitle,
	FieldType,
	FieldIndex,
	FieldEnumValues,
	FieldHelpText,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the CheckListItemDefinition type.
var ForeignKeys = []string{
	"work_order_type_check_list_definitions",
}

var (
	mixin       = schema.CheckListItemDefinition{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.CheckListItemDefinition{}.Fields()

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
