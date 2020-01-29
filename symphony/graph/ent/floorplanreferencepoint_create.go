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
	"github.com/facebookincubator/symphony/graph/ent/floorplanreferencepoint"
)

// FloorPlanReferencePointCreate is the builder for creating a FloorPlanReferencePoint entity.
type FloorPlanReferencePointCreate struct {
	config
	create_time *time.Time
	update_time *time.Time
	x           *int
	y           *int
	latitude    *float64
	longitude   *float64
}

// SetCreateTime sets the create_time field.
func (fprpc *FloorPlanReferencePointCreate) SetCreateTime(t time.Time) *FloorPlanReferencePointCreate {
	fprpc.create_time = &t
	return fprpc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (fprpc *FloorPlanReferencePointCreate) SetNillableCreateTime(t *time.Time) *FloorPlanReferencePointCreate {
	if t != nil {
		fprpc.SetCreateTime(*t)
	}
	return fprpc
}

// SetUpdateTime sets the update_time field.
func (fprpc *FloorPlanReferencePointCreate) SetUpdateTime(t time.Time) *FloorPlanReferencePointCreate {
	fprpc.update_time = &t
	return fprpc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (fprpc *FloorPlanReferencePointCreate) SetNillableUpdateTime(t *time.Time) *FloorPlanReferencePointCreate {
	if t != nil {
		fprpc.SetUpdateTime(*t)
	}
	return fprpc
}

// SetX sets the x field.
func (fprpc *FloorPlanReferencePointCreate) SetX(i int) *FloorPlanReferencePointCreate {
	fprpc.x = &i
	return fprpc
}

// SetY sets the y field.
func (fprpc *FloorPlanReferencePointCreate) SetY(i int) *FloorPlanReferencePointCreate {
	fprpc.y = &i
	return fprpc
}

// SetLatitude sets the latitude field.
func (fprpc *FloorPlanReferencePointCreate) SetLatitude(f float64) *FloorPlanReferencePointCreate {
	fprpc.latitude = &f
	return fprpc
}

// SetLongitude sets the longitude field.
func (fprpc *FloorPlanReferencePointCreate) SetLongitude(f float64) *FloorPlanReferencePointCreate {
	fprpc.longitude = &f
	return fprpc
}

// Save creates the FloorPlanReferencePoint in the database.
func (fprpc *FloorPlanReferencePointCreate) Save(ctx context.Context) (*FloorPlanReferencePoint, error) {
	if fprpc.create_time == nil {
		v := floorplanreferencepoint.DefaultCreateTime()
		fprpc.create_time = &v
	}
	if fprpc.update_time == nil {
		v := floorplanreferencepoint.DefaultUpdateTime()
		fprpc.update_time = &v
	}
	if fprpc.x == nil {
		return nil, errors.New("ent: missing required field \"x\"")
	}
	if fprpc.y == nil {
		return nil, errors.New("ent: missing required field \"y\"")
	}
	if fprpc.latitude == nil {
		return nil, errors.New("ent: missing required field \"latitude\"")
	}
	if fprpc.longitude == nil {
		return nil, errors.New("ent: missing required field \"longitude\"")
	}
	return fprpc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (fprpc *FloorPlanReferencePointCreate) SaveX(ctx context.Context) *FloorPlanReferencePoint {
	v, err := fprpc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (fprpc *FloorPlanReferencePointCreate) sqlSave(ctx context.Context) (*FloorPlanReferencePoint, error) {
	var (
		fprp  = &FloorPlanReferencePoint{config: fprpc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: floorplanreferencepoint.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: floorplanreferencepoint.FieldID,
			},
		}
	)
	if value := fprpc.create_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: floorplanreferencepoint.FieldCreateTime,
		})
		fprp.CreateTime = *value
	}
	if value := fprpc.update_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: floorplanreferencepoint.FieldUpdateTime,
		})
		fprp.UpdateTime = *value
	}
	if value := fprpc.x; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanreferencepoint.FieldX,
		})
		fprp.X = *value
	}
	if value := fprpc.y; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanreferencepoint.FieldY,
		})
		fprp.Y = *value
	}
	if value := fprpc.latitude; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: floorplanreferencepoint.FieldLatitude,
		})
		fprp.Latitude = *value
	}
	if value := fprpc.longitude; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: floorplanreferencepoint.FieldLongitude,
		})
		fprp.Longitude = *value
	}
	if err := sqlgraph.CreateNode(ctx, fprpc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	fprp.ID = strconv.FormatInt(id, 10)
	return fprp, nil
}
