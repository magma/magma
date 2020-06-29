// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package propertytype

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the propertytype type in the database.
	Label = "property_type"
	// FieldID holds the string denoting the id field in the database.
	FieldID                 = "id"                   // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime         = "create_time"          // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime         = "update_time"          // FieldType holds the string denoting the type vertex property in the database.
	FieldType               = "type"                 // FieldName holds the string denoting the name vertex property in the database.
	FieldName               = "name"                 // FieldExternalID holds the string denoting the external_id vertex property in the database.
	FieldExternalID         = "external_id"          // FieldIndex holds the string denoting the index vertex property in the database.
	FieldIndex              = "index"                // FieldCategory holds the string denoting the category vertex property in the database.
	FieldCategory           = "category"             // FieldIntVal holds the string denoting the int_val vertex property in the database.
	FieldIntVal             = "int_val"              // FieldBoolVal holds the string denoting the bool_val vertex property in the database.
	FieldBoolVal            = "bool_val"             // FieldFloatVal holds the string denoting the float_val vertex property in the database.
	FieldFloatVal           = "float_val"            // FieldLatitudeVal holds the string denoting the latitude_val vertex property in the database.
	FieldLatitudeVal        = "latitude_val"         // FieldLongitudeVal holds the string denoting the longitude_val vertex property in the database.
	FieldLongitudeVal       = "longitude_val"        // FieldStringVal holds the string denoting the string_val vertex property in the database.
	FieldStringVal          = "string_val"           // FieldRangeFromVal holds the string denoting the range_from_val vertex property in the database.
	FieldRangeFromVal       = "range_from_val"       // FieldRangeToVal holds the string denoting the range_to_val vertex property in the database.
	FieldRangeToVal         = "range_to_val"         // FieldIsInstanceProperty holds the string denoting the is_instance_property vertex property in the database.
	FieldIsInstanceProperty = "is_instance_property" // FieldEditable holds the string denoting the editable vertex property in the database.
	FieldEditable           = "editable"             // FieldMandatory holds the string denoting the mandatory vertex property in the database.
	FieldMandatory          = "mandatory"            // FieldDeleted holds the string denoting the deleted vertex property in the database.
	FieldDeleted            = "deleted"              // FieldNodeType holds the string denoting the nodetype vertex property in the database.
	FieldNodeType           = "node_type"

	// EdgeProperties holds the string denoting the properties edge name in mutations.
	EdgeProperties = "properties"
	// EdgeLocationType holds the string denoting the location_type edge name in mutations.
	EdgeLocationType = "location_type"
	// EdgeEquipmentPortType holds the string denoting the equipment_port_type edge name in mutations.
	EdgeEquipmentPortType = "equipment_port_type"
	// EdgeLinkEquipmentPortType holds the string denoting the link_equipment_port_type edge name in mutations.
	EdgeLinkEquipmentPortType = "link_equipment_port_type"
	// EdgeEquipmentType holds the string denoting the equipment_type edge name in mutations.
	EdgeEquipmentType = "equipment_type"
	// EdgeServiceType holds the string denoting the service_type edge name in mutations.
	EdgeServiceType = "service_type"
	// EdgeWorkOrderType holds the string denoting the work_order_type edge name in mutations.
	EdgeWorkOrderType = "work_order_type"
	// EdgeWorkOrderTemplate holds the string denoting the work_order_template edge name in mutations.
	EdgeWorkOrderTemplate = "work_order_template"
	// EdgeProjectType holds the string denoting the project_type edge name in mutations.
	EdgeProjectType = "project_type"

	// Table holds the table name of the propertytype in the database.
	Table = "property_types"
	// PropertiesTable is the table the holds the properties relation/edge.
	PropertiesTable = "properties"
	// PropertiesInverseTable is the table name for the Property entity.
	// It exists in this package in order to avoid circular dependency with the "property" package.
	PropertiesInverseTable = "properties"
	// PropertiesColumn is the table column denoting the properties relation/edge.
	PropertiesColumn = "property_type"
	// LocationTypeTable is the table the holds the location_type relation/edge.
	LocationTypeTable = "property_types"
	// LocationTypeInverseTable is the table name for the LocationType entity.
	// It exists in this package in order to avoid circular dependency with the "locationtype" package.
	LocationTypeInverseTable = "location_types"
	// LocationTypeColumn is the table column denoting the location_type relation/edge.
	LocationTypeColumn = "location_type_property_types"
	// EquipmentPortTypeTable is the table the holds the equipment_port_type relation/edge.
	EquipmentPortTypeTable = "property_types"
	// EquipmentPortTypeInverseTable is the table name for the EquipmentPortType entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentporttype" package.
	EquipmentPortTypeInverseTable = "equipment_port_types"
	// EquipmentPortTypeColumn is the table column denoting the equipment_port_type relation/edge.
	EquipmentPortTypeColumn = "equipment_port_type_property_types"
	// LinkEquipmentPortTypeTable is the table the holds the link_equipment_port_type relation/edge.
	LinkEquipmentPortTypeTable = "property_types"
	// LinkEquipmentPortTypeInverseTable is the table name for the EquipmentPortType entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentporttype" package.
	LinkEquipmentPortTypeInverseTable = "equipment_port_types"
	// LinkEquipmentPortTypeColumn is the table column denoting the link_equipment_port_type relation/edge.
	LinkEquipmentPortTypeColumn = "equipment_port_type_link_property_types"
	// EquipmentTypeTable is the table the holds the equipment_type relation/edge.
	EquipmentTypeTable = "property_types"
	// EquipmentTypeInverseTable is the table name for the EquipmentType entity.
	// It exists in this package in order to avoid circular dependency with the "equipmenttype" package.
	EquipmentTypeInverseTable = "equipment_types"
	// EquipmentTypeColumn is the table column denoting the equipment_type relation/edge.
	EquipmentTypeColumn = "equipment_type_property_types"
	// ServiceTypeTable is the table the holds the service_type relation/edge.
	ServiceTypeTable = "property_types"
	// ServiceTypeInverseTable is the table name for the ServiceType entity.
	// It exists in this package in order to avoid circular dependency with the "servicetype" package.
	ServiceTypeInverseTable = "service_types"
	// ServiceTypeColumn is the table column denoting the service_type relation/edge.
	ServiceTypeColumn = "service_type_property_types"
	// WorkOrderTypeTable is the table the holds the work_order_type relation/edge.
	WorkOrderTypeTable = "property_types"
	// WorkOrderTypeInverseTable is the table name for the WorkOrderType entity.
	// It exists in this package in order to avoid circular dependency with the "workordertype" package.
	WorkOrderTypeInverseTable = "work_order_types"
	// WorkOrderTypeColumn is the table column denoting the work_order_type relation/edge.
	WorkOrderTypeColumn = "work_order_type_property_types"
	// WorkOrderTemplateTable is the table the holds the work_order_template relation/edge.
	WorkOrderTemplateTable = "property_types"
	// WorkOrderTemplateInverseTable is the table name for the WorkOrderTemplate entity.
	// It exists in this package in order to avoid circular dependency with the "workordertemplate" package.
	WorkOrderTemplateInverseTable = "work_order_templates"
	// WorkOrderTemplateColumn is the table column denoting the work_order_template relation/edge.
	WorkOrderTemplateColumn = "work_order_template_property_types"
	// ProjectTypeTable is the table the holds the project_type relation/edge.
	ProjectTypeTable = "property_types"
	// ProjectTypeInverseTable is the table name for the ProjectType entity.
	// It exists in this package in order to avoid circular dependency with the "projecttype" package.
	ProjectTypeInverseTable = "project_types"
	// ProjectTypeColumn is the table column denoting the project_type relation/edge.
	ProjectTypeColumn = "project_type_properties"
)

// Columns holds all SQL columns for propertytype fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldType,
	FieldName,
	FieldExternalID,
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
	FieldNodeType,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the PropertyType type.
var ForeignKeys = []string{
	"equipment_port_type_property_types",
	"equipment_port_type_link_property_types",
	"equipment_type_property_types",
	"location_type_property_types",
	"project_type_properties",
	"service_type_property_types",
	"work_order_template_property_types",
	"work_order_type_property_types",
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
	// DefaultIsInstanceProperty holds the default value on creation for the is_instance_property field.
	DefaultIsInstanceProperty bool
	// DefaultEditable holds the default value on creation for the editable field.
	DefaultEditable bool
	// DefaultMandatory holds the default value on creation for the mandatory field.
	DefaultMandatory bool
	// DefaultDeleted holds the default value on creation for the deleted field.
	DefaultDeleted bool
)
