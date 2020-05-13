// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package checklistitemdefinition

import (
	"fmt"
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the checklistitemdefinition type in the database.
	Label = "check_list_item_definition"
	// FieldID holds the string denoting the id field in the database.
	FieldID                     = "id"                        // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime             = "create_time"               // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime             = "update_time"               // FieldTitle holds the string denoting the title vertex property in the database.
	FieldTitle                  = "title"                     // FieldType holds the string denoting the type vertex property in the database.
	FieldType                   = "type"                      // FieldIndex holds the string denoting the index vertex property in the database.
	FieldIndex                  = "index"                     // FieldEnumValues holds the string denoting the enum_values vertex property in the database.
	FieldEnumValues             = "enum_values"               // FieldEnumSelectionModeValue holds the string denoting the enum_selection_mode_value vertex property in the database.
	FieldEnumSelectionModeValue = "enum_selection_mode_value" // FieldHelpText holds the string denoting the help_text vertex property in the database.
	FieldHelpText               = "help_text"

	// EdgeCheckListCategoryDefinition holds the string denoting the check_list_category_definition edge name in mutations.
	EdgeCheckListCategoryDefinition = "check_list_category_definition"

	// Table holds the table name of the checklistitemdefinition in the database.
	Table = "check_list_item_definitions"
	// CheckListCategoryDefinitionTable is the table the holds the check_list_category_definition relation/edge.
	CheckListCategoryDefinitionTable = "check_list_item_definitions"
	// CheckListCategoryDefinitionInverseTable is the table name for the CheckListCategoryDefinition entity.
	// It exists in this package in order to avoid circular dependency with the "checklistcategorydefinition" package.
	CheckListCategoryDefinitionInverseTable = "check_list_category_definitions"
	// CheckListCategoryDefinitionColumn is the table column denoting the check_list_category_definition relation/edge.
	CheckListCategoryDefinitionColumn = "check_list_category_definition_check_list_item_definitions"
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
	FieldEnumSelectionModeValue,
	FieldHelpText,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the CheckListItemDefinition type.
var ForeignKeys = []string{
	"check_list_category_definition_check_list_item_definitions",
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/facebookincubator/symphony/graph/ent/runtime"
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

// EnumSelectionModeValue defines the type for the enum_selection_mode_value enum field.
type EnumSelectionModeValue string

// EnumSelectionModeValue values.
const (
	EnumSelectionModeValueSingle   EnumSelectionModeValue = "single"
	EnumSelectionModeValueMultiple EnumSelectionModeValue = "multiple"
)

func (s EnumSelectionModeValue) String() string {
	return string(s)
}

// EnumSelectionModeValueValidator is a validator for the "esmv" field enum values. It is called by the builders before save.
func EnumSelectionModeValueValidator(esmv EnumSelectionModeValue) error {
	switch esmv {
	case EnumSelectionModeValueSingle, EnumSelectionModeValueMultiple:
		return nil
	default:
		return fmt.Errorf("checklistitemdefinition: invalid enum value for enum_selection_mode_value field: %q", esmv)
	}
}
