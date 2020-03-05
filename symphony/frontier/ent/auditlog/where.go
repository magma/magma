// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package auditlog

import (
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/frontier/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
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
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
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
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// ActingUserID applies equality check predicate on the "acting_user_id" field. It's identical to ActingUserIDEQ.
func ActingUserID(v int) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldActingUserID), v))
	})
}

// Organization applies equality check predicate on the "organization" field. It's identical to OrganizationEQ.
func Organization(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldOrganization), v))
	})
}

// MutationType applies equality check predicate on the "mutation_type" field. It's identical to MutationTypeEQ.
func MutationType(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldMutationType), v))
	})
}

// ObjectID applies equality check predicate on the "object_id" field. It's identical to ObjectIDEQ.
func ObjectID(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldObjectID), v))
	})
}

// ObjectType applies equality check predicate on the "object_type" field. It's identical to ObjectTypeEQ.
func ObjectType(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldObjectType), v))
	})
}

// ObjectDisplayName applies equality check predicate on the "object_display_name" field. It's identical to ObjectDisplayNameEQ.
func ObjectDisplayName(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldObjectDisplayName), v))
	})
}

// URL applies equality check predicate on the "url" field. It's identical to URLEQ.
func URL(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldURL), v))
	})
}

// IPAddress applies equality check predicate on the "ip_address" field. It's identical to IPAddressEQ.
func IPAddress(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIPAddress), v))
	})
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStatus), v))
	})
}

// StatusCode applies equality check predicate on the "status_code" field. It's identical to StatusCodeEQ.
func StatusCode(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStatusCode), v))
	})
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldUpdatedAt), v...))
	})
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldUpdatedAt), v...))
	})
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdatedAt), v))
	})
}

// ActingUserIDEQ applies the EQ predicate on the "acting_user_id" field.
func ActingUserIDEQ(v int) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldActingUserID), v))
	})
}

// ActingUserIDNEQ applies the NEQ predicate on the "acting_user_id" field.
func ActingUserIDNEQ(v int) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldActingUserID), v))
	})
}

// ActingUserIDIn applies the In predicate on the "acting_user_id" field.
func ActingUserIDIn(vs ...int) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldActingUserID), v...))
	})
}

// ActingUserIDNotIn applies the NotIn predicate on the "acting_user_id" field.
func ActingUserIDNotIn(vs ...int) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldActingUserID), v...))
	})
}

// ActingUserIDGT applies the GT predicate on the "acting_user_id" field.
func ActingUserIDGT(v int) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldActingUserID), v))
	})
}

// ActingUserIDGTE applies the GTE predicate on the "acting_user_id" field.
func ActingUserIDGTE(v int) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldActingUserID), v))
	})
}

// ActingUserIDLT applies the LT predicate on the "acting_user_id" field.
func ActingUserIDLT(v int) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldActingUserID), v))
	})
}

// ActingUserIDLTE applies the LTE predicate on the "acting_user_id" field.
func ActingUserIDLTE(v int) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldActingUserID), v))
	})
}

// OrganizationEQ applies the EQ predicate on the "organization" field.
func OrganizationEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldOrganization), v))
	})
}

// OrganizationNEQ applies the NEQ predicate on the "organization" field.
func OrganizationNEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldOrganization), v))
	})
}

// OrganizationIn applies the In predicate on the "organization" field.
func OrganizationIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldOrganization), v...))
	})
}

// OrganizationNotIn applies the NotIn predicate on the "organization" field.
func OrganizationNotIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldOrganization), v...))
	})
}

// OrganizationGT applies the GT predicate on the "organization" field.
func OrganizationGT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldOrganization), v))
	})
}

// OrganizationGTE applies the GTE predicate on the "organization" field.
func OrganizationGTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldOrganization), v))
	})
}

// OrganizationLT applies the LT predicate on the "organization" field.
func OrganizationLT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldOrganization), v))
	})
}

// OrganizationLTE applies the LTE predicate on the "organization" field.
func OrganizationLTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldOrganization), v))
	})
}

// OrganizationContains applies the Contains predicate on the "organization" field.
func OrganizationContains(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldOrganization), v))
	})
}

// OrganizationHasPrefix applies the HasPrefix predicate on the "organization" field.
func OrganizationHasPrefix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldOrganization), v))
	})
}

// OrganizationHasSuffix applies the HasSuffix predicate on the "organization" field.
func OrganizationHasSuffix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldOrganization), v))
	})
}

// OrganizationEqualFold applies the EqualFold predicate on the "organization" field.
func OrganizationEqualFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldOrganization), v))
	})
}

// OrganizationContainsFold applies the ContainsFold predicate on the "organization" field.
func OrganizationContainsFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldOrganization), v))
	})
}

