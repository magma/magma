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
	"strconv"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
)

// EquipmentPortQuery is the builder for querying EquipmentPort entities.
type EquipmentPortQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.EquipmentPort
	// eager-loading edges.
	withDefinition *EquipmentPortDefinitionQuery
	withParent     *EquipmentQuery
	withLink       *LinkQuery
	withProperties *PropertyQuery
	withEndpoints  *ServiceEndpointQuery
	withFKs        bool
	// intermediate query.
	sql *sql.Selector
}

// Where adds a new predicate for the builder.
func (epq *EquipmentPortQuery) Where(ps ...predicate.EquipmentPort) *EquipmentPortQuery {
	epq.predicates = append(epq.predicates, ps...)
	return epq
}

// Limit adds a limit step to the query.
func (epq *EquipmentPortQuery) Limit(limit int) *EquipmentPortQuery {
	epq.limit = &limit
	return epq
}

// Offset adds an offset step to the query.
func (epq *EquipmentPortQuery) Offset(offset int) *EquipmentPortQuery {
	epq.offset = &offset
	return epq
}

// Order adds an order step to the query.
func (epq *EquipmentPortQuery) Order(o ...Order) *EquipmentPortQuery {
	epq.order = append(epq.order, o...)
	return epq
}

// QueryDefinition chains the current query on the definition edge.
func (epq *EquipmentPortQuery) QueryDefinition() *EquipmentPortDefinitionQuery {
	query := &EquipmentPortDefinitionQuery{config: epq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentport.Table, equipmentport.FieldID, epq.sqlQuery()),
		sqlgraph.To(equipmentportdefinition.Table, equipmentportdefinition.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, equipmentport.DefinitionTable, equipmentport.DefinitionColumn),
	)
	query.sql = sqlgraph.SetNeighbors(epq.driver.Dialect(), step)
	return query
}

// QueryParent chains the current query on the parent edge.
func (epq *EquipmentPortQuery) QueryParent() *EquipmentQuery {
	query := &EquipmentQuery{config: epq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentport.Table, equipmentport.FieldID, epq.sqlQuery()),
		sqlgraph.To(equipment.Table, equipment.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, equipmentport.ParentTable, equipmentport.ParentColumn),
	)
	query.sql = sqlgraph.SetNeighbors(epq.driver.Dialect(), step)
	return query
}

// QueryLink chains the current query on the link edge.
func (epq *EquipmentPortQuery) QueryLink() *LinkQuery {
	query := &LinkQuery{config: epq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentport.Table, equipmentport.FieldID, epq.sqlQuery()),
		sqlgraph.To(link.Table, link.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, equipmentport.LinkTable, equipmentport.LinkColumn),
	)
	query.sql = sqlgraph.SetNeighbors(epq.driver.Dialect(), step)
	return query
}

// QueryProperties chains the current query on the properties edge.
func (epq *EquipmentPortQuery) QueryProperties() *PropertyQuery {
	query := &PropertyQuery{config: epq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentport.Table, equipmentport.FieldID, epq.sqlQuery()),
		sqlgraph.To(property.Table, property.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, equipmentport.PropertiesTable, equipmentport.PropertiesColumn),
	)
	query.sql = sqlgraph.SetNeighbors(epq.driver.Dialect(), step)
	return query
}

// QueryEndpoints chains the current query on the endpoints edge.
func (epq *EquipmentPortQuery) QueryEndpoints() *ServiceEndpointQuery {
	query := &ServiceEndpointQuery{config: epq.config}
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentport.Table, equipmentport.FieldID, epq.sqlQuery()),
		sqlgraph.To(serviceendpoint.Table, serviceendpoint.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, equipmentport.EndpointsTable, equipmentport.EndpointsColumn),
	)
	query.sql = sqlgraph.SetNeighbors(epq.driver.Dialect(), step)
	return query
}

