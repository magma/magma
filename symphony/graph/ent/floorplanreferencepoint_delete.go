// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/floorplanreferencepoint"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// FloorPlanReferencePointDelete is the builder for deleting a FloorPlanReferencePoint entity.
type FloorPlanReferencePointDelete struct {
	config
	predicates []predicate.FloorPlanReferencePoint
}

// Where adds a new predicate to the delete builder.
func (fprpd *FloorPlanReferencePointDelete) Where(ps ...predicate.FloorPlanReferencePoint) *FloorPlanReferencePointDelete {
	fprpd.predicates = append(fprpd.predicates, ps...)
	return fprpd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (fprpd *FloorPlanReferencePointDelete) Exec(ctx context.Context) (int, error) {
	return fprpd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (fprpd *FloorPlanReferencePointDelete) ExecX(ctx context.Context) int {
	n, err := fprpd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (fprpd *FloorPlanReferencePointDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(fprpd.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(floorplanreferencepoint.Table))
	for _, p := range fprpd.predicates {
		p(selector)
	}
	query, args := builder.Delete(floorplanreferencepoint.Table).FromSelect(selector).Query()
	if err := fprpd.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// FloorPlanReferencePointDeleteOne is the builder for deleting a single FloorPlanReferencePoint entity.
type FloorPlanReferencePointDeleteOne struct {
	fprpd *FloorPlanReferencePointDelete
}

// Exec executes the deletion query.
func (fprpdo *FloorPlanReferencePointDeleteOne) Exec(ctx context.Context) error {
	n, err := fprpdo.fprpd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{floorplanreferencepoint.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (fprpdo *FloorPlanReferencePointDeleteOne) ExecX(ctx context.Context) {
	fprpdo.fprpd.ExecX(ctx)
}
