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
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: equipment.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipment.FieldID,
			},
		},
	}
	if ps := ed.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, ed.driver, _spec)
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
		return &NotFoundError{equipment.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (edo *EquipmentDeleteOne) ExecX(ctx context.Context) {
	edo.ed.ExecX(ctx)
}
