// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/facebookincubator/symphony/frontier/ent/auditlog"
	"github.com/facebookincubator/symphony/frontier/ent/schema"

	"github.com/facebookincubator/symphony/frontier/ent/tenant"

	"github.com/facebookincubator/symphony/frontier/ent/token"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/frontier/ent/user"
)

// The init function reads all schema descriptors with runtime
// code (default values, validators or hooks) and stitches it
// to their package variables.
func init() {
	auditlogMixin := schema.AuditLog{}.Mixin()
	auditlogMixinFields := [...][]ent.Field{
		auditlogMixin[0].Fields(),
	}
	auditlogFields := schema.AuditLog{}.Fields()
	_ = auditlogFields
	// auditlogDescCreatedAt is the schema descriptor for created_at field.
	auditlogDescCreatedAt := auditlogMixinFields[0][0].Descriptor()
	// auditlog.DefaultCreatedAt holds the default value on creation for the created_at field.
	auditlog.DefaultCreatedAt = auditlogDescCreatedAt.Default.(func() time.Time)
	// auditlogDescUpdatedAt is the schema descriptor for updated_at field.
	auditlogDescUpdatedAt := auditlogMixinFields[0][1].Descriptor()
	// auditlog.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	auditlog.DefaultUpdatedAt = auditlogDescUpdatedAt.Default.(func() time.Time)
	// auditlog.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	auditlog.UpdateDefaultUpdatedAt = auditlogDescUpdatedAt.UpdateDefault.(func() time.Time)
	tenantMixin := schema.Tenant{}.Mixin()
	tenantMixinFields := [...][]ent.Field{
		tenantMixin[0].Fields(),
	}
	tenantFields := schema.Tenant{}.Fields()
	_ = tenantFields
	// tenantDescCreatedAt is the schema descriptor for created_at field.
	tenantDescCreatedAt := tenantMixinFields[0][0].Descriptor()
	// tenant.DefaultCreatedAt holds the default value on creation for the created_at field.
	tenant.DefaultCreatedAt = tenantDescCreatedAt.Default.(func() time.Time)
	// tenantDescUpdatedAt is the schema descriptor for updated_at field.
	tenantDescUpdatedAt := tenantMixinFields[0][1].Descriptor()
	// tenant.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	tenant.DefaultUpdatedAt = tenantDescUpdatedAt.Default.(func() time.Time)
	// tenant.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	tenant.UpdateDefaultUpdatedAt = tenantDescUpdatedAt.UpdateDefault.(func() time.Time)
	// tenantDescName is the schema descriptor for name field.
	tenantDescName := tenantFields[0].Descriptor()
	// tenant.NameValidator is a validator for the "name" field. It is called by the builders before save.
	tenant.NameValidator = tenantDescName.Validators[0].(func(string) error)
	// tenantDescSSOCert is the schema descriptor for SSOCert field.
	tenantDescSSOCert := tenantFields[4].Descriptor()
	// tenant.DefaultSSOCert holds the default value on creation for the SSOCert field.
	tenant.DefaultSSOCert = tenantDescSSOCert.Default.(string)
	// tenantDescSSOEntryPoint is the schema descriptor for SSOEntryPoint field.
	tenantDescSSOEntryPoint := tenantFields[5].Descriptor()
	// tenant.DefaultSSOEntryPoint holds the default value on creation for the SSOEntryPoint field.
	tenant.DefaultSSOEntryPoint = tenantDescSSOEntryPoint.Default.(string)
	// tenantDescSSOIssuer is the schema descriptor for SSOIssuer field.
	tenantDescSSOIssuer := tenantFields[6].Descriptor()
	// tenant.DefaultSSOIssuer holds the default value on creation for the SSOIssuer field.
	tenant.DefaultSSOIssuer = tenantDescSSOIssuer.Default.(string)
	tokenMixin := schema.Token{}.Mixin()
	tokenMixinFields := [...][]ent.Field{
		tokenMixin[0].Fields(),
	}
	tokenFields := schema.Token{}.Fields()
	_ = tokenFields
	// tokenDescCreatedAt is the schema descriptor for created_at field.
	tokenDescCreatedAt := tokenMixinFields[0][0].Descriptor()
	// token.DefaultCreatedAt holds the default value on creation for the created_at field.
	token.DefaultCreatedAt = tokenDescCreatedAt.Default.(func() time.Time)
	// tokenDescUpdatedAt is the schema descriptor for updated_at field.
	tokenDescUpdatedAt := tokenMixinFields[0][1].Descriptor()
	// token.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	token.DefaultUpdatedAt = tokenDescUpdatedAt.Default.(func() time.Time)
	// token.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	token.UpdateDefaultUpdatedAt = tokenDescUpdatedAt.UpdateDefault.(func() time.Time)
	// tokenDescValue is the schema descriptor for value field.
	tokenDescValue := tokenFields[0].Descriptor()
	// token.ValueValidator is a validator for the "value" field. It is called by the builders before save.
	token.ValueValidator = tokenDescValue.Validators[0].(func(string) error)
	userMixin := schema.User{}.Mixin()
	userMixinFields := [...][]ent.Field{
		userMixin[0].Fields(),
	}
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescCreatedAt is the schema descriptor for created_at field.
	userDescCreatedAt := userMixinFields[0][0].Descriptor()
	// user.DefaultCreatedAt holds the default value on creation for the created_at field.
	user.DefaultCreatedAt = userDescCreatedAt.Default.(func() time.Time)
	// userDescUpdatedAt is the schema descriptor for updated_at field.
	userDescUpdatedAt := userMixinFields[0][1].Descriptor()
	// user.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	user.DefaultUpdatedAt = userDescUpdatedAt.Default.(func() time.Time)
	// user.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	user.UpdateDefaultUpdatedAt = userDescUpdatedAt.UpdateDefault.(func() time.Time)
	// userDescEmail is the schema descriptor for email field.
	userDescEmail := userFields[0].Descriptor()
	// user.EmailValidator is a validator for the "email" field. It is called by the builders before save.
	user.EmailValidator = userDescEmail.Validators[0].(func(string) error)
	// userDescPassword is the schema descriptor for password field.
	userDescPassword := userFields[1].Descriptor()
	// user.PasswordValidator is a validator for the "password" field. It is called by the builders before save.
	user.PasswordValidator = userDescPassword.Validators[0].(func(string) error)
	// userDescRole is the schema descriptor for role field.
	userDescRole := userFields[2].Descriptor()
	// user.DefaultRole holds the default value on creation for the role field.
	user.DefaultRole = userDescRole.Default.(int)
	// user.RoleValidator is a validator for the "role" field. It is called by the builders before save.
	user.RoleValidator = userDescRole.Validators[0].(func(int) error)
	// userDescTenant is the schema descriptor for tenant field.
	userDescTenant := userFields[3].Descriptor()
	// user.DefaultTenant holds the default value on creation for the tenant field.
	user.DefaultTenant = userDescTenant.Default.(string)
}
