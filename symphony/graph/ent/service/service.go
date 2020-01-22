// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package service

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the service type in the database.
	Label = "service"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldName holds the string denoting the name vertex property in the database.
	FieldName = "name"
	// FieldExternalID holds the string denoting the external_id vertex property in the database.
	FieldExternalID = "external_id"
	// FieldStatus holds the string denoting the status vertex property in the database.
	FieldStatus = "status"

	// Table holds the table name of the service in the database.
	Table = "services"
	// TypeTable is the table the holds the type relation/edge.
	TypeTable = "services"
	// TypeInverseTable is the table name for the ServiceType entity.
	// It exists in this package in order to avoid circular dependency with the "servicetype" package.
	TypeInverseTable = "service_types"
	// TypeColumn is the table column denoting the type relation/edge.
	TypeColumn = "type_id"
	// DownstreamTable is the table the holds the downstream relation/edge. The primary key declared below.
	DownstreamTable = "service_upstream"
	// UpstreamTable is the table the holds the upstream relation/edge. The primary key declared below.
	UpstreamTable = "service_upstream"
	// PropertiesTable is the table the holds the properties relation/edge.
	PropertiesTable = "properties"
	// PropertiesInverseTable is the table name for the Property entity.
	// It exists in this package in order to avoid circular dependency with the "property" package.
	PropertiesInverseTable = "properties"
	// PropertiesColumn is the table column denoting the properties relation/edge.
	PropertiesColumn = "service_id"
	// LinksTable is the table the holds the links relation/edge. The primary key declared below.
	LinksTable = "service_links"
	// LinksInverseTable is the table name for the Link entity.
	// It exists in this package in order to avoid circular dependency with the "link" package.
	LinksInverseTable = "links"
	// CustomerTable is the table the holds the customer relation/edge. The primary key declared below.
	CustomerTable = "service_customer"
	// CustomerInverseTable is the table name for the Customer entity.
	// It exists in this package in order to avoid circular dependency with the "customer" package.
	CustomerInverseTable = "customers"
	// EndpointsTable is the table the holds the endpoints relation/edge.
	EndpointsTable = "service_endpoints"
	// EndpointsInverseTable is the table name for the ServiceEndpoint entity.
	// It exists in this package in order to avoid circular dependency with the "serviceendpoint" package.
	EndpointsInverseTable = "service_endpoints"
	// EndpointsColumn is the table column denoting the endpoints relation/edge.
	EndpointsColumn = "service_id"
)

// Columns holds all SQL columns for service fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldExternalID,
	FieldStatus,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the Service type.
var ForeignKeys = []string{
	"type_id",
}

var (
	// DownstreamPrimaryKey and DownstreamColumn2 are the table columns denoting the
	// primary key for the downstream relation (M2M).
	DownstreamPrimaryKey = []string{"service_id", "downstream_id"}
	// UpstreamPrimaryKey and UpstreamColumn2 are the table columns denoting the
	// primary key for the upstream relation (M2M).
	UpstreamPrimaryKey = []string{"service_id", "downstream_id"}
	// LinksPrimaryKey and LinksColumn2 are the table columns denoting the
	// primary key for the links relation (M2M).
	LinksPrimaryKey = []string{"service_id", "link_id"}
	// CustomerPrimaryKey and CustomerColumn2 are the table columns denoting the
	// primary key for the customer relation (M2M).
	CustomerPrimaryKey = []string{"service_id", "customer_id"}
)

var (
	mixin       = schema.Service{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.Service{}.Fields()

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

	// descName is the schema descriptor for name field.
	descName = fields[0].Descriptor()
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator = descName.Validators[0].(func(string) error)

	// descExternalID is the schema descriptor for external_id field.
	descExternalID = fields[1].Descriptor()
	// ExternalIDValidator is a validator for the "external_id" field. It is called by the builders before save.
	ExternalIDValidator = descExternalID.Validators[0].(func(string) error)
)
