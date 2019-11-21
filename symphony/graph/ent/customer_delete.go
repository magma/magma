// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/customer"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// CustomerDelete is the builder for deleting a Customer entity.
type CustomerDelete struct {
	config
	predicates []predicate.Customer
}

// Where adds a new predicate to the delete builder.
func (cd *CustomerDelete) Where(ps ...predicate.Customer) *CustomerDelete {
	cd.predicates = append(cd.predicates, ps...)
	return cd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (cd *CustomerDelete) Exec(ctx context.Context) (int, error) {
	return cd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (cd *CustomerDelete) ExecX(ctx context.Context) int {
	n, err := cd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (cd *CustomerDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(cd.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(customer.Table))
	for _, p := range cd.predicates {
		p(selector)
	}
	query, args := builder.Delete(customer.Table).FromSelect(selector).Query()
	if err := cd.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// CustomerDeleteOne is the builder for deleting a single Customer entity.
type CustomerDeleteOne struct {
	cd *CustomerDelete
}

// Exec executes the deletion query.
func (cdo *CustomerDeleteOne) Exec(ctx context.Context) error {
	n, err := cdo.cd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{customer.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (cdo *CustomerDeleteOne) ExecX(ctx context.Context) {
	cdo.cd.ExecX(ctx)
}
