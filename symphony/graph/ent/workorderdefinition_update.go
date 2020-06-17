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
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// WorkOrderDefinitionUpdate is the builder for updating WorkOrderDefinition entities.
type WorkOrderDefinitionUpdate struct {
	config
	hooks      []Hook
	mutation   *WorkOrderDefinitionMutation
	predicates []predicate.WorkOrderDefinition
}

// Where adds a new predicate for the builder.
func (wodu *WorkOrderDefinitionUpdate) Where(ps ...predicate.WorkOrderDefinition) *WorkOrderDefinitionUpdate {
	wodu.predicates = append(wodu.predicates, ps...)
	return wodu
}

// SetIndex sets the index field.
func (wodu *WorkOrderDefinitionUpdate) SetIndex(i int) *WorkOrderDefinitionUpdate {
	wodu.mutation.ResetIndex()
	wodu.mutation.SetIndex(i)
	return wodu
}

// SetNillableIndex sets the index field if the given value is not nil.
func (wodu *WorkOrderDefinitionUpdate) SetNillableIndex(i *int) *WorkOrderDefinitionUpdate {
	if i != nil {
		wodu.SetIndex(*i)
	}
	return wodu
}

// AddIndex adds i to index.
func (wodu *WorkOrderDefinitionUpdate) AddIndex(i int) *WorkOrderDefinitionUpdate {
	wodu.mutation.AddIndex(i)
	return wodu
}

// ClearIndex clears the value of index.
func (wodu *WorkOrderDefinitionUpdate) ClearIndex() *WorkOrderDefinitionUpdate {
	wodu.mutation.ClearIndex()
	return wodu
}

// SetTypeID sets the type edge to WorkOrderType by id.
func (wodu *WorkOrderDefinitionUpdate) SetTypeID(id int) *WorkOrderDefinitionUpdate {
	wodu.mutation.SetTypeID(id)
	return wodu
}

// SetNillableTypeID sets the type edge to WorkOrderType by id if the given value is not nil.
func (wodu *WorkOrderDefinitionUpdate) SetNillableTypeID(id *int) *WorkOrderDefinitionUpdate {
	if id != nil {
		wodu = wodu.SetTypeID(*id)
	}
	return wodu
}

// SetType sets the type edge to WorkOrderType.
func (wodu *WorkOrderDefinitionUpdate) SetType(w *WorkOrderType) *WorkOrderDefinitionUpdate {
	return wodu.SetTypeID(w.ID)
}

// SetProjectTypeID sets the project_type edge to ProjectType by id.
func (wodu *WorkOrderDefinitionUpdate) SetProjectTypeID(id int) *WorkOrderDefinitionUpdate {
	wodu.mutation.SetProjectTypeID(id)
	return wodu
}

// SetNillableProjectTypeID sets the project_type edge to ProjectType by id if the given value is not nil.
func (wodu *WorkOrderDefinitionUpdate) SetNillableProjectTypeID(id *int) *WorkOrderDefinitionUpdate {
	if id != nil {
		wodu = wodu.SetProjectTypeID(*id)
	}
	return wodu
}

// SetProjectType sets the project_type edge to ProjectType.
func (wodu *WorkOrderDefinitionUpdate) SetProjectType(p *ProjectType) *WorkOrderDefinitionUpdate {
	return wodu.SetProjectTypeID(p.ID)
}

// ClearType clears the type edge to WorkOrderType.
func (wodu *WorkOrderDefinitionUpdate) ClearType() *WorkOrderDefinitionUpdate {
	wodu.mutation.ClearType()
	return wodu
}

