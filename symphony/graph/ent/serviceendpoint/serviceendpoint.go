// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package serviceendpoint

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the serviceendpoint type in the database.
	Label = "service_endpoint"
	// FieldID holds the string denoting the id field in the database.
	FieldID         = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"

	// EdgePort holds the string denoting the port edge name in mutations.
	EdgePort = "port"
	// EdgeEquipment holds the string denoting the equipment edge name in mutations.
	EdgeEquipment = "equipment"
	// EdgeService holds the string denoting the service edge name in mutations.
	EdgeService = "service"
	// EdgeDefinition holds the string denoting the definition edge name in mutations.
	EdgeDefinition = "definition"

	// Table holds the table name of the serviceendpoint in the database.
	Table = "service_endpoints"
	// PortTable is the table the holds the port relation/edge.
	PortTable = "service_endpoints"
	// PortInverseTable is the table name for the EquipmentPort entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentport" package.
	PortInverseTable = "equipment_ports"
	// PortColumn is the table column denoting the port relation/edge.
	PortColumn = "service_endpoint_port"
	// EquipmentTable is the table the holds the equipment relation/edge.
	EquipmentTable = "service_endpoints"
	// EquipmentInverseTable is the table name for the Equipment entity.
	// It exists in this package in order to avoid circular dependency with the "equipment" package.
	EquipmentInverseTable = "equipment"
	// EquipmentColumn is the table column denoting the equipment relation/edge.
	EquipmentColumn = "service_endpoint_equipment"
	// ServiceTable is the table the holds the service relation/edge.
	ServiceTable = "service_endpoints"
	// ServiceInverseTable is the table name for the Service entity.
	// It exists in this package in order to avoid circular dependency with the "service" package.
	ServiceInverseTable = "services"
	// ServiceColumn is the table column denoting the service relation/edge.
	ServiceColumn = "service_endpoints"
	// DefinitionTable is the table the holds the definition relation/edge.
	DefinitionTable = "service_endpoints"
	// DefinitionInverseTable is the table name for the ServiceEndpointDefinition entity.
	// It exists in this package in order to avoid circular dependency with the "serviceendpointdefinition" package.
	DefinitionInverseTable = "service_endpoint_definitions"
	// DefinitionColumn is the table column denoting the definition relation/edge.
	DefinitionColumn = "service_endpoint_definition_endpoints"
)

// Columns holds all SQL columns for serviceendpoint fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the ServiceEndpoint type.
var ForeignKeys = []string{
	"service_endpoints",
	"service_endpoint_port",
	"service_endpoint_equipment",
	"service_endpoint_definition_endpoints",
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/facebookincubator/symphony/graph/ent/runtime"
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
