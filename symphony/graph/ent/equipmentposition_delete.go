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
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// EquipmentPositionDelete is the builder for deleting a EquipmentPosition entity.
type EquipmentPositionDelete struct {
	config
	predicates []predicate.EquipmentPosition
}

// Where adds a new predicate to the delete builder.
func (epd *EquipmentPositionDelete) Where(ps ...predicate.EquipmentPosition) *EquipmentPositionDelete {
	epd.predicates = append(epd.predicates, ps...)
	return epd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (epd *EquipmentPositionDelete) Exec(ctx context.Context) (int, error) {
	return epd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (epd *EquipmentPositionDelete) ExecX(ctx context.Context) int {
	n, err := epd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (epd *EquipmentPositionDelete) sqlExec(ctx context.Context) (int, error) {
	spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: equipmentposition.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: equipmentposition.FieldID,
			},
		},
	}
	if ps := epd.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, epd.driver, spec)
}

// EquipmentPositionDeleteOne is the builder for deleting a single EquipmentPosition entity.
type EquipmentPositionDeleteOne struct {
	epd *EquipmentPositionDelete
}

// Exec executes the deletion query.
func (epdo *EquipmentPositionDeleteOne) Exec(ctx context.Context) error {
	n, err := epdo.epd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{equipmentposition.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (epdo *EquipmentPositionDeleteOne) ExecX(ctx context.Context) {
	epdo.epd.ExecX(ctx)
}
