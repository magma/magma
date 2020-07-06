// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package servicetype

import (
	"fmt"
	"io"
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the servicetype type in the database.
	Label = "service_type"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time field in the database.
	FieldUpdateTime = "update_time"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldHasCustomer holds the string denoting the has_customer field in the database.
	FieldHasCustomer = "has_customer"
	// FieldIsDeleted holds the string denoting the is_deleted field in the database.
	FieldIsDeleted = "is_deleted"
	// FieldDiscoveryMethod holds the string denoting the discovery_method field in the database.
	FieldDiscoveryMethod = "discovery_method"

	// EdgeServices holds the string denoting the services edge name in mutations.
	EdgeServices = "services"
	// EdgePropertyTypes holds the string denoting the property_types edge name in mutations.
	EdgePropertyTypes = "property_types"
	// EdgeEndpointDefinitions holds the string denoting the endpoint_definitions edge name in mutations.
	EdgeEndpointDefinitions = "endpoint_definitions"

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
	// EndpointDefinitionsTable is the table the holds the endpoint_definitions relation/edge.
	EndpointDefinitionsTable = "service_endpoint_definitions"
	// EndpointDefinitionsInverseTable is the table name for the ServiceEndpointDefinition entity.
	// It exists in this package in order to avoid circular dependency with the "serviceendpointdefinition" package.
	EndpointDefinitionsInverseTable = "service_endpoint_definitions"
	// EndpointDefinitionsColumn is the table column denoting the endpoint_definitions relation/edge.
	EndpointDefinitionsColumn = "service_type_endpoint_definitions"
)

// Columns holds all SQL columns for servicetype fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldHasCustomer,
	FieldIsDeleted,
	FieldDiscoveryMethod,
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
	// DefaultHasCustomer holds the default value on creation for the has_customer field.
	DefaultHasCustomer bool
	// DefaultIsDeleted holds the default value on creation for the is_deleted field.
	DefaultIsDeleted bool
)

// DiscoveryMethod defines the type for the discovery_method enum field.
type DiscoveryMethod string

// DiscoveryMethod values.
const (
	DiscoveryMethodINVENTORY DiscoveryMethod = "INVENTORY"
)

func (dm DiscoveryMethod) String() string {
	return string(dm)
}

// DiscoveryMethodValidator is a validator for the "dm" field enum values. It is called by the builders before save.
func DiscoveryMethodValidator(dm DiscoveryMethod) error {
	switch dm {
	case DiscoveryMethodINVENTORY:
		return nil
	default:
		return fmt.Errorf("servicetype: invalid enum value for discovery_method field: %q", dm)
	}
}

// MarshalGQL implements graphql.Marshaler interface.
func (dm DiscoveryMethod) MarshalGQL(w io.Writer) {
	writeQuotedStringer(w, dm)
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (dm *DiscoveryMethod) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enum %T must be a string", v)
	}
	*dm = DiscoveryMethod(str)
	if err := DiscoveryMethodValidator(*dm); err != nil {
		return fmt.Errorf("%s is not a valid DiscoveryMethod", str)
	}
	return nil
}

func writeQuotedStringer(w io.Writer, s fmt.Stringer) {
	const quote = '"'
	switch w := w.(type) {
	case io.ByteWriter:
		w.WriteByte(quote)
		defer w.WriteByte(quote)
	default:
		w.Write([]byte{quote})
		defer w.Write([]byte{quote})
	}
	io.WriteString(w, s.String())
}
