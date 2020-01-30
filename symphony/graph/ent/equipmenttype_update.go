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

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentcategory"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
)

// EquipmentTypeUpdate is the builder for updating EquipmentType entities.
type EquipmentTypeUpdate struct {
	config

	update_time                *time.Time
	name                       *string
	port_definitions           map[string]struct{}
	position_definitions       map[string]struct{}
	property_types             map[string]struct{}
	equipment                  map[string]struct{}
	category                   map[string]struct{}
	removedPortDefinitions     map[string]struct{}
	removedPositionDefinitions map[string]struct{}
	removedPropertyTypes       map[string]struct{}
	removedEquipment           map[string]struct{}
	clearedCategory            bool
	predicates                 []predicate.EquipmentType
}

// Where adds a new predicate for the builder.
func (etu *EquipmentTypeUpdate) Where(ps ...predicate.EquipmentType) *EquipmentTypeUpdate {
	etu.predicates = append(etu.predicates, ps...)
	return etu
}

// SetName sets the name field.
func (etu *EquipmentTypeUpdate) SetName(s string) *EquipmentTypeUpdate {
	etu.name = &s
	return etu
}

// AddPortDefinitionIDs adds the port_definitions edge to EquipmentPortDefinition by ids.
func (etu *EquipmentTypeUpdate) AddPortDefinitionIDs(ids ...string) *EquipmentTypeUpdate {
	if etu.port_definitions == nil {
		etu.port_definitions = make(map[string]struct{})
	}
	for i := range ids {
		etu.port_definitions[ids[i]] = struct{}{}
	}
	return etu
}

// AddPortDefinitions adds the port_definitions edges to EquipmentPortDefinition.
func (etu *EquipmentTypeUpdate) AddPortDefinitions(e ...*EquipmentPortDefinition) *EquipmentTypeUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.AddPortDefinitionIDs(ids...)
}

// AddPositionDefinitionIDs adds the position_definitions edge to EquipmentPositionDefinition by ids.
func (etu *EquipmentTypeUpdate) AddPositionDefinitionIDs(ids ...string) *EquipmentTypeUpdate {
	if etu.position_definitions == nil {
		etu.position_definitions = make(map[string]struct{})
	}
	for i := range ids {
		etu.position_definitions[ids[i]] = struct{}{}
	}
	return etu
}

// AddPositionDefinitions adds the position_definitions edges to EquipmentPositionDefinition.
func (etu *EquipmentTypeUpdate) AddPositionDefinitions(e ...*EquipmentPositionDefinition) *EquipmentTypeUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.AddPositionDefinitionIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (etu *EquipmentTypeUpdate) AddPropertyTypeIDs(ids ...string) *EquipmentTypeUpdate {
	if etu.property_types == nil {
		etu.property_types = make(map[string]struct{})
	}
	for i := range ids {
		etu.property_types[ids[i]] = struct{}{}
	}
	return etu
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (etu *EquipmentTypeUpdate) AddPropertyTypes(p ...*PropertyType) *EquipmentTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return etu.AddPropertyTypeIDs(ids...)
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (etu *EquipmentTypeUpdate) AddEquipmentIDs(ids ...string) *EquipmentTypeUpdate {
	if etu.equipment == nil {
		etu.equipment = make(map[string]struct{})
	}
	for i := range ids {
		etu.equipment[ids[i]] = struct{}{}
	}
	return etu
}

// AddEquipment adds the equipment edges to Equipment.
func (etu *EquipmentTypeUpdate) AddEquipment(e ...*Equipment) *EquipmentTypeUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.AddEquipmentIDs(ids...)
}

// SetCategoryID sets the category edge to EquipmentCategory by id.
func (etu *EquipmentTypeUpdate) SetCategoryID(id string) *EquipmentTypeUpdate {
	if etu.category == nil {
		etu.category = make(map[string]struct{})
	}
	etu.category[id] = struct{}{}
	return etu
}

// SetNillableCategoryID sets the category edge to EquipmentCategory by id if the given value is not nil.
func (etu *EquipmentTypeUpdate) SetNillableCategoryID(id *string) *EquipmentTypeUpdate {
	if id != nil {
		etu = etu.SetCategoryID(*id)
	}
	return etu
}

// SetCategory sets the category edge to EquipmentCategory.
func (etu *EquipmentTypeUpdate) SetCategory(e *EquipmentCategory) *EquipmentTypeUpdate {
	return etu.SetCategoryID(e.ID)
}

