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
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatecategory"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatequestion"
)

// SurveyTemplateQuestionUpdate is the builder for updating SurveyTemplateQuestion entities.
type SurveyTemplateQuestionUpdate struct {
	config
	hooks      []Hook
	mutation   *SurveyTemplateQuestionMutation
	predicates []predicate.SurveyTemplateQuestion
}

// Where adds a new predicate for the builder.
func (stqu *SurveyTemplateQuestionUpdate) Where(ps ...predicate.SurveyTemplateQuestion) *SurveyTemplateQuestionUpdate {
	stqu.predicates = append(stqu.predicates, ps...)
	return stqu
}

// SetQuestionTitle sets the question_title field.
func (stqu *SurveyTemplateQuestionUpdate) SetQuestionTitle(s string) *SurveyTemplateQuestionUpdate {
	stqu.mutation.SetQuestionTitle(s)
	return stqu
}

// SetQuestionDescription sets the question_description field.
func (stqu *SurveyTemplateQuestionUpdate) SetQuestionDescription(s string) *SurveyTemplateQuestionUpdate {
	stqu.mutation.SetQuestionDescription(s)
	return stqu
}

// SetQuestionType sets the question_type field.
func (stqu *SurveyTemplateQuestionUpdate) SetQuestionType(s string) *SurveyTemplateQuestionUpdate {
	stqu.mutation.SetQuestionType(s)
	return stqu
}

// SetIndex sets the index field.
func (stqu *SurveyTemplateQuestionUpdate) SetIndex(i int) *SurveyTemplateQuestionUpdate {
	stqu.mutation.ResetIndex()
	stqu.mutation.SetIndex(i)
	return stqu
}

// AddIndex adds i to index.
func (stqu *SurveyTemplateQuestionUpdate) AddIndex(i int) *SurveyTemplateQuestionUpdate {
	stqu.mutation.AddIndex(i)
	return stqu
}

// SetCategoryID sets the category edge to SurveyTemplateCategory by id.
func (stqu *SurveyTemplateQuestionUpdate) SetCategoryID(id int) *SurveyTemplateQuestionUpdate {
	stqu.mutation.SetCategoryID(id)
	return stqu
}

