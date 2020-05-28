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
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/project"
	"github.com/facebookincubator/symphony/pkg/ent/projecttype"
	"github.com/facebookincubator/symphony/pkg/ent/propertytype"
	"github.com/facebookincubator/symphony/pkg/ent/workorderdefinition"
)

// ProjectTypeQuery is the builder for querying ProjectType entities.
type ProjectTypeQuery struct {
	config
	limit      *int
	offset     *int
	order      []OrderFunc
	unique     []string
	predicates []predicate.ProjectType
	// eager-loading edges.
	withProjects   *ProjectQuery
	withProperties *PropertyTypeQuery
	withWorkOrders *WorkOrderDefinitionQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (ptq *ProjectTypeQuery) Where(ps ...predicate.ProjectType) *ProjectTypeQuery {
	ptq.predicates = append(ptq.predicates, ps...)
	return ptq
}

// Limit adds a limit step to the query.
func (ptq *ProjectTypeQuery) Limit(limit int) *ProjectTypeQuery {
	ptq.limit = &limit
	return ptq
}

// Offset adds an offset step to the query.
func (ptq *ProjectTypeQuery) Offset(offset int) *ProjectTypeQuery {
	ptq.offset = &offset
	return ptq
}

// Order adds an order step to the query.
func (ptq *ProjectTypeQuery) Order(o ...OrderFunc) *ProjectTypeQuery {
	ptq.order = append(ptq.order, o...)
	return ptq
}

// QueryProjects chains the current query on the projects edge.
func (ptq *ProjectTypeQuery) QueryProjects() *ProjectQuery {
	query := &ProjectQuery{config: ptq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(projecttype.Table, projecttype.FieldID, ptq.sqlQuery()),
			sqlgraph.To(project.Table, project.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, projecttype.ProjectsTable, projecttype.ProjectsColumn),
		)
		fromU = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryProperties chains the current query on the properties edge.
func (ptq *ProjectTypeQuery) QueryProperties() *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: ptq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(projecttype.Table, projecttype.FieldID, ptq.sqlQuery()),
			sqlgraph.To(propertytype.Table, propertytype.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, projecttype.PropertiesTable, projecttype.PropertiesColumn),
		)
		fromU = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryWorkOrders chains the current query on the work_orders edge.
func (ptq *ProjectTypeQuery) QueryWorkOrders() *WorkOrderDefinitionQuery {
	query := &WorkOrderDefinitionQuery{config: ptq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(projecttype.Table, projecttype.FieldID, ptq.sqlQuery()),
			sqlgraph.To(workorderdefinition.Table, workorderdefinition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, projecttype.WorkOrdersTable, projecttype.WorkOrdersColumn),
		)
		fromU = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first ProjectType entity in the query. Returns *NotFoundError when no projecttype was found.
func (ptq *ProjectTypeQuery) First(ctx context.Context) (*ProjectType, error) {
	pts, err := ptq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(pts) == 0 {
		return nil, &NotFoundError{projecttype.Label}
	}
	return pts[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (ptq *ProjectTypeQuery) FirstX(ctx context.Context) *ProjectType {
	pt, err := ptq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return pt
}

// FirstID returns the first ProjectType id in the query. Returns *NotFoundError when no id was found.
func (ptq *ProjectTypeQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ptq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{projecttype.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (ptq *ProjectTypeQuery) FirstXID(ctx context.Context) int {
	id, err := ptq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only ProjectType entity in the query, returns an error if not exactly one entity was returned.
func (ptq *ProjectTypeQuery) Only(ctx context.Context) (*ProjectType, error) {
	pts, err := ptq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(pts) {
	case 1:
		return pts[0], nil
	case 0:
		return nil, &NotFoundError{projecttype.Label}
	default:
		return nil, &NotSingularError{projecttype.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (ptq *ProjectTypeQuery) OnlyX(ctx context.Context) *ProjectType {
	pt, err := ptq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return pt
}

// OnlyID returns the only ProjectType id in the query, returns an error if not exactly one id was returned.
func (ptq *ProjectTypeQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ptq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{projecttype.Label}
	default:
		err = &NotSingularError{projecttype.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (ptq *ProjectTypeQuery) OnlyXID(ctx context.Context) int {
	id, err := ptq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ProjectTypes.
func (ptq *ProjectTypeQuery) All(ctx context.Context) ([]*ProjectType, error) {
	if err := ptq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return ptq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (ptq *ProjectTypeQuery) AllX(ctx context.Context) []*ProjectType {
	pts, err := ptq.All(ctx)
	if err != nil {
		panic(err)
	}
	return pts
}

// IDs executes the query and returns a list of ProjectType ids.
func (ptq *ProjectTypeQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := ptq.Select(projecttype.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (ptq *ProjectTypeQuery) IDsX(ctx context.Context) []int {
	ids, err := ptq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (ptq *ProjectTypeQuery) Count(ctx context.Context) (int, error) {
	if err := ptq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return ptq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (ptq *ProjectTypeQuery) CountX(ctx context.Context) int {
	count, err := ptq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (ptq *ProjectTypeQuery) Exist(ctx context.Context) (bool, error) {
	if err := ptq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return ptq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (ptq *ProjectTypeQuery) ExistX(ctx context.Context) bool {
	exist, err := ptq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (ptq *ProjectTypeQuery) Clone() *ProjectTypeQuery {
	return &ProjectTypeQuery{
		config:     ptq.config,
		limit:      ptq.limit,
		offset:     ptq.offset,
		order:      append([]OrderFunc{}, ptq.order...),
		unique:     append([]string{}, ptq.unique...),
		predicates: append([]predicate.ProjectType{}, ptq.predicates...),
		// clone intermediate query.
		sql:  ptq.sql.Clone(),
		path: ptq.path,
	}
}

//  WithProjects tells the query-builder to eager-loads the nodes that are connected to
// the "projects" edge. The optional arguments used to configure the query builder of the edge.
func (ptq *ProjectTypeQuery) WithProjects(opts ...func(*ProjectQuery)) *ProjectTypeQuery {
	query := &ProjectQuery{config: ptq.config}
	for _, opt := range opts {
		opt(query)
	}
	ptq.withProjects = query
	return ptq
}

//  WithProperties tells the query-builder to eager-loads the nodes that are connected to
// the "properties" edge. The optional arguments used to configure the query builder of the edge.
func (ptq *ProjectTypeQuery) WithProperties(opts ...func(*PropertyTypeQuery)) *ProjectTypeQuery {
	query := &PropertyTypeQuery{config: ptq.config}
	for _, opt := range opts {
		opt(query)
	}
	ptq.withProperties = query
	return ptq
}

//  WithWorkOrders tells the query-builder to eager-loads the nodes that are connected to
// the "work_orders" edge. The optional arguments used to configure the query builder of the edge.
func (ptq *ProjectTypeQuery) WithWorkOrders(opts ...func(*WorkOrderDefinitionQuery)) *ProjectTypeQuery {
	query := &WorkOrderDefinitionQuery{config: ptq.config}
	for _, opt := range opts {
		opt(query)
	}
	ptq.withWorkOrders = query
	return ptq
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
//	client.ProjectType.Query().
//		GroupBy(projecttype.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (ptq *ProjectTypeQuery) GroupBy(field string, fields ...string) *ProjectTypeGroupBy {
	group := &ProjectTypeGroupBy{config: ptq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return ptq.sqlQuery(), nil
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
//	client.ProjectType.Query().
//		Select(projecttype.FieldCreateTime).
//		Scan(ctx, &v)
//
func (ptq *ProjectTypeQuery) Select(field string, fields ...string) *ProjectTypeSelect {
	selector := &ProjectTypeSelect{config: ptq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return ptq.sqlQuery(), nil
	}
	return selector
}

func (ptq *ProjectTypeQuery) prepareQuery(ctx context.Context) error {
	if ptq.path != nil {
		prev, err := ptq.path(ctx)
		if err != nil {
			return err
		}
		ptq.sql = prev
	}
	if err := projecttype.Policy.EvalQuery(ctx, ptq); err != nil {
		return err
	}
	return nil
}

func (ptq *ProjectTypeQuery) sqlAll(ctx context.Context) ([]*ProjectType, error) {
	var (
		nodes       = []*ProjectType{}
		_spec       = ptq.querySpec()
		loadedTypes = [3]bool{
			ptq.withProjects != nil,
			ptq.withProperties != nil,
			ptq.withWorkOrders != nil,
		}
	)
	_spec.ScanValues = func() []interface{} {
		node := &ProjectType{config: ptq.config}
		nodes = append(nodes, node)
		values := node.scanValues()
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
	if err := sqlgraph.QueryNodes(ctx, ptq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := ptq.withProjects; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*ProjectType)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.Project(func(s *sql.Selector) {
			s.Where(sql.InValues(projecttype.ProjectsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.project_type_projects
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "project_type_projects" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "project_type_projects" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Projects = append(node.Edges.Projects, n)
		}
	}

	if query := ptq.withProperties; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*ProjectType)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.PropertyType(func(s *sql.Selector) {
			s.Where(sql.InValues(projecttype.PropertiesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.project_type_properties
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "project_type_properties" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "project_type_properties" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Properties = append(node.Edges.Properties, n)
		}
	}

	if query := ptq.withWorkOrders; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*ProjectType)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.WorkOrderDefinition(func(s *sql.Selector) {
			s.Where(sql.InValues(projecttype.WorkOrdersColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.project_type_work_orders
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "project_type_work_orders" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "project_type_work_orders" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.WorkOrders = append(node.Edges.WorkOrders, n)
		}
	}

	return nodes, nil
}

func (ptq *ProjectTypeQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := ptq.querySpec()
	return sqlgraph.CountNodes(ctx, ptq.driver, _spec)
}

func (ptq *ProjectTypeQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := ptq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (ptq *ProjectTypeQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   projecttype.Table,
			Columns: projecttype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: projecttype.FieldID,
			},
		},
		From:   ptq.sql,
		Unique: true,
	}
	if ps := ptq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := ptq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := ptq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := ptq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (ptq *ProjectTypeQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(ptq.driver.Dialect())
	t1 := builder.Table(projecttype.Table)
	selector := builder.Select(t1.Columns(projecttype.Columns...)...).From(t1)
	if ptq.sql != nil {
		selector = ptq.sql
		selector.Select(selector.Columns(projecttype.Columns...)...)
	}
	for _, p := range ptq.predicates {
		p(selector)
	}
	for _, p := range ptq.order {
		p(selector)
	}
	if offset := ptq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := ptq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ProjectTypeGroupBy is the builder for group-by ProjectType entities.
type ProjectTypeGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ptgb *ProjectTypeGroupBy) Aggregate(fns ...AggregateFunc) *ProjectTypeGroupBy {
	ptgb.fns = append(ptgb.fns, fns...)
	return ptgb
}

// Scan applies the group-by query and scan the result into the given value.
func (ptgb *ProjectTypeGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := ptgb.path(ctx)
	if err != nil {
		return err
	}
	ptgb.sql = query
	return ptgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ptgb *ProjectTypeGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := ptgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (ptgb *ProjectTypeGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(ptgb.fields) > 1 {
		return nil, errors.New("ent: ProjectTypeGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := ptgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ptgb *ProjectTypeGroupBy) StringsX(ctx context.Context) []string {
	v, err := ptgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (ptgb *ProjectTypeGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(ptgb.fields) > 1 {
		return nil, errors.New("ent: ProjectTypeGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := ptgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ptgb *ProjectTypeGroupBy) IntsX(ctx context.Context) []int {
	v, err := ptgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (ptgb *ProjectTypeGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(ptgb.fields) > 1 {
		return nil, errors.New("ent: ProjectTypeGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := ptgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ptgb *ProjectTypeGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := ptgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (ptgb *ProjectTypeGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(ptgb.fields) > 1 {
		return nil, errors.New("ent: ProjectTypeGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := ptgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ptgb *ProjectTypeGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := ptgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ptgb *ProjectTypeGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ptgb.sqlQuery().Query()
	if err := ptgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ptgb *ProjectTypeGroupBy) sqlQuery() *sql.Selector {
	selector := ptgb.sql
	columns := make([]string, 0, len(ptgb.fields)+len(ptgb.fns))
	columns = append(columns, ptgb.fields...)
	for _, fn := range ptgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(ptgb.fields...)
}

// ProjectTypeSelect is the builder for select fields of ProjectType entities.
type ProjectTypeSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (pts *ProjectTypeSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := pts.path(ctx)
	if err != nil {
		return err
	}
	pts.sql = query
	return pts.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (pts *ProjectTypeSelect) ScanX(ctx context.Context, v interface{}) {
	if err := pts.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (pts *ProjectTypeSelect) Strings(ctx context.Context) ([]string, error) {
	if len(pts.fields) > 1 {
		return nil, errors.New("ent: ProjectTypeSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := pts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (pts *ProjectTypeSelect) StringsX(ctx context.Context) []string {
	v, err := pts.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (pts *ProjectTypeSelect) Ints(ctx context.Context) ([]int, error) {
	if len(pts.fields) > 1 {
		return nil, errors.New("ent: ProjectTypeSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := pts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (pts *ProjectTypeSelect) IntsX(ctx context.Context) []int {
	v, err := pts.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (pts *ProjectTypeSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(pts.fields) > 1 {
		return nil, errors.New("ent: ProjectTypeSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := pts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (pts *ProjectTypeSelect) Float64sX(ctx context.Context) []float64 {
	v, err := pts.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (pts *ProjectTypeSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(pts.fields) > 1 {
		return nil, errors.New("ent: ProjectTypeSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := pts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (pts *ProjectTypeSelect) BoolsX(ctx context.Context) []bool {
	v, err := pts.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (pts *ProjectTypeSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := pts.sqlQuery().Query()
	if err := pts.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (pts *ProjectTypeSelect) sqlQuery() sql.Querier {
	selector := pts.sql
	selector.Select(selector.Columns(pts.fields...)...)
	return selector
}
