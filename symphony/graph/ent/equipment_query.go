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
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// EquipmentQuery is the builder for querying Equipment entities.
type EquipmentQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.Equipment
	// eager-loading edges.
	withType           *EquipmentTypeQuery
	withLocation       *LocationQuery
	withParentPosition *EquipmentPositionQuery
	withPositions      *EquipmentPositionQuery
	withPorts          *EquipmentPortQuery
	withWorkOrder      *WorkOrderQuery
	withProperties     *PropertyQuery
	withFiles          *FileQuery
	withHyperlinks     *HyperlinkQuery
	withEndpoints      *ServiceEndpointQuery
	withFKs            bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (eq *EquipmentQuery) Where(ps ...predicate.Equipment) *EquipmentQuery {
	eq.predicates = append(eq.predicates, ps...)
	return eq
}

// Limit adds a limit step to the query.
func (eq *EquipmentQuery) Limit(limit int) *EquipmentQuery {
	eq.limit = &limit
	return eq
}

// Offset adds an offset step to the query.
func (eq *EquipmentQuery) Offset(offset int) *EquipmentQuery {
	eq.offset = &offset
	return eq
}

// Order adds an order step to the query.
func (eq *EquipmentQuery) Order(o ...Order) *EquipmentQuery {
	eq.order = append(eq.order, o...)
	return eq
}

// QueryType chains the current query on the type edge.
func (eq *EquipmentQuery) QueryType() *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: eq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
			sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, equipment.TypeTable, equipment.TypeColumn),
		)
		fromU = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryLocation chains the current query on the location edge.