// MutationTypeEQ applies the EQ predicate on the "mutation_type" field.
func MutationTypeEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldMutationType), v))
	})
}

// MutationTypeNEQ applies the NEQ predicate on the "mutation_type" field.
func MutationTypeNEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldMutationType), v))
	})
}

// MutationTypeIn applies the In predicate on the "mutation_type" field.
func MutationTypeIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldMutationType), v...))
	})
}

// MutationTypeNotIn applies the NotIn predicate on the "mutation_type" field.
func MutationTypeNotIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldMutationType), v...))
	})
}

// MutationTypeGT applies the GT predicate on the "mutation_type" field.
func MutationTypeGT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldMutationType), v))
	})
}

// MutationTypeGTE applies the GTE predicate on the "mutation_type" field.
func MutationTypeGTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldMutationType), v))
	})
}

// MutationTypeLT applies the LT predicate on the "mutation_type" field.
func MutationTypeLT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldMutationType), v))
	})
}

// MutationTypeLTE applies the LTE predicate on the "mutation_type" field.
func MutationTypeLTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldMutationType), v))
	})
}

// MutationTypeContains applies the Contains predicate on the "mutation_type" field.
func MutationTypeContains(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldMutationType), v))
	})
}

// MutationTypeHasPrefix applies the HasPrefix predicate on the "mutation_type" field.
func MutationTypeHasPrefix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldMutationType), v))
	})
}

// MutationTypeHasSuffix applies the HasSuffix predicate on the "mutation_type" field.
func MutationTypeHasSuffix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldMutationType), v))
	})
}

// MutationTypeEqualFold applies the EqualFold predicate on the "mutation_type" field.
func MutationTypeEqualFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldMutationType), v))
	})
}

// MutationTypeContainsFold applies the ContainsFold predicate on the "mutation_type" field.
func MutationTypeContainsFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldMutationType), v))
	})
}

// ObjectIDEQ applies the EQ predicate on the "object_id" field.
func ObjectIDEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldObjectID), v))
	})
}

// ObjectIDNEQ applies the NEQ predicate on the "object_id" field.
func ObjectIDNEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldObjectID), v))
	})
}

// ObjectIDIn applies the In predicate on the "object_id" field.
func ObjectIDIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldObjectID), v...))
	})
}

// ObjectIDNotIn applies the NotIn predicate on the "object_id" field.
func ObjectIDNotIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldObjectID), v...))
	})
}

// ObjectIDGT applies the GT predicate on the "object_id" field.
func ObjectIDGT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldObjectID), v))
	})
}

// ObjectIDGTE applies the GTE predicate on the "object_id" field.
func ObjectIDGTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldObjectID), v))
	})
}

// ObjectIDLT applies the LT predicate on the "object_id" field.
func ObjectIDLT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldObjectID), v))
	})
}

// ObjectIDLTE applies the LTE predicate on the "object_id" field.
func ObjectIDLTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldObjectID), v))
	})
}

// ObjectIDContains applies the Contains predicate on the "object_id" field.
func ObjectIDContains(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldObjectID), v))
	})
}

// ObjectIDHasPrefix applies the HasPrefix predicate on the "object_id" field.
func ObjectIDHasPrefix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldObjectID), v))
	})
}

// ObjectIDHasSuffix applies the HasSuffix predicate on the "object_id" field.
func ObjectIDHasSuffix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldObjectID), v))
	})
}

// ObjectIDEqualFold applies the EqualFold predicate on the "object_id" field.
func ObjectIDEqualFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldObjectID), v))
	})
}

// ObjectIDContainsFold applies the ContainsFold predicate on the "object_id" field.
func ObjectIDContainsFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldObjectID), v))
	})
}

// ObjectTypeEQ applies the EQ predicate on the "object_type" field.
func ObjectTypeEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldObjectType), v))
	})
}

// ObjectTypeNEQ applies the NEQ predicate on the "object_type" field.
func ObjectTypeNEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldObjectType), v))
	})
}

// ObjectTypeIn applies the In predicate on the "object_type" field.
func ObjectTypeIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldObjectType), v...))
	})
}

// ObjectTypeNotIn applies the NotIn predicate on the "object_type" field.
func ObjectTypeNotIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldObjectType), v...))
	})
}

// ObjectTypeGT applies the GT predicate on the "object_type" field.
func ObjectTypeGT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldObjectType), v))
	})
}

// ObjectTypeGTE applies the GTE predicate on the "object_type" field.
func ObjectTypeGTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldObjectType), v))
	})
}

// ObjectTypeLT applies the LT predicate on the "object_type" field.
func ObjectTypeLT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldObjectType), v))
	})
}

