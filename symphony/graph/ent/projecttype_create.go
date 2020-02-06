// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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
	create_time *time.Time
	update_time *time.Time
	name        *string
	description *string
	projects    map[string]struct{}
	properties  map[string]struct{}
	work_orders map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (ptc *ProjectTypeCreate) SetCreateTime(t time.Time) *ProjectTypeCreate {
	ptc.create_time = &t
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
	ptc.update_time = &t
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
	ptc.name = &s
	return ptc
}

// SetDescription sets the description field.
func (ptc *ProjectTypeCreate) SetDescription(s string) *ProjectTypeCreate {
	ptc.description = &s
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
func (ptc *ProjectTypeCreate) AddProjectIDs(ids ...string) *ProjectTypeCreate {
	if ptc.projects == nil {
		ptc.projects = make(map[string]struct{})
	}
	for i := range ids {
		ptc.projects[ids[i]] = struct{}{}
	}
	return ptc
}

// AddProjects adds the projects edges to Project.
func (ptc *ProjectTypeCreate) AddProjects(p ...*Project) *ProjectTypeCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptc.AddProjectIDs(ids...)
}

// AddPropertyIDs adds the properties edge to PropertyType by ids.
func (ptc *ProjectTypeCreate) AddPropertyIDs(ids ...string) *ProjectTypeCreate {
	if ptc.properties == nil {
		ptc.properties = make(map[string]struct{})
	}
	for i := range ids {
		ptc.properties[ids[i]] = struct{}{}
	}
	return ptc
}

// AddProperties adds the properties edges to PropertyType.
func (ptc *ProjectTypeCreate) AddProperties(p ...*PropertyType) *ProjectTypeCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptc.AddPropertyIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrderDefinition by ids.
func (ptc *ProjectTypeCreate) AddWorkOrderIDs(ids ...string) *ProjectTypeCreate {
	if ptc.work_orders == nil {
		ptc.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		ptc.work_orders[ids[i]] = struct{}{}
	}
	return ptc
}

// AddWorkOrders adds the work_orders edges to WorkOrderDefinition.
func (ptc *ProjectTypeCreate) AddWorkOrders(w ...*WorkOrderDefinition) *ProjectTypeCreate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return ptc.AddWorkOrderIDs(ids...)
}

// Save creates the ProjectType in the database.
func (ptc *ProjectTypeCreate) Save(ctx context.Context) (*ProjectType, error) {
	if ptc.create_time == nil {
		v := projecttype.DefaultCreateTime()
		ptc.create_time = &v
	}
	if ptc.update_time == nil {
		v := projecttype.DefaultUpdateTime()
		ptc.update_time = &v
	}
	if ptc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if err := projecttype.NameValidator(*ptc.name); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
	}
	return ptc.sqlSave(ctx)
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
				Type:   field.TypeString,
				Column: projecttype.FieldID,
			},
		}
	)
	if value := ptc.create_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: projecttype.FieldCreateTime,
		})
		pt.CreateTime = *value
	}
	if value := ptc.update_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: projecttype.FieldUpdateTime,
		})
		pt.UpdateTime = *value
	}
	if value := ptc.name; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: projecttype.FieldName,
		})
		pt.Name = *value
	}
	if value := ptc.description; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: projecttype.FieldDescription,
		})
		pt.Description = value
	}
	if nodes := ptc.projects; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   projecttype.ProjectsTable,
			Columns: []string{projecttype.ProjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: project.FieldID,
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
	if nodes := ptc.properties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   projecttype.PropertiesTable,
			Columns: []string{projecttype.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: propertytype.FieldID,
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
	if nodes := ptc.work_orders; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   projecttype.WorkOrdersTable,
			Columns: []string{projecttype.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorderdefinition.FieldID,
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
	if err := sqlgraph.CreateNode(ctx, ptc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	pt.ID = strconv.FormatInt(id, 10)
	return pt, nil
}
