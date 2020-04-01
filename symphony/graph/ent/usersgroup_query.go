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
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/usersgroup"
)

// UsersGroupQuery is the builder for querying UsersGroup entities.
type UsersGroupQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.UsersGroup
	// eager-loading edges.
	withMembers *UserQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (ugq *UsersGroupQuery) Where(ps ...predicate.UsersGroup) *UsersGroupQuery {
	ugq.predicates = append(ugq.predicates, ps...)
	return ugq
}

// Limit adds a limit step to the query.
func (ugq *UsersGroupQuery) Limit(limit int) *UsersGroupQuery {
	ugq.limit = &limit
	return ugq
}

// Offset adds an offset step to the query.
func (ugq *UsersGroupQuery) Offset(offset int) *UsersGroupQuery {
	ugq.offset = &offset
	return ugq
}

// Order adds an order step to the query.
func (ugq *UsersGroupQuery) Order(o ...Order) *UsersGroupQuery {
	ugq.order = append(ugq.order, o...)
	return ugq
}

// QueryMembers chains the current query on the members edge.
func (ugq *UsersGroupQuery) QueryMembers() *UserQuery {
	query := &UserQuery{config: ugq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ugq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(usersgroup.Table, usersgroup.FieldID, ugq.sqlQuery()),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, usersgroup.MembersTable, usersgroup.MembersPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(ugq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first UsersGroup entity in the query. Returns *NotFoundError when no usersgroup was found.
func (ugq *UsersGroupQuery) First(ctx context.Context) (*UsersGroup, error) {
	ugs, err := ugq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(ugs) == 0 {
		return nil, &NotFoundError{usersgroup.Label}
	}
	return ugs[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (ugq *UsersGroupQuery) FirstX(ctx context.Context) *UsersGroup {
	ug, err := ugq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return ug
}

// FirstID returns the first UsersGroup id in the query. Returns *NotFoundError when no id was found.
func (ugq *UsersGroupQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ugq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{usersgroup.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (ugq *UsersGroupQuery) FirstXID(ctx context.Context) int {
	id, err := ugq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only UsersGroup entity in the query, returns an error if not exactly one entity was returned.
func (ugq *UsersGroupQuery) Only(ctx context.Context) (*UsersGroup, error) {
	ugs, err := ugq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(ugs) {
	case 1:
		return ugs[0], nil
	case 0:
		return nil, &NotFoundError{usersgroup.Label}
	default:
		return nil, &NotSingularError{usersgroup.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (ugq *UsersGroupQuery) OnlyX(ctx context.Context) *UsersGroup {
	ug, err := ugq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return ug
}

// OnlyID returns the only UsersGroup id in the query, returns an error if not exactly one id was returned.
func (ugq *UsersGroupQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ugq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{usersgroup.Label}
	default:
		err = &NotSingularError{usersgroup.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (ugq *UsersGroupQuery) OnlyXID(ctx context.Context) int {
	id, err := ugq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of UsersGroups.
func (ugq *UsersGroupQuery) All(ctx context.Context) ([]*UsersGroup, error) {
	if err := ugq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return ugq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (ugq *UsersGroupQuery) AllX(ctx context.Context) []*UsersGroup {
	ugs, err := ugq.All(ctx)
	if err != nil {
		panic(err)
	}
	return ugs
}

// IDs executes the query and returns a list of UsersGroup ids.
func (ugq *UsersGroupQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := ugq.Select(usersgroup.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (ugq *UsersGroupQuery) IDsX(ctx context.Context) []int {
	ids, err := ugq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (ugq *UsersGroupQuery) Count(ctx context.Context) (int, error) {
	if err := ugq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return ugq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (ugq *UsersGroupQuery) CountX(ctx context.Context) int {
	count, err := ugq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (ugq *UsersGroupQuery) Exist(ctx context.Context) (bool, error) {
	if err := ugq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return ugq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (ugq *UsersGroupQuery) ExistX(ctx context.Context) bool {
	exist, err := ugq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (ugq *UsersGroupQuery) Clone() *UsersGroupQuery {
	return &UsersGroupQuery{
		config:     ugq.config,
		limit:      ugq.limit,
		offset:     ugq.offset,
		order:      append([]Order{}, ugq.order...),
		unique:     append([]string{}, ugq.unique...),
		predicates: append([]predicate.UsersGroup{}, ugq.predicates...),
		// clone intermediate query.
		sql:  ugq.sql.Clone(),
		path: ugq.path,
	}
}

//  WithMembers tells the query-builder to eager-loads the nodes that are connected to
// the "members" edge. The optional arguments used to configure the query builder of the edge.
func (ugq *UsersGroupQuery) WithMembers(opts ...func(*UserQuery)) *UsersGroupQuery {
	query := &UserQuery{config: ugq.config}
	for _, opt := range opts {
		opt(query)
	}
	ugq.withMembers = query
	return ugq
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
//	client.UsersGroup.Query().
//		GroupBy(usersgroup.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (ugq *UsersGroupQuery) GroupBy(field string, fields ...string) *UsersGroupGroupBy {
	group := &UsersGroupGroupBy{config: ugq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := ugq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return ugq.sqlQuery(), nil
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
//	client.UsersGroup.Query().
//		Select(usersgroup.FieldCreateTime).
//		Scan(ctx, &v)
//
func (ugq *UsersGroupQuery) Select(field string, fields ...string) *UsersGroupSelect {
	selector := &UsersGroupSelect{config: ugq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := ugq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return ugq.sqlQuery(), nil
	}
	return selector
}

func (ugq *UsersGroupQuery) prepareQuery(ctx context.Context) error {
	if ugq.path != nil {
		prev, err := ugq.path(ctx)
		if err != nil {
			return err
		}
		ugq.sql = prev
	}
	return nil
}

func (ugq *UsersGroupQuery) sqlAll(ctx context.Context) ([]*UsersGroup, error) {
	var (
		nodes       = []*UsersGroup{}
		_spec       = ugq.querySpec()
		loadedTypes = [1]bool{
			ugq.withMembers != nil,
		}
	)
	_spec.ScanValues = func() []interface{} {
		node := &UsersGroup{config: ugq.config}
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
	if err := sqlgraph.QueryNodes(ctx, ugq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := ugq.withMembers; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		ids := make(map[int]*UsersGroup, len(nodes))
		for _, node := range nodes {
			ids[node.ID] = node
			fks = append(fks, node.ID)
		}
		var (
			edgeids []int
			edges   = make(map[int][]*UsersGroup)
		)
		_spec := &sqlgraph.EdgeQuerySpec{
			Edge: &sqlgraph.EdgeSpec{
				Inverse: false,
				Table:   usersgroup.MembersTable,
				Columns: usersgroup.MembersPrimaryKey,
			},
			Predicate: func(s *sql.Selector) {
				s.Where(sql.InValues(usersgroup.MembersPrimaryKey[0], fks...))
			},

			ScanValues: func() [2]interface{} {
				return [2]interface{}{&sql.NullInt64{}, &sql.NullInt64{}}
			},
			Assign: func(out, in interface{}) error {
				eout, ok := out.(*sql.NullInt64)
				if !ok || eout == nil {
					return fmt.Errorf("unexpected id value for edge-out")
				}
				ein, ok := in.(*sql.NullInt64)
				if !ok || ein == nil {
					return fmt.Errorf("unexpected id value for edge-in")
				}
				outValue := int(eout.Int64)
				inValue := int(ein.Int64)
				node, ok := ids[outValue]
				if !ok {
					return fmt.Errorf("unexpected node id in edges: %v", outValue)
				}
				edgeids = append(edgeids, inValue)
				edges[inValue] = append(edges[inValue], node)
				return nil
			},
		}
		if err := sqlgraph.QueryEdges(ctx, ugq.driver, _spec); err != nil {
			return nil, fmt.Errorf(`query edges "members": %v`, err)
		}
		query.Where(user.IDIn(edgeids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := edges[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected "members" node returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Members = append(nodes[i].Edges.Members, n)
			}
		}
	}

	return nodes, nil
}

func (ugq *UsersGroupQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := ugq.querySpec()
	return sqlgraph.CountNodes(ctx, ugq.driver, _spec)
}

func (ugq *UsersGroupQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := ugq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (ugq *UsersGroupQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   usersgroup.Table,
			Columns: usersgroup.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: usersgroup.FieldID,
			},
		},
		From:   ugq.sql,
		Unique: true,
	}
	if ps := ugq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := ugq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := ugq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := ugq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (ugq *UsersGroupQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(ugq.driver.Dialect())
	t1 := builder.Table(usersgroup.Table)
	selector := builder.Select(t1.Columns(usersgroup.Columns...)...).From(t1)
	if ugq.sql != nil {
		selector = ugq.sql
		selector.Select(selector.Columns(usersgroup.Columns...)...)
	}
	for _, p := range ugq.predicates {
		p(selector)
	}
	for _, p := range ugq.order {
		p(selector)
	}
	if offset := ugq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := ugq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// UsersGroupGroupBy is the builder for group-by UsersGroup entities.
type UsersGroupGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (uggb *UsersGroupGroupBy) Aggregate(fns ...Aggregate) *UsersGroupGroupBy {
	uggb.fns = append(uggb.fns, fns...)
	return uggb
}

// Scan applies the group-by query and scan the result into the given value.
func (uggb *UsersGroupGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := uggb.path(ctx)
	if err != nil {
		return err
	}
	uggb.sql = query
	return uggb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (uggb *UsersGroupGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := uggb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (uggb *UsersGroupGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(uggb.fields) > 1 {
		return nil, errors.New("ent: UsersGroupGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := uggb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (uggb *UsersGroupGroupBy) StringsX(ctx context.Context) []string {
	v, err := uggb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (uggb *UsersGroupGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(uggb.fields) > 1 {
		return nil, errors.New("ent: UsersGroupGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := uggb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (uggb *UsersGroupGroupBy) IntsX(ctx context.Context) []int {
	v, err := uggb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (uggb *UsersGroupGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(uggb.fields) > 1 {
		return nil, errors.New("ent: UsersGroupGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := uggb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (uggb *UsersGroupGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := uggb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (uggb *UsersGroupGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(uggb.fields) > 1 {
		return nil, errors.New("ent: UsersGroupGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := uggb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (uggb *UsersGroupGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := uggb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (uggb *UsersGroupGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := uggb.sqlQuery().Query()
	if err := uggb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (uggb *UsersGroupGroupBy) sqlQuery() *sql.Selector {
	selector := uggb.sql
	columns := make([]string, 0, len(uggb.fields)+len(uggb.fns))
	columns = append(columns, uggb.fields...)
	for _, fn := range uggb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(uggb.fields...)
}

// UsersGroupSelect is the builder for select fields of UsersGroup entities.
type UsersGroupSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (ugs *UsersGroupSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := ugs.path(ctx)
	if err != nil {
		return err
	}
	ugs.sql = query
	return ugs.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ugs *UsersGroupSelect) ScanX(ctx context.Context, v interface{}) {
	if err := ugs.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (ugs *UsersGroupSelect) Strings(ctx context.Context) ([]string, error) {
	if len(ugs.fields) > 1 {
		return nil, errors.New("ent: UsersGroupSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := ugs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ugs *UsersGroupSelect) StringsX(ctx context.Context) []string {
	v, err := ugs.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (ugs *UsersGroupSelect) Ints(ctx context.Context) ([]int, error) {
	if len(ugs.fields) > 1 {
		return nil, errors.New("ent: UsersGroupSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := ugs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ugs *UsersGroupSelect) IntsX(ctx context.Context) []int {
	v, err := ugs.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (ugs *UsersGroupSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(ugs.fields) > 1 {
		return nil, errors.New("ent: UsersGroupSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := ugs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ugs *UsersGroupSelect) Float64sX(ctx context.Context) []float64 {
	v, err := ugs.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (ugs *UsersGroupSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(ugs.fields) > 1 {
		return nil, errors.New("ent: UsersGroupSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := ugs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ugs *UsersGroupSelect) BoolsX(ctx context.Context) []bool {
	v, err := ugs.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ugs *UsersGroupSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ugs.sqlQuery().Query()
	if err := ugs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ugs *UsersGroupSelect) sqlQuery() sql.Querier {
	selector := ugs.sql
	selector.Select(selector.Columns(ugs.fields...)...)
	return selector
}
