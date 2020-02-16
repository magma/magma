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
	"github.com/facebookincubator/symphony/graph/ent/floorplanreferencepoint"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// FloorPlanReferencePointUpdate is the builder for updating FloorPlanReferencePoint entities.
type FloorPlanReferencePointUpdate struct {
	config

	update_time  *time.Time
	x            *int
	addx         *int
	y            *int
	addy         *int
	latitude     *float64
	addlatitude  *float64
	longitude    *float64
	addlongitude *float64
	predicates   []predicate.FloorPlanReferencePoint
}

// Where adds a new predicate for the builder.
func (fprpu *FloorPlanReferencePointUpdate) Where(ps ...predicate.FloorPlanReferencePoint) *FloorPlanReferencePointUpdate {
	fprpu.predicates = append(fprpu.predicates, ps...)
	return fprpu
}

// SetX sets the x field.
func (fprpu *FloorPlanReferencePointUpdate) SetX(i int) *FloorPlanReferencePointUpdate {
	fprpu.x = &i
	fprpu.addx = nil
	return fprpu
}

// AddX adds i to x.
func (fprpu *FloorPlanReferencePointUpdate) AddX(i int) *FloorPlanReferencePointUpdate {
	if fprpu.addx == nil {
		fprpu.addx = &i
	} else {
		*fprpu.addx += i
	}
	return fprpu
}

// SetY sets the y field.
func (fprpu *FloorPlanReferencePointUpdate) SetY(i int) *FloorPlanReferencePointUpdate {
	fprpu.y = &i
	fprpu.addy = nil
	return fprpu
}

// AddY adds i to y.
func (fprpu *FloorPlanReferencePointUpdate) AddY(i int) *FloorPlanReferencePointUpdate {
	if fprpu.addy == nil {
		fprpu.addy = &i
	} else {
		*fprpu.addy += i
	}
	return fprpu
}

// SetLatitude sets the latitude field.
func (fprpu *FloorPlanReferencePointUpdate) SetLatitude(f float64) *FloorPlanReferencePointUpdate {
	fprpu.latitude = &f
	fprpu.addlatitude = nil
	return fprpu
}

// AddLatitude adds f to latitude.
func (fprpu *FloorPlanReferencePointUpdate) AddLatitude(f float64) *FloorPlanReferencePointUpdate {
	if fprpu.addlatitude == nil {
		fprpu.addlatitude = &f
	} else {
		*fprpu.addlatitude += f
	}
	return fprpu
}

// SetLongitude sets the longitude field.
func (fprpu *FloorPlanReferencePointUpdate) SetLongitude(f float64) *FloorPlanReferencePointUpdate {
	fprpu.longitude = &f
	fprpu.addlongitude = nil
	return fprpu
}

// AddLongitude adds f to longitude.
func (fprpu *FloorPlanReferencePointUpdate) AddLongitude(f float64) *FloorPlanReferencePointUpdate {
	if fprpu.addlongitude == nil {
		fprpu.addlongitude = &f
	} else {
		*fprpu.addlongitude += f
	}
	return fprpu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (fprpu *FloorPlanReferencePointUpdate) Save(ctx context.Context) (int, error) {
	if fprpu.update_time == nil {
		v := floorplanreferencepoint.UpdateDefaultUpdateTime()
		fprpu.update_time = &v
	}
	return fprpu.sqlSave(ctx)
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
				Type:   field.TypeString,
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
	if value := fprpu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: floorplanreferencepoint.FieldUpdateTime,
		})
	}
	if value := fprpu.x; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanreferencepoint.FieldX,
		})
	}
	if value := fprpu.addx; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanreferencepoint.FieldX,
		})
	}
	if value := fprpu.y; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanreferencepoint.FieldY,
		})
	}
	if value := fprpu.addy; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanreferencepoint.FieldY,
		})
	}
	if value := fprpu.latitude; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: floorplanreferencepoint.FieldLatitude,
		})
	}
	if value := fprpu.addlatitude; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: floorplanreferencepoint.FieldLatitude,
		})
	}
	if value := fprpu.longitude; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: floorplanreferencepoint.FieldLongitude,
		})
	}
	if value := fprpu.addlongitude; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
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
	id string

	update_time  *time.Time
	x            *int
	addx         *int
	y            *int
	addy         *int
	latitude     *float64
	addlatitude  *float64
	longitude    *float64
	addlongitude *float64
}

