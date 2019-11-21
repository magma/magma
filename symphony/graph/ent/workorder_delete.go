// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// WorkOrderDelete is the builder for deleting a WorkOrder entity.
type WorkOrderDelete struct {
	config
	predicates []predicate.WorkOrder
}

// Where adds a new predicate to the delete builder.
func (wod *WorkOrderDelete) Where(ps ...predicate.WorkOrder) *WorkOrderDelete {
	wod.predicates = append(wod.predicates, ps...)
	return wod
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (wod *WorkOrderDelete) Exec(ctx context.Context) (int, error) {
	return wod.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (wod *WorkOrderDelete) ExecX(ctx context.Context) int {
	n, err := wod.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (wod *WorkOrderDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(wod.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(workorder.Table))
	for _, p := range wod.predicates {
		p(selector)
	}
	query, args := builder.Delete(workorder.Table).FromSelect(selector).Query()
	if err := wod.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// WorkOrderDeleteOne is the builder for deleting a single WorkOrder entity.
type WorkOrderDeleteOne struct {
	wod *WorkOrderDelete
}

// Exec executes the deletion query.
func (wodo *WorkOrderDeleteOne) Exec(ctx context.Context) error {
	n, err := wodo.wod.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{workorder.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (wodo *WorkOrderDeleteOne) ExecX(ctx context.Context) {
	wodo.wod.ExecX(ctx)
}
