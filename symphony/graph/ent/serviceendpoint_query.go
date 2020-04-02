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
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
)

// ServiceEndpointQuery is the builder for querying ServiceEndpoint entities.
type ServiceEndpointQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.ServiceEndpoint
	// eager-loading edges.
	withPort    *EquipmentPortQuery
	withService *ServiceQuery
	withFKs     bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (seq *ServiceEndpointQuery) Where(ps ...predicate.ServiceEndpoint) *ServiceEndpointQuery {
	seq.predicates = append(seq.predicates, ps...)
	return seq
}

// Limit adds a limit step to the query.
func (seq *ServiceEndpointQuery) Limit(limit int) *ServiceEndpointQuery {
	seq.limit = &limit
	return seq
}

// Offset adds an offset step to the query.
func (seq *ServiceEndpointQuery) Offset(offset int) *ServiceEndpointQuery {
	seq.offset = &offset
	return seq
}

// Order adds an order step to the query.
func (seq *ServiceEndpointQuery) Order(o ...Order) *ServiceEndpointQuery {
	seq.order = append(seq.order, o...)
	return seq
}

// QueryPort chains the current query on the port edge.
func (seq *ServiceEndpointQuery) QueryPort() *EquipmentPortQuery {
	query := &EquipmentPortQuery{config: seq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := seq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(serviceendpoint.Table, serviceendpoint.FieldID, seq.sqlQuery()),
			sqlgraph.To(equipmentport.Table, equipmentport.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, serviceendpoint.PortTable, serviceendpoint.PortColumn),
		)
		fromU = sqlgraph.SetNeighbors(seq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryService chains the current query on the service edge.
func (seq *ServiceEndpointQuery) QueryService() *ServiceQuery {
	query := &ServiceQuery{config: seq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := seq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(serviceendpoint.Table, serviceendpoint.FieldID, seq.sqlQuery()),
			sqlgraph.To(service.Table, service.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, serviceendpoint.ServiceTable, serviceendpoint.ServiceColumn),
		)
		fromU = sqlgraph.SetNeighbors(seq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first ServiceEndpoint entity in the query. Returns *NotFoundError when no serviceendpoint was found.
func (seq *ServiceEndpointQuery) First(ctx context.Context) (*ServiceEndpoint, error) {
	ses, err := seq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(ses) == 0 {
		return nil, &NotFoundError{serviceendpoint.Label}
	}
	return ses[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (seq *ServiceEndpointQuery) FirstX(ctx context.Context) *ServiceEndpoint {
	se, err := seq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return se
}

// FirstID returns the first ServiceEndpoint id in the query. Returns *NotFoundError when no id was found.
func (seq *ServiceEndpointQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = seq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{serviceendpoint.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (seq *ServiceEndpointQuery) FirstXID(ctx context.Context) int {
	id, err := seq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only ServiceEndpoint entity in the query, returns an error if not exactly one entity was returned.
func (seq *ServiceEndpointQuery) Only(ctx context.Context) (*ServiceEndpoint, error) {
	ses, err := seq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(ses) {
	case 1:
		return ses[0], nil
	case 0:
		return nil, &NotFoundError{serviceendpoint.Label}
	default:
		return nil, &NotSingularError{serviceendpoint.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (seq *ServiceEndpointQuery) OnlyX(ctx context.Context) *ServiceEndpoint {
	se, err := seq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return se
}

// OnlyID returns the only ServiceEndpoint id in the query, returns an error if not exactly one id was returned.
func (seq *ServiceEndpointQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = seq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{serviceendpoint.Label}
	default:
		err = &NotSingularError{serviceendpoint.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (seq *ServiceEndpointQuery) OnlyXID(ctx context.Context) int {
	id, err := seq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ServiceEndpoints.
func (seq *ServiceEndpointQuery) All(ctx context.Context) ([]*ServiceEndpoint, error) {
	if err := seq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return seq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (seq *ServiceEndpointQuery) AllX(ctx context.Context) []*ServiceEndpoint {
	ses, err := seq.All(ctx)
	if err != nil {
		panic(err)
	}
	return ses
}

// IDs executes the query and returns a list of ServiceEndpoint ids.
func (seq *ServiceEndpointQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := seq.Select(serviceendpoint.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (seq *ServiceEndpointQuery) IDsX(ctx context.Context) []int {
	ids, err := seq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (seq *ServiceEndpointQuery) Count(ctx context.Context) (int, error) {
	if err := seq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return seq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (seq *ServiceEndpointQuery) CountX(ctx context.Context) int {
	count, err := seq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (seq *ServiceEndpointQuery) Exist(ctx context.Context) (bool, error) {
	if err := seq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return seq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (seq *ServiceEndpointQuery) ExistX(ctx context.Context) bool {
	exist, err := seq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (seq *ServiceEndpointQuery) Clone() *ServiceEndpointQuery {
	return &ServiceEndpointQuery{
		config:     seq.config,
		limit:      seq.limit,
		offset:     seq.offset,
		order:      append([]Order{}, seq.order...),
		unique:     append([]string{}, seq.unique...),
		predicates: append([]predicate.ServiceEndpoint{}, seq.predicates...),
		// clone intermediate query.
		sql:  seq.sql.Clone(),
		path: seq.path,
	}
}

//  WithPort tells the query-builder to eager-loads the nodes that are connected to
// the "port" edge. The optional arguments used to configure the query builder of the edge.
func (seq *ServiceEndpointQuery) WithPort(opts ...func(*EquipmentPortQuery)) *ServiceEndpointQuery {
	query := &EquipmentPortQuery{config: seq.config}
	for _, opt := range opts {
		opt(query)
	}
	seq.withPort = query
	return seq
}

//  WithService tells the query-builder to eager-loads the nodes that are connected to
// the "service" edge. The optional arguments used to configure the query builder of the edge.
func (seq *ServiceEndpointQuery) WithService(opts ...func(*ServiceQuery)) *ServiceEndpointQuery {
	query := &ServiceQuery{config: seq.config}
	for _, opt := range opts {
		opt(query)
	}
	seq.withService = query
	return seq
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
//	client.ServiceEndpoint.Query().
//		GroupBy(serviceendpoint.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (seq *ServiceEndpointQuery) GroupBy(field string, fields ...string) *ServiceEndpointGroupBy {
	group := &ServiceEndpointGroupBy{config: seq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := seq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return seq.sqlQuery(), nil
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
//	client.ServiceEndpoint.Query().
//		Select(serviceendpoint.FieldCreateTime).
//		Scan(ctx, &v)
//
func (seq *ServiceEndpointQuery) Select(field string, fields ...string) *ServiceEndpointSelect {
	selector := &ServiceEndpointSelect{config: seq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := seq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return seq.sqlQuery(), nil
	}
	return selector
}

func (seq *ServiceEndpointQuery) prepareQuery(ctx context.Context) error {
	if seq.path != nil {
		prev, err := seq.path(ctx)
		if err != nil {
			return err
		}
		seq.sql = prev
	}
	return nil
}

func (seq *ServiceEndpointQuery) sqlAll(ctx context.Context) ([]*ServiceEndpoint, error) {
	var (
		nodes       = []*ServiceEndpoint{}
		withFKs     = seq.withFKs
		_spec       = seq.querySpec()
		loadedTypes = [2]bool{
			seq.withPort != nil,
			seq.withService != nil,
		}
	)
	if seq.withPort != nil || seq.withService != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, serviceendpoint.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &ServiceEndpoint{config: seq.config}
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
	if err := sqlgraph.QueryNodes(ctx, seq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := seq.withPort; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*ServiceEndpoint)
		for i := range nodes {
			if fk := nodes[i].service_endpoint_port; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(equipmentport.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "service_endpoint_port" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Port = n
			}
		}
	}

	if query := seq.withService; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*ServiceEndpoint)
		for i := range nodes {
			if fk := nodes[i].service_endpoints; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(service.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "service_endpoints" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Service = n
			}
		}
	}

	return nodes, nil
}

func (seq *ServiceEndpointQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := seq.querySpec()
	return sqlgraph.CountNodes(ctx, seq.driver, _spec)
}

func (seq *ServiceEndpointQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := seq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (seq *ServiceEndpointQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   serviceendpoint.Table,
			Columns: serviceendpoint.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: serviceendpoint.FieldID,
			},
		},
		From:   seq.sql,
		Unique: true,
	}
	if ps := seq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := seq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := seq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := seq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (seq *ServiceEndpointQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(seq.driver.Dialect())
	t1 := builder.Table(serviceendpoint.Table)
	selector := builder.Select(t1.Columns(serviceendpoint.Columns...)...).From(t1)
	if seq.sql != nil {
		selector = seq.sql
		selector.Select(selector.Columns(serviceendpoint.Columns...)...)
	}
	for _, p := range seq.predicates {
		p(selector)
	}
	for _, p := range seq.order {
		p(selector)
	}
	if offset := seq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := seq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ServiceEndpointGroupBy is the builder for group-by ServiceEndpoint entities.
type ServiceEndpointGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (segb *ServiceEndpointGroupBy) Aggregate(fns ...Aggregate) *ServiceEndpointGroupBy {
	segb.fns = append(segb.fns, fns...)
	return segb
}

// Scan applies the group-by query and scan the result into the given value.
func (segb *ServiceEndpointGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := segb.path(ctx)
	if err != nil {
		return err
	}
	segb.sql = query
	return segb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (segb *ServiceEndpointGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := segb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (segb *ServiceEndpointGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(segb.fields) > 1 {
		return nil, errors.New("ent: ServiceEndpointGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := segb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (segb *ServiceEndpointGroupBy) StringsX(ctx context.Context) []string {
	v, err := segb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (segb *ServiceEndpointGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(segb.fields) > 1 {
		return nil, errors.New("ent: ServiceEndpointGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := segb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (segb *ServiceEndpointGroupBy) IntsX(ctx context.Context) []int {
	v, err := segb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (segb *ServiceEndpointGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(segb.fields) > 1 {
		return nil, errors.New("ent: ServiceEndpointGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := segb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (segb *ServiceEndpointGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := segb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (segb *ServiceEndpointGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(segb.fields) > 1 {
		return nil, errors.New("ent: ServiceEndpointGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := segb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (segb *ServiceEndpointGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := segb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (segb *ServiceEndpointGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := segb.sqlQuery().Query()
	if err := segb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (segb *ServiceEndpointGroupBy) sqlQuery() *sql.Selector {
	selector := segb.sql
	columns := make([]string, 0, len(segb.fields)+len(segb.fns))
	columns = append(columns, segb.fields...)
	for _, fn := range segb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(segb.fields...)
}

// ServiceEndpointSelect is the builder for select fields of ServiceEndpoint entities.
type ServiceEndpointSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (ses *ServiceEndpointSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := ses.path(ctx)
	if err != nil {
		return err
	}
	ses.sql = query
	return ses.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ses *ServiceEndpointSelect) ScanX(ctx context.Context, v interface{}) {
	if err := ses.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (ses *ServiceEndpointSelect) Strings(ctx context.Context) ([]string, error) {
	if len(ses.fields) > 1 {
		return nil, errors.New("ent: ServiceEndpointSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := ses.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ses *ServiceEndpointSelect) StringsX(ctx context.Context) []string {
	v, err := ses.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (ses *ServiceEndpointSelect) Ints(ctx context.Context) ([]int, error) {
	if len(ses.fields) > 1 {
		return nil, errors.New("ent: ServiceEndpointSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := ses.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ses *ServiceEndpointSelect) IntsX(ctx context.Context) []int {
	v, err := ses.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (ses *ServiceEndpointSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(ses.fields) > 1 {
		return nil, errors.New("ent: ServiceEndpointSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := ses.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ses *ServiceEndpointSelect) Float64sX(ctx context.Context) []float64 {
	v, err := ses.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (ses *ServiceEndpointSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(ses.fields) > 1 {
		return nil, errors.New("ent: ServiceEndpointSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := ses.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ses *ServiceEndpointSelect) BoolsX(ctx context.Context) []bool {
	v, err := ses.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ses *ServiceEndpointSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ses.sqlQuery().Query()
	if err := ses.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ses *ServiceEndpointSelect) sqlQuery() sql.Querier {
	selector := ses.sql
	selector.Select(selector.Columns(ses.fields...)...)
	return selector
}
