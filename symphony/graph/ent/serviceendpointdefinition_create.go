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
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpointdefinition"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
)

// ServiceEndpointDefinitionCreate is the builder for creating a ServiceEndpointDefinition entity.
type ServiceEndpointDefinitionCreate struct {
	config
	mutation *ServiceEndpointDefinitionMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (sedc *ServiceEndpointDefinitionCreate) SetCreateTime(t time.Time) *ServiceEndpointDefinitionCreate {
	sedc.mutation.SetCreateTime(t)
	return sedc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (sedc *ServiceEndpointDefinitionCreate) SetNillableCreateTime(t *time.Time) *ServiceEndpointDefinitionCreate {
	if t != nil {
		sedc.SetCreateTime(*t)
	}
	return sedc
}

// SetUpdateTime sets the update_time field.
func (sedc *ServiceEndpointDefinitionCreate) SetUpdateTime(t time.Time) *ServiceEndpointDefinitionCreate {
	sedc.mutation.SetUpdateTime(t)
	return sedc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (sedc *ServiceEndpointDefinitionCreate) SetNillableUpdateTime(t *time.Time) *ServiceEndpointDefinitionCreate {
	if t != nil {
		sedc.SetUpdateTime(*t)
	}
	return sedc
}

// SetRole sets the role field.
func (sedc *ServiceEndpointDefinitionCreate) SetRole(s string) *ServiceEndpointDefinitionCreate {
	sedc.mutation.SetRole(s)
	return sedc
}

// SetNillableRole sets the role field if the given value is not nil.
func (sedc *ServiceEndpointDefinitionCreate) SetNillableRole(s *string) *ServiceEndpointDefinitionCreate {
	if s != nil {
		sedc.SetRole(*s)
	}
	return sedc
}

// SetName sets the name field.
func (sedc *ServiceEndpointDefinitionCreate) SetName(s string) *ServiceEndpointDefinitionCreate {
	sedc.mutation.SetName(s)
	return sedc
}

// SetIndex sets the index field.
func (sedc *ServiceEndpointDefinitionCreate) SetIndex(i int) *ServiceEndpointDefinitionCreate {
	sedc.mutation.SetIndex(i)
	return sedc
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (sedc *ServiceEndpointDefinitionCreate) AddEndpointIDs(ids ...int) *ServiceEndpointDefinitionCreate {
	sedc.mutation.AddEndpointIDs(ids...)
	return sedc
}

// AddEndpoints adds the endpoints edges to ServiceEndpoint.
func (sedc *ServiceEndpointDefinitionCreate) AddEndpoints(s ...*ServiceEndpoint) *ServiceEndpointDefinitionCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sedc.AddEndpointIDs(ids...)
}

// SetServiceTypeID sets the service_type edge to ServiceType by id.
func (sedc *ServiceEndpointDefinitionCreate) SetServiceTypeID(id int) *ServiceEndpointDefinitionCreate {
	sedc.mutation.SetServiceTypeID(id)
	return sedc
}

// SetNillableServiceTypeID sets the service_type edge to ServiceType by id if the given value is not nil.
func (sedc *ServiceEndpointDefinitionCreate) SetNillableServiceTypeID(id *int) *ServiceEndpointDefinitionCreate {
	if id != nil {
		sedc = sedc.SetServiceTypeID(*id)
	}
	return sedc
}

// SetServiceType sets the service_type edge to ServiceType.
func (sedc *ServiceEndpointDefinitionCreate) SetServiceType(s *ServiceType) *ServiceEndpointDefinitionCreate {
	return sedc.SetServiceTypeID(s.ID)
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (sedc *ServiceEndpointDefinitionCreate) SetEquipmentTypeID(id int) *ServiceEndpointDefinitionCreate {
	sedc.mutation.SetEquipmentTypeID(id)
	return sedc
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (sedc *ServiceEndpointDefinitionCreate) SetNillableEquipmentTypeID(id *int) *ServiceEndpointDefinitionCreate {
	if id != nil {
		sedc = sedc.SetEquipmentTypeID(*id)
	}
	return sedc
}

// SetEquipmentType sets the equipment_type edge to EquipmentType.
func (sedc *ServiceEndpointDefinitionCreate) SetEquipmentType(e *EquipmentType) *ServiceEndpointDefinitionCreate {
	return sedc.SetEquipmentTypeID(e.ID)
}

// Save creates the ServiceEndpointDefinition in the database.
func (sedc *ServiceEndpointDefinitionCreate) Save(ctx context.Context) (*ServiceEndpointDefinition, error) {
	if _, ok := sedc.mutation.CreateTime(); !ok {
		v := serviceendpointdefinition.DefaultCreateTime()
		sedc.mutation.SetCreateTime(v)
	}
	if _, ok := sedc.mutation.UpdateTime(); !ok {
		v := serviceendpointdefinition.DefaultUpdateTime()
		sedc.mutation.SetUpdateTime(v)
	}
	if _, ok := sedc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if v, ok := sedc.mutation.Name(); ok {
		if err := serviceendpointdefinition.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if _, ok := sedc.mutation.Index(); !ok {
		return nil, errors.New("ent: missing required field \"index\"")
	}
	var (
		err  error
		node *ServiceEndpointDefinition
	)
	if len(sedc.hooks) == 0 {
		node, err = sedc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ServiceEndpointDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			sedc.mutation = mutation
			node, err = sedc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(sedc.hooks) - 1; i >= 0; i-- {
			mut = sedc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, sedc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (sedc *ServiceEndpointDefinitionCreate) SaveX(ctx context.Context) *ServiceEndpointDefinition {
	v, err := sedc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (sedc *ServiceEndpointDefinitionCreate) sqlSave(ctx context.Context) (*ServiceEndpointDefinition, error) {
	var (
		sed   = &ServiceEndpointDefinition{config: sedc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: serviceendpointdefinition.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: serviceendpointdefinition.FieldID,
			},
		}
	)
	if value, ok := sedc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: serviceendpointdefinition.FieldCreateTime,
		})
		sed.CreateTime = value
	}
	if value, ok := sedc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: serviceendpointdefinition.FieldUpdateTime,
		})
		sed.UpdateTime = value
	}
	if value, ok := sedc.mutation.Role(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: serviceendpointdefinition.FieldRole,
		})
		sed.Role = value
	}
	if value, ok := sedc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: serviceendpointdefinition.FieldName,
		})
		sed.Name = value
	}
	if value, ok := sedc.mutation.Index(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: serviceendpointdefinition.FieldIndex,
		})
		sed.Index = value
	}
	if nodes := sedc.mutation.EndpointsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   serviceendpointdefinition.EndpointsTable,
			Columns: []string{serviceendpointdefinition.EndpointsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: serviceendpoint.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := sedc.mutation.ServiceTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   serviceendpointdefinition.ServiceTypeTable,
			Columns: []string{serviceendpointdefinition.ServiceTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: servicetype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := sedc.mutation.EquipmentTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   serviceendpointdefinition.EquipmentTypeTable,
			Columns: []string{serviceendpointdefinition.EquipmentTypeColumn},
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
	if err := sqlgraph.CreateNode(ctx, sedc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	sed.ID = int(id)
	return sed, nil
}
