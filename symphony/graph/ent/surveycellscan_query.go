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
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
)

// SurveyCellScanQuery is the builder for querying SurveyCellScan entities.
type SurveyCellScanQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.SurveyCellScan
	// eager-loading edges.
	withSurveyQuestion *SurveyQuestionQuery
	withLocation       *LocationQuery
	withFKs            bool
	// intermediate query.
	sql *sql.Selector
}

// Where adds a new predicate for the builder.
func (scsq *SurveyCellScanQuery) Where(ps ...predicate.SurveyCellScan) *SurveyCellScanQuery {
	scsq.predicates = append(scsq.predicates, ps...)
	return scsq
}

// Limit adds a limit step to the query.
func (scsq *SurveyCellScanQuery) Limit(limit int) *SurveyCellScanQuery {
	scsq.limit = &limit
	return scsq
}

// Offset adds an offset step to the query.
func (scsq *SurveyCellScanQuery) Offset(offset int) *SurveyCellScanQuery {
	scsq.offset = &offset
	return scsq
}

// Order adds an order step to the query.
func (scsq *SurveyCellScanQuery) Order(o ...Order) *SurveyCellScanQuery {
	scsq.order = append(scsq.order, o...)
	return scsq
}

// QuerySurveyQuestion chains the current query on the survey_question edge.
func (scsq *SurveyCellScanQuery) QuerySurveyQuestion() *SurveyQuestionQuery {
	query := &SurveyQuestionQuery{config: scsq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(surveycellscan.Table, surveycellscan.FieldID, scsq.sqlQuery()),
		sqlgraph.To(surveyquestion.Table, surveyquestion.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, surveycellscan.SurveyQuestionTable, surveycellscan.SurveyQuestionColumn),
	)
	query.sql = sqlgraph.SetNeighbors(scsq.driver.Dialect(), step)
	return query
}

// QueryLocation chains the current query on the location edge.
func (scsq *SurveyCellScanQuery) QueryLocation() *LocationQuery {
	query := &LocationQuery{config: scsq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(surveycellscan.Table, surveycellscan.FieldID, scsq.sqlQuery()),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, surveycellscan.LocationTable, surveycellscan.LocationColumn),
	)
	query.sql = sqlgraph.SetNeighbors(scsq.driver.Dialect(), step)
	return query
}

