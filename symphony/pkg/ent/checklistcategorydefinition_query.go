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
	"github.com/facebookincubator/symphony/pkg/ent/checklistcategorydefinition"
	"github.com/facebookincubator/symphony/pkg/ent/checklistitemdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/workordertype"
)

// CheckListCategoryDefinitionQuery is the builder for querying CheckListCategoryDefinition entities.
type CheckListCategoryDefinitionQuery struct {
	config
	limit      *int
	offset     *int
	order      []OrderFunc
	unique     []string
	predicates []predicate.CheckListCategoryDefinition
	// eager-loading edges.
	withCheckListItemDefinitions *CheckListItemDefinitionQuery
	withWorkOrderType            *WorkOrderTypeQuery
	withFKs                      bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (clcdq *CheckListCategoryDefinitionQuery) Where(ps ...predicate.CheckListCategoryDefinition) *CheckListCategoryDefinitionQuery {
	clcdq.predicates = append(clcdq.predicates, ps...)
	return clcdq
}

// Limit adds a limit step to the query.
func (clcdq *CheckListCategoryDefinitionQuery) Limit(limit int) *CheckListCategoryDefinitionQuery {
	clcdq.limit = &limit
	return clcdq
}

// Offset adds an offset step to the query.
func (clcdq *CheckListCategoryDefinitionQuery) Offset(offset int) *CheckListCategoryDefinitionQuery {
	clcdq.offset = &offset
	return clcdq
}

// Order adds an order step to the query.
func (clcdq *CheckListCategoryDefinitionQuery) Order(o ...OrderFunc) *CheckListCategoryDefinitionQuery {
	clcdq.order = append(clcdq.order, o...)
	return clcdq
}

// QueryCheckListItemDefinitions chains the current query on the check_list_item_definitions edge.
func (clcdq *CheckListCategoryDefinitionQuery) QueryCheckListItemDefinitions() *CheckListItemDefinitionQuery {
	query := &CheckListItemDefinitionQuery{config: clcdq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := clcdq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(checklistcategorydefinition.Table, checklistcategorydefinition.FieldID, clcdq.sqlQuery()),
			sqlgraph.To(checklistitemdefinition.Table, checklistitemdefinition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, checklistcategorydefinition.CheckListItemDefinitionsTable, checklistcategorydefinition.CheckListItemDefinitionsColumn),
		)
		fromU = sqlgraph.SetNeighbors(clcdq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryWorkOrderType chains the current query on the work_order_type edge.
func (clcdq *CheckListCategoryDefinitionQuery) QueryWorkOrderType() *WorkOrderTypeQuery {
	query := &WorkOrderTypeQuery{config: clcdq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := clcdq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(checklistcategorydefinition.Table, checklistcategorydefinition.FieldID, clcdq.sqlQuery()),
			sqlgraph.To(workordertype.Table, workordertype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, checklistcategorydefinition.WorkOrderTypeTable, checklistcategorydefinition.WorkOrderTypeColumn),
		)
		fromU = sqlgraph.SetNeighbors(clcdq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first CheckListCategoryDefinition entity in the query. Returns *NotFoundError when no checklistcategorydefinition was found.
func (clcdq *CheckListCategoryDefinitionQuery) First(ctx context.Context) (*CheckListCategoryDefinition, error) {
	clcds, err := clcdq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(clcds) == 0 {
		return nil, &NotFoundError{checklistcategorydefinition.Label}
	}
	return clcds[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (clcdq *CheckListCategoryDefinitionQuery) FirstX(ctx context.Context) *CheckListCategoryDefinition {
	clcd, err := clcdq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return clcd
}

// FirstID returns the first CheckListCategoryDefinition id in the query. Returns *NotFoundError when no id was found.
func (clcdq *CheckListCategoryDefinitionQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = clcdq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{checklistcategorydefinition.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (clcdq *CheckListCategoryDefinitionQuery) FirstXID(ctx context.Context) int {
	id, err := clcdq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only CheckListCategoryDefinition entity in the query, returns an error if not exactly one entity was returned.
func (clcdq *CheckListCategoryDefinitionQuery) Only(ctx context.Context) (*CheckListCategoryDefinition, error) {
	clcds, err := clcdq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(clcds) {
	case 1:
		return clcds[0], nil
	case 0:
		return nil, &NotFoundError{checklistcategorydefinition.Label}
	default:
		return nil, &NotSingularError{checklistcategorydefinition.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (clcdq *CheckListCategoryDefinitionQuery) OnlyX(ctx context.Context) *CheckListCategoryDefinition {
	clcd, err := clcdq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return clcd
}

// OnlyID returns the only CheckListCategoryDefinition id in the query, returns an error if not exactly one id was returned.
func (clcdq *CheckListCategoryDefinitionQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = clcdq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{checklistcategorydefinition.Label}
	default:
		err = &NotSingularError{checklistcategorydefinition.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (clcdq *CheckListCategoryDefinitionQuery) OnlyXID(ctx context.Context) int {
	id, err := clcdq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of CheckListCategoryDefinitions.
func (clcdq *CheckListCategoryDefinitionQuery) All(ctx context.Context) ([]*CheckListCategoryDefinition, error) {
	if err := clcdq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return clcdq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (clcdq *CheckListCategoryDefinitionQuery) AllX(ctx context.Context) []*CheckListCategoryDefinition {
	clcds, err := clcdq.All(ctx)
	if err != nil {
		panic(err)
	}
	return clcds
}

// IDs executes the query and returns a list of CheckListCategoryDefinition ids.
func (clcdq *CheckListCategoryDefinitionQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := clcdq.Select(checklistcategorydefinition.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (clcdq *CheckListCategoryDefinitionQuery) IDsX(ctx context.Context) []int {
	ids, err := clcdq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (clcdq *CheckListCategoryDefinitionQuery) Count(ctx context.Context) (int, error) {
	if err := clcdq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return clcdq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (clcdq *CheckListCategoryDefinitionQuery) CountX(ctx context.Context) int {
	count, err := clcdq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (clcdq *CheckListCategoryDefinitionQuery) Exist(ctx context.Context) (bool, error) {
	if err := clcdq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return clcdq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (clcdq *CheckListCategoryDefinitionQuery) ExistX(ctx context.Context) bool {
	exist, err := clcdq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (clcdq *CheckListCategoryDefinitionQuery) Clone() *CheckListCategoryDefinitionQuery {
	return &CheckListCategoryDefinitionQuery{
		config:     clcdq.config,
		limit:      clcdq.limit,
		offset:     clcdq.offset,
		order:      append([]OrderFunc{}, clcdq.order...),
		unique:     append([]string{}, clcdq.unique...),
		predicates: append([]predicate.CheckListCategoryDefinition{}, clcdq.predicates...),
		// clone intermediate query.
		sql:  clcdq.sql.Clone(),
		path: clcdq.path,
	}
}

//  WithCheckListItemDefinitions tells the query-builder to eager-loads the nodes that are connected to
// the "check_list_item_definitions" edge. The optional arguments used to configure the query builder of the edge.
func (clcdq *CheckListCategoryDefinitionQuery) WithCheckListItemDefinitions(opts ...func(*CheckListItemDefinitionQuery)) *CheckListCategoryDefinitionQuery {
	query := &CheckListItemDefinitionQuery{config: clcdq.config}
	for _, opt := range opts {
		opt(query)
	}
	clcdq.withCheckListItemDefinitions = query
	return clcdq
}

//  WithWorkOrderType tells the query-builder to eager-loads the nodes that are connected to
// the "work_order_type" edge. The optional arguments used to configure the query builder of the edge.
func (clcdq *CheckListCategoryDefinitionQuery) WithWorkOrderType(opts ...func(*WorkOrderTypeQuery)) *CheckListCategoryDefinitionQuery {
	query := &WorkOrderTypeQuery{config: clcdq.config}
	for _, opt := range opts {
		opt(query)
	}
	clcdq.withWorkOrderType = query
	return clcdq
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
//	client.CheckListCategoryDefinition.Query().
//		GroupBy(checklistcategorydefinition.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (clcdq *CheckListCategoryDefinitionQuery) GroupBy(field string, fields ...string) *CheckListCategoryDefinitionGroupBy {
	group := &CheckListCategoryDefinitionGroupBy{config: clcdq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := clcdq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return clcdq.sqlQuery(), nil
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
//	client.CheckListCategoryDefinition.Query().
//		Select(checklistcategorydefinition.FieldCreateTime).
//		Scan(ctx, &v)
//
func (clcdq *CheckListCategoryDefinitionQuery) Select(field string, fields ...string) *CheckListCategoryDefinitionSelect {
	selector := &CheckListCategoryDefinitionSelect{config: clcdq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := clcdq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return clcdq.sqlQuery(), nil
	}
	return selector
}

func (clcdq *CheckListCategoryDefinitionQuery) prepareQuery(ctx context.Context) error {
	if clcdq.path != nil {
		prev, err := clcdq.path(ctx)
		if err != nil {
			return err
		}
		clcdq.sql = prev
	}
	if err := checklistcategorydefinition.Policy.EvalQuery(ctx, clcdq); err != nil {
		return err
	}
	return nil
}

func (clcdq *CheckListCategoryDefinitionQuery) sqlAll(ctx context.Context) ([]*CheckListCategoryDefinition, error) {
	var (
		nodes       = []*CheckListCategoryDefinition{}
		withFKs     = clcdq.withFKs
		_spec       = clcdq.querySpec()
		loadedTypes = [2]bool{
			clcdq.withCheckListItemDefinitions != nil,
			clcdq.withWorkOrderType != nil,
		}
	)
	if clcdq.withWorkOrderType != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, checklistcategorydefinition.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &CheckListCategoryDefinition{config: clcdq.config}
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
	if err := sqlgraph.QueryNodes(ctx, clcdq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := clcdq.withCheckListItemDefinitions; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*CheckListCategoryDefinition)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.CheckListItemDefinition(func(s *sql.Selector) {
			s.Where(sql.InValues(checklistcategorydefinition.CheckListItemDefinitionsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.check_list_category_definition_check_list_item_definitions
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "check_list_category_definition_check_list_item_definitions" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "check_list_category_definition_check_list_item_definitions" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.CheckListItemDefinitions = append(node.Edges.CheckListItemDefinitions, n)
		}
	}

	if query := clcdq.withWorkOrderType; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*CheckListCategoryDefinition)
		for i := range nodes {
			if fk := nodes[i].work_order_type_check_list_category_definitions; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(workordertype.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_type_check_list_category_definitions" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.WorkOrderType = n
			}
		}
	}

	return nodes, nil
}

func (clcdq *CheckListCategoryDefinitionQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := clcdq.querySpec()
	return sqlgraph.CountNodes(ctx, clcdq.driver, _spec)
}

func (clcdq *CheckListCategoryDefinitionQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := clcdq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (clcdq *CheckListCategoryDefinitionQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   checklistcategorydefinition.Table,
			Columns: checklistcategorydefinition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistcategorydefinition.FieldID,
			},
		},
		From:   clcdq.sql,
		Unique: true,
	}
	if ps := clcdq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := clcdq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := clcdq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := clcdq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (clcdq *CheckListCategoryDefinitionQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(clcdq.driver.Dialect())
	t1 := builder.Table(checklistcategorydefinition.Table)
	selector := builder.Select(t1.Columns(checklistcategorydefinition.Columns...)...).From(t1)
	if clcdq.sql != nil {
		selector = clcdq.sql
		selector.Select(selector.Columns(checklistcategorydefinition.Columns...)...)
	}
	for _, p := range clcdq.predicates {
		p(selector)
	}
	for _, p := range clcdq.order {
		p(selector)
	}
	if offset := clcdq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := clcdq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// CheckListCategoryDefinitionGroupBy is the builder for group-by CheckListCategoryDefinition entities.
type CheckListCategoryDefinitionGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (clcdgb *CheckListCategoryDefinitionGroupBy) Aggregate(fns ...AggregateFunc) *CheckListCategoryDefinitionGroupBy {
	clcdgb.fns = append(clcdgb.fns, fns...)
	return clcdgb
}

// Scan applies the group-by query and scan the result into the given value.
func (clcdgb *CheckListCategoryDefinitionGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := clcdgb.path(ctx)
	if err != nil {
		return err
	}
	clcdgb.sql = query
	return clcdgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (clcdgb *CheckListCategoryDefinitionGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := clcdgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (clcdgb *CheckListCategoryDefinitionGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(clcdgb.fields) > 1 {
		return nil, errors.New("ent: CheckListCategoryDefinitionGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := clcdgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (clcdgb *CheckListCategoryDefinitionGroupBy) StringsX(ctx context.Context) []string {
	v, err := clcdgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (clcdgb *CheckListCategoryDefinitionGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(clcdgb.fields) > 1 {
		return nil, errors.New("ent: CheckListCategoryDefinitionGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := clcdgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (clcdgb *CheckListCategoryDefinitionGroupBy) IntsX(ctx context.Context) []int {
	v, err := clcdgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (clcdgb *CheckListCategoryDefinitionGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(clcdgb.fields) > 1 {
		return nil, errors.New("ent: CheckListCategoryDefinitionGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := clcdgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (clcdgb *CheckListCategoryDefinitionGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := clcdgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (clcdgb *CheckListCategoryDefinitionGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(clcdgb.fields) > 1 {
		return nil, errors.New("ent: CheckListCategoryDefinitionGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := clcdgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (clcdgb *CheckListCategoryDefinitionGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := clcdgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (clcdgb *CheckListCategoryDefinitionGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := clcdgb.sqlQuery().Query()
	if err := clcdgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (clcdgb *CheckListCategoryDefinitionGroupBy) sqlQuery() *sql.Selector {
	selector := clcdgb.sql
	columns := make([]string, 0, len(clcdgb.fields)+len(clcdgb.fns))
	columns = append(columns, clcdgb.fields...)
	for _, fn := range clcdgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(clcdgb.fields...)
}

// CheckListCategoryDefinitionSelect is the builder for select fields of CheckListCategoryDefinition entities.
type CheckListCategoryDefinitionSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (clcds *CheckListCategoryDefinitionSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := clcds.path(ctx)
	if err != nil {
		return err
	}
	clcds.sql = query
	return clcds.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (clcds *CheckListCategoryDefinitionSelect) ScanX(ctx context.Context, v interface{}) {
	if err := clcds.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (clcds *CheckListCategoryDefinitionSelect) Strings(ctx context.Context) ([]string, error) {
	if len(clcds.fields) > 1 {
		return nil, errors.New("ent: CheckListCategoryDefinitionSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := clcds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (clcds *CheckListCategoryDefinitionSelect) StringsX(ctx context.Context) []string {
	v, err := clcds.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (clcds *CheckListCategoryDefinitionSelect) Ints(ctx context.Context) ([]int, error) {
	if len(clcds.fields) > 1 {
		return nil, errors.New("ent: CheckListCategoryDefinitionSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := clcds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (clcds *CheckListCategoryDefinitionSelect) IntsX(ctx context.Context) []int {
	v, err := clcds.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (clcds *CheckListCategoryDefinitionSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(clcds.fields) > 1 {
		return nil, errors.New("ent: CheckListCategoryDefinitionSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := clcds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (clcds *CheckListCategoryDefinitionSelect) Float64sX(ctx context.Context) []float64 {
	v, err := clcds.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (clcds *CheckListCategoryDefinitionSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(clcds.fields) > 1 {
		return nil, errors.New("ent: CheckListCategoryDefinitionSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := clcds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (clcds *CheckListCategoryDefinitionSelect) BoolsX(ctx context.Context) []bool {
	v, err := clcds.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (clcds *CheckListCategoryDefinitionSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := clcds.sqlQuery().Query()
	if err := clcds.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (clcds *CheckListCategoryDefinitionSelect) sqlQuery() sql.Querier {
	selector := clcds.sql
	selector.Select(selector.Columns(clcds.fields...)...)
	return selector
}
