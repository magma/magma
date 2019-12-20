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

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
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
		stc  = &SurveyTemplateCategory{config: stcc.config}
		spec = &sqlgraph.CreateSpec{
			Table: surveytemplatecategory.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: surveytemplatecategory.FieldID,
			},
		}
	)
	if value := stcc.create_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveytemplatecategory.FieldCreateTime,
		})
		stc.CreateTime = *value
	}
	if value := stcc.update_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveytemplatecategory.FieldUpdateTime,
		})
		stc.UpdateTime = *value
	}
	if value := stcc.category_title; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveytemplatecategory.FieldCategoryTitle,
		})
		stc.CategoryTitle = *value
	}
	if value := stcc.category_description; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveytemplatecategory.FieldCategoryDescription,
		})
		stc.CategoryDescription = *value
	}
	if nodes := stcc.survey_template_questions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   surveytemplatecategory.SurveyTemplateQuestionsTable,
			Columns: []string{surveytemplatecategory.SurveyTemplateQuestionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveytemplatequestion.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges = append(spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, stcc.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := spec.ID.Value.(int64)
	stc.ID = strconv.FormatInt(id, 10)
	return stc, nil
}
