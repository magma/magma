// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package workorder

import (
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id string) predicate.WorkOrder {
	return predicate.WorkOrder(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.EQ(s.C(FieldID), id))
		},
	)
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	},
	)
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.NEQ(s.C(FieldID), id))
	},
	)
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
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
	},
	)
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
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
	},
	)
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GT(s.C(FieldID), id))
	},
	)
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GTE(s.C(FieldID), id))
	},
	)
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LT(s.C(FieldID), id))
	},
	)
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LTE(s.C(FieldID), id))
	},
	)
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	},
	)
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	},
	)
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	},
	)
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStatus), v))
	},
	)
}

// Priority applies equality check predicate on the "priority" field. It's identical to PriorityEQ.
func Priority(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPriority), v))
	},
	)
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDescription), v))
	},
	)
}

// OwnerName applies equality check predicate on the "owner_name" field. It's identical to OwnerNameEQ.
func OwnerName(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldOwnerName), v))
	},
	)
}

// InstallDate applies equality check predicate on the "install_date" field. It's identical to InstallDateEQ.
func InstallDate(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldInstallDate), v))
	},
	)
}

// CreationDate applies equality check predicate on the "creation_date" field. It's identical to CreationDateEQ.
func CreationDate(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreationDate), v))
	},
	)
}

// Assignee applies equality check predicate on the "assignee" field. It's identical to AssigneeEQ.
func Assignee(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldAssignee), v))
	},
	)
}

// Index applies equality check predicate on the "index" field. It's identical to IndexEQ.
func Index(v int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	},
	)
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	},
	)
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	},
	)
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
	},
	)
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
	},
	)
}

// CreateTimeGT applies the GT predicate on the "create_time" field.
func CreateTimeGT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	},
	)
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	},
	)
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	},
	)
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	},
	)
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	},
	)
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	},
	)
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
	},
	)
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
	},
	)
}

// UpdateTimeGT applies the GT predicate on the "update_time" field.
func UpdateTimeGT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	},
	)
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	},
	)
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	},
	)
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	},
	)
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	},
	)
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	},
	)
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
	},
	)
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
	},
	)
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	},
	)
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	},
	)
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	},
	)
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	},
	)
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	},
	)
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	},
	)
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	},
	)
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	},
	)
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	},
	)
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStatus), v))
	},
	)
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldStatus), v))
	},
	)
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
	},
	)
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
	},
	)
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldStatus), v))
	},
	)
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldStatus), v))
	},
	)
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldStatus), v))
	},
	)
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldStatus), v))
	},
	)
}

// StatusContains applies the Contains predicate on the "status" field.
func StatusContains(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldStatus), v))
	},
	)
}

// StatusHasPrefix applies the HasPrefix predicate on the "status" field.
func StatusHasPrefix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldStatus), v))
	},
	)
}

// StatusHasSuffix applies the HasSuffix predicate on the "status" field.
func StatusHasSuffix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldStatus), v))
	},
	)
}

// StatusEqualFold applies the EqualFold predicate on the "status" field.
func StatusEqualFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldStatus), v))
	},
	)
}

// StatusContainsFold applies the ContainsFold predicate on the "status" field.
func StatusContainsFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldStatus), v))
	},
	)
}

// PriorityEQ applies the EQ predicate on the "priority" field.
func PriorityEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPriority), v))
	},
	)
}

// PriorityNEQ applies the NEQ predicate on the "priority" field.
func PriorityNEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldPriority), v))
	},
	)
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
	},
	)
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
	},
	)
}

// PriorityGT applies the GT predicate on the "priority" field.
func PriorityGT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldPriority), v))
	},
	)
}

// PriorityGTE applies the GTE predicate on the "priority" field.
func PriorityGTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldPriority), v))
	},
	)
}

// PriorityLT applies the LT predicate on the "priority" field.
func PriorityLT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldPriority), v))
	},
	)
}

// PriorityLTE applies the LTE predicate on the "priority" field.
func PriorityLTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldPriority), v))
	},
	)
}

// PriorityContains applies the Contains predicate on the "priority" field.
func PriorityContains(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldPriority), v))
	},
	)
}

// PriorityHasPrefix applies the HasPrefix predicate on the "priority" field.
func PriorityHasPrefix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldPriority), v))
	},
	)
}

// PriorityHasSuffix applies the HasSuffix predicate on the "priority" field.
func PriorityHasSuffix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldPriority), v))
	},
	)
}

// PriorityEqualFold applies the EqualFold predicate on the "priority" field.
func PriorityEqualFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldPriority), v))
	},
	)
}

