// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/facebookincubator/symphony/pkg/ent-contrib/entgqlgen/internal/todo/ent/todo"
	"github.com/ugorji/go/codec"
)

// PageInfo of a connection type.
type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *Cursor `json:"startCursor"`
	EndCursor       *Cursor `json:"endCursor"`
}

// Cursor of an edge type.
type Cursor struct {
	ID int
}

// ErrInvalidPagination error is returned when paginating with invalid parameters.
var ErrInvalidPagination = errors.New("ent: invalid pagination parameters")

var quote = []byte(`"`)

// MarshalGQL implements graphql.Marshaler interface.
func (c Cursor) MarshalGQL(w io.Writer) {
	w.Write(quote)
	defer w.Write(quote)
	wc := base64.NewEncoder(base64.StdEncoding, w)
	defer wc.Close()
	_ = codec.NewEncoder(wc, &codec.MsgpackHandle{}).Encode(c)
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (c *Cursor) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("%T is not a string", v)
	}
	if err := codec.NewDecoder(
		base64.NewDecoder(
			base64.StdEncoding,
			strings.NewReader(s),
		),
		&codec.MsgpackHandle{},
	).Decode(c); err != nil {
		return fmt.Errorf("decode cursor: %w", err)
	}
	return nil
}

// TodoEdge is the edge representation of Todo.
type TodoEdge struct {
	Node   *Todo  `json:"node"`
	Cursor Cursor `json:"cursor"`
}

// TodoConnection is the connection containing edges to Todo.
type TodoConnection struct {
	Edges    []*TodoEdge `json:"edges"`
	PageInfo PageInfo    `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to Todo.
func (t *TodoQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*TodoConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &TodoConnection{
				Edges: []*TodoEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &TodoConnection{
				Edges: []*TodoEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		t = t.Where(todo.IDGT(after.ID))
	}
	if before != nil {
		t = t.Where(todo.IDLT(before.ID))
	}
	if first != nil {
		t = t.Order(Asc(todo.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		t = t.Order(Desc(todo.FieldID)).Limit(*last + 1)
	}
	t = t.collectConnectionFields(ctx)

	nodes, err := t.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &TodoConnection{
			Edges: []*TodoEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn TodoConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*TodoEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &TodoEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (t *TodoQuery) collectConnectionFields(ctx context.Context) *TodoQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		t = t.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return t
}

func fieldForPath(ctx context.Context, path ...string) *graphql.CollectedField {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return nil
	}
	oc := graphql.GetOperationContext(ctx)
	field := fc.Field

walk:
	for _, name := range path {
		for _, f := range graphql.CollectFields(oc, field.Selections, nil) {
			if f.Name == name {
				field = f
				continue walk
			}
		}
		return nil
	}
	return &field
}
