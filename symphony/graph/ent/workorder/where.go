// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package workorder

import (
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
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
func IDGT(id int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStatus), v))
	})
}

// Priority applies equality check predicate on the "priority" field. It's identical to PriorityEQ.
func Priority(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPriority), v))
	})
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDescription), v))
	})
}

// InstallDate applies equality check predicate on the "install_date" field. It's identical to InstallDateEQ.
func InstallDate(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldInstallDate), v))
	})
}

// CreationDate applies equality check predicate on the "creation_date" field. It's identical to CreationDateEQ.
func CreationDate(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreationDate), v))
	})
}

// Index applies equality check predicate on the "index" field. It's identical to IndexEQ.
func Index(v int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	})
}

// CloseDate applies equality check predicate on the "close_date" field. It's identical to CloseDateEQ.
func CloseDate(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCloseDate), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
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
func NameNotIn(vs ...string) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
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
func NameGT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStatus), v))
	})
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldStatus), v))
	})
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...string) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
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
func StatusNotIn(vs ...string) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
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
func StatusGT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldStatus), v))
	})
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldStatus), v))
	})
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldStatus), v))
	})
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldStatus), v))
	})
}

// StatusContains applies the Contains predicate on the "status" field.
func StatusContains(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldStatus), v))
	})
}

// StatusHasPrefix applies the HasPrefix predicate on the "status" field.
func StatusHasPrefix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldStatus), v))
	})
}

// StatusHasSuffix applies the HasSuffix predicate on the "status" field.
func StatusHasSuffix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldStatus), v))
	})
}

// StatusEqualFold applies the EqualFold predicate on the "status" field.
func StatusEqualFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldStatus), v))
	})
}

// StatusContainsFold applies the ContainsFold predicate on the "status" field.
func StatusContainsFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldStatus), v))
	})
}

// PriorityEQ applies the EQ predicate on the "priority" field.
func PriorityEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPriority), v))
	})
}

// PriorityNEQ applies the NEQ predicate on the "priority" field.
func PriorityNEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldPriority), v))
	})
}

// PriorityIn applies the In predicate on the "priority" field.
func PriorityIn(vs ...string) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldPriority), v...))
	})
}

// PriorityNotIn applies the NotIn predicate on the "priority" field.
func PriorityNotIn(vs ...string) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldPriority), v...))
	})
}

// PriorityGT applies the GT predicate on the "priority" field.
func PriorityGT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldPriority), v))
	})
}

// PriorityGTE applies the GTE predicate on the "priority" field.
func PriorityGTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldPriority), v))
	})
}

// PriorityLT applies the LT predicate on the "priority" field.
func PriorityLT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldPriority), v))
	})
}

// PriorityLTE applies the LTE predicate on the "priority" field.
func PriorityLTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldPriority), v))
	})
}

// PriorityContains applies the Contains predicate on the "priority" field.
func PriorityContains(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldPriority), v))
	})
}

// PriorityHasPrefix applies the HasPrefix predicate on the "priority" field.
func PriorityHasPrefix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldPriority), v))
	})
}

// PriorityHasSuffix applies the HasSuffix predicate on the "priority" field.
func PriorityHasSuffix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldPriority), v))
	})
}

// PriorityEqualFold applies the EqualFold predicate on the "priority" field.
func PriorityEqualFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldPriority), v))
	})
}

// PriorityContainsFold applies the ContainsFold predicate on the "priority" field.
func PriorityContainsFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldPriority), v))
	})
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDescription), v))
	})
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldDescription), v))
	})
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldDescription), v...))
	})
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldDescription), v...))
	})
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldDescription), v))
	})
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldDescription), v))
	})
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldDescription), v))
	})
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldDescription), v))
	})
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldDescription), v))
	})
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldDescription), v))
	})
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldDescription), v))
	})
}

// DescriptionIsNil applies the IsNil predicate on the "description" field.
func DescriptionIsNil() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldDescription)))
	})
}

// DescriptionNotNil applies the NotNil predicate on the "description" field.
func DescriptionNotNil() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldDescription)))
	})
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldDescription), v))
	})
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldDescription), v))
	})
}

// InstallDateEQ applies the EQ predicate on the "install_date" field.
func InstallDateEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldInstallDate), v))
	})
}

// InstallDateNEQ applies the NEQ predicate on the "install_date" field.
func InstallDateNEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldInstallDate), v))
	})
}

// InstallDateIn applies the In predicate on the "install_date" field.
func InstallDateIn(vs ...time.Time) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldInstallDate), v...))
	})
}

// InstallDateNotIn applies the NotIn predicate on the "install_date" field.
func InstallDateNotIn(vs ...time.Time) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldInstallDate), v...))
	})
}

