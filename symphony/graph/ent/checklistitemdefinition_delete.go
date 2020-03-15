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
	"github.com/facebookincubator/symphony/graph/ent/checklistitemdefinition"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// CheckListItemDefinitionDelete is the builder for deleting a CheckListItemDefinition entity.
type CheckListItemDefinitionDelete struct {
	config
	predicates []predicate.CheckListItemDefinition
}

// Where adds a new predicate to the delete builder.
func (clidd *CheckListItemDefinitionDelete) Where(ps ...predicate.CheckListItemDefinition) *CheckListItemDefinitionDelete {
	clidd.predicates = append(clidd.predicates, ps...)
	return clidd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (clidd *CheckListItemDefinitionDelete) Exec(ctx context.Context) (int, error) {
	return clidd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (clidd *CheckListItemDefinitionDelete) ExecX(ctx context.Context) int {
	n, err := clidd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (clidd *CheckListItemDefinitionDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: checklistitemdefinition.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistitemdefinition.FieldID,
			},
		},
	}
	if ps := clidd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, clidd.driver, _spec)
}

// CheckListItemDefinitionDeleteOne is the builder for deleting a single CheckListItemDefinition entity.
type CheckListItemDefinitionDeleteOne struct {
	clidd *CheckListItemDefinitionDelete
}

// Exec executes the deletion query.
func (cliddo *CheckListItemDefinitionDeleteOne) Exec(ctx context.Context) error {
	n, err := cliddo.clidd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{checklistitemdefinition.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (cliddo *CheckListItemDefinitionDeleteOne) ExecX(ctx context.Context) {
	cliddo.clidd.ExecX(ctx)
}
