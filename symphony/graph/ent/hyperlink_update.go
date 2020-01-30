// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// HyperlinkUpdate is the builder for updating Hyperlink entities.
type HyperlinkUpdate struct {
	config

	update_time   *time.Time
	url           *string
	name          *string
	clearname     bool
	category      *string
	clearcategory bool
	predicates    []predicate.Hyperlink
}

// Where adds a new predicate for the builder.
func (hu *HyperlinkUpdate) Where(ps ...predicate.Hyperlink) *HyperlinkUpdate {
	hu.predicates = append(hu.predicates, ps...)
	return hu
}

// SetURL sets the url field.
func (hu *HyperlinkUpdate) SetURL(s string) *HyperlinkUpdate {
	hu.url = &s
	return hu
}

// SetName sets the name field.
func (hu *HyperlinkUpdate) SetName(s string) *HyperlinkUpdate {
	hu.name = &s
	return hu
}

// SetNillableName sets the name field if the given value is not nil.
func (hu *HyperlinkUpdate) SetNillableName(s *string) *HyperlinkUpdate {
	if s != nil {
		hu.SetName(*s)
	}
	return hu
}

// ClearName clears the value of name.
func (hu *HyperlinkUpdate) ClearName() *HyperlinkUpdate {
	hu.name = nil
	hu.clearname = true
	return hu
}

// SetCategory sets the category field.
func (hu *HyperlinkUpdate) SetCategory(s string) *HyperlinkUpdate {
	hu.category = &s
	return hu
}

// SetNillableCategory sets the category field if the given value is not nil.
func (hu *HyperlinkUpdate) SetNillableCategory(s *string) *HyperlinkUpdate {
	if s != nil {
		hu.SetCategory(*s)
	}
	return hu
}

// ClearCategory clears the value of category.
func (hu *HyperlinkUpdate) ClearCategory() *HyperlinkUpdate {
	hu.category = nil
	hu.clearcategory = true
	return hu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (hu *HyperlinkUpdate) Save(ctx context.Context) (int, error) {
	if hu.update_time == nil {
		v := hyperlink.UpdateDefaultUpdateTime()
		hu.update_time = &v
	}
	return hu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (hu *HyperlinkUpdate) SaveX(ctx context.Context) int {
	affected, err := hu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (hu *HyperlinkUpdate) Exec(ctx context.Context) error {
	_, err := hu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (hu *HyperlinkUpdate) ExecX(ctx context.Context) {
	if err := hu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (hu *HyperlinkUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   hyperlink.Table,
			Columns: hyperlink.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: hyperlink.FieldID,
			},
		},
	}
	if ps := hu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := hu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: hyperlink.FieldUpdateTime,
		})
	}
	if value := hu.url; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: hyperlink.FieldURL,
		})
	}
	if value := hu.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: hyperlink.FieldName,
		})
	}
	if hu.clearname {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: hyperlink.FieldName,
		})
	}
	if value := hu.category; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: hyperlink.FieldCategory,
		})
	}
	if hu.clearcategory {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: hyperlink.FieldCategory,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, hu.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// HyperlinkUpdateOne is the builder for updating a single Hyperlink entity.
type HyperlinkUpdateOne struct {
	config
	id string

	update_time   *time.Time
	url           *string
	name          *string
	clearname     bool
	category      *string
	clearcategory bool
}

// SetURL sets the url field.
func (huo *HyperlinkUpdateOne) SetURL(s string) *HyperlinkUpdateOne {
	huo.url = &s
	return huo
}

// SetName sets the name field.
func (huo *HyperlinkUpdateOne) SetName(s string) *HyperlinkUpdateOne {
	huo.name = &s
	return huo
}

// SetNillableName sets the name field if the given value is not nil.
func (huo *HyperlinkUpdateOne) SetNillableName(s *string) *HyperlinkUpdateOne {
	if s != nil {
		huo.SetName(*s)
	}
	return huo
}

// ClearName clears the value of name.
func (huo *HyperlinkUpdateOne) ClearName() *HyperlinkUpdateOne {
	huo.name = nil
	huo.clearname = true
	return huo
}

// SetCategory sets the category field.
func (huo *HyperlinkUpdateOne) SetCategory(s string) *HyperlinkUpdateOne {
	huo.category = &s
	return huo
}

// SetNillableCategory sets the category field if the given value is not nil.
func (huo *HyperlinkUpdateOne) SetNillableCategory(s *string) *HyperlinkUpdateOne {
	if s != nil {
		huo.SetCategory(*s)
	}
	return huo
}

// ClearCategory clears the value of category.
func (huo *HyperlinkUpdateOne) ClearCategory() *HyperlinkUpdateOne {
	huo.category = nil
	huo.clearcategory = true
	return huo
}

// Save executes the query and returns the updated entity.
func (huo *HyperlinkUpdateOne) Save(ctx context.Context) (*Hyperlink, error) {
	if huo.update_time == nil {
		v := hyperlink.UpdateDefaultUpdateTime()
		huo.update_time = &v
	}
	return huo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (huo *HyperlinkUpdateOne) SaveX(ctx context.Context) *Hyperlink {
	h, err := huo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return h
}

// Exec executes the query on the entity.
func (huo *HyperlinkUpdateOne) Exec(ctx context.Context) error {
	_, err := huo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (huo *HyperlinkUpdateOne) ExecX(ctx context.Context) {
	if err := huo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (huo *HyperlinkUpdateOne) sqlSave(ctx context.Context) (h *Hyperlink, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   hyperlink.Table,
			Columns: hyperlink.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  huo.id,
				Type:   field.TypeString,
				Column: hyperlink.FieldID,
			},
		},
	}
	if value := huo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: hyperlink.FieldUpdateTime,
		})
	}
	if value := huo.url; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: hyperlink.FieldURL,
		})
	}
	if value := huo.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: hyperlink.FieldName,
		})
	}
	if huo.clearname {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: hyperlink.FieldName,
		})
	}
	if value := huo.category; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: hyperlink.FieldCategory,
		})
	}
	if huo.clearcategory {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: hyperlink.FieldCategory,
		})
	}
	h = &Hyperlink{config: huo.config}
	_spec.Assign = h.assignValues
	_spec.ScanValues = h.scanValues()
	if err = sqlgraph.UpdateNode(ctx, huo.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return h, nil
}
