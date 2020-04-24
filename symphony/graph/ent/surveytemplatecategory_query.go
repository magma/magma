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

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatecategory"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatequestion"
)

// SurveyTemplateCategoryQuery is the builder for querying SurveyTemplateCategory entities.
type SurveyTemplateCategoryQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.SurveyTemplateCategory
	// eager-loading edges.
	withSurveyTemplateQuestions *SurveyTemplateQuestionQuery
	withFKs                     bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (stcq *SurveyTemplateCategoryQuery) Where(ps ...predicate.SurveyTemplateCategory) *SurveyTemplateCategoryQuery {
	stcq.predicates = append(stcq.predicates, ps...)
	return stcq
}

// Limit adds a limit step to the query.
func (stcq *SurveyTemplateCategoryQuery) Limit(limit int) *SurveyTemplateCategoryQuery {
	stcq.limit = &limit
	return stcq
}

// Offset adds an offset step to the query.
func (stcq *SurveyTemplateCategoryQuery) Offset(offset int) *SurveyTemplateCategoryQuery {
	stcq.offset = &offset
	return stcq
}

// Order adds an order step to the query.
func (stcq *SurveyTemplateCategoryQuery) Order(o ...Order) *SurveyTemplateCategoryQuery {
	stcq.order = append(stcq.order, o...)
	return stcq
}

