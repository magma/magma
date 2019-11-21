// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/technician"
)

// TechnicianDelete is the builder for deleting a Technician entity.
type TechnicianDelete struct {
	config
	predicates []predicate.Technician
}

// Where adds a new predicate to the delete builder.
func (td *TechnicianDelete) Where(ps ...predicate.Technician) *TechnicianDelete {
	td.predicates = append(td.predicates, ps...)
	return td
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (td *TechnicianDelete) Exec(ctx context.Context) (int, error) {
	return td.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (td *TechnicianDelete) ExecX(ctx context.Context) int {
	n, err := td.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (td *TechnicianDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(td.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(technician.Table))
	for _, p := range td.predicates {
		p(selector)
	}
	query, args := builder.Delete(technician.Table).FromSelect(selector).Query()
	if err := td.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// TechnicianDeleteOne is the builder for deleting a single Technician entity.
type TechnicianDeleteOne struct {
	td *TechnicianDelete
}

// Exec executes the deletion query.
func (tdo *TechnicianDeleteOne) Exec(ctx context.Context) error {
	n, err := tdo.td.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{technician.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (tdo *TechnicianDeleteOne) ExecX(ctx context.Context) {
	tdo.td.ExecX(ctx)
}
