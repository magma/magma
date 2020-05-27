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
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpointdefinition"
)

// ServiceEndpointCreate is the builder for creating a ServiceEndpoint entity.
type ServiceEndpointCreate struct {
	config
	mutation *ServiceEndpointMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (sec *ServiceEndpointCreate) SetCreateTime(t time.Time) *ServiceEndpointCreate {
	sec.mutation.SetCreateTime(t)
	return sec
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (sec *ServiceEndpointCreate) SetNillableCreateTime(t *time.Time) *ServiceEndpointCreate {
	if t != nil {
		sec.SetCreateTime(*t)
	}
	return sec
}

// SetUpdateTime sets the update_time field.
func (sec *ServiceEndpointCreate) SetUpdateTime(t time.Time) *ServiceEndpointCreate {
	sec.mutation.SetUpdateTime(t)
	return sec
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (sec *ServiceEndpointCreate) SetNillableUpdateTime(t *time.Time) *ServiceEndpointCreate {
	if t != nil {
		sec.SetUpdateTime(*t)
	}
	return sec
}

// SetPortID sets the port edge to EquipmentPort by id.
func (sec *ServiceEndpointCreate) SetPortID(id int) *ServiceEndpointCreate {
	sec.mutation.SetPortID(id)
	return sec
}

// SetNillablePortID sets the port edge to EquipmentPort by id if the given value is not nil.
func (sec *ServiceEndpointCreate) SetNillablePortID(id *int) *ServiceEndpointCreate {
	if id != nil {
		sec = sec.SetPortID(*id)
	}
	return sec
}

// SetPort sets the port edge to EquipmentPort.
func (sec *ServiceEndpointCreate) SetPort(e *EquipmentPort) *ServiceEndpointCreate {
	return sec.SetPortID(e.ID)
}

// SetEquipmentID sets the equipment edge to Equipment by id.
func (sec *ServiceEndpointCreate) SetEquipmentID(id int) *ServiceEndpointCreate {
	sec.mutation.SetEquipmentID(id)
	return sec
}

// SetEquipment sets the equipment edge to Equipment.
func (sec *ServiceEndpointCreate) SetEquipment(e *Equipment) *ServiceEndpointCreate {
	return sec.SetEquipmentID(e.ID)
}

// SetServiceID sets the service edge to Service by id.
func (sec *ServiceEndpointCreate) SetServiceID(id int) *ServiceEndpointCreate {
	sec.mutation.SetServiceID(id)
	return sec
}

// SetService sets the service edge to Service.
func (sec *ServiceEndpointCreate) SetService(s *Service) *ServiceEndpointCreate {
	return sec.SetServiceID(s.ID)
}

// SetDefinitionID sets the definition edge to ServiceEndpointDefinition by id.
func (sec *ServiceEndpointCreate) SetDefinitionID(id int) *ServiceEndpointCreate {
	sec.mutation.SetDefinitionID(id)
	return sec
}

// SetNillableDefinitionID sets the definition edge to ServiceEndpointDefinition by id if the given value is not nil.
func (sec *ServiceEndpointCreate) SetNillableDefinitionID(id *int) *ServiceEndpointCreate {
	if id != nil {
		sec = sec.SetDefinitionID(*id)
	}
	return sec
}

// SetDefinition sets the definition edge to ServiceEndpointDefinition.
func (sec *ServiceEndpointCreate) SetDefinition(s *ServiceEndpointDefinition) *ServiceEndpointCreate {
	return sec.SetDefinitionID(s.ID)
}

// Save creates the ServiceEndpoint in the database.
func (sec *ServiceEndpointCreate) Save(ctx context.Context) (*ServiceEndpoint, error) {
	if _, ok := sec.mutation.CreateTime(); !ok {
		v := serviceendpoint.DefaultCreateTime()
		sec.mutation.SetCreateTime(v)
	}
	if _, ok := sec.mutation.UpdateTime(); !ok {
		v := serviceendpoint.DefaultUpdateTime()
		sec.mutation.SetUpdateTime(v)
	}
	if _, ok := sec.mutation.EquipmentID(); !ok {
		return nil, errors.New("ent: missing required edge \"equipment\"")
	}
	if _, ok := sec.mutation.ServiceID(); !ok {
		return nil, errors.New("ent: missing required edge \"service\"")
	}
	var (
		err  error
		node *ServiceEndpoint
	)
	if len(sec.hooks) == 0 {
		node, err = sec.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ServiceEndpointMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			sec.mutation = mutation
			node, err = sec.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(sec.hooks) - 1; i >= 0; i-- {
			mut = sec.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, sec.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (sec *ServiceEndpointCreate) SaveX(ctx context.Context) *ServiceEndpoint {
	v, err := sec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (sec *ServiceEndpointCreate) sqlSave(ctx context.Context) (*ServiceEndpoint, error) {
	var (
		se    = &ServiceEndpoint{config: sec.config}
		_spec = &sqlgraph.CreateSpec{
			Table: serviceendpoint.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: serviceendpoint.FieldID,
			},
		}
	)
	if value, ok := sec.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: serviceendpoint.FieldCreateTime,
		})
		se.CreateTime = value
	}
	if value, ok := sec.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: serviceendpoint.FieldUpdateTime,
		})
		se.UpdateTime = value
	}
	if nodes := sec.mutation.PortIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   serviceendpoint.PortTable,
			Columns: []string{serviceendpoint.PortColumn},
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
	if nodes := sec.mutation.EquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   serviceendpoint.EquipmentTable,
			Columns: []string{serviceendpoint.EquipmentColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := sec.mutation.ServiceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   serviceendpoint.ServiceTable,
			Columns: []string{serviceendpoint.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: service.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := sec.mutation.DefinitionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   serviceendpoint.DefinitionTable,
			Columns: []string{serviceendpoint.DefinitionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: serviceendpointdefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, sec.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	se.ID = int(id)
	return se, nil
}
