// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package workorderdefinition

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the workorderdefinition type in the database.
	Label = "work_order_definition"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldIndex holds the string denoting the index vertex property in the database.
	FieldIndex = "index"

	// Table holds the table name of the workorderdefinition in the database.
	Table = "work_order_definitions"
	// TypeTable is the table the holds the type relation/edge.
	TypeTable = "work_order_definitions"
	// TypeInverseTable is the table name for the WorkOrderType entity.
	// It exists in this package in order to avoid circular dependency with the "workordertype" package.
	TypeInverseTable = "work_order_types"
	// TypeColumn is the table column denoting the type relation/edge.
	TypeColumn = "work_order_definition_type"
	// ProjectTypeTable is the table the holds the project_type relation/edge.
	ProjectTypeTable = "work_order_definitions"
	// ProjectTypeInverseTable is the table name for the ProjectType entity.
	// It exists in this package in order to avoid circular dependency with the "projecttype" package.
	ProjectTypeInverseTable = "project_types"
	// ProjectTypeColumn is the table column denoting the project_type relation/edge.
	ProjectTypeColumn = "project_type_work_orders"
)

// Columns holds all SQL columns for workorderdefinition fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldIndex,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the WorkOrderDefinition type.
var ForeignKeys = []string{
	"project_type_work_orders",
	"work_order_definition_type",
}

var (
	mixin       = schema.WorkOrderDefinition{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.WorkOrderDefinition{}.Fields()

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
