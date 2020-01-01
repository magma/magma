// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package equipment

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the equipment type in the database.
	Label = "equipment"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldName holds the string denoting the name vertex property in the database.
	FieldName = "name"
	// FieldFutureState holds the string denoting the future_state vertex property in the database.
	FieldFutureState = "future_state"
	// FieldDeviceID holds the string denoting the device_id vertex property in the database.
	FieldDeviceID = "device_id"
	// FieldExternalID holds the string denoting the external_id vertex property in the database.
	FieldExternalID = "external_id"

	// Table holds the table name of the equipment in the database.
	Table = "equipment"
	// TypeTable is the table the holds the type relation/edge.
	TypeTable = "equipment"
	// TypeInverseTable is the table name for the EquipmentType entity.
	// It exists in this package in order to avoid circular dependency with the "equipmenttype" package.
	TypeInverseTable = "equipment_types"
	// TypeColumn is the table column denoting the type relation/edge.
	TypeColumn = "type_id"
	// LocationTable is the table the holds the location relation/edge.
	LocationTable = "equipment"
	// LocationInverseTable is the table name for the Location entity.
	// It exists in this package in order to avoid circular dependency with the "location" package.
	LocationInverseTable = "locations"
	// LocationColumn is the table column denoting the location relation/edge.
	LocationColumn = "location_id"
	// ParentPositionTable is the table the holds the parent_position relation/edge.
	ParentPositionTable = "equipment"
	// ParentPositionInverseTable is the table name for the EquipmentPosition entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentposition" package.
	ParentPositionInverseTable = "equipment_positions"
	// ParentPositionColumn is the table column denoting the parent_position relation/edge.
	ParentPositionColumn = "parent_position_id"
	// PositionsTable is the table the holds the positions relation/edge.
	PositionsTable = "equipment_positions"
	// PositionsInverseTable is the table name for the EquipmentPosition entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentposition" package.
	PositionsInverseTable = "equipment_positions"
	// PositionsColumn is the table column denoting the positions relation/edge.
	PositionsColumn = "parent_id"
	// PortsTable is the table the holds the ports relation/edge.
	PortsTable = "equipment_ports"
	// PortsInverseTable is the table name for the EquipmentPort entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentport" package.
	PortsInverseTable = "equipment_ports"
	// PortsColumn is the table column denoting the ports relation/edge.
	PortsColumn = "parent_id"
	// WorkOrderTable is the table the holds the work_order relation/edge.
	WorkOrderTable = "equipment"
	// WorkOrderInverseTable is the table name for the WorkOrder entity.
	// It exists in this package in order to avoid circular dependency with the "workorder" package.
	WorkOrderInverseTable = "work_orders"
	// WorkOrderColumn is the table column denoting the work_order relation/edge.
	WorkOrderColumn = "work_order_id"
	// PropertiesTable is the table the holds the properties relation/edge.
	PropertiesTable = "properties"
	// PropertiesInverseTable is the table name for the Property entity.
	// It exists in this package in order to avoid circular dependency with the "property" package.
	PropertiesInverseTable = "properties"
	// PropertiesColumn is the table column denoting the properties relation/edge.
	PropertiesColumn = "equipment_id"
	// FilesTable is the table the holds the files relation/edge.
	FilesTable = "files"
	// FilesInverseTable is the table name for the File entity.
	// It exists in this package in order to avoid circular dependency with the "file" package.
	FilesInverseTable = "files"
	// FilesColumn is the table column denoting the files relation/edge.
	FilesColumn = "equipment_file_id"
)

// Columns holds all SQL columns are equipment fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldFutureState,
	FieldDeviceID,
	FieldExternalID,
}

var (
	mixin       = schema.Equipment{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.Equipment{}.Fields()

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
)
