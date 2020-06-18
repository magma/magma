// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/ent/checklistcategorydefinition"
	"github.com/facebookincubator/symphony/pkg/ent/checklistitemdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/workordertype"
)

// CheckListCategoryDefinitionCreate is the builder for creating a CheckListCategoryDefinition entity.
type CheckListCategoryDefinitionCreate struct {
	config
	mutation *CheckListCategoryDefinitionMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (clcdc *CheckListCategoryDefinitionCreate) SetCreateTime(t time.Time) *CheckListCategoryDefinitionCreate {
	clcdc.mutation.SetCreateTime(t)
	return clcdc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (clcdc *CheckListCategoryDefinitionCreate) SetNillableCreateTime(t *time.Time) *CheckListCategoryDefinitionCreate {
	if t != nil {
		clcdc.SetCreateTime(*t)
	}
	return clcdc
}

// SetUpdateTime sets the update_time field.
func (clcdc *CheckListCategoryDefinitionCreate) SetUpdateTime(t time.Time) *CheckListCategoryDefinitionCreate {
	clcdc.mutation.SetUpdateTime(t)
	return clcdc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (clcdc *CheckListCategoryDefinitionCreate) SetNillableUpdateTime(t *time.Time) *CheckListCategoryDefinitionCreate {
	if t != nil {
		clcdc.SetUpdateTime(*t)
	}
	return clcdc
}

// SetTitle sets the title field.
func (clcdc *CheckListCategoryDefinitionCreate) SetTitle(s string) *CheckListCategoryDefinitionCreate {
	clcdc.mutation.SetTitle(s)
	return clcdc
}

// SetDescription sets the description field.
func (clcdc *CheckListCategoryDefinitionCreate) SetDescription(s string) *CheckListCategoryDefinitionCreate {
	clcdc.mutation.SetDescription(s)
	return clcdc
}

// SetNillableDescription sets the description field if the given value is not nil.
func (clcdc *CheckListCategoryDefinitionCreate) SetNillableDescription(s *string) *CheckListCategoryDefinitionCreate {
	if s != nil {
		clcdc.SetDescription(*s)
	}
	return clcdc
}

// AddCheckListItemDefinitionIDs adds the check_list_item_definitions edge to CheckListItemDefinition by ids.
func (clcdc *CheckListCategoryDefinitionCreate) AddCheckListItemDefinitionIDs(ids ...int) *CheckListCategoryDefinitionCreate {
	clcdc.mutation.AddCheckListItemDefinitionIDs(ids...)
	return clcdc
}

// AddCheckListItemDefinitions adds the check_list_item_definitions edges to CheckListItemDefinition.
func (clcdc *CheckListCategoryDefinitionCreate) AddCheckListItemDefinitions(c ...*CheckListItemDefinition) *CheckListCategoryDefinitionCreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return clcdc.AddCheckListItemDefinitionIDs(ids...)
}

// SetWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id.
func (clcdc *CheckListCategoryDefinitionCreate) SetWorkOrderTypeID(id int) *CheckListCategoryDefinitionCreate {
	clcdc.mutation.SetWorkOrderTypeID(id)
	return clcdc
}

// SetWorkOrderType sets the work_order_type edge to WorkOrderType.
func (clcdc *CheckListCategoryDefinitionCreate) SetWorkOrderType(w *WorkOrderType) *CheckListCategoryDefinitionCreate {
	return clcdc.SetWorkOrderTypeID(w.ID)
}

// Save creates the CheckListCategoryDefinition in the database.
func (clcdc *CheckListCategoryDefinitionCreate) Save(ctx context.Context) (*CheckListCategoryDefinition, error) {
	if _, ok := clcdc.mutation.CreateTime(); !ok {
		v := checklistcategorydefinition.DefaultCreateTime()
		clcdc.mutation.SetCreateTime(v)
	}
	if _, ok := clcdc.mutation.UpdateTime(); !ok {
		v := checklistcategorydefinition.DefaultUpdateTime()
		clcdc.mutation.SetUpdateTime(v)
	}
	if _, ok := clcdc.mutation.Title(); !ok {
		return nil, errors.New("ent: missing required field \"title\"")
	}
	if v, ok := clcdc.mutation.Title(); ok {
		if err := checklistcategorydefinition.TitleValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"title\": %v", err)
		}
	}
	if _, ok := clcdc.mutation.WorkOrderTypeID(); !ok {
		return nil, errors.New("ent: missing required edge \"work_order_type\"")
	}
	var (
		err  error
		node *CheckListCategoryDefinition
	)
	if len(clcdc.hooks) == 0 {
		node, err = clcdc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CheckListCategoryDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			clcdc.mutation = mutation
			node, err = clcdc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(clcdc.hooks) - 1; i >= 0; i-- {
			mut = clcdc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, clcdc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (clcdc *CheckListCategoryDefinitionCreate) SaveX(ctx context.Context) *CheckListCategoryDefinition {
	v, err := clcdc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (clcdc *CheckListCategoryDefinitionCreate) sqlSave(ctx context.Context) (*CheckListCategoryDefinition, error) {
	var (
		clcd  = &CheckListCategoryDefinition{config: clcdc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: checklistcategorydefinition.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistcategorydefinition.FieldID,
			},
		}
	)
	if value, ok := clcdc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: checklistcategorydefinition.FieldCreateTime,
		})
		clcd.CreateTime = value
	}
	if value, ok := clcdc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: checklistcategorydefinition.FieldUpdateTime,
		})
		clcd.UpdateTime = value
	}
	if value, ok := clcdc.mutation.Title(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistcategorydefinition.FieldTitle,
		})
		clcd.Title = value
	}
	if value, ok := clcdc.mutation.Description(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistcategorydefinition.FieldDescription,
		})
		clcd.Description = value
	}
	if nodes := clcdc.mutation.CheckListItemDefinitionsIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := clcdc.mutation.WorkOrderTypeIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, clcdc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	clcd.ID = int(id)
	return clcd, nil
}
