// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package technician

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the technician type in the database.
	Label = "technician"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldName holds the string denoting the name vertex property in the database.
	FieldName = "name"
	// FieldEmail holds the string denoting the email vertex property in the database.
	FieldEmail = "email"

	// Table holds the table name of the technician in the database.
	Table = "technicians"
	// WorkOrdersTable is the table the holds the work_orders relation/edge.
	WorkOrdersTable = "work_orders"
	// WorkOrdersInverseTable is the table name for the WorkOrder entity.
	// It exists in this package in order to avoid circular dependency with the "workorder" package.
	WorkOrdersInverseTable = "work_orders"
	// WorkOrdersColumn is the table column denoting the work_orders relation/edge.
	WorkOrdersColumn = "technician_id"
)

// Columns holds all SQL columns for technician fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldEmail,
}

var (
	mixin       = schema.Technician{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.Technician{}.Fields()

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

	// descName is the schema descriptor for name field.
	descName = fields[0].Descriptor()
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator = descName.Validators[0].(func(string) error)

	// descEmail is the schema descriptor for email field.
	descEmail = fields[1].Descriptor()
	// EmailValidator is a validator for the "email" field. It is called by the builders before save.
	EmailValidator = descEmail.Validators[0].(func(string) error)
)
