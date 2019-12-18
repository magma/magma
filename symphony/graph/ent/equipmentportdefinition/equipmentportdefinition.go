// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package equipmentportdefinition

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the equipmentportdefinition type in the database.
	Label = "equipment_port_definition"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldName holds the string denoting the name vertex property in the database.
	FieldName = "name"
	// FieldIndex holds the string denoting the index vertex property in the database.
	FieldIndex = "index"
	// FieldBandwidth holds the string denoting the bandwidth vertex property in the database.
	FieldBandwidth = "bandwidth"
	// FieldVisibilityLabel holds the string denoting the visibility_label vertex property in the database.
	FieldVisibilityLabel = "visibility_label"

	// Table holds the table name of the equipmentportdefinition in the database.
	Table = "equipment_port_definitions"
	// EquipmentPortTypeTable is the table the holds the equipment_port_type relation/edge.
	EquipmentPortTypeTable = "equipment_port_definitions"
	// EquipmentPortTypeInverseTable is the table name for the EquipmentPortType entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentporttype" package.
	EquipmentPortTypeInverseTable = "equipment_port_types"
	// EquipmentPortTypeColumn is the table column denoting the equipment_port_type relation/edge.
	EquipmentPortTypeColumn = "equipment_port_type_id"
	// PortsTable is the table the holds the ports relation/edge.
	PortsTable = "equipment_ports"
	// PortsInverseTable is the table name for the EquipmentPort entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentport" package.
	PortsInverseTable = "equipment_ports"
	// PortsColumn is the table column denoting the ports relation/edge.
	PortsColumn = "definition_id"
	// EquipmentTypeTable is the table the holds the equipment_type relation/edge.
	EquipmentTypeTable = "equipment_port_definitions"
	// EquipmentTypeInverseTable is the table name for the EquipmentType entity.
	// It exists in this package in order to avoid circular dependency with the "equipmenttype" package.
	EquipmentTypeInverseTable = "equipment_types"
	// EquipmentTypeColumn is the table column denoting the equipment_type relation/edge.
	EquipmentTypeColumn = "equipment_type_id"
)

// Columns holds all SQL columns are equipmentportdefinition fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldIndex,
	FieldBandwidth,
	FieldVisibilityLabel,
}

var (
	mixin       = schema.EquipmentPortDefinition{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.EquipmentPortDefinition{}.Fields()

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
