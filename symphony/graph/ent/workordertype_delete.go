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
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// WorkOrderTypeDelete is the builder for deleting a WorkOrderType entity.
type WorkOrderTypeDelete struct {
	config
	predicates []predicate.WorkOrderType
}

// Where adds a new predicate to the delete builder.
func (wotd *WorkOrderTypeDelete) Where(ps ...predicate.WorkOrderType) *WorkOrderTypeDelete {
	wotd.predicates = append(wotd.predicates, ps...)
	return wotd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (wotd *WorkOrderTypeDelete) Exec(ctx context.Context) (int, error) {
	return wotd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (wotd *WorkOrderTypeDelete) ExecX(ctx context.Context) int {
	n, err := wotd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (wotd *WorkOrderTypeDelete) sqlExec(ctx context.Context) (int, error) {
	spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: workordertype.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: workordertype.FieldID,
			},
		},
	}
	if ps := wotd.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, wotd.driver, spec)
}

// WorkOrderTypeDeleteOne is the builder for deleting a single WorkOrderType entity.
type WorkOrderTypeDeleteOne struct {
	wotd *WorkOrderTypeDelete
}

// Exec executes the deletion query.
func (wotdo *WorkOrderTypeDeleteOne) Exec(ctx context.Context) error {
	n, err := wotdo.wotd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{workordertype.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (wotdo *WorkOrderTypeDeleteOne) ExecX(ctx context.Context) {
	wotdo.wotd.ExecX(ctx)
}
