// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/ent/equipment"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentcategory"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/pkg/ent/equipmenttype"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/propertytype"
	"github.com/facebookincubator/symphony/pkg/ent/serviceendpointdefinition"
)

// EquipmentTypeUpdate is the builder for updating EquipmentType entities.
type EquipmentTypeUpdate struct {
	config
	hooks      []Hook
	mutation   *EquipmentTypeMutation
	predicates []predicate.EquipmentType
}

// Where adds a new predicate for the builder.
func (etu *EquipmentTypeUpdate) Where(ps ...predicate.EquipmentType) *EquipmentTypeUpdate {
	etu.predicates = append(etu.predicates, ps...)
	return etu
}

// SetName sets the name field.
func (etu *EquipmentTypeUpdate) SetName(s string) *EquipmentTypeUpdate {
	etu.mutation.SetName(s)
	return etu
}

// AddPortDefinitionIDs adds the port_definitions edge to EquipmentPortDefinition by ids.
func (etu *EquipmentTypeUpdate) AddPortDefinitionIDs(ids ...int) *EquipmentTypeUpdate {
	etu.mutation.AddPortDefinitionIDs(ids...)
	return etu
}

// AddPortDefinitions adds the port_definitions edges to EquipmentPortDefinition.
func (etu *EquipmentTypeUpdate) AddPortDefinitions(e ...*EquipmentPortDefinition) *EquipmentTypeUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.AddPortDefinitionIDs(ids...)
}

// AddPositionDefinitionIDs adds the position_definitions edge to EquipmentPositionDefinition by ids.
func (etu *EquipmentTypeUpdate) AddPositionDefinitionIDs(ids ...int) *EquipmentTypeUpdate {
	etu.mutation.AddPositionDefinitionIDs(ids...)
	return etu
}

// AddPositionDefinitions adds the position_definitions edges to EquipmentPositionDefinition.
func (etu *EquipmentTypeUpdate) AddPositionDefinitions(e ...*EquipmentPositionDefinition) *EquipmentTypeUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.AddPositionDefinitionIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (etu *EquipmentTypeUpdate) AddPropertyTypeIDs(ids ...int) *EquipmentTypeUpdate {
	etu.mutation.AddPropertyTypeIDs(ids...)
	return etu
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (etu *EquipmentTypeUpdate) AddPropertyTypes(p ...*PropertyType) *EquipmentTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return etu.AddPropertyTypeIDs(ids...)
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (etu *EquipmentTypeUpdate) AddEquipmentIDs(ids ...int) *EquipmentTypeUpdate {
	etu.mutation.AddEquipmentIDs(ids...)
	return etu
}

// AddEquipment adds the equipment edges to Equipment.
func (etu *EquipmentTypeUpdate) AddEquipment(e ...*Equipment) *EquipmentTypeUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.AddEquipmentIDs(ids...)
}

// SetCategoryID sets the category edge to EquipmentCategory by id.
func (etu *EquipmentTypeUpdate) SetCategoryID(id int) *EquipmentTypeUpdate {
	etu.mutation.SetCategoryID(id)
	return etu
}

// SetNillableCategoryID sets the category edge to EquipmentCategory by id if the given value is not nil.
func (etu *EquipmentTypeUpdate) SetNillableCategoryID(id *int) *EquipmentTypeUpdate {
	if id != nil {
		etu = etu.SetCategoryID(*id)
	}
	return etu
}

// SetCategory sets the category edge to EquipmentCategory.
func (etu *EquipmentTypeUpdate) SetCategory(e *EquipmentCategory) *EquipmentTypeUpdate {
	return etu.SetCategoryID(e.ID)
}

// AddServiceEndpointDefinitionIDs adds the service_endpoint_definitions edge to ServiceEndpointDefinition by ids.
func (etu *EquipmentTypeUpdate) AddServiceEndpointDefinitionIDs(ids ...int) *EquipmentTypeUpdate {
	etu.mutation.AddServiceEndpointDefinitionIDs(ids...)
	return etu
}

