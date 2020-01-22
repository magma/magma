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
	"strconv"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatecategory"
)

// LocationTypeQuery is the builder for querying LocationType entities.
type LocationTypeQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.LocationType
	// eager-loading edges.
	withLocations                *LocationQuery
	withPropertyTypes            *PropertyTypeQuery
	withSurveyTemplateCategories *SurveyTemplateCategoryQuery
	// intermediate query.
	sql *sql.Selector
}

// Where adds a new predicate for the builder.
func (ltq *LocationTypeQuery) Where(ps ...predicate.LocationType) *LocationTypeQuery {
	ltq.predicates = append(ltq.predicates, ps...)
	return ltq
}

// Limit adds a limit step to the query.
func (ltq *LocationTypeQuery) Limit(limit int) *LocationTypeQuery {
	ltq.limit = &limit
	return ltq
}

// Offset adds an offset step to the query.
func (ltq *LocationTypeQuery) Offset(offset int) *LocationTypeQuery {
	ltq.offset = &offset
	return ltq
}

// Order adds an order step to the query.
func (ltq *LocationTypeQuery) Order(o ...Order) *LocationTypeQuery {
	ltq.order = append(ltq.order, o...)
	return ltq
}

// QueryLocations chains the current query on the locations edge.
func (ltq *LocationTypeQuery) QueryLocations() *LocationQuery {
	query := &LocationQuery{config: ltq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(locationtype.Table, locationtype.FieldID, ltq.sqlQuery()),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, locationtype.LocationsTable, locationtype.LocationsColumn),
	)
	query.sql = sqlgraph.SetNeighbors(ltq.driver.Dialect(), step)
	return query
}

// QueryPropertyTypes chains the current query on the property_types edge.
func (ltq *LocationTypeQuery) QueryPropertyTypes() *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: ltq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(locationtype.Table, locationtype.FieldID, ltq.sqlQuery()),
		sqlgraph.To(propertytype.Table, propertytype.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, locationtype.PropertyTypesTable, locationtype.PropertyTypesColumn),
	)
	query.sql = sqlgraph.SetNeighbors(ltq.driver.Dialect(), step)
	return query
}

// QuerySurveyTemplateCategories chains the current query on the survey_template_categories edge.
func (ltq *LocationTypeQuery) QuerySurveyTemplateCategories() *SurveyTemplateCategoryQuery {
	query := &SurveyTemplateCategoryQuery{config: ltq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(locationtype.Table, locationtype.FieldID, ltq.sqlQuery()),
		sqlgraph.To(surveytemplatecategory.Table, surveytemplatecategory.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, locationtype.SurveyTemplateCategoriesTable, locationtype.SurveyTemplateCategoriesColumn),
	)
	query.sql = sqlgraph.SetNeighbors(ltq.driver.Dialect(), step)
	return query
}

