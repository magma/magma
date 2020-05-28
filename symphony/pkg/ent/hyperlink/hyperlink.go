// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package hyperlink

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the hyperlink type in the database.
	Label = "hyperlink"
	// FieldID holds the string denoting the id field in the database.
	FieldID         = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time" // FieldURL holds the string denoting the url vertex property in the database.
	FieldURL        = "url"         // FieldName holds the string denoting the name vertex property in the database.
	FieldName       = "name"        // FieldCategory holds the string denoting the category vertex property in the database.
	FieldCategory   = "category"

	// EdgeEquipment holds the string denoting the equipment edge name in mutations.
	EdgeEquipment = "equipment"
	// EdgeLocation holds the string denoting the location edge name in mutations.
	EdgeLocation = "location"
	// EdgeWorkOrder holds the string denoting the work_order edge name in mutations.
	EdgeWorkOrder = "work_order"

	// Table holds the table name of the hyperlink in the database.
	Table = "hyperlinks"
	// EquipmentTable is the table the holds the equipment relation/edge.
	EquipmentTable = "hyperlinks"
	// EquipmentInverseTable is the table name for the Equipment entity.
	// It exists in this package in order to avoid circular dependency with the "equipment" package.
	EquipmentInverseTable = "equipment"
	// EquipmentColumn is the table column denoting the equipment relation/edge.
	EquipmentColumn = "equipment_hyperlinks"
	// LocationTable is the table the holds the location relation/edge.
	LocationTable = "hyperlinks"
	// LocationInverseTable is the table name for the Location entity.
	// It exists in this package in order to avoid circular dependency with the "location" package.
	LocationInverseTable = "locations"
	// LocationColumn is the table column denoting the location relation/edge.
	LocationColumn = "location_hyperlinks"
	// WorkOrderTable is the table the holds the work_order relation/edge.
	WorkOrderTable = "hyperlinks"
	// WorkOrderInverseTable is the table name for the WorkOrder entity.
	// It exists in this package in order to avoid circular dependency with the "workorder" package.
	WorkOrderInverseTable = "work_orders"
	// WorkOrderColumn is the table column denoting the work_order relation/edge.
	WorkOrderColumn = "work_order_hyperlinks"
)

// Columns holds all SQL columns for hyperlink fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldURL,
	FieldName,
	FieldCategory,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the Hyperlink type.
var ForeignKeys = []string{
	"equipment_hyperlinks",
	"location_hyperlinks",
	"work_order_hyperlinks",
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
