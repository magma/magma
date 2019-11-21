// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// LinkDelete is the builder for deleting a Link entity.
type LinkDelete struct {
	config
	predicates []predicate.Link
}

// Where adds a new predicate to the delete builder.
func (ld *LinkDelete) Where(ps ...predicate.Link) *LinkDelete {
	ld.predicates = append(ld.predicates, ps...)
	return ld
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ld *LinkDelete) Exec(ctx context.Context) (int, error) {
	return ld.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (ld *LinkDelete) ExecX(ctx context.Context) int {
	n, err := ld.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ld *LinkDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(ld.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(link.Table))
	for _, p := range ld.predicates {
		p(selector)
	}
	query, args := builder.Delete(link.Table).FromSelect(selector).Query()
	if err := ld.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// LinkDeleteOne is the builder for deleting a single Link entity.
type LinkDeleteOne struct {
	ld *LinkDelete
}

// Exec executes the deletion query.
func (ldo *LinkDeleteOne) Exec(ctx context.Context) error {
	n, err := ldo.ld.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{link.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ldo *LinkDeleteOne) ExecX(ctx context.Context) {
	ldo.ld.ExecX(ctx)
}
