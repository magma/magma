// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/project"
	"github.com/facebookincubator/symphony/pkg/ent/projecttype"
	"github.com/facebookincubator/symphony/pkg/ent/propertytype"
	"github.com/facebookincubator/symphony/pkg/ent/workorderdefinition"
)

// ProjectTypeUpdate is the builder for updating ProjectType entities.
type ProjectTypeUpdate struct {
	config
	hooks      []Hook
	mutation   *ProjectTypeMutation
	predicates []predicate.ProjectType
}

// Where adds a new predicate for the builder.
func (ptu *ProjectTypeUpdate) Where(ps ...predicate.ProjectType) *ProjectTypeUpdate {
	ptu.predicates = append(ptu.predicates, ps...)
	return ptu
}

// SetName sets the name field.
func (ptu *ProjectTypeUpdate) SetName(s string) *ProjectTypeUpdate {
	ptu.mutation.SetName(s)
	return ptu
}

// SetDescription sets the description field.
func (ptu *ProjectTypeUpdate) SetDescription(s string) *ProjectTypeUpdate {
	ptu.mutation.SetDescription(s)
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
	ptu.mutation.ClearDescription()
	return ptu
}

// AddProjectIDs adds the projects edge to Project by ids.
func (ptu *ProjectTypeUpdate) AddProjectIDs(ids ...int) *ProjectTypeUpdate {
	ptu.mutation.AddProjectIDs(ids...)
	return ptu
}

// AddProjects adds the projects edges to Project.
func (ptu *ProjectTypeUpdate) AddProjects(p ...*Project) *ProjectTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptu.AddProjectIDs(ids...)
}

// AddPropertyIDs adds the properties edge to PropertyType by ids.
func (ptu *ProjectTypeUpdate) AddPropertyIDs(ids ...int) *ProjectTypeUpdate {
	ptu.mutation.AddPropertyIDs(ids...)
	return ptu
}

// AddProperties adds the properties edges to PropertyType.
func (ptu *ProjectTypeUpdate) AddProperties(p ...*PropertyType) *ProjectTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptu.AddPropertyIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrderDefinition by ids.
func (ptu *ProjectTypeUpdate) AddWorkOrderIDs(ids ...int) *ProjectTypeUpdate {
	ptu.mutation.AddWorkOrderIDs(ids...)
	return ptu
}

// AddWorkOrders adds the work_orders edges to WorkOrderDefinition.
func (ptu *ProjectTypeUpdate) AddWorkOrders(w ...*WorkOrderDefinition) *ProjectTypeUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return ptu.AddWorkOrderIDs(ids...)
}

// RemoveProjectIDs removes the projects edge to Project by ids.
func (ptu *ProjectTypeUpdate) RemoveProjectIDs(ids ...int) *ProjectTypeUpdate {
	ptu.mutation.RemoveProjectIDs(ids...)
	return ptu
}

// RemoveProjects removes projects edges to Project.
func (ptu *ProjectTypeUpdate) RemoveProjects(p ...*Project) *ProjectTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptu.RemoveProjectIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to PropertyType by ids.
func (ptu *ProjectTypeUpdate) RemovePropertyIDs(ids ...int) *ProjectTypeUpdate {
	ptu.mutation.RemovePropertyIDs(ids...)
	return ptu
}

// RemoveProperties removes properties edges to PropertyType.
func (ptu *ProjectTypeUpdate) RemoveProperties(p ...*PropertyType) *ProjectTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptu.RemovePropertyIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrderDefinition by ids.
func (ptu *ProjectTypeUpdate) RemoveWorkOrderIDs(ids ...int) *ProjectTypeUpdate {
	ptu.mutation.RemoveWorkOrderIDs(ids...)
	return ptu
}

// RemoveWorkOrders removes work_orders edges to WorkOrderDefinition.
func (ptu *ProjectTypeUpdate) RemoveWorkOrders(w ...*WorkOrderDefinition) *ProjectTypeUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return ptu.RemoveWorkOrderIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (ptu *ProjectTypeUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := ptu.mutation.UpdateTime(); !ok {
		v := projecttype.UpdateDefaultUpdateTime()
		ptu.mutation.SetUpdateTime(v)
	}
	if v, ok := ptu.mutation.Name(); ok {
		if err := projecttype.NameValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}

	var (
		err      error
		affected int
	)
	if len(ptu.hooks) == 0 {
		affected, err = ptu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ProjectTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ptu.mutation = mutation
			affected, err = ptu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(ptu.hooks) - 1; i >= 0; i-- {
			mut = ptu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ptu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
				Type:   field.TypeInt,
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
	if value, ok := ptu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: projecttype.FieldUpdateTime,
		})
	}
	if value, ok := ptu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: projecttype.FieldName,
		})
	}
	if value, ok := ptu.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: projecttype.FieldDescription,
		})
	}
	if ptu.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: projecttype.FieldDescription,
		})
	}
	if nodes := ptu.mutation.RemovedProjectsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptu.mutation.ProjectsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := ptu.mutation.RemovedPropertiesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptu.mutation.PropertiesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := ptu.mutation.RemovedWorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptu.mutation.WorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ptu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{projecttype.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// ProjectTypeUpdateOne is the builder for updating a single ProjectType entity.
type ProjectTypeUpdateOne struct {
	config
	hooks    []Hook
	mutation *ProjectTypeMutation
}

// SetName sets the name field.
func (ptuo *ProjectTypeUpdateOne) SetName(s string) *ProjectTypeUpdateOne {
	ptuo.mutation.SetName(s)
	return ptuo
}

// SetDescription sets the description field.
func (ptuo *ProjectTypeUpdateOne) SetDescription(s string) *ProjectTypeUpdateOne {
	ptuo.mutation.SetDescription(s)
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
	ptuo.mutation.ClearDescription()
	return ptuo
}

// AddProjectIDs adds the projects edge to Project by ids.
func (ptuo *ProjectTypeUpdateOne) AddProjectIDs(ids ...int) *ProjectTypeUpdateOne {
	ptuo.mutation.AddProjectIDs(ids...)
	return ptuo
}

// AddProjects adds the projects edges to Project.
func (ptuo *ProjectTypeUpdateOne) AddProjects(p ...*Project) *ProjectTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptuo.AddProjectIDs(ids...)
}

