// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// EquipmentUpdate is the builder for updating Equipment entities.
type EquipmentUpdate struct {
	config
	hooks      []Hook
	mutation   *EquipmentMutation
	predicates []predicate.Equipment
}

// Where adds a new predicate for the builder.
func (eu *EquipmentUpdate) Where(ps ...predicate.Equipment) *EquipmentUpdate {
	eu.predicates = append(eu.predicates, ps...)
	return eu
}

// SetName sets the name field.
func (eu *EquipmentUpdate) SetName(s string) *EquipmentUpdate {
	eu.mutation.SetName(s)
	return eu
}

// SetFutureState sets the future_state field.
func (eu *EquipmentUpdate) SetFutureState(s string) *EquipmentUpdate {
	eu.mutation.SetFutureState(s)
	return eu
}

// SetNillableFutureState sets the future_state field if the given value is not nil.
func (eu *EquipmentUpdate) SetNillableFutureState(s *string) *EquipmentUpdate {
	if s != nil {
		eu.SetFutureState(*s)
	}
	return eu
}

// ClearFutureState clears the value of future_state.
func (eu *EquipmentUpdate) ClearFutureState() *EquipmentUpdate {
	eu.mutation.ClearFutureState()
	return eu
}

// SetDeviceID sets the device_id field.
func (eu *EquipmentUpdate) SetDeviceID(s string) *EquipmentUpdate {
	eu.mutation.SetDeviceID(s)
	return eu
}

// SetNillableDeviceID sets the device_id field if the given value is not nil.
func (eu *EquipmentUpdate) SetNillableDeviceID(s *string) *EquipmentUpdate {
	if s != nil {
		eu.SetDeviceID(*s)
	}
	return eu
}

// ClearDeviceID clears the value of device_id.
func (eu *EquipmentUpdate) ClearDeviceID() *EquipmentUpdate {
	eu.mutation.ClearDeviceID()
	return eu
}

