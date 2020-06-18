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
	"github.com/facebookincubator/symphony/pkg/ent/checklistcategorydefinition"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
)

// CheckListCategoryDefinitionDelete is the builder for deleting a CheckListCategoryDefinition entity.
type CheckListCategoryDefinitionDelete struct {
	config
	hooks      []Hook
	mutation   *CheckListCategoryDefinitionMutation
	predicates []predicate.CheckListCategoryDefinition
}

// Where adds a new predicate to the delete builder.
func (clcdd *CheckListCategoryDefinitionDelete) Where(ps ...predicate.CheckListCategoryDefinition) *CheckListCategoryDefinitionDelete {
	clcdd.predicates = append(clcdd.predicates, ps...)
	return clcdd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (clcdd *CheckListCategoryDefinitionDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(clcdd.hooks) == 0 {
		affected, err = clcdd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CheckListCategoryDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			clcdd.mutation = mutation
			affected, err = clcdd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(clcdd.hooks) - 1; i >= 0; i-- {
			mut = clcdd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, clcdd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (clcdd *CheckListCategoryDefinitionDelete) ExecX(ctx context.Context) int {
	n, err := clcdd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (clcdd *CheckListCategoryDefinitionDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: checklistcategorydefinition.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistcategorydefinition.FieldID,
			},
		},
	}
	if ps := clcdd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, clcdd.driver, _spec)
}

// CheckListCategoryDefinitionDeleteOne is the builder for deleting a single CheckListCategoryDefinition entity.
type CheckListCategoryDefinitionDeleteOne struct {
	clcdd *CheckListCategoryDefinitionDelete
}

// Exec executes the deletion query.
func (clcddo *CheckListCategoryDefinitionDeleteOne) Exec(ctx context.Context) error {
	n, err := clcddo.clcdd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{checklistcategorydefinition.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (clcddo *CheckListCategoryDefinitionDeleteOne) ExecX(ctx context.Context) {
	clcddo.clcdd.ExecX(ctx)
}