// RemovePortDefinitionIDs removes the port_definitions edge to EquipmentPortDefinition by ids.
func (etu *EquipmentTypeUpdate) RemovePortDefinitionIDs(ids ...string) *EquipmentTypeUpdate {
	if etu.removedPortDefinitions == nil {
		etu.removedPortDefinitions = make(map[string]struct{})
	}
	for i := range ids {
		etu.removedPortDefinitions[ids[i]] = struct{}{}
	}
	return etu
}

// RemovePortDefinitions removes port_definitions edges to EquipmentPortDefinition.
func (etu *EquipmentTypeUpdate) RemovePortDefinitions(e ...*EquipmentPortDefinition) *EquipmentTypeUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.RemovePortDefinitionIDs(ids...)
}

// RemovePositionDefinitionIDs removes the position_definitions edge to EquipmentPositionDefinition by ids.
func (etu *EquipmentTypeUpdate) RemovePositionDefinitionIDs(ids ...string) *EquipmentTypeUpdate {
	if etu.removedPositionDefinitions == nil {
		etu.removedPositionDefinitions = make(map[string]struct{})
	}
	for i := range ids {
		etu.removedPositionDefinitions[ids[i]] = struct{}{}
	}
	return etu
}

// RemovePositionDefinitions removes position_definitions edges to EquipmentPositionDefinition.
func (etu *EquipmentTypeUpdate) RemovePositionDefinitions(e ...*EquipmentPositionDefinition) *EquipmentTypeUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.RemovePositionDefinitionIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (etu *EquipmentTypeUpdate) RemovePropertyTypeIDs(ids ...string) *EquipmentTypeUpdate {
	if etu.removedPropertyTypes == nil {
		etu.removedPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		etu.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return etu
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (etu *EquipmentTypeUpdate) RemovePropertyTypes(p ...*PropertyType) *EquipmentTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return etu.RemovePropertyTypeIDs(ids...)
}

// RemoveEquipmentIDs removes the equipment edge to Equipment by ids.
func (etu *EquipmentTypeUpdate) RemoveEquipmentIDs(ids ...string) *EquipmentTypeUpdate {
	if etu.removedEquipment == nil {
		etu.removedEquipment = make(map[string]struct{})
	}
	for i := range ids {
		etu.removedEquipment[ids[i]] = struct{}{}
	}
	return etu
}

// RemoveEquipment removes equipment edges to Equipment.
func (etu *EquipmentTypeUpdate) RemoveEquipment(e ...*Equipment) *EquipmentTypeUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.RemoveEquipmentIDs(ids...)
}

// ClearCategory clears the category edge to EquipmentCategory.
func (etu *EquipmentTypeUpdate) ClearCategory() *EquipmentTypeUpdate {
	etu.clearedCategory = true
	return etu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (etu *EquipmentTypeUpdate) Save(ctx context.Context) (int, error) {
	if etu.update_time == nil {
		v := equipmenttype.UpdateDefaultUpdateTime()
		etu.update_time = &v
	}
	if len(etu.category) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"category\"")
	}
	return etu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (etu *EquipmentTypeUpdate) SaveX(ctx context.Context) int {
	affected, err := etu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (etu *EquipmentTypeUpdate) Exec(ctx context.Context) error {
	_, err := etu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (etu *EquipmentTypeUpdate) ExecX(ctx context.Context) {
	if err := etu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (etu *EquipmentTypeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmenttype.Table,
			Columns: equipmenttype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: equipmenttype.FieldID,
			},
		},
	}
	if ps := etu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := etu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmenttype.FieldUpdateTime,
		})
	}
	if value := etu.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipmenttype.FieldName,
		})
	}
	if nodes := etu.removedPortDefinitions; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etu.port_definitions; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := etu.removedPositionDefinitions; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etu.position_definitions; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := etu.removedPropertyTypes; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etu.property_types; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := etu.removedEquipment; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etu.equipment; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if etu.clearedCategory {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etu.category; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, etu.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// EquipmentTypeUpdateOne is the builder for updating a single EquipmentType entity.
type EquipmentTypeUpdateOne struct {
	config
	id string

	update_time                *time.Time
	name                       *string
	port_definitions           map[string]struct{}
	position_definitions       map[string]struct{}
	property_types             map[string]struct{}
	equipment                  map[string]struct{}
	category                   map[string]struct{}
	removedPortDefinitions     map[string]struct{}
	removedPositionDefinitions map[string]struct{}
	removedPropertyTypes       map[string]struct{}
	removedEquipment           map[string]struct{}
	clearedCategory            bool
}

// SetName sets the name field.
func (etuo *EquipmentTypeUpdateOne) SetName(s string) *EquipmentTypeUpdateOne {
	etuo.name = &s
	return etuo
}

// AddPortDefinitionIDs adds the port_definitions edge to EquipmentPortDefinition by ids.
func (etuo *EquipmentTypeUpdateOne) AddPortDefinitionIDs(ids ...string) *EquipmentTypeUpdateOne {
	if etuo.port_definitions == nil {
		etuo.port_definitions = make(map[string]struct{})
	}
	for i := range ids {
		etuo.port_definitions[ids[i]] = struct{}{}
	}
	return etuo
}

// AddPortDefinitions adds the port_definitions edges to EquipmentPortDefinition.
func (etuo *EquipmentTypeUpdateOne) AddPortDefinitions(e ...*EquipmentPortDefinition) *EquipmentTypeUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.AddPortDefinitionIDs(ids...)
}

