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
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
)

// SurveyQuestionDelete is the builder for deleting a SurveyQuestion entity.
type SurveyQuestionDelete struct {
	config
	hooks      []Hook
	mutation   *SurveyQuestionMutation
	predicates []predicate.SurveyQuestion
}

// Where adds a new predicate to the delete builder.
func (sqd *SurveyQuestionDelete) Where(ps ...predicate.SurveyQuestion) *SurveyQuestionDelete {
	sqd.predicates = append(sqd.predicates, ps...)
	return sqd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (sqd *SurveyQuestionDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(sqd.hooks) == 0 {
		affected, err = sqd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyQuestionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			sqd.mutation = mutation
			affected, err = sqd.sqlExec(ctx)
			return affected, err
		})
		for i := len(sqd.hooks); i > 0; i-- {
			mut = sqd.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, sqd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (sqd *SurveyQuestionDelete) ExecX(ctx context.Context) int {
	n, err := sqd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (sqd *SurveyQuestionDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: surveyquestion.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: surveyquestion.FieldID,
			},
		},
	}
	if ps := sqd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, sqd.driver, _spec)
}

// SurveyQuestionDeleteOne is the builder for deleting a single SurveyQuestion entity.
type SurveyQuestionDeleteOne struct {
	sqd *SurveyQuestionDelete
}

// Exec executes the deletion query.
func (sqdo *SurveyQuestionDeleteOne) Exec(ctx context.Context) error {
	n, err := sqdo.sqd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{surveyquestion.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (sqdo *SurveyQuestionDeleteOne) ExecX(ctx context.Context) {
	sqdo.sqd.ExecX(ctx)
}
