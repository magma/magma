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
	"github.com/facebookincubator/symphony/pkg/ent/checklistitemdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
)

// CheckListItemDefinitionDelete is the builder for deleting a CheckListItemDefinition entity.
type CheckListItemDefinitionDelete struct {
	config
	hooks      []Hook
	mutation   *CheckListItemDefinitionMutation
	predicates []predicate.CheckListItemDefinition
}

// Where adds a new predicate to the delete builder.
func (clidd *CheckListItemDefinitionDelete) Where(ps ...predicate.CheckListItemDefinition) *CheckListItemDefinitionDelete {
	clidd.predicates = append(clidd.predicates, ps...)
	return clidd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (clidd *CheckListItemDefinitionDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(clidd.hooks) == 0 {
		affected, err = clidd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CheckListItemDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			clidd.mutation = mutation
			affected, err = clidd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(clidd.hooks) - 1; i >= 0; i-- {
			mut = clidd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, clidd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (clidd *CheckListItemDefinitionDelete) ExecX(ctx context.Context) int {
	n, err := clidd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (clidd *CheckListItemDefinitionDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: checklistitemdefinition.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistitemdefinition.FieldID,
			},
		},
	}
	if ps := clidd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, clidd.driver, _spec)
}

// CheckListItemDefinitionDeleteOne is the builder for deleting a single CheckListItemDefinition entity.
type CheckListItemDefinitionDeleteOne struct {
	clidd *CheckListItemDefinitionDelete
}

// Exec executes the deletion query.
func (cliddo *CheckListItemDefinitionDeleteOne) Exec(ctx context.Context) error {
	n, err := cliddo.clidd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{checklistitemdefinition.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (cliddo *CheckListItemDefinitionDeleteOne) ExecX(ctx context.Context) {
	cliddo.clidd.ExecX(ctx)
}
