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
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// EquipmentPortTypeDelete is the builder for deleting a EquipmentPortType entity.
type EquipmentPortTypeDelete struct {
	config
	predicates []predicate.EquipmentPortType
}

// Where adds a new predicate to the delete builder.
func (eptd *EquipmentPortTypeDelete) Where(ps ...predicate.EquipmentPortType) *EquipmentPortTypeDelete {
	eptd.predicates = append(eptd.predicates, ps...)
	return eptd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (eptd *EquipmentPortTypeDelete) Exec(ctx context.Context) (int, error) {
	return eptd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (eptd *EquipmentPortTypeDelete) ExecX(ctx context.Context) int {
	n, err := eptd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (eptd *EquipmentPortTypeDelete) sqlExec(ctx context.Context) (int, error) {
	spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: equipmentporttype.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: equipmentporttype.FieldID,
			},
		},
	}
	if ps := eptd.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, eptd.driver, spec)
}

// EquipmentPortTypeDeleteOne is the builder for deleting a single EquipmentPortType entity.
type EquipmentPortTypeDeleteOne struct {
	eptd *EquipmentPortTypeDelete
}

// Exec executes the deletion query.
func (eptdo *EquipmentPortTypeDeleteOne) Exec(ctx context.Context) error {
	n, err := eptdo.eptd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{equipmentporttype.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (eptdo *EquipmentPortTypeDeleteOne) ExecX(ctx context.Context) {
	eptdo.eptd.ExecX(ctx)
}
