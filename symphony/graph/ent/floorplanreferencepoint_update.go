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
	"github.com/facebookincubator/symphony/graph/ent/floorplanreferencepoint"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// FloorPlanReferencePointUpdate is the builder for updating FloorPlanReferencePoint entities.
type FloorPlanReferencePointUpdate struct {
	config
	hooks      []Hook
	mutation   *FloorPlanReferencePointMutation
	predicates []predicate.FloorPlanReferencePoint
}

// Where adds a new predicate for the builder.
func (fprpu *FloorPlanReferencePointUpdate) Where(ps ...predicate.FloorPlanReferencePoint) *FloorPlanReferencePointUpdate {
	fprpu.predicates = append(fprpu.predicates, ps...)
	return fprpu
}

// SetX sets the x field.
func (fprpu *FloorPlanReferencePointUpdate) SetX(i int) *FloorPlanReferencePointUpdate {
	fprpu.mutation.ResetX()
	fprpu.mutation.SetX(i)
	return fprpu
}

// AddX adds i to x.
func (fprpu *FloorPlanReferencePointUpdate) AddX(i int) *FloorPlanReferencePointUpdate {
	fprpu.mutation.AddX(i)
	return fprpu
}

// SetY sets the y field.
func (fprpu *FloorPlanReferencePointUpdate) SetY(i int) *FloorPlanReferencePointUpdate {
	fprpu.mutation.ResetY()
	fprpu.mutation.SetY(i)
	return fprpu
}

// AddY adds i to y.
func (fprpu *FloorPlanReferencePointUpdate) AddY(i int) *FloorPlanReferencePointUpdate {
	fprpu.mutation.AddY(i)
	return fprpu
}

// SetLatitude sets the latitude field.
func (fprpu *FloorPlanReferencePointUpdate) SetLatitude(f float64) *FloorPlanReferencePointUpdate {
	fprpu.mutation.ResetLatitude()
	fprpu.mutation.SetLatitude(f)
	return fprpu
}

// AddLatitude adds f to latitude.
func (fprpu *FloorPlanReferencePointUpdate) AddLatitude(f float64) *FloorPlanReferencePointUpdate {
	fprpu.mutation.AddLatitude(f)
	return fprpu
}

// SetLongitude sets the longitude field.
func (fprpu *FloorPlanReferencePointUpdate) SetLongitude(f float64) *FloorPlanReferencePointUpdate {
	fprpu.mutation.ResetLongitude()
	fprpu.mutation.SetLongitude(f)
	return fprpu
}

