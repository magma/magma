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
	"github.com/facebookincubator/symphony/graph/ent/floorplanscale"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// FloorPlanScaleUpdate is the builder for updating FloorPlanScale entities.
type FloorPlanScaleUpdate struct {
	config
	hooks      []Hook
	mutation   *FloorPlanScaleMutation
	predicates []predicate.FloorPlanScale
}

// Where adds a new predicate for the builder.
func (fpsu *FloorPlanScaleUpdate) Where(ps ...predicate.FloorPlanScale) *FloorPlanScaleUpdate {
	fpsu.predicates = append(fpsu.predicates, ps...)
	return fpsu
}

// SetReferencePoint1X sets the reference_point1_x field.
func (fpsu *FloorPlanScaleUpdate) SetReferencePoint1X(i int) *FloorPlanScaleUpdate {
	fpsu.mutation.ResetReferencePoint1X()
	fpsu.mutation.SetReferencePoint1X(i)
	return fpsu
}

// AddReferencePoint1X adds i to reference_point1_x.
func (fpsu *FloorPlanScaleUpdate) AddReferencePoint1X(i int) *FloorPlanScaleUpdate {
	fpsu.mutation.AddReferencePoint1X(i)
	return fpsu
}

// SetReferencePoint1Y sets the reference_point1_y field.
func (fpsu *FloorPlanScaleUpdate) SetReferencePoint1Y(i int) *FloorPlanScaleUpdate {
	fpsu.mutation.ResetReferencePoint1Y()
	fpsu.mutation.SetReferencePoint1Y(i)
	return fpsu
}

// AddReferencePoint1Y adds i to reference_point1_y.
func (fpsu *FloorPlanScaleUpdate) AddReferencePoint1Y(i int) *FloorPlanScaleUpdate {
	fpsu.mutation.AddReferencePoint1Y(i)
	return fpsu
}

// SetReferencePoint2X sets the reference_point2_x field.
func (fpsu *FloorPlanScaleUpdate) SetReferencePoint2X(i int) *FloorPlanScaleUpdate {
	fpsu.mutation.ResetReferencePoint2X()
	fpsu.mutation.SetReferencePoint2X(i)
	return fpsu
}

// AddReferencePoint2X adds i to reference_point2_x.
func (fpsu *FloorPlanScaleUpdate) AddReferencePoint2X(i int) *FloorPlanScaleUpdate {
	fpsu.mutation.AddReferencePoint2X(i)
	return fpsu
}

// SetReferencePoint2Y sets the reference_point2_y field.
func (fpsu *FloorPlanScaleUpdate) SetReferencePoint2Y(i int) *FloorPlanScaleUpdate {
	fpsu.mutation.ResetReferencePoint2Y()
	fpsu.mutation.SetReferencePoint2Y(i)
	return fpsu
}

// AddReferencePoint2Y adds i to reference_point2_y.
func (fpsu *FloorPlanScaleUpdate) AddReferencePoint2Y(i int) *FloorPlanScaleUpdate {
	fpsu.mutation.AddReferencePoint2Y(i)
	return fpsu
}

// SetScaleInMeters sets the scale_in_meters field.
func (fpsu *FloorPlanScaleUpdate) SetScaleInMeters(f float64) *FloorPlanScaleUpdate {
	fpsu.mutation.ResetScaleInMeters()
	fpsu.mutation.SetScaleInMeters(f)
	return fpsu
}

