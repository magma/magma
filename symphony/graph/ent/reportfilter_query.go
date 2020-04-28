// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/reportfilter"
)

// ReportFilterQuery is the builder for querying ReportFilter entities.
type ReportFilterQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.ReportFilter
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (rfq *ReportFilterQuery) Where(ps ...predicate.ReportFilter) *ReportFilterQuery {
	rfq.predicates = append(rfq.predicates, ps...)
	return rfq
}

// Limit adds a limit step to the query.
func (rfq *ReportFilterQuery) Limit(limit int) *ReportFilterQuery {
	rfq.limit = &limit
	return rfq
}

// Offset adds an offset step to the query.
func (rfq *ReportFilterQuery) Offset(offset int) *ReportFilterQuery {
	rfq.offset = &offset
	return rfq
}

// Order adds an order step to the query.
func (rfq *ReportFilterQuery) Order(o ...Order) *ReportFilterQuery {
	rfq.order = append(rfq.order, o...)
	return rfq
}

// First returns the first ReportFilter entity in the query. Returns *NotFoundError when no reportfilter was found.
func (rfq *ReportFilterQuery) First(ctx context.Context) (*ReportFilter, error) {
	rves, err := rfq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(rves) == 0 {
		return nil, &NotFoundError{reportfilter.Label}
	}
	return rves[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (rfq *ReportFilterQuery) FirstX(ctx context.Context) *ReportFilter {
	rf, err := rfq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return rf
}

// FirstID returns the first ReportFilter id in the query. Returns *NotFoundError when no id was found.
func (rfq *ReportFilterQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = rfq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{reportfilter.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (rfq *ReportFilterQuery) FirstXID(ctx context.Context) int {
	id, err := rfq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only ReportFilter entity in the query, returns an error if not exactly one entity was returned.
func (rfq *ReportFilterQuery) Only(ctx context.Context) (*ReportFilter, error) {
	rves, err := rfq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(rves) {
	case 1:
		return rves[0], nil
	case 0:
		return nil, &NotFoundError{reportfilter.Label}
	default:
		return nil, &NotSingularError{reportfilter.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (rfq *ReportFilterQuery) OnlyX(ctx context.Context) *ReportFilter {
	rf, err := rfq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return rf
}

// OnlyID returns the only ReportFilter id in the query, returns an error if not exactly one id was returned.
func (rfq *ReportFilterQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = rfq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{reportfilter.Label}
	default:
		err = &NotSingularError{reportfilter.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (rfq *ReportFilterQuery) OnlyXID(ctx context.Context) int {
	id, err := rfq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ReportFilters.
func (rfq *ReportFilterQuery) All(ctx context.Context) ([]*ReportFilter, error) {
	if err := rfq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return rfq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (rfq *ReportFilterQuery) AllX(ctx context.Context) []*ReportFilter {
	rves, err := rfq.All(ctx)
	if err != nil {
		panic(err)
	}
	return rves
}

// IDs executes the query and returns a list of ReportFilter ids.
func (rfq *ReportFilterQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := rfq.Select(reportfilter.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (rfq *ReportFilterQuery) IDsX(ctx context.Context) []int {
	ids, err := rfq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (rfq *ReportFilterQuery) Count(ctx context.Context) (int, error) {
	if err := rfq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return rfq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (rfq *ReportFilterQuery) CountX(ctx context.Context) int {
	count, err := rfq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (rfq *ReportFilterQuery) Exist(ctx context.Context) (bool, error) {
	if err := rfq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return rfq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (rfq *ReportFilterQuery) ExistX(ctx context.Context) bool {
	exist, err := rfq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (rfq *ReportFilterQuery) Clone() *ReportFilterQuery {
	return &ReportFilterQuery{
		config:     rfq.config,
		limit:      rfq.limit,
		offset:     rfq.offset,
		order:      append([]Order{}, rfq.order...),
		unique:     append([]string{}, rfq.unique...),
		predicates: append([]predicate.ReportFilter{}, rfq.predicates...),
		// clone intermediate query.
		sql:  rfq.sql.Clone(),
		path: rfq.path,
	}
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
//	client.ReportFilter.Query().
//		GroupBy(reportfilter.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (rfq *ReportFilterQuery) GroupBy(field string, fields ...string) *ReportFilterGroupBy {
	group := &ReportFilterGroupBy{config: rfq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := rfq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return rfq.sqlQuery(), nil
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
//	client.ReportFilter.Query().
//		Select(reportfilter.FieldCreateTime).
//		Scan(ctx, &v)
//
func (rfq *ReportFilterQuery) Select(field string, fields ...string) *ReportFilterSelect {
	selector := &ReportFilterSelect{config: rfq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := rfq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return rfq.sqlQuery(), nil
	}
	return selector
}

func (rfq *ReportFilterQuery) prepareQuery(ctx context.Context) error {
	if rfq.path != nil {
		prev, err := rfq.path(ctx)
		if err != nil {
			return err
		}
		rfq.sql = prev
	}
	if err := reportfilter.Policy.EvalQuery(ctx, rfq); err != nil {
		return err
	}
	return nil
}

func (rfq *ReportFilterQuery) sqlAll(ctx context.Context) ([]*ReportFilter, error) {
	var (
		nodes = []*ReportFilter{}
		_spec = rfq.querySpec()
	)
	_spec.ScanValues = func() []interface{} {
		node := &ReportFilter{config: rfq.config}
		nodes = append(nodes, node)
		values := node.scanValues()
		return values
	}
	_spec.Assign = func(values ...interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("ent: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		return node.assignValues(values...)
	}
	if err := sqlgraph.QueryNodes(ctx, rfq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (rfq *ReportFilterQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := rfq.querySpec()
	return sqlgraph.CountNodes(ctx, rfq.driver, _spec)
}

func (rfq *ReportFilterQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := rfq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (rfq *ReportFilterQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   reportfilter.Table,
			Columns: reportfilter.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: reportfilter.FieldID,
			},
		},
		From:   rfq.sql,
		Unique: true,
	}
	if ps := rfq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := rfq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := rfq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := rfq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (rfq *ReportFilterQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(rfq.driver.Dialect())
	t1 := builder.Table(reportfilter.Table)
	selector := builder.Select(t1.Columns(reportfilter.Columns...)...).From(t1)
	if rfq.sql != nil {
		selector = rfq.sql
		selector.Select(selector.Columns(reportfilter.Columns...)...)
	}
	for _, p := range rfq.predicates {
		p(selector)
	}
	for _, p := range rfq.order {
		p(selector)
	}
	if offset := rfq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := rfq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ReportFilterGroupBy is the builder for group-by ReportFilter entities.
type ReportFilterGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (rfgb *ReportFilterGroupBy) Aggregate(fns ...Aggregate) *ReportFilterGroupBy {
	rfgb.fns = append(rfgb.fns, fns...)
	return rfgb
}

// Scan applies the group-by query and scan the result into the given value.
func (rfgb *ReportFilterGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := rfgb.path(ctx)
	if err != nil {
		return err
	}
	rfgb.sql = query
	return rfgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (rfgb *ReportFilterGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := rfgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (rfgb *ReportFilterGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(rfgb.fields) > 1 {
		return nil, errors.New("ent: ReportFilterGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := rfgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (rfgb *ReportFilterGroupBy) StringsX(ctx context.Context) []string {
	v, err := rfgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (rfgb *ReportFilterGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(rfgb.fields) > 1 {
		return nil, errors.New("ent: ReportFilterGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := rfgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (rfgb *ReportFilterGroupBy) IntsX(ctx context.Context) []int {
	v, err := rfgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (rfgb *ReportFilterGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(rfgb.fields) > 1 {
		return nil, errors.New("ent: ReportFilterGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := rfgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (rfgb *ReportFilterGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := rfgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (rfgb *ReportFilterGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(rfgb.fields) > 1 {
		return nil, errors.New("ent: ReportFilterGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := rfgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (rfgb *ReportFilterGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := rfgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (rfgb *ReportFilterGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := rfgb.sqlQuery().Query()
	if err := rfgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (rfgb *ReportFilterGroupBy) sqlQuery() *sql.Selector {
	selector := rfgb.sql
	columns := make([]string, 0, len(rfgb.fields)+len(rfgb.fns))
	columns = append(columns, rfgb.fields...)
	for _, fn := range rfgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(rfgb.fields...)
}

// ReportFilterSelect is the builder for select fields of ReportFilter entities.
type ReportFilterSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (rfs *ReportFilterSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := rfs.path(ctx)
	if err != nil {
		return err
	}
	rfs.sql = query
	return rfs.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (rfs *ReportFilterSelect) ScanX(ctx context.Context, v interface{}) {
	if err := rfs.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (rfs *ReportFilterSelect) Strings(ctx context.Context) ([]string, error) {
	if len(rfs.fields) > 1 {
		return nil, errors.New("ent: ReportFilterSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := rfs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (rfs *ReportFilterSelect) StringsX(ctx context.Context) []string {
	v, err := rfs.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (rfs *ReportFilterSelect) Ints(ctx context.Context) ([]int, error) {
	if len(rfs.fields) > 1 {
		return nil, errors.New("ent: ReportFilterSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := rfs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (rfs *ReportFilterSelect) IntsX(ctx context.Context) []int {
	v, err := rfs.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (rfs *ReportFilterSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(rfs.fields) > 1 {
		return nil, errors.New("ent: ReportFilterSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := rfs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (rfs *ReportFilterSelect) Float64sX(ctx context.Context) []float64 {
	v, err := rfs.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (rfs *ReportFilterSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(rfs.fields) > 1 {
		return nil, errors.New("ent: ReportFilterSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := rfs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (rfs *ReportFilterSelect) BoolsX(ctx context.Context) []bool {
	v, err := rfs.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (rfs *ReportFilterSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := rfs.sqlQuery().Query()
	if err := rfs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (rfs *ReportFilterSelect) sqlQuery() sql.Querier {
	selector := rfs.sql
	selector.Select(selector.Columns(rfs.fields...)...)
	return selector
}
