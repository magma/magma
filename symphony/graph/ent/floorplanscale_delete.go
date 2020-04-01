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
	"github.com/facebookincubator/symphony/graph/ent/floorplanscale"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// FloorPlanScaleDelete is the builder for deleting a FloorPlanScale entity.
type FloorPlanScaleDelete struct {
	config
	hooks      []Hook
	mutation   *FloorPlanScaleMutation
	predicates []predicate.FloorPlanScale
}

// Where adds a new predicate to the delete builder.
func (fpsd *FloorPlanScaleDelete) Where(ps ...predicate.FloorPlanScale) *FloorPlanScaleDelete {
	fpsd.predicates = append(fpsd.predicates, ps...)
	return fpsd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (fpsd *FloorPlanScaleDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(fpsd.hooks) == 0 {
		affected, err = fpsd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FloorPlanScaleMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fpsd.mutation = mutation
			affected, err = fpsd.sqlExec(ctx)
			return affected, err
		})
		for i := len(fpsd.hooks) - 1; i >= 0; i-- {
			mut = fpsd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fpsd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (fpsd *FloorPlanScaleDelete) ExecX(ctx context.Context) int {
	n, err := fpsd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (fpsd *FloorPlanScaleDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: floorplanscale.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: floorplanscale.FieldID,
			},
		},
	}
	if ps := fpsd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, fpsd.driver, _spec)
}

// FloorPlanScaleDeleteOne is the builder for deleting a single FloorPlanScale entity.
type FloorPlanScaleDeleteOne struct {
	fpsd *FloorPlanScaleDelete
}

// Exec executes the deletion query.
func (fpsdo *FloorPlanScaleDeleteOne) Exec(ctx context.Context) error {
	n, err := fpsdo.fpsd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{floorplanscale.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (fpsdo *FloorPlanScaleDeleteOne) ExecX(ctx context.Context) {
	fpsdo.fpsd.ExecX(ctx)
}
