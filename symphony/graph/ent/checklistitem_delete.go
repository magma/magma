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
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// CheckListItemDelete is the builder for deleting a CheckListItem entity.
type CheckListItemDelete struct {
	config
	hooks      []Hook
	mutation   *CheckListItemMutation
	predicates []predicate.CheckListItem
}

// Where adds a new predicate to the delete builder.
func (clid *CheckListItemDelete) Where(ps ...predicate.CheckListItem) *CheckListItemDelete {
	clid.predicates = append(clid.predicates, ps...)
	return clid
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (clid *CheckListItemDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(clid.hooks) == 0 {
		affected, err = clid.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CheckListItemMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			clid.mutation = mutation
			affected, err = clid.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(clid.hooks) - 1; i >= 0; i-- {
			mut = clid.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, clid.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (clid *CheckListItemDelete) ExecX(ctx context.Context) int {
	n, err := clid.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (clid *CheckListItemDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: checklistitem.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistitem.FieldID,
			},
		},
	}
	if ps := clid.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, clid.driver, _spec)
}

// CheckListItemDeleteOne is the builder for deleting a single CheckListItem entity.
type CheckListItemDeleteOne struct {
	clid *CheckListItemDelete
}

// Exec executes the deletion query.
func (clido *CheckListItemDeleteOne) Exec(ctx context.Context) error {
	n, err := clido.clid.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{checklistitem.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (clido *CheckListItemDeleteOne) ExecX(ctx context.Context) {
	clido.clid.ExecX(ctx)
}