// SetExternalID sets the external_id field.
func (eu *EquipmentUpdate) SetExternalID(s string) *EquipmentUpdate {
	eu.mutation.SetExternalID(s)
	return eu
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (eu *EquipmentUpdate) SetNillableExternalID(s *string) *EquipmentUpdate {
	if s != nil {
		eu.SetExternalID(*s)
	}
	return eu
}

// ClearExternalID clears the value of external_id.
func (eu *EquipmentUpdate) ClearExternalID() *EquipmentUpdate {
	eu.mutation.ClearExternalID()
	return eu
}

// SetTypeID sets the type edge to EquipmentType by id.
func (eu *EquipmentUpdate) SetTypeID(id int) *EquipmentUpdate {
	eu.mutation.SetTypeID(id)
	return eu
}

// SetType sets the type edge to EquipmentType.
func (eu *EquipmentUpdate) SetType(e *EquipmentType) *EquipmentUpdate {
	return eu.SetTypeID(e.ID)
}

// SetLocationID sets the location edge to Location by id.
func (eu *EquipmentUpdate) SetLocationID(id int) *EquipmentUpdate {
	eu.mutation.SetLocationID(id)
	return eu
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (eu *EquipmentUpdate) SetNillableLocationID(id *int) *EquipmentUpdate {
	if id != nil {
		eu = eu.SetLocationID(*id)
	}
	return eu
}

// SetLocation sets the location edge to Location.
func (eu *EquipmentUpdate) SetLocation(l *Location) *EquipmentUpdate {
	return eu.SetLocationID(l.ID)
}

// SetParentPositionID sets the parent_position edge to EquipmentPosition by id.
func (eu *EquipmentUpdate) SetParentPositionID(id int) *EquipmentUpdate {
	eu.mutation.SetParentPositionID(id)
	return eu
}

// SetNillableParentPositionID sets the parent_position edge to EquipmentPosition by id if the given value is not nil.
func (eu *EquipmentUpdate) SetNillableParentPositionID(id *int) *EquipmentUpdate {
	if id != nil {
		eu = eu.SetParentPositionID(*id)
	}
	return eu
}

// SetParentPosition sets the parent_position edge to EquipmentPosition.
func (eu *EquipmentUpdate) SetParentPosition(e *EquipmentPosition) *EquipmentUpdate {
	return eu.SetParentPositionID(e.ID)
}

// AddPositionIDs adds the positions edge to EquipmentPosition by ids.
func (eu *EquipmentUpdate) AddPositionIDs(ids ...int) *EquipmentUpdate {
	eu.mutation.AddPositionIDs(ids...)
	return eu
}

// AddPositions adds the positions edges to EquipmentPosition.
func (eu *EquipmentUpdate) AddPositions(e ...*EquipmentPosition) *EquipmentUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eu.AddPositionIDs(ids...)
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (eu *EquipmentUpdate) AddPortIDs(ids ...int) *EquipmentUpdate {
	eu.mutation.AddPortIDs(ids...)
	return eu
}

// AddPorts adds the ports edges to EquipmentPort.
func (eu *EquipmentUpdate) AddPorts(e ...*EquipmentPort) *EquipmentUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eu.AddPortIDs(ids...)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (eu *EquipmentUpdate) SetWorkOrderID(id int) *EquipmentUpdate {
	eu.mutation.SetWorkOrderID(id)
	return eu
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (eu *EquipmentUpdate) SetNillableWorkOrderID(id *int) *EquipmentUpdate {
	if id != nil {
		eu = eu.SetWorkOrderID(*id)
	}
	return eu
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (eu *EquipmentUpdate) SetWorkOrder(w *WorkOrder) *EquipmentUpdate {
	return eu.SetWorkOrderID(w.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (eu *EquipmentUpdate) AddPropertyIDs(ids ...int) *EquipmentUpdate {
	eu.mutation.AddPropertyIDs(ids...)
	return eu
}

// AddProperties adds the properties edges to Property.
func (eu *EquipmentUpdate) AddProperties(p ...*Property) *EquipmentUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eu.AddPropertyIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (eu *EquipmentUpdate) AddFileIDs(ids ...int) *EquipmentUpdate {
	eu.mutation.AddFileIDs(ids...)
	return eu
}

// AddFiles adds the files edges to File.
func (eu *EquipmentUpdate) AddFiles(f ...*File) *EquipmentUpdate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return eu.AddFileIDs(ids...)
}

// AddHyperlinkIDs adds the hyperlinks edge to Hyperlink by ids.
func (eu *EquipmentUpdate) AddHyperlinkIDs(ids ...int) *EquipmentUpdate {
	eu.mutation.AddHyperlinkIDs(ids...)
	return eu
}

// AddHyperlinks adds the hyperlinks edges to Hyperlink.
func (eu *EquipmentUpdate) AddHyperlinks(h ...*Hyperlink) *EquipmentUpdate {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return eu.AddHyperlinkIDs(ids...)
}

// ClearType clears the type edge to EquipmentType.
func (eu *EquipmentUpdate) ClearType() *EquipmentUpdate {
	eu.mutation.ClearType()
	return eu
}

// ClearLocation clears the location edge to Location.
func (eu *EquipmentUpdate) ClearLocation() *EquipmentUpdate {
	eu.mutation.ClearLocation()
	return eu
}

// ClearParentPosition clears the parent_position edge to EquipmentPosition.
func (eu *EquipmentUpdate) ClearParentPosition() *EquipmentUpdate {
	eu.mutation.ClearParentPosition()
	return eu
}

// RemovePositionIDs removes the positions edge to EquipmentPosition by ids.
func (eu *EquipmentUpdate) RemovePositionIDs(ids ...int) *EquipmentUpdate {
	eu.mutation.RemovePositionIDs(ids...)
	return eu
}

// RemovePositions removes positions edges to EquipmentPosition.
func (eu *EquipmentUpdate) RemovePositions(e ...*EquipmentPosition) *EquipmentUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eu.RemovePositionIDs(ids...)
}

// RemovePortIDs removes the ports edge to EquipmentPort by ids.
func (eu *EquipmentUpdate) RemovePortIDs(ids ...int) *EquipmentUpdate {
	eu.mutation.RemovePortIDs(ids...)
	return eu
}

// RemovePorts removes ports edges to EquipmentPort.
func (eu *EquipmentUpdate) RemovePorts(e ...*EquipmentPort) *EquipmentUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eu.RemovePortIDs(ids...)
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (eu *EquipmentUpdate) ClearWorkOrder() *EquipmentUpdate {
	eu.mutation.ClearWorkOrder()
	return eu
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (eu *EquipmentUpdate) RemovePropertyIDs(ids ...int) *EquipmentUpdate {
	eu.mutation.RemovePropertyIDs(ids...)
	return eu
}

// RemoveProperties removes properties edges to Property.
func (eu *EquipmentUpdate) RemoveProperties(p ...*Property) *EquipmentUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eu.RemovePropertyIDs(ids...)
}

// RemoveFileIDs removes the files edge to File by ids.
func (eu *EquipmentUpdate) RemoveFileIDs(ids ...int) *EquipmentUpdate {
	eu.mutation.RemoveFileIDs(ids...)
	return eu
}

// RemoveFiles removes files edges to File.
func (eu *EquipmentUpdate) RemoveFiles(f ...*File) *EquipmentUpdate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return eu.RemoveFileIDs(ids...)
}

// RemoveHyperlinkIDs removes the hyperlinks edge to Hyperlink by ids.
func (eu *EquipmentUpdate) RemoveHyperlinkIDs(ids ...int) *EquipmentUpdate {
	eu.mutation.RemoveHyperlinkIDs(ids...)
	return eu
}

// RemoveHyperlinks removes hyperlinks edges to Hyperlink.
func (eu *EquipmentUpdate) RemoveHyperlinks(h ...*Hyperlink) *EquipmentUpdate {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return eu.RemoveHyperlinkIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (eu *EquipmentUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := eu.mutation.UpdateTime(); !ok {
		v := equipment.UpdateDefaultUpdateTime()
		eu.mutation.SetUpdateTime(v)
	}
	if v, ok := eu.mutation.Name(); ok {
		if err := equipment.NameValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if v, ok := eu.mutation.DeviceID(); ok {
		if err := equipment.DeviceIDValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"device_id\": %v", err)
		}
	}

	if _, ok := eu.mutation.TypeID(); eu.mutation.TypeCleared() && !ok {
		return 0, errors.New("ent: clearing a unique edge \"type\"")
	}

	var (
		err      error
		affected int
	)
	if len(eu.hooks) == 0 {
		affected, err = eu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			eu.mutation = mutation
			affected, err = eu.sqlSave(ctx)
			return affected, err
		})
		for i := len(eu.hooks); i > 0; i-- {
			mut = eu.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, eu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (eu *EquipmentUpdate) SaveX(ctx context.Context) int {
	affected, err := eu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (eu *EquipmentUpdate) Exec(ctx context.Context) error {
	_, err := eu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eu *EquipmentUpdate) ExecX(ctx context.Context) {
	if err := eu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (eu *EquipmentUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipment.Table,
			Columns: equipment.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipment.FieldID,
			},
		},
	}
	if ps := eu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := eu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipment.FieldUpdateTime,
		})
	}
	if value, ok := eu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipment.FieldName,
		})
	}
	if value, ok := eu.mutation.FutureState(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipment.FieldFutureState,
		})
	}
	if eu.mutation.FutureStateCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipment.FieldFutureState,
		})
	}
	if value, ok := eu.mutation.DeviceID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipment.FieldDeviceID,
		})
	}
	if eu.mutation.DeviceIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipment.FieldDeviceID,
		})
	}
	if value, ok := eu.mutation.ExternalID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipment.FieldExternalID,
		})
	}
	if eu.mutation.ExternalIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipment.FieldExternalID,
		})
	}
	if eu.mutation.TypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipment.TypeTable,
			Columns: []string{equipment.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmenttype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.TypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipment.TypeTable,
			Columns: []string{equipment.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmenttype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if eu.mutation.LocationCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipment.LocationTable,
			Columns: []string{equipment.LocationColumn},
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
	if nodes := eu.mutation.LocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipment.LocationTable,
			Columns: []string{equipment.LocationColumn},
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
	if eu.mutation.ParentPositionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   equipment.ParentPositionTable,
			Columns: []string{equipment.ParentPositionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentposition.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.ParentPositionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   equipment.ParentPositionTable,
			Columns: []string{equipment.ParentPositionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentposition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := eu.mutation.RemovedPositionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.PositionsTable,
			Columns: []string{equipment.PositionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentposition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.PositionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.PositionsTable,
			Columns: []string{equipment.PositionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentposition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := eu.mutation.RemovedPortsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.PortsTable,
			Columns: []string{equipment.PortsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentport.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.PortsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.PortsTable,
			Columns: []string{equipment.PortsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentport.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if eu.mutation.WorkOrderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipment.WorkOrderTable,
			Columns: []string{equipment.WorkOrderColumn},
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
	if nodes := eu.mutation.WorkOrderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipment.WorkOrderTable,
			Columns: []string{equipment.WorkOrderColumn},
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
	if nodes := eu.mutation.RemovedPropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.PropertiesTable,
			Columns: []string{equipment.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.PropertiesTable,
			Columns: []string{equipment.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := eu.mutation.RemovedFilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.FilesTable,
			Columns: []string{equipment.FilesColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.FilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.FilesTable,
			Columns: []string{equipment.FilesColumn},
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
	if nodes := eu.mutation.RemovedHyperlinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.HyperlinksTable,
			Columns: []string{equipment.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: hyperlink.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.HyperlinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.HyperlinksTable,
			Columns: []string{equipment.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: hyperlink.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, eu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{equipment.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// EquipmentUpdateOne is the builder for updating a single Equipment entity.
type EquipmentUpdateOne struct {
	config
	hooks    []Hook
	mutation *EquipmentMutation
}

// SetName sets the name field.
func (euo *EquipmentUpdateOne) SetName(s string) *EquipmentUpdateOne {
	euo.mutation.SetName(s)
	return euo
}

// SetFutureState sets the future_state field.
func (euo *EquipmentUpdateOne) SetFutureState(s string) *EquipmentUpdateOne {
	euo.mutation.SetFutureState(s)
	return euo
}

// SetNillableFutureState sets the future_state field if the given value is not nil.
func (euo *EquipmentUpdateOne) SetNillableFutureState(s *string) *EquipmentUpdateOne {
	if s != nil {
		euo.SetFutureState(*s)
	}
	return euo
}

// ClearFutureState clears the value of future_state.
func (euo *EquipmentUpdateOne) ClearFutureState() *EquipmentUpdateOne {
	euo.mutation.ClearFutureState()
	return euo
}

// SetDeviceID sets the device_id field.
func (euo *EquipmentUpdateOne) SetDeviceID(s string) *EquipmentUpdateOne {
	euo.mutation.SetDeviceID(s)
	return euo
}

// SetNillableDeviceID sets the device_id field if the given value is not nil.
func (euo *EquipmentUpdateOne) SetNillableDeviceID(s *string) *EquipmentUpdateOne {
	if s != nil {
		euo.SetDeviceID(*s)
	}
	return euo
}

// ClearDeviceID clears the value of device_id.
func (euo *EquipmentUpdateOne) ClearDeviceID() *EquipmentUpdateOne {
	euo.mutation.ClearDeviceID()
	return euo
}

// SetExternalID sets the external_id field.
func (euo *EquipmentUpdateOne) SetExternalID(s string) *EquipmentUpdateOne {
	euo.mutation.SetExternalID(s)
	return euo
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (euo *EquipmentUpdateOne) SetNillableExternalID(s *string) *EquipmentUpdateOne {
	if s != nil {
		euo.SetExternalID(*s)
	}
	return euo
}

// ClearExternalID clears the value of external_id.
func (euo *EquipmentUpdateOne) ClearExternalID() *EquipmentUpdateOne {
	euo.mutation.ClearExternalID()
	return euo
}

// SetTypeID sets the type edge to EquipmentType by id.
func (euo *EquipmentUpdateOne) SetTypeID(id int) *EquipmentUpdateOne {
	euo.mutation.SetTypeID(id)
	return euo
}

// SetType sets the type edge to EquipmentType.
func (euo *EquipmentUpdateOne) SetType(e *EquipmentType) *EquipmentUpdateOne {
	return euo.SetTypeID(e.ID)
}

// SetLocationID sets the location edge to Location by id.
func (euo *EquipmentUpdateOne) SetLocationID(id int) *EquipmentUpdateOne {
	euo.mutation.SetLocationID(id)
	return euo
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (euo *EquipmentUpdateOne) SetNillableLocationID(id *int) *EquipmentUpdateOne {
	if id != nil {
		euo = euo.SetLocationID(*id)
	}
	return euo
}

// SetLocation sets the location edge to Location.
func (euo *EquipmentUpdateOne) SetLocation(l *Location) *EquipmentUpdateOne {
	return euo.SetLocationID(l.ID)
}

// SetParentPositionID sets the parent_position edge to EquipmentPosition by id.
func (euo *EquipmentUpdateOne) SetParentPositionID(id int) *EquipmentUpdateOne {
	euo.mutation.SetParentPositionID(id)
	return euo
}

// SetNillableParentPositionID sets the parent_position edge to EquipmentPosition by id if the given value is not nil.
func (euo *EquipmentUpdateOne) SetNillableParentPositionID(id *int) *EquipmentUpdateOne {
	if id != nil {
		euo = euo.SetParentPositionID(*id)
	}
	return euo
}

// SetParentPosition sets the parent_position edge to EquipmentPosition.
func (euo *EquipmentUpdateOne) SetParentPosition(e *EquipmentPosition) *EquipmentUpdateOne {
	return euo.SetParentPositionID(e.ID)
}

// AddPositionIDs adds the positions edge to EquipmentPosition by ids.
func (euo *EquipmentUpdateOne) AddPositionIDs(ids ...int) *EquipmentUpdateOne {
	euo.mutation.AddPositionIDs(ids...)
	return euo
}

// AddPositions adds the positions edges to EquipmentPosition.
func (euo *EquipmentUpdateOne) AddPositions(e ...*EquipmentPosition) *EquipmentUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return euo.AddPositionIDs(ids...)
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (euo *EquipmentUpdateOne) AddPortIDs(ids ...int) *EquipmentUpdateOne {
	euo.mutation.AddPortIDs(ids...)
	return euo
}

// AddPorts adds the ports edges to EquipmentPort.
func (euo *EquipmentUpdateOne) AddPorts(e ...*EquipmentPort) *EquipmentUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return euo.AddPortIDs(ids...)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (euo *EquipmentUpdateOne) SetWorkOrderID(id int) *EquipmentUpdateOne {
	euo.mutation.SetWorkOrderID(id)
	return euo
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (euo *EquipmentUpdateOne) SetNillableWorkOrderID(id *int) *EquipmentUpdateOne {
	if id != nil {
		euo = euo.SetWorkOrderID(*id)
	}
	return euo
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (euo *EquipmentUpdateOne) SetWorkOrder(w *WorkOrder) *EquipmentUpdateOne {
	return euo.SetWorkOrderID(w.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (euo *EquipmentUpdateOne) AddPropertyIDs(ids ...int) *EquipmentUpdateOne {
	euo.mutation.AddPropertyIDs(ids...)
	return euo
}

// AddProperties adds the properties edges to Property.
func (euo *EquipmentUpdateOne) AddProperties(p ...*Property) *EquipmentUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return euo.AddPropertyIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (euo *EquipmentUpdateOne) AddFileIDs(ids ...int) *EquipmentUpdateOne {
	euo.mutation.AddFileIDs(ids...)
	return euo
}

// AddFiles adds the files edges to File.
func (euo *EquipmentUpdateOne) AddFiles(f ...*File) *EquipmentUpdateOne {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return euo.AddFileIDs(ids...)
}

// AddHyperlinkIDs adds the hyperlinks edge to Hyperlink by ids.
func (euo *EquipmentUpdateOne) AddHyperlinkIDs(ids ...int) *EquipmentUpdateOne {
	euo.mutation.AddHyperlinkIDs(ids...)
	return euo
}

// AddHyperlinks adds the hyperlinks edges to Hyperlink.
func (euo *EquipmentUpdateOne) AddHyperlinks(h ...*Hyperlink) *EquipmentUpdateOne {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return euo.AddHyperlinkIDs(ids...)
}

// ClearType clears the type edge to EquipmentType.
func (euo *EquipmentUpdateOne) ClearType() *EquipmentUpdateOne {
	euo.mutation.ClearType()
	return euo
}

// ClearLocation clears the location edge to Location.
func (euo *EquipmentUpdateOne) ClearLocation() *EquipmentUpdateOne {
	euo.mutation.ClearLocation()
	return euo
}

// ClearParentPosition clears the parent_position edge to EquipmentPosition.
func (euo *EquipmentUpdateOne) ClearParentPosition() *EquipmentUpdateOne {
	euo.mutation.ClearParentPosition()
	return euo
}

// RemovePositionIDs removes the positions edge to EquipmentPosition by ids.
func (euo *EquipmentUpdateOne) RemovePositionIDs(ids ...int) *EquipmentUpdateOne {
	euo.mutation.RemovePositionIDs(ids...)
	return euo
}

// RemovePositions removes positions edges to EquipmentPosition.
func (euo *EquipmentUpdateOne) RemovePositions(e ...*EquipmentPosition) *EquipmentUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return euo.RemovePositionIDs(ids...)
}

// RemovePortIDs removes the ports edge to EquipmentPort by ids.
func (euo *EquipmentUpdateOne) RemovePortIDs(ids ...int) *EquipmentUpdateOne {
	euo.mutation.RemovePortIDs(ids...)
	return euo
}

// RemovePorts removes ports edges to EquipmentPort.
func (euo *EquipmentUpdateOne) RemovePorts(e ...*EquipmentPort) *EquipmentUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return euo.RemovePortIDs(ids...)
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (euo *EquipmentUpdateOne) ClearWorkOrder() *EquipmentUpdateOne {
	euo.mutation.ClearWorkOrder()
	return euo
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (euo *EquipmentUpdateOne) RemovePropertyIDs(ids ...int) *EquipmentUpdateOne {
	euo.mutation.RemovePropertyIDs(ids...)
	return euo
}

// RemoveProperties removes properties edges to Property.
func (euo *EquipmentUpdateOne) RemoveProperties(p ...*Property) *EquipmentUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return euo.RemovePropertyIDs(ids...)
}

// RemoveFileIDs removes the files edge to File by ids.
func (euo *EquipmentUpdateOne) RemoveFileIDs(ids ...int) *EquipmentUpdateOne {
	euo.mutation.RemoveFileIDs(ids...)
	return euo
}

// RemoveFiles removes files edges to File.
func (euo *EquipmentUpdateOne) RemoveFiles(f ...*File) *EquipmentUpdateOne {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return euo.RemoveFileIDs(ids...)
}

// RemoveHyperlinkIDs removes the hyperlinks edge to Hyperlink by ids.
func (euo *EquipmentUpdateOne) RemoveHyperlinkIDs(ids ...int) *EquipmentUpdateOne {
	euo.mutation.RemoveHyperlinkIDs(ids...)
	return euo
}

// RemoveHyperlinks removes hyperlinks edges to Hyperlink.
func (euo *EquipmentUpdateOne) RemoveHyperlinks(h ...*Hyperlink) *EquipmentUpdateOne {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return euo.RemoveHyperlinkIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (euo *EquipmentUpdateOne) Save(ctx context.Context) (*Equipment, error) {
	if _, ok := euo.mutation.UpdateTime(); !ok {
		v := equipment.UpdateDefaultUpdateTime()
		euo.mutation.SetUpdateTime(v)
	}
	if v, ok := euo.mutation.Name(); ok {
		if err := equipment.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if v, ok := euo.mutation.DeviceID(); ok {
		if err := equipment.DeviceIDValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"device_id\": %v", err)
		}
	}

	if _, ok := euo.mutation.TypeID(); euo.mutation.TypeCleared() && !ok {
		return nil, errors.New("ent: clearing a unique edge \"type\"")
	}

	var (
		err  error
		node *Equipment
	)
	if len(euo.hooks) == 0 {
		node, err = euo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			euo.mutation = mutation
			node, err = euo.sqlSave(ctx)
			return node, err
		})
		for i := len(euo.hooks); i > 0; i-- {
			mut = euo.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, euo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (euo *EquipmentUpdateOne) SaveX(ctx context.Context) *Equipment {
	e, err := euo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return e
}

// Exec executes the query on the entity.
func (euo *EquipmentUpdateOne) Exec(ctx context.Context) error {
	_, err := euo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (euo *EquipmentUpdateOne) ExecX(ctx context.Context) {
	if err := euo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (euo *EquipmentUpdateOne) sqlSave(ctx context.Context) (e *Equipment, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipment.Table,
			Columns: equipment.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipment.FieldID,
			},
		},
	}
	id, ok := euo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing Equipment.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := euo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipment.FieldUpdateTime,
		})
	}
	if value, ok := euo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipment.FieldName,
		})
	}
	if value, ok := euo.mutation.FutureState(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipment.FieldFutureState,
		})
	}
	if euo.mutation.FutureStateCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipment.FieldFutureState,
		})
	}
	if value, ok := euo.mutation.DeviceID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipment.FieldDeviceID,
		})
	}
	if euo.mutation.DeviceIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipment.FieldDeviceID,
		})
	}
	if value, ok := euo.mutation.ExternalID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipment.FieldExternalID,
		})
	}
	if euo.mutation.ExternalIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipment.FieldExternalID,
		})
	}
	if euo.mutation.TypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipment.TypeTable,
			Columns: []string{equipment.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmenttype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.TypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipment.TypeTable,
			Columns: []string{equipment.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmenttype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if euo.mutation.LocationCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipment.LocationTable,
			Columns: []string{equipment.LocationColumn},
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
	if nodes := euo.mutation.LocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipment.LocationTable,
			Columns: []string{equipment.LocationColumn},
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
	if euo.mutation.ParentPositionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   equipment.ParentPositionTable,
			Columns: []string{equipment.ParentPositionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentposition.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.ParentPositionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   equipment.ParentPositionTable,
			Columns: []string{equipment.ParentPositionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentposition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := euo.mutation.RemovedPositionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.PositionsTable,
			Columns: []string{equipment.PositionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentposition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.PositionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.PositionsTable,
			Columns: []string{equipment.PositionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentposition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := euo.mutation.RemovedPortsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.PortsTable,
			Columns: []string{equipment.PortsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentport.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.PortsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.PortsTable,
			Columns: []string{equipment.PortsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentport.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if euo.mutation.WorkOrderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipment.WorkOrderTable,
			Columns: []string{equipment.WorkOrderColumn},
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
	if nodes := euo.mutation.WorkOrderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipment.WorkOrderTable,
			Columns: []string{equipment.WorkOrderColumn},
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
	if nodes := euo.mutation.RemovedPropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.PropertiesTable,
			Columns: []string{equipment.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.PropertiesTable,
			Columns: []string{equipment.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := euo.mutation.RemovedFilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.FilesTable,
			Columns: []string{equipment.FilesColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.FilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.FilesTable,
			Columns: []string{equipment.FilesColumn},
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
	if nodes := euo.mutation.RemovedHyperlinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.HyperlinksTable,
			Columns: []string{equipment.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: hyperlink.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.HyperlinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.HyperlinksTable,
			Columns: []string{equipment.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: hyperlink.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	e = &Equipment{config: euo.config}
	_spec.Assign = e.assignValues
	_spec.ScanValues = e.scanValues()
	if err = sqlgraph.UpdateNode(ctx, euo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{equipment.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return e, nil
}
