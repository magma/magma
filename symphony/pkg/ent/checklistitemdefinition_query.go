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
	"github.com/facebookincubator/symphony/pkg/ent/checklistcategorydefinition"
	"github.com/facebookincubator/symphony/pkg/ent/checklistitemdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
)

// CheckListItemDefinitionQuery is the builder for querying CheckListItemDefinition entities.
type CheckListItemDefinitionQuery struct {
	config
	limit      *int
	offset     *int
	order      []OrderFunc
	unique     []string
	predicates []predicate.CheckListItemDefinition
	// eager-loading edges.
	withCheckListCategoryDefinition *CheckListCategoryDefinitionQuery
	withFKs                         bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (clidq *CheckListItemDefinitionQuery) Where(ps ...predicate.CheckListItemDefinition) *CheckListItemDefinitionQuery {
	clidq.predicates = append(clidq.predicates, ps...)
	return clidq
}

// Limit adds a limit step to the query.
func (clidq *CheckListItemDefinitionQuery) Limit(limit int) *CheckListItemDefinitionQuery {
	clidq.limit = &limit
	return clidq
}

// Offset adds an offset step to the query.
func (clidq *CheckListItemDefinitionQuery) Offset(offset int) *CheckListItemDefinitionQuery {
	clidq.offset = &offset
	return clidq
}

// Order adds an order step to the query.
func (clidq *CheckListItemDefinitionQuery) Order(o ...OrderFunc) *CheckListItemDefinitionQuery {
	clidq.order = append(clidq.order, o...)
	return clidq
}

// QueryCheckListCategoryDefinition chains the current query on the check_list_category_definition edge.
func (clidq *CheckListItemDefinitionQuery) QueryCheckListCategoryDefinition() *CheckListCategoryDefinitionQuery {
	query := &CheckListCategoryDefinitionQuery{config: clidq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := clidq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(checklistitemdefinition.Table, checklistitemdefinition.FieldID, clidq.sqlQuery()),
			sqlgraph.To(checklistcategorydefinition.Table, checklistcategorydefinition.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, checklistitemdefinition.CheckListCategoryDefinitionTable, checklistitemdefinition.CheckListCategoryDefinitionColumn),
		)
		fromU = sqlgraph.SetNeighbors(clidq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first CheckListItemDefinition entity in the query. Returns *NotFoundError when no checklistitemdefinition was found.
func (clidq *CheckListItemDefinitionQuery) First(ctx context.Context) (*CheckListItemDefinition, error) {
	clids, err := clidq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(clids) == 0 {
		return nil, &NotFoundError{checklistitemdefinition.Label}
	}
	return clids[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (clidq *CheckListItemDefinitionQuery) FirstX(ctx context.Context) *CheckListItemDefinition {
	clid, err := clidq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return clid
}

// FirstID returns the first CheckListItemDefinition id in the query. Returns *NotFoundError when no id was found.
func (clidq *CheckListItemDefinitionQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = clidq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{checklistitemdefinition.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (clidq *CheckListItemDefinitionQuery) FirstXID(ctx context.Context) int {
	id, err := clidq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only CheckListItemDefinition entity in the query, returns an error if not exactly one entity was returned.
func (clidq *CheckListItemDefinitionQuery) Only(ctx context.Context) (*CheckListItemDefinition, error) {
	clids, err := clidq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(clids) {
	case 1:
		return clids[0], nil
	case 0:
		return nil, &NotFoundError{checklistitemdefinition.Label}
	default:
		return nil, &NotSingularError{checklistitemdefinition.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (clidq *CheckListItemDefinitionQuery) OnlyX(ctx context.Context) *CheckListItemDefinition {
	clid, err := clidq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return clid
}

// OnlyID returns the only CheckListItemDefinition id in the query, returns an error if not exactly one id was returned.
func (clidq *CheckListItemDefinitionQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = clidq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{checklistitemdefinition.Label}
	default:
		err = &NotSingularError{checklistitemdefinition.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (clidq *CheckListItemDefinitionQuery) OnlyXID(ctx context.Context) int {
	id, err := clidq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of CheckListItemDefinitions.
func (clidq *CheckListItemDefinitionQuery) All(ctx context.Context) ([]*CheckListItemDefinition, error) {
	if err := clidq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return clidq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (clidq *CheckListItemDefinitionQuery) AllX(ctx context.Context) []*CheckListItemDefinition {
	clids, err := clidq.All(ctx)
	if err != nil {
		panic(err)
	}
	return clids
}

// IDs executes the query and returns a list of CheckListItemDefinition ids.
func (clidq *CheckListItemDefinitionQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := clidq.Select(checklistitemdefinition.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (clidq *CheckListItemDefinitionQuery) IDsX(ctx context.Context) []int {
	ids, err := clidq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (clidq *CheckListItemDefinitionQuery) Count(ctx context.Context) (int, error) {
	if err := clidq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return clidq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (clidq *CheckListItemDefinitionQuery) CountX(ctx context.Context) int {
	count, err := clidq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (clidq *CheckListItemDefinitionQuery) Exist(ctx context.Context) (bool, error) {
	if err := clidq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return clidq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (clidq *CheckListItemDefinitionQuery) ExistX(ctx context.Context) bool {
	exist, err := clidq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (clidq *CheckListItemDefinitionQuery) Clone() *CheckListItemDefinitionQuery {
	return &CheckListItemDefinitionQuery{
		config:     clidq.config,
		limit:      clidq.limit,
		offset:     clidq.offset,
		order:      append([]OrderFunc{}, clidq.order...),
		unique:     append([]string{}, clidq.unique...),
		predicates: append([]predicate.CheckListItemDefinition{}, clidq.predicates...),
		// clone intermediate query.
		sql:  clidq.sql.Clone(),
		path: clidq.path,
	}
}

//  WithCheckListCategoryDefinition tells the query-builder to eager-loads the nodes that are connected to
// the "check_list_category_definition" edge. The optional arguments used to configure the query builder of the edge.
func (clidq *CheckListItemDefinitionQuery) WithCheckListCategoryDefinition(opts ...func(*CheckListCategoryDefinitionQuery)) *CheckListItemDefinitionQuery {
	query := &CheckListCategoryDefinitionQuery{config: clidq.config}
	for _, opt := range opts {
		opt(query)
	}
	clidq.withCheckListCategoryDefinition = query
	return clidq
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
//	client.CheckListItemDefinition.Query().
//		GroupBy(checklistitemdefinition.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (clidq *CheckListItemDefinitionQuery) GroupBy(field string, fields ...string) *CheckListItemDefinitionGroupBy {
	group := &CheckListItemDefinitionGroupBy{config: clidq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := clidq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return clidq.sqlQuery(), nil
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
//	client.CheckListItemDefinition.Query().
//		Select(checklistitemdefinition.FieldCreateTime).
//		Scan(ctx, &v)
//
func (clidq *CheckListItemDefinitionQuery) Select(field string, fields ...string) *CheckListItemDefinitionSelect {
	selector := &CheckListItemDefinitionSelect{config: clidq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := clidq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return clidq.sqlQuery(), nil
	}
	return selector
}

func (clidq *CheckListItemDefinitionQuery) prepareQuery(ctx context.Context) error {
	if clidq.path != nil {
		prev, err := clidq.path(ctx)
		if err != nil {
			return err
		}
		clidq.sql = prev
	}
	if err := checklistitemdefinition.Policy.EvalQuery(ctx, clidq); err != nil {
		return err
	}
	return nil
}

func (clidq *CheckListItemDefinitionQuery) sqlAll(ctx context.Context) ([]*CheckListItemDefinition, error) {
	var (
		nodes       = []*CheckListItemDefinition{}
		withFKs     = clidq.withFKs
		_spec       = clidq.querySpec()
		loadedTypes = [1]bool{
			clidq.withCheckListCategoryDefinition != nil,
		}
	)
	if clidq.withCheckListCategoryDefinition != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, checklistitemdefinition.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &CheckListItemDefinition{config: clidq.config}
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
	if err := sqlgraph.QueryNodes(ctx, clidq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := clidq.withCheckListCategoryDefinition; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*CheckListItemDefinition)
		for i := range nodes {
			if fk := nodes[i].check_list_category_definition_check_list_item_definitions; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(checklistcategorydefinition.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "check_list_category_definition_check_list_item_definitions" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.CheckListCategoryDefinition = n
			}
		}
	}

	return nodes, nil
}

func (clidq *CheckListItemDefinitionQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := clidq.querySpec()
	return sqlgraph.CountNodes(ctx, clidq.driver, _spec)
}

func (clidq *CheckListItemDefinitionQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := clidq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (clidq *CheckListItemDefinitionQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   checklistitemdefinition.Table,
			Columns: checklistitemdefinition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistitemdefinition.FieldID,
			},
		},
		From:   clidq.sql,
		Unique: true,
	}
	if ps := clidq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := clidq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := clidq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := clidq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (clidq *CheckListItemDefinitionQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(clidq.driver.Dialect())
	t1 := builder.Table(checklistitemdefinition.Table)
	selector := builder.Select(t1.Columns(checklistitemdefinition.Columns...)...).From(t1)
	if clidq.sql != nil {
		selector = clidq.sql
		selector.Select(selector.Columns(checklistitemdefinition.Columns...)...)
	}
	for _, p := range clidq.predicates {
		p(selector)
	}
	for _, p := range clidq.order {
		p(selector)
	}
	if offset := clidq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := clidq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// CheckListItemDefinitionGroupBy is the builder for group-by CheckListItemDefinition entities.
type CheckListItemDefinitionGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (clidgb *CheckListItemDefinitionGroupBy) Aggregate(fns ...AggregateFunc) *CheckListItemDefinitionGroupBy {
	clidgb.fns = append(clidgb.fns, fns...)
	return clidgb
}

// Scan applies the group-by query and scan the result into the given value.
func (clidgb *CheckListItemDefinitionGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := clidgb.path(ctx)
	if err != nil {
		return err
	}
	clidgb.sql = query
	return clidgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (clidgb *CheckListItemDefinitionGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := clidgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (clidgb *CheckListItemDefinitionGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(clidgb.fields) > 1 {
		return nil, errors.New("ent: CheckListItemDefinitionGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := clidgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (clidgb *CheckListItemDefinitionGroupBy) StringsX(ctx context.Context) []string {
	v, err := clidgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (clidgb *CheckListItemDefinitionGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(clidgb.fields) > 1 {
		return nil, errors.New("ent: CheckListItemDefinitionGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := clidgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (clidgb *CheckListItemDefinitionGroupBy) IntsX(ctx context.Context) []int {
	v, err := clidgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (clidgb *CheckListItemDefinitionGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(clidgb.fields) > 1 {
		return nil, errors.New("ent: CheckListItemDefinitionGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := clidgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (clidgb *CheckListItemDefinitionGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := clidgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (clidgb *CheckListItemDefinitionGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(clidgb.fields) > 1 {
		return nil, errors.New("ent: CheckListItemDefinitionGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := clidgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (clidgb *CheckListItemDefinitionGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := clidgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (clidgb *CheckListItemDefinitionGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := clidgb.sqlQuery().Query()
	if err := clidgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (clidgb *CheckListItemDefinitionGroupBy) sqlQuery() *sql.Selector {
	selector := clidgb.sql
	columns := make([]string, 0, len(clidgb.fields)+len(clidgb.fns))
	columns = append(columns, clidgb.fields...)
	for _, fn := range clidgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(clidgb.fields...)
}

// CheckListItemDefinitionSelect is the builder for select fields of CheckListItemDefinition entities.
type CheckListItemDefinitionSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (clids *CheckListItemDefinitionSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := clids.path(ctx)
	if err != nil {
		return err
	}
	clids.sql = query
	return clids.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (clids *CheckListItemDefinitionSelect) ScanX(ctx context.Context, v interface{}) {
	if err := clids.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (clids *CheckListItemDefinitionSelect) Strings(ctx context.Context) ([]string, error) {
	if len(clids.fields) > 1 {
		return nil, errors.New("ent: CheckListItemDefinitionSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := clids.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (clids *CheckListItemDefinitionSelect) StringsX(ctx context.Context) []string {
	v, err := clids.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (clids *CheckListItemDefinitionSelect) Ints(ctx context.Context) ([]int, error) {
	if len(clids.fields) > 1 {
		return nil, errors.New("ent: CheckListItemDefinitionSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := clids.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (clids *CheckListItemDefinitionSelect) IntsX(ctx context.Context) []int {
	v, err := clids.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (clids *CheckListItemDefinitionSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(clids.fields) > 1 {
		return nil, errors.New("ent: CheckListItemDefinitionSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := clids.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (clids *CheckListItemDefinitionSelect) Float64sX(ctx context.Context) []float64 {
	v, err := clids.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (clids *CheckListItemDefinitionSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(clids.fields) > 1 {
		return nil, errors.New("ent: CheckListItemDefinitionSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := clids.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (clids *CheckListItemDefinitionSelect) BoolsX(ctx context.Context) []bool {
	v, err := clids.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (clids *CheckListItemDefinitionSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := clids.sqlQuery().Query()
	if err := clids.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (clids *CheckListItemDefinitionSelect) sqlQuery() sql.Querier {
	selector := clids.sql
	selector.Select(selector.Columns(clids.fields...)...)
	return selector
}
