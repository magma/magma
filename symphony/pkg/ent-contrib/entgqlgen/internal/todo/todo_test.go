// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package todo

import (
	"strconv"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/pkg/ent-contrib/entgqlgen/internal/todo/ent"
	"github.com/facebookincubator/symphony/pkg/ent-contrib/entgqlgen/internal/todo/ent/enttest"
	"github.com/facebookincubator/symphony/pkg/ent-contrib/entgqlgen/internal/todo/ent/migrate"
	"github.com/facebookincubator/symphony/pkg/testdb"
	"github.com/stretchr/testify/suite"
)

type todoTestSuite struct {
	suite.Suite
	*client.Client
}

const (
	queryAll = `query {
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
	}`
	maxTodos = 32
)

func (s *todoTestSuite) SetupTest() {
	db, name, err := testdb.Open()
	s.Require().NoError(err)
	db.SetMaxOpenConns(1)

	ec := enttest.NewClient(s.T(),
		enttest.WithOptions(ent.Driver(sql.OpenDB(name, db))),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	s.Client = client.New(
		handler.NewDefaultServer(
			NewExecutableSchema(New(ec)),
		),
	)

	var (
		rsp struct {
			CreateTodo struct {
				ID string
			}
		}
		root = 1
	)
	for i := 1; i <= maxTodos; i++ {
		id := strconv.Itoa(i)
		err := s.Post(
			`mutation($text: String!, $parent: ID) { createTodo(todo:{text: $text, parent: $parent}) { id } }`,
			&rsp, client.Var("text", id), client.Var("parent", func() *int {
				if i == root {
					return nil
				}
				if i%2 != 0 {
					return pointer.ToInt(i - 2)
				}
				return &root
			}()),
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
			StartCursor     *string
			EndCursor       *string
		}
	}
}

func (s *todoTestSuite) TestQueryEmpty() {
	{
		var rsp struct{ ClearTodos int }
		err := s.Post(`mutation { clearTodos }`, &rsp)
		s.Require().NoError(err)
		s.Require().Equal(maxTodos, rsp.ClearTodos)
	}
	var rsp response
	err := s.Post(queryAll, &rsp)
	s.Require().NoError(err)
	s.Assert().Empty(rsp.Todos.Edges)
	s.Assert().False(rsp.Todos.PageInfo.HasNextPage)
	s.Assert().False(rsp.Todos.PageInfo.HasPreviousPage)
	s.Assert().Nil(rsp.Todos.PageInfo.StartCursor)
	s.Assert().Nil(rsp.Todos.PageInfo.EndCursor)
}

func (s *todoTestSuite) TestQueryAll() {
	var rsp response
	err := s.Post(queryAll, &rsp)
	s.Require().NoError(err)

	s.Require().Len(rsp.Todos.Edges, maxTodos)
	s.Assert().False(rsp.Todos.PageInfo.HasNextPage)
	s.Assert().False(rsp.Todos.PageInfo.HasPreviousPage)
	s.Assert().Equal(
		rsp.Todos.Edges[0].Cursor,
		*rsp.Todos.PageInfo.StartCursor,
	)
	s.Assert().Equal(
		rsp.Todos.Edges[len(rsp.Todos.Edges)-1].Cursor,
		*rsp.Todos.PageInfo.EndCursor,
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
		after interface{}
		rsp   response
		id    = 1
	)
	for i := 0; i < maxTodos/first; i++ {
		err := s.Post(query, &rsp,
			client.Var("after", after),
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
		after = rsp.Todos.PageInfo.EndCursor
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

	after = rsp.Todos.PageInfo.EndCursor
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
		before interface{}
		rsp    response
		id     = maxTodos
	)
	for i := 0; i < maxTodos/last; i++ {
		err := s.Post(query, &rsp,
			client.Var("before", before),
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
		before = rsp.Todos.PageInfo.StartCursor
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

	before = rsp.Todos.PageInfo.StartCursor
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

func (s *todoTestSuite) TestNodeCollection() {
	const (
		query = `query($id: ID!) {
			todo: node(id: $id) {
				... on Todo {
					parent {
						text
						parent {
							text
						}
					}
					children {
						text
						children {
							text
						}
					}
				}
			}
		}`
	)
	var rsp struct {
		Todo struct {
			Parent *struct {
				Text   string
				Parent *struct {
					Text string
				}
			}
			Children []struct {
				Text     string
				Children []struct {
					Text string
				}
			}
		}
	}
	err := s.Post(query, &rsp, client.Var("id", 1))
	s.Require().NoError(err)
	s.Assert().Nil(rsp.Todo.Parent)
	s.Assert().Len(rsp.Todo.Children, maxTodos/2+1)
	s.Assert().Condition(func() bool {
		for _, child := range rsp.Todo.Children {
			if child.Text == "3" {
				s.Require().Len(child.Children, 1)
				s.Assert().Equal("5", child.Children[0].Text)
				return true
			}
		}
		return false
	})

	err = s.Post(query, &rsp, client.Var("id", 4))
	s.Require().NoError(err)
	s.Require().NotNil(rsp.Todo.Parent)
	s.Assert().Equal("1", rsp.Todo.Parent.Text)
	s.Assert().Empty(rsp.Todo.Children)

	err = s.Post(query, &rsp, client.Var("id", 5))
	s.Require().NoError(err)
	s.Require().NotNil(rsp.Todo.Parent)
	s.Assert().Equal("3", rsp.Todo.Parent.Text)
	s.Require().NotNil(rsp.Todo.Parent.Parent)
	s.Assert().Equal("1", rsp.Todo.Parent.Parent.Text)
	s.Require().Len(rsp.Todo.Children, 1)
	s.Assert().Equal("7", rsp.Todo.Children[0].Text)
}

func (s *todoTestSuite) TestConnCollection() {
	const (
		query = `query {
			todos {
				edges {
					node {
						id
						parent {
							id
						}
						children {
							id
						}
					}
				}
			}
		}`
	)
	var rsp struct {
		Todos struct {
			Edges []struct {
				Node struct {
					ID     string
					Parent *struct {
						ID string
					}
					Children []struct {
						ID string
					}
				}
			}
		}
	}

	err := s.Post(query, &rsp)
	s.Require().NoError(err)
	s.Require().Len(rsp.Todos.Edges, maxTodos)

	for i, edge := range rsp.Todos.Edges {
		switch {
		case i == 0:
			s.Assert().Nil(edge.Node.Parent)
			s.Assert().Len(edge.Node.Children, maxTodos/2+1)
		case i%2 == 0:
			s.Require().NotNil(edge.Node.Parent)
			id, err := strconv.Atoi(edge.Node.Parent.ID)
			s.Require().NoError(err)
			s.Assert().Equal(i-1, id)
			if i < len(rsp.Todos.Edges)-2 {
				s.Assert().Len(edge.Node.Children, 1)
			} else {
				s.Assert().Empty(edge.Node.Children)
			}
		case i%2 != 0:
			s.Require().NotNil(edge.Node.Parent)
			s.Assert().Equal("1", edge.Node.Parent.ID)
			s.Assert().Empty(edge.Node.Children)
		}
	}
}
