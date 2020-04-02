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
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// WorkOrderTypeDelete is the builder for deleting a WorkOrderType entity.
type WorkOrderTypeDelete struct {
	config
	hooks      []Hook
	mutation   *WorkOrderTypeMutation
	predicates []predicate.WorkOrderType
}

// Where adds a new predicate to the delete builder.
func (wotd *WorkOrderTypeDelete) Where(ps ...predicate.WorkOrderType) *WorkOrderTypeDelete {
	wotd.predicates = append(wotd.predicates, ps...)
	return wotd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (wotd *WorkOrderTypeDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(wotd.hooks) == 0 {
		affected, err = wotd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*WorkOrderTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			wotd.mutation = mutation
			affected, err = wotd.sqlExec(ctx)
			return affected, err
		})
		for i := len(wotd.hooks) - 1; i >= 0; i-- {
			mut = wotd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, wotd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (wotd *WorkOrderTypeDelete) ExecX(ctx context.Context) int {
	n, err := wotd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (wotd *WorkOrderTypeDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: workordertype.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: workordertype.FieldID,
			},
		},
	}
	if ps := wotd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, wotd.driver, _spec)
}

// WorkOrderTypeDeleteOne is the builder for deleting a single WorkOrderType entity.
type WorkOrderTypeDeleteOne struct {
	wotd *WorkOrderTypeDelete
}

// Exec executes the deletion query.
func (wotdo *WorkOrderTypeDeleteOne) Exec(ctx context.Context) error {
	n, err := wotdo.wotd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{workordertype.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (wotdo *WorkOrderTypeDeleteOne) ExecX(ctx context.Context) {
	wotdo.wotd.ExecX(ctx)
}
