// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
)

// SurveyCellScanUpdate is the builder for updating SurveyCellScan entities.
type SurveyCellScanUpdate struct {
	config

	update_time                  *time.Time
	network_type                 *string
	signal_strength              *int
	addsignal_strength           *int
	timestamp                    *time.Time
	cleartimestamp               bool
	base_station_id              *string
	clearbase_station_id         bool
	network_id                   *string
	clearnetwork_id              bool
	system_id                    *string
	clearsystem_id               bool
	cell_id                      *string
	clearcell_id                 bool
	location_area_code           *string
	clearlocation_area_code      bool
	mobile_country_code          *string
	clearmobile_country_code     bool
	mobile_network_code          *string
	clearmobile_network_code     bool
	primary_scrambling_code      *string
	clearprimary_scrambling_code bool
	operator                     *string
	clearoperator                bool
	arfcn                        *int
	addarfcn                     *int
	cleararfcn                   bool
	physical_cell_id             *string
	clearphysical_cell_id        bool
	tracking_area_code           *string
	cleartracking_area_code      bool
	timing_advance               *int
	addtiming_advance            *int
	cleartiming_advance          bool
	earfcn                       *int
	addearfcn                    *int
	clearearfcn                  bool
	uarfcn                       *int
	adduarfcn                    *int
	clearuarfcn                  bool
	latitude                     *float64
	addlatitude                  *float64
	clearlatitude                bool
	longitude                    *float64
	addlongitude                 *float64
	clearlongitude               bool
	survey_question              map[string]struct{}
	location                     map[string]struct{}
	clearedSurveyQuestion        bool
	clearedLocation              bool
	predicates                   []predicate.SurveyCellScan
}

// Where adds a new predicate for the builder.
func (scsu *SurveyCellScanUpdate) Where(ps ...predicate.SurveyCellScan) *SurveyCellScanUpdate {
	scsu.predicates = append(scsu.predicates, ps...)
	return scsu
}

// SetNetworkType sets the network_type field.
func (scsu *SurveyCellScanUpdate) SetNetworkType(s string) *SurveyCellScanUpdate {
	scsu.network_type = &s
	return scsu
}

// SetSignalStrength sets the signal_strength field.
func (scsu *SurveyCellScanUpdate) SetSignalStrength(i int) *SurveyCellScanUpdate {
	scsu.signal_strength = &i
	scsu.addsignal_strength = nil
	return scsu
}

// AddSignalStrength adds i to signal_strength.
func (scsu *SurveyCellScanUpdate) AddSignalStrength(i int) *SurveyCellScanUpdate {
	if scsu.addsignal_strength == nil {
		scsu.addsignal_strength = &i
	} else {
		*scsu.addsignal_strength += i
	}
	return scsu
}

// SetTimestamp sets the timestamp field.
func (scsu *SurveyCellScanUpdate) SetTimestamp(t time.Time) *SurveyCellScanUpdate {
	scsu.timestamp = &t
	return scsu
}

// SetNillableTimestamp sets the timestamp field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableTimestamp(t *time.Time) *SurveyCellScanUpdate {
	if t != nil {
		scsu.SetTimestamp(*t)
	}
	return scsu
}

// ClearTimestamp clears the value of timestamp.
func (scsu *SurveyCellScanUpdate) ClearTimestamp() *SurveyCellScanUpdate {
	scsu.timestamp = nil
	scsu.cleartimestamp = true
	return scsu
}

// SetBaseStationID sets the base_station_id field.
func (scsu *SurveyCellScanUpdate) SetBaseStationID(s string) *SurveyCellScanUpdate {
	scsu.base_station_id = &s
	return scsu
}

// SetNillableBaseStationID sets the base_station_id field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableBaseStationID(s *string) *SurveyCellScanUpdate {
	if s != nil {
		scsu.SetBaseStationID(*s)
	}
	return scsu
}

// ClearBaseStationID clears the value of base_station_id.
func (scsu *SurveyCellScanUpdate) ClearBaseStationID() *SurveyCellScanUpdate {
	scsu.base_station_id = nil
	scsu.clearbase_station_id = true
	return scsu
}

// SetNetworkID sets the network_id field.
func (scsu *SurveyCellScanUpdate) SetNetworkID(s string) *SurveyCellScanUpdate {
	scsu.network_id = &s
	return scsu
}

// SetNillableNetworkID sets the network_id field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableNetworkID(s *string) *SurveyCellScanUpdate {
	if s != nil {
		scsu.SetNetworkID(*s)
	}
	return scsu
}

// ClearNetworkID clears the value of network_id.
func (scsu *SurveyCellScanUpdate) ClearNetworkID() *SurveyCellScanUpdate {
	scsu.network_id = nil
	scsu.clearnetwork_id = true
	return scsu
}

// SetSystemID sets the system_id field.
func (scsu *SurveyCellScanUpdate) SetSystemID(s string) *SurveyCellScanUpdate {
	scsu.system_id = &s
	return scsu
}

// SetNillableSystemID sets the system_id field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableSystemID(s *string) *SurveyCellScanUpdate {
	if s != nil {
		scsu.SetSystemID(*s)
	}
	return scsu
}

// ClearSystemID clears the value of system_id.
func (scsu *SurveyCellScanUpdate) ClearSystemID() *SurveyCellScanUpdate {
	scsu.system_id = nil
	scsu.clearsystem_id = true
	return scsu
}

// SetCellID sets the cell_id field.
func (scsu *SurveyCellScanUpdate) SetCellID(s string) *SurveyCellScanUpdate {
	scsu.cell_id = &s
	return scsu
}

// SetNillableCellID sets the cell_id field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableCellID(s *string) *SurveyCellScanUpdate {
	if s != nil {
		scsu.SetCellID(*s)
	}
	return scsu
}

// ClearCellID clears the value of cell_id.
func (scsu *SurveyCellScanUpdate) ClearCellID() *SurveyCellScanUpdate {
	scsu.cell_id = nil
	scsu.clearcell_id = true
	return scsu
}

// SetLocationAreaCode sets the location_area_code field.
func (scsu *SurveyCellScanUpdate) SetLocationAreaCode(s string) *SurveyCellScanUpdate {
	scsu.location_area_code = &s
	return scsu
}

