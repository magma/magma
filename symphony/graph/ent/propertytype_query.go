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
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// PropertyTypeQuery is the builder for querying PropertyType entities.
type PropertyTypeQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.PropertyType
	// intermediate query.
	sql *sql.Selector
}

// Where adds a new predicate for the builder.
func (ptq *PropertyTypeQuery) Where(ps ...predicate.PropertyType) *PropertyTypeQuery {
	ptq.predicates = append(ptq.predicates, ps...)
	return ptq
}

// Limit adds a limit step to the query.
func (ptq *PropertyTypeQuery) Limit(limit int) *PropertyTypeQuery {
	ptq.limit = &limit
	return ptq
}

// Offset adds an offset step to the query.
func (ptq *PropertyTypeQuery) Offset(offset int) *PropertyTypeQuery {
	ptq.offset = &offset
	return ptq
}

// Order adds an order step to the query.
func (ptq *PropertyTypeQuery) Order(o ...Order) *PropertyTypeQuery {
	ptq.order = append(ptq.order, o...)
	return ptq
}

// QueryProperties chains the current query on the properties edge.
func (ptq *PropertyTypeQuery) QueryProperties() *PropertyQuery {
	query := &PropertyQuery{config: ptq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(propertytype.Table, propertytype.FieldID, ptq.sqlQuery()),
		sqlgraph.To(property.Table, property.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, propertytype.PropertiesTable, propertytype.PropertiesColumn),
	)
	query.sql = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
	return query
}

// QueryLocationType chains the current query on the location_type edge.
func (ptq *PropertyTypeQuery) QueryLocationType() *LocationTypeQuery {
	query := &LocationTypeQuery{config: ptq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(propertytype.Table, propertytype.FieldID, ptq.sqlQuery()),
		sqlgraph.To(locationtype.Table, locationtype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, propertytype.LocationTypeTable, propertytype.LocationTypeColumn),
	)
	query.sql = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
	return query
}

// QueryEquipmentPortType chains the current query on the equipment_port_type edge.
func (ptq *PropertyTypeQuery) QueryEquipmentPortType() *EquipmentPortTypeQuery {
	query := &EquipmentPortTypeQuery{config: ptq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(propertytype.Table, propertytype.FieldID, ptq.sqlQuery()),
		sqlgraph.To(equipmentporttype.Table, equipmentporttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, propertytype.EquipmentPortTypeTable, propertytype.EquipmentPortTypeColumn),
	)
	query.sql = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
	return query
}

// QueryLinkEquipmentPortType chains the current query on the link_equipment_port_type edge.
func (ptq *PropertyTypeQuery) QueryLinkEquipmentPortType() *EquipmentPortTypeQuery {
	query := &EquipmentPortTypeQuery{config: ptq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(propertytype.Table, propertytype.FieldID, ptq.sqlQuery()),
		sqlgraph.To(equipmentporttype.Table, equipmentporttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, propertytype.LinkEquipmentPortTypeTable, propertytype.LinkEquipmentPortTypeColumn),
	)
	query.sql = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
	return query
}

// QueryEquipmentType chains the current query on the equipment_type edge.
func (ptq *PropertyTypeQuery) QueryEquipmentType() *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: ptq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(propertytype.Table, propertytype.FieldID, ptq.sqlQuery()),
		sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, propertytype.EquipmentTypeTable, propertytype.EquipmentTypeColumn),
	)
	query.sql = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
	return query
}

// QueryServiceType chains the current query on the service_type edge.
func (ptq *PropertyTypeQuery) QueryServiceType() *ServiceTypeQuery {
	query := &ServiceTypeQuery{config: ptq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(propertytype.Table, propertytype.FieldID, ptq.sqlQuery()),
		sqlgraph.To(servicetype.Table, servicetype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, propertytype.ServiceTypeTable, propertytype.ServiceTypeColumn),
	)
	query.sql = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
	return query
}

// QueryWorkOrderType chains the current query on the work_order_type edge.
func (ptq *PropertyTypeQuery) QueryWorkOrderType() *WorkOrderTypeQuery {
	query := &WorkOrderTypeQuery{config: ptq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(propertytype.Table, propertytype.FieldID, ptq.sqlQuery()),
		sqlgraph.To(workordertype.Table, workordertype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, propertytype.WorkOrderTypeTable, propertytype.WorkOrderTypeColumn),
	)
	query.sql = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
	return query
}

// QueryProjectType chains the current query on the project_type edge.
func (ptq *PropertyTypeQuery) QueryProjectType() *ProjectTypeQuery {
	query := &ProjectTypeQuery{config: ptq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(propertytype.Table, propertytype.FieldID, ptq.sqlQuery()),
		sqlgraph.To(projecttype.Table, projecttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, propertytype.ProjectTypeTable, propertytype.ProjectTypeColumn),
	)
	query.sql = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
	return query
}

