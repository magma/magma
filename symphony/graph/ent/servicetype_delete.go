// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
)

// ServiceTypeDelete is the builder for deleting a ServiceType entity.
type ServiceTypeDelete struct {
	config
	predicates []predicate.ServiceType
}

// Where adds a new predicate to the delete builder.
func (std *ServiceTypeDelete) Where(ps ...predicate.ServiceType) *ServiceTypeDelete {
	std.predicates = append(std.predicates, ps...)
	return std
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (std *ServiceTypeDelete) Exec(ctx context.Context) (int, error) {
	return std.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (std *ServiceTypeDelete) ExecX(ctx context.Context) int {
	n, err := std.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (std *ServiceTypeDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(std.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(servicetype.Table))
	for _, p := range std.predicates {
		p(selector)
	}
	query, args := builder.Delete(servicetype.Table).FromSelect(selector).Query()
	if err := std.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// ServiceTypeDeleteOne is the builder for deleting a single ServiceType entity.
type ServiceTypeDeleteOne struct {
	std *ServiceTypeDelete
}

// Exec executes the deletion query.
func (stdo *ServiceTypeDeleteOne) Exec(ctx context.Context) error {
	n, err := stdo.std.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{servicetype.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (stdo *ServiceTypeDeleteOne) ExecX(ctx context.Context) {
	stdo.std.ExecX(ctx)
}
