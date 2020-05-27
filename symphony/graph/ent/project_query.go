// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"math"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/comment"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// ProjectQuery is the builder for querying Project entities.
type ProjectQuery struct {
	config
	limit      *int
	offset     *int
	order      []OrderFunc
	unique     []string
	predicates []predicate.Project
	// eager-loading edges.
	withType       *ProjectTypeQuery
	withLocation   *LocationQuery
	withComments   *CommentQuery
	withWorkOrders *WorkOrderQuery
	withProperties *PropertyQuery
	withCreator    *UserQuery
	withFKs        bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (pq *ProjectQuery) Where(ps ...predicate.Project) *ProjectQuery {
	pq.predicates = append(pq.predicates, ps...)
	return pq
}

// Limit adds a limit step to the query.
func (pq *ProjectQuery) Limit(limit int) *ProjectQuery {
	pq.limit = &limit
	return pq
}

// Offset adds an offset step to the query.
func (pq *ProjectQuery) Offset(offset int) *ProjectQuery {
	pq.offset = &offset
	return pq
}

// Order adds an order step to the query.
func (pq *ProjectQuery) Order(o ...OrderFunc) *ProjectQuery {
	pq.order = append(pq.order, o...)
	return pq
}

// QueryType chains the current query on the type edge.
func (pq *ProjectQuery) QueryType() *ProjectTypeQuery {
	query := &ProjectTypeQuery{config: pq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(project.Table, project.FieldID, pq.sqlQuery()),
			sqlgraph.To(projecttype.Table, projecttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, project.TypeTable, project.TypeColumn),
		)
		fromU = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryLocation chains the current query on the location edge.
func (pq *ProjectQuery) QueryLocation() *LocationQuery {
	query := &LocationQuery{config: pq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(project.Table, project.FieldID, pq.sqlQuery()),
			sqlgraph.To(location.Table, location.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, project.LocationTable, project.LocationColumn),
		)
		fromU = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryComments chains the current query on the comments edge.
func (pq *ProjectQuery) QueryComments() *CommentQuery {
	query := &CommentQuery{config: pq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(project.Table, project.FieldID, pq.sqlQuery()),
			sqlgraph.To(comment.Table, comment.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, project.CommentsTable, project.CommentsColumn),
		)
		fromU = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryWorkOrders chains the current query on the work_orders edge.
func (pq *ProjectQuery) QueryWorkOrders() *WorkOrderQuery {
	query := &WorkOrderQuery{config: pq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(project.Table, project.FieldID, pq.sqlQuery()),
			sqlgraph.To(workorder.Table, workorder.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, project.WorkOrdersTable, project.WorkOrdersColumn),
		)
		fromU = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryProperties chains the current query on the properties edge.
func (pq *ProjectQuery) QueryProperties() *PropertyQuery {
	query := &PropertyQuery{config: pq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(project.Table, project.FieldID, pq.sqlQuery()),
			sqlgraph.To(property.Table, property.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, project.PropertiesTable, project.PropertiesColumn),
		)
		fromU = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryCreator chains the current query on the creator edge.
func (pq *ProjectQuery) QueryCreator() *UserQuery {
	query := &UserQuery{config: pq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(project.Table, project.FieldID, pq.sqlQuery()),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, project.CreatorTable, project.CreatorColumn),
		)
		fromU = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Project entity in the query. Returns *NotFoundError when no project was found.
func (pq *ProjectQuery) First(ctx context.Context) (*Project, error) {
	prs, err := pq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(prs) == 0 {
		return nil, &NotFoundError{project.Label}
	}
	return prs[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (pq *ProjectQuery) FirstX(ctx context.Context) *Project {
	pr, err := pq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return pr
}

// FirstID returns the first Project id in the query. Returns *NotFoundError when no id was found.
func (pq *ProjectQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = pq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{project.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (pq *ProjectQuery) FirstXID(ctx context.Context) int {
	id, err := pq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only Project entity in the query, returns an error if not exactly one entity was returned.
func (pq *ProjectQuery) Only(ctx context.Context) (*Project, error) {
	prs, err := pq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(prs) {
	case 1:
		return prs[0], nil
	case 0:
		return nil, &NotFoundError{project.Label}
	default:
		return nil, &NotSingularError{project.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (pq *ProjectQuery) OnlyX(ctx context.Context) *Project {
	pr, err := pq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return pr
}

// OnlyID returns the only Project id in the query, returns an error if not exactly one id was returned.
func (pq *ProjectQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = pq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{project.Label}
	default:
		err = &NotSingularError{project.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (pq *ProjectQuery) OnlyXID(ctx context.Context) int {
	id, err := pq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Projects.
func (pq *ProjectQuery) All(ctx context.Context) ([]*Project, error) {
	if err := pq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return pq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (pq *ProjectQuery) AllX(ctx context.Context) []*Project {
	prs, err := pq.All(ctx)
	if err != nil {
		panic(err)
	}
	return prs
}

// IDs executes the query and returns a list of Project ids.
func (pq *ProjectQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := pq.Select(project.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (pq *ProjectQuery) IDsX(ctx context.Context) []int {
	ids, err := pq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (pq *ProjectQuery) Count(ctx context.Context) (int, error) {
	if err := pq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return pq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (pq *ProjectQuery) CountX(ctx context.Context) int {
	count, err := pq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (pq *ProjectQuery) Exist(ctx context.Context) (bool, error) {
	if err := pq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return pq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (pq *ProjectQuery) ExistX(ctx context.Context) bool {
	exist, err := pq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (pq *ProjectQuery) Clone() *ProjectQuery {
	return &ProjectQuery{
		config:     pq.config,
		limit:      pq.limit,
		offset:     pq.offset,
		order:      append([]OrderFunc{}, pq.order...),
		unique:     append([]string{}, pq.unique...),
		predicates: append([]predicate.Project{}, pq.predicates...),
		// clone intermediate query.
		sql:  pq.sql.Clone(),
		path: pq.path,
	}
}

//  WithType tells the query-builder to eager-loads the nodes that are connected to
// the "type" edge. The optional arguments used to configure the query builder of the edge.
func (pq *ProjectQuery) WithType(opts ...func(*ProjectTypeQuery)) *ProjectQuery {
	query := &ProjectTypeQuery{config: pq.config}
	for _, opt := range opts {
		opt(query)
	}
	pq.withType = query
	return pq
}

//  WithLocation tells the query-builder to eager-loads the nodes that are connected to
// the "location" edge. The optional arguments used to configure the query builder of the edge.
func (pq *ProjectQuery) WithLocation(opts ...func(*LocationQuery)) *ProjectQuery {
	query := &LocationQuery{config: pq.config}
	for _, opt := range opts {
		opt(query)
	}
	pq.withLocation = query
	return pq
}

//  WithComments tells the query-builder to eager-loads the nodes that are connected to
// the "comments" edge. The optional arguments used to configure the query builder of the edge.
func (pq *ProjectQuery) WithComments(opts ...func(*CommentQuery)) *ProjectQuery {
	query := &CommentQuery{config: pq.config}
	for _, opt := range opts {
		opt(query)
	}
	pq.withComments = query
	return pq
}

//  WithWorkOrders tells the query-builder to eager-loads the nodes that are connected to
// the "work_orders" edge. The optional arguments used to configure the query builder of the edge.
func (pq *ProjectQuery) WithWorkOrders(opts ...func(*WorkOrderQuery)) *ProjectQuery {
	query := &WorkOrderQuery{config: pq.config}
	for _, opt := range opts {
		opt(query)
	}
	pq.withWorkOrders = query
	return pq
}

//  WithProperties tells the query-builder to eager-loads the nodes that are connected to
// the "properties" edge. The optional arguments used to configure the query builder of the edge.
func (pq *ProjectQuery) WithProperties(opts ...func(*PropertyQuery)) *ProjectQuery {
	query := &PropertyQuery{config: pq.config}
	for _, opt := range opts {
		opt(query)
	}
	pq.withProperties = query
	return pq
}

//  WithCreator tells the query-builder to eager-loads the nodes that are connected to
// the "creator" edge. The optional arguments used to configure the query builder of the edge.
func (pq *ProjectQuery) WithCreator(opts ...func(*UserQuery)) *ProjectQuery {
	query := &UserQuery{config: pq.config}
	for _, opt := range opts {
		opt(query)
	}
	pq.withCreator = query
	return pq
}

// GroupBy used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		CreateTime time.Time `json:"create_time,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Project.Query().
//		GroupBy(project.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (pq *ProjectQuery) GroupBy(field string, fields ...string) *ProjectGroupBy {
	group := &ProjectGroupBy{config: pq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := pq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return pq.sqlQuery(), nil
	}
	return group
}

// Select one or more fields from the given query.
//
// Example:
//
//	var v []struct {
//		CreateTime time.Time `json:"create_time,omitempty"`
//	}
//
//	client.Project.Query().
//		Select(project.FieldCreateTime).
//		Scan(ctx, &v)
//
func (pq *ProjectQuery) Select(field string, fields ...string) *ProjectSelect {
	selector := &ProjectSelect{config: pq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := pq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return pq.sqlQuery(), nil
	}
	return selector
}

func (pq *ProjectQuery) prepareQuery(ctx context.Context) error {
	if pq.path != nil {
		prev, err := pq.path(ctx)
		if err != nil {
			return err
		}
		pq.sql = prev
	}
	if err := project.Policy.EvalQuery(ctx, pq); err != nil {
		return err
	}
	return nil
}

func (pq *ProjectQuery) sqlAll(ctx context.Context) ([]*Project, error) {
	var (
		nodes       = []*Project{}
		withFKs     = pq.withFKs
		_spec       = pq.querySpec()
		loadedTypes = [6]bool{
			pq.withType != nil,
			pq.withLocation != nil,
			pq.withComments != nil,
			pq.withWorkOrders != nil,
			pq.withProperties != nil,
			pq.withCreator != nil,
		}
	)
	if pq.withType != nil || pq.withLocation != nil || pq.withCreator != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, project.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &Project{config: pq.config}
		nodes = append(nodes, node)
		values := node.scanValues()
		if withFKs {
			values = append(values, node.fkValues()...)
		}
		return values
	}
	_spec.Assign = func(values ...interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("ent: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(values...)
	}
	if err := sqlgraph.QueryNodes(ctx, pq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := pq.withType; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*Project)
		for i := range nodes {
			if fk := nodes[i].project_type_projects; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(projecttype.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "project_type_projects" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Type = n
			}
		}
	}

	if query := pq.withLocation; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*Project)
		for i := range nodes {
			if fk := nodes[i].project_location; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(location.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "project_location" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Location = n
			}
		}
	}

	if query := pq.withComments; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*Project)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.Comment(func(s *sql.Selector) {
			s.Where(sql.InValues(project.CommentsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.project_comments
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "project_comments" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "project_comments" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Comments = append(node.Edges.Comments, n)
		}
	}

	if query := pq.withWorkOrders; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*Project)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.WorkOrder(func(s *sql.Selector) {
			s.Where(sql.InValues(project.WorkOrdersColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.project_work_orders
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "project_work_orders" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "project_work_orders" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.WorkOrders = append(node.Edges.WorkOrders, n)
		}
	}

	if query := pq.withProperties; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*Project)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.Property(func(s *sql.Selector) {
			s.Where(sql.InValues(project.PropertiesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.project_properties
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "project_properties" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "project_properties" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Properties = append(node.Edges.Properties, n)
		}
	}

	if query := pq.withCreator; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*Project)
		for i := range nodes {
			if fk := nodes[i].project_creator; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(user.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "project_creator" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Creator = n
			}
		}
	}

	return nodes, nil
}

func (pq *ProjectQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := pq.querySpec()
	return sqlgraph.CountNodes(ctx, pq.driver, _spec)
}

func (pq *ProjectQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := pq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (pq *ProjectQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   project.Table,
			Columns: project.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: project.FieldID,
			},
		},
		From:   pq.sql,
		Unique: true,
	}
	if ps := pq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := pq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := pq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := pq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (pq *ProjectQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(pq.driver.Dialect())
	t1 := builder.Table(project.Table)
	selector := builder.Select(t1.Columns(project.Columns...)...).From(t1)
	if pq.sql != nil {
		selector = pq.sql
		selector.Select(selector.Columns(project.Columns...)...)
	}
	for _, p := range pq.predicates {
		p(selector)
	}
	for _, p := range pq.order {
		p(selector)
	}
	if offset := pq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := pq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ProjectGroupBy is the builder for group-by Project entities.
type ProjectGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (pgb *ProjectGroupBy) Aggregate(fns ...AggregateFunc) *ProjectGroupBy {
	pgb.fns = append(pgb.fns, fns...)
	return pgb
}

// Scan applies the group-by query and scan the result into the given value.
func (pgb *ProjectGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := pgb.path(ctx)
	if err != nil {
		return err
	}
	pgb.sql = query
	return pgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (pgb *ProjectGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := pgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (pgb *ProjectGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(pgb.fields) > 1 {
		return nil, errors.New("ent: ProjectGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := pgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (pgb *ProjectGroupBy) StringsX(ctx context.Context) []string {
	v, err := pgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (pgb *ProjectGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(pgb.fields) > 1 {
		return nil, errors.New("ent: ProjectGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := pgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (pgb *ProjectGroupBy) IntsX(ctx context.Context) []int {
	v, err := pgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (pgb *ProjectGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(pgb.fields) > 1 {
		return nil, errors.New("ent: ProjectGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := pgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (pgb *ProjectGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := pgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (pgb *ProjectGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(pgb.fields) > 1 {
		return nil, errors.New("ent: ProjectGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := pgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (pgb *ProjectGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := pgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (pgb *ProjectGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := pgb.sqlQuery().Query()
	if err := pgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (pgb *ProjectGroupBy) sqlQuery() *sql.Selector {
	selector := pgb.sql
	columns := make([]string, 0, len(pgb.fields)+len(pgb.fns))
	columns = append(columns, pgb.fields...)
	for _, fn := range pgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(pgb.fields...)
}

// ProjectSelect is the builder for select fields of Project entities.
type ProjectSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (ps *ProjectSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := ps.path(ctx)
	if err != nil {
		return err
	}
	ps.sql = query
	return ps.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ps *ProjectSelect) ScanX(ctx context.Context, v interface{}) {
	if err := ps.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (ps *ProjectSelect) Strings(ctx context.Context) ([]string, error) {
	if len(ps.fields) > 1 {
		return nil, errors.New("ent: ProjectSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := ps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ps *ProjectSelect) StringsX(ctx context.Context) []string {
	v, err := ps.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (ps *ProjectSelect) Ints(ctx context.Context) ([]int, error) {
	if len(ps.fields) > 1 {
		return nil, errors.New("ent: ProjectSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := ps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ps *ProjectSelect) IntsX(ctx context.Context) []int {
	v, err := ps.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (ps *ProjectSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(ps.fields) > 1 {
		return nil, errors.New("ent: ProjectSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := ps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ps *ProjectSelect) Float64sX(ctx context.Context) []float64 {
	v, err := ps.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (ps *ProjectSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(ps.fields) > 1 {
		return nil, errors.New("ent: ProjectSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := ps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ps *ProjectSelect) BoolsX(ctx context.Context) []bool {
	v, err := ps.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ps *ProjectSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ps.sqlQuery().Query()
	if err := ps.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ps *ProjectSelect) sqlQuery() sql.Querier {
	selector := ps.sql
	selector.Select(selector.Columns(ps.fields...)...)
	return selector
}
