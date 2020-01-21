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
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/floorplan"
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// LocationUpdate is the builder for updating Location entities.
type LocationUpdate struct {
	config

	update_time             *time.Time
	name                    *string
	external_id             *string
	clearexternal_id        bool
	latitude                *float64
	addlatitude             *float64
	longitude               *float64
	addlongitude            *float64
	site_survey_needed      *bool
	clearsite_survey_needed bool
	_type                   map[string]struct{}
	parent                  map[string]struct{}
	children                map[string]struct{}
	files                   map[string]struct{}
	hyperlinks              map[string]struct{}
	equipment               map[string]struct{}
	properties              map[string]struct{}
	survey                  map[string]struct{}
	wifi_scan               map[string]struct{}
	cell_scan               map[string]struct{}
	work_orders             map[string]struct{}
	floor_plans             map[string]struct{}
	clearedType             bool
	clearedParent           bool
	removedChildren         map[string]struct{}
	removedFiles            map[string]struct{}
	removedHyperlinks       map[string]struct{}
	removedEquipment        map[string]struct{}
	removedProperties       map[string]struct{}
	removedSurvey           map[string]struct{}
	removedWifiScan         map[string]struct{}
	removedCellScan         map[string]struct{}
	removedWorkOrders       map[string]struct{}
	removedFloorPlans       map[string]struct{}
	predicates              []predicate.Location
}

// Where adds a new predicate for the builder.
func (lu *LocationUpdate) Where(ps ...predicate.Location) *LocationUpdate {
	lu.predicates = append(lu.predicates, ps...)
	return lu
}

// SetName sets the name field.
func (lu *LocationUpdate) SetName(s string) *LocationUpdate {
	lu.name = &s
	return lu
}

