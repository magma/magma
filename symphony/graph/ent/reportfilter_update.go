// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/reportfilter"
)

// ReportFilterUpdate is the builder for updating ReportFilter entities.
type ReportFilterUpdate struct {
	config

	update_time *time.Time
	name        *string
	entity      *reportfilter.Entity
	filters     *string
	predicates  []predicate.ReportFilter
}

// Where adds a new predicate for the builder.
func (rfu *ReportFilterUpdate) Where(ps ...predicate.ReportFilter) *ReportFilterUpdate {
	rfu.predicates = append(rfu.predicates, ps...)
	return rfu
}

// SetName sets the name field.
func (rfu *ReportFilterUpdate) SetName(s string) *ReportFilterUpdate {
	rfu.name = &s
	return rfu
}

// SetEntity sets the entity field.
func (rfu *ReportFilterUpdate) SetEntity(r reportfilter.Entity) *ReportFilterUpdate {
	rfu.entity = &r
	return rfu
}

// SetFilters sets the filters field.
func (rfu *ReportFilterUpdate) SetFilters(s string) *ReportFilterUpdate {
	rfu.filters = &s
	return rfu
}

// SetNillableFilters sets the filters field if the given value is not nil.
func (rfu *ReportFilterUpdate) SetNillableFilters(s *string) *ReportFilterUpdate {
	if s != nil {
		rfu.SetFilters(*s)
	}
	return rfu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (rfu *ReportFilterUpdate) Save(ctx context.Context) (int, error) {
	if rfu.update_time == nil {
		v := reportfilter.UpdateDefaultUpdateTime()
		rfu.update_time = &v
	}
	if rfu.name != nil {
		if err := reportfilter.NameValidator(*rfu.name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if rfu.entity != nil {
		if err := reportfilter.EntityValidator(*rfu.entity); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"entity\": %v", err)
		}
	}
	return rfu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (rfu *ReportFilterUpdate) SaveX(ctx context.Context) int {
	affected, err := rfu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (rfu *ReportFilterUpdate) Exec(ctx context.Context) error {
	_, err := rfu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rfu *ReportFilterUpdate) ExecX(ctx context.Context) {
	if err := rfu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (rfu *ReportFilterUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   reportfilter.Table,
			Columns: reportfilter.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: reportfilter.FieldID,
			},
		},
	}
	if ps := rfu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := rfu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: reportfilter.FieldUpdateTime,
		})
	}
	if value := rfu.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: reportfilter.FieldName,
		})
	}
	if value := rfu.entity; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  *value,
			Column: reportfilter.FieldEntity,
		})
	}
	if value := rfu.filters; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: reportfilter.FieldFilters,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, rfu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{reportfilter.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// ReportFilterUpdateOne is the builder for updating a single ReportFilter entity.
type ReportFilterUpdateOne struct {
	config
	id int

	update_time *time.Time
	name        *string
	entity      *reportfilter.Entity
	filters     *string
}

// SetName sets the name field.
func (rfuo *ReportFilterUpdateOne) SetName(s string) *ReportFilterUpdateOne {
	rfuo.name = &s
	return rfuo
}

// SetEntity sets the entity field.
func (rfuo *ReportFilterUpdateOne) SetEntity(r reportfilter.Entity) *ReportFilterUpdateOne {
	rfuo.entity = &r
	return rfuo
}

// SetFilters sets the filters field.
func (rfuo *ReportFilterUpdateOne) SetFilters(s string) *ReportFilterUpdateOne {
	rfuo.filters = &s
	return rfuo
}

// SetNillableFilters sets the filters field if the given value is not nil.
func (rfuo *ReportFilterUpdateOne) SetNillableFilters(s *string) *ReportFilterUpdateOne {
	if s != nil {
		rfuo.SetFilters(*s)
	}
	return rfuo
}

// Save executes the query and returns the updated entity.
func (rfuo *ReportFilterUpdateOne) Save(ctx context.Context) (*ReportFilter, error) {
	if rfuo.update_time == nil {
		v := reportfilter.UpdateDefaultUpdateTime()
		rfuo.update_time = &v
	}
	if rfuo.name != nil {
		if err := reportfilter.NameValidator(*rfuo.name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if rfuo.entity != nil {
		if err := reportfilter.EntityValidator(*rfuo.entity); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"entity\": %v", err)
		}
	}
	return rfuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (rfuo *ReportFilterUpdateOne) SaveX(ctx context.Context) *ReportFilter {
	rf, err := rfuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return rf
}

// Exec executes the query on the entity.
func (rfuo *ReportFilterUpdateOne) Exec(ctx context.Context) error {
	_, err := rfuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rfuo *ReportFilterUpdateOne) ExecX(ctx context.Context) {
	if err := rfuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (rfuo *ReportFilterUpdateOne) sqlSave(ctx context.Context) (rf *ReportFilter, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   reportfilter.Table,
			Columns: reportfilter.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  rfuo.id,
				Type:   field.TypeInt,
				Column: reportfilter.FieldID,
			},
		},
	}
	if value := rfuo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: reportfilter.FieldUpdateTime,
		})
	}
	if value := rfuo.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: reportfilter.FieldName,
		})
	}
	if value := rfuo.entity; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  *value,
			Column: reportfilter.FieldEntity,
		})
	}
	if value := rfuo.filters; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: reportfilter.FieldFilters,
		})
	}
	rf = &ReportFilter{config: rfuo.config}
	_spec.Assign = rf.assignValues
	_spec.ScanValues = rf.scanValues()
	if err = sqlgraph.UpdateNode(ctx, rfuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{reportfilter.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return rf, nil
}
