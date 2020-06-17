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
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatequestion"
)

// SurveyTemplateQuestionDelete is the builder for deleting a SurveyTemplateQuestion entity.
type SurveyTemplateQuestionDelete struct {
	config
	hooks      []Hook
	mutation   *SurveyTemplateQuestionMutation
	predicates []predicate.SurveyTemplateQuestion
}

// Where adds a new predicate to the delete builder.
func (stqd *SurveyTemplateQuestionDelete) Where(ps ...predicate.SurveyTemplateQuestion) *SurveyTemplateQuestionDelete {
	stqd.predicates = append(stqd.predicates, ps...)
	return stqd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (stqd *SurveyTemplateQuestionDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(stqd.hooks) == 0 {
		affected, err = stqd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyTemplateQuestionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			stqd.mutation = mutation
			affected, err = stqd.sqlExec(ctx)
			return affected, err
		})
		for i := len(stqd.hooks) - 1; i >= 0; i-- {
			mut = stqd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, stqd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (stqd *SurveyTemplateQuestionDelete) ExecX(ctx context.Context) int {
	n, err := stqd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (stqd *SurveyTemplateQuestionDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: surveytemplatequestion.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: surveytemplatequestion.FieldID,
			},
		},
	}
	if ps := stqd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, stqd.driver, _spec)
}

// SurveyTemplateQuestionDeleteOne is the builder for deleting a single SurveyTemplateQuestion entity.
type SurveyTemplateQuestionDeleteOne struct {
	stqd *SurveyTemplateQuestionDelete
}

// Exec executes the deletion query.
func (stqdo *SurveyTemplateQuestionDeleteOne) Exec(ctx context.Context) error {
	n, err := stqdo.stqd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{surveytemplatequestion.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (stqdo *SurveyTemplateQuestionDeleteOne) ExecX(ctx context.Context) {
	stqdo.stqd.ExecX(ctx)
}
