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
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/workordertemplate"
)

// WorkOrderTemplateDelete is the builder for deleting a WorkOrderTemplate entity.
type WorkOrderTemplateDelete struct {
	config
	hooks      []Hook
	mutation   *WorkOrderTemplateMutation
	predicates []predicate.WorkOrderTemplate
}

// Where adds a new predicate to the delete builder.
func (wotd *WorkOrderTemplateDelete) Where(ps ...predicate.WorkOrderTemplate) *WorkOrderTemplateDelete {
	wotd.predicates = append(wotd.predicates, ps...)
	return wotd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (wotd *WorkOrderTemplateDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(wotd.hooks) == 0 {
		affected, err = wotd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*WorkOrderTemplateMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			wotd.mutation = mutation
			affected, err = wotd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(wotd.hooks) - 1; i >= 0; i-- {
			mut = wotd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, wotd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (wotd *WorkOrderTemplateDelete) ExecX(ctx context.Context) int {
	n, err := wotd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (wotd *WorkOrderTemplateDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: workordertemplate.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: workordertemplate.FieldID,
			},
		},
	}
	if ps := wotd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, wotd.driver, _spec)
}

// WorkOrderTemplateDeleteOne is the builder for deleting a single WorkOrderTemplate entity.
type WorkOrderTemplateDeleteOne struct {
	wotd *WorkOrderTemplateDelete
}

// Exec executes the deletion query.
func (wotdo *WorkOrderTemplateDeleteOne) Exec(ctx context.Context) error {
	n, err := wotdo.wotd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{workordertemplate.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (wotdo *WorkOrderTemplateDeleteOne) ExecX(ctx context.Context) {
	wotdo.wotd.ExecX(ctx)
}
