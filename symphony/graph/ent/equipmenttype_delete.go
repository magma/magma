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
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// EquipmentTypeDelete is the builder for deleting a EquipmentType entity.
type EquipmentTypeDelete struct {
	config
	predicates []predicate.EquipmentType
}

// Where adds a new predicate to the delete builder.
func (etd *EquipmentTypeDelete) Where(ps ...predicate.EquipmentType) *EquipmentTypeDelete {
	etd.predicates = append(etd.predicates, ps...)
	return etd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (etd *EquipmentTypeDelete) Exec(ctx context.Context) (int, error) {
	return etd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (etd *EquipmentTypeDelete) ExecX(ctx context.Context) int {
	n, err := etd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (etd *EquipmentTypeDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: equipmenttype.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: equipmenttype.FieldID,
			},
		},
	}
	if ps := etd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, etd.driver, _spec)
}

// EquipmentTypeDeleteOne is the builder for deleting a single EquipmentType entity.
type EquipmentTypeDeleteOne struct {
	etd *EquipmentTypeDelete
}

// Exec executes the deletion query.
func (etdo *EquipmentTypeDeleteOne) Exec(ctx context.Context) error {
	n, err := etdo.etd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{equipmenttype.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (etdo *EquipmentTypeDeleteOne) ExecX(ctx context.Context) {
	etdo.etd.ExecX(ctx)
}
