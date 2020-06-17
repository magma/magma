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
	"github.com/facebookincubator/symphony/graph/ent/usersgroup"
)

// UsersGroupDelete is the builder for deleting a UsersGroup entity.
type UsersGroupDelete struct {
	config
	hooks      []Hook
	mutation   *UsersGroupMutation
	predicates []predicate.UsersGroup
}

// Where adds a new predicate to the delete builder.
func (ugd *UsersGroupDelete) Where(ps ...predicate.UsersGroup) *UsersGroupDelete {
	ugd.predicates = append(ugd.predicates, ps...)
	return ugd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ugd *UsersGroupDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(ugd.hooks) == 0 {
		affected, err = ugd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*UsersGroupMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ugd.mutation = mutation
			affected, err = ugd.sqlExec(ctx)
			return affected, err
		})
		for i := len(ugd.hooks) - 1; i >= 0; i-- {
			mut = ugd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ugd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (ugd *UsersGroupDelete) ExecX(ctx context.Context) int {
	n, err := ugd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ugd *UsersGroupDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: usersgroup.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: usersgroup.FieldID,
			},
		},
	}
	if ps := ugd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, ugd.driver, _spec)
}

// UsersGroupDeleteOne is the builder for deleting a single UsersGroup entity.
type UsersGroupDeleteOne struct {
	ugd *UsersGroupDelete
}

// Exec executes the deletion query.
func (ugdo *UsersGroupDeleteOne) Exec(ctx context.Context) error {
	n, err := ugdo.ugd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{usersgroup.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ugdo *UsersGroupDeleteOne) ExecX(ctx context.Context) {
	ugdo.ugd.ExecX(ctx)
}