// SetNillableLocationAreaCode sets the location_area_code field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableLocationAreaCode(s *string) *SurveyCellScanUpdate {
	if s != nil {
		scsu.SetLocationAreaCode(*s)
	}
	return scsu
}

// ClearLocationAreaCode clears the value of location_area_code.
func (scsu *SurveyCellScanUpdate) ClearLocationAreaCode() *SurveyCellScanUpdate {
	scsu.location_area_code = nil
	scsu.clearlocation_area_code = true
	return scsu
}

// SetMobileCountryCode sets the mobile_country_code field.
func (scsu *SurveyCellScanUpdate) SetMobileCountryCode(s string) *SurveyCellScanUpdate {
	scsu.mobile_country_code = &s
	return scsu
}

// SetNillableMobileCountryCode sets the mobile_country_code field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableMobileCountryCode(s *string) *SurveyCellScanUpdate {
	if s != nil {
		scsu.SetMobileCountryCode(*s)
	}
	return scsu
}

// ClearMobileCountryCode clears the value of mobile_country_code.
func (scsu *SurveyCellScanUpdate) ClearMobileCountryCode() *SurveyCellScanUpdate {
	scsu.mobile_country_code = nil
	scsu.clearmobile_country_code = true
	return scsu
}

// SetMobileNetworkCode sets the mobile_network_code field.
func (scsu *SurveyCellScanUpdate) SetMobileNetworkCode(s string) *SurveyCellScanUpdate {
	scsu.mobile_network_code = &s
	return scsu
}

// SetNillableMobileNetworkCode sets the mobile_network_code field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableMobileNetworkCode(s *string) *SurveyCellScanUpdate {
	if s != nil {
		scsu.SetMobileNetworkCode(*s)
	}
	return scsu
}

// ClearMobileNetworkCode clears the value of mobile_network_code.
func (scsu *SurveyCellScanUpdate) ClearMobileNetworkCode() *SurveyCellScanUpdate {
	scsu.mobile_network_code = nil
	scsu.clearmobile_network_code = true
	return scsu
}

// SetPrimaryScramblingCode sets the primary_scrambling_code field.
func (scsu *SurveyCellScanUpdate) SetPrimaryScramblingCode(s string) *SurveyCellScanUpdate {
	scsu.primary_scrambling_code = &s
	return scsu
}

// SetNillablePrimaryScramblingCode sets the primary_scrambling_code field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillablePrimaryScramblingCode(s *string) *SurveyCellScanUpdate {
	if s != nil {
		scsu.SetPrimaryScramblingCode(*s)
	}
	return scsu
}

// ClearPrimaryScramblingCode clears the value of primary_scrambling_code.
func (scsu *SurveyCellScanUpdate) ClearPrimaryScramblingCode() *SurveyCellScanUpdate {
	scsu.primary_scrambling_code = nil
	scsu.clearprimary_scrambling_code = true
	return scsu
}

// SetOperator sets the operator field.
func (scsu *SurveyCellScanUpdate) SetOperator(s string) *SurveyCellScanUpdate {
	scsu.operator = &s
	return scsu
}

// SetNillableOperator sets the operator field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableOperator(s *string) *SurveyCellScanUpdate {
	if s != nil {
		scsu.SetOperator(*s)
	}
	return scsu
}

// ClearOperator clears the value of operator.
func (scsu *SurveyCellScanUpdate) ClearOperator() *SurveyCellScanUpdate {
	scsu.operator = nil
	scsu.clearoperator = true
	return scsu
}

// SetArfcn sets the arfcn field.
func (scsu *SurveyCellScanUpdate) SetArfcn(i int) *SurveyCellScanUpdate {
	scsu.arfcn = &i
	scsu.addarfcn = nil
	return scsu
}

// SetNillableArfcn sets the arfcn field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableArfcn(i *int) *SurveyCellScanUpdate {
	if i != nil {
		scsu.SetArfcn(*i)
	}
	return scsu
}

// AddArfcn adds i to arfcn.
func (scsu *SurveyCellScanUpdate) AddArfcn(i int) *SurveyCellScanUpdate {
	if scsu.addarfcn == nil {
		scsu.addarfcn = &i
	} else {
		*scsu.addarfcn += i
	}
	return scsu
}

// ClearArfcn clears the value of arfcn.
func (scsu *SurveyCellScanUpdate) ClearArfcn() *SurveyCellScanUpdate {
	scsu.arfcn = nil
	scsu.cleararfcn = true
	return scsu
}

// SetPhysicalCellID sets the physical_cell_id field.
func (scsu *SurveyCellScanUpdate) SetPhysicalCellID(s string) *SurveyCellScanUpdate {
	scsu.physical_cell_id = &s
	return scsu
}

// SetNillablePhysicalCellID sets the physical_cell_id field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillablePhysicalCellID(s *string) *SurveyCellScanUpdate {
	if s != nil {
		scsu.SetPhysicalCellID(*s)
	}
	return scsu
}

// ClearPhysicalCellID clears the value of physical_cell_id.
func (scsu *SurveyCellScanUpdate) ClearPhysicalCellID() *SurveyCellScanUpdate {
	scsu.physical_cell_id = nil
	scsu.clearphysical_cell_id = true
	return scsu
}

// SetTrackingAreaCode sets the tracking_area_code field.
func (scsu *SurveyCellScanUpdate) SetTrackingAreaCode(s string) *SurveyCellScanUpdate {
	scsu.tracking_area_code = &s
	return scsu
}

// SetNillableTrackingAreaCode sets the tracking_area_code field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableTrackingAreaCode(s *string) *SurveyCellScanUpdate {
	if s != nil {
		scsu.SetTrackingAreaCode(*s)
	}
	return scsu
}

// ClearTrackingAreaCode clears the value of tracking_area_code.
func (scsu *SurveyCellScanUpdate) ClearTrackingAreaCode() *SurveyCellScanUpdate {
	scsu.tracking_area_code = nil
	scsu.cleartracking_area_code = true
	return scsu
}

// SetTimingAdvance sets the timing_advance field.
func (scsu *SurveyCellScanUpdate) SetTimingAdvance(i int) *SurveyCellScanUpdate {
	scsu.timing_advance = &i
	scsu.addtiming_advance = nil
	return scsu
}

// SetNillableTimingAdvance sets the timing_advance field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableTimingAdvance(i *int) *SurveyCellScanUpdate {
	if i != nil {
		scsu.SetTimingAdvance(*i)
	}
	return scsu
}

