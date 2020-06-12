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
	"github.com/facebookincubator/symphony/pkg/ent/checklistcategory"
	"github.com/facebookincubator/symphony/pkg/ent/checklistitem"
	"github.com/facebookincubator/symphony/pkg/ent/workorder"
)

// CheckListCategoryCreate is the builder for creating a CheckListCategory entity.
type CheckListCategoryCreate struct {
	config
	mutation *CheckListCategoryMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (clcc *CheckListCategoryCreate) SetCreateTime(t time.Time) *CheckListCategoryCreate {
	clcc.mutation.SetCreateTime(t)
	return clcc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (clcc *CheckListCategoryCreate) SetNillableCreateTime(t *time.Time) *CheckListCategoryCreate {
	if t != nil {
		clcc.SetCreateTime(*t)
	}
	return clcc
}

// SetUpdateTime sets the update_time field.
func (clcc *CheckListCategoryCreate) SetUpdateTime(t time.Time) *CheckListCategoryCreate {
	clcc.mutation.SetUpdateTime(t)
	return clcc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (clcc *CheckListCategoryCreate) SetNillableUpdateTime(t *time.Time) *CheckListCategoryCreate {
	if t != nil {
		clcc.SetUpdateTime(*t)
	}
	return clcc
}

// SetTitle sets the title field.
func (clcc *CheckListCategoryCreate) SetTitle(s string) *CheckListCategoryCreate {
	clcc.mutation.SetTitle(s)
	return clcc
}

// SetDescription sets the description field.
func (clcc *CheckListCategoryCreate) SetDescription(s string) *CheckListCategoryCreate {
	clcc.mutation.SetDescription(s)
	return clcc
}

// SetNillableDescription sets the description field if the given value is not nil.
func (clcc *CheckListCategoryCreate) SetNillableDescription(s *string) *CheckListCategoryCreate {
	if s != nil {
		clcc.SetDescription(*s)
	}
	return clcc
}

// AddCheckListItemIDs adds the check_list_items edge to CheckListItem by ids.
func (clcc *CheckListCategoryCreate) AddCheckListItemIDs(ids ...int) *CheckListCategoryCreate {
	clcc.mutation.AddCheckListItemIDs(ids...)
	return clcc
}

// AddCheckListItems adds the check_list_items edges to CheckListItem.
func (clcc *CheckListCategoryCreate) AddCheckListItems(c ...*CheckListItem) *CheckListCategoryCreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return clcc.AddCheckListItemIDs(ids...)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (clcc *CheckListCategoryCreate) SetWorkOrderID(id int) *CheckListCategoryCreate {
	clcc.mutation.SetWorkOrderID(id)
	return clcc
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (clcc *CheckListCategoryCreate) SetWorkOrder(w *WorkOrder) *CheckListCategoryCreate {
	return clcc.SetWorkOrderID(w.ID)
}

// Save creates the CheckListCategory in the database.
func (clcc *CheckListCategoryCreate) Save(ctx context.Context) (*CheckListCategory, error) {
	if _, ok := clcc.mutation.CreateTime(); !ok {
		v := checklistcategory.DefaultCreateTime()
		clcc.mutation.SetCreateTime(v)
	}
	if _, ok := clcc.mutation.UpdateTime(); !ok {
		v := checklistcategory.DefaultUpdateTime()
		clcc.mutation.SetUpdateTime(v)
	}
	if _, ok := clcc.mutation.Title(); !ok {
		return nil, errors.New("ent: missing required field \"title\"")
	}
	if _, ok := clcc.mutation.WorkOrderID(); !ok {
		return nil, errors.New("ent: missing required edge \"work_order\"")
	}
	var (
		err  error
		node *CheckListCategory
	)
	if len(clcc.hooks) == 0 {
		node, err = clcc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CheckListCategoryMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			clcc.mutation = mutation
			node, err = clcc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(clcc.hooks) - 1; i >= 0; i-- {
			mut = clcc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, clcc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (clcc *CheckListCategoryCreate) SaveX(ctx context.Context) *CheckListCategory {
	v, err := clcc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (clcc *CheckListCategoryCreate) sqlSave(ctx context.Context) (*CheckListCategory, error) {
	var (
		clc   = &CheckListCategory{config: clcc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: checklistcategory.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistcategory.FieldID,
			},
		}
	)
	if value, ok := clcc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: checklistcategory.FieldCreateTime,
		})
		clc.CreateTime = value
	}
	if value, ok := clcc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: checklistcategory.FieldUpdateTime,
		})
		clc.UpdateTime = value
	}
	if value, ok := clcc.mutation.Title(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistcategory.FieldTitle,
		})
		clc.Title = value
	}
	if value, ok := clcc.mutation.Description(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistcategory.FieldDescription,
		})
		clc.Description = value
	}
	if nodes := clcc.mutation.CheckListItemsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   checklistcategory.CheckListItemsTable,
			Columns: []string{checklistcategory.CheckListItemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistitem.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := clcc.mutation.WorkOrderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistcategory.WorkOrderTable,
			Columns: []string{checklistcategory.WorkOrderColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, clcc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	clc.ID = int(id)
	return clc, nil
}
