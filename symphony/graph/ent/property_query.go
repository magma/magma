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
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// PropertyQuery is the builder for querying Property entities.
type PropertyQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.Property
	// intermediate query.
	sql *sql.Selector
}

// Where adds a new predicate for the builder.
func (pq *PropertyQuery) Where(ps ...predicate.Property) *PropertyQuery {
	pq.predicates = append(pq.predicates, ps...)
	return pq
}

// Limit adds a limit step to the query.
func (pq *PropertyQuery) Limit(limit int) *PropertyQuery {
	pq.limit = &limit
	return pq
}

// Offset adds an offset step to the query.
func (pq *PropertyQuery) Offset(offset int) *PropertyQuery {
	pq.offset = &offset
	return pq
}

// Order adds an order step to the query.
func (pq *PropertyQuery) Order(o ...Order) *PropertyQuery {
	pq.order = append(pq.order, o...)
	return pq
}

// QueryType chains the current query on the type edge.
func (pq *PropertyQuery) QueryType() *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: pq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, pq.sqlQuery()),
		sqlgraph.To(propertytype.Table, propertytype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, property.TypeTable, property.TypeColumn),
	)
	query.sql = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
	return query
}

// QueryLocation chains the current query on the location edge.
func (pq *PropertyQuery) QueryLocation() *LocationQuery {
	query := &LocationQuery{config: pq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, pq.sqlQuery()),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, property.LocationTable, property.LocationColumn),
	)
	query.sql = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
	return query
}

// QueryEquipment chains the current query on the equipment edge.
func (pq *PropertyQuery) QueryEquipment() *EquipmentQuery {
	query := &EquipmentQuery{config: pq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, pq.sqlQuery()),
		sqlgraph.To(equipment.Table, equipment.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, property.EquipmentTable, property.EquipmentColumn),
	)
	query.sql = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
	return query
}

// QueryService chains the current query on the service edge.
func (pq *PropertyQuery) QueryService() *ServiceQuery {
	query := &ServiceQuery{config: pq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, pq.sqlQuery()),
		sqlgraph.To(service.Table, service.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, property.ServiceTable, property.ServiceColumn),
	)
	query.sql = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
	return query
}

// QueryEquipmentPort chains the current query on the equipment_port edge.
func (pq *PropertyQuery) QueryEquipmentPort() *EquipmentPortQuery {
	query := &EquipmentPortQuery{config: pq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, pq.sqlQuery()),
		sqlgraph.To(equipmentport.Table, equipmentport.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, property.EquipmentPortTable, property.EquipmentPortColumn),
	)
	query.sql = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
	return query
}

// QueryLink chains the current query on the link edge.
func (pq *PropertyQuery) QueryLink() *LinkQuery {
	query := &LinkQuery{config: pq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, pq.sqlQuery()),
		sqlgraph.To(link.Table, link.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, property.LinkTable, property.LinkColumn),
	)
	query.sql = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
	return query
}

// QueryWorkOrder chains the current query on the work_order edge.
func (pq *PropertyQuery) QueryWorkOrder() *WorkOrderQuery {
	query := &WorkOrderQuery{config: pq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, pq.sqlQuery()),
		sqlgraph.To(workorder.Table, workorder.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, property.WorkOrderTable, property.WorkOrderColumn),
	)
	query.sql = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
	return query
}

// QueryProject chains the current query on the project edge.
func (pq *PropertyQuery) QueryProject() *ProjectQuery {
	query := &ProjectQuery{config: pq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, pq.sqlQuery()),
		sqlgraph.To(project.Table, project.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, property.ProjectTable, property.ProjectColumn),
	)
	query.sql = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
	return query
}

// QueryEquipmentValue chains the current query on the equipment_value edge.
func (pq *PropertyQuery) QueryEquipmentValue() *EquipmentQuery {
	query := &EquipmentQuery{config: pq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, pq.sqlQuery()),
		sqlgraph.To(equipment.Table, equipment.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, property.EquipmentValueTable, property.EquipmentValueColumn),
	)
	query.sql = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
	return query
}

// QueryLocationValue chains the current query on the location_value edge.
func (pq *PropertyQuery) QueryLocationValue() *LocationQuery {
	query := &LocationQuery{config: pq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, pq.sqlQuery()),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, property.LocationValueTable, property.LocationValueColumn),
	)
	query.sql = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
	return query
}