// ClearProjectType clears the project_type edge to ProjectType.
func (wodu *WorkOrderDefinitionUpdate) ClearProjectType() *WorkOrderDefinitionUpdate {
	wodu.mutation.ClearProjectType()
	return wodu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (wodu *WorkOrderDefinitionUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := wodu.mutation.UpdateTime(); !ok {
		v := workorderdefinition.UpdateDefaultUpdateTime()
		wodu.mutation.SetUpdateTime(v)
	}

	var (
		err      error
		affected int
	)
	if len(wodu.hooks) == 0 {
		affected, err = wodu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*WorkOrderDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			wodu.mutation = mutation
			affected, err = wodu.sqlSave(ctx)
			return affected, err
		})
		for i := len(wodu.hooks) - 1; i >= 0; i-- {
			mut = wodu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, wodu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (wodu *WorkOrderDefinitionUpdate) SaveX(ctx context.Context) int {
	affected, err := wodu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (wodu *WorkOrderDefinitionUpdate) Exec(ctx context.Context) error {
	_, err := wodu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wodu *WorkOrderDefinitionUpdate) ExecX(ctx context.Context) {
	if err := wodu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (wodu *WorkOrderDefinitionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   workorderdefinition.Table,
			Columns: workorderdefinition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: workorderdefinition.FieldID,
			},
		},
	}
	if ps := wodu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := wodu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: workorderdefinition.FieldUpdateTime,
		})
	}
	if value, ok := wodu.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: workorderdefinition.FieldIndex,
		})
	}
	if value, ok := wodu.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: workorderdefinition.FieldIndex,
		})
	}
	if wodu.mutation.IndexCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: workorderdefinition.FieldIndex,
		})
	}
	if wodu.mutation.TypeCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wodu.mutation.TypeIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if wodu.mutation.ProjectTypeCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wodu.mutation.ProjectTypeIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, wodu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{workorderdefinition.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// WorkOrderDefinitionUpdateOne is the builder for updating a single WorkOrderDefinition entity.
type WorkOrderDefinitionUpdateOne struct {
	config
	hooks    []Hook
	mutation *WorkOrderDefinitionMutation
}

// SetIndex sets the index field.
func (woduo *WorkOrderDefinitionUpdateOne) SetIndex(i int) *WorkOrderDefinitionUpdateOne {
	woduo.mutation.ResetIndex()
	woduo.mutation.SetIndex(i)
	return woduo
}

// SetNillableIndex sets the index field if the given value is not nil.
func (woduo *WorkOrderDefinitionUpdateOne) SetNillableIndex(i *int) *WorkOrderDefinitionUpdateOne {
	if i != nil {
		woduo.SetIndex(*i)
	}
	return woduo
}

// AddIndex adds i to index.
func (woduo *WorkOrderDefinitionUpdateOne) AddIndex(i int) *WorkOrderDefinitionUpdateOne {
	woduo.mutation.AddIndex(i)
	return woduo
}

// ClearIndex clears the value of index.
func (woduo *WorkOrderDefinitionUpdateOne) ClearIndex() *WorkOrderDefinitionUpdateOne {
	woduo.mutation.ClearIndex()
	return woduo
}

// SetTypeID sets the type edge to WorkOrderType by id.
func (woduo *WorkOrderDefinitionUpdateOne) SetTypeID(id int) *WorkOrderDefinitionUpdateOne {
	woduo.mutation.SetTypeID(id)
	return woduo
}

// SetNillableTypeID sets the type edge to WorkOrderType by id if the given value is not nil.
func (woduo *WorkOrderDefinitionUpdateOne) SetNillableTypeID(id *int) *WorkOrderDefinitionUpdateOne {
	if id != nil {
		woduo = woduo.SetTypeID(*id)
	}
	return woduo
}

// SetType sets the type edge to WorkOrderType.
func (woduo *WorkOrderDefinitionUpdateOne) SetType(w *WorkOrderType) *WorkOrderDefinitionUpdateOne {
	return woduo.SetTypeID(w.ID)
}

// SetProjectTypeID sets the project_type edge to ProjectType by id.
func (woduo *WorkOrderDefinitionUpdateOne) SetProjectTypeID(id int) *WorkOrderDefinitionUpdateOne {
	woduo.mutation.SetProjectTypeID(id)
	return woduo
}

// SetNillableProjectTypeID sets the project_type edge to ProjectType by id if the given value is not nil.
func (woduo *WorkOrderDefinitionUpdateOne) SetNillableProjectTypeID(id *int) *WorkOrderDefinitionUpdateOne {
	if id != nil {
		woduo = woduo.SetProjectTypeID(*id)
	}
	return woduo
}

// SetProjectType sets the project_type edge to ProjectType.
func (woduo *WorkOrderDefinitionUpdateOne) SetProjectType(p *ProjectType) *WorkOrderDefinitionUpdateOne {
	return woduo.SetProjectTypeID(p.ID)
}

// ClearType clears the type edge to WorkOrderType.
func (woduo *WorkOrderDefinitionUpdateOne) ClearType() *WorkOrderDefinitionUpdateOne {
	woduo.mutation.ClearType()
	return woduo
}

// ClearProjectType clears the project_type edge to ProjectType.
func (woduo *WorkOrderDefinitionUpdateOne) ClearProjectType() *WorkOrderDefinitionUpdateOne {
	woduo.mutation.ClearProjectType()
	return woduo
}

// Save executes the query and returns the updated entity.
func (woduo *WorkOrderDefinitionUpdateOne) Save(ctx context.Context) (*WorkOrderDefinition, error) {
	if _, ok := woduo.mutation.UpdateTime(); !ok {
		v := workorderdefinition.UpdateDefaultUpdateTime()
		woduo.mutation.SetUpdateTime(v)
	}

	var (
		err  error
		node *WorkOrderDefinition
	)
	if len(woduo.hooks) == 0 {
		node, err = woduo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*WorkOrderDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			woduo.mutation = mutation
			node, err = woduo.sqlSave(ctx)
			return node, err
		})
		for i := len(woduo.hooks) - 1; i >= 0; i-- {
			mut = woduo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, woduo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (woduo *WorkOrderDefinitionUpdateOne) SaveX(ctx context.Context) *WorkOrderDefinition {
	wod, err := woduo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return wod
}

// Exec executes the query on the entity.
func (woduo *WorkOrderDefinitionUpdateOne) Exec(ctx context.Context) error {
	_, err := woduo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (woduo *WorkOrderDefinitionUpdateOne) ExecX(ctx context.Context) {
	if err := woduo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (woduo *WorkOrderDefinitionUpdateOne) sqlSave(ctx context.Context) (wod *WorkOrderDefinition, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   workorderdefinition.Table,
			Columns: workorderdefinition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: workorderdefinition.FieldID,
			},
		},
	}
	id, ok := woduo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing WorkOrderDefinition.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := woduo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: workorderdefinition.FieldUpdateTime,
		})
	}
	if value, ok := woduo.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: workorderdefinition.FieldIndex,
		})
	}
	if value, ok := woduo.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: workorderdefinition.FieldIndex,
		})
	}
	if woduo.mutation.IndexCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: workorderdefinition.FieldIndex,
		})
	}
	if woduo.mutation.TypeCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := woduo.mutation.TypeIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if woduo.mutation.ProjectTypeCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := woduo.mutation.ProjectTypeIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	wod = &WorkOrderDefinition{config: woduo.config}
	_spec.Assign = wod.assignValues
	_spec.ScanValues = wod.scanValues()
	if err = sqlgraph.UpdateNode(ctx, woduo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{workorderdefinition.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return wod, nil
}
