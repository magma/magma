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
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
)

// ServiceEndpointDelete is the builder for deleting a ServiceEndpoint entity.
type ServiceEndpointDelete struct {
	config
	hooks      []Hook
	mutation   *ServiceEndpointMutation
	predicates []predicate.ServiceEndpoint
}

// Where adds a new predicate to the delete builder.
func (sed *ServiceEndpointDelete) Where(ps ...predicate.ServiceEndpoint) *ServiceEndpointDelete {
	sed.predicates = append(sed.predicates, ps...)
	return sed
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (sed *ServiceEndpointDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(sed.hooks) == 0 {
		affected, err = sed.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ServiceEndpointMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			sed.mutation = mutation
			affected, err = sed.sqlExec(ctx)
			return affected, err
		})
		for i := len(sed.hooks); i > 0; i-- {
			mut = sed.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, sed.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (sed *ServiceEndpointDelete) ExecX(ctx context.Context) int {
	n, err := sed.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (sed *ServiceEndpointDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: serviceendpoint.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: serviceendpoint.FieldID,
			},
		},
	}
	if ps := sed.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, sed.driver, _spec)
}

// ServiceEndpointDeleteOne is the builder for deleting a single ServiceEndpoint entity.
type ServiceEndpointDeleteOne struct {
	sed *ServiceEndpointDelete
}

// Exec executes the deletion query.
func (sedo *ServiceEndpointDeleteOne) Exec(ctx context.Context) error {
	n, err := sedo.sed.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{serviceendpoint.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (sedo *ServiceEndpointDeleteOne) ExecX(ctx context.Context) {
	sedo.sed.ExecX(ctx)
}
