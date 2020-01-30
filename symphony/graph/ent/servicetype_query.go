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
	"strconv"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
)

// ServiceTypeQuery is the builder for querying ServiceType entities.
type ServiceTypeQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.ServiceType
	// eager-loading edges.
	withServices      *ServiceQuery
	withPropertyTypes *PropertyTypeQuery
	// intermediate query.
	sql *sql.Selector
}

// Where adds a new predicate for the builder.
func (stq *ServiceTypeQuery) Where(ps ...predicate.ServiceType) *ServiceTypeQuery {
	stq.predicates = append(stq.predicates, ps...)
	return stq
}

// Limit adds a limit step to the query.
func (stq *ServiceTypeQuery) Limit(limit int) *ServiceTypeQuery {
	stq.limit = &limit
	return stq
}

// Offset adds an offset step to the query.
func (stq *ServiceTypeQuery) Offset(offset int) *ServiceTypeQuery {
	stq.offset = &offset
	return stq
}

// Order adds an order step to the query.
func (stq *ServiceTypeQuery) Order(o ...Order) *ServiceTypeQuery {
	stq.order = append(stq.order, o...)
	return stq
}

// QueryServices chains the current query on the services edge.
func (stq *ServiceTypeQuery) QueryServices() *ServiceQuery {
	query := &ServiceQuery{config: stq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(servicetype.Table, servicetype.FieldID, stq.sqlQuery()),
		sqlgraph.To(service.Table, service.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, servicetype.ServicesTable, servicetype.ServicesColumn),
	)
	query.sql = sqlgraph.SetNeighbors(stq.driver.Dialect(), step)
	return query
}

// QueryPropertyTypes chains the current query on the property_types edge.
func (stq *ServiceTypeQuery) QueryPropertyTypes() *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: stq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(servicetype.Table, servicetype.FieldID, stq.sqlQuery()),
		sqlgraph.To(propertytype.Table, propertytype.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, servicetype.PropertyTypesTable, servicetype.PropertyTypesColumn),
	)
	query.sql = sqlgraph.SetNeighbors(stq.driver.Dialect(), step)
	return query
}

