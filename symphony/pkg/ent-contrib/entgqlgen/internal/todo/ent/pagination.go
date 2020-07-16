// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/facebookincubator/symphony/pkg/ent-contrib/entgqlgen/internal/todo/ent/todo"
	"github.com/vektah/gqlparser/v2/gqlerror"
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
	ID int `json:"id"`
}

// MarshalGQL implements graphql.Marshaler interface.
func (c Cursor) MarshalGQL(w io.Writer) {
	const quote = '"'
	switch w := w.(type) {
	case io.ByteWriter:
		w.WriteByte(quote)
		defer w.WriteByte(quote)
	default:
		w.Write([]byte{quote})
		defer w.Write([]byte{quote})
	}
	wc := base64.NewEncoder(base64.StdEncoding, w)
	defer wc.Close()
	_ = json.NewEncoder(wc).Encode(c)
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (c *Cursor) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("%T is not a string", v)
	}
	if err := json.NewDecoder(
		base64.NewDecoder(
			base64.StdEncoding,
			strings.NewReader(s),
		),
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
	Edges      []*TodoEdge `json:"edges"`
	PageInfo   PageInfo    `json:"pageInfo"`
	TotalCount int         `json:"totalCount"`
}

// Paginate executes the query and returns a relay based cursor connection to Todo.
func (t *TodoQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*TodoConnection, error) {
	if first != nil && last != nil {
		return nil, gqlerror.Errorf("Passing both `first` and `last` to paginate a connection is not supported.")
	}

	conn := &TodoConnection{Edges: []*TodoEdge{}}
	if first != nil {
		if *first == 0 {
			return conn, nil
		} else if *first < 0 {
			return nil, gqlerror.Errorf("`first` on a connection cannot be less than zero.")
		}
	}
	if last != nil {
		if *last == 0 {
			return conn, nil
		} else if *last < 0 {
			return nil, gqlerror.Errorf("`last` on a connection cannot be less than zero.")
		}
	}

	if field := fieldForPath(ctx, "totalCount"); field != nil {
		count, err := t.Clone().Count(ctx)
		if err != nil {
			return nil, err
		}
		conn.TotalCount = count
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
		return conn, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

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

	return conn, nil
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
