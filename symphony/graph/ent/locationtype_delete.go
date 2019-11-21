// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// LocationTypeDelete is the builder for deleting a LocationType entity.
type LocationTypeDelete struct {
	config
	predicates []predicate.LocationType
}

// Where adds a new predicate to the delete builder.
func (ltd *LocationTypeDelete) Where(ps ...predicate.LocationType) *LocationTypeDelete {
	ltd.predicates = append(ltd.predicates, ps...)
	return ltd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ltd *LocationTypeDelete) Exec(ctx context.Context) (int, error) {
	return ltd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (ltd *LocationTypeDelete) ExecX(ctx context.Context) int {
	n, err := ltd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ltd *LocationTypeDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(ltd.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(locationtype.Table))
	for _, p := range ltd.predicates {
		p(selector)
	}
	query, args := builder.Delete(locationtype.Table).FromSelect(selector).Query()
	if err := ltd.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// LocationTypeDeleteOne is the builder for deleting a single LocationType entity.
type LocationTypeDeleteOne struct {
	ltd *LocationTypeDelete
}

// Exec executes the deletion query.
func (ltdo *LocationTypeDeleteOne) Exec(ctx context.Context) error {
	n, err := ltdo.ltd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{locationtype.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ltdo *LocationTypeDeleteOne) ExecX(ctx context.Context) {
	ltdo.ltd.ExecX(ctx)
}