// InstallDateGT applies the GT predicate on the "install_date" field.
func InstallDateGT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldInstallDate), v))
	})
}

// InstallDateGTE applies the GTE predicate on the "install_date" field.
func InstallDateGTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldInstallDate), v))
	})
}

// InstallDateLT applies the LT predicate on the "install_date" field.
func InstallDateLT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldInstallDate), v))
	})
}

// InstallDateLTE applies the LTE predicate on the "install_date" field.
func InstallDateLTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldInstallDate), v))
	})
}

// InstallDateIsNil applies the IsNil predicate on the "install_date" field.
func InstallDateIsNil() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldInstallDate)))
	})
}

// InstallDateNotNil applies the NotNil predicate on the "install_date" field.
func InstallDateNotNil() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldInstallDate)))
	})
}

// CreationDateEQ applies the EQ predicate on the "creation_date" field.
func CreationDateEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreationDate), v))
	})
}

// CreationDateNEQ applies the NEQ predicate on the "creation_date" field.
func CreationDateNEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreationDate), v))
	})
}

// CreationDateIn applies the In predicate on the "creation_date" field.
func CreationDateIn(vs ...time.Time) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCreationDate), v...))
	})
}

// CreationDateNotIn applies the NotIn predicate on the "creation_date" field.
func CreationDateNotIn(vs ...time.Time) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCreationDate), v...))
	})
}

// CreationDateGT applies the GT predicate on the "creation_date" field.
func CreationDateGT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreationDate), v))
	})
}

// CreationDateGTE applies the GTE predicate on the "creation_date" field.
func CreationDateGTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreationDate), v))
	})
}

// CreationDateLT applies the LT predicate on the "creation_date" field.
func CreationDateLT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreationDate), v))
	})
}

// CreationDateLTE applies the LTE predicate on the "creation_date" field.
func CreationDateLTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreationDate), v))
	})
}

// IndexEQ applies the EQ predicate on the "index" field.
func IndexEQ(v int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	})
}

// IndexNEQ applies the NEQ predicate on the "index" field.
func IndexNEQ(v int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldIndex), v))
	})
}

// IndexIn applies the In predicate on the "index" field.
func IndexIn(vs ...int) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
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
func IndexNotIn(vs ...int) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
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
func IndexGT(v int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldIndex), v))
	})
}

// IndexGTE applies the GTE predicate on the "index" field.
func IndexGTE(v int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldIndex), v))
	})
}

// IndexLT applies the LT predicate on the "index" field.
func IndexLT(v int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldIndex), v))
	})
}

// IndexLTE applies the LTE predicate on the "index" field.
func IndexLTE(v int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldIndex), v))
	})
}

// IndexIsNil applies the IsNil predicate on the "index" field.
func IndexIsNil() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldIndex)))
	})
}

// IndexNotNil applies the NotNil predicate on the "index" field.
func IndexNotNil() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldIndex)))
	})
}

// CloseDateEQ applies the EQ predicate on the "close_date" field.
func CloseDateEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCloseDate), v))
	})
}

// CloseDateNEQ applies the NEQ predicate on the "close_date" field.
func CloseDateNEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCloseDate), v))
	})
}

// CloseDateIn applies the In predicate on the "close_date" field.
func CloseDateIn(vs ...time.Time) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCloseDate), v...))
	})
}

// CloseDateNotIn applies the NotIn predicate on the "close_date" field.
func CloseDateNotIn(vs ...time.Time) predicate.WorkOrder {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.WorkOrder(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCloseDate), v...))
	})
}

// CloseDateGT applies the GT predicate on the "close_date" field.
func CloseDateGT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCloseDate), v))
	})
}

// CloseDateGTE applies the GTE predicate on the "close_date" field.
func CloseDateGTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCloseDate), v))
	})
}

// CloseDateLT applies the LT predicate on the "close_date" field.
func CloseDateLT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCloseDate), v))
	})
}

// CloseDateLTE applies the LTE predicate on the "close_date" field.
func CloseDateLTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCloseDate), v))
	})
}

// CloseDateIsNil applies the IsNil predicate on the "close_date" field.
func CloseDateIsNil() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldCloseDate)))
	})
}

// CloseDateNotNil applies the NotNil predicate on the "close_date" field.
func CloseDateNotNil() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldCloseDate)))
	})
}

// HasType applies the HasEdge predicate on the "type" edge.
func HasType() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, TypeTable, TypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTypeWith applies the HasEdge predicate on the "type" edge with a given conditions (other predicates).
func HasTypeWith(preds ...predicate.WorkOrderType) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
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

