// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"

	"github.com/facebookincubator/symphony/pkg/graphql/relay"
	"github.com/facebookincubator/symphony/pkg/graphql/relay/internal/todo/ent/todo"
)

// ErrInvalidPagination error is returned when paginating with invalid parameters.
var ErrInvalidPagination = errors.New("ent: invalid pagination parameters")

// TodoEdge is the edge representation of Todo.
type TodoEdge struct {
	Node   *Todo         `json:"node"`
	Cursor *relay.Cursor `json:"cursor"`
}

// TodoConnection is the connection containing edges to Todo.
type TodoConnection struct {
	Edges    []*TodoEdge     `json:"edges"`
	PageInfo *relay.PageInfo `json:"pageInfo"`
}

func newTodoConnection() *TodoConnection {
	return &TodoConnection{
		Edges:    []*TodoEdge{},
		PageInfo: &relay.PageInfo{},
	}
}

// Paginate executes the query and returns a relay based cursor connection to Todo.
func (t *TodoQuery) Paginate(ctx context.Context, after *relay.Cursor, first *int, before *relay.Cursor, last *int) (*TodoConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return newTodoConnection(), nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return newTodoConnection(), nil
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

	nodes, err := t.All(ctx)
	if err != nil || len(nodes) == 0 {
		return newTodoConnection(), err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	info := &relay.PageInfo{}
	if first != nil && len(nodes) > *first {
		info.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		info.HasPreviousPage = true
		nodes = nodes[1:]
	}
	edges := make([]*TodoEdge, len(nodes))
	for i, node := range nodes {
		edges[i] = &TodoEdge{
			Node: node,
			Cursor: &relay.Cursor{
				ID: node.ID,
			},
		}
	}
	info.StartCursor = edges[0].Cursor
	info.EndCursor = edges[len(edges)-1].Cursor

	return &TodoConnection{
		Edges:    edges,
		PageInfo: info,
	}, nil
}