// AddServiceEndpointDefinitions adds the service_endpoint_definitions edges to ServiceEndpointDefinition.
func (etu *EquipmentTypeUpdate) AddServiceEndpointDefinitions(s ...*ServiceEndpointDefinition) *EquipmentTypeUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return etu.AddServiceEndpointDefinitionIDs(ids...)
}

// RemovePortDefinitionIDs removes the port_definitions edge to EquipmentPortDefinition by ids.
func (etu *EquipmentTypeUpdate) RemovePortDefinitionIDs(ids ...int) *EquipmentTypeUpdate {
	etu.mutation.RemovePortDefinitionIDs(ids...)
	return etu
}

// RemovePortDefinitions removes port_definitions edges to EquipmentPortDefinition.
func (etu *EquipmentTypeUpdate) RemovePortDefinitions(e ...*EquipmentPortDefinition) *EquipmentTypeUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.RemovePortDefinitionIDs(ids...)
}

// RemovePositionDefinitionIDs removes the position_definitions edge to EquipmentPositionDefinition by ids.
func (etu *EquipmentTypeUpdate) RemovePositionDefinitionIDs(ids ...int) *EquipmentTypeUpdate {
	etu.mutation.RemovePositionDefinitionIDs(ids...)
	return etu
}

// RemovePositionDefinitions removes position_definitions edges to EquipmentPositionDefinition.
func (etu *EquipmentTypeUpdate) RemovePositionDefinitions(e ...*EquipmentPositionDefinition) *EquipmentTypeUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.RemovePositionDefinitionIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (etu *EquipmentTypeUpdate) RemovePropertyTypeIDs(ids ...int) *EquipmentTypeUpdate {
	etu.mutation.RemovePropertyTypeIDs(ids...)
	return etu
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (etu *EquipmentTypeUpdate) RemovePropertyTypes(p ...*PropertyType) *EquipmentTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return etu.RemovePropertyTypeIDs(ids...)
}

// RemoveEquipmentIDs removes the equipment edge to Equipment by ids.
func (etu *EquipmentTypeUpdate) RemoveEquipmentIDs(ids ...int) *EquipmentTypeUpdate {
	etu.mutation.RemoveEquipmentIDs(ids...)
	return etu
}

// RemoveEquipment removes equipment edges to Equipment.
func (etu *EquipmentTypeUpdate) RemoveEquipment(e ...*Equipment) *EquipmentTypeUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.RemoveEquipmentIDs(ids...)
}

// ClearCategory clears the category edge to EquipmentCategory.
func (etu *EquipmentTypeUpdate) ClearCategory() *EquipmentTypeUpdate {
	etu.mutation.ClearCategory()
	return etu
}

// RemoveServiceEndpointDefinitionIDs removes the service_endpoint_definitions edge to ServiceEndpointDefinition by ids.
func (etu *EquipmentTypeUpdate) RemoveServiceEndpointDefinitionIDs(ids ...int) *EquipmentTypeUpdate {
	etu.mutation.RemoveServiceEndpointDefinitionIDs(ids...)
	return etu
}

