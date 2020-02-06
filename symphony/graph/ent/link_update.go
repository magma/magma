// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// LinkUpdate is the builder for updating Link entities.
type LinkUpdate struct {
	config

	update_time       *time.Time
	future_state      *string
	clearfuture_state bool
	ports             map[string]struct{}
	work_order        map[string]struct{}
	properties        map[string]struct{}
	service           map[string]struct{}
	removedPorts      map[string]struct{}
	clearedWorkOrder  bool
	removedProperties map[string]struct{}
	removedService    map[string]struct{}
	predicates        []predicate.Link
}

// Where adds a new predicate for the builder.
func (lu *LinkUpdate) Where(ps ...predicate.Link) *LinkUpdate {
	lu.predicates = append(lu.predicates, ps...)
	return lu
}

// SetFutureState sets the future_state field.
func (lu *LinkUpdate) SetFutureState(s string) *LinkUpdate {
	lu.future_state = &s
	return lu
}

// SetNillableFutureState sets the future_state field if the given value is not nil.
func (lu *LinkUpdate) SetNillableFutureState(s *string) *LinkUpdate {
	if s != nil {
		lu.SetFutureState(*s)
	}
	return lu
}

// ClearFutureState clears the value of future_state.
func (lu *LinkUpdate) ClearFutureState() *LinkUpdate {
	lu.future_state = nil
	lu.clearfuture_state = true
	return lu
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (lu *LinkUpdate) AddPortIDs(ids ...string) *LinkUpdate {
	if lu.ports == nil {
		lu.ports = make(map[string]struct{})
	}
	for i := range ids {
		lu.ports[ids[i]] = struct{}{}
	}
	return lu
}

// AddPorts adds the ports edges to EquipmentPort.
func (lu *LinkUpdate) AddPorts(e ...*EquipmentPort) *LinkUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return lu.AddPortIDs(ids...)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (lu *LinkUpdate) SetWorkOrderID(id string) *LinkUpdate {
	if lu.work_order == nil {
		lu.work_order = make(map[string]struct{})
	}
	lu.work_order[id] = struct{}{}
	return lu
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (lu *LinkUpdate) SetNillableWorkOrderID(id *string) *LinkUpdate {
	if id != nil {
		lu = lu.SetWorkOrderID(*id)
	}
	return lu
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (lu *LinkUpdate) SetWorkOrder(w *WorkOrder) *LinkUpdate {
	return lu.SetWorkOrderID(w.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (lu *LinkUpdate) AddPropertyIDs(ids ...string) *LinkUpdate {
	if lu.properties == nil {
		lu.properties = make(map[string]struct{})
	}
	for i := range ids {
		lu.properties[ids[i]] = struct{}{}
	}
	return lu
}

// AddProperties adds the properties edges to Property.
func (lu *LinkUpdate) AddProperties(p ...*Property) *LinkUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return lu.AddPropertyIDs(ids...)
}

// AddServiceIDs adds the service edge to Service by ids.
func (lu *LinkUpdate) AddServiceIDs(ids ...string) *LinkUpdate {
	if lu.service == nil {
		lu.service = make(map[string]struct{})
	}
	for i := range ids {
		lu.service[ids[i]] = struct{}{}
	}
	return lu
}

// AddService adds the service edges to Service.
func (lu *LinkUpdate) AddService(s ...*Service) *LinkUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lu.AddServiceIDs(ids...)
}

// RemovePortIDs removes the ports edge to EquipmentPort by ids.
func (lu *LinkUpdate) RemovePortIDs(ids ...string) *LinkUpdate {
	if lu.removedPorts == nil {
		lu.removedPorts = make(map[string]struct{})
	}
	for i := range ids {
		lu.removedPorts[ids[i]] = struct{}{}
	}
	return lu
}

// RemovePorts removes ports edges to EquipmentPort.
func (lu *LinkUpdate) RemovePorts(e ...*EquipmentPort) *LinkUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return lu.RemovePortIDs(ids...)
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (lu *LinkUpdate) ClearWorkOrder() *LinkUpdate {
	lu.clearedWorkOrder = true
	return lu
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (lu *LinkUpdate) RemovePropertyIDs(ids ...string) *LinkUpdate {
	if lu.removedProperties == nil {
		lu.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		lu.removedProperties[ids[i]] = struct{}{}
	}
	return lu
}

// RemoveProperties removes properties edges to Property.
func (lu *LinkUpdate) RemoveProperties(p ...*Property) *LinkUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return lu.RemovePropertyIDs(ids...)
}

// RemoveServiceIDs removes the service edge to Service by ids.
func (lu *LinkUpdate) RemoveServiceIDs(ids ...string) *LinkUpdate {
	if lu.removedService == nil {
		lu.removedService = make(map[string]struct{})
	}
	for i := range ids {
		lu.removedService[ids[i]] = struct{}{}
	}
	return lu
}

// RemoveService removes service edges to Service.
func (lu *LinkUpdate) RemoveService(s ...*Service) *LinkUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lu.RemoveServiceIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (lu *LinkUpdate) Save(ctx context.Context) (int, error) {
	if lu.update_time == nil {
		v := link.UpdateDefaultUpdateTime()
		lu.update_time = &v
	}
	if len(lu.work_order) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"work_order\"")
	}
	return lu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (lu *LinkUpdate) SaveX(ctx context.Context) int {
	affected, err := lu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (lu *LinkUpdate) Exec(ctx context.Context) error {
	_, err := lu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lu *LinkUpdate) ExecX(ctx context.Context) {
	if err := lu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (lu *LinkUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   link.Table,
			Columns: link.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: link.FieldID,
			},
		},
	}
	if ps := lu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := lu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: link.FieldUpdateTime,
		})
	}
	if value := lu.future_state; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: link.FieldFutureState,
		})
	}
	if lu.clearfuture_state {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: link.FieldFutureState,
		})
	}
	if nodes := lu.removedPorts; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   link.PortsTable,
			Columns: []string{link.PortsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentport.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.ports; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   link.PortsTable,
			Columns: []string{link.PortsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentport.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if lu.clearedWorkOrder {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   link.WorkOrderTable,
			Columns: []string{link.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.work_order; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   link.WorkOrderTable,
			Columns: []string{link.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := lu.removedProperties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   link.PropertiesTable,
			Columns: []string{link.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.properties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   link.PropertiesTable,
			Columns: []string{link.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := lu.removedService; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   link.ServiceTable,
			Columns: link.ServicePrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.service; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   link.ServiceTable,
			Columns: link.ServicePrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, lu.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// LinkUpdateOne is the builder for updating a single Link entity.
type LinkUpdateOne struct {
	config
	id string

	update_time       *time.Time
	future_state      *string
	clearfuture_state bool
	ports             map[string]struct{}
	work_order        map[string]struct{}
	properties        map[string]struct{}
	service           map[string]struct{}
	removedPorts      map[string]struct{}
	clearedWorkOrder  bool
	removedProperties map[string]struct{}
	removedService    map[string]struct{}
}

// SetFutureState sets the future_state field.
func (luo *LinkUpdateOne) SetFutureState(s string) *LinkUpdateOne {
	luo.future_state = &s
	return luo
}

// SetNillableFutureState sets the future_state field if the given value is not nil.
func (luo *LinkUpdateOne) SetNillableFutureState(s *string) *LinkUpdateOne {
	if s != nil {
		luo.SetFutureState(*s)
	}
	return luo
}

// ClearFutureState clears the value of future_state.
func (luo *LinkUpdateOne) ClearFutureState() *LinkUpdateOne {
	luo.future_state = nil
	luo.clearfuture_state = true
	return luo
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (luo *LinkUpdateOne) AddPortIDs(ids ...string) *LinkUpdateOne {
	if luo.ports == nil {
		luo.ports = make(map[string]struct{})
	}
	for i := range ids {
		luo.ports[ids[i]] = struct{}{}
	}
	return luo
}

// AddPorts adds the ports edges to EquipmentPort.
func (luo *LinkUpdateOne) AddPorts(e ...*EquipmentPort) *LinkUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return luo.AddPortIDs(ids...)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (luo *LinkUpdateOne) SetWorkOrderID(id string) *LinkUpdateOne {
	if luo.work_order == nil {
		luo.work_order = make(map[string]struct{})
	}
	luo.work_order[id] = struct{}{}
	return luo
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (luo *LinkUpdateOne) SetNillableWorkOrderID(id *string) *LinkUpdateOne {
	if id != nil {
		luo = luo.SetWorkOrderID(*id)
	}
	return luo
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (luo *LinkUpdateOne) SetWorkOrder(w *WorkOrder) *LinkUpdateOne {
	return luo.SetWorkOrderID(w.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (luo *LinkUpdateOne) AddPropertyIDs(ids ...string) *LinkUpdateOne {
	if luo.properties == nil {
		luo.properties = make(map[string]struct{})
	}
	for i := range ids {
		luo.properties[ids[i]] = struct{}{}
	}
	return luo
}

// AddProperties adds the properties edges to Property.
func (luo *LinkUpdateOne) AddProperties(p ...*Property) *LinkUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return luo.AddPropertyIDs(ids...)
}

// AddServiceIDs adds the service edge to Service by ids.
func (luo *LinkUpdateOne) AddServiceIDs(ids ...string) *LinkUpdateOne {
	if luo.service == nil {
		luo.service = make(map[string]struct{})
	}
	for i := range ids {
		luo.service[ids[i]] = struct{}{}
	}
	return luo
}

// AddService adds the service edges to Service.
func (luo *LinkUpdateOne) AddService(s ...*Service) *LinkUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return luo.AddServiceIDs(ids...)
}

// RemovePortIDs removes the ports edge to EquipmentPort by ids.
func (luo *LinkUpdateOne) RemovePortIDs(ids ...string) *LinkUpdateOne {
	if luo.removedPorts == nil {
		luo.removedPorts = make(map[string]struct{})
	}
	for i := range ids {
		luo.removedPorts[ids[i]] = struct{}{}
	}
	return luo
}

// RemovePorts removes ports edges to EquipmentPort.
func (luo *LinkUpdateOne) RemovePorts(e ...*EquipmentPort) *LinkUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return luo.RemovePortIDs(ids...)
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (luo *LinkUpdateOne) ClearWorkOrder() *LinkUpdateOne {
	luo.clearedWorkOrder = true
	return luo
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (luo *LinkUpdateOne) RemovePropertyIDs(ids ...string) *LinkUpdateOne {
	if luo.removedProperties == nil {
		luo.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		luo.removedProperties[ids[i]] = struct{}{}
	}
	return luo
}

// RemoveProperties removes properties edges to Property.
func (luo *LinkUpdateOne) RemoveProperties(p ...*Property) *LinkUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return luo.RemovePropertyIDs(ids...)
}

// RemoveServiceIDs removes the service edge to Service by ids.
func (luo *LinkUpdateOne) RemoveServiceIDs(ids ...string) *LinkUpdateOne {
	if luo.removedService == nil {
		luo.removedService = make(map[string]struct{})
	}
	for i := range ids {
		luo.removedService[ids[i]] = struct{}{}
	}
	return luo
}

// RemoveService removes service edges to Service.
func (luo *LinkUpdateOne) RemoveService(s ...*Service) *LinkUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return luo.RemoveServiceIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (luo *LinkUpdateOne) Save(ctx context.Context) (*Link, error) {
	if luo.update_time == nil {
		v := link.UpdateDefaultUpdateTime()
		luo.update_time = &v
	}
	if len(luo.work_order) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"work_order\"")
	}
	return luo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (luo *LinkUpdateOne) SaveX(ctx context.Context) *Link {
	l, err := luo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return l
}

// Exec executes the query on the entity.
func (luo *LinkUpdateOne) Exec(ctx context.Context) error {
	_, err := luo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (luo *LinkUpdateOne) ExecX(ctx context.Context) {
	if err := luo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (luo *LinkUpdateOne) sqlSave(ctx context.Context) (l *Link, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   link.Table,
			Columns: link.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  luo.id,
				Type:   field.TypeString,
				Column: link.FieldID,
			},
		},
	}
	if value := luo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: link.FieldUpdateTime,
		})
	}
	if value := luo.future_state; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: link.FieldFutureState,
		})
	}
	if luo.clearfuture_state {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: link.FieldFutureState,
		})
	}
	if nodes := luo.removedPorts; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   link.PortsTable,
			Columns: []string{link.PortsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentport.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.ports; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   link.PortsTable,
			Columns: []string{link.PortsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentport.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if luo.clearedWorkOrder {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   link.WorkOrderTable,
			Columns: []string{link.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.work_order; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   link.WorkOrderTable,
			Columns: []string{link.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := luo.removedProperties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   link.PropertiesTable,
			Columns: []string{link.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.properties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   link.PropertiesTable,
			Columns: []string{link.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := luo.removedService; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   link.ServiceTable,
			Columns: link.ServicePrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.service; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   link.ServiceTable,
			Columns: link.ServicePrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	l = &Link{config: luo.config}
	_spec.Assign = l.assignValues
	_spec.ScanValues = l.scanValues()
	if err = sqlgraph.UpdateNode(ctx, luo.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return l, nil
}
