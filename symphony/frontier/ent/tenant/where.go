// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package tenant

import (
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/frontier/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.Tenant {
	return predicate.Tenant(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldID), id))
		},
	)
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	},
	)
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	},
	)
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
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
func IDGT(id int) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	},
	)
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	},
	)
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	},
	)
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	},
	)
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	},
	)
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	},
	)
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	},
	)
}

// SSOCert applies equality check predicate on the "SSOCert" field. It's identical to SSOCertEQ.
func SSOCert(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSSOCert), v))
	},
	)
}

// SSOEntryPoint applies equality check predicate on the "SSOEntryPoint" field. It's identical to SSOEntryPointEQ.
func SSOEntryPoint(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSSOEntryPoint), v))
	},
	)
}

// SSOIssuer applies equality check predicate on the "SSOIssuer" field. It's identical to SSOIssuerEQ.
func SSOIssuer(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSSOIssuer), v))
	},
	)
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	},
	)
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreatedAt), v))
	},
	)
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Tenant {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Tenant(func(s *sql.Selector) {
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
func CreatedAtNotIn(vs ...time.Time) predicate.Tenant {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Tenant(func(s *sql.Selector) {
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
func CreatedAtGT(v time.Time) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreatedAt), v))
	},
	)
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreatedAt), v))
	},
	)
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreatedAt), v))
	},
	)
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreatedAt), v))
	},
	)
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	},
	)
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdatedAt), v))
	},
	)
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Tenant {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Tenant(func(s *sql.Selector) {
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
func UpdatedAtNotIn(vs ...time.Time) predicate.Tenant {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Tenant(func(s *sql.Selector) {
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
func UpdatedAtGT(v time.Time) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdatedAt), v))
	},
	)
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdatedAt), v))
	},
	)
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdatedAt), v))
	},
	)
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdatedAt), v))
	},
	)
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	},
	)
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	},
	)
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Tenant {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Tenant(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldName), v...))
	},
	)
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Tenant {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Tenant(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldName), v...))
	},
	)
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	},
	)
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	},
	)
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	},
	)
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	},
	)
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	},
	)
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	},
	)
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	},
	)
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	},
	)
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	},
	)
}

// TabsIsNil applies the IsNil predicate on the "tabs" field.
func TabsIsNil() predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldTabs)))
	},
	)
}

// TabsNotNil applies the NotNil predicate on the "tabs" field.
func TabsNotNil() predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldTabs)))
	},
	)
}

// SSOCertEQ applies the EQ predicate on the "SSOCert" field.
func SSOCertEQ(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSSOCert), v))
	},
	)
}

// SSOCertNEQ applies the NEQ predicate on the "SSOCert" field.
func SSOCertNEQ(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldSSOCert), v))
	},
	)
}

// SSOCertIn applies the In predicate on the "SSOCert" field.
func SSOCertIn(vs ...string) predicate.Tenant {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Tenant(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldSSOCert), v...))
	},
	)
}

// SSOCertNotIn applies the NotIn predicate on the "SSOCert" field.
func SSOCertNotIn(vs ...string) predicate.Tenant {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Tenant(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldSSOCert), v...))
	},
	)
}

// SSOCertGT applies the GT predicate on the "SSOCert" field.
func SSOCertGT(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldSSOCert), v))
	},
	)
}

// SSOCertGTE applies the GTE predicate on the "SSOCert" field.
func SSOCertGTE(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldSSOCert), v))
	},
	)
}

// SSOCertLT applies the LT predicate on the "SSOCert" field.
func SSOCertLT(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldSSOCert), v))
	},
	)
}

// SSOCertLTE applies the LTE predicate on the "SSOCert" field.
func SSOCertLTE(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldSSOCert), v))
	},
	)
}

// SSOCertContains applies the Contains predicate on the "SSOCert" field.
func SSOCertContains(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldSSOCert), v))
	},
	)
}

// SSOCertHasPrefix applies the HasPrefix predicate on the "SSOCert" field.
func SSOCertHasPrefix(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldSSOCert), v))
	},
	)
}

// SSOCertHasSuffix applies the HasSuffix predicate on the "SSOCert" field.
func SSOCertHasSuffix(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldSSOCert), v))
	},
	)
}

// SSOCertEqualFold applies the EqualFold predicate on the "SSOCert" field.
func SSOCertEqualFold(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldSSOCert), v))
	},
	)
}

// SSOCertContainsFold applies the ContainsFold predicate on the "SSOCert" field.
func SSOCertContainsFold(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldSSOCert), v))
	},
	)
}

// SSOEntryPointEQ applies the EQ predicate on the "SSOEntryPoint" field.
func SSOEntryPointEQ(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSSOEntryPoint), v))
	},
	)
}

// SSOEntryPointNEQ applies the NEQ predicate on the "SSOEntryPoint" field.
func SSOEntryPointNEQ(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldSSOEntryPoint), v))
	},
	)
}

