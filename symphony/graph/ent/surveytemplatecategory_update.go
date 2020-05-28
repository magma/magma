// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatecategory"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatequestion"
)

// SurveyTemplateCategoryUpdate is the builder for updating SurveyTemplateCategory entities.
type SurveyTemplateCategoryUpdate struct {
	config
	hooks      []Hook
	mutation   *SurveyTemplateCategoryMutation
	predicates []predicate.SurveyTemplateCategory
}

// Where adds a new predicate for the builder.
func (stcu *SurveyTemplateCategoryUpdate) Where(ps ...predicate.SurveyTemplateCategory) *SurveyTemplateCategoryUpdate {
	stcu.predicates = append(stcu.predicates, ps...)
	return stcu
}

// SetCategoryTitle sets the category_title field.
func (stcu *SurveyTemplateCategoryUpdate) SetCategoryTitle(s string) *SurveyTemplateCategoryUpdate {
	stcu.mutation.SetCategoryTitle(s)
	return stcu
}

// SetCategoryDescription sets the category_description field.
func (stcu *SurveyTemplateCategoryUpdate) SetCategoryDescription(s string) *SurveyTemplateCategoryUpdate {
	stcu.mutation.SetCategoryDescription(s)
	return stcu
}

// AddSurveyTemplateQuestionIDs adds the survey_template_questions edge to SurveyTemplateQuestion by ids.
func (stcu *SurveyTemplateCategoryUpdate) AddSurveyTemplateQuestionIDs(ids ...int) *SurveyTemplateCategoryUpdate {
	stcu.mutation.AddSurveyTemplateQuestionIDs(ids...)
	return stcu
}

// AddSurveyTemplateQuestions adds the survey_template_questions edges to SurveyTemplateQuestion.
func (stcu *SurveyTemplateCategoryUpdate) AddSurveyTemplateQuestions(s ...*SurveyTemplateQuestion) *SurveyTemplateCategoryUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stcu.AddSurveyTemplateQuestionIDs(ids...)
}

// SetLocationTypeID sets the location_type edge to LocationType by id.
func (stcu *SurveyTemplateCategoryUpdate) SetLocationTypeID(id int) *SurveyTemplateCategoryUpdate {
	stcu.mutation.SetLocationTypeID(id)
	return stcu
}

// SetNillableLocationTypeID sets the location_type edge to LocationType by id if the given value is not nil.
func (stcu *SurveyTemplateCategoryUpdate) SetNillableLocationTypeID(id *int) *SurveyTemplateCategoryUpdate {
	if id != nil {
		stcu = stcu.SetLocationTypeID(*id)
	}
	return stcu
}

// SetLocationType sets the location_type edge to LocationType.
func (stcu *SurveyTemplateCategoryUpdate) SetLocationType(l *LocationType) *SurveyTemplateCategoryUpdate {
	return stcu.SetLocationTypeID(l.ID)
}

// RemoveSurveyTemplateQuestionIDs removes the survey_template_questions edge to SurveyTemplateQuestion by ids.
func (stcu *SurveyTemplateCategoryUpdate) RemoveSurveyTemplateQuestionIDs(ids ...int) *SurveyTemplateCategoryUpdate {
	stcu.mutation.RemoveSurveyTemplateQuestionIDs(ids...)
	return stcu
}

// RemoveSurveyTemplateQuestions removes survey_template_questions edges to SurveyTemplateQuestion.
func (stcu *SurveyTemplateCategoryUpdate) RemoveSurveyTemplateQuestions(s ...*SurveyTemplateQuestion) *SurveyTemplateCategoryUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stcu.RemoveSurveyTemplateQuestionIDs(ids...)
}

