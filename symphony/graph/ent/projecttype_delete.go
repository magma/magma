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
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
)

// ProjectTypeDelete is the builder for deleting a ProjectType entity.
type ProjectTypeDelete struct {
	config
	predicates []predicate.ProjectType
}

// Where adds a new predicate to the delete builder.
func (ptd *ProjectTypeDelete) Where(ps ...predicate.ProjectType) *ProjectTypeDelete {
	ptd.predicates = append(ptd.predicates, ps...)
	return ptd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ptd *ProjectTypeDelete) Exec(ctx context.Context) (int, error) {
	return ptd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (ptd *ProjectTypeDelete) ExecX(ctx context.Context) int {
	n, err := ptd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ptd *ProjectTypeDelete) sqlExec(ctx context.Context) (int, error) {
	spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: projecttype.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: projecttype.FieldID,
			},
		},
	}
	if ps := ptd.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, ptd.driver, spec)
}

// ProjectTypeDeleteOne is the builder for deleting a single ProjectType entity.
type ProjectTypeDeleteOne struct {
	ptd *ProjectTypeDelete
}

// Exec executes the deletion query.
func (ptdo *ProjectTypeDeleteOne) Exec(ctx context.Context) error {
	n, err := ptdo.ptd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{projecttype.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ptdo *ProjectTypeDeleteOne) ExecX(ctx context.Context) {
	ptdo.ptd.ExecX(ctx)
}
