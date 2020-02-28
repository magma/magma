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
	"github.com/facebookincubator/symphony/graph/ent/floorplanreferencepoint"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// FloorPlanReferencePointQuery is the builder for querying FloorPlanReferencePoint entities.
type FloorPlanReferencePointQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.FloorPlanReferencePoint
	// intermediate query.
	sql *sql.Selector
}

// Where adds a new predicate for the builder.
func (fprpq *FloorPlanReferencePointQuery) Where(ps ...predicate.FloorPlanReferencePoint) *FloorPlanReferencePointQuery {
	fprpq.predicates = append(fprpq.predicates, ps...)
	return fprpq
}

// Limit adds a limit step to the query.
func (fprpq *FloorPlanReferencePointQuery) Limit(limit int) *FloorPlanReferencePointQuery {
	fprpq.limit = &limit
	return fprpq
}

// Offset adds an offset step to the query.
func (fprpq *FloorPlanReferencePointQuery) Offset(offset int) *FloorPlanReferencePointQuery {
	fprpq.offset = &offset
	return fprpq
}

// Order adds an order step to the query.
func (fprpq *FloorPlanReferencePointQuery) Order(o ...Order) *FloorPlanReferencePointQuery {
	fprpq.order = append(fprpq.order, o...)
	return fprpq
}

