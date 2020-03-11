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
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: location.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: location.FieldID,
			},
		},
	}
	if ps := ld.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, ld.driver, _spec)
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
		return &NotFoundError{location.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ldo *LocationDeleteOne) ExecX(ctx context.Context) {
	ldo.ld.ExecX(ctx)
}