// SetExternalID sets the external_id field.
func (lu *LocationUpdate) SetExternalID(s string) *LocationUpdate {
	lu.external_id = &s
	return lu
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (lu *LocationUpdate) SetNillableExternalID(s *string) *LocationUpdate {
	if s != nil {
		lu.SetExternalID(*s)
	}
	return lu
}

// ClearExternalID clears the value of external_id.
func (lu *LocationUpdate) ClearExternalID() *LocationUpdate {
	lu.external_id = nil
	lu.clearexternal_id = true
	return lu
}

// SetLatitude sets the latitude field.
func (lu *LocationUpdate) SetLatitude(f float64) *LocationUpdate {
	lu.latitude = &f
	lu.addlatitude = nil
	return lu
}

// SetNillableLatitude sets the latitude field if the given value is not nil.
func (lu *LocationUpdate) SetNillableLatitude(f *float64) *LocationUpdate {
	if f != nil {
		lu.SetLatitude(*f)
	}
	return lu
}

// AddLatitude adds f to latitude.
func (lu *LocationUpdate) AddLatitude(f float64) *LocationUpdate {
	if lu.addlatitude == nil {
		lu.addlatitude = &f
	} else {
		*lu.addlatitude += f
	}
	return lu
}

// SetLongitude sets the longitude field.
func (lu *LocationUpdate) SetLongitude(f float64) *LocationUpdate {
	lu.longitude = &f
	lu.addlongitude = nil
	return lu
}

// SetNillableLongitude sets the longitude field if the given value is not nil.
func (lu *LocationUpdate) SetNillableLongitude(f *float64) *LocationUpdate {
	if f != nil {
		lu.SetLongitude(*f)
	}
	return lu
}

// AddLongitude adds f to longitude.
func (lu *LocationUpdate) AddLongitude(f float64) *LocationUpdate {
	if lu.addlongitude == nil {
		lu.addlongitude = &f
	} else {
		*lu.addlongitude += f
	}
	return lu
}

// SetSiteSurveyNeeded sets the site_survey_needed field.
func (lu *LocationUpdate) SetSiteSurveyNeeded(b bool) *LocationUpdate {
	lu.site_survey_needed = &b
	return lu
}

// SetNillableSiteSurveyNeeded sets the site_survey_needed field if the given value is not nil.
func (lu *LocationUpdate) SetNillableSiteSurveyNeeded(b *bool) *LocationUpdate {
	if b != nil {
		lu.SetSiteSurveyNeeded(*b)
	}
	return lu
}

// ClearSiteSurveyNeeded clears the value of site_survey_needed.
func (lu *LocationUpdate) ClearSiteSurveyNeeded() *LocationUpdate {
	lu.site_survey_needed = nil
	lu.clearsite_survey_needed = true
	return lu
}

// SetTypeID sets the type edge to LocationType by id.
func (lu *LocationUpdate) SetTypeID(id string) *LocationUpdate {
	if lu._type == nil {
		lu._type = make(map[string]struct{})
	}
	lu._type[id] = struct{}{}
	return lu
}

// SetType sets the type edge to LocationType.
func (lu *LocationUpdate) SetType(l *LocationType) *LocationUpdate {
	return lu.SetTypeID(l.ID)
}

// SetParentID sets the parent edge to Location by id.
func (lu *LocationUpdate) SetParentID(id string) *LocationUpdate {
	if lu.parent == nil {
		lu.parent = make(map[string]struct{})
	}
	lu.parent[id] = struct{}{}
	return lu
}

// SetNillableParentID sets the parent edge to Location by id if the given value is not nil.
func (lu *LocationUpdate) SetNillableParentID(id *string) *LocationUpdate {
	if id != nil {
		lu = lu.SetParentID(*id)
	}
	return lu
}

// SetParent sets the parent edge to Location.
func (lu *LocationUpdate) SetParent(l *Location) *LocationUpdate {
	return lu.SetParentID(l.ID)
}

// AddChildIDs adds the children edge to Location by ids.
func (lu *LocationUpdate) AddChildIDs(ids ...string) *LocationUpdate {
	if lu.children == nil {
		lu.children = make(map[string]struct{})
	}
	for i := range ids {
		lu.children[ids[i]] = struct{}{}
	}
	return lu
}

// AddChildren adds the children edges to Location.
func (lu *LocationUpdate) AddChildren(l ...*Location) *LocationUpdate {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return lu.AddChildIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (lu *LocationUpdate) AddFileIDs(ids ...string) *LocationUpdate {
	if lu.files == nil {
		lu.files = make(map[string]struct{})
	}
	for i := range ids {
		lu.files[ids[i]] = struct{}{}
	}
	return lu
}

// AddFiles adds the files edges to File.
func (lu *LocationUpdate) AddFiles(f ...*File) *LocationUpdate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return lu.AddFileIDs(ids...)
}

// AddHyperlinkIDs adds the hyperlinks edge to Hyperlink by ids.
func (lu *LocationUpdate) AddHyperlinkIDs(ids ...string) *LocationUpdate {
	if lu.hyperlinks == nil {
		lu.hyperlinks = make(map[string]struct{})
	}
	for i := range ids {
		lu.hyperlinks[ids[i]] = struct{}{}
	}
	return lu
}

// AddHyperlinks adds the hyperlinks edges to Hyperlink.
func (lu *LocationUpdate) AddHyperlinks(h ...*Hyperlink) *LocationUpdate {
	ids := make([]string, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return lu.AddHyperlinkIDs(ids...)
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (lu *LocationUpdate) AddEquipmentIDs(ids ...string) *LocationUpdate {
	if lu.equipment == nil {
		lu.equipment = make(map[string]struct{})
	}
	for i := range ids {
		lu.equipment[ids[i]] = struct{}{}
	}
	return lu
}

// AddEquipment adds the equipment edges to Equipment.
func (lu *LocationUpdate) AddEquipment(e ...*Equipment) *LocationUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return lu.AddEquipmentIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (lu *LocationUpdate) AddPropertyIDs(ids ...string) *LocationUpdate {
	if lu.properties == nil {
		lu.properties = make(map[string]struct{})
	}
	for i := range ids {
		lu.properties[ids[i]] = struct{}{}
	}
	return lu
}

// AddProperties adds the properties edges to Property.
func (lu *LocationUpdate) AddProperties(p ...*Property) *LocationUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return lu.AddPropertyIDs(ids...)
}

// AddSurveyIDs adds the survey edge to Survey by ids.
func (lu *LocationUpdate) AddSurveyIDs(ids ...string) *LocationUpdate {
	if lu.survey == nil {
		lu.survey = make(map[string]struct{})
	}
	for i := range ids {
		lu.survey[ids[i]] = struct{}{}
	}
	return lu
}

// AddSurvey adds the survey edges to Survey.
func (lu *LocationUpdate) AddSurvey(s ...*Survey) *LocationUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lu.AddSurveyIDs(ids...)
}

// AddWifiScanIDs adds the wifi_scan edge to SurveyWiFiScan by ids.
func (lu *LocationUpdate) AddWifiScanIDs(ids ...string) *LocationUpdate {
	if lu.wifi_scan == nil {
		lu.wifi_scan = make(map[string]struct{})
	}
	for i := range ids {
		lu.wifi_scan[ids[i]] = struct{}{}
	}
	return lu
}

// AddWifiScan adds the wifi_scan edges to SurveyWiFiScan.
func (lu *LocationUpdate) AddWifiScan(s ...*SurveyWiFiScan) *LocationUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lu.AddWifiScanIDs(ids...)
}

// AddCellScanIDs adds the cell_scan edge to SurveyCellScan by ids.
func (lu *LocationUpdate) AddCellScanIDs(ids ...string) *LocationUpdate {
	if lu.cell_scan == nil {
		lu.cell_scan = make(map[string]struct{})
	}
	for i := range ids {
		lu.cell_scan[ids[i]] = struct{}{}
	}
	return lu
}

// AddCellScan adds the cell_scan edges to SurveyCellScan.
func (lu *LocationUpdate) AddCellScan(s ...*SurveyCellScan) *LocationUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lu.AddCellScanIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (lu *LocationUpdate) AddWorkOrderIDs(ids ...string) *LocationUpdate {
	if lu.work_orders == nil {
		lu.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		lu.work_orders[ids[i]] = struct{}{}
	}
	return lu
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (lu *LocationUpdate) AddWorkOrders(w ...*WorkOrder) *LocationUpdate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return lu.AddWorkOrderIDs(ids...)
}

// AddFloorPlanIDs adds the floor_plans edge to FloorPlan by ids.
func (lu *LocationUpdate) AddFloorPlanIDs(ids ...string) *LocationUpdate {
	if lu.floor_plans == nil {
		lu.floor_plans = make(map[string]struct{})
	}
	for i := range ids {
		lu.floor_plans[ids[i]] = struct{}{}
	}
	return lu
}

// AddFloorPlans adds the floor_plans edges to FloorPlan.
func (lu *LocationUpdate) AddFloorPlans(f ...*FloorPlan) *LocationUpdate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return lu.AddFloorPlanIDs(ids...)
}

