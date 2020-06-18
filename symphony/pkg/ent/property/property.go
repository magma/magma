// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package property

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the property type in the database.
	Label = "property"
	// FieldID holds the string denoting the id field in the database.
	FieldID           = "id"             // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime   = "create_time"    // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime   = "update_time"    // FieldIntVal holds the string denoting the int_val vertex property in the database.
	FieldIntVal       = "int_val"        // FieldBoolVal holds the string denoting the bool_val vertex property in the database.
	FieldBoolVal      = "bool_val"       // FieldFloatVal holds the string denoting the float_val vertex property in the database.
	FieldFloatVal     = "float_val"      // FieldLatitudeVal holds the string denoting the latitude_val vertex property in the database.
	FieldLatitudeVal  = "latitude_val"   // FieldLongitudeVal holds the string denoting the longitude_val vertex property in the database.
	FieldLongitudeVal = "longitude_val"  // FieldRangeFromVal holds the string denoting the range_from_val vertex property in the database.
	FieldRangeFromVal = "range_from_val" // FieldRangeToVal holds the string denoting the range_to_val vertex property in the database.
	FieldRangeToVal   = "range_to_val"   // FieldStringVal holds the string denoting the string_val vertex property in the database.
	FieldStringVal    = "string_val"

	// EdgeType holds the string denoting the type edge name in mutations.
	EdgeType = "type"
	// EdgeLocation holds the string denoting the location edge name in mutations.
	EdgeLocation = "location"
	// EdgeEquipment holds the string denoting the equipment edge name in mutations.
	EdgeEquipment = "equipment"
	// EdgeService holds the string denoting the service edge name in mutations.
	EdgeService = "service"
	// EdgeEquipmentPort holds the string denoting the equipment_port edge name in mutations.
	EdgeEquipmentPort = "equipment_port"
	// EdgeLink holds the string denoting the link edge name in mutations.
	EdgeLink = "link"
	// EdgeWorkOrder holds the string denoting the work_order edge name in mutations.
	EdgeWorkOrder = "work_order"
	// EdgeProject holds the string denoting the project edge name in mutations.
	EdgeProject = "project"
	// EdgeEquipmentValue holds the string denoting the equipment_value edge name in mutations.
	EdgeEquipmentValue = "equipment_value"
	// EdgeLocationValue holds the string denoting the location_value edge name in mutations.
	EdgeLocationValue = "location_value"
	// EdgeServiceValue holds the string denoting the service_value edge name in mutations.
	EdgeServiceValue = "service_value"
	// EdgeWorkOrderValue holds the string denoting the work_order_value edge name in mutations.
	EdgeWorkOrderValue = "work_order_value"
	// EdgeUserValue holds the string denoting the user_value edge name in mutations.
	EdgeUserValue = "user_value"

	// Table holds the table name of the property in the database.
	Table = "properties"
	// TypeTable is the table the holds the type relation/edge.
	TypeTable = "properties"
	// TypeInverseTable is the table name for the PropertyType entity.
	// It exists in this package in order to avoid circular dependency with the "propertytype" package.
	TypeInverseTable = "property_types"
	// TypeColumn is the table column denoting the type relation/edge.
	TypeColumn = "property_type"
	// LocationTable is the table the holds the location relation/edge.
	LocationTable = "properties"
	// LocationInverseTable is the table name for the Location entity.
	// It exists in this package in order to avoid circular dependency with the "location" package.
	LocationInverseTable = "locations"
	// LocationColumn is the table column denoting the location relation/edge.
	LocationColumn = "location_properties"
	// EquipmentTable is the table the holds the equipment relation/edge.
	EquipmentTable = "properties"
	// EquipmentInverseTable is the table name for the Equipment entity.
	// It exists in this package in order to avoid circular dependency with the "equipment" package.
	EquipmentInverseTable = "equipment"
	// EquipmentColumn is the table column denoting the equipment relation/edge.
	EquipmentColumn = "equipment_properties"
	// ServiceTable is the table the holds the service relation/edge.
	ServiceTable = "properties"
	// ServiceInverseTable is the table name for the Service entity.
	// It exists in this package in order to avoid circular dependency with the "service" package.
	ServiceInverseTable = "services"
	// ServiceColumn is the table column denoting the service relation/edge.
	ServiceColumn = "service_properties"
	// EquipmentPortTable is the table the holds the equipment_port relation/edge.
	EquipmentPortTable = "properties"
	// EquipmentPortInverseTable is the table name for the EquipmentPort entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentport" package.
	EquipmentPortInverseTable = "equipment_ports"
	// EquipmentPortColumn is the table column denoting the equipment_port relation/edge.
	EquipmentPortColumn = "equipment_port_properties"
	// LinkTable is the table the holds the link relation/edge.
	LinkTable = "properties"
	// LinkInverseTable is the table name for the Link entity.
	// It exists in this package in order to avoid circular dependency with the "link" package.
	LinkInverseTable = "links"
	// LinkColumn is the table column denoting the link relation/edge.
	LinkColumn = "link_properties"
	// WorkOrderTable is the table the holds the work_order relation/edge.
	WorkOrderTable = "properties"
	// WorkOrderInverseTable is the table name for the WorkOrder entity.
	// It exists in this package in order to avoid circular dependency with the "workorder" package.
	WorkOrderInverseTable = "work_orders"
	// WorkOrderColumn is the table column denoting the work_order relation/edge.
	WorkOrderColumn = "work_order_properties"
	// ProjectTable is the table the holds the project relation/edge.
	ProjectTable = "properties"
	// ProjectInverseTable is the table name for the Project entity.
	// It exists in this package in order to avoid circular dependency with the "project" package.
	ProjectInverseTable = "projects"
	// ProjectColumn is the table column denoting the project relation/edge.
	ProjectColumn = "project_properties"
	// EquipmentValueTable is the table the holds the equipment_value relation/edge.
	EquipmentValueTable = "properties"
	// EquipmentValueInverseTable is the table name for the Equipment entity.
	// It exists in this package in order to avoid circular dependency with the "equipment" package.
	EquipmentValueInverseTable = "equipment"
	// EquipmentValueColumn is the table column denoting the equipment_value relation/edge.
	EquipmentValueColumn = "property_equipment_value"
	// LocationValueTable is the table the holds the location_value relation/edge.
	LocationValueTable = "properties"
	// LocationValueInverseTable is the table name for the Location entity.
	// It exists in this package in order to avoid circular dependency with the "location" package.
	LocationValueInverseTable = "locations"
	// LocationValueColumn is the table column denoting the location_value relation/edge.
	LocationValueColumn = "property_location_value"
	// ServiceValueTable is the table the holds the service_value relation/edge.
	ServiceValueTable = "properties"
	// ServiceValueInverseTable is the table name for the Service entity.
	// It exists in this package in order to avoid circular dependency with the "service" package.
	ServiceValueInverseTable = "services"
	// ServiceValueColumn is the table column denoting the service_value relation/edge.
	ServiceValueColumn = "property_service_value"
	// WorkOrderValueTable is the table the holds the work_order_value relation/edge.
	WorkOrderValueTable = "properties"
	// WorkOrderValueInverseTable is the table name for the WorkOrder entity.
	// It exists in this package in order to avoid circular dependency with the "workorder" package.
	WorkOrderValueInverseTable = "work_orders"
	// WorkOrderValueColumn is the table column denoting the work_order_value relation/edge.
	WorkOrderValueColumn = "property_work_order_value"
	// UserValueTable is the table the holds the user_value relation/edge.
	UserValueTable = "properties"
	// UserValueInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserValueInverseTable = "users"
	// UserValueColumn is the table column denoting the user_value relation/edge.
	UserValueColumn = "property_user_value"
)

// Columns holds all SQL columns for property fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldIntVal,
	FieldBoolVal,
	FieldFloatVal,
	FieldLatitudeVal,
	FieldLongitudeVal,
	FieldRangeFromVal,
	FieldRangeToVal,
	FieldStringVal,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the Property type.
var ForeignKeys = []string{
	"equipment_properties",
	"equipment_port_properties",
	"link_properties",
	"location_properties",
	"project_properties",
	"property_type",
	"property_equipment_value",
	"property_location_value",
	"property_service_value",
	"property_work_order_value",
	"property_user_value",
	"service_properties",
	"work_order_properties",
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