// First returns the first ServiceType entity in the query. Returns *NotFoundError when no servicetype was found.
func (stq *ServiceTypeQuery) First(ctx context.Context) (*ServiceType, error) {
	sts, err := stq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(sts) == 0 {
		return nil, &NotFoundError{servicetype.Label}
	}
	return sts[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (stq *ServiceTypeQuery) FirstX(ctx context.Context) *ServiceType {
	st, err := stq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return st
}

// FirstID returns the first ServiceType id in the query. Returns *NotFoundError when no id was found.
func (stq *ServiceTypeQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = stq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{servicetype.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (stq *ServiceTypeQuery) FirstXID(ctx context.Context) string {
	id, err := stq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only ServiceType entity in the query, returns an error if not exactly one entity was returned.
func (stq *ServiceTypeQuery) Only(ctx context.Context) (*ServiceType, error) {
	sts, err := stq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(sts) {
	case 1:
		return sts[0], nil
	case 0:
		return nil, &NotFoundError{servicetype.Label}
	default:
		return nil, &NotSingularError{servicetype.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (stq *ServiceTypeQuery) OnlyX(ctx context.Context) *ServiceType {
	st, err := stq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return st
}

// OnlyID returns the only ServiceType id in the query, returns an error if not exactly one id was returned.
func (stq *ServiceTypeQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = stq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{servicetype.Label}
	default:
		err = &NotSingularError{servicetype.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (stq *ServiceTypeQuery) OnlyXID(ctx context.Context) string {
	id, err := stq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ServiceTypes.
func (stq *ServiceTypeQuery) All(ctx context.Context) ([]*ServiceType, error) {
	return stq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (stq *ServiceTypeQuery) AllX(ctx context.Context) []*ServiceType {
	sts, err := stq.All(ctx)
	if err != nil {
		panic(err)
	}
	return sts
}

// IDs executes the query and returns a list of ServiceType ids.
func (stq *ServiceTypeQuery) IDs(ctx context.Context) ([]string, error) {
	var ids []string
	if err := stq.Select(servicetype.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (stq *ServiceTypeQuery) IDsX(ctx context.Context) []string {
	ids, err := stq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (stq *ServiceTypeQuery) Count(ctx context.Context) (int, error) {
	return stq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (stq *ServiceTypeQuery) CountX(ctx context.Context) int {
	count, err := stq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (stq *ServiceTypeQuery) Exist(ctx context.Context) (bool, error) {
	return stq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (stq *ServiceTypeQuery) ExistX(ctx context.Context) bool {
	exist, err := stq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (stq *ServiceTypeQuery) Clone() *ServiceTypeQuery {
	return &ServiceTypeQuery{
		config:     stq.config,
		limit:      stq.limit,
		offset:     stq.offset,
		order:      append([]Order{}, stq.order...),
		unique:     append([]string{}, stq.unique...),
		predicates: append([]predicate.ServiceType{}, stq.predicates...),
		// clone intermediate query.
		sql: stq.sql.Clone(),
	}
}

//  WithServices tells the query-builder to eager-loads the nodes that are connected to
// the "services" edge. The optional arguments used to configure the query builder of the edge.
func (stq *ServiceTypeQuery) WithServices(opts ...func(*ServiceQuery)) *ServiceTypeQuery {
	query := &ServiceQuery{config: stq.config}
	for _, opt := range opts {
		opt(query)
	}
	stq.withServices = query
	return stq
}

//  WithPropertyTypes tells the query-builder to eager-loads the nodes that are connected to
// the "property_types" edge. The optional arguments used to configure the query builder of the edge.
func (stq *ServiceTypeQuery) WithPropertyTypes(opts ...func(*PropertyTypeQuery)) *ServiceTypeQuery {
	query := &PropertyTypeQuery{config: stq.config}
	for _, opt := range opts {
		opt(query)
	}
	stq.withPropertyTypes = query
	return stq
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
//	client.ServiceType.Query().
//		GroupBy(servicetype.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (stq *ServiceTypeQuery) GroupBy(field string, fields ...string) *ServiceTypeGroupBy {
	group := &ServiceTypeGroupBy{config: stq.config}
	group.fields = append([]string{field}, fields...)
	group.sql = stq.sqlQuery()
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
//	client.ServiceType.Query().
//		Select(servicetype.FieldCreateTime).
//		Scan(ctx, &v)
//
func (stq *ServiceTypeQuery) Select(field string, fields ...string) *ServiceTypeSelect {
	selector := &ServiceTypeSelect{config: stq.config}
	selector.fields = append([]string{field}, fields...)
	selector.sql = stq.sqlQuery()
	return selector
}

func (stq *ServiceTypeQuery) sqlAll(ctx context.Context) ([]*ServiceType, error) {
	var (
		nodes       = []*ServiceType{}
		_spec       = stq.querySpec()
		loadedTypes = [2]bool{
			stq.withServices != nil,
			stq.withPropertyTypes != nil,
		}
	)
	_spec.ScanValues = func() []interface{} {
		node := &ServiceType{config: stq.config}
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
	if err := sqlgraph.QueryNodes(ctx, stq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := stq.withServices; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[string]*ServiceType)
		for i := range nodes {
			id, err := strconv.Atoi(nodes[i].ID)
			if err != nil {
				return nil, err
			}
			fks = append(fks, id)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.Service(func(s *sql.Selector) {
			s.Where(sql.InValues(servicetype.ServicesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.type_id
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "type_id" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "type_id" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Services = append(node.Edges.Services, n)
		}
	}

	if query := stq.withPropertyTypes; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[string]*ServiceType)
		for i := range nodes {
			id, err := strconv.Atoi(nodes[i].ID)
			if err != nil {
				return nil, err
			}
			fks = append(fks, id)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.PropertyType(func(s *sql.Selector) {
			s.Where(sql.InValues(servicetype.PropertyTypesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.service_type_id
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "service_type_id" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "service_type_id" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.PropertyTypes = append(node.Edges.PropertyTypes, n)
		}
	}

	return nodes, nil
}

func (stq *ServiceTypeQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := stq.querySpec()
	return sqlgraph.CountNodes(ctx, stq.driver, _spec)
}

func (stq *ServiceTypeQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := stq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (stq *ServiceTypeQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   servicetype.Table,
			Columns: servicetype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: servicetype.FieldID,
			},
		},
		From:   stq.sql,
		Unique: true,
	}
	if ps := stq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := stq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := stq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := stq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (stq *ServiceTypeQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(stq.driver.Dialect())
	t1 := builder.Table(servicetype.Table)
	selector := builder.Select(t1.Columns(servicetype.Columns...)...).From(t1)
	if stq.sql != nil {
		selector = stq.sql
		selector.Select(selector.Columns(servicetype.Columns...)...)
	}
	for _, p := range stq.predicates {
		p(selector)
	}
	for _, p := range stq.order {
		p(selector)
	}
	if offset := stq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := stq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ServiceTypeGroupBy is the builder for group-by ServiceType entities.
type ServiceTypeGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query.
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the group-by query.
func (stgb *ServiceTypeGroupBy) Aggregate(fns ...Aggregate) *ServiceTypeGroupBy {
	stgb.fns = append(stgb.fns, fns...)
	return stgb
}

// Scan applies the group-by query and scan the result into the given value.
func (stgb *ServiceTypeGroupBy) Scan(ctx context.Context, v interface{}) error {
	return stgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (stgb *ServiceTypeGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := stgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (stgb *ServiceTypeGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(stgb.fields) > 1 {
		return nil, errors.New("ent: ServiceTypeGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := stgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (stgb *ServiceTypeGroupBy) StringsX(ctx context.Context) []string {
	v, err := stgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (stgb *ServiceTypeGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(stgb.fields) > 1 {
		return nil, errors.New("ent: ServiceTypeGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := stgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (stgb *ServiceTypeGroupBy) IntsX(ctx context.Context) []int {
	v, err := stgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (stgb *ServiceTypeGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(stgb.fields) > 1 {
		return nil, errors.New("ent: ServiceTypeGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := stgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (stgb *ServiceTypeGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := stgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (stgb *ServiceTypeGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(stgb.fields) > 1 {
		return nil, errors.New("ent: ServiceTypeGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := stgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (stgb *ServiceTypeGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := stgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (stgb *ServiceTypeGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := stgb.sqlQuery().Query()
	if err := stgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (stgb *ServiceTypeGroupBy) sqlQuery() *sql.Selector {
	selector := stgb.sql
	columns := make([]string, 0, len(stgb.fields)+len(stgb.fns))
	columns = append(columns, stgb.fields...)
	for _, fn := range stgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(stgb.fields...)
}

// ServiceTypeSelect is the builder for select fields of ServiceType entities.
type ServiceTypeSelect struct {
	config
	fields []string
	// intermediate queries.
	sql *sql.Selector
}

// Scan applies the selector query and scan the result into the given value.
func (sts *ServiceTypeSelect) Scan(ctx context.Context, v interface{}) error {
	return sts.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (sts *ServiceTypeSelect) ScanX(ctx context.Context, v interface{}) {
	if err := sts.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (sts *ServiceTypeSelect) Strings(ctx context.Context) ([]string, error) {
	if len(sts.fields) > 1 {
		return nil, errors.New("ent: ServiceTypeSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := sts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (sts *ServiceTypeSelect) StringsX(ctx context.Context) []string {
	v, err := sts.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (sts *ServiceTypeSelect) Ints(ctx context.Context) ([]int, error) {
	if len(sts.fields) > 1 {
		return nil, errors.New("ent: ServiceTypeSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := sts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (sts *ServiceTypeSelect) IntsX(ctx context.Context) []int {
	v, err := sts.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (sts *ServiceTypeSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(sts.fields) > 1 {
		return nil, errors.New("ent: ServiceTypeSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := sts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (sts *ServiceTypeSelect) Float64sX(ctx context.Context) []float64 {
	v, err := sts.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (sts *ServiceTypeSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(sts.fields) > 1 {
		return nil, errors.New("ent: ServiceTypeSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := sts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (sts *ServiceTypeSelect) BoolsX(ctx context.Context) []bool {
	v, err := sts.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (sts *ServiceTypeSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := sts.sqlQuery().Query()
	if err := sts.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (sts *ServiceTypeSelect) sqlQuery() sql.Querier {
	selector := sts.sql
	selector.Select(selector.Columns(sts.fields...)...)
	return selector
}
