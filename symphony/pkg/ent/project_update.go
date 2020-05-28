// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/ent/comment"
	"github.com/facebookincubator/symphony/pkg/ent/location"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/project"
	"github.com/facebookincubator/symphony/pkg/ent/projecttype"
	"github.com/facebookincubator/symphony/pkg/ent/property"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/ent/workorder"
)

// ProjectUpdate is the builder for updating Project entities.
type ProjectUpdate struct {
	config
	hooks      []Hook
	mutation   *ProjectMutation
	predicates []predicate.Project
}

// Where adds a new predicate for the builder.
func (pu *ProjectUpdate) Where(ps ...predicate.Project) *ProjectUpdate {
	pu.predicates = append(pu.predicates, ps...)
	return pu
}

// SetName sets the name field.
func (pu *ProjectUpdate) SetName(s string) *ProjectUpdate {
	pu.mutation.SetName(s)
	return pu
}

// SetDescription sets the description field.
func (pu *ProjectUpdate) SetDescription(s string) *ProjectUpdate {
	pu.mutation.SetDescription(s)
	return pu
}

// SetNillableDescription sets the description field if the given value is not nil.
func (pu *ProjectUpdate) SetNillableDescription(s *string) *ProjectUpdate {
	if s != nil {
		pu.SetDescription(*s)
	}
	return pu
}

// ClearDescription clears the value of description.
func (pu *ProjectUpdate) ClearDescription() *ProjectUpdate {
	pu.mutation.ClearDescription()
	return pu
}

// SetTypeID sets the type edge to ProjectType by id.
func (pu *ProjectUpdate) SetTypeID(id int) *ProjectUpdate {
	pu.mutation.SetTypeID(id)
	return pu
}

// SetType sets the type edge to ProjectType.
func (pu *ProjectUpdate) SetType(p *ProjectType) *ProjectUpdate {
	return pu.SetTypeID(p.ID)
}

// SetLocationID sets the location edge to Location by id.
func (pu *ProjectUpdate) SetLocationID(id int) *ProjectUpdate {
	pu.mutation.SetLocationID(id)
	return pu
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (pu *ProjectUpdate) SetNillableLocationID(id *int) *ProjectUpdate {
	if id != nil {
		pu = pu.SetLocationID(*id)
	}
	return pu
}

// SetLocation sets the location edge to Location.
func (pu *ProjectUpdate) SetLocation(l *Location) *ProjectUpdate {
	return pu.SetLocationID(l.ID)
}

// AddCommentIDs adds the comments edge to Comment by ids.
func (pu *ProjectUpdate) AddCommentIDs(ids ...int) *ProjectUpdate {
	pu.mutation.AddCommentIDs(ids...)
	return pu
}

// AddComments adds the comments edges to Comment.
func (pu *ProjectUpdate) AddComments(c ...*Comment) *ProjectUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return pu.AddCommentIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (pu *ProjectUpdate) AddWorkOrderIDs(ids ...int) *ProjectUpdate {
	pu.mutation.AddWorkOrderIDs(ids...)
	return pu
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (pu *ProjectUpdate) AddWorkOrders(w ...*WorkOrder) *ProjectUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return pu.AddWorkOrderIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (pu *ProjectUpdate) AddPropertyIDs(ids ...int) *ProjectUpdate {
	pu.mutation.AddPropertyIDs(ids...)
	return pu
}

// AddProperties adds the properties edges to Property.
func (pu *ProjectUpdate) AddProperties(p ...*Property) *ProjectUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pu.AddPropertyIDs(ids...)
}

// SetCreatorID sets the creator edge to User by id.
func (pu *ProjectUpdate) SetCreatorID(id int) *ProjectUpdate {
	pu.mutation.SetCreatorID(id)
	return pu
}

// SetNillableCreatorID sets the creator edge to User by id if the given value is not nil.
func (pu *ProjectUpdate) SetNillableCreatorID(id *int) *ProjectUpdate {
	if id != nil {
		pu = pu.SetCreatorID(*id)
	}
	return pu
}

// SetCreator sets the creator edge to User.
func (pu *ProjectUpdate) SetCreator(u *User) *ProjectUpdate {
	return pu.SetCreatorID(u.ID)
}

// ClearType clears the type edge to ProjectType.
func (pu *ProjectUpdate) ClearType() *ProjectUpdate {
	pu.mutation.ClearType()
	return pu
}

// ClearLocation clears the location edge to Location.
func (pu *ProjectUpdate) ClearLocation() *ProjectUpdate {
	pu.mutation.ClearLocation()
	return pu
}

// RemoveCommentIDs removes the comments edge to Comment by ids.
func (pu *ProjectUpdate) RemoveCommentIDs(ids ...int) *ProjectUpdate {
	pu.mutation.RemoveCommentIDs(ids...)
	return pu
}

// RemoveComments removes comments edges to Comment.
func (pu *ProjectUpdate) RemoveComments(c ...*Comment) *ProjectUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return pu.RemoveCommentIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (pu *ProjectUpdate) RemoveWorkOrderIDs(ids ...int) *ProjectUpdate {
	pu.mutation.RemoveWorkOrderIDs(ids...)
	return pu
}

// RemoveWorkOrders removes work_orders edges to WorkOrder.
func (pu *ProjectUpdate) RemoveWorkOrders(w ...*WorkOrder) *ProjectUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return pu.RemoveWorkOrderIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (pu *ProjectUpdate) RemovePropertyIDs(ids ...int) *ProjectUpdate {
	pu.mutation.RemovePropertyIDs(ids...)
	return pu
}

// RemoveProperties removes properties edges to Property.
func (pu *ProjectUpdate) RemoveProperties(p ...*Property) *ProjectUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pu.RemovePropertyIDs(ids...)
}

