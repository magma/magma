// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package todo

import (
	"context"
	"strconv"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/handler"
	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/pkg/graphql/relay/internal/todo/ent"
	"github.com/facebookincubator/symphony/pkg/testdb"
	"github.com/stretchr/testify/suite"
)

type todoTestSuite struct {
	suite.Suite
	*client.Client
}

const maxTodos = 32

func (s *todoTestSuite) SetupTest() {
	db, name, err := testdb.Open()
	s.Require().NoError(err)
	db.SetMaxOpenConns(1)

	ec := ent.NewClient(ent.Driver(sql.OpenDB(name, db)))
	err = ec.Schema.Create(context.Background())
	s.Require().NoError(err)

	s.Client = client.New(handler.GraphQL(
		NewExecutableSchema(New(ec)),
	))

	var rsp struct {
		CreateTodo struct {
			ID string
		}
	}
	for i := 1; i <= maxTodos; i++ {
		id := strconv.Itoa(i)
		err := s.Post(
			`mutation($text: String!) { createTodo(todo:{text: $text}) { id } }`,
			&rsp, client.Var("text", id),
		)
		s.Require().NoError(err)
		s.Require().Equal(id, rsp.CreateTodo.ID)
	}
}

func TestTodo(t *testing.T) {
	suite.Run(t, &todoTestSuite{})
}

type response struct {
	Todos struct {
		Edges []struct {
			Node struct {
				ID string
			}
			Cursor string
		}
		PageInfo struct {
			HasNextPage     bool
			HasPreviousPage bool
			StartCursor     string
			EndCursor       string
		}
	}
}

func (s *todoTestSuite) TestQueryAll() {
	var rsp response
	err := s.Post(`query {
		todos {
			edges {
				node {
					id
				}
				cursor
			}
			pageInfo {
				hasNextPage
				hasPreviousPage
				startCursor
				endCursor
			}
		}
	}`, &rsp)
	s.Require().NoError(err)

	s.Require().Len(rsp.Todos.Edges, maxTodos)
	s.Assert().False(rsp.Todos.PageInfo.HasNextPage)
	s.Assert().False(rsp.Todos.PageInfo.HasPreviousPage)
	s.Assert().Equal(
		rsp.Todos.Edges[0].Cursor,
		rsp.Todos.PageInfo.StartCursor,
	)
	s.Assert().Equal(
		rsp.Todos.Edges[len(rsp.Todos.Edges)-1].Cursor,
		rsp.Todos.PageInfo.EndCursor,
	)
	for i, edge := range rsp.Todos.Edges {
		s.Assert().Equal(strconv.Itoa(i+1), edge.Node.ID)
		s.Assert().NotEmpty(edge.Cursor)
	}
}

func (s *todoTestSuite) TestPageForward() {
	const (
		query = `query($after: Cursor, $first: Int) {
			todos(after: $after, first: $first) {
				edges {
					node {
						id
					}
					cursor
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}`
		first = 5
	)
	var (
		rsp response
		id  = 1
	)
	for i := 0; i < maxTodos/first; i++ {
		err := s.Post(query, &rsp,
			client.Var("after", func() interface{} {
				if i > 0 {
					return rsp.Todos.PageInfo.EndCursor
				}
				return nil
			}()),
			client.Var("first", first),
		)
		s.Require().NoError(err)
		s.Require().Len(rsp.Todos.Edges, first)
		s.Assert().True(rsp.Todos.PageInfo.HasNextPage)
		s.Assert().NotEmpty(rsp.Todos.PageInfo.EndCursor)

		for _, edge := range rsp.Todos.Edges {
			s.Assert().Equal(strconv.Itoa(id), edge.Node.ID)
			s.Assert().NotEmpty(edge.Cursor)
			id++
		}
	}

	err := s.Post(query, &rsp,
		client.Var("after", rsp.Todos.PageInfo.EndCursor),
		client.Var("first", first),
	)
	s.Require().NoError(err)
	s.Require().NotEmpty(rsp.Todos.Edges)
	s.Assert().Len(rsp.Todos.Edges, maxTodos%first)
	s.Assert().False(rsp.Todos.PageInfo.HasNextPage)
	s.Assert().NotEmpty(rsp.Todos.PageInfo.EndCursor)

	for _, edge := range rsp.Todos.Edges {
		s.Assert().Equal(strconv.Itoa(id), edge.Node.ID)
		s.Assert().NotEmpty(edge.Cursor)
		id++
	}

	after := rsp.Todos.PageInfo.EndCursor
	rsp = response{}
	err = s.Post(query, &rsp,
		client.Var("after", after),
		client.Var("first", first),
	)
	s.Require().NoError(err)
	s.Assert().Empty(rsp.Todos.Edges)
	s.Assert().Empty(rsp.Todos.PageInfo.EndCursor)
	s.Assert().False(rsp.Todos.PageInfo.HasNextPage)
}

func (s *todoTestSuite) TestPageBackwards() {
	const (
		query = `query($before: Cursor, $last: Int) {
			todos(before: $before, last: $last) {
				edges {
					node {
						id
					}
					cursor
				}
				pageInfo {
					hasPreviousPage
					startCursor
				}
			}
		}`
		last = 7
	)
	var (
		rsp response
		id  = maxTodos
	)
	for i := 0; i < maxTodos/last; i++ {
		err := s.Post(query, &rsp,
			client.Var("before", func() interface{} {
				if i > 0 {
					return rsp.Todos.PageInfo.StartCursor
				}
				return nil
			}()),
			client.Var("last", last),
		)
		s.Require().NoError(err)
		s.Require().Len(rsp.Todos.Edges, last)
		s.Assert().True(rsp.Todos.PageInfo.HasPreviousPage)
		s.Assert().NotEmpty(rsp.Todos.PageInfo.StartCursor)

		for i := len(rsp.Todos.Edges) - 1; i >= 0; i-- {
			edge := &rsp.Todos.Edges[i]
			s.Assert().Equal(strconv.Itoa(id), edge.Node.ID)
			s.Assert().NotEmpty(edge.Cursor)
			id--
		}
	}

	err := s.Post(query, &rsp,
		client.Var("before", rsp.Todos.PageInfo.StartCursor),
		client.Var("last", last),
	)
	s.Require().NoError(err)
	s.Require().NotEmpty(rsp.Todos.Edges)
	s.Assert().Len(rsp.Todos.Edges, maxTodos%last)
	s.Assert().False(rsp.Todos.PageInfo.HasPreviousPage)
	s.Assert().NotEmpty(rsp.Todos.PageInfo.StartCursor)

	for i := len(rsp.Todos.Edges) - 1; i >= 0; i-- {
		edge := &rsp.Todos.Edges[i]
		s.Assert().Equal(strconv.Itoa(id), edge.Node.ID)
		s.Assert().NotEmpty(edge.Cursor)
		id--
	}
	s.Assert().Zero(id)

	before := rsp.Todos.PageInfo.StartCursor
	rsp = response{}
	err = s.Post(query, &rsp,
		client.Var("before", before),
		client.Var("last", last),
	)
	s.Require().NoError(err)
	s.Assert().Empty(rsp.Todos.Edges)
	s.Assert().Empty(rsp.Todos.PageInfo.StartCursor)
	s.Assert().False(rsp.Todos.PageInfo.HasPreviousPage)
}
