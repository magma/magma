// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// EquipmentPortDelete is the builder for deleting a EquipmentPort entity.
type EquipmentPortDelete struct {
	config
	predicates []predicate.EquipmentPort
}

// Where adds a new predicate to the delete builder.
func (epd *EquipmentPortDelete) Where(ps ...predicate.EquipmentPort) *EquipmentPortDelete {
	epd.predicates = append(epd.predicates, ps...)
	return epd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (epd *EquipmentPortDelete) Exec(ctx context.Context) (int, error) {
	return epd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (epd *EquipmentPortDelete) ExecX(ctx context.Context) int {
	n, err := epd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (epd *EquipmentPortDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(epd.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(equipmentport.Table))
	for _, p := range epd.predicates {
		p(selector)
	}
	query, args := builder.Delete(equipmentport.Table).FromSelect(selector).Query()
	if err := epd.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// EquipmentPortDeleteOne is the builder for deleting a single EquipmentPort entity.
type EquipmentPortDeleteOne struct {
	epd *EquipmentPortDelete
}

// Exec executes the deletion query.
func (epdo *EquipmentPortDeleteOne) Exec(ctx context.Context) error {
	n, err := epdo.epd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{equipmentport.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (epdo *EquipmentPortDeleteOne) ExecX(ctx context.Context) {
	epdo.epd.ExecX(ctx)
}
