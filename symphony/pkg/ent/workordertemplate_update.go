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
	"github.com/facebookincubator/symphony/pkg/ent/checklistcategorydefinition"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/propertytype"
	"github.com/facebookincubator/symphony/pkg/ent/workordertemplate"
	"github.com/facebookincubator/symphony/pkg/ent/workordertype"
)

// WorkOrderTemplateUpdate is the builder for updating WorkOrderTemplate entities.
type WorkOrderTemplateUpdate struct {
	config
	hooks      []Hook
	mutation   *WorkOrderTemplateMutation
	predicates []predicate.WorkOrderTemplate
}

// Where adds a new predicate for the builder.
func (wotu *WorkOrderTemplateUpdate) Where(ps ...predicate.WorkOrderTemplate) *WorkOrderTemplateUpdate {
	wotu.predicates = append(wotu.predicates, ps...)
	return wotu
}

// SetName sets the name field.
func (wotu *WorkOrderTemplateUpdate) SetName(s string) *WorkOrderTemplateUpdate {
	wotu.mutation.SetName(s)
	return wotu
}

// SetDescription sets the description field.
func (wotu *WorkOrderTemplateUpdate) SetDescription(s string) *WorkOrderTemplateUpdate {
	wotu.mutation.SetDescription(s)
	return wotu
}

// SetNillableDescription sets the description field if the given value is not nil.
func (wotu *WorkOrderTemplateUpdate) SetNillableDescription(s *string) *WorkOrderTemplateUpdate {
	if s != nil {
		wotu.SetDescription(*s)
	}
	return wotu
}

// ClearDescription clears the value of description.
func (wotu *WorkOrderTemplateUpdate) ClearDescription() *WorkOrderTemplateUpdate {
	wotu.mutation.ClearDescription()
	return wotu
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (wotu *WorkOrderTemplateUpdate) AddPropertyTypeIDs(ids ...int) *WorkOrderTemplateUpdate {
	wotu.mutation.AddPropertyTypeIDs(ids...)
	return wotu
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (wotu *WorkOrderTemplateUpdate) AddPropertyTypes(p ...*PropertyType) *WorkOrderTemplateUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wotu.AddPropertyTypeIDs(ids...)
}

// AddCheckListCategoryDefinitionIDs adds the check_list_category_definitions edge to CheckListCategoryDefinition by ids.
func (wotu *WorkOrderTemplateUpdate) AddCheckListCategoryDefinitionIDs(ids ...int) *WorkOrderTemplateUpdate {
	wotu.mutation.AddCheckListCategoryDefinitionIDs(ids...)
	return wotu
}

// AddCheckListCategoryDefinitions adds the check_list_category_definitions edges to CheckListCategoryDefinition.
func (wotu *WorkOrderTemplateUpdate) AddCheckListCategoryDefinitions(c ...*CheckListCategoryDefinition) *WorkOrderTemplateUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotu.AddCheckListCategoryDefinitionIDs(ids...)
}

// SetTypeID sets the type edge to WorkOrderType by id.
func (wotu *WorkOrderTemplateUpdate) SetTypeID(id int) *WorkOrderTemplateUpdate {
	wotu.mutation.SetTypeID(id)
	return wotu
}

// SetNillableTypeID sets the type edge to WorkOrderType by id if the given value is not nil.
func (wotu *WorkOrderTemplateUpdate) SetNillableTypeID(id *int) *WorkOrderTemplateUpdate {
	if id != nil {
		wotu = wotu.SetTypeID(*id)
	}
	return wotu
}

// SetType sets the type edge to WorkOrderType.
func (wotu *WorkOrderTemplateUpdate) SetType(w *WorkOrderType) *WorkOrderTemplateUpdate {
	return wotu.SetTypeID(w.ID)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (wotu *WorkOrderTemplateUpdate) RemovePropertyTypeIDs(ids ...int) *WorkOrderTemplateUpdate {
	wotu.mutation.RemovePropertyTypeIDs(ids...)
	return wotu
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (wotu *WorkOrderTemplateUpdate) RemovePropertyTypes(p ...*PropertyType) *WorkOrderTemplateUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wotu.RemovePropertyTypeIDs(ids...)
}

// RemoveCheckListCategoryDefinitionIDs removes the check_list_category_definitions edge to CheckListCategoryDefinition by ids.
func (wotu *WorkOrderTemplateUpdate) RemoveCheckListCategoryDefinitionIDs(ids ...int) *WorkOrderTemplateUpdate {
	wotu.mutation.RemoveCheckListCategoryDefinitionIDs(ids...)
	return wotu
}

// RemoveCheckListCategoryDefinitions removes check_list_category_definitions edges to CheckListCategoryDefinition.
func (wotu *WorkOrderTemplateUpdate) RemoveCheckListCategoryDefinitions(c ...*CheckListCategoryDefinition) *WorkOrderTemplateUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotu.RemoveCheckListCategoryDefinitionIDs(ids...)
}

