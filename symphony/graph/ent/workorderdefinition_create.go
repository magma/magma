// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// WorkOrderDefinitionCreate is the builder for creating a WorkOrderDefinition entity.
type WorkOrderDefinitionCreate struct {
	config
	mutation *WorkOrderDefinitionMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (wodc *WorkOrderDefinitionCreate) SetCreateTime(t time.Time) *WorkOrderDefinitionCreate {
	wodc.mutation.SetCreateTime(t)
	return wodc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (wodc *WorkOrderDefinitionCreate) SetNillableCreateTime(t *time.Time) *WorkOrderDefinitionCreate {
	if t != nil {
		wodc.SetCreateTime(*t)
	}
	return wodc
}

// SetUpdateTime sets the update_time field.
func (wodc *WorkOrderDefinitionCreate) SetUpdateTime(t time.Time) *WorkOrderDefinitionCreate {
	wodc.mutation.SetUpdateTime(t)
	return wodc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (wodc *WorkOrderDefinitionCreate) SetNillableUpdateTime(t *time.Time) *WorkOrderDefinitionCreate {
	if t != nil {
		wodc.SetUpdateTime(*t)
	}
	return wodc
}

// SetIndex sets the index field.
func (wodc *WorkOrderDefinitionCreate) SetIndex(i int) *WorkOrderDefinitionCreate {
	wodc.mutation.SetIndex(i)
	return wodc
}

// SetNillableIndex sets the index field if the given value is not nil.
func (wodc *WorkOrderDefinitionCreate) SetNillableIndex(i *int) *WorkOrderDefinitionCreate {
	if i != nil {
		wodc.SetIndex(*i)
	}
	return wodc
}

// SetTypeID sets the type edge to WorkOrderType by id.
func (wodc *WorkOrderDefinitionCreate) SetTypeID(id int) *WorkOrderDefinitionCreate {
	wodc.mutation.SetTypeID(id)
	return wodc
}

// SetNillableTypeID sets the type edge to WorkOrderType by id if the given value is not nil.
func (wodc *WorkOrderDefinitionCreate) SetNillableTypeID(id *int) *WorkOrderDefinitionCreate {
	if id != nil {
		wodc = wodc.SetTypeID(*id)
	}
	return wodc
}

// SetType sets the type edge to WorkOrderType.
func (wodc *WorkOrderDefinitionCreate) SetType(w *WorkOrderType) *WorkOrderDefinitionCreate {
	return wodc.SetTypeID(w.ID)
}

// SetProjectTypeID sets the project_type edge to ProjectType by id.
func (wodc *WorkOrderDefinitionCreate) SetProjectTypeID(id int) *WorkOrderDefinitionCreate {
	wodc.mutation.SetProjectTypeID(id)
	return wodc
}

// SetNillableProjectTypeID sets the project_type edge to ProjectType by id if the given value is not nil.
func (wodc *WorkOrderDefinitionCreate) SetNillableProjectTypeID(id *int) *WorkOrderDefinitionCreate {
	if id != nil {
		wodc = wodc.SetProjectTypeID(*id)
	}
	return wodc
}

// SetProjectType sets the project_type edge to ProjectType.
func (wodc *WorkOrderDefinitionCreate) SetProjectType(p *ProjectType) *WorkOrderDefinitionCreate {
	return wodc.SetProjectTypeID(p.ID)
}

// Save creates the WorkOrderDefinition in the database.
func (wodc *WorkOrderDefinitionCreate) Save(ctx context.Context) (*WorkOrderDefinition, error) {
	if _, ok := wodc.mutation.CreateTime(); !ok {
		v := workorderdefinition.DefaultCreateTime()
		wodc.mutation.SetCreateTime(v)
	}
	if _, ok := wodc.mutation.UpdateTime(); !ok {
		v := workorderdefinition.DefaultUpdateTime()
		wodc.mutation.SetUpdateTime(v)
	}
	var (
		err  error
		node *WorkOrderDefinition
	)
	if len(wodc.hooks) == 0 {
		node, err = wodc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*WorkOrderDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			wodc.mutation = mutation
			node, err = wodc.sqlSave(ctx)
			return node, err
		})
		for i := len(wodc.hooks); i > 0; i-- {
			mut = wodc.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, wodc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (wodc *WorkOrderDefinitionCreate) SaveX(ctx context.Context) *WorkOrderDefinition {
	v, err := wodc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (wodc *WorkOrderDefinitionCreate) sqlSave(ctx context.Context) (*WorkOrderDefinition, error) {
	var (
		wod   = &WorkOrderDefinition{config: wodc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: workorderdefinition.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: workorderdefinition.FieldID,
			},
		}
	)
	if value, ok := wodc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: workorderdefinition.FieldCreateTime,
		})
		wod.CreateTime = value
	}
	if value, ok := wodc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: workorderdefinition.FieldUpdateTime,
		})
		wod.UpdateTime = value
	}
	if value, ok := wodc.mutation.Index(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: workorderdefinition.FieldIndex,
		})
		wod.Index = value
	}
	if nodes := wodc.mutation.TypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorderdefinition.TypeTable,
			Columns: []string{workorderdefinition.TypeColumn},
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
	if nodes := wodc.mutation.ProjectTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   workorderdefinition.ProjectTypeTable,
			Columns: []string{workorderdefinition.ProjectTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: projecttype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, wodc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	wod.ID = int(id)
	return wod, nil
}