// First returns the first PropertyType entity in the query. Returns *ErrNotFound when no propertytype was found.
func (ptq *PropertyTypeQuery) First(ctx context.Context) (*PropertyType, error) {
	pts, err := ptq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(pts) == 0 {
		return nil, &ErrNotFound{propertytype.Label}
	}
	return pts[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (ptq *PropertyTypeQuery) FirstX(ctx context.Context) *PropertyType {
	pt, err := ptq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return pt
}

// FirstID returns the first PropertyType id in the query. Returns *ErrNotFound when no id was found.
func (ptq *PropertyTypeQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = ptq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &ErrNotFound{propertytype.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (ptq *PropertyTypeQuery) FirstXID(ctx context.Context) string {
	id, err := ptq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only PropertyType entity in the query, returns an error if not exactly one entity was returned.
func (ptq *PropertyTypeQuery) Only(ctx context.Context) (*PropertyType, error) {
	pts, err := ptq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(pts) {
	case 1:
		return pts[0], nil
	case 0:
		return nil, &ErrNotFound{propertytype.Label}
	default:
		return nil, &ErrNotSingular{propertytype.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (ptq *PropertyTypeQuery) OnlyX(ctx context.Context) *PropertyType {
	pt, err := ptq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return pt
}

// OnlyID returns the only PropertyType id in the query, returns an error if not exactly one id was returned.
func (ptq *PropertyTypeQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = ptq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &ErrNotFound{propertytype.Label}
	default:
		err = &ErrNotSingular{propertytype.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (ptq *PropertyTypeQuery) OnlyXID(ctx context.Context) string {
	id, err := ptq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of PropertyTypes.
func (ptq *PropertyTypeQuery) All(ctx context.Context) ([]*PropertyType, error) {
	return ptq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (ptq *PropertyTypeQuery) AllX(ctx context.Context) []*PropertyType {
	pts, err := ptq.All(ctx)
	if err != nil {
		panic(err)
	}
	return pts
}

// IDs executes the query and returns a list of PropertyType ids.
func (ptq *PropertyTypeQuery) IDs(ctx context.Context) ([]string, error) {
	var ids []string
	if err := ptq.Select(propertytype.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (ptq *PropertyTypeQuery) IDsX(ctx context.Context) []string {
	ids, err := ptq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (ptq *PropertyTypeQuery) Count(ctx context.Context) (int, error) {
	return ptq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (ptq *PropertyTypeQuery) CountX(ctx context.Context) int {
	count, err := ptq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (ptq *PropertyTypeQuery) Exist(ctx context.Context) (bool, error) {
	return ptq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (ptq *PropertyTypeQuery) ExistX(ctx context.Context) bool {
	exist, err := ptq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (ptq *PropertyTypeQuery) Clone() *PropertyTypeQuery {
	return &PropertyTypeQuery{
		config:     ptq.config,
		limit:      ptq.limit,
		offset:     ptq.offset,
		order:      append([]Order{}, ptq.order...),
		unique:     append([]string{}, ptq.unique...),
		predicates: append([]predicate.PropertyType{}, ptq.predicates...),
		// clone intermediate query.
		sql: ptq.sql.Clone(),
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
//	client.PropertyType.Query().
//		GroupBy(propertytype.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (ptq *PropertyTypeQuery) GroupBy(field string, fields ...string) *PropertyTypeGroupBy {
	group := &PropertyTypeGroupBy{config: ptq.config}
	group.fields = append([]string{field}, fields...)
	group.sql = ptq.sqlQuery()
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
//	client.PropertyType.Query().
//		Select(propertytype.FieldCreateTime).
//		Scan(ctx, &v)
//
func (ptq *PropertyTypeQuery) Select(field string, fields ...string) *PropertyTypeSelect {
	selector := &PropertyTypeSelect{config: ptq.config}
	selector.fields = append([]string{field}, fields...)
	selector.sql = ptq.sqlQuery()
	return selector
}

func (ptq *PropertyTypeQuery) sqlAll(ctx context.Context) ([]*PropertyType, error) {
	var (
		nodes []*PropertyType
		spec  = ptq.querySpec()
	)
	spec.ScanValues = func() []interface{} {
		node := &PropertyType{config: ptq.config}
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
	if err := sqlgraph.QueryNodes(ctx, ptq.driver, spec); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (ptq *PropertyTypeQuery) sqlCount(ctx context.Context) (int, error) {
	spec := ptq.querySpec()
	return sqlgraph.CountNodes(ctx, ptq.driver, spec)
}

func (ptq *PropertyTypeQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := ptq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (ptq *PropertyTypeQuery) querySpec() *sqlgraph.QuerySpec {
	spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   propertytype.Table,
			Columns: propertytype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: propertytype.FieldID,
			},
		},
		From:   ptq.sql,
		Unique: true,
	}
	if ps := ptq.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := ptq.limit; limit != nil {
		spec.Limit = *limit
	}
	if offset := ptq.offset; offset != nil {
		spec.Offset = *offset
	}
	if ps := ptq.order; len(ps) > 0 {
		spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return spec
}

func (ptq *PropertyTypeQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(ptq.driver.Dialect())
	t1 := builder.Table(propertytype.Table)
	selector := builder.Select(t1.Columns(propertytype.Columns...)...).From(t1)
	if ptq.sql != nil {
		selector = ptq.sql
		selector.Select(selector.Columns(propertytype.Columns...)...)
	}
	for _, p := range ptq.predicates {
		p(selector)
	}
	for _, p := range ptq.order {
		p(selector)
	}
	if offset := ptq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := ptq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// PropertyTypeGroupBy is the builder for group-by PropertyType entities.
type PropertyTypeGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query.
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ptgb *PropertyTypeGroupBy) Aggregate(fns ...Aggregate) *PropertyTypeGroupBy {
	ptgb.fns = append(ptgb.fns, fns...)
	return ptgb
}

// Scan applies the group-by query and scan the result into the given value.
func (ptgb *PropertyTypeGroupBy) Scan(ctx context.Context, v interface{}) error {
	return ptgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ptgb *PropertyTypeGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := ptgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (ptgb *PropertyTypeGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(ptgb.fields) > 1 {
		return nil, errors.New("ent: PropertyTypeGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := ptgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ptgb *PropertyTypeGroupBy) StringsX(ctx context.Context) []string {
	v, err := ptgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (ptgb *PropertyTypeGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(ptgb.fields) > 1 {
		return nil, errors.New("ent: PropertyTypeGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := ptgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ptgb *PropertyTypeGroupBy) IntsX(ctx context.Context) []int {
	v, err := ptgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (ptgb *PropertyTypeGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(ptgb.fields) > 1 {
		return nil, errors.New("ent: PropertyTypeGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := ptgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ptgb *PropertyTypeGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := ptgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (ptgb *PropertyTypeGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(ptgb.fields) > 1 {
		return nil, errors.New("ent: PropertyTypeGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := ptgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ptgb *PropertyTypeGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := ptgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ptgb *PropertyTypeGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ptgb.sqlQuery().Query()
	if err := ptgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ptgb *PropertyTypeGroupBy) sqlQuery() *sql.Selector {
	selector := ptgb.sql
	columns := make([]string, 0, len(ptgb.fields)+len(ptgb.fns))
	columns = append(columns, ptgb.fields...)
	for _, fn := range ptgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(ptgb.fields...)
}

// PropertyTypeSelect is the builder for select fields of PropertyType entities.
type PropertyTypeSelect struct {
	config
	fields []string
	// intermediate queries.
	sql *sql.Selector
}

// Scan applies the selector query and scan the result into the given value.
func (pts *PropertyTypeSelect) Scan(ctx context.Context, v interface{}) error {
	return pts.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (pts *PropertyTypeSelect) ScanX(ctx context.Context, v interface{}) {
	if err := pts.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (pts *PropertyTypeSelect) Strings(ctx context.Context) ([]string, error) {
	if len(pts.fields) > 1 {
		return nil, errors.New("ent: PropertyTypeSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := pts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (pts *PropertyTypeSelect) StringsX(ctx context.Context) []string {
	v, err := pts.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (pts *PropertyTypeSelect) Ints(ctx context.Context) ([]int, error) {
	if len(pts.fields) > 1 {
		return nil, errors.New("ent: PropertyTypeSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := pts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (pts *PropertyTypeSelect) IntsX(ctx context.Context) []int {
	v, err := pts.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (pts *PropertyTypeSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(pts.fields) > 1 {
		return nil, errors.New("ent: PropertyTypeSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := pts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (pts *PropertyTypeSelect) Float64sX(ctx context.Context) []float64 {
	v, err := pts.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (pts *PropertyTypeSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(pts.fields) > 1 {
		return nil, errors.New("ent: PropertyTypeSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := pts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (pts *PropertyTypeSelect) BoolsX(ctx context.Context) []bool {
	v, err := pts.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (pts *PropertyTypeSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := pts.sqlQuery().Query()
	if err := pts.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (pts *PropertyTypeSelect) sqlQuery() sql.Querier {
	selector := pts.sql
	selector.Select(selector.Columns(pts.fields...)...)
	return selector
}
