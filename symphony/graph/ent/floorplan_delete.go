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
	"github.com/facebookincubator/symphony/graph/ent/floorplan"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// FloorPlanDelete is the builder for deleting a FloorPlan entity.
type FloorPlanDelete struct {
	config
	predicates []predicate.FloorPlan
}

// Where adds a new predicate to the delete builder.
func (fpd *FloorPlanDelete) Where(ps ...predicate.FloorPlan) *FloorPlanDelete {
	fpd.predicates = append(fpd.predicates, ps...)
	return fpd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (fpd *FloorPlanDelete) Exec(ctx context.Context) (int, error) {
	return fpd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (fpd *FloorPlanDelete) ExecX(ctx context.Context) int {
	n, err := fpd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (fpd *FloorPlanDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: floorplan.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: floorplan.FieldID,
			},
		},
	}
	if ps := fpd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, fpd.driver, _spec)
}

// FloorPlanDeleteOne is the builder for deleting a single FloorPlan entity.
type FloorPlanDeleteOne struct {
	fpd *FloorPlanDelete
}

// Exec executes the deletion query.
func (fpdo *FloorPlanDeleteOne) Exec(ctx context.Context) error {
	n, err := fpdo.fpd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{floorplan.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (fpdo *FloorPlanDeleteOne) ExecX(ctx context.Context) {
	fpdo.fpd.ExecX(ctx)
}
