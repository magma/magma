// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package equipmentporttype

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the equipmentporttype type in the database.
	Label = "equipment_port_type"
	// FieldID holds the string denoting the id field in the database.
	FieldID         = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time" // FieldName holds the string denoting the name vertex property in the database.
	FieldName       = "name"

	// EdgePropertyTypes holds the string denoting the property_types edge name in mutations.
	EdgePropertyTypes = "property_types"
	// EdgeLinkPropertyTypes holds the string denoting the link_property_types edge name in mutations.
	EdgeLinkPropertyTypes = "link_property_types"
	// EdgePortDefinitions holds the string denoting the port_definitions edge name in mutations.
	EdgePortDefinitions = "port_definitions"

	// Table holds the table name of the equipmentporttype in the database.
	Table = "equipment_port_types"
	// PropertyTypesTable is the table the holds the property_types relation/edge.
	PropertyTypesTable = "property_types"
	// PropertyTypesInverseTable is the table name for the PropertyType entity.
	// It exists in this package in order to avoid circular dependency with the "propertytype" package.
	PropertyTypesInverseTable = "property_types"
	// PropertyTypesColumn is the table column denoting the property_types relation/edge.
	PropertyTypesColumn = "equipment_port_type_property_types"
	// LinkPropertyTypesTable is the table the holds the link_property_types relation/edge.
	LinkPropertyTypesTable = "property_types"
	// LinkPropertyTypesInverseTable is the table name for the PropertyType entity.
	// It exists in this package in order to avoid circular dependency with the "propertytype" package.
	LinkPropertyTypesInverseTable = "property_types"
	// LinkPropertyTypesColumn is the table column denoting the link_property_types relation/edge.
	LinkPropertyTypesColumn = "equipment_port_type_link_property_types"
	// PortDefinitionsTable is the table the holds the port_definitions relation/edge.
	PortDefinitionsTable = "equipment_port_definitions"
	// PortDefinitionsInverseTable is the table name for the EquipmentPortDefinition entity.
	// It exists in this package in order to avoid circular dependency with the "equipmentportdefinition" package.
	PortDefinitionsInverseTable = "equipment_port_definitions"
	// PortDefinitionsColumn is the table column denoting the port_definitions relation/edge.
	PortDefinitionsColumn = "equipment_port_definition_equipment_port_type"
)

// Columns holds all SQL columns for equipmentporttype fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
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
