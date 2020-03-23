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
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// CheckListItemQuery is the builder for querying CheckListItem entities.
type CheckListItemQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.CheckListItem
	// eager-loading edges.
	withFiles     *FileQuery
	withWorkOrder *WorkOrderQuery
	withFKs       bool
	// intermediate query.
	sql *sql.Selector
}

// Where adds a new predicate for the builder.
func (cliq *CheckListItemQuery) Where(ps ...predicate.CheckListItem) *CheckListItemQuery {
	cliq.predicates = append(cliq.predicates, ps...)
	return cliq
}

// Limit adds a limit step to the query.
func (cliq *CheckListItemQuery) Limit(limit int) *CheckListItemQuery {
	cliq.limit = &limit
	return cliq
}

// Offset adds an offset step to the query.
func (cliq *CheckListItemQuery) Offset(offset int) *CheckListItemQuery {
	cliq.offset = &offset
	return cliq
}

// Order adds an order step to the query.
func (cliq *CheckListItemQuery) Order(o ...Order) *CheckListItemQuery {
	cliq.order = append(cliq.order, o...)
	return cliq
}

// QueryFiles chains the current query on the files edge.
func (cliq *CheckListItemQuery) QueryFiles() *FileQuery {
	query := &FileQuery{config: cliq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(checklistitem.Table, checklistitem.FieldID, cliq.sqlQuery()),
		sqlgraph.To(file.Table, file.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, checklistitem.FilesTable, checklistitem.FilesColumn),
	)
	query.sql = sqlgraph.SetNeighbors(cliq.driver.Dialect(), step)
	return query
}

// QueryWorkOrder chains the current query on the work_order edge.
func (cliq *CheckListItemQuery) QueryWorkOrder() *WorkOrderQuery {
	query := &WorkOrderQuery{config: cliq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(checklistitem.Table, checklistitem.FieldID, cliq.sqlQuery()),
		sqlgraph.To(workorder.Table, workorder.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, checklistitem.WorkOrderTable, checklistitem.WorkOrderColumn),
	)
	query.sql = sqlgraph.SetNeighbors(cliq.driver.Dialect(), step)
	return query
}

