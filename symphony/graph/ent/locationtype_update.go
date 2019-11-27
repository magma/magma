// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatecategory"
)

// LocationTypeUpdate is the builder for updating LocationType entities.
type LocationTypeUpdate struct {
	config

	update_time                     *time.Time
	site                            *bool
	name                            *string
	map_type                        *string
	clearmap_type                   bool
	map_zoom_level                  *int
	addmap_zoom_level               *int
	clearmap_zoom_level             bool
	index                           *int
	addindex                        *int
	locations                       map[string]struct{}
	property_types                  map[string]struct{}
	survey_template_categories      map[string]struct{}
	removedLocations                map[string]struct{}
	removedPropertyTypes            map[string]struct{}
	removedSurveyTemplateCategories map[string]struct{}
	predicates                      []predicate.LocationType
}

// Where adds a new predicate for the builder.
func (ltu *LocationTypeUpdate) Where(ps ...predicate.LocationType) *LocationTypeUpdate {
	ltu.predicates = append(ltu.predicates, ps...)
	return ltu
}

// SetSite sets the site field.
func (ltu *LocationTypeUpdate) SetSite(b bool) *LocationTypeUpdate {
	ltu.site = &b
	return ltu
}

// SetNillableSite sets the site field if the given value is not nil.
func (ltu *LocationTypeUpdate) SetNillableSite(b *bool) *LocationTypeUpdate {
	if b != nil {
		ltu.SetSite(*b)
	}
	return ltu
}

// SetName sets the name field.
func (ltu *LocationTypeUpdate) SetName(s string) *LocationTypeUpdate {
	ltu.name = &s
	return ltu
}

// SetMapType sets the map_type field.
func (ltu *LocationTypeUpdate) SetMapType(s string) *LocationTypeUpdate {
	ltu.map_type = &s
	return ltu
}

// SetNillableMapType sets the map_type field if the given value is not nil.
func (ltu *LocationTypeUpdate) SetNillableMapType(s *string) *LocationTypeUpdate {
	if s != nil {
		ltu.SetMapType(*s)
	}
	return ltu
}

// ClearMapType clears the value of map_type.
func (ltu *LocationTypeUpdate) ClearMapType() *LocationTypeUpdate {
	ltu.map_type = nil
	ltu.clearmap_type = true
	return ltu
}

// SetMapZoomLevel sets the map_zoom_level field.
func (ltu *LocationTypeUpdate) SetMapZoomLevel(i int) *LocationTypeUpdate {
	ltu.map_zoom_level = &i
	ltu.addmap_zoom_level = nil
	return ltu
}

// SetNillableMapZoomLevel sets the map_zoom_level field if the given value is not nil.
func (ltu *LocationTypeUpdate) SetNillableMapZoomLevel(i *int) *LocationTypeUpdate {
	if i != nil {
		ltu.SetMapZoomLevel(*i)
	}
	return ltu
}

// AddMapZoomLevel adds i to map_zoom_level.
func (ltu *LocationTypeUpdate) AddMapZoomLevel(i int) *LocationTypeUpdate {
	if ltu.addmap_zoom_level == nil {
		ltu.addmap_zoom_level = &i
	} else {
		*ltu.addmap_zoom_level += i
	}
	return ltu
}

// ClearMapZoomLevel clears the value of map_zoom_level.
func (ltu *LocationTypeUpdate) ClearMapZoomLevel() *LocationTypeUpdate {
	ltu.map_zoom_level = nil
	ltu.clearmap_zoom_level = true
	return ltu
}

// SetIndex sets the index field.
func (ltu *LocationTypeUpdate) SetIndex(i int) *LocationTypeUpdate {
	ltu.index = &i
	ltu.addindex = nil
	return ltu
}

// SetNillableIndex sets the index field if the given value is not nil.
func (ltu *LocationTypeUpdate) SetNillableIndex(i *int) *LocationTypeUpdate {
	if i != nil {
		ltu.SetIndex(*i)
	}
	return ltu
}

// AddIndex adds i to index.
func (ltu *LocationTypeUpdate) AddIndex(i int) *LocationTypeUpdate {
	if ltu.addindex == nil {
		ltu.addindex = &i
	} else {
		*ltu.addindex += i
	}
	return ltu
}