// First returns the first FloorPlanReferencePoint entity in the query. Returns *NotFoundError when no floorplanreferencepoint was found.
func (fprpq *FloorPlanReferencePointQuery) First(ctx context.Context) (*FloorPlanReferencePoint, error) {
	fprps, err := fprpq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(fprps) == 0 {
		return nil, &NotFoundError{floorplanreferencepoint.Label}
	}
	return fprps[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (fprpq *FloorPlanReferencePointQuery) FirstX(ctx context.Context) *FloorPlanReferencePoint {
	fprp, err := fprpq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return fprp
}

// FirstID returns the first FloorPlanReferencePoint id in the query. Returns *NotFoundError when no id was found.
func (fprpq *FloorPlanReferencePointQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = fprpq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{floorplanreferencepoint.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (fprpq *FloorPlanReferencePointQuery) FirstXID(ctx context.Context) int {
	id, err := fprpq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only FloorPlanReferencePoint entity in the query, returns an error if not exactly one entity was returned.
func (fprpq *FloorPlanReferencePointQuery) Only(ctx context.Context) (*FloorPlanReferencePoint, error) {
	fprps, err := fprpq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(fprps) {
	case 1:
		return fprps[0], nil
	case 0:
		return nil, &NotFoundError{floorplanreferencepoint.Label}
	default:
		return nil, &NotSingularError{floorplanreferencepoint.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (fprpq *FloorPlanReferencePointQuery) OnlyX(ctx context.Context) *FloorPlanReferencePoint {
	fprp, err := fprpq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return fprp
}

// OnlyID returns the only FloorPlanReferencePoint id in the query, returns an error if not exactly one id was returned.
func (fprpq *FloorPlanReferencePointQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = fprpq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{floorplanreferencepoint.Label}
	default:
		err = &NotSingularError{floorplanreferencepoint.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (fprpq *FloorPlanReferencePointQuery) OnlyXID(ctx context.Context) int {
	id, err := fprpq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of FloorPlanReferencePoints.
func (fprpq *FloorPlanReferencePointQuery) All(ctx context.Context) ([]*FloorPlanReferencePoint, error) {
	return fprpq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (fprpq *FloorPlanReferencePointQuery) AllX(ctx context.Context) []*FloorPlanReferencePoint {
	fprps, err := fprpq.All(ctx)
	if err != nil {
		panic(err)
	}
	return fprps
}

// IDs executes the query and returns a list of FloorPlanReferencePoint ids.
func (fprpq *FloorPlanReferencePointQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := fprpq.Select(floorplanreferencepoint.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (fprpq *FloorPlanReferencePointQuery) IDsX(ctx context.Context) []int {
	ids, err := fprpq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (fprpq *FloorPlanReferencePointQuery) Count(ctx context.Context) (int, error) {
	return fprpq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (fprpq *FloorPlanReferencePointQuery) CountX(ctx context.Context) int {
	count, err := fprpq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (fprpq *FloorPlanReferencePointQuery) Exist(ctx context.Context) (bool, error) {
	return fprpq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (fprpq *FloorPlanReferencePointQuery) ExistX(ctx context.Context) bool {
	exist, err := fprpq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (fprpq *FloorPlanReferencePointQuery) Clone() *FloorPlanReferencePointQuery {
	return &FloorPlanReferencePointQuery{
		config:     fprpq.config,
		limit:      fprpq.limit,
		offset:     fprpq.offset,
		order:      append([]Order{}, fprpq.order...),
		unique:     append([]string{}, fprpq.unique...),
		predicates: append([]predicate.FloorPlanReferencePoint{}, fprpq.predicates...),
		// clone intermediate query.
		sql: fprpq.sql.Clone(),
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
//	client.FloorPlanReferencePoint.Query().
//		GroupBy(floorplanreferencepoint.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (fprpq *FloorPlanReferencePointQuery) GroupBy(field string, fields ...string) *FloorPlanReferencePointGroupBy {
	group := &FloorPlanReferencePointGroupBy{config: fprpq.config}
	group.fields = append([]string{field}, fields...)
	group.sql = fprpq.sqlQuery()
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
//	client.FloorPlanReferencePoint.Query().
//		Select(floorplanreferencepoint.FieldCreateTime).
//		Scan(ctx, &v)
//
func (fprpq *FloorPlanReferencePointQuery) Select(field string, fields ...string) *FloorPlanReferencePointSelect {
	selector := &FloorPlanReferencePointSelect{config: fprpq.config}
	selector.fields = append([]string{field}, fields...)
	selector.sql = fprpq.sqlQuery()
	return selector
}

func (fprpq *FloorPlanReferencePointQuery) sqlAll(ctx context.Context) ([]*FloorPlanReferencePoint, error) {
	var (
		nodes = []*FloorPlanReferencePoint{}
		_spec = fprpq.querySpec()
	)
	_spec.ScanValues = func() []interface{} {
		node := &FloorPlanReferencePoint{config: fprpq.config}
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
	if err := sqlgraph.QueryNodes(ctx, fprpq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (fprpq *FloorPlanReferencePointQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := fprpq.querySpec()
	return sqlgraph.CountNodes(ctx, fprpq.driver, _spec)
}

func (fprpq *FloorPlanReferencePointQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := fprpq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (fprpq *FloorPlanReferencePointQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   floorplanreferencepoint.Table,
			Columns: floorplanreferencepoint.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: floorplanreferencepoint.FieldID,
			},
		},
		From:   fprpq.sql,
		Unique: true,
	}
	if ps := fprpq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := fprpq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := fprpq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := fprpq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (fprpq *FloorPlanReferencePointQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(fprpq.driver.Dialect())
	t1 := builder.Table(floorplanreferencepoint.Table)
	selector := builder.Select(t1.Columns(floorplanreferencepoint.Columns...)...).From(t1)
	if fprpq.sql != nil {
		selector = fprpq.sql
		selector.Select(selector.Columns(floorplanreferencepoint.Columns...)...)
	}
	for _, p := range fprpq.predicates {
		p(selector)
	}
	for _, p := range fprpq.order {
		p(selector)
	}
	if offset := fprpq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := fprpq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// FloorPlanReferencePointGroupBy is the builder for group-by FloorPlanReferencePoint entities.
type FloorPlanReferencePointGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query.
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the group-by query.
func (fprpgb *FloorPlanReferencePointGroupBy) Aggregate(fns ...Aggregate) *FloorPlanReferencePointGroupBy {
	fprpgb.fns = append(fprpgb.fns, fns...)
	return fprpgb
}

// Scan applies the group-by query and scan the result into the given value.
func (fprpgb *FloorPlanReferencePointGroupBy) Scan(ctx context.Context, v interface{}) error {
	return fprpgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (fprpgb *FloorPlanReferencePointGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := fprpgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (fprpgb *FloorPlanReferencePointGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(fprpgb.fields) > 1 {
		return nil, errors.New("ent: FloorPlanReferencePointGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := fprpgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (fprpgb *FloorPlanReferencePointGroupBy) StringsX(ctx context.Context) []string {
	v, err := fprpgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (fprpgb *FloorPlanReferencePointGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(fprpgb.fields) > 1 {
		return nil, errors.New("ent: FloorPlanReferencePointGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := fprpgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (fprpgb *FloorPlanReferencePointGroupBy) IntsX(ctx context.Context) []int {
	v, err := fprpgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (fprpgb *FloorPlanReferencePointGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(fprpgb.fields) > 1 {
		return nil, errors.New("ent: FloorPlanReferencePointGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := fprpgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (fprpgb *FloorPlanReferencePointGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := fprpgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (fprpgb *FloorPlanReferencePointGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(fprpgb.fields) > 1 {
		return nil, errors.New("ent: FloorPlanReferencePointGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := fprpgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (fprpgb *FloorPlanReferencePointGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := fprpgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (fprpgb *FloorPlanReferencePointGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := fprpgb.sqlQuery().Query()
	if err := fprpgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (fprpgb *FloorPlanReferencePointGroupBy) sqlQuery() *sql.Selector {
	selector := fprpgb.sql
	columns := make([]string, 0, len(fprpgb.fields)+len(fprpgb.fns))
	columns = append(columns, fprpgb.fields...)
	for _, fn := range fprpgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(fprpgb.fields...)
}

// FloorPlanReferencePointSelect is the builder for select fields of FloorPlanReferencePoint entities.
type FloorPlanReferencePointSelect struct {
	config
	fields []string
	// intermediate queries.
	sql *sql.Selector
}

// Scan applies the selector query and scan the result into the given value.
func (fprps *FloorPlanReferencePointSelect) Scan(ctx context.Context, v interface{}) error {
	return fprps.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (fprps *FloorPlanReferencePointSelect) ScanX(ctx context.Context, v interface{}) {
	if err := fprps.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (fprps *FloorPlanReferencePointSelect) Strings(ctx context.Context) ([]string, error) {
	if len(fprps.fields) > 1 {
		return nil, errors.New("ent: FloorPlanReferencePointSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := fprps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (fprps *FloorPlanReferencePointSelect) StringsX(ctx context.Context) []string {
	v, err := fprps.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (fprps *FloorPlanReferencePointSelect) Ints(ctx context.Context) ([]int, error) {
	if len(fprps.fields) > 1 {
		return nil, errors.New("ent: FloorPlanReferencePointSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := fprps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (fprps *FloorPlanReferencePointSelect) IntsX(ctx context.Context) []int {
	v, err := fprps.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (fprps *FloorPlanReferencePointSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(fprps.fields) > 1 {
		return nil, errors.New("ent: FloorPlanReferencePointSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := fprps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (fprps *FloorPlanReferencePointSelect) Float64sX(ctx context.Context) []float64 {
	v, err := fprps.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (fprps *FloorPlanReferencePointSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(fprps.fields) > 1 {
		return nil, errors.New("ent: FloorPlanReferencePointSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := fprps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (fprps *FloorPlanReferencePointSelect) BoolsX(ctx context.Context) []bool {
	v, err := fprps.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (fprps *FloorPlanReferencePointSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := fprps.sqlQuery().Query()
	if err := fprps.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (fprps *FloorPlanReferencePointSelect) sqlQuery() sql.Querier {
	selector := fprps.sql
	selector.Select(selector.Columns(fprps.fields...)...)
	return selector
}
