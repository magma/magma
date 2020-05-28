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
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
)

// WorkOrderDefinitionDelete is the builder for deleting a WorkOrderDefinition entity.
type WorkOrderDefinitionDelete struct {
	config
	hooks      []Hook
	mutation   *WorkOrderDefinitionMutation
	predicates []predicate.WorkOrderDefinition
}

// Where adds a new predicate to the delete builder.
func (wodd *WorkOrderDefinitionDelete) Where(ps ...predicate.WorkOrderDefinition) *WorkOrderDefinitionDelete {
	wodd.predicates = append(wodd.predicates, ps...)
	return wodd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (wodd *WorkOrderDefinitionDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(wodd.hooks) == 0 {
		affected, err = wodd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*WorkOrderDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			wodd.mutation = mutation
			affected, err = wodd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(wodd.hooks) - 1; i >= 0; i-- {
			mut = wodd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, wodd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (wodd *WorkOrderDefinitionDelete) ExecX(ctx context.Context) int {
	n, err := wodd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (wodd *WorkOrderDefinitionDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: workorderdefinition.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: workorderdefinition.FieldID,
			},
		},
	}
	if ps := wodd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, wodd.driver, _spec)
}

// WorkOrderDefinitionDeleteOne is the builder for deleting a single WorkOrderDefinition entity.
type WorkOrderDefinitionDeleteOne struct {
	wodd *WorkOrderDefinitionDelete
}

// Exec executes the deletion query.
func (woddo *WorkOrderDefinitionDeleteOne) Exec(ctx context.Context) error {
	n, err := woddo.wodd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{workorderdefinition.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (woddo *WorkOrderDefinitionDeleteOne) ExecX(ctx context.Context) {
	woddo.wodd.ExecX(ctx)
}