// ClearType clears the type edge to WorkOrderType.
func (wotu *WorkOrderTemplateUpdate) ClearType() *WorkOrderTemplateUpdate {
	wotu.mutation.ClearType()
	return wotu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (wotu *WorkOrderTemplateUpdate) Save(ctx context.Context) (int, error) {

	var (
		err      error
		affected int
	)
	if len(wotu.hooks) == 0 {
		affected, err = wotu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*WorkOrderTemplateMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			wotu.mutation = mutation
			affected, err = wotu.sqlSave(ctx)
			mutation.done = true
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
func (wotu *WorkOrderTemplateUpdate) SaveX(ctx context.Context) int {
	affected, err := wotu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (wotu *WorkOrderTemplateUpdate) Exec(ctx context.Context) error {
	_, err := wotu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wotu *WorkOrderTemplateUpdate) ExecX(ctx context.Context) {
	if err := wotu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (wotu *WorkOrderTemplateUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   workordertemplate.Table,
			Columns: workordertemplate.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: workordertemplate.FieldID,
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
	if value, ok := wotu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workordertemplate.FieldName,
		})
	}
	if value, ok := wotu.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workordertemplate.FieldDescription,
		})
	}
	if wotu.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: workordertemplate.FieldDescription,
		})
	}
	if nodes := wotu.mutation.RemovedPropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertemplate.PropertyTypesTable,
			Columns: []string{workordertemplate.PropertyTypesColumn},
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
			Table:   workordertemplate.PropertyTypesTable,
			Columns: []string{workordertemplate.PropertyTypesColumn},
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
	if nodes := wotu.mutation.RemovedCheckListCategoryDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertemplate.CheckListCategoryDefinitionsTable,
			Columns: []string{workordertemplate.CheckListCategoryDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistcategorydefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotu.mutation.CheckListCategoryDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertemplate.CheckListCategoryDefinitionsTable,
			Columns: []string{workordertemplate.CheckListCategoryDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistcategorydefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if wotu.mutation.TypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workordertemplate.TypeTable,
			Columns: []string{workordertemplate.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workordertype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotu.mutation.TypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workordertemplate.TypeTable,
			Columns: []string{workordertemplate.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workordertype.FieldID,
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
			err = &NotFoundError{workordertemplate.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// WorkOrderTemplateUpdateOne is the builder for updating a single WorkOrderTemplate entity.
type WorkOrderTemplateUpdateOne struct {
	config
	hooks    []Hook
	mutation *WorkOrderTemplateMutation
}

// SetName sets the name field.
func (wotuo *WorkOrderTemplateUpdateOne) SetName(s string) *WorkOrderTemplateUpdateOne {
	wotuo.mutation.SetName(s)
	return wotuo
}

// SetDescription sets the description field.
func (wotuo *WorkOrderTemplateUpdateOne) SetDescription(s string) *WorkOrderTemplateUpdateOne {
	wotuo.mutation.SetDescription(s)
	return wotuo
}

// SetNillableDescription sets the description field if the given value is not nil.
func (wotuo *WorkOrderTemplateUpdateOne) SetNillableDescription(s *string) *WorkOrderTemplateUpdateOne {
	if s != nil {
		wotuo.SetDescription(*s)
	}
	return wotuo
}

// ClearDescription clears the value of description.
func (wotuo *WorkOrderTemplateUpdateOne) ClearDescription() *WorkOrderTemplateUpdateOne {
	wotuo.mutation.ClearDescription()
	return wotuo
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (wotuo *WorkOrderTemplateUpdateOne) AddPropertyTypeIDs(ids ...int) *WorkOrderTemplateUpdateOne {
	wotuo.mutation.AddPropertyTypeIDs(ids...)
	return wotuo
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (wotuo *WorkOrderTemplateUpdateOne) AddPropertyTypes(p ...*PropertyType) *WorkOrderTemplateUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wotuo.AddPropertyTypeIDs(ids...)
}

// AddCheckListCategoryDefinitionIDs adds the check_list_category_definitions edge to CheckListCategoryDefinition by ids.
func (wotuo *WorkOrderTemplateUpdateOne) AddCheckListCategoryDefinitionIDs(ids ...int) *WorkOrderTemplateUpdateOne {
	wotuo.mutation.AddCheckListCategoryDefinitionIDs(ids...)
	return wotuo
}

// AddCheckListCategoryDefinitions adds the check_list_category_definitions edges to CheckListCategoryDefinition.
func (wotuo *WorkOrderTemplateUpdateOne) AddCheckListCategoryDefinitions(c ...*CheckListCategoryDefinition) *WorkOrderTemplateUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotuo.AddCheckListCategoryDefinitionIDs(ids...)
}

// SetTypeID sets the type edge to WorkOrderType by id.
func (wotuo *WorkOrderTemplateUpdateOne) SetTypeID(id int) *WorkOrderTemplateUpdateOne {
	wotuo.mutation.SetTypeID(id)
	return wotuo
}

// SetNillableTypeID sets the type edge to WorkOrderType by id if the given value is not nil.
func (wotuo *WorkOrderTemplateUpdateOne) SetNillableTypeID(id *int) *WorkOrderTemplateUpdateOne {
	if id != nil {
		wotuo = wotuo.SetTypeID(*id)
	}
	return wotuo
}

// SetType sets the type edge to WorkOrderType.
func (wotuo *WorkOrderTemplateUpdateOne) SetType(w *WorkOrderType) *WorkOrderTemplateUpdateOne {
	return wotuo.SetTypeID(w.ID)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (wotuo *WorkOrderTemplateUpdateOne) RemovePropertyTypeIDs(ids ...int) *WorkOrderTemplateUpdateOne {
	wotuo.mutation.RemovePropertyTypeIDs(ids...)
	return wotuo
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (wotuo *WorkOrderTemplateUpdateOne) RemovePropertyTypes(p ...*PropertyType) *WorkOrderTemplateUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wotuo.RemovePropertyTypeIDs(ids...)
}

// RemoveCheckListCategoryDefinitionIDs removes the check_list_category_definitions edge to CheckListCategoryDefinition by ids.
func (wotuo *WorkOrderTemplateUpdateOne) RemoveCheckListCategoryDefinitionIDs(ids ...int) *WorkOrderTemplateUpdateOne {
	wotuo.mutation.RemoveCheckListCategoryDefinitionIDs(ids...)
	return wotuo
}

// RemoveCheckListCategoryDefinitions removes check_list_category_definitions edges to CheckListCategoryDefinition.
func (wotuo *WorkOrderTemplateUpdateOne) RemoveCheckListCategoryDefinitions(c ...*CheckListCategoryDefinition) *WorkOrderTemplateUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotuo.RemoveCheckListCategoryDefinitionIDs(ids...)
}

// ClearType clears the type edge to WorkOrderType.
func (wotuo *WorkOrderTemplateUpdateOne) ClearType() *WorkOrderTemplateUpdateOne {
	wotuo.mutation.ClearType()
	return wotuo
}

// Save executes the query and returns the updated entity.
func (wotuo *WorkOrderTemplateUpdateOne) Save(ctx context.Context) (*WorkOrderTemplate, error) {

	var (
		err  error
		node *WorkOrderTemplate
	)
	if len(wotuo.hooks) == 0 {
		node, err = wotuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*WorkOrderTemplateMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			wotuo.mutation = mutation
			node, err = wotuo.sqlSave(ctx)
			mutation.done = true
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
func (wotuo *WorkOrderTemplateUpdateOne) SaveX(ctx context.Context) *WorkOrderTemplate {
	wot, err := wotuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return wot
}

// Exec executes the query on the entity.
func (wotuo *WorkOrderTemplateUpdateOne) Exec(ctx context.Context) error {
	_, err := wotuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wotuo *WorkOrderTemplateUpdateOne) ExecX(ctx context.Context) {
	if err := wotuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (wotuo *WorkOrderTemplateUpdateOne) sqlSave(ctx context.Context) (wot *WorkOrderTemplate, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   workordertemplate.Table,
			Columns: workordertemplate.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: workordertemplate.FieldID,
			},
		},
	}
	id, ok := wotuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing WorkOrderTemplate.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := wotuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workordertemplate.FieldName,
		})
	}
	if value, ok := wotuo.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workordertemplate.FieldDescription,
		})
	}
	if wotuo.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: workordertemplate.FieldDescription,
		})
	}
	if nodes := wotuo.mutation.RemovedPropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertemplate.PropertyTypesTable,
			Columns: []string{workordertemplate.PropertyTypesColumn},
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
			Table:   workordertemplate.PropertyTypesTable,
			Columns: []string{workordertemplate.PropertyTypesColumn},
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
	if nodes := wotuo.mutation.RemovedCheckListCategoryDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertemplate.CheckListCategoryDefinitionsTable,
			Columns: []string{workordertemplate.CheckListCategoryDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistcategorydefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotuo.mutation.CheckListCategoryDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertemplate.CheckListCategoryDefinitionsTable,
			Columns: []string{workordertemplate.CheckListCategoryDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistcategorydefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if wotuo.mutation.TypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workordertemplate.TypeTable,
			Columns: []string{workordertemplate.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workordertype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wotuo.mutation.TypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workordertemplate.TypeTable,
			Columns: []string{workordertemplate.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workordertype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	wot = &WorkOrderTemplate{config: wotuo.config}
	_spec.Assign = wot.assignValues
	_spec.ScanValues = wot.scanValues()
	if err = sqlgraph.UpdateNode(ctx, wotuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{workordertemplate.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return wot, nil
}
