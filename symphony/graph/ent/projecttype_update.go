// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
)

// ProjectTypeUpdate is the builder for updating ProjectType entities.
type ProjectTypeUpdate struct {
	config

	update_time       *time.Time
	name              *string
	description       *string
	cleardescription  bool
	projects          map[string]struct{}
	properties        map[string]struct{}
	work_orders       map[string]struct{}
	removedProjects   map[string]struct{}
	removedProperties map[string]struct{}
	removedWorkOrders map[string]struct{}
	predicates        []predicate.ProjectType
}

// Where adds a new predicate for the builder.
func (ptu *ProjectTypeUpdate) Where(ps ...predicate.ProjectType) *ProjectTypeUpdate {
	ptu.predicates = append(ptu.predicates, ps...)
	return ptu
}

// SetName sets the name field.
func (ptu *ProjectTypeUpdate) SetName(s string) *ProjectTypeUpdate {
	ptu.name = &s
	return ptu
}

// SetDescription sets the description field.
func (ptu *ProjectTypeUpdate) SetDescription(s string) *ProjectTypeUpdate {
	ptu.description = &s
	return ptu
}

// SetNillableDescription sets the description field if the given value is not nil.
func (ptu *ProjectTypeUpdate) SetNillableDescription(s *string) *ProjectTypeUpdate {
	if s != nil {
		ptu.SetDescription(*s)
	}
	return ptu
}

// ClearDescription clears the value of description.
func (ptu *ProjectTypeUpdate) ClearDescription() *ProjectTypeUpdate {
	ptu.description = nil
	ptu.cleardescription = true
	return ptu
}

// AddProjectIDs adds the projects edge to Project by ids.
func (ptu *ProjectTypeUpdate) AddProjectIDs(ids ...string) *ProjectTypeUpdate {
	if ptu.projects == nil {
		ptu.projects = make(map[string]struct{})
	}
	for i := range ids {
		ptu.projects[ids[i]] = struct{}{}
	}
	return ptu
}

// AddProjects adds the projects edges to Project.
func (ptu *ProjectTypeUpdate) AddProjects(p ...*Project) *ProjectTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptu.AddProjectIDs(ids...)
}

// AddPropertyIDs adds the properties edge to PropertyType by ids.
func (ptu *ProjectTypeUpdate) AddPropertyIDs(ids ...string) *ProjectTypeUpdate {
	if ptu.properties == nil {
		ptu.properties = make(map[string]struct{})
	}
	for i := range ids {
		ptu.properties[ids[i]] = struct{}{}
	}
	return ptu
}

// AddProperties adds the properties edges to PropertyType.
func (ptu *ProjectTypeUpdate) AddProperties(p ...*PropertyType) *ProjectTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptu.AddPropertyIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrderDefinition by ids.
func (ptu *ProjectTypeUpdate) AddWorkOrderIDs(ids ...string) *ProjectTypeUpdate {
	if ptu.work_orders == nil {
		ptu.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		ptu.work_orders[ids[i]] = struct{}{}
	}
	return ptu
}

// AddWorkOrders adds the work_orders edges to WorkOrderDefinition.
func (ptu *ProjectTypeUpdate) AddWorkOrders(w ...*WorkOrderDefinition) *ProjectTypeUpdate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return ptu.AddWorkOrderIDs(ids...)
}

// RemoveProjectIDs removes the projects edge to Project by ids.
func (ptu *ProjectTypeUpdate) RemoveProjectIDs(ids ...string) *ProjectTypeUpdate {
	if ptu.removedProjects == nil {
		ptu.removedProjects = make(map[string]struct{})
	}
	for i := range ids {
		ptu.removedProjects[ids[i]] = struct{}{}
	}
	return ptu
}

// RemoveProjects removes projects edges to Project.
func (ptu *ProjectTypeUpdate) RemoveProjects(p ...*Project) *ProjectTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptu.RemoveProjectIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to PropertyType by ids.
func (ptu *ProjectTypeUpdate) RemovePropertyIDs(ids ...string) *ProjectTypeUpdate {
	if ptu.removedProperties == nil {
		ptu.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		ptu.removedProperties[ids[i]] = struct{}{}
	}
	return ptu
}

