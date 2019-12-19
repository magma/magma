// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package user

import (
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/frontier/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.User {
	return predicate.User(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldID), id))
		},
	)
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	},
	)
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	},
	)
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	},
	)
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	},
	)
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	},
	)
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	},
	)
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	},
	)
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	},
	)
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	},
	)
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	},
	)
}

// Email applies equality check predicate on the "email" field. It's identical to EmailEQ.
func Email(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEmail), v))
	},
	)
}

// Password applies equality check predicate on the "password" field. It's identical to PasswordEQ.
func Password(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPassword), v))
	},
	)
}

// Role applies equality check predicate on the "role" field. It's identical to RoleEQ.
func Role(v int) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRole), v))
	},
	)
}

// Tenant applies equality check predicate on the "tenant" field. It's identical to TenantEQ.
func Tenant(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTenant), v))
	},
	)
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	},
	)
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreatedAt), v))
	},
	)
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.User {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.User(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCreatedAt), v...))
	},
	)
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.User {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.User(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCreatedAt), v...))
	},
	)
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreatedAt), v))
	},
	)
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreatedAt), v))
	},
	)
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreatedAt), v))
	},
	)
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreatedAt), v))
	},
	)
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	},
	)
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdatedAt), v))
	},
	)
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.User {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.User(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldUpdatedAt), v...))
	},
	)
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.User {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.User(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldUpdatedAt), v...))
	},
	)
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdatedAt), v))
	},
	)
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdatedAt), v))
	},
	)
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdatedAt), v))
	},
	)
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdatedAt), v))
	},
	)
}

// EmailEQ applies the EQ predicate on the "email" field.
func EmailEQ(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEmail), v))
	},
	)
}

// EmailNEQ applies the NEQ predicate on the "email" field.
func EmailNEQ(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldEmail), v))
	},
	)
}

// EmailIn applies the In predicate on the "email" field.
func EmailIn(vs ...string) predicate.User {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.User(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldEmail), v...))
	},
	)
}

// EmailNotIn applies the NotIn predicate on the "email" field.
func EmailNotIn(vs ...string) predicate.User {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.User(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldEmail), v...))
	},
	)
}

// EmailGT applies the GT predicate on the "email" field.
func EmailGT(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldEmail), v))
	},
	)
}

// EmailGTE applies the GTE predicate on the "email" field.
func EmailGTE(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldEmail), v))
	},
	)
}

// EmailLT applies the LT predicate on the "email" field.
func EmailLT(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldEmail), v))
	},
	)
}

// EmailLTE applies the LTE predicate on the "email" field.
func EmailLTE(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldEmail), v))
	},
	)
}

// EmailContains applies the Contains predicate on the "email" field.
func EmailContains(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldEmail), v))
	},
	)
}

// EmailHasPrefix applies the HasPrefix predicate on the "email" field.
func EmailHasPrefix(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldEmail), v))
	},
	)
}

// EmailHasSuffix applies the HasSuffix predicate on the "email" field.
func EmailHasSuffix(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldEmail), v))
	},
	)
}

// EmailEqualFold applies the EqualFold predicate on the "email" field.
func EmailEqualFold(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldEmail), v))
	},
	)
}

// EmailContainsFold applies the ContainsFold predicate on the "email" field.
func EmailContainsFold(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldEmail), v))
	},
	)
}

// PasswordEQ applies the EQ predicate on the "password" field.
func PasswordEQ(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPassword), v))
	},
	)
}

// PasswordNEQ applies the NEQ predicate on the "password" field.
func PasswordNEQ(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldPassword), v))
	},
	)
}

// PasswordIn applies the In predicate on the "password" field.
func PasswordIn(vs ...string) predicate.User {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.User(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldPassword), v...))
	},
	)
}

// PasswordNotIn applies the NotIn predicate on the "password" field.
func PasswordNotIn(vs ...string) predicate.User {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.User(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldPassword), v...))
	},
	)
}

// PasswordGT applies the GT predicate on the "password" field.
func PasswordGT(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldPassword), v))
	},
	)
}

// PasswordGTE applies the GTE predicate on the "password" field.
func PasswordGTE(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldPassword), v))
	},
	)
}

// PasswordLT applies the LT predicate on the "password" field.
func PasswordLT(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldPassword), v))
	},
	)
}

// PasswordLTE applies the LTE predicate on the "password" field.
func PasswordLTE(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldPassword), v))
	},
	)
}

// PasswordContains applies the Contains predicate on the "password" field.
func PasswordContains(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldPassword), v))
	},
	)
}

// PasswordHasPrefix applies the HasPrefix predicate on the "password" field.
func PasswordHasPrefix(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldPassword), v))
	},
	)
}

