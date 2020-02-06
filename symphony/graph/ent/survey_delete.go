// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/survey"
)

// SurveyDelete is the builder for deleting a Survey entity.
type SurveyDelete struct {
	config
	predicates []predicate.Survey
}

// Where adds a new predicate to the delete builder.
func (sd *SurveyDelete) Where(ps ...predicate.Survey) *SurveyDelete {
	sd.predicates = append(sd.predicates, ps...)
	return sd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (sd *SurveyDelete) Exec(ctx context.Context) (int, error) {
	return sd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (sd *SurveyDelete) ExecX(ctx context.Context) int {
	n, err := sd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (sd *SurveyDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: survey.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: survey.FieldID,
			},
		},
	}
	if ps := sd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, sd.driver, _spec)
}

// SurveyDeleteOne is the builder for deleting a single Survey entity.
type SurveyDeleteOne struct {
	sd *SurveyDelete
}

// Exec executes the deletion query.
func (sdo *SurveyDeleteOne) Exec(ctx context.Context) error {
	n, err := sdo.sd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{survey.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (sdo *SurveyDeleteOne) ExecX(ctx context.Context) {
	sdo.sd.ExecX(ctx)
}
