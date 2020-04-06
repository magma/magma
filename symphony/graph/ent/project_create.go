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
	"github.com/facebookincubator/symphony/graph/ent/comment"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// ProjectCreate is the builder for creating a Project entity.
type ProjectCreate struct {
	config
	mutation *ProjectMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (pc *ProjectCreate) SetCreateTime(t time.Time) *ProjectCreate {
	pc.mutation.SetCreateTime(t)
	return pc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (pc *ProjectCreate) SetNillableCreateTime(t *time.Time) *ProjectCreate {
	if t != nil {
		pc.SetCreateTime(*t)
	}
	return pc
}

// SetUpdateTime sets the update_time field.
func (pc *ProjectCreate) SetUpdateTime(t time.Time) *ProjectCreate {
	pc.mutation.SetUpdateTime(t)
	return pc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (pc *ProjectCreate) SetNillableUpdateTime(t *time.Time) *ProjectCreate {
	if t != nil {
		pc.SetUpdateTime(*t)
	}
	return pc
}

// SetName sets the name field.
func (pc *ProjectCreate) SetName(s string) *ProjectCreate {
	pc.mutation.SetName(s)
	return pc
}

// SetDescription sets the description field.
func (pc *ProjectCreate) SetDescription(s string) *ProjectCreate {
	pc.mutation.SetDescription(s)
	return pc
}

// SetNillableDescription sets the description field if the given value is not nil.
func (pc *ProjectCreate) SetNillableDescription(s *string) *ProjectCreate {
	if s != nil {
		pc.SetDescription(*s)
	}
	return pc
}

// SetTypeID sets the type edge to ProjectType by id.
func (pc *ProjectCreate) SetTypeID(id int) *ProjectCreate {
	pc.mutation.SetTypeID(id)
	return pc
}

// SetType sets the type edge to ProjectType.
func (pc *ProjectCreate) SetType(p *ProjectType) *ProjectCreate {
	return pc.SetTypeID(p.ID)
}

// SetLocationID sets the location edge to Location by id.
func (pc *ProjectCreate) SetLocationID(id int) *ProjectCreate {
	pc.mutation.SetLocationID(id)
	return pc
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (pc *ProjectCreate) SetNillableLocationID(id *int) *ProjectCreate {
	if id != nil {
		pc = pc.SetLocationID(*id)
	}
	return pc
}

// SetLocation sets the location edge to Location.
func (pc *ProjectCreate) SetLocation(l *Location) *ProjectCreate {
	return pc.SetLocationID(l.ID)
}

// AddCommentIDs adds the comments edge to Comment by ids.
func (pc *ProjectCreate) AddCommentIDs(ids ...int) *ProjectCreate {
	pc.mutation.AddCommentIDs(ids...)
	return pc
}

// AddComments adds the comments edges to Comment.
func (pc *ProjectCreate) AddComments(c ...*Comment) *ProjectCreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return pc.AddCommentIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (pc *ProjectCreate) AddWorkOrderIDs(ids ...int) *ProjectCreate {
	pc.mutation.AddWorkOrderIDs(ids...)
	return pc
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (pc *ProjectCreate) AddWorkOrders(w ...*WorkOrder) *ProjectCreate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return pc.AddWorkOrderIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (pc *ProjectCreate) AddPropertyIDs(ids ...int) *ProjectCreate {
	pc.mutation.AddPropertyIDs(ids...)
	return pc
}

// AddProperties adds the properties edges to Property.
func (pc *ProjectCreate) AddProperties(p ...*Property) *ProjectCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pc.AddPropertyIDs(ids...)
}

// SetCreatorID sets the creator edge to User by id.
func (pc *ProjectCreate) SetCreatorID(id int) *ProjectCreate {
	pc.mutation.SetCreatorID(id)
	return pc
}

// SetNillableCreatorID sets the creator edge to User by id if the given value is not nil.
func (pc *ProjectCreate) SetNillableCreatorID(id *int) *ProjectCreate {
	if id != nil {
		pc = pc.SetCreatorID(*id)
	}
	return pc
}

// SetCreator sets the creator edge to User.
func (pc *ProjectCreate) SetCreator(u *User) *ProjectCreate {
	return pc.SetCreatorID(u.ID)
}

// Save creates the Project in the database.
func (pc *ProjectCreate) Save(ctx context.Context) (*Project, error) {
	if _, ok := pc.mutation.CreateTime(); !ok {
		v := project.DefaultCreateTime()
		pc.mutation.SetCreateTime(v)
	}
	if _, ok := pc.mutation.UpdateTime(); !ok {
		v := project.DefaultUpdateTime()
		pc.mutation.SetUpdateTime(v)
	}
	if _, ok := pc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if v, ok := pc.mutation.Name(); ok {
		if err := project.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if _, ok := pc.mutation.TypeID(); !ok {
		return nil, errors.New("ent: missing required edge \"type\"")
	}
	var (
		err  error
		node *Project
	)
	if len(pc.hooks) == 0 {
		node, err = pc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ProjectMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			pc.mutation = mutation
			node, err = pc.sqlSave(ctx)
			return node, err
		})
		for i := len(pc.hooks); i > 0; i-- {
			mut = pc.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, pc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (pc *ProjectCreate) SaveX(ctx context.Context) *Project {
	v, err := pc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (pc *ProjectCreate) sqlSave(ctx context.Context) (*Project, error) {
	var (
		pr    = &Project{config: pc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: project.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: project.FieldID,
			},
		}
	)
	if value, ok := pc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: project.FieldCreateTime,
		})
		pr.CreateTime = value
	}
	if value, ok := pc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: project.FieldUpdateTime,
		})
		pr.UpdateTime = value
	}
	if value, ok := pc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: project.FieldName,
		})
		pr.Name = value
	}
	if value, ok := pc.mutation.Description(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: project.FieldDescription,
		})
		pr.Description = &value
	}
	if nodes := pc.mutation.TypeIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.LocationIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.CommentsIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.WorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.PropertiesIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.CreatorIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, pc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	pr.ID = int(id)
	return pr, nil
}
