// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatequestion"
)

// SurveyTemplateQuestionCreate is the builder for creating a SurveyTemplateQuestion entity.
type SurveyTemplateQuestionCreate struct {
	config
	create_time          *time.Time
	update_time          *time.Time
	question_title       *string
	question_description *string
	question_type        *string
	index                *int
	category             map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (stqc *SurveyTemplateQuestionCreate) SetCreateTime(t time.Time) *SurveyTemplateQuestionCreate {
	stqc.create_time = &t
	return stqc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (stqc *SurveyTemplateQuestionCreate) SetNillableCreateTime(t *time.Time) *SurveyTemplateQuestionCreate {
	if t != nil {
		stqc.SetCreateTime(*t)
	}
	return stqc
}

// SetUpdateTime sets the update_time field.
func (stqc *SurveyTemplateQuestionCreate) SetUpdateTime(t time.Time) *SurveyTemplateQuestionCreate {
	stqc.update_time = &t
	return stqc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (stqc *SurveyTemplateQuestionCreate) SetNillableUpdateTime(t *time.Time) *SurveyTemplateQuestionCreate {
	if t != nil {
		stqc.SetUpdateTime(*t)
	}
	return stqc
}

// SetQuestionTitle sets the question_title field.
func (stqc *SurveyTemplateQuestionCreate) SetQuestionTitle(s string) *SurveyTemplateQuestionCreate {
	stqc.question_title = &s
	return stqc
}

// SetQuestionDescription sets the question_description field.
func (stqc *SurveyTemplateQuestionCreate) SetQuestionDescription(s string) *SurveyTemplateQuestionCreate {
	stqc.question_description = &s
	return stqc
}

// SetQuestionType sets the question_type field.
func (stqc *SurveyTemplateQuestionCreate) SetQuestionType(s string) *SurveyTemplateQuestionCreate {
	stqc.question_type = &s
	return stqc
}

// SetIndex sets the index field.
func (stqc *SurveyTemplateQuestionCreate) SetIndex(i int) *SurveyTemplateQuestionCreate {
	stqc.index = &i
	return stqc
}

// SetCategoryID sets the category edge to SurveyTemplateCategory by id.
func (stqc *SurveyTemplateQuestionCreate) SetCategoryID(id string) *SurveyTemplateQuestionCreate {
	if stqc.category == nil {
		stqc.category = make(map[string]struct{})
	}
	stqc.category[id] = struct{}{}
	return stqc
}

// SetNillableCategoryID sets the category edge to SurveyTemplateCategory by id if the given value is not nil.
func (stqc *SurveyTemplateQuestionCreate) SetNillableCategoryID(id *string) *SurveyTemplateQuestionCreate {
	if id != nil {
		stqc = stqc.SetCategoryID(*id)
	}
	return stqc
}

// SetCategory sets the category edge to SurveyTemplateCategory.
func (stqc *SurveyTemplateQuestionCreate) SetCategory(s *SurveyTemplateCategory) *SurveyTemplateQuestionCreate {
	return stqc.SetCategoryID(s.ID)
}

// Save creates the SurveyTemplateQuestion in the database.
func (stqc *SurveyTemplateQuestionCreate) Save(ctx context.Context) (*SurveyTemplateQuestion, error) {
	if stqc.create_time == nil {
		v := surveytemplatequestion.DefaultCreateTime()
		stqc.create_time = &v
	}
	if stqc.update_time == nil {
		v := surveytemplatequestion.DefaultUpdateTime()
		stqc.update_time = &v
	}
	if stqc.question_title == nil {
		return nil, errors.New("ent: missing required field \"question_title\"")
	}
	if stqc.question_description == nil {
		return nil, errors.New("ent: missing required field \"question_description\"")
	}
	if stqc.question_type == nil {
		return nil, errors.New("ent: missing required field \"question_type\"")
	}
	if stqc.index == nil {
		return nil, errors.New("ent: missing required field \"index\"")
	}
	if len(stqc.category) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"category\"")
	}
	return stqc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (stqc *SurveyTemplateQuestionCreate) SaveX(ctx context.Context) *SurveyTemplateQuestion {
	v, err := stqc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (stqc *SurveyTemplateQuestionCreate) sqlSave(ctx context.Context) (*SurveyTemplateQuestion, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(stqc.driver.Dialect())
		stq     = &SurveyTemplateQuestion{config: stqc.config}
	)
	tx, err := stqc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(surveytemplatequestion.Table).Default()
	if value := stqc.create_time; value != nil {
		insert.Set(surveytemplatequestion.FieldCreateTime, *value)
		stq.CreateTime = *value
	}
	if value := stqc.update_time; value != nil {
		insert.Set(surveytemplatequestion.FieldUpdateTime, *value)
		stq.UpdateTime = *value
	}
	if value := stqc.question_title; value != nil {
		insert.Set(surveytemplatequestion.FieldQuestionTitle, *value)
		stq.QuestionTitle = *value
	}
	if value := stqc.question_description; value != nil {
		insert.Set(surveytemplatequestion.FieldQuestionDescription, *value)
		stq.QuestionDescription = *value
	}
	if value := stqc.question_type; value != nil {
		insert.Set(surveytemplatequestion.FieldQuestionType, *value)
		stq.QuestionType = *value
	}
	if value := stqc.index; value != nil {
		insert.Set(surveytemplatequestion.FieldIndex, *value)
		stq.Index = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(surveytemplatequestion.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	stq.ID = strconv.FormatInt(id, 10)
	if len(stqc.category) > 0 {
		for eid := range stqc.category {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(surveytemplatequestion.CategoryTable).
				Set(surveytemplatequestion.CategoryColumn, eid).
				Where(sql.EQ(surveytemplatequestion.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return stq, nil
}
