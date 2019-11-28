// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatecategory"
)

// LocationTypeCreate is the builder for creating a LocationType entity.
type LocationTypeCreate struct {
	config
	create_time                *time.Time
	update_time                *time.Time
	site                       *bool
	name                       *string
	map_type                   *string
	map_zoom_level             *int
	index                      *int
	locations                  map[string]struct{}
	property_types             map[string]struct{}
	survey_template_categories map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (ltc *LocationTypeCreate) SetCreateTime(t time.Time) *LocationTypeCreate {
	ltc.create_time = &t
	return ltc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (ltc *LocationTypeCreate) SetNillableCreateTime(t *time.Time) *LocationTypeCreate {
	if t != nil {
		ltc.SetCreateTime(*t)
	}
	return ltc
}

// SetUpdateTime sets the update_time field.
func (ltc *LocationTypeCreate) SetUpdateTime(t time.Time) *LocationTypeCreate {
	ltc.update_time = &t
	return ltc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (ltc *LocationTypeCreate) SetNillableUpdateTime(t *time.Time) *LocationTypeCreate {
	if t != nil {
		ltc.SetUpdateTime(*t)
	}
	return ltc
}

// SetSite sets the site field.
func (ltc *LocationTypeCreate) SetSite(b bool) *LocationTypeCreate {
	ltc.site = &b
	return ltc
}

// SetNillableSite sets the site field if the given value is not nil.
func (ltc *LocationTypeCreate) SetNillableSite(b *bool) *LocationTypeCreate {
	if b != nil {
		ltc.SetSite(*b)
	}
	return ltc
}

// SetName sets the name field.
func (ltc *LocationTypeCreate) SetName(s string) *LocationTypeCreate {
	ltc.name = &s
	return ltc
}

// SetMapType sets the map_type field.
func (ltc *LocationTypeCreate) SetMapType(s string) *LocationTypeCreate {
	ltc.map_type = &s
	return ltc
}

// SetNillableMapType sets the map_type field if the given value is not nil.
func (ltc *LocationTypeCreate) SetNillableMapType(s *string) *LocationTypeCreate {
	if s != nil {
		ltc.SetMapType(*s)
	}
	return ltc
}

// SetMapZoomLevel sets the map_zoom_level field.
func (ltc *LocationTypeCreate) SetMapZoomLevel(i int) *LocationTypeCreate {
	ltc.map_zoom_level = &i
	return ltc
}

// SetNillableMapZoomLevel sets the map_zoom_level field if the given value is not nil.
func (ltc *LocationTypeCreate) SetNillableMapZoomLevel(i *int) *LocationTypeCreate {
	if i != nil {
		ltc.SetMapZoomLevel(*i)
	}
	return ltc
}

// SetIndex sets the index field.
func (ltc *LocationTypeCreate) SetIndex(i int) *LocationTypeCreate {
	ltc.index = &i
	return ltc
}

// SetNillableIndex sets the index field if the given value is not nil.
func (ltc *LocationTypeCreate) SetNillableIndex(i *int) *LocationTypeCreate {
	if i != nil {
		ltc.SetIndex(*i)
	}
	return ltc
}

// AddLocationIDs adds the locations edge to Location by ids.
func (ltc *LocationTypeCreate) AddLocationIDs(ids ...string) *LocationTypeCreate {
	if ltc.locations == nil {
		ltc.locations = make(map[string]struct{})
	}
	for i := range ids {
		ltc.locations[ids[i]] = struct{}{}
	}
	return ltc
}

// AddLocations adds the locations edges to Location.
func (ltc *LocationTypeCreate) AddLocations(l ...*Location) *LocationTypeCreate {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return ltc.AddLocationIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (ltc *LocationTypeCreate) AddPropertyTypeIDs(ids ...string) *LocationTypeCreate {
	if ltc.property_types == nil {
		ltc.property_types = make(map[string]struct{})
	}
	for i := range ids {
		ltc.property_types[ids[i]] = struct{}{}
	}
	return ltc
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (ltc *LocationTypeCreate) AddPropertyTypes(p ...*PropertyType) *LocationTypeCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ltc.AddPropertyTypeIDs(ids...)
}

// AddSurveyTemplateCategoryIDs adds the survey_template_categories edge to SurveyTemplateCategory by ids.
func (ltc *LocationTypeCreate) AddSurveyTemplateCategoryIDs(ids ...string) *LocationTypeCreate {
	if ltc.survey_template_categories == nil {
		ltc.survey_template_categories = make(map[string]struct{})
	}
	for i := range ids {
		ltc.survey_template_categories[ids[i]] = struct{}{}
	}
	return ltc
}

// AddSurveyTemplateCategories adds the survey_template_categories edges to SurveyTemplateCategory.
func (ltc *LocationTypeCreate) AddSurveyTemplateCategories(s ...*SurveyTemplateCategory) *LocationTypeCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return ltc.AddSurveyTemplateCategoryIDs(ids...)
}

// Save creates the LocationType in the database.
func (ltc *LocationTypeCreate) Save(ctx context.Context) (*LocationType, error) {
	if ltc.create_time == nil {
		v := locationtype.DefaultCreateTime()
		ltc.create_time = &v
	}
	if ltc.update_time == nil {
		v := locationtype.DefaultUpdateTime()
		ltc.update_time = &v
	}
	if ltc.site == nil {
		v := locationtype.DefaultSite
		ltc.site = &v
	}
	if ltc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if ltc.map_zoom_level == nil {
		v := locationtype.DefaultMapZoomLevel
		ltc.map_zoom_level = &v
	}
	if ltc.index == nil {
		v := locationtype.DefaultIndex
		ltc.index = &v
	}
	return ltc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (ltc *LocationTypeCreate) SaveX(ctx context.Context) *LocationType {
	v, err := ltc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ltc *LocationTypeCreate) sqlSave(ctx context.Context) (*LocationType, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(ltc.driver.Dialect())
		lt      = &LocationType{config: ltc.config}
	)
	tx, err := ltc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(locationtype.Table).Default()
	if value := ltc.create_time; value != nil {
		insert.Set(locationtype.FieldCreateTime, *value)
		lt.CreateTime = *value
	}
	if value := ltc.update_time; value != nil {
		insert.Set(locationtype.FieldUpdateTime, *value)
		lt.UpdateTime = *value
	}
	if value := ltc.site; value != nil {
		insert.Set(locationtype.FieldSite, *value)
		lt.Site = *value
	}
	if value := ltc.name; value != nil {
		insert.Set(locationtype.FieldName, *value)
		lt.Name = *value
	}
	if value := ltc.map_type; value != nil {
		insert.Set(locationtype.FieldMapType, *value)
		lt.MapType = *value
	}
	if value := ltc.map_zoom_level; value != nil {
		insert.Set(locationtype.FieldMapZoomLevel, *value)
		lt.MapZoomLevel = *value
	}
	if value := ltc.index; value != nil {
		insert.Set(locationtype.FieldIndex, *value)
		lt.Index = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(locationtype.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	lt.ID = strconv.FormatInt(id, 10)
	if len(ltc.locations) > 0 {
		p := sql.P()
		for eid := range ltc.locations {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
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
		if int(affected) < len(ltc.locations) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"locations\" %v already connected to a different \"LocationType\"", keys(ltc.locations))})
		}
	}
	if len(ltc.property_types) > 0 {
		p := sql.P()
		for eid := range ltc.property_types {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
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
		if int(affected) < len(ltc.property_types) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"property_types\" %v already connected to a different \"LocationType\"", keys(ltc.property_types))})
		}
	}
	if len(ltc.survey_template_categories) > 0 {
		p := sql.P()
		for eid := range ltc.survey_template_categories {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
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
		if int(affected) < len(ltc.survey_template_categories) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"survey_template_categories\" %v already connected to a different \"LocationType\"", keys(ltc.survey_template_categories))})
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return lt, nil
}