// First returns the first CheckListItem entity in the query. Returns *NotFoundError when no checklistitem was found.
func (cliq *CheckListItemQuery) First(ctx context.Context) (*CheckListItem, error) {
	clis, err := cliq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(clis) == 0 {
		return nil, &NotFoundError{checklistitem.Label}
	}
	return clis[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (cliq *CheckListItemQuery) FirstX(ctx context.Context) *CheckListItem {
	cli, err := cliq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return cli
}

// FirstID returns the first CheckListItem id in the query. Returns *NotFoundError when no id was found.
func (cliq *CheckListItemQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = cliq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{checklistitem.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (cliq *CheckListItemQuery) FirstXID(ctx context.Context) int {
	id, err := cliq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only CheckListItem entity in the query, returns an error if not exactly one entity was returned.
func (cliq *CheckListItemQuery) Only(ctx context.Context) (*CheckListItem, error) {
	clis, err := cliq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(clis) {
	case 1:
		return clis[0], nil
	case 0:
		return nil, &NotFoundError{checklistitem.Label}
	default:
		return nil, &NotSingularError{checklistitem.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (cliq *CheckListItemQuery) OnlyX(ctx context.Context) *CheckListItem {
	cli, err := cliq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return cli
}

// OnlyID returns the only CheckListItem id in the query, returns an error if not exactly one id was returned.
func (cliq *CheckListItemQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = cliq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{checklistitem.Label}
	default:
		err = &NotSingularError{checklistitem.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (cliq *CheckListItemQuery) OnlyXID(ctx context.Context) int {
	id, err := cliq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of CheckListItems.
func (cliq *CheckListItemQuery) All(ctx context.Context) ([]*CheckListItem, error) {
	return cliq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (cliq *CheckListItemQuery) AllX(ctx context.Context) []*CheckListItem {
	clis, err := cliq.All(ctx)
	if err != nil {
		panic(err)
	}
	return clis
}

// IDs executes the query and returns a list of CheckListItem ids.
func (cliq *CheckListItemQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := cliq.Select(checklistitem.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (cliq *CheckListItemQuery) IDsX(ctx context.Context) []int {
	ids, err := cliq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (cliq *CheckListItemQuery) Count(ctx context.Context) (int, error) {
	return cliq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (cliq *CheckListItemQuery) CountX(ctx context.Context) int {
	count, err := cliq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (cliq *CheckListItemQuery) Exist(ctx context.Context) (bool, error) {
	return cliq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (cliq *CheckListItemQuery) ExistX(ctx context.Context) bool {
	exist, err := cliq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (cliq *CheckListItemQuery) Clone() *CheckListItemQuery {
	return &CheckListItemQuery{
		config:     cliq.config,
		limit:      cliq.limit,
		offset:     cliq.offset,
		order:      append([]Order{}, cliq.order...),
		unique:     append([]string{}, cliq.unique...),
		predicates: append([]predicate.CheckListItem{}, cliq.predicates...),
		// clone intermediate query.
		sql: cliq.sql.Clone(),
	}
}

//  WithFiles tells the query-builder to eager-loads the nodes that are connected to
// the "files" edge. The optional arguments used to configure the query builder of the edge.
func (cliq *CheckListItemQuery) WithFiles(opts ...func(*FileQuery)) *CheckListItemQuery {
	query := &FileQuery{config: cliq.config}
	for _, opt := range opts {
		opt(query)
	}
	cliq.withFiles = query
	return cliq
}

//  WithWorkOrder tells the query-builder to eager-loads the nodes that are connected to
// the "work_order" edge. The optional arguments used to configure the query builder of the edge.
func (cliq *CheckListItemQuery) WithWorkOrder(opts ...func(*WorkOrderQuery)) *CheckListItemQuery {
	query := &WorkOrderQuery{config: cliq.config}
	for _, opt := range opts {
		opt(query)
	}
	cliq.withWorkOrder = query
	return cliq
}

// GroupBy used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Title string `json:"title,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.CheckListItem.Query().
//		GroupBy(checklistitem.FieldTitle).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (cliq *CheckListItemQuery) GroupBy(field string, fields ...string) *CheckListItemGroupBy {
	group := &CheckListItemGroupBy{config: cliq.config}
	group.fields = append([]string{field}, fields...)
	group.sql = cliq.sqlQuery()
	return group
}

// Select one or more fields from the given query.
//
// Example:
//
//	var v []struct {
//		Title string `json:"title,omitempty"`
//	}
//
//	client.CheckListItem.Query().
//		Select(checklistitem.FieldTitle).
//		Scan(ctx, &v)
//
func (cliq *CheckListItemQuery) Select(field string, fields ...string) *CheckListItemSelect {
	selector := &CheckListItemSelect{config: cliq.config}
	selector.fields = append([]string{field}, fields...)
	selector.sql = cliq.sqlQuery()
	return selector
}

func (cliq *CheckListItemQuery) sqlAll(ctx context.Context) ([]*CheckListItem, error) {
	var (
		nodes       = []*CheckListItem{}
		withFKs     = cliq.withFKs
		_spec       = cliq.querySpec()
		loadedTypes = [2]bool{
			cliq.withFiles != nil,
			cliq.withWorkOrder != nil,
		}
	)
	if cliq.withWorkOrder != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, checklistitem.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &CheckListItem{config: cliq.config}
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
	if err := sqlgraph.QueryNodes(ctx, cliq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := cliq.withFiles; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*CheckListItem)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.File(func(s *sql.Selector) {
			s.Where(sql.InValues(checklistitem.FilesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.check_list_item_files
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "check_list_item_files" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "check_list_item_files" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Files = append(node.Edges.Files, n)
		}
	}

	if query := cliq.withWorkOrder; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*CheckListItem)
		for i := range nodes {
			if fk := nodes[i].work_order_check_list_items; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(workorder.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_check_list_items" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.WorkOrder = n
			}
		}
	}

	return nodes, nil
}

func (cliq *CheckListItemQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := cliq.querySpec()
	return sqlgraph.CountNodes(ctx, cliq.driver, _spec)
}

func (cliq *CheckListItemQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := cliq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (cliq *CheckListItemQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   checklistitem.Table,
			Columns: checklistitem.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistitem.FieldID,
			},
		},
		From:   cliq.sql,
		Unique: true,
	}
	if ps := cliq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := cliq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := cliq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := cliq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (cliq *CheckListItemQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(cliq.driver.Dialect())
	t1 := builder.Table(checklistitem.Table)
	selector := builder.Select(t1.Columns(checklistitem.Columns...)...).From(t1)
	if cliq.sql != nil {
		selector = cliq.sql
		selector.Select(selector.Columns(checklistitem.Columns...)...)
	}
	for _, p := range cliq.predicates {
		p(selector)
	}
	for _, p := range cliq.order {
		p(selector)
	}
	if offset := cliq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := cliq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// CheckListItemGroupBy is the builder for group-by CheckListItem entities.
type CheckListItemGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query.
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the group-by query.
func (cligb *CheckListItemGroupBy) Aggregate(fns ...Aggregate) *CheckListItemGroupBy {
	cligb.fns = append(cligb.fns, fns...)
	return cligb
}

// Scan applies the group-by query and scan the result into the given value.
func (cligb *CheckListItemGroupBy) Scan(ctx context.Context, v interface{}) error {
	return cligb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (cligb *CheckListItemGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := cligb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (cligb *CheckListItemGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(cligb.fields) > 1 {
		return nil, errors.New("ent: CheckListItemGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := cligb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (cligb *CheckListItemGroupBy) StringsX(ctx context.Context) []string {
	v, err := cligb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (cligb *CheckListItemGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(cligb.fields) > 1 {
		return nil, errors.New("ent: CheckListItemGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := cligb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (cligb *CheckListItemGroupBy) IntsX(ctx context.Context) []int {
	v, err := cligb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (cligb *CheckListItemGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(cligb.fields) > 1 {
		return nil, errors.New("ent: CheckListItemGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := cligb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (cligb *CheckListItemGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := cligb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (cligb *CheckListItemGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(cligb.fields) > 1 {
		return nil, errors.New("ent: CheckListItemGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := cligb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (cligb *CheckListItemGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := cligb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (cligb *CheckListItemGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := cligb.sqlQuery().Query()
	if err := cligb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (cligb *CheckListItemGroupBy) sqlQuery() *sql.Selector {
	selector := cligb.sql
	columns := make([]string, 0, len(cligb.fields)+len(cligb.fns))
	columns = append(columns, cligb.fields...)
	for _, fn := range cligb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(cligb.fields...)
}

// CheckListItemSelect is the builder for select fields of CheckListItem entities.
type CheckListItemSelect struct {
	config
	fields []string
	// intermediate queries.
	sql *sql.Selector
}

// Scan applies the selector query and scan the result into the given value.
func (clis *CheckListItemSelect) Scan(ctx context.Context, v interface{}) error {
	return clis.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (clis *CheckListItemSelect) ScanX(ctx context.Context, v interface{}) {
	if err := clis.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (clis *CheckListItemSelect) Strings(ctx context.Context) ([]string, error) {
	if len(clis.fields) > 1 {
		return nil, errors.New("ent: CheckListItemSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := clis.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (clis *CheckListItemSelect) StringsX(ctx context.Context) []string {
	v, err := clis.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (clis *CheckListItemSelect) Ints(ctx context.Context) ([]int, error) {
	if len(clis.fields) > 1 {
		return nil, errors.New("ent: CheckListItemSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := clis.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (clis *CheckListItemSelect) IntsX(ctx context.Context) []int {
	v, err := clis.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (clis *CheckListItemSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(clis.fields) > 1 {
		return nil, errors.New("ent: CheckListItemSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := clis.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (clis *CheckListItemSelect) Float64sX(ctx context.Context) []float64 {
	v, err := clis.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (clis *CheckListItemSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(clis.fields) > 1 {
		return nil, errors.New("ent: CheckListItemSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := clis.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (clis *CheckListItemSelect) BoolsX(ctx context.Context) []bool {
	v, err := clis.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (clis *CheckListItemSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := clis.sqlQuery().Query()
	if err := clis.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (clis *CheckListItemSelect) sqlQuery() sql.Querier {
	selector := clis.sql
	selector.Select(selector.Columns(clis.fields...)...)
	return selector
}