// HasEquipment applies the HasEdge predicate on the "equipment" edge.
func HasEquipment() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, EquipmentTable, EquipmentColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEquipmentWith applies the HasEdge predicate on the "equipment" edge with a given conditions (other predicates).
func HasEquipmentWith(preds ...predicate.Equipment) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, EquipmentTable, EquipmentColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasLinks applies the HasEdge predicate on the "links" edge.
func HasLinks() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LinksTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, LinksTable, LinksColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLinksWith applies the HasEdge predicate on the "links" edge with a given conditions (other predicates).
func HasLinksWith(preds ...predicate.Link) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LinksInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, LinksTable, LinksColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasFiles applies the HasEdge predicate on the "files" edge.
func HasFiles() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(FilesTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, FilesTable, FilesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasFilesWith applies the HasEdge predicate on the "files" edge with a given conditions (other predicates).
func HasFilesWith(preds ...predicate.File) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(FilesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, FilesTable, FilesColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasHyperlinks applies the HasEdge predicate on the "hyperlinks" edge.
func HasHyperlinks() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(HyperlinksTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, HyperlinksTable, HyperlinksColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasHyperlinksWith applies the HasEdge predicate on the "hyperlinks" edge with a given conditions (other predicates).
func HasHyperlinksWith(preds ...predicate.Hyperlink) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(HyperlinksInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, HyperlinksTable, HyperlinksColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasLocation applies the HasEdge predicate on the "location" edge.
func HasLocation() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LocationTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, LocationTable, LocationColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLocationWith applies the HasEdge predicate on the "location" edge with a given conditions (other predicates).
func HasLocationWith(preds ...predicate.Location) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LocationInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, LocationTable, LocationColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasComments applies the HasEdge predicate on the "comments" edge.
func HasComments() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CommentsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, CommentsTable, CommentsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCommentsWith applies the HasEdge predicate on the "comments" edge with a given conditions (other predicates).
func HasCommentsWith(preds ...predicate.Comment) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CommentsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, CommentsTable, CommentsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProperties applies the HasEdge predicate on the "properties" edge.
func HasProperties() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PropertiesTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PropertiesTable, PropertiesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPropertiesWith applies the HasEdge predicate on the "properties" edge with a given conditions (other predicates).
func HasPropertiesWith(preds ...predicate.Property) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PropertiesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PropertiesTable, PropertiesColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasCheckListCategories applies the HasEdge predicate on the "check_list_categories" edge.
func HasCheckListCategories() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CheckListCategoriesTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, CheckListCategoriesTable, CheckListCategoriesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCheckListCategoriesWith applies the HasEdge predicate on the "check_list_categories" edge with a given conditions (other predicates).
func HasCheckListCategoriesWith(preds ...predicate.CheckListCategory) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CheckListCategoriesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, CheckListCategoriesTable, CheckListCategoriesColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasCheckListItems applies the HasEdge predicate on the "check_list_items" edge.
func HasCheckListItems() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CheckListItemsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, CheckListItemsTable, CheckListItemsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCheckListItemsWith applies the HasEdge predicate on the "check_list_items" edge with a given conditions (other predicates).
func HasCheckListItemsWith(preds ...predicate.CheckListItem) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CheckListItemsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, CheckListItemsTable, CheckListItemsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasTechnician applies the HasEdge predicate on the "technician" edge.
func HasTechnician() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TechnicianTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, TechnicianTable, TechnicianColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTechnicianWith applies the HasEdge predicate on the "technician" edge with a given conditions (other predicates).
func HasTechnicianWith(preds ...predicate.Technician) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TechnicianInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, TechnicianTable, TechnicianColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProject applies the HasEdge predicate on the "project" edge.
func HasProject() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProjectTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProjectTable, ProjectColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProjectWith applies the HasEdge predicate on the "project" edge with a given conditions (other predicates).
func HasProjectWith(preds ...predicate.Project) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProjectInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProjectTable, ProjectColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasOwner applies the HasEdge predicate on the "owner" edge.
func HasOwner() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(OwnerTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, OwnerTable, OwnerColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasOwnerWith applies the HasEdge predicate on the "owner" edge with a given conditions (other predicates).
func HasOwnerWith(preds ...predicate.User) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(OwnerInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, OwnerTable, OwnerColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasAssignee applies the HasEdge predicate on the "assignee" edge.
func HasAssignee() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(AssigneeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, AssigneeTable, AssigneeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasAssigneeWith applies the HasEdge predicate on the "assignee" edge with a given conditions (other predicates).
func HasAssigneeWith(preds ...predicate.User) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(AssigneeInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, AssigneeTable, AssigneeColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.WorkOrder) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.WorkOrder) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
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
func Not(p predicate.WorkOrder) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		p(s.Not())
	})
}
