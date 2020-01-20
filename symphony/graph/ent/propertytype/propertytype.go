// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package propertytype

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the propertytype type in the database.
	Label = "property_type"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldType holds the string denoting the type vertex property in the database.
	FieldType = "type"
	// FieldName holds the string denoting the name vertex property in the database.
	FieldName = "name"
	// FieldIndex holds the string denoting the index vertex property in the database.
	FieldIndex = "index"
	// FieldCategory holds the string denoting the category vertex property in the database.
	FieldCategory = "category"
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
	// FieldStringVal holds the string denoting the string_val vertex property in the database.
	FieldStringVal = "string_val"
	// FieldRangeFromVal holds the string denoting the range_from_val vertex property in the database.
	FieldRangeFromVal = "range_from_val"
	// FieldRangeToVal holds the string denoting the range_to_val vertex property in the database.
	FieldRangeToVal = "range_to_val"
	// FieldIsInstanceProperty holds the string denoting the is_instance_property vertex property in the database.
	FieldIsInstanceProperty = "is_instance_property"
	// FieldEditable holds the string denoting the editable vertex property in the database.
	FieldEditable = "editable"
	// FieldMandatory holds the string denoting the mandatory vertex property in the database.
	FieldMandatory = "mandatory"
	// FieldDeleted holds the string denoting the deleted vertex property in the database.
	FieldDeleted = "deleted"

	// Table holds the table name of the propertytype in the database.
	Table = "property_types"
	// PropertiesTable is the table the holds the properties relation/edge.
	PropertiesTable = "properties"
	// PropertiesInverseTable is the table name for the Property entity.
	// It exists in this package in order to avoid circular dependency with the "property" package.
	PropertiesInverseTable = "properties"
	// PropertiesColumn is the table column denoting the properties relation/edge.
	PropertiesColumn = "type_id"
	// LocationTypeTable is the table the holds the location_type relation/edge.
	LocationTypeTable = "property_types"
	// LocationTypeInverseTable is the table name for the LocationType entity.
	// It exists in this package in order to avoid circular dependency with the "locationtype" package.
	LocationTypeInverseTable = "location_types"
	// LocationTypeColumn is the table column denoting the location_type relation/edge.
	LocationTypeColumn = "location_type_id"
	// EquipmentPortTypeTable is the table the holds the equipment_port_type relation/edge.
	EquipmentPortTypeTable = "property_types"
	// EquipmentPortTypeInverseTable is the table name for the EquipmentPortType entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentporttype" package.
	EquipmentPortTypeInverseTable = "equipment_port_types"
	// EquipmentPortTypeColumn is the table column denoting the equipment_port_type relation/edge.
	EquipmentPortTypeColumn = "equipment_port_type_id"
	// LinkEquipmentPortTypeTable is the table the holds the link_equipment_port_type relation/edge.
	LinkEquipmentPortTypeTable = "property_types"
	// LinkEquipmentPortTypeInverseTable is the table name for the EquipmentPortType entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentporttype" package.
	LinkEquipmentPortTypeInverseTable = "equipment_port_types"
	// LinkEquipmentPortTypeColumn is the table column denoting the link_equipment_port_type relation/edge.
	LinkEquipmentPortTypeColumn = "link_equipment_port_type_id"
	// EquipmentTypeTable is the table the holds the equipment_type relation/edge.
	EquipmentTypeTable = "property_types"
	// EquipmentTypeInverseTable is the table name for the EquipmentType entity.
	// It exists in this package in order to avoid circular dependency with the "equipmenttype" package.
	EquipmentTypeInverseTable = "equipment_types"
	// EquipmentTypeColumn is the table column denoting the equipment_type relation/edge.
	EquipmentTypeColumn = "equipment_type_id"
	// ServiceTypeTable is the table the holds the service_type relation/edge.
	ServiceTypeTable = "property_types"
	// ServiceTypeInverseTable is the table name for the ServiceType entity.
	// It exists in this package in order to avoid circular dependency with the "servicetype" package.
	ServiceTypeInverseTable = "service_types"
	// ServiceTypeColumn is the table column denoting the service_type relation/edge.
	ServiceTypeColumn = "service_type_id"
	// WorkOrderTypeTable is the table the holds the work_order_type relation/edge.
	WorkOrderTypeTable = "property_types"
	// WorkOrderTypeInverseTable is the table name for the WorkOrderType entity.
	// It exists in this package in order to avoid circular dependency with the "workordertype" package.
	WorkOrderTypeInverseTable = "work_order_types"
	// WorkOrderTypeColumn is the table column denoting the work_order_type relation/edge.
	WorkOrderTypeColumn = "work_order_type_id"
	// ProjectTypeTable is the table the holds the project_type relation/edge.
	ProjectTypeTable = "property_types"
	// ProjectTypeInverseTable is the table name for the ProjectType entity.
	// It exists in this package in order to avoid circular dependency with the "projecttype" package.
	ProjectTypeInverseTable = "project_types"
	// ProjectTypeColumn is the table column denoting the project_type relation/edge.
	ProjectTypeColumn = "project_type_id"
)

// Columns holds all SQL columns are propertytype fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldType,
	FieldName,
	FieldIndex,
	FieldCategory,
	FieldIntVal,
	FieldBoolVal,
	FieldFloatVal,
	FieldLatitudeVal,
	FieldLongitudeVal,
	FieldStringVal,
	FieldRangeFromVal,
	FieldRangeToVal,
	FieldIsInstanceProperty,
	FieldEditable,
	FieldMandatory,
	FieldDeleted,
}

var (
	mixin       = schema.PropertyType{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.PropertyType{}.Fields()

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

	// descIsInstanceProperty is the schema descriptor for is_instance_property field.
	descIsInstanceProperty = fields[12].Descriptor()
	// DefaultIsInstanceProperty holds the default value on creation for the is_instance_property field.
	DefaultIsInstanceProperty = descIsInstanceProperty.Default.(bool)

	// descEditable is the schema descriptor for editable field.
	descEditable = fields[13].Descriptor()
	// DefaultEditable holds the default value on creation for the editable field.
	DefaultEditable = descEditable.Default.(bool)

	// descMandatory is the schema descriptor for mandatory field.
	descMandatory = fields[14].Descriptor()
	// DefaultMandatory holds the default value on creation for the mandatory field.
	DefaultMandatory = descMandatory.Default.(bool)

	// descDeleted is the schema descriptor for deleted field.
	descDeleted = fields[15].Descriptor()
	// DefaultDeleted holds the default value on creation for the deleted field.
	DefaultDeleted = descDeleted.Default.(bool)
)
