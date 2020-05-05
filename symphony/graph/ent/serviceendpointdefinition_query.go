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
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpointdefinition"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
)

// ServiceEndpointDefinitionQuery is the builder for querying ServiceEndpointDefinition entities.
type ServiceEndpointDefinitionQuery struct {
	config
	limit      *int
	offset     *int
	order      []OrderFunc
	unique     []string
	predicates []predicate.ServiceEndpointDefinition
	// eager-loading edges.
	withEndpoints     *ServiceEndpointQuery
	withServiceType   *ServiceTypeQuery
	withEquipmentType *EquipmentTypeQuery
	withFKs           bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (sedq *ServiceEndpointDefinitionQuery) Where(ps ...predicate.ServiceEndpointDefinition) *ServiceEndpointDefinitionQuery {
	sedq.predicates = append(sedq.predicates, ps...)
	return sedq
}

// Limit adds a limit step to the query.
func (sedq *ServiceEndpointDefinitionQuery) Limit(limit int) *ServiceEndpointDefinitionQuery {
	sedq.limit = &limit
	return sedq
}

// Offset adds an offset step to the query.
func (sedq *ServiceEndpointDefinitionQuery) Offset(offset int) *ServiceEndpointDefinitionQuery {
	sedq.offset = &offset
	return sedq
}

// Order adds an order step to the query.
func (sedq *ServiceEndpointDefinitionQuery) Order(o ...OrderFunc) *ServiceEndpointDefinitionQuery {
	sedq.order = append(sedq.order, o...)
	return sedq
}

// QueryEndpoints chains the current query on the endpoints edge.
func (sedq *ServiceEndpointDefinitionQuery) QueryEndpoints() *ServiceEndpointQuery {
	query := &ServiceEndpointQuery{config: sedq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := sedq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(serviceendpointdefinition.Table, serviceendpointdefinition.FieldID, sedq.sqlQuery()),
			sqlgraph.To(serviceendpoint.Table, serviceendpoint.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, serviceendpointdefinition.EndpointsTable, serviceendpointdefinition.EndpointsColumn),
		)
		fromU = sqlgraph.SetNeighbors(sedq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryServiceType chains the current query on the service_type edge.
func (sedq *ServiceEndpointDefinitionQuery) QueryServiceType() *ServiceTypeQuery {
	query := &ServiceTypeQuery{config: sedq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := sedq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(serviceendpointdefinition.Table, serviceendpointdefinition.FieldID, sedq.sqlQuery()),
			sqlgraph.To(servicetype.Table, servicetype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, serviceendpointdefinition.ServiceTypeTable, serviceendpointdefinition.ServiceTypeColumn),
		)
		fromU = sqlgraph.SetNeighbors(sedq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryEquipmentType chains the current query on the equipment_type edge.
func (sedq *ServiceEndpointDefinitionQuery) QueryEquipmentType() *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: sedq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := sedq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(serviceendpointdefinition.Table, serviceendpointdefinition.FieldID, sedq.sqlQuery()),
			sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, serviceendpointdefinition.EquipmentTypeTable, serviceendpointdefinition.EquipmentTypeColumn),
		)
		fromU = sqlgraph.SetNeighbors(sedq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first ServiceEndpointDefinition entity in the query. Returns *NotFoundError when no serviceendpointdefinition was found.
func (sedq *ServiceEndpointDefinitionQuery) First(ctx context.Context) (*ServiceEndpointDefinition, error) {
	seds, err := sedq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(seds) == 0 {
		return nil, &NotFoundError{serviceendpointdefinition.Label}
	}
	return seds[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (sedq *ServiceEndpointDefinitionQuery) FirstX(ctx context.Context) *ServiceEndpointDefinition {
	sed, err := sedq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return sed
}

// FirstID returns the first ServiceEndpointDefinition id in the query. Returns *NotFoundError when no id was found.
func (sedq *ServiceEndpointDefinitionQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = sedq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{serviceendpointdefinition.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (sedq *ServiceEndpointDefinitionQuery) FirstXID(ctx context.Context) int {
	id, err := sedq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only ServiceEndpointDefinition entity in the query, returns an error if not exactly one entity was returned.
func (sedq *ServiceEndpointDefinitionQuery) Only(ctx context.Context) (*ServiceEndpointDefinition, error) {
	seds, err := sedq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(seds) {
	case 1:
		return seds[0], nil
	case 0:
		return nil, &NotFoundError{serviceendpointdefinition.Label}
	default:
		return nil, &NotSingularError{serviceendpointdefinition.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (sedq *ServiceEndpointDefinitionQuery) OnlyX(ctx context.Context) *ServiceEndpointDefinition {
	sed, err := sedq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return sed
}

// OnlyID returns the only ServiceEndpointDefinition id in the query, returns an error if not exactly one id was returned.
func (sedq *ServiceEndpointDefinitionQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = sedq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{serviceendpointdefinition.Label}
	default:
		err = &NotSingularError{serviceendpointdefinition.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (sedq *ServiceEndpointDefinitionQuery) OnlyXID(ctx context.Context) int {
	id, err := sedq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ServiceEndpointDefinitions.
func (sedq *ServiceEndpointDefinitionQuery) All(ctx context.Context) ([]*ServiceEndpointDefinition, error) {
	if err := sedq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return sedq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (sedq *ServiceEndpointDefinitionQuery) AllX(ctx context.Context) []*ServiceEndpointDefinition {
	seds, err := sedq.All(ctx)
	if err != nil {
		panic(err)
	}
	return seds
}

// IDs executes the query and returns a list of ServiceEndpointDefinition ids.
func (sedq *ServiceEndpointDefinitionQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := sedq.Select(serviceendpointdefinition.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (sedq *ServiceEndpointDefinitionQuery) IDsX(ctx context.Context) []int {
	ids, err := sedq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (sedq *ServiceEndpointDefinitionQuery) Count(ctx context.Context) (int, error) {
	if err := sedq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return sedq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (sedq *ServiceEndpointDefinitionQuery) CountX(ctx context.Context) int {
	count, err := sedq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (sedq *ServiceEndpointDefinitionQuery) Exist(ctx context.Context) (bool, error) {
	if err := sedq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return sedq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (sedq *ServiceEndpointDefinitionQuery) ExistX(ctx context.Context) bool {
	exist, err := sedq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (sedq *ServiceEndpointDefinitionQuery) Clone() *ServiceEndpointDefinitionQuery {
	return &ServiceEndpointDefinitionQuery{
		config:     sedq.config,
		limit:      sedq.limit,
		offset:     sedq.offset,
		order:      append([]OrderFunc{}, sedq.order...),
		unique:     append([]string{}, sedq.unique...),
		predicates: append([]predicate.ServiceEndpointDefinition{}, sedq.predicates...),
		// clone intermediate query.
		sql:  sedq.sql.Clone(),
		path: sedq.path,
	}
}

//  WithEndpoints tells the query-builder to eager-loads the nodes that are connected to
// the "endpoints" edge. The optional arguments used to configure the query builder of the edge.
func (sedq *ServiceEndpointDefinitionQuery) WithEndpoints(opts ...func(*ServiceEndpointQuery)) *ServiceEndpointDefinitionQuery {
	query := &ServiceEndpointQuery{config: sedq.config}
	for _, opt := range opts {
		opt(query)
	}
	sedq.withEndpoints = query
	return sedq
}

//  WithServiceType tells the query-builder to eager-loads the nodes that are connected to
// the "service_type" edge. The optional arguments used to configure the query builder of the edge.
func (sedq *ServiceEndpointDefinitionQuery) WithServiceType(opts ...func(*ServiceTypeQuery)) *ServiceEndpointDefinitionQuery {
	query := &ServiceTypeQuery{config: sedq.config}
	for _, opt := range opts {
		opt(query)
	}
	sedq.withServiceType = query
	return sedq
}

//  WithEquipmentType tells the query-builder to eager-loads the nodes that are connected to
// the "equipment_type" edge. The optional arguments used to configure the query builder of the edge.
func (sedq *ServiceEndpointDefinitionQuery) WithEquipmentType(opts ...func(*EquipmentTypeQuery)) *ServiceEndpointDefinitionQuery {
	query := &EquipmentTypeQuery{config: sedq.config}
	for _, opt := range opts {
		opt(query)
	}
	sedq.withEquipmentType = query
	return sedq
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
//	client.ServiceEndpointDefinition.Query().
//		GroupBy(serviceendpointdefinition.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (sedq *ServiceEndpointDefinitionQuery) GroupBy(field string, fields ...string) *ServiceEndpointDefinitionGroupBy {
	group := &ServiceEndpointDefinitionGroupBy{config: sedq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := sedq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return sedq.sqlQuery(), nil
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
//	client.ServiceEndpointDefinition.Query().
//		Select(serviceendpointdefinition.FieldCreateTime).
//		Scan(ctx, &v)
//
func (sedq *ServiceEndpointDefinitionQuery) Select(field string, fields ...string) *ServiceEndpointDefinitionSelect {
	selector := &ServiceEndpointDefinitionSelect{config: sedq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := sedq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return sedq.sqlQuery(), nil
	}
	return selector
}

func (sedq *ServiceEndpointDefinitionQuery) prepareQuery(ctx context.Context) error {
	if sedq.path != nil {
		prev, err := sedq.path(ctx)
		if err != nil {
			return err
		}
		sedq.sql = prev
	}
	if err := serviceendpointdefinition.Policy.EvalQuery(ctx, sedq); err != nil {
		return err
	}
	return nil
}

func (sedq *ServiceEndpointDefinitionQuery) sqlAll(ctx context.Context) ([]*ServiceEndpointDefinition, error) {
	var (
		nodes       = []*ServiceEndpointDefinition{}
		withFKs     = sedq.withFKs
		_spec       = sedq.querySpec()
		loadedTypes = [3]bool{
			sedq.withEndpoints != nil,
			sedq.withServiceType != nil,
			sedq.withEquipmentType != nil,
		}
	)
	if sedq.withServiceType != nil || sedq.withEquipmentType != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, serviceendpointdefinition.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &ServiceEndpointDefinition{config: sedq.config}
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
	if err := sqlgraph.QueryNodes(ctx, sedq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := sedq.withEndpoints; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*ServiceEndpointDefinition)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.ServiceEndpoint(func(s *sql.Selector) {
			s.Where(sql.InValues(serviceendpointdefinition.EndpointsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.service_endpoint_definition_endpoints
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "service_endpoint_definition_endpoints" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "service_endpoint_definition_endpoints" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Endpoints = append(node.Edges.Endpoints, n)
		}
	}

	if query := sedq.withServiceType; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*ServiceEndpointDefinition)
		for i := range nodes {
			if fk := nodes[i].service_type_endpoint_definitions; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(servicetype.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "service_type_endpoint_definitions" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.ServiceType = n
			}
		}
	}

	if query := sedq.withEquipmentType; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*ServiceEndpointDefinition)
		for i := range nodes {
			if fk := nodes[i].equipment_type_service_endpoint_definitions; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(equipmenttype.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_type_service_endpoint_definitions" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.EquipmentType = n
			}
		}
	}

	return nodes, nil
}

func (sedq *ServiceEndpointDefinitionQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := sedq.querySpec()
	return sqlgraph.CountNodes(ctx, sedq.driver, _spec)
}

func (sedq *ServiceEndpointDefinitionQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := sedq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (sedq *ServiceEndpointDefinitionQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   serviceendpointdefinition.Table,
			Columns: serviceendpointdefinition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: serviceendpointdefinition.FieldID,
			},
		},
		From:   sedq.sql,
		Unique: true,
	}
	if ps := sedq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := sedq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := sedq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := sedq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (sedq *ServiceEndpointDefinitionQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(sedq.driver.Dialect())
	t1 := builder.Table(serviceendpointdefinition.Table)
	selector := builder.Select(t1.Columns(serviceendpointdefinition.Columns...)...).From(t1)
	if sedq.sql != nil {
		selector = sedq.sql
		selector.Select(selector.Columns(serviceendpointdefinition.Columns...)...)
	}
	for _, p := range sedq.predicates {
		p(selector)
	}
	for _, p := range sedq.order {
		p(selector)
	}
	if offset := sedq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := sedq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ServiceEndpointDefinitionGroupBy is the builder for group-by ServiceEndpointDefinition entities.
type ServiceEndpointDefinitionGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (sedgb *ServiceEndpointDefinitionGroupBy) Aggregate(fns ...AggregateFunc) *ServiceEndpointDefinitionGroupBy {
	sedgb.fns = append(sedgb.fns, fns...)
	return sedgb
}

// Scan applies the group-by query and scan the result into the given value.
func (sedgb *ServiceEndpointDefinitionGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := sedgb.path(ctx)
	if err != nil {
		return err
	}
	sedgb.sql = query
	return sedgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (sedgb *ServiceEndpointDefinitionGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := sedgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (sedgb *ServiceEndpointDefinitionGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(sedgb.fields) > 1 {
		return nil, errors.New("ent: ServiceEndpointDefinitionGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := sedgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (sedgb *ServiceEndpointDefinitionGroupBy) StringsX(ctx context.Context) []string {
	v, err := sedgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (sedgb *ServiceEndpointDefinitionGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(sedgb.fields) > 1 {
		return nil, errors.New("ent: ServiceEndpointDefinitionGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := sedgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (sedgb *ServiceEndpointDefinitionGroupBy) IntsX(ctx context.Context) []int {
	v, err := sedgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (sedgb *ServiceEndpointDefinitionGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(sedgb.fields) > 1 {
		return nil, errors.New("ent: ServiceEndpointDefinitionGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := sedgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (sedgb *ServiceEndpointDefinitionGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := sedgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (sedgb *ServiceEndpointDefinitionGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(sedgb.fields) > 1 {
		return nil, errors.New("ent: ServiceEndpointDefinitionGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := sedgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (sedgb *ServiceEndpointDefinitionGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := sedgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (sedgb *ServiceEndpointDefinitionGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := sedgb.sqlQuery().Query()
	if err := sedgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (sedgb *ServiceEndpointDefinitionGroupBy) sqlQuery() *sql.Selector {
	selector := sedgb.sql
	columns := make([]string, 0, len(sedgb.fields)+len(sedgb.fns))
	columns = append(columns, sedgb.fields...)
	for _, fn := range sedgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(sedgb.fields...)
}

// ServiceEndpointDefinitionSelect is the builder for select fields of ServiceEndpointDefinition entities.
type ServiceEndpointDefinitionSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (seds *ServiceEndpointDefinitionSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := seds.path(ctx)
	if err != nil {
		return err
	}
	seds.sql = query
	return seds.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (seds *ServiceEndpointDefinitionSelect) ScanX(ctx context.Context, v interface{}) {
	if err := seds.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (seds *ServiceEndpointDefinitionSelect) Strings(ctx context.Context) ([]string, error) {
	if len(seds.fields) > 1 {
		return nil, errors.New("ent: ServiceEndpointDefinitionSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := seds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (seds *ServiceEndpointDefinitionSelect) StringsX(ctx context.Context) []string {
	v, err := seds.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (seds *ServiceEndpointDefinitionSelect) Ints(ctx context.Context) ([]int, error) {
	if len(seds.fields) > 1 {
		return nil, errors.New("ent: ServiceEndpointDefinitionSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := seds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (seds *ServiceEndpointDefinitionSelect) IntsX(ctx context.Context) []int {
	v, err := seds.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (seds *ServiceEndpointDefinitionSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(seds.fields) > 1 {
		return nil, errors.New("ent: ServiceEndpointDefinitionSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := seds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (seds *ServiceEndpointDefinitionSelect) Float64sX(ctx context.Context) []float64 {
	v, err := seds.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (seds *ServiceEndpointDefinitionSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(seds.fields) > 1 {
		return nil, errors.New("ent: ServiceEndpointDefinitionSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := seds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (seds *ServiceEndpointDefinitionSelect) BoolsX(ctx context.Context) []bool {
	v, err := seds.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (seds *ServiceEndpointDefinitionSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := seds.sqlQuery().Query()
	if err := seds.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (seds *ServiceEndpointDefinitionSelect) sqlQuery() sql.Querier {
	selector := seds.sql
	selector.Select(selector.Columns(seds.fields...)...)
	return selector
}
