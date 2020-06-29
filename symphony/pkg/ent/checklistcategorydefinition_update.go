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
	"github.com/facebookincubator/symphony/pkg/ent/checklistitemdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/workordertemplate"
	"github.com/facebookincubator/symphony/pkg/ent/workordertype"
)

// CheckListCategoryDefinitionUpdate is the builder for updating CheckListCategoryDefinition entities.
type CheckListCategoryDefinitionUpdate struct {
	config
	hooks      []Hook
	mutation   *CheckListCategoryDefinitionMutation
	predicates []predicate.CheckListCategoryDefinition
}

// Where adds a new predicate for the builder.
func (clcdu *CheckListCategoryDefinitionUpdate) Where(ps ...predicate.CheckListCategoryDefinition) *CheckListCategoryDefinitionUpdate {
	clcdu.predicates = append(clcdu.predicates, ps...)
	return clcdu
}

// SetTitle sets the title field.
func (clcdu *CheckListCategoryDefinitionUpdate) SetTitle(s string) *CheckListCategoryDefinitionUpdate {
	clcdu.mutation.SetTitle(s)
	return clcdu
}

// SetDescription sets the description field.
func (clcdu *CheckListCategoryDefinitionUpdate) SetDescription(s string) *CheckListCategoryDefinitionUpdate {
	clcdu.mutation.SetDescription(s)
	return clcdu
}

// SetNillableDescription sets the description field if the given value is not nil.
func (clcdu *CheckListCategoryDefinitionUpdate) SetNillableDescription(s *string) *CheckListCategoryDefinitionUpdate {
	if s != nil {
		clcdu.SetDescription(*s)
	}
	return clcdu
}

// ClearDescription clears the value of description.
func (clcdu *CheckListCategoryDefinitionUpdate) ClearDescription() *CheckListCategoryDefinitionUpdate {
	clcdu.mutation.ClearDescription()
	return clcdu
}

// AddCheckListItemDefinitionIDs adds the check_list_item_definitions edge to CheckListItemDefinition by ids.
func (clcdu *CheckListCategoryDefinitionUpdate) AddCheckListItemDefinitionIDs(ids ...int) *CheckListCategoryDefinitionUpdate {
	clcdu.mutation.AddCheckListItemDefinitionIDs(ids...)
	return clcdu
}

// AddCheckListItemDefinitions adds the check_list_item_definitions edges to CheckListItemDefinition.
func (clcdu *CheckListCategoryDefinitionUpdate) AddCheckListItemDefinitions(c ...*CheckListItemDefinition) *CheckListCategoryDefinitionUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return clcdu.AddCheckListItemDefinitionIDs(ids...)
}

// SetWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id.
func (clcdu *CheckListCategoryDefinitionUpdate) SetWorkOrderTypeID(id int) *CheckListCategoryDefinitionUpdate {
	clcdu.mutation.SetWorkOrderTypeID(id)
	return clcdu
}

// SetNillableWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id if the given value is not nil.
func (clcdu *CheckListCategoryDefinitionUpdate) SetNillableWorkOrderTypeID(id *int) *CheckListCategoryDefinitionUpdate {
	if id != nil {
		clcdu = clcdu.SetWorkOrderTypeID(*id)
	}
	return clcdu
}

// SetWorkOrderType sets the work_order_type edge to WorkOrderType.
func (clcdu *CheckListCategoryDefinitionUpdate) SetWorkOrderType(w *WorkOrderType) *CheckListCategoryDefinitionUpdate {
	return clcdu.SetWorkOrderTypeID(w.ID)
}

// SetWorkOrderTemplateID sets the work_order_template edge to WorkOrderTemplate by id.
func (clcdu *CheckListCategoryDefinitionUpdate) SetWorkOrderTemplateID(id int) *CheckListCategoryDefinitionUpdate {
	clcdu.mutation.SetWorkOrderTemplateID(id)
	return clcdu
}

// SetNillableWorkOrderTemplateID sets the work_order_template edge to WorkOrderTemplate by id if the given value is not nil.
func (clcdu *CheckListCategoryDefinitionUpdate) SetNillableWorkOrderTemplateID(id *int) *CheckListCategoryDefinitionUpdate {
	if id != nil {
		clcdu = clcdu.SetWorkOrderTemplateID(*id)
	}
	return clcdu
}

// SetWorkOrderTemplate sets the work_order_template edge to WorkOrderTemplate.
func (clcdu *CheckListCategoryDefinitionUpdate) SetWorkOrderTemplate(w *WorkOrderTemplate) *CheckListCategoryDefinitionUpdate {
	return clcdu.SetWorkOrderTemplateID(w.ID)
}

