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
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
)

// SurveyWiFiScanQuery is the builder for querying SurveyWiFiScan entities.
type SurveyWiFiScanQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.SurveyWiFiScan
	// eager-loading edges.
	withSurveyQuestion *SurveyQuestionQuery
	withLocation       *LocationQuery
	withFKs            bool
	// intermediate query.
	sql *sql.Selector
}

// Where adds a new predicate for the builder.
func (swfsq *SurveyWiFiScanQuery) Where(ps ...predicate.SurveyWiFiScan) *SurveyWiFiScanQuery {
	swfsq.predicates = append(swfsq.predicates, ps...)
	return swfsq
}

// Limit adds a limit step to the query.
func (swfsq *SurveyWiFiScanQuery) Limit(limit int) *SurveyWiFiScanQuery {
	swfsq.limit = &limit
	return swfsq
}

// Offset adds an offset step to the query.
func (swfsq *SurveyWiFiScanQuery) Offset(offset int) *SurveyWiFiScanQuery {
	swfsq.offset = &offset
	return swfsq
}

// Order adds an order step to the query.
func (swfsq *SurveyWiFiScanQuery) Order(o ...Order) *SurveyWiFiScanQuery {
	swfsq.order = append(swfsq.order, o...)
	return swfsq
}

// QuerySurveyQuestion chains the current query on the survey_question edge.
func (swfsq *SurveyWiFiScanQuery) QuerySurveyQuestion() *SurveyQuestionQuery {
	query := &SurveyQuestionQuery{config: swfsq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(surveywifiscan.Table, surveywifiscan.FieldID, swfsq.sqlQuery()),
		sqlgraph.To(surveyquestion.Table, surveyquestion.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, surveywifiscan.SurveyQuestionTable, surveywifiscan.SurveyQuestionColumn),
	)
	query.sql = sqlgraph.SetNeighbors(swfsq.driver.Dialect(), step)
	return query
}

// QueryLocation chains the current query on the location edge.
func (swfsq *SurveyWiFiScanQuery) QueryLocation() *LocationQuery {
	query := &LocationQuery{config: swfsq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(surveywifiscan.Table, surveywifiscan.FieldID, swfsq.sqlQuery()),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, surveywifiscan.LocationTable, surveywifiscan.LocationColumn),
	)
	query.sql = sqlgraph.SetNeighbors(swfsq.driver.Dialect(), step)
	return query
}

