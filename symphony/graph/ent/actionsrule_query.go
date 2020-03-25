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
	"github.com/facebookincubator/symphony/graph/ent/actionsrule"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ActionsRuleQuery is the builder for querying ActionsRule entities.
type ActionsRuleQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.ActionsRule
	// intermediate query.
	sql *sql.Selector
}

// Where adds a new predicate for the builder.
func (arq *ActionsRuleQuery) Where(ps ...predicate.ActionsRule) *ActionsRuleQuery {
	arq.predicates = append(arq.predicates, ps...)
	return arq
}

// Limit adds a limit step to the query.
func (arq *ActionsRuleQuery) Limit(limit int) *ActionsRuleQuery {
	arq.limit = &limit
	return arq
}

// Offset adds an offset step to the query.
func (arq *ActionsRuleQuery) Offset(offset int) *ActionsRuleQuery {
	arq.offset = &offset
	return arq
}

// Order adds an order step to the query.
func (arq *ActionsRuleQuery) Order(o ...Order) *ActionsRuleQuery {
	arq.order = append(arq.order, o...)
	return arq
}

// First returns the first ActionsRule entity in the query. Returns *NotFoundError when no actionsrule was found.
func (arq *ActionsRuleQuery) First(ctx context.Context) (*ActionsRule, error) {
	ars, err := arq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(ars) == 0 {
		return nil, &NotFoundError{actionsrule.Label}
	}
	return ars[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (arq *ActionsRuleQuery) FirstX(ctx context.Context) *ActionsRule {
	ar, err := arq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return ar
}

// FirstID returns the first ActionsRule id in the query. Returns *NotFoundError when no id was found.
func (arq *ActionsRuleQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = arq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{actionsrule.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (arq *ActionsRuleQuery) FirstXID(ctx context.Context) int {
	id, err := arq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only ActionsRule entity in the query, returns an error if not exactly one entity was returned.
func (arq *ActionsRuleQuery) Only(ctx context.Context) (*ActionsRule, error) {
	ars, err := arq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(ars) {
	case 1:
		return ars[0], nil
	case 0:
		return nil, &NotFoundError{actionsrule.Label}
	default:
		return nil, &NotSingularError{actionsrule.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (arq *ActionsRuleQuery) OnlyX(ctx context.Context) *ActionsRule {
	ar, err := arq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return ar
}

// OnlyID returns the only ActionsRule id in the query, returns an error if not exactly one id was returned.
func (arq *ActionsRuleQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = arq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{actionsrule.Label}
	default:
		err = &NotSingularError{actionsrule.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (arq *ActionsRuleQuery) OnlyXID(ctx context.Context) int {
	id, err := arq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ActionsRules.
func (arq *ActionsRuleQuery) All(ctx context.Context) ([]*ActionsRule, error) {
	return arq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (arq *ActionsRuleQuery) AllX(ctx context.Context) []*ActionsRule {
	ars, err := arq.All(ctx)
	if err != nil {
		panic(err)
	}
	return ars
}

// IDs executes the query and returns a list of ActionsRule ids.
func (arq *ActionsRuleQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := arq.Select(actionsrule.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (arq *ActionsRuleQuery) IDsX(ctx context.Context) []int {
	ids, err := arq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (arq *ActionsRuleQuery) Count(ctx context.Context) (int, error) {
	return arq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (arq *ActionsRuleQuery) CountX(ctx context.Context) int {
	count, err := arq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (arq *ActionsRuleQuery) Exist(ctx context.Context) (bool, error) {
	return arq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (arq *ActionsRuleQuery) ExistX(ctx context.Context) bool {
	exist, err := arq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (arq *ActionsRuleQuery) Clone() *ActionsRuleQuery {
	return &ActionsRuleQuery{
		config:     arq.config,
		limit:      arq.limit,
		offset:     arq.offset,
		order:      append([]Order{}, arq.order...),
		unique:     append([]string{}, arq.unique...),
		predicates: append([]predicate.ActionsRule{}, arq.predicates...),
		// clone intermediate query.
		sql: arq.sql.Clone(),
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
//	client.ActionsRule.Query().
//		GroupBy(actionsrule.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (arq *ActionsRuleQuery) GroupBy(field string, fields ...string) *ActionsRuleGroupBy {
	group := &ActionsRuleGroupBy{config: arq.config}
	group.fields = append([]string{field}, fields...)
	group.sql = arq.sqlQuery()
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
//	client.ActionsRule.Query().
//		Select(actionsrule.FieldCreateTime).
//		Scan(ctx, &v)
//
func (arq *ActionsRuleQuery) Select(field string, fields ...string) *ActionsRuleSelect {
	selector := &ActionsRuleSelect{config: arq.config}
	selector.fields = append([]string{field}, fields...)
	selector.sql = arq.sqlQuery()
	return selector
}

func (arq *ActionsRuleQuery) sqlAll(ctx context.Context) ([]*ActionsRule, error) {
	var (
		nodes = []*ActionsRule{}
		_spec = arq.querySpec()
	)
	_spec.ScanValues = func() []interface{} {
		node := &ActionsRule{config: arq.config}
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
	if err := sqlgraph.QueryNodes(ctx, arq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (arq *ActionsRuleQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := arq.querySpec()
	return sqlgraph.CountNodes(ctx, arq.driver, _spec)
}

func (arq *ActionsRuleQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := arq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (arq *ActionsRuleQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   actionsrule.Table,
			Columns: actionsrule.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: actionsrule.FieldID,
			},
		},
		From:   arq.sql,
		Unique: true,
	}
	if ps := arq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := arq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := arq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := arq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (arq *ActionsRuleQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(arq.driver.Dialect())
	t1 := builder.Table(actionsrule.Table)
	selector := builder.Select(t1.Columns(actionsrule.Columns...)...).From(t1)
	if arq.sql != nil {
		selector = arq.sql
		selector.Select(selector.Columns(actionsrule.Columns...)...)
	}
	for _, p := range arq.predicates {
		p(selector)
	}
	for _, p := range arq.order {
		p(selector)
	}
	if offset := arq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := arq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ActionsRuleGroupBy is the builder for group-by ActionsRule entities.
type ActionsRuleGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query.
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the group-by query.
func (argb *ActionsRuleGroupBy) Aggregate(fns ...Aggregate) *ActionsRuleGroupBy {
	argb.fns = append(argb.fns, fns...)
	return argb
}

// Scan applies the group-by query and scan the result into the given value.
func (argb *ActionsRuleGroupBy) Scan(ctx context.Context, v interface{}) error {
	return argb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (argb *ActionsRuleGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := argb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (argb *ActionsRuleGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(argb.fields) > 1 {
		return nil, errors.New("ent: ActionsRuleGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := argb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (argb *ActionsRuleGroupBy) StringsX(ctx context.Context) []string {
	v, err := argb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (argb *ActionsRuleGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(argb.fields) > 1 {
		return nil, errors.New("ent: ActionsRuleGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := argb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (argb *ActionsRuleGroupBy) IntsX(ctx context.Context) []int {
	v, err := argb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (argb *ActionsRuleGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(argb.fields) > 1 {
		return nil, errors.New("ent: ActionsRuleGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := argb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (argb *ActionsRuleGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := argb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (argb *ActionsRuleGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(argb.fields) > 1 {
		return nil, errors.New("ent: ActionsRuleGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := argb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (argb *ActionsRuleGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := argb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (argb *ActionsRuleGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := argb.sqlQuery().Query()
	if err := argb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (argb *ActionsRuleGroupBy) sqlQuery() *sql.Selector {
	selector := argb.sql
	columns := make([]string, 0, len(argb.fields)+len(argb.fns))
	columns = append(columns, argb.fields...)
	for _, fn := range argb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(argb.fields...)
}

// ActionsRuleSelect is the builder for select fields of ActionsRule entities.
type ActionsRuleSelect struct {
	config
	fields []string
	// intermediate queries.
	sql *sql.Selector
}

// Scan applies the selector query and scan the result into the given value.
func (ars *ActionsRuleSelect) Scan(ctx context.Context, v interface{}) error {
	return ars.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ars *ActionsRuleSelect) ScanX(ctx context.Context, v interface{}) {
	if err := ars.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (ars *ActionsRuleSelect) Strings(ctx context.Context) ([]string, error) {
	if len(ars.fields) > 1 {
		return nil, errors.New("ent: ActionsRuleSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := ars.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ars *ActionsRuleSelect) StringsX(ctx context.Context) []string {
	v, err := ars.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (ars *ActionsRuleSelect) Ints(ctx context.Context) ([]int, error) {
	if len(ars.fields) > 1 {
		return nil, errors.New("ent: ActionsRuleSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := ars.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ars *ActionsRuleSelect) IntsX(ctx context.Context) []int {
	v, err := ars.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (ars *ActionsRuleSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(ars.fields) > 1 {
		return nil, errors.New("ent: ActionsRuleSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := ars.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ars *ActionsRuleSelect) Float64sX(ctx context.Context) []float64 {
	v, err := ars.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (ars *ActionsRuleSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(ars.fields) > 1 {
		return nil, errors.New("ent: ActionsRuleSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := ars.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ars *ActionsRuleSelect) BoolsX(ctx context.Context) []bool {
	v, err := ars.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ars *ActionsRuleSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ars.sqlQuery().Query()
	if err := ars.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ars *ActionsRuleSelect) sqlQuery() sql.Querier {
	selector := ars.sql
	selector.Select(selector.Columns(ars.fields...)...)
	return selector
}
