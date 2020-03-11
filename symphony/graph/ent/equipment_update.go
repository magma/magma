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

	update_time           *time.Time
	name                  *string
	future_state          *string
	clearfuture_state     bool
	device_id             *string
	cleardevice_id        bool
	external_id           *string
	clearexternal_id      bool
	_type                 map[int]struct{}
	location              map[int]struct{}
	parent_position       map[int]struct{}
	positions             map[int]struct{}
	ports                 map[int]struct{}
	work_order            map[int]struct{}
	properties            map[int]struct{}
	files                 map[int]struct{}
	hyperlinks            map[int]struct{}
	clearedType           bool
	clearedLocation       bool
	clearedParentPosition bool
	removedPositions      map[int]struct{}
	removedPorts          map[int]struct{}
	clearedWorkOrder      bool
	removedProperties     map[int]struct{}
	removedFiles          map[int]struct{}
	removedHyperlinks     map[int]struct{}
	predicates            []predicate.Equipment
}

// Where adds a new predicate for the builder.
func (eu *EquipmentUpdate) Where(ps ...predicate.Equipment) *EquipmentUpdate {
	eu.predicates = append(eu.predicates, ps...)
	return eu
}

// SetName sets the name field.
func (eu *EquipmentUpdate) SetName(s string) *EquipmentUpdate {
	eu.name = &s
	return eu
}

// SetFutureState sets the future_state field.
func (eu *EquipmentUpdate) SetFutureState(s string) *EquipmentUpdate {
	eu.future_state = &s
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
	eu.future_state = nil
	eu.clearfuture_state = true
	return eu
}

// SetDeviceID sets the device_id field.
func (eu *EquipmentUpdate) SetDeviceID(s string) *EquipmentUpdate {
	eu.device_id = &s
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
	eu.device_id = nil
	eu.cleardevice_id = true
	return eu
}

// SetExternalID sets the external_id field.
func (eu *EquipmentUpdate) SetExternalID(s string) *EquipmentUpdate {
	eu.external_id = &s
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
	eu.external_id = nil
	eu.clearexternal_id = true
	return eu
}

// SetTypeID sets the type edge to EquipmentType by id.
func (eu *EquipmentUpdate) SetTypeID(id int) *EquipmentUpdate {
	if eu._type == nil {
		eu._type = make(map[int]struct{})
	}
	eu._type[id] = struct{}{}
	return eu
}

// SetType sets the type edge to EquipmentType.
func (eu *EquipmentUpdate) SetType(e *EquipmentType) *EquipmentUpdate {
	return eu.SetTypeID(e.ID)
}

// SetLocationID sets the location edge to Location by id.
func (eu *EquipmentUpdate) SetLocationID(id int) *EquipmentUpdate {
	if eu.location == nil {
		eu.location = make(map[int]struct{})
	}
	eu.location[id] = struct{}{}
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
	if eu.parent_position == nil {
		eu.parent_position = make(map[int]struct{})
	}
	eu.parent_position[id] = struct{}{}
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
	if eu.positions == nil {
		eu.positions = make(map[int]struct{})
	}
	for i := range ids {
		eu.positions[ids[i]] = struct{}{}
	}
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
	if eu.ports == nil {
		eu.ports = make(map[int]struct{})
	}
	for i := range ids {
		eu.ports[ids[i]] = struct{}{}
	}
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
	if eu.work_order == nil {
		eu.work_order = make(map[int]struct{})
	}
	eu.work_order[id] = struct{}{}
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
	if eu.properties == nil {
		eu.properties = make(map[int]struct{})
	}
	for i := range ids {
		eu.properties[ids[i]] = struct{}{}
	}
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
	if eu.files == nil {
		eu.files = make(map[int]struct{})
	}
	for i := range ids {
		eu.files[ids[i]] = struct{}{}
	}
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
	if eu.hyperlinks == nil {
		eu.hyperlinks = make(map[int]struct{})
	}
	for i := range ids {
		eu.hyperlinks[ids[i]] = struct{}{}
	}
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
	eu.clearedType = true
	return eu
}

// ClearLocation clears the location edge to Location.
func (eu *EquipmentUpdate) ClearLocation() *EquipmentUpdate {
	eu.clearedLocation = true
	return eu
}

// ClearParentPosition clears the parent_position edge to EquipmentPosition.
func (eu *EquipmentUpdate) ClearParentPosition() *EquipmentUpdate {
	eu.clearedParentPosition = true
	return eu
}

