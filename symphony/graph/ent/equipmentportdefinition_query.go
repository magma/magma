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
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// EquipmentPortDefinitionQuery is the builder for querying EquipmentPortDefinition entities.
type EquipmentPortDefinitionQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.EquipmentPortDefinition
	// intermediate query.
	sql *sql.Selector
}

// Where adds a new predicate for the builder.
func (epdq *EquipmentPortDefinitionQuery) Where(ps ...predicate.EquipmentPortDefinition) *EquipmentPortDefinitionQuery {
	epdq.predicates = append(epdq.predicates, ps...)
	return epdq
}

// Limit adds a limit step to the query.
func (epdq *EquipmentPortDefinitionQuery) Limit(limit int) *EquipmentPortDefinitionQuery {
	epdq.limit = &limit
	return epdq
}

// Offset adds an offset step to the query.
func (epdq *EquipmentPortDefinitionQuery) Offset(offset int) *EquipmentPortDefinitionQuery {
	epdq.offset = &offset
	return epdq
}

// Order adds an order step to the query.
func (epdq *EquipmentPortDefinitionQuery) Order(o ...Order) *EquipmentPortDefinitionQuery {
	epdq.order = append(epdq.order, o...)
	return epdq
}

// QueryEquipmentPortType chains the current query on the equipment_port_type edge.
func (epdq *EquipmentPortDefinitionQuery) QueryEquipmentPortType() *EquipmentPortTypeQuery {
	query := &EquipmentPortTypeQuery{config: epdq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentportdefinition.Table, equipmentportdefinition.FieldID, epdq.sqlQuery()),
		sqlgraph.To(equipmentporttype.Table, equipmentporttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, equipmentportdefinition.EquipmentPortTypeTable, equipmentportdefinition.EquipmentPortTypeColumn),
	)
	query.sql = sqlgraph.SetNeighbors(epdq.driver.Dialect(), step)
	return query
}

// QueryPorts chains the current query on the ports edge.
func (epdq *EquipmentPortDefinitionQuery) QueryPorts() *EquipmentPortQuery {
	query := &EquipmentPortQuery{config: epdq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentportdefinition.Table, equipmentportdefinition.FieldID, epdq.sqlQuery()),
		sqlgraph.To(equipmentport.Table, equipmentport.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, equipmentportdefinition.PortsTable, equipmentportdefinition.PortsColumn),
	)
	query.sql = sqlgraph.SetNeighbors(epdq.driver.Dialect(), step)
	return query
}

// QueryEquipmentType chains the current query on the equipment_type edge.
func (epdq *EquipmentPortDefinitionQuery) QueryEquipmentType() *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: epdq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentportdefinition.Table, equipmentportdefinition.FieldID, epdq.sqlQuery()),
		sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, equipmentportdefinition.EquipmentTypeTable, equipmentportdefinition.EquipmentTypeColumn),
	)
	query.sql = sqlgraph.SetNeighbors(epdq.driver.Dialect(), step)
	return query
}

