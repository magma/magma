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
	"github.com/facebookincubator/symphony/graph/ent/actionsrule"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ActionsRuleDelete is the builder for deleting a ActionsRule entity.
type ActionsRuleDelete struct {
	config
	hooks      []Hook
	mutation   *ActionsRuleMutation
	predicates []predicate.ActionsRule
}

// Where adds a new predicate to the delete builder.
func (ard *ActionsRuleDelete) Where(ps ...predicate.ActionsRule) *ActionsRuleDelete {
	ard.predicates = append(ard.predicates, ps...)
	return ard
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ard *ActionsRuleDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(ard.hooks) == 0 {
		affected, err = ard.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ActionsRuleMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ard.mutation = mutation
			affected, err = ard.sqlExec(ctx)
			return affected, err
		})
		for i := len(ard.hooks) - 1; i >= 0; i-- {
			mut = ard.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ard.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (ard *ActionsRuleDelete) ExecX(ctx context.Context) int {
	n, err := ard.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ard *ActionsRuleDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: actionsrule.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: actionsrule.FieldID,
			},
		},
	}
	if ps := ard.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, ard.driver, _spec)
}

// ActionsRuleDeleteOne is the builder for deleting a single ActionsRule entity.
type ActionsRuleDeleteOne struct {
	ard *ActionsRuleDelete
}

// Exec executes the deletion query.
func (ardo *ActionsRuleDeleteOne) Exec(ctx context.Context) error {
	n, err := ardo.ard.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{actionsrule.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ardo *ActionsRuleDeleteOne) ExecX(ctx context.Context) {
	ardo.ard.ExecX(ctx)
}
