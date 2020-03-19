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
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// CheckListItemUpdate is the builder for updating CheckListItem entities.
type CheckListItemUpdate struct {
	config
	hooks      []Hook
	mutation   *CheckListItemMutation
	predicates []predicate.CheckListItem
}

// Where adds a new predicate for the builder.
func (cliu *CheckListItemUpdate) Where(ps ...predicate.CheckListItem) *CheckListItemUpdate {
	cliu.predicates = append(cliu.predicates, ps...)
	return cliu
}

// SetTitle sets the title field.
func (cliu *CheckListItemUpdate) SetTitle(s string) *CheckListItemUpdate {
	cliu.mutation.SetTitle(s)
	return cliu
}

// SetType sets the type field.
func (cliu *CheckListItemUpdate) SetType(s string) *CheckListItemUpdate {
	cliu.mutation.SetType(s)
	return cliu
}

// SetIndex sets the index field.
func (cliu *CheckListItemUpdate) SetIndex(i int) *CheckListItemUpdate {
	cliu.mutation.ResetIndex()
	cliu.mutation.SetIndex(i)
	return cliu
}

// SetNillableIndex sets the index field if the given value is not nil.
func (cliu *CheckListItemUpdate) SetNillableIndex(i *int) *CheckListItemUpdate {
	if i != nil {
		cliu.SetIndex(*i)
	}
	return cliu
}

// AddIndex adds i to index.
func (cliu *CheckListItemUpdate) AddIndex(i int) *CheckListItemUpdate {
	cliu.mutation.AddIndex(i)
	return cliu
}

// ClearIndex clears the value of index.
func (cliu *CheckListItemUpdate) ClearIndex() *CheckListItemUpdate {
	cliu.mutation.ClearIndex()
	return cliu
}

// SetChecked sets the checked field.
func (cliu *CheckListItemUpdate) SetChecked(b bool) *CheckListItemUpdate {
	cliu.mutation.SetChecked(b)
	return cliu
}

// SetNillableChecked sets the checked field if the given value is not nil.
func (cliu *CheckListItemUpdate) SetNillableChecked(b *bool) *CheckListItemUpdate {
	if b != nil {
		cliu.SetChecked(*b)
	}
	return cliu
}

// ClearChecked clears the value of checked.
func (cliu *CheckListItemUpdate) ClearChecked() *CheckListItemUpdate {
	cliu.mutation.ClearChecked()
	return cliu
}

// SetStringVal sets the string_val field.
func (cliu *CheckListItemUpdate) SetStringVal(s string) *CheckListItemUpdate {
	cliu.mutation.SetStringVal(s)
	return cliu
}

// SetNillableStringVal sets the string_val field if the given value is not nil.
func (cliu *CheckListItemUpdate) SetNillableStringVal(s *string) *CheckListItemUpdate {
	if s != nil {
		cliu.SetStringVal(*s)
	}
	return cliu
}

// ClearStringVal clears the value of string_val.
func (cliu *CheckListItemUpdate) ClearStringVal() *CheckListItemUpdate {
	cliu.mutation.ClearStringVal()
	return cliu
}

// SetEnumValues sets the enum_values field.
func (cliu *CheckListItemUpdate) SetEnumValues(s string) *CheckListItemUpdate {
	cliu.mutation.SetEnumValues(s)
	return cliu
}

// SetNillableEnumValues sets the enum_values field if the given value is not nil.
func (cliu *CheckListItemUpdate) SetNillableEnumValues(s *string) *CheckListItemUpdate {
	if s != nil {
		cliu.SetEnumValues(*s)
	}
	return cliu
}

// ClearEnumValues clears the value of enum_values.
func (cliu *CheckListItemUpdate) ClearEnumValues() *CheckListItemUpdate {
	cliu.mutation.ClearEnumValues()
	return cliu
}

// SetHelpText sets the help_text field.
func (cliu *CheckListItemUpdate) SetHelpText(s string) *CheckListItemUpdate {
	cliu.mutation.SetHelpText(s)
	return cliu
}

