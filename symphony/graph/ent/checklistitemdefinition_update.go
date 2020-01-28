// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"strconv"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/checklistitemdefinition"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// CheckListItemDefinitionUpdate is the builder for updating CheckListItemDefinition entities.
type CheckListItemDefinitionUpdate struct {
	config
	title                *string
	_type                *string
	index                *int
	addindex             *int
	clearindex           bool
	enum_values          *string
	clearenum_values     bool
	help_text            *string
	clearhelp_text       bool
	work_order_type      map[string]struct{}
	clearedWorkOrderType bool
	predicates           []predicate.CheckListItemDefinition
}

// Where adds a new predicate for the builder.
func (clidu *CheckListItemDefinitionUpdate) Where(ps ...predicate.CheckListItemDefinition) *CheckListItemDefinitionUpdate {
	clidu.predicates = append(clidu.predicates, ps...)
	return clidu
}

// SetTitle sets the title field.
func (clidu *CheckListItemDefinitionUpdate) SetTitle(s string) *CheckListItemDefinitionUpdate {
	clidu.title = &s
	return clidu
}

// SetType sets the type field.
func (clidu *CheckListItemDefinitionUpdate) SetType(s string) *CheckListItemDefinitionUpdate {
	clidu._type = &s
	return clidu
}

// SetIndex sets the index field.
func (clidu *CheckListItemDefinitionUpdate) SetIndex(i int) *CheckListItemDefinitionUpdate {
	clidu.index = &i
	clidu.addindex = nil
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
	if clidu.addindex == nil {
		clidu.addindex = &i
	} else {
		*clidu.addindex += i
	}
	return clidu
}

// ClearIndex clears the value of index.
func (clidu *CheckListItemDefinitionUpdate) ClearIndex() *CheckListItemDefinitionUpdate {
	clidu.index = nil
	clidu.clearindex = true
	return clidu
}

// SetEnumValues sets the enum_values field.
func (clidu *CheckListItemDefinitionUpdate) SetEnumValues(s string) *CheckListItemDefinitionUpdate {
	clidu.enum_values = &s
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
	clidu.enum_values = nil
	clidu.clearenum_values = true
	return clidu
}

// SetHelpText sets the help_text field.
func (clidu *CheckListItemDefinitionUpdate) SetHelpText(s string) *CheckListItemDefinitionUpdate {
	clidu.help_text = &s
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
	clidu.help_text = nil
	clidu.clearhelp_text = true
	return clidu
}

// SetWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id.
func (clidu *CheckListItemDefinitionUpdate) SetWorkOrderTypeID(id string) *CheckListItemDefinitionUpdate {
	if clidu.work_order_type == nil {
		clidu.work_order_type = make(map[string]struct{})
	}
	clidu.work_order_type[id] = struct{}{}
	return clidu
}

// SetNillableWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id if the given value is not nil.
func (clidu *CheckListItemDefinitionUpdate) SetNillableWorkOrderTypeID(id *string) *CheckListItemDefinitionUpdate {
	if id != nil {
		clidu = clidu.SetWorkOrderTypeID(*id)
	}
	return clidu
}

// SetWorkOrderType sets the work_order_type edge to WorkOrderType.
func (clidu *CheckListItemDefinitionUpdate) SetWorkOrderType(w *WorkOrderType) *CheckListItemDefinitionUpdate {
	return clidu.SetWorkOrderTypeID(w.ID)
}