// AddPropertyIDs adds the properties edge to PropertyType by ids.
func (ptuo *ProjectTypeUpdateOne) AddPropertyIDs(ids ...int) *ProjectTypeUpdateOne {
	ptuo.mutation.AddPropertyIDs(ids...)
	return ptuo
}

// AddProperties adds the properties edges to PropertyType.
func (ptuo *ProjectTypeUpdateOne) AddProperties(p ...*PropertyType) *ProjectTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptuo.AddPropertyIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrderDefinition by ids.
func (ptuo *ProjectTypeUpdateOne) AddWorkOrderIDs(ids ...int) *ProjectTypeUpdateOne {
	ptuo.mutation.AddWorkOrderIDs(ids...)
	return ptuo
}

// AddWorkOrders adds the work_orders edges to WorkOrderDefinition.
func (ptuo *ProjectTypeUpdateOne) AddWorkOrders(w ...*WorkOrderDefinition) *ProjectTypeUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return ptuo.AddWorkOrderIDs(ids...)
}

// RemoveProjectIDs removes the projects edge to Project by ids.
func (ptuo *ProjectTypeUpdateOne) RemoveProjectIDs(ids ...int) *ProjectTypeUpdateOne {
	ptuo.mutation.RemoveProjectIDs(ids...)
	return ptuo
}

// RemoveProjects removes projects edges to Project.
func (ptuo *ProjectTypeUpdateOne) RemoveProjects(p ...*Project) *ProjectTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptuo.RemoveProjectIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to PropertyType by ids.
func (ptuo *ProjectTypeUpdateOne) RemovePropertyIDs(ids ...int) *ProjectTypeUpdateOne {
	ptuo.mutation.RemovePropertyIDs(ids...)
	return ptuo
}

// RemoveProperties removes properties edges to PropertyType.
func (ptuo *ProjectTypeUpdateOne) RemoveProperties(p ...*PropertyType) *ProjectTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptuo.RemovePropertyIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrderDefinition by ids.
func (ptuo *ProjectTypeUpdateOne) RemoveWorkOrderIDs(ids ...int) *ProjectTypeUpdateOne {
	ptuo.mutation.RemoveWorkOrderIDs(ids...)
	return ptuo
}

// RemoveWorkOrders removes work_orders edges to WorkOrderDefinition.
func (ptuo *ProjectTypeUpdateOne) RemoveWorkOrders(w ...*WorkOrderDefinition) *ProjectTypeUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return ptuo.RemoveWorkOrderIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (ptuo *ProjectTypeUpdateOne) Save(ctx context.Context) (*ProjectType, error) {
	if _, ok := ptuo.mutation.UpdateTime(); !ok {
		v := projecttype.UpdateDefaultUpdateTime()
		ptuo.mutation.SetUpdateTime(v)
	}
	if v, ok := ptuo.mutation.Name(); ok {
		if err := projecttype.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}

	var (
		err  error
		node *ProjectType
	)
	if len(ptuo.hooks) == 0 {
		node, err = ptuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ProjectTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ptuo.mutation = mutation
			node, err = ptuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(ptuo.hooks) - 1; i >= 0; i-- {
			mut = ptuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ptuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
				Type:   field.TypeInt,
				Column: projecttype.FieldID,
			},
		},
	}
	id, ok := ptuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing ProjectType.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := ptuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: projecttype.FieldUpdateTime,
		})
	}
	if value, ok := ptuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: projecttype.FieldName,
		})
	}
	if value, ok := ptuo.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: projecttype.FieldDescription,
		})
	}
	if ptuo.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: projecttype.FieldDescription,
		})
	}
	if nodes := ptuo.mutation.RemovedProjectsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptuo.mutation.ProjectsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := ptuo.mutation.RemovedPropertiesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptuo.mutation.PropertiesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := ptuo.mutation.RemovedWorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptuo.mutation.WorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	pt = &ProjectType{config: ptuo.config}
	_spec.Assign = pt.assignValues
	_spec.ScanValues = pt.scanValues()
	if err = sqlgraph.UpdateNode(ctx, ptuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{projecttype.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return pt, nil
}