// AddTimingAdvance adds i to timing_advance.
func (scsu *SurveyCellScanUpdate) AddTimingAdvance(i int) *SurveyCellScanUpdate {
	if scsu.addtiming_advance == nil {
		scsu.addtiming_advance = &i
	} else {
		*scsu.addtiming_advance += i
	}
	return scsu
}

// ClearTimingAdvance clears the value of timing_advance.
func (scsu *SurveyCellScanUpdate) ClearTimingAdvance() *SurveyCellScanUpdate {
	scsu.timing_advance = nil
	scsu.cleartiming_advance = true
	return scsu
}

// SetEarfcn sets the earfcn field.
func (scsu *SurveyCellScanUpdate) SetEarfcn(i int) *SurveyCellScanUpdate {
	scsu.earfcn = &i
	scsu.addearfcn = nil
	return scsu
}

// SetNillableEarfcn sets the earfcn field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableEarfcn(i *int) *SurveyCellScanUpdate {
	if i != nil {
		scsu.SetEarfcn(*i)
	}
	return scsu
}

// AddEarfcn adds i to earfcn.
func (scsu *SurveyCellScanUpdate) AddEarfcn(i int) *SurveyCellScanUpdate {
	if scsu.addearfcn == nil {
		scsu.addearfcn = &i
	} else {
		*scsu.addearfcn += i
	}
	return scsu
}

// ClearEarfcn clears the value of earfcn.
func (scsu *SurveyCellScanUpdate) ClearEarfcn() *SurveyCellScanUpdate {
	scsu.earfcn = nil
	scsu.clearearfcn = true
	return scsu
}

// SetUarfcn sets the uarfcn field.
func (scsu *SurveyCellScanUpdate) SetUarfcn(i int) *SurveyCellScanUpdate {
	scsu.uarfcn = &i
	scsu.adduarfcn = nil
	return scsu
}

// SetNillableUarfcn sets the uarfcn field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableUarfcn(i *int) *SurveyCellScanUpdate {
	if i != nil {
		scsu.SetUarfcn(*i)
	}
	return scsu
}

// AddUarfcn adds i to uarfcn.
func (scsu *SurveyCellScanUpdate) AddUarfcn(i int) *SurveyCellScanUpdate {
	if scsu.adduarfcn == nil {
		scsu.adduarfcn = &i
	} else {
		*scsu.adduarfcn += i
	}
	return scsu
}

// ClearUarfcn clears the value of uarfcn.
func (scsu *SurveyCellScanUpdate) ClearUarfcn() *SurveyCellScanUpdate {
	scsu.uarfcn = nil
	scsu.clearuarfcn = true
	return scsu
}

// SetLatitude sets the latitude field.
func (scsu *SurveyCellScanUpdate) SetLatitude(f float64) *SurveyCellScanUpdate {
	scsu.latitude = &f
	scsu.addlatitude = nil
	return scsu
}

// SetNillableLatitude sets the latitude field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableLatitude(f *float64) *SurveyCellScanUpdate {
	if f != nil {
		scsu.SetLatitude(*f)
	}
	return scsu
}

// AddLatitude adds f to latitude.
func (scsu *SurveyCellScanUpdate) AddLatitude(f float64) *SurveyCellScanUpdate {
	if scsu.addlatitude == nil {
		scsu.addlatitude = &f
	} else {
		*scsu.addlatitude += f
	}
	return scsu
}

// ClearLatitude clears the value of latitude.
func (scsu *SurveyCellScanUpdate) ClearLatitude() *SurveyCellScanUpdate {
	scsu.latitude = nil
	scsu.clearlatitude = true
	return scsu
}

// SetLongitude sets the longitude field.
func (scsu *SurveyCellScanUpdate) SetLongitude(f float64) *SurveyCellScanUpdate {
	scsu.longitude = &f
	scsu.addlongitude = nil
	return scsu
}

// SetNillableLongitude sets the longitude field if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableLongitude(f *float64) *SurveyCellScanUpdate {
	if f != nil {
		scsu.SetLongitude(*f)
	}
	return scsu
}

// AddLongitude adds f to longitude.
func (scsu *SurveyCellScanUpdate) AddLongitude(f float64) *SurveyCellScanUpdate {
	if scsu.addlongitude == nil {
		scsu.addlongitude = &f
	} else {
		*scsu.addlongitude += f
	}
	return scsu
}

// ClearLongitude clears the value of longitude.
func (scsu *SurveyCellScanUpdate) ClearLongitude() *SurveyCellScanUpdate {
	scsu.longitude = nil
	scsu.clearlongitude = true
	return scsu
}

// SetSurveyQuestionID sets the survey_question edge to SurveyQuestion by id.
func (scsu *SurveyCellScanUpdate) SetSurveyQuestionID(id string) *SurveyCellScanUpdate {
	if scsu.survey_question == nil {
		scsu.survey_question = make(map[string]struct{})
	}
	scsu.survey_question[id] = struct{}{}
	return scsu
}

// SetNillableSurveyQuestionID sets the survey_question edge to SurveyQuestion by id if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableSurveyQuestionID(id *string) *SurveyCellScanUpdate {
	if id != nil {
		scsu = scsu.SetSurveyQuestionID(*id)
	}
	return scsu
}

// SetSurveyQuestion sets the survey_question edge to SurveyQuestion.
func (scsu *SurveyCellScanUpdate) SetSurveyQuestion(s *SurveyQuestion) *SurveyCellScanUpdate {
	return scsu.SetSurveyQuestionID(s.ID)
}

