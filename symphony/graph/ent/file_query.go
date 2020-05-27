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
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/floorplan"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// FileQuery is the builder for querying File entities.
type FileQuery struct {
	config
	limit      *int
	offset     *int
	order      []OrderFunc
	unique     []string
	predicates []predicate.File
	// eager-loading edges.
	withLocation            *LocationQuery
	withEquipment           *EquipmentQuery
	withUser                *UserQuery
	withWorkOrder           *WorkOrderQuery
	withChecklistItem       *CheckListItemQuery
	withSurvey              *SurveyQuery
	withFloorPlan           *FloorPlanQuery
	withPhotoSurveyQuestion *SurveyQuestionQuery
	withSurveyQuestion      *SurveyQuestionQuery
	withFKs                 bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (fq *FileQuery) Where(ps ...predicate.File) *FileQuery {
	fq.predicates = append(fq.predicates, ps...)
	return fq
}

// Limit adds a limit step to the query.
func (fq *FileQuery) Limit(limit int) *FileQuery {
	fq.limit = &limit
	return fq
}

// Offset adds an offset step to the query.
func (fq *FileQuery) Offset(offset int) *FileQuery {
	fq.offset = &offset
	return fq
}

// Order adds an order step to the query.
func (fq *FileQuery) Order(o ...OrderFunc) *FileQuery {
	fq.order = append(fq.order, o...)
	return fq
}

// QueryLocation chains the current query on the location edge.
func (fq *FileQuery) QueryLocation() *LocationQuery {
	query := &LocationQuery{config: fq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := fq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(file.Table, file.FieldID, fq.sqlQuery()),
			sqlgraph.To(location.Table, location.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, file.LocationTable, file.LocationColumn),
		)
		fromU = sqlgraph.SetNeighbors(fq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryEquipment chains the current query on the equipment edge.
func (fq *FileQuery) QueryEquipment() *EquipmentQuery {
	query := &EquipmentQuery{config: fq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := fq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(file.Table, file.FieldID, fq.sqlQuery()),
			sqlgraph.To(equipment.Table, equipment.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, file.EquipmentTable, file.EquipmentColumn),
		)
		fromU = sqlgraph.SetNeighbors(fq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryUser chains the current query on the user edge.
func (fq *FileQuery) QueryUser() *UserQuery {
	query := &UserQuery{config: fq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := fq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(file.Table, file.FieldID, fq.sqlQuery()),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, file.UserTable, file.UserColumn),
		)
		fromU = sqlgraph.SetNeighbors(fq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryWorkOrder chains the current query on the work_order edge.
func (fq *FileQuery) QueryWorkOrder() *WorkOrderQuery {
	query := &WorkOrderQuery{config: fq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := fq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(file.Table, file.FieldID, fq.sqlQuery()),
			sqlgraph.To(workorder.Table, workorder.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, file.WorkOrderTable, file.WorkOrderColumn),
		)
		fromU = sqlgraph.SetNeighbors(fq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryChecklistItem chains the current query on the checklist_item edge.
func (fq *FileQuery) QueryChecklistItem() *CheckListItemQuery {
	query := &CheckListItemQuery{config: fq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := fq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(file.Table, file.FieldID, fq.sqlQuery()),
			sqlgraph.To(checklistitem.Table, checklistitem.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, file.ChecklistItemTable, file.ChecklistItemColumn),
		)
		fromU = sqlgraph.SetNeighbors(fq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QuerySurvey chains the current query on the survey edge.
func (fq *FileQuery) QuerySurvey() *SurveyQuery {
	query := &SurveyQuery{config: fq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := fq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(file.Table, file.FieldID, fq.sqlQuery()),
			sqlgraph.To(survey.Table, survey.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, file.SurveyTable, file.SurveyColumn),
		)
		fromU = sqlgraph.SetNeighbors(fq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryFloorPlan chains the current query on the floor_plan edge.
func (fq *FileQuery) QueryFloorPlan() *FloorPlanQuery {
	query := &FloorPlanQuery{config: fq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := fq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(file.Table, file.FieldID, fq.sqlQuery()),
			sqlgraph.To(floorplan.Table, floorplan.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, file.FloorPlanTable, file.FloorPlanColumn),
		)
		fromU = sqlgraph.SetNeighbors(fq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryPhotoSurveyQuestion chains the current query on the photo_survey_question edge.
func (fq *FileQuery) QueryPhotoSurveyQuestion() *SurveyQuestionQuery {
	query := &SurveyQuestionQuery{config: fq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := fq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(file.Table, file.FieldID, fq.sqlQuery()),
			sqlgraph.To(surveyquestion.Table, surveyquestion.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, file.PhotoSurveyQuestionTable, file.PhotoSurveyQuestionColumn),
		)
		fromU = sqlgraph.SetNeighbors(fq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QuerySurveyQuestion chains the current query on the survey_question edge.
func (fq *FileQuery) QuerySurveyQuestion() *SurveyQuestionQuery {
	query := &SurveyQuestionQuery{config: fq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := fq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(file.Table, file.FieldID, fq.sqlQuery()),
			sqlgraph.To(surveyquestion.Table, surveyquestion.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, file.SurveyQuestionTable, file.SurveyQuestionColumn),
		)
		fromU = sqlgraph.SetNeighbors(fq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first File entity in the query. Returns *NotFoundError when no file was found.
func (fq *FileQuery) First(ctx context.Context) (*File, error) {
	fs, err := fq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(fs) == 0 {
		return nil, &NotFoundError{file.Label}
	}
	return fs[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (fq *FileQuery) FirstX(ctx context.Context) *File {
	f, err := fq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return f
}

// FirstID returns the first File id in the query. Returns *NotFoundError when no id was found.
func (fq *FileQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = fq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{file.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (fq *FileQuery) FirstXID(ctx context.Context) int {
	id, err := fq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only File entity in the query, returns an error if not exactly one entity was returned.
func (fq *FileQuery) Only(ctx context.Context) (*File, error) {
	fs, err := fq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(fs) {
	case 1:
		return fs[0], nil
	case 0:
		return nil, &NotFoundError{file.Label}
	default:
		return nil, &NotSingularError{file.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (fq *FileQuery) OnlyX(ctx context.Context) *File {
	f, err := fq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return f
}

// OnlyID returns the only File id in the query, returns an error if not exactly one id was returned.
func (fq *FileQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = fq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{file.Label}
	default:
		err = &NotSingularError{file.Label}
	}
	return
}

// OnlyXID is like OnlyID, but panics if an error occurs.
func (fq *FileQuery) OnlyXID(ctx context.Context) int {
	id, err := fq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Files.
func (fq *FileQuery) All(ctx context.Context) ([]*File, error) {
	if err := fq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return fq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (fq *FileQuery) AllX(ctx context.Context) []*File {
	fs, err := fq.All(ctx)
	if err != nil {
		panic(err)
	}
	return fs
}

// IDs executes the query and returns a list of File ids.
func (fq *FileQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := fq.Select(file.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (fq *FileQuery) IDsX(ctx context.Context) []int {
	ids, err := fq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (fq *FileQuery) Count(ctx context.Context) (int, error) {
	if err := fq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return fq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (fq *FileQuery) CountX(ctx context.Context) int {
	count, err := fq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (fq *FileQuery) Exist(ctx context.Context) (bool, error) {
	if err := fq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return fq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (fq *FileQuery) ExistX(ctx context.Context) bool {
	exist, err := fq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (fq *FileQuery) Clone() *FileQuery {
	return &FileQuery{
		config:     fq.config,
		limit:      fq.limit,
		offset:     fq.offset,
		order:      append([]OrderFunc{}, fq.order...),
		unique:     append([]string{}, fq.unique...),
		predicates: append([]predicate.File{}, fq.predicates...),
		// clone intermediate query.
		sql:  fq.sql.Clone(),
		path: fq.path,
	}
}

//  WithLocation tells the query-builder to eager-loads the nodes that are connected to
// the "location" edge. The optional arguments used to configure the query builder of the edge.
func (fq *FileQuery) WithLocation(opts ...func(*LocationQuery)) *FileQuery {
	query := &LocationQuery{config: fq.config}
	for _, opt := range opts {
		opt(query)
	}
	fq.withLocation = query
	return fq
}

//  WithEquipment tells the query-builder to eager-loads the nodes that are connected to
// the "equipment" edge. The optional arguments used to configure the query builder of the edge.
func (fq *FileQuery) WithEquipment(opts ...func(*EquipmentQuery)) *FileQuery {
	query := &EquipmentQuery{config: fq.config}
	for _, opt := range opts {
		opt(query)
	}
	fq.withEquipment = query
	return fq
}

//  WithUser tells the query-builder to eager-loads the nodes that are connected to
// the "user" edge. The optional arguments used to configure the query builder of the edge.
func (fq *FileQuery) WithUser(opts ...func(*UserQuery)) *FileQuery {
	query := &UserQuery{config: fq.config}
	for _, opt := range opts {
		opt(query)
	}
	fq.withUser = query
	return fq
}

//  WithWorkOrder tells the query-builder to eager-loads the nodes that are connected to
// the "work_order" edge. The optional arguments used to configure the query builder of the edge.
func (fq *FileQuery) WithWorkOrder(opts ...func(*WorkOrderQuery)) *FileQuery {
	query := &WorkOrderQuery{config: fq.config}
	for _, opt := range opts {
		opt(query)
	}
	fq.withWorkOrder = query
	return fq
}

//  WithChecklistItem tells the query-builder to eager-loads the nodes that are connected to
// the "checklist_item" edge. The optional arguments used to configure the query builder of the edge.
func (fq *FileQuery) WithChecklistItem(opts ...func(*CheckListItemQuery)) *FileQuery {
	query := &CheckListItemQuery{config: fq.config}
	for _, opt := range opts {
		opt(query)
	}
	fq.withChecklistItem = query
	return fq
}

//  WithSurvey tells the query-builder to eager-loads the nodes that are connected to
// the "survey" edge. The optional arguments used to configure the query builder of the edge.
func (fq *FileQuery) WithSurvey(opts ...func(*SurveyQuery)) *FileQuery {
	query := &SurveyQuery{config: fq.config}
	for _, opt := range opts {
		opt(query)
	}
	fq.withSurvey = query
	return fq
}

//  WithFloorPlan tells the query-builder to eager-loads the nodes that are connected to
// the "floor_plan" edge. The optional arguments used to configure the query builder of the edge.
func (fq *FileQuery) WithFloorPlan(opts ...func(*FloorPlanQuery)) *FileQuery {
	query := &FloorPlanQuery{config: fq.config}
	for _, opt := range opts {
		opt(query)
	}
	fq.withFloorPlan = query
	return fq
}

//  WithPhotoSurveyQuestion tells the query-builder to eager-loads the nodes that are connected to
// the "photo_survey_question" edge. The optional arguments used to configure the query builder of the edge.
func (fq *FileQuery) WithPhotoSurveyQuestion(opts ...func(*SurveyQuestionQuery)) *FileQuery {
	query := &SurveyQuestionQuery{config: fq.config}
	for _, opt := range opts {
		opt(query)
	}
	fq.withPhotoSurveyQuestion = query
	return fq
}

//  WithSurveyQuestion tells the query-builder to eager-loads the nodes that are connected to
// the "survey_question" edge. The optional arguments used to configure the query builder of the edge.
func (fq *FileQuery) WithSurveyQuestion(opts ...func(*SurveyQuestionQuery)) *FileQuery {
	query := &SurveyQuestionQuery{config: fq.config}
	for _, opt := range opts {
		opt(query)
	}
	fq.withSurveyQuestion = query
	return fq
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
//	client.File.Query().
//		GroupBy(file.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (fq *FileQuery) GroupBy(field string, fields ...string) *FileGroupBy {
	group := &FileGroupBy{config: fq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := fq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return fq.sqlQuery(), nil
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
//	client.File.Query().
//		Select(file.FieldCreateTime).
//		Scan(ctx, &v)
//
func (fq *FileQuery) Select(field string, fields ...string) *FileSelect {
	selector := &FileSelect{config: fq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := fq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return fq.sqlQuery(), nil
	}
	return selector
}

func (fq *FileQuery) prepareQuery(ctx context.Context) error {
	if fq.path != nil {
		prev, err := fq.path(ctx)
		if err != nil {
			return err
		}
		fq.sql = prev
	}
	if err := file.Policy.EvalQuery(ctx, fq); err != nil {
		return err
	}
	return nil
}

func (fq *FileQuery) sqlAll(ctx context.Context) ([]*File, error) {
	var (
		nodes       = []*File{}
		withFKs     = fq.withFKs
		_spec       = fq.querySpec()
		loadedTypes = [9]bool{
			fq.withLocation != nil,
			fq.withEquipment != nil,
			fq.withUser != nil,
			fq.withWorkOrder != nil,
			fq.withChecklistItem != nil,
			fq.withSurvey != nil,
			fq.withFloorPlan != nil,
			fq.withPhotoSurveyQuestion != nil,
			fq.withSurveyQuestion != nil,
		}
	)
	if fq.withLocation != nil || fq.withEquipment != nil || fq.withUser != nil || fq.withWorkOrder != nil || fq.withChecklistItem != nil || fq.withSurvey != nil || fq.withFloorPlan != nil || fq.withPhotoSurveyQuestion != nil || fq.withSurveyQuestion != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, file.ForeignKeys...)
	}
	_spec.ScanValues = func() []interface{} {
		node := &File{config: fq.config}
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
	if err := sqlgraph.QueryNodes(ctx, fq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := fq.withLocation; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*File)
		for i := range nodes {
			if fk := nodes[i].location_files; fk != nil {
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
				return nil, fmt.Errorf(`unexpected foreign-key "location_files" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Location = n
			}
		}
	}

	if query := fq.withEquipment; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*File)
		for i := range nodes {
			if fk := nodes[i].equipment_files; fk != nil {
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
				return nil, fmt.Errorf(`unexpected foreign-key "equipment_files" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Equipment = n
			}
		}
	}

	if query := fq.withUser; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*File)
		for i := range nodes {
			if fk := nodes[i].user_profile_photo; fk != nil {
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
				return nil, fmt.Errorf(`unexpected foreign-key "user_profile_photo" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.User = n
			}
		}
	}

	if query := fq.withWorkOrder; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*File)
		for i := range nodes {
			if fk := nodes[i].work_order_files; fk != nil {
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
				return nil, fmt.Errorf(`unexpected foreign-key "work_order_files" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.WorkOrder = n
			}
		}
	}

	if query := fq.withChecklistItem; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*File)
		for i := range nodes {
			if fk := nodes[i].check_list_item_files; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(checklistitem.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "check_list_item_files" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.ChecklistItem = n
			}
		}
	}

	if query := fq.withSurvey; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*File)
		for i := range nodes {
			if fk := nodes[i].survey_source_file; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(survey.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "survey_source_file" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Survey = n
			}
		}
	}

	if query := fq.withFloorPlan; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*File)
		for i := range nodes {
			if fk := nodes[i].floor_plan_image; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(floorplan.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "floor_plan_image" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.FloorPlan = n
			}
		}
	}

	if query := fq.withPhotoSurveyQuestion; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*File)
		for i := range nodes {
			if fk := nodes[i].survey_question_photo_data; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(surveyquestion.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "survey_question_photo_data" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.PhotoSurveyQuestion = n
			}
		}
	}

	if query := fq.withSurveyQuestion; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*File)
		for i := range nodes {
			if fk := nodes[i].survey_question_images; fk != nil {
				ids = append(ids, *fk)
				nodeids[*fk] = append(nodeids[*fk], nodes[i])
			}
		}
		query.Where(surveyquestion.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "survey_question_images" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.SurveyQuestion = n
			}
		}
	}

	return nodes, nil
}

func (fq *FileQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := fq.querySpec()
	return sqlgraph.CountNodes(ctx, fq.driver, _spec)
}

func (fq *FileQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := fq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (fq *FileQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   file.Table,
			Columns: file.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: file.FieldID,
			},
		},
		From:   fq.sql,
		Unique: true,
	}
	if ps := fq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := fq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := fq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := fq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (fq *FileQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(fq.driver.Dialect())
	t1 := builder.Table(file.Table)
	selector := builder.Select(t1.Columns(file.Columns...)...).From(t1)
	if fq.sql != nil {
		selector = fq.sql
		selector.Select(selector.Columns(file.Columns...)...)
	}
	for _, p := range fq.predicates {
		p(selector)
	}
	for _, p := range fq.order {
		p(selector)
	}
	if offset := fq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := fq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// FileGroupBy is the builder for group-by File entities.
type FileGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (fgb *FileGroupBy) Aggregate(fns ...AggregateFunc) *FileGroupBy {
	fgb.fns = append(fgb.fns, fns...)
	return fgb
}

// Scan applies the group-by query and scan the result into the given value.
func (fgb *FileGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := fgb.path(ctx)
	if err != nil {
		return err
	}
	fgb.sql = query
	return fgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (fgb *FileGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := fgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (fgb *FileGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(fgb.fields) > 1 {
		return nil, errors.New("ent: FileGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := fgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (fgb *FileGroupBy) StringsX(ctx context.Context) []string {
	v, err := fgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (fgb *FileGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(fgb.fields) > 1 {
		return nil, errors.New("ent: FileGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := fgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (fgb *FileGroupBy) IntsX(ctx context.Context) []int {
	v, err := fgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (fgb *FileGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(fgb.fields) > 1 {
		return nil, errors.New("ent: FileGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := fgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (fgb *FileGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := fgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (fgb *FileGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(fgb.fields) > 1 {
		return nil, errors.New("ent: FileGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := fgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (fgb *FileGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := fgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (fgb *FileGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := fgb.sqlQuery().Query()
	if err := fgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (fgb *FileGroupBy) sqlQuery() *sql.Selector {
	selector := fgb.sql
	columns := make([]string, 0, len(fgb.fields)+len(fgb.fns))
	columns = append(columns, fgb.fields...)
	for _, fn := range fgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(fgb.fields...)
}

// FileSelect is the builder for select fields of File entities.
type FileSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (fs *FileSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := fs.path(ctx)
	if err != nil {
		return err
	}
	fs.sql = query
	return fs.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (fs *FileSelect) ScanX(ctx context.Context, v interface{}) {
	if err := fs.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (fs *FileSelect) Strings(ctx context.Context) ([]string, error) {
	if len(fs.fields) > 1 {
		return nil, errors.New("ent: FileSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := fs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (fs *FileSelect) StringsX(ctx context.Context) []string {
	v, err := fs.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (fs *FileSelect) Ints(ctx context.Context) ([]int, error) {
	if len(fs.fields) > 1 {
		return nil, errors.New("ent: FileSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := fs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (fs *FileSelect) IntsX(ctx context.Context) []int {
	v, err := fs.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (fs *FileSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(fs.fields) > 1 {
		return nil, errors.New("ent: FileSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := fs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (fs *FileSelect) Float64sX(ctx context.Context) []float64 {
	v, err := fs.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (fs *FileSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(fs.fields) > 1 {
		return nil, errors.New("ent: FileSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := fs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (fs *FileSelect) BoolsX(ctx context.Context) []bool {
	v, err := fs.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (fs *FileSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := fs.sqlQuery().Query()
	if err := fs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (fs *FileSelect) sqlQuery() sql.Querier {
	selector := fs.sql
	selector.Select(selector.Columns(fs.fields...)...)
	return selector
}
