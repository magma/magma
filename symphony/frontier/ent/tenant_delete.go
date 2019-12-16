// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/frontier/ent/predicate"
	"github.com/facebookincubator/symphony/frontier/ent/tenant"
)

// TenantDelete is the builder for deleting a Tenant entity.
type TenantDelete struct {
	config
	predicates []predicate.Tenant
}

// Where adds a new predicate to the delete builder.
func (td *TenantDelete) Where(ps ...predicate.Tenant) *TenantDelete {
	td.predicates = append(td.predicates, ps...)
	return td
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (td *TenantDelete) Exec(ctx context.Context) (int, error) {
	return td.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (td *TenantDelete) ExecX(ctx context.Context) int {
	n, err := td.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (td *TenantDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(td.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(tenant.Table))
	for _, p := range td.predicates {
		p(selector)
	}
	query, args := builder.Delete(tenant.Table).FromSelect(selector).Query()
	if err := td.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// TenantDeleteOne is the builder for deleting a single Tenant entity.
type TenantDeleteOne struct {
	td *TenantDelete
}

// Exec executes the deletion query.
func (tdo *TenantDeleteOne) Exec(ctx context.Context) error {
	n, err := tdo.td.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{tenant.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (tdo *TenantDeleteOne) ExecX(ctx context.Context) {
	tdo.td.ExecX(ctx)
}
