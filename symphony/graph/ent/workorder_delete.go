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
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// WorkOrderDelete is the builder for deleting a WorkOrder entity.
type WorkOrderDelete struct {
	config
	hooks      []Hook
	mutation   *WorkOrderMutation
	predicates []predicate.WorkOrder
}

// Where adds a new predicate to the delete builder.
func (wod *WorkOrderDelete) Where(ps ...predicate.WorkOrder) *WorkOrderDelete {
	wod.predicates = append(wod.predicates, ps...)
	return wod
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (wod *WorkOrderDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(wod.hooks) == 0 {
		affected, err = wod.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*WorkOrderMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			wod.mutation = mutation
			affected, err = wod.sqlExec(ctx)
			return affected, err
		})
		for i := len(wod.hooks); i > 0; i-- {
			mut = wod.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, wod.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (wod *WorkOrderDelete) ExecX(ctx context.Context) int {
	n, err := wod.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (wod *WorkOrderDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: workorder.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: workorder.FieldID,
			},
		},
	}
	if ps := wod.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, wod.driver, _spec)
}

// WorkOrderDeleteOne is the builder for deleting a single WorkOrder entity.
type WorkOrderDeleteOne struct {
	wod *WorkOrderDelete
}

// Exec executes the deletion query.
func (wodo *WorkOrderDeleteOne) Exec(ctx context.Context) error {
	n, err := wodo.wod.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{workorder.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (wodo *WorkOrderDeleteOne) ExecX(ctx context.Context) {
	wodo.wod.ExecX(ctx)
}
