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
	"github.com/facebookincubator/symphony/graph/ent/equipmentcategory"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// EquipmentCategoryQuery is the builder for querying EquipmentCategory entities.
type EquipmentCategoryQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.EquipmentCategory
	// eager-loading edges.
	withTypes *EquipmentTypeQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (ecq *EquipmentCategoryQuery) Where(ps ...predicate.EquipmentCategory) *EquipmentCategoryQuery {
	ecq.predicates = append(ecq.predicates, ps...)
	return ecq
}

// Limit adds a limit step to the query.
func (ecq *EquipmentCategoryQuery) Limit(limit int) *EquipmentCategoryQuery {
	ecq.limit = &limit
	return ecq
}

// Offset adds an offset step to the query.
func (ecq *EquipmentCategoryQuery) Offset(offset int) *EquipmentCategoryQuery {
	ecq.offset = &offset
	return ecq
}

// Order adds an order step to the query.
func (ecq *EquipmentCategoryQuery) Order(o ...Order) *EquipmentCategoryQuery {
	ecq.order = append(ecq.order, o...)
	return ecq
}

// QueryTypes chains the current query on the types edge.
func (ecq *EquipmentCategoryQuery) QueryTypes() *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: ecq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ecq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentcategory.Table, equipmentcategory.FieldID, ecq.sqlQuery()),
			sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, equipmentcategory.TypesTable, equipmentcategory.TypesColumn),
		)
		fromU = sqlgraph.SetNeighbors(ecq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first EquipmentCategory entity in the query. Returns *NotFoundError when no equipmentcategory was found.
func (ecq *EquipmentCategoryQuery) First(ctx context.Context) (*EquipmentCategory, error) {
	ecs, err := ecq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(ecs) == 0 {
		return nil, &NotFoundError{equipmentcategory.Label}
	}
	return ecs[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (ecq *EquipmentCategoryQuery) FirstX(ctx context.Context) *EquipmentCategory {
	ec, err := ecq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return ec
}

// FirstID returns the first EquipmentCategory id in the query. Returns *NotFoundError when no id was found.
func (ecq *EquipmentCategoryQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ecq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{equipmentcategory.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (ecq *EquipmentCategoryQuery) FirstXID(ctx context.Context) int {
	id, err := ecq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only EquipmentCategory entity in the query, returns an error if not exactly one entity was returned.
func (ecq *EquipmentCategoryQuery) Only(ctx context.Context) (*EquipmentCategory, error) {
	ecs, err := ecq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(ecs) {
	case 1:
		return ecs[0], nil
	case 0:
		return nil, &NotFoundError{equipmentcategory.Label}
	default:
		return nil, &NotSingularError{equipmentcategory.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (ecq *EquipmentCategoryQuery) OnlyX(ctx context.Context) *EquipmentCategory {
	ec, err := ecq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return ec
}

// OnlyID returns the only EquipmentCategory id in the query, returns an error if not exactly one id was returned.
func (ecq *EquipmentCategoryQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ecq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{equipmentcategory.Label}
	default:
		err = &NotSingularError{equipmentcategory.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (ecq *EquipmentCategoryQuery) OnlyXID(ctx context.Context) int {
	id, err := ecq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of EquipmentCategories.
func (ecq *EquipmentCategoryQuery) All(ctx context.Context) ([]*EquipmentCategory, error) {
	if err := ecq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return ecq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (ecq *EquipmentCategoryQuery) AllX(ctx context.Context) []*EquipmentCategory {
	ecs, err := ecq.All(ctx)
	if err != nil {
		panic(err)
	}
	return ecs
}

// IDs executes the query and returns a list of EquipmentCategory ids.
func (ecq *EquipmentCategoryQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := ecq.Select(equipmentcategory.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (ecq *EquipmentCategoryQuery) IDsX(ctx context.Context) []int {
	ids, err := ecq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (ecq *EquipmentCategoryQuery) Count(ctx context.Context) (int, error) {
	if err := ecq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return ecq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (ecq *EquipmentCategoryQuery) CountX(ctx context.Context) int {
	count, err := ecq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (ecq *EquipmentCategoryQuery) Exist(ctx context.Context) (bool, error) {
	if err := ecq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return ecq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (ecq *EquipmentCategoryQuery) ExistX(ctx context.Context) bool {
	exist, err := ecq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (ecq *EquipmentCategoryQuery) Clone() *EquipmentCategoryQuery {
	return &EquipmentCategoryQuery{
		config:     ecq.config,
		limit:      ecq.limit,
		offset:     ecq.offset,
		order:      append([]Order{}, ecq.order...),
		unique:     append([]string{}, ecq.unique...),
		predicates: append([]predicate.EquipmentCategory{}, ecq.predicates...),
		// clone intermediate query.
		sql:  ecq.sql.Clone(),
		path: ecq.path,
	}
}

//  WithTypes tells the query-builder to eager-loads the nodes that are connected to
// the "types" edge. The optional arguments used to configure the query builder of the edge.
func (ecq *EquipmentCategoryQuery) WithTypes(opts ...func(*EquipmentTypeQuery)) *EquipmentCategoryQuery {
	query := &EquipmentTypeQuery{config: ecq.config}
	for _, opt := range opts {
		opt(query)
	}
	ecq.withTypes = query
	return ecq
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
//	client.EquipmentCategory.Query().
//		GroupBy(equipmentcategory.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (ecq *EquipmentCategoryQuery) GroupBy(field string, fields ...string) *EquipmentCategoryGroupBy {
	group := &EquipmentCategoryGroupBy{config: ecq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := ecq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return ecq.sqlQuery(), nil
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
//	client.EquipmentCategory.Query().
//		Select(equipmentcategory.FieldCreateTime).
//		Scan(ctx, &v)
//
func (ecq *EquipmentCategoryQuery) Select(field string, fields ...string) *EquipmentCategorySelect {
	selector := &EquipmentCategorySelect{config: ecq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := ecq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return ecq.sqlQuery(), nil
	}
	return selector
}

func (ecq *EquipmentCategoryQuery) prepareQuery(ctx context.Context) error {
	if ecq.path != nil {
		prev, err := ecq.path(ctx)
		if err != nil {
			return err
		}
		ecq.sql = prev
	}
	return nil
}

func (ecq *EquipmentCategoryQuery) sqlAll(ctx context.Context) ([]*EquipmentCategory, error) {
	var (
		nodes       = []*EquipmentCategory{}
		_spec       = ecq.querySpec()
		loadedTypes = [1]bool{
			ecq.withTypes != nil,
		}
	)
	_spec.ScanValues = func() []interface{} {
		node := &EquipmentCategory{config: ecq.config}
		nodes = append(nodes, node)
		values := node.scanValues()
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
	if err := sqlgraph.QueryNodes(ctx, ecq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := ecq.withTypes; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*EquipmentCategory)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.EquipmentType(func(s *sql.Selector) {
			s.Where(sql.InValues(equipmentcategory.TypesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.equipment_type_category
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "equipment_type_category" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_type_category" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Types = append(node.Edges.Types, n)
		}
	}

	return nodes, nil
}

func (ecq *EquipmentCategoryQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := ecq.querySpec()
	return sqlgraph.CountNodes(ctx, ecq.driver, _spec)
}

func (ecq *EquipmentCategoryQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := ecq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (ecq *EquipmentCategoryQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentcategory.Table,
			Columns: equipmentcategory.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentcategory.FieldID,
			},
		},
		From:   ecq.sql,
		Unique: true,
	}
	if ps := ecq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := ecq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := ecq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := ecq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (ecq *EquipmentCategoryQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(ecq.driver.Dialect())
	t1 := builder.Table(equipmentcategory.Table)
	selector := builder.Select(t1.Columns(equipmentcategory.Columns...)...).From(t1)
	if ecq.sql != nil {
		selector = ecq.sql
		selector.Select(selector.Columns(equipmentcategory.Columns...)...)
	}
	for _, p := range ecq.predicates {
		p(selector)
	}
	for _, p := range ecq.order {
		p(selector)
	}
	if offset := ecq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := ecq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// EquipmentCategoryGroupBy is the builder for group-by EquipmentCategory entities.
type EquipmentCategoryGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ecgb *EquipmentCategoryGroupBy) Aggregate(fns ...Aggregate) *EquipmentCategoryGroupBy {
	ecgb.fns = append(ecgb.fns, fns...)
	return ecgb
}

// Scan applies the group-by query and scan the result into the given value.
func (ecgb *EquipmentCategoryGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := ecgb.path(ctx)
	if err != nil {
		return err
	}
	ecgb.sql = query
	return ecgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ecgb *EquipmentCategoryGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := ecgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (ecgb *EquipmentCategoryGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(ecgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentCategoryGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := ecgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ecgb *EquipmentCategoryGroupBy) StringsX(ctx context.Context) []string {
	v, err := ecgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (ecgb *EquipmentCategoryGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(ecgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentCategoryGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := ecgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ecgb *EquipmentCategoryGroupBy) IntsX(ctx context.Context) []int {
	v, err := ecgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (ecgb *EquipmentCategoryGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(ecgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentCategoryGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := ecgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ecgb *EquipmentCategoryGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := ecgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (ecgb *EquipmentCategoryGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(ecgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentCategoryGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := ecgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ecgb *EquipmentCategoryGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := ecgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ecgb *EquipmentCategoryGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ecgb.sqlQuery().Query()
	if err := ecgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ecgb *EquipmentCategoryGroupBy) sqlQuery() *sql.Selector {
	selector := ecgb.sql
	columns := make([]string, 0, len(ecgb.fields)+len(ecgb.fns))
	columns = append(columns, ecgb.fields...)
	for _, fn := range ecgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(ecgb.fields...)
}

// EquipmentCategorySelect is the builder for select fields of EquipmentCategory entities.
type EquipmentCategorySelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (ecs *EquipmentCategorySelect) Scan(ctx context.Context, v interface{}) error {
	query, err := ecs.path(ctx)
	if err != nil {
		return err
	}
	ecs.sql = query
	return ecs.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ecs *EquipmentCategorySelect) ScanX(ctx context.Context, v interface{}) {
	if err := ecs.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (ecs *EquipmentCategorySelect) Strings(ctx context.Context) ([]string, error) {
	if len(ecs.fields) > 1 {
		return nil, errors.New("ent: EquipmentCategorySelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := ecs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ecs *EquipmentCategorySelect) StringsX(ctx context.Context) []string {
	v, err := ecs.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (ecs *EquipmentCategorySelect) Ints(ctx context.Context) ([]int, error) {
	if len(ecs.fields) > 1 {
		return nil, errors.New("ent: EquipmentCategorySelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := ecs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ecs *EquipmentCategorySelect) IntsX(ctx context.Context) []int {
	v, err := ecs.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (ecs *EquipmentCategorySelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(ecs.fields) > 1 {
		return nil, errors.New("ent: EquipmentCategorySelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := ecs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ecs *EquipmentCategorySelect) Float64sX(ctx context.Context) []float64 {
	v, err := ecs.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (ecs *EquipmentCategorySelect) Bools(ctx context.Context) ([]bool, error) {
	if len(ecs.fields) > 1 {
		return nil, errors.New("ent: EquipmentCategorySelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := ecs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ecs *EquipmentCategorySelect) BoolsX(ctx context.Context) []bool {
	v, err := ecs.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ecs *EquipmentCategorySelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ecs.sqlQuery().Query()
	if err := ecs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ecs *EquipmentCategorySelect) sqlQuery() sql.Querier {
	selector := ecs.sql
	selector.Select(selector.Columns(ecs.fields...)...)
	return selector
}