// SetLocationID sets the location edge to Location by id.
func (scsu *SurveyCellScanUpdate) SetLocationID(id string) *SurveyCellScanUpdate {
	if scsu.location == nil {
		scsu.location = make(map[string]struct{})
	}
	scsu.location[id] = struct{}{}
	return scsu
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableLocationID(id *string) *SurveyCellScanUpdate {
	if id != nil {
		scsu = scsu.SetLocationID(*id)
	}
	return scsu
}

// SetLocation sets the location edge to Location.
func (scsu *SurveyCellScanUpdate) SetLocation(l *Location) *SurveyCellScanUpdate {
	return scsu.SetLocationID(l.ID)
}

// ClearSurveyQuestion clears the survey_question edge to SurveyQuestion.
func (scsu *SurveyCellScanUpdate) ClearSurveyQuestion() *SurveyCellScanUpdate {
	scsu.clearedSurveyQuestion = true
	return scsu
}

// ClearLocation clears the location edge to Location.
func (scsu *SurveyCellScanUpdate) ClearLocation() *SurveyCellScanUpdate {
	scsu.clearedLocation = true
	return scsu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (scsu *SurveyCellScanUpdate) Save(ctx context.Context) (int, error) {
	if scsu.update_time == nil {
		v := surveycellscan.UpdateDefaultUpdateTime()
		scsu.update_time = &v
	}
	if len(scsu.survey_question) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"survey_question\"")
	}
	if len(scsu.location) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	return scsu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (scsu *SurveyCellScanUpdate) SaveX(ctx context.Context) int {
	affected, err := scsu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (scsu *SurveyCellScanUpdate) Exec(ctx context.Context) error {
	_, err := scsu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scsu *SurveyCellScanUpdate) ExecX(ctx context.Context) {
	if err := scsu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (scsu *SurveyCellScanUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveycellscan.Table,
			Columns: surveycellscan.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: surveycellscan.FieldID,
			},
		},
	}
	if ps := scsu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := scsu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveycellscan.FieldUpdateTime,
		})
	}
	if value := scsu.network_type; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldNetworkType,
		})
	}
	if value := scsu.signal_strength; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldSignalStrength,
		})
	}
	if value := scsu.addsignal_strength; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldSignalStrength,
		})
	}
	if value := scsu.timestamp; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveycellscan.FieldTimestamp,
		})
	}
	if scsu.cleartimestamp {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: surveycellscan.FieldTimestamp,
		})
	}
	if value := scsu.base_station_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldBaseStationID,
		})
	}
	if scsu.clearbase_station_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldBaseStationID,
		})
	}
	if value := scsu.network_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldNetworkID,
		})
	}
	if scsu.clearnetwork_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldNetworkID,
		})
	}
	if value := scsu.system_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldSystemID,
		})
	}
	if scsu.clearsystem_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldSystemID,
		})
	}
	if value := scsu.cell_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldCellID,
		})
	}
	if scsu.clearcell_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldCellID,
		})
	}
	if value := scsu.location_area_code; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldLocationAreaCode,
		})
	}
	if scsu.clearlocation_area_code {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldLocationAreaCode,
		})
	}
	if value := scsu.mobile_country_code; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldMobileCountryCode,
		})
	}
	if scsu.clearmobile_country_code {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldMobileCountryCode,
		})
	}
	if value := scsu.mobile_network_code; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldMobileNetworkCode,
		})
	}
	if scsu.clearmobile_network_code {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldMobileNetworkCode,
		})
	}
	if value := scsu.primary_scrambling_code; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldPrimaryScramblingCode,
		})
	}
	if scsu.clearprimary_scrambling_code {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldPrimaryScramblingCode,
		})
	}
	if value := scsu.operator; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldOperator,
		})
	}
	if scsu.clearoperator {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldOperator,
		})
	}
	if value := scsu.arfcn; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldArfcn,
		})
	}
	if value := scsu.addarfcn; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldArfcn,
		})
	}
	if scsu.cleararfcn {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveycellscan.FieldArfcn,
		})
	}
	if value := scsu.physical_cell_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldPhysicalCellID,
		})
	}
	if scsu.clearphysical_cell_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldPhysicalCellID,
		})
	}
	if value := scsu.tracking_area_code; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldTrackingAreaCode,
		})
	}
	if scsu.cleartracking_area_code {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldTrackingAreaCode,
		})
	}
	if value := scsu.timing_advance; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldTimingAdvance,
		})
	}
	if value := scsu.addtiming_advance; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldTimingAdvance,
		})
	}
	if scsu.cleartiming_advance {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveycellscan.FieldTimingAdvance,
		})
	}
	if value := scsu.earfcn; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldEarfcn,
		})
	}
	if value := scsu.addearfcn; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldEarfcn,
		})
	}
	if scsu.clearearfcn {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveycellscan.FieldEarfcn,
		})
	}
	if value := scsu.uarfcn; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldUarfcn,
		})
	}
	if value := scsu.adduarfcn; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldUarfcn,
		})
	}
	if scsu.clearuarfcn {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveycellscan.FieldUarfcn,
		})
	}
	if value := scsu.latitude; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveycellscan.FieldLatitude,
		})
	}
	if value := scsu.addlatitude; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveycellscan.FieldLatitude,
		})
	}
	if scsu.clearlatitude {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveycellscan.FieldLatitude,
		})
	}
	if value := scsu.longitude; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveycellscan.FieldLongitude,
		})
	}
	if value := scsu.addlongitude; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveycellscan.FieldLongitude,
		})
	}
	if scsu.clearlongitude {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveycellscan.FieldLongitude,
		})
	}
	if scsu.clearedSurveyQuestion {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveycellscan.SurveyQuestionTable,
			Columns: []string{surveycellscan.SurveyQuestionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveyquestion.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := scsu.survey_question; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveycellscan.SurveyQuestionTable,
			Columns: []string{surveycellscan.SurveyQuestionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveyquestion.FieldID,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if scsu.clearedLocation {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveycellscan.LocationTable,
			Columns: []string{surveycellscan.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := scsu.location; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveycellscan.LocationTable,
			Columns: []string{surveycellscan.LocationColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, scsu.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// SurveyCellScanUpdateOne is the builder for updating a single SurveyCellScan entity.
type SurveyCellScanUpdateOne struct {
	config
	id string

	update_time                  *time.Time
	network_type                 *string
	signal_strength              *int
	addsignal_strength           *int
	timestamp                    *time.Time
	cleartimestamp               bool
	base_station_id              *string
	clearbase_station_id         bool
	network_id                   *string
	clearnetwork_id              bool
	system_id                    *string
	clearsystem_id               bool
	cell_id                      *string
	clearcell_id                 bool
	location_area_code           *string
	clearlocation_area_code      bool
	mobile_country_code          *string
	clearmobile_country_code     bool
	mobile_network_code          *string
	clearmobile_network_code     bool
	primary_scrambling_code      *string
	clearprimary_scrambling_code bool
	operator                     *string
	clearoperator                bool
	arfcn                        *int
	addarfcn                     *int
	cleararfcn                   bool
	physical_cell_id             *string
	clearphysical_cell_id        bool
	tracking_area_code           *string
	cleartracking_area_code      bool
	timing_advance               *int
	addtiming_advance            *int
	cleartiming_advance          bool
	earfcn                       *int
	addearfcn                    *int
	clearearfcn                  bool
	uarfcn                       *int
	adduarfcn                    *int
	clearuarfcn                  bool
	latitude                     *float64
	addlatitude                  *float64
	clearlatitude                bool
	longitude                    *float64
	addlongitude                 *float64
	clearlongitude               bool
	survey_question              map[string]struct{}
	location                     map[string]struct{}
	clearedSurveyQuestion        bool
	clearedLocation              bool
}

// SetNetworkType sets the network_type field.
func (scsuo *SurveyCellScanUpdateOne) SetNetworkType(s string) *SurveyCellScanUpdateOne {
	scsuo.network_type = &s
	return scsuo
}

// SetSignalStrength sets the signal_strength field.
func (scsuo *SurveyCellScanUpdateOne) SetSignalStrength(i int) *SurveyCellScanUpdateOne {
	scsuo.signal_strength = &i
	scsuo.addsignal_strength = nil
	return scsuo
}

// AddSignalStrength adds i to signal_strength.
func (scsuo *SurveyCellScanUpdateOne) AddSignalStrength(i int) *SurveyCellScanUpdateOne {
	if scsuo.addsignal_strength == nil {
		scsuo.addsignal_strength = &i
	} else {
		*scsuo.addsignal_strength += i
	}
	return scsuo
}

// SetTimestamp sets the timestamp field.
func (scsuo *SurveyCellScanUpdateOne) SetTimestamp(t time.Time) *SurveyCellScanUpdateOne {
	scsuo.timestamp = &t
	return scsuo
}

// SetNillableTimestamp sets the timestamp field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableTimestamp(t *time.Time) *SurveyCellScanUpdateOne {
	if t != nil {
		scsuo.SetTimestamp(*t)
	}
	return scsuo
}

// ClearTimestamp clears the value of timestamp.
func (scsuo *SurveyCellScanUpdateOne) ClearTimestamp() *SurveyCellScanUpdateOne {
	scsuo.timestamp = nil
	scsuo.cleartimestamp = true
	return scsuo
}

// SetBaseStationID sets the base_station_id field.
func (scsuo *SurveyCellScanUpdateOne) SetBaseStationID(s string) *SurveyCellScanUpdateOne {
	scsuo.base_station_id = &s
	return scsuo
}

// SetNillableBaseStationID sets the base_station_id field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableBaseStationID(s *string) *SurveyCellScanUpdateOne {
	if s != nil {
		scsuo.SetBaseStationID(*s)
	}
	return scsuo
}

// ClearBaseStationID clears the value of base_station_id.
func (scsuo *SurveyCellScanUpdateOne) ClearBaseStationID() *SurveyCellScanUpdateOne {
	scsuo.base_station_id = nil
	scsuo.clearbase_station_id = true
	return scsuo
}

// SetNetworkID sets the network_id field.
func (scsuo *SurveyCellScanUpdateOne) SetNetworkID(s string) *SurveyCellScanUpdateOne {
	scsuo.network_id = &s
	return scsuo
}

// SetNillableNetworkID sets the network_id field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableNetworkID(s *string) *SurveyCellScanUpdateOne {
	if s != nil {
		scsuo.SetNetworkID(*s)
	}
	return scsuo
}

// ClearNetworkID clears the value of network_id.
func (scsuo *SurveyCellScanUpdateOne) ClearNetworkID() *SurveyCellScanUpdateOne {
	scsuo.network_id = nil
	scsuo.clearnetwork_id = true
	return scsuo
}

// SetSystemID sets the system_id field.
func (scsuo *SurveyCellScanUpdateOne) SetSystemID(s string) *SurveyCellScanUpdateOne {
	scsuo.system_id = &s
	return scsuo
}

// SetNillableSystemID sets the system_id field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableSystemID(s *string) *SurveyCellScanUpdateOne {
	if s != nil {
		scsuo.SetSystemID(*s)
	}
	return scsuo
}

// ClearSystemID clears the value of system_id.
func (scsuo *SurveyCellScanUpdateOne) ClearSystemID() *SurveyCellScanUpdateOne {
	scsuo.system_id = nil
	scsuo.clearsystem_id = true
	return scsuo
}

// SetCellID sets the cell_id field.
func (scsuo *SurveyCellScanUpdateOne) SetCellID(s string) *SurveyCellScanUpdateOne {
	scsuo.cell_id = &s
	return scsuo
}

// SetNillableCellID sets the cell_id field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableCellID(s *string) *SurveyCellScanUpdateOne {
	if s != nil {
		scsuo.SetCellID(*s)
	}
	return scsuo
}

// ClearCellID clears the value of cell_id.
func (scsuo *SurveyCellScanUpdateOne) ClearCellID() *SurveyCellScanUpdateOne {
	scsuo.cell_id = nil
	scsuo.clearcell_id = true
	return scsuo
}

// SetLocationAreaCode sets the location_area_code field.
func (scsuo *SurveyCellScanUpdateOne) SetLocationAreaCode(s string) *SurveyCellScanUpdateOne {
	scsuo.location_area_code = &s
	return scsuo
}

// SetNillableLocationAreaCode sets the location_area_code field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableLocationAreaCode(s *string) *SurveyCellScanUpdateOne {
	if s != nil {
		scsuo.SetLocationAreaCode(*s)
	}
	return scsuo
}

