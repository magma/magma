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
	"github.com/facebookincubator/symphony/pkg/ent/serviceendpointdefinition"
)

// ServiceEndpointDefinitionDelete is the builder for deleting a ServiceEndpointDefinition entity.
type ServiceEndpointDefinitionDelete struct {
	config
	hooks      []Hook
	mutation   *ServiceEndpointDefinitionMutation
	predicates []predicate.ServiceEndpointDefinition
}

// Where adds a new predicate to the delete builder.
func (sedd *ServiceEndpointDefinitionDelete) Where(ps ...predicate.ServiceEndpointDefinition) *ServiceEndpointDefinitionDelete {
	sedd.predicates = append(sedd.predicates, ps...)
	return sedd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (sedd *ServiceEndpointDefinitionDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(sedd.hooks) == 0 {
		affected, err = sedd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ServiceEndpointDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			sedd.mutation = mutation
			affected, err = sedd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(sedd.hooks) - 1; i >= 0; i-- {
			mut = sedd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, sedd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (sedd *ServiceEndpointDefinitionDelete) ExecX(ctx context.Context) int {
	n, err := sedd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (sedd *ServiceEndpointDefinitionDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: serviceendpointdefinition.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: serviceendpointdefinition.FieldID,
			},
		},
	}
	if ps := sedd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, sedd.driver, _spec)
}

// ServiceEndpointDefinitionDeleteOne is the builder for deleting a single ServiceEndpointDefinition entity.
type ServiceEndpointDefinitionDeleteOne struct {
	sedd *ServiceEndpointDefinitionDelete
}

// Exec executes the deletion query.
func (seddo *ServiceEndpointDefinitionDeleteOne) Exec(ctx context.Context) error {
	n, err := seddo.sedd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{serviceendpointdefinition.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (seddo *ServiceEndpointDefinitionDeleteOne) ExecX(ctx context.Context) {
	seddo.sedd.ExecX(ctx)
}
