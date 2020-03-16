// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
)

// HyperlinkCreate is the builder for creating a Hyperlink entity.
type HyperlinkCreate struct {
	config
	mutation *HyperlinkMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (hc *HyperlinkCreate) SetCreateTime(t time.Time) *HyperlinkCreate {
	hc.mutation.SetCreateTime(t)
	return hc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (hc *HyperlinkCreate) SetNillableCreateTime(t *time.Time) *HyperlinkCreate {
	if t != nil {
		hc.SetCreateTime(*t)
	}
	return hc
}

// SetUpdateTime sets the update_time field.
func (hc *HyperlinkCreate) SetUpdateTime(t time.Time) *HyperlinkCreate {
	hc.mutation.SetUpdateTime(t)
	return hc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (hc *HyperlinkCreate) SetNillableUpdateTime(t *time.Time) *HyperlinkCreate {
	if t != nil {
		hc.SetUpdateTime(*t)
	}
	return hc
}

// SetURL sets the url field.
func (hc *HyperlinkCreate) SetURL(s string) *HyperlinkCreate {
	hc.mutation.SetURL(s)
	return hc
}

// SetName sets the name field.
func (hc *HyperlinkCreate) SetName(s string) *HyperlinkCreate {
	hc.mutation.SetName(s)
	return hc
}

// SetNillableName sets the name field if the given value is not nil.
func (hc *HyperlinkCreate) SetNillableName(s *string) *HyperlinkCreate {
	if s != nil {
		hc.SetName(*s)
	}
	return hc
}

// SetCategory sets the category field.
func (hc *HyperlinkCreate) SetCategory(s string) *HyperlinkCreate {
	hc.mutation.SetCategory(s)
	return hc
}

// SetNillableCategory sets the category field if the given value is not nil.
func (hc *HyperlinkCreate) SetNillableCategory(s *string) *HyperlinkCreate {
	if s != nil {
		hc.SetCategory(*s)
	}
	return hc
}

// Save creates the Hyperlink in the database.
func (hc *HyperlinkCreate) Save(ctx context.Context) (*Hyperlink, error) {
	if _, ok := hc.mutation.CreateTime(); !ok {
		v := hyperlink.DefaultCreateTime()
		hc.mutation.SetCreateTime(v)
	}
	if _, ok := hc.mutation.UpdateTime(); !ok {
		v := hyperlink.DefaultUpdateTime()
		hc.mutation.SetUpdateTime(v)
	}
	if _, ok := hc.mutation.URL(); !ok {
		return nil, errors.New("ent: missing required field \"url\"")
	}
	var (
		err  error
		node *Hyperlink
	)
	if len(hc.hooks) == 0 {
		node, err = hc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*HyperlinkMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			hc.mutation = mutation
			node, err = hc.sqlSave(ctx)
			return node, err
		})
		for i := len(hc.hooks); i > 0; i-- {
			mut = hc.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, hc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (hc *HyperlinkCreate) SaveX(ctx context.Context) *Hyperlink {
	v, err := hc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (hc *HyperlinkCreate) sqlSave(ctx context.Context) (*Hyperlink, error) {
	var (
		h     = &Hyperlink{config: hc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: hyperlink.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: hyperlink.FieldID,
			},
		}
	)
	if value, ok := hc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: hyperlink.FieldCreateTime,
		})
		h.CreateTime = value
	}
	if value, ok := hc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: hyperlink.FieldUpdateTime,
		})
		h.UpdateTime = value
	}
	if value, ok := hc.mutation.URL(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: hyperlink.FieldURL,
		})
		h.URL = value
	}
	if value, ok := hc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: hyperlink.FieldName,
		})
		h.Name = value
	}
	if value, ok := hc.mutation.Category(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: hyperlink.FieldCategory,
		})
		h.Category = value
	}
	if err := sqlgraph.CreateNode(ctx, hc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	h.ID = int(id)
	return h, nil
}