// ObjectTypeLTE applies the LTE predicate on the "object_type" field.
func ObjectTypeLTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldObjectType), v))
	})
}

// ObjectTypeContains applies the Contains predicate on the "object_type" field.
func ObjectTypeContains(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldObjectType), v))
	})
}

// ObjectTypeHasPrefix applies the HasPrefix predicate on the "object_type" field.
func ObjectTypeHasPrefix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldObjectType), v))
	})
}

// ObjectTypeHasSuffix applies the HasSuffix predicate on the "object_type" field.
func ObjectTypeHasSuffix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldObjectType), v))
	})
}

// ObjectTypeEqualFold applies the EqualFold predicate on the "object_type" field.
func ObjectTypeEqualFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldObjectType), v))
	})
}

// ObjectTypeContainsFold applies the ContainsFold predicate on the "object_type" field.
func ObjectTypeContainsFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldObjectType), v))
	})
}

// ObjectDisplayNameEQ applies the EQ predicate on the "object_display_name" field.
func ObjectDisplayNameEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldObjectDisplayName), v))
	})
}

// ObjectDisplayNameNEQ applies the NEQ predicate on the "object_display_name" field.
func ObjectDisplayNameNEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldObjectDisplayName), v))
	})
}

// ObjectDisplayNameIn applies the In predicate on the "object_display_name" field.
func ObjectDisplayNameIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldObjectDisplayName), v...))
	})
}

// ObjectDisplayNameNotIn applies the NotIn predicate on the "object_display_name" field.
func ObjectDisplayNameNotIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldObjectDisplayName), v...))
	})
}

// ObjectDisplayNameGT applies the GT predicate on the "object_display_name" field.
func ObjectDisplayNameGT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldObjectDisplayName), v))
	})
}

// ObjectDisplayNameGTE applies the GTE predicate on the "object_display_name" field.
func ObjectDisplayNameGTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldObjectDisplayName), v))
	})
}

// ObjectDisplayNameLT applies the LT predicate on the "object_display_name" field.
func ObjectDisplayNameLT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldObjectDisplayName), v))
	})
}

// ObjectDisplayNameLTE applies the LTE predicate on the "object_display_name" field.
func ObjectDisplayNameLTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldObjectDisplayName), v))
	})
}

// ObjectDisplayNameContains applies the Contains predicate on the "object_display_name" field.
func ObjectDisplayNameContains(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldObjectDisplayName), v))
	})
}

// ObjectDisplayNameHasPrefix applies the HasPrefix predicate on the "object_display_name" field.
func ObjectDisplayNameHasPrefix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldObjectDisplayName), v))
	})
}

// ObjectDisplayNameHasSuffix applies the HasSuffix predicate on the "object_display_name" field.
func ObjectDisplayNameHasSuffix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldObjectDisplayName), v))
	})
}

// ObjectDisplayNameEqualFold applies the EqualFold predicate on the "object_display_name" field.
func ObjectDisplayNameEqualFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldObjectDisplayName), v))
	})
}

// ObjectDisplayNameContainsFold applies the ContainsFold predicate on the "object_display_name" field.
func ObjectDisplayNameContainsFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldObjectDisplayName), v))
	})
}

// URLEQ applies the EQ predicate on the "url" field.
func URLEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldURL), v))
	})
}

// URLNEQ applies the NEQ predicate on the "url" field.
func URLNEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldURL), v))
	})
}

// URLIn applies the In predicate on the "url" field.
func URLIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldURL), v...))
	})
}

// URLNotIn applies the NotIn predicate on the "url" field.
func URLNotIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldURL), v...))
	})
}

// URLGT applies the GT predicate on the "url" field.
func URLGT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldURL), v))
	})
}

// URLGTE applies the GTE predicate on the "url" field.
func URLGTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldURL), v))
	})
}

// URLLT applies the LT predicate on the "url" field.
func URLLT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldURL), v))
	})
}

// URLLTE applies the LTE predicate on the "url" field.
func URLLTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldURL), v))
	})
}

// URLContains applies the Contains predicate on the "url" field.
func URLContains(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldURL), v))
	})
}

// URLHasPrefix applies the HasPrefix predicate on the "url" field.
func URLHasPrefix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldURL), v))
	})
}

// URLHasSuffix applies the HasSuffix predicate on the "url" field.
func URLHasSuffix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldURL), v))
	})
}

// URLEqualFold applies the EqualFold predicate on the "url" field.
func URLEqualFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldURL), v))
	})
}

// URLContainsFold applies the ContainsFold predicate on the "url" field.
func URLContainsFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldURL), v))
	})
}

// IPAddressEQ applies the EQ predicate on the "ip_address" field.
func IPAddressEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIPAddress), v))
	})
}

