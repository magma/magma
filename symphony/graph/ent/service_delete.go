// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/service"
)

// ServiceDelete is the builder for deleting a Service entity.
type ServiceDelete struct {
	config
	predicates []predicate.Service
}

// Where adds a new predicate to the delete builder.
func (sd *ServiceDelete) Where(ps ...predicate.Service) *ServiceDelete {
	sd.predicates = append(sd.predicates, ps...)
	return sd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (sd *ServiceDelete) Exec(ctx context.Context) (int, error) {
	return sd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (sd *ServiceDelete) ExecX(ctx context.Context) int {
	n, err := sd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (sd *ServiceDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(sd.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(service.Table))
	for _, p := range sd.predicates {
		p(selector)
	}
	query, args := builder.Delete(service.Table).FromSelect(selector).Query()
	if err := sd.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// ServiceDeleteOne is the builder for deleting a single Service entity.
type ServiceDeleteOne struct {
	sd *ServiceDelete
}

// Exec executes the deletion query.
func (sdo *ServiceDeleteOne) Exec(ctx context.Context) error {
	n, err := sdo.sd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{service.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (sdo *ServiceDeleteOne) ExecX(ctx context.Context) {
	sdo.sd.ExecX(ctx)
}
