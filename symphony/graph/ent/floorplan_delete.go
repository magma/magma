// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
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
	var (
		res     sql.Result
		builder = sql.Dialect(fpd.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(floorplan.Table))
	for _, p := range fpd.predicates {
		p(selector)
	}
	query, args := builder.Delete(floorplan.Table).FromSelect(selector).Query()
	if err := fpd.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
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