// First returns the first EquipmentPortDefinition entity in the query. Returns *ErrNotFound when no equipmentportdefinition was found.
func (epdq *EquipmentPortDefinitionQuery) First(ctx context.Context) (*EquipmentPortDefinition, error) {
	epds, err := epdq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(epds) == 0 {
		return nil, &ErrNotFound{equipmentportdefinition.Label}
	}
	return epds[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (epdq *EquipmentPortDefinitionQuery) FirstX(ctx context.Context) *EquipmentPortDefinition {
	epd, err := epdq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return epd
}

// FirstID returns the first EquipmentPortDefinition id in the query. Returns *ErrNotFound when no id was found.
func (epdq *EquipmentPortDefinitionQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = epdq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &ErrNotFound{equipmentportdefinition.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (epdq *EquipmentPortDefinitionQuery) FirstXID(ctx context.Context) string {
	id, err := epdq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only EquipmentPortDefinition entity in the query, returns an error if not exactly one entity was returned.
func (epdq *EquipmentPortDefinitionQuery) Only(ctx context.Context) (*EquipmentPortDefinition, error) {
	epds, err := epdq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(epds) {
	case 1:
		return epds[0], nil
	case 0:
		return nil, &ErrNotFound{equipmentportdefinition.Label}
	default:
		return nil, &ErrNotSingular{equipmentportdefinition.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (epdq *EquipmentPortDefinitionQuery) OnlyX(ctx context.Context) *EquipmentPortDefinition {
	epd, err := epdq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return epd
}

// OnlyID returns the only EquipmentPortDefinition id in the query, returns an error if not exactly one id was returned.
func (epdq *EquipmentPortDefinitionQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = epdq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &ErrNotFound{equipmentportdefinition.Label}
	default:
		err = &ErrNotSingular{equipmentportdefinition.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (epdq *EquipmentPortDefinitionQuery) OnlyXID(ctx context.Context) string {
	id, err := epdq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of EquipmentPortDefinitions.
func (epdq *EquipmentPortDefinitionQuery) All(ctx context.Context) ([]*EquipmentPortDefinition, error) {
	return epdq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (epdq *EquipmentPortDefinitionQuery) AllX(ctx context.Context) []*EquipmentPortDefinition {
	epds, err := epdq.All(ctx)
	if err != nil {
		panic(err)
	}
	return epds
}

// IDs executes the query and returns a list of EquipmentPortDefinition ids.
func (epdq *EquipmentPortDefinitionQuery) IDs(ctx context.Context) ([]string, error) {
	var ids []string
	if err := epdq.Select(equipmentportdefinition.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (epdq *EquipmentPortDefinitionQuery) IDsX(ctx context.Context) []string {
	ids, err := epdq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (epdq *EquipmentPortDefinitionQuery) Count(ctx context.Context) (int, error) {
	return epdq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (epdq *EquipmentPortDefinitionQuery) CountX(ctx context.Context) int {
	count, err := epdq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (epdq *EquipmentPortDefinitionQuery) Exist(ctx context.Context) (bool, error) {
	return epdq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (epdq *EquipmentPortDefinitionQuery) ExistX(ctx context.Context) bool {
	exist, err := epdq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (epdq *EquipmentPortDefinitionQuery) Clone() *EquipmentPortDefinitionQuery {
	return &EquipmentPortDefinitionQuery{
		config:     epdq.config,
		limit:      epdq.limit,
		offset:     epdq.offset,
		order:      append([]Order{}, epdq.order...),
		unique:     append([]string{}, epdq.unique...),
		predicates: append([]predicate.EquipmentPortDefinition{}, epdq.predicates...),
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
//	client.EquipmentPortDefinition.Query().
//		GroupBy(equipmentportdefinition.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (epdq *EquipmentPortDefinitionQuery) GroupBy(field string, fields ...string) *EquipmentPortDefinitionGroupBy {
	group := &EquipmentPortDefinitionGroupBy{config: epdq.config}
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
//	client.EquipmentPortDefinition.Query().
//		Select(equipmentportdefinition.FieldCreateTime).
//		Scan(ctx, &v)
//
func (epdq *EquipmentPortDefinitionQuery) Select(field string, fields ...string) *EquipmentPortDefinitionSelect {
	selector := &EquipmentPortDefinitionSelect{config: epdq.config}
	selector.fields = append([]string{field}, fields...)
	selector.sql = epdq.sqlQuery()
	return selector
}

func (epdq *EquipmentPortDefinitionQuery) sqlAll(ctx context.Context) ([]*EquipmentPortDefinition, error) {
	var (
		nodes []*EquipmentPortDefinition
		spec  = epdq.querySpec()
	)
	spec.ScanValues = func() []interface{} {
		node := &EquipmentPortDefinition{config: epdq.config}
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

func (epdq *EquipmentPortDefinitionQuery) sqlCount(ctx context.Context) (int, error) {
	spec := epdq.querySpec()
	return sqlgraph.CountNodes(ctx, epdq.driver, spec)
}

func (epdq *EquipmentPortDefinitionQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := epdq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (epdq *EquipmentPortDefinitionQuery) querySpec() *sqlgraph.QuerySpec {
	spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentportdefinition.Table,
			Columns: equipmentportdefinition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: equipmentportdefinition.FieldID,
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

func (epdq *EquipmentPortDefinitionQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(epdq.driver.Dialect())
	t1 := builder.Table(equipmentportdefinition.Table)
	selector := builder.Select(t1.Columns(equipmentportdefinition.Columns...)...).From(t1)
	if epdq.sql != nil {
		selector = epdq.sql
		selector.Select(selector.Columns(equipmentportdefinition.Columns...)...)
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

// EquipmentPortDefinitionGroupBy is the builder for group-by EquipmentPortDefinition entities.
type EquipmentPortDefinitionGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query.
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the group-by query.
func (epdgb *EquipmentPortDefinitionGroupBy) Aggregate(fns ...Aggregate) *EquipmentPortDefinitionGroupBy {
	epdgb.fns = append(epdgb.fns, fns...)
	return epdgb
}

// Scan applies the group-by query and scan the result into the given value.
func (epdgb *EquipmentPortDefinitionGroupBy) Scan(ctx context.Context, v interface{}) error {
	return epdgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (epdgb *EquipmentPortDefinitionGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := epdgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (epdgb *EquipmentPortDefinitionGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(epdgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortDefinitionGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := epdgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (epdgb *EquipmentPortDefinitionGroupBy) StringsX(ctx context.Context) []string {
	v, err := epdgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (epdgb *EquipmentPortDefinitionGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(epdgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortDefinitionGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := epdgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (epdgb *EquipmentPortDefinitionGroupBy) IntsX(ctx context.Context) []int {
	v, err := epdgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (epdgb *EquipmentPortDefinitionGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(epdgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortDefinitionGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := epdgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (epdgb *EquipmentPortDefinitionGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := epdgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (epdgb *EquipmentPortDefinitionGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(epdgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortDefinitionGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := epdgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (epdgb *EquipmentPortDefinitionGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := epdgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (epdgb *EquipmentPortDefinitionGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := epdgb.sqlQuery().Query()
	if err := epdgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (epdgb *EquipmentPortDefinitionGroupBy) sqlQuery() *sql.Selector {
	selector := epdgb.sql
	columns := make([]string, 0, len(epdgb.fields)+len(epdgb.fns))
	columns = append(columns, epdgb.fields...)
	for _, fn := range epdgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(epdgb.fields...)
}

// EquipmentPortDefinitionSelect is the builder for select fields of EquipmentPortDefinition entities.
type EquipmentPortDefinitionSelect struct {
	config
	fields []string
	// intermediate queries.
	sql *sql.Selector
}

// Scan applies the selector query and scan the result into the given value.
func (epds *EquipmentPortDefinitionSelect) Scan(ctx context.Context, v interface{}) error {
	return epds.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (epds *EquipmentPortDefinitionSelect) ScanX(ctx context.Context, v interface{}) {
	if err := epds.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (epds *EquipmentPortDefinitionSelect) Strings(ctx context.Context) ([]string, error) {
	if len(epds.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortDefinitionSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := epds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (epds *EquipmentPortDefinitionSelect) StringsX(ctx context.Context) []string {
	v, err := epds.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (epds *EquipmentPortDefinitionSelect) Ints(ctx context.Context) ([]int, error) {
	if len(epds.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortDefinitionSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := epds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (epds *EquipmentPortDefinitionSelect) IntsX(ctx context.Context) []int {
	v, err := epds.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (epds *EquipmentPortDefinitionSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(epds.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortDefinitionSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := epds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (epds *EquipmentPortDefinitionSelect) Float64sX(ctx context.Context) []float64 {
	v, err := epds.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (epds *EquipmentPortDefinitionSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(epds.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortDefinitionSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := epds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (epds *EquipmentPortDefinitionSelect) BoolsX(ctx context.Context) []bool {
	v, err := epds.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (epds *EquipmentPortDefinitionSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := epds.sqlQuery().Query()
	if err := epds.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (epds *EquipmentPortDefinitionSelect) sqlQuery() sql.Querier {
	selector := epds.sql
	selector.Select(selector.Columns(epds.fields...)...)
	return selector
}
