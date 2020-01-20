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
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// LinkQuery is the builder for querying Link entities.
type LinkQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.Link
	// intermediate query.
	sql *sql.Selector
}

// Where adds a new predicate for the builder.
func (lq *LinkQuery) Where(ps ...predicate.Link) *LinkQuery {
	lq.predicates = append(lq.predicates, ps...)
	return lq
}

// Limit adds a limit step to the query.
func (lq *LinkQuery) Limit(limit int) *LinkQuery {
	lq.limit = &limit
	return lq
}

// Offset adds an offset step to the query.
func (lq *LinkQuery) Offset(offset int) *LinkQuery {
	lq.offset = &offset
	return lq
}

// Order adds an order step to the query.
func (lq *LinkQuery) Order(o ...Order) *LinkQuery {
	lq.order = append(lq.order, o...)
	return lq
}

// QueryPorts chains the current query on the ports edge.
func (lq *LinkQuery) QueryPorts() *EquipmentPortQuery {
	query := &EquipmentPortQuery{config: lq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(link.Table, link.FieldID, lq.sqlQuery()),
		sqlgraph.To(equipmentport.Table, equipmentport.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, link.PortsTable, link.PortsColumn),
	)
	query.sql = sqlgraph.SetNeighbors(lq.driver.Dialect(), step)
	return query
}

// QueryWorkOrder chains the current query on the work_order edge.
func (lq *LinkQuery) QueryWorkOrder() *WorkOrderQuery {
	query := &WorkOrderQuery{config: lq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(link.Table, link.FieldID, lq.sqlQuery()),
		sqlgraph.To(workorder.Table, workorder.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, link.WorkOrderTable, link.WorkOrderColumn),
	)
	query.sql = sqlgraph.SetNeighbors(lq.driver.Dialect(), step)
	return query
}

// QueryProperties chains the current query on the properties edge.
func (lq *LinkQuery) QueryProperties() *PropertyQuery {
	query := &PropertyQuery{config: lq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(link.Table, link.FieldID, lq.sqlQuery()),
		sqlgraph.To(property.Table, property.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, link.PropertiesTable, link.PropertiesColumn),
	)
	query.sql = sqlgraph.SetNeighbors(lq.driver.Dialect(), step)
	return query
}

// QueryService chains the current query on the service edge.
func (lq *LinkQuery) QueryService() *ServiceQuery {
	query := &ServiceQuery{config: lq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(link.Table, link.FieldID, lq.sqlQuery()),
		sqlgraph.To(service.Table, service.FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, link.ServiceTable, link.ServicePrimaryKey...),
	)
	query.sql = sqlgraph.SetNeighbors(lq.driver.Dialect(), step)
	return query
}

