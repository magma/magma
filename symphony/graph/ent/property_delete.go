// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
)

// PropertyDelete is the builder for deleting a Property entity.
type PropertyDelete struct {
	config
	predicates []predicate.Property
}

// Where adds a new predicate to the delete builder.
func (pd *PropertyDelete) Where(ps ...predicate.Property) *PropertyDelete {
	pd.predicates = append(pd.predicates, ps...)
	return pd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (pd *PropertyDelete) Exec(ctx context.Context) (int, error) {
	return pd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (pd *PropertyDelete) ExecX(ctx context.Context) int {
	n, err := pd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (pd *PropertyDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: property.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: property.FieldID,
			},
		},
	}
	if ps := pd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, pd.driver, _spec)
}

// PropertyDeleteOne is the builder for deleting a single Property entity.
type PropertyDeleteOne struct {
	pd *PropertyDelete
}

// Exec executes the deletion query.
func (pdo *PropertyDeleteOne) Exec(ctx context.Context) error {
	n, err := pdo.pd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{property.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (pdo *PropertyDeleteOne) ExecX(ctx context.Context) {
	pdo.pd.ExecX(ctx)
}
