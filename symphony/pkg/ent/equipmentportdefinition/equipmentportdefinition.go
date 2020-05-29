// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package equipmentportdefinition

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the equipmentportdefinition type in the database.
	Label = "equipment_port_definition"
	// FieldID holds the string denoting the id field in the database.
	FieldID              = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime      = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime      = "update_time" // FieldName holds the string denoting the name vertex property in the database.
	FieldName            = "name"        // FieldIndex holds the string denoting the index vertex property in the database.
	FieldIndex           = "index"       // FieldBandwidth holds the string denoting the bandwidth vertex property in the database.
	FieldBandwidth       = "bandwidth"   // FieldVisibilityLabel holds the string denoting the visibility_label vertex property in the database.
	FieldVisibilityLabel = "visibility_label"

	// EdgeEquipmentPortType holds the string denoting the equipment_port_type edge name in mutations.
	EdgeEquipmentPortType = "equipment_port_type"
	// EdgePorts holds the string denoting the ports edge name in mutations.
	EdgePorts = "ports"
	// EdgeEquipmentType holds the string denoting the equipment_type edge name in mutations.
	EdgeEquipmentType = "equipment_type"

	// Table holds the table name of the equipmentportdefinition in the database.
	Table = "equipment_port_definitions"
	// EquipmentPortTypeTable is the table the holds the equipment_port_type relation/edge.
	EquipmentPortTypeTable = "equipment_port_definitions"
	// EquipmentPortTypeInverseTable is the table name for the EquipmentPortType entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentporttype" package.
	EquipmentPortTypeInverseTable = "equipment_port_types"
	// EquipmentPortTypeColumn is the table column denoting the equipment_port_type relation/edge.
	EquipmentPortTypeColumn = "equipment_port_definition_equipment_port_type"
	// PortsTable is the table the holds the ports relation/edge.
	PortsTable = "equipment_ports"
	// PortsInverseTable is the table name for the EquipmentPort entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentport" package.
	PortsInverseTable = "equipment_ports"
	// PortsColumn is the table column denoting the ports relation/edge.
	PortsColumn = "equipment_port_definition"
	// EquipmentTypeTable is the table the holds the equipment_type relation/edge.
	EquipmentTypeTable = "equipment_port_definitions"
	// EquipmentTypeInverseTable is the table name for the EquipmentType entity.
	// It exists in this package in order to avoid circular dependency with the "equipmenttype" package.
	EquipmentTypeInverseTable = "equipment_types"
	// EquipmentTypeColumn is the table column denoting the equipment_type relation/edge.
	EquipmentTypeColumn = "equipment_type_port_definitions"
)

// Columns holds all SQL columns for equipmentportdefinition fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldIndex,
	FieldBandwidth,
	FieldVisibilityLabel,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the EquipmentPortDefinition type.
var ForeignKeys = []string{
	"equipment_port_definition_equipment_port_type",
	"equipment_type_port_definitions",
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
