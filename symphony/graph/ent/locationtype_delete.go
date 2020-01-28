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
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: locationtype.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: locationtype.FieldID,
			},
		},
	}
	if ps := ltd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, ltd.driver, _spec)
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
		return &NotFoundError{locationtype.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ltdo *LocationTypeDeleteOne) ExecX(ctx context.Context) {
	ltdo.ltd.ExecX(ctx)
}