// ClearLocationAreaCode clears the value of location_area_code.
func (scsuo *SurveyCellScanUpdateOne) ClearLocationAreaCode() *SurveyCellScanUpdateOne {
	scsuo.location_area_code = nil
	scsuo.clearlocation_area_code = true
	return scsuo
}

// SetMobileCountryCode sets the mobile_country_code field.
func (scsuo *SurveyCellScanUpdateOne) SetMobileCountryCode(s string) *SurveyCellScanUpdateOne {
	scsuo.mobile_country_code = &s
	return scsuo
}

// SetNillableMobileCountryCode sets the mobile_country_code field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableMobileCountryCode(s *string) *SurveyCellScanUpdateOne {
	if s != nil {
		scsuo.SetMobileCountryCode(*s)
	}
	return scsuo
}

// ClearMobileCountryCode clears the value of mobile_country_code.
func (scsuo *SurveyCellScanUpdateOne) ClearMobileCountryCode() *SurveyCellScanUpdateOne {
	scsuo.mobile_country_code = nil
	scsuo.clearmobile_country_code = true
	return scsuo
}

// SetMobileNetworkCode sets the mobile_network_code field.
func (scsuo *SurveyCellScanUpdateOne) SetMobileNetworkCode(s string) *SurveyCellScanUpdateOne {
	scsuo.mobile_network_code = &s
	return scsuo
}