// AddPositionDefinitionIDs adds the position_definitions edge to EquipmentPositionDefinition by ids.
func (etuo *EquipmentTypeUpdateOne) AddPositionDefinitionIDs(ids ...string) *EquipmentTypeUpdateOne {
	if etuo.position_definitions == nil {
		etuo.position_definitions = make(map[string]struct{})
	}
	for i := range ids {
		etuo.position_definitions[ids[i]] = struct{}{}
	}
	return etuo
}

// AddPositionDefinitions adds the position_definitions edges to EquipmentPositionDefinition.
func (etuo *EquipmentTypeUpdateOne) AddPositionDefinitions(e ...*EquipmentPositionDefinition) *EquipmentTypeUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.AddPositionDefinitionIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (etuo *EquipmentTypeUpdateOne) AddPropertyTypeIDs(ids ...string) *EquipmentTypeUpdateOne {
	if etuo.property_types == nil {
		etuo.property_types = make(map[string]struct{})
	}
	for i := range ids {
		etuo.property_types[ids[i]] = struct{}{}
	}
	return etuo
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (etuo *EquipmentTypeUpdateOne) AddPropertyTypes(p ...*PropertyType) *EquipmentTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return etuo.AddPropertyTypeIDs(ids...)
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (etuo *EquipmentTypeUpdateOne) AddEquipmentIDs(ids ...string) *EquipmentTypeUpdateOne {
	if etuo.equipment == nil {
		etuo.equipment = make(map[string]struct{})
	}
	for i := range ids {
		etuo.equipment[ids[i]] = struct{}{}
	}
	return etuo
}

// AddEquipment adds the equipment edges to Equipment.
func (etuo *EquipmentTypeUpdateOne) AddEquipment(e ...*Equipment) *EquipmentTypeUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.AddEquipmentIDs(ids...)
}

// SetCategoryID sets the category edge to EquipmentCategory by id.
func (etuo *EquipmentTypeUpdateOne) SetCategoryID(id string) *EquipmentTypeUpdateOne {
	if etuo.category == nil {
		etuo.category = make(map[string]struct{})
	}
	etuo.category[id] = struct{}{}
	return etuo
}

// SetNillableCategoryID sets the category edge to EquipmentCategory by id if the given value is not nil.
func (etuo *EquipmentTypeUpdateOne) SetNillableCategoryID(id *string) *EquipmentTypeUpdateOne {
	if id != nil {
		etuo = etuo.SetCategoryID(*id)
	}
	return etuo
}

// SetCategory sets the category edge to EquipmentCategory.
func (etuo *EquipmentTypeUpdateOne) SetCategory(e *EquipmentCategory) *EquipmentTypeUpdateOne {
	return etuo.SetCategoryID(e.ID)
}

// RemovePortDefinitionIDs removes the port_definitions edge to EquipmentPortDefinition by ids.
func (etuo *EquipmentTypeUpdateOne) RemovePortDefinitionIDs(ids ...string) *EquipmentTypeUpdateOne {
	if etuo.removedPortDefinitions == nil {
		etuo.removedPortDefinitions = make(map[string]struct{})
	}
	for i := range ids {
		etuo.removedPortDefinitions[ids[i]] = struct{}{}
	}
	return etuo
}