// AddLongitude adds f to longitude.
func (fprpu *FloorPlanReferencePointUpdate) AddLongitude(f float64) *FloorPlanReferencePointUpdate {
	fprpu.mutation.AddLongitude(f)
	return fprpu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (fprpu *FloorPlanReferencePointUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := fprpu.mutation.UpdateTime(); !ok {
		v := floorplanreferencepoint.UpdateDefaultUpdateTime()
		fprpu.mutation.SetUpdateTime(v)
	}
	var (
		err      error
		affected int
	)
	if len(fprpu.hooks) == 0 {
		affected, err = fprpu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FloorPlanReferencePointMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fprpu.mutation = mutation
			affected, err = fprpu.sqlSave(ctx)
			return affected, err
		})
		for i := len(fprpu.hooks); i > 0; i-- {
			mut = fprpu.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, fprpu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (fprpu *FloorPlanReferencePointUpdate) SaveX(ctx context.Context) int {
	affected, err := fprpu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fprpu *FloorPlanReferencePointUpdate) Exec(ctx context.Context) error {
	_, err := fprpu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fprpu *FloorPlanReferencePointUpdate) ExecX(ctx context.Context) {
	if err := fprpu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fprpu *FloorPlanReferencePointUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   floorplanreferencepoint.Table,
			Columns: floorplanreferencepoint.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: floorplanreferencepoint.FieldID,
			},
		},
	}
	if ps := fprpu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fprpu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: floorplanreferencepoint.FieldUpdateTime,
		})
	}
	if value, ok := fprpu.mutation.X(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanreferencepoint.FieldX,
		})
	}
	if value, ok := fprpu.mutation.AddedX(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanreferencepoint.FieldX,
		})
	}
	if value, ok := fprpu.mutation.Y(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanreferencepoint.FieldY,
		})
	}
	if value, ok := fprpu.mutation.AddedY(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanreferencepoint.FieldY,
		})
	}
	if value, ok := fprpu.mutation.Latitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: floorplanreferencepoint.FieldLatitude,
		})
	}
	if value, ok := fprpu.mutation.AddedLatitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: floorplanreferencepoint.FieldLatitude,
		})
	}
	if value, ok := fprpu.mutation.Longitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: floorplanreferencepoint.FieldLongitude,
		})
	}
	if value, ok := fprpu.mutation.AddedLongitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: floorplanreferencepoint.FieldLongitude,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fprpu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{floorplanreferencepoint.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// FloorPlanReferencePointUpdateOne is the builder for updating a single FloorPlanReferencePoint entity.
type FloorPlanReferencePointUpdateOne struct {
	config
	hooks    []Hook
	mutation *FloorPlanReferencePointMutation
}

// SetX sets the x field.
func (fprpuo *FloorPlanReferencePointUpdateOne) SetX(i int) *FloorPlanReferencePointUpdateOne {
	fprpuo.mutation.ResetX()
	fprpuo.mutation.SetX(i)
	return fprpuo
}

// AddX adds i to x.
func (fprpuo *FloorPlanReferencePointUpdateOne) AddX(i int) *FloorPlanReferencePointUpdateOne {
	fprpuo.mutation.AddX(i)
	return fprpuo
}

// SetY sets the y field.
func (fprpuo *FloorPlanReferencePointUpdateOne) SetY(i int) *FloorPlanReferencePointUpdateOne {
	fprpuo.mutation.ResetY()
	fprpuo.mutation.SetY(i)
	return fprpuo
}

// AddY adds i to y.
func (fprpuo *FloorPlanReferencePointUpdateOne) AddY(i int) *FloorPlanReferencePointUpdateOne {
	fprpuo.mutation.AddY(i)
	return fprpuo
}

// SetLatitude sets the latitude field.
func (fprpuo *FloorPlanReferencePointUpdateOne) SetLatitude(f float64) *FloorPlanReferencePointUpdateOne {
	fprpuo.mutation.ResetLatitude()
	fprpuo.mutation.SetLatitude(f)
	return fprpuo
}

// AddLatitude adds f to latitude.
func (fprpuo *FloorPlanReferencePointUpdateOne) AddLatitude(f float64) *FloorPlanReferencePointUpdateOne {
	fprpuo.mutation.AddLatitude(f)
	return fprpuo
}

// SetLongitude sets the longitude field.
func (fprpuo *FloorPlanReferencePointUpdateOne) SetLongitude(f float64) *FloorPlanReferencePointUpdateOne {
	fprpuo.mutation.ResetLongitude()
	fprpuo.mutation.SetLongitude(f)
	return fprpuo
}

// AddLongitude adds f to longitude.
func (fprpuo *FloorPlanReferencePointUpdateOne) AddLongitude(f float64) *FloorPlanReferencePointUpdateOne {
	fprpuo.mutation.AddLongitude(f)
	return fprpuo
}

// Save executes the query and returns the updated entity.
func (fprpuo *FloorPlanReferencePointUpdateOne) Save(ctx context.Context) (*FloorPlanReferencePoint, error) {
	if _, ok := fprpuo.mutation.UpdateTime(); !ok {
		v := floorplanreferencepoint.UpdateDefaultUpdateTime()
		fprpuo.mutation.SetUpdateTime(v)
	}
	var (
		err  error
		node *FloorPlanReferencePoint
	)
	if len(fprpuo.hooks) == 0 {
		node, err = fprpuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FloorPlanReferencePointMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fprpuo.mutation = mutation
			node, err = fprpuo.sqlSave(ctx)
			return node, err
		})
		for i := len(fprpuo.hooks); i > 0; i-- {
			mut = fprpuo.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, fprpuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (fprpuo *FloorPlanReferencePointUpdateOne) SaveX(ctx context.Context) *FloorPlanReferencePoint {
	fprp, err := fprpuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return fprp
}

// Exec executes the query on the entity.
func (fprpuo *FloorPlanReferencePointUpdateOne) Exec(ctx context.Context) error {
	_, err := fprpuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fprpuo *FloorPlanReferencePointUpdateOne) ExecX(ctx context.Context) {
	if err := fprpuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fprpuo *FloorPlanReferencePointUpdateOne) sqlSave(ctx context.Context) (fprp *FloorPlanReferencePoint, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   floorplanreferencepoint.Table,
			Columns: floorplanreferencepoint.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: floorplanreferencepoint.FieldID,
			},
		},
	}
	id, ok := fprpuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing FloorPlanReferencePoint.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := fprpuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: floorplanreferencepoint.FieldUpdateTime,
		})
	}
	if value, ok := fprpuo.mutation.X(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanreferencepoint.FieldX,
		})
	}
	if value, ok := fprpuo.mutation.AddedX(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanreferencepoint.FieldX,
		})
	}
	if value, ok := fprpuo.mutation.Y(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanreferencepoint.FieldY,
		})
	}
	if value, ok := fprpuo.mutation.AddedY(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: floorplanreferencepoint.FieldY,
		})
	}
	if value, ok := fprpuo.mutation.Latitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: floorplanreferencepoint.FieldLatitude,
		})
	}
	if value, ok := fprpuo.mutation.AddedLatitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: floorplanreferencepoint.FieldLatitude,
		})
	}
	if value, ok := fprpuo.mutation.Longitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: floorplanreferencepoint.FieldLongitude,
		})
	}
	if value, ok := fprpuo.mutation.AddedLongitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: floorplanreferencepoint.FieldLongitude,
		})
	}
	fprp = &FloorPlanReferencePoint{config: fprpuo.config}
	_spec.Assign = fprp.assignValues
	_spec.ScanValues = fprp.scanValues()
	if err = sqlgraph.UpdateNode(ctx, fprpuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{floorplanreferencepoint.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return fprp, nil
}
