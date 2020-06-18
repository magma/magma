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
	"github.com/facebookincubator/symphony/pkg/ent/reportfilter"
)

// ReportFilterUpdate is the builder for updating ReportFilter entities.
type ReportFilterUpdate struct {
	config
	hooks      []Hook
	mutation   *ReportFilterMutation
	predicates []predicate.ReportFilter
}

// Where adds a new predicate for the builder.
func (rfu *ReportFilterUpdate) Where(ps ...predicate.ReportFilter) *ReportFilterUpdate {
	rfu.predicates = append(rfu.predicates, ps...)
	return rfu
}

// SetName sets the name field.
func (rfu *ReportFilterUpdate) SetName(s string) *ReportFilterUpdate {
	rfu.mutation.SetName(s)
	return rfu
}

// SetEntity sets the entity field.
func (rfu *ReportFilterUpdate) SetEntity(r reportfilter.Entity) *ReportFilterUpdate {
	rfu.mutation.SetEntity(r)
	return rfu
}

// SetFilters sets the filters field.
func (rfu *ReportFilterUpdate) SetFilters(s string) *ReportFilterUpdate {
	rfu.mutation.SetFilters(s)
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
	if _, ok := rfu.mutation.UpdateTime(); !ok {
		v := reportfilter.UpdateDefaultUpdateTime()
		rfu.mutation.SetUpdateTime(v)
	}
	if v, ok := rfu.mutation.Name(); ok {
		if err := reportfilter.NameValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if v, ok := rfu.mutation.Entity(); ok {
		if err := reportfilter.EntityValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"entity\": %v", err)
		}
	}
	var (
		err      error
		affected int
	)
	if len(rfu.hooks) == 0 {
		affected, err = rfu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ReportFilterMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			rfu.mutation = mutation
			affected, err = rfu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(rfu.hooks) - 1; i >= 0; i-- {
			mut = rfu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, rfu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
	if value, ok := rfu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: reportfilter.FieldUpdateTime,
		})
	}
	if value, ok := rfu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: reportfilter.FieldName,
		})
	}
	if value, ok := rfu.mutation.Entity(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: reportfilter.FieldEntity,
		})
	}
	if value, ok := rfu.mutation.Filters(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
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
	hooks    []Hook
	mutation *ReportFilterMutation
}

// SetName sets the name field.
func (rfuo *ReportFilterUpdateOne) SetName(s string) *ReportFilterUpdateOne {
	rfuo.mutation.SetName(s)
	return rfuo
}

// SetEntity sets the entity field.
func (rfuo *ReportFilterUpdateOne) SetEntity(r reportfilter.Entity) *ReportFilterUpdateOne {
	rfuo.mutation.SetEntity(r)
	return rfuo
}

// SetFilters sets the filters field.
func (rfuo *ReportFilterUpdateOne) SetFilters(s string) *ReportFilterUpdateOne {
	rfuo.mutation.SetFilters(s)
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
	if _, ok := rfuo.mutation.UpdateTime(); !ok {
		v := reportfilter.UpdateDefaultUpdateTime()
		rfuo.mutation.SetUpdateTime(v)
	}
	if v, ok := rfuo.mutation.Name(); ok {
		if err := reportfilter.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if v, ok := rfuo.mutation.Entity(); ok {
		if err := reportfilter.EntityValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"entity\": %v", err)
		}
	}
	var (
		err  error
		node *ReportFilter
	)
	if len(rfuo.hooks) == 0 {
		node, err = rfuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ReportFilterMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			rfuo.mutation = mutation
			node, err = rfuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(rfuo.hooks) - 1; i >= 0; i-- {
			mut = rfuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, rfuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
				Type:   field.TypeInt,
				Column: reportfilter.FieldID,
			},
		},
	}
	id, ok := rfuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing ReportFilter.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := rfuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: reportfilter.FieldUpdateTime,
		})
	}
	if value, ok := rfuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: reportfilter.FieldName,
		})
	}
	if value, ok := rfuo.mutation.Entity(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: reportfilter.FieldEntity,
		})
	}
	if value, ok := rfuo.mutation.Filters(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
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
