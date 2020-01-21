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
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// EquipmentQuery is the builder for querying Equipment entities.
type EquipmentQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.Equipment
	// intermediate query.
	sql *sql.Selector
}

// Where adds a new predicate for the builder.
func (eq *EquipmentQuery) Where(ps ...predicate.Equipment) *EquipmentQuery {
	eq.predicates = append(eq.predicates, ps...)
	return eq
}

// Limit adds a limit step to the query.
func (eq *EquipmentQuery) Limit(limit int) *EquipmentQuery {
	eq.limit = &limit
	return eq
}

// Offset adds an offset step to the query.
func (eq *EquipmentQuery) Offset(offset int) *EquipmentQuery {
	eq.offset = &offset
	return eq
}

// Order adds an order step to the query.
func (eq *EquipmentQuery) Order(o ...Order) *EquipmentQuery {
	eq.order = append(eq.order, o...)
	return eq
}

// QueryType chains the current query on the type edge.
func (eq *EquipmentQuery) QueryType() *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: eq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
		sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, equipment.TypeTable, equipment.TypeColumn),
	)
	query.sql = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
	return query
}

// QueryLocation chains the current query on the location edge.
func (eq *EquipmentQuery) QueryLocation() *LocationQuery {
	query := &LocationQuery{config: eq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, equipment.LocationTable, equipment.LocationColumn),
	)
	query.sql = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
	return query
}

// QueryParentPosition chains the current query on the parent_position edge.
func (eq *EquipmentQuery) QueryParentPosition() *EquipmentPositionQuery {
	query := &EquipmentPositionQuery{config: eq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
		sqlgraph.To(equipmentposition.Table, equipmentposition.FieldID),
		sqlgraph.Edge(sqlgraph.O2O, true, equipment.ParentPositionTable, equipment.ParentPositionColumn),
	)
	query.sql = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
	return query
}

// QueryPositions chains the current query on the positions edge.
func (eq *EquipmentQuery) QueryPositions() *EquipmentPositionQuery {
	query := &EquipmentPositionQuery{config: eq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
		sqlgraph.To(equipmentposition.Table, equipmentposition.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, equipment.PositionsTable, equipment.PositionsColumn),
	)
	query.sql = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
	return query
}

// QueryPorts chains the current query on the ports edge.
func (eq *EquipmentQuery) QueryPorts() *EquipmentPortQuery {
	query := &EquipmentPortQuery{config: eq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
		sqlgraph.To(equipmentport.Table, equipmentport.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, equipment.PortsTable, equipment.PortsColumn),
	)
	query.sql = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
	return query
}

// QueryWorkOrder chains the current query on the work_order edge.
func (eq *EquipmentQuery) QueryWorkOrder() *WorkOrderQuery {
	query := &WorkOrderQuery{config: eq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
		sqlgraph.To(workorder.Table, workorder.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, equipment.WorkOrderTable, equipment.WorkOrderColumn),
	)
	query.sql = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
	return query
}

// QueryProperties chains the current query on the properties edge.
func (eq *EquipmentQuery) QueryProperties() *PropertyQuery {
	query := &PropertyQuery{config: eq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
		sqlgraph.To(property.Table, property.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, equipment.PropertiesTable, equipment.PropertiesColumn),
	)
	query.sql = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
	return query
}

// QueryFiles chains the current query on the files edge.
func (eq *EquipmentQuery) QueryFiles() *FileQuery {
	query := &FileQuery{config: eq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
		sqlgraph.To(file.Table, file.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, equipment.FilesTable, equipment.FilesColumn),
	)
	query.sql = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
	return query
}

// QueryHyperlinks chains the current query on the hyperlinks edge.
func (eq *EquipmentQuery) QueryHyperlinks() *HyperlinkQuery {
	query := &HyperlinkQuery{config: eq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
		sqlgraph.To(hyperlink.Table, hyperlink.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, equipment.HyperlinksTable, equipment.HyperlinksColumn),
	)
	query.sql = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
	return query
}