// AddScaleInMeters adds f to scale_in_meters.
func (fpsu *FloorPlanScaleUpdate) AddScaleInMeters(f float64) *FloorPlanScaleUpdate {
	fpsu.mutation.AddScaleInMeters(f)
	return fpsu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (fpsu *FloorPlanScaleUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := fpsu.mutation.UpdateTime(); !ok {
		v := floorplanscale.UpdateDefaultUpdateTime()
		fpsu.mutation.SetUpdateTime(v)
	}
	var (
		err      error
		affected int
	)
	if len(fpsu.hooks) == 0 {
		affected, err = fpsu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FloorPlanScaleMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fpsu.mutation = mutation
			affected, err = fpsu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(fpsu.hooks) - 1; i >= 0; i-- {
			mut = fpsu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fpsu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
				Type:   field.TypeInt,
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
	if value, ok := fpsu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: floorplanscale.FieldUpdateTime,
		})
	}
	if value, ok := fpsu.mutation.ReferencePoint1X(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint1X,
		})
	}
	if value, ok := fpsu.mutation.AddedReferencePoint1X(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint1X,
		})
	}
	if value, ok := fpsu.mutation.ReferencePoint1Y(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint1Y,
		})
	}
	if value, ok := fpsu.mutation.AddedReferencePoint1Y(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint1Y,
		})
	}
	if value, ok := fpsu.mutation.ReferencePoint2X(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint2X,
		})
	}
	if value, ok := fpsu.mutation.AddedReferencePoint2X(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint2X,
		})
	}
	if value, ok := fpsu.mutation.ReferencePoint2Y(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint2Y,
		})
	}
	if value, ok := fpsu.mutation.AddedReferencePoint2Y(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint2Y,
		})
	}
	if value, ok := fpsu.mutation.ScaleInMeters(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: floorplanscale.FieldScaleInMeters,
		})
	}
	if value, ok := fpsu.mutation.AddedScaleInMeters(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: floorplanscale.FieldScaleInMeters,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fpsu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{floorplanscale.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// FloorPlanScaleUpdateOne is the builder for updating a single FloorPlanScale entity.
type FloorPlanScaleUpdateOne struct {
	config
	hooks    []Hook
	mutation *FloorPlanScaleMutation
}

// SetReferencePoint1X sets the reference_point1_x field.
func (fpsuo *FloorPlanScaleUpdateOne) SetReferencePoint1X(i int) *FloorPlanScaleUpdateOne {
	fpsuo.mutation.ResetReferencePoint1X()
	fpsuo.mutation.SetReferencePoint1X(i)
	return fpsuo
}

// AddReferencePoint1X adds i to reference_point1_x.
func (fpsuo *FloorPlanScaleUpdateOne) AddReferencePoint1X(i int) *FloorPlanScaleUpdateOne {
	fpsuo.mutation.AddReferencePoint1X(i)
	return fpsuo
}

// SetReferencePoint1Y sets the reference_point1_y field.
func (fpsuo *FloorPlanScaleUpdateOne) SetReferencePoint1Y(i int) *FloorPlanScaleUpdateOne {
	fpsuo.mutation.ResetReferencePoint1Y()
	fpsuo.mutation.SetReferencePoint1Y(i)
	return fpsuo
}

// AddReferencePoint1Y adds i to reference_point1_y.
func (fpsuo *FloorPlanScaleUpdateOne) AddReferencePoint1Y(i int) *FloorPlanScaleUpdateOne {
	fpsuo.mutation.AddReferencePoint1Y(i)
	return fpsuo
}

// SetReferencePoint2X sets the reference_point2_x field.
func (fpsuo *FloorPlanScaleUpdateOne) SetReferencePoint2X(i int) *FloorPlanScaleUpdateOne {
	fpsuo.mutation.ResetReferencePoint2X()
	fpsuo.mutation.SetReferencePoint2X(i)
	return fpsuo
}

// AddReferencePoint2X adds i to reference_point2_x.
func (fpsuo *FloorPlanScaleUpdateOne) AddReferencePoint2X(i int) *FloorPlanScaleUpdateOne {
	fpsuo.mutation.AddReferencePoint2X(i)
	return fpsuo
}

// SetReferencePoint2Y sets the reference_point2_y field.
func (fpsuo *FloorPlanScaleUpdateOne) SetReferencePoint2Y(i int) *FloorPlanScaleUpdateOne {
	fpsuo.mutation.ResetReferencePoint2Y()
	fpsuo.mutation.SetReferencePoint2Y(i)
	return fpsuo
}

// AddReferencePoint2Y adds i to reference_point2_y.
func (fpsuo *FloorPlanScaleUpdateOne) AddReferencePoint2Y(i int) *FloorPlanScaleUpdateOne {
	fpsuo.mutation.AddReferencePoint2Y(i)
	return fpsuo
}

// SetScaleInMeters sets the scale_in_meters field.
func (fpsuo *FloorPlanScaleUpdateOne) SetScaleInMeters(f float64) *FloorPlanScaleUpdateOne {
	fpsuo.mutation.ResetScaleInMeters()
	fpsuo.mutation.SetScaleInMeters(f)
	return fpsuo
}

// AddScaleInMeters adds f to scale_in_meters.
func (fpsuo *FloorPlanScaleUpdateOne) AddScaleInMeters(f float64) *FloorPlanScaleUpdateOne {
	fpsuo.mutation.AddScaleInMeters(f)
	return fpsuo
}

// Save executes the query and returns the updated entity.
func (fpsuo *FloorPlanScaleUpdateOne) Save(ctx context.Context) (*FloorPlanScale, error) {
	if _, ok := fpsuo.mutation.UpdateTime(); !ok {
		v := floorplanscale.UpdateDefaultUpdateTime()
		fpsuo.mutation.SetUpdateTime(v)
	}
	var (
		err  error
		node *FloorPlanScale
	)
	if len(fpsuo.hooks) == 0 {
		node, err = fpsuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FloorPlanScaleMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fpsuo.mutation = mutation
			node, err = fpsuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(fpsuo.hooks) - 1; i >= 0; i-- {
			mut = fpsuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fpsuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
				Type:   field.TypeInt,
				Column: floorplanscale.FieldID,
			},
		},
	}
	id, ok := fpsuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing FloorPlanScale.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := fpsuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: floorplanscale.FieldUpdateTime,
		})
	}
	if value, ok := fpsuo.mutation.ReferencePoint1X(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint1X,
		})
	}
	if value, ok := fpsuo.mutation.AddedReferencePoint1X(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint1X,
		})
	}
	if value, ok := fpsuo.mutation.ReferencePoint1Y(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint1Y,
		})
	}
	if value, ok := fpsuo.mutation.AddedReferencePoint1Y(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint1Y,
		})
	}
	if value, ok := fpsuo.mutation.ReferencePoint2X(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint2X,
		})
	}
	if value, ok := fpsuo.mutation.AddedReferencePoint2X(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint2X,
		})
	}
	if value, ok := fpsuo.mutation.ReferencePoint2Y(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint2Y,
		})
	}
	if value, ok := fpsuo.mutation.AddedReferencePoint2Y(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanscale.FieldReferencePoint2Y,
		})
	}
	if value, ok := fpsuo.mutation.ScaleInMeters(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: floorplanscale.FieldScaleInMeters,
		})
	}
	if value, ok := fpsuo.mutation.AddedScaleInMeters(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: floorplanscale.FieldScaleInMeters,
		})
	}
	fps = &FloorPlanScale{config: fpsuo.config}
	_spec.Assign = fps.assignValues
	_spec.ScanValues = fps.scanValues()
	if err = sqlgraph.UpdateNode(ctx, fpsuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{floorplanscale.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return fps, nil
}