// SetNillableMobileNetworkCode sets the mobile_network_code field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableMobileNetworkCode(s *string) *SurveyCellScanUpdateOne {
	if s != nil {
		scsuo.SetMobileNetworkCode(*s)
	}
	return scsuo
}

// ClearMobileNetworkCode clears the value of mobile_network_code.
func (scsuo *SurveyCellScanUpdateOne) ClearMobileNetworkCode() *SurveyCellScanUpdateOne {
	scsuo.mobile_network_code = nil
	scsuo.clearmobile_network_code = true
	return scsuo
}

// SetPrimaryScramblingCode sets the primary_scrambling_code field.
func (scsuo *SurveyCellScanUpdateOne) SetPrimaryScramblingCode(s string) *SurveyCellScanUpdateOne {
	scsuo.primary_scrambling_code = &s
	return scsuo
}

// SetNillablePrimaryScramblingCode sets the primary_scrambling_code field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillablePrimaryScramblingCode(s *string) *SurveyCellScanUpdateOne {
	if s != nil {
		scsuo.SetPrimaryScramblingCode(*s)
	}
	return scsuo
}

// ClearPrimaryScramblingCode clears the value of primary_scrambling_code.
func (scsuo *SurveyCellScanUpdateOne) ClearPrimaryScramblingCode() *SurveyCellScanUpdateOne {
	scsuo.primary_scrambling_code = nil
	scsuo.clearprimary_scrambling_code = true
	return scsuo
}

// SetOperator sets the operator field.
func (scsuo *SurveyCellScanUpdateOne) SetOperator(s string) *SurveyCellScanUpdateOne {
	scsuo.operator = &s
	return scsuo
}

// SetNillableOperator sets the operator field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableOperator(s *string) *SurveyCellScanUpdateOne {
	if s != nil {
		scsuo.SetOperator(*s)
	}
	return scsuo
}

// ClearOperator clears the value of operator.
func (scsuo *SurveyCellScanUpdateOne) ClearOperator() *SurveyCellScanUpdateOne {
	scsuo.operator = nil
	scsuo.clearoperator = true
	return scsuo
}

// SetArfcn sets the arfcn field.
func (scsuo *SurveyCellScanUpdateOne) SetArfcn(i int) *SurveyCellScanUpdateOne {
	scsuo.arfcn = &i
	scsuo.addarfcn = nil
	return scsuo
}

// SetNillableArfcn sets the arfcn field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableArfcn(i *int) *SurveyCellScanUpdateOne {
	if i != nil {
		scsuo.SetArfcn(*i)
	}
	return scsuo
}

// AddArfcn adds i to arfcn.
func (scsuo *SurveyCellScanUpdateOne) AddArfcn(i int) *SurveyCellScanUpdateOne {
	if scsuo.addarfcn == nil {
		scsuo.addarfcn = &i
	} else {
		*scsuo.addarfcn += i
	}
	return scsuo
}

// ClearArfcn clears the value of arfcn.
func (scsuo *SurveyCellScanUpdateOne) ClearArfcn() *SurveyCellScanUpdateOne {
	scsuo.arfcn = nil
	scsuo.cleararfcn = true
	return scsuo
}

// SetPhysicalCellID sets the physical_cell_id field.
func (scsuo *SurveyCellScanUpdateOne) SetPhysicalCellID(s string) *SurveyCellScanUpdateOne {
	scsuo.physical_cell_id = &s
	return scsuo
}

// SetNillablePhysicalCellID sets the physical_cell_id field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillablePhysicalCellID(s *string) *SurveyCellScanUpdateOne {
	if s != nil {
		scsuo.SetPhysicalCellID(*s)
	}
	return scsuo
}

// ClearPhysicalCellID clears the value of physical_cell_id.
func (scsuo *SurveyCellScanUpdateOne) ClearPhysicalCellID() *SurveyCellScanUpdateOne {
	scsuo.physical_cell_id = nil
	scsuo.clearphysical_cell_id = true
	return scsuo
}

// SetTrackingAreaCode sets the tracking_area_code field.
func (scsuo *SurveyCellScanUpdateOne) SetTrackingAreaCode(s string) *SurveyCellScanUpdateOne {
	scsuo.tracking_area_code = &s
	return scsuo
}

// SetNillableTrackingAreaCode sets the tracking_area_code field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableTrackingAreaCode(s *string) *SurveyCellScanUpdateOne {
	if s != nil {
		scsuo.SetTrackingAreaCode(*s)
	}
	return scsuo
}

// ClearTrackingAreaCode clears the value of tracking_area_code.
func (scsuo *SurveyCellScanUpdateOne) ClearTrackingAreaCode() *SurveyCellScanUpdateOne {
	scsuo.tracking_area_code = nil
	scsuo.cleartracking_area_code = true
	return scsuo
}

// SetTimingAdvance sets the timing_advance field.
func (scsuo *SurveyCellScanUpdateOne) SetTimingAdvance(i int) *SurveyCellScanUpdateOne {
	scsuo.timing_advance = &i
	scsuo.addtiming_advance = nil
	return scsuo
}

// SetNillableTimingAdvance sets the timing_advance field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableTimingAdvance(i *int) *SurveyCellScanUpdateOne {
	if i != nil {
		scsuo.SetTimingAdvance(*i)
	}
	return scsuo
}

// AddTimingAdvance adds i to timing_advance.
func (scsuo *SurveyCellScanUpdateOne) AddTimingAdvance(i int) *SurveyCellScanUpdateOne {
	if scsuo.addtiming_advance == nil {
		scsuo.addtiming_advance = &i
	} else {
		*scsuo.addtiming_advance += i
	}
	return scsuo
}

