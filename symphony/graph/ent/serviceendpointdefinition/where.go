// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package serviceendpointdefinition

import (
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
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
func IDGT(id int) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// Role applies equality check predicate on the "role" field. It's identical to RoleEQ.
func Role(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRole), v))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// Index applies equality check predicate on the "index" field. It's identical to IndexEQ.
func Index(v int) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.ServiceEndpointDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.ServiceEndpointDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.ServiceEndpointDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.ServiceEndpointDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// RoleEQ applies the EQ predicate on the "role" field.
func RoleEQ(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRole), v))
	})
}

// RoleNEQ applies the NEQ predicate on the "role" field.
func RoleNEQ(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldRole), v))
	})
}

// RoleIn applies the In predicate on the "role" field.
func RoleIn(vs ...string) predicate.ServiceEndpointDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldRole), v...))
	})
}

// RoleNotIn applies the NotIn predicate on the "role" field.
func RoleNotIn(vs ...string) predicate.ServiceEndpointDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldRole), v...))
	})
}

// RoleGT applies the GT predicate on the "role" field.
func RoleGT(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldRole), v))
	})
}

// RoleGTE applies the GTE predicate on the "role" field.
func RoleGTE(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldRole), v))
	})
}

// RoleLT applies the LT predicate on the "role" field.
func RoleLT(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldRole), v))
	})
}

// RoleLTE applies the LTE predicate on the "role" field.
func RoleLTE(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldRole), v))
	})
}

// RoleContains applies the Contains predicate on the "role" field.
func RoleContains(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldRole), v))
	})
}

// RoleHasPrefix applies the HasPrefix predicate on the "role" field.
func RoleHasPrefix(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldRole), v))
	})
}

// RoleHasSuffix applies the HasSuffix predicate on the "role" field.
func RoleHasSuffix(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldRole), v))
	})
}

// RoleIsNil applies the IsNil predicate on the "role" field.
func RoleIsNil() predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldRole)))
	})
}

// RoleNotNil applies the NotNil predicate on the "role" field.
func RoleNotNil() predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldRole)))
	})
}

// RoleEqualFold applies the EqualFold predicate on the "role" field.
func RoleEqualFold(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldRole), v))
	})
}

// RoleContainsFold applies the ContainsFold predicate on the "role" field.
func RoleContainsFold(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldRole), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.ServiceEndpointDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
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
func NameNotIn(vs ...string) predicate.ServiceEndpointDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
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
func NameGT(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// IndexEQ applies the EQ predicate on the "index" field.
func IndexEQ(v int) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	})
}

// IndexNEQ applies the NEQ predicate on the "index" field.
func IndexNEQ(v int) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldIndex), v))
	})
}

// IndexIn applies the In predicate on the "index" field.
func IndexIn(vs ...int) predicate.ServiceEndpointDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldIndex), v...))
	})
}

// IndexNotIn applies the NotIn predicate on the "index" field.
func IndexNotIn(vs ...int) predicate.ServiceEndpointDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldIndex), v...))
	})
}

// IndexGT applies the GT predicate on the "index" field.
func IndexGT(v int) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldIndex), v))
	})
}

// IndexGTE applies the GTE predicate on the "index" field.
func IndexGTE(v int) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldIndex), v))
	})
}

// IndexLT applies the LT predicate on the "index" field.
func IndexLT(v int) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldIndex), v))
	})
}

// IndexLTE applies the LTE predicate on the "index" field.
func IndexLTE(v int) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldIndex), v))
	})
}

// HasEndpoints applies the HasEdge predicate on the "endpoints" edge.
func HasEndpoints() predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EndpointsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, EndpointsTable, EndpointsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEndpointsWith applies the HasEdge predicate on the "endpoints" edge with a given conditions (other predicates).
func HasEndpointsWith(preds ...predicate.ServiceEndpoint) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EndpointsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, EndpointsTable, EndpointsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasServiceType applies the HasEdge predicate on the "service_type" edge.
func HasServiceType() predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ServiceTypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ServiceTypeTable, ServiceTypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasServiceTypeWith applies the HasEdge predicate on the "service_type" edge with a given conditions (other predicates).
func HasServiceTypeWith(preds ...predicate.ServiceType) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ServiceTypeInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ServiceTypeTable, ServiceTypeColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEquipmentType applies the HasEdge predicate on the "equipment_type" edge.
func HasEquipmentType() predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentTypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EquipmentTypeTable, EquipmentTypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEquipmentTypeWith applies the HasEdge predicate on the "equipment_type" edge with a given conditions (other predicates).
func HasEquipmentTypeWith(preds ...predicate.EquipmentType) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentTypeInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EquipmentTypeTable, EquipmentTypeColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.ServiceEndpointDefinition) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.ServiceEndpointDefinition) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
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
func Not(p predicate.ServiceEndpointDefinition) predicate.ServiceEndpointDefinition {
	return predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
		p(s.Not())
	})
}
