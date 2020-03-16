// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/checklistcategory"
	"github.com/facebookincubator/symphony/graph/ent/checklistitemdefinition"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// WorkOrderTypeUpdate is the builder for updating WorkOrderType entities.
type WorkOrderTypeUpdate struct {
	config

	update_time                 *time.Time
	name                        *string
	description                 *string
	cleardescription            bool
	work_orders                 map[int]struct{}
	property_types              map[int]struct{}
	definitions                 map[int]struct{}
	check_list_categories       map[int]struct{}
	check_list_definitions      map[int]struct{}
	removedWorkOrders           map[int]struct{}
	removedPropertyTypes        map[int]struct{}
	removedDefinitions          map[int]struct{}
	removedCheckListCategories  map[int]struct{}
	removedCheckListDefinitions map[int]struct{}
	predicates                  []predicate.WorkOrderType
}

// Where adds a new predicate for the builder.
func (wotu *WorkOrderTypeUpdate) Where(ps ...predicate.WorkOrderType) *WorkOrderTypeUpdate {
	wotu.predicates = append(wotu.predicates, ps...)
	return wotu
}

// SetName sets the name field.
func (wotu *WorkOrderTypeUpdate) SetName(s string) *WorkOrderTypeUpdate {
	wotu.name = &s
	return wotu
}

// SetDescription sets the description field.
func (wotu *WorkOrderTypeUpdate) SetDescription(s string) *WorkOrderTypeUpdate {
	wotu.description = &s
	return wotu
}

// SetNillableDescription sets the description field if the given value is not nil.
func (wotu *WorkOrderTypeUpdate) SetNillableDescription(s *string) *WorkOrderTypeUpdate {
	if s != nil {
		wotu.SetDescription(*s)
	}
	return wotu
}

// ClearDescription clears the value of description.
func (wotu *WorkOrderTypeUpdate) ClearDescription() *WorkOrderTypeUpdate {
	wotu.description = nil
	wotu.cleardescription = true
	return wotu
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (wotu *WorkOrderTypeUpdate) AddWorkOrderIDs(ids ...int) *WorkOrderTypeUpdate {
	if wotu.work_orders == nil {
		wotu.work_orders = make(map[int]struct{})
	}
	for i := range ids {
		wotu.work_orders[ids[i]] = struct{}{}
	}
	return wotu
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (wotu *WorkOrderTypeUpdate) AddWorkOrders(w ...*WorkOrder) *WorkOrderTypeUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotu.AddWorkOrderIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (wotu *WorkOrderTypeUpdate) AddPropertyTypeIDs(ids ...int) *WorkOrderTypeUpdate {
	if wotu.property_types == nil {
		wotu.property_types = make(map[int]struct{})
	}
	for i := range ids {
		wotu.property_types[ids[i]] = struct{}{}
	}
	return wotu
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (wotu *WorkOrderTypeUpdate) AddPropertyTypes(p ...*PropertyType) *WorkOrderTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wotu.AddPropertyTypeIDs(ids...)
}

// AddDefinitionIDs adds the definitions edge to WorkOrderDefinition by ids.
func (wotu *WorkOrderTypeUpdate) AddDefinitionIDs(ids ...int) *WorkOrderTypeUpdate {
	if wotu.definitions == nil {
		wotu.definitions = make(map[int]struct{})
	}
	for i := range ids {
		wotu.definitions[ids[i]] = struct{}{}
	}
	return wotu
}

// AddDefinitions adds the definitions edges to WorkOrderDefinition.
func (wotu *WorkOrderTypeUpdate) AddDefinitions(w ...*WorkOrderDefinition) *WorkOrderTypeUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotu.AddDefinitionIDs(ids...)
}

// AddCheckListCategoryIDs adds the check_list_categories edge to CheckListCategory by ids.
func (wotu *WorkOrderTypeUpdate) AddCheckListCategoryIDs(ids ...int) *WorkOrderTypeUpdate {
	if wotu.check_list_categories == nil {
		wotu.check_list_categories = make(map[int]struct{})
	}
	for i := range ids {
		wotu.check_list_categories[ids[i]] = struct{}{}
	}
	return wotu
}

// AddCheckListCategories adds the check_list_categories edges to CheckListCategory.
func (wotu *WorkOrderTypeUpdate) AddCheckListCategories(c ...*CheckListCategory) *WorkOrderTypeUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotu.AddCheckListCategoryIDs(ids...)
}

