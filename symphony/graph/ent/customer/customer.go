// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package customer

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the customer type in the database.
	Label = "customer"
	// FieldID holds the string denoting the id field in the database.
	FieldID         = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time" // FieldName holds the string denoting the name vertex property in the database.
	FieldName       = "name"        // FieldExternalID holds the string denoting the external_id vertex property in the database.
	FieldExternalID = "external_id"

	// EdgeServices holds the string denoting the services edge name in mutations.
	EdgeServices = "services"

	// Table holds the table name of the customer in the database.
	Table = "customers"
	// ServicesTable is the table the holds the services relation/edge. The primary key declared below.
	ServicesTable = "service_customer"
	// ServicesInverseTable is the table name for the Service entity.
	// It exists in this package in order to avoid circular dependency with the "service" package.
	ServicesInverseTable = "services"
)

// Columns holds all SQL columns for customer fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldExternalID,
}

var (
	// ServicesPrimaryKey and ServicesColumn2 are the table columns denoting the
	// primary key for the services relation (M2M).
	ServicesPrimaryKey = []string{"service_id", "customer_id"}
)

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
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// ExternalIDValidator is a validator for the "external_id" field. It is called by the builders before save.
	ExternalIDValidator func(string) error
)