// First returns the first Link entity in the query. Returns *ErrNotFound when no link was found.
func (lq *LinkQuery) First(ctx context.Context) (*Link, error) {
	ls, err := lq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(ls) == 0 {
		return nil, &ErrNotFound{link.Label}
	}
	return ls[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (lq *LinkQuery) FirstX(ctx context.Context) *Link {
	l, err := lq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return l
}

// FirstID returns the first Link id in the query. Returns *ErrNotFound when no id was found.
func (lq *LinkQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = lq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &ErrNotFound{link.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (lq *LinkQuery) FirstXID(ctx context.Context) string {
	id, err := lq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only Link entity in the query, returns an error if not exactly one entity was returned.
func (lq *LinkQuery) Only(ctx context.Context) (*Link, error) {
	ls, err := lq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(ls) {
	case 1:
		return ls[0], nil
	case 0:
		return nil, &ErrNotFound{link.Label}
	default:
		return nil, &ErrNotSingular{link.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (lq *LinkQuery) OnlyX(ctx context.Context) *Link {
	l, err := lq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return l
}

// OnlyID returns the only Link id in the query, returns an error if not exactly one id was returned.
func (lq *LinkQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = lq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &ErrNotFound{link.Label}
	default:
		err = &ErrNotSingular{link.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (lq *LinkQuery) OnlyXID(ctx context.Context) string {
	id, err := lq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Links.
func (lq *LinkQuery) All(ctx context.Context) ([]*Link, error) {
	return lq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (lq *LinkQuery) AllX(ctx context.Context) []*Link {
	ls, err := lq.All(ctx)
	if err != nil {
		panic(err)
	}
	return ls
}

// IDs executes the query and returns a list of Link ids.
func (lq *LinkQuery) IDs(ctx context.Context) ([]string, error) {
	var ids []string
	if err := lq.Select(link.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (lq *LinkQuery) IDsX(ctx context.Context) []string {
	ids, err := lq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (lq *LinkQuery) Count(ctx context.Context) (int, error) {
	return lq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (lq *LinkQuery) CountX(ctx context.Context) int {
	count, err := lq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (lq *LinkQuery) Exist(ctx context.Context) (bool, error) {
	return lq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (lq *LinkQuery) ExistX(ctx context.Context) bool {
	exist, err := lq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (lq *LinkQuery) Clone() *LinkQuery {
	return &LinkQuery{
		config:     lq.config,
		limit:      lq.limit,
		offset:     lq.offset,
		order:      append([]Order{}, lq.order...),
		unique:     append([]string{}, lq.unique...),
		predicates: append([]predicate.Link{}, lq.predicates...),
		// clone intermediate query.
		sql: lq.sql.Clone(),
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
//	client.Link.Query().
//		GroupBy(link.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (lq *LinkQuery) GroupBy(field string, fields ...string) *LinkGroupBy {
	group := &LinkGroupBy{config: lq.config}
	group.fields = append([]string{field}, fields...)
	group.sql = lq.sqlQuery()
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
//	client.Link.Query().
//		Select(link.FieldCreateTime).
//		Scan(ctx, &v)
//
func (lq *LinkQuery) Select(field string, fields ...string) *LinkSelect {
	selector := &LinkSelect{config: lq.config}
	selector.fields = append([]string{field}, fields...)
	selector.sql = lq.sqlQuery()
	return selector
}

func (lq *LinkQuery) sqlAll(ctx context.Context) ([]*Link, error) {
	var (
		nodes []*Link
		spec  = lq.querySpec()
	)
	spec.ScanValues = func() []interface{} {
		node := &Link{config: lq.config}
		nodes = append(nodes, node)
		return node.scanValues()
	}
	spec.Assign = func(values ...interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("ent: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		return node.assignValues(values...)
	}
	if err := sqlgraph.QueryNodes(ctx, lq.driver, spec); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (lq *LinkQuery) sqlCount(ctx context.Context) (int, error) {
	spec := lq.querySpec()
	return sqlgraph.CountNodes(ctx, lq.driver, spec)
}

func (lq *LinkQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := lq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (lq *LinkQuery) querySpec() *sqlgraph.QuerySpec {
	spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   link.Table,
			Columns: link.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: link.FieldID,
			},
		},
		From:   lq.sql,
		Unique: true,
	}
	if ps := lq.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := lq.limit; limit != nil {
		spec.Limit = *limit
	}
	if offset := lq.offset; offset != nil {
		spec.Offset = *offset
	}
	if ps := lq.order; len(ps) > 0 {
		spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return spec
}

func (lq *LinkQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(lq.driver.Dialect())
	t1 := builder.Table(link.Table)
	selector := builder.Select(t1.Columns(link.Columns...)...).From(t1)
	if lq.sql != nil {
		selector = lq.sql
		selector.Select(selector.Columns(link.Columns...)...)
	}
	for _, p := range lq.predicates {
		p(selector)
	}
	for _, p := range lq.order {
		p(selector)
	}
	if offset := lq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := lq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// LinkGroupBy is the builder for group-by Link entities.
type LinkGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query.
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the group-by query.
func (lgb *LinkGroupBy) Aggregate(fns ...Aggregate) *LinkGroupBy {
	lgb.fns = append(lgb.fns, fns...)
	return lgb
}

// Scan applies the group-by query and scan the result into the given value.
func (lgb *LinkGroupBy) Scan(ctx context.Context, v interface{}) error {
	return lgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (lgb *LinkGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := lgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (lgb *LinkGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(lgb.fields) > 1 {
		return nil, errors.New("ent: LinkGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := lgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (lgb *LinkGroupBy) StringsX(ctx context.Context) []string {
	v, err := lgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (lgb *LinkGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(lgb.fields) > 1 {
		return nil, errors.New("ent: LinkGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := lgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (lgb *LinkGroupBy) IntsX(ctx context.Context) []int {
	v, err := lgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (lgb *LinkGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(lgb.fields) > 1 {
		return nil, errors.New("ent: LinkGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := lgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (lgb *LinkGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := lgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (lgb *LinkGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(lgb.fields) > 1 {
		return nil, errors.New("ent: LinkGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := lgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (lgb *LinkGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := lgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (lgb *LinkGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := lgb.sqlQuery().Query()
	if err := lgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (lgb *LinkGroupBy) sqlQuery() *sql.Selector {
	selector := lgb.sql
	columns := make([]string, 0, len(lgb.fields)+len(lgb.fns))
	columns = append(columns, lgb.fields...)
	for _, fn := range lgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(lgb.fields...)
}

// LinkSelect is the builder for select fields of Link entities.
type LinkSelect struct {
	config
	fields []string
	// intermediate queries.
	sql *sql.Selector
}

// Scan applies the selector query and scan the result into the given value.
func (ls *LinkSelect) Scan(ctx context.Context, v interface{}) error {
	return ls.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ls *LinkSelect) ScanX(ctx context.Context, v interface{}) {
	if err := ls.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (ls *LinkSelect) Strings(ctx context.Context) ([]string, error) {
	if len(ls.fields) > 1 {
		return nil, errors.New("ent: LinkSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := ls.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ls *LinkSelect) StringsX(ctx context.Context) []string {
	v, err := ls.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (ls *LinkSelect) Ints(ctx context.Context) ([]int, error) {
	if len(ls.fields) > 1 {
		return nil, errors.New("ent: LinkSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := ls.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ls *LinkSelect) IntsX(ctx context.Context) []int {
	v, err := ls.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (ls *LinkSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(ls.fields) > 1 {
		return nil, errors.New("ent: LinkSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := ls.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ls *LinkSelect) Float64sX(ctx context.Context) []float64 {
	v, err := ls.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (ls *LinkSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(ls.fields) > 1 {
		return nil, errors.New("ent: LinkSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := ls.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ls *LinkSelect) BoolsX(ctx context.Context) []bool {
	v, err := ls.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ls *LinkSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ls.sqlQuery().Query()
	if err := ls.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ls *LinkSelect) sqlQuery() sql.Querier {
	selector := ls.sql
	selector.Select(selector.Columns(ls.fields...)...)
	return selector
}