// AddCheckListDefinitionIDs adds the check_list_definitions edge to CheckListItemDefinition by ids.
func (wotu *WorkOrderTypeUpdate) AddCheckListDefinitionIDs(ids ...int) *WorkOrderTypeUpdate {
	if wotu.check_list_definitions == nil {
		wotu.check_list_definitions = make(map[int]struct{})
	}
	for i := range ids {
		wotu.check_list_definitions[ids[i]] = struct{}{}
	}
	return wotu
}

// AddCheckListDefinitions adds the check_list_definitions edges to CheckListItemDefinition.
func (wotu *WorkOrderTypeUpdate) AddCheckListDefinitions(c ...*CheckListItemDefinition) *WorkOrderTypeUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotu.AddCheckListDefinitionIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (wotu *WorkOrderTypeUpdate) RemoveWorkOrderIDs(ids ...int) *WorkOrderTypeUpdate {
	if wotu.removedWorkOrders == nil {
		wotu.removedWorkOrders = make(map[int]struct{})
	}
	for i := range ids {
		wotu.removedWorkOrders[ids[i]] = struct{}{}
	}
	return wotu
}

// RemoveWorkOrders removes work_orders edges to WorkOrder.
func (wotu *WorkOrderTypeUpdate) RemoveWorkOrders(w ...*WorkOrder) *WorkOrderTypeUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotu.RemoveWorkOrderIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (wotu *WorkOrderTypeUpdate) RemovePropertyTypeIDs(ids ...int) *WorkOrderTypeUpdate {
	if wotu.removedPropertyTypes == nil {
		wotu.removedPropertyTypes = make(map[int]struct{})
	}
	for i := range ids {
		wotu.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return wotu
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (wotu *WorkOrderTypeUpdate) RemovePropertyTypes(p ...*PropertyType) *WorkOrderTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wotu.RemovePropertyTypeIDs(ids...)
}

// RemoveDefinitionIDs removes the definitions edge to WorkOrderDefinition by ids.
func (wotu *WorkOrderTypeUpdate) RemoveDefinitionIDs(ids ...int) *WorkOrderTypeUpdate {
	if wotu.removedDefinitions == nil {
		wotu.removedDefinitions = make(map[int]struct{})
	}
	for i := range ids {
		wotu.removedDefinitions[ids[i]] = struct{}{}
	}
	return wotu
}

// RemoveDefinitions removes definitions edges to WorkOrderDefinition.
func (wotu *WorkOrderTypeUpdate) RemoveDefinitions(w ...*WorkOrderDefinition) *WorkOrderTypeUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotu.RemoveDefinitionIDs(ids...)
}

// RemoveCheckListCategoryIDs removes the check_list_categories edge to CheckListCategory by ids.
func (wotu *WorkOrderTypeUpdate) RemoveCheckListCategoryIDs(ids ...int) *WorkOrderTypeUpdate {
	if wotu.removedCheckListCategories == nil {
		wotu.removedCheckListCategories = make(map[int]struct{})
	}
	for i := range ids {
		wotu.removedCheckListCategories[ids[i]] = struct{}{}
	}
	return wotu
}

// RemoveCheckListCategories removes check_list_categories edges to CheckListCategory.
func (wotu *WorkOrderTypeUpdate) RemoveCheckListCategories(c ...*CheckListCategory) *WorkOrderTypeUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotu.RemoveCheckListCategoryIDs(ids...)
}

// RemoveCheckListDefinitionIDs removes the check_list_definitions edge to CheckListItemDefinition by ids.
func (wotu *WorkOrderTypeUpdate) RemoveCheckListDefinitionIDs(ids ...int) *WorkOrderTypeUpdate {
	if wotu.removedCheckListDefinitions == nil {
		wotu.removedCheckListDefinitions = make(map[int]struct{})
	}
	for i := range ids {
		wotu.removedCheckListDefinitions[ids[i]] = struct{}{}
	}
	return wotu
}