// PasswordHasSuffix applies the HasSuffix predicate on the "password" field.
func PasswordHasSuffix(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldPassword), v))
	},
	)
}

// PasswordEqualFold applies the EqualFold predicate on the "password" field.
func PasswordEqualFold(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldPassword), v))
	},
	)
}

// PasswordContainsFold applies the ContainsFold predicate on the "password" field.
func PasswordContainsFold(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldPassword), v))
	},
	)
}

// RoleEQ applies the EQ predicate on the "role" field.
func RoleEQ(v int) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRole), v))
	},
	)
}

// RoleNEQ applies the NEQ predicate on the "role" field.
func RoleNEQ(v int) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldRole), v))
	},
	)
}

// RoleIn applies the In predicate on the "role" field.
func RoleIn(vs ...int) predicate.User {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.User(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldRole), v...))
	},
	)
}

// RoleNotIn applies the NotIn predicate on the "role" field.
func RoleNotIn(vs ...int) predicate.User {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.User(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldRole), v...))
	},
	)
}

// RoleGT applies the GT predicate on the "role" field.
func RoleGT(v int) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldRole), v))
	},
	)
}

// RoleGTE applies the GTE predicate on the "role" field.
func RoleGTE(v int) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldRole), v))
	},
	)
}

// RoleLT applies the LT predicate on the "role" field.
func RoleLT(v int) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldRole), v))
	},
	)
}

// RoleLTE applies the LTE predicate on the "role" field.
func RoleLTE(v int) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldRole), v))
	},
	)
}

// TenantEQ applies the EQ predicate on the "tenant" field.
func TenantEQ(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTenant), v))
	},
	)
}

// TenantNEQ applies the NEQ predicate on the "tenant" field.
func TenantNEQ(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTenant), v))
	},
	)
}

// TenantIn applies the In predicate on the "tenant" field.
func TenantIn(vs ...string) predicate.User {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.User(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldTenant), v...))
	},
	)
}

// TenantNotIn applies the NotIn predicate on the "tenant" field.
func TenantNotIn(vs ...string) predicate.User {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.User(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldTenant), v...))
	},
	)
}

// TenantGT applies the GT predicate on the "tenant" field.
func TenantGT(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldTenant), v))
	},
	)
}

// TenantGTE applies the GTE predicate on the "tenant" field.
func TenantGTE(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldTenant), v))
	},
	)
}

// TenantLT applies the LT predicate on the "tenant" field.
func TenantLT(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldTenant), v))
	},
	)
}

// TenantLTE applies the LTE predicate on the "tenant" field.
func TenantLTE(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldTenant), v))
	},
	)
}

// TenantContains applies the Contains predicate on the "tenant" field.
func TenantContains(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldTenant), v))
	},
	)
}

// TenantHasPrefix applies the HasPrefix predicate on the "tenant" field.
func TenantHasPrefix(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldTenant), v))
	},
	)
}

// TenantHasSuffix applies the HasSuffix predicate on the "tenant" field.
func TenantHasSuffix(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldTenant), v))
	},
	)
}

// TenantEqualFold applies the EqualFold predicate on the "tenant" field.
func TenantEqualFold(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldTenant), v))
	},
	)
}

// TenantContainsFold applies the ContainsFold predicate on the "tenant" field.
func TenantContainsFold(v string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldTenant), v))
	},
	)
}

// TabsIsNil applies the IsNil predicate on the "tabs" field.
func TabsIsNil() predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldTabs)))
	},
	)
}

// TabsNotNil applies the NotNil predicate on the "tabs" field.
func TabsNotNil() predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldTabs)))
	},
	)
}

// HasTokens applies the HasEdge predicate on the "tokens" edge.
func HasTokens() predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TokensTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, TokensTable, TokensColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	},
	)
}

// HasTokensWith applies the HasEdge predicate on the "tokens" edge with a given conditions (other predicates).
func HasTokensWith(preds ...predicate.Token) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TokensInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, TokensTable, TokensColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	},
	)
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.User) predicate.User {
	return predicate.User(
		func(s *sql.Selector) {
			s1 := s.Clone().SetP(nil)
			for _, p := range predicates {
				p(s1)
			}
			s.Where(s1.P())
		},
	)
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.User) predicate.User {
	return predicate.User(
		func(s *sql.Selector) {
			s1 := s.Clone().SetP(nil)
			for i, p := range predicates {
				if i > 0 {
					s1.Or()
				}
				p(s1)
			}
			s.Where(s1.P())
		},
	)
}

// Not applies the not operator on the given predicate.
func Not(p predicate.User) predicate.User {
	return predicate.User(
		func(s *sql.Selector) {
			p(s.Not())
		},
	)
}
