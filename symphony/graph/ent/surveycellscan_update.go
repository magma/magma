// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
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
	hooks      []Hook
	mutation   *SurveyCellScanMutation
	predicates []predicate.SurveyCellScan
}

// Where adds a new predicate for the builder.
func (scsu *SurveyCellScanUpdate) Where(ps ...predicate.SurveyCellScan) *SurveyCellScanUpdate {
	scsu.predicates = append(scsu.predicates, ps...)
	return scsu
}

// SetNetworkType sets the network_type field.
func (scsu *SurveyCellScanUpdate) SetNetworkType(s string) *SurveyCellScanUpdate {
	scsu.mutation.SetNetworkType(s)
	return scsu
}

// SetSignalStrength sets the signal_strength field.
func (scsu *SurveyCellScanUpdate) SetSignalStrength(i int) *SurveyCellScanUpdate {
	scsu.mutation.ResetSignalStrength()
	scsu.mutation.SetSignalStrength(i)
	return scsu
}

// AddSignalStrength adds i to signal_strength.
func (scsu *SurveyCellScanUpdate) AddSignalStrength(i int) *SurveyCellScanUpdate {
	scsu.mutation.AddSignalStrength(i)
	return scsu
}

// SetTimestamp sets the timestamp field.
func (scsu *SurveyCellScanUpdate) SetTimestamp(t time.Time) *SurveyCellScanUpdate {
	scsu.mutation.SetTimestamp(t)
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
	scsu.mutation.ClearTimestamp()
	return scsu
}

// SetBaseStationID sets the base_station_id field.
func (scsu *SurveyCellScanUpdate) SetBaseStationID(s string) *SurveyCellScanUpdate {
	scsu.mutation.SetBaseStationID(s)
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
	scsu.mutation.ClearBaseStationID()
	return scsu
}

// SetNetworkID sets the network_id field.
func (scsu *SurveyCellScanUpdate) SetNetworkID(s string) *SurveyCellScanUpdate {
	scsu.mutation.SetNetworkID(s)
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
	scsu.mutation.ClearNetworkID()
	return scsu
}

// SetSystemID sets the system_id field.
func (scsu *SurveyCellScanUpdate) SetSystemID(s string) *SurveyCellScanUpdate {
	scsu.mutation.SetSystemID(s)
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
	scsu.mutation.ClearSystemID()
	return scsu
}

// SetCellID sets the cell_id field.
func (scsu *SurveyCellScanUpdate) SetCellID(s string) *SurveyCellScanUpdate {
	scsu.mutation.SetCellID(s)
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
	scsu.mutation.ClearCellID()
	return scsu
}

// SetLocationAreaCode sets the location_area_code field.
func (scsu *SurveyCellScanUpdate) SetLocationAreaCode(s string) *SurveyCellScanUpdate {
	scsu.mutation.SetLocationAreaCode(s)
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
	scsu.mutation.ClearLocationAreaCode()
	return scsu
}

// SetMobileCountryCode sets the mobile_country_code field.
func (scsu *SurveyCellScanUpdate) SetMobileCountryCode(s string) *SurveyCellScanUpdate {
	scsu.mutation.SetMobileCountryCode(s)
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
	scsu.mutation.ClearMobileCountryCode()
	return scsu
}

// SetMobileNetworkCode sets the mobile_network_code field.
func (scsu *SurveyCellScanUpdate) SetMobileNetworkCode(s string) *SurveyCellScanUpdate {
	scsu.mutation.SetMobileNetworkCode(s)
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
	scsu.mutation.ClearMobileNetworkCode()
	return scsu
}

// SetPrimaryScramblingCode sets the primary_scrambling_code field.
func (scsu *SurveyCellScanUpdate) SetPrimaryScramblingCode(s string) *SurveyCellScanUpdate {
	scsu.mutation.SetPrimaryScramblingCode(s)
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
	scsu.mutation.ClearPrimaryScramblingCode()
	return scsu
}

// SetOperator sets the operator field.
func (scsu *SurveyCellScanUpdate) SetOperator(s string) *SurveyCellScanUpdate {
	scsu.mutation.SetOperator(s)
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
	scsu.mutation.ClearOperator()
	return scsu
}

// SetArfcn sets the arfcn field.
func (scsu *SurveyCellScanUpdate) SetArfcn(i int) *SurveyCellScanUpdate {
	scsu.mutation.ResetArfcn()
	scsu.mutation.SetArfcn(i)
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
	scsu.mutation.AddArfcn(i)
	return scsu
}

