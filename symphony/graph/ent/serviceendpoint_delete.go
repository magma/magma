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
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
)

// ServiceEndpointDelete is the builder for deleting a ServiceEndpoint entity.
type ServiceEndpointDelete struct {
	config
	predicates []predicate.ServiceEndpoint
}

// Where adds a new predicate to the delete builder.
func (sed *ServiceEndpointDelete) Where(ps ...predicate.ServiceEndpoint) *ServiceEndpointDelete {
	sed.predicates = append(sed.predicates, ps...)
	return sed
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (sed *ServiceEndpointDelete) Exec(ctx context.Context) (int, error) {
	return sed.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (sed *ServiceEndpointDelete) ExecX(ctx context.Context) int {
	n, err := sed.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (sed *ServiceEndpointDelete) sqlExec(ctx context.Context) (int, error) {
	spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: serviceendpoint.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: serviceendpoint.FieldID,
			},
		},
	}
	if ps := sed.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, sed.driver, spec)
}

// ServiceEndpointDeleteOne is the builder for deleting a single ServiceEndpoint entity.
type ServiceEndpointDeleteOne struct {
	sed *ServiceEndpointDelete
}

// Exec executes the deletion query.
func (sedo *ServiceEndpointDeleteOne) Exec(ctx context.Context) error {
	n, err := sedo.sed.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{serviceendpoint.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (sedo *ServiceEndpointDeleteOne) ExecX(ctx context.Context) {
	sedo.sed.ExecX(ctx)
}
