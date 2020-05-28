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
	"github.com/facebookincubator/symphony/pkg/ent/equipmentport"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
)

// EquipmentPortDelete is the builder for deleting a EquipmentPort entity.
type EquipmentPortDelete struct {
	config
	hooks      []Hook
	mutation   *EquipmentPortMutation
	predicates []predicate.EquipmentPort
}

// Where adds a new predicate to the delete builder.
func (epd *EquipmentPortDelete) Where(ps ...predicate.EquipmentPort) *EquipmentPortDelete {
	epd.predicates = append(epd.predicates, ps...)
	return epd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (epd *EquipmentPortDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(epd.hooks) == 0 {
		affected, err = epd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentPortMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			epd.mutation = mutation
			affected, err = epd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(epd.hooks) - 1; i >= 0; i-- {
			mut = epd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, epd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: equipmentport.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentport.FieldID,
			},
		},
	}
	if ps := epd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, epd.driver, _spec)
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
		return &NotFoundError{equipmentport.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (epdo *EquipmentPortDeleteOne) ExecX(ctx context.Context) {
	epdo.epd.ExecX(ctx)
}
