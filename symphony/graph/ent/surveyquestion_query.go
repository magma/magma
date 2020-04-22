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
	// eager-loading edges.
	withSurvey    *SurveyQuery
	withWifiScan  *SurveyWiFiScanQuery
	withCellScan  *SurveyCellScanQuery
	withPhotoData *FileQuery
	withImages    *FileQuery
	withFKs       bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
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
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := sqq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, sqq.sqlQuery()),
			sqlgraph.To(survey.Table, survey.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, surveyquestion.SurveyTable, surveyquestion.SurveyColumn),
		)
		fromU = sqlgraph.SetNeighbors(sqq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryWifiScan chains the current query on the wifi_scan edge.
func (sqq *SurveyQuestionQuery) QueryWifiScan() *SurveyWiFiScanQuery {
	query := &SurveyWiFiScanQuery{config: sqq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := sqq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, sqq.sqlQuery()),
			sqlgraph.To(surveywifiscan.Table, surveywifiscan.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, surveyquestion.WifiScanTable, surveyquestion.WifiScanColumn),
		)
		fromU = sqlgraph.SetNeighbors(sqq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryCellScan chains the current query on the cell_scan edge.
func (sqq *SurveyQuestionQuery) QueryCellScan() *SurveyCellScanQuery {
	query := &SurveyCellScanQuery{config: sqq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := sqq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, sqq.sqlQuery()),
			sqlgraph.To(surveycellscan.Table, surveycellscan.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, surveyquestion.CellScanTable, surveyquestion.CellScanColumn),
		)
		fromU = sqlgraph.SetNeighbors(sqq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryPhotoData chains the current query on the photo_data edge.
func (sqq *SurveyQuestionQuery) QueryPhotoData() *FileQuery {
	query := &FileQuery{config: sqq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := sqq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, sqq.sqlQuery()),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, surveyquestion.PhotoDataTable, surveyquestion.PhotoDataColumn),
		)
		fromU = sqlgraph.SetNeighbors(sqq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryImages chains the current query on the images edge.
func (sqq *SurveyQuestionQuery) QueryImages() *FileQuery {
	query := &FileQuery{config: sqq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := sqq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, sqq.sqlQuery()),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, surveyquestion.ImagesTable, surveyquestion.ImagesColumn),
		)
		fromU = sqlgraph.SetNeighbors(sqq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first SurveyQuestion entity in the query. Returns *NotFoundError when no surveyquestion was found.
func (sqq *SurveyQuestionQuery) First(ctx context.Context) (*SurveyQuestion, error) {
	sqs, err := sqq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(sqs) == 0 {
		return nil, &NotFoundError{surveyquestion.Label}
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

// FirstID returns the first SurveyQuestion id in the query. Returns *NotFoundError when no id was found.
func (sqq *SurveyQuestionQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = sqq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{surveyquestion.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (sqq *SurveyQuestionQuery) FirstXID(ctx context.Context) int {
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
		return nil, &NotFoundError{surveyquestion.Label}
	default:
		return nil, &NotSingularError{surveyquestion.Label}
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
func (sqq *SurveyQuestionQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = sqq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{surveyquestion.Label}
	default:
		err = &NotSingularError{surveyquestion.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (sqq *SurveyQuestionQuery) OnlyXID(ctx context.Context) int {
	id, err := sqq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of SurveyQuestions.
func (sqq *SurveyQuestionQuery) All(ctx context.Context) ([]*SurveyQuestion, error) {
	if err := sqq.prepareQuery(ctx); err != nil {
		return nil, err
	}
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
func (sqq *SurveyQuestionQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := sqq.Select(surveyquestion.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (sqq *SurveyQuestionQuery) IDsX(ctx context.Context) []int {
	ids, err := sqq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (sqq *SurveyQuestionQuery) Count(ctx context.Context) (int, error) {
	if err := sqq.prepareQuery(ctx); err != nil {
		return 0, err
	}
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
	if err := sqq.prepareQuery(ctx); err != nil {
		return false, err
	}
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
		sql:  sqq.sql.Clone(),
		path: sqq.path,
	}
}

//  WithSurvey tells the query-builder to eager-loads the nodes that are connected to
// the "survey" edge. The optional arguments used to configure the query builder of the edge.
func (sqq *SurveyQuestionQuery) WithSurvey(opts ...func(*SurveyQuery)) *SurveyQuestionQuery {
	query := &SurveyQuery{config: sqq.config}
	for _, opt := range opts {
		opt(query)
	}
	sqq.withSurvey = query
	return sqq
}

//  WithWifiScan tells the query-builder to eager-loads the nodes that are connected to
// the "wifi_scan" edge. The optional arguments used to configure the query builder of the edge.
func (sqq *SurveyQuestionQuery) WithWifiScan(opts ...func(*SurveyWiFiScanQuery)) *SurveyQuestionQuery {
	query := &SurveyWiFiScanQuery{config: sqq.config}
	for _, opt := range opts {
		opt(query)
	}
	sqq.withWifiScan = query
	return sqq
}

//  WithCellScan tells the query-builder to eager-loads the nodes that are connected to
// the "cell_scan" edge. The optional arguments used to configure the query builder of the edge.
func (sqq *SurveyQuestionQuery) WithCellScan(opts ...func(*SurveyCellScanQuery)) *SurveyQuestionQuery {
	query := &SurveyCellScanQuery{config: sqq.config}
	for _, opt := range opts {
		opt(query)
	}
	sqq.withCellScan = query
	return sqq
}

//  WithPhotoData tells the query-builder to eager-loads the nodes that are connected to
// the "photo_data" edge. The optional arguments used to configure the query builder of the edge.
func (sqq *SurveyQuestionQuery) WithPhotoData(opts ...func(*FileQuery)) *SurveyQuestionQuery {
	query := &FileQuery{config: sqq.config}
	for _, opt := range opts {
		opt(query)
	}
	sqq.withPhotoData = query
	return sqq
}

//  WithImages tells the query-builder to eager-loads the nodes that are connected to
// the "images" edge. The optional arguments used to configure the query builder of the edge.
func (sqq *SurveyQuestionQuery) WithImages(opts ...func(*FileQuery)) *SurveyQuestionQuery {
	query := &FileQuery{config: sqq.config}
	for _, opt := range opts {
		opt(query)
	}
	sqq.withImages = query
	return sqq
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
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := sqq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return sqq.sqlQuery(), nil
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
//	client.SurveyQuestion.Query().
//		Select(surveyquestion.FieldCreateTime).
//		Scan(ctx, &v)
//
func (sqq *SurveyQuestionQuery) Select(field string, fields ...string) *SurveyQuestionSelect {
	selector := &SurveyQuestionSelect{config: sqq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := sqq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return sqq.sqlQuery(), nil
	}
	return selector
}

func (sqq *SurveyQuestionQuery) prepareQuery(ctx context.Context) error {
	if sqq.path != nil {
		prev, err := sqq.path(ctx)
		if err != nil {
			return err
		}
		sqq.sql = prev
	}
	if err := surveyquestion.Policy.EvalQuery(ctx, sqq); err != nil {
		return err
	}
	return nil
}

func (sqq *SurveyQuestionQuery) sqlAll(ctx context.Context) ([]*SurveyQuestion, error) {
	var (
		nodes       = []*SurveyQuestion{}
		withFKs     = sqq.withFKs
		_spec       = sqq.querySpec()
		loadedTypes = [5]bool{
			sqq.withSurvey != nil,
			sqq.withWifiScan != nil,
			sqq.withCellScan != nil,
			sqq.withPhotoData != nil,
			sqq.withImages != nil,
		}
	)
	if sqq.withSurvey != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, surveyquestion.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &SurveyQuestion{config: sqq.config}
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
	if err := sqlgraph.QueryNodes(ctx, sqq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := sqq.withSurvey; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*SurveyQuestion)
		for i := range nodes {
			if fk := nodes[i].survey_question_survey; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(survey.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "survey_question_survey" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Survey = n
			}
		}
	}

	if query := sqq.withWifiScan; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*SurveyQuestion)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.SurveyWiFiScan(func(s *sql.Selector) {
			s.Where(sql.InValues(surveyquestion.WifiScanColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.survey_wi_fi_scan_survey_question
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "survey_wi_fi_scan_survey_question" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "survey_wi_fi_scan_survey_question" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.WifiScan = append(node.Edges.WifiScan, n)
		}
	}

	if query := sqq.withCellScan; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*SurveyQuestion)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.SurveyCellScan(func(s *sql.Selector) {
			s.Where(sql.InValues(surveyquestion.CellScanColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.survey_cell_scan_survey_question
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "survey_cell_scan_survey_question" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "survey_cell_scan_survey_question" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.CellScan = append(node.Edges.CellScan, n)
		}
	}

	if query := sqq.withPhotoData; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*SurveyQuestion)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.File(func(s *sql.Selector) {
			s.Where(sql.InValues(surveyquestion.PhotoDataColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.survey_question_photo_data
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "survey_question_photo_data" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "survey_question_photo_data" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.PhotoData = append(node.Edges.PhotoData, n)
		}
	}

	if query := sqq.withImages; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*SurveyQuestion)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.File(func(s *sql.Selector) {
			s.Where(sql.InValues(surveyquestion.ImagesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.survey_question_images
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "survey_question_images" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "survey_question_images" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Images = append(node.Edges.Images, n)
		}
	}

	return nodes, nil
}

func (sqq *SurveyQuestionQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := sqq.querySpec()
	return sqlgraph.CountNodes(ctx, sqq.driver, _spec)
}

func (sqq *SurveyQuestionQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := sqq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (sqq *SurveyQuestionQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveyquestion.Table,
			Columns: surveyquestion.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: surveyquestion.FieldID,
			},
		},
		From:   sqq.sql,
		Unique: true,
	}
	if ps := sqq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := sqq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := sqq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := sqq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
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
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (sqgb *SurveyQuestionGroupBy) Aggregate(fns ...Aggregate) *SurveyQuestionGroupBy {
	sqgb.fns = append(sqgb.fns, fns...)
	return sqgb
}

// Scan applies the group-by query and scan the result into the given value.
func (sqgb *SurveyQuestionGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := sqgb.path(ctx)
	if err != nil {
		return err
	}
	sqgb.sql = query
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
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (sqs *SurveyQuestionSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := sqs.path(ctx)
	if err != nil {
		return err
	}
	sqs.sql = query
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
