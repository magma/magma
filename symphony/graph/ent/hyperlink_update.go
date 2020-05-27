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
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// HyperlinkUpdate is the builder for updating Hyperlink entities.
type HyperlinkUpdate struct {
	config
	hooks      []Hook
	mutation   *HyperlinkMutation
	predicates []predicate.Hyperlink
}

// Where adds a new predicate for the builder.
func (hu *HyperlinkUpdate) Where(ps ...predicate.Hyperlink) *HyperlinkUpdate {
	hu.predicates = append(hu.predicates, ps...)
	return hu
}

// SetURL sets the url field.
func (hu *HyperlinkUpdate) SetURL(s string) *HyperlinkUpdate {
	hu.mutation.SetURL(s)
	return hu
}

// SetName sets the name field.
func (hu *HyperlinkUpdate) SetName(s string) *HyperlinkUpdate {
	hu.mutation.SetName(s)
	return hu
}

// SetNillableName sets the name field if the given value is not nil.
func (hu *HyperlinkUpdate) SetNillableName(s *string) *HyperlinkUpdate {
	if s != nil {
		hu.SetName(*s)
	}
	return hu
}

// ClearName clears the value of name.
func (hu *HyperlinkUpdate) ClearName() *HyperlinkUpdate {
	hu.mutation.ClearName()
	return hu
}

// SetCategory sets the category field.
func (hu *HyperlinkUpdate) SetCategory(s string) *HyperlinkUpdate {
	hu.mutation.SetCategory(s)
	return hu
}

// SetNillableCategory sets the category field if the given value is not nil.
func (hu *HyperlinkUpdate) SetNillableCategory(s *string) *HyperlinkUpdate {
	if s != nil {
		hu.SetCategory(*s)
	}
	return hu
}

// ClearCategory clears the value of category.
func (hu *HyperlinkUpdate) ClearCategory() *HyperlinkUpdate {
	hu.mutation.ClearCategory()
	return hu
}

// SetEquipmentID sets the equipment edge to Equipment by id.
func (hu *HyperlinkUpdate) SetEquipmentID(id int) *HyperlinkUpdate {
	hu.mutation.SetEquipmentID(id)
	return hu
}

// SetNillableEquipmentID sets the equipment edge to Equipment by id if the given value is not nil.
func (hu *HyperlinkUpdate) SetNillableEquipmentID(id *int) *HyperlinkUpdate {
	if id != nil {
		hu = hu.SetEquipmentID(*id)
	}
	return hu
}

// SetEquipment sets the equipment edge to Equipment.
func (hu *HyperlinkUpdate) SetEquipment(e *Equipment) *HyperlinkUpdate {
	return hu.SetEquipmentID(e.ID)
}

