// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/comment"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// CommentUpdate is the builder for updating Comment entities.
type CommentUpdate struct {
	config

	update_time *time.Time
	author_name *string
	text        *string
	predicates  []predicate.Comment
}

// Where adds a new predicate for the builder.
func (cu *CommentUpdate) Where(ps ...predicate.Comment) *CommentUpdate {
	cu.predicates = append(cu.predicates, ps...)
	return cu
}

// SetAuthorName sets the author_name field.
func (cu *CommentUpdate) SetAuthorName(s string) *CommentUpdate {
	cu.author_name = &s
	return cu
}

// SetText sets the text field.
func (cu *CommentUpdate) SetText(s string) *CommentUpdate {
	cu.text = &s
	return cu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (cu *CommentUpdate) Save(ctx context.Context) (int, error) {
	if cu.update_time == nil {
		v := comment.UpdateDefaultUpdateTime()
		cu.update_time = &v
	}
	return cu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (cu *CommentUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *CommentUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *CommentUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cu *CommentUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(cu.driver.Dialect())
		selector = builder.Select(comment.FieldID).From(builder.Table(comment.Table))
	)
	for _, p := range cu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = cu.driver.Query(ctx, query, args, rows); err != nil {
		return 0, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ent: failed reading id: %v", err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return 0, nil
	}

	tx, err := cu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(comment.Table)
	)
	updater = updater.Where(sql.InInts(comment.FieldID, ids...))
	if value := cu.update_time; value != nil {
		updater.Set(comment.FieldUpdateTime, *value)
	}
	if value := cu.author_name; value != nil {
		updater.Set(comment.FieldAuthorName, *value)
	}
	if value := cu.text; value != nil {
		updater.Set(comment.FieldText, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// CommentUpdateOne is the builder for updating a single Comment entity.
type CommentUpdateOne struct {
	config
	id string

	update_time *time.Time
	author_name *string
	text        *string
}

// SetAuthorName sets the author_name field.
func (cuo *CommentUpdateOne) SetAuthorName(s string) *CommentUpdateOne {
	cuo.author_name = &s
	return cuo
}

// SetText sets the text field.
func (cuo *CommentUpdateOne) SetText(s string) *CommentUpdateOne {
	cuo.text = &s
	return cuo
}

// Save executes the query and returns the updated entity.
func (cuo *CommentUpdateOne) Save(ctx context.Context) (*Comment, error) {
	if cuo.update_time == nil {
		v := comment.UpdateDefaultUpdateTime()
		cuo.update_time = &v
	}
	return cuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *CommentUpdateOne) SaveX(ctx context.Context) *Comment {
	c, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return c
}

// Exec executes the query on the entity.
func (cuo *CommentUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *CommentUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cuo *CommentUpdateOne) sqlSave(ctx context.Context) (c *Comment, err error) {
	var (
		builder  = sql.Dialect(cuo.driver.Dialect())
		selector = builder.Select(comment.Columns...).From(builder.Table(comment.Table))
	)
	comment.ID(cuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = cuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		c = &Comment{config: cuo.config}
		if err := c.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into Comment: %v", err)
		}
		id = c.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("Comment with id: %v", cuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one Comment with the same id: %v", cuo.id)
	}

	tx, err := cuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(comment.Table)
	)
	updater = updater.Where(sql.InInts(comment.FieldID, ids...))
	if value := cuo.update_time; value != nil {
		updater.Set(comment.FieldUpdateTime, *value)
		c.UpdateTime = *value
	}
	if value := cuo.author_name; value != nil {
		updater.Set(comment.FieldAuthorName, *value)
		c.AuthorName = *value
	}
	if value := cuo.text; value != nil {
		updater.Set(comment.FieldText, *value)
		c.Text = *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return c, nil
}
