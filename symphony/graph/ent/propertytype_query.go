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
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// PropertyTypeQuery is the builder for querying PropertyType entities.
type PropertyTypeQuery struct {
	config
	limit      *int
	offset     *int
	order      []Order
	unique     []string
	predicates []predicate.PropertyType
	// eager-loading edges.
	withProperties            *PropertyQuery
	withLocationType          *LocationTypeQuery
	withEquipmentPortType     *EquipmentPortTypeQuery
	withLinkEquipmentPortType *EquipmentPortTypeQuery
	withEquipmentType         *EquipmentTypeQuery
	withServiceType           *ServiceTypeQuery
	withWorkOrderType         *WorkOrderTypeQuery
	withProjectType           *ProjectTypeQuery
	withFKs                   bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (ptq *PropertyTypeQuery) Where(ps ...predicate.PropertyType) *PropertyTypeQuery {
	ptq.predicates = append(ptq.predicates, ps...)
	return ptq
}

// Limit adds a limit step to the query.
func (ptq *PropertyTypeQuery) Limit(limit int) *PropertyTypeQuery {
	ptq.limit = &limit
	return ptq
}

// Offset adds an offset step to the query.
func (ptq *PropertyTypeQuery) Offset(offset int) *PropertyTypeQuery {
	ptq.offset = &offset
	return ptq
}

// Order adds an order step to the query.
func (ptq *PropertyTypeQuery) Order(o ...Order) *PropertyTypeQuery {
	ptq.order = append(ptq.order, o...)
	return ptq
}

// QueryProperties chains the current query on the properties edge.
func (ptq *PropertyTypeQuery) QueryProperties() *PropertyQuery {
	query := &PropertyQuery{config: ptq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, ptq.sqlQuery()),
			sqlgraph.To(property.Table, property.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, propertytype.PropertiesTable, propertytype.PropertiesColumn),
		)
		fromU = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryLocationType chains the current query on the location_type edge.
func (ptq *PropertyTypeQuery) QueryLocationType() *LocationTypeQuery {
	query := &LocationTypeQuery{config: ptq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, ptq.sqlQuery()),
			sqlgraph.To(locationtype.Table, locationtype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, propertytype.LocationTypeTable, propertytype.LocationTypeColumn),
		)
		fromU = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryEquipmentPortType chains the current query on the equipment_port_type edge.
func (ptq *PropertyTypeQuery) QueryEquipmentPortType() *EquipmentPortTypeQuery {
	query := &EquipmentPortTypeQuery{config: ptq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, ptq.sqlQuery()),
			sqlgraph.To(equipmentporttype.Table, equipmentporttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, propertytype.EquipmentPortTypeTable, propertytype.EquipmentPortTypeColumn),
		)
		fromU = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryLinkEquipmentPortType chains the current query on the link_equipment_port_type edge.
func (ptq *PropertyTypeQuery) QueryLinkEquipmentPortType() *EquipmentPortTypeQuery {
	query := &EquipmentPortTypeQuery{config: ptq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, ptq.sqlQuery()),
			sqlgraph.To(equipmentporttype.Table, equipmentporttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, propertytype.LinkEquipmentPortTypeTable, propertytype.LinkEquipmentPortTypeColumn),
		)
		fromU = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryEquipmentType chains the current query on the equipment_type edge.
func (ptq *PropertyTypeQuery) QueryEquipmentType() *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: ptq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, ptq.sqlQuery()),
			sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, propertytype.EquipmentTypeTable, propertytype.EquipmentTypeColumn),
		)
		fromU = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryServiceType chains the current query on the service_type edge.
func (ptq *PropertyTypeQuery) QueryServiceType() *ServiceTypeQuery {
	query := &ServiceTypeQuery{config: ptq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, ptq.sqlQuery()),
			sqlgraph.To(servicetype.Table, servicetype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, propertytype.ServiceTypeTable, propertytype.ServiceTypeColumn),
		)
		fromU = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryWorkOrderType chains the current query on the work_order_type edge.