// SetNillableHelpText sets the help_text field if the given value is not nil.
func (cliu *CheckListItemUpdate) SetNillableHelpText(s *string) *CheckListItemUpdate {
	if s != nil {
		cliu.SetHelpText(*s)
	}
	return cliu
}

// ClearHelpText clears the value of help_text.
func (cliu *CheckListItemUpdate) ClearHelpText() *CheckListItemUpdate {
	cliu.mutation.ClearHelpText()
	return cliu
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (cliu *CheckListItemUpdate) SetWorkOrderID(id int) *CheckListItemUpdate {
	cliu.mutation.SetWorkOrderID(id)
	return cliu
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (cliu *CheckListItemUpdate) SetNillableWorkOrderID(id *int) *CheckListItemUpdate {
	if id != nil {
		cliu = cliu.SetWorkOrderID(*id)
	}
	return cliu
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (cliu *CheckListItemUpdate) SetWorkOrder(w *WorkOrder) *CheckListItemUpdate {
	return cliu.SetWorkOrderID(w.ID)
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (cliu *CheckListItemUpdate) ClearWorkOrder() *CheckListItemUpdate {
	cliu.mutation.ClearWorkOrder()
	return cliu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (cliu *CheckListItemUpdate) Save(ctx context.Context) (int, error) {

	var (
		err      error
		affected int
	)
	if len(cliu.hooks) == 0 {
		affected, err = cliu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CheckListItemMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			cliu.mutation = mutation
			affected, err = cliu.sqlSave(ctx)
			return affected, err
		})
		for i := len(cliu.hooks); i > 0; i-- {
			mut = cliu.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, cliu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (cliu *CheckListItemUpdate) SaveX(ctx context.Context) int {
	affected, err := cliu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cliu *CheckListItemUpdate) Exec(ctx context.Context) error {
	_, err := cliu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cliu *CheckListItemUpdate) ExecX(ctx context.Context) {
	if err := cliu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cliu *CheckListItemUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   checklistitem.Table,
			Columns: checklistitem.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistitem.FieldID,
			},
		},
	}
	if ps := cliu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cliu.mutation.Title(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitem.FieldTitle,
		})
	}
	if value, ok := cliu.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitem.FieldType,
		})
	}
	if value, ok := cliu.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: checklistitem.FieldIndex,
		})
	}
	if value, ok := cliu.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: checklistitem.FieldIndex,
		})
	}
	if cliu.mutation.IndexCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: checklistitem.FieldIndex,
		})
	}
	if value, ok := cliu.mutation.Checked(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: checklistitem.FieldChecked,
		})
	}
	if cliu.mutation.CheckedCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: checklistitem.FieldChecked,
		})
	}
	if value, ok := cliu.mutation.StringVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitem.FieldStringVal,
		})
	}
	if cliu.mutation.StringValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistitem.FieldStringVal,
		})
	}
	if value, ok := cliu.mutation.EnumValues(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitem.FieldEnumValues,
		})
	}
	if cliu.mutation.EnumValuesCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistitem.FieldEnumValues,
		})
	}
	if value, ok := cliu.mutation.HelpText(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitem.FieldHelpText,
		})
	}
	if cliu.mutation.HelpTextCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistitem.FieldHelpText,
		})
	}
	if cliu.mutation.WorkOrderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistitem.WorkOrderTable,
			Columns: []string{checklistitem.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cliu.mutation.WorkOrderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistitem.WorkOrderTable,
			Columns: []string{checklistitem.WorkOrderColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, cliu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{checklistitem.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// CheckListItemUpdateOne is the builder for updating a single CheckListItem entity.
type CheckListItemUpdateOne struct {
	config
	hooks    []Hook
	mutation *CheckListItemMutation
}

// SetTitle sets the title field.
func (cliuo *CheckListItemUpdateOne) SetTitle(s string) *CheckListItemUpdateOne {
	cliuo.mutation.SetTitle(s)
	return cliuo
}

// SetType sets the type field.
func (cliuo *CheckListItemUpdateOne) SetType(s string) *CheckListItemUpdateOne {
	cliuo.mutation.SetType(s)
	return cliuo
}

// SetIndex sets the index field.
func (cliuo *CheckListItemUpdateOne) SetIndex(i int) *CheckListItemUpdateOne {
	cliuo.mutation.ResetIndex()
	cliuo.mutation.SetIndex(i)
	return cliuo
}

// SetNillableIndex sets the index field if the given value is not nil.
func (cliuo *CheckListItemUpdateOne) SetNillableIndex(i *int) *CheckListItemUpdateOne {
	if i != nil {
		cliuo.SetIndex(*i)
	}
	return cliuo
}

// AddIndex adds i to index.
func (cliuo *CheckListItemUpdateOne) AddIndex(i int) *CheckListItemUpdateOne {
	cliuo.mutation.AddIndex(i)
	return cliuo
}

// ClearIndex clears the value of index.
func (cliuo *CheckListItemUpdateOne) ClearIndex() *CheckListItemUpdateOne {
	cliuo.mutation.ClearIndex()
	return cliuo
}

// SetChecked sets the checked field.
func (cliuo *CheckListItemUpdateOne) SetChecked(b bool) *CheckListItemUpdateOne {
	cliuo.mutation.SetChecked(b)
	return cliuo
}

// SetNillableChecked sets the checked field if the given value is not nil.
func (cliuo *CheckListItemUpdateOne) SetNillableChecked(b *bool) *CheckListItemUpdateOne {
	if b != nil {
		cliuo.SetChecked(*b)
	}
	return cliuo
}

// ClearChecked clears the value of checked.
func (cliuo *CheckListItemUpdateOne) ClearChecked() *CheckListItemUpdateOne {
	cliuo.mutation.ClearChecked()
	return cliuo
}

// SetStringVal sets the string_val field.
func (cliuo *CheckListItemUpdateOne) SetStringVal(s string) *CheckListItemUpdateOne {
	cliuo.mutation.SetStringVal(s)
	return cliuo
}

// SetNillableStringVal sets the string_val field if the given value is not nil.
func (cliuo *CheckListItemUpdateOne) SetNillableStringVal(s *string) *CheckListItemUpdateOne {
	if s != nil {
		cliuo.SetStringVal(*s)
	}
	return cliuo
}

// ClearStringVal clears the value of string_val.
func (cliuo *CheckListItemUpdateOne) ClearStringVal() *CheckListItemUpdateOne {
	cliuo.mutation.ClearStringVal()
	return cliuo
}

// SetEnumValues sets the enum_values field.
func (cliuo *CheckListItemUpdateOne) SetEnumValues(s string) *CheckListItemUpdateOne {
	cliuo.mutation.SetEnumValues(s)
	return cliuo
}

// SetNillableEnumValues sets the enum_values field if the given value is not nil.
func (cliuo *CheckListItemUpdateOne) SetNillableEnumValues(s *string) *CheckListItemUpdateOne {
	if s != nil {
		cliuo.SetEnumValues(*s)
	}
	return cliuo
}

// ClearEnumValues clears the value of enum_values.
func (cliuo *CheckListItemUpdateOne) ClearEnumValues() *CheckListItemUpdateOne {
	cliuo.mutation.ClearEnumValues()
	return cliuo
}

// SetHelpText sets the help_text field.
func (cliuo *CheckListItemUpdateOne) SetHelpText(s string) *CheckListItemUpdateOne {
	cliuo.mutation.SetHelpText(s)
	return cliuo
}

// SetNillableHelpText sets the help_text field if the given value is not nil.
func (cliuo *CheckListItemUpdateOne) SetNillableHelpText(s *string) *CheckListItemUpdateOne {
	if s != nil {
		cliuo.SetHelpText(*s)
	}
	return cliuo
}

// ClearHelpText clears the value of help_text.
func (cliuo *CheckListItemUpdateOne) ClearHelpText() *CheckListItemUpdateOne {
	cliuo.mutation.ClearHelpText()
	return cliuo
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (cliuo *CheckListItemUpdateOne) SetWorkOrderID(id int) *CheckListItemUpdateOne {
	cliuo.mutation.SetWorkOrderID(id)
	return cliuo
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (cliuo *CheckListItemUpdateOne) SetNillableWorkOrderID(id *int) *CheckListItemUpdateOne {
	if id != nil {
		cliuo = cliuo.SetWorkOrderID(*id)
	}
	return cliuo
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (cliuo *CheckListItemUpdateOne) SetWorkOrder(w *WorkOrder) *CheckListItemUpdateOne {
	return cliuo.SetWorkOrderID(w.ID)
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (cliuo *CheckListItemUpdateOne) ClearWorkOrder() *CheckListItemUpdateOne {
	cliuo.mutation.ClearWorkOrder()
	return cliuo
}

// Save executes the query and returns the updated entity.
func (cliuo *CheckListItemUpdateOne) Save(ctx context.Context) (*CheckListItem, error) {

	var (
		err  error
		node *CheckListItem
	)
	if len(cliuo.hooks) == 0 {
		node, err = cliuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CheckListItemMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			cliuo.mutation = mutation
			node, err = cliuo.sqlSave(ctx)
			return node, err
		})
		for i := len(cliuo.hooks); i > 0; i-- {
			mut = cliuo.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, cliuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (cliuo *CheckListItemUpdateOne) SaveX(ctx context.Context) *CheckListItem {
	cli, err := cliuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return cli
}

// Exec executes the query on the entity.
func (cliuo *CheckListItemUpdateOne) Exec(ctx context.Context) error {
	_, err := cliuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cliuo *CheckListItemUpdateOne) ExecX(ctx context.Context) {
	if err := cliuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cliuo *CheckListItemUpdateOne) sqlSave(ctx context.Context) (cli *CheckListItem, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   checklistitem.Table,
			Columns: checklistitem.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistitem.FieldID,
			},
		},
	}
	id, ok := cliuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing CheckListItem.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := cliuo.mutation.Title(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitem.FieldTitle,
		})
	}
	if value, ok := cliuo.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitem.FieldType,
		})
	}
	if value, ok := cliuo.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: checklistitem.FieldIndex,
		})
	}
	if value, ok := cliuo.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: checklistitem.FieldIndex,
		})
	}
	if cliuo.mutation.IndexCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: checklistitem.FieldIndex,
		})
	}
	if value, ok := cliuo.mutation.Checked(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: checklistitem.FieldChecked,
		})
	}
	if cliuo.mutation.CheckedCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: checklistitem.FieldChecked,
		})
	}
	if value, ok := cliuo.mutation.StringVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitem.FieldStringVal,
		})
	}
	if cliuo.mutation.StringValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistitem.FieldStringVal,
		})
	}
	if value, ok := cliuo.mutation.EnumValues(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitem.FieldEnumValues,
		})
	}
	if cliuo.mutation.EnumValuesCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistitem.FieldEnumValues,
		})
	}
	if value, ok := cliuo.mutation.HelpText(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitem.FieldHelpText,
		})
	}
	if cliuo.mutation.HelpTextCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistitem.FieldHelpText,
		})
	}
	if cliuo.mutation.WorkOrderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistitem.WorkOrderTable,
			Columns: []string{checklistitem.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cliuo.mutation.WorkOrderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistitem.WorkOrderTable,
			Columns: []string{checklistitem.WorkOrderColumn},
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
	cli = &CheckListItem{config: cliuo.config}
	_spec.Assign = cli.assignValues
	_spec.ScanValues = cli.scanValues()
	if err = sqlgraph.UpdateNode(ctx, cliuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{checklistitem.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return cli, nil
}
