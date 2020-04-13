// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package equipment

import (
	"time"
)

const (
	// Label holds the string label denoting the equipment type in the database.
	Label = "equipment"
	// FieldID holds the string denoting the id field in the database.
	FieldID          = "id"           // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime  = "create_time"  // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime  = "update_time"  // FieldName holds the string denoting the name vertex property in the database.
	FieldName        = "name"         // FieldFutureState holds the string denoting the future_state vertex property in the database.
	FieldFutureState = "future_state" // FieldDeviceID holds the string denoting the device_id vertex property in the database.
	FieldDeviceID    = "device_id"    // FieldExternalID holds the string denoting the external_id vertex property in the database.
	FieldExternalID  = "external_id"

	// EdgeType holds the string denoting the type edge name in mutations.
	EdgeType = "type"
	// EdgeLocation holds the string denoting the location edge name in mutations.
	EdgeLocation = "location"
	// EdgeParentPosition holds the string denoting the parent_position edge name in mutations.
	EdgeParentPosition = "parent_position"
	// EdgePositions holds the string denoting the positions edge name in mutations.
	EdgePositions = "positions"
	// EdgePorts holds the string denoting the ports edge name in mutations.
	EdgePorts = "ports"
	// EdgeWorkOrder holds the string denoting the work_order edge name in mutations.
	EdgeWorkOrder = "work_order"
	// EdgeProperties holds the string denoting the properties edge name in mutations.
	EdgeProperties = "properties"
	// EdgeFiles holds the string denoting the files edge name in mutations.
	EdgeFiles = "files"
	// EdgeHyperlinks holds the string denoting the hyperlinks edge name in mutations.
	EdgeHyperlinks = "hyperlinks"
	// EdgeEndpoints holds the string denoting the endpoints edge name in mutations.
	EdgeEndpoints = "endpoints"

	// Table holds the table name of the equipment in the database.
	Table = "equipment"
	// TypeTable is the table the holds the type relation/edge.
	TypeTable = "equipment"
	// TypeInverseTable is the table name for the EquipmentType entity.
	// It exists in this package in order to avoid circular dependency with the "equipmenttype" package.
	TypeInverseTable = "equipment_types"
	// TypeColumn is the table column denoting the type relation/edge.
	TypeColumn = "equipment_type"
	// LocationTable is the table the holds the location relation/edge.
	LocationTable = "equipment"
	// LocationInverseTable is the table name for the Location entity.
	// It exists in this package in order to avoid circular dependency with the "location" package.
	LocationInverseTable = "locations"
	// LocationColumn is the table column denoting the location relation/edge.
	LocationColumn = "location_equipment"
	// ParentPositionTable is the table the holds the parent_position relation/edge.
	ParentPositionTable = "equipment"
	// ParentPositionInverseTable is the table name for the EquipmentPosition entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentposition" package.
	ParentPositionInverseTable = "equipment_positions"
	// ParentPositionColumn is the table column denoting the parent_position relation/edge.
	ParentPositionColumn = "equipment_position_attachment"
	// PositionsTable is the table the holds the positions relation/edge.
	PositionsTable = "equipment_positions"
	// PositionsInverseTable is the table name for the EquipmentPosition entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentposition" package.
	PositionsInverseTable = "equipment_positions"
	// PositionsColumn is the table column denoting the positions relation/edge.
	PositionsColumn = "equipment_positions"
	// PortsTable is the table the holds the ports relation/edge.
	PortsTable = "equipment_ports"
	// PortsInverseTable is the table name for the EquipmentPort entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentport" package.
	PortsInverseTable = "equipment_ports"
	// PortsColumn is the table column denoting the ports relation/edge.
	PortsColumn = "equipment_ports"
	// WorkOrderTable is the table the holds the work_order relation/edge.
	WorkOrderTable = "equipment"
	// WorkOrderInverseTable is the table name for the WorkOrder entity.
	// It exists in this package in order to avoid circular dependency with the "workorder" package.
	WorkOrderInverseTable = "work_orders"
	// WorkOrderColumn is the table column denoting the work_order relation/edge.
	WorkOrderColumn = "equipment_work_order"
	// PropertiesTable is the table the holds the properties relation/edge.
	PropertiesTable = "properties"
	// PropertiesInverseTable is the table name for the Property entity.
	// It exists in this package in order to avoid circular dependency with the "property" package.
	PropertiesInverseTable = "properties"
	// PropertiesColumn is the table column denoting the properties relation/edge.
	PropertiesColumn = "equipment_properties"
	// FilesTable is the table the holds the files relation/edge.
	FilesTable = "files"
	// FilesInverseTable is the table name for the File entity.
	// It exists in this package in order to avoid circular dependency with the "file" package.
	FilesInverseTable = "files"
	// FilesColumn is the table column denoting the files relation/edge.
	FilesColumn = "equipment_files"
	// HyperlinksTable is the table the holds the hyperlinks relation/edge.
	HyperlinksTable = "hyperlinks"
	// HyperlinksInverseTable is the table name for the Hyperlink entity.
	// It exists in this package in order to avoid circular dependency with the "hyperlink" package.
	HyperlinksInverseTable = "hyperlinks"
	// HyperlinksColumn is the table column denoting the hyperlinks relation/edge.
	HyperlinksColumn = "equipment_hyperlinks"
	// EndpointsTable is the table the holds the endpoints relation/edge.
	EndpointsTable = "service_endpoints"
	// EndpointsInverseTable is the table name for the ServiceEndpoint entity.
	// It exists in this package in order to avoid circular dependency with the "serviceendpoint" package.
	EndpointsInverseTable = "service_endpoints"
	// EndpointsColumn is the table column denoting the endpoints relation/edge.
	EndpointsColumn = "service_endpoint_equipment"
)

// Columns holds all SQL columns for equipment fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldFutureState,
	FieldDeviceID,
	FieldExternalID,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the Equipment type.
var ForeignKeys = []string{
	"equipment_type",
	"equipment_work_order",
	"equipment_position_attachment",
	"location_equipment",
}

var (
	// DefaultCreateTime holds the default value on creation for the create_time field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the update_time field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	UpdateDefaultUpdateTime func() time.Time
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DeviceIDValidator is a validator for the "device_id" field. It is called by the builders before save.
	DeviceIDValidator func(string) error
)
