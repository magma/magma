// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/frontier/ent/auditlog"
	"github.com/facebookincubator/symphony/frontier/ent/predicate"
)

// AuditLogDelete is the builder for deleting a AuditLog entity.
type AuditLogDelete struct {
	config
	predicates []predicate.AuditLog
}

// Where adds a new predicate to the delete builder.
func (ald *AuditLogDelete) Where(ps ...predicate.AuditLog) *AuditLogDelete {
	ald.predicates = append(ald.predicates, ps...)
	return ald
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ald *AuditLogDelete) Exec(ctx context.Context) (int, error) {
	return ald.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (ald *AuditLogDelete) ExecX(ctx context.Context) int {
	n, err := ald.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ald *AuditLogDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(ald.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(auditlog.Table))
	for _, p := range ald.predicates {
		p(selector)
	}
	query, args := builder.Delete(auditlog.Table).FromSelect(selector).Query()
	if err := ald.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// AuditLogDeleteOne is the builder for deleting a single AuditLog entity.
type AuditLogDeleteOne struct {
	ald *AuditLogDelete
}

// Exec executes the deletion query.
func (aldo *AuditLogDeleteOne) Exec(ctx context.Context) error {
	n, err := aldo.ald.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{auditlog.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (aldo *AuditLogDeleteOne) ExecX(ctx context.Context) {
	aldo.ald.ExecX(ctx)
}
