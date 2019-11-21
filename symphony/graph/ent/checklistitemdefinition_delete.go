// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/checklistitemdefinition"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// CheckListItemDefinitionDelete is the builder for deleting a CheckListItemDefinition entity.
type CheckListItemDefinitionDelete struct {
	config
	predicates []predicate.CheckListItemDefinition
}

// Where adds a new predicate to the delete builder.
func (clidd *CheckListItemDefinitionDelete) Where(ps ...predicate.CheckListItemDefinition) *CheckListItemDefinitionDelete {
	clidd.predicates = append(clidd.predicates, ps...)
	return clidd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (clidd *CheckListItemDefinitionDelete) Exec(ctx context.Context) (int, error) {
	return clidd.sqlExec(ctx)
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
	var (
		res     sql.Result
		builder = sql.Dialect(clidd.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(checklistitemdefinition.Table))
	for _, p := range clidd.predicates {
		p(selector)
	}
	query, args := builder.Delete(checklistitemdefinition.Table).FromSelect(selector).Query()
	if err := clidd.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
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
		return &ErrNotFound{checklistitemdefinition.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (cliddo *CheckListItemDefinitionDeleteOne) ExecX(ctx context.Context) {
	cliddo.clidd.ExecX(ctx)
}
