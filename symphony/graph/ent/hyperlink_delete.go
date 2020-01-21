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
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// HyperlinkDelete is the builder for deleting a Hyperlink entity.
type HyperlinkDelete struct {
	config
	predicates []predicate.Hyperlink
}

// Where adds a new predicate to the delete builder.
func (hd *HyperlinkDelete) Where(ps ...predicate.Hyperlink) *HyperlinkDelete {
	hd.predicates = append(hd.predicates, ps...)
	return hd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (hd *HyperlinkDelete) Exec(ctx context.Context) (int, error) {
	return hd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (hd *HyperlinkDelete) ExecX(ctx context.Context) int {
	n, err := hd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (hd *HyperlinkDelete) sqlExec(ctx context.Context) (int, error) {
	spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: hyperlink.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: hyperlink.FieldID,
			},
		},
	}
	if ps := hd.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, hd.driver, spec)
}

// HyperlinkDeleteOne is the builder for deleting a single Hyperlink entity.
type HyperlinkDeleteOne struct {
	hd *HyperlinkDelete
}

// Exec executes the deletion query.
func (hdo *HyperlinkDeleteOne) Exec(ctx context.Context) error {
	n, err := hdo.hd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{hyperlink.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (hdo *HyperlinkDeleteOne) ExecX(ctx context.Context) {
	hdo.hd.ExecX(ctx)
}
