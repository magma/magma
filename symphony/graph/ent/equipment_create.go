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
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// EquipmentCreate is the builder for creating a Equipment entity.
type EquipmentCreate struct {
	config
	mutation *EquipmentMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (ec *EquipmentCreate) SetCreateTime(t time.Time) *EquipmentCreate {
	ec.mutation.SetCreateTime(t)
	return ec
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (ec *EquipmentCreate) SetNillableCreateTime(t *time.Time) *EquipmentCreate {
	if t != nil {
		ec.SetCreateTime(*t)
	}
	return ec
}

// SetUpdateTime sets the update_time field.
func (ec *EquipmentCreate) SetUpdateTime(t time.Time) *EquipmentCreate {
	ec.mutation.SetUpdateTime(t)
	return ec
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (ec *EquipmentCreate) SetNillableUpdateTime(t *time.Time) *EquipmentCreate {
	if t != nil {
		ec.SetUpdateTime(*t)
	}
	return ec
}

// SetName sets the name field.
func (ec *EquipmentCreate) SetName(s string) *EquipmentCreate {
	ec.mutation.SetName(s)
	return ec
}

// SetFutureState sets the future_state field.
func (ec *EquipmentCreate) SetFutureState(s string) *EquipmentCreate {
	ec.mutation.SetFutureState(s)
	return ec
}

// SetNillableFutureState sets the future_state field if the given value is not nil.
func (ec *EquipmentCreate) SetNillableFutureState(s *string) *EquipmentCreate {
	if s != nil {
		ec.SetFutureState(*s)
	}
	return ec
}

// SetDeviceID sets the device_id field.
func (ec *EquipmentCreate) SetDeviceID(s string) *EquipmentCreate {
	ec.mutation.SetDeviceID(s)
	return ec
}

// SetNillableDeviceID sets the device_id field if the given value is not nil.
func (ec *EquipmentCreate) SetNillableDeviceID(s *string) *EquipmentCreate {
	if s != nil {
		ec.SetDeviceID(*s)
	}
	return ec
}

// SetExternalID sets the external_id field.
func (ec *EquipmentCreate) SetExternalID(s string) *EquipmentCreate {
	ec.mutation.SetExternalID(s)
	return ec
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (ec *EquipmentCreate) SetNillableExternalID(s *string) *EquipmentCreate {
	if s != nil {
		ec.SetExternalID(*s)
	}
	return ec
}

// SetTypeID sets the type edge to EquipmentType by id.
func (ec *EquipmentCreate) SetTypeID(id int) *EquipmentCreate {
	ec.mutation.SetTypeID(id)
	return ec
}

// SetType sets the type edge to EquipmentType.
func (ec *EquipmentCreate) SetType(e *EquipmentType) *EquipmentCreate {
	return ec.SetTypeID(e.ID)
}

// SetLocationID sets the location edge to Location by id.
func (ec *EquipmentCreate) SetLocationID(id int) *EquipmentCreate {
	ec.mutation.SetLocationID(id)
	return ec
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (ec *EquipmentCreate) SetNillableLocationID(id *int) *EquipmentCreate {
	if id != nil {
		ec = ec.SetLocationID(*id)
	}
	return ec
}

// SetLocation sets the location edge to Location.
func (ec *EquipmentCreate) SetLocation(l *Location) *EquipmentCreate {
	return ec.SetLocationID(l.ID)
}

// SetParentPositionID sets the parent_position edge to EquipmentPosition by id.
func (ec *EquipmentCreate) SetParentPositionID(id int) *EquipmentCreate {
	ec.mutation.SetParentPositionID(id)
	return ec
}

// SetNillableParentPositionID sets the parent_position edge to EquipmentPosition by id if the given value is not nil.
func (ec *EquipmentCreate) SetNillableParentPositionID(id *int) *EquipmentCreate {
	if id != nil {
		ec = ec.SetParentPositionID(*id)
	}
	return ec
}

// SetParentPosition sets the parent_position edge to EquipmentPosition.
func (ec *EquipmentCreate) SetParentPosition(e *EquipmentPosition) *EquipmentCreate {
	return ec.SetParentPositionID(e.ID)
}

// AddPositionIDs adds the positions edge to EquipmentPosition by ids.
func (ec *EquipmentCreate) AddPositionIDs(ids ...int) *EquipmentCreate {
	ec.mutation.AddPositionIDs(ids...)
	return ec
}

// AddPositions adds the positions edges to EquipmentPosition.
func (ec *EquipmentCreate) AddPositions(e ...*EquipmentPosition) *EquipmentCreate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ec.AddPositionIDs(ids...)
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (ec *EquipmentCreate) AddPortIDs(ids ...int) *EquipmentCreate {
	ec.mutation.AddPortIDs(ids...)
	return ec
}

// AddPorts adds the ports edges to EquipmentPort.
func (ec *EquipmentCreate) AddPorts(e ...*EquipmentPort) *EquipmentCreate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ec.AddPortIDs(ids...)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (ec *EquipmentCreate) SetWorkOrderID(id int) *EquipmentCreate {
	ec.mutation.SetWorkOrderID(id)
	return ec
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (ec *EquipmentCreate) SetNillableWorkOrderID(id *int) *EquipmentCreate {
	if id != nil {
		ec = ec.SetWorkOrderID(*id)
	}
	return ec
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (ec *EquipmentCreate) SetWorkOrder(w *WorkOrder) *EquipmentCreate {
	return ec.SetWorkOrderID(w.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (ec *EquipmentCreate) AddPropertyIDs(ids ...int) *EquipmentCreate {
	ec.mutation.AddPropertyIDs(ids...)
	return ec
}

// AddProperties adds the properties edges to Property.
func (ec *EquipmentCreate) AddProperties(p ...*Property) *EquipmentCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ec.AddPropertyIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (ec *EquipmentCreate) AddFileIDs(ids ...int) *EquipmentCreate {
	ec.mutation.AddFileIDs(ids...)
	return ec
}

// AddFiles adds the files edges to File.
func (ec *EquipmentCreate) AddFiles(f ...*File) *EquipmentCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return ec.AddFileIDs(ids...)
}

// AddHyperlinkIDs adds the hyperlinks edge to Hyperlink by ids.
func (ec *EquipmentCreate) AddHyperlinkIDs(ids ...int) *EquipmentCreate {
	ec.mutation.AddHyperlinkIDs(ids...)
	return ec
}

// AddHyperlinks adds the hyperlinks edges to Hyperlink.
func (ec *EquipmentCreate) AddHyperlinks(h ...*Hyperlink) *EquipmentCreate {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return ec.AddHyperlinkIDs(ids...)
}

// Save creates the Equipment in the database.
func (ec *EquipmentCreate) Save(ctx context.Context) (*Equipment, error) {
	if _, ok := ec.mutation.CreateTime(); !ok {
		v := equipment.DefaultCreateTime()
		ec.mutation.SetCreateTime(v)
	}
	if _, ok := ec.mutation.UpdateTime(); !ok {
		v := equipment.DefaultUpdateTime()
		ec.mutation.SetUpdateTime(v)
	}
	if _, ok := ec.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if v, ok := ec.mutation.Name(); ok {
		if err := equipment.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if v, ok := ec.mutation.DeviceID(); ok {
		if err := equipment.DeviceIDValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"device_id\": %v", err)
		}
	}
	if _, ok := ec.mutation.TypeID(); !ok {
		return nil, errors.New("ent: missing required edge \"type\"")
	}
	var (
		err  error
		node *Equipment
	)
	if len(ec.hooks) == 0 {
		node, err = ec.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ec.mutation = mutation
			node, err = ec.sqlSave(ctx)
			return node, err
		})
		for i := len(ec.hooks); i > 0; i-- {
			mut = ec.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, ec.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (ec *EquipmentCreate) SaveX(ctx context.Context) *Equipment {
	v, err := ec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ec *EquipmentCreate) sqlSave(ctx context.Context) (*Equipment, error) {
	var (
		e     = &Equipment{config: ec.config}
		_spec = &sqlgraph.CreateSpec{
			Table: equipment.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipment.FieldID,
			},
		}
	)
	if value, ok := ec.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipment.FieldCreateTime,
		})
		e.CreateTime = value
	}
	if value, ok := ec.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipment.FieldUpdateTime,
		})
		e.UpdateTime = value
	}
	if value, ok := ec.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipment.FieldName,
		})
		e.Name = value
	}
	if value, ok := ec.mutation.FutureState(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipment.FieldFutureState,
		})
		e.FutureState = value
	}
	if value, ok := ec.mutation.DeviceID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipment.FieldDeviceID,
		})
		e.DeviceID = value
	}
	if value, ok := ec.mutation.ExternalID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipment.FieldExternalID,
		})
		e.ExternalID = value
	}
	if nodes := ec.mutation.TypeIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.LocationIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.ParentPositionIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.PositionsIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.PortsIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.WorkOrderIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.PropertiesIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.FilesIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.HyperlinksIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, ec.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	e.ID = int(id)
	return e, nil
}