// SetX sets the x field.
func (fprpuo *FloorPlanReferencePointUpdateOne) SetX(i int) *FloorPlanReferencePointUpdateOne {
	fprpuo.x = &i
	fprpuo.addx = nil
	return fprpuo
}

// AddX adds i to x.
func (fprpuo *FloorPlanReferencePointUpdateOne) AddX(i int) *FloorPlanReferencePointUpdateOne {
	if fprpuo.addx == nil {
		fprpuo.addx = &i
	} else {
		*fprpuo.addx += i
	}
	return fprpuo
}

// SetY sets the y field.
func (fprpuo *FloorPlanReferencePointUpdateOne) SetY(i int) *FloorPlanReferencePointUpdateOne {
	fprpuo.y = &i
	fprpuo.addy = nil
	return fprpuo
}

// AddY adds i to y.
func (fprpuo *FloorPlanReferencePointUpdateOne) AddY(i int) *FloorPlanReferencePointUpdateOne {
	if fprpuo.addy == nil {
		fprpuo.addy = &i
	} else {
		*fprpuo.addy += i
	}
	return fprpuo
}

// SetLatitude sets the latitude field.
func (fprpuo *FloorPlanReferencePointUpdateOne) SetLatitude(f float64) *FloorPlanReferencePointUpdateOne {
	fprpuo.latitude = &f
	fprpuo.addlatitude = nil
	return fprpuo
}

// AddLatitude adds f to latitude.
func (fprpuo *FloorPlanReferencePointUpdateOne) AddLatitude(f float64) *FloorPlanReferencePointUpdateOne {
	if fprpuo.addlatitude == nil {
		fprpuo.addlatitude = &f
	} else {
		*fprpuo.addlatitude += f
	}
	return fprpuo
}

// SetLongitude sets the longitude field.
func (fprpuo *FloorPlanReferencePointUpdateOne) SetLongitude(f float64) *FloorPlanReferencePointUpdateOne {
	fprpuo.longitude = &f
	fprpuo.addlongitude = nil
	return fprpuo
}

// AddLongitude adds f to longitude.
func (fprpuo *FloorPlanReferencePointUpdateOne) AddLongitude(f float64) *FloorPlanReferencePointUpdateOne {
	if fprpuo.addlongitude == nil {
		fprpuo.addlongitude = &f
	} else {
		*fprpuo.addlongitude += f
	}
	return fprpuo
}

// Save executes the query and returns the updated entity.
func (fprpuo *FloorPlanReferencePointUpdateOne) Save(ctx context.Context) (*FloorPlanReferencePoint, error) {
	if fprpuo.update_time == nil {
		v := floorplanreferencepoint.UpdateDefaultUpdateTime()
		fprpuo.update_time = &v
	}
	return fprpuo.sqlSave(ctx)
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
				Value:  fprpuo.id,
				Type:   field.TypeString,
				Column: floorplanreferencepoint.FieldID,
			},
		},
	}
	if value := fprpuo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: floorplanreferencepoint.FieldUpdateTime,
		})
	}
	if value := fprpuo.x; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanreferencepoint.FieldX,
		})
	}
	if value := fprpuo.addx; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanreferencepoint.FieldX,
		})
	}
	if value := fprpuo.y; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanreferencepoint.FieldY,
		})
	}
	if value := fprpuo.addy; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: floorplanreferencepoint.FieldY,
		})
	}
	if value := fprpuo.latitude; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: floorplanreferencepoint.FieldLatitude,
		})
	}
	if value := fprpuo.addlatitude; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: floorplanreferencepoint.FieldLatitude,
		})
	}
	if value := fprpuo.longitude; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: floorplanreferencepoint.FieldLongitude,
		})
	}
	if value := fprpuo.addlongitude; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
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
