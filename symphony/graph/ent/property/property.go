// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package property

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the property type in the database.
	Label = "property"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldIntVal holds the string denoting the int_val vertex property in the database.
	FieldIntVal = "int_val"
	// FieldBoolVal holds the string denoting the bool_val vertex property in the database.
	FieldBoolVal = "bool_val"
	// FieldFloatVal holds the string denoting the float_val vertex property in the database.
	FieldFloatVal = "float_val"
	// FieldLatitudeVal holds the string denoting the latitude_val vertex property in the database.
	FieldLatitudeVal = "latitude_val"
	// FieldLongitudeVal holds the string denoting the longitude_val vertex property in the database.
	FieldLongitudeVal = "longitude_val"
	// FieldRangeFromVal holds the string denoting the range_from_val vertex property in the database.
	FieldRangeFromVal = "range_from_val"
	// FieldRangeToVal holds the string denoting the range_to_val vertex property in the database.
	FieldRangeToVal = "range_to_val"
	// FieldStringVal holds the string denoting the string_val vertex property in the database.
	FieldStringVal = "string_val"

	// Table holds the table name of the property in the database.
	Table = "properties"
	// TypeTable is the table the holds the type relation/edge.
	TypeTable = "properties"
	// TypeInverseTable is the table name for the PropertyType entity.
	// It exists in this package in order to avoid circular dependency with the "propertytype" package.
	TypeInverseTable = "property_types"
	// TypeColumn is the table column denoting the type relation/edge.
	TypeColumn = "type_id"
	// LocationTable is the table the holds the location relation/edge.
	LocationTable = "properties"
	// LocationInverseTable is the table name for the Location entity.
	// It exists in this package in order to avoid circular dependency with the "location" package.
	LocationInverseTable = "locations"
	// LocationColumn is the table column denoting the location relation/edge.
	LocationColumn = "location_id"
	// EquipmentTable is the table the holds the equipment relation/edge.
	EquipmentTable = "properties"
	// EquipmentInverseTable is the table name for the Equipment entity.
	// It exists in this package in order to avoid circular dependency with the "equipment" package.
	EquipmentInverseTable = "equipment"
	// EquipmentColumn is the table column denoting the equipment relation/edge.
	EquipmentColumn = "equipment_id"
	// ServiceTable is the table the holds the service relation/edge.
	ServiceTable = "properties"
	// ServiceInverseTable is the table name for the Service entity.
	// It exists in this package in order to avoid circular dependency with the "service" package.
	ServiceInverseTable = "services"
	// ServiceColumn is the table column denoting the service relation/edge.
	ServiceColumn = "service_id"
	// EquipmentPortTable is the table the holds the equipment_port relation/edge.
	EquipmentPortTable = "properties"
	// EquipmentPortInverseTable is the table name for the EquipmentPort entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentport" package.
	EquipmentPortInverseTable = "equipment_ports"
	// EquipmentPortColumn is the table column denoting the equipment_port relation/edge.
	EquipmentPortColumn = "equipment_port_id"
	// LinkTable is the table the holds the link relation/edge.
	LinkTable = "properties"
	// LinkInverseTable is the table name for the Link entity.
	// It exists in this package in order to avoid circular dependency with the "link" package.
	LinkInverseTable = "links"
	// LinkColumn is the table column denoting the link relation/edge.
	LinkColumn = "link_id"
	// WorkOrderTable is the table the holds the work_order relation/edge.
	WorkOrderTable = "properties"
	// WorkOrderInverseTable is the table name for the WorkOrder entity.
	// It exists in this package in order to avoid circular dependency with the "workorder" package.
	WorkOrderInverseTable = "work_orders"
	// WorkOrderColumn is the table column denoting the work_order relation/edge.
	WorkOrderColumn = "work_order_id"
	// ProjectTable is the table the holds the project relation/edge.
	ProjectTable = "properties"
	// ProjectInverseTable is the table name for the Project entity.
	// It exists in this package in order to avoid circular dependency with the "project" package.
	ProjectInverseTable = "projects"
	// ProjectColumn is the table column denoting the project relation/edge.
	ProjectColumn = "project_id"
	// EquipmentValueTable is the table the holds the equipment_value relation/edge.
	EquipmentValueTable = "properties"
	// EquipmentValueInverseTable is the table name for the Equipment entity.
	// It exists in this package in order to avoid circular dependency with the "equipment" package.
	EquipmentValueInverseTable = "equipment"
	// EquipmentValueColumn is the table column denoting the equipment_value relation/edge.
	EquipmentValueColumn = "property_equipment_value_id"
	// LocationValueTable is the table the holds the location_value relation/edge.
	LocationValueTable = "properties"
	// LocationValueInverseTable is the table name for the Location entity.
	// It exists in this package in order to avoid circular dependency with the "location" package.
	LocationValueInverseTable = "locations"
	// LocationValueColumn is the table column denoting the location_value relation/edge.
	LocationValueColumn = "property_location_value_id"
	// ServiceValueTable is the table the holds the service_value relation/edge.
	ServiceValueTable = "properties"
	// ServiceValueInverseTable is the table name for the Service entity.
	// It exists in this package in order to avoid circular dependency with the "service" package.
	ServiceValueInverseTable = "services"
	// ServiceValueColumn is the table column denoting the service_value relation/edge.
	ServiceValueColumn = "property_service_value_id"
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
	"equipment_id",
	"equipment_port_id",
	"link_id",
	"location_id",
	"project_id",
	"type_id",
	"property_equipment_value_id",
	"property_location_value_id",
	"property_service_value_id",
	"service_id",
	"work_order_id",
}

var (
	mixin       = schema.Property{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.Property{}.Fields()

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
