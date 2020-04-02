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
	"github.com/facebookincubator/symphony/frontier/ent/auditlog"
	"github.com/facebookincubator/symphony/frontier/ent/predicate"
)

// AuditLogQuery is the builder for querying AuditLog entities.
type AuditLogQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.AuditLog
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (alq *AuditLogQuery) Where(ps ...predicate.AuditLog) *AuditLogQuery {
	alq.predicates = append(alq.predicates, ps...)
	return alq
}

// Limit adds a limit step to the query.
func (alq *AuditLogQuery) Limit(limit int) *AuditLogQuery {
	alq.limit = &limit
	return alq
}

// Offset adds an offset step to the query.
func (alq *AuditLogQuery) Offset(offset int) *AuditLogQuery {
	alq.offset = &offset
	return alq
}

// Order adds an order step to the query.
func (alq *AuditLogQuery) Order(o ...Order) *AuditLogQuery {
	alq.order = append(alq.order, o...)
	return alq
}

// First returns the first AuditLog entity in the query. Returns *NotFoundError when no auditlog was found.
func (alq *AuditLogQuery) First(ctx context.Context) (*AuditLog, error) {
	als, err := alq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(als) == 0 {
		return nil, &NotFoundError{auditlog.Label}
	}
	return als[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (alq *AuditLogQuery) FirstX(ctx context.Context) *AuditLog {
	al, err := alq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return al
}

// FirstID returns the first AuditLog id in the query. Returns *NotFoundError when no id was found.
func (alq *AuditLogQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = alq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{auditlog.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (alq *AuditLogQuery) FirstXID(ctx context.Context) int {
	id, err := alq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only AuditLog entity in the query, returns an error if not exactly one entity was returned.
func (alq *AuditLogQuery) Only(ctx context.Context) (*AuditLog, error) {
	als, err := alq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(als) {
	case 1:
		return als[0], nil
	case 0:
		return nil, &NotFoundError{auditlog.Label}
	default:
		return nil, &NotSingularError{auditlog.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (alq *AuditLogQuery) OnlyX(ctx context.Context) *AuditLog {
	al, err := alq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return al
}

// OnlyID returns the only AuditLog id in the query, returns an error if not exactly one id was returned.
func (alq *AuditLogQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = alq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{auditlog.Label}
	default:
		err = &NotSingularError{auditlog.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (alq *AuditLogQuery) OnlyXID(ctx context.Context) int {
	id, err := alq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of AuditLogs.
func (alq *AuditLogQuery) All(ctx context.Context) ([]*AuditLog, error) {
	if err := alq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return alq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (alq *AuditLogQuery) AllX(ctx context.Context) []*AuditLog {
	als, err := alq.All(ctx)
	if err != nil {
		panic(err)
	}
	return als
}

// IDs executes the query and returns a list of AuditLog ids.
func (alq *AuditLogQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := alq.Select(auditlog.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (alq *AuditLogQuery) IDsX(ctx context.Context) []int {
	ids, err := alq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (alq *AuditLogQuery) Count(ctx context.Context) (int, error) {
	if err := alq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return alq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (alq *AuditLogQuery) CountX(ctx context.Context) int {
	count, err := alq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (alq *AuditLogQuery) Exist(ctx context.Context) (bool, error) {
	if err := alq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return alq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (alq *AuditLogQuery) ExistX(ctx context.Context) bool {
	exist, err := alq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (alq *AuditLogQuery) Clone() *AuditLogQuery {
	return &AuditLogQuery{
		config:     alq.config,
		limit:      alq.limit,
		offset:     alq.offset,
		order:      append([]Order{}, alq.order...),
		unique:     append([]string{}, alq.unique...),
		predicates: append([]predicate.AuditLog{}, alq.predicates...),
		// clone intermediate query.
		sql:  alq.sql.Clone(),
		path: alq.path,
	}
}

// GroupBy used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"created_at,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.AuditLog.Query().
//		GroupBy(auditlog.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (alq *AuditLogQuery) GroupBy(field string, fields ...string) *AuditLogGroupBy {
	group := &AuditLogGroupBy{config: alq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := alq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return alq.sqlQuery(), nil
	}
	return group
}

// Select one or more fields from the given query.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"created_at,omitempty"`
//	}
//
//	client.AuditLog.Query().
//		Select(auditlog.FieldCreatedAt).
//		Scan(ctx, &v)
//
func (alq *AuditLogQuery) Select(field string, fields ...string) *AuditLogSelect {
	selector := &AuditLogSelect{config: alq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := alq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return alq.sqlQuery(), nil
	}
	return selector
}

func (alq *AuditLogQuery) prepareQuery(ctx context.Context) error {
	if alq.path != nil {
		prev, err := alq.path(ctx)
		if err != nil {
			return err
		}
		alq.sql = prev
	}
	return nil
}

func (alq *AuditLogQuery) sqlAll(ctx context.Context) ([]*AuditLog, error) {
	var (
		nodes = []*AuditLog{}
		_spec = alq.querySpec()
	)
	_spec.ScanValues = func() []interface{} {
		node := &AuditLog{config: alq.config}
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
	if err := sqlgraph.QueryNodes(ctx, alq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (alq *AuditLogQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := alq.querySpec()
	return sqlgraph.CountNodes(ctx, alq.driver, _spec)
}

func (alq *AuditLogQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := alq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (alq *AuditLogQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   auditlog.Table,
			Columns: auditlog.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: auditlog.FieldID,
			},
		},
		From:   alq.sql,
		Unique: true,
	}
	if ps := alq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := alq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := alq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := alq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (alq *AuditLogQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(alq.driver.Dialect())
	t1 := builder.Table(auditlog.Table)
	selector := builder.Select(t1.Columns(auditlog.Columns...)...).From(t1)
	if alq.sql != nil {
		selector = alq.sql
		selector.Select(selector.Columns(auditlog.Columns...)...)
	}
	for _, p := range alq.predicates {
		p(selector)
	}
	for _, p := range alq.order {
		p(selector)
	}
	if offset := alq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := alq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// AuditLogGroupBy is the builder for group-by AuditLog entities.
type AuditLogGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (algb *AuditLogGroupBy) Aggregate(fns ...Aggregate) *AuditLogGroupBy {
	algb.fns = append(algb.fns, fns...)
	return algb
}

// Scan applies the group-by query and scan the result into the given value.
func (algb *AuditLogGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := algb.path(ctx)
	if err != nil {
		return err
	}
	algb.sql = query
	return algb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (algb *AuditLogGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := algb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (algb *AuditLogGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(algb.fields) > 1 {
		return nil, errors.New("ent: AuditLogGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := algb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (algb *AuditLogGroupBy) StringsX(ctx context.Context) []string {
	v, err := algb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (algb *AuditLogGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(algb.fields) > 1 {
		return nil, errors.New("ent: AuditLogGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := algb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (algb *AuditLogGroupBy) IntsX(ctx context.Context) []int {
	v, err := algb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (algb *AuditLogGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(algb.fields) > 1 {
		return nil, errors.New("ent: AuditLogGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := algb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (algb *AuditLogGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := algb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (algb *AuditLogGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(algb.fields) > 1 {
		return nil, errors.New("ent: AuditLogGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := algb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (algb *AuditLogGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := algb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (algb *AuditLogGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := algb.sqlQuery().Query()
	if err := algb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (algb *AuditLogGroupBy) sqlQuery() *sql.Selector {
	selector := algb.sql
	columns := make([]string, 0, len(algb.fields)+len(algb.fns))
	columns = append(columns, algb.fields...)
	for _, fn := range algb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(algb.fields...)
}

// AuditLogSelect is the builder for select fields of AuditLog entities.
type AuditLogSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (als *AuditLogSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := als.path(ctx)
	if err != nil {
		return err
	}
	als.sql = query
	return als.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (als *AuditLogSelect) ScanX(ctx context.Context, v interface{}) {
	if err := als.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (als *AuditLogSelect) Strings(ctx context.Context) ([]string, error) {
	if len(als.fields) > 1 {
		return nil, errors.New("ent: AuditLogSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := als.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (als *AuditLogSelect) StringsX(ctx context.Context) []string {
	v, err := als.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (als *AuditLogSelect) Ints(ctx context.Context) ([]int, error) {
	if len(als.fields) > 1 {
		return nil, errors.New("ent: AuditLogSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := als.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (als *AuditLogSelect) IntsX(ctx context.Context) []int {
	v, err := als.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (als *AuditLogSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(als.fields) > 1 {
		return nil, errors.New("ent: AuditLogSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := als.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (als *AuditLogSelect) Float64sX(ctx context.Context) []float64 {
	v, err := als.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (als *AuditLogSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(als.fields) > 1 {
		return nil, errors.New("ent: AuditLogSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := als.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (als *AuditLogSelect) BoolsX(ctx context.Context) []bool {
	v, err := als.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (als *AuditLogSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := als.sqlQuery().Query()
	if err := als.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (als *AuditLogSelect) sqlQuery() sql.Querier {
	selector := als.sql
	selector.Select(selector.Columns(als.fields...)...)
	return selector
}
