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
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
)

// SurveyQuestionQuery is the builder for querying SurveyQuestion entities.
type SurveyQuestionQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.SurveyQuestion
	// intermediate query.
	sql *sql.Selector
}

// Where adds a new predicate for the builder.
func (sqq *SurveyQuestionQuery) Where(ps ...predicate.SurveyQuestion) *SurveyQuestionQuery {
	sqq.predicates = append(sqq.predicates, ps...)
	return sqq
}

// Limit adds a limit step to the query.
func (sqq *SurveyQuestionQuery) Limit(limit int) *SurveyQuestionQuery {
	sqq.limit = &limit
	return sqq
}

// Offset adds an offset step to the query.
func (sqq *SurveyQuestionQuery) Offset(offset int) *SurveyQuestionQuery {
	sqq.offset = &offset
	return sqq
}

// Order adds an order step to the query.
func (sqq *SurveyQuestionQuery) Order(o ...Order) *SurveyQuestionQuery {
	sqq.order = append(sqq.order, o...)
	return sqq
}

// QuerySurvey chains the current query on the survey edge.
func (sqq *SurveyQuestionQuery) QuerySurvey() *SurveyQuery {
	query := &SurveyQuery{config: sqq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, sqq.sqlQuery()),
		sqlgraph.To(survey.Table, survey.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, surveyquestion.SurveyTable, surveyquestion.SurveyColumn),
	)
	query.sql = sqlgraph.SetNeighbors(sqq.driver.Dialect(), step)
	return query
}

// QueryWifiScan chains the current query on the wifi_scan edge.
func (sqq *SurveyQuestionQuery) QueryWifiScan() *SurveyWiFiScanQuery {
	query := &SurveyWiFiScanQuery{config: sqq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, sqq.sqlQuery()),
		sqlgraph.To(surveywifiscan.Table, surveywifiscan.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, surveyquestion.WifiScanTable, surveyquestion.WifiScanColumn),
	)
	query.sql = sqlgraph.SetNeighbors(sqq.driver.Dialect(), step)
	return query
}

// QueryCellScan chains the current query on the cell_scan edge.
func (sqq *SurveyQuestionQuery) QueryCellScan() *SurveyCellScanQuery {
	query := &SurveyCellScanQuery{config: sqq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, sqq.sqlQuery()),
		sqlgraph.To(surveycellscan.Table, surveycellscan.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, surveyquestion.CellScanTable, surveyquestion.CellScanColumn),
	)
	query.sql = sqlgraph.SetNeighbors(sqq.driver.Dialect(), step)
	return query
}

// QueryPhotoData chains the current query on the photo_data edge.
func (sqq *SurveyQuestionQuery) QueryPhotoData() *FileQuery {
	query := &FileQuery{config: sqq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, sqq.sqlQuery()),
		sqlgraph.To(file.Table, file.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, surveyquestion.PhotoDataTable, surveyquestion.PhotoDataColumn),
	)
	query.sql = sqlgraph.SetNeighbors(sqq.driver.Dialect(), step)
	return query
}

