// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package checklistcategorydefinition

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the checklistcategorydefinition type in the database.
	Label = "check_list_category_definition"
	// FieldID holds the string denoting the id field in the database.
	FieldID          = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime  = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime  = "update_time" // FieldTitle holds the string denoting the title vertex property in the database.
	FieldTitle       = "title"       // FieldDescription holds the string denoting the description vertex property in the database.
	FieldDescription = "description"

	// EdgeCheckListItemDefinitions holds the string denoting the check_list_item_definitions edge name in mutations.
	EdgeCheckListItemDefinitions = "check_list_item_definitions"
	// EdgeWorkOrderType holds the string denoting the work_order_type edge name in mutations.
	EdgeWorkOrderType = "work_order_type"
	// EdgeWorkOrderTemplate holds the string denoting the work_order_template edge name in mutations.
	EdgeWorkOrderTemplate = "work_order_template"

	// Table holds the table name of the checklistcategorydefinition in the database.
	Table = "check_list_category_definitions"
	// CheckListItemDefinitionsTable is the table the holds the check_list_item_definitions relation/edge.
	CheckListItemDefinitionsTable = "check_list_item_definitions"
	// CheckListItemDefinitionsInverseTable is the table name for the CheckListItemDefinition entity.
	// It exists in this package in order to avoid circular dependency with the "checklistitemdefinition" package.
	CheckListItemDefinitionsInverseTable = "check_list_item_definitions"
	// CheckListItemDefinitionsColumn is the table column denoting the check_list_item_definitions relation/edge.
	CheckListItemDefinitionsColumn = "check_list_category_definition_check_list_item_definitions"
	// WorkOrderTypeTable is the table the holds the work_order_type relation/edge.
	WorkOrderTypeTable = "check_list_category_definitions"
	// WorkOrderTypeInverseTable is the table name for the WorkOrderType entity.
	// It exists in this package in order to avoid circular dependency with the "workordertype" package.
	WorkOrderTypeInverseTable = "work_order_types"
	// WorkOrderTypeColumn is the table column denoting the work_order_type relation/edge.
	WorkOrderTypeColumn = "work_order_type_check_list_category_definitions"
	// WorkOrderTemplateTable is the table the holds the work_order_template relation/edge.
	WorkOrderTemplateTable = "check_list_category_definitions"
	// WorkOrderTemplateInverseTable is the table name for the WorkOrderTemplate entity.
	// It exists in this package in order to avoid circular dependency with the "workordertemplate" package.
	WorkOrderTemplateInverseTable = "work_order_templates"
	// WorkOrderTemplateColumn is the table column denoting the work_order_template relation/edge.
	WorkOrderTemplateColumn = "work_order_template_check_list_category_definitions"
)

// Columns holds all SQL columns for checklistcategorydefinition fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldTitle,
	FieldDescription,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the CheckListCategoryDefinition type.
var ForeignKeys = []string{
	"work_order_template_check_list_category_definitions",
	"work_order_type_check_list_category_definitions",
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
	// TitleValidator is a validator for the "title" field. It is called by the builders before save.
	TitleValidator func(string) error
)
