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
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatecategory"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatequestion"
)

// SurveyTemplateCategoryCreate is the builder for creating a SurveyTemplateCategory entity.
type SurveyTemplateCategoryCreate struct {
	config
	create_time               *time.Time
	update_time               *time.Time
	category_title            *string
	category_description      *string
	survey_template_questions map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (stcc *SurveyTemplateCategoryCreate) SetCreateTime(t time.Time) *SurveyTemplateCategoryCreate {
	stcc.create_time = &t
	return stcc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (stcc *SurveyTemplateCategoryCreate) SetNillableCreateTime(t *time.Time) *SurveyTemplateCategoryCreate {
	if t != nil {
		stcc.SetCreateTime(*t)
	}
	return stcc
}

// SetUpdateTime sets the update_time field.
func (stcc *SurveyTemplateCategoryCreate) SetUpdateTime(t time.Time) *SurveyTemplateCategoryCreate {
	stcc.update_time = &t
	return stcc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (stcc *SurveyTemplateCategoryCreate) SetNillableUpdateTime(t *time.Time) *SurveyTemplateCategoryCreate {
	if t != nil {
		stcc.SetUpdateTime(*t)
	}
	return stcc
}

// SetCategoryTitle sets the category_title field.
func (stcc *SurveyTemplateCategoryCreate) SetCategoryTitle(s string) *SurveyTemplateCategoryCreate {
	stcc.category_title = &s
	return stcc
}

// SetCategoryDescription sets the category_description field.
func (stcc *SurveyTemplateCategoryCreate) SetCategoryDescription(s string) *SurveyTemplateCategoryCreate {
	stcc.category_description = &s
	return stcc
}

// AddSurveyTemplateQuestionIDs adds the survey_template_questions edge to SurveyTemplateQuestion by ids.
func (stcc *SurveyTemplateCategoryCreate) AddSurveyTemplateQuestionIDs(ids ...string) *SurveyTemplateCategoryCreate {
	if stcc.survey_template_questions == nil {
		stcc.survey_template_questions = make(map[string]struct{})
	}
	for i := range ids {
		stcc.survey_template_questions[ids[i]] = struct{}{}
	}
	return stcc
}

// AddSurveyTemplateQuestions adds the survey_template_questions edges to SurveyTemplateQuestion.
func (stcc *SurveyTemplateCategoryCreate) AddSurveyTemplateQuestions(s ...*SurveyTemplateQuestion) *SurveyTemplateCategoryCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stcc.AddSurveyTemplateQuestionIDs(ids...)
}

// Save creates the SurveyTemplateCategory in the database.
func (stcc *SurveyTemplateCategoryCreate) Save(ctx context.Context) (*SurveyTemplateCategory, error) {
	if stcc.create_time == nil {
		v := surveytemplatecategory.DefaultCreateTime()
		stcc.create_time = &v
	}
	if stcc.update_time == nil {
		v := surveytemplatecategory.DefaultUpdateTime()
		stcc.update_time = &v
	}
	if stcc.category_title == nil {
		return nil, errors.New("ent: missing required field \"category_title\"")
	}
	if stcc.category_description == nil {
		return nil, errors.New("ent: missing required field \"category_description\"")
	}
	return stcc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (stcc *SurveyTemplateCategoryCreate) SaveX(ctx context.Context) *SurveyTemplateCategory {
	v, err := stcc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (stcc *SurveyTemplateCategoryCreate) sqlSave(ctx context.Context) (*SurveyTemplateCategory, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(stcc.driver.Dialect())
		stc     = &SurveyTemplateCategory{config: stcc.config}
	)
	tx, err := stcc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(surveytemplatecategory.Table).Default()
	if value := stcc.create_time; value != nil {
		insert.Set(surveytemplatecategory.FieldCreateTime, *value)
		stc.CreateTime = *value
	}
	if value := stcc.update_time; value != nil {
		insert.Set(surveytemplatecategory.FieldUpdateTime, *value)
		stc.UpdateTime = *value
	}
	if value := stcc.category_title; value != nil {
		insert.Set(surveytemplatecategory.FieldCategoryTitle, *value)
		stc.CategoryTitle = *value
	}
	if value := stcc.category_description; value != nil {
		insert.Set(surveytemplatecategory.FieldCategoryDescription, *value)
		stc.CategoryDescription = *value
	}
	id, err := insertLastID(ctx, tx, insert.Returning(surveytemplatecategory.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	stc.ID = strconv.FormatInt(id, 10)
	if len(stcc.survey_template_questions) > 0 {
		p := sql.P()
		for eid := range stcc.survey_template_questions {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
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
		if int(affected) < len(stcc.survey_template_questions) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"survey_template_questions\" %v already connected to a different \"SurveyTemplateCategory\"", keys(stcc.survey_template_questions))})
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return stc, nil
}
