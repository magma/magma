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
	"github.com/facebookincubator/symphony/graph/ent/checklistcategorydefinition"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// WorkOrderTypeQuery is the builder for querying WorkOrderType entities.
type WorkOrderTypeQuery struct {
	config
	limit      *int
	offset     *int
	order      []OrderFunc
	unique     []string
	predicates []predicate.WorkOrderType
	// eager-loading edges.
	withWorkOrders                   *WorkOrderQuery
	withPropertyTypes                *PropertyTypeQuery
	withDefinitions                  *WorkOrderDefinitionQuery
	withCheckListCategoryDefinitions *CheckListCategoryDefinitionQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (wotq *WorkOrderTypeQuery) Where(ps ...predicate.WorkOrderType) *WorkOrderTypeQuery {
	wotq.predicates = append(wotq.predicates, ps...)
	return wotq
}

// Limit adds a limit step to the query.
func (wotq *WorkOrderTypeQuery) Limit(limit int) *WorkOrderTypeQuery {
	wotq.limit = &limit
	return wotq
}

// Offset adds an offset step to the query.
func (wotq *WorkOrderTypeQuery) Offset(offset int) *WorkOrderTypeQuery {
	wotq.offset = &offset
	return wotq
}

// Order adds an order step to the query.
func (wotq *WorkOrderTypeQuery) Order(o ...OrderFunc) *WorkOrderTypeQuery {
	wotq.order = append(wotq.order, o...)
	return wotq
}

// QueryWorkOrders chains the current query on the work_orders edge.
func (wotq *WorkOrderTypeQuery) QueryWorkOrders() *WorkOrderQuery {
	query := &WorkOrderQuery{config: wotq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := wotq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workordertype.Table, workordertype.FieldID, wotq.sqlQuery()),
			sqlgraph.To(workorder.Table, workorder.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, workordertype.WorkOrdersTable, workordertype.WorkOrdersColumn),
		)
		fromU = sqlgraph.SetNeighbors(wotq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryPropertyTypes chains the current query on the property_types edge.
func (wotq *WorkOrderTypeQuery) QueryPropertyTypes() *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: wotq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := wotq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workordertype.Table, workordertype.FieldID, wotq.sqlQuery()),
			sqlgraph.To(propertytype.Table, propertytype.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, workordertype.PropertyTypesTable, workordertype.PropertyTypesColumn),
		)
		fromU = sqlgraph.SetNeighbors(wotq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryDefinitions chains the current query on the definitions edge.
func (wotq *WorkOrderTypeQuery) QueryDefinitions() *WorkOrderDefinitionQuery {
	query := &WorkOrderDefinitionQuery{config: wotq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := wotq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workordertype.Table, workordertype.FieldID, wotq.sqlQuery()),
			sqlgraph.To(workorderdefinition.Table, workorderdefinition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, workordertype.DefinitionsTable, workordertype.DefinitionsColumn),
		)
		fromU = sqlgraph.SetNeighbors(wotq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryCheckListCategoryDefinitions chains the current query on the check_list_category_definitions edge.
func (wotq *WorkOrderTypeQuery) QueryCheckListCategoryDefinitions() *CheckListCategoryDefinitionQuery {
	query := &CheckListCategoryDefinitionQuery{config: wotq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := wotq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workordertype.Table, workordertype.FieldID, wotq.sqlQuery()),
			sqlgraph.To(checklistcategorydefinition.Table, checklistcategorydefinition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, workordertype.CheckListCategoryDefinitionsTable, workordertype.CheckListCategoryDefinitionsColumn),
		)
		fromU = sqlgraph.SetNeighbors(wotq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first WorkOrderType entity in the query. Returns *NotFoundError when no workordertype was found.
func (wotq *WorkOrderTypeQuery) First(ctx context.Context) (*WorkOrderType, error) {
	wots, err := wotq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(wots) == 0 {
		return nil, &NotFoundError{workordertype.Label}
	}
	return wots[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (wotq *WorkOrderTypeQuery) FirstX(ctx context.Context) *WorkOrderType {
	wot, err := wotq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return wot
}

// FirstID returns the first WorkOrderType id in the query. Returns *NotFoundError when no id was found.
func (wotq *WorkOrderTypeQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = wotq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{workordertype.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (wotq *WorkOrderTypeQuery) FirstXID(ctx context.Context) int {
	id, err := wotq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only WorkOrderType entity in the query, returns an error if not exactly one entity was returned.
func (wotq *WorkOrderTypeQuery) Only(ctx context.Context) (*WorkOrderType, error) {
	wots, err := wotq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(wots) {
	case 1:
		return wots[0], nil
	case 0:
		return nil, &NotFoundError{workordertype.Label}
	default:
		return nil, &NotSingularError{workordertype.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (wotq *WorkOrderTypeQuery) OnlyX(ctx context.Context) *WorkOrderType {
	wot, err := wotq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return wot
}

// OnlyID returns the only WorkOrderType id in the query, returns an error if not exactly one id was returned.
func (wotq *WorkOrderTypeQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = wotq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{workordertype.Label}
	default:
		err = &NotSingularError{workordertype.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (wotq *WorkOrderTypeQuery) OnlyXID(ctx context.Context) int {
	id, err := wotq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of WorkOrderTypes.
func (wotq *WorkOrderTypeQuery) All(ctx context.Context) ([]*WorkOrderType, error) {
	if err := wotq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return wotq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (wotq *WorkOrderTypeQuery) AllX(ctx context.Context) []*WorkOrderType {
	wots, err := wotq.All(ctx)
	if err != nil {
		panic(err)
	}
	return wots
}

// IDs executes the query and returns a list of WorkOrderType ids.
func (wotq *WorkOrderTypeQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := wotq.Select(workordertype.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (wotq *WorkOrderTypeQuery) IDsX(ctx context.Context) []int {
	ids, err := wotq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (wotq *WorkOrderTypeQuery) Count(ctx context.Context) (int, error) {
	if err := wotq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return wotq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (wotq *WorkOrderTypeQuery) CountX(ctx context.Context) int {
	count, err := wotq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (wotq *WorkOrderTypeQuery) Exist(ctx context.Context) (bool, error) {
	if err := wotq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return wotq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (wotq *WorkOrderTypeQuery) ExistX(ctx context.Context) bool {
	exist, err := wotq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (wotq *WorkOrderTypeQuery) Clone() *WorkOrderTypeQuery {
	return &WorkOrderTypeQuery{
		config:     wotq.config,
		limit:      wotq.limit,
		offset:     wotq.offset,
		order:      append([]OrderFunc{}, wotq.order...),
		unique:     append([]string{}, wotq.unique...),
		predicates: append([]predicate.WorkOrderType{}, wotq.predicates...),
		// clone intermediate query.
		sql:  wotq.sql.Clone(),
		path: wotq.path,
	}
}

//  WithWorkOrders tells the query-builder to eager-loads the nodes that are connected to
// the "work_orders" edge. The optional arguments used to configure the query builder of the edge.
func (wotq *WorkOrderTypeQuery) WithWorkOrders(opts ...func(*WorkOrderQuery)) *WorkOrderTypeQuery {
	query := &WorkOrderQuery{config: wotq.config}
	for _, opt := range opts {
		opt(query)
	}
	wotq.withWorkOrders = query
	return wotq
}

//  WithPropertyTypes tells the query-builder to eager-loads the nodes that are connected to
// the "property_types" edge. The optional arguments used to configure the query builder of the edge.
func (wotq *WorkOrderTypeQuery) WithPropertyTypes(opts ...func(*PropertyTypeQuery)) *WorkOrderTypeQuery {
	query := &PropertyTypeQuery{config: wotq.config}
	for _, opt := range opts {
		opt(query)
	}
	wotq.withPropertyTypes = query
	return wotq
}

//  WithDefinitions tells the query-builder to eager-loads the nodes that are connected to
// the "definitions" edge. The optional arguments used to configure the query builder of the edge.
func (wotq *WorkOrderTypeQuery) WithDefinitions(opts ...func(*WorkOrderDefinitionQuery)) *WorkOrderTypeQuery {
	query := &WorkOrderDefinitionQuery{config: wotq.config}
	for _, opt := range opts {
		opt(query)
	}
	wotq.withDefinitions = query
	return wotq
}

//  WithCheckListCategoryDefinitions tells the query-builder to eager-loads the nodes that are connected to
// the "check_list_category_definitions" edge. The optional arguments used to configure the query builder of the edge.
func (wotq *WorkOrderTypeQuery) WithCheckListCategoryDefinitions(opts ...func(*CheckListCategoryDefinitionQuery)) *WorkOrderTypeQuery {
	query := &CheckListCategoryDefinitionQuery{config: wotq.config}
	for _, opt := range opts {
		opt(query)
	}
	wotq.withCheckListCategoryDefinitions = query
	return wotq
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
//	client.WorkOrderType.Query().
//		GroupBy(workordertype.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (wotq *WorkOrderTypeQuery) GroupBy(field string, fields ...string) *WorkOrderTypeGroupBy {
	group := &WorkOrderTypeGroupBy{config: wotq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := wotq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return wotq.sqlQuery(), nil
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
//	client.WorkOrderType.Query().
//		Select(workordertype.FieldCreateTime).
//		Scan(ctx, &v)
//
func (wotq *WorkOrderTypeQuery) Select(field string, fields ...string) *WorkOrderTypeSelect {
	selector := &WorkOrderTypeSelect{config: wotq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := wotq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return wotq.sqlQuery(), nil
	}
	return selector
}

func (wotq *WorkOrderTypeQuery) prepareQuery(ctx context.Context) error {
	if wotq.path != nil {
		prev, err := wotq.path(ctx)
		if err != nil {
			return err
		}
		wotq.sql = prev
	}
	if err := workordertype.Policy.EvalQuery(ctx, wotq); err != nil {
		return err
	}
	return nil
}

func (wotq *WorkOrderTypeQuery) sqlAll(ctx context.Context) ([]*WorkOrderType, error) {
	var (
		nodes       = []*WorkOrderType{}
		_spec       = wotq.querySpec()
		loadedTypes = [4]bool{
			wotq.withWorkOrders != nil,
			wotq.withPropertyTypes != nil,
			wotq.withDefinitions != nil,
			wotq.withCheckListCategoryDefinitions != nil,
		}
	)
	_spec.ScanValues = func() []interface{} {
		node := &WorkOrderType{config: wotq.config}
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
	if err := sqlgraph.QueryNodes(ctx, wotq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := wotq.withWorkOrders; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*WorkOrderType)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.WorkOrder(func(s *sql.Selector) {
			s.Where(sql.InValues(workordertype.WorkOrdersColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.work_order_type
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "work_order_type" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_type" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.WorkOrders = append(node.Edges.WorkOrders, n)
		}
	}

	if query := wotq.withPropertyTypes; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*WorkOrderType)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.PropertyType(func(s *sql.Selector) {
			s.Where(sql.InValues(workordertype.PropertyTypesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.work_order_type_property_types
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "work_order_type_property_types" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_type_property_types" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.PropertyTypes = append(node.Edges.PropertyTypes, n)
		}
	}

	if query := wotq.withDefinitions; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*WorkOrderType)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.WorkOrderDefinition(func(s *sql.Selector) {
			s.Where(sql.InValues(workordertype.DefinitionsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.work_order_definition_type
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "work_order_definition_type" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_definition_type" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Definitions = append(node.Edges.Definitions, n)
		}
	}

	if query := wotq.withCheckListCategoryDefinitions; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*WorkOrderType)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.CheckListCategoryDefinition(func(s *sql.Selector) {
			s.Where(sql.InValues(workordertype.CheckListCategoryDefinitionsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.work_order_type_check_list_category_definitions
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "work_order_type_check_list_category_definitions" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_type_check_list_category_definitions" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.CheckListCategoryDefinitions = append(node.Edges.CheckListCategoryDefinitions, n)
		}
	}

	return nodes, nil
}

func (wotq *WorkOrderTypeQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := wotq.querySpec()
	return sqlgraph.CountNodes(ctx, wotq.driver, _spec)
}

func (wotq *WorkOrderTypeQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := wotq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (wotq *WorkOrderTypeQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   workordertype.Table,
			Columns: workordertype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: workordertype.FieldID,
			},
		},
		From:   wotq.sql,
		Unique: true,
	}
	if ps := wotq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := wotq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := wotq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := wotq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (wotq *WorkOrderTypeQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(wotq.driver.Dialect())
	t1 := builder.Table(workordertype.Table)
	selector := builder.Select(t1.Columns(workordertype.Columns...)...).From(t1)
	if wotq.sql != nil {
		selector = wotq.sql
		selector.Select(selector.Columns(workordertype.Columns...)...)
	}
	for _, p := range wotq.predicates {
		p(selector)
	}
	for _, p := range wotq.order {
		p(selector)
	}
	if offset := wotq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := wotq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// WorkOrderTypeGroupBy is the builder for group-by WorkOrderType entities.
type WorkOrderTypeGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (wotgb *WorkOrderTypeGroupBy) Aggregate(fns ...AggregateFunc) *WorkOrderTypeGroupBy {
	wotgb.fns = append(wotgb.fns, fns...)
	return wotgb
}

// Scan applies the group-by query and scan the result into the given value.
func (wotgb *WorkOrderTypeGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := wotgb.path(ctx)
	if err != nil {
		return err
	}
	wotgb.sql = query
	return wotgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (wotgb *WorkOrderTypeGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := wotgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (wotgb *WorkOrderTypeGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(wotgb.fields) > 1 {
		return nil, errors.New("ent: WorkOrderTypeGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := wotgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (wotgb *WorkOrderTypeGroupBy) StringsX(ctx context.Context) []string {
	v, err := wotgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (wotgb *WorkOrderTypeGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(wotgb.fields) > 1 {
		return nil, errors.New("ent: WorkOrderTypeGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := wotgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (wotgb *WorkOrderTypeGroupBy) IntsX(ctx context.Context) []int {
	v, err := wotgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (wotgb *WorkOrderTypeGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(wotgb.fields) > 1 {
		return nil, errors.New("ent: WorkOrderTypeGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := wotgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (wotgb *WorkOrderTypeGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := wotgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (wotgb *WorkOrderTypeGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(wotgb.fields) > 1 {
		return nil, errors.New("ent: WorkOrderTypeGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := wotgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (wotgb *WorkOrderTypeGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := wotgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (wotgb *WorkOrderTypeGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := wotgb.sqlQuery().Query()
	if err := wotgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (wotgb *WorkOrderTypeGroupBy) sqlQuery() *sql.Selector {
	selector := wotgb.sql
	columns := make([]string, 0, len(wotgb.fields)+len(wotgb.fns))
	columns = append(columns, wotgb.fields...)
	for _, fn := range wotgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(wotgb.fields...)
}

// WorkOrderTypeSelect is the builder for select fields of WorkOrderType entities.
type WorkOrderTypeSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (wots *WorkOrderTypeSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := wots.path(ctx)
	if err != nil {
		return err
	}
	wots.sql = query
	return wots.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (wots *WorkOrderTypeSelect) ScanX(ctx context.Context, v interface{}) {
	if err := wots.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (wots *WorkOrderTypeSelect) Strings(ctx context.Context) ([]string, error) {
	if len(wots.fields) > 1 {
		return nil, errors.New("ent: WorkOrderTypeSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := wots.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (wots *WorkOrderTypeSelect) StringsX(ctx context.Context) []string {
	v, err := wots.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (wots *WorkOrderTypeSelect) Ints(ctx context.Context) ([]int, error) {
	if len(wots.fields) > 1 {
		return nil, errors.New("ent: WorkOrderTypeSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := wots.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (wots *WorkOrderTypeSelect) IntsX(ctx context.Context) []int {
	v, err := wots.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (wots *WorkOrderTypeSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(wots.fields) > 1 {
		return nil, errors.New("ent: WorkOrderTypeSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := wots.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (wots *WorkOrderTypeSelect) Float64sX(ctx context.Context) []float64 {
	v, err := wots.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (wots *WorkOrderTypeSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(wots.fields) > 1 {
		return nil, errors.New("ent: WorkOrderTypeSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := wots.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (wots *WorkOrderTypeSelect) BoolsX(ctx context.Context) []bool {
	v, err := wots.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (wots *WorkOrderTypeSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := wots.sqlQuery().Query()
	if err := wots.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (wots *WorkOrderTypeSelect) sqlQuery() sql.Querier {
	selector := wots.sql
	selector.Select(selector.Columns(wots.fields...)...)
	return selector
}
