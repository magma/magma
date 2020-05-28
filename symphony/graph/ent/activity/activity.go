// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package activity

import (
	"fmt"
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the activity type in the database.
	Label = "activity"
	// FieldID holds the string denoting the id field in the database.
	FieldID           = "id"            // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime   = "create_time"   // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime   = "update_time"   // FieldChangedField holds the string denoting the changed_field vertex property in the database.
	FieldChangedField = "changed_field" // FieldIsCreate holds the string denoting the is_create vertex property in the database.
	FieldIsCreate     = "is_create"     // FieldOldValue holds the string denoting the old_value vertex property in the database.
	FieldOldValue     = "old_value"     // FieldNewValue holds the string denoting the new_value vertex property in the database.
	FieldNewValue     = "new_value"

	// EdgeAuthor holds the string denoting the author edge name in mutations.
	EdgeAuthor = "author"
	// EdgeWorkOrder holds the string denoting the work_order edge name in mutations.
	EdgeWorkOrder = "work_order"

	// Table holds the table name of the activity in the database.
	Table = "activities"
	// AuthorTable is the table the holds the author relation/edge.
	AuthorTable = "activities"
	// AuthorInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	AuthorInverseTable = "users"
	// AuthorColumn is the table column denoting the author relation/edge.
	AuthorColumn = "activity_author"
	// WorkOrderTable is the table the holds the work_order relation/edge.
	WorkOrderTable = "activities"
	// WorkOrderInverseTable is the table name for the WorkOrder entity.
	// It exists in this package in order to avoid circular dependency with the "workorder" package.
	WorkOrderInverseTable = "work_orders"
	// WorkOrderColumn is the table column denoting the work_order relation/edge.
	WorkOrderColumn = "work_order_activities"
)

// Columns holds all SQL columns for activity fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldChangedField,
	FieldIsCreate,
	FieldOldValue,
	FieldNewValue,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the Activity type.
var ForeignKeys = []string{
	"activity_author",
	"work_order_activities",
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
	// DefaultIsCreate holds the default value on creation for the is_create field.
	DefaultIsCreate bool
)

// ChangedField defines the type for the changed_field enum field.
type ChangedField string

// ChangedField values.
const (
	ChangedFieldSTATUS       ChangedField = "STATUS"
	ChangedFieldPRIORITY     ChangedField = "PRIORITY"
	ChangedFieldASSIGNEE     ChangedField = "ASSIGNEE"
	ChangedFieldCREATIONDATE ChangedField = "CREATION_DATE"
	ChangedFieldOWNER        ChangedField = "OWNER"
)

func (s ChangedField) String() string {
	return string(s)
}

// ChangedFieldValidator is a validator for the "cf" field enum values. It is called by the builders before save.
func ChangedFieldValidator(cf ChangedField) error {
	switch cf {
	case ChangedFieldSTATUS, ChangedFieldPRIORITY, ChangedFieldASSIGNEE, ChangedFieldCREATIONDATE, ChangedFieldOWNER:
		return nil
	default:
		return fmt.Errorf("activity: invalid enum value for changed_field field: %q", cf)
	}
}