// SSOEntryPointIn applies the In predicate on the "SSOEntryPoint" field.
func SSOEntryPointIn(vs ...string) predicate.Tenant {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Tenant(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldSSOEntryPoint), v...))
	},
	)
}

// SSOEntryPointNotIn applies the NotIn predicate on the "SSOEntryPoint" field.
func SSOEntryPointNotIn(vs ...string) predicate.Tenant {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Tenant(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldSSOEntryPoint), v...))
	},
	)
}

// SSOEntryPointGT applies the GT predicate on the "SSOEntryPoint" field.
func SSOEntryPointGT(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldSSOEntryPoint), v))
	},
	)
}

// SSOEntryPointGTE applies the GTE predicate on the "SSOEntryPoint" field.
func SSOEntryPointGTE(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldSSOEntryPoint), v))
	},
	)
}

// SSOEntryPointLT applies the LT predicate on the "SSOEntryPoint" field.
func SSOEntryPointLT(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldSSOEntryPoint), v))
	},
	)
}

// SSOEntryPointLTE applies the LTE predicate on the "SSOEntryPoint" field.
func SSOEntryPointLTE(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldSSOEntryPoint), v))
	},
	)
}

// SSOEntryPointContains applies the Contains predicate on the "SSOEntryPoint" field.
func SSOEntryPointContains(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldSSOEntryPoint), v))
	},
	)
}

// SSOEntryPointHasPrefix applies the HasPrefix predicate on the "SSOEntryPoint" field.
func SSOEntryPointHasPrefix(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldSSOEntryPoint), v))
	},
	)
}

// SSOEntryPointHasSuffix applies the HasSuffix predicate on the "SSOEntryPoint" field.
func SSOEntryPointHasSuffix(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldSSOEntryPoint), v))
	},
	)
}

// SSOEntryPointEqualFold applies the EqualFold predicate on the "SSOEntryPoint" field.
func SSOEntryPointEqualFold(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldSSOEntryPoint), v))
	},
	)
}

// SSOEntryPointContainsFold applies the ContainsFold predicate on the "SSOEntryPoint" field.
func SSOEntryPointContainsFold(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldSSOEntryPoint), v))
	},
	)
}

// SSOIssuerEQ applies the EQ predicate on the "SSOIssuer" field.
func SSOIssuerEQ(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSSOIssuer), v))
	},
	)
}

// SSOIssuerNEQ applies the NEQ predicate on the "SSOIssuer" field.
func SSOIssuerNEQ(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldSSOIssuer), v))
	},
	)
}

// SSOIssuerIn applies the In predicate on the "SSOIssuer" field.
func SSOIssuerIn(vs ...string) predicate.Tenant {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Tenant(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldSSOIssuer), v...))
	},
	)
}

// SSOIssuerNotIn applies the NotIn predicate on the "SSOIssuer" field.
func SSOIssuerNotIn(vs ...string) predicate.Tenant {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Tenant(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldSSOIssuer), v...))
	},
	)
}

// SSOIssuerGT applies the GT predicate on the "SSOIssuer" field.
func SSOIssuerGT(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldSSOIssuer), v))
	},
	)
}

// SSOIssuerGTE applies the GTE predicate on the "SSOIssuer" field.
func SSOIssuerGTE(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldSSOIssuer), v))
	},
	)
}

// SSOIssuerLT applies the LT predicate on the "SSOIssuer" field.
func SSOIssuerLT(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldSSOIssuer), v))
	},
	)
}

// SSOIssuerLTE applies the LTE predicate on the "SSOIssuer" field.
func SSOIssuerLTE(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldSSOIssuer), v))
	},
	)
}

// SSOIssuerContains applies the Contains predicate on the "SSOIssuer" field.
func SSOIssuerContains(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldSSOIssuer), v))
	},
	)
}

// SSOIssuerHasPrefix applies the HasPrefix predicate on the "SSOIssuer" field.
func SSOIssuerHasPrefix(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldSSOIssuer), v))
	},
	)
}

// SSOIssuerHasSuffix applies the HasSuffix predicate on the "SSOIssuer" field.
func SSOIssuerHasSuffix(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldSSOIssuer), v))
	},
	)
}

// SSOIssuerEqualFold applies the EqualFold predicate on the "SSOIssuer" field.
func SSOIssuerEqualFold(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldSSOIssuer), v))
	},
	)
}

// SSOIssuerContainsFold applies the ContainsFold predicate on the "SSOIssuer" field.
func SSOIssuerContainsFold(v string) predicate.Tenant {
	return predicate.Tenant(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldSSOIssuer), v))
	},
	)
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.Tenant) predicate.Tenant {
	return predicate.Tenant(
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
func Or(predicates ...predicate.Tenant) predicate.Tenant {
	return predicate.Tenant(
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
func Not(p predicate.Tenant) predicate.Tenant {
	return predicate.Tenant(
		func(s *sql.Selector) {
			p(s.Not())
		},
	)
}
