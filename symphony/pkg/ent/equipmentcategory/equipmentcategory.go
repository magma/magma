// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package equipmentcategory

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the equipmentcategory type in the database.
	Label = "equipment_category"
	// FieldID holds the string denoting the id field in the database.
	FieldID         = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time" // FieldName holds the string denoting the name vertex property in the database.
	FieldName       = "name"

	// EdgeTypes holds the string denoting the types edge name in mutations.
	EdgeTypes = "types"

	// Table holds the table name of the equipmentcategory in the database.
	Table = "equipment_categories"
	// TypesTable is the table the holds the types relation/edge.
	TypesTable = "equipment_types"
	// TypesInverseTable is the table name for the EquipmentType entity.
	// It exists in this package in order to avoid circular dependency with the "equipmenttype" package.
	TypesInverseTable = "equipment_types"
	// TypesColumn is the table column denoting the types relation/edge.
	TypesColumn = "equipment_type_category"
)

// Columns holds all SQL columns for equipmentcategory fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
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
