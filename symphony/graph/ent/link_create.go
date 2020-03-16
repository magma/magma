// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// LinkCreate is the builder for creating a Link entity.
type LinkCreate struct {
	config
	create_time  *time.Time
	update_time  *time.Time
	future_state *string
	ports        map[int]struct{}
	work_order   map[int]struct{}
	properties   map[int]struct{}
	service      map[int]struct{}
}

// SetCreateTime sets the create_time field.
func (lc *LinkCreate) SetCreateTime(t time.Time) *LinkCreate {
	lc.create_time = &t
	return lc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (lc *LinkCreate) SetNillableCreateTime(t *time.Time) *LinkCreate {
	if t != nil {
		lc.SetCreateTime(*t)
	}
	return lc
}

// SetUpdateTime sets the update_time field.
func (lc *LinkCreate) SetUpdateTime(t time.Time) *LinkCreate {
	lc.update_time = &t
	return lc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (lc *LinkCreate) SetNillableUpdateTime(t *time.Time) *LinkCreate {
	if t != nil {
		lc.SetUpdateTime(*t)
	}
	return lc
}

// SetFutureState sets the future_state field.
func (lc *LinkCreate) SetFutureState(s string) *LinkCreate {
	lc.future_state = &s
	return lc
}

// SetNillableFutureState sets the future_state field if the given value is not nil.
func (lc *LinkCreate) SetNillableFutureState(s *string) *LinkCreate {
	if s != nil {
		lc.SetFutureState(*s)
	}
	return lc
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (lc *LinkCreate) AddPortIDs(ids ...int) *LinkCreate {
	if lc.ports == nil {
		lc.ports = make(map[int]struct{})
	}
	for i := range ids {
		lc.ports[ids[i]] = struct{}{}
	}
	return lc
}

// AddPorts adds the ports edges to EquipmentPort.
func (lc *LinkCreate) AddPorts(e ...*EquipmentPort) *LinkCreate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return lc.AddPortIDs(ids...)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (lc *LinkCreate) SetWorkOrderID(id int) *LinkCreate {
	if lc.work_order == nil {
		lc.work_order = make(map[int]struct{})
	}
	lc.work_order[id] = struct{}{}
	return lc
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (lc *LinkCreate) SetNillableWorkOrderID(id *int) *LinkCreate {
	if id != nil {
		lc = lc.SetWorkOrderID(*id)
	}
	return lc
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (lc *LinkCreate) SetWorkOrder(w *WorkOrder) *LinkCreate {
	return lc.SetWorkOrderID(w.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (lc *LinkCreate) AddPropertyIDs(ids ...int) *LinkCreate {
	if lc.properties == nil {
		lc.properties = make(map[int]struct{})
	}
	for i := range ids {
		lc.properties[ids[i]] = struct{}{}
	}
	return lc
}

// AddProperties adds the properties edges to Property.
func (lc *LinkCreate) AddProperties(p ...*Property) *LinkCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return lc.AddPropertyIDs(ids...)
}

// AddServiceIDs adds the service edge to Service by ids.
func (lc *LinkCreate) AddServiceIDs(ids ...int) *LinkCreate {
	if lc.service == nil {
		lc.service = make(map[int]struct{})
	}
	for i := range ids {
		lc.service[ids[i]] = struct{}{}
	}
	return lc
}

// AddService adds the service edges to Service.
func (lc *LinkCreate) AddService(s ...*Service) *LinkCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lc.AddServiceIDs(ids...)
}

// Save creates the Link in the database.
func (lc *LinkCreate) Save(ctx context.Context) (*Link, error) {
	if lc.create_time == nil {
		v := link.DefaultCreateTime()
		lc.create_time = &v
	}
	if lc.update_time == nil {
		v := link.DefaultUpdateTime()
		lc.update_time = &v
	}
	if len(lc.work_order) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"work_order\"")
	}
	return lc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (lc *LinkCreate) SaveX(ctx context.Context) *Link {
	v, err := lc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (lc *LinkCreate) sqlSave(ctx context.Context) (*Link, error) {
	var (
		l     = &Link{config: lc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: link.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: link.FieldID,
			},
		}
	)
	if value := lc.create_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: link.FieldCreateTime,
		})
		l.CreateTime = *value
	}
	if value := lc.update_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: link.FieldUpdateTime,
		})
		l.UpdateTime = *value
	}
	if value := lc.future_state; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: link.FieldFutureState,
		})
		l.FutureState = *value
	}
	if nodes := lc.ports; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   link.PortsTable,
			Columns: []string{link.PortsColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lc.work_order; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   link.WorkOrderTable,
			Columns: []string{link.WorkOrderColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lc.properties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   link.PropertiesTable,
			Columns: []string{link.PropertiesColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lc.service; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   link.ServiceTable,
			Columns: link.ServicePrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: service.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, lc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	l.ID = int(id)
	return l, nil
}
