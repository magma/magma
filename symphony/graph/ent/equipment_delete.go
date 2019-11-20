// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// EquipmentDelete is the builder for deleting a Equipment entity.
type EquipmentDelete struct {
	config
	predicates []predicate.Equipment
}

// Where adds a new predicate to the delete builder.
func (ed *EquipmentDelete) Where(ps ...predicate.Equipment) *EquipmentDelete {
	ed.predicates = append(ed.predicates, ps...)
	return ed
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ed *EquipmentDelete) Exec(ctx context.Context) (int, error) {
	return ed.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (ed *EquipmentDelete) ExecX(ctx context.Context) int {
	n, err := ed.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ed *EquipmentDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(ed.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(equipment.Table))
	for _, p := range ed.predicates {
		p(selector)
	}
	query, args := builder.Delete(equipment.Table).FromSelect(selector).Query()
	if err := ed.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// EquipmentDeleteOne is the builder for deleting a single Equipment entity.
type EquipmentDeleteOne struct {
	ed *EquipmentDelete
}

// Exec executes the deletion query.
func (edo *EquipmentDeleteOne) Exec(ctx context.Context) error {
	n, err := edo.ed.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{equipment.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (edo *EquipmentDeleteOne) ExecX(ctx context.Context) {
	edo.ed.ExecX(ctx)
}
