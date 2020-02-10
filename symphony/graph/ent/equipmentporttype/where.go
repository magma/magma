// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package equipmentporttype

import (
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i], _ = strconv.Atoi(ids[i])
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i], _ = strconv.Atoi(ids[i])
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.EquipmentPortType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortType(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.EquipmentPortType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortType(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.EquipmentPortType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortType(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.EquipmentPortType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortType(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.EquipmentPortType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortType(func(s *sql.Selector) {
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
func NameNotIn(vs ...string) predicate.EquipmentPortType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortType(func(s *sql.Selector) {
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
func NameGT(v string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// HasPropertyTypes applies the HasEdge predicate on the "property_types" edge.
func HasPropertyTypes() predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PropertyTypesTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PropertyTypesTable, PropertyTypesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPropertyTypesWith applies the HasEdge predicate on the "property_types" edge with a given conditions (other predicates).
func HasPropertyTypesWith(preds ...predicate.PropertyType) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
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

// HasLinkPropertyTypes applies the HasEdge predicate on the "link_property_types" edge.
func HasLinkPropertyTypes() predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LinkPropertyTypesTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, LinkPropertyTypesTable, LinkPropertyTypesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLinkPropertyTypesWith applies the HasEdge predicate on the "link_property_types" edge with a given conditions (other predicates).
func HasLinkPropertyTypesWith(preds ...predicate.PropertyType) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LinkPropertyTypesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, LinkPropertyTypesTable, LinkPropertyTypesColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPortDefinitions applies the HasEdge predicate on the "port_definitions" edge.
func HasPortDefinitions() predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PortDefinitionsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, PortDefinitionsTable, PortDefinitionsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPortDefinitionsWith applies the HasEdge predicate on the "port_definitions" edge with a given conditions (other predicates).
func HasPortDefinitionsWith(preds ...predicate.EquipmentPortDefinition) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PortDefinitionsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, PortDefinitionsTable, PortDefinitionsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.EquipmentPortType) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.EquipmentPortType) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
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
func Not(p predicate.EquipmentPortType) predicate.EquipmentPortType {
	return predicate.EquipmentPortType(func(s *sql.Selector) {
		p(s.Not())
	})
}
