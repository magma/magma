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
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// EquipmentPositionDefinitionDelete is the builder for deleting a EquipmentPositionDefinition entity.
type EquipmentPositionDefinitionDelete struct {
	config
	predicates []predicate.EquipmentPositionDefinition
}

// Where adds a new predicate to the delete builder.
func (epdd *EquipmentPositionDefinitionDelete) Where(ps ...predicate.EquipmentPositionDefinition) *EquipmentPositionDefinitionDelete {
	epdd.predicates = append(epdd.predicates, ps...)
	return epdd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (epdd *EquipmentPositionDefinitionDelete) Exec(ctx context.Context) (int, error) {
	return epdd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (epdd *EquipmentPositionDefinitionDelete) ExecX(ctx context.Context) int {
	n, err := epdd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (epdd *EquipmentPositionDefinitionDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: equipmentpositiondefinition.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: equipmentpositiondefinition.FieldID,
			},
		},
	}
	if ps := epdd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, epdd.driver, _spec)
}

// EquipmentPositionDefinitionDeleteOne is the builder for deleting a single EquipmentPositionDefinition entity.
type EquipmentPositionDefinitionDeleteOne struct {
	epdd *EquipmentPositionDefinitionDelete
}

// Exec executes the deletion query.
func (epddo *EquipmentPositionDefinitionDeleteOne) Exec(ctx context.Context) error {
	n, err := epddo.epdd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{equipmentpositiondefinition.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (epddo *EquipmentPositionDefinitionDeleteOne) ExecX(ctx context.Context) {
	epddo.epdd.ExecX(ctx)
}