// First returns the first EquipmentPort entity in the query. Returns *NotFoundError when no equipmentport was found.
func (epq *EquipmentPortQuery) First(ctx context.Context) (*EquipmentPort, error) {
	eps, err := epq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(eps) == 0 {
		return nil, &NotFoundError{equipmentport.Label}
	}
	return eps[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (epq *EquipmentPortQuery) FirstX(ctx context.Context) *EquipmentPort {
	ep, err := epq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return ep
}

// FirstID returns the first EquipmentPort id in the query. Returns *NotFoundError when no id was found.
func (epq *EquipmentPortQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = epq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{equipmentport.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (epq *EquipmentPortQuery) FirstXID(ctx context.Context) string {
	id, err := epq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only EquipmentPort entity in the query, returns an error if not exactly one entity was returned.
func (epq *EquipmentPortQuery) Only(ctx context.Context) (*EquipmentPort, error) {
	eps, err := epq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(eps) {
	case 1:
		return eps[0], nil
	case 0:
		return nil, &NotFoundError{equipmentport.Label}
	default:
		return nil, &NotSingularError{equipmentport.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (epq *EquipmentPortQuery) OnlyX(ctx context.Context) *EquipmentPort {
	ep, err := epq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return ep
}

// OnlyID returns the only EquipmentPort id in the query, returns an error if not exactly one id was returned.
func (epq *EquipmentPortQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = epq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{equipmentport.Label}
	default:
		err = &NotSingularError{equipmentport.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (epq *EquipmentPortQuery) OnlyXID(ctx context.Context) string {
	id, err := epq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of EquipmentPorts.
func (epq *EquipmentPortQuery) All(ctx context.Context) ([]*EquipmentPort, error) {
	return epq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (epq *EquipmentPortQuery) AllX(ctx context.Context) []*EquipmentPort {
	eps, err := epq.All(ctx)
	if err != nil {
		panic(err)
	}
	return eps
}

// IDs executes the query and returns a list of EquipmentPort ids.
func (epq *EquipmentPortQuery) IDs(ctx context.Context) ([]string, error) {
	var ids []string
	if err := epq.Select(equipmentport.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (epq *EquipmentPortQuery) IDsX(ctx context.Context) []string {
	ids, err := epq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (epq *EquipmentPortQuery) Count(ctx context.Context) (int, error) {
	return epq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (epq *EquipmentPortQuery) CountX(ctx context.Context) int {
	count, err := epq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (epq *EquipmentPortQuery) Exist(ctx context.Context) (bool, error) {
	return epq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (epq *EquipmentPortQuery) ExistX(ctx context.Context) bool {
	exist, err := epq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (epq *EquipmentPortQuery) Clone() *EquipmentPortQuery {
	return &EquipmentPortQuery{
		config:     epq.config,
		limit:      epq.limit,
		offset:     epq.offset,
		order:      append([]Order{}, epq.order...),
		unique:     append([]string{}, epq.unique...),
		predicates: append([]predicate.EquipmentPort{}, epq.predicates...),
		// clone intermediate query.
		sql: epq.sql.Clone(),
	}
}

//  WithDefinition tells the query-builder to eager-loads the nodes that are connected to
// the "definition" edge. The optional arguments used to configure the query builder of the edge.
func (epq *EquipmentPortQuery) WithDefinition(opts ...func(*EquipmentPortDefinitionQuery)) *EquipmentPortQuery {
	query := &EquipmentPortDefinitionQuery{config: epq.config}
	for _, opt := range opts {
		opt(query)
	}
	epq.withDefinition = query
	return epq
}

//  WithParent tells the query-builder to eager-loads the nodes that are connected to
// the "parent" edge. The optional arguments used to configure the query builder of the edge.
func (epq *EquipmentPortQuery) WithParent(opts ...func(*EquipmentQuery)) *EquipmentPortQuery {
	query := &EquipmentQuery{config: epq.config}
	for _, opt := range opts {
		opt(query)
	}
	epq.withParent = query
	return epq
}

//  WithLink tells the query-builder to eager-loads the nodes that are connected to
// the "link" edge. The optional arguments used to configure the query builder of the edge.
func (epq *EquipmentPortQuery) WithLink(opts ...func(*LinkQuery)) *EquipmentPortQuery {
	query := &LinkQuery{config: epq.config}
	for _, opt := range opts {
		opt(query)
	}
	epq.withLink = query
	return epq
}

//  WithProperties tells the query-builder to eager-loads the nodes that are connected to
// the "properties" edge. The optional arguments used to configure the query builder of the edge.
func (epq *EquipmentPortQuery) WithProperties(opts ...func(*PropertyQuery)) *EquipmentPortQuery {
	query := &PropertyQuery{config: epq.config}
	for _, opt := range opts {
		opt(query)
	}
	epq.withProperties = query
	return epq
}

//  WithEndpoints tells the query-builder to eager-loads the nodes that are connected to
// the "endpoints" edge. The optional arguments used to configure the query builder of the edge.
func (epq *EquipmentPortQuery) WithEndpoints(opts ...func(*ServiceEndpointQuery)) *EquipmentPortQuery {
	query := &ServiceEndpointQuery{config: epq.config}
	for _, opt := range opts {
		opt(query)
	}
	epq.withEndpoints = query
	return epq
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
//	client.EquipmentPort.Query().
//		GroupBy(equipmentport.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (epq *EquipmentPortQuery) GroupBy(field string, fields ...string) *EquipmentPortGroupBy {
	group := &EquipmentPortGroupBy{config: epq.config}
	group.fields = append([]string{field}, fields...)
	group.sql = epq.sqlQuery()
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
//	client.EquipmentPort.Query().
//		Select(equipmentport.FieldCreateTime).
//		Scan(ctx, &v)
//
func (epq *EquipmentPortQuery) Select(field string, fields ...string) *EquipmentPortSelect {
	selector := &EquipmentPortSelect{config: epq.config}
	selector.fields = append([]string{field}, fields...)
	selector.sql = epq.sqlQuery()
	return selector
}

func (epq *EquipmentPortQuery) sqlAll(ctx context.Context) ([]*EquipmentPort, error) {
	var (
		nodes       = []*EquipmentPort{}
		withFKs     = epq.withFKs
		_spec       = epq.querySpec()
		loadedTypes = [5]bool{
			epq.withDefinition != nil,
			epq.withParent != nil,
			epq.withLink != nil,
			epq.withProperties != nil,
			epq.withEndpoints != nil,
		}
	)
	if epq.withDefinition != nil || epq.withParent != nil || epq.withLink != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, equipmentport.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &EquipmentPort{config: epq.config}
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
	if err := sqlgraph.QueryNodes(ctx, epq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := epq.withDefinition; query != nil {
		ids := make([]string, 0, len(nodes))
		nodeids := make(map[string][]*EquipmentPort)
		for i := range nodes {
			if fk := nodes[i].definition_id; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(equipmentportdefinition.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "definition_id" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Definition = n
			}
		}
	}

	if query := epq.withParent; query != nil {
		ids := make([]string, 0, len(nodes))
		nodeids := make(map[string][]*EquipmentPort)
		for i := range nodes {
			if fk := nodes[i].parent_id; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(equipment.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "parent_id" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Parent = n
			}
		}
	}

	if query := epq.withLink; query != nil {
		ids := make([]string, 0, len(nodes))
		nodeids := make(map[string][]*EquipmentPort)
		for i := range nodes {
			if fk := nodes[i].link_id; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(link.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "link_id" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Link = n
			}
		}
	}

	if query := epq.withProperties; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[string]*EquipmentPort)
		for i := range nodes {
			id, err := strconv.Atoi(nodes[i].ID)
			if err != nil {
				return nil, err
			}
			fks = append(fks, id)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.Property(func(s *sql.Selector) {
			s.Where(sql.InValues(equipmentport.PropertiesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.equipment_port_id
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "equipment_port_id" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_port_id" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Properties = append(node.Edges.Properties, n)
		}
	}

	if query := epq.withEndpoints; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[string]*EquipmentPort)
		for i := range nodes {
			id, err := strconv.Atoi(nodes[i].ID)
			if err != nil {
				return nil, err
			}
			fks = append(fks, id)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.ServiceEndpoint(func(s *sql.Selector) {
			s.Where(sql.InValues(equipmentport.EndpointsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.port_id
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "port_id" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "port_id" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Endpoints = append(node.Edges.Endpoints, n)
		}
	}

	return nodes, nil
}

func (epq *EquipmentPortQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := epq.querySpec()
	return sqlgraph.CountNodes(ctx, epq.driver, _spec)
}

func (epq *EquipmentPortQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := epq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (epq *EquipmentPortQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentport.Table,
			Columns: equipmentport.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: equipmentport.FieldID,
			},
		},
		From:   epq.sql,
		Unique: true,
	}
	if ps := epq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := epq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := epq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := epq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (epq *EquipmentPortQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(epq.driver.Dialect())
	t1 := builder.Table(equipmentport.Table)
	selector := builder.Select(t1.Columns(equipmentport.Columns...)...).From(t1)
	if epq.sql != nil {
		selector = epq.sql
		selector.Select(selector.Columns(equipmentport.Columns...)...)
	}
	for _, p := range epq.predicates {
		p(selector)
	}
	for _, p := range epq.order {
		p(selector)
	}
	if offset := epq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := epq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// EquipmentPortGroupBy is the builder for group-by EquipmentPort entities.
type EquipmentPortGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query.
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the group-by query.
func (epgb *EquipmentPortGroupBy) Aggregate(fns ...Aggregate) *EquipmentPortGroupBy {
	epgb.fns = append(epgb.fns, fns...)
	return epgb
}

// Scan applies the group-by query and scan the result into the given value.
func (epgb *EquipmentPortGroupBy) Scan(ctx context.Context, v interface{}) error {
	return epgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (epgb *EquipmentPortGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := epgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (epgb *EquipmentPortGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(epgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := epgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (epgb *EquipmentPortGroupBy) StringsX(ctx context.Context) []string {
	v, err := epgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (epgb *EquipmentPortGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(epgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := epgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (epgb *EquipmentPortGroupBy) IntsX(ctx context.Context) []int {
	v, err := epgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (epgb *EquipmentPortGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(epgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := epgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (epgb *EquipmentPortGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := epgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (epgb *EquipmentPortGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(epgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := epgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (epgb *EquipmentPortGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := epgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (epgb *EquipmentPortGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := epgb.sqlQuery().Query()
	if err := epgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (epgb *EquipmentPortGroupBy) sqlQuery() *sql.Selector {
	selector := epgb.sql
	columns := make([]string, 0, len(epgb.fields)+len(epgb.fns))
	columns = append(columns, epgb.fields...)
	for _, fn := range epgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(epgb.fields...)
}

// EquipmentPortSelect is the builder for select fields of EquipmentPort entities.
type EquipmentPortSelect struct {
	config
	fields []string
	// intermediate queries.
	sql *sql.Selector
}

// Scan applies the selector query and scan the result into the given value.
func (eps *EquipmentPortSelect) Scan(ctx context.Context, v interface{}) error {
	return eps.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (eps *EquipmentPortSelect) ScanX(ctx context.Context, v interface{}) {
	if err := eps.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (eps *EquipmentPortSelect) Strings(ctx context.Context) ([]string, error) {
	if len(eps.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := eps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (eps *EquipmentPortSelect) StringsX(ctx context.Context) []string {
	v, err := eps.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (eps *EquipmentPortSelect) Ints(ctx context.Context) ([]int, error) {
	if len(eps.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := eps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (eps *EquipmentPortSelect) IntsX(ctx context.Context) []int {
	v, err := eps.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (eps *EquipmentPortSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(eps.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := eps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (eps *EquipmentPortSelect) Float64sX(ctx context.Context) []float64 {
	v, err := eps.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (eps *EquipmentPortSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(eps.fields) > 1 {
		return nil, errors.New("ent: EquipmentPortSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := eps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (eps *EquipmentPortSelect) BoolsX(ctx context.Context) []bool {
	v, err := eps.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (eps *EquipmentPortSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := eps.sqlQuery().Query()
	if err := eps.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (eps *EquipmentPortSelect) sqlQuery() sql.Querier {
	selector := eps.sql
	selector.Select(selector.Columns(eps.fields...)...)
	return selector
}