// QueryServiceValue chains the current query on the service_value edge.
func (pq *PropertyQuery) QueryServiceValue() *ServiceQuery {
	query := &ServiceQuery{config: pq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, pq.sqlQuery()),
		sqlgraph.To(service.Table, service.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, property.ServiceValueTable, property.ServiceValueColumn),
	)
	query.sql = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
	return query
}

// First returns the first Property entity in the query. Returns *ErrNotFound when no property was found.
func (pq *PropertyQuery) First(ctx context.Context) (*Property, error) {
	prs, err := pq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(prs) == 0 {
		return nil, &ErrNotFound{property.Label}
	}
	return prs[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (pq *PropertyQuery) FirstX(ctx context.Context) *Property {
	pr, err := pq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return pr
}

// FirstID returns the first Property id in the query. Returns *ErrNotFound when no id was found.
func (pq *PropertyQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = pq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &ErrNotFound{property.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (pq *PropertyQuery) FirstXID(ctx context.Context) string {
	id, err := pq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only Property entity in the query, returns an error if not exactly one entity was returned.
func (pq *PropertyQuery) Only(ctx context.Context) (*Property, error) {
	prs, err := pq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(prs) {
	case 1:
		return prs[0], nil
	case 0:
		return nil, &ErrNotFound{property.Label}
	default:
		return nil, &ErrNotSingular{property.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (pq *PropertyQuery) OnlyX(ctx context.Context) *Property {
	pr, err := pq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return pr
}

// OnlyID returns the only Property id in the query, returns an error if not exactly one id was returned.
func (pq *PropertyQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = pq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &ErrNotFound{property.Label}
	default:
		err = &ErrNotSingular{property.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (pq *PropertyQuery) OnlyXID(ctx context.Context) string {
	id, err := pq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Properties.
func (pq *PropertyQuery) All(ctx context.Context) ([]*Property, error) {
	return pq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (pq *PropertyQuery) AllX(ctx context.Context) []*Property {
	prs, err := pq.All(ctx)
	if err != nil {
		panic(err)
	}
	return prs
}

// IDs executes the query and returns a list of Property ids.
func (pq *PropertyQuery) IDs(ctx context.Context) ([]string, error) {
	var ids []string
	if err := pq.Select(property.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (pq *PropertyQuery) IDsX(ctx context.Context) []string {
	ids, err := pq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (pq *PropertyQuery) Count(ctx context.Context) (int, error) {
	return pq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (pq *PropertyQuery) CountX(ctx context.Context) int {
	count, err := pq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (pq *PropertyQuery) Exist(ctx context.Context) (bool, error) {
	return pq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (pq *PropertyQuery) ExistX(ctx context.Context) bool {
	exist, err := pq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (pq *PropertyQuery) Clone() *PropertyQuery {
	return &PropertyQuery{
		config:     pq.config,
		limit:      pq.limit,
		offset:     pq.offset,
		order:      append([]Order{}, pq.order...),
		unique:     append([]string{}, pq.unique...),
		predicates: append([]predicate.Property{}, pq.predicates...),
		// clone intermediate query.
		sql: pq.sql.Clone(),
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
//	client.Property.Query().
//		GroupBy(property.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (pq *PropertyQuery) GroupBy(field string, fields ...string) *PropertyGroupBy {
	group := &PropertyGroupBy{config: pq.config}
	group.fields = append([]string{field}, fields...)
	group.sql = pq.sqlQuery()
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
//	client.Property.Query().
//		Select(property.FieldCreateTime).
//		Scan(ctx, &v)
//
func (pq *PropertyQuery) Select(field string, fields ...string) *PropertySelect {
	selector := &PropertySelect{config: pq.config}
	selector.fields = append([]string{field}, fields...)
	selector.sql = pq.sqlQuery()
	return selector
}

func (pq *PropertyQuery) sqlAll(ctx context.Context) ([]*Property, error) {
	var (
		nodes []*Property
		spec  = pq.querySpec()
	)
	spec.ScanValues = func() []interface{} {
		node := &Property{config: pq.config}
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
	if err := sqlgraph.QueryNodes(ctx, pq.driver, spec); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (pq *PropertyQuery) sqlCount(ctx context.Context) (int, error) {
	spec := pq.querySpec()
	return sqlgraph.CountNodes(ctx, pq.driver, spec)
}

func (pq *PropertyQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := pq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (pq *PropertyQuery) querySpec() *sqlgraph.QuerySpec {
	spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   property.Table,
			Columns: property.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: property.FieldID,
			},
		},
		From:   pq.sql,
		Unique: true,
	}
	if ps := pq.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := pq.limit; limit != nil {
		spec.Limit = *limit
	}
	if offset := pq.offset; offset != nil {
		spec.Offset = *offset
	}
	if ps := pq.order; len(ps) > 0 {
		spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return spec
}

func (pq *PropertyQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(pq.driver.Dialect())
	t1 := builder.Table(property.Table)
	selector := builder.Select(t1.Columns(property.Columns...)...).From(t1)
	if pq.sql != nil {
		selector = pq.sql
		selector.Select(selector.Columns(property.Columns...)...)
	}
	for _, p := range pq.predicates {
		p(selector)
	}
	for _, p := range pq.order {
		p(selector)
	}
	if offset := pq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := pq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// PropertyGroupBy is the builder for group-by Property entities.
type PropertyGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query.
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the group-by query.
func (pgb *PropertyGroupBy) Aggregate(fns ...Aggregate) *PropertyGroupBy {
	pgb.fns = append(pgb.fns, fns...)
	return pgb
}

// Scan applies the group-by query and scan the result into the given value.
func (pgb *PropertyGroupBy) Scan(ctx context.Context, v interface{}) error {
	return pgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (pgb *PropertyGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := pgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (pgb *PropertyGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(pgb.fields) > 1 {
		return nil, errors.New("ent: PropertyGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := pgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (pgb *PropertyGroupBy) StringsX(ctx context.Context) []string {
	v, err := pgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (pgb *PropertyGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(pgb.fields) > 1 {
		return nil, errors.New("ent: PropertyGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := pgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (pgb *PropertyGroupBy) IntsX(ctx context.Context) []int {
	v, err := pgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (pgb *PropertyGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(pgb.fields) > 1 {
		return nil, errors.New("ent: PropertyGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := pgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (pgb *PropertyGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := pgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (pgb *PropertyGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(pgb.fields) > 1 {
		return nil, errors.New("ent: PropertyGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := pgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (pgb *PropertyGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := pgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (pgb *PropertyGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := pgb.sqlQuery().Query()
	if err := pgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (pgb *PropertyGroupBy) sqlQuery() *sql.Selector {
	selector := pgb.sql
	columns := make([]string, 0, len(pgb.fields)+len(pgb.fns))
	columns = append(columns, pgb.fields...)
	for _, fn := range pgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(pgb.fields...)
}

// PropertySelect is the builder for select fields of Property entities.
type PropertySelect struct {
	config
	fields []string
	// intermediate queries.
	sql *sql.Selector
}

// Scan applies the selector query and scan the result into the given value.
func (ps *PropertySelect) Scan(ctx context.Context, v interface{}) error {
	return ps.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ps *PropertySelect) ScanX(ctx context.Context, v interface{}) {
	if err := ps.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (ps *PropertySelect) Strings(ctx context.Context) ([]string, error) {
	if len(ps.fields) > 1 {
		return nil, errors.New("ent: PropertySelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := ps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ps *PropertySelect) StringsX(ctx context.Context) []string {
	v, err := ps.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (ps *PropertySelect) Ints(ctx context.Context) ([]int, error) {
	if len(ps.fields) > 1 {
		return nil, errors.New("ent: PropertySelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := ps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ps *PropertySelect) IntsX(ctx context.Context) []int {
	v, err := ps.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (ps *PropertySelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(ps.fields) > 1 {
		return nil, errors.New("ent: PropertySelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := ps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ps *PropertySelect) Float64sX(ctx context.Context) []float64 {
	v, err := ps.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (ps *PropertySelect) Bools(ctx context.Context) ([]bool, error) {
	if len(ps.fields) > 1 {
		return nil, errors.New("ent: PropertySelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := ps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ps *PropertySelect) BoolsX(ctx context.Context) []bool {
	v, err := ps.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ps *PropertySelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ps.sqlQuery().Query()
	if err := ps.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ps *PropertySelect) sqlQuery() sql.Querier {
	selector := ps.sql
	selector.Select(selector.Columns(ps.fields...)...)
	return selector
}
