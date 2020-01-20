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
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// EquipmentPositionDefinitionQuery is the builder for querying EquipmentPositionDefinition entities.
type EquipmentPositionDefinitionQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.EquipmentPositionDefinition
	// intermediate query.
	sql *sql.Selector
}

// Where adds a new predicate for the builder.
func (epdq *EquipmentPositionDefinitionQuery) Where(ps ...predicate.EquipmentPositionDefinition) *EquipmentPositionDefinitionQuery {
	epdq.predicates = append(epdq.predicates, ps...)
	return epdq
}

// Limit adds a limit step to the query.
func (epdq *EquipmentPositionDefinitionQuery) Limit(limit int) *EquipmentPositionDefinitionQuery {
	epdq.limit = &limit
	return epdq
}

// Offset adds an offset step to the query.
func (epdq *EquipmentPositionDefinitionQuery) Offset(offset int) *EquipmentPositionDefinitionQuery {
	epdq.offset = &offset
	return epdq
}

// Order adds an order step to the query.
func (epdq *EquipmentPositionDefinitionQuery) Order(o ...Order) *EquipmentPositionDefinitionQuery {
	epdq.order = append(epdq.order, o...)
	return epdq
}

// QueryPositions chains the current query on the positions edge.
func (epdq *EquipmentPositionDefinitionQuery) QueryPositions() *EquipmentPositionQuery {
	query := &EquipmentPositionQuery{config: epdq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentpositiondefinition.Table, equipmentpositiondefinition.FieldID, epdq.sqlQuery()),
		sqlgraph.To(equipmentposition.Table, equipmentposition.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, equipmentpositiondefinition.PositionsTable, equipmentpositiondefinition.PositionsColumn),
	)
	query.sql = sqlgraph.SetNeighbors(epdq.driver.Dialect(), step)
	return query
}

// QueryEquipmentType chains the current query on the equipment_type edge.
func (epdq *EquipmentPositionDefinitionQuery) QueryEquipmentType() *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: epdq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentpositiondefinition.Table, equipmentpositiondefinition.FieldID, epdq.sqlQuery()),
		sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, equipmentpositiondefinition.EquipmentTypeTable, equipmentpositiondefinition.EquipmentTypeColumn),
	)
	query.sql = sqlgraph.SetNeighbors(epdq.driver.Dialect(), step)
	return query
}

