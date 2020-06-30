// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package workordertemplate

import (
	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the workordertemplate type in the database.
	Label = "work_order_template"
	// FieldID holds the string denoting the id field in the database.
	FieldID          = "id"   // FieldName holds the string denoting the name vertex property in the database.
	FieldName        = "name" // FieldDescription holds the string denoting the description vertex property in the database.
	FieldDescription = "description"

	// EdgePropertyTypes holds the string denoting the property_types edge name in mutations.
	EdgePropertyTypes = "property_types"
	// EdgeCheckListCategoryDefinitions holds the string denoting the check_list_category_definitions edge name in mutations.
	EdgeCheckListCategoryDefinitions = "check_list_category_definitions"
	// EdgeType holds the string denoting the type edge name in mutations.
	EdgeType = "type"

	// Table holds the table name of the workordertemplate in the database.
	Table = "work_order_templates"
	// PropertyTypesTable is the table the holds the property_types relation/edge.
	PropertyTypesTable = "property_types"
	// PropertyTypesInverseTable is the table name for the PropertyType entity.
	// It exists in this package in order to avoid circular dependency with the "propertytype" package.
	PropertyTypesInverseTable = "property_types"
	// PropertyTypesColumn is the table column denoting the property_types relation/edge.
	PropertyTypesColumn = "work_order_template_property_types"
	// CheckListCategoryDefinitionsTable is the table the holds the check_list_category_definitions relation/edge.
	CheckListCategoryDefinitionsTable = "check_list_category_definitions"
	// CheckListCategoryDefinitionsInverseTable is the table name for the CheckListCategoryDefinition entity.
	// It exists in this package in order to avoid circular dependency with the "checklistcategorydefinition" package.
	CheckListCategoryDefinitionsInverseTable = "check_list_category_definitions"
	// CheckListCategoryDefinitionsColumn is the table column denoting the check_list_category_definitions relation/edge.
	CheckListCategoryDefinitionsColumn = "work_order_template_check_list_category_definitions"
	// TypeTable is the table the holds the type relation/edge.
	TypeTable = "work_order_templates"
	// TypeInverseTable is the table name for the WorkOrderType entity.
	// It exists in this package in order to avoid circular dependency with the "workordertype" package.
	TypeInverseTable = "work_order_types"
	// TypeColumn is the table column denoting the type relation/edge.
	TypeColumn = "work_order_template_type"
)

// Columns holds all SQL columns for workordertemplate fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldDescription,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the WorkOrderTemplate type.
var ForeignKeys = []string{
	"work_order_template_type",
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
)
