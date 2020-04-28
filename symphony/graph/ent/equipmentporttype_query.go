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
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
)

// EquipmentPortTypeQuery is the builder for querying EquipmentPortType entities.
type EquipmentPortTypeQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.EquipmentPortType
	// eager-loading edges.
	withPropertyTypes     *PropertyTypeQuery
	withLinkPropertyTypes *PropertyTypeQuery
	withPortDefinitions   *EquipmentPortDefinitionQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (eptq *EquipmentPortTypeQuery) Where(ps ...predicate.EquipmentPortType) *EquipmentPortTypeQuery {
	eptq.predicates = append(eptq.predicates, ps...)
	return eptq
}

// Limit adds a limit step to the query.
func (eptq *EquipmentPortTypeQuery) Limit(limit int) *EquipmentPortTypeQuery {
	eptq.limit = &limit
	return eptq
}

// Offset adds an offset step to the query.
func (eptq *EquipmentPortTypeQuery) Offset(offset int) *EquipmentPortTypeQuery {
	eptq.offset = &offset
	return eptq
}

// Order adds an order step to the query.
func (eptq *EquipmentPortTypeQuery) Order(o ...Order) *EquipmentPortTypeQuery {
	eptq.order = append(eptq.order, o...)
	return eptq
}

// QueryPropertyTypes chains the current query on the property_types edge.
func (eptq *EquipmentPortTypeQuery) QueryPropertyTypes() *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: eptq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentporttype.Table, equipmentporttype.FieldID, eptq.sqlQuery()),
			sqlgraph.To(propertytype.Table, propertytype.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipmentporttype.PropertyTypesTable, equipmentporttype.PropertyTypesColumn),
		)
		fromU = sqlgraph.SetNeighbors(eptq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryLinkPropertyTypes chains the current query on the link_property_types edge.
func (eptq *EquipmentPortTypeQuery) QueryLinkPropertyTypes() *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: eptq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentporttype.Table, equipmentporttype.FieldID, eptq.sqlQuery()),
			sqlgraph.To(propertytype.Table, propertytype.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipmentporttype.LinkPropertyTypesTable, equipmentporttype.LinkPropertyTypesColumn),
		)
		fromU = sqlgraph.SetNeighbors(eptq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryPortDefinitions chains the current query on the port_definitions edge.
func (eptq *EquipmentPortTypeQuery) QueryPortDefinitions() *EquipmentPortDefinitionQuery {
	query := &EquipmentPortDefinitionQuery{config: eptq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentporttype.Table, equipmentporttype.FieldID, eptq.sqlQuery()),
			sqlgraph.To(equipmentportdefinition.Table, equipmentportdefinition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, equipmentporttype.PortDefinitionsTable, equipmentporttype.PortDefinitionsColumn),
		)
		fromU = sqlgraph.SetNeighbors(eptq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first EquipmentPortType entity in the query. Returns *NotFoundError when no equipmentporttype was found.
func (eptq *EquipmentPortTypeQuery) First(ctx context.Context) (*EquipmentPortType, error) {
	epts, err := eptq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(epts) == 0 {
		return nil, &NotFoundError{equipmentporttype.Label}
	}
	return epts[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (eptq *EquipmentPortTypeQuery) FirstX(ctx context.Context) *EquipmentPortType {
	ept, err := eptq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return ept
}

// FirstID returns the first EquipmentPortType id in the query. Returns *NotFoundError when no id was found.
func (eptq *EquipmentPortTypeQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = eptq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{equipmentporttype.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (eptq *EquipmentPortTypeQuery) FirstXID(ctx context.Context) int {
	id, err := eptq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only EquipmentPortType entity in the query, returns an error if not exactly one entity was returned.
func (eptq *EquipmentPortTypeQuery) Only(ctx context.Context) (*EquipmentPortType, error) {
	epts, err := eptq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(epts) {
	case 1:
		return epts[0], nil
	case 0:
		return nil, &NotFoundError{equipmentporttype.Label}
	default:
		return nil, &NotSingularError{equipmentporttype.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (eptq *EquipmentPortTypeQuery) OnlyX(ctx context.Context) *EquipmentPortType {
	ept, err := eptq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return ept
}

// OnlyID returns the only EquipmentPortType id in the query, returns an error if not exactly one id was returned.
func (eptq *EquipmentPortTypeQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = eptq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{equipmentporttype.Label}
	default:
		err = &NotSingularError{equipmentporttype.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (eptq *EquipmentPortTypeQuery) OnlyXID(ctx context.Context) int {
	id, err := eptq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of EquipmentPortTypes.
func (eptq *EquipmentPortTypeQuery) All(ctx context.Context) ([]*EquipmentPortType, error) {
	if err := eptq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return eptq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (eptq *EquipmentPortTypeQuery) AllX(ctx context.Context) []*EquipmentPortType {
	epts, err := eptq.All(ctx)
	if err != nil {
		panic(err)
	}
	return epts
}

// IDs executes the query and returns a list of EquipmentPortType ids.
func (eptq *EquipmentPortTypeQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := eptq.Select(equipmentporttype.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (eptq *EquipmentPortTypeQuery) IDsX(ctx context.Context) []int {
	ids, err := eptq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (eptq *EquipmentPortTypeQuery) Count(ctx context.Context) (int, error) {
	if err := eptq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return eptq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (eptq *EquipmentPortTypeQuery) CountX(ctx context.Context) int {
	count, err := eptq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (eptq *EquipmentPortTypeQuery) Exist(ctx context.Context) (bool, error) {
	if err := eptq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return eptq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (eptq *EquipmentPortTypeQuery) ExistX(ctx context.Context) bool {
	exist, err := eptq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (eptq *EquipmentPortTypeQuery) Clone() *EquipmentPortTypeQuery {
	return &EquipmentPortTypeQuery{
		config:     eptq.config,
		limit:      eptq.limit,
		offset:     eptq.offset,
		order:      append([]Order{}, eptq.order...),
		unique:     append([]string{}, eptq.unique...),
		predicates: append([]predicate.EquipmentPortType{}, eptq.predicates...),
		// clone intermediate query.
		sql:  eptq.sql.Clone(),
		path: eptq.path,
	}
}

//  WithPropertyTypes tells the query-builder to eager-loads the nodes that are connected to
// the "property_types" edge. The optional arguments used to configure the query builder of the edge.
func (eptq *EquipmentPortTypeQuery) WithPropertyTypes(opts ...func(*PropertyTypeQuery)) *EquipmentPortTypeQuery {
	query := &PropertyTypeQuery{config: eptq.config}
	for _, opt := range opts {
		opt(query)
	}
	eptq.withPropertyTypes = query
	return eptq
}

//  WithLinkPropertyTypes tells the query-builder to eager-loads the nodes that are connected to
// the "link_property_types" edge. The optional arguments used to configure the query builder of the edge.
func (eptq *EquipmentPortTypeQuery) WithLinkPropertyTypes(opts ...func(*PropertyTypeQuery)) *EquipmentPortTypeQuery {
	query := &PropertyTypeQuery{config: eptq.config}
	for _, opt := range opts {
		opt(query)
	}
	eptq.withLinkPropertyTypes = query
	return eptq
}

//  WithPortDefinitions tells the query-builder to eager-loads the nodes that are connected to
// the "port_definitions" edge. The optional arguments used to configure the query builder of the edge.
func (eptq *EquipmentPortTypeQuery) WithPortDefinitions(opts ...func(*EquipmentPortDefinitionQuery)) *EquipmentPortTypeQuery {
	query := &EquipmentPortDefinitionQuery{config: eptq.config}
	for _, opt := range opts {
		opt(query)
	}
	eptq.withPortDefinitions = query
	return eptq
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
//	client.EquipmentPortType.Query().
//		GroupBy(equipmentporttype.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (eptq *EquipmentPortTypeQuery) GroupBy(field string, fields ...string) *EquipmentPortTypeGroupBy {
	group := &EquipmentPortTypeGroupBy{config: eptq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := eptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return eptq.sqlQuery(), nil
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
//	client.EquipmentPortType.Query().
//		Select(equipmentporttype.FieldCreateTime).
//		Scan(ctx, &v)
//
func (eptq *EquipmentPortTypeQuery) Select(field string, fields ...string) *EquipmentPortTypeSelect {
	selector := &EquipmentPortTypeSelect{config: eptq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := eptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return eptq.sqlQuery(), nil
	}
	return selector
}

func (eptq *EquipmentPortTypeQuery) prepareQuery(ctx context.Context) error {
	if eptq.path != nil {
		prev, err := eptq.path(ctx)
		if err != nil {
			return err
		}
		eptq.sql = prev
	}
	if err := equipmentporttype.Policy.EvalQuery(ctx, eptq); err != nil {
		return err
	}
	return nil
}

func (eptq *EquipmentPortTypeQuery) sqlAll(ctx context.Context) ([]*EquipmentPortType, error) {
	var (
		nodes       = []*EquipmentPortType{}
		_spec       = eptq.querySpec()
		loadedTypes = [3]bool{
			eptq.withPropertyTypes != nil,
			eptq.withLinkPropertyTypes != nil,
			eptq.withPortDefinitions != nil,
		}
	)
	_spec.ScanValues = func() []interface{} {
		node := &EquipmentPortType{config: eptq.config}
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
	if err := sqlgraph.QueryNodes(ctx, eptq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := eptq.withPropertyTypes; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*EquipmentPortType)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.PropertyType(func(s *sql.Selector) {
			s.Where(sql.InValues(equipmentporttype.PropertyTypesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.equipment_port_type_property_types
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "equipment_port_type_property_types" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_port_type_property_types" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.PropertyTypes = append(node.Edges.PropertyTypes, n)
		}
	}

	if query := eptq.withLinkPropertyTypes; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*EquipmentPortType)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.PropertyType(func(s *sql.Selector) {
			s.Where(sql.InValues(equipmentporttype.LinkPropertyTypesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.equipment_port_type_link_property_types
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "equipment_port_type_link_property_types" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_port_type_link_property_types" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.LinkPropertyTypes = append(node.Edges.LinkPropertyTypes, n)
		}
	}

	if query := eptq.withPortDefinitions; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*EquipmentPortType)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.EquipmentPortDefinition(func(s *sql.Selector) {
			s.Where(sql.InValues(equipmentporttype.PortDefinitionsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.equipment_port_definition_equipment_port_type
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "equipment_port_definition_equipment_port_type" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_port_definition_equipment_port_type" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.PortDefinitions = append(node.Edges.PortDefinitions, n)
		}
	}

	return nodes, nil
}

func (eptq *EquipmentPortTypeQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := eptq.querySpec()
	return sqlgraph.CountNodes(ctx, eptq.driver, _spec)
}

func (eptq *EquipmentPortTypeQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := eptq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (eptq *EquipmentPortTypeQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentporttype.Table,
			Columns: equipmentporttype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentporttype.FieldID,
			},
		},
		From:   eptq.sql,
		Unique: true,
	}
	if ps := eptq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := eptq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := eptq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := eptq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (eptq *EquipmentPortTypeQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(eptq.driver.Dialect())
	t1 := builder.Table(equipmentporttype.Table)
	selector := builder.Select(t1.Columns(equipmentporttype.Columns...)...).From(t1)
	if eptq.sql != nil {
		selector = eptq.sql
		selector.Select(selector.Columns(equipmentporttype.Columns...)...)
	}
	for _, p := range eptq.predicates {
		p(selector)
	}
	for _, p := range eptq.order {
		p(selector)
	}
	if offset := eptq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := eptq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// EquipmentPortTypeGroupBy is the builder for group-by EquipmentPortType entities.
type EquipmentPortTypeGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (eptgb *EquipmentPortTypeGroupBy) Aggregate(fns ...Aggregate) *EquipmentPortTypeGroupBy {
	eptgb.fns = append(eptgb.fns, fns...)
	return eptgb
}

// Scan applies the group-by query and scan the result into the given value.
func (eptgb *EquipmentPortTypeGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := eptgb.path(ctx)
	if err != nil {
		return err
	}
	eptgb.sql = query
	return eptgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (eptgb *EquipmentPortTypeGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := eptgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (eptgb *EquipmentPortTypeGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(eptgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortTypeGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := eptgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (eptgb *EquipmentPortTypeGroupBy) StringsX(ctx context.Context) []string {
	v, err := eptgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (eptgb *EquipmentPortTypeGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(eptgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortTypeGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := eptgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (eptgb *EquipmentPortTypeGroupBy) IntsX(ctx context.Context) []int {
	v, err := eptgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (eptgb *EquipmentPortTypeGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(eptgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortTypeGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := eptgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (eptgb *EquipmentPortTypeGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := eptgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (eptgb *EquipmentPortTypeGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(eptgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortTypeGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := eptgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (eptgb *EquipmentPortTypeGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := eptgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (eptgb *EquipmentPortTypeGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := eptgb.sqlQuery().Query()
	if err := eptgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (eptgb *EquipmentPortTypeGroupBy) sqlQuery() *sql.Selector {
	selector := eptgb.sql
	columns := make([]string, 0, len(eptgb.fields)+len(eptgb.fns))
	columns = append(columns, eptgb.fields...)
	for _, fn := range eptgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(eptgb.fields...)
}

// EquipmentPortTypeSelect is the builder for select fields of EquipmentPortType entities.
type EquipmentPortTypeSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (epts *EquipmentPortTypeSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := epts.path(ctx)
	if err != nil {
		return err
	}
	epts.sql = query
	return epts.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (epts *EquipmentPortTypeSelect) ScanX(ctx context.Context, v interface{}) {
	if err := epts.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (epts *EquipmentPortTypeSelect) Strings(ctx context.Context) ([]string, error) {
	if len(epts.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortTypeSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := epts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (epts *EquipmentPortTypeSelect) StringsX(ctx context.Context) []string {
	v, err := epts.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (epts *EquipmentPortTypeSelect) Ints(ctx context.Context) ([]int, error) {
	if len(epts.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortTypeSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := epts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (epts *EquipmentPortTypeSelect) IntsX(ctx context.Context) []int {
	v, err := epts.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (epts *EquipmentPortTypeSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(epts.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortTypeSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := epts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (epts *EquipmentPortTypeSelect) Float64sX(ctx context.Context) []float64 {
	v, err := epts.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (epts *EquipmentPortTypeSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(epts.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortTypeSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := epts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (epts *EquipmentPortTypeSelect) BoolsX(ctx context.Context) []bool {
	v, err := epts.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (epts *EquipmentPortTypeSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := epts.sqlQuery().Query()
	if err := epts.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (epts *EquipmentPortTypeSelect) sqlQuery() sql.Querier {
	selector := epts.sql
	selector.Select(selector.Columns(epts.fields...)...)
	return selector
}
