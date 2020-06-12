// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package workorderdefinition

import (
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
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
func IDGT(id int) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// Index applies equality check predicate on the "index" field. It's identical to IndexEQ.
func Index(v int) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.WorkOrderDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.WorkOrderDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.WorkOrderDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.WorkOrderDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// IndexEQ applies the EQ predicate on the "index" field.
func IndexEQ(v int) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	})
}

// IndexNEQ applies the NEQ predicate on the "index" field.
func IndexNEQ(v int) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldIndex), v))
	})
}

// IndexIn applies the In predicate on the "index" field.
func IndexIn(vs ...int) predicate.WorkOrderDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
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
func IndexNotIn(vs ...int) predicate.WorkOrderDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
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
func IndexGT(v int) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldIndex), v))
	})
}

// IndexGTE applies the GTE predicate on the "index" field.
func IndexGTE(v int) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldIndex), v))
	})
}

// IndexLT applies the LT predicate on the "index" field.
func IndexLT(v int) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldIndex), v))
	})
}

// IndexLTE applies the LTE predicate on the "index" field.
func IndexLTE(v int) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldIndex), v))
	})
}

// IndexIsNil applies the IsNil predicate on the "index" field.
func IndexIsNil() predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldIndex)))
	})
}

// IndexNotNil applies the NotNil predicate on the "index" field.
func IndexNotNil() predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldIndex)))
	})
}

// HasType applies the HasEdge predicate on the "type" edge.
func HasType() predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, TypeTable, TypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTypeWith applies the HasEdge predicate on the "type" edge with a given conditions (other predicates).
func HasTypeWith(preds ...predicate.WorkOrderType) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TypeInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, TypeTable, TypeColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProjectType applies the HasEdge predicate on the "project_type" edge.
func HasProjectType() predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProjectTypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProjectTypeTable, ProjectTypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProjectTypeWith applies the HasEdge predicate on the "project_type" edge with a given conditions (other predicates).
func HasProjectTypeWith(preds ...predicate.ProjectType) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProjectTypeInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProjectTypeTable, ProjectTypeColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.WorkOrderDefinition) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.WorkOrderDefinition) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
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
func Not(p predicate.WorkOrderDefinition) predicate.WorkOrderDefinition {
	return predicate.WorkOrderDefinition(func(s *sql.Selector) {
		p(s.Not())
	})
}