// AddLocationIDs adds the locations edge to Location by ids.
func (ltu *LocationTypeUpdate) AddLocationIDs(ids ...string) *LocationTypeUpdate {
	if ltu.locations == nil {
		ltu.locations = make(map[string]struct{})
	}
	for i := range ids {
		ltu.locations[ids[i]] = struct{}{}
	}
	return ltu
}

// AddLocations adds the locations edges to Location.
func (ltu *LocationTypeUpdate) AddLocations(l ...*Location) *LocationTypeUpdate {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return ltu.AddLocationIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (ltu *LocationTypeUpdate) AddPropertyTypeIDs(ids ...string) *LocationTypeUpdate {
	if ltu.property_types == nil {
		ltu.property_types = make(map[string]struct{})
	}
	for i := range ids {
		ltu.property_types[ids[i]] = struct{}{}
	}
	return ltu
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (ltu *LocationTypeUpdate) AddPropertyTypes(p ...*PropertyType) *LocationTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ltu.AddPropertyTypeIDs(ids...)
}

// AddSurveyTemplateCategoryIDs adds the survey_template_categories edge to SurveyTemplateCategory by ids.
func (ltu *LocationTypeUpdate) AddSurveyTemplateCategoryIDs(ids ...string) *LocationTypeUpdate {
	if ltu.survey_template_categories == nil {
		ltu.survey_template_categories = make(map[string]struct{})
	}
	for i := range ids {
		ltu.survey_template_categories[ids[i]] = struct{}{}
	}
	return ltu
}

// AddSurveyTemplateCategories adds the survey_template_categories edges to SurveyTemplateCategory.
func (ltu *LocationTypeUpdate) AddSurveyTemplateCategories(s ...*SurveyTemplateCategory) *LocationTypeUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return ltu.AddSurveyTemplateCategoryIDs(ids...)
}

// RemoveLocationIDs removes the locations edge to Location by ids.
func (ltu *LocationTypeUpdate) RemoveLocationIDs(ids ...string) *LocationTypeUpdate {
	if ltu.removedLocations == nil {
		ltu.removedLocations = make(map[string]struct{})
	}
	for i := range ids {
		ltu.removedLocations[ids[i]] = struct{}{}
	}
	return ltu
}

