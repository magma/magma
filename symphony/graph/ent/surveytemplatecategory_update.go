// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatecategory"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatequestion"
)

// SurveyTemplateCategoryUpdate is the builder for updating SurveyTemplateCategory entities.
type SurveyTemplateCategoryUpdate struct {
	config

	update_time                    *time.Time
	category_title                 *string
	category_description           *string
	survey_template_questions      map[string]struct{}
	removedSurveyTemplateQuestions map[string]struct{}
	predicates                     []predicate.SurveyTemplateCategory
}

// Where adds a new predicate for the builder.
func (stcu *SurveyTemplateCategoryUpdate) Where(ps ...predicate.SurveyTemplateCategory) *SurveyTemplateCategoryUpdate {
	stcu.predicates = append(stcu.predicates, ps...)
	return stcu
}

// SetCategoryTitle sets the category_title field.
func (stcu *SurveyTemplateCategoryUpdate) SetCategoryTitle(s string) *SurveyTemplateCategoryUpdate {
	stcu.category_title = &s
	return stcu
}

// SetCategoryDescription sets the category_description field.
func (stcu *SurveyTemplateCategoryUpdate) SetCategoryDescription(s string) *SurveyTemplateCategoryUpdate {
	stcu.category_description = &s
	return stcu
}

// AddSurveyTemplateQuestionIDs adds the survey_template_questions edge to SurveyTemplateQuestion by ids.
func (stcu *SurveyTemplateCategoryUpdate) AddSurveyTemplateQuestionIDs(ids ...string) *SurveyTemplateCategoryUpdate {
	if stcu.survey_template_questions == nil {
		stcu.survey_template_questions = make(map[string]struct{})
	}
	for i := range ids {
		stcu.survey_template_questions[ids[i]] = struct{}{}
	}
	return stcu
}

// AddSurveyTemplateQuestions adds the survey_template_questions edges to SurveyTemplateQuestion.
func (stcu *SurveyTemplateCategoryUpdate) AddSurveyTemplateQuestions(s ...*SurveyTemplateQuestion) *SurveyTemplateCategoryUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stcu.AddSurveyTemplateQuestionIDs(ids...)
}

// RemoveSurveyTemplateQuestionIDs removes the survey_template_questions edge to SurveyTemplateQuestion by ids.
func (stcu *SurveyTemplateCategoryUpdate) RemoveSurveyTemplateQuestionIDs(ids ...string) *SurveyTemplateCategoryUpdate {
	if stcu.removedSurveyTemplateQuestions == nil {
		stcu.removedSurveyTemplateQuestions = make(map[string]struct{})
	}
	for i := range ids {
		stcu.removedSurveyTemplateQuestions[ids[i]] = struct{}{}
	}
	return stcu
}

