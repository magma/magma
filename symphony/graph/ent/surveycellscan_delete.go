// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"
)

// SurveyCellScanDelete is the builder for deleting a SurveyCellScan entity.
type SurveyCellScanDelete struct {
	config
	predicates []predicate.SurveyCellScan
}

// Where adds a new predicate to the delete builder.
func (scsd *SurveyCellScanDelete) Where(ps ...predicate.SurveyCellScan) *SurveyCellScanDelete {
	scsd.predicates = append(scsd.predicates, ps...)
	return scsd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (scsd *SurveyCellScanDelete) Exec(ctx context.Context) (int, error) {
	return scsd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (scsd *SurveyCellScanDelete) ExecX(ctx context.Context) int {
	n, err := scsd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (scsd *SurveyCellScanDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(scsd.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(surveycellscan.Table))
	for _, p := range scsd.predicates {
		p(selector)
	}
	query, args := builder.Delete(surveycellscan.Table).FromSelect(selector).Query()
	if err := scsd.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// SurveyCellScanDeleteOne is the builder for deleting a single SurveyCellScan entity.
type SurveyCellScanDeleteOne struct {
	scsd *SurveyCellScanDelete
}

// Exec executes the deletion query.
func (scsdo *SurveyCellScanDeleteOne) Exec(ctx context.Context) error {
	n, err := scsdo.scsd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{surveycellscan.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (scsdo *SurveyCellScanDeleteOne) ExecX(ctx context.Context) {
	scsdo.scsd.ExecX(ctx)
}
