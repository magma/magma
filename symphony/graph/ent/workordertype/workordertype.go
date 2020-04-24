// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package workordertype

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the workordertype type in the database.
	Label = "work_order_type"
	// FieldID holds the string denoting the id field in the database.
	FieldID          = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime  = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime  = "update_time" // FieldName holds the string denoting the name vertex property in the database.
	FieldName        = "name"        // FieldDescription holds the string denoting the description vertex property in the database.
	FieldDescription = "description"

	// EdgeWorkOrders holds the string denoting the work_orders edge name in mutations.
	EdgeWorkOrders = "work_orders"
	// EdgePropertyTypes holds the string denoting the property_types edge name in mutations.
	EdgePropertyTypes = "property_types"
	// EdgeDefinitions holds the string denoting the definitions edge name in mutations.
	EdgeDefinitions = "definitions"
	// EdgeCheckListCategories holds the string denoting the check_list_categories edge name in mutations.
	EdgeCheckListCategories = "check_list_categories"
	// EdgeCheckListDefinitions holds the string denoting the check_list_definitions edge name in mutations.
	EdgeCheckListDefinitions = "check_list_definitions"

	// Table holds the table name of the workordertype in the database.
	Table = "work_order_types"
	// WorkOrdersTable is the table the holds the work_orders relation/edge.
	WorkOrdersTable = "work_orders"
	// WorkOrdersInverseTable is the table name for the WorkOrder entity.
	// It exists in this package in order to avoid circular dependency with the "workorder" package.
	WorkOrdersInverseTable = "work_orders"
	// WorkOrdersColumn is the table column denoting the work_orders relation/edge.
	WorkOrdersColumn = "work_order_type"
	// PropertyTypesTable is the table the holds the property_types relation/edge.
	PropertyTypesTable = "property_types"
	// PropertyTypesInverseTable is the table name for the PropertyType entity.
	// It exists in this package in order to avoid circular dependency with the "propertytype" package.
	PropertyTypesInverseTable = "property_types"
	// PropertyTypesColumn is the table column denoting the property_types relation/edge.
	PropertyTypesColumn = "work_order_type_property_types"
	// DefinitionsTable is the table the holds the definitions relation/edge.
	DefinitionsTable = "work_order_definitions"
	// DefinitionsInverseTable is the table name for the WorkOrderDefinition entity.
	// It exists in this package in order to avoid circular dependency with the "workorderdefinition" package.
	DefinitionsInverseTable = "work_order_definitions"
	// DefinitionsColumn is the table column denoting the definitions relation/edge.
	DefinitionsColumn = "work_order_definition_type"
	// CheckListCategoriesTable is the table the holds the check_list_categories relation/edge.
	CheckListCategoriesTable = "check_list_categories"
	// CheckListCategoriesInverseTable is the table name for the CheckListCategory entity.
	// It exists in this package in order to avoid circular dependency with the "checklistcategory" package.
	CheckListCategoriesInverseTable = "check_list_categories"
	// CheckListCategoriesColumn is the table column denoting the check_list_categories relation/edge.
	CheckListCategoriesColumn = "work_order_type_check_list_categories"
	// CheckListDefinitionsTable is the table the holds the check_list_definitions relation/edge.
	CheckListDefinitionsTable = "check_list_item_definitions"
	// CheckListDefinitionsInverseTable is the table name for the CheckListItemDefinition entity.
	// It exists in this package in order to avoid circular dependency with the "checklistitemdefinition" package.
	CheckListDefinitionsInverseTable = "check_list_item_definitions"
	// CheckListDefinitionsColumn is the table column denoting the check_list_definitions relation/edge.
	CheckListDefinitionsColumn = "work_order_type_check_list_definitions"
)

// Columns holds all SQL columns for workordertype fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldDescription,
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
