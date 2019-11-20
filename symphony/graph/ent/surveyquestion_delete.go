// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
)

// SurveyQuestionDelete is the builder for deleting a SurveyQuestion entity.
type SurveyQuestionDelete struct {
	config
	predicates []predicate.SurveyQuestion
}

// Where adds a new predicate to the delete builder.
func (sqd *SurveyQuestionDelete) Where(ps ...predicate.SurveyQuestion) *SurveyQuestionDelete {
	sqd.predicates = append(sqd.predicates, ps...)
	return sqd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (sqd *SurveyQuestionDelete) Exec(ctx context.Context) (int, error) {
	return sqd.sqlExec(ctx)
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
	var (
		res     sql.Result
		builder = sql.Dialect(sqd.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(surveyquestion.Table))
	for _, p := range sqd.predicates {
		p(selector)
	}
	query, args := builder.Delete(surveyquestion.Table).FromSelect(selector).Query()
	if err := sqd.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
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
		return &ErrNotFound{surveyquestion.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (sqdo *SurveyQuestionDeleteOne) ExecX(ctx context.Context) {
	sqdo.sqd.ExecX(ctx)
}
