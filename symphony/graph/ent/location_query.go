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
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/floorplan"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// LocationQuery is the builder for querying Location entities.
type LocationQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.Location
	// intermediate query.
	sql *sql.Selector
}

// Where adds a new predicate for the builder.
func (lq *LocationQuery) Where(ps ...predicate.Location) *LocationQuery {
	lq.predicates = append(lq.predicates, ps...)
	return lq
}

// Limit adds a limit step to the query.
func (lq *LocationQuery) Limit(limit int) *LocationQuery {
	lq.limit = &limit
	return lq
}

// Offset adds an offset step to the query.
func (lq *LocationQuery) Offset(offset int) *LocationQuery {
	lq.offset = &offset
	return lq
}

// Order adds an order step to the query.
func (lq *LocationQuery) Order(o ...Order) *LocationQuery {
	lq.order = append(lq.order, o...)
	return lq
}

// QueryType chains the current query on the type edge.
func (lq *LocationQuery) QueryType() *LocationTypeQuery {
	query := &LocationTypeQuery{config: lq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, lq.sqlQuery()),
		sqlgraph.To(locationtype.Table, locationtype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, location.TypeTable, location.TypeColumn),
	)
	query.sql = sqlgraph.SetNeighbors(lq.driver.Dialect(), step)
	return query
}

// QueryParent chains the current query on the parent edge.
func (lq *LocationQuery) QueryParent() *LocationQuery {
	query := &LocationQuery{config: lq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, lq.sqlQuery()),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, location.ParentTable, location.ParentColumn),
	)
	query.sql = sqlgraph.SetNeighbors(lq.driver.Dialect(), step)
	return query
}

// QueryChildren chains the current query on the children edge.
func (lq *LocationQuery) QueryChildren() *LocationQuery {
	query := &LocationQuery{config: lq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, lq.sqlQuery()),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, location.ChildrenTable, location.ChildrenColumn),
	)
	query.sql = sqlgraph.SetNeighbors(lq.driver.Dialect(), step)
	return query
}

// QueryFiles chains the current query on the files edge.
func (lq *LocationQuery) QueryFiles() *FileQuery {
	query := &FileQuery{config: lq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, lq.sqlQuery()),
		sqlgraph.To(file.Table, file.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, location.FilesTable, location.FilesColumn),
	)
	query.sql = sqlgraph.SetNeighbors(lq.driver.Dialect(), step)
	return query
}

// QueryEquipment chains the current query on the equipment edge.
func (lq *LocationQuery) QueryEquipment() *EquipmentQuery {
	query := &EquipmentQuery{config: lq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, lq.sqlQuery()),
		sqlgraph.To(equipment.Table, equipment.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, location.EquipmentTable, location.EquipmentColumn),
	)
	query.sql = sqlgraph.SetNeighbors(lq.driver.Dialect(), step)
	return query
}

// QueryProperties chains the current query on the properties edge.
func (lq *LocationQuery) QueryProperties() *PropertyQuery {
	query := &PropertyQuery{config: lq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, lq.sqlQuery()),
		sqlgraph.To(property.Table, property.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, location.PropertiesTable, location.PropertiesColumn),
	)
	query.sql = sqlgraph.SetNeighbors(lq.driver.Dialect(), step)
	return query
}

// QuerySurvey chains the current query on the survey edge.
func (lq *LocationQuery) QuerySurvey() *SurveyQuery {
	query := &SurveyQuery{config: lq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, lq.sqlQuery()),
		sqlgraph.To(survey.Table, survey.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, location.SurveyTable, location.SurveyColumn),
	)
	query.sql = sqlgraph.SetNeighbors(lq.driver.Dialect(), step)
	return query
}

// QueryWifiScan chains the current query on the wifi_scan edge.
func (lq *LocationQuery) QueryWifiScan() *SurveyWiFiScanQuery {
	query := &SurveyWiFiScanQuery{config: lq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, lq.sqlQuery()),
		sqlgraph.To(surveywifiscan.Table, surveywifiscan.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, location.WifiScanTable, location.WifiScanColumn),
	)
	query.sql = sqlgraph.SetNeighbors(lq.driver.Dialect(), step)
	return query
}

