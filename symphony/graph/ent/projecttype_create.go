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
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
)

// ProjectTypeCreate is the builder for creating a ProjectType entity.
type ProjectTypeCreate struct {
	config
	mutation *ProjectTypeMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (ptc *ProjectTypeCreate) SetCreateTime(t time.Time) *ProjectTypeCreate {
	ptc.mutation.SetCreateTime(t)
	return ptc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (ptc *ProjectTypeCreate) SetNillableCreateTime(t *time.Time) *ProjectTypeCreate {
	if t != nil {
		ptc.SetCreateTime(*t)
	}
	return ptc
}

// SetUpdateTime sets the update_time field.
func (ptc *ProjectTypeCreate) SetUpdateTime(t time.Time) *ProjectTypeCreate {
	ptc.mutation.SetUpdateTime(t)
	return ptc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (ptc *ProjectTypeCreate) SetNillableUpdateTime(t *time.Time) *ProjectTypeCreate {
	if t != nil {
		ptc.SetUpdateTime(*t)
	}
	return ptc
}

// SetName sets the name field.
func (ptc *ProjectTypeCreate) SetName(s string) *ProjectTypeCreate {
	ptc.mutation.SetName(s)
	return ptc
}

// SetDescription sets the description field.
func (ptc *ProjectTypeCreate) SetDescription(s string) *ProjectTypeCreate {
	ptc.mutation.SetDescription(s)
	return ptc
}

// SetNillableDescription sets the description field if the given value is not nil.
func (ptc *ProjectTypeCreate) SetNillableDescription(s *string) *ProjectTypeCreate {
	if s != nil {
		ptc.SetDescription(*s)
	}
	return ptc
}

// AddProjectIDs adds the projects edge to Project by ids.
func (ptc *ProjectTypeCreate) AddProjectIDs(ids ...int) *ProjectTypeCreate {
	ptc.mutation.AddProjectIDs(ids...)
	return ptc
}

// AddProjects adds the projects edges to Project.
func (ptc *ProjectTypeCreate) AddProjects(p ...*Project) *ProjectTypeCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptc.AddProjectIDs(ids...)
}

// AddPropertyIDs adds the properties edge to PropertyType by ids.
func (ptc *ProjectTypeCreate) AddPropertyIDs(ids ...int) *ProjectTypeCreate {
	ptc.mutation.AddPropertyIDs(ids...)
	return ptc
}

// AddProperties adds the properties edges to PropertyType.
func (ptc *ProjectTypeCreate) AddProperties(p ...*PropertyType) *ProjectTypeCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptc.AddPropertyIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrderDefinition by ids.
func (ptc *ProjectTypeCreate) AddWorkOrderIDs(ids ...int) *ProjectTypeCreate {
	ptc.mutation.AddWorkOrderIDs(ids...)
	return ptc
}

// AddWorkOrders adds the work_orders edges to WorkOrderDefinition.
func (ptc *ProjectTypeCreate) AddWorkOrders(w ...*WorkOrderDefinition) *ProjectTypeCreate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return ptc.AddWorkOrderIDs(ids...)
}

// Save creates the ProjectType in the database.
func (ptc *ProjectTypeCreate) Save(ctx context.Context) (*ProjectType, error) {
	if _, ok := ptc.mutation.CreateTime(); !ok {
		v := projecttype.DefaultCreateTime()
		ptc.mutation.SetCreateTime(v)
	}
	if _, ok := ptc.mutation.UpdateTime(); !ok {
		v := projecttype.DefaultUpdateTime()
		ptc.mutation.SetUpdateTime(v)
	}
	if _, ok := ptc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if v, ok := ptc.mutation.Name(); ok {
		if err := projecttype.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	var (
		err  error
		node *ProjectType
	)
	if len(ptc.hooks) == 0 {
		node, err = ptc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ProjectTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ptc.mutation = mutation
			node, err = ptc.sqlSave(ctx)
			return node, err
		})
		for i := len(ptc.hooks) - 1; i >= 0; i-- {
			mut = ptc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ptc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (ptc *ProjectTypeCreate) SaveX(ctx context.Context) *ProjectType {
	v, err := ptc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ptc *ProjectTypeCreate) sqlSave(ctx context.Context) (*ProjectType, error) {
	var (
		pt    = &ProjectType{config: ptc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: projecttype.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: projecttype.FieldID,
			},
		}
	)
	if value, ok := ptc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: projecttype.FieldCreateTime,
		})
		pt.CreateTime = value
	}
	if value, ok := ptc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: projecttype.FieldUpdateTime,
		})
		pt.UpdateTime = value
	}
	if value, ok := ptc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: projecttype.FieldName,
		})
		pt.Name = value
	}
	if value, ok := ptc.mutation.Description(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: projecttype.FieldDescription,
		})
		pt.Description = &value
	}
	if nodes := ptc.mutation.ProjectsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   projecttype.ProjectsTable,
			Columns: []string{projecttype.ProjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: project.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ptc.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   projecttype.PropertiesTable,
			Columns: []string{projecttype.PropertiesColumn},
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
	if nodes := ptc.mutation.WorkOrdersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   projecttype.WorkOrdersTable,
			Columns: []string{projecttype.WorkOrdersColumn},
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
	if err := sqlgraph.CreateNode(ctx, ptc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	pt.ID = int(id)
	return pt, nil
}
