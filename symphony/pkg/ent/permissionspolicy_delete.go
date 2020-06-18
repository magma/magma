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
	"github.com/facebookincubator/symphony/pkg/ent/permissionspolicy"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
)

// PermissionsPolicyDelete is the builder for deleting a PermissionsPolicy entity.
type PermissionsPolicyDelete struct {
	config
	hooks      []Hook
	mutation   *PermissionsPolicyMutation
	predicates []predicate.PermissionsPolicy
}

// Where adds a new predicate to the delete builder.
func (ppd *PermissionsPolicyDelete) Where(ps ...predicate.PermissionsPolicy) *PermissionsPolicyDelete {
	ppd.predicates = append(ppd.predicates, ps...)
	return ppd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ppd *PermissionsPolicyDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(ppd.hooks) == 0 {
		affected, err = ppd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PermissionsPolicyMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ppd.mutation = mutation
			affected, err = ppd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(ppd.hooks) - 1; i >= 0; i-- {
			mut = ppd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ppd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (ppd *PermissionsPolicyDelete) ExecX(ctx context.Context) int {
	n, err := ppd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ppd *PermissionsPolicyDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: permissionspolicy.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: permissionspolicy.FieldID,
			},
		},
	}
	if ps := ppd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, ppd.driver, _spec)
}

// PermissionsPolicyDeleteOne is the builder for deleting a single PermissionsPolicy entity.
type PermissionsPolicyDeleteOne struct {
	ppd *PermissionsPolicyDelete
}

// Exec executes the deletion query.
func (ppdo *PermissionsPolicyDeleteOne) Exec(ctx context.Context) error {
	n, err := ppdo.ppd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{permissionspolicy.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ppdo *PermissionsPolicyDeleteOne) ExecX(ctx context.Context) {
	ppdo.ppd.ExecX(ctx)
}
