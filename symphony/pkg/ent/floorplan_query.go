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
	"github.com/facebookincubator/symphony/pkg/ent/file"
	"github.com/facebookincubator/symphony/pkg/ent/floorplan"
	"github.com/facebookincubator/symphony/pkg/ent/floorplanreferencepoint"
	"github.com/facebookincubator/symphony/pkg/ent/floorplanscale"
	"github.com/facebookincubator/symphony/pkg/ent/location"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
)

// FloorPlanQuery is the builder for querying FloorPlan entities.
type FloorPlanQuery struct {
	config
	limit      *int
	offset     *int
	order      []OrderFunc
	unique     []string
	predicates []predicate.FloorPlan
	// eager-loading edges.
	withLocation       *LocationQuery
	withReferencePoint *FloorPlanReferencePointQuery
	withScale          *FloorPlanScaleQuery
	withImage          *FileQuery
	withFKs            bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (fpq *FloorPlanQuery) Where(ps ...predicate.FloorPlan) *FloorPlanQuery {
	fpq.predicates = append(fpq.predicates, ps...)
	return fpq
}

// Limit adds a limit step to the query.
func (fpq *FloorPlanQuery) Limit(limit int) *FloorPlanQuery {
	fpq.limit = &limit
	return fpq
}

// Offset adds an offset step to the query.
func (fpq *FloorPlanQuery) Offset(offset int) *FloorPlanQuery {
	fpq.offset = &offset
	return fpq
}

// Order adds an order step to the query.
func (fpq *FloorPlanQuery) Order(o ...OrderFunc) *FloorPlanQuery {
	fpq.order = append(fpq.order, o...)
	return fpq
}

// QueryLocation chains the current query on the location edge.
func (fpq *FloorPlanQuery) QueryLocation() *LocationQuery {
	query := &LocationQuery{config: fpq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := fpq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(floorplan.Table, floorplan.FieldID, fpq.sqlQuery()),
			sqlgraph.To(location.Table, location.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, floorplan.LocationTable, floorplan.LocationColumn),
		)
		fromU = sqlgraph.SetNeighbors(fpq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryReferencePoint chains the current query on the reference_point edge.
func (fpq *FloorPlanQuery) QueryReferencePoint() *FloorPlanReferencePointQuery {
	query := &FloorPlanReferencePointQuery{config: fpq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := fpq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(floorplan.Table, floorplan.FieldID, fpq.sqlQuery()),
			sqlgraph.To(floorplanreferencepoint.Table, floorplanreferencepoint.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, floorplan.ReferencePointTable, floorplan.ReferencePointColumn),
		)
		fromU = sqlgraph.SetNeighbors(fpq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryScale chains the current query on the scale edge.
func (fpq *FloorPlanQuery) QueryScale() *FloorPlanScaleQuery {
	query := &FloorPlanScaleQuery{config: fpq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := fpq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(floorplan.Table, floorplan.FieldID, fpq.sqlQuery()),
			sqlgraph.To(floorplanscale.Table, floorplanscale.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, floorplan.ScaleTable, floorplan.ScaleColumn),
		)
		fromU = sqlgraph.SetNeighbors(fpq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryImage chains the current query on the image edge.
func (fpq *FloorPlanQuery) QueryImage() *FileQuery {
	query := &FileQuery{config: fpq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := fpq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(floorplan.Table, floorplan.FieldID, fpq.sqlQuery()),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, floorplan.ImageTable, floorplan.ImageColumn),
		)
		fromU = sqlgraph.SetNeighbors(fpq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first FloorPlan entity in the query. Returns *NotFoundError when no floorplan was found.
func (fpq *FloorPlanQuery) First(ctx context.Context) (*FloorPlan, error) {
	fps, err := fpq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(fps) == 0 {
		return nil, &NotFoundError{floorplan.Label}
	}
	return fps[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (fpq *FloorPlanQuery) FirstX(ctx context.Context) *FloorPlan {
	fp, err := fpq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return fp
}

// FirstID returns the first FloorPlan id in the query. Returns *NotFoundError when no id was found.
func (fpq *FloorPlanQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = fpq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{floorplan.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (fpq *FloorPlanQuery) FirstXID(ctx context.Context) int {
	id, err := fpq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only FloorPlan entity in the query, returns an error if not exactly one entity was returned.
func (fpq *FloorPlanQuery) Only(ctx context.Context) (*FloorPlan, error) {
	fps, err := fpq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(fps) {
	case 1:
		return fps[0], nil
	case 0:
		return nil, &NotFoundError{floorplan.Label}
	default:
		return nil, &NotSingularError{floorplan.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (fpq *FloorPlanQuery) OnlyX(ctx context.Context) *FloorPlan {
	fp, err := fpq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return fp
}

// OnlyID returns the only FloorPlan id in the query, returns an error if not exactly one id was returned.
func (fpq *FloorPlanQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = fpq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{floorplan.Label}
	default:
		err = &NotSingularError{floorplan.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (fpq *FloorPlanQuery) OnlyXID(ctx context.Context) int {
	id, err := fpq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of FloorPlans.
func (fpq *FloorPlanQuery) All(ctx context.Context) ([]*FloorPlan, error) {
	if err := fpq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return fpq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (fpq *FloorPlanQuery) AllX(ctx context.Context) []*FloorPlan {
	fps, err := fpq.All(ctx)
	if err != nil {
		panic(err)
	}
	return fps
}

// IDs executes the query and returns a list of FloorPlan ids.
func (fpq *FloorPlanQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := fpq.Select(floorplan.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (fpq *FloorPlanQuery) IDsX(ctx context.Context) []int {
	ids, err := fpq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (fpq *FloorPlanQuery) Count(ctx context.Context) (int, error) {
	if err := fpq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return fpq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (fpq *FloorPlanQuery) CountX(ctx context.Context) int {
	count, err := fpq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (fpq *FloorPlanQuery) Exist(ctx context.Context) (bool, error) {
	if err := fpq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return fpq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (fpq *FloorPlanQuery) ExistX(ctx context.Context) bool {
	exist, err := fpq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (fpq *FloorPlanQuery) Clone() *FloorPlanQuery {
	return &FloorPlanQuery{
		config:     fpq.config,
		limit:      fpq.limit,
		offset:     fpq.offset,
		order:      append([]OrderFunc{}, fpq.order...),
		unique:     append([]string{}, fpq.unique...),
		predicates: append([]predicate.FloorPlan{}, fpq.predicates...),
		// clone intermediate query.
		sql:  fpq.sql.Clone(),
		path: fpq.path,
	}
}

//  WithLocation tells the query-builder to eager-loads the nodes that are connected to
// the "location" edge. The optional arguments used to configure the query builder of the edge.
func (fpq *FloorPlanQuery) WithLocation(opts ...func(*LocationQuery)) *FloorPlanQuery {
	query := &LocationQuery{config: fpq.config}
	for _, opt := range opts {
		opt(query)
	}
	fpq.withLocation = query
	return fpq
}

//  WithReferencePoint tells the query-builder to eager-loads the nodes that are connected to
// the "reference_point" edge. The optional arguments used to configure the query builder of the edge.
func (fpq *FloorPlanQuery) WithReferencePoint(opts ...func(*FloorPlanReferencePointQuery)) *FloorPlanQuery {
	query := &FloorPlanReferencePointQuery{config: fpq.config}
	for _, opt := range opts {
		opt(query)
	}
	fpq.withReferencePoint = query
	return fpq
}

//  WithScale tells the query-builder to eager-loads the nodes that are connected to
// the "scale" edge. The optional arguments used to configure the query builder of the edge.
func (fpq *FloorPlanQuery) WithScale(opts ...func(*FloorPlanScaleQuery)) *FloorPlanQuery {
	query := &FloorPlanScaleQuery{config: fpq.config}
	for _, opt := range opts {
		opt(query)
	}
	fpq.withScale = query
	return fpq
}

//  WithImage tells the query-builder to eager-loads the nodes that are connected to
// the "image" edge. The optional arguments used to configure the query builder of the edge.
func (fpq *FloorPlanQuery) WithImage(opts ...func(*FileQuery)) *FloorPlanQuery {
	query := &FileQuery{config: fpq.config}
	for _, opt := range opts {
		opt(query)
	}
	fpq.withImage = query
	return fpq
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
//	client.FloorPlan.Query().
//		GroupBy(floorplan.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (fpq *FloorPlanQuery) GroupBy(field string, fields ...string) *FloorPlanGroupBy {
	group := &FloorPlanGroupBy{config: fpq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := fpq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return fpq.sqlQuery(), nil
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
//	client.FloorPlan.Query().
//		Select(floorplan.FieldCreateTime).
//		Scan(ctx, &v)
//
func (fpq *FloorPlanQuery) Select(field string, fields ...string) *FloorPlanSelect {
	selector := &FloorPlanSelect{config: fpq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := fpq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return fpq.sqlQuery(), nil
	}
	return selector
}

func (fpq *FloorPlanQuery) prepareQuery(ctx context.Context) error {
	if fpq.path != nil {
		prev, err := fpq.path(ctx)
		if err != nil {
			return err
		}
		fpq.sql = prev
	}
	if err := floorplan.Policy.EvalQuery(ctx, fpq); err != nil {
		return err
	}
	return nil
}

func (fpq *FloorPlanQuery) sqlAll(ctx context.Context) ([]*FloorPlan, error) {
	var (
		nodes       = []*FloorPlan{}
		withFKs     = fpq.withFKs
		_spec       = fpq.querySpec()
		loadedTypes = [4]bool{
			fpq.withLocation != nil,
			fpq.withReferencePoint != nil,
			fpq.withScale != nil,
			fpq.withImage != nil,
		}
	)
	if fpq.withLocation != nil || fpq.withReferencePoint != nil || fpq.withScale != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, floorplan.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &FloorPlan{config: fpq.config}
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
	if err := sqlgraph.QueryNodes(ctx, fpq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := fpq.withLocation; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*FloorPlan)
		for i := range nodes {
			if fk := nodes[i].floor_plan_location; fk != nil {
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
				return nil, fmt.Errorf(`unexpected foreign-key "floor_plan_location" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Location = n
			}
		}
	}

	if query := fpq.withReferencePoint; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*FloorPlan)
		for i := range nodes {
			if fk := nodes[i].floor_plan_reference_point; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(floorplanreferencepoint.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "floor_plan_reference_point" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.ReferencePoint = n
			}
		}
	}

	if query := fpq.withScale; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*FloorPlan)
		for i := range nodes {
			if fk := nodes[i].floor_plan_scale; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(floorplanscale.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "floor_plan_scale" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Scale = n
			}
		}
	}

	if query := fpq.withImage; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*FloorPlan)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.File(func(s *sql.Selector) {
			s.Where(sql.InValues(floorplan.ImageColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.floor_plan_image
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "floor_plan_image" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "floor_plan_image" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Image = n
		}
	}

	return nodes, nil
}

func (fpq *FloorPlanQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := fpq.querySpec()
	return sqlgraph.CountNodes(ctx, fpq.driver, _spec)
}

func (fpq *FloorPlanQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := fpq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (fpq *FloorPlanQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   floorplan.Table,
			Columns: floorplan.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: floorplan.FieldID,
			},
		},
		From:   fpq.sql,
		Unique: true,
	}
	if ps := fpq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := fpq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := fpq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := fpq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (fpq *FloorPlanQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(fpq.driver.Dialect())
	t1 := builder.Table(floorplan.Table)
	selector := builder.Select(t1.Columns(floorplan.Columns...)...).From(t1)
	if fpq.sql != nil {
		selector = fpq.sql
		selector.Select(selector.Columns(floorplan.Columns...)...)
	}
	for _, p := range fpq.predicates {
		p(selector)
	}
	for _, p := range fpq.order {
		p(selector)
	}
	if offset := fpq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := fpq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// FloorPlanGroupBy is the builder for group-by FloorPlan entities.
type FloorPlanGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (fpgb *FloorPlanGroupBy) Aggregate(fns ...AggregateFunc) *FloorPlanGroupBy {
	fpgb.fns = append(fpgb.fns, fns...)
	return fpgb
}

// Scan applies the group-by query and scan the result into the given value.
func (fpgb *FloorPlanGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := fpgb.path(ctx)
	if err != nil {
		return err
	}
	fpgb.sql = query
	return fpgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (fpgb *FloorPlanGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := fpgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (fpgb *FloorPlanGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(fpgb.fields) > 1 {
		return nil, errors.New("ent: FloorPlanGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := fpgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (fpgb *FloorPlanGroupBy) StringsX(ctx context.Context) []string {
	v, err := fpgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (fpgb *FloorPlanGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(fpgb.fields) > 1 {
		return nil, errors.New("ent: FloorPlanGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := fpgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (fpgb *FloorPlanGroupBy) IntsX(ctx context.Context) []int {
	v, err := fpgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (fpgb *FloorPlanGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(fpgb.fields) > 1 {
		return nil, errors.New("ent: FloorPlanGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := fpgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (fpgb *FloorPlanGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := fpgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (fpgb *FloorPlanGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(fpgb.fields) > 1 {
		return nil, errors.New("ent: FloorPlanGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := fpgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (fpgb *FloorPlanGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := fpgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (fpgb *FloorPlanGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := fpgb.sqlQuery().Query()
	if err := fpgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (fpgb *FloorPlanGroupBy) sqlQuery() *sql.Selector {
	selector := fpgb.sql
	columns := make([]string, 0, len(fpgb.fields)+len(fpgb.fns))
	columns = append(columns, fpgb.fields...)
	for _, fn := range fpgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(fpgb.fields...)
}

// FloorPlanSelect is the builder for select fields of FloorPlan entities.
type FloorPlanSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (fps *FloorPlanSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := fps.path(ctx)
	if err != nil {
		return err
	}
	fps.sql = query
	return fps.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (fps *FloorPlanSelect) ScanX(ctx context.Context, v interface{}) {
	if err := fps.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (fps *FloorPlanSelect) Strings(ctx context.Context) ([]string, error) {
	if len(fps.fields) > 1 {
		return nil, errors.New("ent: FloorPlanSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := fps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (fps *FloorPlanSelect) StringsX(ctx context.Context) []string {
	v, err := fps.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (fps *FloorPlanSelect) Ints(ctx context.Context) ([]int, error) {
	if len(fps.fields) > 1 {
		return nil, errors.New("ent: FloorPlanSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := fps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (fps *FloorPlanSelect) IntsX(ctx context.Context) []int {
	v, err := fps.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (fps *FloorPlanSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(fps.fields) > 1 {
		return nil, errors.New("ent: FloorPlanSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := fps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (fps *FloorPlanSelect) Float64sX(ctx context.Context) []float64 {
	v, err := fps.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (fps *FloorPlanSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(fps.fields) > 1 {
		return nil, errors.New("ent: FloorPlanSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := fps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (fps *FloorPlanSelect) BoolsX(ctx context.Context) []bool {
	v, err := fps.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (fps *FloorPlanSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := fps.sqlQuery().Query()
	if err := fps.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (fps *FloorPlanSelect) sqlQuery() sql.Querier {
	selector := fps.sql
	selector.Select(selector.Columns(fps.fields...)...)
	return selector
}