// ClearLocationType clears the location_type edge to LocationType.
func (stcu *SurveyTemplateCategoryUpdate) ClearLocationType() *SurveyTemplateCategoryUpdate {
	stcu.mutation.ClearLocationType()
	return stcu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (stcu *SurveyTemplateCategoryUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := stcu.mutation.UpdateTime(); !ok {
		v := surveytemplatecategory.UpdateDefaultUpdateTime()
		stcu.mutation.SetUpdateTime(v)
	}

	var (
		err      error
		affected int
	)
	if len(stcu.hooks) == 0 {
		affected, err = stcu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyTemplateCategoryMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			stcu.mutation = mutation
			affected, err = stcu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(stcu.hooks) - 1; i >= 0; i-- {
			mut = stcu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, stcu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveytemplatecategory.Table,
			Columns: surveytemplatecategory.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: surveytemplatecategory.FieldID,
			},
		},
	}
	if ps := stcu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := stcu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveytemplatecategory.FieldUpdateTime,
		})
	}
	if value, ok := stcu.mutation.CategoryTitle(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveytemplatecategory.FieldCategoryTitle,
		})
	}
	if value, ok := stcu.mutation.CategoryDescription(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveytemplatecategory.FieldCategoryDescription,
		})
	}
	if nodes := stcu.mutation.RemovedSurveyTemplateQuestionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   surveytemplatecategory.SurveyTemplateQuestionsTable,
			Columns: []string{surveytemplatecategory.SurveyTemplateQuestionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveytemplatequestion.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := stcu.mutation.SurveyTemplateQuestionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   surveytemplatecategory.SurveyTemplateQuestionsTable,
			Columns: []string{surveytemplatecategory.SurveyTemplateQuestionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveytemplatequestion.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if stcu.mutation.LocationTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   surveytemplatecategory.LocationTypeTable,
			Columns: []string{surveytemplatecategory.LocationTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: locationtype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := stcu.mutation.LocationTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   surveytemplatecategory.LocationTypeTable,
			Columns: []string{surveytemplatecategory.LocationTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: locationtype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, stcu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{surveytemplatecategory.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// SurveyTemplateCategoryUpdateOne is the builder for updating a single SurveyTemplateCategory entity.
type SurveyTemplateCategoryUpdateOne struct {
	config
	hooks    []Hook
	mutation *SurveyTemplateCategoryMutation
}

// SetCategoryTitle sets the category_title field.
func (stcuo *SurveyTemplateCategoryUpdateOne) SetCategoryTitle(s string) *SurveyTemplateCategoryUpdateOne {
	stcuo.mutation.SetCategoryTitle(s)
	return stcuo
}

// SetCategoryDescription sets the category_description field.
func (stcuo *SurveyTemplateCategoryUpdateOne) SetCategoryDescription(s string) *SurveyTemplateCategoryUpdateOne {
	stcuo.mutation.SetCategoryDescription(s)
	return stcuo
}

// AddSurveyTemplateQuestionIDs adds the survey_template_questions edge to SurveyTemplateQuestion by ids.
func (stcuo *SurveyTemplateCategoryUpdateOne) AddSurveyTemplateQuestionIDs(ids ...int) *SurveyTemplateCategoryUpdateOne {
	stcuo.mutation.AddSurveyTemplateQuestionIDs(ids...)
	return stcuo
}

// AddSurveyTemplateQuestions adds the survey_template_questions edges to SurveyTemplateQuestion.
func (stcuo *SurveyTemplateCategoryUpdateOne) AddSurveyTemplateQuestions(s ...*SurveyTemplateQuestion) *SurveyTemplateCategoryUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stcuo.AddSurveyTemplateQuestionIDs(ids...)
}

// SetLocationTypeID sets the location_type edge to LocationType by id.
func (stcuo *SurveyTemplateCategoryUpdateOne) SetLocationTypeID(id int) *SurveyTemplateCategoryUpdateOne {
	stcuo.mutation.SetLocationTypeID(id)
	return stcuo
}

// SetNillableLocationTypeID sets the location_type edge to LocationType by id if the given value is not nil.
func (stcuo *SurveyTemplateCategoryUpdateOne) SetNillableLocationTypeID(id *int) *SurveyTemplateCategoryUpdateOne {
	if id != nil {
		stcuo = stcuo.SetLocationTypeID(*id)
	}
	return stcuo
}

// SetLocationType sets the location_type edge to LocationType.
func (stcuo *SurveyTemplateCategoryUpdateOne) SetLocationType(l *LocationType) *SurveyTemplateCategoryUpdateOne {
	return stcuo.SetLocationTypeID(l.ID)
}

// RemoveSurveyTemplateQuestionIDs removes the survey_template_questions edge to SurveyTemplateQuestion by ids.
func (stcuo *SurveyTemplateCategoryUpdateOne) RemoveSurveyTemplateQuestionIDs(ids ...int) *SurveyTemplateCategoryUpdateOne {
	stcuo.mutation.RemoveSurveyTemplateQuestionIDs(ids...)
	return stcuo
}

// RemoveSurveyTemplateQuestions removes survey_template_questions edges to SurveyTemplateQuestion.
func (stcuo *SurveyTemplateCategoryUpdateOne) RemoveSurveyTemplateQuestions(s ...*SurveyTemplateQuestion) *SurveyTemplateCategoryUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stcuo.RemoveSurveyTemplateQuestionIDs(ids...)
}

// ClearLocationType clears the location_type edge to LocationType.
func (stcuo *SurveyTemplateCategoryUpdateOne) ClearLocationType() *SurveyTemplateCategoryUpdateOne {
	stcuo.mutation.ClearLocationType()
	return stcuo
}

// Save executes the query and returns the updated entity.
func (stcuo *SurveyTemplateCategoryUpdateOne) Save(ctx context.Context) (*SurveyTemplateCategory, error) {
	if _, ok := stcuo.mutation.UpdateTime(); !ok {
		v := surveytemplatecategory.UpdateDefaultUpdateTime()
		stcuo.mutation.SetUpdateTime(v)
	}

	var (
		err  error
		node *SurveyTemplateCategory
	)
	if len(stcuo.hooks) == 0 {
		node, err = stcuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyTemplateCategoryMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			stcuo.mutation = mutation
			node, err = stcuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(stcuo.hooks) - 1; i >= 0; i-- {
			mut = stcuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, stcuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveytemplatecategory.Table,
			Columns: surveytemplatecategory.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: surveytemplatecategory.FieldID,
			},
		},
	}
	id, ok := stcuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing SurveyTemplateCategory.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := stcuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveytemplatecategory.FieldUpdateTime,
		})
	}
	if value, ok := stcuo.mutation.CategoryTitle(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveytemplatecategory.FieldCategoryTitle,
		})
	}
	if value, ok := stcuo.mutation.CategoryDescription(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveytemplatecategory.FieldCategoryDescription,
		})
	}
	if nodes := stcuo.mutation.RemovedSurveyTemplateQuestionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   surveytemplatecategory.SurveyTemplateQuestionsTable,
			Columns: []string{surveytemplatecategory.SurveyTemplateQuestionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveytemplatequestion.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := stcuo.mutation.SurveyTemplateQuestionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   surveytemplatecategory.SurveyTemplateQuestionsTable,
			Columns: []string{surveytemplatecategory.SurveyTemplateQuestionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveytemplatequestion.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if stcuo.mutation.LocationTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   surveytemplatecategory.LocationTypeTable,
			Columns: []string{surveytemplatecategory.LocationTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: locationtype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := stcuo.mutation.LocationTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   surveytemplatecategory.LocationTypeTable,
			Columns: []string{surveytemplatecategory.LocationTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: locationtype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	stc = &SurveyTemplateCategory{config: stcuo.config}
	_spec.Assign = stc.assignValues
	_spec.ScanValues = stc.scanValues()
	if err = sqlgraph.UpdateNode(ctx, stcuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{surveytemplatecategory.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return stc, nil
}
