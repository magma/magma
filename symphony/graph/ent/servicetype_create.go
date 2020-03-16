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
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
)

// ServiceTypeCreate is the builder for creating a ServiceType entity.
type ServiceTypeCreate struct {
	config
	create_time    *time.Time
	update_time    *time.Time
	name           *string
	has_customer   *bool
	services       map[int]struct{}
	property_types map[int]struct{}
}

// SetCreateTime sets the create_time field.
func (stc *ServiceTypeCreate) SetCreateTime(t time.Time) *ServiceTypeCreate {
	stc.create_time = &t
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
	stc.update_time = &t
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
	stc.name = &s
	return stc
}

// SetHasCustomer sets the has_customer field.
func (stc *ServiceTypeCreate) SetHasCustomer(b bool) *ServiceTypeCreate {
	stc.has_customer = &b
	return stc
}

// SetNillableHasCustomer sets the has_customer field if the given value is not nil.
func (stc *ServiceTypeCreate) SetNillableHasCustomer(b *bool) *ServiceTypeCreate {
	if b != nil {
		stc.SetHasCustomer(*b)
	}
	return stc
}

// AddServiceIDs adds the services edge to Service by ids.
func (stc *ServiceTypeCreate) AddServiceIDs(ids ...int) *ServiceTypeCreate {
	if stc.services == nil {
		stc.services = make(map[int]struct{})
	}
	for i := range ids {
		stc.services[ids[i]] = struct{}{}
	}
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
	if stc.property_types == nil {
		stc.property_types = make(map[int]struct{})
	}
	for i := range ids {
		stc.property_types[ids[i]] = struct{}{}
	}
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

// Save creates the ServiceType in the database.
func (stc *ServiceTypeCreate) Save(ctx context.Context) (*ServiceType, error) {
	if stc.create_time == nil {
		v := servicetype.DefaultCreateTime()
		stc.create_time = &v
	}
	if stc.update_time == nil {
		v := servicetype.DefaultUpdateTime()
		stc.update_time = &v
	}
	if stc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if stc.has_customer == nil {
		v := servicetype.DefaultHasCustomer
		stc.has_customer = &v
	}
	return stc.sqlSave(ctx)
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
	if value := stc.create_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: servicetype.FieldCreateTime,
		})
		st.CreateTime = *value
	}
	if value := stc.update_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: servicetype.FieldUpdateTime,
		})
		st.UpdateTime = *value
	}
	if value := stc.name; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: servicetype.FieldName,
		})
		st.Name = *value
	}
	if value := stc.has_customer; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: servicetype.FieldHasCustomer,
		})
		st.HasCustomer = *value
	}
	if nodes := stc.services; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := stc.property_types; len(nodes) > 0 {
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
		for k, _ := range nodes {
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
