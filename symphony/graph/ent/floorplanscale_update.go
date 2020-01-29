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
	"github.com/facebookincubator/symphony/graph/ent/floorplanscale"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// FloorPlanScaleUpdate is the builder for updating FloorPlanScale entities.
type FloorPlanScaleUpdate struct {
	config

	update_time           *time.Time
	reference_point1_x    *int
	addreference_point1_x *int
	reference_point1_y    *int
	addreference_point1_y *int
	reference_point2_x    *int
	addreference_point2_x *int
	reference_point2_y    *int
	addreference_point2_y *int
	scale_in_meters       *float64
	addscale_in_meters    *float64
	predicates            []predicate.FloorPlanScale
}

// Where adds a new predicate for the builder.
func (fpsu *FloorPlanScaleUpdate) Where(ps ...predicate.FloorPlanScale) *FloorPlanScaleUpdate {
	fpsu.predicates = append(fpsu.predicates, ps...)
	return fpsu
}

// SetReferencePoint1X sets the reference_point1_x field.
func (fpsu *FloorPlanScaleUpdate) SetReferencePoint1X(i int) *FloorPlanScaleUpdate {
	fpsu.reference_point1_x = &i
	fpsu.addreference_point1_x = nil
	return fpsu
}

// AddReferencePoint1X adds i to reference_point1_x.
func (fpsu *FloorPlanScaleUpdate) AddReferencePoint1X(i int) *FloorPlanScaleUpdate {
	if fpsu.addreference_point1_x == nil {
		fpsu.addreference_point1_x = &i
	} else {
		*fpsu.addreference_point1_x += i
	}
	return fpsu
}

// SetReferencePoint1Y sets the reference_point1_y field.
func (fpsu *FloorPlanScaleUpdate) SetReferencePoint1Y(i int) *FloorPlanScaleUpdate {
	fpsu.reference_point1_y = &i
	fpsu.addreference_point1_y = nil
	return fpsu
}

// AddReferencePoint1Y adds i to reference_point1_y.
func (fpsu *FloorPlanScaleUpdate) AddReferencePoint1Y(i int) *FloorPlanScaleUpdate {
	if fpsu.addreference_point1_y == nil {
		fpsu.addreference_point1_y = &i
	} else {
		*fpsu.addreference_point1_y += i
	}
	return fpsu
}

// SetReferencePoint2X sets the reference_point2_x field.
func (fpsu *FloorPlanScaleUpdate) SetReferencePoint2X(i int) *FloorPlanScaleUpdate {
	fpsu.reference_point2_x = &i
	fpsu.addreference_point2_x = nil
	return fpsu
}

// AddReferencePoint2X adds i to reference_point2_x.
func (fpsu *FloorPlanScaleUpdate) AddReferencePoint2X(i int) *FloorPlanScaleUpdate {
	if fpsu.addreference_point2_x == nil {
		fpsu.addreference_point2_x = &i
	} else {
		*fpsu.addreference_point2_x += i
	}
	return fpsu
}

// SetReferencePoint2Y sets the reference_point2_y field.
func (fpsu *FloorPlanScaleUpdate) SetReferencePoint2Y(i int) *FloorPlanScaleUpdate {
	fpsu.reference_point2_y = &i
	fpsu.addreference_point2_y = nil
	return fpsu
}

// AddReferencePoint2Y adds i to reference_point2_y.
func (fpsu *FloorPlanScaleUpdate) AddReferencePoint2Y(i int) *FloorPlanScaleUpdate {
	if fpsu.addreference_point2_y == nil {
		fpsu.addreference_point2_y = &i
	} else {
		*fpsu.addreference_point2_y += i
	}
	return fpsu
}

// SetScaleInMeters sets the scale_in_meters field.
func (fpsu *FloorPlanScaleUpdate) SetScaleInMeters(f float64) *FloorPlanScaleUpdate {
	fpsu.scale_in_meters = &f
	fpsu.addscale_in_meters = nil
	return fpsu
}

