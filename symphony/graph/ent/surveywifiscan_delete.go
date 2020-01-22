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
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
)

// SurveyWiFiScanDelete is the builder for deleting a SurveyWiFiScan entity.
type SurveyWiFiScanDelete struct {
	config
	predicates []predicate.SurveyWiFiScan
}

// Where adds a new predicate to the delete builder.
func (swfsd *SurveyWiFiScanDelete) Where(ps ...predicate.SurveyWiFiScan) *SurveyWiFiScanDelete {
	swfsd.predicates = append(swfsd.predicates, ps...)
	return swfsd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (swfsd *SurveyWiFiScanDelete) Exec(ctx context.Context) (int, error) {
	return swfsd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (swfsd *SurveyWiFiScanDelete) ExecX(ctx context.Context) int {
	n, err := swfsd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (swfsd *SurveyWiFiScanDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: surveywifiscan.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: surveywifiscan.FieldID,
			},
		},
	}
	if ps := swfsd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, swfsd.driver, _spec)
}

// SurveyWiFiScanDeleteOne is the builder for deleting a single SurveyWiFiScan entity.
type SurveyWiFiScanDeleteOne struct {
	swfsd *SurveyWiFiScanDelete
}

// Exec executes the deletion query.
func (swfsdo *SurveyWiFiScanDeleteOne) Exec(ctx context.Context) error {
	n, err := swfsdo.swfsd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{surveywifiscan.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (swfsdo *SurveyWiFiScanDeleteOne) ExecX(ctx context.Context) {
	swfsdo.swfsd.ExecX(ctx)
}
