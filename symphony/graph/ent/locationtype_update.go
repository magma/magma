// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
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
	locations                       map[int]struct{}
	property_types                  map[int]struct{}
	survey_template_categories      map[int]struct{}
	removedLocations                map[int]struct{}
	removedPropertyTypes            map[int]struct{}
	removedSurveyTemplateCategories map[int]struct{}
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
func (ltu *LocationTypeUpdate) AddLocationIDs(ids ...int) *LocationTypeUpdate {
	if ltu.locations == nil {
		ltu.locations = make(map[int]struct{})
	}
	for i := range ids {
		ltu.locations[ids[i]] = struct{}{}
	}
	return ltu
}

// AddLocations adds the locations edges to Location.
func (ltu *LocationTypeUpdate) AddLocations(l ...*Location) *LocationTypeUpdate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return ltu.AddLocationIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (ltu *LocationTypeUpdate) AddPropertyTypeIDs(ids ...int) *LocationTypeUpdate {
	if ltu.property_types == nil {
		ltu.property_types = make(map[int]struct{})
	}
	for i := range ids {
		ltu.property_types[ids[i]] = struct{}{}
	}
	return ltu
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (ltu *LocationTypeUpdate) AddPropertyTypes(p ...*PropertyType) *LocationTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ltu.AddPropertyTypeIDs(ids...)
}

// AddSurveyTemplateCategoryIDs adds the survey_template_categories edge to SurveyTemplateCategory by ids.
func (ltu *LocationTypeUpdate) AddSurveyTemplateCategoryIDs(ids ...int) *LocationTypeUpdate {
	if ltu.survey_template_categories == nil {
		ltu.survey_template_categories = make(map[int]struct{})
	}
	for i := range ids {
		ltu.survey_template_categories[ids[i]] = struct{}{}
	}
	return ltu
}

// AddSurveyTemplateCategories adds the survey_template_categories edges to SurveyTemplateCategory.
func (ltu *LocationTypeUpdate) AddSurveyTemplateCategories(s ...*SurveyTemplateCategory) *LocationTypeUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return ltu.AddSurveyTemplateCategoryIDs(ids...)
}

// RemoveLocationIDs removes the locations edge to Location by ids.
func (ltu *LocationTypeUpdate) RemoveLocationIDs(ids ...int) *LocationTypeUpdate {
	if ltu.removedLocations == nil {
		ltu.removedLocations = make(map[int]struct{})
	}
	for i := range ids {
		ltu.removedLocations[ids[i]] = struct{}{}
	}
	return ltu
}