// RemoveProperties removes properties edges to PropertyType.
func (ptu *ProjectTypeUpdate) RemoveProperties(p ...*PropertyType) *ProjectTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptu.RemovePropertyIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrderDefinition by ids.
func (ptu *ProjectTypeUpdate) RemoveWorkOrderIDs(ids ...string) *ProjectTypeUpdate {
	if ptu.removedWorkOrders == nil {
		ptu.removedWorkOrders = make(map[string]struct{})
	}
	for i := range ids {
		ptu.removedWorkOrders[ids[i]] = struct{}{}
	}
	return ptu
}

// RemoveWorkOrders removes work_orders edges to WorkOrderDefinition.
func (ptu *ProjectTypeUpdate) RemoveWorkOrders(w ...*WorkOrderDefinition) *ProjectTypeUpdate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return ptu.RemoveWorkOrderIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (ptu *ProjectTypeUpdate) Save(ctx context.Context) (int, error) {
	if ptu.update_time == nil {
		v := projecttype.UpdateDefaultUpdateTime()
		ptu.update_time = &v
	}
	if ptu.name != nil {
		if err := projecttype.NameValidator(*ptu.name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	return ptu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (ptu *ProjectTypeUpdate) SaveX(ctx context.Context) int {
	affected, err := ptu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ptu *ProjectTypeUpdate) Exec(ctx context.Context) error {
	_, err := ptu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ptu *ProjectTypeUpdate) ExecX(ctx context.Context) {
	if err := ptu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ptu *ProjectTypeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   projecttype.Table,
			Columns: projecttype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: projecttype.FieldID,
			},
		},
	}
	if ps := ptu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := ptu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: projecttype.FieldUpdateTime,
		})
	}
	if value := ptu.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: projecttype.FieldName,
		})
	}
	if value := ptu.description; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: projecttype.FieldDescription,
		})
	}
	if ptu.cleardescription {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: projecttype.FieldDescription,
		})
	}
	if nodes := ptu.removedProjects; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptu.projects; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := ptu.removedProperties; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptu.properties; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := ptu.removedWorkOrders; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptu.work_orders; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ptu.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// ProjectTypeUpdateOne is the builder for updating a single ProjectType entity.
type ProjectTypeUpdateOne struct {
	config
	id string

	update_time       *time.Time
	name              *string
	description       *string
	cleardescription  bool
	projects          map[string]struct{}
	properties        map[string]struct{}
	work_orders       map[string]struct{}
	removedProjects   map[string]struct{}
	removedProperties map[string]struct{}
	removedWorkOrders map[string]struct{}
}

// SetName sets the name field.
func (ptuo *ProjectTypeUpdateOne) SetName(s string) *ProjectTypeUpdateOne {
	ptuo.name = &s
	return ptuo
}

// SetDescription sets the description field.
func (ptuo *ProjectTypeUpdateOne) SetDescription(s string) *ProjectTypeUpdateOne {
	ptuo.description = &s
	return ptuo
}

// SetNillableDescription sets the description field if the given value is not nil.
func (ptuo *ProjectTypeUpdateOne) SetNillableDescription(s *string) *ProjectTypeUpdateOne {
	if s != nil {
		ptuo.SetDescription(*s)
	}
	return ptuo
}

// ClearDescription clears the value of description.
func (ptuo *ProjectTypeUpdateOne) ClearDescription() *ProjectTypeUpdateOne {
	ptuo.description = nil
	ptuo.cleardescription = true
	return ptuo
}

// AddProjectIDs adds the projects edge to Project by ids.
func (ptuo *ProjectTypeUpdateOne) AddProjectIDs(ids ...string) *ProjectTypeUpdateOne {
	if ptuo.projects == nil {
		ptuo.projects = make(map[string]struct{})
	}
	for i := range ids {
		ptuo.projects[ids[i]] = struct{}{}
	}
	return ptuo
}

// AddProjects adds the projects edges to Project.
func (ptuo *ProjectTypeUpdateOne) AddProjects(p ...*Project) *ProjectTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptuo.AddProjectIDs(ids...)
}

// AddPropertyIDs adds the properties edge to PropertyType by ids.
func (ptuo *ProjectTypeUpdateOne) AddPropertyIDs(ids ...string) *ProjectTypeUpdateOne {
	if ptuo.properties == nil {
		ptuo.properties = make(map[string]struct{})
	}
	for i := range ids {
		ptuo.properties[ids[i]] = struct{}{}
	}
	return ptuo
}

