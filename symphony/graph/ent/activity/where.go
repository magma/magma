// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package activity

import (
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
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
func IDGT(id int) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// IsCreate applies equality check predicate on the "is_create" field. It's identical to IsCreateEQ.
func IsCreate(v bool) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIsCreate), v))
	})
}

// OldValue applies equality check predicate on the "old_value" field. It's identical to OldValueEQ.
func OldValue(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldOldValue), v))
	})
}

// NewValue applies equality check predicate on the "new_value" field. It's identical to NewValueEQ.
func NewValue(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldNewValue), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.Activity {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Activity(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.Activity {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Activity(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.Activity {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Activity(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.Activity {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Activity(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// ChangedFieldEQ applies the EQ predicate on the "changed_field" field.
func ChangedFieldEQ(v ChangedField) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldChangedField), v))
	})
}

// ChangedFieldNEQ applies the NEQ predicate on the "changed_field" field.
func ChangedFieldNEQ(v ChangedField) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldChangedField), v))
	})
}

// ChangedFieldIn applies the In predicate on the "changed_field" field.
func ChangedFieldIn(vs ...ChangedField) predicate.Activity {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Activity(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldChangedField), v...))
	})
}

// ChangedFieldNotIn applies the NotIn predicate on the "changed_field" field.
func ChangedFieldNotIn(vs ...ChangedField) predicate.Activity {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Activity(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldChangedField), v...))
	})
}

// IsCreateEQ applies the EQ predicate on the "is_create" field.
func IsCreateEQ(v bool) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIsCreate), v))
	})
}

// IsCreateNEQ applies the NEQ predicate on the "is_create" field.
func IsCreateNEQ(v bool) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldIsCreate), v))
	})
}

// OldValueEQ applies the EQ predicate on the "old_value" field.
func OldValueEQ(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldOldValue), v))
	})
}

// OldValueNEQ applies the NEQ predicate on the "old_value" field.
func OldValueNEQ(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldOldValue), v))
	})
}

// OldValueIn applies the In predicate on the "old_value" field.
func OldValueIn(vs ...string) predicate.Activity {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Activity(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldOldValue), v...))
	})
}

// OldValueNotIn applies the NotIn predicate on the "old_value" field.
func OldValueNotIn(vs ...string) predicate.Activity {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Activity(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldOldValue), v...))
	})
}

// OldValueGT applies the GT predicate on the "old_value" field.
func OldValueGT(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldOldValue), v))
	})
}

// OldValueGTE applies the GTE predicate on the "old_value" field.
func OldValueGTE(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldOldValue), v))
	})
}

// OldValueLT applies the LT predicate on the "old_value" field.
func OldValueLT(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldOldValue), v))
	})
}

// OldValueLTE applies the LTE predicate on the "old_value" field.
func OldValueLTE(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldOldValue), v))
	})
}

// OldValueContains applies the Contains predicate on the "old_value" field.
func OldValueContains(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldOldValue), v))
	})
}

// OldValueHasPrefix applies the HasPrefix predicate on the "old_value" field.
func OldValueHasPrefix(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldOldValue), v))
	})
}

// OldValueHasSuffix applies the HasSuffix predicate on the "old_value" field.
func OldValueHasSuffix(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldOldValue), v))
	})
}

// OldValueIsNil applies the IsNil predicate on the "old_value" field.
func OldValueIsNil() predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldOldValue)))
	})
}

// OldValueNotNil applies the NotNil predicate on the "old_value" field.
func OldValueNotNil() predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldOldValue)))
	})
}

// OldValueEqualFold applies the EqualFold predicate on the "old_value" field.
func OldValueEqualFold(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldOldValue), v))
	})
}

// OldValueContainsFold applies the ContainsFold predicate on the "old_value" field.
func OldValueContainsFold(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldOldValue), v))
	})
}

// NewValueEQ applies the EQ predicate on the "new_value" field.
func NewValueEQ(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldNewValue), v))
	})
}

// NewValueNEQ applies the NEQ predicate on the "new_value" field.
func NewValueNEQ(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldNewValue), v))
	})
}

// NewValueIn applies the In predicate on the "new_value" field.
func NewValueIn(vs ...string) predicate.Activity {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Activity(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldNewValue), v...))
	})
}

// NewValueNotIn applies the NotIn predicate on the "new_value" field.
func NewValueNotIn(vs ...string) predicate.Activity {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Activity(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldNewValue), v...))
	})
}

// NewValueGT applies the GT predicate on the "new_value" field.
func NewValueGT(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldNewValue), v))
	})
}

// NewValueGTE applies the GTE predicate on the "new_value" field.
func NewValueGTE(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldNewValue), v))
	})
}

// NewValueLT applies the LT predicate on the "new_value" field.
func NewValueLT(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldNewValue), v))
	})
}

// NewValueLTE applies the LTE predicate on the "new_value" field.
func NewValueLTE(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldNewValue), v))
	})
}

// NewValueContains applies the Contains predicate on the "new_value" field.
func NewValueContains(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldNewValue), v))
	})
}

// NewValueHasPrefix applies the HasPrefix predicate on the "new_value" field.
func NewValueHasPrefix(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldNewValue), v))
	})
}

// NewValueHasSuffix applies the HasSuffix predicate on the "new_value" field.
func NewValueHasSuffix(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldNewValue), v))
	})
}

// NewValueIsNil applies the IsNil predicate on the "new_value" field.
func NewValueIsNil() predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldNewValue)))
	})
}

// NewValueNotNil applies the NotNil predicate on the "new_value" field.
func NewValueNotNil() predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldNewValue)))
	})
}

// NewValueEqualFold applies the EqualFold predicate on the "new_value" field.
func NewValueEqualFold(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldNewValue), v))
	})
}

// NewValueContainsFold applies the ContainsFold predicate on the "new_value" field.
func NewValueContainsFold(v string) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldNewValue), v))
	})
}

// HasAuthor applies the HasEdge predicate on the "author" edge.
func HasAuthor() predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(AuthorTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, AuthorTable, AuthorColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasAuthorWith applies the HasEdge predicate on the "author" edge with a given conditions (other predicates).
func HasAuthorWith(preds ...predicate.User) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(AuthorInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, AuthorTable, AuthorColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasWorkOrder applies the HasEdge predicate on the "work_order" edge.
func HasWorkOrder() predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(WorkOrderTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, WorkOrderTable, WorkOrderColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasWorkOrderWith applies the HasEdge predicate on the "work_order" edge with a given conditions (other predicates).
func HasWorkOrderWith(preds ...predicate.WorkOrder) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(WorkOrderInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, WorkOrderTable, WorkOrderColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.Activity) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.Activity) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
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
func Not(p predicate.Activity) predicate.Activity {
	return predicate.Activity(func(s *sql.Selector) {
		p(s.Not())
	})
}