// RemoveServiceEndpointDefinitions removes service_endpoint_definitions edges to ServiceEndpointDefinition.
func (etu *EquipmentTypeUpdate) RemoveServiceEndpointDefinitions(s ...*ServiceEndpointDefinition) *EquipmentTypeUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return etu.RemoveServiceEndpointDefinitionIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (etu *EquipmentTypeUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := etu.mutation.UpdateTime(); !ok {
		v := equipmenttype.UpdateDefaultUpdateTime()
		etu.mutation.SetUpdateTime(v)
	}

	var (
		err      error
		affected int
	)
	if len(etu.hooks) == 0 {
		affected, err = etu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			etu.mutation = mutation
			affected, err = etu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(etu.hooks) - 1; i >= 0; i-- {
			mut = etu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, etu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
				Type:   field.TypeInt,
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
	if value, ok := etu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmenttype.FieldUpdateTime,
		})
	}
	if value, ok := etu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmenttype.FieldName,
		})
	}
	if nodes := etu.mutation.RemovedPortDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.PortDefinitionsTable,
			Columns: []string{equipmenttype.PortDefinitionsColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etu.mutation.PortDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.PortDefinitionsTable,
			Columns: []string{equipmenttype.PortDefinitionsColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := etu.mutation.RemovedPositionDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.PositionDefinitionsTable,
			Columns: []string{equipmenttype.PositionDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentpositiondefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etu.mutation.PositionDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.PositionDefinitionsTable,
			Columns: []string{equipmenttype.PositionDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentpositiondefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := etu.mutation.RemovedPropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.PropertyTypesTable,
			Columns: []string{equipmenttype.PropertyTypesColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etu.mutation.PropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.PropertyTypesTable,
			Columns: []string{equipmenttype.PropertyTypesColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := etu.mutation.RemovedEquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmenttype.EquipmentTable,
			Columns: []string{equipmenttype.EquipmentColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etu.mutation.EquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmenttype.EquipmentTable,
			Columns: []string{equipmenttype.EquipmentColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if etu.mutation.CategoryCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmenttype.CategoryTable,
			Columns: []string{equipmenttype.CategoryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentcategory.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etu.mutation.CategoryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmenttype.CategoryTable,
			Columns: []string{equipmenttype.CategoryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentcategory.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := etu.mutation.RemovedServiceEndpointDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.ServiceEndpointDefinitionsTable,
			Columns: []string{equipmenttype.ServiceEndpointDefinitionsColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etu.mutation.ServiceEndpointDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.ServiceEndpointDefinitionsTable,
			Columns: []string{equipmenttype.ServiceEndpointDefinitionsColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, etu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{equipmenttype.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// EquipmentTypeUpdateOne is the builder for updating a single EquipmentType entity.
type EquipmentTypeUpdateOne struct {
	config
	hooks    []Hook
	mutation *EquipmentTypeMutation
}

// SetName sets the name field.
func (etuo *EquipmentTypeUpdateOne) SetName(s string) *EquipmentTypeUpdateOne {
	etuo.mutation.SetName(s)
	return etuo
}

// AddPortDefinitionIDs adds the port_definitions edge to EquipmentPortDefinition by ids.
func (etuo *EquipmentTypeUpdateOne) AddPortDefinitionIDs(ids ...int) *EquipmentTypeUpdateOne {
	etuo.mutation.AddPortDefinitionIDs(ids...)
	return etuo
}

// AddPortDefinitions adds the port_definitions edges to EquipmentPortDefinition.
func (etuo *EquipmentTypeUpdateOne) AddPortDefinitions(e ...*EquipmentPortDefinition) *EquipmentTypeUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.AddPortDefinitionIDs(ids...)
}

// AddPositionDefinitionIDs adds the position_definitions edge to EquipmentPositionDefinition by ids.
func (etuo *EquipmentTypeUpdateOne) AddPositionDefinitionIDs(ids ...int) *EquipmentTypeUpdateOne {
	etuo.mutation.AddPositionDefinitionIDs(ids...)
	return etuo
}

// AddPositionDefinitions adds the position_definitions edges to EquipmentPositionDefinition.
func (etuo *EquipmentTypeUpdateOne) AddPositionDefinitions(e ...*EquipmentPositionDefinition) *EquipmentTypeUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.AddPositionDefinitionIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (etuo *EquipmentTypeUpdateOne) AddPropertyTypeIDs(ids ...int) *EquipmentTypeUpdateOne {
	etuo.mutation.AddPropertyTypeIDs(ids...)
	return etuo
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (etuo *EquipmentTypeUpdateOne) AddPropertyTypes(p ...*PropertyType) *EquipmentTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return etuo.AddPropertyTypeIDs(ids...)
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (etuo *EquipmentTypeUpdateOne) AddEquipmentIDs(ids ...int) *EquipmentTypeUpdateOne {
	etuo.mutation.AddEquipmentIDs(ids...)
	return etuo
}

// AddEquipment adds the equipment edges to Equipment.
func (etuo *EquipmentTypeUpdateOne) AddEquipment(e ...*Equipment) *EquipmentTypeUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.AddEquipmentIDs(ids...)
}

// SetCategoryID sets the category edge to EquipmentCategory by id.
func (etuo *EquipmentTypeUpdateOne) SetCategoryID(id int) *EquipmentTypeUpdateOne {
	etuo.mutation.SetCategoryID(id)
	return etuo
}

// SetNillableCategoryID sets the category edge to EquipmentCategory by id if the given value is not nil.
func (etuo *EquipmentTypeUpdateOne) SetNillableCategoryID(id *int) *EquipmentTypeUpdateOne {
	if id != nil {
		etuo = etuo.SetCategoryID(*id)
	}
	return etuo
}

// SetCategory sets the category edge to EquipmentCategory.
func (etuo *EquipmentTypeUpdateOne) SetCategory(e *EquipmentCategory) *EquipmentTypeUpdateOne {
	return etuo.SetCategoryID(e.ID)
}

// AddServiceEndpointDefinitionIDs adds the service_endpoint_definitions edge to ServiceEndpointDefinition by ids.
func (etuo *EquipmentTypeUpdateOne) AddServiceEndpointDefinitionIDs(ids ...int) *EquipmentTypeUpdateOne {
	etuo.mutation.AddServiceEndpointDefinitionIDs(ids...)
	return etuo
}

// AddServiceEndpointDefinitions adds the service_endpoint_definitions edges to ServiceEndpointDefinition.
func (etuo *EquipmentTypeUpdateOne) AddServiceEndpointDefinitions(s ...*ServiceEndpointDefinition) *EquipmentTypeUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return etuo.AddServiceEndpointDefinitionIDs(ids...)
}

// RemovePortDefinitionIDs removes the port_definitions edge to EquipmentPortDefinition by ids.
func (etuo *EquipmentTypeUpdateOne) RemovePortDefinitionIDs(ids ...int) *EquipmentTypeUpdateOne {
	etuo.mutation.RemovePortDefinitionIDs(ids...)
	return etuo
}

// RemovePortDefinitions removes port_definitions edges to EquipmentPortDefinition.
func (etuo *EquipmentTypeUpdateOne) RemovePortDefinitions(e ...*EquipmentPortDefinition) *EquipmentTypeUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.RemovePortDefinitionIDs(ids...)
}

// RemovePositionDefinitionIDs removes the position_definitions edge to EquipmentPositionDefinition by ids.
func (etuo *EquipmentTypeUpdateOne) RemovePositionDefinitionIDs(ids ...int) *EquipmentTypeUpdateOne {
	etuo.mutation.RemovePositionDefinitionIDs(ids...)
	return etuo
}

// RemovePositionDefinitions removes position_definitions edges to EquipmentPositionDefinition.
func (etuo *EquipmentTypeUpdateOne) RemovePositionDefinitions(e ...*EquipmentPositionDefinition) *EquipmentTypeUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.RemovePositionDefinitionIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (etuo *EquipmentTypeUpdateOne) RemovePropertyTypeIDs(ids ...int) *EquipmentTypeUpdateOne {
	etuo.mutation.RemovePropertyTypeIDs(ids...)
	return etuo
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (etuo *EquipmentTypeUpdateOne) RemovePropertyTypes(p ...*PropertyType) *EquipmentTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return etuo.RemovePropertyTypeIDs(ids...)
}

// RemoveEquipmentIDs removes the equipment edge to Equipment by ids.
func (etuo *EquipmentTypeUpdateOne) RemoveEquipmentIDs(ids ...int) *EquipmentTypeUpdateOne {
	etuo.mutation.RemoveEquipmentIDs(ids...)
	return etuo
}

// RemoveEquipment removes equipment edges to Equipment.
func (etuo *EquipmentTypeUpdateOne) RemoveEquipment(e ...*Equipment) *EquipmentTypeUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.RemoveEquipmentIDs(ids...)
}

// ClearCategory clears the category edge to EquipmentCategory.
func (etuo *EquipmentTypeUpdateOne) ClearCategory() *EquipmentTypeUpdateOne {
	etuo.mutation.ClearCategory()
	return etuo
}

// RemoveServiceEndpointDefinitionIDs removes the service_endpoint_definitions edge to ServiceEndpointDefinition by ids.
func (etuo *EquipmentTypeUpdateOne) RemoveServiceEndpointDefinitionIDs(ids ...int) *EquipmentTypeUpdateOne {
	etuo.mutation.RemoveServiceEndpointDefinitionIDs(ids...)
	return etuo
}

// RemoveServiceEndpointDefinitions removes service_endpoint_definitions edges to ServiceEndpointDefinition.
func (etuo *EquipmentTypeUpdateOne) RemoveServiceEndpointDefinitions(s ...*ServiceEndpointDefinition) *EquipmentTypeUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return etuo.RemoveServiceEndpointDefinitionIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (etuo *EquipmentTypeUpdateOne) Save(ctx context.Context) (*EquipmentType, error) {
	if _, ok := etuo.mutation.UpdateTime(); !ok {
		v := equipmenttype.UpdateDefaultUpdateTime()
		etuo.mutation.SetUpdateTime(v)
	}

	var (
		err  error
		node *EquipmentType
	)
	if len(etuo.hooks) == 0 {
		node, err = etuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			etuo.mutation = mutation
			node, err = etuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(etuo.hooks) - 1; i >= 0; i-- {
			mut = etuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, etuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
				Type:   field.TypeInt,
				Column: equipmenttype.FieldID,
			},
		},
	}
	id, ok := etuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing EquipmentType.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := etuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmenttype.FieldUpdateTime,
		})
	}
	if value, ok := etuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmenttype.FieldName,
		})
	}
	if nodes := etuo.mutation.RemovedPortDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.PortDefinitionsTable,
			Columns: []string{equipmenttype.PortDefinitionsColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etuo.mutation.PortDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.PortDefinitionsTable,
			Columns: []string{equipmenttype.PortDefinitionsColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := etuo.mutation.RemovedPositionDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.PositionDefinitionsTable,
			Columns: []string{equipmenttype.PositionDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentpositiondefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etuo.mutation.PositionDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.PositionDefinitionsTable,
			Columns: []string{equipmenttype.PositionDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentpositiondefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := etuo.mutation.RemovedPropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.PropertyTypesTable,
			Columns: []string{equipmenttype.PropertyTypesColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etuo.mutation.PropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.PropertyTypesTable,
			Columns: []string{equipmenttype.PropertyTypesColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := etuo.mutation.RemovedEquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmenttype.EquipmentTable,
			Columns: []string{equipmenttype.EquipmentColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etuo.mutation.EquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmenttype.EquipmentTable,
			Columns: []string{equipmenttype.EquipmentColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if etuo.mutation.CategoryCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmenttype.CategoryTable,
			Columns: []string{equipmenttype.CategoryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentcategory.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etuo.mutation.CategoryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmenttype.CategoryTable,
			Columns: []string{equipmenttype.CategoryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentcategory.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := etuo.mutation.RemovedServiceEndpointDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.ServiceEndpointDefinitionsTable,
			Columns: []string{equipmenttype.ServiceEndpointDefinitionsColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := etuo.mutation.ServiceEndpointDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmenttype.ServiceEndpointDefinitionsTable,
			Columns: []string{equipmenttype.ServiceEndpointDefinitionsColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	et = &EquipmentType{config: etuo.config}
	_spec.Assign = et.assignValues
	_spec.ScanValues = et.scanValues()
	if err = sqlgraph.UpdateNode(ctx, etuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{equipmenttype.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return et, nil
}
