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
	"github.com/facebookincubator/symphony/graph/ent/reportfilter"
)

// ReportFilterDelete is the builder for deleting a ReportFilter entity.
type ReportFilterDelete struct {
	config
	hooks      []Hook
	mutation   *ReportFilterMutation
	predicates []predicate.ReportFilter
}

// Where adds a new predicate to the delete builder.
func (rfd *ReportFilterDelete) Where(ps ...predicate.ReportFilter) *ReportFilterDelete {
	rfd.predicates = append(rfd.predicates, ps...)
	return rfd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (rfd *ReportFilterDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(rfd.hooks) == 0 {
		affected, err = rfd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ReportFilterMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			rfd.mutation = mutation
			affected, err = rfd.sqlExec(ctx)
			return affected, err
		})
		for i := len(rfd.hooks); i > 0; i-- {
			mut = rfd.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, rfd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (rfd *ReportFilterDelete) ExecX(ctx context.Context) int {
	n, err := rfd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (rfd *ReportFilterDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: reportfilter.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: reportfilter.FieldID,
			},
		},
	}
	if ps := rfd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, rfd.driver, _spec)
}

// ReportFilterDeleteOne is the builder for deleting a single ReportFilter entity.
type ReportFilterDeleteOne struct {
	rfd *ReportFilterDelete
}

// Exec executes the deletion query.
func (rfdo *ReportFilterDeleteOne) Exec(ctx context.Context) error {
	n, err := rfdo.rfd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{reportfilter.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (rfdo *ReportFilterDeleteOne) ExecX(ctx context.Context) {
	rfdo.rfd.ExecX(ctx)
}
