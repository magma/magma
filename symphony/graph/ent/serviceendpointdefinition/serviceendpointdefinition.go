// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package serviceendpointdefinition

import (
	"time"
)

const (
	// Label holds the string label denoting the serviceendpointdefinition type in the database.
	Label = "service_endpoint_definition"
	// FieldID holds the string denoting the id field in the database.
	FieldID         = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time" // FieldRole holds the string denoting the role vertex property in the database.
	FieldRole       = "role"        // FieldName holds the string denoting the name vertex property in the database.
	FieldName       = "name"        // FieldIndex holds the string denoting the index vertex property in the database.
	FieldIndex      = "index"

	// EdgeEndpoints holds the string denoting the endpoints edge name in mutations.
	EdgeEndpoints = "endpoints"
	// EdgeServiceType holds the string denoting the service_type edge name in mutations.
	EdgeServiceType = "service_type"
	// EdgeEquipmentType holds the string denoting the equipment_type edge name in mutations.
	EdgeEquipmentType = "equipment_type"

	// Table holds the table name of the serviceendpointdefinition in the database.
	Table = "service_endpoint_definitions"
	// EndpointsTable is the table the holds the endpoints relation/edge.
	EndpointsTable = "service_endpoints"
	// EndpointsInverseTable is the table name for the ServiceEndpoint entity.
	// It exists in this package in order to avoid circular dependency with the "serviceendpoint" package.
	EndpointsInverseTable = "service_endpoints"
	// EndpointsColumn is the table column denoting the endpoints relation/edge.
	EndpointsColumn = "service_endpoint_definition_endpoints"
	// ServiceTypeTable is the table the holds the service_type relation/edge.
	ServiceTypeTable = "service_endpoint_definitions"
	// ServiceTypeInverseTable is the table name for the ServiceType entity.
	// It exists in this package in order to avoid circular dependency with the "servicetype" package.
	ServiceTypeInverseTable = "service_types"
	// ServiceTypeColumn is the table column denoting the service_type relation/edge.
	ServiceTypeColumn = "service_type_endpoint_definitions"
	// EquipmentTypeTable is the table the holds the equipment_type relation/edge.
	EquipmentTypeTable = "service_endpoint_definitions"
	// EquipmentTypeInverseTable is the table name for the EquipmentType entity.
	// It exists in this package in order to avoid circular dependency with the "equipmenttype" package.
	EquipmentTypeInverseTable = "equipment_types"
	// EquipmentTypeColumn is the table column denoting the equipment_type relation/edge.
	EquipmentTypeColumn = "equipment_type_service_endpoint_definitions"
)

// Columns holds all SQL columns for serviceendpointdefinition fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldRole,
	FieldName,
	FieldIndex,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the ServiceEndpointDefinition type.
var ForeignKeys = []string{
	"equipment_type_service_endpoint_definitions",
	"service_type_endpoint_definitions",
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
)