// RemovePositionIDs removes the positions edge to EquipmentPosition by ids.
func (eu *EquipmentUpdate) RemovePositionIDs(ids ...int) *EquipmentUpdate {
	if eu.removedPositions == nil {
		eu.removedPositions = make(map[int]struct{})
	}
	for i := range ids {
		eu.removedPositions[ids[i]] = struct{}{}
	}
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
	if eu.removedPorts == nil {
		eu.removedPorts = make(map[int]struct{})
	}
	for i := range ids {
		eu.removedPorts[ids[i]] = struct{}{}
	}
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
	eu.clearedWorkOrder = true
	return eu
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (eu *EquipmentUpdate) RemovePropertyIDs(ids ...int) *EquipmentUpdate {
	if eu.removedProperties == nil {
		eu.removedProperties = make(map[int]struct{})
	}
	for i := range ids {
		eu.removedProperties[ids[i]] = struct{}{}
	}
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
	if eu.removedFiles == nil {
		eu.removedFiles = make(map[int]struct{})
	}
	for i := range ids {
		eu.removedFiles[ids[i]] = struct{}{}
	}
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
	if eu.removedHyperlinks == nil {
		eu.removedHyperlinks = make(map[int]struct{})
	}
	for i := range ids {
		eu.removedHyperlinks[ids[i]] = struct{}{}
	}
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
	if eu.update_time == nil {
		v := equipment.UpdateDefaultUpdateTime()
		eu.update_time = &v
	}
	if eu.name != nil {
		if err := equipment.NameValidator(*eu.name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if eu.device_id != nil {
		if err := equipment.DeviceIDValidator(*eu.device_id); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"device_id\": %v", err)
		}
	}
	if len(eu._type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if eu.clearedType && eu._type == nil {
		return 0, errors.New("ent: clearing a unique edge \"type\"")
	}
	if len(eu.location) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	if len(eu.parent_position) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"parent_position\"")
	}
	if len(eu.work_order) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"work_order\"")
	}
	return eu.sqlSave(ctx)
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
	if value := eu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipment.FieldUpdateTime,
		})
	}
	if value := eu.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipment.FieldName,
		})
	}
	if value := eu.future_state; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipment.FieldFutureState,
		})
	}
	if eu.clearfuture_state {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipment.FieldFutureState,
		})
	}
	if value := eu.device_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipment.FieldDeviceID,
		})
	}
	if eu.cleardevice_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipment.FieldDeviceID,
		})
	}
	if value := eu.external_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipment.FieldExternalID,
		})
	}
	if eu.clearexternal_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipment.FieldExternalID,
		})
	}
	if eu.clearedType {
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
	if nodes := eu._type; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if eu.clearedLocation {
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
	if nodes := eu.location; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if eu.clearedParentPosition {
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
	if nodes := eu.parent_position; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := eu.removedPositions; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.positions; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := eu.removedPorts; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.ports; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if eu.clearedWorkOrder {
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
	if nodes := eu.work_order; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := eu.removedProperties; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.properties; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := eu.removedFiles; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.files; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := eu.removedHyperlinks; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.hyperlinks; len(nodes) > 0 {
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
		for k, _ := range nodes {
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
	id int

	update_time           *time.Time
	name                  *string
	future_state          *string
	clearfuture_state     bool
	device_id             *string
	cleardevice_id        bool
	external_id           *string
	clearexternal_id      bool
	_type                 map[int]struct{}
	location              map[int]struct{}
	parent_position       map[int]struct{}
	positions             map[int]struct{}
	ports                 map[int]struct{}
	work_order            map[int]struct{}
	properties            map[int]struct{}
	files                 map[int]struct{}
	hyperlinks            map[int]struct{}
	clearedType           bool
	clearedLocation       bool
	clearedParentPosition bool
	removedPositions      map[int]struct{}
	removedPorts          map[int]struct{}
	clearedWorkOrder      bool
	removedProperties     map[int]struct{}
	removedFiles          map[int]struct{}
	removedHyperlinks     map[int]struct{}
}

// SetName sets the name field.
func (euo *EquipmentUpdateOne) SetName(s string) *EquipmentUpdateOne {
	euo.name = &s
	return euo
}

// SetFutureState sets the future_state field.
func (euo *EquipmentUpdateOne) SetFutureState(s string) *EquipmentUpdateOne {
	euo.future_state = &s
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
	euo.future_state = nil
	euo.clearfuture_state = true
	return euo
}

// SetDeviceID sets the device_id field.
func (euo *EquipmentUpdateOne) SetDeviceID(s string) *EquipmentUpdateOne {
	euo.device_id = &s
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
	euo.device_id = nil
	euo.cleardevice_id = true
	return euo
}

// SetExternalID sets the external_id field.
func (euo *EquipmentUpdateOne) SetExternalID(s string) *EquipmentUpdateOne {
	euo.external_id = &s
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
	euo.external_id = nil
	euo.clearexternal_id = true
	return euo
}

// SetTypeID sets the type edge to EquipmentType by id.
func (euo *EquipmentUpdateOne) SetTypeID(id int) *EquipmentUpdateOne {
	if euo._type == nil {
		euo._type = make(map[int]struct{})
	}
	euo._type[id] = struct{}{}
	return euo
}

// SetType sets the type edge to EquipmentType.
func (euo *EquipmentUpdateOne) SetType(e *EquipmentType) *EquipmentUpdateOne {
	return euo.SetTypeID(e.ID)
}

// SetLocationID sets the location edge to Location by id.
func (euo *EquipmentUpdateOne) SetLocationID(id int) *EquipmentUpdateOne {
	if euo.location == nil {
		euo.location = make(map[int]struct{})
	}
	euo.location[id] = struct{}{}
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
	if euo.parent_position == nil {
		euo.parent_position = make(map[int]struct{})
	}
	euo.parent_position[id] = struct{}{}
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
	if euo.positions == nil {
		euo.positions = make(map[int]struct{})
	}
	for i := range ids {
		euo.positions[ids[i]] = struct{}{}
	}
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
	if euo.ports == nil {
		euo.ports = make(map[int]struct{})
	}
	for i := range ids {
		euo.ports[ids[i]] = struct{}{}
	}
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
	if euo.work_order == nil {
		euo.work_order = make(map[int]struct{})
	}
	euo.work_order[id] = struct{}{}
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
	if euo.properties == nil {
		euo.properties = make(map[int]struct{})
	}
	for i := range ids {
		euo.properties[ids[i]] = struct{}{}
	}
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
	if euo.files == nil {
		euo.files = make(map[int]struct{})
	}
	for i := range ids {
		euo.files[ids[i]] = struct{}{}
	}
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
	if euo.hyperlinks == nil {
		euo.hyperlinks = make(map[int]struct{})
	}
	for i := range ids {
		euo.hyperlinks[ids[i]] = struct{}{}
	}
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
	euo.clearedType = true
	return euo
}

// ClearLocation clears the location edge to Location.
func (euo *EquipmentUpdateOne) ClearLocation() *EquipmentUpdateOne {
	euo.clearedLocation = true
	return euo
}

// ClearParentPosition clears the parent_position edge to EquipmentPosition.
func (euo *EquipmentUpdateOne) ClearParentPosition() *EquipmentUpdateOne {
	euo.clearedParentPosition = true
	return euo
}

// RemovePositionIDs removes the positions edge to EquipmentPosition by ids.
func (euo *EquipmentUpdateOne) RemovePositionIDs(ids ...int) *EquipmentUpdateOne {
	if euo.removedPositions == nil {
		euo.removedPositions = make(map[int]struct{})
	}
	for i := range ids {
		euo.removedPositions[ids[i]] = struct{}{}
	}
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
	if euo.removedPorts == nil {
		euo.removedPorts = make(map[int]struct{})
	}
	for i := range ids {
		euo.removedPorts[ids[i]] = struct{}{}
	}
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
	euo.clearedWorkOrder = true
	return euo
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (euo *EquipmentUpdateOne) RemovePropertyIDs(ids ...int) *EquipmentUpdateOne {
	if euo.removedProperties == nil {
		euo.removedProperties = make(map[int]struct{})
	}
	for i := range ids {
		euo.removedProperties[ids[i]] = struct{}{}
	}
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
	if euo.removedFiles == nil {
		euo.removedFiles = make(map[int]struct{})
	}
	for i := range ids {
		euo.removedFiles[ids[i]] = struct{}{}
	}
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
	if euo.removedHyperlinks == nil {
		euo.removedHyperlinks = make(map[int]struct{})
	}
	for i := range ids {
		euo.removedHyperlinks[ids[i]] = struct{}{}
	}
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
	if euo.update_time == nil {
		v := equipment.UpdateDefaultUpdateTime()
		euo.update_time = &v
	}
	if euo.name != nil {
		if err := equipment.NameValidator(*euo.name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if euo.device_id != nil {
		if err := equipment.DeviceIDValidator(*euo.device_id); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"device_id\": %v", err)
		}
	}
	if len(euo._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if euo.clearedType && euo._type == nil {
		return nil, errors.New("ent: clearing a unique edge \"type\"")
	}
	if len(euo.location) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	if len(euo.parent_position) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"parent_position\"")
	}
	if len(euo.work_order) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"work_order\"")
	}
	return euo.sqlSave(ctx)
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
				Value:  euo.id,
				Type:   field.TypeInt,
				Column: equipment.FieldID,
			},
		},
	}
	if value := euo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipment.FieldUpdateTime,
		})
	}
	if value := euo.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipment.FieldName,
		})
	}
	if value := euo.future_state; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipment.FieldFutureState,
		})
	}
	if euo.clearfuture_state {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipment.FieldFutureState,
		})
	}
	if value := euo.device_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipment.FieldDeviceID,
		})
	}
	if euo.cleardevice_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipment.FieldDeviceID,
		})
	}
	if value := euo.external_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipment.FieldExternalID,
		})
	}
	if euo.clearexternal_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipment.FieldExternalID,
		})
	}
	if euo.clearedType {
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
	if nodes := euo._type; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if euo.clearedLocation {
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
	if nodes := euo.location; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if euo.clearedParentPosition {
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
	if nodes := euo.parent_position; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := euo.removedPositions; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.positions; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := euo.removedPorts; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.ports; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if euo.clearedWorkOrder {
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
	if nodes := euo.work_order; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := euo.removedProperties; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.properties; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := euo.removedFiles; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.files; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := euo.removedHyperlinks; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.hyperlinks; len(nodes) > 0 {
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
		for k, _ := range nodes {
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