// First returns the first SurveyQuestion entity in the query. Returns *ErrNotFound when no surveyquestion was found.
func (sqq *SurveyQuestionQuery) First(ctx context.Context) (*SurveyQuestion, error) {
	sqs, err := sqq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(sqs) == 0 {
		return nil, &ErrNotFound{surveyquestion.Label}
	}
	return sqs[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (sqq *SurveyQuestionQuery) FirstX(ctx context.Context) *SurveyQuestion {
	sq, err := sqq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return sq
}

// FirstID returns the first SurveyQuestion id in the query. Returns *ErrNotFound when no id was found.
func (sqq *SurveyQuestionQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = sqq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &ErrNotFound{surveyquestion.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (sqq *SurveyQuestionQuery) FirstXID(ctx context.Context) string {
	id, err := sqq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only SurveyQuestion entity in the query, returns an error if not exactly one entity was returned.
func (sqq *SurveyQuestionQuery) Only(ctx context.Context) (*SurveyQuestion, error) {
	sqs, err := sqq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(sqs) {
	case 1:
		return sqs[0], nil
	case 0:
		return nil, &ErrNotFound{surveyquestion.Label}
	default:
		return nil, &ErrNotSingular{surveyquestion.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (sqq *SurveyQuestionQuery) OnlyX(ctx context.Context) *SurveyQuestion {
	sq, err := sqq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return sq
}

// OnlyID returns the only SurveyQuestion id in the query, returns an error if not exactly one id was returned.
func (sqq *SurveyQuestionQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = sqq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &ErrNotFound{surveyquestion.Label}
	default:
		err = &ErrNotSingular{surveyquestion.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (sqq *SurveyQuestionQuery) OnlyXID(ctx context.Context) string {
	id, err := sqq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of SurveyQuestions.
func (sqq *SurveyQuestionQuery) All(ctx context.Context) ([]*SurveyQuestion, error) {
	return sqq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (sqq *SurveyQuestionQuery) AllX(ctx context.Context) []*SurveyQuestion {
	sqs, err := sqq.All(ctx)
	if err != nil {
		panic(err)
	}
	return sqs
}

// IDs executes the query and returns a list of SurveyQuestion ids.
func (sqq *SurveyQuestionQuery) IDs(ctx context.Context) ([]string, error) {
	var ids []string
	if err := sqq.Select(surveyquestion.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (sqq *SurveyQuestionQuery) IDsX(ctx context.Context) []string {
	ids, err := sqq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (sqq *SurveyQuestionQuery) Count(ctx context.Context) (int, error) {
	return sqq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (sqq *SurveyQuestionQuery) CountX(ctx context.Context) int {
	count, err := sqq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (sqq *SurveyQuestionQuery) Exist(ctx context.Context) (bool, error) {
	return sqq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (sqq *SurveyQuestionQuery) ExistX(ctx context.Context) bool {
	exist, err := sqq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (sqq *SurveyQuestionQuery) Clone() *SurveyQuestionQuery {
	return &SurveyQuestionQuery{
		config:     sqq.config,
		limit:      sqq.limit,
		offset:     sqq.offset,
		order:      append([]Order{}, sqq.order...),
		unique:     append([]string{}, sqq.unique...),
		predicates: append([]predicate.SurveyQuestion{}, sqq.predicates...),
		// clone intermediate query.
		sql: sqq.sql.Clone(),
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
//	client.SurveyQuestion.Query().
//		GroupBy(surveyquestion.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (sqq *SurveyQuestionQuery) GroupBy(field string, fields ...string) *SurveyQuestionGroupBy {
	group := &SurveyQuestionGroupBy{config: sqq.config}
	group.fields = append([]string{field}, fields...)
	group.sql = sqq.sqlQuery()
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
//	client.SurveyQuestion.Query().
//		Select(surveyquestion.FieldCreateTime).
//		Scan(ctx, &v)
//
func (sqq *SurveyQuestionQuery) Select(field string, fields ...string) *SurveyQuestionSelect {
	selector := &SurveyQuestionSelect{config: sqq.config}
	selector.fields = append([]string{field}, fields...)
	selector.sql = sqq.sqlQuery()
	return selector
}

func (sqq *SurveyQuestionQuery) sqlAll(ctx context.Context) ([]*SurveyQuestion, error) {
	var (
		nodes []*SurveyQuestion
		spec  = sqq.querySpec()
	)
	spec.ScanValues = func() []interface{} {
		node := &SurveyQuestion{config: sqq.config}
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
	if err := sqlgraph.QueryNodes(ctx, sqq.driver, spec); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (sqq *SurveyQuestionQuery) sqlCount(ctx context.Context) (int, error) {
	spec := sqq.querySpec()
	return sqlgraph.CountNodes(ctx, sqq.driver, spec)
}

func (sqq *SurveyQuestionQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := sqq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (sqq *SurveyQuestionQuery) querySpec() *sqlgraph.QuerySpec {
	spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveyquestion.Table,
			Columns: surveyquestion.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: surveyquestion.FieldID,
			},
		},
		From:   sqq.sql,
		Unique: true,
	}
	if ps := sqq.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := sqq.limit; limit != nil {
		spec.Limit = *limit
	}
	if offset := sqq.offset; offset != nil {
		spec.Offset = *offset
	}
	if ps := sqq.order; len(ps) > 0 {
		spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return spec
}

func (sqq *SurveyQuestionQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(sqq.driver.Dialect())
	t1 := builder.Table(surveyquestion.Table)
	selector := builder.Select(t1.Columns(surveyquestion.Columns...)...).From(t1)
	if sqq.sql != nil {
		selector = sqq.sql
		selector.Select(selector.Columns(surveyquestion.Columns...)...)
	}
	for _, p := range sqq.predicates {
		p(selector)
	}
	for _, p := range sqq.order {
		p(selector)
	}
	if offset := sqq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := sqq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// SurveyQuestionGroupBy is the builder for group-by SurveyQuestion entities.
type SurveyQuestionGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query.
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the group-by query.
func (sqgb *SurveyQuestionGroupBy) Aggregate(fns ...Aggregate) *SurveyQuestionGroupBy {
	sqgb.fns = append(sqgb.fns, fns...)
	return sqgb
}

// Scan applies the group-by query and scan the result into the given value.
func (sqgb *SurveyQuestionGroupBy) Scan(ctx context.Context, v interface{}) error {
	return sqgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (sqgb *SurveyQuestionGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := sqgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (sqgb *SurveyQuestionGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(sqgb.fields) > 1 {
		return nil, errors.New("ent: SurveyQuestionGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := sqgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (sqgb *SurveyQuestionGroupBy) StringsX(ctx context.Context) []string {
	v, err := sqgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (sqgb *SurveyQuestionGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(sqgb.fields) > 1 {
		return nil, errors.New("ent: SurveyQuestionGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := sqgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (sqgb *SurveyQuestionGroupBy) IntsX(ctx context.Context) []int {
	v, err := sqgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (sqgb *SurveyQuestionGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(sqgb.fields) > 1 {
		return nil, errors.New("ent: SurveyQuestionGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := sqgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (sqgb *SurveyQuestionGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := sqgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (sqgb *SurveyQuestionGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(sqgb.fields) > 1 {
		return nil, errors.New("ent: SurveyQuestionGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := sqgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (sqgb *SurveyQuestionGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := sqgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (sqgb *SurveyQuestionGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := sqgb.sqlQuery().Query()
	if err := sqgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (sqgb *SurveyQuestionGroupBy) sqlQuery() *sql.Selector {
	selector := sqgb.sql
	columns := make([]string, 0, len(sqgb.fields)+len(sqgb.fns))
	columns = append(columns, sqgb.fields...)
	for _, fn := range sqgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(sqgb.fields...)
}

// SurveyQuestionSelect is the builder for select fields of SurveyQuestion entities.
type SurveyQuestionSelect struct {
	config
	fields []string
	// intermediate queries.
	sql *sql.Selector
}

// Scan applies the selector query and scan the result into the given value.
func (sqs *SurveyQuestionSelect) Scan(ctx context.Context, v interface{}) error {
	return sqs.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (sqs *SurveyQuestionSelect) ScanX(ctx context.Context, v interface{}) {
	if err := sqs.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (sqs *SurveyQuestionSelect) Strings(ctx context.Context) ([]string, error) {
	if len(sqs.fields) > 1 {
		return nil, errors.New("ent: SurveyQuestionSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := sqs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (sqs *SurveyQuestionSelect) StringsX(ctx context.Context) []string {
	v, err := sqs.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (sqs *SurveyQuestionSelect) Ints(ctx context.Context) ([]int, error) {
	if len(sqs.fields) > 1 {
		return nil, errors.New("ent: SurveyQuestionSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := sqs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (sqs *SurveyQuestionSelect) IntsX(ctx context.Context) []int {
	v, err := sqs.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (sqs *SurveyQuestionSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(sqs.fields) > 1 {
		return nil, errors.New("ent: SurveyQuestionSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := sqs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (sqs *SurveyQuestionSelect) Float64sX(ctx context.Context) []float64 {
	v, err := sqs.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (sqs *SurveyQuestionSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(sqs.fields) > 1 {
		return nil, errors.New("ent: SurveyQuestionSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := sqs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (sqs *SurveyQuestionSelect) BoolsX(ctx context.Context) []bool {
	v, err := sqs.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (sqs *SurveyQuestionSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := sqs.sqlQuery().Query()
	if err := sqs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (sqs *SurveyQuestionSelect) sqlQuery() sql.Querier {
	selector := sqs.sql
	selector.Select(selector.Columns(sqs.fields...)...)
	return selector
}