// AddScaleInMeters adds f to scale_in_meters.
func (fpsu *FloorPlanScaleUpdate) AddScaleInMeters(f float64) *FloorPlanScaleUpdate {
	if fpsu.addscale_in_meters == nil {
		fpsu.addscale_in_meters = &f
	} else {
		*fpsu.addscale_in_meters += f
	}
	return fpsu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (fpsu *FloorPlanScaleUpdate) Save(ctx context.Context) (int, error) {
	if fpsu.update_time == nil {
		v := floorplanscale.UpdateDefaultUpdateTime()
		fpsu.update_time = &v
	}
	return fpsu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (fpsu *FloorPlanScaleUpdate) SaveX(ctx context.Context) int {
	affected, err := fpsu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fpsu *FloorPlanScaleUpdate) Exec(ctx context.Context) error {
	_, err := fpsu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fpsu *FloorPlanScaleUpdate) ExecX(ctx context.Context) {
	if err := fpsu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fpsu *FloorPlanScaleUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   floorplanscale.Table,
			Columns: floorplanscale.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: floorplanscale.FieldID,
			},
		},
	}
	if ps := fpsu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := fpsu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: floorplanscale.FieldUpdateTime,
		})
	}
	if value := fpsu.reference_point1_x; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint1X,
		})
	}
	if value := fpsu.addreference_point1_x; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint1X,
		})
	}
	if value := fpsu.reference_point1_y; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint1Y,
		})
	}
	if value := fpsu.addreference_point1_y; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint1Y,
		})
	}
	if value := fpsu.reference_point2_x; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint2X,
		})
	}
	if value := fpsu.addreference_point2_x; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint2X,
		})
	}
	if value := fpsu.reference_point2_y; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint2Y,
		})
	}
	if value := fpsu.addreference_point2_y; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint2Y,
		})
	}
	if value := fpsu.scale_in_meters; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: floorplanscale.FieldScaleInMeters,
		})
	}
	if value := fpsu.addscale_in_meters; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: floorplanscale.FieldScaleInMeters,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fpsu.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// FloorPlanScaleUpdateOne is the builder for updating a single FloorPlanScale entity.
type FloorPlanScaleUpdateOne struct {
	config
	id string

	update_time           *time.Time
	reference_point1_x    *int
	addreference_point1_x *int
	reference_point1_y    *int
	addreference_point1_y *int
	reference_point2_x    *int
	addreference_point2_x *int
	reference_point2_y    *int
	addreference_point2_y *int
	scale_in_meters       *float64
	addscale_in_meters    *float64
}

// SetReferencePoint1X sets the reference_point1_x field.
func (fpsuo *FloorPlanScaleUpdateOne) SetReferencePoint1X(i int) *FloorPlanScaleUpdateOne {
	fpsuo.reference_point1_x = &i
	fpsuo.addreference_point1_x = nil
	return fpsuo
}

// AddReferencePoint1X adds i to reference_point1_x.
func (fpsuo *FloorPlanScaleUpdateOne) AddReferencePoint1X(i int) *FloorPlanScaleUpdateOne {
	if fpsuo.addreference_point1_x == nil {
		fpsuo.addreference_point1_x = &i
	} else {
		*fpsuo.addreference_point1_x += i
	}
	return fpsuo
}

// SetReferencePoint1Y sets the reference_point1_y field.
func (fpsuo *FloorPlanScaleUpdateOne) SetReferencePoint1Y(i int) *FloorPlanScaleUpdateOne {
	fpsuo.reference_point1_y = &i
	fpsuo.addreference_point1_y = nil
	return fpsuo
}

// AddReferencePoint1Y adds i to reference_point1_y.
func (fpsuo *FloorPlanScaleUpdateOne) AddReferencePoint1Y(i int) *FloorPlanScaleUpdateOne {
	if fpsuo.addreference_point1_y == nil {
		fpsuo.addreference_point1_y = &i
	} else {
		*fpsuo.addreference_point1_y += i
	}
	return fpsuo
}

// SetReferencePoint2X sets the reference_point2_x field.
func (fpsuo *FloorPlanScaleUpdateOne) SetReferencePoint2X(i int) *FloorPlanScaleUpdateOne {
	fpsuo.reference_point2_x = &i
	fpsuo.addreference_point2_x = nil
	return fpsuo
}

