// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatecategory"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatequestion"
)

// SurveyTemplateQuestionUpdate is the builder for updating SurveyTemplateQuestion entities.
type SurveyTemplateQuestionUpdate struct {
	config

	update_time          *time.Time
	question_title       *string
	question_description *string
	question_type        *string
	index                *int
	addindex             *int
	category             map[string]struct{}
	clearedCategory      bool
	predicates           []predicate.SurveyTemplateQuestion
}

// Where adds a new predicate for the builder.
func (stqu *SurveyTemplateQuestionUpdate) Where(ps ...predicate.SurveyTemplateQuestion) *SurveyTemplateQuestionUpdate {
	stqu.predicates = append(stqu.predicates, ps...)
	return stqu
}

// SetQuestionTitle sets the question_title field.
func (stqu *SurveyTemplateQuestionUpdate) SetQuestionTitle(s string) *SurveyTemplateQuestionUpdate {
	stqu.question_title = &s
	return stqu
}

// SetQuestionDescription sets the question_description field.
func (stqu *SurveyTemplateQuestionUpdate) SetQuestionDescription(s string) *SurveyTemplateQuestionUpdate {
	stqu.question_description = &s
	return stqu
}

// SetQuestionType sets the question_type field.
func (stqu *SurveyTemplateQuestionUpdate) SetQuestionType(s string) *SurveyTemplateQuestionUpdate {
	stqu.question_type = &s
	return stqu
}

// SetIndex sets the index field.
func (stqu *SurveyTemplateQuestionUpdate) SetIndex(i int) *SurveyTemplateQuestionUpdate {
	stqu.index = &i
	stqu.addindex = nil
	return stqu
}

// AddIndex adds i to index.
func (stqu *SurveyTemplateQuestionUpdate) AddIndex(i int) *SurveyTemplateQuestionUpdate {
	if stqu.addindex == nil {
		stqu.addindex = &i
	} else {
		*stqu.addindex += i
	}
	return stqu
}

// SetCategoryID sets the category edge to SurveyTemplateCategory by id.
func (stqu *SurveyTemplateQuestionUpdate) SetCategoryID(id string) *SurveyTemplateQuestionUpdate {
	if stqu.category == nil {
		stqu.category = make(map[string]struct{})
	}
	stqu.category[id] = struct{}{}
	return stqu
}

// SetNillableCategoryID sets the category edge to SurveyTemplateCategory by id if the given value is not nil.
func (stqu *SurveyTemplateQuestionUpdate) SetNillableCategoryID(id *string) *SurveyTemplateQuestionUpdate {
	if id != nil {
		stqu = stqu.SetCategoryID(*id)
	}
	return stqu
}

// SetCategory sets the category edge to SurveyTemplateCategory.
func (stqu *SurveyTemplateQuestionUpdate) SetCategory(s *SurveyTemplateCategory) *SurveyTemplateQuestionUpdate {
	return stqu.SetCategoryID(s.ID)
}

