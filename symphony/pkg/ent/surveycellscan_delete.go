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
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/surveycellscan"
)

// SurveyCellScanDelete is the builder for deleting a SurveyCellScan entity.
type SurveyCellScanDelete struct {
	config
	hooks      []Hook
	mutation   *SurveyCellScanMutation
	predicates []predicate.SurveyCellScan
}

// Where adds a new predicate to the delete builder.
func (scsd *SurveyCellScanDelete) Where(ps ...predicate.SurveyCellScan) *SurveyCellScanDelete {
	scsd.predicates = append(scsd.predicates, ps...)
	return scsd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (scsd *SurveyCellScanDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(scsd.hooks) == 0 {
		affected, err = scsd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyCellScanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			scsd.mutation = mutation
			affected, err = scsd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(scsd.hooks) - 1; i >= 0; i-- {
			mut = scsd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, scsd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (scsd *SurveyCellScanDelete) ExecX(ctx context.Context) int {
	n, err := scsd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (scsd *SurveyCellScanDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: surveycellscan.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: surveycellscan.FieldID,
			},
		},
	}
	if ps := scsd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, scsd.driver, _spec)
}

// SurveyCellScanDeleteOne is the builder for deleting a single SurveyCellScan entity.
type SurveyCellScanDeleteOne struct {
	scsd *SurveyCellScanDelete
}

// Exec executes the deletion query.
func (scsdo *SurveyCellScanDeleteOne) Exec(ctx context.Context) error {
	n, err := scsdo.scsd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{surveycellscan.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (scsdo *SurveyCellScanDeleteOne) ExecX(ctx context.Context) {
	scsdo.scsd.ExecX(ctx)
}