// AddProperties adds the properties edges to PropertyType.
func (ptuo *ProjectTypeUpdateOne) AddProperties(p ...*PropertyType) *ProjectTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptuo.AddPropertyIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrderDefinition by ids.
func (ptuo *ProjectTypeUpdateOne) AddWorkOrderIDs(ids ...string) *ProjectTypeUpdateOne {
	if ptuo.work_orders == nil {
		ptuo.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		ptuo.work_orders[ids[i]] = struct{}{}
	}
	return ptuo
}

// AddWorkOrders adds the work_orders edges to WorkOrderDefinition.
func (ptuo *ProjectTypeUpdateOne) AddWorkOrders(w ...*WorkOrderDefinition) *ProjectTypeUpdateOne {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return ptuo.AddWorkOrderIDs(ids...)
}

// RemoveProjectIDs removes the projects edge to Project by ids.
func (ptuo *ProjectTypeUpdateOne) RemoveProjectIDs(ids ...string) *ProjectTypeUpdateOne {
	if ptuo.removedProjects == nil {
		ptuo.removedProjects = make(map[string]struct{})
	}
	for i := range ids {
		ptuo.removedProjects[ids[i]] = struct{}{}
	}
	return ptuo
}

// RemoveProjects removes projects edges to Project.
func (ptuo *ProjectTypeUpdateOne) RemoveProjects(p ...*Project) *ProjectTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptuo.RemoveProjectIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to PropertyType by ids.
func (ptuo *ProjectTypeUpdateOne) RemovePropertyIDs(ids ...string) *ProjectTypeUpdateOne {
	if ptuo.removedProperties == nil {
		ptuo.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		ptuo.removedProperties[ids[i]] = struct{}{}
	}
	return ptuo
}

// RemoveProperties removes properties edges to PropertyType.
func (ptuo *ProjectTypeUpdateOne) RemoveProperties(p ...*PropertyType) *ProjectTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptuo.RemovePropertyIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrderDefinition by ids.
func (ptuo *ProjectTypeUpdateOne) RemoveWorkOrderIDs(ids ...string) *ProjectTypeUpdateOne {
	if ptuo.removedWorkOrders == nil {
		ptuo.removedWorkOrders = make(map[string]struct{})
	}
	for i := range ids {
		ptuo.removedWorkOrders[ids[i]] = struct{}{}
	}
	return ptuo
}

// RemoveWorkOrders removes work_orders edges to WorkOrderDefinition.
func (ptuo *ProjectTypeUpdateOne) RemoveWorkOrders(w ...*WorkOrderDefinition) *ProjectTypeUpdateOne {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return ptuo.RemoveWorkOrderIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (ptuo *ProjectTypeUpdateOne) Save(ctx context.Context) (*ProjectType, error) {
	if ptuo.update_time == nil {
		v := projecttype.UpdateDefaultUpdateTime()
		ptuo.update_time = &v
	}
	if ptuo.name != nil {
		if err := projecttype.NameValidator(*ptuo.name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	return ptuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (ptuo *ProjectTypeUpdateOne) SaveX(ctx context.Context) *ProjectType {
	pt, err := ptuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return pt
}

// Exec executes the query on the entity.
func (ptuo *ProjectTypeUpdateOne) Exec(ctx context.Context) error {
	_, err := ptuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ptuo *ProjectTypeUpdateOne) ExecX(ctx context.Context) {
	if err := ptuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ptuo *ProjectTypeUpdateOne) sqlSave(ctx context.Context) (pt *ProjectType, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   projecttype.Table,
			Columns: projecttype.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  ptuo.id,
				Type:   field.TypeString,
				Column: projecttype.FieldID,
			},
		},
	}
	if value := ptuo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: projecttype.FieldUpdateTime,
		})
	}
	if value := ptuo.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: projecttype.FieldName,
		})
	}
	if value := ptuo.description; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: projecttype.FieldDescription,
		})
	}
	if ptuo.cleardescription {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: projecttype.FieldDescription,
		})
	}
	if nodes := ptuo.removedProjects; len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptuo.projects; len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := ptuo.removedProperties; len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptuo.properties; len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := ptuo.removedWorkOrders; len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptuo.work_orders; len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	pt = &ProjectType{config: ptuo.config}
	_spec.Assign = pt.assignValues
	_spec.ScanValues = pt.scanValues()
	if err = sqlgraph.UpdateNode(ctx, ptuo.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return pt, nil
}
