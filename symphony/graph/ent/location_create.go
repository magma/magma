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
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/floorplan"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// LocationCreate is the builder for creating a Location entity.
type LocationCreate struct {
	config
	create_time        *time.Time
	update_time        *time.Time
	name               *string
	external_id        *string
	latitude           *float64
	longitude          *float64
	site_survey_needed *bool
	_type              map[string]struct{}
	parent             map[string]struct{}
	children           map[string]struct{}
	files              map[string]struct{}
	equipment          map[string]struct{}
	properties         map[string]struct{}
	survey             map[string]struct{}
	wifi_scan          map[string]struct{}
	cell_scan          map[string]struct{}
	work_orders        map[string]struct{}
	floor_plans        map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (lc *LocationCreate) SetCreateTime(t time.Time) *LocationCreate {
	lc.create_time = &t
	return lc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (lc *LocationCreate) SetNillableCreateTime(t *time.Time) *LocationCreate {
	if t != nil {
		lc.SetCreateTime(*t)
	}
	return lc
}

// SetUpdateTime sets the update_time field.
func (lc *LocationCreate) SetUpdateTime(t time.Time) *LocationCreate {
	lc.update_time = &t
	return lc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (lc *LocationCreate) SetNillableUpdateTime(t *time.Time) *LocationCreate {
	if t != nil {
		lc.SetUpdateTime(*t)
	}
	return lc
}

// SetName sets the name field.
func (lc *LocationCreate) SetName(s string) *LocationCreate {
	lc.name = &s
	return lc
}

// SetExternalID sets the external_id field.
func (lc *LocationCreate) SetExternalID(s string) *LocationCreate {
	lc.external_id = &s
	return lc
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (lc *LocationCreate) SetNillableExternalID(s *string) *LocationCreate {
	if s != nil {
		lc.SetExternalID(*s)
	}
	return lc
}

// SetLatitude sets the latitude field.
func (lc *LocationCreate) SetLatitude(f float64) *LocationCreate {
	lc.latitude = &f
	return lc
}

// SetNillableLatitude sets the latitude field if the given value is not nil.
func (lc *LocationCreate) SetNillableLatitude(f *float64) *LocationCreate {
	if f != nil {
		lc.SetLatitude(*f)
	}
	return lc
}

// SetLongitude sets the longitude field.
func (lc *LocationCreate) SetLongitude(f float64) *LocationCreate {
	lc.longitude = &f
	return lc
}

// SetNillableLongitude sets the longitude field if the given value is not nil.
func (lc *LocationCreate) SetNillableLongitude(f *float64) *LocationCreate {
	if f != nil {
		lc.SetLongitude(*f)
	}
	return lc
}

// SetSiteSurveyNeeded sets the site_survey_needed field.
func (lc *LocationCreate) SetSiteSurveyNeeded(b bool) *LocationCreate {
	lc.site_survey_needed = &b
	return lc
}

// SetNillableSiteSurveyNeeded sets the site_survey_needed field if the given value is not nil.
func (lc *LocationCreate) SetNillableSiteSurveyNeeded(b *bool) *LocationCreate {
	if b != nil {
		lc.SetSiteSurveyNeeded(*b)
	}
	return lc
}

// SetTypeID sets the type edge to LocationType by id.
func (lc *LocationCreate) SetTypeID(id string) *LocationCreate {
	if lc._type == nil {
		lc._type = make(map[string]struct{})
	}
	lc._type[id] = struct{}{}
	return lc
}

// SetType sets the type edge to LocationType.
func (lc *LocationCreate) SetType(l *LocationType) *LocationCreate {
	return lc.SetTypeID(l.ID)
}

// SetParentID sets the parent edge to Location by id.
func (lc *LocationCreate) SetParentID(id string) *LocationCreate {
	if lc.parent == nil {
		lc.parent = make(map[string]struct{})
	}
	lc.parent[id] = struct{}{}
	return lc
}

// SetNillableParentID sets the parent edge to Location by id if the given value is not nil.
func (lc *LocationCreate) SetNillableParentID(id *string) *LocationCreate {
	if id != nil {
		lc = lc.SetParentID(*id)
	}
	return lc
}

// SetParent sets the parent edge to Location.
func (lc *LocationCreate) SetParent(l *Location) *LocationCreate {
	return lc.SetParentID(l.ID)
}

// AddChildIDs adds the children edge to Location by ids.
func (lc *LocationCreate) AddChildIDs(ids ...string) *LocationCreate {
	if lc.children == nil {
		lc.children = make(map[string]struct{})
	}
	for i := range ids {
		lc.children[ids[i]] = struct{}{}
	}
	return lc
}

// AddChildren adds the children edges to Location.
func (lc *LocationCreate) AddChildren(l ...*Location) *LocationCreate {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return lc.AddChildIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (lc *LocationCreate) AddFileIDs(ids ...string) *LocationCreate {
	if lc.files == nil {
		lc.files = make(map[string]struct{})
	}
	for i := range ids {
		lc.files[ids[i]] = struct{}{}
	}
	return lc
}

// AddFiles adds the files edges to File.
func (lc *LocationCreate) AddFiles(f ...*File) *LocationCreate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return lc.AddFileIDs(ids...)
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (lc *LocationCreate) AddEquipmentIDs(ids ...string) *LocationCreate {
	if lc.equipment == nil {
		lc.equipment = make(map[string]struct{})
	}
	for i := range ids {
		lc.equipment[ids[i]] = struct{}{}
	}
	return lc
}

// AddEquipment adds the equipment edges to Equipment.
func (lc *LocationCreate) AddEquipment(e ...*Equipment) *LocationCreate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return lc.AddEquipmentIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (lc *LocationCreate) AddPropertyIDs(ids ...string) *LocationCreate {
	if lc.properties == nil {
		lc.properties = make(map[string]struct{})
	}
	for i := range ids {
		lc.properties[ids[i]] = struct{}{}
	}
	return lc
}

// AddProperties adds the properties edges to Property.
func (lc *LocationCreate) AddProperties(p ...*Property) *LocationCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return lc.AddPropertyIDs(ids...)
}

// AddSurveyIDs adds the survey edge to Survey by ids.
func (lc *LocationCreate) AddSurveyIDs(ids ...string) *LocationCreate {
	if lc.survey == nil {
		lc.survey = make(map[string]struct{})
	}
	for i := range ids {
		lc.survey[ids[i]] = struct{}{}
	}
	return lc
}

// AddSurvey adds the survey edges to Survey.
func (lc *LocationCreate) AddSurvey(s ...*Survey) *LocationCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lc.AddSurveyIDs(ids...)
}

// AddWifiScanIDs adds the wifi_scan edge to SurveyWiFiScan by ids.
func (lc *LocationCreate) AddWifiScanIDs(ids ...string) *LocationCreate {
	if lc.wifi_scan == nil {
		lc.wifi_scan = make(map[string]struct{})
	}
	for i := range ids {
		lc.wifi_scan[ids[i]] = struct{}{}
	}
	return lc
}

// AddWifiScan adds the wifi_scan edges to SurveyWiFiScan.
func (lc *LocationCreate) AddWifiScan(s ...*SurveyWiFiScan) *LocationCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lc.AddWifiScanIDs(ids...)
}

// AddCellScanIDs adds the cell_scan edge to SurveyCellScan by ids.
func (lc *LocationCreate) AddCellScanIDs(ids ...string) *LocationCreate {
	if lc.cell_scan == nil {
		lc.cell_scan = make(map[string]struct{})
	}
	for i := range ids {
		lc.cell_scan[ids[i]] = struct{}{}
	}
	return lc
}

// AddCellScan adds the cell_scan edges to SurveyCellScan.
func (lc *LocationCreate) AddCellScan(s ...*SurveyCellScan) *LocationCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lc.AddCellScanIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (lc *LocationCreate) AddWorkOrderIDs(ids ...string) *LocationCreate {
	if lc.work_orders == nil {
		lc.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		lc.work_orders[ids[i]] = struct{}{}
	}
	return lc
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (lc *LocationCreate) AddWorkOrders(w ...*WorkOrder) *LocationCreate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return lc.AddWorkOrderIDs(ids...)
}

// AddFloorPlanIDs adds the floor_plans edge to FloorPlan by ids.
func (lc *LocationCreate) AddFloorPlanIDs(ids ...string) *LocationCreate {
	if lc.floor_plans == nil {
		lc.floor_plans = make(map[string]struct{})
	}
	for i := range ids {
		lc.floor_plans[ids[i]] = struct{}{}
	}
	return lc
}

// AddFloorPlans adds the floor_plans edges to FloorPlan.
func (lc *LocationCreate) AddFloorPlans(f ...*FloorPlan) *LocationCreate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return lc.AddFloorPlanIDs(ids...)
}

// Save creates the Location in the database.
func (lc *LocationCreate) Save(ctx context.Context) (*Location, error) {
	if lc.create_time == nil {
		v := location.DefaultCreateTime()
		lc.create_time = &v
	}
	if lc.update_time == nil {
		v := location.DefaultUpdateTime()
		lc.update_time = &v
	}
	if lc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if err := location.NameValidator(*lc.name); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
	}
	if lc.latitude == nil {
		v := location.DefaultLatitude
		lc.latitude = &v
	}
	if err := location.LatitudeValidator(*lc.latitude); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"latitude\": %v", err)
	}
	if lc.longitude == nil {
		v := location.DefaultLongitude
		lc.longitude = &v
	}
	if err := location.LongitudeValidator(*lc.longitude); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"longitude\": %v", err)
	}
	if lc.site_survey_needed == nil {
		v := location.DefaultSiteSurveyNeeded
		lc.site_survey_needed = &v
	}
	if len(lc._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if lc._type == nil {
		return nil, errors.New("ent: missing required edge \"type\"")
	}
	if len(lc.parent) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"parent\"")
	}
	return lc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (lc *LocationCreate) SaveX(ctx context.Context) *Location {
	v, err := lc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (lc *LocationCreate) sqlSave(ctx context.Context) (*Location, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(lc.driver.Dialect())
		l       = &Location{config: lc.config}
	)
	tx, err := lc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(location.Table).Default()
	if value := lc.create_time; value != nil {
		insert.Set(location.FieldCreateTime, *value)
		l.CreateTime = *value
	}
	if value := lc.update_time; value != nil {
		insert.Set(location.FieldUpdateTime, *value)
		l.UpdateTime = *value
	}
	if value := lc.name; value != nil {
		insert.Set(location.FieldName, *value)
		l.Name = *value
	}
	if value := lc.external_id; value != nil {
		insert.Set(location.FieldExternalID, *value)
		l.ExternalID = *value
	}
	if value := lc.latitude; value != nil {
		insert.Set(location.FieldLatitude, *value)
		l.Latitude = *value
	}
	if value := lc.longitude; value != nil {
		insert.Set(location.FieldLongitude, *value)
		l.Longitude = *value
	}
	if value := lc.site_survey_needed; value != nil {
		insert.Set(location.FieldSiteSurveyNeeded, *value)
		l.SiteSurveyNeeded = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(location.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	l.ID = strconv.FormatInt(id, 10)
	if len(lc._type) > 0 {
		for eid := range lc._type {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(location.TypeTable).
				Set(location.TypeColumn, eid).
				Where(sql.EQ(location.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(lc.parent) > 0 {
		for eid := range lc.parent {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(location.ParentTable).
				Set(location.ParentColumn, eid).
				Where(sql.EQ(location.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(lc.children) > 0 {
		p := sql.P()
		for eid := range lc.children {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(location.FieldID, eid)
		}
		query, args := builder.Update(location.ChildrenTable).
			Set(location.ChildrenColumn, id).
			Where(sql.And(p, sql.IsNull(location.ChildrenColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(lc.children) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"children\" %v already connected to a different \"Location\"", keys(lc.children))})
		}
	}
	if len(lc.files) > 0 {
		p := sql.P()
		for eid := range lc.files {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(file.FieldID, eid)
		}
		query, args := builder.Update(location.FilesTable).
			Set(location.FilesColumn, id).
			Where(sql.And(p, sql.IsNull(location.FilesColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(lc.files) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"files\" %v already connected to a different \"Location\"", keys(lc.files))})
		}
	}
	if len(lc.equipment) > 0 {
		p := sql.P()
		for eid := range lc.equipment {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(equipment.FieldID, eid)
		}
		query, args := builder.Update(location.EquipmentTable).
			Set(location.EquipmentColumn, id).
			Where(sql.And(p, sql.IsNull(location.EquipmentColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(lc.equipment) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"equipment\" %v already connected to a different \"Location\"", keys(lc.equipment))})
		}
	}
	if len(lc.properties) > 0 {
		p := sql.P()
		for eid := range lc.properties {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(property.FieldID, eid)
		}
		query, args := builder.Update(location.PropertiesTable).
			Set(location.PropertiesColumn, id).
			Where(sql.And(p, sql.IsNull(location.PropertiesColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(lc.properties) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"Location\"", keys(lc.properties))})
		}
	}
	if len(lc.survey) > 0 {
		p := sql.P()
		for eid := range lc.survey {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(survey.FieldID, eid)
		}
		query, args := builder.Update(location.SurveyTable).
			Set(location.SurveyColumn, id).
			Where(sql.And(p, sql.IsNull(location.SurveyColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(lc.survey) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"survey\" %v already connected to a different \"Location\"", keys(lc.survey))})
		}
	}
	if len(lc.wifi_scan) > 0 {
		p := sql.P()
		for eid := range lc.wifi_scan {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(surveywifiscan.FieldID, eid)
		}
		query, args := builder.Update(location.WifiScanTable).
			Set(location.WifiScanColumn, id).
			Where(sql.And(p, sql.IsNull(location.WifiScanColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(lc.wifi_scan) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"wifi_scan\" %v already connected to a different \"Location\"", keys(lc.wifi_scan))})
		}
	}
	if len(lc.cell_scan) > 0 {
		p := sql.P()
		for eid := range lc.cell_scan {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(surveycellscan.FieldID, eid)
		}
		query, args := builder.Update(location.CellScanTable).
			Set(location.CellScanColumn, id).
			Where(sql.And(p, sql.IsNull(location.CellScanColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(lc.cell_scan) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"cell_scan\" %v already connected to a different \"Location\"", keys(lc.cell_scan))})
		}
	}
	if len(lc.work_orders) > 0 {
		p := sql.P()
		for eid := range lc.work_orders {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(workorder.FieldID, eid)
		}
		query, args := builder.Update(location.WorkOrdersTable).
			Set(location.WorkOrdersColumn, id).
			Where(sql.And(p, sql.IsNull(location.WorkOrdersColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(lc.work_orders) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"work_orders\" %v already connected to a different \"Location\"", keys(lc.work_orders))})
		}
	}
	if len(lc.floor_plans) > 0 {
		p := sql.P()
		for eid := range lc.floor_plans {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(floorplan.FieldID, eid)
		}
		query, args := builder.Update(location.FloorPlansTable).
			Set(location.FloorPlansColumn, id).
			Where(sql.And(p, sql.IsNull(location.FloorPlansColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(lc.floor_plans) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"floor_plans\" %v already connected to a different \"Location\"", keys(lc.floor_plans))})
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return l, nil
}
