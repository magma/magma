// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package checklistcategorydefinition

import (
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
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
func IDGT(id int) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// Title applies equality check predicate on the "title" field. It's identical to TitleEQ.
func Title(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTitle), v))
	})
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDescription), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.CheckListCategoryDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.CheckListCategoryDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.CheckListCategoryDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.CheckListCategoryDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// TitleEQ applies the EQ predicate on the "title" field.
func TitleEQ(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTitle), v))
	})
}

// TitleNEQ applies the NEQ predicate on the "title" field.
func TitleNEQ(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTitle), v))
	})
}

// TitleIn applies the In predicate on the "title" field.
func TitleIn(vs ...string) predicate.CheckListCategoryDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldTitle), v...))
	})
}

// TitleNotIn applies the NotIn predicate on the "title" field.
func TitleNotIn(vs ...string) predicate.CheckListCategoryDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldTitle), v...))
	})
}

// TitleGT applies the GT predicate on the "title" field.
func TitleGT(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldTitle), v))
	})
}

// TitleGTE applies the GTE predicate on the "title" field.
func TitleGTE(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldTitle), v))
	})
}

// TitleLT applies the LT predicate on the "title" field.
func TitleLT(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldTitle), v))
	})
}

// TitleLTE applies the LTE predicate on the "title" field.
func TitleLTE(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldTitle), v))
	})
}

// TitleContains applies the Contains predicate on the "title" field.
func TitleContains(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldTitle), v))
	})
}

// TitleHasPrefix applies the HasPrefix predicate on the "title" field.
func TitleHasPrefix(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldTitle), v))
	})
}

// TitleHasSuffix applies the HasSuffix predicate on the "title" field.
func TitleHasSuffix(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldTitle), v))
	})
}

// TitleEqualFold applies the EqualFold predicate on the "title" field.
func TitleEqualFold(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldTitle), v))
	})
}

// TitleContainsFold applies the ContainsFold predicate on the "title" field.
func TitleContainsFold(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldTitle), v))
	})
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDescription), v))
	})
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldDescription), v))
	})
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.CheckListCategoryDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
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
func DescriptionNotIn(vs ...string) predicate.CheckListCategoryDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
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
func DescriptionGT(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldDescription), v))
	})
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldDescription), v))
	})
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldDescription), v))
	})
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldDescription), v))
	})
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldDescription), v))
	})
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldDescription), v))
	})
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldDescription), v))
	})
}

// DescriptionIsNil applies the IsNil predicate on the "description" field.
func DescriptionIsNil() predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldDescription)))
	})
}

// DescriptionNotNil applies the NotNil predicate on the "description" field.
func DescriptionNotNil() predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldDescription)))
	})
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldDescription), v))
	})
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldDescription), v))
	})
}

// HasCheckListItemDefinitions applies the HasEdge predicate on the "check_list_item_definitions" edge.
func HasCheckListItemDefinitions() predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CheckListItemDefinitionsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, CheckListItemDefinitionsTable, CheckListItemDefinitionsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCheckListItemDefinitionsWith applies the HasEdge predicate on the "check_list_item_definitions" edge with a given conditions (other predicates).
func HasCheckListItemDefinitionsWith(preds ...predicate.CheckListItemDefinition) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CheckListItemDefinitionsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, CheckListItemDefinitionsTable, CheckListItemDefinitionsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasWorkOrderType applies the HasEdge predicate on the "work_order_type" edge.
func HasWorkOrderType() predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(WorkOrderTypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, WorkOrderTypeTable, WorkOrderTypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasWorkOrderTypeWith applies the HasEdge predicate on the "work_order_type" edge with a given conditions (other predicates).
func HasWorkOrderTypeWith(preds ...predicate.WorkOrderType) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(WorkOrderTypeInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, WorkOrderTypeTable, WorkOrderTypeColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasWorkOrderTemplate applies the HasEdge predicate on the "work_order_template" edge.
func HasWorkOrderTemplate() predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(WorkOrderTemplateTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, WorkOrderTemplateTable, WorkOrderTemplateColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasWorkOrderTemplateWith applies the HasEdge predicate on the "work_order_template" edge with a given conditions (other predicates).
func HasWorkOrderTemplateWith(preds ...predicate.WorkOrderTemplate) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(WorkOrderTemplateInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, WorkOrderTemplateTable, WorkOrderTemplateColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.CheckListCategoryDefinition) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.CheckListCategoryDefinition) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
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
func Not(p predicate.CheckListCategoryDefinition) predicate.CheckListCategoryDefinition {
	return predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
		p(s.Not())
	})
}
