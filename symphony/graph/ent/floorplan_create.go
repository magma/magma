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
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/floorplan"
	"github.com/facebookincubator/symphony/graph/ent/floorplanreferencepoint"
	"github.com/facebookincubator/symphony/graph/ent/floorplanscale"
	"github.com/facebookincubator/symphony/graph/ent/location"
)

// FloorPlanCreate is the builder for creating a FloorPlan entity.
type FloorPlanCreate struct {
	config
	mutation *FloorPlanMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (fpc *FloorPlanCreate) SetCreateTime(t time.Time) *FloorPlanCreate {
	fpc.mutation.SetCreateTime(t)
	return fpc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (fpc *FloorPlanCreate) SetNillableCreateTime(t *time.Time) *FloorPlanCreate {
	if t != nil {
		fpc.SetCreateTime(*t)
	}
	return fpc
}

// SetUpdateTime sets the update_time field.
func (fpc *FloorPlanCreate) SetUpdateTime(t time.Time) *FloorPlanCreate {
	fpc.mutation.SetUpdateTime(t)
	return fpc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (fpc *FloorPlanCreate) SetNillableUpdateTime(t *time.Time) *FloorPlanCreate {
	if t != nil {
		fpc.SetUpdateTime(*t)
	}
	return fpc
}

// SetName sets the name field.
func (fpc *FloorPlanCreate) SetName(s string) *FloorPlanCreate {
	fpc.mutation.SetName(s)
	return fpc
}

// SetLocationID sets the location edge to Location by id.
func (fpc *FloorPlanCreate) SetLocationID(id int) *FloorPlanCreate {
	fpc.mutation.SetLocationID(id)
	return fpc
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (fpc *FloorPlanCreate) SetNillableLocationID(id *int) *FloorPlanCreate {
	if id != nil {
		fpc = fpc.SetLocationID(*id)
	}
	return fpc
}

// SetLocation sets the location edge to Location.
func (fpc *FloorPlanCreate) SetLocation(l *Location) *FloorPlanCreate {
	return fpc.SetLocationID(l.ID)
}

// SetReferencePointID sets the reference_point edge to FloorPlanReferencePoint by id.
func (fpc *FloorPlanCreate) SetReferencePointID(id int) *FloorPlanCreate {
	fpc.mutation.SetReferencePointID(id)
	return fpc
}

// SetNillableReferencePointID sets the reference_point edge to FloorPlanReferencePoint by id if the given value is not nil.
func (fpc *FloorPlanCreate) SetNillableReferencePointID(id *int) *FloorPlanCreate {
	if id != nil {
		fpc = fpc.SetReferencePointID(*id)
	}
	return fpc
}

// SetReferencePoint sets the reference_point edge to FloorPlanReferencePoint.
func (fpc *FloorPlanCreate) SetReferencePoint(f *FloorPlanReferencePoint) *FloorPlanCreate {
	return fpc.SetReferencePointID(f.ID)
}

// SetScaleID sets the scale edge to FloorPlanScale by id.
func (fpc *FloorPlanCreate) SetScaleID(id int) *FloorPlanCreate {
	fpc.mutation.SetScaleID(id)
	return fpc
}

// SetNillableScaleID sets the scale edge to FloorPlanScale by id if the given value is not nil.
func (fpc *FloorPlanCreate) SetNillableScaleID(id *int) *FloorPlanCreate {
	if id != nil {
		fpc = fpc.SetScaleID(*id)
	}
	return fpc
}

// SetScale sets the scale edge to FloorPlanScale.
func (fpc *FloorPlanCreate) SetScale(f *FloorPlanScale) *FloorPlanCreate {
	return fpc.SetScaleID(f.ID)
}

// SetImageID sets the image edge to File by id.
func (fpc *FloorPlanCreate) SetImageID(id int) *FloorPlanCreate {
	fpc.mutation.SetImageID(id)
	return fpc
}

// SetNillableImageID sets the image edge to File by id if the given value is not nil.
func (fpc *FloorPlanCreate) SetNillableImageID(id *int) *FloorPlanCreate {
	if id != nil {
		fpc = fpc.SetImageID(*id)
	}
	return fpc
}

// SetImage sets the image edge to File.
func (fpc *FloorPlanCreate) SetImage(f *File) *FloorPlanCreate {
	return fpc.SetImageID(f.ID)
}

// Save creates the FloorPlan in the database.
func (fpc *FloorPlanCreate) Save(ctx context.Context) (*FloorPlan, error) {
	if _, ok := fpc.mutation.CreateTime(); !ok {
		v := floorplan.DefaultCreateTime()
		fpc.mutation.SetCreateTime(v)
	}
	if _, ok := fpc.mutation.UpdateTime(); !ok {
		v := floorplan.DefaultUpdateTime()
		fpc.mutation.SetUpdateTime(v)
	}
	if _, ok := fpc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	var (
		err  error
		node *FloorPlan
	)
	if len(fpc.hooks) == 0 {
		node, err = fpc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FloorPlanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fpc.mutation = mutation
			node, err = fpc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(fpc.hooks) - 1; i >= 0; i-- {
			mut = fpc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fpc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (fpc *FloorPlanCreate) SaveX(ctx context.Context) *FloorPlan {
	v, err := fpc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (fpc *FloorPlanCreate) sqlSave(ctx context.Context) (*FloorPlan, error) {
	var (
		fp    = &FloorPlan{config: fpc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: floorplan.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: floorplan.FieldID,
			},
		}
	)
	if value, ok := fpc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: floorplan.FieldCreateTime,
		})
		fp.CreateTime = value
	}
	if value, ok := fpc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: floorplan.FieldUpdateTime,
		})
		fp.UpdateTime = value
	}
	if value, ok := fpc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: floorplan.FieldName,
		})
		fp.Name = value
	}
	if nodes := fpc.mutation.LocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   floorplan.LocationTable,
			Columns: []string{floorplan.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fpc.mutation.ReferencePointIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   floorplan.ReferencePointTable,
			Columns: []string{floorplan.ReferencePointColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: floorplanreferencepoint.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fpc.mutation.ScaleIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   floorplan.ScaleTable,
			Columns: []string{floorplan.ScaleColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: floorplanscale.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fpc.mutation.ImageIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   floorplan.ImageTable,
			Columns: []string{floorplan.ImageColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: file.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, fpc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	fp.ID = int(id)
	return fp, nil
}
