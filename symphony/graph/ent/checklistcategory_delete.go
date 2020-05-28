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
	"github.com/facebookincubator/symphony/graph/ent/checklistcategory"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// CheckListCategoryDelete is the builder for deleting a CheckListCategory entity.
type CheckListCategoryDelete struct {
	config
	hooks      []Hook
	mutation   *CheckListCategoryMutation
	predicates []predicate.CheckListCategory
}

// Where adds a new predicate to the delete builder.
func (clcd *CheckListCategoryDelete) Where(ps ...predicate.CheckListCategory) *CheckListCategoryDelete {
	clcd.predicates = append(clcd.predicates, ps...)
	return clcd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (clcd *CheckListCategoryDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(clcd.hooks) == 0 {
		affected, err = clcd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CheckListCategoryMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			clcd.mutation = mutation
			affected, err = clcd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(clcd.hooks) - 1; i >= 0; i-- {
			mut = clcd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, clcd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (clcd *CheckListCategoryDelete) ExecX(ctx context.Context) int {
	n, err := clcd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (clcd *CheckListCategoryDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: checklistcategory.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistcategory.FieldID,
			},
		},
	}
	if ps := clcd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, clcd.driver, _spec)
}

// CheckListCategoryDeleteOne is the builder for deleting a single CheckListCategory entity.
type CheckListCategoryDeleteOne struct {
	clcd *CheckListCategoryDelete
}

// Exec executes the deletion query.
func (clcdo *CheckListCategoryDeleteOne) Exec(ctx context.Context) error {
	n, err := clcdo.clcd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{checklistcategory.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (clcdo *CheckListCategoryDeleteOne) ExecX(ctx context.Context) {
	clcdo.clcd.ExecX(ctx)
}
