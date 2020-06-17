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
	"github.com/facebookincubator/symphony/graph/ent/checklistcategory"
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// CheckListCategoryQuery is the builder for querying CheckListCategory entities.
type CheckListCategoryQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.CheckListCategory
	// eager-loading edges.
	withCheckListItems *CheckListItemQuery
	withFKs            bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (clcq *CheckListCategoryQuery) Where(ps ...predicate.CheckListCategory) *CheckListCategoryQuery {
	clcq.predicates = append(clcq.predicates, ps...)
	return clcq
}

// Limit adds a limit step to the query.
func (clcq *CheckListCategoryQuery) Limit(limit int) *CheckListCategoryQuery {
	clcq.limit = &limit
	return clcq
}

// Offset adds an offset step to the query.
func (clcq *CheckListCategoryQuery) Offset(offset int) *CheckListCategoryQuery {
	clcq.offset = &offset
	return clcq
}

// Order adds an order step to the query.
func (clcq *CheckListCategoryQuery) Order(o ...Order) *CheckListCategoryQuery {
	clcq.order = append(clcq.order, o...)
	return clcq
}

// QueryCheckListItems chains the current query on the check_list_items edge.
func (clcq *CheckListCategoryQuery) QueryCheckListItems() *CheckListItemQuery {
	query := &CheckListItemQuery{config: clcq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := clcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(checklistcategory.Table, checklistcategory.FieldID, clcq.sqlQuery()),
			sqlgraph.To(checklistitem.Table, checklistitem.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, checklistcategory.CheckListItemsTable, checklistcategory.CheckListItemsColumn),
		)
		fromU = sqlgraph.SetNeighbors(clcq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first CheckListCategory entity in the query. Returns *NotFoundError when no checklistcategory was found.
func (clcq *CheckListCategoryQuery) First(ctx context.Context) (*CheckListCategory, error) {
	clcs, err := clcq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(clcs) == 0 {
		return nil, &NotFoundError{checklistcategory.Label}
	}
	return clcs[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (clcq *CheckListCategoryQuery) FirstX(ctx context.Context) *CheckListCategory {
	clc, err := clcq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return clc
}

// FirstID returns the first CheckListCategory id in the query. Returns *NotFoundError when no id was found.
func (clcq *CheckListCategoryQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = clcq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{checklistcategory.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (clcq *CheckListCategoryQuery) FirstXID(ctx context.Context) int {
	id, err := clcq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only CheckListCategory entity in the query, returns an error if not exactly one entity was returned.
func (clcq *CheckListCategoryQuery) Only(ctx context.Context) (*CheckListCategory, error) {
	clcs, err := clcq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(clcs) {
	case 1:
		return clcs[0], nil
	case 0:
		return nil, &NotFoundError{checklistcategory.Label}
	default:
		return nil, &NotSingularError{checklistcategory.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (clcq *CheckListCategoryQuery) OnlyX(ctx context.Context) *CheckListCategory {
	clc, err := clcq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return clc
}

// OnlyID returns the only CheckListCategory id in the query, returns an error if not exactly one id was returned.
func (clcq *CheckListCategoryQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = clcq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{checklistcategory.Label}
	default:
		err = &NotSingularError{checklistcategory.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (clcq *CheckListCategoryQuery) OnlyXID(ctx context.Context) int {
	id, err := clcq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of CheckListCategories.
func (clcq *CheckListCategoryQuery) All(ctx context.Context) ([]*CheckListCategory, error) {
	if err := clcq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return clcq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (clcq *CheckListCategoryQuery) AllX(ctx context.Context) []*CheckListCategory {
	clcs, err := clcq.All(ctx)
	if err != nil {
		panic(err)
	}
	return clcs
}

// IDs executes the query and returns a list of CheckListCategory ids.
func (clcq *CheckListCategoryQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := clcq.Select(checklistcategory.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (clcq *CheckListCategoryQuery) IDsX(ctx context.Context) []int {
	ids, err := clcq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (clcq *CheckListCategoryQuery) Count(ctx context.Context) (int, error) {
	if err := clcq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return clcq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (clcq *CheckListCategoryQuery) CountX(ctx context.Context) int {
	count, err := clcq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (clcq *CheckListCategoryQuery) Exist(ctx context.Context) (bool, error) {
	if err := clcq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return clcq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (clcq *CheckListCategoryQuery) ExistX(ctx context.Context) bool {
	exist, err := clcq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (clcq *CheckListCategoryQuery) Clone() *CheckListCategoryQuery {
	return &CheckListCategoryQuery{
		config:     clcq.config,
		limit:      clcq.limit,
		offset:     clcq.offset,
		order:      append([]Order{}, clcq.order...),
		unique:     append([]string{}, clcq.unique...),
		predicates: append([]predicate.CheckListCategory{}, clcq.predicates...),
		// clone intermediate query.
		sql:  clcq.sql.Clone(),
		path: clcq.path,
	}
}

//  WithCheckListItems tells the query-builder to eager-loads the nodes that are connected to
// the "check_list_items" edge. The optional arguments used to configure the query builder of the edge.
func (clcq *CheckListCategoryQuery) WithCheckListItems(opts ...func(*CheckListItemQuery)) *CheckListCategoryQuery {
	query := &CheckListItemQuery{config: clcq.config}
	for _, opt := range opts {
		opt(query)
	}
	clcq.withCheckListItems = query
	return clcq
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
//	client.CheckListCategory.Query().
//		GroupBy(checklistcategory.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (clcq *CheckListCategoryQuery) GroupBy(field string, fields ...string) *CheckListCategoryGroupBy {
	group := &CheckListCategoryGroupBy{config: clcq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := clcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return clcq.sqlQuery(), nil
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
//	client.CheckListCategory.Query().
//		Select(checklistcategory.FieldCreateTime).
//		Scan(ctx, &v)
//
func (clcq *CheckListCategoryQuery) Select(field string, fields ...string) *CheckListCategorySelect {
	selector := &CheckListCategorySelect{config: clcq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := clcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return clcq.sqlQuery(), nil
	}
	return selector
}

func (clcq *CheckListCategoryQuery) prepareQuery(ctx context.Context) error {
	if clcq.path != nil {
		prev, err := clcq.path(ctx)
		if err != nil {
			return err
		}
		clcq.sql = prev
	}
	return nil
}

func (clcq *CheckListCategoryQuery) sqlAll(ctx context.Context) ([]*CheckListCategory, error) {
	var (
		nodes       = []*CheckListCategory{}
		withFKs     = clcq.withFKs
		_spec       = clcq.querySpec()
		loadedTypes = [1]bool{
			clcq.withCheckListItems != nil,
		}
	)
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, checklistcategory.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &CheckListCategory{config: clcq.config}
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
	if err := sqlgraph.QueryNodes(ctx, clcq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := clcq.withCheckListItems; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*CheckListCategory)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.CheckListItem(func(s *sql.Selector) {
			s.Where(sql.InValues(checklistcategory.CheckListItemsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.check_list_category_check_list_items
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "check_list_category_check_list_items" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "check_list_category_check_list_items" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.CheckListItems = append(node.Edges.CheckListItems, n)
		}
	}

	return nodes, nil
}

func (clcq *CheckListCategoryQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := clcq.querySpec()
	return sqlgraph.CountNodes(ctx, clcq.driver, _spec)
}

func (clcq *CheckListCategoryQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := clcq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (clcq *CheckListCategoryQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   checklistcategory.Table,
			Columns: checklistcategory.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistcategory.FieldID,
			},
		},
		From:   clcq.sql,
		Unique: true,
	}
	if ps := clcq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := clcq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := clcq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := clcq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (clcq *CheckListCategoryQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(clcq.driver.Dialect())
	t1 := builder.Table(checklistcategory.Table)
	selector := builder.Select(t1.Columns(checklistcategory.Columns...)...).From(t1)
	if clcq.sql != nil {
		selector = clcq.sql
		selector.Select(selector.Columns(checklistcategory.Columns...)...)
	}
	for _, p := range clcq.predicates {
		p(selector)
	}
	for _, p := range clcq.order {
		p(selector)
	}
	if offset := clcq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := clcq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// CheckListCategoryGroupBy is the builder for group-by CheckListCategory entities.
type CheckListCategoryGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (clcgb *CheckListCategoryGroupBy) Aggregate(fns ...Aggregate) *CheckListCategoryGroupBy {
	clcgb.fns = append(clcgb.fns, fns...)
	return clcgb
}

// Scan applies the group-by query and scan the result into the given value.
func (clcgb *CheckListCategoryGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := clcgb.path(ctx)
	if err != nil {
		return err
	}
	clcgb.sql = query
	return clcgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (clcgb *CheckListCategoryGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := clcgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (clcgb *CheckListCategoryGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(clcgb.fields) > 1 {
		return nil, errors.New("ent: CheckListCategoryGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := clcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (clcgb *CheckListCategoryGroupBy) StringsX(ctx context.Context) []string {
	v, err := clcgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (clcgb *CheckListCategoryGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(clcgb.fields) > 1 {
		return nil, errors.New("ent: CheckListCategoryGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := clcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (clcgb *CheckListCategoryGroupBy) IntsX(ctx context.Context) []int {
	v, err := clcgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (clcgb *CheckListCategoryGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(clcgb.fields) > 1 {
		return nil, errors.New("ent: CheckListCategoryGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := clcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (clcgb *CheckListCategoryGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := clcgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (clcgb *CheckListCategoryGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(clcgb.fields) > 1 {
		return nil, errors.New("ent: CheckListCategoryGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := clcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (clcgb *CheckListCategoryGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := clcgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (clcgb *CheckListCategoryGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := clcgb.sqlQuery().Query()
	if err := clcgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (clcgb *CheckListCategoryGroupBy) sqlQuery() *sql.Selector {
	selector := clcgb.sql
	columns := make([]string, 0, len(clcgb.fields)+len(clcgb.fns))
	columns = append(columns, clcgb.fields...)
	for _, fn := range clcgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(clcgb.fields...)
}

// CheckListCategorySelect is the builder for select fields of CheckListCategory entities.
type CheckListCategorySelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (clcs *CheckListCategorySelect) Scan(ctx context.Context, v interface{}) error {
	query, err := clcs.path(ctx)
	if err != nil {
		return err
	}
	clcs.sql = query
	return clcs.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (clcs *CheckListCategorySelect) ScanX(ctx context.Context, v interface{}) {
	if err := clcs.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (clcs *CheckListCategorySelect) Strings(ctx context.Context) ([]string, error) {
	if len(clcs.fields) > 1 {
		return nil, errors.New("ent: CheckListCategorySelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := clcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (clcs *CheckListCategorySelect) StringsX(ctx context.Context) []string {
	v, err := clcs.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (clcs *CheckListCategorySelect) Ints(ctx context.Context) ([]int, error) {
	if len(clcs.fields) > 1 {
		return nil, errors.New("ent: CheckListCategorySelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := clcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (clcs *CheckListCategorySelect) IntsX(ctx context.Context) []int {
	v, err := clcs.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (clcs *CheckListCategorySelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(clcs.fields) > 1 {
		return nil, errors.New("ent: CheckListCategorySelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := clcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (clcs *CheckListCategorySelect) Float64sX(ctx context.Context) []float64 {
	v, err := clcs.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (clcs *CheckListCategorySelect) Bools(ctx context.Context) ([]bool, error) {
	if len(clcs.fields) > 1 {
		return nil, errors.New("ent: CheckListCategorySelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := clcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (clcs *CheckListCategorySelect) BoolsX(ctx context.Context) []bool {
	v, err := clcs.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (clcs *CheckListCategorySelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := clcs.sqlQuery().Query()
	if err := clcs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (clcs *CheckListCategorySelect) sqlQuery() sql.Querier {
	selector := clcs.sql
	selector.Select(selector.Columns(clcs.fields...)...)
	return selector
}
