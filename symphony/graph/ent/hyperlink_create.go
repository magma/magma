// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
)

// HyperlinkCreate is the builder for creating a Hyperlink entity.
type HyperlinkCreate struct {
	config
	create_time *time.Time
	update_time *time.Time
	url         *string
	name        *string
	category    *string
}

// SetCreateTime sets the create_time field.
func (hc *HyperlinkCreate) SetCreateTime(t time.Time) *HyperlinkCreate {
	hc.create_time = &t
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
	hc.update_time = &t
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
	hc.url = &s
	return hc
}

// SetName sets the name field.
func (hc *HyperlinkCreate) SetName(s string) *HyperlinkCreate {
	hc.name = &s
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
	hc.category = &s
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
	if hc.create_time == nil {
		v := hyperlink.DefaultCreateTime()
		hc.create_time = &v
	}
	if hc.update_time == nil {
		v := hyperlink.DefaultUpdateTime()
		hc.update_time = &v
	}
	if hc.url == nil {
		return nil, errors.New("ent: missing required field \"url\"")
	}
	return hc.sqlSave(ctx)
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
		h    = &Hyperlink{config: hc.config}
		spec = &sqlgraph.CreateSpec{
			Table: hyperlink.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: hyperlink.FieldID,
			},
		}
	)
	if value := hc.create_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: hyperlink.FieldCreateTime,
		})
		h.CreateTime = *value
	}
	if value := hc.update_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: hyperlink.FieldUpdateTime,
		})
		h.UpdateTime = *value
	}
	if value := hc.url; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: hyperlink.FieldURL,
		})
		h.URL = *value
	}
	if value := hc.name; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: hyperlink.FieldName,
		})
		h.Name = *value
	}
	if value := hc.category; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: hyperlink.FieldCategory,
		})
		h.Category = *value
	}
	if err := sqlgraph.CreateNode(ctx, hc.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := spec.ID.Value.(int64)
	h.ID = strconv.FormatInt(id, 10)
	return h, nil
}