// ClearTimingAdvance clears the value of timing_advance.
func (scsuo *SurveyCellScanUpdateOne) ClearTimingAdvance() *SurveyCellScanUpdateOne {
	scsuo.timing_advance = nil
	scsuo.cleartiming_advance = true
	return scsuo
}

// SetEarfcn sets the earfcn field.
func (scsuo *SurveyCellScanUpdateOne) SetEarfcn(i int) *SurveyCellScanUpdateOne {
	scsuo.earfcn = &i
	scsuo.addearfcn = nil
	return scsuo
}

// SetNillableEarfcn sets the earfcn field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableEarfcn(i *int) *SurveyCellScanUpdateOne {
	if i != nil {
		scsuo.SetEarfcn(*i)
	}
	return scsuo
}

// AddEarfcn adds i to earfcn.
func (scsuo *SurveyCellScanUpdateOne) AddEarfcn(i int) *SurveyCellScanUpdateOne {
	if scsuo.addearfcn == nil {
		scsuo.addearfcn = &i
	} else {
		*scsuo.addearfcn += i
	}
	return scsuo
}

// ClearEarfcn clears the value of earfcn.
func (scsuo *SurveyCellScanUpdateOne) ClearEarfcn() *SurveyCellScanUpdateOne {
	scsuo.earfcn = nil
	scsuo.clearearfcn = true
	return scsuo
}

// SetUarfcn sets the uarfcn field.
func (scsuo *SurveyCellScanUpdateOne) SetUarfcn(i int) *SurveyCellScanUpdateOne {
	scsuo.uarfcn = &i
	scsuo.adduarfcn = nil
	return scsuo
}

// SetNillableUarfcn sets the uarfcn field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableUarfcn(i *int) *SurveyCellScanUpdateOne {
	if i != nil {
		scsuo.SetUarfcn(*i)
	}
	return scsuo
}

// AddUarfcn adds i to uarfcn.
func (scsuo *SurveyCellScanUpdateOne) AddUarfcn(i int) *SurveyCellScanUpdateOne {
	if scsuo.adduarfcn == nil {
		scsuo.adduarfcn = &i
	} else {
		*scsuo.adduarfcn += i
	}
	return scsuo
}

// ClearUarfcn clears the value of uarfcn.
func (scsuo *SurveyCellScanUpdateOne) ClearUarfcn() *SurveyCellScanUpdateOne {
	scsuo.uarfcn = nil
	scsuo.clearuarfcn = true
	return scsuo
}

// SetLatitude sets the latitude field.
func (scsuo *SurveyCellScanUpdateOne) SetLatitude(f float64) *SurveyCellScanUpdateOne {
	scsuo.latitude = &f
	scsuo.addlatitude = nil
	return scsuo
}

// SetNillableLatitude sets the latitude field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableLatitude(f *float64) *SurveyCellScanUpdateOne {
	if f != nil {
		scsuo.SetLatitude(*f)
	}
	return scsuo
}

// AddLatitude adds f to latitude.
func (scsuo *SurveyCellScanUpdateOne) AddLatitude(f float64) *SurveyCellScanUpdateOne {
	if scsuo.addlatitude == nil {
		scsuo.addlatitude = &f
	} else {
		*scsuo.addlatitude += f
	}
	return scsuo
}

// ClearLatitude clears the value of latitude.
func (scsuo *SurveyCellScanUpdateOne) ClearLatitude() *SurveyCellScanUpdateOne {
	scsuo.latitude = nil
	scsuo.clearlatitude = true
	return scsuo
}

// SetLongitude sets the longitude field.
func (scsuo *SurveyCellScanUpdateOne) SetLongitude(f float64) *SurveyCellScanUpdateOne {
	scsuo.longitude = &f
	scsuo.addlongitude = nil
	return scsuo
}

// SetNillableLongitude sets the longitude field if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableLongitude(f *float64) *SurveyCellScanUpdateOne {
	if f != nil {
		scsuo.SetLongitude(*f)
	}
	return scsuo
}

// AddLongitude adds f to longitude.
func (scsuo *SurveyCellScanUpdateOne) AddLongitude(f float64) *SurveyCellScanUpdateOne {
	if scsuo.addlongitude == nil {
		scsuo.addlongitude = &f
	} else {
		*scsuo.addlongitude += f
	}
	return scsuo
}

// ClearLongitude clears the value of longitude.
func (scsuo *SurveyCellScanUpdateOne) ClearLongitude() *SurveyCellScanUpdateOne {
	scsuo.longitude = nil
	scsuo.clearlongitude = true
	return scsuo
}

// SetSurveyQuestionID sets the survey_question edge to SurveyQuestion by id.
func (scsuo *SurveyCellScanUpdateOne) SetSurveyQuestionID(id string) *SurveyCellScanUpdateOne {
	if scsuo.survey_question == nil {
		scsuo.survey_question = make(map[string]struct{})
	}
	scsuo.survey_question[id] = struct{}{}
	return scsuo
}

// SetNillableSurveyQuestionID sets the survey_question edge to SurveyQuestion by id if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableSurveyQuestionID(id *string) *SurveyCellScanUpdateOne {
	if id != nil {
		scsuo = scsuo.SetSurveyQuestionID(*id)
	}
	return scsuo
}

// SetSurveyQuestion sets the survey_question edge to SurveyQuestion.
func (scsuo *SurveyCellScanUpdateOne) SetSurveyQuestion(s *SurveyQuestion) *SurveyCellScanUpdateOne {
	return scsuo.SetSurveyQuestionID(s.ID)
}

