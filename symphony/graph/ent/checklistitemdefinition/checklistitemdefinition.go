// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package checklistitemdefinition

const (
	// Label holds the string label denoting the checklistitemdefinition type in the database.
	Label = "check_list_item_definition"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
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
	WorkOrderTypeColumn = "work_order_type_id"
)

// Columns holds all SQL columns are checklistitemdefinition fields.
var Columns = []string{
	FieldID,
	FieldTitle,
	FieldType,
	FieldIndex,
	FieldEnumValues,
	FieldHelpText,
}
