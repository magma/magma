// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package servicetype

import (
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
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
func IDGT(id int) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// HasCustomer applies equality check predicate on the "has_customer" field. It's identical to HasCustomerEQ.
func HasCustomer(v bool) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldHasCustomer), v))
	})
}

// IsDeleted applies equality check predicate on the "is_deleted" field. It's identical to IsDeletedEQ.
func IsDeleted(v bool) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIsDeleted), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.ServiceType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCreateTime), v...))
	})
}

// CreateTimeNotIn applies the NotIn predicate on the "create_time" field.
func CreateTimeNotIn(vs ...time.Time) predicate.ServiceType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCreateTime), v...))
	})
}

// CreateTimeGT applies the GT predicate on the "create_time" field.
func CreateTimeGT(v time.Time) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.ServiceType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldUpdateTime), v...))
	})
}

// UpdateTimeNotIn applies the NotIn predicate on the "update_time" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.ServiceType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldUpdateTime), v...))
	})
}

// UpdateTimeGT applies the GT predicate on the "update_time" field.
func UpdateTimeGT(v time.Time) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.ServiceType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldName), v...))
	})
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.ServiceType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldName), v...))
	})
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// HasCustomerEQ applies the EQ predicate on the "has_customer" field.
func HasCustomerEQ(v bool) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldHasCustomer), v))
	})
}

// HasCustomerNEQ applies the NEQ predicate on the "has_customer" field.
func HasCustomerNEQ(v bool) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldHasCustomer), v))
	})
}

// IsDeletedEQ applies the EQ predicate on the "is_deleted" field.
func IsDeletedEQ(v bool) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIsDeleted), v))
	})
}

// IsDeletedNEQ applies the NEQ predicate on the "is_deleted" field.
func IsDeletedNEQ(v bool) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldIsDeleted), v))
	})
}

// DiscoveryMethodEQ applies the EQ predicate on the "discovery_method" field.
func DiscoveryMethodEQ(v DiscoveryMethod) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDiscoveryMethod), v))
	})
}

// DiscoveryMethodNEQ applies the NEQ predicate on the "discovery_method" field.
func DiscoveryMethodNEQ(v DiscoveryMethod) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldDiscoveryMethod), v))
	})
}

// DiscoveryMethodIn applies the In predicate on the "discovery_method" field.
func DiscoveryMethodIn(vs ...DiscoveryMethod) predicate.ServiceType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldDiscoveryMethod), v...))
	})
}

// DiscoveryMethodNotIn applies the NotIn predicate on the "discovery_method" field.
func DiscoveryMethodNotIn(vs ...DiscoveryMethod) predicate.ServiceType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldDiscoveryMethod), v...))
	})
}

// DiscoveryMethodIsNil applies the IsNil predicate on the "discovery_method" field.
func DiscoveryMethodIsNil() predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldDiscoveryMethod)))
	})
}

// DiscoveryMethodNotNil applies the NotNil predicate on the "discovery_method" field.
func DiscoveryMethodNotNil() predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldDiscoveryMethod)))
	})
}

// HasServices applies the HasEdge predicate on the "services" edge.
func HasServices() predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ServicesTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, ServicesTable, ServicesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasServicesWith applies the HasEdge predicate on the "services" edge with a given conditions (other predicates).
func HasServicesWith(preds ...predicate.Service) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ServicesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, ServicesTable, ServicesColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPropertyTypes applies the HasEdge predicate on the "property_types" edge.
func HasPropertyTypes() predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PropertyTypesTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PropertyTypesTable, PropertyTypesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPropertyTypesWith applies the HasEdge predicate on the "property_types" edge with a given conditions (other predicates).
func HasPropertyTypesWith(preds ...predicate.PropertyType) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PropertyTypesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PropertyTypesTable, PropertyTypesColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEndpointDefinitions applies the HasEdge predicate on the "endpoint_definitions" edge.
func HasEndpointDefinitions() predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EndpointDefinitionsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, EndpointDefinitionsTable, EndpointDefinitionsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEndpointDefinitionsWith applies the HasEdge predicate on the "endpoint_definitions" edge with a given conditions (other predicates).
func HasEndpointDefinitionsWith(preds ...predicate.ServiceEndpointDefinition) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EndpointDefinitionsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, EndpointDefinitionsTable, EndpointDefinitionsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.ServiceType) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.ServiceType) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
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
func Not(p predicate.ServiceType) predicate.ServiceType {
	return predicate.ServiceType(func(s *sql.Selector) {
		p(s.Not())
	})
}
