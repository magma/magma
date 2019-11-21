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
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatecategory"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatequestion"
)

// SurveyTemplateQuestionQuery is the builder for querying SurveyTemplateQuestion entities.
type SurveyTemplateQuestionQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.SurveyTemplateQuestion
	// intermediate queries.
	sql *sql.Selector
}

// Where adds a new predicate for the builder.
func (stqq *SurveyTemplateQuestionQuery) Where(ps ...predicate.SurveyTemplateQuestion) *SurveyTemplateQuestionQuery {
	stqq.predicates = append(stqq.predicates, ps...)
	return stqq
}

// Limit adds a limit step to the query.
func (stqq *SurveyTemplateQuestionQuery) Limit(limit int) *SurveyTemplateQuestionQuery {
	stqq.limit = &limit
	return stqq
}

// Offset adds an offset step to the query.
func (stqq *SurveyTemplateQuestionQuery) Offset(offset int) *SurveyTemplateQuestionQuery {
	stqq.offset = &offset
	return stqq
}

// Order adds an order step to the query.
func (stqq *SurveyTemplateQuestionQuery) Order(o ...Order) *SurveyTemplateQuestionQuery {
	stqq.order = append(stqq.order, o...)
	return stqq
}

// QueryCategory chains the current query on the category edge.
func (stqq *SurveyTemplateQuestionQuery) QueryCategory() *SurveyTemplateCategoryQuery {
	query := &SurveyTemplateCategoryQuery{config: stqq.config}

	builder := sql.Dialect(stqq.driver.Dialect())
	t1 := builder.Table(surveytemplatecategory.Table)
	t2 := stqq.sqlQuery()
	t2.Select(t2.C(surveytemplatequestion.CategoryColumn))
	query.sql = builder.Select(t1.Columns(surveytemplatecategory.Columns...)...).
		From(t1).
		Join(t2).
		On(t1.C(surveytemplatecategory.FieldID), t2.C(surveytemplatequestion.CategoryColumn))
	return query
}