// ClearArfcn clears the value of arfcn.
func (scsu *SurveyCellScanUpdate) ClearArfcn() *SurveyCellScanUpdate {
	scsu.mutation.ClearArfcn()
	return scsu
}

// SetPhysicalCellID sets the physical_cell_id field.
func (scsu *SurveyCellScanUpdate) SetPhysicalCellID(s string) *SurveyCellScanUpdate {
	scsu.mutation.SetPhysicalCellID(s)
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
	scsu.mutation.ClearPhysicalCellID()
	return scsu
}

// SetTrackingAreaCode sets the tracking_area_code field.
func (scsu *SurveyCellScanUpdate) SetTrackingAreaCode(s string) *SurveyCellScanUpdate {
	scsu.mutation.SetTrackingAreaCode(s)
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
	scsu.mutation.ClearTrackingAreaCode()
	return scsu
}

// SetTimingAdvance sets the timing_advance field.
func (scsu *SurveyCellScanUpdate) SetTimingAdvance(i int) *SurveyCellScanUpdate {
	scsu.mutation.ResetTimingAdvance()
	scsu.mutation.SetTimingAdvance(i)
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
	scsu.mutation.AddTimingAdvance(i)
	return scsu
}

// ClearTimingAdvance clears the value of timing_advance.
func (scsu *SurveyCellScanUpdate) ClearTimingAdvance() *SurveyCellScanUpdate {
	scsu.mutation.ClearTimingAdvance()
	return scsu
}

// SetEarfcn sets the earfcn field.
func (scsu *SurveyCellScanUpdate) SetEarfcn(i int) *SurveyCellScanUpdate {
	scsu.mutation.ResetEarfcn()
	scsu.mutation.SetEarfcn(i)
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
	scsu.mutation.AddEarfcn(i)
	return scsu
}

// ClearEarfcn clears the value of earfcn.
func (scsu *SurveyCellScanUpdate) ClearEarfcn() *SurveyCellScanUpdate {
	scsu.mutation.ClearEarfcn()
	return scsu
}

// SetUarfcn sets the uarfcn field.
func (scsu *SurveyCellScanUpdate) SetUarfcn(i int) *SurveyCellScanUpdate {
	scsu.mutation.ResetUarfcn()
	scsu.mutation.SetUarfcn(i)
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
	scsu.mutation.AddUarfcn(i)
	return scsu
}

// ClearUarfcn clears the value of uarfcn.
func (scsu *SurveyCellScanUpdate) ClearUarfcn() *SurveyCellScanUpdate {
	scsu.mutation.ClearUarfcn()
	return scsu
}

// SetLatitude sets the latitude field.
func (scsu *SurveyCellScanUpdate) SetLatitude(f float64) *SurveyCellScanUpdate {
	scsu.mutation.ResetLatitude()
	scsu.mutation.SetLatitude(f)
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
	scsu.mutation.AddLatitude(f)
	return scsu
}

// ClearLatitude clears the value of latitude.
func (scsu *SurveyCellScanUpdate) ClearLatitude() *SurveyCellScanUpdate {
	scsu.mutation.ClearLatitude()
	return scsu
}

// SetLongitude sets the longitude field.
func (scsu *SurveyCellScanUpdate) SetLongitude(f float64) *SurveyCellScanUpdate {
	scsu.mutation.ResetLongitude()
	scsu.mutation.SetLongitude(f)
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
	scsu.mutation.AddLongitude(f)
	return scsu
}

// ClearLongitude clears the value of longitude.
func (scsu *SurveyCellScanUpdate) ClearLongitude() *SurveyCellScanUpdate {
	scsu.mutation.ClearLongitude()
	return scsu
}

// SetSurveyQuestionID sets the survey_question edge to SurveyQuestion by id.
func (scsu *SurveyCellScanUpdate) SetSurveyQuestionID(id int) *SurveyCellScanUpdate {
	scsu.mutation.SetSurveyQuestionID(id)
	return scsu
}

