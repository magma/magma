// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
)

// EquipmentPortTypeUpdate is the builder for updating EquipmentPortType entities.
type EquipmentPortTypeUpdate struct {
	config

	update_time              *time.Time
	name                     *string
	property_types           map[string]struct{}
	link_property_types      map[string]struct{}
	port_definitions         map[string]struct{}
	removedPropertyTypes     map[string]struct{}
	removedLinkPropertyTypes map[string]struct{}
	removedPortDefinitions   map[string]struct{}
	predicates               []predicate.EquipmentPortType
}

// Where adds a new predicate for the builder.
func (eptu *EquipmentPortTypeUpdate) Where(ps ...predicate.EquipmentPortType) *EquipmentPortTypeUpdate {
	eptu.predicates = append(eptu.predicates, ps...)
	return eptu
}

// SetName sets the name field.
func (eptu *EquipmentPortTypeUpdate) SetName(s string) *EquipmentPortTypeUpdate {
	eptu.name = &s
	return eptu
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (eptu *EquipmentPortTypeUpdate) AddPropertyTypeIDs(ids ...string) *EquipmentPortTypeUpdate {
	if eptu.property_types == nil {
		eptu.property_types = make(map[string]struct{})
	}
	for i := range ids {
		eptu.property_types[ids[i]] = struct{}{}
	}
	return eptu
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (eptu *EquipmentPortTypeUpdate) AddPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptu.AddPropertyTypeIDs(ids...)
}

// AddLinkPropertyTypeIDs adds the link_property_types edge to PropertyType by ids.
func (eptu *EquipmentPortTypeUpdate) AddLinkPropertyTypeIDs(ids ...string) *EquipmentPortTypeUpdate {
	if eptu.link_property_types == nil {
		eptu.link_property_types = make(map[string]struct{})
	}
	for i := range ids {
		eptu.link_property_types[ids[i]] = struct{}{}
	}
	return eptu
}

// AddLinkPropertyTypes adds the link_property_types edges to PropertyType.
func (eptu *EquipmentPortTypeUpdate) AddLinkPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptu.AddLinkPropertyTypeIDs(ids...)
}

// AddPortDefinitionIDs adds the port_definitions edge to EquipmentPortDefinition by ids.
func (eptu *EquipmentPortTypeUpdate) AddPortDefinitionIDs(ids ...string) *EquipmentPortTypeUpdate {
	if eptu.port_definitions == nil {
		eptu.port_definitions = make(map[string]struct{})
	}
	for i := range ids {
		eptu.port_definitions[ids[i]] = struct{}{}
	}
	return eptu
}

// AddPortDefinitions adds the port_definitions edges to EquipmentPortDefinition.
func (eptu *EquipmentPortTypeUpdate) AddPortDefinitions(e ...*EquipmentPortDefinition) *EquipmentPortTypeUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eptu.AddPortDefinitionIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (eptu *EquipmentPortTypeUpdate) RemovePropertyTypeIDs(ids ...string) *EquipmentPortTypeUpdate {
	if eptu.removedPropertyTypes == nil {
		eptu.removedPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		eptu.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return eptu
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (eptu *EquipmentPortTypeUpdate) RemovePropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptu.RemovePropertyTypeIDs(ids...)
}

// RemoveLinkPropertyTypeIDs removes the link_property_types edge to PropertyType by ids.
func (eptu *EquipmentPortTypeUpdate) RemoveLinkPropertyTypeIDs(ids ...string) *EquipmentPortTypeUpdate {
	if eptu.removedLinkPropertyTypes == nil {
		eptu.removedLinkPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		eptu.removedLinkPropertyTypes[ids[i]] = struct{}{}
	}
	return eptu
}

// RemoveLinkPropertyTypes removes link_property_types edges to PropertyType.
func (eptu *EquipmentPortTypeUpdate) RemoveLinkPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptu.RemoveLinkPropertyTypeIDs(ids...)
}

// RemovePortDefinitionIDs removes the port_definitions edge to EquipmentPortDefinition by ids.
func (eptu *EquipmentPortTypeUpdate) RemovePortDefinitionIDs(ids ...string) *EquipmentPortTypeUpdate {
	if eptu.removedPortDefinitions == nil {
		eptu.removedPortDefinitions = make(map[string]struct{})
	}
	for i := range ids {
		eptu.removedPortDefinitions[ids[i]] = struct{}{}
	}
	return eptu
}