// First returns the first Equipment entity in the query. Returns *ErrNotFound when no equipment was found.
func (eq *EquipmentQuery) First(ctx context.Context) (*Equipment, error) {
	es, err := eq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(es) == 0 {
		return nil, &ErrNotFound{equipment.Label}
	}
	return es[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (eq *EquipmentQuery) FirstX(ctx context.Context) *Equipment {
	e, err := eq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return e
}

// FirstID returns the first Equipment id in the query. Returns *ErrNotFound when no id was found.
func (eq *EquipmentQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = eq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &ErrNotFound{equipment.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (eq *EquipmentQuery) FirstXID(ctx context.Context) string {
	id, err := eq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only Equipment entity in the query, returns an error if not exactly one entity was returned.
func (eq *EquipmentQuery) Only(ctx context.Context) (*Equipment, error) {
	es, err := eq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(es) {
	case 1:
		return es[0], nil
	case 0:
		return nil, &ErrNotFound{equipment.Label}
	default:
		return nil, &ErrNotSingular{equipment.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (eq *EquipmentQuery) OnlyX(ctx context.Context) *Equipment {
	e, err := eq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return e
}

// OnlyID returns the only Equipment id in the query, returns an error if not exactly one id was returned.
func (eq *EquipmentQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = eq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &ErrNotFound{equipment.Label}
	default:
		err = &ErrNotSingular{equipment.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (eq *EquipmentQuery) OnlyXID(ctx context.Context) string {
	id, err := eq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of EquipmentSlice.
func (eq *EquipmentQuery) All(ctx context.Context) ([]*Equipment, error) {
	return eq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (eq *EquipmentQuery) AllX(ctx context.Context) []*Equipment {
	es, err := eq.All(ctx)
	if err != nil {
		panic(err)
	}
	return es
}

// IDs executes the query and returns a list of Equipment ids.
func (eq *EquipmentQuery) IDs(ctx context.Context) ([]string, error) {
	var ids []string
	if err := eq.Select(equipment.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (eq *EquipmentQuery) IDsX(ctx context.Context) []string {
	ids, err := eq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (eq *EquipmentQuery) Count(ctx context.Context) (int, error) {
	return eq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (eq *EquipmentQuery) CountX(ctx context.Context) int {
	count, err := eq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (eq *EquipmentQuery) Exist(ctx context.Context) (bool, error) {
	return eq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (eq *EquipmentQuery) ExistX(ctx context.Context) bool {
	exist, err := eq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (eq *EquipmentQuery) Clone() *EquipmentQuery {
	return &EquipmentQuery{
		config:     eq.config,
		limit:      eq.limit,
		offset:     eq.offset,
		order:      append([]Order{}, eq.order...),
		unique:     append([]string{}, eq.unique...),
		predicates: append([]predicate.Equipment{}, eq.predicates...),
		// clone intermediate query.
		sql: eq.sql.Clone(),
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
//	client.Equipment.Query().
//		GroupBy(equipment.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (eq *EquipmentQuery) GroupBy(field string, fields ...string) *EquipmentGroupBy {
	group := &EquipmentGroupBy{config: eq.config}
	group.fields = append([]string{field}, fields...)
	group.sql = eq.sqlQuery()
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
//	client.Equipment.Query().
//		Select(equipment.FieldCreateTime).
//		Scan(ctx, &v)
//
func (eq *EquipmentQuery) Select(field string, fields ...string) *EquipmentSelect {
	selector := &EquipmentSelect{config: eq.config}
	selector.fields = append([]string{field}, fields...)
	selector.sql = eq.sqlQuery()
	return selector
}

func (eq *EquipmentQuery) sqlAll(ctx context.Context) ([]*Equipment, error) {
	var (
		nodes []*Equipment
		spec  = eq.querySpec()
	)
	spec.ScanValues = func() []interface{} {
		node := &Equipment{config: eq.config}
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
	if err := sqlgraph.QueryNodes(ctx, eq.driver, spec); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (eq *EquipmentQuery) sqlCount(ctx context.Context) (int, error) {
	spec := eq.querySpec()
	return sqlgraph.CountNodes(ctx, eq.driver, spec)
}

func (eq *EquipmentQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := eq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (eq *EquipmentQuery) querySpec() *sqlgraph.QuerySpec {
	spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipment.Table,
			Columns: equipment.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: equipment.FieldID,
			},
		},
		From:   eq.sql,
		Unique: true,
	}
	if ps := eq.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := eq.limit; limit != nil {
		spec.Limit = *limit
	}
	if offset := eq.offset; offset != nil {
		spec.Offset = *offset
	}
	if ps := eq.order; len(ps) > 0 {
		spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return spec
}

func (eq *EquipmentQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(eq.driver.Dialect())
	t1 := builder.Table(equipment.Table)
	selector := builder.Select(t1.Columns(equipment.Columns...)...).From(t1)
	if eq.sql != nil {
		selector = eq.sql
		selector.Select(selector.Columns(equipment.Columns...)...)
	}
	for _, p := range eq.predicates {
		p(selector)
	}
	for _, p := range eq.order {
		p(selector)
	}
	if offset := eq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := eq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// EquipmentGroupBy is the builder for group-by Equipment entities.
type EquipmentGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query.
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the group-by query.
func (egb *EquipmentGroupBy) Aggregate(fns ...Aggregate) *EquipmentGroupBy {
	egb.fns = append(egb.fns, fns...)
	return egb
}

// Scan applies the group-by query and scan the result into the given value.
func (egb *EquipmentGroupBy) Scan(ctx context.Context, v interface{}) error {
	return egb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (egb *EquipmentGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := egb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (egb *EquipmentGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(egb.fields) > 1 {
		return nil, errors.New("ent: EquipmentGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := egb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (egb *EquipmentGroupBy) StringsX(ctx context.Context) []string {
	v, err := egb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (egb *EquipmentGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(egb.fields) > 1 {
		return nil, errors.New("ent: EquipmentGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := egb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (egb *EquipmentGroupBy) IntsX(ctx context.Context) []int {
	v, err := egb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (egb *EquipmentGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(egb.fields) > 1 {
		return nil, errors.New("ent: EquipmentGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := egb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (egb *EquipmentGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := egb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (egb *EquipmentGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(egb.fields) > 1 {
		return nil, errors.New("ent: EquipmentGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := egb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (egb *EquipmentGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := egb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (egb *EquipmentGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := egb.sqlQuery().Query()
	if err := egb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (egb *EquipmentGroupBy) sqlQuery() *sql.Selector {
	selector := egb.sql
	columns := make([]string, 0, len(egb.fields)+len(egb.fns))
	columns = append(columns, egb.fields...)
	for _, fn := range egb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(egb.fields...)
}

// EquipmentSelect is the builder for select fields of Equipment entities.
type EquipmentSelect struct {
	config
	fields []string
	// intermediate queries.
	sql *sql.Selector
}

// Scan applies the selector query and scan the result into the given value.
func (es *EquipmentSelect) Scan(ctx context.Context, v interface{}) error {
	return es.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (es *EquipmentSelect) ScanX(ctx context.Context, v interface{}) {
	if err := es.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (es *EquipmentSelect) Strings(ctx context.Context) ([]string, error) {
	if len(es.fields) > 1 {
		return nil, errors.New("ent: EquipmentSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := es.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (es *EquipmentSelect) StringsX(ctx context.Context) []string {
	v, err := es.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (es *EquipmentSelect) Ints(ctx context.Context) ([]int, error) {
	if len(es.fields) > 1 {
		return nil, errors.New("ent: EquipmentSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := es.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (es *EquipmentSelect) IntsX(ctx context.Context) []int {
	v, err := es.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (es *EquipmentSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(es.fields) > 1 {
		return nil, errors.New("ent: EquipmentSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := es.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (es *EquipmentSelect) Float64sX(ctx context.Context) []float64 {
	v, err := es.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (es *EquipmentSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(es.fields) > 1 {
		return nil, errors.New("ent: EquipmentSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := es.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (es *EquipmentSelect) BoolsX(ctx context.Context) []bool {
	v, err := es.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (es *EquipmentSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := es.sqlQuery().Query()
	if err := es.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (es *EquipmentSelect) sqlQuery() sql.Querier {
	selector := es.sql
	selector.Select(selector.Columns(es.fields...)...)
	return selector
}