// QueryCellScan chains the current query on the cell_scan edge.
func (lq *LocationQuery) QueryCellScan() *SurveyCellScanQuery {
	query := &SurveyCellScanQuery{config: lq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, lq.sqlQuery()),
		sqlgraph.To(surveycellscan.Table, surveycellscan.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, location.CellScanTable, location.CellScanColumn),
	)
	query.sql = sqlgraph.SetNeighbors(lq.driver.Dialect(), step)
	return query
}

// QueryWorkOrders chains the current query on the work_orders edge.
func (lq *LocationQuery) QueryWorkOrders() *WorkOrderQuery {
	query := &WorkOrderQuery{config: lq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, lq.sqlQuery()),
		sqlgraph.To(workorder.Table, workorder.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, location.WorkOrdersTable, location.WorkOrdersColumn),
	)
	query.sql = sqlgraph.SetNeighbors(lq.driver.Dialect(), step)
	return query
}

// QueryFloorPlans chains the current query on the floor_plans edge.
func (lq *LocationQuery) QueryFloorPlans() *FloorPlanQuery {
	query := &FloorPlanQuery{config: lq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, lq.sqlQuery()),
		sqlgraph.To(floorplan.Table, floorplan.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, location.FloorPlansTable, location.FloorPlansColumn),
	)
	query.sql = sqlgraph.SetNeighbors(lq.driver.Dialect(), step)
	return query
}

