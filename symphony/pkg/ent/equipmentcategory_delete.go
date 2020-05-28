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
	"github.com/facebookincubator/symphony/pkg/ent/equipmentcategory"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
)

// EquipmentCategoryDelete is the builder for deleting a EquipmentCategory entity.
type EquipmentCategoryDelete struct {
	config
	hooks      []Hook
	mutation   *EquipmentCategoryMutation
	predicates []predicate.EquipmentCategory
}

// Where adds a new predicate to the delete builder.
func (ecd *EquipmentCategoryDelete) Where(ps ...predicate.EquipmentCategory) *EquipmentCategoryDelete {
	ecd.predicates = append(ecd.predicates, ps...)
	return ecd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ecd *EquipmentCategoryDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(ecd.hooks) == 0 {
		affected, err = ecd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentCategoryMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ecd.mutation = mutation
			affected, err = ecd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(ecd.hooks) - 1; i >= 0; i-- {
			mut = ecd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ecd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (ecd *EquipmentCategoryDelete) ExecX(ctx context.Context) int {
	n, err := ecd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ecd *EquipmentCategoryDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: equipmentcategory.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentcategory.FieldID,
			},
		},
	}
	if ps := ecd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, ecd.driver, _spec)
}

// EquipmentCategoryDeleteOne is the builder for deleting a single EquipmentCategory entity.
type EquipmentCategoryDeleteOne struct {
	ecd *EquipmentCategoryDelete
}

// Exec executes the deletion query.
func (ecdo *EquipmentCategoryDeleteOne) Exec(ctx context.Context) error {
	n, err := ecdo.ecd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{equipmentcategory.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ecdo *EquipmentCategoryDeleteOne) ExecX(ctx context.Context) {
	ecdo.ecd.ExecX(ctx)
}