// PriorityContainsFold applies the ContainsFold predicate on the "priority" field.
func PriorityContainsFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldPriority), v))
	},
	)
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDescription), v))
	},
	)
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldDescription), v))
	},
	)
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
	},
	)
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
	},
	)
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldDescription), v))
	},
	)
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldDescription), v))
	},
	)
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldDescription), v))
	},
	)
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldDescription), v))
	},
	)
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldDescription), v))
	},
	)
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldDescription), v))
	},
	)
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldDescription), v))
	},
	)
}

// DescriptionIsNil applies the IsNil predicate on the "description" field.
func DescriptionIsNil() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldDescription)))
	},
	)
}

// DescriptionNotNil applies the NotNil predicate on the "description" field.
func DescriptionNotNil() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldDescription)))
	},
	)
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldDescription), v))
	},
	)
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldDescription), v))
	},
	)
}

// OwnerNameEQ applies the EQ predicate on the "owner_name" field.
func OwnerNameEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldOwnerName), v))
	},
	)
}

// OwnerNameNEQ applies the NEQ predicate on the "owner_name" field.
func OwnerNameNEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldOwnerName), v))
	},
	)
}

// OwnerNameIn applies the In predicate on the "owner_name" field.
func OwnerNameIn(vs ...string) predicate.WorkOrder {
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
		s.Where(sql.In(s.C(FieldOwnerName), v...))
	},
	)
}

// OwnerNameNotIn applies the NotIn predicate on the "owner_name" field.
func OwnerNameNotIn(vs ...string) predicate.WorkOrder {
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
		s.Where(sql.NotIn(s.C(FieldOwnerName), v...))
	},
	)
}

// OwnerNameGT applies the GT predicate on the "owner_name" field.
func OwnerNameGT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldOwnerName), v))
	},
	)
}

// OwnerNameGTE applies the GTE predicate on the "owner_name" field.
func OwnerNameGTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldOwnerName), v))
	},
	)
}

// OwnerNameLT applies the LT predicate on the "owner_name" field.
func OwnerNameLT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldOwnerName), v))
	},
	)
}

// OwnerNameLTE applies the LTE predicate on the "owner_name" field.
func OwnerNameLTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldOwnerName), v))
	},
	)
}

// OwnerNameContains applies the Contains predicate on the "owner_name" field.
func OwnerNameContains(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldOwnerName), v))
	},
	)
}

// OwnerNameHasPrefix applies the HasPrefix predicate on the "owner_name" field.
func OwnerNameHasPrefix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldOwnerName), v))
	},
	)
}

// OwnerNameHasSuffix applies the HasSuffix predicate on the "owner_name" field.
func OwnerNameHasSuffix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldOwnerName), v))
	},
	)
}

// OwnerNameEqualFold applies the EqualFold predicate on the "owner_name" field.
func OwnerNameEqualFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldOwnerName), v))
	},
	)
}

// OwnerNameContainsFold applies the ContainsFold predicate on the "owner_name" field.
func OwnerNameContainsFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldOwnerName), v))
	},
	)
}

// InstallDateEQ applies the EQ predicate on the "install_date" field.
func InstallDateEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldInstallDate), v))
	},
	)
}

// InstallDateNEQ applies the NEQ predicate on the "install_date" field.
func InstallDateNEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldInstallDate), v))
	},
	)
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
	},
	)
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
	},
	)
}

// InstallDateGT applies the GT predicate on the "install_date" field.
func InstallDateGT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldInstallDate), v))
	},
	)
}

// InstallDateGTE applies the GTE predicate on the "install_date" field.
func InstallDateGTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldInstallDate), v))
	},
	)
}

// InstallDateLT applies the LT predicate on the "install_date" field.
func InstallDateLT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldInstallDate), v))
	},
	)
}

// InstallDateLTE applies the LTE predicate on the "install_date" field.
func InstallDateLTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldInstallDate), v))
	},
	)
}

// InstallDateIsNil applies the IsNil predicate on the "install_date" field.
func InstallDateIsNil() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldInstallDate)))
	},
	)
}

// InstallDateNotNil applies the NotNil predicate on the "install_date" field.
func InstallDateNotNil() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldInstallDate)))
	},
	)
}

// CreationDateEQ applies the EQ predicate on the "creation_date" field.
func CreationDateEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreationDate), v))
	},
	)
}

// CreationDateNEQ applies the NEQ predicate on the "creation_date" field.
func CreationDateNEQ(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreationDate), v))
	},
	)
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
	},
	)
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
	},
	)
}