// RemoveCheckListItemDefinitionIDs removes the check_list_item_definitions edge to CheckListItemDefinition by ids.
func (clcdu *CheckListCategoryDefinitionUpdate) RemoveCheckListItemDefinitionIDs(ids ...int) *CheckListCategoryDefinitionUpdate {
	clcdu.mutation.RemoveCheckListItemDefinitionIDs(ids...)
	return clcdu
}

// RemoveCheckListItemDefinitions removes check_list_item_definitions edges to CheckListItemDefinition.
func (clcdu *CheckListCategoryDefinitionUpdate) RemoveCheckListItemDefinitions(c ...*CheckListItemDefinition) *CheckListCategoryDefinitionUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return clcdu.RemoveCheckListItemDefinitionIDs(ids...)
}

// ClearWorkOrderType clears the work_order_type edge to WorkOrderType.
func (clcdu *CheckListCategoryDefinitionUpdate) ClearWorkOrderType() *CheckListCategoryDefinitionUpdate {
	clcdu.mutation.ClearWorkOrderType()
	return clcdu
}

// ClearWorkOrderTemplate clears the work_order_template edge to WorkOrderTemplate.
func (clcdu *CheckListCategoryDefinitionUpdate) ClearWorkOrderTemplate() *CheckListCategoryDefinitionUpdate {
	clcdu.mutation.ClearWorkOrderTemplate()
	return clcdu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (clcdu *CheckListCategoryDefinitionUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := clcdu.mutation.UpdateTime(); !ok {
		v := checklistcategorydefinition.UpdateDefaultUpdateTime()
		clcdu.mutation.SetUpdateTime(v)
	}
	if v, ok := clcdu.mutation.Title(); ok {
		if err := checklistcategorydefinition.TitleValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"title\": %v", err)
		}
	}

	var (
		err      error
		affected int
	)
	if len(clcdu.hooks) == 0 {
		affected, err = clcdu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CheckListCategoryDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			clcdu.mutation = mutation
			affected, err = clcdu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(clcdu.hooks) - 1; i >= 0; i-- {
			mut = clcdu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, clcdu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (clcdu *CheckListCategoryDefinitionUpdate) SaveX(ctx context.Context) int {
	affected, err := clcdu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (clcdu *CheckListCategoryDefinitionUpdate) Exec(ctx context.Context) error {
	_, err := clcdu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (clcdu *CheckListCategoryDefinitionUpdate) ExecX(ctx context.Context) {
	if err := clcdu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (clcdu *CheckListCategoryDefinitionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   checklistcategorydefinition.Table,
			Columns: checklistcategorydefinition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistcategorydefinition.FieldID,
			},
		},
	}
	if ps := clcdu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := clcdu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: checklistcategorydefinition.FieldUpdateTime,
		})
	}
	if value, ok := clcdu.mutation.Title(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistcategorydefinition.FieldTitle,
		})
	}
	if value, ok := clcdu.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistcategorydefinition.FieldDescription,
		})
	}
	if clcdu.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistcategorydefinition.FieldDescription,
		})
	}
	if nodes := clcdu.mutation.RemovedCheckListItemDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   checklistcategorydefinition.CheckListItemDefinitionsTable,
			Columns: []string{checklistcategorydefinition.CheckListItemDefinitionsColumn},
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
	if nodes := clcdu.mutation.CheckListItemDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   checklistcategorydefinition.CheckListItemDefinitionsTable,
			Columns: []string{checklistcategorydefinition.CheckListItemDefinitionsColumn},
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
	if clcdu.mutation.WorkOrderTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistcategorydefinition.WorkOrderTypeTable,
			Columns: []string{checklistcategorydefinition.WorkOrderTypeColumn},
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
	if nodes := clcdu.mutation.WorkOrderTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistcategorydefinition.WorkOrderTypeTable,
			Columns: []string{checklistcategorydefinition.WorkOrderTypeColumn},
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
	if clcdu.mutation.WorkOrderTemplateCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistcategorydefinition.WorkOrderTemplateTable,
			Columns: []string{checklistcategorydefinition.WorkOrderTemplateColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workordertemplate.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := clcdu.mutation.WorkOrderTemplateIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistcategorydefinition.WorkOrderTemplateTable,
			Columns: []string{checklistcategorydefinition.WorkOrderTemplateColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workordertemplate.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, clcdu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{checklistcategorydefinition.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// CheckListCategoryDefinitionUpdateOne is the builder for updating a single CheckListCategoryDefinition entity.
type CheckListCategoryDefinitionUpdateOne struct {
	config
	hooks    []Hook
	mutation *CheckListCategoryDefinitionMutation
}

// SetTitle sets the title field.
func (clcduo *CheckListCategoryDefinitionUpdateOne) SetTitle(s string) *CheckListCategoryDefinitionUpdateOne {
	clcduo.mutation.SetTitle(s)
	return clcduo
}

// SetDescription sets the description field.
func (clcduo *CheckListCategoryDefinitionUpdateOne) SetDescription(s string) *CheckListCategoryDefinitionUpdateOne {
	clcduo.mutation.SetDescription(s)
	return clcduo
}

// SetNillableDescription sets the description field if the given value is not nil.
func (clcduo *CheckListCategoryDefinitionUpdateOne) SetNillableDescription(s *string) *CheckListCategoryDefinitionUpdateOne {
	if s != nil {
		clcduo.SetDescription(*s)
	}
	return clcduo
}

// ClearDescription clears the value of description.
func (clcduo *CheckListCategoryDefinitionUpdateOne) ClearDescription() *CheckListCategoryDefinitionUpdateOne {
	clcduo.mutation.ClearDescription()
	return clcduo
}

// AddCheckListItemDefinitionIDs adds the check_list_item_definitions edge to CheckListItemDefinition by ids.
func (clcduo *CheckListCategoryDefinitionUpdateOne) AddCheckListItemDefinitionIDs(ids ...int) *CheckListCategoryDefinitionUpdateOne {
	clcduo.mutation.AddCheckListItemDefinitionIDs(ids...)
	return clcduo
}

// AddCheckListItemDefinitions adds the check_list_item_definitions edges to CheckListItemDefinition.
func (clcduo *CheckListCategoryDefinitionUpdateOne) AddCheckListItemDefinitions(c ...*CheckListItemDefinition) *CheckListCategoryDefinitionUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return clcduo.AddCheckListItemDefinitionIDs(ids...)
}

// SetWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id.
func (clcduo *CheckListCategoryDefinitionUpdateOne) SetWorkOrderTypeID(id int) *CheckListCategoryDefinitionUpdateOne {
	clcduo.mutation.SetWorkOrderTypeID(id)
	return clcduo
}

// SetNillableWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id if the given value is not nil.
func (clcduo *CheckListCategoryDefinitionUpdateOne) SetNillableWorkOrderTypeID(id *int) *CheckListCategoryDefinitionUpdateOne {
	if id != nil {
		clcduo = clcduo.SetWorkOrderTypeID(*id)
	}
	return clcduo
}

// SetWorkOrderType sets the work_order_type edge to WorkOrderType.
func (clcduo *CheckListCategoryDefinitionUpdateOne) SetWorkOrderType(w *WorkOrderType) *CheckListCategoryDefinitionUpdateOne {
	return clcduo.SetWorkOrderTypeID(w.ID)
}

// SetWorkOrderTemplateID sets the work_order_template edge to WorkOrderTemplate by id.
func (clcduo *CheckListCategoryDefinitionUpdateOne) SetWorkOrderTemplateID(id int) *CheckListCategoryDefinitionUpdateOne {
	clcduo.mutation.SetWorkOrderTemplateID(id)
	return clcduo
}

// SetNillableWorkOrderTemplateID sets the work_order_template edge to WorkOrderTemplate by id if the given value is not nil.
func (clcduo *CheckListCategoryDefinitionUpdateOne) SetNillableWorkOrderTemplateID(id *int) *CheckListCategoryDefinitionUpdateOne {
	if id != nil {
		clcduo = clcduo.SetWorkOrderTemplateID(*id)
	}
	return clcduo
}

// SetWorkOrderTemplate sets the work_order_template edge to WorkOrderTemplate.
func (clcduo *CheckListCategoryDefinitionUpdateOne) SetWorkOrderTemplate(w *WorkOrderTemplate) *CheckListCategoryDefinitionUpdateOne {
	return clcduo.SetWorkOrderTemplateID(w.ID)
}

// RemoveCheckListItemDefinitionIDs removes the check_list_item_definitions edge to CheckListItemDefinition by ids.
func (clcduo *CheckListCategoryDefinitionUpdateOne) RemoveCheckListItemDefinitionIDs(ids ...int) *CheckListCategoryDefinitionUpdateOne {
	clcduo.mutation.RemoveCheckListItemDefinitionIDs(ids...)
	return clcduo
}

// RemoveCheckListItemDefinitions removes check_list_item_definitions edges to CheckListItemDefinition.
func (clcduo *CheckListCategoryDefinitionUpdateOne) RemoveCheckListItemDefinitions(c ...*CheckListItemDefinition) *CheckListCategoryDefinitionUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return clcduo.RemoveCheckListItemDefinitionIDs(ids...)
}

// ClearWorkOrderType clears the work_order_type edge to WorkOrderType.
func (clcduo *CheckListCategoryDefinitionUpdateOne) ClearWorkOrderType() *CheckListCategoryDefinitionUpdateOne {
	clcduo.mutation.ClearWorkOrderType()
	return clcduo
}

// ClearWorkOrderTemplate clears the work_order_template edge to WorkOrderTemplate.
func (clcduo *CheckListCategoryDefinitionUpdateOne) ClearWorkOrderTemplate() *CheckListCategoryDefinitionUpdateOne {
	clcduo.mutation.ClearWorkOrderTemplate()
	return clcduo
}

// Save executes the query and returns the updated entity.
func (clcduo *CheckListCategoryDefinitionUpdateOne) Save(ctx context.Context) (*CheckListCategoryDefinition, error) {
	if _, ok := clcduo.mutation.UpdateTime(); !ok {
		v := checklistcategorydefinition.UpdateDefaultUpdateTime()
		clcduo.mutation.SetUpdateTime(v)
	}
	if v, ok := clcduo.mutation.Title(); ok {
		if err := checklistcategorydefinition.TitleValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"title\": %v", err)
		}
	}

	var (
		err  error
		node *CheckListCategoryDefinition
	)
	if len(clcduo.hooks) == 0 {
		node, err = clcduo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CheckListCategoryDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			clcduo.mutation = mutation
			node, err = clcduo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(clcduo.hooks) - 1; i >= 0; i-- {
			mut = clcduo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, clcduo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (clcduo *CheckListCategoryDefinitionUpdateOne) SaveX(ctx context.Context) *CheckListCategoryDefinition {
	clcd, err := clcduo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return clcd
}

// Exec executes the query on the entity.
func (clcduo *CheckListCategoryDefinitionUpdateOne) Exec(ctx context.Context) error {
	_, err := clcduo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (clcduo *CheckListCategoryDefinitionUpdateOne) ExecX(ctx context.Context) {
	if err := clcduo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (clcduo *CheckListCategoryDefinitionUpdateOne) sqlSave(ctx context.Context) (clcd *CheckListCategoryDefinition, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   checklistcategorydefinition.Table,
			Columns: checklistcategorydefinition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistcategorydefinition.FieldID,
			},
		},
	}
	id, ok := clcduo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing CheckListCategoryDefinition.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := clcduo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: checklistcategorydefinition.FieldUpdateTime,
		})
	}
	if value, ok := clcduo.mutation.Title(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistcategorydefinition.FieldTitle,
		})
	}
	if value, ok := clcduo.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistcategorydefinition.FieldDescription,
		})
	}
	if clcduo.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistcategorydefinition.FieldDescription,
		})
	}
	if nodes := clcduo.mutation.RemovedCheckListItemDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   checklistcategorydefinition.CheckListItemDefinitionsTable,
			Columns: []string{checklistcategorydefinition.CheckListItemDefinitionsColumn},
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
	if nodes := clcduo.mutation.CheckListItemDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   checklistcategorydefinition.CheckListItemDefinitionsTable,
			Columns: []string{checklistcategorydefinition.CheckListItemDefinitionsColumn},
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
	if clcduo.mutation.WorkOrderTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistcategorydefinition.WorkOrderTypeTable,
			Columns: []string{checklistcategorydefinition.WorkOrderTypeColumn},
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
	if nodes := clcduo.mutation.WorkOrderTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistcategorydefinition.WorkOrderTypeTable,
			Columns: []string{checklistcategorydefinition.WorkOrderTypeColumn},
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
	if clcduo.mutation.WorkOrderTemplateCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistcategorydefinition.WorkOrderTemplateTable,
			Columns: []string{checklistcategorydefinition.WorkOrderTemplateColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workordertemplate.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := clcduo.mutation.WorkOrderTemplateIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistcategorydefinition.WorkOrderTemplateTable,
			Columns: []string{checklistcategorydefinition.WorkOrderTemplateColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workordertemplate.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	clcd = &CheckListCategoryDefinition{config: clcduo.config}
	_spec.Assign = clcd.assignValues
	_spec.ScanValues = clcd.scanValues()
	if err = sqlgraph.UpdateNode(ctx, clcduo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{checklistcategorydefinition.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return clcd, nil
}