// SetNillableCategoryID sets the category edge to SurveyTemplateCategory by id if the given value is not nil.
func (stqu *SurveyTemplateQuestionUpdate) SetNillableCategoryID(id *int) *SurveyTemplateQuestionUpdate {
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
	stqu.mutation.ClearCategory()
	return stqu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (stqu *SurveyTemplateQuestionUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := stqu.mutation.UpdateTime(); !ok {
		v := surveytemplatequestion.UpdateDefaultUpdateTime()
		stqu.mutation.SetUpdateTime(v)
	}

	var (
		err      error
		affected int
	)
	if len(stqu.hooks) == 0 {
		affected, err = stqu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyTemplateQuestionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			stqu.mutation = mutation
			affected, err = stqu.sqlSave(ctx)
			return affected, err
		})
		for i := len(stqu.hooks); i > 0; i-- {
			mut = stqu.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, stqu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveytemplatequestion.Table,
			Columns: surveytemplatequestion.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: surveytemplatequestion.FieldID,
			},
		},
	}
	if ps := stqu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := stqu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveytemplatequestion.FieldUpdateTime,
		})
	}
	if value, ok := stqu.mutation.QuestionTitle(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveytemplatequestion.FieldQuestionTitle,
		})
	}
	if value, ok := stqu.mutation.QuestionDescription(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveytemplatequestion.FieldQuestionDescription,
		})
	}
	if value, ok := stqu.mutation.QuestionType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveytemplatequestion.FieldQuestionType,
		})
	}
	if value, ok := stqu.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveytemplatequestion.FieldIndex,
		})
	}
	if value, ok := stqu.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveytemplatequestion.FieldIndex,
		})
	}
	if stqu.mutation.CategoryCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   surveytemplatequestion.CategoryTable,
			Columns: []string{surveytemplatequestion.CategoryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveytemplatecategory.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := stqu.mutation.CategoryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   surveytemplatequestion.CategoryTable,
			Columns: []string{surveytemplatequestion.CategoryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveytemplatecategory.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, stqu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{surveytemplatequestion.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// SurveyTemplateQuestionUpdateOne is the builder for updating a single SurveyTemplateQuestion entity.
type SurveyTemplateQuestionUpdateOne struct {
	config
	hooks    []Hook
	mutation *SurveyTemplateQuestionMutation
}

// SetQuestionTitle sets the question_title field.
func (stquo *SurveyTemplateQuestionUpdateOne) SetQuestionTitle(s string) *SurveyTemplateQuestionUpdateOne {
	stquo.mutation.SetQuestionTitle(s)
	return stquo
}

// SetQuestionDescription sets the question_description field.
func (stquo *SurveyTemplateQuestionUpdateOne) SetQuestionDescription(s string) *SurveyTemplateQuestionUpdateOne {
	stquo.mutation.SetQuestionDescription(s)
	return stquo
}

// SetQuestionType sets the question_type field.
func (stquo *SurveyTemplateQuestionUpdateOne) SetQuestionType(s string) *SurveyTemplateQuestionUpdateOne {
	stquo.mutation.SetQuestionType(s)
	return stquo
}

// SetIndex sets the index field.
func (stquo *SurveyTemplateQuestionUpdateOne) SetIndex(i int) *SurveyTemplateQuestionUpdateOne {
	stquo.mutation.ResetIndex()
	stquo.mutation.SetIndex(i)
	return stquo
}

// AddIndex adds i to index.
func (stquo *SurveyTemplateQuestionUpdateOne) AddIndex(i int) *SurveyTemplateQuestionUpdateOne {
	stquo.mutation.AddIndex(i)
	return stquo
}

// SetCategoryID sets the category edge to SurveyTemplateCategory by id.
func (stquo *SurveyTemplateQuestionUpdateOne) SetCategoryID(id int) *SurveyTemplateQuestionUpdateOne {
	stquo.mutation.SetCategoryID(id)
	return stquo
}

// SetNillableCategoryID sets the category edge to SurveyTemplateCategory by id if the given value is not nil.
func (stquo *SurveyTemplateQuestionUpdateOne) SetNillableCategoryID(id *int) *SurveyTemplateQuestionUpdateOne {
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
	stquo.mutation.ClearCategory()
	return stquo
}

// Save executes the query and returns the updated entity.
func (stquo *SurveyTemplateQuestionUpdateOne) Save(ctx context.Context) (*SurveyTemplateQuestion, error) {
	if _, ok := stquo.mutation.UpdateTime(); !ok {
		v := surveytemplatequestion.UpdateDefaultUpdateTime()
		stquo.mutation.SetUpdateTime(v)
	}

	var (
		err  error
		node *SurveyTemplateQuestion
	)
	if len(stquo.hooks) == 0 {
		node, err = stquo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyTemplateQuestionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			stquo.mutation = mutation
			node, err = stquo.sqlSave(ctx)
			return node, err
		})
		for i := len(stquo.hooks); i > 0; i-- {
			mut = stquo.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, stquo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveytemplatequestion.Table,
			Columns: surveytemplatequestion.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: surveytemplatequestion.FieldID,
			},
		},
	}
	id, ok := stquo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing SurveyTemplateQuestion.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := stquo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveytemplatequestion.FieldUpdateTime,
		})
	}
	if value, ok := stquo.mutation.QuestionTitle(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveytemplatequestion.FieldQuestionTitle,
		})
	}
	if value, ok := stquo.mutation.QuestionDescription(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveytemplatequestion.FieldQuestionDescription,
		})
	}
	if value, ok := stquo.mutation.QuestionType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveytemplatequestion.FieldQuestionType,
		})
	}
	if value, ok := stquo.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveytemplatequestion.FieldIndex,
		})
	}
	if value, ok := stquo.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveytemplatequestion.FieldIndex,
		})
	}
	if stquo.mutation.CategoryCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   surveytemplatequestion.CategoryTable,
			Columns: []string{surveytemplatequestion.CategoryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveytemplatecategory.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := stquo.mutation.CategoryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   surveytemplatequestion.CategoryTable,
			Columns: []string{surveytemplatequestion.CategoryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveytemplatecategory.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	stq = &SurveyTemplateQuestion{config: stquo.config}
	_spec.Assign = stq.assignValues
	_spec.ScanValues = stq.scanValues()
	if err = sqlgraph.UpdateNode(ctx, stquo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{surveytemplatequestion.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return stq, nil
}
