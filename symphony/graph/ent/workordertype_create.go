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
	"github.com/facebookincubator/symphony/graph/ent/checklistcategorydefinition"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// WorkOrderTypeCreate is the builder for creating a WorkOrderType entity.
type WorkOrderTypeCreate struct {
	config
	mutation *WorkOrderTypeMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (wotc *WorkOrderTypeCreate) SetCreateTime(t time.Time) *WorkOrderTypeCreate {
	wotc.mutation.SetCreateTime(t)
	return wotc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (wotc *WorkOrderTypeCreate) SetNillableCreateTime(t *time.Time) *WorkOrderTypeCreate {
	if t != nil {
		wotc.SetCreateTime(*t)
	}
	return wotc
}

// SetUpdateTime sets the update_time field.
func (wotc *WorkOrderTypeCreate) SetUpdateTime(t time.Time) *WorkOrderTypeCreate {
	wotc.mutation.SetUpdateTime(t)
	return wotc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (wotc *WorkOrderTypeCreate) SetNillableUpdateTime(t *time.Time) *WorkOrderTypeCreate {
	if t != nil {
		wotc.SetUpdateTime(*t)
	}
	return wotc
}

// SetName sets the name field.
func (wotc *WorkOrderTypeCreate) SetName(s string) *WorkOrderTypeCreate {
	wotc.mutation.SetName(s)
	return wotc
}

// SetDescription sets the description field.
func (wotc *WorkOrderTypeCreate) SetDescription(s string) *WorkOrderTypeCreate {
	wotc.mutation.SetDescription(s)
	return wotc
}

// SetNillableDescription sets the description field if the given value is not nil.
func (wotc *WorkOrderTypeCreate) SetNillableDescription(s *string) *WorkOrderTypeCreate {
	if s != nil {
		wotc.SetDescription(*s)
	}
	return wotc
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (wotc *WorkOrderTypeCreate) AddWorkOrderIDs(ids ...int) *WorkOrderTypeCreate {
	wotc.mutation.AddWorkOrderIDs(ids...)
	return wotc
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (wotc *WorkOrderTypeCreate) AddWorkOrders(w ...*WorkOrder) *WorkOrderTypeCreate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotc.AddWorkOrderIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (wotc *WorkOrderTypeCreate) AddPropertyTypeIDs(ids ...int) *WorkOrderTypeCreate {
	wotc.mutation.AddPropertyTypeIDs(ids...)
	return wotc
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (wotc *WorkOrderTypeCreate) AddPropertyTypes(p ...*PropertyType) *WorkOrderTypeCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wotc.AddPropertyTypeIDs(ids...)
}

// AddDefinitionIDs adds the definitions edge to WorkOrderDefinition by ids.
func (wotc *WorkOrderTypeCreate) AddDefinitionIDs(ids ...int) *WorkOrderTypeCreate {
	wotc.mutation.AddDefinitionIDs(ids...)
	return wotc
}

// AddDefinitions adds the definitions edges to WorkOrderDefinition.
func (wotc *WorkOrderTypeCreate) AddDefinitions(w ...*WorkOrderDefinition) *WorkOrderTypeCreate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotc.AddDefinitionIDs(ids...)
}

// AddCheckListCategoryDefinitionIDs adds the check_list_category_definitions edge to CheckListCategoryDefinition by ids.
func (wotc *WorkOrderTypeCreate) AddCheckListCategoryDefinitionIDs(ids ...int) *WorkOrderTypeCreate {
	wotc.mutation.AddCheckListCategoryDefinitionIDs(ids...)
	return wotc
}

// AddCheckListCategoryDefinitions adds the check_list_category_definitions edges to CheckListCategoryDefinition.
func (wotc *WorkOrderTypeCreate) AddCheckListCategoryDefinitions(c ...*CheckListCategoryDefinition) *WorkOrderTypeCreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotc.AddCheckListCategoryDefinitionIDs(ids...)
}

// Save creates the WorkOrderType in the database.
func (wotc *WorkOrderTypeCreate) Save(ctx context.Context) (*WorkOrderType, error) {
	if _, ok := wotc.mutation.CreateTime(); !ok {
		v := workordertype.DefaultCreateTime()
		wotc.mutation.SetCreateTime(v)
	}
	if _, ok := wotc.mutation.UpdateTime(); !ok {
		v := workordertype.DefaultUpdateTime()
		wotc.mutation.SetUpdateTime(v)
	}
	if _, ok := wotc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	var (
		err  error
		node *WorkOrderType
	)
	if len(wotc.hooks) == 0 {
		node, err = wotc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*WorkOrderTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			wotc.mutation = mutation
			node, err = wotc.sqlSave(ctx)
			return node, err
		})
		for i := len(wotc.hooks) - 1; i >= 0; i-- {
			mut = wotc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, wotc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (wotc *WorkOrderTypeCreate) SaveX(ctx context.Context) *WorkOrderType {
	v, err := wotc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (wotc *WorkOrderTypeCreate) sqlSave(ctx context.Context) (*WorkOrderType, error) {
	var (
		wot   = &WorkOrderType{config: wotc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: workordertype.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: workordertype.FieldID,
			},
		}
	)
	if value, ok := wotc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: workordertype.FieldCreateTime,
		})
		wot.CreateTime = value
	}
	if value, ok := wotc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: workordertype.FieldUpdateTime,
		})
		wot.UpdateTime = value
	}
	if value, ok := wotc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workordertype.FieldName,
		})
		wot.Name = value
	}
	if value, ok := wotc.mutation.Description(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workordertype.FieldDescription,
		})
		wot.Description = value
	}
	if nodes := wotc.mutation.WorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := wotc.mutation.PropertyTypesIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := wotc.mutation.DefinitionsIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := wotc.mutation.CheckListCategoryDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workordertype.CheckListCategoryDefinitionsTable,
			Columns: []string{workordertype.CheckListCategoryDefinitionsColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, wotc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	wot.ID = int(id)
	return wot, nil
}
