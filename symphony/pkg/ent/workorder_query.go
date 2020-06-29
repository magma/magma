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
	"github.com/facebookincubator/symphony/pkg/ent/activity"
	"github.com/facebookincubator/symphony/pkg/ent/checklistcategory"
	"github.com/facebookincubator/symphony/pkg/ent/comment"
	"github.com/facebookincubator/symphony/pkg/ent/equipment"
	"github.com/facebookincubator/symphony/pkg/ent/file"
	"github.com/facebookincubator/symphony/pkg/ent/hyperlink"
	"github.com/facebookincubator/symphony/pkg/ent/link"
	"github.com/facebookincubator/symphony/pkg/ent/location"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/project"
	"github.com/facebookincubator/symphony/pkg/ent/property"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/ent/workorder"
	"github.com/facebookincubator/symphony/pkg/ent/workordertemplate"
	"github.com/facebookincubator/symphony/pkg/ent/workordertype"
)

// WorkOrderQuery is the builder for querying WorkOrder entities.
type WorkOrderQuery struct {
	config
	limit      *int
	offset     *int
	order      []OrderFunc
	unique     []string
	predicates []predicate.WorkOrder
	// eager-loading edges.
	withType                *WorkOrderTypeQuery
	withTemplate            *WorkOrderTemplateQuery
	withEquipment           *EquipmentQuery
	withLinks               *LinkQuery
	withFiles               *FileQuery
	withHyperlinks          *HyperlinkQuery
	withLocation            *LocationQuery
	withComments            *CommentQuery
	withActivities          *ActivityQuery
	withProperties          *PropertyQuery
	withCheckListCategories *CheckListCategoryQuery
	withProject             *ProjectQuery
	withOwner               *UserQuery
	withAssignee            *UserQuery
	withFKs                 bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (woq *WorkOrderQuery) Where(ps ...predicate.WorkOrder) *WorkOrderQuery {
	woq.predicates = append(woq.predicates, ps...)
	return woq
}

// Limit adds a limit step to the query.
func (woq *WorkOrderQuery) Limit(limit int) *WorkOrderQuery {
	woq.limit = &limit
	return woq
}

// Offset adds an offset step to the query.
func (woq *WorkOrderQuery) Offset(offset int) *WorkOrderQuery {
	woq.offset = &offset
	return woq
}

// Order adds an order step to the query.
func (woq *WorkOrderQuery) Order(o ...OrderFunc) *WorkOrderQuery {
	woq.order = append(woq.order, o...)
	return woq
}

// QueryType chains the current query on the type edge.
func (woq *WorkOrderQuery) QueryType() *WorkOrderTypeQuery {
	query := &WorkOrderTypeQuery{config: woq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := woq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, woq.sqlQuery()),
			sqlgraph.To(workordertype.Table, workordertype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, workorder.TypeTable, workorder.TypeColumn),
		)
		fromU = sqlgraph.SetNeighbors(woq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryTemplate chains the current query on the template edge.
func (woq *WorkOrderQuery) QueryTemplate() *WorkOrderTemplateQuery {
	query := &WorkOrderTemplateQuery{config: woq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := woq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, woq.sqlQuery()),
			sqlgraph.To(workordertemplate.Table, workordertemplate.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, workorder.TemplateTable, workorder.TemplateColumn),
		)
		fromU = sqlgraph.SetNeighbors(woq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryEquipment chains the current query on the equipment edge.
func (woq *WorkOrderQuery) QueryEquipment() *EquipmentQuery {
	query := &EquipmentQuery{config: woq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := woq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, woq.sqlQuery()),
			sqlgraph.To(equipment.Table, equipment.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, workorder.EquipmentTable, workorder.EquipmentColumn),
		)
		fromU = sqlgraph.SetNeighbors(woq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryLinks chains the current query on the links edge.
func (woq *WorkOrderQuery) QueryLinks() *LinkQuery {
	query := &LinkQuery{config: woq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := woq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, woq.sqlQuery()),
			sqlgraph.To(link.Table, link.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, workorder.LinksTable, workorder.LinksColumn),
		)
		fromU = sqlgraph.SetNeighbors(woq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryFiles chains the current query on the files edge.
func (woq *WorkOrderQuery) QueryFiles() *FileQuery {
	query := &FileQuery{config: woq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := woq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, woq.sqlQuery()),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, workorder.FilesTable, workorder.FilesColumn),
		)
		fromU = sqlgraph.SetNeighbors(woq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryHyperlinks chains the current query on the hyperlinks edge.
func (woq *WorkOrderQuery) QueryHyperlinks() *HyperlinkQuery {
	query := &HyperlinkQuery{config: woq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := woq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, woq.sqlQuery()),
			sqlgraph.To(hyperlink.Table, hyperlink.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, workorder.HyperlinksTable, workorder.HyperlinksColumn),
		)
		fromU = sqlgraph.SetNeighbors(woq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryLocation chains the current query on the location edge.
func (woq *WorkOrderQuery) QueryLocation() *LocationQuery {
	query := &LocationQuery{config: woq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := woq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, woq.sqlQuery()),
			sqlgraph.To(location.Table, location.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, workorder.LocationTable, workorder.LocationColumn),
		)
		fromU = sqlgraph.SetNeighbors(woq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryComments chains the current query on the comments edge.
func (woq *WorkOrderQuery) QueryComments() *CommentQuery {
	query := &CommentQuery{config: woq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := woq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, woq.sqlQuery()),
			sqlgraph.To(comment.Table, comment.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, workorder.CommentsTable, workorder.CommentsColumn),
		)
		fromU = sqlgraph.SetNeighbors(woq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryActivities chains the current query on the activities edge.
func (woq *WorkOrderQuery) QueryActivities() *ActivityQuery {
	query := &ActivityQuery{config: woq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := woq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, woq.sqlQuery()),
			sqlgraph.To(activity.Table, activity.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, workorder.ActivitiesTable, workorder.ActivitiesColumn),
		)
		fromU = sqlgraph.SetNeighbors(woq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryProperties chains the current query on the properties edge.
func (woq *WorkOrderQuery) QueryProperties() *PropertyQuery {
	query := &PropertyQuery{config: woq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := woq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, woq.sqlQuery()),
			sqlgraph.To(property.Table, property.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, workorder.PropertiesTable, workorder.PropertiesColumn),
		)
		fromU = sqlgraph.SetNeighbors(woq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryCheckListCategories chains the current query on the check_list_categories edge.
func (woq *WorkOrderQuery) QueryCheckListCategories() *CheckListCategoryQuery {
	query := &CheckListCategoryQuery{config: woq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := woq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, woq.sqlQuery()),
			sqlgraph.To(checklistcategory.Table, checklistcategory.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, workorder.CheckListCategoriesTable, workorder.CheckListCategoriesColumn),
		)
		fromU = sqlgraph.SetNeighbors(woq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryProject chains the current query on the project edge.
func (woq *WorkOrderQuery) QueryProject() *ProjectQuery {
	query := &ProjectQuery{config: woq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := woq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, woq.sqlQuery()),
			sqlgraph.To(project.Table, project.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, workorder.ProjectTable, workorder.ProjectColumn),
		)
		fromU = sqlgraph.SetNeighbors(woq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryOwner chains the current query on the owner edge.
func (woq *WorkOrderQuery) QueryOwner() *UserQuery {
	query := &UserQuery{config: woq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := woq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, woq.sqlQuery()),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, workorder.OwnerTable, workorder.OwnerColumn),
		)
		fromU = sqlgraph.SetNeighbors(woq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryAssignee chains the current query on the assignee edge.
func (woq *WorkOrderQuery) QueryAssignee() *UserQuery {
	query := &UserQuery{config: woq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := woq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, woq.sqlQuery()),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, workorder.AssigneeTable, workorder.AssigneeColumn),
		)
		fromU = sqlgraph.SetNeighbors(woq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first WorkOrder entity in the query. Returns *NotFoundError when no workorder was found.
func (woq *WorkOrderQuery) First(ctx context.Context) (*WorkOrder, error) {
	wos, err := woq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(wos) == 0 {
		return nil, &NotFoundError{workorder.Label}
	}
	return wos[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (woq *WorkOrderQuery) FirstX(ctx context.Context) *WorkOrder {
	wo, err := woq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return wo
}

// FirstID returns the first WorkOrder id in the query. Returns *NotFoundError when no id was found.
func (woq *WorkOrderQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = woq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{workorder.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (woq *WorkOrderQuery) FirstXID(ctx context.Context) int {
	id, err := woq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only WorkOrder entity in the query, returns an error if not exactly one entity was returned.
func (woq *WorkOrderQuery) Only(ctx context.Context) (*WorkOrder, error) {
	wos, err := woq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(wos) {
	case 1:
		return wos[0], nil
	case 0:
		return nil, &NotFoundError{workorder.Label}
	default:
		return nil, &NotSingularError{workorder.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (woq *WorkOrderQuery) OnlyX(ctx context.Context) *WorkOrder {
	wo, err := woq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return wo
}

// OnlyID returns the only WorkOrder id in the query, returns an error if not exactly one id was returned.
func (woq *WorkOrderQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = woq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{workorder.Label}
	default:
		err = &NotSingularError{workorder.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (woq *WorkOrderQuery) OnlyXID(ctx context.Context) int {
	id, err := woq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of WorkOrders.
func (woq *WorkOrderQuery) All(ctx context.Context) ([]*WorkOrder, error) {
	if err := woq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return woq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (woq *WorkOrderQuery) AllX(ctx context.Context) []*WorkOrder {
	wos, err := woq.All(ctx)
	if err != nil {
		panic(err)
	}
	return wos
}

// IDs executes the query and returns a list of WorkOrder ids.
func (woq *WorkOrderQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := woq.Select(workorder.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (woq *WorkOrderQuery) IDsX(ctx context.Context) []int {
	ids, err := woq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (woq *WorkOrderQuery) Count(ctx context.Context) (int, error) {
	if err := woq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return woq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (woq *WorkOrderQuery) CountX(ctx context.Context) int {
	count, err := woq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (woq *WorkOrderQuery) Exist(ctx context.Context) (bool, error) {
	if err := woq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return woq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (woq *WorkOrderQuery) ExistX(ctx context.Context) bool {
	exist, err := woq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (woq *WorkOrderQuery) Clone() *WorkOrderQuery {
	return &WorkOrderQuery{
		config:     woq.config,
		limit:      woq.limit,
		offset:     woq.offset,
		order:      append([]OrderFunc{}, woq.order...),
		unique:     append([]string{}, woq.unique...),
		predicates: append([]predicate.WorkOrder{}, woq.predicates...),
		// clone intermediate query.
		sql:  woq.sql.Clone(),
		path: woq.path,
	}
}

//  WithType tells the query-builder to eager-loads the nodes that are connected to
// the "type" edge. The optional arguments used to configure the query builder of the edge.
func (woq *WorkOrderQuery) WithType(opts ...func(*WorkOrderTypeQuery)) *WorkOrderQuery {
	query := &WorkOrderTypeQuery{config: woq.config}
	for _, opt := range opts {
		opt(query)
	}
	woq.withType = query
	return woq
}

//  WithTemplate tells the query-builder to eager-loads the nodes that are connected to
// the "template" edge. The optional arguments used to configure the query builder of the edge.
func (woq *WorkOrderQuery) WithTemplate(opts ...func(*WorkOrderTemplateQuery)) *WorkOrderQuery {
	query := &WorkOrderTemplateQuery{config: woq.config}
	for _, opt := range opts {
		opt(query)
	}
	woq.withTemplate = query
	return woq
}

//  WithEquipment tells the query-builder to eager-loads the nodes that are connected to
// the "equipment" edge. The optional arguments used to configure the query builder of the edge.
func (woq *WorkOrderQuery) WithEquipment(opts ...func(*EquipmentQuery)) *WorkOrderQuery {
	query := &EquipmentQuery{config: woq.config}
	for _, opt := range opts {
		opt(query)
	}
	woq.withEquipment = query
	return woq
}

//  WithLinks tells the query-builder to eager-loads the nodes that are connected to
// the "links" edge. The optional arguments used to configure the query builder of the edge.
func (woq *WorkOrderQuery) WithLinks(opts ...func(*LinkQuery)) *WorkOrderQuery {
	query := &LinkQuery{config: woq.config}
	for _, opt := range opts {
		opt(query)
	}
	woq.withLinks = query
	return woq
}

//  WithFiles tells the query-builder to eager-loads the nodes that are connected to
// the "files" edge. The optional arguments used to configure the query builder of the edge.
func (woq *WorkOrderQuery) WithFiles(opts ...func(*FileQuery)) *WorkOrderQuery {
	query := &FileQuery{config: woq.config}
	for _, opt := range opts {
		opt(query)
	}
	woq.withFiles = query
	return woq
}

//  WithHyperlinks tells the query-builder to eager-loads the nodes that are connected to
// the "hyperlinks" edge. The optional arguments used to configure the query builder of the edge.
func (woq *WorkOrderQuery) WithHyperlinks(opts ...func(*HyperlinkQuery)) *WorkOrderQuery {
	query := &HyperlinkQuery{config: woq.config}
	for _, opt := range opts {
		opt(query)
	}
	woq.withHyperlinks = query
	return woq
}

//  WithLocation tells the query-builder to eager-loads the nodes that are connected to
// the "location" edge. The optional arguments used to configure the query builder of the edge.
func (woq *WorkOrderQuery) WithLocation(opts ...func(*LocationQuery)) *WorkOrderQuery {
	query := &LocationQuery{config: woq.config}
	for _, opt := range opts {
		opt(query)
	}
	woq.withLocation = query
	return woq
}

//  WithComments tells the query-builder to eager-loads the nodes that are connected to
// the "comments" edge. The optional arguments used to configure the query builder of the edge.
func (woq *WorkOrderQuery) WithComments(opts ...func(*CommentQuery)) *WorkOrderQuery {
	query := &CommentQuery{config: woq.config}
	for _, opt := range opts {
		opt(query)
	}
	woq.withComments = query
	return woq
}

//  WithActivities tells the query-builder to eager-loads the nodes that are connected to
// the "activities" edge. The optional arguments used to configure the query builder of the edge.
func (woq *WorkOrderQuery) WithActivities(opts ...func(*ActivityQuery)) *WorkOrderQuery {
	query := &ActivityQuery{config: woq.config}
	for _, opt := range opts {
		opt(query)
	}
	woq.withActivities = query
	return woq
}

//  WithProperties tells the query-builder to eager-loads the nodes that are connected to
// the "properties" edge. The optional arguments used to configure the query builder of the edge.
func (woq *WorkOrderQuery) WithProperties(opts ...func(*PropertyQuery)) *WorkOrderQuery {
	query := &PropertyQuery{config: woq.config}
	for _, opt := range opts {
		opt(query)
	}
	woq.withProperties = query
	return woq
}

//  WithCheckListCategories tells the query-builder to eager-loads the nodes that are connected to
// the "check_list_categories" edge. The optional arguments used to configure the query builder of the edge.
func (woq *WorkOrderQuery) WithCheckListCategories(opts ...func(*CheckListCategoryQuery)) *WorkOrderQuery {
	query := &CheckListCategoryQuery{config: woq.config}
	for _, opt := range opts {
		opt(query)
	}
	woq.withCheckListCategories = query
	return woq
}

//  WithProject tells the query-builder to eager-loads the nodes that are connected to
// the "project" edge. The optional arguments used to configure the query builder of the edge.
func (woq *WorkOrderQuery) WithProject(opts ...func(*ProjectQuery)) *WorkOrderQuery {
	query := &ProjectQuery{config: woq.config}
	for _, opt := range opts {
		opt(query)
	}
	woq.withProject = query
	return woq
}

//  WithOwner tells the query-builder to eager-loads the nodes that are connected to
// the "owner" edge. The optional arguments used to configure the query builder of the edge.
func (woq *WorkOrderQuery) WithOwner(opts ...func(*UserQuery)) *WorkOrderQuery {
	query := &UserQuery{config: woq.config}
	for _, opt := range opts {
		opt(query)
	}
	woq.withOwner = query
	return woq
}

//  WithAssignee tells the query-builder to eager-loads the nodes that are connected to
// the "assignee" edge. The optional arguments used to configure the query builder of the edge.
func (woq *WorkOrderQuery) WithAssignee(opts ...func(*UserQuery)) *WorkOrderQuery {
	query := &UserQuery{config: woq.config}
	for _, opt := range opts {
		opt(query)
	}
	woq.withAssignee = query
	return woq
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
//	client.WorkOrder.Query().
//		GroupBy(workorder.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (woq *WorkOrderQuery) GroupBy(field string, fields ...string) *WorkOrderGroupBy {
	group := &WorkOrderGroupBy{config: woq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := woq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return woq.sqlQuery(), nil
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
//	client.WorkOrder.Query().
//		Select(workorder.FieldCreateTime).
//		Scan(ctx, &v)
//
func (woq *WorkOrderQuery) Select(field string, fields ...string) *WorkOrderSelect {
	selector := &WorkOrderSelect{config: woq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := woq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return woq.sqlQuery(), nil
	}
	return selector
}

func (woq *WorkOrderQuery) prepareQuery(ctx context.Context) error {
	if woq.path != nil {
		prev, err := woq.path(ctx)
		if err != nil {
			return err
		}
		woq.sql = prev
	}
	if err := workorder.Policy.EvalQuery(ctx, woq); err != nil {
		return err
	}
	return nil
}

func (woq *WorkOrderQuery) sqlAll(ctx context.Context) ([]*WorkOrder, error) {
	var (
		nodes       = []*WorkOrder{}
		withFKs     = woq.withFKs
		_spec       = woq.querySpec()
		loadedTypes = [14]bool{
			woq.withType != nil,
			woq.withTemplate != nil,
			woq.withEquipment != nil,
			woq.withLinks != nil,
			woq.withFiles != nil,
			woq.withHyperlinks != nil,
			woq.withLocation != nil,
			woq.withComments != nil,
			woq.withActivities != nil,
			woq.withProperties != nil,
			woq.withCheckListCategories != nil,
			woq.withProject != nil,
			woq.withOwner != nil,
			woq.withAssignee != nil,
		}
	)
	if woq.withType != nil || woq.withTemplate != nil || woq.withLocation != nil || woq.withProject != nil || woq.withOwner != nil || woq.withAssignee != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, workorder.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &WorkOrder{config: woq.config}
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
	if err := sqlgraph.QueryNodes(ctx, woq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := woq.withType; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*WorkOrder)
		for i := range nodes {
			if fk := nodes[i].work_order_type; fk != nil {
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
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_type" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Type = n
			}
		}
	}

	if query := woq.withTemplate; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*WorkOrder)
		for i := range nodes {
			if fk := nodes[i].work_order_template; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(workordertemplate.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_template" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Template = n
			}
		}
	}

	if query := woq.withEquipment; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*WorkOrder)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.Equipment(func(s *sql.Selector) {
			s.Where(sql.InValues(workorder.EquipmentColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.equipment_work_order
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "equipment_work_order" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_work_order" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Equipment = append(node.Edges.Equipment, n)
		}
	}

	if query := woq.withLinks; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*WorkOrder)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.Link(func(s *sql.Selector) {
			s.Where(sql.InValues(workorder.LinksColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.link_work_order
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "link_work_order" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "link_work_order" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Links = append(node.Edges.Links, n)
		}
	}

	if query := woq.withFiles; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*WorkOrder)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.File(func(s *sql.Selector) {
			s.Where(sql.InValues(workorder.FilesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.work_order_files
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "work_order_files" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_files" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Files = append(node.Edges.Files, n)
		}
	}

	if query := woq.withHyperlinks; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*WorkOrder)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.Hyperlink(func(s *sql.Selector) {
			s.Where(sql.InValues(workorder.HyperlinksColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.work_order_hyperlinks
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "work_order_hyperlinks" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_hyperlinks" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Hyperlinks = append(node.Edges.Hyperlinks, n)
		}
	}

	if query := woq.withLocation; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*WorkOrder)
		for i := range nodes {
			if fk := nodes[i].work_order_location; fk != nil {
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
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_location" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Location = n
			}
		}
	}

	if query := woq.withComments; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*WorkOrder)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.Comment(func(s *sql.Selector) {
			s.Where(sql.InValues(workorder.CommentsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.work_order_comments
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "work_order_comments" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_comments" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Comments = append(node.Edges.Comments, n)
		}
	}

	if query := woq.withActivities; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*WorkOrder)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.Activity(func(s *sql.Selector) {
			s.Where(sql.InValues(workorder.ActivitiesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.work_order_activities
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "work_order_activities" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_activities" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Activities = append(node.Edges.Activities, n)
		}
	}

	if query := woq.withProperties; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*WorkOrder)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.Property(func(s *sql.Selector) {
			s.Where(sql.InValues(workorder.PropertiesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.work_order_properties
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "work_order_properties" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_properties" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Properties = append(node.Edges.Properties, n)
		}
	}

	if query := woq.withCheckListCategories; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*WorkOrder)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.CheckListCategory(func(s *sql.Selector) {
			s.Where(sql.InValues(workorder.CheckListCategoriesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.work_order_check_list_categories
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "work_order_check_list_categories" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_check_list_categories" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.CheckListCategories = append(node.Edges.CheckListCategories, n)
		}
	}

	if query := woq.withProject; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*WorkOrder)
		for i := range nodes {
			if fk := nodes[i].project_work_orders; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(project.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "project_work_orders" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Project = n
			}
		}
	}

	if query := woq.withOwner; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*WorkOrder)
		for i := range nodes {
			if fk := nodes[i].work_order_owner; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(user.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_owner" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Owner = n
			}
		}
	}

	if query := woq.withAssignee; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*WorkOrder)
		for i := range nodes {
			if fk := nodes[i].work_order_assignee; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(user.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_assignee" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Assignee = n
			}
		}
	}

	return nodes, nil
}

func (woq *WorkOrderQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := woq.querySpec()
	return sqlgraph.CountNodes(ctx, woq.driver, _spec)
}

func (woq *WorkOrderQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := woq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (woq *WorkOrderQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   workorder.Table,
			Columns: workorder.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: workorder.FieldID,
			},
		},
		From:   woq.sql,
		Unique: true,
	}
	if ps := woq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := woq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := woq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := woq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (woq *WorkOrderQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(woq.driver.Dialect())
	t1 := builder.Table(workorder.Table)
	selector := builder.Select(t1.Columns(workorder.Columns...)...).From(t1)
	if woq.sql != nil {
		selector = woq.sql
		selector.Select(selector.Columns(workorder.Columns...)...)
	}
	for _, p := range woq.predicates {
		p(selector)
	}
	for _, p := range woq.order {
		p(selector)
	}
	if offset := woq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := woq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// WorkOrderGroupBy is the builder for group-by WorkOrder entities.
type WorkOrderGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (wogb *WorkOrderGroupBy) Aggregate(fns ...AggregateFunc) *WorkOrderGroupBy {
	wogb.fns = append(wogb.fns, fns...)
	return wogb
}

// Scan applies the group-by query and scan the result into the given value.
func (wogb *WorkOrderGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := wogb.path(ctx)
	if err != nil {
		return err
	}
	wogb.sql = query
	return wogb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (wogb *WorkOrderGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := wogb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (wogb *WorkOrderGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(wogb.fields) > 1 {
		return nil, errors.New("ent: WorkOrderGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := wogb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (wogb *WorkOrderGroupBy) StringsX(ctx context.Context) []string {
	v, err := wogb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (wogb *WorkOrderGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(wogb.fields) > 1 {
		return nil, errors.New("ent: WorkOrderGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := wogb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (wogb *WorkOrderGroupBy) IntsX(ctx context.Context) []int {
	v, err := wogb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (wogb *WorkOrderGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(wogb.fields) > 1 {
		return nil, errors.New("ent: WorkOrderGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := wogb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (wogb *WorkOrderGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := wogb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (wogb *WorkOrderGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(wogb.fields) > 1 {
		return nil, errors.New("ent: WorkOrderGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := wogb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (wogb *WorkOrderGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := wogb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (wogb *WorkOrderGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := wogb.sqlQuery().Query()
	if err := wogb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (wogb *WorkOrderGroupBy) sqlQuery() *sql.Selector {
	selector := wogb.sql
	columns := make([]string, 0, len(wogb.fields)+len(wogb.fns))
	columns = append(columns, wogb.fields...)
	for _, fn := range wogb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(wogb.fields...)
}

// WorkOrderSelect is the builder for select fields of WorkOrder entities.
type WorkOrderSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (wos *WorkOrderSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := wos.path(ctx)
	if err != nil {
		return err
	}
	wos.sql = query
	return wos.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (wos *WorkOrderSelect) ScanX(ctx context.Context, v interface{}) {
	if err := wos.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (wos *WorkOrderSelect) Strings(ctx context.Context) ([]string, error) {
	if len(wos.fields) > 1 {
		return nil, errors.New("ent: WorkOrderSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := wos.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (wos *WorkOrderSelect) StringsX(ctx context.Context) []string {
	v, err := wos.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (wos *WorkOrderSelect) Ints(ctx context.Context) ([]int, error) {
	if len(wos.fields) > 1 {
		return nil, errors.New("ent: WorkOrderSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := wos.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (wos *WorkOrderSelect) IntsX(ctx context.Context) []int {
	v, err := wos.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (wos *WorkOrderSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(wos.fields) > 1 {
		return nil, errors.New("ent: WorkOrderSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := wos.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (wos *WorkOrderSelect) Float64sX(ctx context.Context) []float64 {
	v, err := wos.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (wos *WorkOrderSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(wos.fields) > 1 {
		return nil, errors.New("ent: WorkOrderSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := wos.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (wos *WorkOrderSelect) BoolsX(ctx context.Context) []bool {
	v, err := wos.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (wos *WorkOrderSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := wos.sqlQuery().Query()
	if err := wos.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (wos *WorkOrderSelect) sqlQuery() sql.Querier {
	selector := wos.sql
	selector.Select(selector.Columns(wos.fields...)...)
	return selector
}
