// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package tenant

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/frontier/ent/schema"
)

const (
	// Label holds the string label denoting the tenant type in the database.
	Label = "tenant"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at vertex property in the database.
	FieldCreatedAt = "createdAt"
	// FieldUpdatedAt holds the string denoting the updated_at vertex property in the database.
	FieldUpdatedAt = "updatedAt"
	// FieldName holds the string denoting the name vertex property in the database.
	FieldName = "name"
	// FieldDomains holds the string denoting the domains vertex property in the database.
	FieldDomains = "customDomains"
	// FieldNetworks holds the string denoting the networks vertex property in the database.
	FieldNetworks = "networkIDs"
	// FieldTabs holds the string denoting the tabs vertex property in the database.
	FieldTabs = "tabs"
	// FieldSSOCert holds the string denoting the ssocert vertex property in the database.
	FieldSSOCert = "ssoCert"
	// FieldSSOEntryPoint holds the string denoting the ssoentrypoint vertex property in the database.
	FieldSSOEntryPoint = "ssoEntrypoint"
	// FieldSSOIssuer holds the string denoting the ssoissuer vertex property in the database.
	FieldSSOIssuer = "ssoIssuer"

	// Table holds the table name of the tenant in the database.
	Table = "Organizations"
)

// Columns holds all SQL columns for tenant fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldName,
	FieldDomains,
	FieldNetworks,
	FieldTabs,
	FieldSSOCert,
	FieldSSOEntryPoint,
	FieldSSOIssuer,
}

var (
	mixin       = schema.Tenant{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.Tenant{}.Fields()

	// descCreatedAt is the schema descriptor for created_at field.
	descCreatedAt = mixinFields[0][0].Descriptor()
	// DefaultCreatedAt holds the default value on creation for the created_at field.
	DefaultCreatedAt = descCreatedAt.Default.(func() time.Time)

	// descUpdatedAt is the schema descriptor for updated_at field.
	descUpdatedAt = mixinFields[0][1].Descriptor()
	// DefaultUpdatedAt holds the default value on creation for the updated_at field.
	DefaultUpdatedAt = descUpdatedAt.Default.(func() time.Time)
	// UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	UpdateDefaultUpdatedAt = descUpdatedAt.UpdateDefault.(func() time.Time)

	// descName is the schema descriptor for name field.
	descName = fields[0].Descriptor()
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator = descName.Validators[0].(func(string) error)

	// descSSOCert is the schema descriptor for SSOCert field.
	descSSOCert = fields[4].Descriptor()
	// DefaultSSOCert holds the default value on creation for the SSOCert field.
	DefaultSSOCert = descSSOCert.Default.(string)

	// descSSOEntryPoint is the schema descriptor for SSOEntryPoint field.
	descSSOEntryPoint = fields[5].Descriptor()
	// DefaultSSOEntryPoint holds the default value on creation for the SSOEntryPoint field.
	DefaultSSOEntryPoint = descSSOEntryPoint.Default.(string)

	// descSSOIssuer is the schema descriptor for SSOIssuer field.
	descSSOIssuer = fields[6].Descriptor()
	// DefaultSSOIssuer holds the default value on creation for the SSOIssuer field.
	DefaultSSOIssuer = descSSOIssuer.Default.(string)
)