// SetNillableSurveyQuestionID sets the survey_question edge to SurveyQuestion by id if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableSurveyQuestionID(id *int) *SurveyCellScanUpdate {
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
func (scsu *SurveyCellScanUpdate) SetLocationID(id int) *SurveyCellScanUpdate {
	scsu.mutation.SetLocationID(id)
	return scsu
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (scsu *SurveyCellScanUpdate) SetNillableLocationID(id *int) *SurveyCellScanUpdate {
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
	scsu.mutation.ClearSurveyQuestion()
	return scsu
}

// ClearLocation clears the location edge to Location.
func (scsu *SurveyCellScanUpdate) ClearLocation() *SurveyCellScanUpdate {
	scsu.mutation.ClearLocation()
	return scsu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (scsu *SurveyCellScanUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := scsu.mutation.UpdateTime(); !ok {
		v := surveycellscan.UpdateDefaultUpdateTime()
		scsu.mutation.SetUpdateTime(v)
	}

	var (
		err      error
		affected int
	)
	if len(scsu.hooks) == 0 {
		affected, err = scsu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyCellScanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			scsu.mutation = mutation
			affected, err = scsu.sqlSave(ctx)
			return affected, err
		})
		for i := len(scsu.hooks) - 1; i >= 0; i-- {
			mut = scsu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, scsu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
				Type:   field.TypeInt,
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
	if value, ok := scsu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveycellscan.FieldUpdateTime,
		})
	}
	if value, ok := scsu.mutation.NetworkType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldNetworkType,
		})
	}
	if value, ok := scsu.mutation.SignalStrength(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldSignalStrength,
		})
	}
	if value, ok := scsu.mutation.AddedSignalStrength(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldSignalStrength,
		})
	}
	if value, ok := scsu.mutation.Timestamp(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveycellscan.FieldTimestamp,
		})
	}
	if scsu.mutation.TimestampCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: surveycellscan.FieldTimestamp,
		})
	}
	if value, ok := scsu.mutation.BaseStationID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldBaseStationID,
		})
	}
	if scsu.mutation.BaseStationIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldBaseStationID,
		})
	}
	if value, ok := scsu.mutation.NetworkID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldNetworkID,
		})
	}
	if scsu.mutation.NetworkIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldNetworkID,
		})
	}
	if value, ok := scsu.mutation.SystemID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldSystemID,
		})
	}
	if scsu.mutation.SystemIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldSystemID,
		})
	}
	if value, ok := scsu.mutation.CellID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldCellID,
		})
	}
	if scsu.mutation.CellIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldCellID,
		})
	}
	if value, ok := scsu.mutation.LocationAreaCode(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldLocationAreaCode,
		})
	}
	if scsu.mutation.LocationAreaCodeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldLocationAreaCode,
		})
	}
	if value, ok := scsu.mutation.MobileCountryCode(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldMobileCountryCode,
		})
	}
	if scsu.mutation.MobileCountryCodeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldMobileCountryCode,
		})
	}
	if value, ok := scsu.mutation.MobileNetworkCode(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldMobileNetworkCode,
		})
	}
	if scsu.mutation.MobileNetworkCodeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldMobileNetworkCode,
		})
	}
	if value, ok := scsu.mutation.PrimaryScramblingCode(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldPrimaryScramblingCode,
		})
	}
	if scsu.mutation.PrimaryScramblingCodeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldPrimaryScramblingCode,
		})
	}
	if value, ok := scsu.mutation.Operator(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldOperator,
		})
	}
	if scsu.mutation.OperatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldOperator,
		})
	}
	if value, ok := scsu.mutation.Arfcn(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldArfcn,
		})
	}
	if value, ok := scsu.mutation.AddedArfcn(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldArfcn,
		})
	}
	if scsu.mutation.ArfcnCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveycellscan.FieldArfcn,
		})
	}
	if value, ok := scsu.mutation.PhysicalCellID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldPhysicalCellID,
		})
	}
	if scsu.mutation.PhysicalCellIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldPhysicalCellID,
		})
	}
	if value, ok := scsu.mutation.TrackingAreaCode(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldTrackingAreaCode,
		})
	}
	if scsu.mutation.TrackingAreaCodeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldTrackingAreaCode,
		})
	}
	if value, ok := scsu.mutation.TimingAdvance(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldTimingAdvance,
		})
	}
	if value, ok := scsu.mutation.AddedTimingAdvance(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldTimingAdvance,
		})
	}
	if scsu.mutation.TimingAdvanceCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveycellscan.FieldTimingAdvance,
		})
	}
	if value, ok := scsu.mutation.Earfcn(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldEarfcn,
		})
	}
	if value, ok := scsu.mutation.AddedEarfcn(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldEarfcn,
		})
	}
	if scsu.mutation.EarfcnCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveycellscan.FieldEarfcn,
		})
	}
	if value, ok := scsu.mutation.Uarfcn(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldUarfcn,
		})
	}
	if value, ok := scsu.mutation.AddedUarfcn(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldUarfcn,
		})
	}
	if scsu.mutation.UarfcnCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveycellscan.FieldUarfcn,
		})
	}
	if value, ok := scsu.mutation.Latitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveycellscan.FieldLatitude,
		})
	}
	if value, ok := scsu.mutation.AddedLatitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveycellscan.FieldLatitude,
		})
	}
	if scsu.mutation.LatitudeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveycellscan.FieldLatitude,
		})
	}
	if value, ok := scsu.mutation.Longitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveycellscan.FieldLongitude,
		})
	}
	if value, ok := scsu.mutation.AddedLongitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveycellscan.FieldLongitude,
		})
	}
	if scsu.mutation.LongitudeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveycellscan.FieldLongitude,
		})
	}
	if scsu.mutation.SurveyQuestionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveycellscan.SurveyQuestionTable,
			Columns: []string{surveycellscan.SurveyQuestionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveyquestion.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := scsu.mutation.SurveyQuestionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveycellscan.SurveyQuestionTable,
			Columns: []string{surveycellscan.SurveyQuestionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveyquestion.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if scsu.mutation.LocationCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveycellscan.LocationTable,
			Columns: []string{surveycellscan.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := scsu.mutation.LocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveycellscan.LocationTable,
			Columns: []string{surveycellscan.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, scsu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{surveycellscan.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// SurveyCellScanUpdateOne is the builder for updating a single SurveyCellScan entity.
type SurveyCellScanUpdateOne struct {
	config
	hooks    []Hook
	mutation *SurveyCellScanMutation
}

// SetNetworkType sets the network_type field.
func (scsuo *SurveyCellScanUpdateOne) SetNetworkType(s string) *SurveyCellScanUpdateOne {
	scsuo.mutation.SetNetworkType(s)
	return scsuo
}

// SetSignalStrength sets the signal_strength field.
func (scsuo *SurveyCellScanUpdateOne) SetSignalStrength(i int) *SurveyCellScanUpdateOne {
	scsuo.mutation.ResetSignalStrength()
	scsuo.mutation.SetSignalStrength(i)
	return scsuo
}

// AddSignalStrength adds i to signal_strength.
func (scsuo *SurveyCellScanUpdateOne) AddSignalStrength(i int) *SurveyCellScanUpdateOne {
	scsuo.mutation.AddSignalStrength(i)
	return scsuo
}

// SetTimestamp sets the timestamp field.
func (scsuo *SurveyCellScanUpdateOne) SetTimestamp(t time.Time) *SurveyCellScanUpdateOne {
	scsuo.mutation.SetTimestamp(t)
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
	scsuo.mutation.ClearTimestamp()
	return scsuo
}

// SetBaseStationID sets the base_station_id field.
func (scsuo *SurveyCellScanUpdateOne) SetBaseStationID(s string) *SurveyCellScanUpdateOne {
	scsuo.mutation.SetBaseStationID(s)
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
	scsuo.mutation.ClearBaseStationID()
	return scsuo
}

// SetNetworkID sets the network_id field.
func (scsuo *SurveyCellScanUpdateOne) SetNetworkID(s string) *SurveyCellScanUpdateOne {
	scsuo.mutation.SetNetworkID(s)
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
	scsuo.mutation.ClearNetworkID()
	return scsuo
}

// SetSystemID sets the system_id field.
func (scsuo *SurveyCellScanUpdateOne) SetSystemID(s string) *SurveyCellScanUpdateOne {
	scsuo.mutation.SetSystemID(s)
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
	scsuo.mutation.ClearSystemID()
	return scsuo
}

// SetCellID sets the cell_id field.
func (scsuo *SurveyCellScanUpdateOne) SetCellID(s string) *SurveyCellScanUpdateOne {
	scsuo.mutation.SetCellID(s)
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
	scsuo.mutation.ClearCellID()
	return scsuo
}

// SetLocationAreaCode sets the location_area_code field.
func (scsuo *SurveyCellScanUpdateOne) SetLocationAreaCode(s string) *SurveyCellScanUpdateOne {
	scsuo.mutation.SetLocationAreaCode(s)
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
	scsuo.mutation.ClearLocationAreaCode()
	return scsuo
}

// SetMobileCountryCode sets the mobile_country_code field.
func (scsuo *SurveyCellScanUpdateOne) SetMobileCountryCode(s string) *SurveyCellScanUpdateOne {
	scsuo.mutation.SetMobileCountryCode(s)
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
	scsuo.mutation.ClearMobileCountryCode()
	return scsuo
}

// SetMobileNetworkCode sets the mobile_network_code field.
func (scsuo *SurveyCellScanUpdateOne) SetMobileNetworkCode(s string) *SurveyCellScanUpdateOne {
	scsuo.mutation.SetMobileNetworkCode(s)
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
	scsuo.mutation.ClearMobileNetworkCode()
	return scsuo
}

// SetPrimaryScramblingCode sets the primary_scrambling_code field.
func (scsuo *SurveyCellScanUpdateOne) SetPrimaryScramblingCode(s string) *SurveyCellScanUpdateOne {
	scsuo.mutation.SetPrimaryScramblingCode(s)
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
	scsuo.mutation.ClearPrimaryScramblingCode()
	return scsuo
}

// SetOperator sets the operator field.
func (scsuo *SurveyCellScanUpdateOne) SetOperator(s string) *SurveyCellScanUpdateOne {
	scsuo.mutation.SetOperator(s)
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
	scsuo.mutation.ClearOperator()
	return scsuo
}

// SetArfcn sets the arfcn field.
func (scsuo *SurveyCellScanUpdateOne) SetArfcn(i int) *SurveyCellScanUpdateOne {
	scsuo.mutation.ResetArfcn()
	scsuo.mutation.SetArfcn(i)
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
	scsuo.mutation.AddArfcn(i)
	return scsuo
}

// ClearArfcn clears the value of arfcn.
func (scsuo *SurveyCellScanUpdateOne) ClearArfcn() *SurveyCellScanUpdateOne {
	scsuo.mutation.ClearArfcn()
	return scsuo
}

// SetPhysicalCellID sets the physical_cell_id field.
func (scsuo *SurveyCellScanUpdateOne) SetPhysicalCellID(s string) *SurveyCellScanUpdateOne {
	scsuo.mutation.SetPhysicalCellID(s)
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
	scsuo.mutation.ClearPhysicalCellID()
	return scsuo
}

// SetTrackingAreaCode sets the tracking_area_code field.
func (scsuo *SurveyCellScanUpdateOne) SetTrackingAreaCode(s string) *SurveyCellScanUpdateOne {
	scsuo.mutation.SetTrackingAreaCode(s)
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
	scsuo.mutation.ClearTrackingAreaCode()
	return scsuo
}

// SetTimingAdvance sets the timing_advance field.
func (scsuo *SurveyCellScanUpdateOne) SetTimingAdvance(i int) *SurveyCellScanUpdateOne {
	scsuo.mutation.ResetTimingAdvance()
	scsuo.mutation.SetTimingAdvance(i)
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
	scsuo.mutation.AddTimingAdvance(i)
	return scsuo
}

// ClearTimingAdvance clears the value of timing_advance.
func (scsuo *SurveyCellScanUpdateOne) ClearTimingAdvance() *SurveyCellScanUpdateOne {
	scsuo.mutation.ClearTimingAdvance()
	return scsuo
}

// SetEarfcn sets the earfcn field.
func (scsuo *SurveyCellScanUpdateOne) SetEarfcn(i int) *SurveyCellScanUpdateOne {
	scsuo.mutation.ResetEarfcn()
	scsuo.mutation.SetEarfcn(i)
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
	scsuo.mutation.AddEarfcn(i)
	return scsuo
}

// ClearEarfcn clears the value of earfcn.
func (scsuo *SurveyCellScanUpdateOne) ClearEarfcn() *SurveyCellScanUpdateOne {
	scsuo.mutation.ClearEarfcn()
	return scsuo
}

// SetUarfcn sets the uarfcn field.
func (scsuo *SurveyCellScanUpdateOne) SetUarfcn(i int) *SurveyCellScanUpdateOne {
	scsuo.mutation.ResetUarfcn()
	scsuo.mutation.SetUarfcn(i)
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
	scsuo.mutation.AddUarfcn(i)
	return scsuo
}

// ClearUarfcn clears the value of uarfcn.
func (scsuo *SurveyCellScanUpdateOne) ClearUarfcn() *SurveyCellScanUpdateOne {
	scsuo.mutation.ClearUarfcn()
	return scsuo
}

// SetLatitude sets the latitude field.
func (scsuo *SurveyCellScanUpdateOne) SetLatitude(f float64) *SurveyCellScanUpdateOne {
	scsuo.mutation.ResetLatitude()
	scsuo.mutation.SetLatitude(f)
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
	scsuo.mutation.AddLatitude(f)
	return scsuo
}

// ClearLatitude clears the value of latitude.
func (scsuo *SurveyCellScanUpdateOne) ClearLatitude() *SurveyCellScanUpdateOne {
	scsuo.mutation.ClearLatitude()
	return scsuo
}

// SetLongitude sets the longitude field.
func (scsuo *SurveyCellScanUpdateOne) SetLongitude(f float64) *SurveyCellScanUpdateOne {
	scsuo.mutation.ResetLongitude()
	scsuo.mutation.SetLongitude(f)
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
	scsuo.mutation.AddLongitude(f)
	return scsuo
}

// ClearLongitude clears the value of longitude.
func (scsuo *SurveyCellScanUpdateOne) ClearLongitude() *SurveyCellScanUpdateOne {
	scsuo.mutation.ClearLongitude()
	return scsuo
}

// SetSurveyQuestionID sets the survey_question edge to SurveyQuestion by id.
func (scsuo *SurveyCellScanUpdateOne) SetSurveyQuestionID(id int) *SurveyCellScanUpdateOne {
	scsuo.mutation.SetSurveyQuestionID(id)
	return scsuo
}

// SetNillableSurveyQuestionID sets the survey_question edge to SurveyQuestion by id if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableSurveyQuestionID(id *int) *SurveyCellScanUpdateOne {
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
func (scsuo *SurveyCellScanUpdateOne) SetLocationID(id int) *SurveyCellScanUpdateOne {
	scsuo.mutation.SetLocationID(id)
	return scsuo
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (scsuo *SurveyCellScanUpdateOne) SetNillableLocationID(id *int) *SurveyCellScanUpdateOne {
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
	scsuo.mutation.ClearSurveyQuestion()
	return scsuo
}

// ClearLocation clears the location edge to Location.
func (scsuo *SurveyCellScanUpdateOne) ClearLocation() *SurveyCellScanUpdateOne {
	scsuo.mutation.ClearLocation()
	return scsuo
}

// Save executes the query and returns the updated entity.
func (scsuo *SurveyCellScanUpdateOne) Save(ctx context.Context) (*SurveyCellScan, error) {
	if _, ok := scsuo.mutation.UpdateTime(); !ok {
		v := surveycellscan.UpdateDefaultUpdateTime()
		scsuo.mutation.SetUpdateTime(v)
	}

	var (
		err  error
		node *SurveyCellScan
	)
	if len(scsuo.hooks) == 0 {
		node, err = scsuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyCellScanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			scsuo.mutation = mutation
			node, err = scsuo.sqlSave(ctx)
			return node, err
		})
		for i := len(scsuo.hooks) - 1; i >= 0; i-- {
			mut = scsuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, scsuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
				Type:   field.TypeInt,
				Column: surveycellscan.FieldID,
			},
		},
	}
	id, ok := scsuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing SurveyCellScan.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := scsuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveycellscan.FieldUpdateTime,
		})
	}
	if value, ok := scsuo.mutation.NetworkType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldNetworkType,
		})
	}
	if value, ok := scsuo.mutation.SignalStrength(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldSignalStrength,
		})
	}
	if value, ok := scsuo.mutation.AddedSignalStrength(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldSignalStrength,
		})
	}
	if value, ok := scsuo.mutation.Timestamp(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveycellscan.FieldTimestamp,
		})
	}
	if scsuo.mutation.TimestampCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: surveycellscan.FieldTimestamp,
		})
	}
	if value, ok := scsuo.mutation.BaseStationID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldBaseStationID,
		})
	}
	if scsuo.mutation.BaseStationIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldBaseStationID,
		})
	}
	if value, ok := scsuo.mutation.NetworkID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldNetworkID,
		})
	}
	if scsuo.mutation.NetworkIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldNetworkID,
		})
	}
	if value, ok := scsuo.mutation.SystemID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldSystemID,
		})
	}
	if scsuo.mutation.SystemIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldSystemID,
		})
	}
	if value, ok := scsuo.mutation.CellID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldCellID,
		})
	}
	if scsuo.mutation.CellIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldCellID,
		})
	}
	if value, ok := scsuo.mutation.LocationAreaCode(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldLocationAreaCode,
		})
	}
	if scsuo.mutation.LocationAreaCodeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldLocationAreaCode,
		})
	}
	if value, ok := scsuo.mutation.MobileCountryCode(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldMobileCountryCode,
		})
	}
	if scsuo.mutation.MobileCountryCodeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldMobileCountryCode,
		})
	}
	if value, ok := scsuo.mutation.MobileNetworkCode(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldMobileNetworkCode,
		})
	}
	if scsuo.mutation.MobileNetworkCodeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldMobileNetworkCode,
		})
	}
	if value, ok := scsuo.mutation.PrimaryScramblingCode(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldPrimaryScramblingCode,
		})
	}
	if scsuo.mutation.PrimaryScramblingCodeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldPrimaryScramblingCode,
		})
	}
	if value, ok := scsuo.mutation.Operator(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldOperator,
		})
	}
	if scsuo.mutation.OperatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldOperator,
		})
	}
	if value, ok := scsuo.mutation.Arfcn(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldArfcn,
		})
	}
	if value, ok := scsuo.mutation.AddedArfcn(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldArfcn,
		})
	}
	if scsuo.mutation.ArfcnCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveycellscan.FieldArfcn,
		})
	}
	if value, ok := scsuo.mutation.PhysicalCellID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldPhysicalCellID,
		})
	}
	if scsuo.mutation.PhysicalCellIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldPhysicalCellID,
		})
	}
	if value, ok := scsuo.mutation.TrackingAreaCode(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldTrackingAreaCode,
		})
	}
	if scsuo.mutation.TrackingAreaCodeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveycellscan.FieldTrackingAreaCode,
		})
	}
	if value, ok := scsuo.mutation.TimingAdvance(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldTimingAdvance,
		})
	}
	if value, ok := scsuo.mutation.AddedTimingAdvance(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldTimingAdvance,
		})
	}
	if scsuo.mutation.TimingAdvanceCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveycellscan.FieldTimingAdvance,
		})
	}
	if value, ok := scsuo.mutation.Earfcn(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldEarfcn,
		})
	}
	if value, ok := scsuo.mutation.AddedEarfcn(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldEarfcn,
		})
	}
	if scsuo.mutation.EarfcnCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveycellscan.FieldEarfcn,
		})
	}
	if value, ok := scsuo.mutation.Uarfcn(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldUarfcn,
		})
	}
	if value, ok := scsuo.mutation.AddedUarfcn(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldUarfcn,
		})
	}
	if scsuo.mutation.UarfcnCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveycellscan.FieldUarfcn,
		})
	}
	if value, ok := scsuo.mutation.Latitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveycellscan.FieldLatitude,
		})
	}
	if value, ok := scsuo.mutation.AddedLatitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveycellscan.FieldLatitude,
		})
	}
	if scsuo.mutation.LatitudeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveycellscan.FieldLatitude,
		})
	}
	if value, ok := scsuo.mutation.Longitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveycellscan.FieldLongitude,
		})
	}
	if value, ok := scsuo.mutation.AddedLongitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveycellscan.FieldLongitude,
		})
	}
	if scsuo.mutation.LongitudeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveycellscan.FieldLongitude,
		})
	}
	if scsuo.mutation.SurveyQuestionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveycellscan.SurveyQuestionTable,
			Columns: []string{surveycellscan.SurveyQuestionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveyquestion.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := scsuo.mutation.SurveyQuestionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveycellscan.SurveyQuestionTable,
			Columns: []string{surveycellscan.SurveyQuestionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveyquestion.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if scsuo.mutation.LocationCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveycellscan.LocationTable,
			Columns: []string{surveycellscan.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := scsuo.mutation.LocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveycellscan.LocationTable,
			Columns: []string{surveycellscan.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	scs = &SurveyCellScan{config: scsuo.config}
	_spec.Assign = scs.assignValues
	_spec.ScanValues = scs.scanValues()
	if err = sqlgraph.UpdateNode(ctx, scsuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{surveycellscan.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return scs, nil
}