// First returns the first SurveyTemplateQuestion entity in the query. Returns *ErrNotFound when no surveytemplatequestion was found.
func (stqq *SurveyTemplateQuestionQuery) First(ctx context.Context) (*SurveyTemplateQuestion, error) {
	stqs, err := stqq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(stqs) == 0 {
		return nil, &ErrNotFound{surveytemplatequestion.Label}
	}
	return stqs[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (stqq *SurveyTemplateQuestionQuery) FirstX(ctx context.Context) *SurveyTemplateQuestion {
	stq, err := stqq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return stq
}

// FirstID returns the first SurveyTemplateQuestion id in the query. Returns *ErrNotFound when no id was found.
func (stqq *SurveyTemplateQuestionQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = stqq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &ErrNotFound{surveytemplatequestion.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (stqq *SurveyTemplateQuestionQuery) FirstXID(ctx context.Context) string {
	id, err := stqq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only SurveyTemplateQuestion entity in the query, returns an error if not exactly one entity was returned.
func (stqq *SurveyTemplateQuestionQuery) Only(ctx context.Context) (*SurveyTemplateQuestion, error) {
	stqs, err := stqq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(stqs) {
	case 1:
		return stqs[0], nil
	case 0:
		return nil, &ErrNotFound{surveytemplatequestion.Label}
	default:
		return nil, &ErrNotSingular{surveytemplatequestion.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (stqq *SurveyTemplateQuestionQuery) OnlyX(ctx context.Context) *SurveyTemplateQuestion {
	stq, err := stqq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return stq
}

// OnlyID returns the only SurveyTemplateQuestion id in the query, returns an error if not exactly one id was returned.
func (stqq *SurveyTemplateQuestionQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = stqq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &ErrNotFound{surveytemplatequestion.Label}
	default:
		err = &ErrNotSingular{surveytemplatequestion.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (stqq *SurveyTemplateQuestionQuery) OnlyXID(ctx context.Context) string {
	id, err := stqq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of SurveyTemplateQuestions.
func (stqq *SurveyTemplateQuestionQuery) All(ctx context.Context) ([]*SurveyTemplateQuestion, error) {
	return stqq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (stqq *SurveyTemplateQuestionQuery) AllX(ctx context.Context) []*SurveyTemplateQuestion {
	stqs, err := stqq.All(ctx)
	if err != nil {
		panic(err)
	}
	return stqs
}

// IDs executes the query and returns a list of SurveyTemplateQuestion ids.
func (stqq *SurveyTemplateQuestionQuery) IDs(ctx context.Context) ([]string, error) {
	var ids []string
	if err := stqq.Select(surveytemplatequestion.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (stqq *SurveyTemplateQuestionQuery) IDsX(ctx context.Context) []string {
	ids, err := stqq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (stqq *SurveyTemplateQuestionQuery) Count(ctx context.Context) (int, error) {
	return stqq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (stqq *SurveyTemplateQuestionQuery) CountX(ctx context.Context) int {
	count, err := stqq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (stqq *SurveyTemplateQuestionQuery) Exist(ctx context.Context) (bool, error) {
	return stqq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (stqq *SurveyTemplateQuestionQuery) ExistX(ctx context.Context) bool {
	exist, err := stqq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (stqq *SurveyTemplateQuestionQuery) Clone() *SurveyTemplateQuestionQuery {
	return &SurveyTemplateQuestionQuery{
		config:     stqq.config,
		limit:      stqq.limit,
		offset:     stqq.offset,
		order:      append([]Order{}, stqq.order...),
		unique:     append([]string{}, stqq.unique...),
		predicates: append([]predicate.SurveyTemplateQuestion{}, stqq.predicates...),
		// clone intermediate queries.
		sql: stqq.sql.Clone(),
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
//	client.SurveyTemplateQuestion.Query().
//		GroupBy(surveytemplatequestion.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (stqq *SurveyTemplateQuestionQuery) GroupBy(field string, fields ...string) *SurveyTemplateQuestionGroupBy {
	group := &SurveyTemplateQuestionGroupBy{config: stqq.config}
	group.fields = append([]string{field}, fields...)
	group.sql = stqq.sqlQuery()
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
//	client.SurveyTemplateQuestion.Query().
//		Select(surveytemplatequestion.FieldCreateTime).
//		Scan(ctx, &v)
//
func (stqq *SurveyTemplateQuestionQuery) Select(field string, fields ...string) *SurveyTemplateQuestionSelect {
	selector := &SurveyTemplateQuestionSelect{config: stqq.config}
	selector.fields = append([]string{field}, fields...)
	selector.sql = stqq.sqlQuery()
	return selector
}

func (stqq *SurveyTemplateQuestionQuery) sqlAll(ctx context.Context) ([]*SurveyTemplateQuestion, error) {
	rows := &sql.Rows{}
	selector := stqq.sqlQuery()
	if unique := stqq.unique; len(unique) == 0 {
		selector.Distinct()
	}
	query, args := selector.Query()
	if err := stqq.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()
	var stqs SurveyTemplateQuestions
	if err := stqs.FromRows(rows); err != nil {
		return nil, err
	}
	stqs.config(stqq.config)
	return stqs, nil
}

func (stqq *SurveyTemplateQuestionQuery) sqlCount(ctx context.Context) (int, error) {
	rows := &sql.Rows{}
	selector := stqq.sqlQuery()
	unique := []string{surveytemplatequestion.FieldID}
	if len(stqq.unique) > 0 {
		unique = stqq.unique
	}
	selector.Count(sql.Distinct(selector.Columns(unique...)...))
	query, args := selector.Query()
	if err := stqq.driver.Query(ctx, query, args, rows); err != nil {
		return 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, errors.New("ent: no rows found")
	}
	var n int
	if err := rows.Scan(&n); err != nil {
		return 0, fmt.Errorf("ent: failed reading count: %v", err)
	}
	return n, nil
}

func (stqq *SurveyTemplateQuestionQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := stqq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (stqq *SurveyTemplateQuestionQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(stqq.driver.Dialect())
	t1 := builder.Table(surveytemplatequestion.Table)
	selector := builder.Select(t1.Columns(surveytemplatequestion.Columns...)...).From(t1)
	if stqq.sql != nil {
		selector = stqq.sql
		selector.Select(selector.Columns(surveytemplatequestion.Columns...)...)
	}
	for _, p := range stqq.predicates {
		p(selector)
	}
	for _, p := range stqq.order {
		p(selector)
	}
	if offset := stqq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := stqq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// SurveyTemplateQuestionGroupBy is the builder for group-by SurveyTemplateQuestion entities.
type SurveyTemplateQuestionGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate queries.
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the group-by query.
func (stqgb *SurveyTemplateQuestionGroupBy) Aggregate(fns ...Aggregate) *SurveyTemplateQuestionGroupBy {
	stqgb.fns = append(stqgb.fns, fns...)
	return stqgb
}

// Scan applies the group-by query and scan the result into the given value.
func (stqgb *SurveyTemplateQuestionGroupBy) Scan(ctx context.Context, v interface{}) error {
	return stqgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (stqgb *SurveyTemplateQuestionGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := stqgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (stqgb *SurveyTemplateQuestionGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(stqgb.fields) > 1 {
		return nil, errors.New("ent: SurveyTemplateQuestionGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := stqgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (stqgb *SurveyTemplateQuestionGroupBy) StringsX(ctx context.Context) []string {
	v, err := stqgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (stqgb *SurveyTemplateQuestionGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(stqgb.fields) > 1 {
		return nil, errors.New("ent: SurveyTemplateQuestionGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := stqgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (stqgb *SurveyTemplateQuestionGroupBy) IntsX(ctx context.Context) []int {
	v, err := stqgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (stqgb *SurveyTemplateQuestionGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(stqgb.fields) > 1 {
		return nil, errors.New("ent: SurveyTemplateQuestionGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := stqgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (stqgb *SurveyTemplateQuestionGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := stqgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (stqgb *SurveyTemplateQuestionGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(stqgb.fields) > 1 {
		return nil, errors.New("ent: SurveyTemplateQuestionGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := stqgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (stqgb *SurveyTemplateQuestionGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := stqgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (stqgb *SurveyTemplateQuestionGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := stqgb.sqlQuery().Query()
	if err := stqgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (stqgb *SurveyTemplateQuestionGroupBy) sqlQuery() *sql.Selector {
	selector := stqgb.sql
	columns := make([]string, 0, len(stqgb.fields)+len(stqgb.fns))
	columns = append(columns, stqgb.fields...)
	for _, fn := range stqgb.fns {
		columns = append(columns, fn.SQL(selector))
	}
	return selector.Select(columns...).GroupBy(stqgb.fields...)
}

// SurveyTemplateQuestionSelect is the builder for select fields of SurveyTemplateQuestion entities.
type SurveyTemplateQuestionSelect struct {
	config
	fields []string
	// intermediate queries.
	sql *sql.Selector
}

// Scan applies the selector query and scan the result into the given value.
func (stqs *SurveyTemplateQuestionSelect) Scan(ctx context.Context, v interface{}) error {
	return stqs.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (stqs *SurveyTemplateQuestionSelect) ScanX(ctx context.Context, v interface{}) {
	if err := stqs.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (stqs *SurveyTemplateQuestionSelect) Strings(ctx context.Context) ([]string, error) {
	if len(stqs.fields) > 1 {
		return nil, errors.New("ent: SurveyTemplateQuestionSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := stqs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (stqs *SurveyTemplateQuestionSelect) StringsX(ctx context.Context) []string {
	v, err := stqs.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (stqs *SurveyTemplateQuestionSelect) Ints(ctx context.Context) ([]int, error) {
	if len(stqs.fields) > 1 {
		return nil, errors.New("ent: SurveyTemplateQuestionSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := stqs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (stqs *SurveyTemplateQuestionSelect) IntsX(ctx context.Context) []int {
	v, err := stqs.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (stqs *SurveyTemplateQuestionSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(stqs.fields) > 1 {
		return nil, errors.New("ent: SurveyTemplateQuestionSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := stqs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (stqs *SurveyTemplateQuestionSelect) Float64sX(ctx context.Context) []float64 {
	v, err := stqs.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (stqs *SurveyTemplateQuestionSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(stqs.fields) > 1 {
		return nil, errors.New("ent: SurveyTemplateQuestionSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := stqs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (stqs *SurveyTemplateQuestionSelect) BoolsX(ctx context.Context) []bool {
	v, err := stqs.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (stqs *SurveyTemplateQuestionSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := stqs.sqlQuery().Query()
	if err := stqs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (stqs *SurveyTemplateQuestionSelect) sqlQuery() sql.Querier {
	view := "surveytemplatequestion_view"
	return sql.Dialect(stqs.driver.Dialect()).
		Select(stqs.fields...).From(stqs.sql.As(view))
}
