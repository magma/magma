// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package servicetype

import (
	"time"
)

const (
	// Label holds the string label denoting the servicetype type in the database.
	Label = "service_type"
	// FieldID holds the string denoting the id field in the database.
	FieldID          = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime  = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime  = "update_time" // FieldName holds the string denoting the name vertex property in the database.
	FieldName        = "name"        // FieldHasCustomer holds the string denoting the has_customer vertex property in the database.
	FieldHasCustomer = "has_customer"

	// EdgeServices holds the string denoting the services edge name in mutations.
	EdgeServices = "services"
	// EdgePropertyTypes holds the string denoting the property_types edge name in mutations.
	EdgePropertyTypes = "property_types"

	// Table holds the table name of the servicetype in the database.
	Table = "service_types"
	// ServicesTable is the table the holds the services relation/edge.
	ServicesTable = "services"
	// ServicesInverseTable is the table name for the Service entity.
	// It exists in this package in order to avoid circular dependency with the "service" package.
	ServicesInverseTable = "services"
	// ServicesColumn is the table column denoting the services relation/edge.
	ServicesColumn = "service_type"
	// PropertyTypesTable is the table the holds the property_types relation/edge.
	PropertyTypesTable = "property_types"
	// PropertyTypesInverseTable is the table name for the PropertyType entity.
	// It exists in this package in order to avoid circular dependency with the "propertytype" package.
	PropertyTypesInverseTable = "property_types"
	// PropertyTypesColumn is the table column denoting the property_types relation/edge.
	PropertyTypesColumn = "service_type_property_types"
)

// Columns holds all SQL columns for servicetype fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldHasCustomer,
}

var (
	// DefaultCreateTime holds the default value on creation for the create_time field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the update_time field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	UpdateDefaultUpdateTime func() time.Time
	// DefaultHasCustomer holds the default value on creation for the has_customer field.
	DefaultHasCustomer bool
)
