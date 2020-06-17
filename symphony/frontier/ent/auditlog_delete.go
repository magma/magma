// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/frontier/ent/auditlog"
	"github.com/facebookincubator/symphony/frontier/ent/predicate"
)

// AuditLogDelete is the builder for deleting a AuditLog entity.
type AuditLogDelete struct {
	config
	hooks      []Hook
	mutation   *AuditLogMutation
	predicates []predicate.AuditLog
}

// Where adds a new predicate to the delete builder.
func (ald *AuditLogDelete) Where(ps ...predicate.AuditLog) *AuditLogDelete {
	ald.predicates = append(ald.predicates, ps...)
	return ald
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ald *AuditLogDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(ald.hooks) == 0 {
		affected, err = ald.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AuditLogMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ald.mutation = mutation
			affected, err = ald.sqlExec(ctx)
			return affected, err
		})
		for i := len(ald.hooks) - 1; i >= 0; i-- {
			mut = ald.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ald.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: auditlog.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: auditlog.FieldID,
			},
		},
	}
	if ps := ald.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, ald.driver, _spec)
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
		return &NotFoundError{auditlog.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (aldo *AuditLogDeleteOne) ExecX(ctx context.Context) {
	aldo.ald.ExecX(ctx)
}
