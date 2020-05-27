// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/checklistcategorydefinition"
	"github.com/facebookincubator/symphony/graph/ent/checklistitemdefinition"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// CheckListItemDefinitionUpdate is the builder for updating CheckListItemDefinition entities.
type CheckListItemDefinitionUpdate struct {
	config
	hooks      []Hook
	mutation   *CheckListItemDefinitionMutation
	predicates []predicate.CheckListItemDefinition
}

// Where adds a new predicate for the builder.
func (clidu *CheckListItemDefinitionUpdate) Where(ps ...predicate.CheckListItemDefinition) *CheckListItemDefinitionUpdate {
	clidu.predicates = append(clidu.predicates, ps...)
	return clidu
}

// SetTitle sets the title field.
func (clidu *CheckListItemDefinitionUpdate) SetTitle(s string) *CheckListItemDefinitionUpdate {
	clidu.mutation.SetTitle(s)
	return clidu
}

// SetType sets the type field.
func (clidu *CheckListItemDefinitionUpdate) SetType(s string) *CheckListItemDefinitionUpdate {
	clidu.mutation.SetType(s)
	return clidu
}

// SetIndex sets the index field.
func (clidu *CheckListItemDefinitionUpdate) SetIndex(i int) *CheckListItemDefinitionUpdate {
	clidu.mutation.ResetIndex()
	clidu.mutation.SetIndex(i)
	return clidu
}

// SetNillableIndex sets the index field if the given value is not nil.
func (clidu *CheckListItemDefinitionUpdate) SetNillableIndex(i *int) *CheckListItemDefinitionUpdate {
	if i != nil {
		clidu.SetIndex(*i)
	}
	return clidu
}

// AddIndex adds i to index.
func (clidu *CheckListItemDefinitionUpdate) AddIndex(i int) *CheckListItemDefinitionUpdate {
	clidu.mutation.AddIndex(i)
	return clidu
}

// ClearIndex clears the value of index.
func (clidu *CheckListItemDefinitionUpdate) ClearIndex() *CheckListItemDefinitionUpdate {
	clidu.mutation.ClearIndex()
	return clidu
}

// SetEnumValues sets the enum_values field.
func (clidu *CheckListItemDefinitionUpdate) SetEnumValues(s string) *CheckListItemDefinitionUpdate {
	clidu.mutation.SetEnumValues(s)
	return clidu
}

// SetNillableEnumValues sets the enum_values field if the given value is not nil.
func (clidu *CheckListItemDefinitionUpdate) SetNillableEnumValues(s *string) *CheckListItemDefinitionUpdate {
	if s != nil {
		clidu.SetEnumValues(*s)
	}
	return clidu
}

// ClearEnumValues clears the value of enum_values.
func (clidu *CheckListItemDefinitionUpdate) ClearEnumValues() *CheckListItemDefinitionUpdate {
	clidu.mutation.ClearEnumValues()
	return clidu
}

// SetEnumSelectionModeValue sets the enum_selection_mode_value field.
func (clidu *CheckListItemDefinitionUpdate) SetEnumSelectionModeValue(csmv checklistitemdefinition.EnumSelectionModeValue) *CheckListItemDefinitionUpdate {
	clidu.mutation.SetEnumSelectionModeValue(csmv)
	return clidu
}

// SetNillableEnumSelectionModeValue sets the enum_selection_mode_value field if the given value is not nil.
func (clidu *CheckListItemDefinitionUpdate) SetNillableEnumSelectionModeValue(csmv *checklistitemdefinition.EnumSelectionModeValue) *CheckListItemDefinitionUpdate {
	if csmv != nil {
		clidu.SetEnumSelectionModeValue(*csmv)
	}
	return clidu
}

// ClearEnumSelectionModeValue clears the value of enum_selection_mode_value.
func (clidu *CheckListItemDefinitionUpdate) ClearEnumSelectionModeValue() *CheckListItemDefinitionUpdate {
	clidu.mutation.ClearEnumSelectionModeValue()
	return clidu
}

// SetHelpText sets the help_text field.
func (clidu *CheckListItemDefinitionUpdate) SetHelpText(s string) *CheckListItemDefinitionUpdate {
	clidu.mutation.SetHelpText(s)
	return clidu
}

// SetNillableHelpText sets the help_text field if the given value is not nil.
func (clidu *CheckListItemDefinitionUpdate) SetNillableHelpText(s *string) *CheckListItemDefinitionUpdate {
	if s != nil {
		clidu.SetHelpText(*s)
	}
	return clidu
}

// ClearHelpText clears the value of help_text.
func (clidu *CheckListItemDefinitionUpdate) ClearHelpText() *CheckListItemDefinitionUpdate {
	clidu.mutation.ClearHelpText()
	return clidu
}