// First returns the first SurveyWiFiScan entity in the query. Returns *NotFoundError when no surveywifiscan was found.
func (swfsq *SurveyWiFiScanQuery) First(ctx context.Context) (*SurveyWiFiScan, error) {
	swfsSlice, err := swfsq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(swfsSlice) == 0 {
		return nil, &NotFoundError{surveywifiscan.Label}
	}
	return swfsSlice[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (swfsq *SurveyWiFiScanQuery) FirstX(ctx context.Context) *SurveyWiFiScan {
	swfs, err := swfsq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return swfs
}

// FirstID returns the first SurveyWiFiScan id in the query. Returns *NotFoundError when no id was found.
func (swfsq *SurveyWiFiScanQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = swfsq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{surveywifiscan.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (swfsq *SurveyWiFiScanQuery) FirstXID(ctx context.Context) string {
	id, err := swfsq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only SurveyWiFiScan entity in the query, returns an error if not exactly one entity was returned.
func (swfsq *SurveyWiFiScanQuery) Only(ctx context.Context) (*SurveyWiFiScan, error) {
	swfsSlice, err := swfsq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(swfsSlice) {
	case 1:
		return swfsSlice[0], nil
	case 0:
		return nil, &NotFoundError{surveywifiscan.Label}
	default:
		return nil, &NotSingularError{surveywifiscan.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (swfsq *SurveyWiFiScanQuery) OnlyX(ctx context.Context) *SurveyWiFiScan {
	swfs, err := swfsq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return swfs
}

// OnlyID returns the only SurveyWiFiScan id in the query, returns an error if not exactly one id was returned.
func (swfsq *SurveyWiFiScanQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = swfsq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{surveywifiscan.Label}
	default:
		err = &NotSingularError{surveywifiscan.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (swfsq *SurveyWiFiScanQuery) OnlyXID(ctx context.Context) string {
	id, err := swfsq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of SurveyWiFiScans.
func (swfsq *SurveyWiFiScanQuery) All(ctx context.Context) ([]*SurveyWiFiScan, error) {
	return swfsq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (swfsq *SurveyWiFiScanQuery) AllX(ctx context.Context) []*SurveyWiFiScan {
	swfsSlice, err := swfsq.All(ctx)
	if err != nil {
		panic(err)
	}
	return swfsSlice
}

// IDs executes the query and returns a list of SurveyWiFiScan ids.
func (swfsq *SurveyWiFiScanQuery) IDs(ctx context.Context) ([]string, error) {
	var ids []string
	if err := swfsq.Select(surveywifiscan.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (swfsq *SurveyWiFiScanQuery) IDsX(ctx context.Context) []string {
	ids, err := swfsq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (swfsq *SurveyWiFiScanQuery) Count(ctx context.Context) (int, error) {
	return swfsq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (swfsq *SurveyWiFiScanQuery) CountX(ctx context.Context) int {
	count, err := swfsq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (swfsq *SurveyWiFiScanQuery) Exist(ctx context.Context) (bool, error) {
	return swfsq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (swfsq *SurveyWiFiScanQuery) ExistX(ctx context.Context) bool {
	exist, err := swfsq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (swfsq *SurveyWiFiScanQuery) Clone() *SurveyWiFiScanQuery {
	return &SurveyWiFiScanQuery{
		config:     swfsq.config,
		limit:      swfsq.limit,
		offset:     swfsq.offset,
		order:      append([]Order{}, swfsq.order...),
		unique:     append([]string{}, swfsq.unique...),
		predicates: append([]predicate.SurveyWiFiScan{}, swfsq.predicates...),
		// clone intermediate query.
		sql: swfsq.sql.Clone(),
	}
}

//  WithSurveyQuestion tells the query-builder to eager-loads the nodes that are connected to
// the "survey_question" edge. The optional arguments used to configure the query builder of the edge.
func (swfsq *SurveyWiFiScanQuery) WithSurveyQuestion(opts ...func(*SurveyQuestionQuery)) *SurveyWiFiScanQuery {
	query := &SurveyQuestionQuery{config: swfsq.config}
	for _, opt := range opts {
		opt(query)
	}
	swfsq.withSurveyQuestion = query
	return swfsq
}

//  WithLocation tells the query-builder to eager-loads the nodes that are connected to
// the "location" edge. The optional arguments used to configure the query builder of the edge.
func (swfsq *SurveyWiFiScanQuery) WithLocation(opts ...func(*LocationQuery)) *SurveyWiFiScanQuery {
	query := &LocationQuery{config: swfsq.config}
	for _, opt := range opts {
		opt(query)
	}
	swfsq.withLocation = query
	return swfsq
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
//	client.SurveyWiFiScan.Query().
//		GroupBy(surveywifiscan.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (swfsq *SurveyWiFiScanQuery) GroupBy(field string, fields ...string) *SurveyWiFiScanGroupBy {
	group := &SurveyWiFiScanGroupBy{config: swfsq.config}
	group.fields = append([]string{field}, fields...)
	group.sql = swfsq.sqlQuery()
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
//	client.SurveyWiFiScan.Query().
//		Select(surveywifiscan.FieldCreateTime).
//		Scan(ctx, &v)
//
func (swfsq *SurveyWiFiScanQuery) Select(field string, fields ...string) *SurveyWiFiScanSelect {
	selector := &SurveyWiFiScanSelect{config: swfsq.config}
	selector.fields = append([]string{field}, fields...)
	selector.sql = swfsq.sqlQuery()
	return selector
}

func (swfsq *SurveyWiFiScanQuery) sqlAll(ctx context.Context) ([]*SurveyWiFiScan, error) {
	var (
		nodes   []*SurveyWiFiScan = []*SurveyWiFiScan{}
		withFKs                   = swfsq.withFKs
		_spec                     = swfsq.querySpec()
	)
	if swfsq.withSurveyQuestion != nil || swfsq.withLocation != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, surveywifiscan.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &SurveyWiFiScan{config: swfsq.config}
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
		return node.assignValues(values...)
	}
	if err := sqlgraph.QueryNodes(ctx, swfsq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := swfsq.withSurveyQuestion; query != nil {
		ids := make([]string, 0, len(nodes))
		nodeids := make(map[string][]*SurveyWiFiScan)
		for i := range nodes {
			if fk := nodes[i].survey_question_id; fk != nil {
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
				return nil, fmt.Errorf(`unexpected foreign-key "survey_question_id" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.SurveyQuestion = n
			}
		}
	}

	if query := swfsq.withLocation; query != nil {
		ids := make([]string, 0, len(nodes))
		nodeids := make(map[string][]*SurveyWiFiScan)
		for i := range nodes {
			if fk := nodes[i].location_id; fk != nil {
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
				return nil, fmt.Errorf(`unexpected foreign-key "location_id" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Location = n
			}
		}
	}

	return nodes, nil
}

func (swfsq *SurveyWiFiScanQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := swfsq.querySpec()
	return sqlgraph.CountNodes(ctx, swfsq.driver, _spec)
}

func (swfsq *SurveyWiFiScanQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := swfsq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (swfsq *SurveyWiFiScanQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveywifiscan.Table,
			Columns: surveywifiscan.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: surveywifiscan.FieldID,
			},
		},
		From:   swfsq.sql,
		Unique: true,
	}
	if ps := swfsq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := swfsq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := swfsq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := swfsq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (swfsq *SurveyWiFiScanQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(swfsq.driver.Dialect())
	t1 := builder.Table(surveywifiscan.Table)
	selector := builder.Select(t1.Columns(surveywifiscan.Columns...)...).From(t1)
	if swfsq.sql != nil {
		selector = swfsq.sql
		selector.Select(selector.Columns(surveywifiscan.Columns...)...)
	}
	for _, p := range swfsq.predicates {
		p(selector)
	}
	for _, p := range swfsq.order {
		p(selector)
	}
	if offset := swfsq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := swfsq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// SurveyWiFiScanGroupBy is the builder for group-by SurveyWiFiScan entities.
type SurveyWiFiScanGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query.
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the group-by query.
func (swfsgb *SurveyWiFiScanGroupBy) Aggregate(fns ...Aggregate) *SurveyWiFiScanGroupBy {
	swfsgb.fns = append(swfsgb.fns, fns...)
	return swfsgb
}

// Scan applies the group-by query and scan the result into the given value.
func (swfsgb *SurveyWiFiScanGroupBy) Scan(ctx context.Context, v interface{}) error {
	return swfsgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (swfsgb *SurveyWiFiScanGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := swfsgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (swfsgb *SurveyWiFiScanGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(swfsgb.fields) > 1 {
		return nil, errors.New("ent: SurveyWiFiScanGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := swfsgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (swfsgb *SurveyWiFiScanGroupBy) StringsX(ctx context.Context) []string {
	v, err := swfsgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (swfsgb *SurveyWiFiScanGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(swfsgb.fields) > 1 {
		return nil, errors.New("ent: SurveyWiFiScanGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := swfsgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (swfsgb *SurveyWiFiScanGroupBy) IntsX(ctx context.Context) []int {
	v, err := swfsgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (swfsgb *SurveyWiFiScanGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(swfsgb.fields) > 1 {
		return nil, errors.New("ent: SurveyWiFiScanGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := swfsgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (swfsgb *SurveyWiFiScanGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := swfsgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (swfsgb *SurveyWiFiScanGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(swfsgb.fields) > 1 {
		return nil, errors.New("ent: SurveyWiFiScanGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := swfsgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (swfsgb *SurveyWiFiScanGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := swfsgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (swfsgb *SurveyWiFiScanGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := swfsgb.sqlQuery().Query()
	if err := swfsgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (swfsgb *SurveyWiFiScanGroupBy) sqlQuery() *sql.Selector {
	selector := swfsgb.sql
	columns := make([]string, 0, len(swfsgb.fields)+len(swfsgb.fns))
	columns = append(columns, swfsgb.fields...)
	for _, fn := range swfsgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(swfsgb.fields...)
}

// SurveyWiFiScanSelect is the builder for select fields of SurveyWiFiScan entities.
type SurveyWiFiScanSelect struct {
	config
	fields []string
	// intermediate queries.
	sql *sql.Selector
}

// Scan applies the selector query and scan the result into the given value.
func (swfss *SurveyWiFiScanSelect) Scan(ctx context.Context, v interface{}) error {
	return swfss.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (swfss *SurveyWiFiScanSelect) ScanX(ctx context.Context, v interface{}) {
	if err := swfss.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (swfss *SurveyWiFiScanSelect) Strings(ctx context.Context) ([]string, error) {
	if len(swfss.fields) > 1 {
		return nil, errors.New("ent: SurveyWiFiScanSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := swfss.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (swfss *SurveyWiFiScanSelect) StringsX(ctx context.Context) []string {
	v, err := swfss.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (swfss *SurveyWiFiScanSelect) Ints(ctx context.Context) ([]int, error) {
	if len(swfss.fields) > 1 {
		return nil, errors.New("ent: SurveyWiFiScanSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := swfss.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (swfss *SurveyWiFiScanSelect) IntsX(ctx context.Context) []int {
	v, err := swfss.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (swfss *SurveyWiFiScanSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(swfss.fields) > 1 {
		return nil, errors.New("ent: SurveyWiFiScanSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := swfss.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (swfss *SurveyWiFiScanSelect) Float64sX(ctx context.Context) []float64 {
	v, err := swfss.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (swfss *SurveyWiFiScanSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(swfss.fields) > 1 {
		return nil, errors.New("ent: SurveyWiFiScanSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := swfss.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (swfss *SurveyWiFiScanSelect) BoolsX(ctx context.Context) []bool {
	v, err := swfss.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (swfss *SurveyWiFiScanSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := swfss.sqlQuery().Query()
	if err := swfss.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (swfss *SurveyWiFiScanSelect) sqlQuery() sql.Querier {
	selector := swfss.sql
	selector.Select(selector.Columns(swfss.fields...)...)
	return selector
}
