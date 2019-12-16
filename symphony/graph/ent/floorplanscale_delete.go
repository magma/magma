// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/floorplanscale"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// FloorPlanScaleDelete is the builder for deleting a FloorPlanScale entity.
type FloorPlanScaleDelete struct {
	config
	predicates []predicate.FloorPlanScale
}

// Where adds a new predicate to the delete builder.
func (fpsd *FloorPlanScaleDelete) Where(ps ...predicate.FloorPlanScale) *FloorPlanScaleDelete {
	fpsd.predicates = append(fpsd.predicates, ps...)
	return fpsd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (fpsd *FloorPlanScaleDelete) Exec(ctx context.Context) (int, error) {
	return fpsd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (fpsd *FloorPlanScaleDelete) ExecX(ctx context.Context) int {
	n, err := fpsd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (fpsd *FloorPlanScaleDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(fpsd.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(floorplanscale.Table))
	for _, p := range fpsd.predicates {
		p(selector)
	}
	query, args := builder.Delete(floorplanscale.Table).FromSelect(selector).Query()
	if err := fpsd.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// FloorPlanScaleDeleteOne is the builder for deleting a single FloorPlanScale entity.
type FloorPlanScaleDeleteOne struct {
	fpsd *FloorPlanScaleDelete
}

// Exec executes the deletion query.
func (fpsdo *FloorPlanScaleDeleteOne) Exec(ctx context.Context) error {
	n, err := fpsdo.fpsd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{floorplanscale.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (fpsdo *FloorPlanScaleDeleteOne) ExecX(ctx context.Context) {
	fpsdo.fpsd.ExecX(ctx)
}
