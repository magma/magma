// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package equipmentport

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the equipmentport type in the database.
	Label = "equipment_port"
	// FieldID holds the string denoting the id field in the database.
	FieldID         = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"

	// EdgeDefinition holds the string denoting the definition edge name in mutations.
	EdgeDefinition = "definition"
	// EdgeParent holds the string denoting the parent edge name in mutations.
	EdgeParent = "parent"
	// EdgeLink holds the string denoting the link edge name in mutations.
	EdgeLink = "link"
	// EdgeProperties holds the string denoting the properties edge name in mutations.
	EdgeProperties = "properties"
	// EdgeEndpoints holds the string denoting the endpoints edge name in mutations.
	EdgeEndpoints = "endpoints"

	// Table holds the table name of the equipmentport in the database.
	Table = "equipment_ports"
	// DefinitionTable is the table the holds the definition relation/edge.
	DefinitionTable = "equipment_ports"
	// DefinitionInverseTable is the table name for the EquipmentPortDefinition entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentportdefinition" package.
	DefinitionInverseTable = "equipment_port_definitions"
	// DefinitionColumn is the table column denoting the definition relation/edge.
	DefinitionColumn = "equipment_port_definition"
	// ParentTable is the table the holds the parent relation/edge.
	ParentTable = "equipment_ports"
	// ParentInverseTable is the table name for the Equipment entity.
	// It exists in this package in order to avoid circular dependency with the "equipment" package.
	ParentInverseTable = "equipment"
	// ParentColumn is the table column denoting the parent relation/edge.
	ParentColumn = "equipment_ports"
	// LinkTable is the table the holds the link relation/edge.
	LinkTable = "equipment_ports"
	// LinkInverseTable is the table name for the Link entity.
	// It exists in this package in order to avoid circular dependency with the "link" package.
	LinkInverseTable = "links"
	// LinkColumn is the table column denoting the link relation/edge.
	LinkColumn = "equipment_port_link"
	// PropertiesTable is the table the holds the properties relation/edge.
	PropertiesTable = "properties"
	// PropertiesInverseTable is the table name for the Property entity.
	// It exists in this package in order to avoid circular dependency with the "property" package.
	PropertiesInverseTable = "properties"
	// PropertiesColumn is the table column denoting the properties relation/edge.
	PropertiesColumn = "equipment_port_properties"
	// EndpointsTable is the table the holds the endpoints relation/edge.
	EndpointsTable = "service_endpoints"
	// EndpointsInverseTable is the table name for the ServiceEndpoint entity.
	// It exists in this package in order to avoid circular dependency with the "serviceendpoint" package.
	EndpointsInverseTable = "service_endpoints"
	// EndpointsColumn is the table column denoting the endpoints relation/edge.
	EndpointsColumn = "service_endpoint_port"
)

// Columns holds all SQL columns for equipmentport fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the EquipmentPort type.
var ForeignKeys = []string{
	"equipment_ports",
	"equipment_port_definition",
	"equipment_port_link",
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