// SetLocationID sets the location edge to Location by id.
func (hu *HyperlinkUpdate) SetLocationID(id int) *HyperlinkUpdate {
	hu.mutation.SetLocationID(id)
	return hu
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (hu *HyperlinkUpdate) SetNillableLocationID(id *int) *HyperlinkUpdate {
	if id != nil {
		hu = hu.SetLocationID(*id)
	}
	return hu
}

// SetLocation sets the location edge to Location.
func (hu *HyperlinkUpdate) SetLocation(l *Location) *HyperlinkUpdate {
	return hu.SetLocationID(l.ID)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (hu *HyperlinkUpdate) SetWorkOrderID(id int) *HyperlinkUpdate {
	hu.mutation.SetWorkOrderID(id)
	return hu
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (hu *HyperlinkUpdate) SetNillableWorkOrderID(id *int) *HyperlinkUpdate {
	if id != nil {
		hu = hu.SetWorkOrderID(*id)
	}
	return hu
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (hu *HyperlinkUpdate) SetWorkOrder(w *WorkOrder) *HyperlinkUpdate {
	return hu.SetWorkOrderID(w.ID)
}

// ClearEquipment clears the equipment edge to Equipment.
func (hu *HyperlinkUpdate) ClearEquipment() *HyperlinkUpdate {
	hu.mutation.ClearEquipment()
	return hu
}

// ClearLocation clears the location edge to Location.
func (hu *HyperlinkUpdate) ClearLocation() *HyperlinkUpdate {
	hu.mutation.ClearLocation()
	return hu
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (hu *HyperlinkUpdate) ClearWorkOrder() *HyperlinkUpdate {
	hu.mutation.ClearWorkOrder()
	return hu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (hu *HyperlinkUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := hu.mutation.UpdateTime(); !ok {
		v := hyperlink.UpdateDefaultUpdateTime()
		hu.mutation.SetUpdateTime(v)
	}

	var (
		err      error
		affected int
	)
	if len(hu.hooks) == 0 {
		affected, err = hu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*HyperlinkMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			hu.mutation = mutation
			affected, err = hu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(hu.hooks) - 1; i >= 0; i-- {
			mut = hu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, hu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (hu *HyperlinkUpdate) SaveX(ctx context.Context) int {
	affected, err := hu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (hu *HyperlinkUpdate) Exec(ctx context.Context) error {
	_, err := hu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (hu *HyperlinkUpdate) ExecX(ctx context.Context) {
	if err := hu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (hu *HyperlinkUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   hyperlink.Table,
			Columns: hyperlink.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: hyperlink.FieldID,
			},
		},
	}
	if ps := hu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := hu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: hyperlink.FieldUpdateTime,
		})
	}
	if value, ok := hu.mutation.URL(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: hyperlink.FieldURL,
		})
	}
	if value, ok := hu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: hyperlink.FieldName,
		})
	}
	if hu.mutation.NameCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: hyperlink.FieldName,
		})
	}
	if value, ok := hu.mutation.Category(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: hyperlink.FieldCategory,
		})
	}
	if hu.mutation.CategoryCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: hyperlink.FieldCategory,
		})
	}
	if hu.mutation.EquipmentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   hyperlink.EquipmentTable,
			Columns: []string{hyperlink.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := hu.mutation.EquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   hyperlink.EquipmentTable,
			Columns: []string{hyperlink.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if hu.mutation.LocationCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   hyperlink.LocationTable,
			Columns: []string{hyperlink.LocationColumn},
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
	if nodes := hu.mutation.LocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   hyperlink.LocationTable,
			Columns: []string{hyperlink.LocationColumn},
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
	if hu.mutation.WorkOrderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   hyperlink.WorkOrderTable,
			Columns: []string{hyperlink.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := hu.mutation.WorkOrderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   hyperlink.WorkOrderTable,
			Columns: []string{hyperlink.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, hu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{hyperlink.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// HyperlinkUpdateOne is the builder for updating a single Hyperlink entity.
type HyperlinkUpdateOne struct {
	config
	hooks    []Hook
	mutation *HyperlinkMutation
}

// SetURL sets the url field.
func (huo *HyperlinkUpdateOne) SetURL(s string) *HyperlinkUpdateOne {
	huo.mutation.SetURL(s)
	return huo
}

// SetName sets the name field.
func (huo *HyperlinkUpdateOne) SetName(s string) *HyperlinkUpdateOne {
	huo.mutation.SetName(s)
	return huo
}

// SetNillableName sets the name field if the given value is not nil.
func (huo *HyperlinkUpdateOne) SetNillableName(s *string) *HyperlinkUpdateOne {
	if s != nil {
		huo.SetName(*s)
	}
	return huo
}

// ClearName clears the value of name.
func (huo *HyperlinkUpdateOne) ClearName() *HyperlinkUpdateOne {
	huo.mutation.ClearName()
	return huo
}

// SetCategory sets the category field.
func (huo *HyperlinkUpdateOne) SetCategory(s string) *HyperlinkUpdateOne {
	huo.mutation.SetCategory(s)
	return huo
}

// SetNillableCategory sets the category field if the given value is not nil.
func (huo *HyperlinkUpdateOne) SetNillableCategory(s *string) *HyperlinkUpdateOne {
	if s != nil {
		huo.SetCategory(*s)
	}
	return huo
}

// ClearCategory clears the value of category.
func (huo *HyperlinkUpdateOne) ClearCategory() *HyperlinkUpdateOne {
	huo.mutation.ClearCategory()
	return huo
}

// SetEquipmentID sets the equipment edge to Equipment by id.
func (huo *HyperlinkUpdateOne) SetEquipmentID(id int) *HyperlinkUpdateOne {
	huo.mutation.SetEquipmentID(id)
	return huo
}

// SetNillableEquipmentID sets the equipment edge to Equipment by id if the given value is not nil.
func (huo *HyperlinkUpdateOne) SetNillableEquipmentID(id *int) *HyperlinkUpdateOne {
	if id != nil {
		huo = huo.SetEquipmentID(*id)
	}
	return huo
}

// SetEquipment sets the equipment edge to Equipment.
func (huo *HyperlinkUpdateOne) SetEquipment(e *Equipment) *HyperlinkUpdateOne {
	return huo.SetEquipmentID(e.ID)
}

// SetLocationID sets the location edge to Location by id.
func (huo *HyperlinkUpdateOne) SetLocationID(id int) *HyperlinkUpdateOne {
	huo.mutation.SetLocationID(id)
	return huo
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (huo *HyperlinkUpdateOne) SetNillableLocationID(id *int) *HyperlinkUpdateOne {
	if id != nil {
		huo = huo.SetLocationID(*id)
	}
	return huo
}

// SetLocation sets the location edge to Location.
func (huo *HyperlinkUpdateOne) SetLocation(l *Location) *HyperlinkUpdateOne {
	return huo.SetLocationID(l.ID)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (huo *HyperlinkUpdateOne) SetWorkOrderID(id int) *HyperlinkUpdateOne {
	huo.mutation.SetWorkOrderID(id)
	return huo
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (huo *HyperlinkUpdateOne) SetNillableWorkOrderID(id *int) *HyperlinkUpdateOne {
	if id != nil {
		huo = huo.SetWorkOrderID(*id)
	}
	return huo
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (huo *HyperlinkUpdateOne) SetWorkOrder(w *WorkOrder) *HyperlinkUpdateOne {
	return huo.SetWorkOrderID(w.ID)
}

// ClearEquipment clears the equipment edge to Equipment.
func (huo *HyperlinkUpdateOne) ClearEquipment() *HyperlinkUpdateOne {
	huo.mutation.ClearEquipment()
	return huo
}

// ClearLocation clears the location edge to Location.
func (huo *HyperlinkUpdateOne) ClearLocation() *HyperlinkUpdateOne {
	huo.mutation.ClearLocation()
	return huo
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (huo *HyperlinkUpdateOne) ClearWorkOrder() *HyperlinkUpdateOne {
	huo.mutation.ClearWorkOrder()
	return huo
}

// Save executes the query and returns the updated entity.
func (huo *HyperlinkUpdateOne) Save(ctx context.Context) (*Hyperlink, error) {
	if _, ok := huo.mutation.UpdateTime(); !ok {
		v := hyperlink.UpdateDefaultUpdateTime()
		huo.mutation.SetUpdateTime(v)
	}

	var (
		err  error
		node *Hyperlink
	)
	if len(huo.hooks) == 0 {
		node, err = huo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*HyperlinkMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			huo.mutation = mutation
			node, err = huo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(huo.hooks) - 1; i >= 0; i-- {
			mut = huo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, huo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (huo *HyperlinkUpdateOne) SaveX(ctx context.Context) *Hyperlink {
	h, err := huo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return h
}

// Exec executes the query on the entity.
func (huo *HyperlinkUpdateOne) Exec(ctx context.Context) error {
	_, err := huo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (huo *HyperlinkUpdateOne) ExecX(ctx context.Context) {
	if err := huo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (huo *HyperlinkUpdateOne) sqlSave(ctx context.Context) (h *Hyperlink, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   hyperlink.Table,
			Columns: hyperlink.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: hyperlink.FieldID,
			},
		},
	}
	id, ok := huo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing Hyperlink.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := huo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: hyperlink.FieldUpdateTime,
		})
	}
	if value, ok := huo.mutation.URL(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: hyperlink.FieldURL,
		})
	}
	if value, ok := huo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: hyperlink.FieldName,
		})
	}
	if huo.mutation.NameCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: hyperlink.FieldName,
		})
	}
	if value, ok := huo.mutation.Category(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: hyperlink.FieldCategory,
		})
	}
	if huo.mutation.CategoryCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: hyperlink.FieldCategory,
		})
	}
	if huo.mutation.EquipmentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   hyperlink.EquipmentTable,
			Columns: []string{hyperlink.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := huo.mutation.EquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   hyperlink.EquipmentTable,
			Columns: []string{hyperlink.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if huo.mutation.LocationCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   hyperlink.LocationTable,
			Columns: []string{hyperlink.LocationColumn},
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
	if nodes := huo.mutation.LocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   hyperlink.LocationTable,
			Columns: []string{hyperlink.LocationColumn},
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
	if huo.mutation.WorkOrderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   hyperlink.WorkOrderTable,
			Columns: []string{hyperlink.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := huo.mutation.WorkOrderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   hyperlink.WorkOrderTable,
			Columns: []string{hyperlink.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	h = &Hyperlink{config: huo.config}
	_spec.Assign = h.assignValues
	_spec.ScanValues = h.scanValues()
	if err = sqlgraph.UpdateNode(ctx, huo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{hyperlink.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return h, nil
}