// RemoveLocations removes locations edges to Location.
func (ltu *LocationTypeUpdate) RemoveLocations(l ...*Location) *LocationTypeUpdate {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return ltu.RemoveLocationIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (ltu *LocationTypeUpdate) RemovePropertyTypeIDs(ids ...string) *LocationTypeUpdate {
	if ltu.removedPropertyTypes == nil {
		ltu.removedPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		ltu.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return ltu
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (ltu *LocationTypeUpdate) RemovePropertyTypes(p ...*PropertyType) *LocationTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ltu.RemovePropertyTypeIDs(ids...)
}

// RemoveSurveyTemplateCategoryIDs removes the survey_template_categories edge to SurveyTemplateCategory by ids.
func (ltu *LocationTypeUpdate) RemoveSurveyTemplateCategoryIDs(ids ...string) *LocationTypeUpdate {
	if ltu.removedSurveyTemplateCategories == nil {
		ltu.removedSurveyTemplateCategories = make(map[string]struct{})
	}
	for i := range ids {
		ltu.removedSurveyTemplateCategories[ids[i]] = struct{}{}
	}
	return ltu
}

// RemoveSurveyTemplateCategories removes survey_template_categories edges to SurveyTemplateCategory.
func (ltu *LocationTypeUpdate) RemoveSurveyTemplateCategories(s ...*SurveyTemplateCategory) *LocationTypeUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return ltu.RemoveSurveyTemplateCategoryIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (ltu *LocationTypeUpdate) Save(ctx context.Context) (int, error) {
	if ltu.update_time == nil {
		v := locationtype.UpdateDefaultUpdateTime()
		ltu.update_time = &v
	}
	return ltu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (ltu *LocationTypeUpdate) SaveX(ctx context.Context) int {
	affected, err := ltu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ltu *LocationTypeUpdate) Exec(ctx context.Context) error {
	_, err := ltu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ltu *LocationTypeUpdate) ExecX(ctx context.Context) {
	if err := ltu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ltu *LocationTypeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(ltu.driver.Dialect())
		selector = builder.Select(locationtype.FieldID).From(builder.Table(locationtype.Table))
	)
	for _, p := range ltu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = ltu.driver.Query(ctx, query, args, rows); err != nil {
		return 0, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ent: failed reading id: %v", err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return 0, nil
	}

	tx, err := ltu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(locationtype.Table)
	)
	updater = updater.Where(sql.InInts(locationtype.FieldID, ids...))
	if value := ltu.update_time; value != nil {
		updater.Set(locationtype.FieldUpdateTime, *value)
	}
	if value := ltu.site; value != nil {
		updater.Set(locationtype.FieldSite, *value)
	}
	if value := ltu.name; value != nil {
		updater.Set(locationtype.FieldName, *value)
	}
	if value := ltu.map_type; value != nil {
		updater.Set(locationtype.FieldMapType, *value)
	}
	if ltu.clearmap_type {
		updater.SetNull(locationtype.FieldMapType)
	}
	if value := ltu.map_zoom_level; value != nil {
		updater.Set(locationtype.FieldMapZoomLevel, *value)
	}
	if value := ltu.addmap_zoom_level; value != nil {
		updater.Add(locationtype.FieldMapZoomLevel, *value)
	}
	if ltu.clearmap_zoom_level {
		updater.SetNull(locationtype.FieldMapZoomLevel)
	}
	if value := ltu.index; value != nil {
		updater.Set(locationtype.FieldIndex, *value)
	}
	if value := ltu.addindex; value != nil {
		updater.Add(locationtype.FieldIndex, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ltu.removedLocations) > 0 {
		eids := make([]int, len(ltu.removedLocations))
		for eid := range ltu.removedLocations {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(locationtype.LocationsTable).
			SetNull(locationtype.LocationsColumn).
			Where(sql.InInts(locationtype.LocationsColumn, ids...)).
			Where(sql.InInts(location.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ltu.locations) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range ltu.locations {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(location.FieldID, eid)
			}
			query, args := builder.Update(locationtype.LocationsTable).
				Set(locationtype.LocationsColumn, id).
				Where(sql.And(p, sql.IsNull(locationtype.LocationsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(ltu.locations) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"locations\" %v already connected to a different \"LocationType\"", keys(ltu.locations))})
			}
		}
	}
	if len(ltu.removedPropertyTypes) > 0 {
		eids := make([]int, len(ltu.removedPropertyTypes))
		for eid := range ltu.removedPropertyTypes {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(locationtype.PropertyTypesTable).
			SetNull(locationtype.PropertyTypesColumn).
			Where(sql.InInts(locationtype.PropertyTypesColumn, ids...)).
			Where(sql.InInts(propertytype.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ltu.property_types) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range ltu.property_types {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(propertytype.FieldID, eid)
			}
			query, args := builder.Update(locationtype.PropertyTypesTable).
				Set(locationtype.PropertyTypesColumn, id).
				Where(sql.And(p, sql.IsNull(locationtype.PropertyTypesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(ltu.property_types) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"property_types\" %v already connected to a different \"LocationType\"", keys(ltu.property_types))})
			}
		}
	}
	if len(ltu.removedSurveyTemplateCategories) > 0 {
		eids := make([]int, len(ltu.removedSurveyTemplateCategories))
		for eid := range ltu.removedSurveyTemplateCategories {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(locationtype.SurveyTemplateCategoriesTable).
			SetNull(locationtype.SurveyTemplateCategoriesColumn).
			Where(sql.InInts(locationtype.SurveyTemplateCategoriesColumn, ids...)).
			Where(sql.InInts(surveytemplatecategory.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ltu.survey_template_categories) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range ltu.survey_template_categories {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(surveytemplatecategory.FieldID, eid)
			}
			query, args := builder.Update(locationtype.SurveyTemplateCategoriesTable).
				Set(locationtype.SurveyTemplateCategoriesColumn, id).
				Where(sql.And(p, sql.IsNull(locationtype.SurveyTemplateCategoriesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(ltu.survey_template_categories) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"survey_template_categories\" %v already connected to a different \"LocationType\"", keys(ltu.survey_template_categories))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// LocationTypeUpdateOne is the builder for updating a single LocationType entity.
type LocationTypeUpdateOne struct {
	config
	id string

	update_time                     *time.Time
	site                            *bool
	name                            *string
	map_type                        *string
	clearmap_type                   bool
	map_zoom_level                  *int
	addmap_zoom_level               *int
	clearmap_zoom_level             bool
	index                           *int
	addindex                        *int
	locations                       map[string]struct{}
	property_types                  map[string]struct{}
	survey_template_categories      map[string]struct{}
	removedLocations                map[string]struct{}
	removedPropertyTypes            map[string]struct{}
	removedSurveyTemplateCategories map[string]struct{}
}

// SetSite sets the site field.
func (ltuo *LocationTypeUpdateOne) SetSite(b bool) *LocationTypeUpdateOne {
	ltuo.site = &b
	return ltuo
}

// SetNillableSite sets the site field if the given value is not nil.
func (ltuo *LocationTypeUpdateOne) SetNillableSite(b *bool) *LocationTypeUpdateOne {
	if b != nil {
		ltuo.SetSite(*b)
	}
	return ltuo
}

// SetName sets the name field.
func (ltuo *LocationTypeUpdateOne) SetName(s string) *LocationTypeUpdateOne {
	ltuo.name = &s
	return ltuo
}

// SetMapType sets the map_type field.
func (ltuo *LocationTypeUpdateOne) SetMapType(s string) *LocationTypeUpdateOne {
	ltuo.map_type = &s
	return ltuo
}

// SetNillableMapType sets the map_type field if the given value is not nil.
func (ltuo *LocationTypeUpdateOne) SetNillableMapType(s *string) *LocationTypeUpdateOne {
	if s != nil {
		ltuo.SetMapType(*s)
	}
	return ltuo
}

// ClearMapType clears the value of map_type.
func (ltuo *LocationTypeUpdateOne) ClearMapType() *LocationTypeUpdateOne {
	ltuo.map_type = nil
	ltuo.clearmap_type = true
	return ltuo
}

// SetMapZoomLevel sets the map_zoom_level field.
func (ltuo *LocationTypeUpdateOne) SetMapZoomLevel(i int) *LocationTypeUpdateOne {
	ltuo.map_zoom_level = &i
	ltuo.addmap_zoom_level = nil
	return ltuo
}

// SetNillableMapZoomLevel sets the map_zoom_level field if the given value is not nil.
func (ltuo *LocationTypeUpdateOne) SetNillableMapZoomLevel(i *int) *LocationTypeUpdateOne {
	if i != nil {
		ltuo.SetMapZoomLevel(*i)
	}
	return ltuo
}

// AddMapZoomLevel adds i to map_zoom_level.
func (ltuo *LocationTypeUpdateOne) AddMapZoomLevel(i int) *LocationTypeUpdateOne {
	if ltuo.addmap_zoom_level == nil {
		ltuo.addmap_zoom_level = &i
	} else {
		*ltuo.addmap_zoom_level += i
	}
	return ltuo
}

// ClearMapZoomLevel clears the value of map_zoom_level.
func (ltuo *LocationTypeUpdateOne) ClearMapZoomLevel() *LocationTypeUpdateOne {
	ltuo.map_zoom_level = nil
	ltuo.clearmap_zoom_level = true
	return ltuo
}

// SetIndex sets the index field.
func (ltuo *LocationTypeUpdateOne) SetIndex(i int) *LocationTypeUpdateOne {
	ltuo.index = &i
	ltuo.addindex = nil
	return ltuo
}

// SetNillableIndex sets the index field if the given value is not nil.
func (ltuo *LocationTypeUpdateOne) SetNillableIndex(i *int) *LocationTypeUpdateOne {
	if i != nil {
		ltuo.SetIndex(*i)
	}
	return ltuo
}

// AddIndex adds i to index.
func (ltuo *LocationTypeUpdateOne) AddIndex(i int) *LocationTypeUpdateOne {
	if ltuo.addindex == nil {
		ltuo.addindex = &i
	} else {
		*ltuo.addindex += i
	}
	return ltuo
}

// AddLocationIDs adds the locations edge to Location by ids.
func (ltuo *LocationTypeUpdateOne) AddLocationIDs(ids ...string) *LocationTypeUpdateOne {
	if ltuo.locations == nil {
		ltuo.locations = make(map[string]struct{})
	}
	for i := range ids {
		ltuo.locations[ids[i]] = struct{}{}
	}
	return ltuo
}

// AddLocations adds the locations edges to Location.
func (ltuo *LocationTypeUpdateOne) AddLocations(l ...*Location) *LocationTypeUpdateOne {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return ltuo.AddLocationIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (ltuo *LocationTypeUpdateOne) AddPropertyTypeIDs(ids ...string) *LocationTypeUpdateOne {
	if ltuo.property_types == nil {
		ltuo.property_types = make(map[string]struct{})
	}
	for i := range ids {
		ltuo.property_types[ids[i]] = struct{}{}
	}
	return ltuo
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (ltuo *LocationTypeUpdateOne) AddPropertyTypes(p ...*PropertyType) *LocationTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ltuo.AddPropertyTypeIDs(ids...)
}

// AddSurveyTemplateCategoryIDs adds the survey_template_categories edge to SurveyTemplateCategory by ids.
func (ltuo *LocationTypeUpdateOne) AddSurveyTemplateCategoryIDs(ids ...string) *LocationTypeUpdateOne {
	if ltuo.survey_template_categories == nil {
		ltuo.survey_template_categories = make(map[string]struct{})
	}
	for i := range ids {
		ltuo.survey_template_categories[ids[i]] = struct{}{}
	}
	return ltuo
}

// AddSurveyTemplateCategories adds the survey_template_categories edges to SurveyTemplateCategory.
func (ltuo *LocationTypeUpdateOne) AddSurveyTemplateCategories(s ...*SurveyTemplateCategory) *LocationTypeUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return ltuo.AddSurveyTemplateCategoryIDs(ids...)
}

// RemoveLocationIDs removes the locations edge to Location by ids.
func (ltuo *LocationTypeUpdateOne) RemoveLocationIDs(ids ...string) *LocationTypeUpdateOne {
	if ltuo.removedLocations == nil {
		ltuo.removedLocations = make(map[string]struct{})
	}
	for i := range ids {
		ltuo.removedLocations[ids[i]] = struct{}{}
	}
	return ltuo
}

// RemoveLocations removes locations edges to Location.
func (ltuo *LocationTypeUpdateOne) RemoveLocations(l ...*Location) *LocationTypeUpdateOne {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return ltuo.RemoveLocationIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (ltuo *LocationTypeUpdateOne) RemovePropertyTypeIDs(ids ...string) *LocationTypeUpdateOne {
	if ltuo.removedPropertyTypes == nil {
		ltuo.removedPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		ltuo.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return ltuo
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (ltuo *LocationTypeUpdateOne) RemovePropertyTypes(p ...*PropertyType) *LocationTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ltuo.RemovePropertyTypeIDs(ids...)
}

// RemoveSurveyTemplateCategoryIDs removes the survey_template_categories edge to SurveyTemplateCategory by ids.
func (ltuo *LocationTypeUpdateOne) RemoveSurveyTemplateCategoryIDs(ids ...string) *LocationTypeUpdateOne {
	if ltuo.removedSurveyTemplateCategories == nil {
		ltuo.removedSurveyTemplateCategories = make(map[string]struct{})
	}
	for i := range ids {
		ltuo.removedSurveyTemplateCategories[ids[i]] = struct{}{}
	}
	return ltuo
}

// RemoveSurveyTemplateCategories removes survey_template_categories edges to SurveyTemplateCategory.
func (ltuo *LocationTypeUpdateOne) RemoveSurveyTemplateCategories(s ...*SurveyTemplateCategory) *LocationTypeUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return ltuo.RemoveSurveyTemplateCategoryIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (ltuo *LocationTypeUpdateOne) Save(ctx context.Context) (*LocationType, error) {
	if ltuo.update_time == nil {
		v := locationtype.UpdateDefaultUpdateTime()
		ltuo.update_time = &v
	}
	return ltuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (ltuo *LocationTypeUpdateOne) SaveX(ctx context.Context) *LocationType {
	lt, err := ltuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return lt
}

// Exec executes the query on the entity.
func (ltuo *LocationTypeUpdateOne) Exec(ctx context.Context) error {
	_, err := ltuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ltuo *LocationTypeUpdateOne) ExecX(ctx context.Context) {
	if err := ltuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ltuo *LocationTypeUpdateOne) sqlSave(ctx context.Context) (lt *LocationType, err error) {
	var (
		builder  = sql.Dialect(ltuo.driver.Dialect())
		selector = builder.Select(locationtype.Columns...).From(builder.Table(locationtype.Table))
	)
	locationtype.ID(ltuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = ltuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		lt = &LocationType{config: ltuo.config}
		if err := lt.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into LocationType: %v", err)
		}
		id = lt.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("LocationType with id: %v", ltuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one LocationType with the same id: %v", ltuo.id)
	}

	tx, err := ltuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(locationtype.Table)
	)
	updater = updater.Where(sql.InInts(locationtype.FieldID, ids...))
	if value := ltuo.update_time; value != nil {
		updater.Set(locationtype.FieldUpdateTime, *value)
		lt.UpdateTime = *value
	}
	if value := ltuo.site; value != nil {
		updater.Set(locationtype.FieldSite, *value)
		lt.Site = *value
	}
	if value := ltuo.name; value != nil {
		updater.Set(locationtype.FieldName, *value)
		lt.Name = *value
	}
	if value := ltuo.map_type; value != nil {
		updater.Set(locationtype.FieldMapType, *value)
		lt.MapType = *value
	}
	if ltuo.clearmap_type {
		var value string
		lt.MapType = value
		updater.SetNull(locationtype.FieldMapType)
	}
	if value := ltuo.map_zoom_level; value != nil {
		updater.Set(locationtype.FieldMapZoomLevel, *value)
		lt.MapZoomLevel = *value
	}
	if value := ltuo.addmap_zoom_level; value != nil {
		updater.Add(locationtype.FieldMapZoomLevel, *value)
		lt.MapZoomLevel += *value
	}
	if ltuo.clearmap_zoom_level {
		var value int
		lt.MapZoomLevel = value
		updater.SetNull(locationtype.FieldMapZoomLevel)
	}
	if value := ltuo.index; value != nil {
		updater.Set(locationtype.FieldIndex, *value)
		lt.Index = *value
	}
	if value := ltuo.addindex; value != nil {
		updater.Add(locationtype.FieldIndex, *value)
		lt.Index += *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ltuo.removedLocations) > 0 {
		eids := make([]int, len(ltuo.removedLocations))
		for eid := range ltuo.removedLocations {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(locationtype.LocationsTable).
			SetNull(locationtype.LocationsColumn).
			Where(sql.InInts(locationtype.LocationsColumn, ids...)).
			Where(sql.InInts(location.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ltuo.locations) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range ltuo.locations {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(location.FieldID, eid)
			}
			query, args := builder.Update(locationtype.LocationsTable).
				Set(locationtype.LocationsColumn, id).
				Where(sql.And(p, sql.IsNull(locationtype.LocationsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(ltuo.locations) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"locations\" %v already connected to a different \"LocationType\"", keys(ltuo.locations))})
			}
		}
	}
	if len(ltuo.removedPropertyTypes) > 0 {
		eids := make([]int, len(ltuo.removedPropertyTypes))
		for eid := range ltuo.removedPropertyTypes {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(locationtype.PropertyTypesTable).
			SetNull(locationtype.PropertyTypesColumn).
			Where(sql.InInts(locationtype.PropertyTypesColumn, ids...)).
			Where(sql.InInts(propertytype.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ltuo.property_types) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range ltuo.property_types {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(propertytype.FieldID, eid)
			}
			query, args := builder.Update(locationtype.PropertyTypesTable).
				Set(locationtype.PropertyTypesColumn, id).
				Where(sql.And(p, sql.IsNull(locationtype.PropertyTypesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(ltuo.property_types) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"property_types\" %v already connected to a different \"LocationType\"", keys(ltuo.property_types))})
			}
		}
	}
	if len(ltuo.removedSurveyTemplateCategories) > 0 {
		eids := make([]int, len(ltuo.removedSurveyTemplateCategories))
		for eid := range ltuo.removedSurveyTemplateCategories {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(locationtype.SurveyTemplateCategoriesTable).
			SetNull(locationtype.SurveyTemplateCategoriesColumn).
			Where(sql.InInts(locationtype.SurveyTemplateCategoriesColumn, ids...)).
			Where(sql.InInts(surveytemplatecategory.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ltuo.survey_template_categories) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range ltuo.survey_template_categories {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(surveytemplatecategory.FieldID, eid)
			}
			query, args := builder.Update(locationtype.SurveyTemplateCategoriesTable).
				Set(locationtype.SurveyTemplateCategoriesColumn, id).
				Where(sql.And(p, sql.IsNull(locationtype.SurveyTemplateCategoriesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(ltuo.survey_template_categories) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"survey_template_categories\" %v already connected to a different \"LocationType\"", keys(ltuo.survey_template_categories))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return lt, nil
}
