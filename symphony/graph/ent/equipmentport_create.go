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
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
)

// EquipmentPortCreate is the builder for creating a EquipmentPort entity.
type EquipmentPortCreate struct {
	config
	mutation *EquipmentPortMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (epc *EquipmentPortCreate) SetCreateTime(t time.Time) *EquipmentPortCreate {
	epc.mutation.SetCreateTime(t)
	return epc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (epc *EquipmentPortCreate) SetNillableCreateTime(t *time.Time) *EquipmentPortCreate {
	if t != nil {
		epc.SetCreateTime(*t)
	}
	return epc
}

// SetUpdateTime sets the update_time field.
func (epc *EquipmentPortCreate) SetUpdateTime(t time.Time) *EquipmentPortCreate {
	epc.mutation.SetUpdateTime(t)
	return epc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (epc *EquipmentPortCreate) SetNillableUpdateTime(t *time.Time) *EquipmentPortCreate {
	if t != nil {
		epc.SetUpdateTime(*t)
	}
	return epc
}

// SetDefinitionID sets the definition edge to EquipmentPortDefinition by id.
func (epc *EquipmentPortCreate) SetDefinitionID(id int) *EquipmentPortCreate {
	epc.mutation.SetDefinitionID(id)
	return epc
}

// SetDefinition sets the definition edge to EquipmentPortDefinition.
func (epc *EquipmentPortCreate) SetDefinition(e *EquipmentPortDefinition) *EquipmentPortCreate {
	return epc.SetDefinitionID(e.ID)
}

// SetParentID sets the parent edge to Equipment by id.
func (epc *EquipmentPortCreate) SetParentID(id int) *EquipmentPortCreate {
	epc.mutation.SetParentID(id)
	return epc
}

// SetNillableParentID sets the parent edge to Equipment by id if the given value is not nil.
func (epc *EquipmentPortCreate) SetNillableParentID(id *int) *EquipmentPortCreate {
	if id != nil {
		epc = epc.SetParentID(*id)
	}
	return epc
}

// SetParent sets the parent edge to Equipment.
func (epc *EquipmentPortCreate) SetParent(e *Equipment) *EquipmentPortCreate {
	return epc.SetParentID(e.ID)
}

// SetLinkID sets the link edge to Link by id.
func (epc *EquipmentPortCreate) SetLinkID(id int) *EquipmentPortCreate {
	epc.mutation.SetLinkID(id)
	return epc
}

// SetNillableLinkID sets the link edge to Link by id if the given value is not nil.
func (epc *EquipmentPortCreate) SetNillableLinkID(id *int) *EquipmentPortCreate {
	if id != nil {
		epc = epc.SetLinkID(*id)
	}
	return epc
}

// SetLink sets the link edge to Link.
func (epc *EquipmentPortCreate) SetLink(l *Link) *EquipmentPortCreate {
	return epc.SetLinkID(l.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (epc *EquipmentPortCreate) AddPropertyIDs(ids ...int) *EquipmentPortCreate {
	epc.mutation.AddPropertyIDs(ids...)
	return epc
}

// AddProperties adds the properties edges to Property.
func (epc *EquipmentPortCreate) AddProperties(p ...*Property) *EquipmentPortCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return epc.AddPropertyIDs(ids...)
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (epc *EquipmentPortCreate) AddEndpointIDs(ids ...int) *EquipmentPortCreate {
	epc.mutation.AddEndpointIDs(ids...)
	return epc
}

// AddEndpoints adds the endpoints edges to ServiceEndpoint.
func (epc *EquipmentPortCreate) AddEndpoints(s ...*ServiceEndpoint) *EquipmentPortCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return epc.AddEndpointIDs(ids...)
}

// Save creates the EquipmentPort in the database.
func (epc *EquipmentPortCreate) Save(ctx context.Context) (*EquipmentPort, error) {
	if _, ok := epc.mutation.CreateTime(); !ok {
		v := equipmentport.DefaultCreateTime()
		epc.mutation.SetCreateTime(v)
	}
	if _, ok := epc.mutation.UpdateTime(); !ok {
		v := equipmentport.DefaultUpdateTime()
		epc.mutation.SetUpdateTime(v)
	}
	if _, ok := epc.mutation.DefinitionID(); !ok {
		return nil, errors.New("ent: missing required edge \"definition\"")
	}
	var (
		err  error
		node *EquipmentPort
	)
	if len(epc.hooks) == 0 {
		node, err = epc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentPortMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			epc.mutation = mutation
			node, err = epc.sqlSave(ctx)
			return node, err
		})
		for i := len(epc.hooks); i > 0; i-- {
			mut = epc.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, epc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (epc *EquipmentPortCreate) SaveX(ctx context.Context) *EquipmentPort {
	v, err := epc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (epc *EquipmentPortCreate) sqlSave(ctx context.Context) (*EquipmentPort, error) {
	var (
		ep    = &EquipmentPort{config: epc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: equipmentport.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentport.FieldID,
			},
		}
	)
	if value, ok := epc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentport.FieldCreateTime,
		})
		ep.CreateTime = value
	}
	if value, ok := epc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentport.FieldUpdateTime,
		})
		ep.UpdateTime = value
	}
	if nodes := epc.mutation.DefinitionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentport.DefinitionTable,
			Columns: []string{equipmentport.DefinitionColumn},
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
	if nodes := epc.mutation.ParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipmentport.ParentTable,
			Columns: []string{equipmentport.ParentColumn},
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
	if nodes := epc.mutation.LinkIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentport.LinkTable,
			Columns: []string{equipmentport.LinkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: link.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := epc.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentport.PropertiesTable,
			Columns: []string{equipmentport.PropertiesColumn},
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
	if nodes := epc.mutation.EndpointsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentport.EndpointsTable,
			Columns: []string{equipmentport.EndpointsColumn},
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
	if err := sqlgraph.CreateNode(ctx, epc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	ep.ID = int(id)
	return ep, nil
}