// ClearType clears the type edge to LocationType.
func (lu *LocationUpdate) ClearType() *LocationUpdate {
	lu.clearedType = true
	return lu
}

// ClearParent clears the parent edge to Location.
func (lu *LocationUpdate) ClearParent() *LocationUpdate {
	lu.clearedParent = true
	return lu
}

// RemoveChildIDs removes the children edge to Location by ids.
func (lu *LocationUpdate) RemoveChildIDs(ids ...string) *LocationUpdate {
	if lu.removedChildren == nil {
		lu.removedChildren = make(map[string]struct{})
	}
	for i := range ids {
		lu.removedChildren[ids[i]] = struct{}{}
	}
	return lu
}

// RemoveChildren removes children edges to Location.
func (lu *LocationUpdate) RemoveChildren(l ...*Location) *LocationUpdate {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return lu.RemoveChildIDs(ids...)
}

// RemoveFileIDs removes the files edge to File by ids.
func (lu *LocationUpdate) RemoveFileIDs(ids ...string) *LocationUpdate {
	if lu.removedFiles == nil {
		lu.removedFiles = make(map[string]struct{})
	}
	for i := range ids {
		lu.removedFiles[ids[i]] = struct{}{}
	}
	return lu
}

// RemoveFiles removes files edges to File.
func (lu *LocationUpdate) RemoveFiles(f ...*File) *LocationUpdate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return lu.RemoveFileIDs(ids...)
}

// RemoveHyperlinkIDs removes the hyperlinks edge to Hyperlink by ids.
func (lu *LocationUpdate) RemoveHyperlinkIDs(ids ...string) *LocationUpdate {
	if lu.removedHyperlinks == nil {
		lu.removedHyperlinks = make(map[string]struct{})
	}
	for i := range ids {
		lu.removedHyperlinks[ids[i]] = struct{}{}
	}
	return lu
}

// RemoveHyperlinks removes hyperlinks edges to Hyperlink.
func (lu *LocationUpdate) RemoveHyperlinks(h ...*Hyperlink) *LocationUpdate {
	ids := make([]string, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return lu.RemoveHyperlinkIDs(ids...)
}

// RemoveEquipmentIDs removes the equipment edge to Equipment by ids.
func (lu *LocationUpdate) RemoveEquipmentIDs(ids ...string) *LocationUpdate {
	if lu.removedEquipment == nil {
		lu.removedEquipment = make(map[string]struct{})
	}
	for i := range ids {
		lu.removedEquipment[ids[i]] = struct{}{}
	}
	return lu
}

// RemoveEquipment removes equipment edges to Equipment.
func (lu *LocationUpdate) RemoveEquipment(e ...*Equipment) *LocationUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return lu.RemoveEquipmentIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (lu *LocationUpdate) RemovePropertyIDs(ids ...string) *LocationUpdate {
	if lu.removedProperties == nil {
		lu.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		lu.removedProperties[ids[i]] = struct{}{}
	}
	return lu
}

// RemoveProperties removes properties edges to Property.
func (lu *LocationUpdate) RemoveProperties(p ...*Property) *LocationUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return lu.RemovePropertyIDs(ids...)
}

// RemoveSurveyIDs removes the survey edge to Survey by ids.
func (lu *LocationUpdate) RemoveSurveyIDs(ids ...string) *LocationUpdate {
	if lu.removedSurvey == nil {
		lu.removedSurvey = make(map[string]struct{})
	}
	for i := range ids {
		lu.removedSurvey[ids[i]] = struct{}{}
	}
	return lu
}

// RemoveSurvey removes survey edges to Survey.
func (lu *LocationUpdate) RemoveSurvey(s ...*Survey) *LocationUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lu.RemoveSurveyIDs(ids...)
}

// RemoveWifiScanIDs removes the wifi_scan edge to SurveyWiFiScan by ids.
func (lu *LocationUpdate) RemoveWifiScanIDs(ids ...string) *LocationUpdate {
	if lu.removedWifiScan == nil {
		lu.removedWifiScan = make(map[string]struct{})
	}
	for i := range ids {
		lu.removedWifiScan[ids[i]] = struct{}{}
	}
	return lu
}

// RemoveWifiScan removes wifi_scan edges to SurveyWiFiScan.
func (lu *LocationUpdate) RemoveWifiScan(s ...*SurveyWiFiScan) *LocationUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lu.RemoveWifiScanIDs(ids...)
}

// RemoveCellScanIDs removes the cell_scan edge to SurveyCellScan by ids.
func (lu *LocationUpdate) RemoveCellScanIDs(ids ...string) *LocationUpdate {
	if lu.removedCellScan == nil {
		lu.removedCellScan = make(map[string]struct{})
	}
	for i := range ids {
		lu.removedCellScan[ids[i]] = struct{}{}
	}
	return lu
}

// RemoveCellScan removes cell_scan edges to SurveyCellScan.
func (lu *LocationUpdate) RemoveCellScan(s ...*SurveyCellScan) *LocationUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lu.RemoveCellScanIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (lu *LocationUpdate) RemoveWorkOrderIDs(ids ...string) *LocationUpdate {
	if lu.removedWorkOrders == nil {
		lu.removedWorkOrders = make(map[string]struct{})
	}
	for i := range ids {
		lu.removedWorkOrders[ids[i]] = struct{}{}
	}
	return lu
}

