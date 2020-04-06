// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package surveytemplatequestion

import (
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
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
func IDGT(id int) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// QuestionTitle applies equality check predicate on the "question_title" field. It's identical to QuestionTitleEQ.
func QuestionTitle(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldQuestionTitle), v))
	})
}

// QuestionDescription applies equality check predicate on the "question_description" field. It's identical to QuestionDescriptionEQ.
func QuestionDescription(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldQuestionDescription), v))
	})
}

// QuestionType applies equality check predicate on the "question_type" field. It's identical to QuestionTypeEQ.
func QuestionType(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldQuestionType), v))
	})
}

// Index applies equality check predicate on the "index" field. It's identical to IndexEQ.
func Index(v int) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.SurveyTemplateQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.SurveyTemplateQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.SurveyTemplateQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.SurveyTemplateQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// QuestionTitleEQ applies the EQ predicate on the "question_title" field.
func QuestionTitleEQ(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldQuestionTitle), v))
	})
}

// QuestionTitleNEQ applies the NEQ predicate on the "question_title" field.
func QuestionTitleNEQ(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldQuestionTitle), v))
	})
}

// QuestionTitleIn applies the In predicate on the "question_title" field.
func QuestionTitleIn(vs ...string) predicate.SurveyTemplateQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldQuestionTitle), v...))
	})
}

// QuestionTitleNotIn applies the NotIn predicate on the "question_title" field.
func QuestionTitleNotIn(vs ...string) predicate.SurveyTemplateQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldQuestionTitle), v...))
	})
}

// QuestionTitleGT applies the GT predicate on the "question_title" field.
func QuestionTitleGT(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldQuestionTitle), v))
	})
}

// QuestionTitleGTE applies the GTE predicate on the "question_title" field.
func QuestionTitleGTE(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldQuestionTitle), v))
	})
}

// QuestionTitleLT applies the LT predicate on the "question_title" field.
func QuestionTitleLT(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldQuestionTitle), v))
	})
}

// QuestionTitleLTE applies the LTE predicate on the "question_title" field.
func QuestionTitleLTE(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldQuestionTitle), v))
	})
}

// QuestionTitleContains applies the Contains predicate on the "question_title" field.
func QuestionTitleContains(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldQuestionTitle), v))
	})
}

// QuestionTitleHasPrefix applies the HasPrefix predicate on the "question_title" field.
func QuestionTitleHasPrefix(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldQuestionTitle), v))
	})
}

// QuestionTitleHasSuffix applies the HasSuffix predicate on the "question_title" field.
func QuestionTitleHasSuffix(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldQuestionTitle), v))
	})
}

// QuestionTitleEqualFold applies the EqualFold predicate on the "question_title" field.
func QuestionTitleEqualFold(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldQuestionTitle), v))
	})
}

// QuestionTitleContainsFold applies the ContainsFold predicate on the "question_title" field.
func QuestionTitleContainsFold(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldQuestionTitle), v))
	})
}

// QuestionDescriptionEQ applies the EQ predicate on the "question_description" field.
func QuestionDescriptionEQ(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldQuestionDescription), v))
	})
}

// QuestionDescriptionNEQ applies the NEQ predicate on the "question_description" field.
func QuestionDescriptionNEQ(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldQuestionDescription), v))
	})
}

// QuestionDescriptionIn applies the In predicate on the "question_description" field.
func QuestionDescriptionIn(vs ...string) predicate.SurveyTemplateQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldQuestionDescription), v...))
	})
}

// QuestionDescriptionNotIn applies the NotIn predicate on the "question_description" field.
func QuestionDescriptionNotIn(vs ...string) predicate.SurveyTemplateQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldQuestionDescription), v...))
	})
}

// QuestionDescriptionGT applies the GT predicate on the "question_description" field.
func QuestionDescriptionGT(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldQuestionDescription), v))
	})
}

// QuestionDescriptionGTE applies the GTE predicate on the "question_description" field.
func QuestionDescriptionGTE(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldQuestionDescription), v))
	})
}

