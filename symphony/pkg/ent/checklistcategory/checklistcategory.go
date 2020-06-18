// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package checklistcategory

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the checklistcategory type in the database.
	Label = "check_list_category"
	// FieldID holds the string denoting the id field in the database.
	FieldID          = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime  = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime  = "update_time" // FieldTitle holds the string denoting the title vertex property in the database.
	FieldTitle       = "title"       // FieldDescription holds the string denoting the description vertex property in the database.
	FieldDescription = "description"

	// EdgeCheckListItems holds the string denoting the check_list_items edge name in mutations.
	EdgeCheckListItems = "check_list_items"
	// EdgeWorkOrder holds the string denoting the work_order edge name in mutations.
	EdgeWorkOrder = "work_order"

	// Table holds the table name of the checklistcategory in the database.
	Table = "check_list_categories"
	// CheckListItemsTable is the table the holds the check_list_items relation/edge.
	CheckListItemsTable = "check_list_items"
	// CheckListItemsInverseTable is the table name for the CheckListItem entity.
	// It exists in this package in order to avoid circular dependency with the "checklistitem" package.
	CheckListItemsInverseTable = "check_list_items"
	// CheckListItemsColumn is the table column denoting the check_list_items relation/edge.
	CheckListItemsColumn = "check_list_category_check_list_items"
	// WorkOrderTable is the table the holds the work_order relation/edge.
	WorkOrderTable = "check_list_categories"
	// WorkOrderInverseTable is the table name for the WorkOrder entity.
	// It exists in this package in order to avoid circular dependency with the "workorder" package.
	WorkOrderInverseTable = "work_orders"
	// WorkOrderColumn is the table column denoting the work_order relation/edge.
	WorkOrderColumn = "work_order_check_list_categories"
)

// Columns holds all SQL columns for checklistcategory fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldTitle,
	FieldDescription,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the CheckListCategory type.
var ForeignKeys = []string{
	"work_order_check_list_categories",
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
