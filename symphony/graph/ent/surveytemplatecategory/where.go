// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package surveytemplatecategory

import (
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
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
func IDNotIn(ids ...string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
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
func IDGT(id string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// CategoryTitle applies equality check predicate on the "category_title" field. It's identical to CategoryTitleEQ.
func CategoryTitle(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCategoryTitle), v))
	})
}

// CategoryDescription applies equality check predicate on the "category_description" field. It's identical to CategoryDescriptionEQ.
func CategoryDescription(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCategoryDescription), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.SurveyTemplateCategory {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.SurveyTemplateCategory {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.SurveyTemplateCategory {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.SurveyTemplateCategory {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// CategoryTitleEQ applies the EQ predicate on the "category_title" field.
func CategoryTitleEQ(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCategoryTitle), v))
	})
}

// CategoryTitleNEQ applies the NEQ predicate on the "category_title" field.
func CategoryTitleNEQ(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCategoryTitle), v))
	})
}

// CategoryTitleIn applies the In predicate on the "category_title" field.
func CategoryTitleIn(vs ...string) predicate.SurveyTemplateCategory {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCategoryTitle), v...))
	})
}

// CategoryTitleNotIn applies the NotIn predicate on the "category_title" field.
func CategoryTitleNotIn(vs ...string) predicate.SurveyTemplateCategory {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCategoryTitle), v...))
	})
}

// CategoryTitleGT applies the GT predicate on the "category_title" field.
func CategoryTitleGT(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCategoryTitle), v))
	})
}

// CategoryTitleGTE applies the GTE predicate on the "category_title" field.
func CategoryTitleGTE(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCategoryTitle), v))
	})
}

// CategoryTitleLT applies the LT predicate on the "category_title" field.
func CategoryTitleLT(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCategoryTitle), v))
	})
}

// CategoryTitleLTE applies the LTE predicate on the "category_title" field.
func CategoryTitleLTE(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCategoryTitle), v))
	})
}

// CategoryTitleContains applies the Contains predicate on the "category_title" field.
func CategoryTitleContains(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldCategoryTitle), v))
	})
}

// CategoryTitleHasPrefix applies the HasPrefix predicate on the "category_title" field.
func CategoryTitleHasPrefix(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldCategoryTitle), v))
	})
}

// CategoryTitleHasSuffix applies the HasSuffix predicate on the "category_title" field.
func CategoryTitleHasSuffix(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldCategoryTitle), v))
	})
}

// CategoryTitleEqualFold applies the EqualFold predicate on the "category_title" field.
func CategoryTitleEqualFold(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldCategoryTitle), v))
	})
}

// CategoryTitleContainsFold applies the ContainsFold predicate on the "category_title" field.
func CategoryTitleContainsFold(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldCategoryTitle), v))
	})
}

// CategoryDescriptionEQ applies the EQ predicate on the "category_description" field.
func CategoryDescriptionEQ(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCategoryDescription), v))
	})
}

// CategoryDescriptionNEQ applies the NEQ predicate on the "category_description" field.
func CategoryDescriptionNEQ(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCategoryDescription), v))
	})
}

// CategoryDescriptionIn applies the In predicate on the "category_description" field.
func CategoryDescriptionIn(vs ...string) predicate.SurveyTemplateCategory {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCategoryDescription), v...))
	})
}

// CategoryDescriptionNotIn applies the NotIn predicate on the "category_description" field.
func CategoryDescriptionNotIn(vs ...string) predicate.SurveyTemplateCategory {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCategoryDescription), v...))
	})
}

// CategoryDescriptionGT applies the GT predicate on the "category_description" field.
func CategoryDescriptionGT(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCategoryDescription), v))
	})
}

// CategoryDescriptionGTE applies the GTE predicate on the "category_description" field.
func CategoryDescriptionGTE(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCategoryDescription), v))
	})
}

// CategoryDescriptionLT applies the LT predicate on the "category_description" field.
func CategoryDescriptionLT(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCategoryDescription), v))
	})
}

// CategoryDescriptionLTE applies the LTE predicate on the "category_description" field.
func CategoryDescriptionLTE(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCategoryDescription), v))
	})
}

// CategoryDescriptionContains applies the Contains predicate on the "category_description" field.
func CategoryDescriptionContains(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldCategoryDescription), v))
	})
}

// CategoryDescriptionHasPrefix applies the HasPrefix predicate on the "category_description" field.
func CategoryDescriptionHasPrefix(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldCategoryDescription), v))
	})
}

// CategoryDescriptionHasSuffix applies the HasSuffix predicate on the "category_description" field.
func CategoryDescriptionHasSuffix(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldCategoryDescription), v))
	})
}

// CategoryDescriptionEqualFold applies the EqualFold predicate on the "category_description" field.
func CategoryDescriptionEqualFold(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldCategoryDescription), v))
	})
}

// CategoryDescriptionContainsFold applies the ContainsFold predicate on the "category_description" field.
func CategoryDescriptionContainsFold(v string) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldCategoryDescription), v))
	})
}

// HasSurveyTemplateQuestions applies the HasEdge predicate on the "survey_template_questions" edge.
func HasSurveyTemplateQuestions() predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(SurveyTemplateQuestionsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, SurveyTemplateQuestionsTable, SurveyTemplateQuestionsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSurveyTemplateQuestionsWith applies the HasEdge predicate on the "survey_template_questions" edge with a given conditions (other predicates).
func HasSurveyTemplateQuestionsWith(preds ...predicate.SurveyTemplateQuestion) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(SurveyTemplateQuestionsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, SurveyTemplateQuestionsTable, SurveyTemplateQuestionsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.SurveyTemplateCategory) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.SurveyTemplateCategory) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
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
func Not(p predicate.SurveyTemplateCategory) predicate.SurveyTemplateCategory {
	return predicate.SurveyTemplateCategory(func(s *sql.Selector) {
		p(s.Not())
	})
}
