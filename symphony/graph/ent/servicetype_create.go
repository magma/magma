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
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpointdefinition"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
)

// ServiceTypeCreate is the builder for creating a ServiceType entity.
type ServiceTypeCreate struct {
	config
	mutation *ServiceTypeMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (stc *ServiceTypeCreate) SetCreateTime(t time.Time) *ServiceTypeCreate {
	stc.mutation.SetCreateTime(t)
	return stc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (stc *ServiceTypeCreate) SetNillableCreateTime(t *time.Time) *ServiceTypeCreate {
	if t != nil {
		stc.SetCreateTime(*t)
	}
	return stc
}

// SetUpdateTime sets the update_time field.
func (stc *ServiceTypeCreate) SetUpdateTime(t time.Time) *ServiceTypeCreate {
	stc.mutation.SetUpdateTime(t)
	return stc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (stc *ServiceTypeCreate) SetNillableUpdateTime(t *time.Time) *ServiceTypeCreate {
	if t != nil {
		stc.SetUpdateTime(*t)
	}
	return stc
}

// SetName sets the name field.
func (stc *ServiceTypeCreate) SetName(s string) *ServiceTypeCreate {
	stc.mutation.SetName(s)
	return stc
}

// SetHasCustomer sets the has_customer field.
func (stc *ServiceTypeCreate) SetHasCustomer(b bool) *ServiceTypeCreate {
	stc.mutation.SetHasCustomer(b)
	return stc
}

// SetNillableHasCustomer sets the has_customer field if the given value is not nil.
func (stc *ServiceTypeCreate) SetNillableHasCustomer(b *bool) *ServiceTypeCreate {
	if b != nil {
		stc.SetHasCustomer(*b)
	}
	return stc
}

// SetIsDeleted sets the is_deleted field.
func (stc *ServiceTypeCreate) SetIsDeleted(b bool) *ServiceTypeCreate {
	stc.mutation.SetIsDeleted(b)
	return stc
}

// SetNillableIsDeleted sets the is_deleted field if the given value is not nil.
func (stc *ServiceTypeCreate) SetNillableIsDeleted(b *bool) *ServiceTypeCreate {
	if b != nil {
		stc.SetIsDeleted(*b)
	}
	return stc
}

// SetDiscoveryMethod sets the discovery_method field.
func (stc *ServiceTypeCreate) SetDiscoveryMethod(sm servicetype.DiscoveryMethod) *ServiceTypeCreate {
	stc.mutation.SetDiscoveryMethod(sm)
	return stc
}

// SetNillableDiscoveryMethod sets the discovery_method field if the given value is not nil.
func (stc *ServiceTypeCreate) SetNillableDiscoveryMethod(sm *servicetype.DiscoveryMethod) *ServiceTypeCreate {
	if sm != nil {
		stc.SetDiscoveryMethod(*sm)
	}
	return stc
}

// AddServiceIDs adds the services edge to Service by ids.
func (stc *ServiceTypeCreate) AddServiceIDs(ids ...int) *ServiceTypeCreate {
	stc.mutation.AddServiceIDs(ids...)
	return stc
}

// AddServices adds the services edges to Service.
func (stc *ServiceTypeCreate) AddServices(s ...*Service) *ServiceTypeCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stc.AddServiceIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (stc *ServiceTypeCreate) AddPropertyTypeIDs(ids ...int) *ServiceTypeCreate {
	stc.mutation.AddPropertyTypeIDs(ids...)
	return stc
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (stc *ServiceTypeCreate) AddPropertyTypes(p ...*PropertyType) *ServiceTypeCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return stc.AddPropertyTypeIDs(ids...)
}

// AddEndpointDefinitionIDs adds the endpoint_definitions edge to ServiceEndpointDefinition by ids.
func (stc *ServiceTypeCreate) AddEndpointDefinitionIDs(ids ...int) *ServiceTypeCreate {
	stc.mutation.AddEndpointDefinitionIDs(ids...)
	return stc
}

// AddEndpointDefinitions adds the endpoint_definitions edges to ServiceEndpointDefinition.
func (stc *ServiceTypeCreate) AddEndpointDefinitions(s ...*ServiceEndpointDefinition) *ServiceTypeCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stc.AddEndpointDefinitionIDs(ids...)
}

// Save creates the ServiceType in the database.
func (stc *ServiceTypeCreate) Save(ctx context.Context) (*ServiceType, error) {
	if _, ok := stc.mutation.CreateTime(); !ok {
		v := servicetype.DefaultCreateTime()
		stc.mutation.SetCreateTime(v)
	}
	if _, ok := stc.mutation.UpdateTime(); !ok {
		v := servicetype.DefaultUpdateTime()
		stc.mutation.SetUpdateTime(v)
	}
	if _, ok := stc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if _, ok := stc.mutation.HasCustomer(); !ok {
		v := servicetype.DefaultHasCustomer
		stc.mutation.SetHasCustomer(v)
	}
	if _, ok := stc.mutation.IsDeleted(); !ok {
		v := servicetype.DefaultIsDeleted
		stc.mutation.SetIsDeleted(v)
	}
	if v, ok := stc.mutation.DiscoveryMethod(); ok {
		if err := servicetype.DiscoveryMethodValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"discovery_method\": %v", err)
		}
	}
	var (
		err  error
		node *ServiceType
	)
	if len(stc.hooks) == 0 {
		node, err = stc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ServiceTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			stc.mutation = mutation
			node, err = stc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(stc.hooks) - 1; i >= 0; i-- {
			mut = stc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, stc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (stc *ServiceTypeCreate) SaveX(ctx context.Context) *ServiceType {
	v, err := stc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (stc *ServiceTypeCreate) sqlSave(ctx context.Context) (*ServiceType, error) {
	var (
		st    = &ServiceType{config: stc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: servicetype.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: servicetype.FieldID,
			},
		}
	)
	if value, ok := stc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: servicetype.FieldCreateTime,
		})
		st.CreateTime = value
	}
	if value, ok := stc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: servicetype.FieldUpdateTime,
		})
		st.UpdateTime = value
	}
	if value, ok := stc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: servicetype.FieldName,
		})
		st.Name = value
	}
	if value, ok := stc.mutation.HasCustomer(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: servicetype.FieldHasCustomer,
		})
		st.HasCustomer = value
	}
	if value, ok := stc.mutation.IsDeleted(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: servicetype.FieldIsDeleted,
		})
		st.IsDeleted = value
	}
	if value, ok := stc.mutation.DiscoveryMethod(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: servicetype.FieldDiscoveryMethod,
		})
		st.DiscoveryMethod = value
	}
	if nodes := stc.mutation.ServicesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   servicetype.ServicesTable,
			Columns: []string{servicetype.ServicesColumn},
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
	if nodes := stc.mutation.PropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   servicetype.PropertyTypesTable,
			Columns: []string{servicetype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: propertytype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := stc.mutation.EndpointDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   servicetype.EndpointDefinitionsTable,
			Columns: []string{servicetype.EndpointDefinitionsColumn},
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
	if err := sqlgraph.CreateNode(ctx, stc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	st.ID = int(id)
	return st, nil
}