// QuestionDescriptionLT applies the LT predicate on the "question_description" field.
func QuestionDescriptionLT(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldQuestionDescription), v))
	})
}

// QuestionDescriptionLTE applies the LTE predicate on the "question_description" field.
func QuestionDescriptionLTE(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldQuestionDescription), v))
	})
}

// QuestionDescriptionContains applies the Contains predicate on the "question_description" field.
func QuestionDescriptionContains(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldQuestionDescription), v))
	})
}

// QuestionDescriptionHasPrefix applies the HasPrefix predicate on the "question_description" field.
func QuestionDescriptionHasPrefix(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldQuestionDescription), v))
	})
}

// QuestionDescriptionHasSuffix applies the HasSuffix predicate on the "question_description" field.
func QuestionDescriptionHasSuffix(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldQuestionDescription), v))
	})
}

// QuestionDescriptionEqualFold applies the EqualFold predicate on the "question_description" field.
func QuestionDescriptionEqualFold(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldQuestionDescription), v))
	})
}

// QuestionDescriptionContainsFold applies the ContainsFold predicate on the "question_description" field.
func QuestionDescriptionContainsFold(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldQuestionDescription), v))
	})
}

// QuestionTypeEQ applies the EQ predicate on the "question_type" field.
func QuestionTypeEQ(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeNEQ applies the NEQ predicate on the "question_type" field.
func QuestionTypeNEQ(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeIn applies the In predicate on the "question_type" field.
func QuestionTypeIn(vs ...string) predicate.SurveyTemplateQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldQuestionType), v...))
	})
}

// QuestionTypeNotIn applies the NotIn predicate on the "question_type" field.
func QuestionTypeNotIn(vs ...string) predicate.SurveyTemplateQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldQuestionType), v...))
	})
}

// QuestionTypeGT applies the GT predicate on the "question_type" field.
func QuestionTypeGT(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeGTE applies the GTE predicate on the "question_type" field.
func QuestionTypeGTE(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeLT applies the LT predicate on the "question_type" field.
func QuestionTypeLT(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeLTE applies the LTE predicate on the "question_type" field.
func QuestionTypeLTE(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeContains applies the Contains predicate on the "question_type" field.
func QuestionTypeContains(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeHasPrefix applies the HasPrefix predicate on the "question_type" field.
func QuestionTypeHasPrefix(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeHasSuffix applies the HasSuffix predicate on the "question_type" field.
func QuestionTypeHasSuffix(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeEqualFold applies the EqualFold predicate on the "question_type" field.
func QuestionTypeEqualFold(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeContainsFold applies the ContainsFold predicate on the "question_type" field.
func QuestionTypeContainsFold(v string) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldQuestionType), v))
	})
}

// IndexEQ applies the EQ predicate on the "index" field.
func IndexEQ(v int) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	})
}

// IndexNEQ applies the NEQ predicate on the "index" field.
func IndexNEQ(v int) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldIndex), v))
	})
}

// IndexIn applies the In predicate on the "index" field.
func IndexIn(vs ...int) predicate.SurveyTemplateQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
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
func IndexNotIn(vs ...int) predicate.SurveyTemplateQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
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
func IndexGT(v int) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldIndex), v))
	})
}

// IndexGTE applies the GTE predicate on the "index" field.
func IndexGTE(v int) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldIndex), v))
	})
}

// IndexLT applies the LT predicate on the "index" field.
func IndexLT(v int) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldIndex), v))
	})
}

// IndexLTE applies the LTE predicate on the "index" field.
func IndexLTE(v int) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldIndex), v))
	})
}

// HasCategory applies the HasEdge predicate on the "category" edge.
func HasCategory() predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CategoryTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, CategoryTable, CategoryColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCategoryWith applies the HasEdge predicate on the "category" edge with a given conditions (other predicates).
func HasCategoryWith(preds ...predicate.SurveyTemplateCategory) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CategoryInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, CategoryTable, CategoryColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.SurveyTemplateQuestion) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.SurveyTemplateQuestion) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
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
func Not(p predicate.SurveyTemplateQuestion) predicate.SurveyTemplateQuestion {
	return predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
		p(s.Not())
	})
}