// AddReferencePoint2X adds i to reference_point2_x.
func (fpsuo *FloorPlanScaleUpdateOne) AddReferencePoint2X(i int) *FloorPlanScaleUpdateOne {
	if fpsuo.addreference_point2_x == nil {
		fpsuo.addreference_point2_x = &i
	} else {
		*fpsuo.addreference_point2_x += i
	}
	return fpsuo
}

// SetReferencePoint2Y sets the reference_point2_y field.
func (fpsuo *FloorPlanScaleUpdateOne) SetReferencePoint2Y(i int) *FloorPlanScaleUpdateOne {
	fpsuo.reference_point2_y = &i
	fpsuo.addreference_point2_y = nil
	return fpsuo
}

// AddReferencePoint2Y adds i to reference_point2_y.
func (fpsuo *FloorPlanScaleUpdateOne) AddReferencePoint2Y(i int) *FloorPlanScaleUpdateOne {
	if fpsuo.addreference_point2_y == nil {
		fpsuo.addreference_point2_y = &i
	} else {
		*fpsuo.addreference_point2_y += i
	}
	return fpsuo
}

// SetScaleInMeters sets the scale_in_meters field.
func (fpsuo *FloorPlanScaleUpdateOne) SetScaleInMeters(f float64) *FloorPlanScaleUpdateOne {
	fpsuo.scale_in_meters = &f
	fpsuo.addscale_in_meters = nil
	return fpsuo
}

// AddScaleInMeters adds f to scale_in_meters.
func (fpsuo *FloorPlanScaleUpdateOne) AddScaleInMeters(f float64) *FloorPlanScaleUpdateOne {
	if fpsuo.addscale_in_meters == nil {
		fpsuo.addscale_in_meters = &f
	} else {
		*fpsuo.addscale_in_meters += f
	}
	return fpsuo
}

// Save executes the query and returns the updated entity.
func (fpsuo *FloorPlanScaleUpdateOne) Save(ctx context.Context) (*FloorPlanScale, error) {
	if fpsuo.update_time == nil {
		v := floorplanscale.UpdateDefaultUpdateTime()
		fpsuo.update_time = &v
	}
	return fpsuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (fpsuo *FloorPlanScaleUpdateOne) SaveX(ctx context.Context) *FloorPlanScale {
	fps, err := fpsuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return fps
}

// Exec executes the query on the entity.
func (fpsuo *FloorPlanScaleUpdateOne) Exec(ctx context.Context) error {
	_, err := fpsuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fpsuo *FloorPlanScaleUpdateOne) ExecX(ctx context.Context) {
	if err := fpsuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fpsuo *FloorPlanScaleUpdateOne) sqlSave(ctx context.Context) (fps *FloorPlanScale, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   floorplanscale.Table,
			Columns: floorplanscale.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  fpsuo.id,
				Type:   field.TypeString,
				Column: floorplanscale.FieldID,
			},
		},
	}
	if value := fpsuo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: floorplanscale.FieldUpdateTime,
		})
	}
	if value := fpsuo.reference_point1_x; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint1X,
		})
	}
	if value := fpsuo.addreference_point1_x; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint1X,
		})
	}
	if value := fpsuo.reference_point1_y; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint1Y,
		})
	}
	if value := fpsuo.addreference_point1_y; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint1Y,
		})
	}
	if value := fpsuo.reference_point2_x; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint2X,
		})
	}
	if value := fpsuo.addreference_point2_x; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint2X,
		})
	}
	if value := fpsuo.reference_point2_y; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint2Y,
		})
	}
	if value := fpsuo.addreference_point2_y; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanscale.FieldReferencePoint2Y,
		})
	}
	if value := fpsuo.scale_in_meters; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: floorplanscale.FieldScaleInMeters,
		})
	}
	if value := fpsuo.addscale_in_meters; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: floorplanscale.FieldScaleInMeters,
		})
	}
	fps = &FloorPlanScale{config: fpsuo.config}
	_spec.Assign = fps.assignValues
	_spec.ScanValues = fps.scanValues()
	if err = sqlgraph.UpdateNode(ctx, fpsuo.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return fps, nil
}