// SetCheckListCategoryDefinitionID sets the check_list_category_definition edge to CheckListCategoryDefinition by id.
func (clidu *CheckListItemDefinitionUpdate) SetCheckListCategoryDefinitionID(id int) *CheckListItemDefinitionUpdate {
	clidu.mutation.SetCheckListCategoryDefinitionID(id)
	return clidu
}

// SetCheckListCategoryDefinition sets the check_list_category_definition edge to CheckListCategoryDefinition.
func (clidu *CheckListItemDefinitionUpdate) SetCheckListCategoryDefinition(c *CheckListCategoryDefinition) *CheckListItemDefinitionUpdate {
	return clidu.SetCheckListCategoryDefinitionID(c.ID)
}

// ClearCheckListCategoryDefinition clears the check_list_category_definition edge to CheckListCategoryDefinition.
func (clidu *CheckListItemDefinitionUpdate) ClearCheckListCategoryDefinition() *CheckListItemDefinitionUpdate {
	clidu.mutation.ClearCheckListCategoryDefinition()
	return clidu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (clidu *CheckListItemDefinitionUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := clidu.mutation.UpdateTime(); !ok {
		v := checklistitemdefinition.UpdateDefaultUpdateTime()
		clidu.mutation.SetUpdateTime(v)
	}
	if v, ok := clidu.mutation.EnumSelectionModeValue(); ok {
		if err := checklistitemdefinition.EnumSelectionModeValueValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"enum_selection_mode_value\": %v", err)
		}
	}

	if _, ok := clidu.mutation.CheckListCategoryDefinitionID(); clidu.mutation.CheckListCategoryDefinitionCleared() && !ok {
		return 0, errors.New("ent: clearing a unique edge \"check_list_category_definition\"")
	}
	var (
		err      error
		affected int
	)
	if len(clidu.hooks) == 0 {
		affected, err = clidu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CheckListItemDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			clidu.mutation = mutation
			affected, err = clidu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(clidu.hooks) - 1; i >= 0; i-- {
			mut = clidu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, clidu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (clidu *CheckListItemDefinitionUpdate) SaveX(ctx context.Context) int {
	affected, err := clidu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (clidu *CheckListItemDefinitionUpdate) Exec(ctx context.Context) error {
	_, err := clidu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (clidu *CheckListItemDefinitionUpdate) ExecX(ctx context.Context) {
	if err := clidu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (clidu *CheckListItemDefinitionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   checklistitemdefinition.Table,
			Columns: checklistitemdefinition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistitemdefinition.FieldID,
			},
		},
	}
	if ps := clidu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := clidu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: checklistitemdefinition.FieldUpdateTime,
		})
	}
	if value, ok := clidu.mutation.Title(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitemdefinition.FieldTitle,
		})
	}
	if value, ok := clidu.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitemdefinition.FieldType,
		})
	}
	if value, ok := clidu.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: checklistitemdefinition.FieldIndex,
		})
	}
	if value, ok := clidu.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: checklistitemdefinition.FieldIndex,
		})
	}
	if clidu.mutation.IndexCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: checklistitemdefinition.FieldIndex,
		})
	}
	if value, ok := clidu.mutation.EnumValues(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitemdefinition.FieldEnumValues,
		})
	}
	if clidu.mutation.EnumValuesCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistitemdefinition.FieldEnumValues,
		})
	}
	if value, ok := clidu.mutation.EnumSelectionModeValue(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: checklistitemdefinition.FieldEnumSelectionModeValue,
		})
	}
	if clidu.mutation.EnumSelectionModeValueCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Column: checklistitemdefinition.FieldEnumSelectionModeValue,
		})
	}
	if value, ok := clidu.mutation.HelpText(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitemdefinition.FieldHelpText,
		})
	}
	if clidu.mutation.HelpTextCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistitemdefinition.FieldHelpText,
		})
	}
	if clidu.mutation.CheckListCategoryDefinitionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistitemdefinition.CheckListCategoryDefinitionTable,
			Columns: []string{checklistitemdefinition.CheckListCategoryDefinitionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistcategorydefinition.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := clidu.mutation.CheckListCategoryDefinitionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistitemdefinition.CheckListCategoryDefinitionTable,
			Columns: []string{checklistitemdefinition.CheckListCategoryDefinitionColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, clidu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{checklistitemdefinition.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// CheckListItemDefinitionUpdateOne is the builder for updating a single CheckListItemDefinition entity.
type CheckListItemDefinitionUpdateOne struct {
	config
	hooks    []Hook
	mutation *CheckListItemDefinitionMutation
}

// SetTitle sets the title field.
func (cliduo *CheckListItemDefinitionUpdateOne) SetTitle(s string) *CheckListItemDefinitionUpdateOne {
	cliduo.mutation.SetTitle(s)
	return cliduo
}

// SetType sets the type field.
func (cliduo *CheckListItemDefinitionUpdateOne) SetType(s string) *CheckListItemDefinitionUpdateOne {
	cliduo.mutation.SetType(s)
	return cliduo
}

// SetIndex sets the index field.
func (cliduo *CheckListItemDefinitionUpdateOne) SetIndex(i int) *CheckListItemDefinitionUpdateOne {
	cliduo.mutation.ResetIndex()
	cliduo.mutation.SetIndex(i)
	return cliduo
}

// SetNillableIndex sets the index field if the given value is not nil.
func (cliduo *CheckListItemDefinitionUpdateOne) SetNillableIndex(i *int) *CheckListItemDefinitionUpdateOne {
	if i != nil {
		cliduo.SetIndex(*i)
	}
	return cliduo
}

// AddIndex adds i to index.
func (cliduo *CheckListItemDefinitionUpdateOne) AddIndex(i int) *CheckListItemDefinitionUpdateOne {
	cliduo.mutation.AddIndex(i)
	return cliduo
}

// ClearIndex clears the value of index.
func (cliduo *CheckListItemDefinitionUpdateOne) ClearIndex() *CheckListItemDefinitionUpdateOne {
	cliduo.mutation.ClearIndex()
	return cliduo
}

// SetEnumValues sets the enum_values field.
func (cliduo *CheckListItemDefinitionUpdateOne) SetEnumValues(s string) *CheckListItemDefinitionUpdateOne {
	cliduo.mutation.SetEnumValues(s)
	return cliduo
}

// SetNillableEnumValues sets the enum_values field if the given value is not nil.
func (cliduo *CheckListItemDefinitionUpdateOne) SetNillableEnumValues(s *string) *CheckListItemDefinitionUpdateOne {
	if s != nil {
		cliduo.SetEnumValues(*s)
	}
	return cliduo
}

// ClearEnumValues clears the value of enum_values.
func (cliduo *CheckListItemDefinitionUpdateOne) ClearEnumValues() *CheckListItemDefinitionUpdateOne {
	cliduo.mutation.ClearEnumValues()
	return cliduo
}

// SetEnumSelectionModeValue sets the enum_selection_mode_value field.
func (cliduo *CheckListItemDefinitionUpdateOne) SetEnumSelectionModeValue(csmv checklistitemdefinition.EnumSelectionModeValue) *CheckListItemDefinitionUpdateOne {
	cliduo.mutation.SetEnumSelectionModeValue(csmv)
	return cliduo
}

// SetNillableEnumSelectionModeValue sets the enum_selection_mode_value field if the given value is not nil.
func (cliduo *CheckListItemDefinitionUpdateOne) SetNillableEnumSelectionModeValue(csmv *checklistitemdefinition.EnumSelectionModeValue) *CheckListItemDefinitionUpdateOne {
	if csmv != nil {
		cliduo.SetEnumSelectionModeValue(*csmv)
	}
	return cliduo
}

// ClearEnumSelectionModeValue clears the value of enum_selection_mode_value.
func (cliduo *CheckListItemDefinitionUpdateOne) ClearEnumSelectionModeValue() *CheckListItemDefinitionUpdateOne {
	cliduo.mutation.ClearEnumSelectionModeValue()
	return cliduo
}

// SetHelpText sets the help_text field.
func (cliduo *CheckListItemDefinitionUpdateOne) SetHelpText(s string) *CheckListItemDefinitionUpdateOne {
	cliduo.mutation.SetHelpText(s)
	return cliduo
}

// SetNillableHelpText sets the help_text field if the given value is not nil.
func (cliduo *CheckListItemDefinitionUpdateOne) SetNillableHelpText(s *string) *CheckListItemDefinitionUpdateOne {
	if s != nil {
		cliduo.SetHelpText(*s)
	}
	return cliduo
}

// ClearHelpText clears the value of help_text.
func (cliduo *CheckListItemDefinitionUpdateOne) ClearHelpText() *CheckListItemDefinitionUpdateOne {
	cliduo.mutation.ClearHelpText()
	return cliduo
}

// SetCheckListCategoryDefinitionID sets the check_list_category_definition edge to CheckListCategoryDefinition by id.
func (cliduo *CheckListItemDefinitionUpdateOne) SetCheckListCategoryDefinitionID(id int) *CheckListItemDefinitionUpdateOne {
	cliduo.mutation.SetCheckListCategoryDefinitionID(id)
	return cliduo
}

// SetCheckListCategoryDefinition sets the check_list_category_definition edge to CheckListCategoryDefinition.
func (cliduo *CheckListItemDefinitionUpdateOne) SetCheckListCategoryDefinition(c *CheckListCategoryDefinition) *CheckListItemDefinitionUpdateOne {
	return cliduo.SetCheckListCategoryDefinitionID(c.ID)
}

// ClearCheckListCategoryDefinition clears the check_list_category_definition edge to CheckListCategoryDefinition.
func (cliduo *CheckListItemDefinitionUpdateOne) ClearCheckListCategoryDefinition() *CheckListItemDefinitionUpdateOne {
	cliduo.mutation.ClearCheckListCategoryDefinition()
	return cliduo
}

// Save executes the query and returns the updated entity.
func (cliduo *CheckListItemDefinitionUpdateOne) Save(ctx context.Context) (*CheckListItemDefinition, error) {
	if _, ok := cliduo.mutation.UpdateTime(); !ok {
		v := checklistitemdefinition.UpdateDefaultUpdateTime()
		cliduo.mutation.SetUpdateTime(v)
	}
	if v, ok := cliduo.mutation.EnumSelectionModeValue(); ok {
		if err := checklistitemdefinition.EnumSelectionModeValueValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"enum_selection_mode_value\": %v", err)
		}
	}

	if _, ok := cliduo.mutation.CheckListCategoryDefinitionID(); cliduo.mutation.CheckListCategoryDefinitionCleared() && !ok {
		return nil, errors.New("ent: clearing a unique edge \"check_list_category_definition\"")
	}
	var (
		err  error
		node *CheckListItemDefinition
	)
	if len(cliduo.hooks) == 0 {
		node, err = cliduo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CheckListItemDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			cliduo.mutation = mutation
			node, err = cliduo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(cliduo.hooks) - 1; i >= 0; i-- {
			mut = cliduo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, cliduo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (cliduo *CheckListItemDefinitionUpdateOne) SaveX(ctx context.Context) *CheckListItemDefinition {
	clid, err := cliduo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return clid
}

// Exec executes the query on the entity.
func (cliduo *CheckListItemDefinitionUpdateOne) Exec(ctx context.Context) error {
	_, err := cliduo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cliduo *CheckListItemDefinitionUpdateOne) ExecX(ctx context.Context) {
	if err := cliduo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cliduo *CheckListItemDefinitionUpdateOne) sqlSave(ctx context.Context) (clid *CheckListItemDefinition, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   checklistitemdefinition.Table,
			Columns: checklistitemdefinition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistitemdefinition.FieldID,
			},
		},
	}
	id, ok := cliduo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing CheckListItemDefinition.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := cliduo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: checklistitemdefinition.FieldUpdateTime,
		})
	}
	if value, ok := cliduo.mutation.Title(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitemdefinition.FieldTitle,
		})
	}
	if value, ok := cliduo.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitemdefinition.FieldType,
		})
	}
	if value, ok := cliduo.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: checklistitemdefinition.FieldIndex,
		})
	}
	if value, ok := cliduo.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: checklistitemdefinition.FieldIndex,
		})
	}
	if cliduo.mutation.IndexCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: checklistitemdefinition.FieldIndex,
		})
	}
	if value, ok := cliduo.mutation.EnumValues(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitemdefinition.FieldEnumValues,
		})
	}
	if cliduo.mutation.EnumValuesCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistitemdefinition.FieldEnumValues,
		})
	}
	if value, ok := cliduo.mutation.EnumSelectionModeValue(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: checklistitemdefinition.FieldEnumSelectionModeValue,
		})
	}
	if cliduo.mutation.EnumSelectionModeValueCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Column: checklistitemdefinition.FieldEnumSelectionModeValue,
		})
	}
	if value, ok := cliduo.mutation.HelpText(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistitemdefinition.FieldHelpText,
		})
	}
	if cliduo.mutation.HelpTextCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistitemdefinition.FieldHelpText,
		})
	}
	if cliduo.mutation.CheckListCategoryDefinitionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistitemdefinition.CheckListCategoryDefinitionTable,
			Columns: []string{checklistitemdefinition.CheckListCategoryDefinitionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistcategorydefinition.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cliduo.mutation.CheckListCategoryDefinitionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistitemdefinition.CheckListCategoryDefinitionTable,
			Columns: []string{checklistitemdefinition.CheckListCategoryDefinitionColumn},
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
	clid = &CheckListItemDefinition{config: cliduo.config}
	_spec.Assign = clid.assignValues
	_spec.ScanValues = clid.scanValues()
	if err = sqlgraph.UpdateNode(ctx, cliduo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{checklistitemdefinition.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return clid, nil
}
