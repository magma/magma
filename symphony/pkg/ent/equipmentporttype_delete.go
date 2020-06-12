// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
)

// EquipmentPortTypeDelete is the builder for deleting a EquipmentPortType entity.
type EquipmentPortTypeDelete struct {
	config
	hooks      []Hook
	mutation   *EquipmentPortTypeMutation
	predicates []predicate.EquipmentPortType
}

// Where adds a new predicate to the delete builder.
func (eptd *EquipmentPortTypeDelete) Where(ps ...predicate.EquipmentPortType) *EquipmentPortTypeDelete {
	eptd.predicates = append(eptd.predicates, ps...)
	return eptd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (eptd *EquipmentPortTypeDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(eptd.hooks) == 0 {
		affected, err = eptd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentPortTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			eptd.mutation = mutation
			affected, err = eptd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(eptd.hooks) - 1; i >= 0; i-- {
			mut = eptd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, eptd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: equipmentporttype.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentporttype.FieldID,
			},
		},
	}
	if ps := eptd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, eptd.driver, _spec)
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
		return &NotFoundError{equipmentporttype.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (eptdo *EquipmentPortTypeDeleteOne) ExecX(ctx context.Context) {
	eptdo.eptd.ExecX(ctx)
}