// QuerySurveyTemplateQuestions chains the current query on the survey_template_questions edge.
func (stcq *SurveyTemplateCategoryQuery) QuerySurveyTemplateQuestions() *SurveyTemplateQuestionQuery {
	query := &SurveyTemplateQuestionQuery{config: stcq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := stcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(surveytemplatecategory.Table, surveytemplatecategory.FieldID, stcq.sqlQuery()),
			sqlgraph.To(surveytemplatequestion.Table, surveytemplatequestion.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, surveytemplatecategory.SurveyTemplateQuestionsTable, surveytemplatecategory.SurveyTemplateQuestionsColumn),
		)
		fromU = sqlgraph.SetNeighbors(stcq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first SurveyTemplateCategory entity in the query. Returns *NotFoundError when no surveytemplatecategory was found.
func (stcq *SurveyTemplateCategoryQuery) First(ctx context.Context) (*SurveyTemplateCategory, error) {
	stcs, err := stcq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(stcs) == 0 {
		return nil, &NotFoundError{surveytemplatecategory.Label}
	}
	return stcs[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (stcq *SurveyTemplateCategoryQuery) FirstX(ctx context.Context) *SurveyTemplateCategory {
	stc, err := stcq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return stc
}

// FirstID returns the first SurveyTemplateCategory id in the query. Returns *NotFoundError when no id was found.
func (stcq *SurveyTemplateCategoryQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = stcq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{surveytemplatecategory.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (stcq *SurveyTemplateCategoryQuery) FirstXID(ctx context.Context) int {
	id, err := stcq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only SurveyTemplateCategory entity in the query, returns an error if not exactly one entity was returned.
func (stcq *SurveyTemplateCategoryQuery) Only(ctx context.Context) (*SurveyTemplateCategory, error) {
	stcs, err := stcq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(stcs) {
	case 1:
		return stcs[0], nil
	case 0:
		return nil, &NotFoundError{surveytemplatecategory.Label}
	default:
		return nil, &NotSingularError{surveytemplatecategory.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (stcq *SurveyTemplateCategoryQuery) OnlyX(ctx context.Context) *SurveyTemplateCategory {
	stc, err := stcq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return stc
}

// OnlyID returns the only SurveyTemplateCategory id in the query, returns an error if not exactly one id was returned.
func (stcq *SurveyTemplateCategoryQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = stcq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{surveytemplatecategory.Label}
	default:
		err = &NotSingularError{surveytemplatecategory.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (stcq *SurveyTemplateCategoryQuery) OnlyXID(ctx context.Context) int {
	id, err := stcq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of SurveyTemplateCategories.
func (stcq *SurveyTemplateCategoryQuery) All(ctx context.Context) ([]*SurveyTemplateCategory, error) {
	if err := stcq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return stcq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (stcq *SurveyTemplateCategoryQuery) AllX(ctx context.Context) []*SurveyTemplateCategory {
	stcs, err := stcq.All(ctx)
	if err != nil {
		panic(err)
	}
	return stcs
}

// IDs executes the query and returns a list of SurveyTemplateCategory ids.
func (stcq *SurveyTemplateCategoryQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := stcq.Select(surveytemplatecategory.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (stcq *SurveyTemplateCategoryQuery) IDsX(ctx context.Context) []int {
	ids, err := stcq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (stcq *SurveyTemplateCategoryQuery) Count(ctx context.Context) (int, error) {
	if err := stcq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return stcq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (stcq *SurveyTemplateCategoryQuery) CountX(ctx context.Context) int {
	count, err := stcq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (stcq *SurveyTemplateCategoryQuery) Exist(ctx context.Context) (bool, error) {
	if err := stcq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return stcq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (stcq *SurveyTemplateCategoryQuery) ExistX(ctx context.Context) bool {
	exist, err := stcq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (stcq *SurveyTemplateCategoryQuery) Clone() *SurveyTemplateCategoryQuery {
	return &SurveyTemplateCategoryQuery{
		config:     stcq.config,
		limit:      stcq.limit,
		offset:     stcq.offset,
		order:      append([]Order{}, stcq.order...),
		unique:     append([]string{}, stcq.unique...),
		predicates: append([]predicate.SurveyTemplateCategory{}, stcq.predicates...),
		// clone intermediate query.
		sql:  stcq.sql.Clone(),
		path: stcq.path,
	}
}

//  WithSurveyTemplateQuestions tells the query-builder to eager-loads the nodes that are connected to
// the "survey_template_questions" edge. The optional arguments used to configure the query builder of the edge.
func (stcq *SurveyTemplateCategoryQuery) WithSurveyTemplateQuestions(opts ...func(*SurveyTemplateQuestionQuery)) *SurveyTemplateCategoryQuery {
	query := &SurveyTemplateQuestionQuery{config: stcq.config}
	for _, opt := range opts {
		opt(query)
	}
	stcq.withSurveyTemplateQuestions = query
	return stcq
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
//	client.SurveyTemplateCategory.Query().
//		GroupBy(surveytemplatecategory.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (stcq *SurveyTemplateCategoryQuery) GroupBy(field string, fields ...string) *SurveyTemplateCategoryGroupBy {
	group := &SurveyTemplateCategoryGroupBy{config: stcq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := stcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return stcq.sqlQuery(), nil
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
//	client.SurveyTemplateCategory.Query().
//		Select(surveytemplatecategory.FieldCreateTime).
//		Scan(ctx, &v)
//
func (stcq *SurveyTemplateCategoryQuery) Select(field string, fields ...string) *SurveyTemplateCategorySelect {
	selector := &SurveyTemplateCategorySelect{config: stcq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := stcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return stcq.sqlQuery(), nil
	}
	return selector
}

func (stcq *SurveyTemplateCategoryQuery) prepareQuery(ctx context.Context) error {
	if stcq.path != nil {
		prev, err := stcq.path(ctx)
		if err != nil {
			return err
		}
		stcq.sql = prev
	}
	if err := surveytemplatecategory.Policy.EvalQuery(ctx, stcq); err != nil {
		return err
	}
	return nil
}

func (stcq *SurveyTemplateCategoryQuery) sqlAll(ctx context.Context) ([]*SurveyTemplateCategory, error) {
	var (
		nodes       = []*SurveyTemplateCategory{}
		withFKs     = stcq.withFKs
		_spec       = stcq.querySpec()
		loadedTypes = [1]bool{
			stcq.withSurveyTemplateQuestions != nil,
		}
	)
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, surveytemplatecategory.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &SurveyTemplateCategory{config: stcq.config}
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
	if err := sqlgraph.QueryNodes(ctx, stcq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := stcq.withSurveyTemplateQuestions; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*SurveyTemplateCategory)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.SurveyTemplateQuestion(func(s *sql.Selector) {
			s.Where(sql.InValues(surveytemplatecategory.SurveyTemplateQuestionsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.survey_template_category_survey_template_questions
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "survey_template_category_survey_template_questions" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "survey_template_category_survey_template_questions" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.SurveyTemplateQuestions = append(node.Edges.SurveyTemplateQuestions, n)
		}
	}

	return nodes, nil
}

func (stcq *SurveyTemplateCategoryQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := stcq.querySpec()
	return sqlgraph.CountNodes(ctx, stcq.driver, _spec)
}

func (stcq *SurveyTemplateCategoryQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := stcq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (stcq *SurveyTemplateCategoryQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveytemplatecategory.Table,
			Columns: surveytemplatecategory.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: surveytemplatecategory.FieldID,
			},
		},
		From:   stcq.sql,
		Unique: true,
	}
	if ps := stcq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := stcq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := stcq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := stcq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (stcq *SurveyTemplateCategoryQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(stcq.driver.Dialect())
	t1 := builder.Table(surveytemplatecategory.Table)
	selector := builder.Select(t1.Columns(surveytemplatecategory.Columns...)...).From(t1)
	if stcq.sql != nil {
		selector = stcq.sql
		selector.Select(selector.Columns(surveytemplatecategory.Columns...)...)
	}
	for _, p := range stcq.predicates {
		p(selector)
	}
	for _, p := range stcq.order {
		p(selector)
	}
	if offset := stcq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := stcq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// SurveyTemplateCategoryGroupBy is the builder for group-by SurveyTemplateCategory entities.
type SurveyTemplateCategoryGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (stcgb *SurveyTemplateCategoryGroupBy) Aggregate(fns ...Aggregate) *SurveyTemplateCategoryGroupBy {
	stcgb.fns = append(stcgb.fns, fns...)
	return stcgb
}

// Scan applies the group-by query and scan the result into the given value.
func (stcgb *SurveyTemplateCategoryGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := stcgb.path(ctx)
	if err != nil {
		return err
	}
	stcgb.sql = query
	return stcgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (stcgb *SurveyTemplateCategoryGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := stcgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (stcgb *SurveyTemplateCategoryGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(stcgb.fields) > 1 {
		return nil, errors.New("ent: SurveyTemplateCategoryGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := stcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (stcgb *SurveyTemplateCategoryGroupBy) StringsX(ctx context.Context) []string {
	v, err := stcgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (stcgb *SurveyTemplateCategoryGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(stcgb.fields) > 1 {
		return nil, errors.New("ent: SurveyTemplateCategoryGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := stcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (stcgb *SurveyTemplateCategoryGroupBy) IntsX(ctx context.Context) []int {
	v, err := stcgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (stcgb *SurveyTemplateCategoryGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(stcgb.fields) > 1 {
		return nil, errors.New("ent: SurveyTemplateCategoryGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := stcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (stcgb *SurveyTemplateCategoryGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := stcgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (stcgb *SurveyTemplateCategoryGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(stcgb.fields) > 1 {
		return nil, errors.New("ent: SurveyTemplateCategoryGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := stcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (stcgb *SurveyTemplateCategoryGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := stcgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (stcgb *SurveyTemplateCategoryGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := stcgb.sqlQuery().Query()
	if err := stcgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (stcgb *SurveyTemplateCategoryGroupBy) sqlQuery() *sql.Selector {
	selector := stcgb.sql
	columns := make([]string, 0, len(stcgb.fields)+len(stcgb.fns))
	columns = append(columns, stcgb.fields...)
	for _, fn := range stcgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(stcgb.fields...)
}

// SurveyTemplateCategorySelect is the builder for select fields of SurveyTemplateCategory entities.
type SurveyTemplateCategorySelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (stcs *SurveyTemplateCategorySelect) Scan(ctx context.Context, v interface{}) error {
	query, err := stcs.path(ctx)
	if err != nil {
		return err
	}
	stcs.sql = query
	return stcs.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (stcs *SurveyTemplateCategorySelect) ScanX(ctx context.Context, v interface{}) {
	if err := stcs.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (stcs *SurveyTemplateCategorySelect) Strings(ctx context.Context) ([]string, error) {
	if len(stcs.fields) > 1 {
		return nil, errors.New("ent: SurveyTemplateCategorySelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := stcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (stcs *SurveyTemplateCategorySelect) StringsX(ctx context.Context) []string {
	v, err := stcs.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (stcs *SurveyTemplateCategorySelect) Ints(ctx context.Context) ([]int, error) {
	if len(stcs.fields) > 1 {
		return nil, errors.New("ent: SurveyTemplateCategorySelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := stcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (stcs *SurveyTemplateCategorySelect) IntsX(ctx context.Context) []int {
	v, err := stcs.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (stcs *SurveyTemplateCategorySelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(stcs.fields) > 1 {
		return nil, errors.New("ent: SurveyTemplateCategorySelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := stcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (stcs *SurveyTemplateCategorySelect) Float64sX(ctx context.Context) []float64 {
	v, err := stcs.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (stcs *SurveyTemplateCategorySelect) Bools(ctx context.Context) ([]bool, error) {
	if len(stcs.fields) > 1 {
		return nil, errors.New("ent: SurveyTemplateCategorySelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := stcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (stcs *SurveyTemplateCategorySelect) BoolsX(ctx context.Context) []bool {
	v, err := stcs.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (stcs *SurveyTemplateCategorySelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := stcs.sqlQuery().Query()
	if err := stcs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (stcs *SurveyTemplateCategorySelect) sqlQuery() sql.Querier {
	selector := stcs.sql
	selector.Select(selector.Columns(stcs.fields...)...)
	return selector
}