// First returns the first SurveyCellScan entity in the query. Returns *NotFoundError when no surveycellscan was found.
func (scsq *SurveyCellScanQuery) First(ctx context.Context) (*SurveyCellScan, error) {
	scsSlice, err := scsq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(scsSlice) == 0 {
		return nil, &NotFoundError{surveycellscan.Label}
	}
	return scsSlice[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (scsq *SurveyCellScanQuery) FirstX(ctx context.Context) *SurveyCellScan {
	scs, err := scsq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return scs
}

// FirstID returns the first SurveyCellScan id in the query. Returns *NotFoundError when no id was found.
func (scsq *SurveyCellScanQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = scsq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{surveycellscan.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (scsq *SurveyCellScanQuery) FirstXID(ctx context.Context) int {
	id, err := scsq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only SurveyCellScan entity in the query, returns an error if not exactly one entity was returned.
func (scsq *SurveyCellScanQuery) Only(ctx context.Context) (*SurveyCellScan, error) {
	scsSlice, err := scsq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(scsSlice) {
	case 1:
		return scsSlice[0], nil
	case 0:
		return nil, &NotFoundError{surveycellscan.Label}
	default:
		return nil, &NotSingularError{surveycellscan.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (scsq *SurveyCellScanQuery) OnlyX(ctx context.Context) *SurveyCellScan {
	scs, err := scsq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return scs
}

// OnlyID returns the only SurveyCellScan id in the query, returns an error if not exactly one id was returned.
func (scsq *SurveyCellScanQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = scsq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{surveycellscan.Label}
	default:
		err = &NotSingularError{surveycellscan.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (scsq *SurveyCellScanQuery) OnlyXID(ctx context.Context) int {
	id, err := scsq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of SurveyCellScans.
func (scsq *SurveyCellScanQuery) All(ctx context.Context) ([]*SurveyCellScan, error) {
	return scsq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (scsq *SurveyCellScanQuery) AllX(ctx context.Context) []*SurveyCellScan {
	scsSlice, err := scsq.All(ctx)
	if err != nil {
		panic(err)
	}
	return scsSlice
}

// IDs executes the query and returns a list of SurveyCellScan ids.
func (scsq *SurveyCellScanQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := scsq.Select(surveycellscan.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (scsq *SurveyCellScanQuery) IDsX(ctx context.Context) []int {
	ids, err := scsq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (scsq *SurveyCellScanQuery) Count(ctx context.Context) (int, error) {
	return scsq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (scsq *SurveyCellScanQuery) CountX(ctx context.Context) int {
	count, err := scsq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (scsq *SurveyCellScanQuery) Exist(ctx context.Context) (bool, error) {
	return scsq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (scsq *SurveyCellScanQuery) ExistX(ctx context.Context) bool {
	exist, err := scsq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (scsq *SurveyCellScanQuery) Clone() *SurveyCellScanQuery {
	return &SurveyCellScanQuery{
		config:     scsq.config,
		limit:      scsq.limit,
		offset:     scsq.offset,
		order:      append([]Order{}, scsq.order...),
		unique:     append([]string{}, scsq.unique...),
		predicates: append([]predicate.SurveyCellScan{}, scsq.predicates...),
		// clone intermediate query.
		sql: scsq.sql.Clone(),
	}
}

//  WithSurveyQuestion tells the query-builder to eager-loads the nodes that are connected to
// the "survey_question" edge. The optional arguments used to configure the query builder of the edge.
func (scsq *SurveyCellScanQuery) WithSurveyQuestion(opts ...func(*SurveyQuestionQuery)) *SurveyCellScanQuery {
	query := &SurveyQuestionQuery{config: scsq.config}
	for _, opt := range opts {
		opt(query)
	}
	scsq.withSurveyQuestion = query
	return scsq
}

//  WithLocation tells the query-builder to eager-loads the nodes that are connected to
// the "location" edge. The optional arguments used to configure the query builder of the edge.
func (scsq *SurveyCellScanQuery) WithLocation(opts ...func(*LocationQuery)) *SurveyCellScanQuery {
	query := &LocationQuery{config: scsq.config}
	for _, opt := range opts {
		opt(query)
	}
	scsq.withLocation = query
	return scsq
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
//	client.SurveyCellScan.Query().
//		GroupBy(surveycellscan.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (scsq *SurveyCellScanQuery) GroupBy(field string, fields ...string) *SurveyCellScanGroupBy {
	group := &SurveyCellScanGroupBy{config: scsq.config}
	group.fields = append([]string{field}, fields...)
	group.sql = scsq.sqlQuery()
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
//	client.SurveyCellScan.Query().
//		Select(surveycellscan.FieldCreateTime).
//		Scan(ctx, &v)
//
func (scsq *SurveyCellScanQuery) Select(field string, fields ...string) *SurveyCellScanSelect {
	selector := &SurveyCellScanSelect{config: scsq.config}
	selector.fields = append([]string{field}, fields...)
	selector.sql = scsq.sqlQuery()
	return selector
}

func (scsq *SurveyCellScanQuery) sqlAll(ctx context.Context) ([]*SurveyCellScan, error) {
	var (
		nodes       = []*SurveyCellScan{}
		withFKs     = scsq.withFKs
		_spec       = scsq.querySpec()
		loadedTypes = [2]bool{
			scsq.withSurveyQuestion != nil,
			scsq.withLocation != nil,
		}
	)
	if scsq.withSurveyQuestion != nil || scsq.withLocation != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, surveycellscan.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &SurveyCellScan{config: scsq.config}
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
	if err := sqlgraph.QueryNodes(ctx, scsq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := scsq.withSurveyQuestion; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*SurveyCellScan)
		for i := range nodes {
			if fk := nodes[i].survey_cell_scan_survey_question; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(surveyquestion.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "survey_cell_scan_survey_question" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.SurveyQuestion = n
			}
		}
	}

	if query := scsq.withLocation; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*SurveyCellScan)
		for i := range nodes {
			if fk := nodes[i].survey_cell_scan_location; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(location.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "survey_cell_scan_location" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Location = n
			}
		}
	}

	return nodes, nil
}

func (scsq *SurveyCellScanQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := scsq.querySpec()
	return sqlgraph.CountNodes(ctx, scsq.driver, _spec)
}

func (scsq *SurveyCellScanQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := scsq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (scsq *SurveyCellScanQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveycellscan.Table,
			Columns: surveycellscan.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: surveycellscan.FieldID,
			},
		},
		From:   scsq.sql,
		Unique: true,
	}
	if ps := scsq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := scsq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := scsq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := scsq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (scsq *SurveyCellScanQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(scsq.driver.Dialect())
	t1 := builder.Table(surveycellscan.Table)
	selector := builder.Select(t1.Columns(surveycellscan.Columns...)...).From(t1)
	if scsq.sql != nil {
		selector = scsq.sql
		selector.Select(selector.Columns(surveycellscan.Columns...)...)
	}
	for _, p := range scsq.predicates {
		p(selector)
	}
	for _, p := range scsq.order {
		p(selector)
	}
	if offset := scsq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := scsq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// SurveyCellScanGroupBy is the builder for group-by SurveyCellScan entities.
type SurveyCellScanGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query.
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the group-by query.
func (scsgb *SurveyCellScanGroupBy) Aggregate(fns ...Aggregate) *SurveyCellScanGroupBy {
	scsgb.fns = append(scsgb.fns, fns...)
	return scsgb
}

// Scan applies the group-by query and scan the result into the given value.
func (scsgb *SurveyCellScanGroupBy) Scan(ctx context.Context, v interface{}) error {
	return scsgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (scsgb *SurveyCellScanGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := scsgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (scsgb *SurveyCellScanGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(scsgb.fields) > 1 {
		return nil, errors.New("ent: SurveyCellScanGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := scsgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (scsgb *SurveyCellScanGroupBy) StringsX(ctx context.Context) []string {
	v, err := scsgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (scsgb *SurveyCellScanGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(scsgb.fields) > 1 {
		return nil, errors.New("ent: SurveyCellScanGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := scsgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (scsgb *SurveyCellScanGroupBy) IntsX(ctx context.Context) []int {
	v, err := scsgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (scsgb *SurveyCellScanGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(scsgb.fields) > 1 {
		return nil, errors.New("ent: SurveyCellScanGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := scsgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (scsgb *SurveyCellScanGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := scsgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (scsgb *SurveyCellScanGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(scsgb.fields) > 1 {
		return nil, errors.New("ent: SurveyCellScanGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := scsgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (scsgb *SurveyCellScanGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := scsgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (scsgb *SurveyCellScanGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := scsgb.sqlQuery().Query()
	if err := scsgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (scsgb *SurveyCellScanGroupBy) sqlQuery() *sql.Selector {
	selector := scsgb.sql
	columns := make([]string, 0, len(scsgb.fields)+len(scsgb.fns))
	columns = append(columns, scsgb.fields...)
	for _, fn := range scsgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(scsgb.fields...)
}

// SurveyCellScanSelect is the builder for select fields of SurveyCellScan entities.
type SurveyCellScanSelect struct {
	config
	fields []string
	// intermediate queries.
	sql *sql.Selector
}

// Scan applies the selector query and scan the result into the given value.
func (scss *SurveyCellScanSelect) Scan(ctx context.Context, v interface{}) error {
	return scss.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (scss *SurveyCellScanSelect) ScanX(ctx context.Context, v interface{}) {
	if err := scss.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (scss *SurveyCellScanSelect) Strings(ctx context.Context) ([]string, error) {
	if len(scss.fields) > 1 {
		return nil, errors.New("ent: SurveyCellScanSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := scss.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (scss *SurveyCellScanSelect) StringsX(ctx context.Context) []string {
	v, err := scss.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (scss *SurveyCellScanSelect) Ints(ctx context.Context) ([]int, error) {
	if len(scss.fields) > 1 {
		return nil, errors.New("ent: SurveyCellScanSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := scss.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (scss *SurveyCellScanSelect) IntsX(ctx context.Context) []int {
	v, err := scss.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (scss *SurveyCellScanSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(scss.fields) > 1 {
		return nil, errors.New("ent: SurveyCellScanSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := scss.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (scss *SurveyCellScanSelect) Float64sX(ctx context.Context) []float64 {
	v, err := scss.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (scss *SurveyCellScanSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(scss.fields) > 1 {
		return nil, errors.New("ent: SurveyCellScanSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := scss.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (scss *SurveyCellScanSelect) BoolsX(ctx context.Context) []bool {
	v, err := scss.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (scss *SurveyCellScanSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := scss.sqlQuery().Query()
	if err := scss.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (scss *SurveyCellScanSelect) sqlQuery() sql.Querier {
	selector := scss.sql
	selector.Select(selector.Columns(scss.fields...)...)
	return selector
}
