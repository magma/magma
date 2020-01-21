// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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
	create_time     *time.Time
	update_time     *time.Time
	name            *string
	future_state    *string
	device_id       *string
	external_id     *string
	_type           map[string]struct{}
	location        map[string]struct{}
	parent_position map[string]struct{}
	positions       map[string]struct{}
	ports           map[string]struct{}
	work_order      map[string]struct{}
	properties      map[string]struct{}
	files           map[string]struct{}
	hyperlinks      map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (ec *EquipmentCreate) SetCreateTime(t time.Time) *EquipmentCreate {
	ec.create_time = &t
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
	ec.update_time = &t
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
	ec.name = &s
	return ec
}

// SetFutureState sets the future_state field.
func (ec *EquipmentCreate) SetFutureState(s string) *EquipmentCreate {
	ec.future_state = &s
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
	ec.device_id = &s
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
	ec.external_id = &s
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
func (ec *EquipmentCreate) SetTypeID(id string) *EquipmentCreate {
	if ec._type == nil {
		ec._type = make(map[string]struct{})
	}
	ec._type[id] = struct{}{}
	return ec
}

// SetType sets the type edge to EquipmentType.
func (ec *EquipmentCreate) SetType(e *EquipmentType) *EquipmentCreate {
	return ec.SetTypeID(e.ID)
}

// SetLocationID sets the location edge to Location by id.
func (ec *EquipmentCreate) SetLocationID(id string) *EquipmentCreate {
	if ec.location == nil {
		ec.location = make(map[string]struct{})
	}
	ec.location[id] = struct{}{}
	return ec
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (ec *EquipmentCreate) SetNillableLocationID(id *string) *EquipmentCreate {
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
func (ec *EquipmentCreate) SetParentPositionID(id string) *EquipmentCreate {
	if ec.parent_position == nil {
		ec.parent_position = make(map[string]struct{})
	}
	ec.parent_position[id] = struct{}{}
	return ec
}

// SetNillableParentPositionID sets the parent_position edge to EquipmentPosition by id if the given value is not nil.
func (ec *EquipmentCreate) SetNillableParentPositionID(id *string) *EquipmentCreate {
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
func (ec *EquipmentCreate) AddPositionIDs(ids ...string) *EquipmentCreate {
	if ec.positions == nil {
		ec.positions = make(map[string]struct{})
	}
	for i := range ids {
		ec.positions[ids[i]] = struct{}{}
	}
	return ec
}

// AddPositions adds the positions edges to EquipmentPosition.
func (ec *EquipmentCreate) AddPositions(e ...*EquipmentPosition) *EquipmentCreate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ec.AddPositionIDs(ids...)
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (ec *EquipmentCreate) AddPortIDs(ids ...string) *EquipmentCreate {
	if ec.ports == nil {
		ec.ports = make(map[string]struct{})
	}
	for i := range ids {
		ec.ports[ids[i]] = struct{}{}
	}
	return ec
}

// AddPorts adds the ports edges to EquipmentPort.
func (ec *EquipmentCreate) AddPorts(e ...*EquipmentPort) *EquipmentCreate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ec.AddPortIDs(ids...)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (ec *EquipmentCreate) SetWorkOrderID(id string) *EquipmentCreate {
	if ec.work_order == nil {
		ec.work_order = make(map[string]struct{})
	}
	ec.work_order[id] = struct{}{}
	return ec
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (ec *EquipmentCreate) SetNillableWorkOrderID(id *string) *EquipmentCreate {
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
func (ec *EquipmentCreate) AddPropertyIDs(ids ...string) *EquipmentCreate {
	if ec.properties == nil {
		ec.properties = make(map[string]struct{})
	}
	for i := range ids {
		ec.properties[ids[i]] = struct{}{}
	}
	return ec
}

// AddProperties adds the properties edges to Property.
func (ec *EquipmentCreate) AddProperties(p ...*Property) *EquipmentCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ec.AddPropertyIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (ec *EquipmentCreate) AddFileIDs(ids ...string) *EquipmentCreate {
	if ec.files == nil {
		ec.files = make(map[string]struct{})
	}
	for i := range ids {
		ec.files[ids[i]] = struct{}{}
	}
	return ec
}

// AddFiles adds the files edges to File.
func (ec *EquipmentCreate) AddFiles(f ...*File) *EquipmentCreate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return ec.AddFileIDs(ids...)
}

// AddHyperlinkIDs adds the hyperlinks edge to Hyperlink by ids.
func (ec *EquipmentCreate) AddHyperlinkIDs(ids ...string) *EquipmentCreate {
	if ec.hyperlinks == nil {
		ec.hyperlinks = make(map[string]struct{})
	}
	for i := range ids {
		ec.hyperlinks[ids[i]] = struct{}{}
	}
	return ec
}

// AddHyperlinks adds the hyperlinks edges to Hyperlink.
func (ec *EquipmentCreate) AddHyperlinks(h ...*Hyperlink) *EquipmentCreate {
	ids := make([]string, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return ec.AddHyperlinkIDs(ids...)
}

// Save creates the Equipment in the database.
func (ec *EquipmentCreate) Save(ctx context.Context) (*Equipment, error) {
	if ec.create_time == nil {
		v := equipment.DefaultCreateTime()
		ec.create_time = &v
	}
	if ec.update_time == nil {
		v := equipment.DefaultUpdateTime()
		ec.update_time = &v
	}
	if ec.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if err := equipment.NameValidator(*ec.name); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
	}
	if len(ec._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if ec._type == nil {
		return nil, errors.New("ent: missing required edge \"type\"")
	}
	if len(ec.location) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	if len(ec.parent_position) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"parent_position\"")
	}
	if len(ec.work_order) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"work_order\"")
	}
	return ec.sqlSave(ctx)
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
		e    = &Equipment{config: ec.config}
		spec = &sqlgraph.CreateSpec{
			Table: equipment.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: equipment.FieldID,
			},
		}
	)
	if value := ec.create_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipment.FieldCreateTime,
		})
		e.CreateTime = *value
	}
	if value := ec.update_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipment.FieldUpdateTime,
		})
		e.UpdateTime = *value
	}
	if value := ec.name; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipment.FieldName,
		})
		e.Name = *value
	}
	if value := ec.future_state; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipment.FieldFutureState,
		})
		e.FutureState = *value
	}
	if value := ec.device_id; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipment.FieldDeviceID,
		})
		e.DeviceID = *value
	}
	if value := ec.external_id; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipment.FieldExternalID,
		})
		e.ExternalID = *value
	}
	if nodes := ec._type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipment.TypeTable,
			Columns: []string{equipment.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmenttype.FieldID,
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
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := ec.location; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipment.LocationTable,
			Columns: []string{equipment.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
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
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := ec.parent_position; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   equipment.ParentPositionTable,
			Columns: []string{equipment.ParentPositionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentposition.FieldID,
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
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := ec.positions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.PositionsTable,
			Columns: []string{equipment.PositionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentposition.FieldID,
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
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := ec.ports; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.PortsTable,
			Columns: []string{equipment.PortsColumn},
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
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := ec.work_order; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipment.WorkOrderTable,
			Columns: []string{equipment.WorkOrderColumn},
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
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := ec.properties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.PropertiesTable,
			Columns: []string{equipment.PropertiesColumn},
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
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := ec.files; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.FilesTable,
			Columns: []string{equipment.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: file.FieldID,
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
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := ec.hyperlinks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipment.HyperlinksTable,
			Columns: []string{equipment.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: hyperlink.FieldID,
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
		spec.Edges = append(spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, ec.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := spec.ID.Value.(int64)
	e.ID = strconv.FormatInt(id, 10)
	return e, nil
}
