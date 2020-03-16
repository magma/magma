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
	"github.com/facebookincubator/symphony/graph/ent/reportfilter"
)

// ReportFilterCreate is the builder for creating a ReportFilter entity.
type ReportFilterCreate struct {
	config
	create_time *time.Time
	update_time *time.Time
	name        *string
	entity      *reportfilter.Entity
	filters     *string
}

// SetCreateTime sets the create_time field.
func (rfc *ReportFilterCreate) SetCreateTime(t time.Time) *ReportFilterCreate {
	rfc.create_time = &t
	return rfc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (rfc *ReportFilterCreate) SetNillableCreateTime(t *time.Time) *ReportFilterCreate {
	if t != nil {
		rfc.SetCreateTime(*t)
	}
	return rfc
}

// SetUpdateTime sets the update_time field.
func (rfc *ReportFilterCreate) SetUpdateTime(t time.Time) *ReportFilterCreate {
	rfc.update_time = &t
	return rfc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (rfc *ReportFilterCreate) SetNillableUpdateTime(t *time.Time) *ReportFilterCreate {
	if t != nil {
		rfc.SetUpdateTime(*t)
	}
	return rfc
}

// SetName sets the name field.
func (rfc *ReportFilterCreate) SetName(s string) *ReportFilterCreate {
	rfc.name = &s
	return rfc
}

// SetEntity sets the entity field.
func (rfc *ReportFilterCreate) SetEntity(r reportfilter.Entity) *ReportFilterCreate {
	rfc.entity = &r
	return rfc
}

// SetFilters sets the filters field.
func (rfc *ReportFilterCreate) SetFilters(s string) *ReportFilterCreate {
	rfc.filters = &s
	return rfc
}

// SetNillableFilters sets the filters field if the given value is not nil.
func (rfc *ReportFilterCreate) SetNillableFilters(s *string) *ReportFilterCreate {
	if s != nil {
		rfc.SetFilters(*s)
	}
	return rfc
}

// Save creates the ReportFilter in the database.
func (rfc *ReportFilterCreate) Save(ctx context.Context) (*ReportFilter, error) {
	if rfc.create_time == nil {
		v := reportfilter.DefaultCreateTime()
		rfc.create_time = &v
	}
	if rfc.update_time == nil {
		v := reportfilter.DefaultUpdateTime()
		rfc.update_time = &v
	}
	if rfc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if err := reportfilter.NameValidator(*rfc.name); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
	}
	if rfc.entity == nil {
		return nil, errors.New("ent: missing required field \"entity\"")
	}
	if err := reportfilter.EntityValidator(*rfc.entity); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"entity\": %v", err)
	}
	if rfc.filters == nil {
		v := reportfilter.DefaultFilters
		rfc.filters = &v
	}
	return rfc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (rfc *ReportFilterCreate) SaveX(ctx context.Context) *ReportFilter {
	v, err := rfc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (rfc *ReportFilterCreate) sqlSave(ctx context.Context) (*ReportFilter, error) {
	var (
		rf    = &ReportFilter{config: rfc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: reportfilter.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: reportfilter.FieldID,
			},
		}
	)
	if value := rfc.create_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: reportfilter.FieldCreateTime,
		})
		rf.CreateTime = *value
	}
	if value := rfc.update_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: reportfilter.FieldUpdateTime,
		})
		rf.UpdateTime = *value
	}
	if value := rfc.name; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: reportfilter.FieldName,
		})
		rf.Name = *value
	}
	if value := rfc.entity; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  *value,
			Column: reportfilter.FieldEntity,
		})
		rf.Entity = *value
	}
	if value := rfc.filters; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: reportfilter.FieldFilters,
		})
		rf.Filters = *value
	}
	if err := sqlgraph.CreateNode(ctx, rfc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	rf.ID = int(id)
	return rf, nil
}
