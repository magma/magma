// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
)

// WorkOrderDefinitionDelete is the builder for deleting a WorkOrderDefinition entity.
type WorkOrderDefinitionDelete struct {
	config
	predicates []predicate.WorkOrderDefinition
}

// Where adds a new predicate to the delete builder.
func (wodd *WorkOrderDefinitionDelete) Where(ps ...predicate.WorkOrderDefinition) *WorkOrderDefinitionDelete {
	wodd.predicates = append(wodd.predicates, ps...)
	return wodd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (wodd *WorkOrderDefinitionDelete) Exec(ctx context.Context) (int, error) {
	return wodd.sqlExec(ctx)
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
	var (
		res     sql.Result
		builder = sql.Dialect(wodd.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(workorderdefinition.Table))
	for _, p := range wodd.predicates {
		p(selector)
	}
	query, args := builder.Delete(workorderdefinition.Table).FromSelect(selector).Query()
	if err := wodd.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
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
		return &ErrNotFound{workorderdefinition.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (woddo *WorkOrderDefinitionDeleteOne) ExecX(ctx context.Context) {
	woddo.wodd.ExecX(ctx)
}