// RemoveWorkOrders removes work_orders edges to WorkOrder.
func (lu *LocationUpdate) RemoveWorkOrders(w ...*WorkOrder) *LocationUpdate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return lu.RemoveWorkOrderIDs(ids...)
}

// RemoveFloorPlanIDs removes the floor_plans edge to FloorPlan by ids.
func (lu *LocationUpdate) RemoveFloorPlanIDs(ids ...string) *LocationUpdate {
	if lu.removedFloorPlans == nil {
		lu.removedFloorPlans = make(map[string]struct{})
	}
	for i := range ids {
		lu.removedFloorPlans[ids[i]] = struct{}{}
	}
	return lu
}

// RemoveFloorPlans removes floor_plans edges to FloorPlan.
func (lu *LocationUpdate) RemoveFloorPlans(f ...*FloorPlan) *LocationUpdate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return lu.RemoveFloorPlanIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (lu *LocationUpdate) Save(ctx context.Context) (int, error) {
	if lu.update_time == nil {
		v := location.UpdateDefaultUpdateTime()
		lu.update_time = &v
	}
	if lu.name != nil {
		if err := location.NameValidator(*lu.name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if lu.latitude != nil {
		if err := location.LatitudeValidator(*lu.latitude); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"latitude\": %v", err)
		}
	}
	if lu.longitude != nil {
		if err := location.LongitudeValidator(*lu.longitude); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"longitude\": %v", err)
		}
	}
	if len(lu._type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if lu.clearedType && lu._type == nil {
		return 0, errors.New("ent: clearing a unique edge \"type\"")
	}
	if len(lu.parent) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"parent\"")
	}
	return lu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (lu *LocationUpdate) SaveX(ctx context.Context) int {
	affected, err := lu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (lu *LocationUpdate) Exec(ctx context.Context) error {
	_, err := lu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lu *LocationUpdate) ExecX(ctx context.Context) {
	if err := lu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (lu *LocationUpdate) sqlSave(ctx context.Context) (n int, err error) {
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   location.Table,
			Columns: location.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: location.FieldID,
			},
		},
	}
	if ps := lu.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := lu.update_time; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: location.FieldUpdateTime,
		})
	}
	if value := lu.name; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: location.FieldName,
		})
	}
	if value := lu.external_id; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: location.FieldExternalID,
		})
	}
	if lu.clearexternal_id {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: location.FieldExternalID,
		})
	}
	if value := lu.latitude; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: location.FieldLatitude,
		})
	}
	if value := lu.addlatitude; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: location.FieldLatitude,
		})
	}
	if value := lu.longitude; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: location.FieldLongitude,
		})
	}
	if value := lu.addlongitude; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: location.FieldLongitude,
		})
	}
	if value := lu.site_survey_needed; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: location.FieldSiteSurveyNeeded,
		})
	}
	if lu.clearsite_survey_needed {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: location.FieldSiteSurveyNeeded,
		})
	}
	if lu.clearedType {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   location.TypeTable,
			Columns: []string{location.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: locationtype.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := lu._type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   location.TypeTable,
			Columns: []string{location.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: locationtype.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if lu.clearedParent {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   location.ParentTable,
			Columns: []string{location.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := lu.parent; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   location.ParentTable,
			Columns: []string{location.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := lu.removedChildren; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.ChildrenTable,
			Columns: []string{location.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := lu.children; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.ChildrenTable,
			Columns: []string{location.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := lu.removedFiles; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.FilesTable,
			Columns: []string{location.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: file.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := lu.files; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.FilesTable,
			Columns: []string{location.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: file.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := lu.removedHyperlinks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.HyperlinksTable,
			Columns: []string{location.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: hyperlink.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := lu.hyperlinks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.HyperlinksTable,
			Columns: []string{location.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: hyperlink.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := lu.removedEquipment; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.EquipmentTable,
			Columns: []string{location.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipment.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := lu.equipment; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.EquipmentTable,
			Columns: []string{location.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipment.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := lu.removedProperties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.PropertiesTable,
			Columns: []string{location.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := lu.properties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.PropertiesTable,
			Columns: []string{location.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := lu.removedSurvey; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.SurveyTable,
			Columns: []string{location.SurveyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: survey.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := lu.survey; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.SurveyTable,
			Columns: []string{location.SurveyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: survey.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := lu.removedWifiScan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WifiScanTable,
			Columns: []string{location.WifiScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveywifiscan.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := lu.wifi_scan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WifiScanTable,
			Columns: []string{location.WifiScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveywifiscan.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := lu.removedCellScan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.CellScanTable,
			Columns: []string{location.CellScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveycellscan.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := lu.cell_scan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.CellScanTable,
			Columns: []string{location.CellScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveycellscan.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := lu.removedWorkOrders; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WorkOrdersTable,
			Columns: []string{location.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := lu.work_orders; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WorkOrdersTable,
			Columns: []string{location.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := lu.removedFloorPlans; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.FloorPlansTable,
			Columns: []string{location.FloorPlansColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: floorplan.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := lu.floor_plans; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.FloorPlansTable,
			Columns: []string{location.FloorPlansColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: floorplan.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, lu.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// LocationUpdateOne is the builder for updating a single Location entity.
type LocationUpdateOne struct {
	config
	id string

	update_time             *time.Time
	name                    *string
	external_id             *string
	clearexternal_id        bool
	latitude                *float64
	addlatitude             *float64
	longitude               *float64
	addlongitude            *float64
	site_survey_needed      *bool
	clearsite_survey_needed bool
	_type                   map[string]struct{}
	parent                  map[string]struct{}
	children                map[string]struct{}
	files                   map[string]struct{}
	hyperlinks              map[string]struct{}
	equipment               map[string]struct{}
	properties              map[string]struct{}
	survey                  map[string]struct{}
	wifi_scan               map[string]struct{}
	cell_scan               map[string]struct{}
	work_orders             map[string]struct{}
	floor_plans             map[string]struct{}
	clearedType             bool
	clearedParent           bool
	removedChildren         map[string]struct{}
	removedFiles            map[string]struct{}
	removedHyperlinks       map[string]struct{}
	removedEquipment        map[string]struct{}
	removedProperties       map[string]struct{}
	removedSurvey           map[string]struct{}
	removedWifiScan         map[string]struct{}
	removedCellScan         map[string]struct{}
	removedWorkOrders       map[string]struct{}
	removedFloorPlans       map[string]struct{}
}

// SetName sets the name field.
func (luo *LocationUpdateOne) SetName(s string) *LocationUpdateOne {
	luo.name = &s
	return luo
}

// SetExternalID sets the external_id field.
func (luo *LocationUpdateOne) SetExternalID(s string) *LocationUpdateOne {
	luo.external_id = &s
	return luo
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (luo *LocationUpdateOne) SetNillableExternalID(s *string) *LocationUpdateOne {
	if s != nil {
		luo.SetExternalID(*s)
	}
	return luo
}

// ClearExternalID clears the value of external_id.
func (luo *LocationUpdateOne) ClearExternalID() *LocationUpdateOne {
	luo.external_id = nil
	luo.clearexternal_id = true
	return luo
}

// SetLatitude sets the latitude field.
func (luo *LocationUpdateOne) SetLatitude(f float64) *LocationUpdateOne {
	luo.latitude = &f
	luo.addlatitude = nil
	return luo
}

// SetNillableLatitude sets the latitude field if the given value is not nil.
func (luo *LocationUpdateOne) SetNillableLatitude(f *float64) *LocationUpdateOne {
	if f != nil {
		luo.SetLatitude(*f)
	}
	return luo
}

// AddLatitude adds f to latitude.
func (luo *LocationUpdateOne) AddLatitude(f float64) *LocationUpdateOne {
	if luo.addlatitude == nil {
		luo.addlatitude = &f
	} else {
		*luo.addlatitude += f
	}
	return luo
}

// SetLongitude sets the longitude field.
func (luo *LocationUpdateOne) SetLongitude(f float64) *LocationUpdateOne {
	luo.longitude = &f
	luo.addlongitude = nil
	return luo
}

// SetNillableLongitude sets the longitude field if the given value is not nil.
func (luo *LocationUpdateOne) SetNillableLongitude(f *float64) *LocationUpdateOne {
	if f != nil {
		luo.SetLongitude(*f)
	}
	return luo
}

// AddLongitude adds f to longitude.
func (luo *LocationUpdateOne) AddLongitude(f float64) *LocationUpdateOne {
	if luo.addlongitude == nil {
		luo.addlongitude = &f
	} else {
		*luo.addlongitude += f
	}
	return luo
}

// SetSiteSurveyNeeded sets the site_survey_needed field.
func (luo *LocationUpdateOne) SetSiteSurveyNeeded(b bool) *LocationUpdateOne {
	luo.site_survey_needed = &b
	return luo
}

// SetNillableSiteSurveyNeeded sets the site_survey_needed field if the given value is not nil.
func (luo *LocationUpdateOne) SetNillableSiteSurveyNeeded(b *bool) *LocationUpdateOne {
	if b != nil {
		luo.SetSiteSurveyNeeded(*b)
	}
	return luo
}

// ClearSiteSurveyNeeded clears the value of site_survey_needed.
func (luo *LocationUpdateOne) ClearSiteSurveyNeeded() *LocationUpdateOne {
	luo.site_survey_needed = nil
	luo.clearsite_survey_needed = true
	return luo
}

// SetTypeID sets the type edge to LocationType by id.
func (luo *LocationUpdateOne) SetTypeID(id string) *LocationUpdateOne {
	if luo._type == nil {
		luo._type = make(map[string]struct{})
	}
	luo._type[id] = struct{}{}
	return luo
}

// SetType sets the type edge to LocationType.
func (luo *LocationUpdateOne) SetType(l *LocationType) *LocationUpdateOne {
	return luo.SetTypeID(l.ID)
}

// SetParentID sets the parent edge to Location by id.
func (luo *LocationUpdateOne) SetParentID(id string) *LocationUpdateOne {
	if luo.parent == nil {
		luo.parent = make(map[string]struct{})
	}
	luo.parent[id] = struct{}{}
	return luo
}

// SetNillableParentID sets the parent edge to Location by id if the given value is not nil.
func (luo *LocationUpdateOne) SetNillableParentID(id *string) *LocationUpdateOne {
	if id != nil {
		luo = luo.SetParentID(*id)
	}
	return luo
}

// SetParent sets the parent edge to Location.
func (luo *LocationUpdateOne) SetParent(l *Location) *LocationUpdateOne {
	return luo.SetParentID(l.ID)
}

// AddChildIDs adds the children edge to Location by ids.
func (luo *LocationUpdateOne) AddChildIDs(ids ...string) *LocationUpdateOne {
	if luo.children == nil {
		luo.children = make(map[string]struct{})
	}
	for i := range ids {
		luo.children[ids[i]] = struct{}{}
	}
	return luo
}

// AddChildren adds the children edges to Location.
func (luo *LocationUpdateOne) AddChildren(l ...*Location) *LocationUpdateOne {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return luo.AddChildIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (luo *LocationUpdateOne) AddFileIDs(ids ...string) *LocationUpdateOne {
	if luo.files == nil {
		luo.files = make(map[string]struct{})
	}
	for i := range ids {
		luo.files[ids[i]] = struct{}{}
	}
	return luo
}

// AddFiles adds the files edges to File.
func (luo *LocationUpdateOne) AddFiles(f ...*File) *LocationUpdateOne {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return luo.AddFileIDs(ids...)
}

// AddHyperlinkIDs adds the hyperlinks edge to Hyperlink by ids.
func (luo *LocationUpdateOne) AddHyperlinkIDs(ids ...string) *LocationUpdateOne {
	if luo.hyperlinks == nil {
		luo.hyperlinks = make(map[string]struct{})
	}
	for i := range ids {
		luo.hyperlinks[ids[i]] = struct{}{}
	}
	return luo
}

// AddHyperlinks adds the hyperlinks edges to Hyperlink.
func (luo *LocationUpdateOne) AddHyperlinks(h ...*Hyperlink) *LocationUpdateOne {
	ids := make([]string, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return luo.AddHyperlinkIDs(ids...)
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (luo *LocationUpdateOne) AddEquipmentIDs(ids ...string) *LocationUpdateOne {
	if luo.equipment == nil {
		luo.equipment = make(map[string]struct{})
	}
	for i := range ids {
		luo.equipment[ids[i]] = struct{}{}
	}
	return luo
}

// AddEquipment adds the equipment edges to Equipment.
func (luo *LocationUpdateOne) AddEquipment(e ...*Equipment) *LocationUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return luo.AddEquipmentIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (luo *LocationUpdateOne) AddPropertyIDs(ids ...string) *LocationUpdateOne {
	if luo.properties == nil {
		luo.properties = make(map[string]struct{})
	}
	for i := range ids {
		luo.properties[ids[i]] = struct{}{}
	}
	return luo
}

// AddProperties adds the properties edges to Property.
func (luo *LocationUpdateOne) AddProperties(p ...*Property) *LocationUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return luo.AddPropertyIDs(ids...)
}

// AddSurveyIDs adds the survey edge to Survey by ids.
func (luo *LocationUpdateOne) AddSurveyIDs(ids ...string) *LocationUpdateOne {
	if luo.survey == nil {
		luo.survey = make(map[string]struct{})
	}
	for i := range ids {
		luo.survey[ids[i]] = struct{}{}
	}
	return luo
}

// AddSurvey adds the survey edges to Survey.
func (luo *LocationUpdateOne) AddSurvey(s ...*Survey) *LocationUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return luo.AddSurveyIDs(ids...)
}

// AddWifiScanIDs adds the wifi_scan edge to SurveyWiFiScan by ids.
func (luo *LocationUpdateOne) AddWifiScanIDs(ids ...string) *LocationUpdateOne {
	if luo.wifi_scan == nil {
		luo.wifi_scan = make(map[string]struct{})
	}
	for i := range ids {
		luo.wifi_scan[ids[i]] = struct{}{}
	}
	return luo
}

// AddWifiScan adds the wifi_scan edges to SurveyWiFiScan.
func (luo *LocationUpdateOne) AddWifiScan(s ...*SurveyWiFiScan) *LocationUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return luo.AddWifiScanIDs(ids...)
}

// AddCellScanIDs adds the cell_scan edge to SurveyCellScan by ids.
func (luo *LocationUpdateOne) AddCellScanIDs(ids ...string) *LocationUpdateOne {
	if luo.cell_scan == nil {
		luo.cell_scan = make(map[string]struct{})
	}
	for i := range ids {
		luo.cell_scan[ids[i]] = struct{}{}
	}
	return luo
}

// AddCellScan adds the cell_scan edges to SurveyCellScan.
func (luo *LocationUpdateOne) AddCellScan(s ...*SurveyCellScan) *LocationUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return luo.AddCellScanIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (luo *LocationUpdateOne) AddWorkOrderIDs(ids ...string) *LocationUpdateOne {
	if luo.work_orders == nil {
		luo.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		luo.work_orders[ids[i]] = struct{}{}
	}
	return luo
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (luo *LocationUpdateOne) AddWorkOrders(w ...*WorkOrder) *LocationUpdateOne {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return luo.AddWorkOrderIDs(ids...)
}

// AddFloorPlanIDs adds the floor_plans edge to FloorPlan by ids.
func (luo *LocationUpdateOne) AddFloorPlanIDs(ids ...string) *LocationUpdateOne {
	if luo.floor_plans == nil {
		luo.floor_plans = make(map[string]struct{})
	}
	for i := range ids {
		luo.floor_plans[ids[i]] = struct{}{}
	}
	return luo
}

// AddFloorPlans adds the floor_plans edges to FloorPlan.
func (luo *LocationUpdateOne) AddFloorPlans(f ...*FloorPlan) *LocationUpdateOne {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return luo.AddFloorPlanIDs(ids...)
}

// ClearType clears the type edge to LocationType.
func (luo *LocationUpdateOne) ClearType() *LocationUpdateOne {
	luo.clearedType = true
	return luo
}

// ClearParent clears the parent edge to Location.
func (luo *LocationUpdateOne) ClearParent() *LocationUpdateOne {
	luo.clearedParent = true
	return luo
}

// RemoveChildIDs removes the children edge to Location by ids.
func (luo *LocationUpdateOne) RemoveChildIDs(ids ...string) *LocationUpdateOne {
	if luo.removedChildren == nil {
		luo.removedChildren = make(map[string]struct{})
	}
	for i := range ids {
		luo.removedChildren[ids[i]] = struct{}{}
	}
	return luo
}

// RemoveChildren removes children edges to Location.
func (luo *LocationUpdateOne) RemoveChildren(l ...*Location) *LocationUpdateOne {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return luo.RemoveChildIDs(ids...)
}

// RemoveFileIDs removes the files edge to File by ids.
func (luo *LocationUpdateOne) RemoveFileIDs(ids ...string) *LocationUpdateOne {
	if luo.removedFiles == nil {
		luo.removedFiles = make(map[string]struct{})
	}
	for i := range ids {
		luo.removedFiles[ids[i]] = struct{}{}
	}
	return luo
}

// RemoveFiles removes files edges to File.
func (luo *LocationUpdateOne) RemoveFiles(f ...*File) *LocationUpdateOne {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return luo.RemoveFileIDs(ids...)
}

// RemoveHyperlinkIDs removes the hyperlinks edge to Hyperlink by ids.
func (luo *LocationUpdateOne) RemoveHyperlinkIDs(ids ...string) *LocationUpdateOne {
	if luo.removedHyperlinks == nil {
		luo.removedHyperlinks = make(map[string]struct{})
	}
	for i := range ids {
		luo.removedHyperlinks[ids[i]] = struct{}{}
	}
	return luo
}

// RemoveHyperlinks removes hyperlinks edges to Hyperlink.
func (luo *LocationUpdateOne) RemoveHyperlinks(h ...*Hyperlink) *LocationUpdateOne {
	ids := make([]string, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return luo.RemoveHyperlinkIDs(ids...)
}

// RemoveEquipmentIDs removes the equipment edge to Equipment by ids.
func (luo *LocationUpdateOne) RemoveEquipmentIDs(ids ...string) *LocationUpdateOne {
	if luo.removedEquipment == nil {
		luo.removedEquipment = make(map[string]struct{})
	}
	for i := range ids {
		luo.removedEquipment[ids[i]] = struct{}{}
	}
	return luo
}

// RemoveEquipment removes equipment edges to Equipment.
func (luo *LocationUpdateOne) RemoveEquipment(e ...*Equipment) *LocationUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return luo.RemoveEquipmentIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (luo *LocationUpdateOne) RemovePropertyIDs(ids ...string) *LocationUpdateOne {
	if luo.removedProperties == nil {
		luo.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		luo.removedProperties[ids[i]] = struct{}{}
	}
	return luo
}

// RemoveProperties removes properties edges to Property.
func (luo *LocationUpdateOne) RemoveProperties(p ...*Property) *LocationUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return luo.RemovePropertyIDs(ids...)
}

// RemoveSurveyIDs removes the survey edge to Survey by ids.
func (luo *LocationUpdateOne) RemoveSurveyIDs(ids ...string) *LocationUpdateOne {
	if luo.removedSurvey == nil {
		luo.removedSurvey = make(map[string]struct{})
	}
	for i := range ids {
		luo.removedSurvey[ids[i]] = struct{}{}
	}
	return luo
}

// RemoveSurvey removes survey edges to Survey.
func (luo *LocationUpdateOne) RemoveSurvey(s ...*Survey) *LocationUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return luo.RemoveSurveyIDs(ids...)
}

// RemoveWifiScanIDs removes the wifi_scan edge to SurveyWiFiScan by ids.
func (luo *LocationUpdateOne) RemoveWifiScanIDs(ids ...string) *LocationUpdateOne {
	if luo.removedWifiScan == nil {
		luo.removedWifiScan = make(map[string]struct{})
	}
	for i := range ids {
		luo.removedWifiScan[ids[i]] = struct{}{}
	}
	return luo
}

// RemoveWifiScan removes wifi_scan edges to SurveyWiFiScan.
func (luo *LocationUpdateOne) RemoveWifiScan(s ...*SurveyWiFiScan) *LocationUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return luo.RemoveWifiScanIDs(ids...)
}

// RemoveCellScanIDs removes the cell_scan edge to SurveyCellScan by ids.
func (luo *LocationUpdateOne) RemoveCellScanIDs(ids ...string) *LocationUpdateOne {
	if luo.removedCellScan == nil {
		luo.removedCellScan = make(map[string]struct{})
	}
	for i := range ids {
		luo.removedCellScan[ids[i]] = struct{}{}
	}
	return luo
}

// RemoveCellScan removes cell_scan edges to SurveyCellScan.
func (luo *LocationUpdateOne) RemoveCellScan(s ...*SurveyCellScan) *LocationUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return luo.RemoveCellScanIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (luo *LocationUpdateOne) RemoveWorkOrderIDs(ids ...string) *LocationUpdateOne {
	if luo.removedWorkOrders == nil {
		luo.removedWorkOrders = make(map[string]struct{})
	}
	for i := range ids {
		luo.removedWorkOrders[ids[i]] = struct{}{}
	}
	return luo
}

// RemoveWorkOrders removes work_orders edges to WorkOrder.
func (luo *LocationUpdateOne) RemoveWorkOrders(w ...*WorkOrder) *LocationUpdateOne {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return luo.RemoveWorkOrderIDs(ids...)
}

// RemoveFloorPlanIDs removes the floor_plans edge to FloorPlan by ids.
func (luo *LocationUpdateOne) RemoveFloorPlanIDs(ids ...string) *LocationUpdateOne {
	if luo.removedFloorPlans == nil {
		luo.removedFloorPlans = make(map[string]struct{})
	}
	for i := range ids {
		luo.removedFloorPlans[ids[i]] = struct{}{}
	}
	return luo
}

// RemoveFloorPlans removes floor_plans edges to FloorPlan.
func (luo *LocationUpdateOne) RemoveFloorPlans(f ...*FloorPlan) *LocationUpdateOne {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return luo.RemoveFloorPlanIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (luo *LocationUpdateOne) Save(ctx context.Context) (*Location, error) {
	if luo.update_time == nil {
		v := location.UpdateDefaultUpdateTime()
		luo.update_time = &v
	}
	if luo.name != nil {
		if err := location.NameValidator(*luo.name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if luo.latitude != nil {
		if err := location.LatitudeValidator(*luo.latitude); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"latitude\": %v", err)
		}
	}
	if luo.longitude != nil {
		if err := location.LongitudeValidator(*luo.longitude); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"longitude\": %v", err)
		}
	}
	if len(luo._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if luo.clearedType && luo._type == nil {
		return nil, errors.New("ent: clearing a unique edge \"type\"")
	}
	if len(luo.parent) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"parent\"")
	}
	return luo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (luo *LocationUpdateOne) SaveX(ctx context.Context) *Location {
	l, err := luo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return l
}

// Exec executes the query on the entity.
func (luo *LocationUpdateOne) Exec(ctx context.Context) error {
	_, err := luo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (luo *LocationUpdateOne) ExecX(ctx context.Context) {
	if err := luo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (luo *LocationUpdateOne) sqlSave(ctx context.Context) (l *Location, err error) {
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   location.Table,
			Columns: location.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  luo.id,
				Type:   field.TypeString,
				Column: location.FieldID,
			},
		},
	}
	if value := luo.update_time; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: location.FieldUpdateTime,
		})
	}
	if value := luo.name; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: location.FieldName,
		})
	}
	if value := luo.external_id; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: location.FieldExternalID,
		})
	}
	if luo.clearexternal_id {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: location.FieldExternalID,
		})
	}
	if value := luo.latitude; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: location.FieldLatitude,
		})
	}
	if value := luo.addlatitude; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: location.FieldLatitude,
		})
	}
	if value := luo.longitude; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: location.FieldLongitude,
		})
	}
	if value := luo.addlongitude; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: location.FieldLongitude,
		})
	}
	if value := luo.site_survey_needed; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: location.FieldSiteSurveyNeeded,
		})
	}
	if luo.clearsite_survey_needed {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: location.FieldSiteSurveyNeeded,
		})
	}
	if luo.clearedType {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   location.TypeTable,
			Columns: []string{location.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: locationtype.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := luo._type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   location.TypeTable,
			Columns: []string{location.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: locationtype.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if luo.clearedParent {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   location.ParentTable,
			Columns: []string{location.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := luo.parent; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   location.ParentTable,
			Columns: []string{location.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := luo.removedChildren; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.ChildrenTable,
			Columns: []string{location.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := luo.children; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.ChildrenTable,
			Columns: []string{location.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := luo.removedFiles; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.FilesTable,
			Columns: []string{location.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: file.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := luo.files; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.FilesTable,
			Columns: []string{location.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: file.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := luo.removedHyperlinks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.HyperlinksTable,
			Columns: []string{location.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: hyperlink.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := luo.hyperlinks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.HyperlinksTable,
			Columns: []string{location.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: hyperlink.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := luo.removedEquipment; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.EquipmentTable,
			Columns: []string{location.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipment.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := luo.equipment; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.EquipmentTable,
			Columns: []string{location.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipment.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := luo.removedProperties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.PropertiesTable,
			Columns: []string{location.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := luo.properties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.PropertiesTable,
			Columns: []string{location.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := luo.removedSurvey; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.SurveyTable,
			Columns: []string{location.SurveyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: survey.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := luo.survey; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.SurveyTable,
			Columns: []string{location.SurveyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: survey.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := luo.removedWifiScan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WifiScanTable,
			Columns: []string{location.WifiScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveywifiscan.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := luo.wifi_scan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WifiScanTable,
			Columns: []string{location.WifiScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveywifiscan.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := luo.removedCellScan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.CellScanTable,
			Columns: []string{location.CellScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveycellscan.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := luo.cell_scan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.CellScanTable,
			Columns: []string{location.CellScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveycellscan.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := luo.removedWorkOrders; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WorkOrdersTable,
			Columns: []string{location.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := luo.work_orders; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WorkOrdersTable,
			Columns: []string{location.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := luo.removedFloorPlans; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.FloorPlansTable,
			Columns: []string{location.FloorPlansColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: floorplan.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := luo.floor_plans; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.FloorPlansTable,
			Columns: []string{location.FloorPlansColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: floorplan.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	l = &Location{config: luo.config}
	spec.Assign = l.assignValues
	spec.ScanValues = l.scanValues()
	if err = sqlgraph.UpdateNode(ctx, luo.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return l, nil
}
