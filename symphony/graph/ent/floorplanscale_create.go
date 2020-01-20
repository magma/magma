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
	"github.com/facebookincubator/symphony/graph/ent/floorplanscale"
)

// FloorPlanScaleCreate is the builder for creating a FloorPlanScale entity.
type FloorPlanScaleCreate struct {
	config
	create_time        *time.Time
	update_time        *time.Time
	reference_point1_x *int
	reference_point1_y *int
	reference_point2_x *int
	reference_point2_y *int
	scale_in_meters    *float64
}

// SetCreateTime sets the create_time field.
func (fpsc *FloorPlanScaleCreate) SetCreateTime(t time.Time) *FloorPlanScaleCreate {
	fpsc.create_time = &t
	return fpsc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (fpsc *FloorPlanScaleCreate) SetNillableCreateTime(t *time.Time) *FloorPlanScaleCreate {
	if t != nil {
		fpsc.SetCreateTime(*t)
	}
	return fpsc
}

// SetUpdateTime sets the update_time field.
func (fpsc *FloorPlanScaleCreate) SetUpdateTime(t time.Time) *FloorPlanScaleCreate {
	fpsc.update_time = &t
	return fpsc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (fpsc *FloorPlanScaleCreate) SetNillableUpdateTime(t *time.Time) *FloorPlanScaleCreate {
	if t != nil {
		fpsc.SetUpdateTime(*t)
	}
	return fpsc
}

// SetReferencePoint1X sets the reference_point1_x field.
func (fpsc *FloorPlanScaleCreate) SetReferencePoint1X(i int) *FloorPlanScaleCreate {
	fpsc.reference_point1_x = &i
	return fpsc
}

// SetReferencePoint1Y sets the reference_point1_y field.
func (fpsc *FloorPlanScaleCreate) SetReferencePoint1Y(i int) *FloorPlanScaleCreate {
	fpsc.reference_point1_y = &i
	return fpsc
}

// SetReferencePoint2X sets the reference_point2_x field.
func (fpsc *FloorPlanScaleCreate) SetReferencePoint2X(i int) *FloorPlanScaleCreate {
	fpsc.reference_point2_x = &i
	return fpsc
}

// SetReferencePoint2Y sets the reference_point2_y field.
func (fpsc *FloorPlanScaleCreate) SetReferencePoint2Y(i int) *FloorPlanScaleCreate {
	fpsc.reference_point2_y = &i
	return fpsc
}

// SetScaleInMeters sets the scale_in_meters field.
func (fpsc *FloorPlanScaleCreate) SetScaleInMeters(f float64) *FloorPlanScaleCreate {
	fpsc.scale_in_meters = &f
	return fpsc
}

// Save creates the FloorPlanScale in the database.
func (fpsc *FloorPlanScaleCreate) Save(ctx context.Context) (*FloorPlanScale, error) {
	if fpsc.create_time == nil {
		v := floorplanscale.DefaultCreateTime()
		fpsc.create_time = &v
	}
	if fpsc.update_time == nil {
		v := floorplanscale.DefaultUpdateTime()
		fpsc.update_time = &v
	}
	if fpsc.reference_point1_x == nil {
		return nil, errors.New("ent: missing required field \"reference_point1_x\"")
	}
	if fpsc.reference_point1_y == nil {
		return nil, errors.New("ent: missing required field \"reference_point1_y\"")
	}
	if fpsc.reference_point2_x == nil {
		return nil, errors.New("ent: missing required field \"reference_point2_x\"")
	}
	if fpsc.reference_point2_y == nil {
		return nil, errors.New("ent: missing required field \"reference_point2_y\"")
	}
	if fpsc.scale_in_meters == nil {
		return nil, errors.New("ent: missing required field \"scale_in_meters\"")
	}
	return fpsc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (fpsc *FloorPlanScaleCreate) SaveX(ctx context.Context) *FloorPlanScale {
	v, err := fpsc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (fpsc *FloorPlanScaleCreate) sqlSave(ctx context.Context) (*FloorPlanScale, error) {
	var (
		fps  = &FloorPlanScale{config: fpsc.config}
		spec = &sqlgraph.CreateSpec{
			Table: floorplanscale.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: floorplanscale.FieldID,
			},
		}
	)
	if value := fpsc.create_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: floorplanscale.FieldCreateTime,
		})
		fps.CreateTime = *value
	}
	if value := fpsc.update_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: floorplanscale.FieldUpdateTime,
		})
		fps.UpdateTime = *value
	}
	if value := fpsc.reference_point1_x; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint1X,
		})
		fps.ReferencePoint1X = *value
	}
	if value := fpsc.reference_point1_y; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint1Y,
		})
		fps.ReferencePoint1Y = *value
	}
	if value := fpsc.reference_point2_x; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint2X,
		})
		fps.ReferencePoint2X = *value
	}
	if value := fpsc.reference_point2_y; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint2Y,
		})
		fps.ReferencePoint2Y = *value
	}
	if value := fpsc.scale_in_meters; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: floorplanscale.FieldScaleInMeters,
		})
		fps.ScaleInMeters = *value
	}
	if err := sqlgraph.CreateNode(ctx, fpsc.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := spec.ID.Value.(int64)
	fps.ID = strconv.FormatInt(id, 10)
	return fps, nil
}