// RemoveCheckListDefinitions removes check_list_definitions edges to CheckListItemDefinition.
func (wotu *WorkOrderTypeUpdate) RemoveCheckListDefinitions(c ...*CheckListItemDefinition) *WorkOrderTypeUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotu.RemoveCheckListDefinitionIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (wotu *WorkOrderTypeUpdate) Save(ctx context.Context) (int, error) {
	if wotu.update_time == nil {
		v := workordertype.UpdateDefaultUpdateTime()
		wotu.update_time = &v
	}
	return wotu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (wotu *WorkOrderTypeUpdate) SaveX(ctx context.Context) int {
	affected, err := wotu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (wotu *WorkOrderTypeUpdate) Exec(ctx context.Context) error {
	_, err := wotu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wotu *WorkOrderTypeUpdate) ExecX(ctx context.Context) {
	if err := wotu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (wotu *WorkOrderTypeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   workordertype.Table,
			Columns: workordertype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: workordertype.FieldID,
			},
		},
	}
	if ps := wotu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := wotu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: workordertype.FieldUpdateTime,
		})
	}
	if value := wotu.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: workordertype.FieldName,
		})
	}
	if value := wotu.description; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: workordertype.FieldDescription,
		})
	}
	if wotu.cleardescription {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: workordertype.FieldDescription,
		})
	}
	if nodes := wotu.removedWorkOrders; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workordertype.WorkOrdersTable,
			Columns: []string{workordertype.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotu.work_orders; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workordertype.WorkOrdersTable,
			Columns: []string{workordertype.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wotu.removedPropertyTypes; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertype.PropertyTypesTable,
			Columns: []string{workordertype.PropertyTypesColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotu.property_types; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertype.PropertyTypesTable,
			Columns: []string{workordertype.PropertyTypesColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wotu.removedDefinitions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workordertype.DefinitionsTable,
			Columns: []string{workordertype.DefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorderdefinition.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotu.definitions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workordertype.DefinitionsTable,
			Columns: []string{workordertype.DefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorderdefinition.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wotu.removedCheckListCategories; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertype.CheckListCategoriesTable,
			Columns: []string{workordertype.CheckListCategoriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistcategory.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotu.check_list_categories; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertype.CheckListCategoriesTable,
			Columns: []string{workordertype.CheckListCategoriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistcategory.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wotu.removedCheckListDefinitions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertype.CheckListDefinitionsTable,
			Columns: []string{workordertype.CheckListDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistitemdefinition.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotu.check_list_definitions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertype.CheckListDefinitionsTable,
			Columns: []string{workordertype.CheckListDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistitemdefinition.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, wotu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{workordertype.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// WorkOrderTypeUpdateOne is the builder for updating a single WorkOrderType entity.
type WorkOrderTypeUpdateOne struct {
	config
	id int

	update_time                 *time.Time
	name                        *string
	description                 *string
	cleardescription            bool
	work_orders                 map[int]struct{}
	property_types              map[int]struct{}
	definitions                 map[int]struct{}
	check_list_categories       map[int]struct{}
	check_list_definitions      map[int]struct{}
	removedWorkOrders           map[int]struct{}
	removedPropertyTypes        map[int]struct{}
	removedDefinitions          map[int]struct{}
	removedCheckListCategories  map[int]struct{}
	removedCheckListDefinitions map[int]struct{}
}

// SetName sets the name field.
func (wotuo *WorkOrderTypeUpdateOne) SetName(s string) *WorkOrderTypeUpdateOne {
	wotuo.name = &s
	return wotuo
}

// SetDescription sets the description field.
func (wotuo *WorkOrderTypeUpdateOne) SetDescription(s string) *WorkOrderTypeUpdateOne {
	wotuo.description = &s
	return wotuo
}

// SetNillableDescription sets the description field if the given value is not nil.
func (wotuo *WorkOrderTypeUpdateOne) SetNillableDescription(s *string) *WorkOrderTypeUpdateOne {
	if s != nil {
		wotuo.SetDescription(*s)
	}
	return wotuo
}

// ClearDescription clears the value of description.
func (wotuo *WorkOrderTypeUpdateOne) ClearDescription() *WorkOrderTypeUpdateOne {
	wotuo.description = nil
	wotuo.cleardescription = true
	return wotuo
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (wotuo *WorkOrderTypeUpdateOne) AddWorkOrderIDs(ids ...int) *WorkOrderTypeUpdateOne {
	if wotuo.work_orders == nil {
		wotuo.work_orders = make(map[int]struct{})
	}
	for i := range ids {
		wotuo.work_orders[ids[i]] = struct{}{}
	}
	return wotuo
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (wotuo *WorkOrderTypeUpdateOne) AddWorkOrders(w ...*WorkOrder) *WorkOrderTypeUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotuo.AddWorkOrderIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (wotuo *WorkOrderTypeUpdateOne) AddPropertyTypeIDs(ids ...int) *WorkOrderTypeUpdateOne {
	if wotuo.property_types == nil {
		wotuo.property_types = make(map[int]struct{})
	}
	for i := range ids {
		wotuo.property_types[ids[i]] = struct{}{}
	}
	return wotuo
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (wotuo *WorkOrderTypeUpdateOne) AddPropertyTypes(p ...*PropertyType) *WorkOrderTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wotuo.AddPropertyTypeIDs(ids...)
}

// AddDefinitionIDs adds the definitions edge to WorkOrderDefinition by ids.
func (wotuo *WorkOrderTypeUpdateOne) AddDefinitionIDs(ids ...int) *WorkOrderTypeUpdateOne {
	if wotuo.definitions == nil {
		wotuo.definitions = make(map[int]struct{})
	}
	for i := range ids {
		wotuo.definitions[ids[i]] = struct{}{}
	}
	return wotuo
}

// AddDefinitions adds the definitions edges to WorkOrderDefinition.
func (wotuo *WorkOrderTypeUpdateOne) AddDefinitions(w ...*WorkOrderDefinition) *WorkOrderTypeUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotuo.AddDefinitionIDs(ids...)
}

// AddCheckListCategoryIDs adds the check_list_categories edge to CheckListCategory by ids.
func (wotuo *WorkOrderTypeUpdateOne) AddCheckListCategoryIDs(ids ...int) *WorkOrderTypeUpdateOne {
	if wotuo.check_list_categories == nil {
		wotuo.check_list_categories = make(map[int]struct{})
	}
	for i := range ids {
		wotuo.check_list_categories[ids[i]] = struct{}{}
	}
	return wotuo
}

// AddCheckListCategories adds the check_list_categories edges to CheckListCategory.
func (wotuo *WorkOrderTypeUpdateOne) AddCheckListCategories(c ...*CheckListCategory) *WorkOrderTypeUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotuo.AddCheckListCategoryIDs(ids...)
}

// AddCheckListDefinitionIDs adds the check_list_definitions edge to CheckListItemDefinition by ids.
func (wotuo *WorkOrderTypeUpdateOne) AddCheckListDefinitionIDs(ids ...int) *WorkOrderTypeUpdateOne {
	if wotuo.check_list_definitions == nil {
		wotuo.check_list_definitions = make(map[int]struct{})
	}
	for i := range ids {
		wotuo.check_list_definitions[ids[i]] = struct{}{}
	}
	return wotuo
}

// AddCheckListDefinitions adds the check_list_definitions edges to CheckListItemDefinition.
func (wotuo *WorkOrderTypeUpdateOne) AddCheckListDefinitions(c ...*CheckListItemDefinition) *WorkOrderTypeUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotuo.AddCheckListDefinitionIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (wotuo *WorkOrderTypeUpdateOne) RemoveWorkOrderIDs(ids ...int) *WorkOrderTypeUpdateOne {
	if wotuo.removedWorkOrders == nil {
		wotuo.removedWorkOrders = make(map[int]struct{})
	}
	for i := range ids {
		wotuo.removedWorkOrders[ids[i]] = struct{}{}
	}
	return wotuo
}

// RemoveWorkOrders removes work_orders edges to WorkOrder.
func (wotuo *WorkOrderTypeUpdateOne) RemoveWorkOrders(w ...*WorkOrder) *WorkOrderTypeUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotuo.RemoveWorkOrderIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (wotuo *WorkOrderTypeUpdateOne) RemovePropertyTypeIDs(ids ...int) *WorkOrderTypeUpdateOne {
	if wotuo.removedPropertyTypes == nil {
		wotuo.removedPropertyTypes = make(map[int]struct{})
	}
	for i := range ids {
		wotuo.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return wotuo
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (wotuo *WorkOrderTypeUpdateOne) RemovePropertyTypes(p ...*PropertyType) *WorkOrderTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wotuo.RemovePropertyTypeIDs(ids...)
}

// RemoveDefinitionIDs removes the definitions edge to WorkOrderDefinition by ids.
func (wotuo *WorkOrderTypeUpdateOne) RemoveDefinitionIDs(ids ...int) *WorkOrderTypeUpdateOne {
	if wotuo.removedDefinitions == nil {
		wotuo.removedDefinitions = make(map[int]struct{})
	}
	for i := range ids {
		wotuo.removedDefinitions[ids[i]] = struct{}{}
	}
	return wotuo
}

// RemoveDefinitions removes definitions edges to WorkOrderDefinition.
func (wotuo *WorkOrderTypeUpdateOne) RemoveDefinitions(w ...*WorkOrderDefinition) *WorkOrderTypeUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotuo.RemoveDefinitionIDs(ids...)
}

// RemoveCheckListCategoryIDs removes the check_list_categories edge to CheckListCategory by ids.
func (wotuo *WorkOrderTypeUpdateOne) RemoveCheckListCategoryIDs(ids ...int) *WorkOrderTypeUpdateOne {
	if wotuo.removedCheckListCategories == nil {
		wotuo.removedCheckListCategories = make(map[int]struct{})
	}
	for i := range ids {
		wotuo.removedCheckListCategories[ids[i]] = struct{}{}
	}
	return wotuo
}

// RemoveCheckListCategories removes check_list_categories edges to CheckListCategory.
func (wotuo *WorkOrderTypeUpdateOne) RemoveCheckListCategories(c ...*CheckListCategory) *WorkOrderTypeUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotuo.RemoveCheckListCategoryIDs(ids...)
}

// RemoveCheckListDefinitionIDs removes the check_list_definitions edge to CheckListItemDefinition by ids.
func (wotuo *WorkOrderTypeUpdateOne) RemoveCheckListDefinitionIDs(ids ...int) *WorkOrderTypeUpdateOne {
	if wotuo.removedCheckListDefinitions == nil {
		wotuo.removedCheckListDefinitions = make(map[int]struct{})
	}
	for i := range ids {
		wotuo.removedCheckListDefinitions[ids[i]] = struct{}{}
	}
	return wotuo
}

// RemoveCheckListDefinitions removes check_list_definitions edges to CheckListItemDefinition.
func (wotuo *WorkOrderTypeUpdateOne) RemoveCheckListDefinitions(c ...*CheckListItemDefinition) *WorkOrderTypeUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotuo.RemoveCheckListDefinitionIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (wotuo *WorkOrderTypeUpdateOne) Save(ctx context.Context) (*WorkOrderType, error) {
	if wotuo.update_time == nil {
		v := workordertype.UpdateDefaultUpdateTime()
		wotuo.update_time = &v
	}
	return wotuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (wotuo *WorkOrderTypeUpdateOne) SaveX(ctx context.Context) *WorkOrderType {
	wot, err := wotuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return wot
}

// Exec executes the query on the entity.
func (wotuo *WorkOrderTypeUpdateOne) Exec(ctx context.Context) error {
	_, err := wotuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wotuo *WorkOrderTypeUpdateOne) ExecX(ctx context.Context) {
	if err := wotuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (wotuo *WorkOrderTypeUpdateOne) sqlSave(ctx context.Context) (wot *WorkOrderType, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   workordertype.Table,
			Columns: workordertype.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  wotuo.id,
				Type:   field.TypeInt,
				Column: workordertype.FieldID,
			},
		},
	}
	if value := wotuo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: workordertype.FieldUpdateTime,
		})
	}
	if value := wotuo.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: workordertype.FieldName,
		})
	}
	if value := wotuo.description; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: workordertype.FieldDescription,
		})
	}
	if wotuo.cleardescription {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: workordertype.FieldDescription,
		})
	}
	if nodes := wotuo.removedWorkOrders; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workordertype.WorkOrdersTable,
			Columns: []string{workordertype.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotuo.work_orders; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workordertype.WorkOrdersTable,
			Columns: []string{workordertype.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wotuo.removedPropertyTypes; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertype.PropertyTypesTable,
			Columns: []string{workordertype.PropertyTypesColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotuo.property_types; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertype.PropertyTypesTable,
			Columns: []string{workordertype.PropertyTypesColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wotuo.removedDefinitions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workordertype.DefinitionsTable,
			Columns: []string{workordertype.DefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorderdefinition.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotuo.definitions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workordertype.DefinitionsTable,
			Columns: []string{workordertype.DefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorderdefinition.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wotuo.removedCheckListCategories; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertype.CheckListCategoriesTable,
			Columns: []string{workordertype.CheckListCategoriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistcategory.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotuo.check_list_categories; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertype.CheckListCategoriesTable,
			Columns: []string{workordertype.CheckListCategoriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistcategory.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wotuo.removedCheckListDefinitions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertype.CheckListDefinitionsTable,
			Columns: []string{workordertype.CheckListDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistitemdefinition.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotuo.check_list_definitions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertype.CheckListDefinitionsTable,
			Columns: []string{workordertype.CheckListDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistitemdefinition.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	wot = &WorkOrderType{config: wotuo.config}
	_spec.Assign = wot.assignValues
	_spec.ScanValues = wot.scanValues()
	if err = sqlgraph.UpdateNode(ctx, wotuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{workordertype.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return wot, nil
}
