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
	equipment               map[string]struct{}
	properties              map[string]struct{}
	survey                  map[string]struct{}
	wifi_scan               map[string]struct{}
	cell_scan               map[string]struct{}
	work_orders             map[string]struct{}
	clearedType             bool
	clearedParent           bool
	removedChildren         map[string]struct{}
	removedFiles            map[string]struct{}
	removedEquipment        map[string]struct{}
	removedProperties       map[string]struct{}
	removedSurvey           map[string]struct{}
	removedWifiScan         map[string]struct{}
	removedCellScan         map[string]struct{}
	removedWorkOrders       map[string]struct{}
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
	var (
		builder  = sql.Dialect(lu.driver.Dialect())
		selector = builder.Select(location.FieldID).From(builder.Table(location.Table))
	)
	for _, p := range lu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = lu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := lu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(location.Table).Where(sql.InInts(location.FieldID, ids...))
	)
	if value := lu.update_time; value != nil {
		updater.Set(location.FieldUpdateTime, *value)
	}
	if value := lu.name; value != nil {
		updater.Set(location.FieldName, *value)
	}
	if value := lu.external_id; value != nil {
		updater.Set(location.FieldExternalID, *value)
	}
	if lu.clearexternal_id {
		updater.SetNull(location.FieldExternalID)
	}
	if value := lu.latitude; value != nil {
		updater.Set(location.FieldLatitude, *value)
	}
	if value := lu.addlatitude; value != nil {
		updater.Add(location.FieldLatitude, *value)
	}
	if value := lu.longitude; value != nil {
		updater.Set(location.FieldLongitude, *value)
	}
	if value := lu.addlongitude; value != nil {
		updater.Add(location.FieldLongitude, *value)
	}
	if value := lu.site_survey_needed; value != nil {
		updater.Set(location.FieldSiteSurveyNeeded, *value)
	}
	if lu.clearsite_survey_needed {
		updater.SetNull(location.FieldSiteSurveyNeeded)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if lu.clearedType {
		query, args := builder.Update(location.TypeTable).
			SetNull(location.TypeColumn).
			Where(sql.InInts(locationtype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(lu._type) > 0 {
		for eid := range lu._type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(location.TypeTable).
				Set(location.TypeColumn, eid).
				Where(sql.InInts(location.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if lu.clearedParent {
		query, args := builder.Update(location.ParentTable).
			SetNull(location.ParentColumn).
			Where(sql.InInts(location.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(lu.parent) > 0 {
		for eid := range lu.parent {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(location.ParentTable).
				Set(location.ParentColumn, eid).
				Where(sql.InInts(location.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if len(lu.removedChildren) > 0 {
		eids := make([]int, len(lu.removedChildren))
		for eid := range lu.removedChildren {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(location.ChildrenTable).
			SetNull(location.ChildrenColumn).
			Where(sql.InInts(location.ChildrenColumn, ids...)).
			Where(sql.InInts(location.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(lu.children) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range lu.children {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(location.FieldID, eid)
			}
			query, args := builder.Update(location.ChildrenTable).
				Set(location.ChildrenColumn, id).
				Where(sql.And(p, sql.IsNull(location.ChildrenColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(lu.children) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"children\" %v already connected to a different \"Location\"", keys(lu.children))})
			}
		}
	}
	if len(lu.removedFiles) > 0 {
		eids := make([]int, len(lu.removedFiles))
		for eid := range lu.removedFiles {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(location.FilesTable).
			SetNull(location.FilesColumn).
			Where(sql.InInts(location.FilesColumn, ids...)).
			Where(sql.InInts(file.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(lu.files) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range lu.files {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(file.FieldID, eid)
			}
			query, args := builder.Update(location.FilesTable).
				Set(location.FilesColumn, id).
				Where(sql.And(p, sql.IsNull(location.FilesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(lu.files) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"files\" %v already connected to a different \"Location\"", keys(lu.files))})
			}
		}
	}
	if len(lu.removedEquipment) > 0 {
		eids := make([]int, len(lu.removedEquipment))
		for eid := range lu.removedEquipment {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(location.EquipmentTable).
			SetNull(location.EquipmentColumn).
			Where(sql.InInts(location.EquipmentColumn, ids...)).
			Where(sql.InInts(equipment.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(lu.equipment) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range lu.equipment {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(equipment.FieldID, eid)
			}
			query, args := builder.Update(location.EquipmentTable).
				Set(location.EquipmentColumn, id).
				Where(sql.And(p, sql.IsNull(location.EquipmentColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(lu.equipment) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"equipment\" %v already connected to a different \"Location\"", keys(lu.equipment))})
			}
		}
	}
	if len(lu.removedProperties) > 0 {
		eids := make([]int, len(lu.removedProperties))
		for eid := range lu.removedProperties {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(location.PropertiesTable).
			SetNull(location.PropertiesColumn).
			Where(sql.InInts(location.PropertiesColumn, ids...)).
			Where(sql.InInts(property.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(lu.properties) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range lu.properties {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(property.FieldID, eid)
			}
			query, args := builder.Update(location.PropertiesTable).
				Set(location.PropertiesColumn, id).
				Where(sql.And(p, sql.IsNull(location.PropertiesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(lu.properties) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"Location\"", keys(lu.properties))})
			}
		}
	}
	if len(lu.removedSurvey) > 0 {
		eids := make([]int, len(lu.removedSurvey))
		for eid := range lu.removedSurvey {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(location.SurveyTable).
			SetNull(location.SurveyColumn).
			Where(sql.InInts(location.SurveyColumn, ids...)).
			Where(sql.InInts(survey.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(lu.survey) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range lu.survey {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(survey.FieldID, eid)
			}
			query, args := builder.Update(location.SurveyTable).
				Set(location.SurveyColumn, id).
				Where(sql.And(p, sql.IsNull(location.SurveyColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(lu.survey) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"survey\" %v already connected to a different \"Location\"", keys(lu.survey))})
			}
		}
	}
	if len(lu.removedWifiScan) > 0 {
		eids := make([]int, len(lu.removedWifiScan))
		for eid := range lu.removedWifiScan {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(location.WifiScanTable).
			SetNull(location.WifiScanColumn).
			Where(sql.InInts(location.WifiScanColumn, ids...)).
			Where(sql.InInts(surveywifiscan.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(lu.wifi_scan) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range lu.wifi_scan {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(surveywifiscan.FieldID, eid)
			}
			query, args := builder.Update(location.WifiScanTable).
				Set(location.WifiScanColumn, id).
				Where(sql.And(p, sql.IsNull(location.WifiScanColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(lu.wifi_scan) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"wifi_scan\" %v already connected to a different \"Location\"", keys(lu.wifi_scan))})
			}
		}
	}
	if len(lu.removedCellScan) > 0 {
		eids := make([]int, len(lu.removedCellScan))
		for eid := range lu.removedCellScan {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(location.CellScanTable).
			SetNull(location.CellScanColumn).
			Where(sql.InInts(location.CellScanColumn, ids...)).
			Where(sql.InInts(surveycellscan.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(lu.cell_scan) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range lu.cell_scan {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(surveycellscan.FieldID, eid)
			}
			query, args := builder.Update(location.CellScanTable).
				Set(location.CellScanColumn, id).
				Where(sql.And(p, sql.IsNull(location.CellScanColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(lu.cell_scan) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"cell_scan\" %v already connected to a different \"Location\"", keys(lu.cell_scan))})
			}
		}
	}
	if len(lu.removedWorkOrders) > 0 {
		eids := make([]int, len(lu.removedWorkOrders))
		for eid := range lu.removedWorkOrders {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(location.WorkOrdersTable).
			SetNull(location.WorkOrdersColumn).
			Where(sql.InInts(location.WorkOrdersColumn, ids...)).
			Where(sql.InInts(workorder.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(lu.work_orders) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range lu.work_orders {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(workorder.FieldID, eid)
			}
			query, args := builder.Update(location.WorkOrdersTable).
				Set(location.WorkOrdersColumn, id).
				Where(sql.And(p, sql.IsNull(location.WorkOrdersColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(lu.work_orders) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"work_orders\" %v already connected to a different \"Location\"", keys(lu.work_orders))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
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
	equipment               map[string]struct{}
	properties              map[string]struct{}
	survey                  map[string]struct{}
	wifi_scan               map[string]struct{}
	cell_scan               map[string]struct{}
	work_orders             map[string]struct{}
	clearedType             bool
	clearedParent           bool
	removedChildren         map[string]struct{}
	removedFiles            map[string]struct{}
	removedEquipment        map[string]struct{}
	removedProperties       map[string]struct{}
	removedSurvey           map[string]struct{}
	removedWifiScan         map[string]struct{}
	removedCellScan         map[string]struct{}
	removedWorkOrders       map[string]struct{}
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
	var (
		builder  = sql.Dialect(luo.driver.Dialect())
		selector = builder.Select(location.Columns...).From(builder.Table(location.Table))
	)
	location.ID(luo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = luo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		l = &Location{config: luo.config}
		if err := l.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into Location: %v", err)
		}
		id = l.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("Location with id: %v", luo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one Location with the same id: %v", luo.id)
	}

	tx, err := luo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(location.Table).Where(sql.InInts(location.FieldID, ids...))
	)
	if value := luo.update_time; value != nil {
		updater.Set(location.FieldUpdateTime, *value)
		l.UpdateTime = *value
	}
	if value := luo.name; value != nil {
		updater.Set(location.FieldName, *value)
		l.Name = *value
	}
	if value := luo.external_id; value != nil {
		updater.Set(location.FieldExternalID, *value)
		l.ExternalID = *value
	}
	if luo.clearexternal_id {
		var value string
		l.ExternalID = value
		updater.SetNull(location.FieldExternalID)
	}
	if value := luo.latitude; value != nil {
		updater.Set(location.FieldLatitude, *value)
		l.Latitude = *value
	}
	if value := luo.addlatitude; value != nil {
		updater.Add(location.FieldLatitude, *value)
		l.Latitude += *value
	}
	if value := luo.longitude; value != nil {
		updater.Set(location.FieldLongitude, *value)
		l.Longitude = *value
	}
	if value := luo.addlongitude; value != nil {
		updater.Add(location.FieldLongitude, *value)
		l.Longitude += *value
	}
	if value := luo.site_survey_needed; value != nil {
		updater.Set(location.FieldSiteSurveyNeeded, *value)
		l.SiteSurveyNeeded = *value
	}
	if luo.clearsite_survey_needed {
		var value bool
		l.SiteSurveyNeeded = value
		updater.SetNull(location.FieldSiteSurveyNeeded)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if luo.clearedType {
		query, args := builder.Update(location.TypeTable).
			SetNull(location.TypeColumn).
			Where(sql.InInts(locationtype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(luo._type) > 0 {
		for eid := range luo._type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(location.TypeTable).
				Set(location.TypeColumn, eid).
				Where(sql.InInts(location.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if luo.clearedParent {
		query, args := builder.Update(location.ParentTable).
			SetNull(location.ParentColumn).
			Where(sql.InInts(location.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(luo.parent) > 0 {
		for eid := range luo.parent {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(location.ParentTable).
				Set(location.ParentColumn, eid).
				Where(sql.InInts(location.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(luo.removedChildren) > 0 {
		eids := make([]int, len(luo.removedChildren))
		for eid := range luo.removedChildren {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(location.ChildrenTable).
			SetNull(location.ChildrenColumn).
			Where(sql.InInts(location.ChildrenColumn, ids...)).
			Where(sql.InInts(location.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(luo.children) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range luo.children {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(luo.children) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"children\" %v already connected to a different \"Location\"", keys(luo.children))})
			}
		}
	}
	if len(luo.removedFiles) > 0 {
		eids := make([]int, len(luo.removedFiles))
		for eid := range luo.removedFiles {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(location.FilesTable).
			SetNull(location.FilesColumn).
			Where(sql.InInts(location.FilesColumn, ids...)).
			Where(sql.InInts(file.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(luo.files) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range luo.files {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(luo.files) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"files\" %v already connected to a different \"Location\"", keys(luo.files))})
			}
		}
	}
	if len(luo.removedEquipment) > 0 {
		eids := make([]int, len(luo.removedEquipment))
		for eid := range luo.removedEquipment {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(location.EquipmentTable).
			SetNull(location.EquipmentColumn).
			Where(sql.InInts(location.EquipmentColumn, ids...)).
			Where(sql.InInts(equipment.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(luo.equipment) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range luo.equipment {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(luo.equipment) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"equipment\" %v already connected to a different \"Location\"", keys(luo.equipment))})
			}
		}
	}
	if len(luo.removedProperties) > 0 {
		eids := make([]int, len(luo.removedProperties))
		for eid := range luo.removedProperties {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(location.PropertiesTable).
			SetNull(location.PropertiesColumn).
			Where(sql.InInts(location.PropertiesColumn, ids...)).
			Where(sql.InInts(property.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(luo.properties) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range luo.properties {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(luo.properties) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"Location\"", keys(luo.properties))})
			}
		}
	}
	if len(luo.removedSurvey) > 0 {
		eids := make([]int, len(luo.removedSurvey))
		for eid := range luo.removedSurvey {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(location.SurveyTable).
			SetNull(location.SurveyColumn).
			Where(sql.InInts(location.SurveyColumn, ids...)).
			Where(sql.InInts(survey.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(luo.survey) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range luo.survey {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(luo.survey) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"survey\" %v already connected to a different \"Location\"", keys(luo.survey))})
			}
		}
	}
	if len(luo.removedWifiScan) > 0 {
		eids := make([]int, len(luo.removedWifiScan))
		for eid := range luo.removedWifiScan {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(location.WifiScanTable).
			SetNull(location.WifiScanColumn).
			Where(sql.InInts(location.WifiScanColumn, ids...)).
			Where(sql.InInts(surveywifiscan.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(luo.wifi_scan) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range luo.wifi_scan {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(luo.wifi_scan) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"wifi_scan\" %v already connected to a different \"Location\"", keys(luo.wifi_scan))})
			}
		}
	}
	if len(luo.removedCellScan) > 0 {
		eids := make([]int, len(luo.removedCellScan))
		for eid := range luo.removedCellScan {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(location.CellScanTable).
			SetNull(location.CellScanColumn).
			Where(sql.InInts(location.CellScanColumn, ids...)).
			Where(sql.InInts(surveycellscan.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(luo.cell_scan) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range luo.cell_scan {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(luo.cell_scan) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"cell_scan\" %v already connected to a different \"Location\"", keys(luo.cell_scan))})
			}
		}
	}
	if len(luo.removedWorkOrders) > 0 {
		eids := make([]int, len(luo.removedWorkOrders))
		for eid := range luo.removedWorkOrders {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(location.WorkOrdersTable).
			SetNull(location.WorkOrdersColumn).
			Where(sql.InInts(location.WorkOrdersColumn, ids...)).
			Where(sql.InInts(workorder.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(luo.work_orders) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range luo.work_orders {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(luo.work_orders) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"work_orders\" %v already connected to a different \"Location\"", keys(luo.work_orders))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return l, nil
}
