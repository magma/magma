// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package actionsrule

import (
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
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
func IDNotIn(ids ...string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
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
func IDGT(id string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// TriggerID applies equality check predicate on the "triggerID" field. It's identical to TriggerIDEQ.
func TriggerID(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTriggerID), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.ActionsRule {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ActionsRule(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.ActionsRule {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ActionsRule(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.ActionsRule {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ActionsRule(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.ActionsRule {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ActionsRule(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.ActionsRule {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ActionsRule(func(s *sql.Selector) {
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
func NameNotIn(vs ...string) predicate.ActionsRule {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ActionsRule(func(s *sql.Selector) {
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
func NameGT(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// TriggerIDEQ applies the EQ predicate on the "triggerID" field.
func TriggerIDEQ(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTriggerID), v))
	})
}

// TriggerIDNEQ applies the NEQ predicate on the "triggerID" field.
func TriggerIDNEQ(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTriggerID), v))
	})
}

// TriggerIDIn applies the In predicate on the "triggerID" field.
func TriggerIDIn(vs ...string) predicate.ActionsRule {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ActionsRule(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldTriggerID), v...))
	})
}

// TriggerIDNotIn applies the NotIn predicate on the "triggerID" field.
func TriggerIDNotIn(vs ...string) predicate.ActionsRule {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ActionsRule(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldTriggerID), v...))
	})
}

// TriggerIDGT applies the GT predicate on the "triggerID" field.
func TriggerIDGT(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldTriggerID), v))
	})
}

// TriggerIDGTE applies the GTE predicate on the "triggerID" field.
func TriggerIDGTE(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldTriggerID), v))
	})
}

// TriggerIDLT applies the LT predicate on the "triggerID" field.
func TriggerIDLT(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldTriggerID), v))
	})
}

// TriggerIDLTE applies the LTE predicate on the "triggerID" field.
func TriggerIDLTE(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldTriggerID), v))
	})
}

// TriggerIDContains applies the Contains predicate on the "triggerID" field.
func TriggerIDContains(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldTriggerID), v))
	})
}

// TriggerIDHasPrefix applies the HasPrefix predicate on the "triggerID" field.
func TriggerIDHasPrefix(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldTriggerID), v))
	})
}

// TriggerIDHasSuffix applies the HasSuffix predicate on the "triggerID" field.
func TriggerIDHasSuffix(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldTriggerID), v))
	})
}

// TriggerIDEqualFold applies the EqualFold predicate on the "triggerID" field.
func TriggerIDEqualFold(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldTriggerID), v))
	})
}

// TriggerIDContainsFold applies the ContainsFold predicate on the "triggerID" field.
func TriggerIDContainsFold(v string) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldTriggerID), v))
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.ActionsRule) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.ActionsRule) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
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
func Not(p predicate.ActionsRule) predicate.ActionsRule {
	return predicate.ActionsRule(func(s *sql.Selector) {
		p(s.Not())
	})
}
