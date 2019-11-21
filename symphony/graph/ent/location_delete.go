// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// LocationDelete is the builder for deleting a Location entity.
type LocationDelete struct {
	config
	predicates []predicate.Location
}

// Where adds a new predicate to the delete builder.
func (ld *LocationDelete) Where(ps ...predicate.Location) *LocationDelete {
	ld.predicates = append(ld.predicates, ps...)
	return ld
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ld *LocationDelete) Exec(ctx context.Context) (int, error) {
	return ld.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (ld *LocationDelete) ExecX(ctx context.Context) int {
	n, err := ld.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ld *LocationDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(ld.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(location.Table))
	for _, p := range ld.predicates {
		p(selector)
	}
	query, args := builder.Delete(location.Table).FromSelect(selector).Query()
	if err := ld.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// LocationDeleteOne is the builder for deleting a single Location entity.
type LocationDeleteOne struct {
	ld *LocationDelete
}

// Exec executes the deletion query.
func (ldo *LocationDeleteOne) Exec(ctx context.Context) error {
	n, err := ldo.ld.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{location.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ldo *LocationDeleteOne) ExecX(ctx context.Context) {
	ldo.ld.ExecX(ctx)
}
