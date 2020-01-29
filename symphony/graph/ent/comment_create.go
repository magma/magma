// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/comment"
)

// CommentCreate is the builder for creating a Comment entity.
type CommentCreate struct {
	config
	create_time *time.Time
	update_time *time.Time
	author_name *string
	text        *string
}

// SetCreateTime sets the create_time field.
func (cc *CommentCreate) SetCreateTime(t time.Time) *CommentCreate {
	cc.create_time = &t
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
	cc.update_time = &t
	return cc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (cc *CommentCreate) SetNillableUpdateTime(t *time.Time) *CommentCreate {
	if t != nil {
		cc.SetUpdateTime(*t)
	}
	return cc
}

// SetAuthorName sets the author_name field.
func (cc *CommentCreate) SetAuthorName(s string) *CommentCreate {
	cc.author_name = &s
	return cc
}

// SetText sets the text field.
func (cc *CommentCreate) SetText(s string) *CommentCreate {
	cc.text = &s
	return cc
}

// Save creates the Comment in the database.
func (cc *CommentCreate) Save(ctx context.Context) (*Comment, error) {
	if cc.create_time == nil {
		v := comment.DefaultCreateTime()
		cc.create_time = &v
	}
	if cc.update_time == nil {
		v := comment.DefaultUpdateTime()
		cc.update_time = &v
	}
	if cc.author_name == nil {
		return nil, errors.New("ent: missing required field \"author_name\"")
	}
	if cc.text == nil {
		return nil, errors.New("ent: missing required field \"text\"")
	}
	return cc.sqlSave(ctx)
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
				Type:   field.TypeString,
				Column: comment.FieldID,
			},
		}
	)
	if value := cc.create_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: comment.FieldCreateTime,
		})
		c.CreateTime = *value
	}
	if value := cc.update_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: comment.FieldUpdateTime,
		})
		c.UpdateTime = *value
	}
	if value := cc.author_name; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: comment.FieldAuthorName,
		})
		c.AuthorName = *value
	}
	if value := cc.text; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: comment.FieldText,
		})
		c.Text = *value
	}
	if err := sqlgraph.CreateNode(ctx, cc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	c.ID = strconv.FormatInt(id, 10)
	return c, nil
}
