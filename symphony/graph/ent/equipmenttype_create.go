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
	"github.com/facebookincubator/symphony/graph/ent/equipmentcategory"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
)

// EquipmentTypeCreate is the builder for creating a EquipmentType entity.
type EquipmentTypeCreate struct {
	config
	create_time          *time.Time
	update_time          *time.Time
	name                 *string
	port_definitions     map[string]struct{}
	position_definitions map[string]struct{}
	property_types       map[string]struct{}
	equipment            map[string]struct{}
	category             map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (etc *EquipmentTypeCreate) SetCreateTime(t time.Time) *EquipmentTypeCreate {
	etc.create_time = &t
	return etc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (etc *EquipmentTypeCreate) SetNillableCreateTime(t *time.Time) *EquipmentTypeCreate {
	if t != nil {
		etc.SetCreateTime(*t)
	}
	return etc
}

// SetUpdateTime sets the update_time field.
func (etc *EquipmentTypeCreate) SetUpdateTime(t time.Time) *EquipmentTypeCreate {
	etc.update_time = &t
	return etc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (etc *EquipmentTypeCreate) SetNillableUpdateTime(t *time.Time) *EquipmentTypeCreate {
	if t != nil {
		etc.SetUpdateTime(*t)
	}
	return etc
}

// SetName sets the name field.
func (etc *EquipmentTypeCreate) SetName(s string) *EquipmentTypeCreate {
	etc.name = &s
	return etc
}

// AddPortDefinitionIDs adds the port_definitions edge to EquipmentPortDefinition by ids.
func (etc *EquipmentTypeCreate) AddPortDefinitionIDs(ids ...string) *EquipmentTypeCreate {
	if etc.port_definitions == nil {
		etc.port_definitions = make(map[string]struct{})
	}
	for i := range ids {
		etc.port_definitions[ids[i]] = struct{}{}
	}
	return etc
}

// AddPortDefinitions adds the port_definitions edges to EquipmentPortDefinition.
func (etc *EquipmentTypeCreate) AddPortDefinitions(e ...*EquipmentPortDefinition) *EquipmentTypeCreate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etc.AddPortDefinitionIDs(ids...)
}

// AddPositionDefinitionIDs adds the position_definitions edge to EquipmentPositionDefinition by ids.
func (etc *EquipmentTypeCreate) AddPositionDefinitionIDs(ids ...string) *EquipmentTypeCreate {
	if etc.position_definitions == nil {
		etc.position_definitions = make(map[string]struct{})
	}
	for i := range ids {
		etc.position_definitions[ids[i]] = struct{}{}
	}
	return etc
}

// AddPositionDefinitions adds the position_definitions edges to EquipmentPositionDefinition.
func (etc *EquipmentTypeCreate) AddPositionDefinitions(e ...*EquipmentPositionDefinition) *EquipmentTypeCreate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etc.AddPositionDefinitionIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (etc *EquipmentTypeCreate) AddPropertyTypeIDs(ids ...string) *EquipmentTypeCreate {
	if etc.property_types == nil {
		etc.property_types = make(map[string]struct{})
	}
	for i := range ids {
		etc.property_types[ids[i]] = struct{}{}
	}
	return etc
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (etc *EquipmentTypeCreate) AddPropertyTypes(p ...*PropertyType) *EquipmentTypeCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return etc.AddPropertyTypeIDs(ids...)
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (etc *EquipmentTypeCreate) AddEquipmentIDs(ids ...string) *EquipmentTypeCreate {
	if etc.equipment == nil {
		etc.equipment = make(map[string]struct{})
	}
	for i := range ids {
		etc.equipment[ids[i]] = struct{}{}
	}
	return etc
}

// AddEquipment adds the equipment edges to Equipment.
func (etc *EquipmentTypeCreate) AddEquipment(e ...*Equipment) *EquipmentTypeCreate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etc.AddEquipmentIDs(ids...)
}

// SetCategoryID sets the category edge to EquipmentCategory by id.
func (etc *EquipmentTypeCreate) SetCategoryID(id string) *EquipmentTypeCreate {
	if etc.category == nil {
		etc.category = make(map[string]struct{})
	}
	etc.category[id] = struct{}{}
	return etc
}

// SetNillableCategoryID sets the category edge to EquipmentCategory by id if the given value is not nil.
func (etc *EquipmentTypeCreate) SetNillableCategoryID(id *string) *EquipmentTypeCreate {
	if id != nil {
		etc = etc.SetCategoryID(*id)
	}
	return etc
}

// SetCategory sets the category edge to EquipmentCategory.
func (etc *EquipmentTypeCreate) SetCategory(e *EquipmentCategory) *EquipmentTypeCreate {
	return etc.SetCategoryID(e.ID)
}

// Save creates the EquipmentType in the database.
func (etc *EquipmentTypeCreate) Save(ctx context.Context) (*EquipmentType, error) {
	if etc.create_time == nil {
		v := equipmenttype.DefaultCreateTime()
		etc.create_time = &v
	}
	if etc.update_time == nil {
		v := equipmenttype.DefaultUpdateTime()
		etc.update_time = &v
	}
	if etc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if len(etc.category) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"category\"")
	}
	return etc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (etc *EquipmentTypeCreate) SaveX(ctx context.Context) *EquipmentType {
	v, err := etc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (etc *EquipmentTypeCreate) sqlSave(ctx context.Context) (*EquipmentType, error) {
	var (
		et   = &EquipmentType{config: etc.config}
		spec = &sqlgraph.CreateSpec{
			Table: equipmenttype.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: equipmenttype.FieldID,
			},
		}
	)
	if value := etc.create_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmenttype.FieldCreateTime,
		})
		et.CreateTime = *value
	}
	if value := etc.update_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmenttype.FieldUpdateTime,
		})
		et.UpdateTime = *value
	}
	if value := etc.name; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipmenttype.FieldName,
		})
		et.Name = *value
	}
	if nodes := etc.port_definitions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.PortDefinitionsTable,
			Columns: []string{equipmenttype.PortDefinitionsColumn},
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
	if nodes := etc.position_definitions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.PositionDefinitionsTable,
			Columns: []string{equipmenttype.PositionDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentpositiondefinition.FieldID,
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
	if nodes := etc.property_types; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.PropertyTypesTable,
			Columns: []string{equipmenttype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: propertytype.FieldID,
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
	if nodes := etc.equipment; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmenttype.EquipmentTable,
			Columns: []string{equipmenttype.EquipmentColumn},
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
	if nodes := etc.category; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmenttype.CategoryTable,
			Columns: []string{equipmenttype.CategoryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentcategory.FieldID,
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
	if err := sqlgraph.CreateNode(ctx, etc.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := spec.ID.Value.(int64)
	et.ID = strconv.FormatInt(id, 10)
	return et, nil
}
