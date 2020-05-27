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
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// CommentCreate is the builder for creating a Comment entity.
type CommentCreate struct {
	config
	mutation *CommentMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (cc *CommentCreate) SetCreateTime(t time.Time) *CommentCreate {
	cc.mutation.SetCreateTime(t)
	return cc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (cc *CommentCreate) SetNillableCreateTime(t *time.Time) *CommentCreate {
	if t != nil {
		cc.SetCreateTime(*t)
	}
	return cc
}

// SetUpdateTime sets the update_time field.
func (cc *CommentCreate) SetUpdateTime(t time.Time) *CommentCreate {
	cc.mutation.SetUpdateTime(t)
	return cc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (cc *CommentCreate) SetNillableUpdateTime(t *time.Time) *CommentCreate {
	if t != nil {
		cc.SetUpdateTime(*t)
	}
	return cc
}

// SetText sets the text field.
func (cc *CommentCreate) SetText(s string) *CommentCreate {
	cc.mutation.SetText(s)
	return cc
}

// SetAuthorID sets the author edge to User by id.
func (cc *CommentCreate) SetAuthorID(id int) *CommentCreate {
	cc.mutation.SetAuthorID(id)
	return cc
}

// SetAuthor sets the author edge to User.
func (cc *CommentCreate) SetAuthor(u *User) *CommentCreate {
	return cc.SetAuthorID(u.ID)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (cc *CommentCreate) SetWorkOrderID(id int) *CommentCreate {
	cc.mutation.SetWorkOrderID(id)
	return cc
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (cc *CommentCreate) SetNillableWorkOrderID(id *int) *CommentCreate {
	if id != nil {
		cc = cc.SetWorkOrderID(*id)
	}
	return cc
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (cc *CommentCreate) SetWorkOrder(w *WorkOrder) *CommentCreate {
	return cc.SetWorkOrderID(w.ID)
}

// SetProjectID sets the project edge to Project by id.
func (cc *CommentCreate) SetProjectID(id int) *CommentCreate {
	cc.mutation.SetProjectID(id)
	return cc
}

// SetNillableProjectID sets the project edge to Project by id if the given value is not nil.
func (cc *CommentCreate) SetNillableProjectID(id *int) *CommentCreate {
	if id != nil {
		cc = cc.SetProjectID(*id)
	}
	return cc
}

// SetProject sets the project edge to Project.
func (cc *CommentCreate) SetProject(p *Project) *CommentCreate {
	return cc.SetProjectID(p.ID)
}

// Save creates the Comment in the database.
func (cc *CommentCreate) Save(ctx context.Context) (*Comment, error) {
	if _, ok := cc.mutation.CreateTime(); !ok {
		v := comment.DefaultCreateTime()
		cc.mutation.SetCreateTime(v)
	}
	if _, ok := cc.mutation.UpdateTime(); !ok {
		v := comment.DefaultUpdateTime()
		cc.mutation.SetUpdateTime(v)
	}
	if _, ok := cc.mutation.Text(); !ok {
		return nil, errors.New("ent: missing required field \"text\"")
	}
	if _, ok := cc.mutation.AuthorID(); !ok {
		return nil, errors.New("ent: missing required edge \"author\"")
	}
	var (
		err  error
		node *Comment
	)
	if len(cc.hooks) == 0 {
		node, err = cc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CommentMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			cc.mutation = mutation
			node, err = cc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(cc.hooks) - 1; i >= 0; i-- {
			mut = cc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, cc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (cc *CommentCreate) SaveX(ctx context.Context) *Comment {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (cc *CommentCreate) sqlSave(ctx context.Context) (*Comment, error) {
	var (
		c     = &Comment{config: cc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: comment.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: comment.FieldID,
			},
		}
	)
	if value, ok := cc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: comment.FieldCreateTime,
		})
		c.CreateTime = value
	}
	if value, ok := cc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: comment.FieldUpdateTime,
		})
		c.UpdateTime = value
	}
	if value, ok := cc.mutation.Text(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: comment.FieldText,
		})
		c.Text = value
	}
	if nodes := cc.mutation.AuthorIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   comment.AuthorTable,
			Columns: []string{comment.AuthorColumn},
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
	if nodes := cc.mutation.WorkOrderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   comment.WorkOrderTable,
			Columns: []string{comment.WorkOrderColumn},
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
	if nodes := cc.mutation.ProjectIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   comment.ProjectTable,
			Columns: []string{comment.ProjectColumn},
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
	if err := sqlgraph.CreateNode(ctx, cc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	c.ID = int(id)
	return c, nil
}
