// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/checklistcategory"
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
)

// CheckListCategoryCreate is the builder for creating a CheckListCategory entity.
type CheckListCategoryCreate struct {
	config
	create_time      *time.Time
	update_time      *time.Time
	title            *string
	description      *string
	check_list_items map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (clcc *CheckListCategoryCreate) SetCreateTime(t time.Time) *CheckListCategoryCreate {
	clcc.create_time = &t
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
	clcc.update_time = &t
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
	clcc.title = &s
	return clcc
}

// SetDescription sets the description field.
func (clcc *CheckListCategoryCreate) SetDescription(s string) *CheckListCategoryCreate {
	clcc.description = &s
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
func (clcc *CheckListCategoryCreate) AddCheckListItemIDs(ids ...string) *CheckListCategoryCreate {
	if clcc.check_list_items == nil {
		clcc.check_list_items = make(map[string]struct{})
	}
	for i := range ids {
		clcc.check_list_items[ids[i]] = struct{}{}
	}
	return clcc
}

// AddCheckListItems adds the check_list_items edges to CheckListItem.
func (clcc *CheckListCategoryCreate) AddCheckListItems(c ...*CheckListItem) *CheckListCategoryCreate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return clcc.AddCheckListItemIDs(ids...)
}

// Save creates the CheckListCategory in the database.
func (clcc *CheckListCategoryCreate) Save(ctx context.Context) (*CheckListCategory, error) {
	if clcc.create_time == nil {
		v := checklistcategory.DefaultCreateTime()
		clcc.create_time = &v
	}
	if clcc.update_time == nil {
		v := checklistcategory.DefaultUpdateTime()
		clcc.update_time = &v
	}
	if clcc.title == nil {
		return nil, errors.New("ent: missing required field \"title\"")
	}
	return clcc.sqlSave(ctx)
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
				Type:   field.TypeString,
				Column: checklistcategory.FieldID,
			},
		}
	)
	if value := clcc.create_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: checklistcategory.FieldCreateTime,
		})
		clc.CreateTime = *value
	}
	if value := clcc.update_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: checklistcategory.FieldUpdateTime,
		})
		clc.UpdateTime = *value
	}
	if value := clcc.title; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: checklistcategory.FieldTitle,
		})
		clc.Title = *value
	}
	if value := clcc.description; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: checklistcategory.FieldDescription,
		})
		clc.Description = *value
	}
	if nodes := clcc.check_list_items; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   checklistcategory.CheckListItemsTable,
			Columns: []string{checklistcategory.CheckListItemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: checklistitem.FieldID,
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, clcc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	clc.ID = strconv.FormatInt(id, 10)
	return clc, nil
}