func (eq *EquipmentQuery) QueryLocation() *LocationQuery {
	query := &LocationQuery{config: eq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
			sqlgraph.To(location.Table, location.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, equipment.LocationTable, equipment.LocationColumn),
		)
		fromU = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryParentPosition chains the current query on the parent_position edge.
func (eq *EquipmentQuery) QueryParentPosition() *EquipmentPositionQuery {
	query := &EquipmentPositionQuery{config: eq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
			sqlgraph.To(equipmentposition.Table, equipmentposition.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, equipment.ParentPositionTable, equipment.ParentPositionColumn),
		)
		fromU = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryPositions chains the current query on the positions edge.
func (eq *EquipmentQuery) QueryPositions() *EquipmentPositionQuery {
	query := &EquipmentPositionQuery{config: eq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
			sqlgraph.To(equipmentposition.Table, equipmentposition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipment.PositionsTable, equipment.PositionsColumn),
		)
		fromU = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryPorts chains the current query on the ports edge.
func (eq *EquipmentQuery) QueryPorts() *EquipmentPortQuery {
	query := &EquipmentPortQuery{config: eq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
			sqlgraph.To(equipmentport.Table, equipmentport.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipment.PortsTable, equipment.PortsColumn),
		)
		fromU = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryWorkOrder chains the current query on the work_order edge.
func (eq *EquipmentQuery) QueryWorkOrder() *WorkOrderQuery {
	query := &WorkOrderQuery{config: eq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
			sqlgraph.To(workorder.Table, workorder.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, equipment.WorkOrderTable, equipment.WorkOrderColumn),
		)
		fromU = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryProperties chains the current query on the properties edge.
func (eq *EquipmentQuery) QueryProperties() *PropertyQuery {
	query := &PropertyQuery{config: eq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
			sqlgraph.To(property.Table, property.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipment.PropertiesTable, equipment.PropertiesColumn),
		)
		fromU = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryFiles chains the current query on the files edge.
func (eq *EquipmentQuery) QueryFiles() *FileQuery {
	query := &FileQuery{config: eq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipment.FilesTable, equipment.FilesColumn),
		)
		fromU = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryHyperlinks chains the current query on the hyperlinks edge.
func (eq *EquipmentQuery) QueryHyperlinks() *HyperlinkQuery {
	query := &HyperlinkQuery{config: eq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
			sqlgraph.To(hyperlink.Table, hyperlink.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipment.HyperlinksTable, equipment.HyperlinksColumn),
		)
		fromU = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryEndpoints chains the current query on the endpoints edge.
func (eq *EquipmentQuery) QueryEndpoints() *ServiceEndpointQuery {
	query := &ServiceEndpointQuery{config: eq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, eq.sqlQuery()),
			sqlgraph.To(serviceendpoint.Table, serviceendpoint.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, equipment.EndpointsTable, equipment.EndpointsColumn),
		)
		fromU = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Equipment entity in the query. Returns *NotFoundError when no equipment was found.
func (eq *EquipmentQuery) First(ctx context.Context) (*Equipment, error) {
	es, err := eq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(es) == 0 {
		return nil, &NotFoundError{equipment.Label}
	}
	return es[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (eq *EquipmentQuery) FirstX(ctx context.Context) *Equipment {
	e, err := eq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return e
}

// FirstID returns the first Equipment id in the query. Returns *NotFoundError when no id was found.
func (eq *EquipmentQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = eq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{equipment.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (eq *EquipmentQuery) FirstXID(ctx context.Context) int {
	id, err := eq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only Equipment entity in the query, returns an error if not exactly one entity was returned.
func (eq *EquipmentQuery) Only(ctx context.Context) (*Equipment, error) {
	es, err := eq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(es) {
	case 1:
		return es[0], nil
	case 0:
		return nil, &NotFoundError{equipment.Label}
	default:
		return nil, &NotSingularError{equipment.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (eq *EquipmentQuery) OnlyX(ctx context.Context) *Equipment {
	e, err := eq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return e
}

// OnlyID returns the only Equipment id in the query, returns an error if not exactly one id was returned.
func (eq *EquipmentQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = eq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{equipment.Label}
	default:
		err = &NotSingularError{equipment.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (eq *EquipmentQuery) OnlyXID(ctx context.Context) int {
	id, err := eq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of EquipmentSlice.
func (eq *EquipmentQuery) All(ctx context.Context) ([]*Equipment, error) {
	if err := eq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return eq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (eq *EquipmentQuery) AllX(ctx context.Context) []*Equipment {
	es, err := eq.All(ctx)
	if err != nil {
		panic(err)
	}
	return es
}

// IDs executes the query and returns a list of Equipment ids.
func (eq *EquipmentQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := eq.Select(equipment.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (eq *EquipmentQuery) IDsX(ctx context.Context) []int {
	ids, err := eq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (eq *EquipmentQuery) Count(ctx context.Context) (int, error) {
	if err := eq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return eq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (eq *EquipmentQuery) CountX(ctx context.Context) int {
	count, err := eq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (eq *EquipmentQuery) Exist(ctx context.Context) (bool, error) {
	if err := eq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return eq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (eq *EquipmentQuery) ExistX(ctx context.Context) bool {
	exist, err := eq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (eq *EquipmentQuery) Clone() *EquipmentQuery {
	return &EquipmentQuery{
		config:     eq.config,
		limit:      eq.limit,
		offset:     eq.offset,
		order:      append([]Order{}, eq.order...),
		unique:     append([]string{}, eq.unique...),
		predicates: append([]predicate.Equipment{}, eq.predicates...),
		// clone intermediate query.
		sql:  eq.sql.Clone(),
		path: eq.path,
	}
}

//  WithType tells the query-builder to eager-loads the nodes that are connected to
// the "type" edge. The optional arguments used to configure the query builder of the edge.
func (eq *EquipmentQuery) WithType(opts ...func(*EquipmentTypeQuery)) *EquipmentQuery {
	query := &EquipmentTypeQuery{config: eq.config}
	for _, opt := range opts {
		opt(query)
	}
	eq.withType = query
	return eq
}

//  WithLocation tells the query-builder to eager-loads the nodes that are connected to
// the "location" edge. The optional arguments used to configure the query builder of the edge.
func (eq *EquipmentQuery) WithLocation(opts ...func(*LocationQuery)) *EquipmentQuery {
	query := &LocationQuery{config: eq.config}
	for _, opt := range opts {
		opt(query)
	}
	eq.withLocation = query
	return eq
}

//  WithParentPosition tells the query-builder to eager-loads the nodes that are connected to
// the "parent_position" edge. The optional arguments used to configure the query builder of the edge.
func (eq *EquipmentQuery) WithParentPosition(opts ...func(*EquipmentPositionQuery)) *EquipmentQuery {
	query := &EquipmentPositionQuery{config: eq.config}
	for _, opt := range opts {
		opt(query)
	}
	eq.withParentPosition = query
	return eq
}

//  WithPositions tells the query-builder to eager-loads the nodes that are connected to
// the "positions" edge. The optional arguments used to configure the query builder of the edge.
func (eq *EquipmentQuery) WithPositions(opts ...func(*EquipmentPositionQuery)) *EquipmentQuery {
	query := &EquipmentPositionQuery{config: eq.config}
	for _, opt := range opts {
		opt(query)
	}
	eq.withPositions = query
	return eq
}

//  WithPorts tells the query-builder to eager-loads the nodes that are connected to
// the "ports" edge. The optional arguments used to configure the query builder of the edge.
func (eq *EquipmentQuery) WithPorts(opts ...func(*EquipmentPortQuery)) *EquipmentQuery {
	query := &EquipmentPortQuery{config: eq.config}
	for _, opt := range opts {
		opt(query)
	}
	eq.withPorts = query
	return eq
}

//  WithWorkOrder tells the query-builder to eager-loads the nodes that are connected to
// the "work_order" edge. The optional arguments used to configure the query builder of the edge.
func (eq *EquipmentQuery) WithWorkOrder(opts ...func(*WorkOrderQuery)) *EquipmentQuery {
	query := &WorkOrderQuery{config: eq.config}
	for _, opt := range opts {
		opt(query)
	}
	eq.withWorkOrder = query
	return eq
}

//  WithProperties tells the query-builder to eager-loads the nodes that are connected to
// the "properties" edge. The optional arguments used to configure the query builder of the edge.
func (eq *EquipmentQuery) WithProperties(opts ...func(*PropertyQuery)) *EquipmentQuery {
	query := &PropertyQuery{config: eq.config}
	for _, opt := range opts {
		opt(query)
	}
	eq.withProperties = query
	return eq
}

//  WithFiles tells the query-builder to eager-loads the nodes that are connected to
// the "files" edge. The optional arguments used to configure the query builder of the edge.
func (eq *EquipmentQuery) WithFiles(opts ...func(*FileQuery)) *EquipmentQuery {
	query := &FileQuery{config: eq.config}
	for _, opt := range opts {
		opt(query)
	}
	eq.withFiles = query
	return eq
}

//  WithHyperlinks tells the query-builder to eager-loads the nodes that are connected to
// the "hyperlinks" edge. The optional arguments used to configure the query builder of the edge.
func (eq *EquipmentQuery) WithHyperlinks(opts ...func(*HyperlinkQuery)) *EquipmentQuery {
	query := &HyperlinkQuery{config: eq.config}
	for _, opt := range opts {
		opt(query)
	}
	eq.withHyperlinks = query
	return eq
}

//  WithEndpoints tells the query-builder to eager-loads the nodes that are connected to
// the "endpoints" edge. The optional arguments used to configure the query builder of the edge.
func (eq *EquipmentQuery) WithEndpoints(opts ...func(*ServiceEndpointQuery)) *EquipmentQuery {
	query := &ServiceEndpointQuery{config: eq.config}
	for _, opt := range opts {
		opt(query)
	}
	eq.withEndpoints = query
	return eq
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
//	client.Equipment.Query().
//		GroupBy(equipment.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (eq *EquipmentQuery) GroupBy(field string, fields ...string) *EquipmentGroupBy {
	group := &EquipmentGroupBy{config: eq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return eq.sqlQuery(), nil
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
//	client.Equipment.Query().
//		Select(equipment.FieldCreateTime).
//		Scan(ctx, &v)
//
func (eq *EquipmentQuery) Select(field string, fields ...string) *EquipmentSelect {
	selector := &EquipmentSelect{config: eq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return eq.sqlQuery(), nil
	}
	return selector
}

func (eq *EquipmentQuery) prepareQuery(ctx context.Context) error {
	if eq.path != nil {
		prev, err := eq.path(ctx)
		if err != nil {
			return err
		}
		eq.sql = prev
	}
	return nil
}

func (eq *EquipmentQuery) sqlAll(ctx context.Context) ([]*Equipment, error) {
	var (
		nodes       = []*Equipment{}
		withFKs     = eq.withFKs
		_spec       = eq.querySpec()
		loadedTypes = [10]bool{
			eq.withType != nil,
			eq.withLocation != nil,
			eq.withParentPosition != nil,
			eq.withPositions != nil,
			eq.withPorts != nil,
			eq.withWorkOrder != nil,
			eq.withProperties != nil,
			eq.withFiles != nil,
			eq.withHyperlinks != nil,
			eq.withEndpoints != nil,
		}
	)
	if eq.withType != nil || eq.withLocation != nil || eq.withParentPosition != nil || eq.withWorkOrder != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, equipment.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &Equipment{config: eq.config}
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
	if err := sqlgraph.QueryNodes(ctx, eq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := eq.withType; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*Equipment)
		for i := range nodes {
			if fk := nodes[i].equipment_type; fk != nil {
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
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_type" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Type = n
			}
		}
	}

	if query := eq.withLocation; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*Equipment)
		for i := range nodes {
			if fk := nodes[i].location_equipment; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(location.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "location_equipment" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Location = n
			}
		}
	}

	if query := eq.withParentPosition; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*Equipment)
		for i := range nodes {
			if fk := nodes[i].equipment_position_attachment; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(equipmentposition.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_position_attachment" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.ParentPosition = n
			}
		}
	}

	if query := eq.withPositions; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*Equipment)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.EquipmentPosition(func(s *sql.Selector) {
			s.Where(sql.InValues(equipment.PositionsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.equipment_positions
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "equipment_positions" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_positions" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Positions = append(node.Edges.Positions, n)
		}
	}

	if query := eq.withPorts; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*Equipment)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.EquipmentPort(func(s *sql.Selector) {
			s.Where(sql.InValues(equipment.PortsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.equipment_ports
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "equipment_ports" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_ports" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Ports = append(node.Edges.Ports, n)
		}
	}

	if query := eq.withWorkOrder; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*Equipment)
		for i := range nodes {
			if fk := nodes[i].equipment_work_order; fk != nil {
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
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_work_order" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.WorkOrder = n
			}
		}
	}

	if query := eq.withProperties; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*Equipment)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.Property(func(s *sql.Selector) {
			s.Where(sql.InValues(equipment.PropertiesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.equipment_properties
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "equipment_properties" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_properties" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Properties = append(node.Edges.Properties, n)
		}
	}

	if query := eq.withFiles; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*Equipment)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.File(func(s *sql.Selector) {
			s.Where(sql.InValues(equipment.FilesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.equipment_files
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "equipment_files" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_files" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Files = append(node.Edges.Files, n)
		}
	}

	if query := eq.withHyperlinks; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*Equipment)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.Hyperlink(func(s *sql.Selector) {
			s.Where(sql.InValues(equipment.HyperlinksColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.equipment_hyperlinks
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "equipment_hyperlinks" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_hyperlinks" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Hyperlinks = append(node.Edges.Hyperlinks, n)
		}
	}

	if query := eq.withEndpoints; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*Equipment)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.ServiceEndpoint(func(s *sql.Selector) {
			s.Where(sql.InValues(equipment.EndpointsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.service_endpoint_equipment
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "service_endpoint_equipment" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "service_endpoint_equipment" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Endpoints = append(node.Edges.Endpoints, n)
		}
	}

	return nodes, nil
}

func (eq *EquipmentQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := eq.querySpec()
	return sqlgraph.CountNodes(ctx, eq.driver, _spec)
}

func (eq *EquipmentQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := eq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (eq *EquipmentQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipment.Table,
			Columns: equipment.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipment.FieldID,
			},
		},
		From:   eq.sql,
		Unique: true,
	}
	if ps := eq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := eq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := eq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := eq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (eq *EquipmentQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(eq.driver.Dialect())
	t1 := builder.Table(equipment.Table)
	selector := builder.Select(t1.Columns(equipment.Columns...)...).From(t1)
	if eq.sql != nil {
		selector = eq.sql
		selector.Select(selector.Columns(equipment.Columns...)...)
	}
	for _, p := range eq.predicates {
		p(selector)
	}
	for _, p := range eq.order {
		p(selector)
	}
	if offset := eq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := eq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// EquipmentGroupBy is the builder for group-by Equipment entities.
type EquipmentGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (egb *EquipmentGroupBy) Aggregate(fns ...Aggregate) *EquipmentGroupBy {
	egb.fns = append(egb.fns, fns...)
	return egb
}

// Scan applies the group-by query and scan the result into the given value.
func (egb *EquipmentGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := egb.path(ctx)
	if err != nil {
		return err
	}
	egb.sql = query
	return egb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (egb *EquipmentGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := egb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (egb *EquipmentGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(egb.fields) > 1 {
		return nil, errors.New("ent: EquipmentGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := egb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (egb *EquipmentGroupBy) StringsX(ctx context.Context) []string {
	v, err := egb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (egb *EquipmentGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(egb.fields) > 1 {
		return nil, errors.New("ent: EquipmentGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := egb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (egb *EquipmentGroupBy) IntsX(ctx context.Context) []int {
	v, err := egb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (egb *EquipmentGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(egb.fields) > 1 {
		return nil, errors.New("ent: EquipmentGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := egb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (egb *EquipmentGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := egb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (egb *EquipmentGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(egb.fields) > 1 {
		return nil, errors.New("ent: EquipmentGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := egb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (egb *EquipmentGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := egb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (egb *EquipmentGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := egb.sqlQuery().Query()
	if err := egb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (egb *EquipmentGroupBy) sqlQuery() *sql.Selector {
	selector := egb.sql
	columns := make([]string, 0, len(egb.fields)+len(egb.fns))
	columns = append(columns, egb.fields...)
	for _, fn := range egb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(egb.fields...)
}

// EquipmentSelect is the builder for select fields of Equipment entities.
type EquipmentSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (es *EquipmentSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := es.path(ctx)
	if err != nil {
		return err
	}
	es.sql = query
	return es.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (es *EquipmentSelect) ScanX(ctx context.Context, v interface{}) {
	if err := es.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (es *EquipmentSelect) Strings(ctx context.Context) ([]string, error) {
	if len(es.fields) > 1 {
		return nil, errors.New("ent: EquipmentSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := es.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (es *EquipmentSelect) StringsX(ctx context.Context) []string {
	v, err := es.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (es *EquipmentSelect) Ints(ctx context.Context) ([]int, error) {
	if len(es.fields) > 1 {
		return nil, errors.New("ent: EquipmentSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := es.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (es *EquipmentSelect) IntsX(ctx context.Context) []int {
	v, err := es.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (es *EquipmentSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(es.fields) > 1 {
		return nil, errors.New("ent: EquipmentSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := es.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (es *EquipmentSelect) Float64sX(ctx context.Context) []float64 {
	v, err := es.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (es *EquipmentSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(es.fields) > 1 {
		return nil, errors.New("ent: EquipmentSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := es.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (es *EquipmentSelect) BoolsX(ctx context.Context) []bool {
	v, err := es.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (es *EquipmentSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := es.sqlQuery().Query()
	if err := es.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (es *EquipmentSelect) sqlQuery() sql.Querier {
	selector := es.sql
	selector.Select(selector.Columns(es.fields...)...)
	return selector
}