func (ptq *PropertyTypeQuery) QueryWorkOrderType() *WorkOrderTypeQuery {
	query := &WorkOrderTypeQuery{config: ptq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, ptq.sqlQuery()),
			sqlgraph.To(workordertype.Table, workordertype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, propertytype.WorkOrderTypeTable, propertytype.WorkOrderTypeColumn),
		)
		fromU = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryProjectType chains the current query on the project_type edge.
func (ptq *PropertyTypeQuery) QueryProjectType() *ProjectTypeQuery {
	query := &ProjectTypeQuery{config: ptq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, ptq.sqlQuery()),
			sqlgraph.To(projecttype.Table, projecttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, propertytype.ProjectTypeTable, propertytype.ProjectTypeColumn),
		)
		fromU = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first PropertyType entity in the query. Returns *NotFoundError when no propertytype was found.
func (ptq *PropertyTypeQuery) First(ctx context.Context) (*PropertyType, error) {
	pts, err := ptq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(pts) == 0 {
		return nil, &NotFoundError{propertytype.Label}
	}
	return pts[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (ptq *PropertyTypeQuery) FirstX(ctx context.Context) *PropertyType {
	pt, err := ptq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return pt
}

// FirstID returns the first PropertyType id in the query. Returns *NotFoundError when no id was found.
func (ptq *PropertyTypeQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ptq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{propertytype.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (ptq *PropertyTypeQuery) FirstXID(ctx context.Context) int {
	id, err := ptq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only PropertyType entity in the query, returns an error if not exactly one entity was returned.
func (ptq *PropertyTypeQuery) Only(ctx context.Context) (*PropertyType, error) {
	pts, err := ptq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(pts) {
	case 1:
		return pts[0], nil
	case 0:
		return nil, &NotFoundError{propertytype.Label}
	default:
		return nil, &NotSingularError{propertytype.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (ptq *PropertyTypeQuery) OnlyX(ctx context.Context) *PropertyType {
	pt, err := ptq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return pt
}

// OnlyID returns the only PropertyType id in the query, returns an error if not exactly one id was returned.
func (ptq *PropertyTypeQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ptq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{propertytype.Label}
	default:
		err = &NotSingularError{propertytype.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (ptq *PropertyTypeQuery) OnlyXID(ctx context.Context) int {
	id, err := ptq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of PropertyTypes.
func (ptq *PropertyTypeQuery) All(ctx context.Context) ([]*PropertyType, error) {
	if err := ptq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return ptq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (ptq *PropertyTypeQuery) AllX(ctx context.Context) []*PropertyType {
	pts, err := ptq.All(ctx)
	if err != nil {
		panic(err)
	}
	return pts
}

// IDs executes the query and returns a list of PropertyType ids.
func (ptq *PropertyTypeQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := ptq.Select(propertytype.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (ptq *PropertyTypeQuery) IDsX(ctx context.Context) []int {
	ids, err := ptq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (ptq *PropertyTypeQuery) Count(ctx context.Context) (int, error) {
	if err := ptq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return ptq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (ptq *PropertyTypeQuery) CountX(ctx context.Context) int {
	count, err := ptq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (ptq *PropertyTypeQuery) Exist(ctx context.Context) (bool, error) {
	if err := ptq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return ptq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (ptq *PropertyTypeQuery) ExistX(ctx context.Context) bool {
	exist, err := ptq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (ptq *PropertyTypeQuery) Clone() *PropertyTypeQuery {
	return &PropertyTypeQuery{
		config:     ptq.config,
		limit:      ptq.limit,
		offset:     ptq.offset,
		order:      append([]Order{}, ptq.order...),
		unique:     append([]string{}, ptq.unique...),
		predicates: append([]predicate.PropertyType{}, ptq.predicates...),
		// clone intermediate query.
		sql:  ptq.sql.Clone(),
		path: ptq.path,
	}
}

//  WithProperties tells the query-builder to eager-loads the nodes that are connected to
// the "properties" edge. The optional arguments used to configure the query builder of the edge.
func (ptq *PropertyTypeQuery) WithProperties(opts ...func(*PropertyQuery)) *PropertyTypeQuery {
	query := &PropertyQuery{config: ptq.config}
	for _, opt := range opts {
		opt(query)
	}
	ptq.withProperties = query
	return ptq
}

//  WithLocationType tells the query-builder to eager-loads the nodes that are connected to
// the "location_type" edge. The optional arguments used to configure the query builder of the edge.
func (ptq *PropertyTypeQuery) WithLocationType(opts ...func(*LocationTypeQuery)) *PropertyTypeQuery {
	query := &LocationTypeQuery{config: ptq.config}
	for _, opt := range opts {
		opt(query)
	}
	ptq.withLocationType = query
	return ptq
}

//  WithEquipmentPortType tells the query-builder to eager-loads the nodes that are connected to
// the "equipment_port_type" edge. The optional arguments used to configure the query builder of the edge.
func (ptq *PropertyTypeQuery) WithEquipmentPortType(opts ...func(*EquipmentPortTypeQuery)) *PropertyTypeQuery {
	query := &EquipmentPortTypeQuery{config: ptq.config}
	for _, opt := range opts {
		opt(query)
	}
	ptq.withEquipmentPortType = query
	return ptq
}

//  WithLinkEquipmentPortType tells the query-builder to eager-loads the nodes that are connected to
// the "link_equipment_port_type" edge. The optional arguments used to configure the query builder of the edge.
func (ptq *PropertyTypeQuery) WithLinkEquipmentPortType(opts ...func(*EquipmentPortTypeQuery)) *PropertyTypeQuery {
	query := &EquipmentPortTypeQuery{config: ptq.config}
	for _, opt := range opts {
		opt(query)
	}
	ptq.withLinkEquipmentPortType = query
	return ptq
}

//  WithEquipmentType tells the query-builder to eager-loads the nodes that are connected to
// the "equipment_type" edge. The optional arguments used to configure the query builder of the edge.
func (ptq *PropertyTypeQuery) WithEquipmentType(opts ...func(*EquipmentTypeQuery)) *PropertyTypeQuery {
	query := &EquipmentTypeQuery{config: ptq.config}
	for _, opt := range opts {
		opt(query)
	}
	ptq.withEquipmentType = query
	return ptq
}

//  WithServiceType tells the query-builder to eager-loads the nodes that are connected to
// the "service_type" edge. The optional arguments used to configure the query builder of the edge.
func (ptq *PropertyTypeQuery) WithServiceType(opts ...func(*ServiceTypeQuery)) *PropertyTypeQuery {
	query := &ServiceTypeQuery{config: ptq.config}
	for _, opt := range opts {
		opt(query)
	}
	ptq.withServiceType = query
	return ptq
}

//  WithWorkOrderType tells the query-builder to eager-loads the nodes that are connected to
// the "work_order_type" edge. The optional arguments used to configure the query builder of the edge.
func (ptq *PropertyTypeQuery) WithWorkOrderType(opts ...func(*WorkOrderTypeQuery)) *PropertyTypeQuery {
	query := &WorkOrderTypeQuery{config: ptq.config}
	for _, opt := range opts {
		opt(query)
	}
	ptq.withWorkOrderType = query
	return ptq
}

//  WithProjectType tells the query-builder to eager-loads the nodes that are connected to
// the "project_type" edge. The optional arguments used to configure the query builder of the edge.
func (ptq *PropertyTypeQuery) WithProjectType(opts ...func(*ProjectTypeQuery)) *PropertyTypeQuery {
	query := &ProjectTypeQuery{config: ptq.config}
	for _, opt := range opts {
		opt(query)
	}
	ptq.withProjectType = query
	return ptq
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
//	client.PropertyType.Query().
//		GroupBy(propertytype.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (ptq *PropertyTypeQuery) GroupBy(field string, fields ...string) *PropertyTypeGroupBy {
	group := &PropertyTypeGroupBy{config: ptq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return ptq.sqlQuery(), nil
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
//	client.PropertyType.Query().
//		Select(propertytype.FieldCreateTime).
//		Scan(ctx, &v)
//
func (ptq *PropertyTypeQuery) Select(field string, fields ...string) *PropertyTypeSelect {
	selector := &PropertyTypeSelect{config: ptq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return ptq.sqlQuery(), nil
	}
	return selector
}

func (ptq *PropertyTypeQuery) prepareQuery(ctx context.Context) error {
	if ptq.path != nil {
		prev, err := ptq.path(ctx)
		if err != nil {
			return err
		}
		ptq.sql = prev
	}
	return nil
}

func (ptq *PropertyTypeQuery) sqlAll(ctx context.Context) ([]*PropertyType, error) {
	var (
		nodes       = []*PropertyType{}
		withFKs     = ptq.withFKs
		_spec       = ptq.querySpec()
		loadedTypes = [8]bool{
			ptq.withProperties != nil,
			ptq.withLocationType != nil,
			ptq.withEquipmentPortType != nil,
			ptq.withLinkEquipmentPortType != nil,
			ptq.withEquipmentType != nil,
			ptq.withServiceType != nil,
			ptq.withWorkOrderType != nil,
			ptq.withProjectType != nil,
		}
	)
	if ptq.withLocationType != nil || ptq.withEquipmentPortType != nil || ptq.withLinkEquipmentPortType != nil || ptq.withEquipmentType != nil || ptq.withServiceType != nil || ptq.withWorkOrderType != nil || ptq.withProjectType != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, propertytype.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &PropertyType{config: ptq.config}
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
	if err := sqlgraph.QueryNodes(ctx, ptq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := ptq.withProperties; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*PropertyType)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.Property(func(s *sql.Selector) {
			s.Where(sql.InValues(propertytype.PropertiesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.property_type
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "property_type" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "property_type" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Properties = append(node.Edges.Properties, n)
		}
	}

	if query := ptq.withLocationType; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*PropertyType)
		for i := range nodes {
			if fk := nodes[i].location_type_property_types; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(locationtype.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "location_type_property_types" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.LocationType = n
			}
		}
	}

	if query := ptq.withEquipmentPortType; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*PropertyType)
		for i := range nodes {
			if fk := nodes[i].equipment_port_type_property_types; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(equipmentporttype.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_port_type_property_types" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.EquipmentPortType = n
			}
		}
	}

	if query := ptq.withLinkEquipmentPortType; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*PropertyType)
		for i := range nodes {
			if fk := nodes[i].equipment_port_type_link_property_types; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(equipmentporttype.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_port_type_link_property_types" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.LinkEquipmentPortType = n
			}
		}
	}

	if query := ptq.withEquipmentType; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*PropertyType)
		for i := range nodes {
			if fk := nodes[i].equipment_type_property_types; fk != nil {
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
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_type_property_types" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.EquipmentType = n
			}
		}
	}

	if query := ptq.withServiceType; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*PropertyType)
		for i := range nodes {
			if fk := nodes[i].service_type_property_types; fk != nil {
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
				return nil, fmt.Errorf(`unexpected foreign-key "service_type_property_types" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.ServiceType = n
			}
		}
	}

	if query := ptq.withWorkOrderType; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*PropertyType)
		for i := range nodes {
			if fk := nodes[i].work_order_type_property_types; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(workordertype.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_type_property_types" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.WorkOrderType = n
			}
		}
	}

	if query := ptq.withProjectType; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*PropertyType)
		for i := range nodes {
			if fk := nodes[i].project_type_properties; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(projecttype.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "project_type_properties" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.ProjectType = n
			}
		}
	}

	return nodes, nil
}

func (ptq *PropertyTypeQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := ptq.querySpec()
	return sqlgraph.CountNodes(ctx, ptq.driver, _spec)
}

func (ptq *PropertyTypeQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := ptq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (ptq *PropertyTypeQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   propertytype.Table,
			Columns: propertytype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: propertytype.FieldID,
			},
		},
		From:   ptq.sql,
		Unique: true,
	}
	if ps := ptq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := ptq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := ptq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := ptq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (ptq *PropertyTypeQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(ptq.driver.Dialect())
	t1 := builder.Table(propertytype.Table)
	selector := builder.Select(t1.Columns(propertytype.Columns...)...).From(t1)
	if ptq.sql != nil {
		selector = ptq.sql
		selector.Select(selector.Columns(propertytype.Columns...)...)
	}
	for _, p := range ptq.predicates {
		p(selector)
	}
	for _, p := range ptq.order {
		p(selector)
	}
	if offset := ptq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := ptq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// PropertyTypeGroupBy is the builder for group-by PropertyType entities.
type PropertyTypeGroupBy struct {
	config
	fields []string
	fns    []Aggregate
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ptgb *PropertyTypeGroupBy) Aggregate(fns ...Aggregate) *PropertyTypeGroupBy {
	ptgb.fns = append(ptgb.fns, fns...)
	return ptgb
}

// Scan applies the group-by query and scan the result into the given value.
func (ptgb *PropertyTypeGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := ptgb.path(ctx)
	if err != nil {
		return err
	}
	ptgb.sql = query
	return ptgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ptgb *PropertyTypeGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := ptgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (ptgb *PropertyTypeGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(ptgb.fields) > 1 {
		return nil, errors.New("ent: PropertyTypeGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := ptgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ptgb *PropertyTypeGroupBy) StringsX(ctx context.Context) []string {
	v, err := ptgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (ptgb *PropertyTypeGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(ptgb.fields) > 1 {
		return nil, errors.New("ent: PropertyTypeGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := ptgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ptgb *PropertyTypeGroupBy) IntsX(ctx context.Context) []int {
	v, err := ptgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (ptgb *PropertyTypeGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(ptgb.fields) > 1 {
		return nil, errors.New("ent: PropertyTypeGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := ptgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ptgb *PropertyTypeGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := ptgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (ptgb *PropertyTypeGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(ptgb.fields) > 1 {
		return nil, errors.New("ent: PropertyTypeGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := ptgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ptgb *PropertyTypeGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := ptgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ptgb *PropertyTypeGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ptgb.sqlQuery().Query()
	if err := ptgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ptgb *PropertyTypeGroupBy) sqlQuery() *sql.Selector {
	selector := ptgb.sql
	columns := make([]string, 0, len(ptgb.fields)+len(ptgb.fns))
	columns = append(columns, ptgb.fields...)
	for _, fn := range ptgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(ptgb.fields...)
}

// PropertyTypeSelect is the builder for select fields of PropertyType entities.
type PropertyTypeSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (pts *PropertyTypeSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := pts.path(ctx)
	if err != nil {
		return err
	}
	pts.sql = query
	return pts.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (pts *PropertyTypeSelect) ScanX(ctx context.Context, v interface{}) {
	if err := pts.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (pts *PropertyTypeSelect) Strings(ctx context.Context) ([]string, error) {
	if len(pts.fields) > 1 {
		return nil, errors.New("ent: PropertyTypeSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := pts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (pts *PropertyTypeSelect) StringsX(ctx context.Context) []string {
	v, err := pts.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (pts *PropertyTypeSelect) Ints(ctx context.Context) ([]int, error) {
	if len(pts.fields) > 1 {
		return nil, errors.New("ent: PropertyTypeSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := pts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (pts *PropertyTypeSelect) IntsX(ctx context.Context) []int {
	v, err := pts.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (pts *PropertyTypeSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(pts.fields) > 1 {
		return nil, errors.New("ent: PropertyTypeSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := pts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (pts *PropertyTypeSelect) Float64sX(ctx context.Context) []float64 {
	v, err := pts.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (pts *PropertyTypeSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(pts.fields) > 1 {
		return nil, errors.New("ent: PropertyTypeSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := pts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (pts *PropertyTypeSelect) BoolsX(ctx context.Context) []bool {
	v, err := pts.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (pts *PropertyTypeSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := pts.sqlQuery().Query()
	if err := pts.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (pts *PropertyTypeSelect) sqlQuery() sql.Querier {
	selector := pts.sql
	selector.Select(selector.Columns(pts.fields...)...)
	return selector
}
