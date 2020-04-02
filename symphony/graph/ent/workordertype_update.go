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
	hooks      []Hook
	mutation   *WorkOrderTypeMutation
	predicates []predicate.WorkOrderType
}

// Where adds a new predicate for the builder.
func (wotu *WorkOrderTypeUpdate) Where(ps ...predicate.WorkOrderType) *WorkOrderTypeUpdate {
	wotu.predicates = append(wotu.predicates, ps...)
	return wotu
}

// SetName sets the name field.
func (wotu *WorkOrderTypeUpdate) SetName(s string) *WorkOrderTypeUpdate {
	wotu.mutation.SetName(s)
	return wotu
}

// SetDescription sets the description field.
func (wotu *WorkOrderTypeUpdate) SetDescription(s string) *WorkOrderTypeUpdate {
	wotu.mutation.SetDescription(s)
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
	wotu.mutation.ClearDescription()
	return wotu
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (wotu *WorkOrderTypeUpdate) AddWorkOrderIDs(ids ...int) *WorkOrderTypeUpdate {
	wotu.mutation.AddWorkOrderIDs(ids...)
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
	wotu.mutation.AddPropertyTypeIDs(ids...)
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
	wotu.mutation.AddDefinitionIDs(ids...)
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
	wotu.mutation.AddCheckListCategoryIDs(ids...)
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
	wotu.mutation.AddCheckListDefinitionIDs(ids...)
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
	wotu.mutation.RemoveWorkOrderIDs(ids...)
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
	wotu.mutation.RemovePropertyTypeIDs(ids...)
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
	wotu.mutation.RemoveDefinitionIDs(ids...)
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
	wotu.mutation.RemoveCheckListCategoryIDs(ids...)
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
	wotu.mutation.RemoveCheckListDefinitionIDs(ids...)
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
	if _, ok := wotu.mutation.UpdateTime(); !ok {
		v := workordertype.UpdateDefaultUpdateTime()
		wotu.mutation.SetUpdateTime(v)
	}

	var (
		err      error
		affected int
	)
	if len(wotu.hooks) == 0 {
		affected, err = wotu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*WorkOrderTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			wotu.mutation = mutation
			affected, err = wotu.sqlSave(ctx)
			return affected, err
		})
		for i := len(wotu.hooks) - 1; i >= 0; i-- {
			mut = wotu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, wotu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
	if value, ok := wotu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: workordertype.FieldUpdateTime,
		})
	}
	if value, ok := wotu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workordertype.FieldName,
		})
	}
	if value, ok := wotu.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workordertype.FieldDescription,
		})
	}
	if wotu.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: workordertype.FieldDescription,
		})
	}
	if nodes := wotu.mutation.RemovedWorkOrdersIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotu.mutation.WorkOrdersIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wotu.mutation.RemovedPropertyTypesIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotu.mutation.PropertyTypesIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wotu.mutation.RemovedDefinitionsIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotu.mutation.DefinitionsIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wotu.mutation.RemovedCheckListCategoriesIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotu.mutation.CheckListCategoriesIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wotu.mutation.RemovedCheckListDefinitionsIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotu.mutation.CheckListDefinitionsIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
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
	hooks    []Hook
	mutation *WorkOrderTypeMutation
}

// SetName sets the name field.
func (wotuo *WorkOrderTypeUpdateOne) SetName(s string) *WorkOrderTypeUpdateOne {
	wotuo.mutation.SetName(s)
	return wotuo
}

// SetDescription sets the description field.
func (wotuo *WorkOrderTypeUpdateOne) SetDescription(s string) *WorkOrderTypeUpdateOne {
	wotuo.mutation.SetDescription(s)
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
	wotuo.mutation.ClearDescription()
	return wotuo
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (wotuo *WorkOrderTypeUpdateOne) AddWorkOrderIDs(ids ...int) *WorkOrderTypeUpdateOne {
	wotuo.mutation.AddWorkOrderIDs(ids...)
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
	wotuo.mutation.AddPropertyTypeIDs(ids...)
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
	wotuo.mutation.AddDefinitionIDs(ids...)
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
	wotuo.mutation.AddCheckListCategoryIDs(ids...)
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
	wotuo.mutation.AddCheckListDefinitionIDs(ids...)
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
	wotuo.mutation.RemoveWorkOrderIDs(ids...)
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
	wotuo.mutation.RemovePropertyTypeIDs(ids...)
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
	wotuo.mutation.RemoveDefinitionIDs(ids...)
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
	wotuo.mutation.RemoveCheckListCategoryIDs(ids...)
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
	wotuo.mutation.RemoveCheckListDefinitionIDs(ids...)
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
	if _, ok := wotuo.mutation.UpdateTime(); !ok {
		v := workordertype.UpdateDefaultUpdateTime()
		wotuo.mutation.SetUpdateTime(v)
	}

	var (
		err  error
		node *WorkOrderType
	)
	if len(wotuo.hooks) == 0 {
		node, err = wotuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*WorkOrderTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			wotuo.mutation = mutation
			node, err = wotuo.sqlSave(ctx)
			return node, err
		})
		for i := len(wotuo.hooks) - 1; i >= 0; i-- {
			mut = wotuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, wotuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
				Type:   field.TypeInt,
				Column: workordertype.FieldID,
			},
		},
	}
	id, ok := wotuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing WorkOrderType.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := wotuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: workordertype.FieldUpdateTime,
		})
	}
	if value, ok := wotuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workordertype.FieldName,
		})
	}
	if value, ok := wotuo.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workordertype.FieldDescription,
		})
	}
	if wotuo.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: workordertype.FieldDescription,
		})
	}
	if nodes := wotuo.mutation.RemovedWorkOrdersIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotuo.mutation.WorkOrdersIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wotuo.mutation.RemovedPropertyTypesIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotuo.mutation.PropertyTypesIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wotuo.mutation.RemovedDefinitionsIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotuo.mutation.DefinitionsIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wotuo.mutation.RemovedCheckListCategoriesIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotuo.mutation.CheckListCategoriesIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wotuo.mutation.RemovedCheckListDefinitionsIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotuo.mutation.CheckListDefinitionsIDs(); len(nodes) > 0 {
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
		for _, k := range nodes {
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
