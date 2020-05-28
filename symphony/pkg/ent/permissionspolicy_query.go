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
	"github.com/facebookincubator/symphony/pkg/ent/permissionspolicy"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/usersgroup"
)

// PermissionsPolicyQuery is the builder for querying PermissionsPolicy entities.
type PermissionsPolicyQuery struct {
	config
	limit      *int
	offset     *int
	order      []OrderFunc
	unique     []string
	predicates []predicate.PermissionsPolicy
	// eager-loading edges.
	withGroups *UsersGroupQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (ppq *PermissionsPolicyQuery) Where(ps ...predicate.PermissionsPolicy) *PermissionsPolicyQuery {
	ppq.predicates = append(ppq.predicates, ps...)
	return ppq
}

// Limit adds a limit step to the query.
func (ppq *PermissionsPolicyQuery) Limit(limit int) *PermissionsPolicyQuery {
	ppq.limit = &limit
	return ppq
}

// Offset adds an offset step to the query.
func (ppq *PermissionsPolicyQuery) Offset(offset int) *PermissionsPolicyQuery {
	ppq.offset = &offset
	return ppq
}

// Order adds an order step to the query.
func (ppq *PermissionsPolicyQuery) Order(o ...OrderFunc) *PermissionsPolicyQuery {
	ppq.order = append(ppq.order, o...)
	return ppq
}

// QueryGroups chains the current query on the groups edge.
func (ppq *PermissionsPolicyQuery) QueryGroups() *UsersGroupQuery {
	query := &UsersGroupQuery{config: ppq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ppq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(permissionspolicy.Table, permissionspolicy.FieldID, ppq.sqlQuery()),
			sqlgraph.To(usersgroup.Table, usersgroup.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, permissionspolicy.GroupsTable, permissionspolicy.GroupsPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(ppq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first PermissionsPolicy entity in the query. Returns *NotFoundError when no permissionspolicy was found.
func (ppq *PermissionsPolicyQuery) First(ctx context.Context) (*PermissionsPolicy, error) {
	pps, err := ppq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(pps) == 0 {
		return nil, &NotFoundError{permissionspolicy.Label}
	}
	return pps[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (ppq *PermissionsPolicyQuery) FirstX(ctx context.Context) *PermissionsPolicy {
	pp, err := ppq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return pp
}

// FirstID returns the first PermissionsPolicy id in the query. Returns *NotFoundError when no id was found.
func (ppq *PermissionsPolicyQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ppq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{permissionspolicy.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (ppq *PermissionsPolicyQuery) FirstXID(ctx context.Context) int {
	id, err := ppq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only PermissionsPolicy entity in the query, returns an error if not exactly one entity was returned.
func (ppq *PermissionsPolicyQuery) Only(ctx context.Context) (*PermissionsPolicy, error) {
	pps, err := ppq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(pps) {
	case 1:
		return pps[0], nil
	case 0:
		return nil, &NotFoundError{permissionspolicy.Label}
	default:
		return nil, &NotSingularError{permissionspolicy.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (ppq *PermissionsPolicyQuery) OnlyX(ctx context.Context) *PermissionsPolicy {
	pp, err := ppq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return pp
}

// OnlyID returns the only PermissionsPolicy id in the query, returns an error if not exactly one id was returned.
func (ppq *PermissionsPolicyQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ppq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{permissionspolicy.Label}
	default:
		err = &NotSingularError{permissionspolicy.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (ppq *PermissionsPolicyQuery) OnlyXID(ctx context.Context) int {
	id, err := ppq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of PermissionsPolicies.
func (ppq *PermissionsPolicyQuery) All(ctx context.Context) ([]*PermissionsPolicy, error) {
	if err := ppq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return ppq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (ppq *PermissionsPolicyQuery) AllX(ctx context.Context) []*PermissionsPolicy {
	pps, err := ppq.All(ctx)
	if err != nil {
		panic(err)
	}
	return pps
}

// IDs executes the query and returns a list of PermissionsPolicy ids.
func (ppq *PermissionsPolicyQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := ppq.Select(permissionspolicy.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (ppq *PermissionsPolicyQuery) IDsX(ctx context.Context) []int {
	ids, err := ppq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (ppq *PermissionsPolicyQuery) Count(ctx context.Context) (int, error) {
	if err := ppq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return ppq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (ppq *PermissionsPolicyQuery) CountX(ctx context.Context) int {
	count, err := ppq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (ppq *PermissionsPolicyQuery) Exist(ctx context.Context) (bool, error) {
	if err := ppq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return ppq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (ppq *PermissionsPolicyQuery) ExistX(ctx context.Context) bool {
	exist, err := ppq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (ppq *PermissionsPolicyQuery) Clone() *PermissionsPolicyQuery {
	return &PermissionsPolicyQuery{
		config:     ppq.config,
		limit:      ppq.limit,
		offset:     ppq.offset,
		order:      append([]OrderFunc{}, ppq.order...),
		unique:     append([]string{}, ppq.unique...),
		predicates: append([]predicate.PermissionsPolicy{}, ppq.predicates...),
		// clone intermediate query.
		sql:  ppq.sql.Clone(),
		path: ppq.path,
	}
}

//  WithGroups tells the query-builder to eager-loads the nodes that are connected to
// the "groups" edge. The optional arguments used to configure the query builder of the edge.
func (ppq *PermissionsPolicyQuery) WithGroups(opts ...func(*UsersGroupQuery)) *PermissionsPolicyQuery {
	query := &UsersGroupQuery{config: ppq.config}
	for _, opt := range opts {
		opt(query)
	}
	ppq.withGroups = query
	return ppq
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
//	client.PermissionsPolicy.Query().
//		GroupBy(permissionspolicy.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (ppq *PermissionsPolicyQuery) GroupBy(field string, fields ...string) *PermissionsPolicyGroupBy {
	group := &PermissionsPolicyGroupBy{config: ppq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := ppq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return ppq.sqlQuery(), nil
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
//	client.PermissionsPolicy.Query().
//		Select(permissionspolicy.FieldCreateTime).
//		Scan(ctx, &v)
//
func (ppq *PermissionsPolicyQuery) Select(field string, fields ...string) *PermissionsPolicySelect {
	selector := &PermissionsPolicySelect{config: ppq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := ppq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return ppq.sqlQuery(), nil
	}
	return selector
}

func (ppq *PermissionsPolicyQuery) prepareQuery(ctx context.Context) error {
	if ppq.path != nil {
		prev, err := ppq.path(ctx)
		if err != nil {
			return err
		}
		ppq.sql = prev
	}
	if err := permissionspolicy.Policy.EvalQuery(ctx, ppq); err != nil {
		return err
	}
	return nil
}

func (ppq *PermissionsPolicyQuery) sqlAll(ctx context.Context) ([]*PermissionsPolicy, error) {
	var (
		nodes       = []*PermissionsPolicy{}
		_spec       = ppq.querySpec()
		loadedTypes = [1]bool{
			ppq.withGroups != nil,
		}
	)
	_spec.ScanValues = func() []interface{} {
		node := &PermissionsPolicy{config: ppq.config}
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
	if err := sqlgraph.QueryNodes(ctx, ppq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := ppq.withGroups; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		ids := make(map[int]*PermissionsPolicy, len(nodes))
		for _, node := range nodes {
			ids[node.ID] = node
			fks = append(fks, node.ID)
		}
		var (
			edgeids []int
			edges   = make(map[int][]*PermissionsPolicy)
		)
		_spec := &sqlgraph.EdgeQuerySpec{
			Edge: &sqlgraph.EdgeSpec{
				Inverse: true,
				Table:   permissionspolicy.GroupsTable,
				Columns: permissionspolicy.GroupsPrimaryKey,
			},
			Predicate: func(s *sql.Selector) {
				s.Where(sql.InValues(permissionspolicy.GroupsPrimaryKey[1], fks...))
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
		if err := sqlgraph.QueryEdges(ctx, ppq.driver, _spec); err != nil {
			return nil, fmt.Errorf(`query edges "groups": %v`, err)
		}
		query.Where(usersgroup.IDIn(edgeids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := edges[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected "groups" node returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Groups = append(nodes[i].Edges.Groups, n)
			}
		}
	}

	return nodes, nil
}

func (ppq *PermissionsPolicyQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := ppq.querySpec()
	return sqlgraph.CountNodes(ctx, ppq.driver, _spec)
}

func (ppq *PermissionsPolicyQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := ppq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (ppq *PermissionsPolicyQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   permissionspolicy.Table,
			Columns: permissionspolicy.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: permissionspolicy.FieldID,
			},
		},
		From:   ppq.sql,
		Unique: true,
	}
	if ps := ppq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := ppq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := ppq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := ppq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (ppq *PermissionsPolicyQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(ppq.driver.Dialect())
	t1 := builder.Table(permissionspolicy.Table)
	selector := builder.Select(t1.Columns(permissionspolicy.Columns...)...).From(t1)
	if ppq.sql != nil {
		selector = ppq.sql
		selector.Select(selector.Columns(permissionspolicy.Columns...)...)
	}
	for _, p := range ppq.predicates {
		p(selector)
	}
	for _, p := range ppq.order {
		p(selector)
	}
	if offset := ppq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := ppq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// PermissionsPolicyGroupBy is the builder for group-by PermissionsPolicy entities.
type PermissionsPolicyGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ppgb *PermissionsPolicyGroupBy) Aggregate(fns ...AggregateFunc) *PermissionsPolicyGroupBy {
	ppgb.fns = append(ppgb.fns, fns...)
	return ppgb
}

// Scan applies the group-by query and scan the result into the given value.
func (ppgb *PermissionsPolicyGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := ppgb.path(ctx)
	if err != nil {
		return err
	}
	ppgb.sql = query
	return ppgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ppgb *PermissionsPolicyGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := ppgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (ppgb *PermissionsPolicyGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(ppgb.fields) > 1 {
		return nil, errors.New("ent: PermissionsPolicyGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := ppgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ppgb *PermissionsPolicyGroupBy) StringsX(ctx context.Context) []string {
	v, err := ppgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (ppgb *PermissionsPolicyGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(ppgb.fields) > 1 {
		return nil, errors.New("ent: PermissionsPolicyGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := ppgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ppgb *PermissionsPolicyGroupBy) IntsX(ctx context.Context) []int {
	v, err := ppgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (ppgb *PermissionsPolicyGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(ppgb.fields) > 1 {
		return nil, errors.New("ent: PermissionsPolicyGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := ppgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ppgb *PermissionsPolicyGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := ppgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (ppgb *PermissionsPolicyGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(ppgb.fields) > 1 {
		return nil, errors.New("ent: PermissionsPolicyGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := ppgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ppgb *PermissionsPolicyGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := ppgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ppgb *PermissionsPolicyGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ppgb.sqlQuery().Query()
	if err := ppgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ppgb *PermissionsPolicyGroupBy) sqlQuery() *sql.Selector {
	selector := ppgb.sql
	columns := make([]string, 0, len(ppgb.fields)+len(ppgb.fns))
	columns = append(columns, ppgb.fields...)
	for _, fn := range ppgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(ppgb.fields...)
}

// PermissionsPolicySelect is the builder for select fields of PermissionsPolicy entities.
type PermissionsPolicySelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (pps *PermissionsPolicySelect) Scan(ctx context.Context, v interface{}) error {
	query, err := pps.path(ctx)
	if err != nil {
		return err
	}
	pps.sql = query
	return pps.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (pps *PermissionsPolicySelect) ScanX(ctx context.Context, v interface{}) {
	if err := pps.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (pps *PermissionsPolicySelect) Strings(ctx context.Context) ([]string, error) {
	if len(pps.fields) > 1 {
		return nil, errors.New("ent: PermissionsPolicySelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := pps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (pps *PermissionsPolicySelect) StringsX(ctx context.Context) []string {
	v, err := pps.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (pps *PermissionsPolicySelect) Ints(ctx context.Context) ([]int, error) {
	if len(pps.fields) > 1 {
		return nil, errors.New("ent: PermissionsPolicySelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := pps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (pps *PermissionsPolicySelect) IntsX(ctx context.Context) []int {
	v, err := pps.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (pps *PermissionsPolicySelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(pps.fields) > 1 {
		return nil, errors.New("ent: PermissionsPolicySelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := pps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (pps *PermissionsPolicySelect) Float64sX(ctx context.Context) []float64 {
	v, err := pps.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (pps *PermissionsPolicySelect) Bools(ctx context.Context) ([]bool, error) {
	if len(pps.fields) > 1 {
		return nil, errors.New("ent: PermissionsPolicySelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := pps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (pps *PermissionsPolicySelect) BoolsX(ctx context.Context) []bool {
	v, err := pps.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (pps *PermissionsPolicySelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := pps.sqlQuery().Query()
	if err := pps.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (pps *PermissionsPolicySelect) sqlQuery() sql.Querier {
	selector := pps.sql
	selector.Select(selector.Columns(pps.fields...)...)
	return selector
}
