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
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/floorplan"
	"github.com/facebookincubator/symphony/graph/ent/floorplanreferencepoint"
	"github.com/facebookincubator/symphony/graph/ent/floorplanscale"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// FloorPlanUpdate is the builder for updating FloorPlan entities.
type FloorPlanUpdate struct {
	config
	hooks      []Hook
	mutation   *FloorPlanMutation
	predicates []predicate.FloorPlan
}

// Where adds a new predicate for the builder.
func (fpu *FloorPlanUpdate) Where(ps ...predicate.FloorPlan) *FloorPlanUpdate {
	fpu.predicates = append(fpu.predicates, ps...)
	return fpu
}

// SetName sets the name field.
func (fpu *FloorPlanUpdate) SetName(s string) *FloorPlanUpdate {
	fpu.mutation.SetName(s)
	return fpu
}

// SetLocationID sets the location edge to Location by id.
func (fpu *FloorPlanUpdate) SetLocationID(id int) *FloorPlanUpdate {
	fpu.mutation.SetLocationID(id)
	return fpu
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (fpu *FloorPlanUpdate) SetNillableLocationID(id *int) *FloorPlanUpdate {
	if id != nil {
		fpu = fpu.SetLocationID(*id)
	}
	return fpu
}

// SetLocation sets the location edge to Location.
func (fpu *FloorPlanUpdate) SetLocation(l *Location) *FloorPlanUpdate {
	return fpu.SetLocationID(l.ID)
}

// SetReferencePointID sets the reference_point edge to FloorPlanReferencePoint by id.
func (fpu *FloorPlanUpdate) SetReferencePointID(id int) *FloorPlanUpdate {
	fpu.mutation.SetReferencePointID(id)
	return fpu
}

// SetNillableReferencePointID sets the reference_point edge to FloorPlanReferencePoint by id if the given value is not nil.
func (fpu *FloorPlanUpdate) SetNillableReferencePointID(id *int) *FloorPlanUpdate {
	if id != nil {
		fpu = fpu.SetReferencePointID(*id)
	}
	return fpu
}

// SetReferencePoint sets the reference_point edge to FloorPlanReferencePoint.
func (fpu *FloorPlanUpdate) SetReferencePoint(f *FloorPlanReferencePoint) *FloorPlanUpdate {
	return fpu.SetReferencePointID(f.ID)
}

// SetScaleID sets the scale edge to FloorPlanScale by id.
func (fpu *FloorPlanUpdate) SetScaleID(id int) *FloorPlanUpdate {
	fpu.mutation.SetScaleID(id)
	return fpu
}

// SetNillableScaleID sets the scale edge to FloorPlanScale by id if the given value is not nil.
func (fpu *FloorPlanUpdate) SetNillableScaleID(id *int) *FloorPlanUpdate {
	if id != nil {
		fpu = fpu.SetScaleID(*id)
	}
	return fpu
}

// SetScale sets the scale edge to FloorPlanScale.
func (fpu *FloorPlanUpdate) SetScale(f *FloorPlanScale) *FloorPlanUpdate {
	return fpu.SetScaleID(f.ID)
}

// SetImageID sets the image edge to File by id.
func (fpu *FloorPlanUpdate) SetImageID(id int) *FloorPlanUpdate {
	fpu.mutation.SetImageID(id)
	return fpu
}

// SetNillableImageID sets the image edge to File by id if the given value is not nil.
func (fpu *FloorPlanUpdate) SetNillableImageID(id *int) *FloorPlanUpdate {
	if id != nil {
		fpu = fpu.SetImageID(*id)
	}
	return fpu
}

// SetImage sets the image edge to File.
func (fpu *FloorPlanUpdate) SetImage(f *File) *FloorPlanUpdate {
	return fpu.SetImageID(f.ID)
}

// ClearLocation clears the location edge to Location.
func (fpu *FloorPlanUpdate) ClearLocation() *FloorPlanUpdate {
	fpu.mutation.ClearLocation()
	return fpu
}

// ClearReferencePoint clears the reference_point edge to FloorPlanReferencePoint.
func (fpu *FloorPlanUpdate) ClearReferencePoint() *FloorPlanUpdate {
	fpu.mutation.ClearReferencePoint()
	return fpu
}

// ClearScale clears the scale edge to FloorPlanScale.
func (fpu *FloorPlanUpdate) ClearScale() *FloorPlanUpdate {
	fpu.mutation.ClearScale()
	return fpu
}

// ClearImage clears the image edge to File.
func (fpu *FloorPlanUpdate) ClearImage() *FloorPlanUpdate {
	fpu.mutation.ClearImage()
	return fpu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (fpu *FloorPlanUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := fpu.mutation.UpdateTime(); !ok {
		v := floorplan.UpdateDefaultUpdateTime()
		fpu.mutation.SetUpdateTime(v)
	}

	var (
		err      error
		affected int
	)
	if len(fpu.hooks) == 0 {
		affected, err = fpu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FloorPlanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fpu.mutation = mutation
			affected, err = fpu.sqlSave(ctx)
			return affected, err
		})
		for i := len(fpu.hooks) - 1; i >= 0; i-- {
			mut = fpu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fpu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (fpu *FloorPlanUpdate) SaveX(ctx context.Context) int {
	affected, err := fpu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fpu *FloorPlanUpdate) Exec(ctx context.Context) error {
	_, err := fpu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fpu *FloorPlanUpdate) ExecX(ctx context.Context) {
	if err := fpu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fpu *FloorPlanUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   floorplan.Table,
			Columns: floorplan.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: floorplan.FieldID,
			},
		},
	}
	if ps := fpu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fpu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: floorplan.FieldUpdateTime,
		})
	}
	if value, ok := fpu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: floorplan.FieldName,
		})
	}
	if fpu.mutation.LocationCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fpu.mutation.LocationIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fpu.mutation.ReferencePointCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fpu.mutation.ReferencePointIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fpu.mutation.ScaleCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fpu.mutation.ScaleIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fpu.mutation.ImageCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fpu.mutation.ImageIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fpu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{floorplan.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// FloorPlanUpdateOne is the builder for updating a single FloorPlan entity.
type FloorPlanUpdateOne struct {
	config
	hooks    []Hook
	mutation *FloorPlanMutation
}

// SetName sets the name field.
func (fpuo *FloorPlanUpdateOne) SetName(s string) *FloorPlanUpdateOne {
	fpuo.mutation.SetName(s)
	return fpuo
}

// SetLocationID sets the location edge to Location by id.
func (fpuo *FloorPlanUpdateOne) SetLocationID(id int) *FloorPlanUpdateOne {
	fpuo.mutation.SetLocationID(id)
	return fpuo
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (fpuo *FloorPlanUpdateOne) SetNillableLocationID(id *int) *FloorPlanUpdateOne {
	if id != nil {
		fpuo = fpuo.SetLocationID(*id)
	}
	return fpuo
}

// SetLocation sets the location edge to Location.
func (fpuo *FloorPlanUpdateOne) SetLocation(l *Location) *FloorPlanUpdateOne {
	return fpuo.SetLocationID(l.ID)
}

// SetReferencePointID sets the reference_point edge to FloorPlanReferencePoint by id.
func (fpuo *FloorPlanUpdateOne) SetReferencePointID(id int) *FloorPlanUpdateOne {
	fpuo.mutation.SetReferencePointID(id)
	return fpuo
}

// SetNillableReferencePointID sets the reference_point edge to FloorPlanReferencePoint by id if the given value is not nil.
func (fpuo *FloorPlanUpdateOne) SetNillableReferencePointID(id *int) *FloorPlanUpdateOne {
	if id != nil {
		fpuo = fpuo.SetReferencePointID(*id)
	}
	return fpuo
}

// SetReferencePoint sets the reference_point edge to FloorPlanReferencePoint.
func (fpuo *FloorPlanUpdateOne) SetReferencePoint(f *FloorPlanReferencePoint) *FloorPlanUpdateOne {
	return fpuo.SetReferencePointID(f.ID)
}

// SetScaleID sets the scale edge to FloorPlanScale by id.
func (fpuo *FloorPlanUpdateOne) SetScaleID(id int) *FloorPlanUpdateOne {
	fpuo.mutation.SetScaleID(id)
	return fpuo
}

// SetNillableScaleID sets the scale edge to FloorPlanScale by id if the given value is not nil.
func (fpuo *FloorPlanUpdateOne) SetNillableScaleID(id *int) *FloorPlanUpdateOne {
	if id != nil {
		fpuo = fpuo.SetScaleID(*id)
	}
	return fpuo
}

// SetScale sets the scale edge to FloorPlanScale.
func (fpuo *FloorPlanUpdateOne) SetScale(f *FloorPlanScale) *FloorPlanUpdateOne {
	return fpuo.SetScaleID(f.ID)
}

// SetImageID sets the image edge to File by id.
func (fpuo *FloorPlanUpdateOne) SetImageID(id int) *FloorPlanUpdateOne {
	fpuo.mutation.SetImageID(id)
	return fpuo
}

// SetNillableImageID sets the image edge to File by id if the given value is not nil.
func (fpuo *FloorPlanUpdateOne) SetNillableImageID(id *int) *FloorPlanUpdateOne {
	if id != nil {
		fpuo = fpuo.SetImageID(*id)
	}
	return fpuo
}

// SetImage sets the image edge to File.
func (fpuo *FloorPlanUpdateOne) SetImage(f *File) *FloorPlanUpdateOne {
	return fpuo.SetImageID(f.ID)
}

// ClearLocation clears the location edge to Location.
func (fpuo *FloorPlanUpdateOne) ClearLocation() *FloorPlanUpdateOne {
	fpuo.mutation.ClearLocation()
	return fpuo
}

// ClearReferencePoint clears the reference_point edge to FloorPlanReferencePoint.
func (fpuo *FloorPlanUpdateOne) ClearReferencePoint() *FloorPlanUpdateOne {
	fpuo.mutation.ClearReferencePoint()
	return fpuo
}

// ClearScale clears the scale edge to FloorPlanScale.
func (fpuo *FloorPlanUpdateOne) ClearScale() *FloorPlanUpdateOne {
	fpuo.mutation.ClearScale()
	return fpuo
}

// ClearImage clears the image edge to File.
func (fpuo *FloorPlanUpdateOne) ClearImage() *FloorPlanUpdateOne {
	fpuo.mutation.ClearImage()
	return fpuo
}

// Save executes the query and returns the updated entity.
func (fpuo *FloorPlanUpdateOne) Save(ctx context.Context) (*FloorPlan, error) {
	if _, ok := fpuo.mutation.UpdateTime(); !ok {
		v := floorplan.UpdateDefaultUpdateTime()
		fpuo.mutation.SetUpdateTime(v)
	}

	var (
		err  error
		node *FloorPlan
	)
	if len(fpuo.hooks) == 0 {
		node, err = fpuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FloorPlanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fpuo.mutation = mutation
			node, err = fpuo.sqlSave(ctx)
			return node, err
		})
		for i := len(fpuo.hooks) - 1; i >= 0; i-- {
			mut = fpuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fpuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (fpuo *FloorPlanUpdateOne) SaveX(ctx context.Context) *FloorPlan {
	fp, err := fpuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return fp
}

// Exec executes the query on the entity.
func (fpuo *FloorPlanUpdateOne) Exec(ctx context.Context) error {
	_, err := fpuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fpuo *FloorPlanUpdateOne) ExecX(ctx context.Context) {
	if err := fpuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fpuo *FloorPlanUpdateOne) sqlSave(ctx context.Context) (fp *FloorPlan, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   floorplan.Table,
			Columns: floorplan.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: floorplan.FieldID,
			},
		},
	}
	id, ok := fpuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing FloorPlan.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := fpuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: floorplan.FieldUpdateTime,
		})
	}
	if value, ok := fpuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: floorplan.FieldName,
		})
	}
	if fpuo.mutation.LocationCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fpuo.mutation.LocationIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fpuo.mutation.ReferencePointCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fpuo.mutation.ReferencePointIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fpuo.mutation.ScaleCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fpuo.mutation.ScaleIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fpuo.mutation.ImageCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fpuo.mutation.ImageIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	fp = &FloorPlan{config: fpuo.config}
	_spec.Assign = fp.assignValues
	_spec.ScanValues = fp.scanValues()
	if err = sqlgraph.UpdateNode(ctx, fpuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{floorplan.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return fp, nil
}
