// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
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
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveytemplatecategory.Table,
			Columns: surveytemplatecategory.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: surveytemplatecategory.FieldID,
			},
		},
	}
	if ps := stcu.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := stcu.update_time; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveytemplatecategory.FieldUpdateTime,
		})
	}
	if value := stcu.category_title; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveytemplatecategory.FieldCategoryTitle,
		})
	}
	if value := stcu.category_description; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveytemplatecategory.FieldCategoryDescription,
		})
	}
	if nodes := stcu.removedSurveyTemplateQuestions; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := stcu.survey_template_questions; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, stcu.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
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
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveytemplatecategory.Table,
			Columns: surveytemplatecategory.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  stcuo.id,
				Type:   field.TypeString,
				Column: surveytemplatecategory.FieldID,
			},
		},
	}
	if value := stcuo.update_time; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveytemplatecategory.FieldUpdateTime,
		})
	}
	if value := stcuo.category_title; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveytemplatecategory.FieldCategoryTitle,
		})
	}
	if value := stcuo.category_description; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveytemplatecategory.FieldCategoryDescription,
		})
	}
	if nodes := stcuo.removedSurveyTemplateQuestions; len(nodes) > 0 {
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
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := stcuo.survey_template_questions; len(nodes) > 0 {
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
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	stc = &SurveyTemplateCategory{config: stcuo.config}
	spec.Assign = stc.assignValues
	spec.ScanValues = stc.scanValues()
	if err = sqlgraph.UpdateNode(ctx, stcuo.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return stc, nil
}
