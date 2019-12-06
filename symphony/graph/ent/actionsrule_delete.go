// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/actionsrule"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ActionsRuleDelete is the builder for deleting a ActionsRule entity.
type ActionsRuleDelete struct {
	config
	predicates []predicate.ActionsRule
}

// Where adds a new predicate to the delete builder.
func (ard *ActionsRuleDelete) Where(ps ...predicate.ActionsRule) *ActionsRuleDelete {
	ard.predicates = append(ard.predicates, ps...)
	return ard
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ard *ActionsRuleDelete) Exec(ctx context.Context) (int, error) {
	return ard.sqlExec(ctx)
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
	var (
		res     sql.Result
		builder = sql.Dialect(ard.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(actionsrule.Table))
	for _, p := range ard.predicates {
		p(selector)
	}
	query, args := builder.Delete(actionsrule.Table).FromSelect(selector).Query()
	if err := ard.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
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
		return &ErrNotFound{actionsrule.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ardo *ActionsRuleDeleteOne) ExecX(ctx context.Context) {
	ardo.ard.ExecX(ctx)
}