// RemovePortDefinitions removes port_definitions edges to EquipmentPortDefinition.
func (eptu *EquipmentPortTypeUpdate) RemovePortDefinitions(e ...*EquipmentPortDefinition) *EquipmentPortTypeUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eptu.RemovePortDefinitionIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (eptu *EquipmentPortTypeUpdate) Save(ctx context.Context) (int, error) {
	if eptu.update_time == nil {
		v := equipmentporttype.UpdateDefaultUpdateTime()
		eptu.update_time = &v
	}
	return eptu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (eptu *EquipmentPortTypeUpdate) SaveX(ctx context.Context) int {
	affected, err := eptu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (eptu *EquipmentPortTypeUpdate) Exec(ctx context.Context) error {
	_, err := eptu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eptu *EquipmentPortTypeUpdate) ExecX(ctx context.Context) {
	if err := eptu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (eptu *EquipmentPortTypeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentporttype.Table,
			Columns: equipmentporttype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: equipmentporttype.FieldID,
			},
		},
	}
	if ps := eptu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := eptu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmentporttype.FieldUpdateTime,
		})
	}
	if value := eptu.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipmentporttype.FieldName,
		})
	}
	if nodes := eptu.removedPropertyTypes; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.PropertyTypesTable,
			Columns: []string{equipmentporttype.PropertyTypesColumn},
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
	if nodes := eptu.property_types; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.PropertyTypesTable,
			Columns: []string{equipmentporttype.PropertyTypesColumn},
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
	if nodes := eptu.removedLinkPropertyTypes; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.LinkPropertyTypesTable,
			Columns: []string{equipmentporttype.LinkPropertyTypesColumn},
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
	if nodes := eptu.link_property_types; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.LinkPropertyTypesTable,
			Columns: []string{equipmentporttype.LinkPropertyTypesColumn},
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
	if nodes := eptu.removedPortDefinitions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentporttype.PortDefinitionsTable,
			Columns: []string{equipmentporttype.PortDefinitionsColumn},
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
	if nodes := eptu.port_definitions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentporttype.PortDefinitionsTable,
			Columns: []string{equipmentporttype.PortDefinitionsColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, eptu.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// EquipmentPortTypeUpdateOne is the builder for updating a single EquipmentPortType entity.
type EquipmentPortTypeUpdateOne struct {
	config
	id string

	update_time              *time.Time
	name                     *string
	property_types           map[string]struct{}
	link_property_types      map[string]struct{}
	port_definitions         map[string]struct{}
	removedPropertyTypes     map[string]struct{}
	removedLinkPropertyTypes map[string]struct{}
	removedPortDefinitions   map[string]struct{}
}

// SetName sets the name field.
func (eptuo *EquipmentPortTypeUpdateOne) SetName(s string) *EquipmentPortTypeUpdateOne {
	eptuo.name = &s
	return eptuo
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (eptuo *EquipmentPortTypeUpdateOne) AddPropertyTypeIDs(ids ...string) *EquipmentPortTypeUpdateOne {
	if eptuo.property_types == nil {
		eptuo.property_types = make(map[string]struct{})
	}
	for i := range ids {
		eptuo.property_types[ids[i]] = struct{}{}
	}
	return eptuo
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (eptuo *EquipmentPortTypeUpdateOne) AddPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptuo.AddPropertyTypeIDs(ids...)
}

// AddLinkPropertyTypeIDs adds the link_property_types edge to PropertyType by ids.
func (eptuo *EquipmentPortTypeUpdateOne) AddLinkPropertyTypeIDs(ids ...string) *EquipmentPortTypeUpdateOne {
	if eptuo.link_property_types == nil {
		eptuo.link_property_types = make(map[string]struct{})
	}
	for i := range ids {
		eptuo.link_property_types[ids[i]] = struct{}{}
	}
	return eptuo
}

// AddLinkPropertyTypes adds the link_property_types edges to PropertyType.
func (eptuo *EquipmentPortTypeUpdateOne) AddLinkPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptuo.AddLinkPropertyTypeIDs(ids...)
}

// AddPortDefinitionIDs adds the port_definitions edge to EquipmentPortDefinition by ids.
func (eptuo *EquipmentPortTypeUpdateOne) AddPortDefinitionIDs(ids ...string) *EquipmentPortTypeUpdateOne {
	if eptuo.port_definitions == nil {
		eptuo.port_definitions = make(map[string]struct{})
	}
	for i := range ids {
		eptuo.port_definitions[ids[i]] = struct{}{}
	}
	return eptuo
}

// AddPortDefinitions adds the port_definitions edges to EquipmentPortDefinition.
func (eptuo *EquipmentPortTypeUpdateOne) AddPortDefinitions(e ...*EquipmentPortDefinition) *EquipmentPortTypeUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eptuo.AddPortDefinitionIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (eptuo *EquipmentPortTypeUpdateOne) RemovePropertyTypeIDs(ids ...string) *EquipmentPortTypeUpdateOne {
	if eptuo.removedPropertyTypes == nil {
		eptuo.removedPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		eptuo.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return eptuo
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (eptuo *EquipmentPortTypeUpdateOne) RemovePropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptuo.RemovePropertyTypeIDs(ids...)
}

// RemoveLinkPropertyTypeIDs removes the link_property_types edge to PropertyType by ids.
func (eptuo *EquipmentPortTypeUpdateOne) RemoveLinkPropertyTypeIDs(ids ...string) *EquipmentPortTypeUpdateOne {
	if eptuo.removedLinkPropertyTypes == nil {
		eptuo.removedLinkPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		eptuo.removedLinkPropertyTypes[ids[i]] = struct{}{}
	}
	return eptuo
}

// RemoveLinkPropertyTypes removes link_property_types edges to PropertyType.
func (eptuo *EquipmentPortTypeUpdateOne) RemoveLinkPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptuo.RemoveLinkPropertyTypeIDs(ids...)
}

// RemovePortDefinitionIDs removes the port_definitions edge to EquipmentPortDefinition by ids.
func (eptuo *EquipmentPortTypeUpdateOne) RemovePortDefinitionIDs(ids ...string) *EquipmentPortTypeUpdateOne {
	if eptuo.removedPortDefinitions == nil {
		eptuo.removedPortDefinitions = make(map[string]struct{})
	}
	for i := range ids {
		eptuo.removedPortDefinitions[ids[i]] = struct{}{}
	}
	return eptuo
}

// RemovePortDefinitions removes port_definitions edges to EquipmentPortDefinition.
func (eptuo *EquipmentPortTypeUpdateOne) RemovePortDefinitions(e ...*EquipmentPortDefinition) *EquipmentPortTypeUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eptuo.RemovePortDefinitionIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (eptuo *EquipmentPortTypeUpdateOne) Save(ctx context.Context) (*EquipmentPortType, error) {
	if eptuo.update_time == nil {
		v := equipmentporttype.UpdateDefaultUpdateTime()
		eptuo.update_time = &v
	}
	return eptuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (eptuo *EquipmentPortTypeUpdateOne) SaveX(ctx context.Context) *EquipmentPortType {
	ept, err := eptuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return ept
}

// Exec executes the query on the entity.
func (eptuo *EquipmentPortTypeUpdateOne) Exec(ctx context.Context) error {
	_, err := eptuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eptuo *EquipmentPortTypeUpdateOne) ExecX(ctx context.Context) {
	if err := eptuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (eptuo *EquipmentPortTypeUpdateOne) sqlSave(ctx context.Context) (ept *EquipmentPortType, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentporttype.Table,
			Columns: equipmentporttype.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  eptuo.id,
				Type:   field.TypeString,
				Column: equipmentporttype.FieldID,
			},
		},
	}
	if value := eptuo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmentporttype.FieldUpdateTime,
		})
	}
	if value := eptuo.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipmentporttype.FieldName,
		})
	}
	if nodes := eptuo.removedPropertyTypes; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.PropertyTypesTable,
			Columns: []string{equipmentporttype.PropertyTypesColumn},
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
	if nodes := eptuo.property_types; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.PropertyTypesTable,
			Columns: []string{equipmentporttype.PropertyTypesColumn},
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
	if nodes := eptuo.removedLinkPropertyTypes; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.LinkPropertyTypesTable,
			Columns: []string{equipmentporttype.LinkPropertyTypesColumn},
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
	if nodes := eptuo.link_property_types; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.LinkPropertyTypesTable,
			Columns: []string{equipmentporttype.LinkPropertyTypesColumn},
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
	if nodes := eptuo.removedPortDefinitions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentporttype.PortDefinitionsTable,
			Columns: []string{equipmentporttype.PortDefinitionsColumn},
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
	if nodes := eptuo.port_definitions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentporttype.PortDefinitionsTable,
			Columns: []string{equipmentporttype.PortDefinitionsColumn},
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
	ept = &EquipmentPortType{config: eptuo.config}
	_spec.Assign = ept.assignValues
	_spec.ScanValues = ept.scanValues()
	if err = sqlgraph.UpdateNode(ctx, eptuo.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return ept, nil
}