// SetLocationID sets the location edge to Location by id.
func (scsuo *SurveyCellScanUpdateOne) SetLocationID(id string) *SurveyCellScanUpdateOne {
	if scsuo.location == nil {
		scsuo.location = make(map[string]struct{})
	}
	scsuo.location[id] = struct{}{}
	return scsuo
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableLocationID(id *string) *SurveyCellScanUpdateOne {
	if id != nil {
		scsuo = scsuo.SetLocationID(*id)
	}
	return scsuo
}

// SetLocation sets the location edge to Location.
func (scsuo *SurveyCellScanUpdateOne) SetLocation(l *Location) *SurveyCellScanUpdateOne {
	return scsuo.SetLocationID(l.ID)
}

// ClearSurveyQuestion clears the survey_question edge to SurveyQuestion.
func (scsuo *SurveyCellScanUpdateOne) ClearSurveyQuestion() *SurveyCellScanUpdateOne {
	scsuo.clearedSurveyQuestion = true
	return scsuo
}

// ClearLocation clears the location edge to Location.
func (scsuo *SurveyCellScanUpdateOne) ClearLocation() *SurveyCellScanUpdateOne {
	scsuo.clearedLocation = true
	return scsuo
}

// Save executes the query and returns the updated entity.
func (scsuo *SurveyCellScanUpdateOne) Save(ctx context.Context) (*SurveyCellScan, error) {
	if scsuo.update_time == nil {
		v := surveycellscan.UpdateDefaultUpdateTime()
		scsuo.update_time = &v
	}
	if len(scsuo.survey_question) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"survey_question\"")
	}
	if len(scsuo.location) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	return scsuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (scsuo *SurveyCellScanUpdateOne) SaveX(ctx context.Context) *SurveyCellScan {
	scs, err := scsuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return scs
}

// Exec executes the query on the entity.
func (scsuo *SurveyCellScanUpdateOne) Exec(ctx context.Context) error {
	_, err := scsuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scsuo *SurveyCellScanUpdateOne) ExecX(ctx context.Context) {
	if err := scsuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (scsuo *SurveyCellScanUpdateOne) sqlSave(ctx context.Context) (scs *SurveyCellScan, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveycellscan.Table,
			Columns: surveycellscan.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  scsuo.id,
				Type:   field.TypeString,
				Column: surveycellscan.FieldID,
			},
		},
	}
	if value := scsuo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveycellscan.FieldUpdateTime,
		})
	}
	if value := scsuo.network_type; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldNetworkType,
		})
	}
	if value := scsuo.signal_strength; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldSignalStrength,
		})
	}
	if value := scsuo.addsignal_strength; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldSignalStrength,
		})
	}
	if value := scsuo.timestamp; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveycellscan.FieldTimestamp,
		})
	}
	if scsuo.cleartimestamp {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: surveycellscan.FieldTimestamp,
		})
	}
	if value := scsuo.base_station_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldBaseStationID,
		})
	}
	if scsuo.clearbase_station_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldBaseStationID,
		})
	}
	if value := scsuo.network_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldNetworkID,
		})
	}
	if scsuo.clearnetwork_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldNetworkID,
		})
	}
	if value := scsuo.system_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldSystemID,
		})
	}
	if scsuo.clearsystem_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldSystemID,
		})
	}
	if value := scsuo.cell_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldCellID,
		})
	}
	if scsuo.clearcell_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldCellID,
		})
	}
	if value := scsuo.location_area_code; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldLocationAreaCode,
		})
	}
	if scsuo.clearlocation_area_code {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldLocationAreaCode,
		})
	}
	if value := scsuo.mobile_country_code; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldMobileCountryCode,
		})
	}
	if scsuo.clearmobile_country_code {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldMobileCountryCode,
		})
	}
	if value := scsuo.mobile_network_code; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldMobileNetworkCode,
		})
	}
	if scsuo.clearmobile_network_code {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldMobileNetworkCode,
		})
	}
	if value := scsuo.primary_scrambling_code; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldPrimaryScramblingCode,
		})
	}
	if scsuo.clearprimary_scrambling_code {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldPrimaryScramblingCode,
		})
	}
	if value := scsuo.operator; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldOperator,
		})
	}
	if scsuo.clearoperator {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldOperator,
		})
	}
	if value := scsuo.arfcn; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldArfcn,
		})
	}
	if value := scsuo.addarfcn; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldArfcn,
		})
	}
	if scsuo.cleararfcn {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveycellscan.FieldArfcn,
		})
	}
	if value := scsuo.physical_cell_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldPhysicalCellID,
		})
	}
	if scsuo.clearphysical_cell_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldPhysicalCellID,
		})
	}
	if value := scsuo.tracking_area_code; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldTrackingAreaCode,
		})
	}
	if scsuo.cleartracking_area_code {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldTrackingAreaCode,
		})
	}
	if value := scsuo.timing_advance; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldTimingAdvance,
		})
	}
	if value := scsuo.addtiming_advance; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldTimingAdvance,
		})
	}
	if scsuo.cleartiming_advance {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveycellscan.FieldTimingAdvance,
		})
	}
	if value := scsuo.earfcn; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldEarfcn,
		})
	}
	if value := scsuo.addearfcn; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldEarfcn,
		})
	}
	if scsuo.clearearfcn {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveycellscan.FieldEarfcn,
		})
	}
	if value := scsuo.uarfcn; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldUarfcn,
		})
	}
	if value := scsuo.adduarfcn; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldUarfcn,
		})
	}
	if scsuo.clearuarfcn {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveycellscan.FieldUarfcn,
		})
	}
	if value := scsuo.latitude; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveycellscan.FieldLatitude,
		})
	}
	if value := scsuo.addlatitude; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveycellscan.FieldLatitude,
		})
	}
	if scsuo.clearlatitude {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveycellscan.FieldLatitude,
		})
	}
	if value := scsuo.longitude; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveycellscan.FieldLongitude,
		})
	}
	if value := scsuo.addlongitude; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveycellscan.FieldLongitude,
		})
	}
	if scsuo.clearlongitude {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveycellscan.FieldLongitude,
		})
	}
	if scsuo.clearedSurveyQuestion {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveycellscan.SurveyQuestionTable,
			Columns: []string{surveycellscan.SurveyQuestionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveyquestion.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := scsuo.survey_question; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveycellscan.SurveyQuestionTable,
			Columns: []string{surveycellscan.SurveyQuestionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveyquestion.FieldID,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if scsuo.clearedLocation {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveycellscan.LocationTable,
			Columns: []string{surveycellscan.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := scsuo.location; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveycellscan.LocationTable,
			Columns: []string{surveycellscan.LocationColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	scs = &SurveyCellScan{config: scsuo.config}
	_spec.Assign = scs.assignValues
	_spec.ScanValues = scs.scanValues()
	if err = sqlgraph.UpdateNode(ctx, scsuo.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return scs, nil
}