// ClearCategory clears the category edge to SurveyTemplateCategory.
func (stqu *SurveyTemplateQuestionUpdate) ClearCategory() *SurveyTemplateQuestionUpdate {
	stqu.clearedCategory = true
	return stqu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (stqu *SurveyTemplateQuestionUpdate) Save(ctx context.Context) (int, error) {
	if stqu.update_time == nil {
		v := surveytemplatequestion.UpdateDefaultUpdateTime()
		stqu.update_time = &v
	}
	if len(stqu.category) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"category\"")
	}
	return stqu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (stqu *SurveyTemplateQuestionUpdate) SaveX(ctx context.Context) int {
	affected, err := stqu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (stqu *SurveyTemplateQuestionUpdate) Exec(ctx context.Context) error {
	_, err := stqu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (stqu *SurveyTemplateQuestionUpdate) ExecX(ctx context.Context) {
	if err := stqu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (stqu *SurveyTemplateQuestionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(stqu.driver.Dialect())
		selector = builder.Select(surveytemplatequestion.FieldID).From(builder.Table(surveytemplatequestion.Table))
	)
	for _, p := range stqu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = stqu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := stqu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(surveytemplatequestion.Table).Where(sql.InInts(surveytemplatequestion.FieldID, ids...))
	)
	if value := stqu.update_time; value != nil {
		updater.Set(surveytemplatequestion.FieldUpdateTime, *value)
	}
	if value := stqu.question_title; value != nil {
		updater.Set(surveytemplatequestion.FieldQuestionTitle, *value)
	}
	if value := stqu.question_description; value != nil {
		updater.Set(surveytemplatequestion.FieldQuestionDescription, *value)
	}
	if value := stqu.question_type; value != nil {
		updater.Set(surveytemplatequestion.FieldQuestionType, *value)
	}
	if value := stqu.index; value != nil {
		updater.Set(surveytemplatequestion.FieldIndex, *value)
	}
	if value := stqu.addindex; value != nil {
		updater.Add(surveytemplatequestion.FieldIndex, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if stqu.clearedCategory {
		query, args := builder.Update(surveytemplatequestion.CategoryTable).
			SetNull(surveytemplatequestion.CategoryColumn).
			Where(sql.InInts(surveytemplatecategory.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(stqu.category) > 0 {
		for eid := range stqu.category {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(surveytemplatequestion.CategoryTable).
				Set(surveytemplatequestion.CategoryColumn, eid).
				Where(sql.InInts(surveytemplatequestion.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// SurveyTemplateQuestionUpdateOne is the builder for updating a single SurveyTemplateQuestion entity.
type SurveyTemplateQuestionUpdateOne struct {
	config
	id string

	update_time          *time.Time
	question_title       *string
	question_description *string
	question_type        *string
	index                *int
	addindex             *int
	category             map[string]struct{}
	clearedCategory      bool
}

// SetQuestionTitle sets the question_title field.
func (stquo *SurveyTemplateQuestionUpdateOne) SetQuestionTitle(s string) *SurveyTemplateQuestionUpdateOne {
	stquo.question_title = &s
	return stquo
}

// SetQuestionDescription sets the question_description field.
func (stquo *SurveyTemplateQuestionUpdateOne) SetQuestionDescription(s string) *SurveyTemplateQuestionUpdateOne {
	stquo.question_description = &s
	return stquo
}

// SetQuestionType sets the question_type field.
func (stquo *SurveyTemplateQuestionUpdateOne) SetQuestionType(s string) *SurveyTemplateQuestionUpdateOne {
	stquo.question_type = &s
	return stquo
}

// SetIndex sets the index field.
func (stquo *SurveyTemplateQuestionUpdateOne) SetIndex(i int) *SurveyTemplateQuestionUpdateOne {
	stquo.index = &i
	stquo.addindex = nil
	return stquo
}

// AddIndex adds i to index.
func (stquo *SurveyTemplateQuestionUpdateOne) AddIndex(i int) *SurveyTemplateQuestionUpdateOne {
	if stquo.addindex == nil {
		stquo.addindex = &i
	} else {
		*stquo.addindex += i
	}
	return stquo
}

// SetCategoryID sets the category edge to SurveyTemplateCategory by id.
func (stquo *SurveyTemplateQuestionUpdateOne) SetCategoryID(id string) *SurveyTemplateQuestionUpdateOne {
	if stquo.category == nil {
		stquo.category = make(map[string]struct{})
	}
	stquo.category[id] = struct{}{}
	return stquo
}

// SetNillableCategoryID sets the category edge to SurveyTemplateCategory by id if the given value is not nil.
func (stquo *SurveyTemplateQuestionUpdateOne) SetNillableCategoryID(id *string) *SurveyTemplateQuestionUpdateOne {
	if id != nil {
		stquo = stquo.SetCategoryID(*id)
	}
	return stquo
}

// SetCategory sets the category edge to SurveyTemplateCategory.
func (stquo *SurveyTemplateQuestionUpdateOne) SetCategory(s *SurveyTemplateCategory) *SurveyTemplateQuestionUpdateOne {
	return stquo.SetCategoryID(s.ID)
}

// ClearCategory clears the category edge to SurveyTemplateCategory.
func (stquo *SurveyTemplateQuestionUpdateOne) ClearCategory() *SurveyTemplateQuestionUpdateOne {
	stquo.clearedCategory = true
	return stquo
}

// Save executes the query and returns the updated entity.
func (stquo *SurveyTemplateQuestionUpdateOne) Save(ctx context.Context) (*SurveyTemplateQuestion, error) {
	if stquo.update_time == nil {
		v := surveytemplatequestion.UpdateDefaultUpdateTime()
		stquo.update_time = &v
	}
	if len(stquo.category) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"category\"")
	}
	return stquo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (stquo *SurveyTemplateQuestionUpdateOne) SaveX(ctx context.Context) *SurveyTemplateQuestion {
	stq, err := stquo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return stq
}

// Exec executes the query on the entity.
func (stquo *SurveyTemplateQuestionUpdateOne) Exec(ctx context.Context) error {
	_, err := stquo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (stquo *SurveyTemplateQuestionUpdateOne) ExecX(ctx context.Context) {
	if err := stquo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (stquo *SurveyTemplateQuestionUpdateOne) sqlSave(ctx context.Context) (stq *SurveyTemplateQuestion, err error) {
	var (
		builder  = sql.Dialect(stquo.driver.Dialect())
		selector = builder.Select(surveytemplatequestion.Columns...).From(builder.Table(surveytemplatequestion.Table))
	)
	surveytemplatequestion.ID(stquo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = stquo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		stq = &SurveyTemplateQuestion{config: stquo.config}
		if err := stq.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into SurveyTemplateQuestion: %v", err)
		}
		id = stq.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("SurveyTemplateQuestion with id: %v", stquo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one SurveyTemplateQuestion with the same id: %v", stquo.id)
	}

	tx, err := stquo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(surveytemplatequestion.Table).Where(sql.InInts(surveytemplatequestion.FieldID, ids...))
	)
	if value := stquo.update_time; value != nil {
		updater.Set(surveytemplatequestion.FieldUpdateTime, *value)
		stq.UpdateTime = *value
	}
	if value := stquo.question_title; value != nil {
		updater.Set(surveytemplatequestion.FieldQuestionTitle, *value)
		stq.QuestionTitle = *value
	}
	if value := stquo.question_description; value != nil {
		updater.Set(surveytemplatequestion.FieldQuestionDescription, *value)
		stq.QuestionDescription = *value
	}
	if value := stquo.question_type; value != nil {
		updater.Set(surveytemplatequestion.FieldQuestionType, *value)
		stq.QuestionType = *value
	}
	if value := stquo.index; value != nil {
		updater.Set(surveytemplatequestion.FieldIndex, *value)
		stq.Index = *value
	}
	if value := stquo.addindex; value != nil {
		updater.Add(surveytemplatequestion.FieldIndex, *value)
		stq.Index += *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if stquo.clearedCategory {
		query, args := builder.Update(surveytemplatequestion.CategoryTable).
			SetNull(surveytemplatequestion.CategoryColumn).
			Where(sql.InInts(surveytemplatecategory.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(stquo.category) > 0 {
		for eid := range stquo.category {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(surveytemplatequestion.CategoryTable).
				Set(surveytemplatequestion.CategoryColumn, eid).
				Where(sql.InInts(surveytemplatequestion.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return stq, nil
}