// RemoveLocations removes locations edges to Location.
func (ltu *LocationTypeUpdate) RemoveLocations(l ...*Location) *LocationTypeUpdate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return ltu.RemoveLocationIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (ltu *LocationTypeUpdate) RemovePropertyTypeIDs(ids ...int) *LocationTypeUpdate {
	if ltu.removedPropertyTypes == nil {
		ltu.removedPropertyTypes = make(map[int]struct{})
	}
	for i := range ids {
		ltu.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return ltu
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (ltu *LocationTypeUpdate) RemovePropertyTypes(p ...*PropertyType) *LocationTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ltu.RemovePropertyTypeIDs(ids...)
}

// RemoveSurveyTemplateCategoryIDs removes the survey_template_categories edge to SurveyTemplateCategory by ids.
func (ltu *LocationTypeUpdate) RemoveSurveyTemplateCategoryIDs(ids ...int) *LocationTypeUpdate {
	if ltu.removedSurveyTemplateCategories == nil {
		ltu.removedSurveyTemplateCategories = make(map[int]struct{})
	}
	for i := range ids {
		ltu.removedSurveyTemplateCategories[ids[i]] = struct{}{}
	}
	return ltu
}

// RemoveSurveyTemplateCategories removes survey_template_categories edges to SurveyTemplateCategory.
func (ltu *LocationTypeUpdate) RemoveSurveyTemplateCategories(s ...*SurveyTemplateCategory) *LocationTypeUpdate {
	ids := make([]int, len(s))
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
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   locationtype.Table,
			Columns: locationtype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: locationtype.FieldID,
			},
		},
	}
	if ps := ltu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := ltu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: locationtype.FieldUpdateTime,
		})
	}
	if value := ltu.site; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: locationtype.FieldSite,
		})
	}
	if value := ltu.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: locationtype.FieldName,
		})
	}
	if value := ltu.map_type; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: locationtype.FieldMapType,
		})
	}
	if ltu.clearmap_type {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: locationtype.FieldMapType,
		})
	}
	if value := ltu.map_zoom_level; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: locationtype.FieldMapZoomLevel,
		})
	}
	if value := ltu.addmap_zoom_level; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: locationtype.FieldMapZoomLevel,
		})
	}
	if ltu.clearmap_zoom_level {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: locationtype.FieldMapZoomLevel,
		})
	}
	if value := ltu.index; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: locationtype.FieldIndex,
		})
	}
	if value := ltu.addindex; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: locationtype.FieldIndex,
		})
	}
	if nodes := ltu.removedLocations; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   locationtype.LocationsTable,
			Columns: []string{locationtype.LocationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ltu.locations; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   locationtype.LocationsTable,
			Columns: []string{locationtype.LocationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := ltu.removedPropertyTypes; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   locationtype.PropertyTypesTable,
			Columns: []string{locationtype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: propertytype.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ltu.property_types; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   locationtype.PropertyTypesTable,
			Columns: []string{locationtype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: propertytype.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := ltu.removedSurveyTemplateCategories; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   locationtype.SurveyTemplateCategoriesTable,
			Columns: []string{locationtype.SurveyTemplateCategoriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveytemplatecategory.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ltu.survey_template_categories; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   locationtype.SurveyTemplateCategoriesTable,
			Columns: []string{locationtype.SurveyTemplateCategoriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveytemplatecategory.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ltu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{locationtype.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// LocationTypeUpdateOne is the builder for updating a single LocationType entity.
type LocationTypeUpdateOne struct {
	config
	id int

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
	locations                       map[int]struct{}
	property_types                  map[int]struct{}
	survey_template_categories      map[int]struct{}
	removedLocations                map[int]struct{}
	removedPropertyTypes            map[int]struct{}
	removedSurveyTemplateCategories map[int]struct{}
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
func (ltuo *LocationTypeUpdateOne) AddLocationIDs(ids ...int) *LocationTypeUpdateOne {
	if ltuo.locations == nil {
		ltuo.locations = make(map[int]struct{})
	}
	for i := range ids {
		ltuo.locations[ids[i]] = struct{}{}
	}
	return ltuo
}

// AddLocations adds the locations edges to Location.
func (ltuo *LocationTypeUpdateOne) AddLocations(l ...*Location) *LocationTypeUpdateOne {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return ltuo.AddLocationIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (ltuo *LocationTypeUpdateOne) AddPropertyTypeIDs(ids ...int) *LocationTypeUpdateOne {
	if ltuo.property_types == nil {
		ltuo.property_types = make(map[int]struct{})
	}
	for i := range ids {
		ltuo.property_types[ids[i]] = struct{}{}
	}
	return ltuo
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (ltuo *LocationTypeUpdateOne) AddPropertyTypes(p ...*PropertyType) *LocationTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ltuo.AddPropertyTypeIDs(ids...)
}

// AddSurveyTemplateCategoryIDs adds the survey_template_categories edge to SurveyTemplateCategory by ids.
func (ltuo *LocationTypeUpdateOne) AddSurveyTemplateCategoryIDs(ids ...int) *LocationTypeUpdateOne {
	if ltuo.survey_template_categories == nil {
		ltuo.survey_template_categories = make(map[int]struct{})
	}
	for i := range ids {
		ltuo.survey_template_categories[ids[i]] = struct{}{}
	}
	return ltuo
}

// AddSurveyTemplateCategories adds the survey_template_categories edges to SurveyTemplateCategory.
func (ltuo *LocationTypeUpdateOne) AddSurveyTemplateCategories(s ...*SurveyTemplateCategory) *LocationTypeUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return ltuo.AddSurveyTemplateCategoryIDs(ids...)
}

// RemoveLocationIDs removes the locations edge to Location by ids.
func (ltuo *LocationTypeUpdateOne) RemoveLocationIDs(ids ...int) *LocationTypeUpdateOne {
	if ltuo.removedLocations == nil {
		ltuo.removedLocations = make(map[int]struct{})
	}
	for i := range ids {
		ltuo.removedLocations[ids[i]] = struct{}{}
	}
	return ltuo
}

// RemoveLocations removes locations edges to Location.
func (ltuo *LocationTypeUpdateOne) RemoveLocations(l ...*Location) *LocationTypeUpdateOne {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return ltuo.RemoveLocationIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (ltuo *LocationTypeUpdateOne) RemovePropertyTypeIDs(ids ...int) *LocationTypeUpdateOne {
	if ltuo.removedPropertyTypes == nil {
		ltuo.removedPropertyTypes = make(map[int]struct{})
	}
	for i := range ids {
		ltuo.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return ltuo
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (ltuo *LocationTypeUpdateOne) RemovePropertyTypes(p ...*PropertyType) *LocationTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ltuo.RemovePropertyTypeIDs(ids...)
}

// RemoveSurveyTemplateCategoryIDs removes the survey_template_categories edge to SurveyTemplateCategory by ids.
func (ltuo *LocationTypeUpdateOne) RemoveSurveyTemplateCategoryIDs(ids ...int) *LocationTypeUpdateOne {
	if ltuo.removedSurveyTemplateCategories == nil {
		ltuo.removedSurveyTemplateCategories = make(map[int]struct{})
	}
	for i := range ids {
		ltuo.removedSurveyTemplateCategories[ids[i]] = struct{}{}
	}
	return ltuo
}

// RemoveSurveyTemplateCategories removes survey_template_categories edges to SurveyTemplateCategory.
func (ltuo *LocationTypeUpdateOne) RemoveSurveyTemplateCategories(s ...*SurveyTemplateCategory) *LocationTypeUpdateOne {
	ids := make([]int, len(s))
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
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   locationtype.Table,
			Columns: locationtype.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  ltuo.id,
				Type:   field.TypeInt,
				Column: locationtype.FieldID,
			},
		},
	}
	if value := ltuo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: locationtype.FieldUpdateTime,
		})
	}
	if value := ltuo.site; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: locationtype.FieldSite,
		})
	}
	if value := ltuo.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: locationtype.FieldName,
		})
	}
	if value := ltuo.map_type; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: locationtype.FieldMapType,
		})
	}
	if ltuo.clearmap_type {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: locationtype.FieldMapType,
		})
	}
	if value := ltuo.map_zoom_level; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: locationtype.FieldMapZoomLevel,
		})
	}
	if value := ltuo.addmap_zoom_level; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: locationtype.FieldMapZoomLevel,
		})
	}
	if ltuo.clearmap_zoom_level {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: locationtype.FieldMapZoomLevel,
		})
	}
	if value := ltuo.index; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: locationtype.FieldIndex,
		})
	}
	if value := ltuo.addindex; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: locationtype.FieldIndex,
		})
	}
	if nodes := ltuo.removedLocations; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   locationtype.LocationsTable,
			Columns: []string{locationtype.LocationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ltuo.locations; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   locationtype.LocationsTable,
			Columns: []string{locationtype.LocationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := ltuo.removedPropertyTypes; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   locationtype.PropertyTypesTable,
			Columns: []string{locationtype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: propertytype.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ltuo.property_types; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   locationtype.PropertyTypesTable,
			Columns: []string{locationtype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: propertytype.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := ltuo.removedSurveyTemplateCategories; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   locationtype.SurveyTemplateCategoriesTable,
			Columns: []string{locationtype.SurveyTemplateCategoriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveytemplatecategory.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ltuo.survey_template_categories; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   locationtype.SurveyTemplateCategoriesTable,
			Columns: []string{locationtype.SurveyTemplateCategoriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveytemplatecategory.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	lt = &LocationType{config: ltuo.config}
	_spec.Assign = lt.assignValues
	_spec.ScanValues = lt.scanValues()
	if err = sqlgraph.UpdateNode(ctx, ltuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{locationtype.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return lt, nil
}