// First returns the first LocationType entity in the query. Returns *ErrNotFound when no locationtype was found.
func (ltq *LocationTypeQuery) First(ctx context.Context) (*LocationType, error) {
	lts, err := ltq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(lts) == 0 {
		return nil, &ErrNotFound{locationtype.Label}
	}
	return lts[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (ltq *LocationTypeQuery) FirstX(ctx context.Context) *LocationType {
	lt, err := ltq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return lt
}

// FirstID returns the first LocationType id in the query. Returns *ErrNotFound when no id was found.
func (ltq *LocationTypeQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = ltq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &ErrNotFound{locationtype.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (ltq *LocationTypeQuery) FirstXID(ctx context.Context) string {
	id, err := ltq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only LocationType entity in the query, returns an error if not exactly one entity was returned.
func (ltq *LocationTypeQuery) Only(ctx context.Context) (*LocationType, error) {
	lts, err := ltq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(lts) {
	case 1:
		return lts[0], nil
	case 0:
		return nil, &ErrNotFound{locationtype.Label}
	default:
		return nil, &ErrNotSingular{locationtype.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (ltq *LocationTypeQuery) OnlyX(ctx context.Context) *LocationType {
	lt, err := ltq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return lt
}

// OnlyID returns the only LocationType id in the query, returns an error if not exactly one id was returned.
func (ltq *LocationTypeQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = ltq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &ErrNotFound{locationtype.Label}
	default:
		err = &ErrNotSingular{locationtype.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (ltq *LocationTypeQuery) OnlyXID(ctx context.Context) string {
	id, err := ltq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of LocationTypes.
func (ltq *LocationTypeQuery) All(ctx context.Context) ([]*LocationType, error) {
	return ltq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (ltq *LocationTypeQuery) AllX(ctx context.Context) []*LocationType {
	lts, err := ltq.All(ctx)
	if err != nil {
		panic(err)
	}
	return lts
}

// IDs executes the query and returns a list of LocationType ids.
func (ltq *LocationTypeQuery) IDs(ctx context.Context) ([]string, error) {
	var ids []string
	if err := ltq.Select(locationtype.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (ltq *LocationTypeQuery) IDsX(ctx context.Context) []string {
	ids, err := ltq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (ltq *LocationTypeQuery) Count(ctx context.Context) (int, error) {
	return ltq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (ltq *LocationTypeQuery) CountX(ctx context.Context) int {
	count, err := ltq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (ltq *LocationTypeQuery) Exist(ctx context.Context) (bool, error) {
	return ltq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (ltq *LocationTypeQuery) ExistX(ctx context.Context) bool {
	exist, err := ltq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (ltq *LocationTypeQuery) Clone() *LocationTypeQuery {
	return &LocationTypeQuery{
		config:     ltq.config,
		limit:      ltq.limit,
		offset:     ltq.offset,
		order:      append([]Order{}, ltq.order...),
		unique:     append([]string{}, ltq.unique...),
		predicates: append([]predicate.LocationType{}, ltq.predicates...),
		// clone intermediate query.
		sql: ltq.sql.Clone(),
	}
}

//  WithLocations tells the query-builder to eager-loads the nodes that are connected to
// the "locations" edge. The optional arguments used to configure the query builder of the edge.
func (ltq *LocationTypeQuery) WithLocations(opts ...func(*LocationQuery)) *LocationTypeQuery {
	query := &LocationQuery{config: ltq.config}
	for _, opt := range opts {
		opt(query)
	}
	ltq.withLocations = query
	return ltq
}

//  WithPropertyTypes tells the query-builder to eager-loads the nodes that are connected to
// the "property_types" edge. The optional arguments used to configure the query builder of the edge.
func (ltq *LocationTypeQuery) WithPropertyTypes(opts ...func(*PropertyTypeQuery)) *LocationTypeQuery {
	query := &PropertyTypeQuery{config: ltq.config}
	for _, opt := range opts {
		opt(query)
	}
	ltq.withPropertyTypes = query
	return ltq
}

//  WithSurveyTemplateCategories tells the query-builder to eager-loads the nodes that are connected to
// the "survey_template_categories" edge. The optional arguments used to configure the query builder of the edge.
func (ltq *LocationTypeQuery) WithSurveyTemplateCategories(opts ...func(*SurveyTemplateCategoryQuery)) *LocationTypeQuery {
	query := &SurveyTemplateCategoryQuery{config: ltq.config}
	for _, opt := range opts {
		opt(query)
	}
	ltq.withSurveyTemplateCategories = query
	return ltq
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
//	client.LocationType.Query().
//		GroupBy(locationtype.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (ltq *LocationTypeQuery) GroupBy(field string, fields ...string) *LocationTypeGroupBy {
	group := &LocationTypeGroupBy{config: ltq.config}
	group.fields = append([]string{field}, fields...)
	group.sql = ltq.sqlQuery()
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
//	client.LocationType.Query().
//		Select(locationtype.FieldCreateTime).
//		Scan(ctx, &v)
//
func (ltq *LocationTypeQuery) Select(field string, fields ...string) *LocationTypeSelect {
	selector := &LocationTypeSelect{config: ltq.config}
	selector.fields = append([]string{field}, fields...)
	selector.sql = ltq.sqlQuery()
	return selector
}

func (ltq *LocationTypeQuery) sqlAll(ctx context.Context) ([]*LocationType, error) {
	var (
		nodes []*LocationType
		_spec = ltq.querySpec()
	)
	_spec.ScanValues = func() []interface{} {
		node := &LocationType{config: ltq.config}
		nodes = append(nodes, node)
		values := node.scanValues()
		return values
	}
	_spec.Assign = func(values ...interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("ent: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		return node.assignValues(values...)
	}
	if err := sqlgraph.QueryNodes(ctx, ltq.driver, _spec); err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := ltq.withLocations; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[string]*LocationType)
		for i := range nodes {
			id, err := strconv.Atoi(nodes[i].ID)
			if err != nil {
				return nil, err
			}
			fks = append(fks, id)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.Location(func(s *sql.Selector) {
			s.Where(sql.InValues(locationtype.LocationsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.type_id
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "type_id" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "type_id" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Locations = append(node.Edges.Locations, n)
		}
	}

	if query := ltq.withPropertyTypes; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[string]*LocationType)
		for i := range nodes {
			id, err := strconv.Atoi(nodes[i].ID)
			if err != nil {
				return nil, err
			}
			fks = append(fks, id)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.PropertyType(func(s *sql.Selector) {
			s.Where(sql.InValues(locationtype.PropertyTypesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.location_type_id
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "location_type_id" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "location_type_id" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.PropertyTypes = append(node.Edges.PropertyTypes, n)
		}
	}

	if query := ltq.withSurveyTemplateCategories; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[string]*LocationType)
		for i := range nodes {
			id, err := strconv.Atoi(nodes[i].ID)
			if err != nil {
				return nil, err
			}
			fks = append(fks, id)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.SurveyTemplateCategory(func(s *sql.Selector) {
			s.Where(sql.InValues(locationtype.SurveyTemplateCategoriesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.location_type_survey_template_category_id
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "location_type_survey_template_category_id" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "location_type_survey_template_category_id" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.SurveyTemplateCategories = append(node.Edges.SurveyTemplateCategories, n)
		}
	}

	return nodes, nil
}

func (ltq *LocationTypeQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := ltq.querySpec()
	return sqlgraph.CountNodes(ctx, ltq.driver, _spec)
}

func (ltq *LocationTypeQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := ltq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (ltq *LocationTypeQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   locationtype.Table,
			Columns: locationtype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: locationtype.FieldID,
			},
		},
		From:   ltq.sql,
		Unique: true,
	}
	if ps := ltq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := ltq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := ltq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := ltq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (ltq *LocationTypeQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(ltq.driver.Dialect())
	t1 := builder.Table(locationtype.Table)
	selector := builder.Select(t1.Columns(locationtype.Columns...)...).From(t1)
	if ltq.sql != nil {
		selector = ltq.sql
		selector.Select(selector.Columns(locationtype.Columns...)...)
	}
	for _, p := range ltq.predicates {
		p(selector)
	}
	for _, p := range ltq.order {
		p(selector)
	}
	if offset := ltq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := ltq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// LocationTypeGroupBy is the builder for group-by LocationType entities.
type LocationTypeGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query.
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ltgb *LocationTypeGroupBy) Aggregate(fns ...Aggregate) *LocationTypeGroupBy {
	ltgb.fns = append(ltgb.fns, fns...)
	return ltgb
}

// Scan applies the group-by query and scan the result into the given value.
func (ltgb *LocationTypeGroupBy) Scan(ctx context.Context, v interface{}) error {
	return ltgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ltgb *LocationTypeGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := ltgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (ltgb *LocationTypeGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(ltgb.fields) > 1 {
		return nil, errors.New("ent: LocationTypeGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := ltgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ltgb *LocationTypeGroupBy) StringsX(ctx context.Context) []string {
	v, err := ltgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (ltgb *LocationTypeGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(ltgb.fields) > 1 {
		return nil, errors.New("ent: LocationTypeGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := ltgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ltgb *LocationTypeGroupBy) IntsX(ctx context.Context) []int {
	v, err := ltgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (ltgb *LocationTypeGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(ltgb.fields) > 1 {
		return nil, errors.New("ent: LocationTypeGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := ltgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ltgb *LocationTypeGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := ltgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (ltgb *LocationTypeGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(ltgb.fields) > 1 {
		return nil, errors.New("ent: LocationTypeGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := ltgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ltgb *LocationTypeGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := ltgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ltgb *LocationTypeGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ltgb.sqlQuery().Query()
	if err := ltgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ltgb *LocationTypeGroupBy) sqlQuery() *sql.Selector {
	selector := ltgb.sql
	columns := make([]string, 0, len(ltgb.fields)+len(ltgb.fns))
	columns = append(columns, ltgb.fields...)
	for _, fn := range ltgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(ltgb.fields...)
}

// LocationTypeSelect is the builder for select fields of LocationType entities.
type LocationTypeSelect struct {
	config
	fields []string
	// intermediate queries.
	sql *sql.Selector
}

// Scan applies the selector query and scan the result into the given value.
func (lts *LocationTypeSelect) Scan(ctx context.Context, v interface{}) error {
	return lts.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (lts *LocationTypeSelect) ScanX(ctx context.Context, v interface{}) {
	if err := lts.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (lts *LocationTypeSelect) Strings(ctx context.Context) ([]string, error) {
	if len(lts.fields) > 1 {
		return nil, errors.New("ent: LocationTypeSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := lts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (lts *LocationTypeSelect) StringsX(ctx context.Context) []string {
	v, err := lts.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (lts *LocationTypeSelect) Ints(ctx context.Context) ([]int, error) {
	if len(lts.fields) > 1 {
		return nil, errors.New("ent: LocationTypeSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := lts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (lts *LocationTypeSelect) IntsX(ctx context.Context) []int {
	v, err := lts.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (lts *LocationTypeSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(lts.fields) > 1 {
		return nil, errors.New("ent: LocationTypeSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := lts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (lts *LocationTypeSelect) Float64sX(ctx context.Context) []float64 {
	v, err := lts.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (lts *LocationTypeSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(lts.fields) > 1 {
		return nil, errors.New("ent: LocationTypeSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := lts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (lts *LocationTypeSelect) BoolsX(ctx context.Context) []bool {
	v, err := lts.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (lts *LocationTypeSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := lts.sqlQuery().Query()
	if err := lts.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (lts *LocationTypeSelect) sqlQuery() sql.Querier {
	selector := lts.sql
	selector.Select(selector.Columns(lts.fields...)...)
	return selector
}