// First returns the first EquipmentPositionDefinition entity in the query. Returns *ErrNotFound when no equipmentpositiondefinition was found.
func (epdq *EquipmentPositionDefinitionQuery) First(ctx context.Context) (*EquipmentPositionDefinition, error) {
	epds, err := epdq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(epds) == 0 {
		return nil, &ErrNotFound{equipmentpositiondefinition.Label}
	}
	return epds[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (epdq *EquipmentPositionDefinitionQuery) FirstX(ctx context.Context) *EquipmentPositionDefinition {
	epd, err := epdq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return epd
}

// FirstID returns the first EquipmentPositionDefinition id in the query. Returns *ErrNotFound when no id was found.
func (epdq *EquipmentPositionDefinitionQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = epdq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &ErrNotFound{equipmentpositiondefinition.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (epdq *EquipmentPositionDefinitionQuery) FirstXID(ctx context.Context) string {
	id, err := epdq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only EquipmentPositionDefinition entity in the query, returns an error if not exactly one entity was returned.
func (epdq *EquipmentPositionDefinitionQuery) Only(ctx context.Context) (*EquipmentPositionDefinition, error) {
	epds, err := epdq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(epds) {
	case 1:
		return epds[0], nil
	case 0:
		return nil, &ErrNotFound{equipmentpositiondefinition.Label}
	default:
		return nil, &ErrNotSingular{equipmentpositiondefinition.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (epdq *EquipmentPositionDefinitionQuery) OnlyX(ctx context.Context) *EquipmentPositionDefinition {
	epd, err := epdq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return epd
}

// OnlyID returns the only EquipmentPositionDefinition id in the query, returns an error if not exactly one id was returned.
func (epdq *EquipmentPositionDefinitionQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = epdq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &ErrNotFound{equipmentpositiondefinition.Label}
	default:
		err = &ErrNotSingular{equipmentpositiondefinition.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (epdq *EquipmentPositionDefinitionQuery) OnlyXID(ctx context.Context) string {
	id, err := epdq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of EquipmentPositionDefinitions.
func (epdq *EquipmentPositionDefinitionQuery) All(ctx context.Context) ([]*EquipmentPositionDefinition, error) {
	return epdq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (epdq *EquipmentPositionDefinitionQuery) AllX(ctx context.Context) []*EquipmentPositionDefinition {
	epds, err := epdq.All(ctx)
	if err != nil {
		panic(err)
	}
	return epds
}

// IDs executes the query and returns a list of EquipmentPositionDefinition ids.
func (epdq *EquipmentPositionDefinitionQuery) IDs(ctx context.Context) ([]string, error) {
	var ids []string
	if err := epdq.Select(equipmentpositiondefinition.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (epdq *EquipmentPositionDefinitionQuery) IDsX(ctx context.Context) []string {
	ids, err := epdq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (epdq *EquipmentPositionDefinitionQuery) Count(ctx context.Context) (int, error) {
	return epdq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (epdq *EquipmentPositionDefinitionQuery) CountX(ctx context.Context) int {
	count, err := epdq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (epdq *EquipmentPositionDefinitionQuery) Exist(ctx context.Context) (bool, error) {
	return epdq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (epdq *EquipmentPositionDefinitionQuery) ExistX(ctx context.Context) bool {
	exist, err := epdq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (epdq *EquipmentPositionDefinitionQuery) Clone() *EquipmentPositionDefinitionQuery {
	return &EquipmentPositionDefinitionQuery{
		config:     epdq.config,
		limit:      epdq.limit,
		offset:     epdq.offset,
		order:      append([]Order{}, epdq.order...),
		unique:     append([]string{}, epdq.unique...),
		predicates: append([]predicate.EquipmentPositionDefinition{}, epdq.predicates...),
		// clone intermediate query.
		sql: epdq.sql.Clone(),
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
//	client.EquipmentPositionDefinition.Query().
//		GroupBy(equipmentpositiondefinition.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (epdq *EquipmentPositionDefinitionQuery) GroupBy(field string, fields ...string) *EquipmentPositionDefinitionGroupBy {
	group := &EquipmentPositionDefinitionGroupBy{config: epdq.config}
	group.fields = append([]string{field}, fields...)
	group.sql = epdq.sqlQuery()
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
//	client.EquipmentPositionDefinition.Query().
//		Select(equipmentpositiondefinition.FieldCreateTime).
//		Scan(ctx, &v)
//
func (epdq *EquipmentPositionDefinitionQuery) Select(field string, fields ...string) *EquipmentPositionDefinitionSelect {
	selector := &EquipmentPositionDefinitionSelect{config: epdq.config}
	selector.fields = append([]string{field}, fields...)
	selector.sql = epdq.sqlQuery()
	return selector
}

func (epdq *EquipmentPositionDefinitionQuery) sqlAll(ctx context.Context) ([]*EquipmentPositionDefinition, error) {
	var (
		nodes []*EquipmentPositionDefinition
		spec  = epdq.querySpec()
	)
	spec.ScanValues = func() []interface{} {
		node := &EquipmentPositionDefinition{config: epdq.config}
		nodes = append(nodes, node)
		return node.scanValues()
	}
	spec.Assign = func(values ...interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("ent: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		return node.assignValues(values...)
	}
	if err := sqlgraph.QueryNodes(ctx, epdq.driver, spec); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (epdq *EquipmentPositionDefinitionQuery) sqlCount(ctx context.Context) (int, error) {
	spec := epdq.querySpec()
	return sqlgraph.CountNodes(ctx, epdq.driver, spec)
}

func (epdq *EquipmentPositionDefinitionQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := epdq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (epdq *EquipmentPositionDefinitionQuery) querySpec() *sqlgraph.QuerySpec {
	spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentpositiondefinition.Table,
			Columns: equipmentpositiondefinition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: equipmentpositiondefinition.FieldID,
			},
		},
		From:   epdq.sql,
		Unique: true,
	}
	if ps := epdq.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := epdq.limit; limit != nil {
		spec.Limit = *limit
	}
	if offset := epdq.offset; offset != nil {
		spec.Offset = *offset
	}
	if ps := epdq.order; len(ps) > 0 {
		spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return spec
}

func (epdq *EquipmentPositionDefinitionQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(epdq.driver.Dialect())
	t1 := builder.Table(equipmentpositiondefinition.Table)
	selector := builder.Select(t1.Columns(equipmentpositiondefinition.Columns...)...).From(t1)
	if epdq.sql != nil {
		selector = epdq.sql
		selector.Select(selector.Columns(equipmentpositiondefinition.Columns...)...)
	}
	for _, p := range epdq.predicates {
		p(selector)
	}
	for _, p := range epdq.order {
		p(selector)
	}
	if offset := epdq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := epdq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// EquipmentPositionDefinitionGroupBy is the builder for group-by EquipmentPositionDefinition entities.
type EquipmentPositionDefinitionGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query.
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the group-by query.
func (epdgb *EquipmentPositionDefinitionGroupBy) Aggregate(fns ...Aggregate) *EquipmentPositionDefinitionGroupBy {
	epdgb.fns = append(epdgb.fns, fns...)
	return epdgb
}

// Scan applies the group-by query and scan the result into the given value.
func (epdgb *EquipmentPositionDefinitionGroupBy) Scan(ctx context.Context, v interface{}) error {
	return epdgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (epdgb *EquipmentPositionDefinitionGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := epdgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (epdgb *EquipmentPositionDefinitionGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(epdgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentPositionDefinitionGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := epdgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (epdgb *EquipmentPositionDefinitionGroupBy) StringsX(ctx context.Context) []string {
	v, err := epdgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (epdgb *EquipmentPositionDefinitionGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(epdgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentPositionDefinitionGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := epdgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (epdgb *EquipmentPositionDefinitionGroupBy) IntsX(ctx context.Context) []int {
	v, err := epdgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (epdgb *EquipmentPositionDefinitionGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(epdgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentPositionDefinitionGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := epdgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (epdgb *EquipmentPositionDefinitionGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := epdgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (epdgb *EquipmentPositionDefinitionGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(epdgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentPositionDefinitionGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := epdgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (epdgb *EquipmentPositionDefinitionGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := epdgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (epdgb *EquipmentPositionDefinitionGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := epdgb.sqlQuery().Query()
	if err := epdgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (epdgb *EquipmentPositionDefinitionGroupBy) sqlQuery() *sql.Selector {
	selector := epdgb.sql
	columns := make([]string, 0, len(epdgb.fields)+len(epdgb.fns))
	columns = append(columns, epdgb.fields...)
	for _, fn := range epdgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(epdgb.fields...)
}

// EquipmentPositionDefinitionSelect is the builder for select fields of EquipmentPositionDefinition entities.
type EquipmentPositionDefinitionSelect struct {
	config
	fields []string
	// intermediate queries.
	sql *sql.Selector
}

// Scan applies the selector query and scan the result into the given value.
func (epds *EquipmentPositionDefinitionSelect) Scan(ctx context.Context, v interface{}) error {
	return epds.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (epds *EquipmentPositionDefinitionSelect) ScanX(ctx context.Context, v interface{}) {
	if err := epds.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (epds *EquipmentPositionDefinitionSelect) Strings(ctx context.Context) ([]string, error) {
	if len(epds.fields) > 1 {
		return nil, errors.New("ent: EquipmentPositionDefinitionSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := epds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (epds *EquipmentPositionDefinitionSelect) StringsX(ctx context.Context) []string {
	v, err := epds.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (epds *EquipmentPositionDefinitionSelect) Ints(ctx context.Context) ([]int, error) {
	if len(epds.fields) > 1 {
		return nil, errors.New("ent: EquipmentPositionDefinitionSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := epds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (epds *EquipmentPositionDefinitionSelect) IntsX(ctx context.Context) []int {
	v, err := epds.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (epds *EquipmentPositionDefinitionSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(epds.fields) > 1 {
		return nil, errors.New("ent: EquipmentPositionDefinitionSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := epds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (epds *EquipmentPositionDefinitionSelect) Float64sX(ctx context.Context) []float64 {
	v, err := epds.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (epds *EquipmentPositionDefinitionSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(epds.fields) > 1 {
		return nil, errors.New("ent: EquipmentPositionDefinitionSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := epds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (epds *EquipmentPositionDefinitionSelect) BoolsX(ctx context.Context) []bool {
	v, err := epds.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (epds *EquipmentPositionDefinitionSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := epds.sqlQuery().Query()
	if err := epds.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (epds *EquipmentPositionDefinitionSelect) sqlQuery() sql.Querier {
	selector := epds.sql
	selector.Select(selector.Columns(epds.fields...)...)
	return selector
}
