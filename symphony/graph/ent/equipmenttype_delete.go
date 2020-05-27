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
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// EquipmentTypeDelete is the builder for deleting a EquipmentType entity.
type EquipmentTypeDelete struct {
	config
	hooks      []Hook
	mutation   *EquipmentTypeMutation
	predicates []predicate.EquipmentType
}

// Where adds a new predicate to the delete builder.
func (etd *EquipmentTypeDelete) Where(ps ...predicate.EquipmentType) *EquipmentTypeDelete {
	etd.predicates = append(etd.predicates, ps...)
	return etd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (etd *EquipmentTypeDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(etd.hooks) == 0 {
		affected, err = etd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			etd.mutation = mutation
			affected, err = etd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(etd.hooks) - 1; i >= 0; i-- {
			mut = etd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, etd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
				Type:   field.TypeInt,
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
		return &NotFoundError{equipmenttype.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (etdo *EquipmentTypeDeleteOne) ExecX(ctx context.Context) {
	etdo.etd.ExecX(ctx)
}
