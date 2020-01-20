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
	create_time *time.Time
	update_time *time.Time
	definition  map[string]struct{}
	parent      map[string]struct{}
	link        map[string]struct{}
	properties  map[string]struct{}
	endpoints   map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (epc *EquipmentPortCreate) SetCreateTime(t time.Time) *EquipmentPortCreate {
	epc.create_time = &t
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
	epc.update_time = &t
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
func (epc *EquipmentPortCreate) SetDefinitionID(id string) *EquipmentPortCreate {
	if epc.definition == nil {
		epc.definition = make(map[string]struct{})
	}
	epc.definition[id] = struct{}{}
	return epc
}

// SetDefinition sets the definition edge to EquipmentPortDefinition.
func (epc *EquipmentPortCreate) SetDefinition(e *EquipmentPortDefinition) *EquipmentPortCreate {
	return epc.SetDefinitionID(e.ID)
}

// SetParentID sets the parent edge to Equipment by id.
func (epc *EquipmentPortCreate) SetParentID(id string) *EquipmentPortCreate {
	if epc.parent == nil {
		epc.parent = make(map[string]struct{})
	}
	epc.parent[id] = struct{}{}
	return epc
}

// SetNillableParentID sets the parent edge to Equipment by id if the given value is not nil.
func (epc *EquipmentPortCreate) SetNillableParentID(id *string) *EquipmentPortCreate {
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
func (epc *EquipmentPortCreate) SetLinkID(id string) *EquipmentPortCreate {
	if epc.link == nil {
		epc.link = make(map[string]struct{})
	}
	epc.link[id] = struct{}{}
	return epc
}

// SetNillableLinkID sets the link edge to Link by id if the given value is not nil.
func (epc *EquipmentPortCreate) SetNillableLinkID(id *string) *EquipmentPortCreate {
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
func (epc *EquipmentPortCreate) AddPropertyIDs(ids ...string) *EquipmentPortCreate {
	if epc.properties == nil {
		epc.properties = make(map[string]struct{})
	}
	for i := range ids {
		epc.properties[ids[i]] = struct{}{}
	}
	return epc
}

// AddProperties adds the properties edges to Property.
func (epc *EquipmentPortCreate) AddProperties(p ...*Property) *EquipmentPortCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return epc.AddPropertyIDs(ids...)
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (epc *EquipmentPortCreate) AddEndpointIDs(ids ...string) *EquipmentPortCreate {
	if epc.endpoints == nil {
		epc.endpoints = make(map[string]struct{})
	}
	for i := range ids {
		epc.endpoints[ids[i]] = struct{}{}
	}
	return epc
}

// AddEndpoints adds the endpoints edges to ServiceEndpoint.
func (epc *EquipmentPortCreate) AddEndpoints(s ...*ServiceEndpoint) *EquipmentPortCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return epc.AddEndpointIDs(ids...)
}

// Save creates the EquipmentPort in the database.
func (epc *EquipmentPortCreate) Save(ctx context.Context) (*EquipmentPort, error) {
	if epc.create_time == nil {
		v := equipmentport.DefaultCreateTime()
		epc.create_time = &v
	}
	if epc.update_time == nil {
		v := equipmentport.DefaultUpdateTime()
		epc.update_time = &v
	}
	if len(epc.definition) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"definition\"")
	}
	if epc.definition == nil {
		return nil, errors.New("ent: missing required edge \"definition\"")
	}
	if len(epc.parent) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"parent\"")
	}
	if len(epc.link) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"link\"")
	}
	return epc.sqlSave(ctx)
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
		ep   = &EquipmentPort{config: epc.config}
		spec = &sqlgraph.CreateSpec{
			Table: equipmentport.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: equipmentport.FieldID,
			},
		}
	)
	if value := epc.create_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmentport.FieldCreateTime,
		})
		ep.CreateTime = *value
	}
	if value := epc.update_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmentport.FieldUpdateTime,
		})
		ep.UpdateTime = *value
	}
	if nodes := epc.definition; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentport.DefinitionTable,
			Columns: []string{equipmentport.DefinitionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentportdefinition.FieldID,
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
	if nodes := epc.parent; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipmentport.ParentTable,
			Columns: []string{equipmentport.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipment.FieldID,
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
	if nodes := epc.link; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentport.LinkTable,
			Columns: []string{equipmentport.LinkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: link.FieldID,
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
	if nodes := epc.properties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentport.PropertiesTable,
			Columns: []string{equipmentport.PropertiesColumn},
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
	if nodes := epc.endpoints; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentport.EndpointsTable,
			Columns: []string{equipmentport.EndpointsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: serviceendpoint.FieldID,
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
	if err := sqlgraph.CreateNode(ctx, epc.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := spec.ID.Value.(int64)
	ep.ID = strconv.FormatInt(id, 10)
	return ep, nil
}
