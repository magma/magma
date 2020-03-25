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
	"github.com/facebookincubator/symphony/graph/ent/floorplanscale"
)

// FloorPlanScaleCreate is the builder for creating a FloorPlanScale entity.
type FloorPlanScaleCreate struct {
	config
	mutation *FloorPlanScaleMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (fpsc *FloorPlanScaleCreate) SetCreateTime(t time.Time) *FloorPlanScaleCreate {
	fpsc.mutation.SetCreateTime(t)
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
	fpsc.mutation.SetUpdateTime(t)
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
	fpsc.mutation.SetReferencePoint1X(i)
	return fpsc
}

// SetReferencePoint1Y sets the reference_point1_y field.
func (fpsc *FloorPlanScaleCreate) SetReferencePoint1Y(i int) *FloorPlanScaleCreate {
	fpsc.mutation.SetReferencePoint1Y(i)
	return fpsc
}

// SetReferencePoint2X sets the reference_point2_x field.
func (fpsc *FloorPlanScaleCreate) SetReferencePoint2X(i int) *FloorPlanScaleCreate {
	fpsc.mutation.SetReferencePoint2X(i)
	return fpsc
}

// SetReferencePoint2Y sets the reference_point2_y field.
func (fpsc *FloorPlanScaleCreate) SetReferencePoint2Y(i int) *FloorPlanScaleCreate {
	fpsc.mutation.SetReferencePoint2Y(i)
	return fpsc
}

// SetScaleInMeters sets the scale_in_meters field.
func (fpsc *FloorPlanScaleCreate) SetScaleInMeters(f float64) *FloorPlanScaleCreate {
	fpsc.mutation.SetScaleInMeters(f)
	return fpsc
}

// Save creates the FloorPlanScale in the database.
func (fpsc *FloorPlanScaleCreate) Save(ctx context.Context) (*FloorPlanScale, error) {
	if _, ok := fpsc.mutation.CreateTime(); !ok {
		v := floorplanscale.DefaultCreateTime()
		fpsc.mutation.SetCreateTime(v)
	}
	if _, ok := fpsc.mutation.UpdateTime(); !ok {
		v := floorplanscale.DefaultUpdateTime()
		fpsc.mutation.SetUpdateTime(v)
	}
	if _, ok := fpsc.mutation.ReferencePoint1X(); !ok {
		return nil, errors.New("ent: missing required field \"reference_point1_x\"")
	}
	if _, ok := fpsc.mutation.ReferencePoint1Y(); !ok {
		return nil, errors.New("ent: missing required field \"reference_point1_y\"")
	}
	if _, ok := fpsc.mutation.ReferencePoint2X(); !ok {
		return nil, errors.New("ent: missing required field \"reference_point2_x\"")
	}
	if _, ok := fpsc.mutation.ReferencePoint2Y(); !ok {
		return nil, errors.New("ent: missing required field \"reference_point2_y\"")
	}
	if _, ok := fpsc.mutation.ScaleInMeters(); !ok {
		return nil, errors.New("ent: missing required field \"scale_in_meters\"")
	}
	var (
		err  error
		node *FloorPlanScale
	)
	if len(fpsc.hooks) == 0 {
		node, err = fpsc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FloorPlanScaleMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fpsc.mutation = mutation
			node, err = fpsc.sqlSave(ctx)
			return node, err
		})
		for i := len(fpsc.hooks); i > 0; i-- {
			mut = fpsc.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, fpsc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
		fps   = &FloorPlanScale{config: fpsc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: floorplanscale.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: floorplanscale.FieldID,
			},
		}
	)
	if value, ok := fpsc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: floorplanscale.FieldCreateTime,
		})
		fps.CreateTime = value
	}
	if value, ok := fpsc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: floorplanscale.FieldUpdateTime,
		})
		fps.UpdateTime = value
	}
	if value, ok := fpsc.mutation.ReferencePoint1X(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint1X,
		})
		fps.ReferencePoint1X = value
	}
	if value, ok := fpsc.mutation.ReferencePoint1Y(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint1Y,
		})
		fps.ReferencePoint1Y = value
	}
	if value, ok := fpsc.mutation.ReferencePoint2X(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint2X,
		})
		fps.ReferencePoint2X = value
	}
	if value, ok := fpsc.mutation.ReferencePoint2Y(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint2Y,
		})
		fps.ReferencePoint2Y = value
	}
	if value, ok := fpsc.mutation.ScaleInMeters(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: floorplanscale.FieldScaleInMeters,
		})
		fps.ScaleInMeters = value
	}
	if err := sqlgraph.CreateNode(ctx, fpsc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	fps.ID = int(id)
	return fps, nil
}
