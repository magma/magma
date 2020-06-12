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
	"github.com/facebookincubator/symphony/pkg/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/pkg/ent/propertytype"
)

// EquipmentPortTypeCreate is the builder for creating a EquipmentPortType entity.
type EquipmentPortTypeCreate struct {
	config
	mutation *EquipmentPortTypeMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (eptc *EquipmentPortTypeCreate) SetCreateTime(t time.Time) *EquipmentPortTypeCreate {
	eptc.mutation.SetCreateTime(t)
	return eptc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (eptc *EquipmentPortTypeCreate) SetNillableCreateTime(t *time.Time) *EquipmentPortTypeCreate {
	if t != nil {
		eptc.SetCreateTime(*t)
	}
	return eptc
}

// SetUpdateTime sets the update_time field.
func (eptc *EquipmentPortTypeCreate) SetUpdateTime(t time.Time) *EquipmentPortTypeCreate {
	eptc.mutation.SetUpdateTime(t)
	return eptc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (eptc *EquipmentPortTypeCreate) SetNillableUpdateTime(t *time.Time) *EquipmentPortTypeCreate {
	if t != nil {
		eptc.SetUpdateTime(*t)
	}
	return eptc
}

// SetName sets the name field.
func (eptc *EquipmentPortTypeCreate) SetName(s string) *EquipmentPortTypeCreate {
	eptc.mutation.SetName(s)
	return eptc
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (eptc *EquipmentPortTypeCreate) AddPropertyTypeIDs(ids ...int) *EquipmentPortTypeCreate {
	eptc.mutation.AddPropertyTypeIDs(ids...)
	return eptc
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (eptc *EquipmentPortTypeCreate) AddPropertyTypes(p ...*PropertyType) *EquipmentPortTypeCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptc.AddPropertyTypeIDs(ids...)
}

// AddLinkPropertyTypeIDs adds the link_property_types edge to PropertyType by ids.
func (eptc *EquipmentPortTypeCreate) AddLinkPropertyTypeIDs(ids ...int) *EquipmentPortTypeCreate {
	eptc.mutation.AddLinkPropertyTypeIDs(ids...)
	return eptc
}

// AddLinkPropertyTypes adds the link_property_types edges to PropertyType.
func (eptc *EquipmentPortTypeCreate) AddLinkPropertyTypes(p ...*PropertyType) *EquipmentPortTypeCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptc.AddLinkPropertyTypeIDs(ids...)
}

// AddPortDefinitionIDs adds the port_definitions edge to EquipmentPortDefinition by ids.
func (eptc *EquipmentPortTypeCreate) AddPortDefinitionIDs(ids ...int) *EquipmentPortTypeCreate {
	eptc.mutation.AddPortDefinitionIDs(ids...)
	return eptc
}

// AddPortDefinitions adds the port_definitions edges to EquipmentPortDefinition.
func (eptc *EquipmentPortTypeCreate) AddPortDefinitions(e ...*EquipmentPortDefinition) *EquipmentPortTypeCreate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eptc.AddPortDefinitionIDs(ids...)
}

// Save creates the EquipmentPortType in the database.
func (eptc *EquipmentPortTypeCreate) Save(ctx context.Context) (*EquipmentPortType, error) {
	if _, ok := eptc.mutation.CreateTime(); !ok {
		v := equipmentporttype.DefaultCreateTime()
		eptc.mutation.SetCreateTime(v)
	}
	if _, ok := eptc.mutation.UpdateTime(); !ok {
		v := equipmentporttype.DefaultUpdateTime()
		eptc.mutation.SetUpdateTime(v)
	}
	if _, ok := eptc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	var (
		err  error
		node *EquipmentPortType
	)
	if len(eptc.hooks) == 0 {
		node, err = eptc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentPortTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			eptc.mutation = mutation
			node, err = eptc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(eptc.hooks) - 1; i >= 0; i-- {
			mut = eptc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, eptc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (eptc *EquipmentPortTypeCreate) SaveX(ctx context.Context) *EquipmentPortType {
	v, err := eptc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (eptc *EquipmentPortTypeCreate) sqlSave(ctx context.Context) (*EquipmentPortType, error) {
	var (
		ept   = &EquipmentPortType{config: eptc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: equipmentporttype.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentporttype.FieldID,
			},
		}
	)
	if value, ok := eptc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentporttype.FieldCreateTime,
		})
		ept.CreateTime = value
	}
	if value, ok := eptc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentporttype.FieldUpdateTime,
		})
		ept.UpdateTime = value
	}
	if value, ok := eptc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentporttype.FieldName,
		})
		ept.Name = value
	}
	if nodes := eptc.mutation.PropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.PropertyTypesTable,
			Columns: []string{equipmentporttype.PropertyTypesColumn},
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
	if nodes := eptc.mutation.LinkPropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.LinkPropertyTypesTable,
			Columns: []string{equipmentporttype.LinkPropertyTypesColumn},
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
	if nodes := eptc.mutation.PortDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentporttype.PortDefinitionsTable,
			Columns: []string{equipmentporttype.PortDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentportdefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, eptc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	ept.ID = int(id)
	return ept, nil
}
