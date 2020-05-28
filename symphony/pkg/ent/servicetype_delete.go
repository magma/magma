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
	"github.com/facebookincubator/symphony/pkg/ent/servicetype"
)

// ServiceTypeDelete is the builder for deleting a ServiceType entity.
type ServiceTypeDelete struct {
	config
	hooks      []Hook
	mutation   *ServiceTypeMutation
	predicates []predicate.ServiceType
}

// Where adds a new predicate to the delete builder.
func (std *ServiceTypeDelete) Where(ps ...predicate.ServiceType) *ServiceTypeDelete {
	std.predicates = append(std.predicates, ps...)
	return std
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (std *ServiceTypeDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(std.hooks) == 0 {
		affected, err = std.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ServiceTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			std.mutation = mutation
			affected, err = std.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(std.hooks) - 1; i >= 0; i-- {
			mut = std.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, std.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (std *ServiceTypeDelete) ExecX(ctx context.Context) int {
	n, err := std.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (std *ServiceTypeDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: servicetype.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: servicetype.FieldID,
			},
		},
	}
	if ps := std.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, std.driver, _spec)
}

// ServiceTypeDeleteOne is the builder for deleting a single ServiceType entity.
type ServiceTypeDeleteOne struct {
	std *ServiceTypeDelete
}

// Exec executes the deletion query.
func (stdo *ServiceTypeDeleteOne) Exec(ctx context.Context) error {
	n, err := stdo.std.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{servicetype.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (stdo *ServiceTypeDeleteOne) ExecX(ctx context.Context) {
	stdo.std.ExecX(ctx)
}