// ClearWorkOrderType clears the work_order_type edge to WorkOrderType.
func (clidu *CheckListItemDefinitionUpdate) ClearWorkOrderType() *CheckListItemDefinitionUpdate {
	clidu.clearedWorkOrderType = true
	return clidu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (clidu *CheckListItemDefinitionUpdate) Save(ctx context.Context) (int, error) {
	if len(clidu.work_order_type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"work_order_type\"")
	}
	return clidu.sqlSave(ctx)
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
				Type:   field.TypeString,
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
	if value := clidu.title; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: checklistitemdefinition.FieldTitle,
		})
	}
	if value := clidu._type; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: checklistitemdefinition.FieldType,
		})
	}
	if value := clidu.index; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: checklistitemdefinition.FieldIndex,
		})
	}
	if value := clidu.addindex; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: checklistitemdefinition.FieldIndex,
		})
	}
	if clidu.clearindex {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: checklistitemdefinition.FieldIndex,
		})
	}
	if value := clidu.enum_values; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: checklistitemdefinition.FieldEnumValues,
		})
	}
	if clidu.clearenum_values {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistitemdefinition.FieldEnumValues,
		})
	}
	if value := clidu.help_text; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: checklistitemdefinition.FieldHelpText,
		})
	}
	if clidu.clearhelp_text {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistitemdefinition.FieldHelpText,
		})
	}
	if clidu.clearedWorkOrderType {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistitemdefinition.WorkOrderTypeTable,
			Columns: []string{checklistitemdefinition.WorkOrderTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workordertype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := clidu.work_order_type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistitemdefinition.WorkOrderTypeTable,
			Columns: []string{checklistitemdefinition.WorkOrderTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workordertype.FieldID,
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
	if n, err = sqlgraph.UpdateNodes(ctx, clidu.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// CheckListItemDefinitionUpdateOne is the builder for updating a single CheckListItemDefinition entity.
type CheckListItemDefinitionUpdateOne struct {
	config
	id                   string
	title                *string
	_type                *string
	index                *int
	addindex             *int
	clearindex           bool
	enum_values          *string
	clearenum_values     bool
	help_text            *string
	clearhelp_text       bool
	work_order_type      map[string]struct{}
	clearedWorkOrderType bool
}

// SetTitle sets the title field.
func (cliduo *CheckListItemDefinitionUpdateOne) SetTitle(s string) *CheckListItemDefinitionUpdateOne {
	cliduo.title = &s
	return cliduo
}

// SetType sets the type field.
func (cliduo *CheckListItemDefinitionUpdateOne) SetType(s string) *CheckListItemDefinitionUpdateOne {
	cliduo._type = &s
	return cliduo
}

// SetIndex sets the index field.
func (cliduo *CheckListItemDefinitionUpdateOne) SetIndex(i int) *CheckListItemDefinitionUpdateOne {
	cliduo.index = &i
	cliduo.addindex = nil
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
	if cliduo.addindex == nil {
		cliduo.addindex = &i
	} else {
		*cliduo.addindex += i
	}
	return cliduo
}

// ClearIndex clears the value of index.
func (cliduo *CheckListItemDefinitionUpdateOne) ClearIndex() *CheckListItemDefinitionUpdateOne {
	cliduo.index = nil
	cliduo.clearindex = true
	return cliduo
}

// SetEnumValues sets the enum_values field.
func (cliduo *CheckListItemDefinitionUpdateOne) SetEnumValues(s string) *CheckListItemDefinitionUpdateOne {
	cliduo.enum_values = &s
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
	cliduo.enum_values = nil
	cliduo.clearenum_values = true
	return cliduo
}

// SetHelpText sets the help_text field.
func (cliduo *CheckListItemDefinitionUpdateOne) SetHelpText(s string) *CheckListItemDefinitionUpdateOne {
	cliduo.help_text = &s
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
	cliduo.help_text = nil
	cliduo.clearhelp_text = true
	return cliduo
}

// SetWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id.
func (cliduo *CheckListItemDefinitionUpdateOne) SetWorkOrderTypeID(id string) *CheckListItemDefinitionUpdateOne {
	if cliduo.work_order_type == nil {
		cliduo.work_order_type = make(map[string]struct{})
	}
	cliduo.work_order_type[id] = struct{}{}
	return cliduo
}

// SetNillableWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id if the given value is not nil.
func (cliduo *CheckListItemDefinitionUpdateOne) SetNillableWorkOrderTypeID(id *string) *CheckListItemDefinitionUpdateOne {
	if id != nil {
		cliduo = cliduo.SetWorkOrderTypeID(*id)
	}
	return cliduo
}

// SetWorkOrderType sets the work_order_type edge to WorkOrderType.
func (cliduo *CheckListItemDefinitionUpdateOne) SetWorkOrderType(w *WorkOrderType) *CheckListItemDefinitionUpdateOne {
	return cliduo.SetWorkOrderTypeID(w.ID)
}

// ClearWorkOrderType clears the work_order_type edge to WorkOrderType.
func (cliduo *CheckListItemDefinitionUpdateOne) ClearWorkOrderType() *CheckListItemDefinitionUpdateOne {
	cliduo.clearedWorkOrderType = true
	return cliduo
}

// Save executes the query and returns the updated entity.
func (cliduo *CheckListItemDefinitionUpdateOne) Save(ctx context.Context) (*CheckListItemDefinition, error) {
	if len(cliduo.work_order_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"work_order_type\"")
	}
	return cliduo.sqlSave(ctx)
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
				Value:  cliduo.id,
				Type:   field.TypeString,
				Column: checklistitemdefinition.FieldID,
			},
		},
	}
	if value := cliduo.title; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: checklistitemdefinition.FieldTitle,
		})
	}
	if value := cliduo._type; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: checklistitemdefinition.FieldType,
		})
	}
	if value := cliduo.index; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: checklistitemdefinition.FieldIndex,
		})
	}
	if value := cliduo.addindex; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: checklistitemdefinition.FieldIndex,
		})
	}
	if cliduo.clearindex {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: checklistitemdefinition.FieldIndex,
		})
	}
	if value := cliduo.enum_values; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: checklistitemdefinition.FieldEnumValues,
		})
	}
	if cliduo.clearenum_values {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistitemdefinition.FieldEnumValues,
		})
	}
	if value := cliduo.help_text; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: checklistitemdefinition.FieldHelpText,
		})
	}
	if cliduo.clearhelp_text {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistitemdefinition.FieldHelpText,
		})
	}
	if cliduo.clearedWorkOrderType {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistitemdefinition.WorkOrderTypeTable,
			Columns: []string{checklistitemdefinition.WorkOrderTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workordertype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cliduo.work_order_type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistitemdefinition.WorkOrderTypeTable,
			Columns: []string{checklistitemdefinition.WorkOrderTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workordertype.FieldID,
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
	clid = &CheckListItemDefinition{config: cliduo.config}
	_spec.Assign = clid.assignValues
	_spec.ScanValues = clid.scanValues()
	if err = sqlgraph.UpdateNode(ctx, cliduo.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return clid, nil
}