// First returns the first Location entity in the query. Returns *ErrNotFound when no location was found.
func (lq *LocationQuery) First(ctx context.Context) (*Location, error) {
	ls, err := lq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(ls) == 0 {
		return nil, &ErrNotFound{location.Label}
	}
	return ls[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (lq *LocationQuery) FirstX(ctx context.Context) *Location {
	l, err := lq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return l
}

// FirstID returns the first Location id in the query. Returns *ErrNotFound when no id was found.
func (lq *LocationQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = lq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &ErrNotFound{location.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (lq *LocationQuery) FirstXID(ctx context.Context) string {
	id, err := lq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only Location entity in the query, returns an error if not exactly one entity was returned.
func (lq *LocationQuery) Only(ctx context.Context) (*Location, error) {
	ls, err := lq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(ls) {
	case 1:
		return ls[0], nil
	case 0:
		return nil, &ErrNotFound{location.Label}
	default:
		return nil, &ErrNotSingular{location.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (lq *LocationQuery) OnlyX(ctx context.Context) *Location {
	l, err := lq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return l
}

// OnlyID returns the only Location id in the query, returns an error if not exactly one id was returned.
func (lq *LocationQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = lq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &ErrNotFound{location.Label}
	default:
		err = &ErrNotSingular{location.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (lq *LocationQuery) OnlyXID(ctx context.Context) string {
	id, err := lq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Locations.
func (lq *LocationQuery) All(ctx context.Context) ([]*Location, error) {
	return lq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (lq *LocationQuery) AllX(ctx context.Context) []*Location {
	ls, err := lq.All(ctx)
	if err != nil {
		panic(err)
	}
	return ls
}

// IDs executes the query and returns a list of Location ids.
func (lq *LocationQuery) IDs(ctx context.Context) ([]string, error) {
	var ids []string
	if err := lq.Select(location.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (lq *LocationQuery) IDsX(ctx context.Context) []string {
	ids, err := lq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (lq *LocationQuery) Count(ctx context.Context) (int, error) {
	return lq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (lq *LocationQuery) CountX(ctx context.Context) int {
	count, err := lq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (lq *LocationQuery) Exist(ctx context.Context) (bool, error) {
	return lq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (lq *LocationQuery) ExistX(ctx context.Context) bool {
	exist, err := lq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (lq *LocationQuery) Clone() *LocationQuery {
	return &LocationQuery{
		config:     lq.config,
		limit:      lq.limit,
		offset:     lq.offset,
		order:      append([]Order{}, lq.order...),
		unique:     append([]string{}, lq.unique...),
		predicates: append([]predicate.Location{}, lq.predicates...),
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
//	client.Location.Query().
//		GroupBy(location.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (lq *LocationQuery) GroupBy(field string, fields ...string) *LocationGroupBy {
	group := &LocationGroupBy{config: lq.config}
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
//	client.Location.Query().
//		Select(location.FieldCreateTime).
//		Scan(ctx, &v)
//
func (lq *LocationQuery) Select(field string, fields ...string) *LocationSelect {
	selector := &LocationSelect{config: lq.config}
	selector.fields = append([]string{field}, fields...)
	selector.sql = lq.sqlQuery()
	return selector
}

func (lq *LocationQuery) sqlAll(ctx context.Context) ([]*Location, error) {
	var (
		nodes []*Location
		spec  = lq.querySpec()
	)
	spec.ScanValues = func() []interface{} {
		node := &Location{config: lq.config}
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

func (lq *LocationQuery) sqlCount(ctx context.Context) (int, error) {
	spec := lq.querySpec()
	return sqlgraph.CountNodes(ctx, lq.driver, spec)
}

func (lq *LocationQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := lq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (lq *LocationQuery) querySpec() *sqlgraph.QuerySpec {
	spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   location.Table,
			Columns: location.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: location.FieldID,
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

func (lq *LocationQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(lq.driver.Dialect())
	t1 := builder.Table(location.Table)
	selector := builder.Select(t1.Columns(location.Columns...)...).From(t1)
	if lq.sql != nil {
		selector = lq.sql
		selector.Select(selector.Columns(location.Columns...)...)
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

// LocationGroupBy is the builder for group-by Location entities.
type LocationGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query.
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the group-by query.
func (lgb *LocationGroupBy) Aggregate(fns ...Aggregate) *LocationGroupBy {
	lgb.fns = append(lgb.fns, fns...)
	return lgb
}

// Scan applies the group-by query and scan the result into the given value.
func (lgb *LocationGroupBy) Scan(ctx context.Context, v interface{}) error {
	return lgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (lgb *LocationGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := lgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (lgb *LocationGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(lgb.fields) > 1 {
		return nil, errors.New("ent: LocationGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := lgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (lgb *LocationGroupBy) StringsX(ctx context.Context) []string {
	v, err := lgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (lgb *LocationGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(lgb.fields) > 1 {
		return nil, errors.New("ent: LocationGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := lgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (lgb *LocationGroupBy) IntsX(ctx context.Context) []int {
	v, err := lgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (lgb *LocationGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(lgb.fields) > 1 {
		return nil, errors.New("ent: LocationGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := lgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (lgb *LocationGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := lgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (lgb *LocationGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(lgb.fields) > 1 {
		return nil, errors.New("ent: LocationGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := lgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (lgb *LocationGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := lgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (lgb *LocationGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := lgb.sqlQuery().Query()
	if err := lgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (lgb *LocationGroupBy) sqlQuery() *sql.Selector {
	selector := lgb.sql
	columns := make([]string, 0, len(lgb.fields)+len(lgb.fns))
	columns = append(columns, lgb.fields...)
	for _, fn := range lgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(lgb.fields...)
}

// LocationSelect is the builder for select fields of Location entities.
type LocationSelect struct {
	config
	fields []string
	// intermediate queries.
	sql *sql.Selector
}

// Scan applies the selector query and scan the result into the given value.
func (ls *LocationSelect) Scan(ctx context.Context, v interface{}) error {
	return ls.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ls *LocationSelect) ScanX(ctx context.Context, v interface{}) {
	if err := ls.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (ls *LocationSelect) Strings(ctx context.Context) ([]string, error) {
	if len(ls.fields) > 1 {
		return nil, errors.New("ent: LocationSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := ls.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ls *LocationSelect) StringsX(ctx context.Context) []string {
	v, err := ls.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (ls *LocationSelect) Ints(ctx context.Context) ([]int, error) {
	if len(ls.fields) > 1 {
		return nil, errors.New("ent: LocationSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := ls.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ls *LocationSelect) IntsX(ctx context.Context) []int {
	v, err := ls.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (ls *LocationSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(ls.fields) > 1 {
		return nil, errors.New("ent: LocationSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := ls.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ls *LocationSelect) Float64sX(ctx context.Context) []float64 {
	v, err := ls.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (ls *LocationSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(ls.fields) > 1 {
		return nil, errors.New("ent: LocationSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := ls.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ls *LocationSelect) BoolsX(ctx context.Context) []bool {
	v, err := ls.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ls *LocationSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ls.sqlQuery().Query()
	if err := ls.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ls *LocationSelect) sqlQuery() sql.Querier {
	selector := ls.sql
	selector.Select(selector.Columns(ls.fields...)...)
	return selector
}