// RemoveSurveyTemplateQuestions removes survey_template_questions edges to SurveyTemplateQuestion.
func (stcu *SurveyTemplateCategoryUpdate) RemoveSurveyTemplateQuestions(s ...*SurveyTemplateQuestion) *SurveyTemplateCategoryUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stcu.RemoveSurveyTemplateQuestionIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (stcu *SurveyTemplateCategoryUpdate) Save(ctx context.Context) (int, error) {
	if stcu.update_time == nil {
		v := surveytemplatecategory.UpdateDefaultUpdateTime()
		stcu.update_time = &v
	}
	return stcu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (stcu *SurveyTemplateCategoryUpdate) SaveX(ctx context.Context) int {
	affected, err := stcu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (stcu *SurveyTemplateCategoryUpdate) Exec(ctx context.Context) error {
	_, err := stcu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (stcu *SurveyTemplateCategoryUpdate) ExecX(ctx context.Context) {
	if err := stcu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (stcu *SurveyTemplateCategoryUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(stcu.driver.Dialect())
		selector = builder.Select(surveytemplatecategory.FieldID).From(builder.Table(surveytemplatecategory.Table))
	)
	for _, p := range stcu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = stcu.driver.Query(ctx, query, args, rows); err != nil {
		return 0, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ent: failed reading id: %v", err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return 0, nil
	}

	tx, err := stcu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(surveytemplatecategory.Table).Where(sql.InInts(surveytemplatecategory.FieldID, ids...))
	)
	if value := stcu.update_time; value != nil {
		updater.Set(surveytemplatecategory.FieldUpdateTime, *value)
	}
	if value := stcu.category_title; value != nil {
		updater.Set(surveytemplatecategory.FieldCategoryTitle, *value)
	}
	if value := stcu.category_description; value != nil {
		updater.Set(surveytemplatecategory.FieldCategoryDescription, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(stcu.removedSurveyTemplateQuestions) > 0 {
		eids := make([]int, len(stcu.removedSurveyTemplateQuestions))
		for eid := range stcu.removedSurveyTemplateQuestions {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(surveytemplatecategory.SurveyTemplateQuestionsTable).
			SetNull(surveytemplatecategory.SurveyTemplateQuestionsColumn).
			Where(sql.InInts(surveytemplatecategory.SurveyTemplateQuestionsColumn, ids...)).
			Where(sql.InInts(surveytemplatequestion.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(stcu.survey_template_questions) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range stcu.survey_template_questions {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(surveytemplatequestion.FieldID, eid)
			}
			query, args := builder.Update(surveytemplatecategory.SurveyTemplateQuestionsTable).
				Set(surveytemplatecategory.SurveyTemplateQuestionsColumn, id).
				Where(sql.And(p, sql.IsNull(surveytemplatecategory.SurveyTemplateQuestionsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(stcu.survey_template_questions) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"survey_template_questions\" %v already connected to a different \"SurveyTemplateCategory\"", keys(stcu.survey_template_questions))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// SurveyTemplateCategoryUpdateOne is the builder for updating a single SurveyTemplateCategory entity.
type SurveyTemplateCategoryUpdateOne struct {
	config
	id string

	update_time                    *time.Time
	category_title                 *string
	category_description           *string
	survey_template_questions      map[string]struct{}
	removedSurveyTemplateQuestions map[string]struct{}
}

// SetCategoryTitle sets the category_title field.
func (stcuo *SurveyTemplateCategoryUpdateOne) SetCategoryTitle(s string) *SurveyTemplateCategoryUpdateOne {
	stcuo.category_title = &s
	return stcuo
}

// SetCategoryDescription sets the category_description field.
func (stcuo *SurveyTemplateCategoryUpdateOne) SetCategoryDescription(s string) *SurveyTemplateCategoryUpdateOne {
	stcuo.category_description = &s
	return stcuo
}

// AddSurveyTemplateQuestionIDs adds the survey_template_questions edge to SurveyTemplateQuestion by ids.
func (stcuo *SurveyTemplateCategoryUpdateOne) AddSurveyTemplateQuestionIDs(ids ...string) *SurveyTemplateCategoryUpdateOne {
	if stcuo.survey_template_questions == nil {
		stcuo.survey_template_questions = make(map[string]struct{})
	}
	for i := range ids {
		stcuo.survey_template_questions[ids[i]] = struct{}{}
	}
	return stcuo
}

// AddSurveyTemplateQuestions adds the survey_template_questions edges to SurveyTemplateQuestion.
func (stcuo *SurveyTemplateCategoryUpdateOne) AddSurveyTemplateQuestions(s ...*SurveyTemplateQuestion) *SurveyTemplateCategoryUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stcuo.AddSurveyTemplateQuestionIDs(ids...)
}

// RemoveSurveyTemplateQuestionIDs removes the survey_template_questions edge to SurveyTemplateQuestion by ids.
func (stcuo *SurveyTemplateCategoryUpdateOne) RemoveSurveyTemplateQuestionIDs(ids ...string) *SurveyTemplateCategoryUpdateOne {
	if stcuo.removedSurveyTemplateQuestions == nil {
		stcuo.removedSurveyTemplateQuestions = make(map[string]struct{})
	}
	for i := range ids {
		stcuo.removedSurveyTemplateQuestions[ids[i]] = struct{}{}
	}
	return stcuo
}

// RemoveSurveyTemplateQuestions removes survey_template_questions edges to SurveyTemplateQuestion.
func (stcuo *SurveyTemplateCategoryUpdateOne) RemoveSurveyTemplateQuestions(s ...*SurveyTemplateQuestion) *SurveyTemplateCategoryUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stcuo.RemoveSurveyTemplateQuestionIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (stcuo *SurveyTemplateCategoryUpdateOne) Save(ctx context.Context) (*SurveyTemplateCategory, error) {
	if stcuo.update_time == nil {
		v := surveytemplatecategory.UpdateDefaultUpdateTime()
		stcuo.update_time = &v
	}
	return stcuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (stcuo *SurveyTemplateCategoryUpdateOne) SaveX(ctx context.Context) *SurveyTemplateCategory {
	stc, err := stcuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return stc
}

// Exec executes the query on the entity.
func (stcuo *SurveyTemplateCategoryUpdateOne) Exec(ctx context.Context) error {
	_, err := stcuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (stcuo *SurveyTemplateCategoryUpdateOne) ExecX(ctx context.Context) {
	if err := stcuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (stcuo *SurveyTemplateCategoryUpdateOne) sqlSave(ctx context.Context) (stc *SurveyTemplateCategory, err error) {
	var (
		builder  = sql.Dialect(stcuo.driver.Dialect())
		selector = builder.Select(surveytemplatecategory.Columns...).From(builder.Table(surveytemplatecategory.Table))
	)
	surveytemplatecategory.ID(stcuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = stcuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		stc = &SurveyTemplateCategory{config: stcuo.config}
		if err := stc.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into SurveyTemplateCategory: %v", err)
		}
		id = stc.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("SurveyTemplateCategory with id: %v", stcuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one SurveyTemplateCategory with the same id: %v", stcuo.id)
	}

	tx, err := stcuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(surveytemplatecategory.Table).Where(sql.InInts(surveytemplatecategory.FieldID, ids...))
	)
	if value := stcuo.update_time; value != nil {
		updater.Set(surveytemplatecategory.FieldUpdateTime, *value)
		stc.UpdateTime = *value
	}
	if value := stcuo.category_title; value != nil {
		updater.Set(surveytemplatecategory.FieldCategoryTitle, *value)
		stc.CategoryTitle = *value
	}
	if value := stcuo.category_description; value != nil {
		updater.Set(surveytemplatecategory.FieldCategoryDescription, *value)
		stc.CategoryDescription = *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(stcuo.removedSurveyTemplateQuestions) > 0 {
		eids := make([]int, len(stcuo.removedSurveyTemplateQuestions))
		for eid := range stcuo.removedSurveyTemplateQuestions {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(surveytemplatecategory.SurveyTemplateQuestionsTable).
			SetNull(surveytemplatecategory.SurveyTemplateQuestionsColumn).
			Where(sql.InInts(surveytemplatecategory.SurveyTemplateQuestionsColumn, ids...)).
			Where(sql.InInts(surveytemplatequestion.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(stcuo.survey_template_questions) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range stcuo.survey_template_questions {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(surveytemplatequestion.FieldID, eid)
			}
			query, args := builder.Update(surveytemplatecategory.SurveyTemplateQuestionsTable).
				Set(surveytemplatecategory.SurveyTemplateQuestionsColumn, id).
				Where(sql.And(p, sql.IsNull(surveytemplatecategory.SurveyTemplateQuestionsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(stcuo.survey_template_questions) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"survey_template_questions\" %v already connected to a different \"SurveyTemplateCategory\"", keys(stcuo.survey_template_questions))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return stc, nil
}
