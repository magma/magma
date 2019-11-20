// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// WorkOrderTypeDelete is the builder for deleting a WorkOrderType entity.
type WorkOrderTypeDelete struct {
	config
	predicates []predicate.WorkOrderType
}

// Where adds a new predicate to the delete builder.
func (wotd *WorkOrderTypeDelete) Where(ps ...predicate.WorkOrderType) *WorkOrderTypeDelete {
	wotd.predicates = append(wotd.predicates, ps...)
	return wotd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (wotd *WorkOrderTypeDelete) Exec(ctx context.Context) (int, error) {
	return wotd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (wotd *WorkOrderTypeDelete) ExecX(ctx context.Context) int {
	n, err := wotd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (wotd *WorkOrderTypeDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(wotd.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(workordertype.Table))
	for _, p := range wotd.predicates {
		p(selector)
	}
	query, args := builder.Delete(workordertype.Table).FromSelect(selector).Query()
	if err := wotd.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// WorkOrderTypeDeleteOne is the builder for deleting a single WorkOrderType entity.
type WorkOrderTypeDeleteOne struct {
	wotd *WorkOrderTypeDelete
}

// Exec executes the deletion query.
func (wotdo *WorkOrderTypeDeleteOne) Exec(ctx context.Context) error {
	n, err := wotdo.wotd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{workordertype.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (wotdo *WorkOrderTypeDeleteOne) ExecX(ctx context.Context) {
	wotdo.wotd.ExecX(ctx)
}
