// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package todo

import (
	"context"
	"errors"

	"github.com/facebookincubator/symphony/pkg/ent-contrib/entgqlgen/internal/todo/ent"
)

func New(client *ent.Client) Config {
	return Config{
		Resolvers: &resolvers{
			client: client,
		},
	}
}

type resolvers struct{ client *ent.Client }

func (r *resolvers) Node(ctx context.Context, id int) (ent.Noder, error) {
	node, err := r.client.Noder(ctx, id)
	if err == nil {
		return node, nil
	}
	var e *ent.NotFoundError
	if errors.As(err, &e) {
		err = nil
	}
	return nil, err
}

func (r *resolvers) Todos(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int) (*ent.TodoConnection, error) {
	return r.client.Todo.Query().
		Paginate(ctx, after, first, before, last)
}

func (r *resolvers) CreateTodo(ctx context.Context, todo TodoInput) (*ent.Todo, error) {
	return r.client.Todo.
		Create().
		SetText(todo.Text).
		Save(ctx)
}

func (r *resolvers) ClearTodos(ctx context.Context) (int, error) {
	return r.client.Todo.
		Delete().
		Exec(ctx)
}

func (r *resolvers) Query() QueryResolver       { return r }
func (r *resolvers) Mutation() MutationResolver { return r }
