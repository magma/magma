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
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentcategory"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpointdefinition"
)

// EquipmentTypeQuery is the builder for querying EquipmentType entities.
type EquipmentTypeQuery struct {
	config
	limit      *int
	offset     *int
	order      []OrderFunc
	unique     []string
	predicates []predicate.EquipmentType
	// eager-loading edges.
	withPortDefinitions            *EquipmentPortDefinitionQuery
	withPositionDefinitions        *EquipmentPositionDefinitionQuery
	withPropertyTypes              *PropertyTypeQuery
	withEquipment                  *EquipmentQuery
	withCategory                   *EquipmentCategoryQuery
	withServiceEndpointDefinitions *ServiceEndpointDefinitionQuery
	withFKs                        bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (etq *EquipmentTypeQuery) Where(ps ...predicate.EquipmentType) *EquipmentTypeQuery {
	etq.predicates = append(etq.predicates, ps...)
	return etq
}

// Limit adds a limit step to the query.
func (etq *EquipmentTypeQuery) Limit(limit int) *EquipmentTypeQuery {
	etq.limit = &limit
	return etq
}

// Offset adds an offset step to the query.
func (etq *EquipmentTypeQuery) Offset(offset int) *EquipmentTypeQuery {
	etq.offset = &offset
	return etq
}

// Order adds an order step to the query.
func (etq *EquipmentTypeQuery) Order(o ...OrderFunc) *EquipmentTypeQuery {
	etq.order = append(etq.order, o...)
	return etq
}

// QueryPortDefinitions chains the current query on the port_definitions edge.
func (etq *EquipmentTypeQuery) QueryPortDefinitions() *EquipmentPortDefinitionQuery {
	query := &EquipmentPortDefinitionQuery{config: etq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := etq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmenttype.Table, equipmenttype.FieldID, etq.sqlQuery()),
			sqlgraph.To(equipmentportdefinition.Table, equipmentportdefinition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipmenttype.PortDefinitionsTable, equipmenttype.PortDefinitionsColumn),
		)
		fromU = sqlgraph.SetNeighbors(etq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryPositionDefinitions chains the current query on the position_definitions edge.
func (etq *EquipmentTypeQuery) QueryPositionDefinitions() *EquipmentPositionDefinitionQuery {
	query := &EquipmentPositionDefinitionQuery{config: etq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := etq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmenttype.Table, equipmenttype.FieldID, etq.sqlQuery()),
			sqlgraph.To(equipmentpositiondefinition.Table, equipmentpositiondefinition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipmenttype.PositionDefinitionsTable, equipmenttype.PositionDefinitionsColumn),
		)
		fromU = sqlgraph.SetNeighbors(etq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryPropertyTypes chains the current query on the property_types edge.
func (etq *EquipmentTypeQuery) QueryPropertyTypes() *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: etq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := etq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmenttype.Table, equipmenttype.FieldID, etq.sqlQuery()),
			sqlgraph.To(propertytype.Table, propertytype.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipmenttype.PropertyTypesTable, equipmenttype.PropertyTypesColumn),
		)
		fromU = sqlgraph.SetNeighbors(etq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryEquipment chains the current query on the equipment edge.
func (etq *EquipmentTypeQuery) QueryEquipment() *EquipmentQuery {
	query := &EquipmentQuery{config: etq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := etq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmenttype.Table, equipmenttype.FieldID, etq.sqlQuery()),
			sqlgraph.To(equipment.Table, equipment.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, equipmenttype.EquipmentTable, equipmenttype.EquipmentColumn),
		)
		fromU = sqlgraph.SetNeighbors(etq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryCategory chains the current query on the category edge.
func (etq *EquipmentTypeQuery) QueryCategory() *EquipmentCategoryQuery {
	query := &EquipmentCategoryQuery{config: etq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := etq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmenttype.Table, equipmenttype.FieldID, etq.sqlQuery()),
			sqlgraph.To(equipmentcategory.Table, equipmentcategory.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, equipmenttype.CategoryTable, equipmenttype.CategoryColumn),
		)
		fromU = sqlgraph.SetNeighbors(etq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryServiceEndpointDefinitions chains the current query on the service_endpoint_definitions edge.
func (etq *EquipmentTypeQuery) QueryServiceEndpointDefinitions() *ServiceEndpointDefinitionQuery {
	query := &ServiceEndpointDefinitionQuery{config: etq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := etq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmenttype.Table, equipmenttype.FieldID, etq.sqlQuery()),
			sqlgraph.To(serviceendpointdefinition.Table, serviceendpointdefinition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipmenttype.ServiceEndpointDefinitionsTable, equipmenttype.ServiceEndpointDefinitionsColumn),
		)
		fromU = sqlgraph.SetNeighbors(etq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first EquipmentType entity in the query. Returns *NotFoundError when no equipmenttype was found.
func (etq *EquipmentTypeQuery) First(ctx context.Context) (*EquipmentType, error) {
	ets, err := etq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(ets) == 0 {
		return nil, &NotFoundError{equipmenttype.Label}
	}
	return ets[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (etq *EquipmentTypeQuery) FirstX(ctx context.Context) *EquipmentType {
	et, err := etq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return et
}

// FirstID returns the first EquipmentType id in the query. Returns *NotFoundError when no id was found.
func (etq *EquipmentTypeQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = etq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{equipmenttype.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (etq *EquipmentTypeQuery) FirstXID(ctx context.Context) int {
	id, err := etq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only EquipmentType entity in the query, returns an error if not exactly one entity was returned.
func (etq *EquipmentTypeQuery) Only(ctx context.Context) (*EquipmentType, error) {
	ets, err := etq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(ets) {
	case 1:
		return ets[0], nil
	case 0:
		return nil, &NotFoundError{equipmenttype.Label}
	default:
		return nil, &NotSingularError{equipmenttype.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (etq *EquipmentTypeQuery) OnlyX(ctx context.Context) *EquipmentType {
	et, err := etq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return et
}

// OnlyID returns the only EquipmentType id in the query, returns an error if not exactly one id was returned.
func (etq *EquipmentTypeQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = etq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{equipmenttype.Label}
	default:
		err = &NotSingularError{equipmenttype.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (etq *EquipmentTypeQuery) OnlyXID(ctx context.Context) int {
	id, err := etq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of EquipmentTypes.
func (etq *EquipmentTypeQuery) All(ctx context.Context) ([]*EquipmentType, error) {
	if err := etq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return etq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (etq *EquipmentTypeQuery) AllX(ctx context.Context) []*EquipmentType {
	ets, err := etq.All(ctx)
	if err != nil {
		panic(err)
	}
	return ets
}

// IDs executes the query and returns a list of EquipmentType ids.
func (etq *EquipmentTypeQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := etq.Select(equipmenttype.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (etq *EquipmentTypeQuery) IDsX(ctx context.Context) []int {
	ids, err := etq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (etq *EquipmentTypeQuery) Count(ctx context.Context) (int, error) {
	if err := etq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return etq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (etq *EquipmentTypeQuery) CountX(ctx context.Context) int {
	count, err := etq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (etq *EquipmentTypeQuery) Exist(ctx context.Context) (bool, error) {
	if err := etq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return etq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (etq *EquipmentTypeQuery) ExistX(ctx context.Context) bool {
	exist, err := etq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (etq *EquipmentTypeQuery) Clone() *EquipmentTypeQuery {
	return &EquipmentTypeQuery{
		config:     etq.config,
		limit:      etq.limit,
		offset:     etq.offset,
		order:      append([]OrderFunc{}, etq.order...),
		unique:     append([]string{}, etq.unique...),
		predicates: append([]predicate.EquipmentType{}, etq.predicates...),
		// clone intermediate query.
		sql:  etq.sql.Clone(),
		path: etq.path,
	}
}

//  WithPortDefinitions tells the query-builder to eager-loads the nodes that are connected to
// the "port_definitions" edge. The optional arguments used to configure the query builder of the edge.
func (etq *EquipmentTypeQuery) WithPortDefinitions(opts ...func(*EquipmentPortDefinitionQuery)) *EquipmentTypeQuery {
	query := &EquipmentPortDefinitionQuery{config: etq.config}
	for _, opt := range opts {
		opt(query)
	}
	etq.withPortDefinitions = query
	return etq
}

//  WithPositionDefinitions tells the query-builder to eager-loads the nodes that are connected to
// the "position_definitions" edge. The optional arguments used to configure the query builder of the edge.
func (etq *EquipmentTypeQuery) WithPositionDefinitions(opts ...func(*EquipmentPositionDefinitionQuery)) *EquipmentTypeQuery {
	query := &EquipmentPositionDefinitionQuery{config: etq.config}
	for _, opt := range opts {
		opt(query)
	}
	etq.withPositionDefinitions = query
	return etq
}

//  WithPropertyTypes tells the query-builder to eager-loads the nodes that are connected to
// the "property_types" edge. The optional arguments used to configure the query builder of the edge.
func (etq *EquipmentTypeQuery) WithPropertyTypes(opts ...func(*PropertyTypeQuery)) *EquipmentTypeQuery {
	query := &PropertyTypeQuery{config: etq.config}
	for _, opt := range opts {
		opt(query)
	}
	etq.withPropertyTypes = query
	return etq
}

//  WithEquipment tells the query-builder to eager-loads the nodes that are connected to
// the "equipment" edge. The optional arguments used to configure the query builder of the edge.
func (etq *EquipmentTypeQuery) WithEquipment(opts ...func(*EquipmentQuery)) *EquipmentTypeQuery {
	query := &EquipmentQuery{config: etq.config}
	for _, opt := range opts {
		opt(query)
	}
	etq.withEquipment = query
	return etq
}

//  WithCategory tells the query-builder to eager-loads the nodes that are connected to
// the "category" edge. The optional arguments used to configure the query builder of the edge.
func (etq *EquipmentTypeQuery) WithCategory(opts ...func(*EquipmentCategoryQuery)) *EquipmentTypeQuery {
	query := &EquipmentCategoryQuery{config: etq.config}
	for _, opt := range opts {
		opt(query)
	}
	etq.withCategory = query
	return etq
}

//  WithServiceEndpointDefinitions tells the query-builder to eager-loads the nodes that are connected to
// the "service_endpoint_definitions" edge. The optional arguments used to configure the query builder of the edge.
func (etq *EquipmentTypeQuery) WithServiceEndpointDefinitions(opts ...func(*ServiceEndpointDefinitionQuery)) *EquipmentTypeQuery {
	query := &ServiceEndpointDefinitionQuery{config: etq.config}
	for _, opt := range opts {
		opt(query)
	}
	etq.withServiceEndpointDefinitions = query
	return etq
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
//	client.EquipmentType.Query().
//		GroupBy(equipmenttype.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (etq *EquipmentTypeQuery) GroupBy(field string, fields ...string) *EquipmentTypeGroupBy {
	group := &EquipmentTypeGroupBy{config: etq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := etq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return etq.sqlQuery(), nil
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
//	client.EquipmentType.Query().
//		Select(equipmenttype.FieldCreateTime).
//		Scan(ctx, &v)
//
func (etq *EquipmentTypeQuery) Select(field string, fields ...string) *EquipmentTypeSelect {
	selector := &EquipmentTypeSelect{config: etq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := etq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return etq.sqlQuery(), nil
	}
	return selector
}

func (etq *EquipmentTypeQuery) prepareQuery(ctx context.Context) error {
	if etq.path != nil {
		prev, err := etq.path(ctx)
		if err != nil {
			return err
		}
		etq.sql = prev
	}
	if err := equipmenttype.Policy.EvalQuery(ctx, etq); err != nil {
		return err
	}
	return nil
}

func (etq *EquipmentTypeQuery) sqlAll(ctx context.Context) ([]*EquipmentType, error) {
	var (
		nodes       = []*EquipmentType{}
		withFKs     = etq.withFKs
		_spec       = etq.querySpec()
		loadedTypes = [6]bool{
			etq.withPortDefinitions != nil,
			etq.withPositionDefinitions != nil,
			etq.withPropertyTypes != nil,
			etq.withEquipment != nil,
			etq.withCategory != nil,
			etq.withServiceEndpointDefinitions != nil,
		}
	)
	if etq.withCategory != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, equipmenttype.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &EquipmentType{config: etq.config}
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
	if err := sqlgraph.QueryNodes(ctx, etq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := etq.withPortDefinitions; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*EquipmentType)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.EquipmentPortDefinition(func(s *sql.Selector) {
			s.Where(sql.InValues(equipmenttype.PortDefinitionsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.equipment_type_port_definitions
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "equipment_type_port_definitions" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_type_port_definitions" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.PortDefinitions = append(node.Edges.PortDefinitions, n)
		}
	}

	if query := etq.withPositionDefinitions; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*EquipmentType)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
			s.Where(sql.InValues(equipmenttype.PositionDefinitionsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.equipment_type_position_definitions
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "equipment_type_position_definitions" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_type_position_definitions" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.PositionDefinitions = append(node.Edges.PositionDefinitions, n)
		}
	}

	if query := etq.withPropertyTypes; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*EquipmentType)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.PropertyType(func(s *sql.Selector) {
			s.Where(sql.InValues(equipmenttype.PropertyTypesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.equipment_type_property_types
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "equipment_type_property_types" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_type_property_types" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.PropertyTypes = append(node.Edges.PropertyTypes, n)
		}
	}

	if query := etq.withEquipment; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*EquipmentType)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.Equipment(func(s *sql.Selector) {
			s.Where(sql.InValues(equipmenttype.EquipmentColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.equipment_type
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "equipment_type" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_type" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Equipment = append(node.Edges.Equipment, n)
		}
	}

	if query := etq.withCategory; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*EquipmentType)
		for i := range nodes {
			if fk := nodes[i].equipment_type_category; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(equipmentcategory.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_type_category" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Category = n
			}
		}
	}

	if query := etq.withServiceEndpointDefinitions; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*EquipmentType)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.ServiceEndpointDefinition(func(s *sql.Selector) {
			s.Where(sql.InValues(equipmenttype.ServiceEndpointDefinitionsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.equipment_type_service_endpoint_definitions
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "equipment_type_service_endpoint_definitions" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_type_service_endpoint_definitions" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.ServiceEndpointDefinitions = append(node.Edges.ServiceEndpointDefinitions, n)
		}
	}

	return nodes, nil
}

func (etq *EquipmentTypeQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := etq.querySpec()
	return sqlgraph.CountNodes(ctx, etq.driver, _spec)
}

func (etq *EquipmentTypeQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := etq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (etq *EquipmentTypeQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmenttype.Table,
			Columns: equipmenttype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmenttype.FieldID,
			},
		},
		From:   etq.sql,
		Unique: true,
	}
	if ps := etq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := etq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := etq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := etq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (etq *EquipmentTypeQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(etq.driver.Dialect())
	t1 := builder.Table(equipmenttype.Table)
	selector := builder.Select(t1.Columns(equipmenttype.Columns...)...).From(t1)
	if etq.sql != nil {
		selector = etq.sql
		selector.Select(selector.Columns(equipmenttype.Columns...)...)
	}
	for _, p := range etq.predicates {
		p(selector)
	}
	for _, p := range etq.order {
		p(selector)
	}
	if offset := etq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := etq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// EquipmentTypeGroupBy is the builder for group-by EquipmentType entities.
type EquipmentTypeGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (etgb *EquipmentTypeGroupBy) Aggregate(fns ...AggregateFunc) *EquipmentTypeGroupBy {
	etgb.fns = append(etgb.fns, fns...)
	return etgb
}

// Scan applies the group-by query and scan the result into the given value.
func (etgb *EquipmentTypeGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := etgb.path(ctx)
	if err != nil {
		return err
	}
	etgb.sql = query
	return etgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (etgb *EquipmentTypeGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := etgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (etgb *EquipmentTypeGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(etgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentTypeGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := etgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (etgb *EquipmentTypeGroupBy) StringsX(ctx context.Context) []string {
	v, err := etgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (etgb *EquipmentTypeGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(etgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentTypeGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := etgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (etgb *EquipmentTypeGroupBy) IntsX(ctx context.Context) []int {
	v, err := etgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (etgb *EquipmentTypeGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(etgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentTypeGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := etgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (etgb *EquipmentTypeGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := etgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (etgb *EquipmentTypeGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(etgb.fields) > 1 {
		return nil, errors.New("ent: EquipmentTypeGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := etgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (etgb *EquipmentTypeGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := etgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (etgb *EquipmentTypeGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := etgb.sqlQuery().Query()
	if err := etgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (etgb *EquipmentTypeGroupBy) sqlQuery() *sql.Selector {
	selector := etgb.sql
	columns := make([]string, 0, len(etgb.fields)+len(etgb.fns))
	columns = append(columns, etgb.fields...)
	for _, fn := range etgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(etgb.fields...)
}

// EquipmentTypeSelect is the builder for select fields of EquipmentType entities.
type EquipmentTypeSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (ets *EquipmentTypeSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := ets.path(ctx)
	if err != nil {
		return err
	}
	ets.sql = query
	return ets.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ets *EquipmentTypeSelect) ScanX(ctx context.Context, v interface{}) {
	if err := ets.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (ets *EquipmentTypeSelect) Strings(ctx context.Context) ([]string, error) {
	if len(ets.fields) > 1 {
		return nil, errors.New("ent: EquipmentTypeSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := ets.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ets *EquipmentTypeSelect) StringsX(ctx context.Context) []string {
	v, err := ets.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (ets *EquipmentTypeSelect) Ints(ctx context.Context) ([]int, error) {
	if len(ets.fields) > 1 {
		return nil, errors.New("ent: EquipmentTypeSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := ets.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ets *EquipmentTypeSelect) IntsX(ctx context.Context) []int {
	v, err := ets.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (ets *EquipmentTypeSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(ets.fields) > 1 {
		return nil, errors.New("ent: EquipmentTypeSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := ets.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ets *EquipmentTypeSelect) Float64sX(ctx context.Context) []float64 {
	v, err := ets.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (ets *EquipmentTypeSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(ets.fields) > 1 {
		return nil, errors.New("ent: EquipmentTypeSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := ets.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ets *EquipmentTypeSelect) BoolsX(ctx context.Context) []bool {
	v, err := ets.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ets *EquipmentTypeSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ets.sqlQuery().Query()
	if err := ets.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ets *EquipmentTypeSelect) sqlQuery() sql.Querier {
	selector := ets.sql
	selector.Select(selector.Columns(ets.fields...)...)
	return selector
}