// ClearCreator clears the creator edge to User.
func (pu *ProjectUpdate) ClearCreator() *ProjectUpdate {
	pu.mutation.ClearCreator()
	return pu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (pu *ProjectUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := pu.mutation.UpdateTime(); !ok {
		v := project.UpdateDefaultUpdateTime()
		pu.mutation.SetUpdateTime(v)
	}
	if v, ok := pu.mutation.Name(); ok {
		if err := project.NameValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}

	if _, ok := pu.mutation.TypeID(); pu.mutation.TypeCleared() && !ok {
		return 0, errors.New("ent: clearing a unique edge \"type\"")
	}

	var (
		err      error
		affected int
	)
	if len(pu.hooks) == 0 {
		affected, err = pu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ProjectMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			pu.mutation = mutation
			affected, err = pu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(pu.hooks) - 1; i >= 0; i-- {
			mut = pu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, pu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (pu *ProjectUpdate) SaveX(ctx context.Context) int {
	affected, err := pu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pu *ProjectUpdate) Exec(ctx context.Context) error {
	_, err := pu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pu *ProjectUpdate) ExecX(ctx context.Context) {
	if err := pu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (pu *ProjectUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   project.Table,
			Columns: project.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: project.FieldID,
			},
		},
	}
	if ps := pu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: project.FieldUpdateTime,
		})
	}
	if value, ok := pu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: project.FieldName,
		})
	}
	if value, ok := pu.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: project.FieldDescription,
		})
	}
	if pu.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: project.FieldDescription,
		})
	}
	if pu.mutation.TypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   project.TypeTable,
			Columns: []string{project.TypeColumn},
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
	if nodes := pu.mutation.TypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   project.TypeTable,
			Columns: []string{project.TypeColumn},
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
	if pu.mutation.LocationCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   project.LocationTable,
			Columns: []string{project.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.LocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   project.LocationTable,
			Columns: []string{project.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := pu.mutation.RemovedCommentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.CommentsTable,
			Columns: []string{project.CommentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: comment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.CommentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.CommentsTable,
			Columns: []string{project.CommentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: comment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := pu.mutation.RemovedWorkOrdersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.WorkOrdersTable,
			Columns: []string{project.WorkOrdersColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.WorkOrdersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.WorkOrdersTable,
			Columns: []string{project.WorkOrdersColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := pu.mutation.RemovedPropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.PropertiesTable,
			Columns: []string{project.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.PropertiesTable,
			Columns: []string{project.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pu.mutation.CreatorCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   project.CreatorTable,
			Columns: []string{project.CreatorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.CreatorIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   project.CreatorTable,
			Columns: []string{project.CreatorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, pu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{project.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// ProjectUpdateOne is the builder for updating a single Project entity.
type ProjectUpdateOne struct {
	config
	hooks    []Hook
	mutation *ProjectMutation
}

// SetName sets the name field.
func (puo *ProjectUpdateOne) SetName(s string) *ProjectUpdateOne {
	puo.mutation.SetName(s)
	return puo
}

// SetDescription sets the description field.
func (puo *ProjectUpdateOne) SetDescription(s string) *ProjectUpdateOne {
	puo.mutation.SetDescription(s)
	return puo
}

// SetNillableDescription sets the description field if the given value is not nil.
func (puo *ProjectUpdateOne) SetNillableDescription(s *string) *ProjectUpdateOne {
	if s != nil {
		puo.SetDescription(*s)
	}
	return puo
}

// ClearDescription clears the value of description.
func (puo *ProjectUpdateOne) ClearDescription() *ProjectUpdateOne {
	puo.mutation.ClearDescription()
	return puo
}

// SetTypeID sets the type edge to ProjectType by id.
func (puo *ProjectUpdateOne) SetTypeID(id int) *ProjectUpdateOne {
	puo.mutation.SetTypeID(id)
	return puo
}

// SetType sets the type edge to ProjectType.
func (puo *ProjectUpdateOne) SetType(p *ProjectType) *ProjectUpdateOne {
	return puo.SetTypeID(p.ID)
}

// SetLocationID sets the location edge to Location by id.
func (puo *ProjectUpdateOne) SetLocationID(id int) *ProjectUpdateOne {
	puo.mutation.SetLocationID(id)
	return puo
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (puo *ProjectUpdateOne) SetNillableLocationID(id *int) *ProjectUpdateOne {
	if id != nil {
		puo = puo.SetLocationID(*id)
	}
	return puo
}

// SetLocation sets the location edge to Location.
func (puo *ProjectUpdateOne) SetLocation(l *Location) *ProjectUpdateOne {
	return puo.SetLocationID(l.ID)
}

// AddCommentIDs adds the comments edge to Comment by ids.
func (puo *ProjectUpdateOne) AddCommentIDs(ids ...int) *ProjectUpdateOne {
	puo.mutation.AddCommentIDs(ids...)
	return puo
}

// AddComments adds the comments edges to Comment.
func (puo *ProjectUpdateOne) AddComments(c ...*Comment) *ProjectUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return puo.AddCommentIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (puo *ProjectUpdateOne) AddWorkOrderIDs(ids ...int) *ProjectUpdateOne {
	puo.mutation.AddWorkOrderIDs(ids...)
	return puo
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (puo *ProjectUpdateOne) AddWorkOrders(w ...*WorkOrder) *ProjectUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return puo.AddWorkOrderIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (puo *ProjectUpdateOne) AddPropertyIDs(ids ...int) *ProjectUpdateOne {
	puo.mutation.AddPropertyIDs(ids...)
	return puo
}

// AddProperties adds the properties edges to Property.
func (puo *ProjectUpdateOne) AddProperties(p ...*Property) *ProjectUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return puo.AddPropertyIDs(ids...)
}

// SetCreatorID sets the creator edge to User by id.
func (puo *ProjectUpdateOne) SetCreatorID(id int) *ProjectUpdateOne {
	puo.mutation.SetCreatorID(id)
	return puo
}

// SetNillableCreatorID sets the creator edge to User by id if the given value is not nil.
func (puo *ProjectUpdateOne) SetNillableCreatorID(id *int) *ProjectUpdateOne {
	if id != nil {
		puo = puo.SetCreatorID(*id)
	}
	return puo
}

// SetCreator sets the creator edge to User.
func (puo *ProjectUpdateOne) SetCreator(u *User) *ProjectUpdateOne {
	return puo.SetCreatorID(u.ID)
}

// ClearType clears the type edge to ProjectType.
func (puo *ProjectUpdateOne) ClearType() *ProjectUpdateOne {
	puo.mutation.ClearType()
	return puo
}

// ClearLocation clears the location edge to Location.
func (puo *ProjectUpdateOne) ClearLocation() *ProjectUpdateOne {
	puo.mutation.ClearLocation()
	return puo
}

// RemoveCommentIDs removes the comments edge to Comment by ids.
func (puo *ProjectUpdateOne) RemoveCommentIDs(ids ...int) *ProjectUpdateOne {
	puo.mutation.RemoveCommentIDs(ids...)
	return puo
}

// RemoveComments removes comments edges to Comment.
func (puo *ProjectUpdateOne) RemoveComments(c ...*Comment) *ProjectUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return puo.RemoveCommentIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (puo *ProjectUpdateOne) RemoveWorkOrderIDs(ids ...int) *ProjectUpdateOne {
	puo.mutation.RemoveWorkOrderIDs(ids...)
	return puo
}

// RemoveWorkOrders removes work_orders edges to WorkOrder.
func (puo *ProjectUpdateOne) RemoveWorkOrders(w ...*WorkOrder) *ProjectUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return puo.RemoveWorkOrderIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (puo *ProjectUpdateOne) RemovePropertyIDs(ids ...int) *ProjectUpdateOne {
	puo.mutation.RemovePropertyIDs(ids...)
	return puo
}

// RemoveProperties removes properties edges to Property.
func (puo *ProjectUpdateOne) RemoveProperties(p ...*Property) *ProjectUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return puo.RemovePropertyIDs(ids...)
}

// ClearCreator clears the creator edge to User.
func (puo *ProjectUpdateOne) ClearCreator() *ProjectUpdateOne {
	puo.mutation.ClearCreator()
	return puo
}

// Save executes the query and returns the updated entity.
func (puo *ProjectUpdateOne) Save(ctx context.Context) (*Project, error) {
	if _, ok := puo.mutation.UpdateTime(); !ok {
		v := project.UpdateDefaultUpdateTime()
		puo.mutation.SetUpdateTime(v)
	}
	if v, ok := puo.mutation.Name(); ok {
		if err := project.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}

	if _, ok := puo.mutation.TypeID(); puo.mutation.TypeCleared() && !ok {
		return nil, errors.New("ent: clearing a unique edge \"type\"")
	}

	var (
		err  error
		node *Project
	)
	if len(puo.hooks) == 0 {
		node, err = puo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ProjectMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			puo.mutation = mutation
			node, err = puo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(puo.hooks) - 1; i >= 0; i-- {
			mut = puo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, puo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (puo *ProjectUpdateOne) SaveX(ctx context.Context) *Project {
	pr, err := puo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return pr
}

// Exec executes the query on the entity.
func (puo *ProjectUpdateOne) Exec(ctx context.Context) error {
	_, err := puo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (puo *ProjectUpdateOne) ExecX(ctx context.Context) {
	if err := puo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (puo *ProjectUpdateOne) sqlSave(ctx context.Context) (pr *Project, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   project.Table,
			Columns: project.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: project.FieldID,
			},
		},
	}
	id, ok := puo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing Project.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := puo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: project.FieldUpdateTime,
		})
	}
	if value, ok := puo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: project.FieldName,
		})
	}
	if value, ok := puo.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: project.FieldDescription,
		})
	}
	if puo.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: project.FieldDescription,
		})
	}
	if puo.mutation.TypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   project.TypeTable,
			Columns: []string{project.TypeColumn},
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
	if nodes := puo.mutation.TypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   project.TypeTable,
			Columns: []string{project.TypeColumn},
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
	if puo.mutation.LocationCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   project.LocationTable,
			Columns: []string{project.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.LocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   project.LocationTable,
			Columns: []string{project.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := puo.mutation.RemovedCommentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.CommentsTable,
			Columns: []string{project.CommentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: comment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.CommentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.CommentsTable,
			Columns: []string{project.CommentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: comment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := puo.mutation.RemovedWorkOrdersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.WorkOrdersTable,
			Columns: []string{project.WorkOrdersColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.WorkOrdersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.WorkOrdersTable,
			Columns: []string{project.WorkOrdersColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := puo.mutation.RemovedPropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.PropertiesTable,
			Columns: []string{project.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.PropertiesTable,
			Columns: []string{project.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if puo.mutation.CreatorCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   project.CreatorTable,
			Columns: []string{project.CreatorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.CreatorIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   project.CreatorTable,
			Columns: []string{project.CreatorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	pr = &Project{config: puo.config}
	_spec.Assign = pr.assignValues
	_spec.ScanValues = pr.scanValues()
	if err = sqlgraph.UpdateNode(ctx, puo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{project.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return pr, nil
}