// CreationDateGT applies the GT predicate on the "creation_date" field.
func CreationDateGT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreationDate), v))
	},
	)
}

// CreationDateGTE applies the GTE predicate on the "creation_date" field.
func CreationDateGTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreationDate), v))
	},
	)
}

// CreationDateLT applies the LT predicate on the "creation_date" field.
func CreationDateLT(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreationDate), v))
	},
	)
}

// CreationDateLTE applies the LTE predicate on the "creation_date" field.
func CreationDateLTE(v time.Time) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreationDate), v))
	},
	)
}

// AssigneeEQ applies the EQ predicate on the "assignee" field.
func AssigneeEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldAssignee), v))
	},
	)
}

// AssigneeNEQ applies the NEQ predicate on the "assignee" field.
func AssigneeNEQ(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldAssignee), v))
	},
	)
}

// AssigneeIn applies the In predicate on the "assignee" field.
func AssigneeIn(vs ...string) predicate.WorkOrder {
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
		s.Where(sql.In(s.C(FieldAssignee), v...))
	},
	)
}

// AssigneeNotIn applies the NotIn predicate on the "assignee" field.
func AssigneeNotIn(vs ...string) predicate.WorkOrder {
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
		s.Where(sql.NotIn(s.C(FieldAssignee), v...))
	},
	)
}

// AssigneeGT applies the GT predicate on the "assignee" field.
func AssigneeGT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldAssignee), v))
	},
	)
}

// AssigneeGTE applies the GTE predicate on the "assignee" field.
func AssigneeGTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldAssignee), v))
	},
	)
}

// AssigneeLT applies the LT predicate on the "assignee" field.
func AssigneeLT(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldAssignee), v))
	},
	)
}

// AssigneeLTE applies the LTE predicate on the "assignee" field.
func AssigneeLTE(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldAssignee), v))
	},
	)
}

// AssigneeContains applies the Contains predicate on the "assignee" field.
func AssigneeContains(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldAssignee), v))
	},
	)
}

// AssigneeHasPrefix applies the HasPrefix predicate on the "assignee" field.
func AssigneeHasPrefix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldAssignee), v))
	},
	)
}

// AssigneeHasSuffix applies the HasSuffix predicate on the "assignee" field.
func AssigneeHasSuffix(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldAssignee), v))
	},
	)
}

// AssigneeIsNil applies the IsNil predicate on the "assignee" field.
func AssigneeIsNil() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldAssignee)))
	},
	)
}

// AssigneeNotNil applies the NotNil predicate on the "assignee" field.
func AssigneeNotNil() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldAssignee)))
	},
	)
}

// AssigneeEqualFold applies the EqualFold predicate on the "assignee" field.
func AssigneeEqualFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldAssignee), v))
	},
	)
}

// AssigneeContainsFold applies the ContainsFold predicate on the "assignee" field.
func AssigneeContainsFold(v string) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldAssignee), v))
	},
	)
}

// IndexEQ applies the EQ predicate on the "index" field.
func IndexEQ(v int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	},
	)
}

// IndexNEQ applies the NEQ predicate on the "index" field.
func IndexNEQ(v int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldIndex), v))
	},
	)
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
	},
	)
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
	},
	)
}

// IndexGT applies the GT predicate on the "index" field.
func IndexGT(v int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldIndex), v))
	},
	)
}

// IndexGTE applies the GTE predicate on the "index" field.
func IndexGTE(v int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldIndex), v))
	},
	)
}

// IndexLT applies the LT predicate on the "index" field.
func IndexLT(v int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldIndex), v))
	},
	)
}

// IndexLTE applies the LTE predicate on the "index" field.
func IndexLTE(v int) predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldIndex), v))
	},
	)
}

// IndexIsNil applies the IsNil predicate on the "index" field.
func IndexIsNil() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldIndex)))
	},
	)
}

// IndexNotNil applies the NotNil predicate on the "index" field.
func IndexNotNil() predicate.WorkOrder {
	return predicate.WorkOrder(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldIndex)))
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
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
	},
	)
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.WorkOrder) predicate.WorkOrder {
	return predicate.WorkOrder(
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
func Or(predicates ...predicate.WorkOrder) predicate.WorkOrder {
	return predicate.WorkOrder(
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
func Not(p predicate.WorkOrder) predicate.WorkOrder {
	return predicate.WorkOrder(
		func(s *sql.Selector) {
			p(s.Not())
		},
	)
}