// RemovePortDefinitions removes port_definitions edges to EquipmentPortDefinition.
func (etuo *EquipmentTypeUpdateOne) RemovePortDefinitions(e ...*EquipmentPortDefinition) *EquipmentTypeUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.RemovePortDefinitionIDs(ids...)
}

// RemovePositionDefinitionIDs removes the position_definitions edge to EquipmentPositionDefinition by ids.
func (etuo *EquipmentTypeUpdateOne) RemovePositionDefinitionIDs(ids ...string) *EquipmentTypeUpdateOne {
	if etuo.removedPositionDefinitions == nil {
		etuo.removedPositionDefinitions = make(map[string]struct{})
	}
	for i := range ids {
		etuo.removedPositionDefinitions[ids[i]] = struct{}{}
	}
	return etuo
}

// RemovePositionDefinitions removes position_definitions edges to EquipmentPositionDefinition.
func (etuo *EquipmentTypeUpdateOne) RemovePositionDefinitions(e ...*EquipmentPositionDefinition) *EquipmentTypeUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.RemovePositionDefinitionIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (etuo *EquipmentTypeUpdateOne) RemovePropertyTypeIDs(ids ...string) *EquipmentTypeUpdateOne {
	if etuo.removedPropertyTypes == nil {
		etuo.removedPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		etuo.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return etuo
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (etuo *EquipmentTypeUpdateOne) RemovePropertyTypes(p ...*PropertyType) *EquipmentTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return etuo.RemovePropertyTypeIDs(ids...)
}

// RemoveEquipmentIDs removes the equipment edge to Equipment by ids.
func (etuo *EquipmentTypeUpdateOne) RemoveEquipmentIDs(ids ...string) *EquipmentTypeUpdateOne {
	if etuo.removedEquipment == nil {
		etuo.removedEquipment = make(map[string]struct{})
	}
	for i := range ids {
		etuo.removedEquipment[ids[i]] = struct{}{}
	}
	return etuo
}

// RemoveEquipment removes equipment edges to Equipment.
func (etuo *EquipmentTypeUpdateOne) RemoveEquipment(e ...*Equipment) *EquipmentTypeUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.RemoveEquipmentIDs(ids...)
}

// ClearCategory clears the category edge to EquipmentCategory.
func (etuo *EquipmentTypeUpdateOne) ClearCategory() *EquipmentTypeUpdateOne {
	etuo.clearedCategory = true
	return etuo
}

// Save executes the query and returns the updated entity.
func (etuo *EquipmentTypeUpdateOne) Save(ctx context.Context) (*EquipmentType, error) {
	if etuo.update_time == nil {
		v := equipmenttype.UpdateDefaultUpdateTime()
		etuo.update_time = &v
	}
	if len(etuo.category) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"category\"")
	}
	return etuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (etuo *EquipmentTypeUpdateOne) SaveX(ctx context.Context) *EquipmentType {
	et, err := etuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return et
}

// Exec executes the query on the entity.
func (etuo *EquipmentTypeUpdateOne) Exec(ctx context.Context) error {
	_, err := etuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (etuo *EquipmentTypeUpdateOne) ExecX(ctx context.Context) {
	if err := etuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (etuo *EquipmentTypeUpdateOne) sqlSave(ctx context.Context) (et *EquipmentType, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmenttype.Table,
			Columns: equipmenttype.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  etuo.id,
				Type:   field.TypeString,
				Column: equipmenttype.FieldID,
			},
		},
	}
	if value := etuo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmenttype.FieldUpdateTime,
		})
	}
	if value := etuo.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipmenttype.FieldName,
		})
	}
	if nodes := etuo.removedPortDefinitions; len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etuo.port_definitions; len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := etuo.removedPositionDefinitions; len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etuo.position_definitions; len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := etuo.removedPropertyTypes; len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etuo.property_types; len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := etuo.removedEquipment; len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etuo.equipment; len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if etuo.clearedCategory {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etuo.category; len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	et = &EquipmentType{config: etuo.config}
	_spec.Assign = et.assignValues
	_spec.ScanValues = et.scanValues()
	if err = sqlgraph.UpdateNode(ctx, etuo.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return et, nil
}