// IPAddressNEQ applies the NEQ predicate on the "ip_address" field.
func IPAddressNEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldIPAddress), v))
	})
}

// IPAddressIn applies the In predicate on the "ip_address" field.
func IPAddressIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldIPAddress), v...))
	})
}

// IPAddressNotIn applies the NotIn predicate on the "ip_address" field.
func IPAddressNotIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldIPAddress), v...))
	})
}

// IPAddressGT applies the GT predicate on the "ip_address" field.
func IPAddressGT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldIPAddress), v))
	})
}

// IPAddressGTE applies the GTE predicate on the "ip_address" field.
func IPAddressGTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldIPAddress), v))
	})
}

// IPAddressLT applies the LT predicate on the "ip_address" field.
func IPAddressLT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldIPAddress), v))
	})
}

// IPAddressLTE applies the LTE predicate on the "ip_address" field.
func IPAddressLTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldIPAddress), v))
	})
}

// IPAddressContains applies the Contains predicate on the "ip_address" field.
func IPAddressContains(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldIPAddress), v))
	})
}

// IPAddressHasPrefix applies the HasPrefix predicate on the "ip_address" field.
func IPAddressHasPrefix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldIPAddress), v))
	})
}

// IPAddressHasSuffix applies the HasSuffix predicate on the "ip_address" field.
func IPAddressHasSuffix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldIPAddress), v))
	})
}

// IPAddressEqualFold applies the EqualFold predicate on the "ip_address" field.
func IPAddressEqualFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldIPAddress), v))
	})
}

// IPAddressContainsFold applies the ContainsFold predicate on the "ip_address" field.
func IPAddressContainsFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldIPAddress), v))
	})
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStatus), v))
	})
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldStatus), v))
	})
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldStatus), v...))
	})
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldStatus), v...))
	})
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldStatus), v))
	})
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldStatus), v))
	})
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldStatus), v))
	})
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldStatus), v))
	})
}

// StatusContains applies the Contains predicate on the "status" field.
func StatusContains(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldStatus), v))
	})
}

// StatusHasPrefix applies the HasPrefix predicate on the "status" field.
func StatusHasPrefix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldStatus), v))
	})
}

// StatusHasSuffix applies the HasSuffix predicate on the "status" field.
func StatusHasSuffix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldStatus), v))
	})
}

// StatusEqualFold applies the EqualFold predicate on the "status" field.
func StatusEqualFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldStatus), v))
	})
}

// StatusContainsFold applies the ContainsFold predicate on the "status" field.
func StatusContainsFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldStatus), v))
	})
}

// StatusCodeEQ applies the EQ predicate on the "status_code" field.
func StatusCodeEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStatusCode), v))
	})
}

// StatusCodeNEQ applies the NEQ predicate on the "status_code" field.
func StatusCodeNEQ(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldStatusCode), v))
	})
}

// StatusCodeIn applies the In predicate on the "status_code" field.
func StatusCodeIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldStatusCode), v...))
	})
}

// StatusCodeNotIn applies the NotIn predicate on the "status_code" field.
func StatusCodeNotIn(vs ...string) predicate.AuditLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AuditLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldStatusCode), v...))
	})
}

// StatusCodeGT applies the GT predicate on the "status_code" field.
func StatusCodeGT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldStatusCode), v))
	})
}

// StatusCodeGTE applies the GTE predicate on the "status_code" field.
func StatusCodeGTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldStatusCode), v))
	})
}

// StatusCodeLT applies the LT predicate on the "status_code" field.
func StatusCodeLT(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldStatusCode), v))
	})
}

// StatusCodeLTE applies the LTE predicate on the "status_code" field.
func StatusCodeLTE(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldStatusCode), v))
	})
}

// StatusCodeContains applies the Contains predicate on the "status_code" field.
func StatusCodeContains(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldStatusCode), v))
	})
}

// StatusCodeHasPrefix applies the HasPrefix predicate on the "status_code" field.
func StatusCodeHasPrefix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldStatusCode), v))
	})
}

// StatusCodeHasSuffix applies the HasSuffix predicate on the "status_code" field.
func StatusCodeHasSuffix(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldStatusCode), v))
	})
}

// StatusCodeEqualFold applies the EqualFold predicate on the "status_code" field.
func StatusCodeEqualFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldStatusCode), v))
	})
}

// StatusCodeContainsFold applies the ContainsFold predicate on the "status_code" field.
func StatusCodeContainsFold(v string) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldStatusCode), v))
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.AuditLog) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.AuditLog) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.AuditLog) predicate.AuditLog {
	return predicate.AuditLog(func(s *sql.Selector) {
		p(s.Not())
	})
}
