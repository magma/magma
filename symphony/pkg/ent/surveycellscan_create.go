// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/ent/checklistitem"
	"github.com/facebookincubator/symphony/pkg/ent/location"
	"github.com/facebookincubator/symphony/pkg/ent/surveycellscan"
	"github.com/facebookincubator/symphony/pkg/ent/surveyquestion"
)

// SurveyCellScanCreate is the builder for creating a SurveyCellScan entity.
type SurveyCellScanCreate struct {
	config
	mutation *SurveyCellScanMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (scsc *SurveyCellScanCreate) SetCreateTime(t time.Time) *SurveyCellScanCreate {
	scsc.mutation.SetCreateTime(t)
	return scsc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableCreateTime(t *time.Time) *SurveyCellScanCreate {
	if t != nil {
		scsc.SetCreateTime(*t)
	}
	return scsc
}

// SetUpdateTime sets the update_time field.
func (scsc *SurveyCellScanCreate) SetUpdateTime(t time.Time) *SurveyCellScanCreate {
	scsc.mutation.SetUpdateTime(t)
	return scsc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableUpdateTime(t *time.Time) *SurveyCellScanCreate {
	if t != nil {
		scsc.SetUpdateTime(*t)
	}
	return scsc
}

// SetNetworkType sets the network_type field.
func (scsc *SurveyCellScanCreate) SetNetworkType(s string) *SurveyCellScanCreate {
	scsc.mutation.SetNetworkType(s)
	return scsc
}

// SetSignalStrength sets the signal_strength field.
func (scsc *SurveyCellScanCreate) SetSignalStrength(i int) *SurveyCellScanCreate {
	scsc.mutation.SetSignalStrength(i)
	return scsc
}

// SetTimestamp sets the timestamp field.
func (scsc *SurveyCellScanCreate) SetTimestamp(t time.Time) *SurveyCellScanCreate {
	scsc.mutation.SetTimestamp(t)
	return scsc
}

// SetNillableTimestamp sets the timestamp field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableTimestamp(t *time.Time) *SurveyCellScanCreate {
	if t != nil {
		scsc.SetTimestamp(*t)
	}
	return scsc
}

// SetBaseStationID sets the base_station_id field.
func (scsc *SurveyCellScanCreate) SetBaseStationID(s string) *SurveyCellScanCreate {
	scsc.mutation.SetBaseStationID(s)
	return scsc
}

// SetNillableBaseStationID sets the base_station_id field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableBaseStationID(s *string) *SurveyCellScanCreate {
	if s != nil {
		scsc.SetBaseStationID(*s)
	}
	return scsc
}

// SetNetworkID sets the network_id field.
func (scsc *SurveyCellScanCreate) SetNetworkID(s string) *SurveyCellScanCreate {
	scsc.mutation.SetNetworkID(s)
	return scsc
}

// SetNillableNetworkID sets the network_id field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableNetworkID(s *string) *SurveyCellScanCreate {
	if s != nil {
		scsc.SetNetworkID(*s)
	}
	return scsc
}

// SetSystemID sets the system_id field.
func (scsc *SurveyCellScanCreate) SetSystemID(s string) *SurveyCellScanCreate {
	scsc.mutation.SetSystemID(s)
	return scsc
}

// SetNillableSystemID sets the system_id field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableSystemID(s *string) *SurveyCellScanCreate {
	if s != nil {
		scsc.SetSystemID(*s)
	}
	return scsc
}

// SetCellID sets the cell_id field.
func (scsc *SurveyCellScanCreate) SetCellID(s string) *SurveyCellScanCreate {
	scsc.mutation.SetCellID(s)
	return scsc
}

// SetNillableCellID sets the cell_id field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableCellID(s *string) *SurveyCellScanCreate {
	if s != nil {
		scsc.SetCellID(*s)
	}
	return scsc
}

// SetLocationAreaCode sets the location_area_code field.
func (scsc *SurveyCellScanCreate) SetLocationAreaCode(s string) *SurveyCellScanCreate {
	scsc.mutation.SetLocationAreaCode(s)
	return scsc
}

// SetNillableLocationAreaCode sets the location_area_code field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableLocationAreaCode(s *string) *SurveyCellScanCreate {
	if s != nil {
		scsc.SetLocationAreaCode(*s)
	}
	return scsc
}

// SetMobileCountryCode sets the mobile_country_code field.
func (scsc *SurveyCellScanCreate) SetMobileCountryCode(s string) *SurveyCellScanCreate {
	scsc.mutation.SetMobileCountryCode(s)
	return scsc
}

// SetNillableMobileCountryCode sets the mobile_country_code field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableMobileCountryCode(s *string) *SurveyCellScanCreate {
	if s != nil {
		scsc.SetMobileCountryCode(*s)
	}
	return scsc
}

// SetMobileNetworkCode sets the mobile_network_code field.
func (scsc *SurveyCellScanCreate) SetMobileNetworkCode(s string) *SurveyCellScanCreate {
	scsc.mutation.SetMobileNetworkCode(s)
	return scsc
}

// SetNillableMobileNetworkCode sets the mobile_network_code field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableMobileNetworkCode(s *string) *SurveyCellScanCreate {
	if s != nil {
		scsc.SetMobileNetworkCode(*s)
	}
	return scsc
}

// SetPrimaryScramblingCode sets the primary_scrambling_code field.
func (scsc *SurveyCellScanCreate) SetPrimaryScramblingCode(s string) *SurveyCellScanCreate {
	scsc.mutation.SetPrimaryScramblingCode(s)
	return scsc
}

// SetNillablePrimaryScramblingCode sets the primary_scrambling_code field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillablePrimaryScramblingCode(s *string) *SurveyCellScanCreate {
	if s != nil {
		scsc.SetPrimaryScramblingCode(*s)
	}
	return scsc
}

// SetOperator sets the operator field.
func (scsc *SurveyCellScanCreate) SetOperator(s string) *SurveyCellScanCreate {
	scsc.mutation.SetOperator(s)
	return scsc
}

// SetNillableOperator sets the operator field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableOperator(s *string) *SurveyCellScanCreate {
	if s != nil {
		scsc.SetOperator(*s)
	}
	return scsc
}

// SetArfcn sets the arfcn field.
func (scsc *SurveyCellScanCreate) SetArfcn(i int) *SurveyCellScanCreate {
	scsc.mutation.SetArfcn(i)
	return scsc
}

// SetNillableArfcn sets the arfcn field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableArfcn(i *int) *SurveyCellScanCreate {
	if i != nil {
		scsc.SetArfcn(*i)
	}
	return scsc
}

// SetPhysicalCellID sets the physical_cell_id field.
func (scsc *SurveyCellScanCreate) SetPhysicalCellID(s string) *SurveyCellScanCreate {
	scsc.mutation.SetPhysicalCellID(s)
	return scsc
}

// SetNillablePhysicalCellID sets the physical_cell_id field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillablePhysicalCellID(s *string) *SurveyCellScanCreate {
	if s != nil {
		scsc.SetPhysicalCellID(*s)
	}
	return scsc
}

// SetTrackingAreaCode sets the tracking_area_code field.
func (scsc *SurveyCellScanCreate) SetTrackingAreaCode(s string) *SurveyCellScanCreate {
	scsc.mutation.SetTrackingAreaCode(s)
	return scsc
}

// SetNillableTrackingAreaCode sets the tracking_area_code field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableTrackingAreaCode(s *string) *SurveyCellScanCreate {
	if s != nil {
		scsc.SetTrackingAreaCode(*s)
	}
	return scsc
}

// SetTimingAdvance sets the timing_advance field.
func (scsc *SurveyCellScanCreate) SetTimingAdvance(i int) *SurveyCellScanCreate {
	scsc.mutation.SetTimingAdvance(i)
	return scsc
}

// SetNillableTimingAdvance sets the timing_advance field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableTimingAdvance(i *int) *SurveyCellScanCreate {
	if i != nil {
		scsc.SetTimingAdvance(*i)
	}
	return scsc
}

// SetEarfcn sets the earfcn field.
func (scsc *SurveyCellScanCreate) SetEarfcn(i int) *SurveyCellScanCreate {
	scsc.mutation.SetEarfcn(i)
	return scsc
}

// SetNillableEarfcn sets the earfcn field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableEarfcn(i *int) *SurveyCellScanCreate {
	if i != nil {
		scsc.SetEarfcn(*i)
	}
	return scsc
}

// SetUarfcn sets the uarfcn field.
func (scsc *SurveyCellScanCreate) SetUarfcn(i int) *SurveyCellScanCreate {
	scsc.mutation.SetUarfcn(i)
	return scsc
}

// SetNillableUarfcn sets the uarfcn field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableUarfcn(i *int) *SurveyCellScanCreate {
	if i != nil {
		scsc.SetUarfcn(*i)
	}
	return scsc
}

// SetLatitude sets the latitude field.
func (scsc *SurveyCellScanCreate) SetLatitude(f float64) *SurveyCellScanCreate {
	scsc.mutation.SetLatitude(f)
	return scsc
}

// SetNillableLatitude sets the latitude field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableLatitude(f *float64) *SurveyCellScanCreate {
	if f != nil {
		scsc.SetLatitude(*f)
	}
	return scsc
}

// SetLongitude sets the longitude field.
func (scsc *SurveyCellScanCreate) SetLongitude(f float64) *SurveyCellScanCreate {
	scsc.mutation.SetLongitude(f)
	return scsc
}

// SetNillableLongitude sets the longitude field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableLongitude(f *float64) *SurveyCellScanCreate {
	if f != nil {
		scsc.SetLongitude(*f)
	}
	return scsc
}

// SetChecklistItemID sets the checklist_item edge to CheckListItem by id.
func (scsc *SurveyCellScanCreate) SetChecklistItemID(id int) *SurveyCellScanCreate {
	scsc.mutation.SetChecklistItemID(id)
	return scsc
}

// SetNillableChecklistItemID sets the checklist_item edge to CheckListItem by id if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableChecklistItemID(id *int) *SurveyCellScanCreate {
	if id != nil {
		scsc = scsc.SetChecklistItemID(*id)
	}
	return scsc
}

// SetChecklistItem sets the checklist_item edge to CheckListItem.
func (scsc *SurveyCellScanCreate) SetChecklistItem(c *CheckListItem) *SurveyCellScanCreate {
	return scsc.SetChecklistItemID(c.ID)
}

// SetSurveyQuestionID sets the survey_question edge to SurveyQuestion by id.
func (scsc *SurveyCellScanCreate) SetSurveyQuestionID(id int) *SurveyCellScanCreate {
	scsc.mutation.SetSurveyQuestionID(id)
	return scsc
}

// SetNillableSurveyQuestionID sets the survey_question edge to SurveyQuestion by id if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableSurveyQuestionID(id *int) *SurveyCellScanCreate {
	if id != nil {
		scsc = scsc.SetSurveyQuestionID(*id)
	}
	return scsc
}

// SetSurveyQuestion sets the survey_question edge to SurveyQuestion.
func (scsc *SurveyCellScanCreate) SetSurveyQuestion(s *SurveyQuestion) *SurveyCellScanCreate {
	return scsc.SetSurveyQuestionID(s.ID)
}

// SetLocationID sets the location edge to Location by id.
func (scsc *SurveyCellScanCreate) SetLocationID(id int) *SurveyCellScanCreate {
	scsc.mutation.SetLocationID(id)
	return scsc
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableLocationID(id *int) *SurveyCellScanCreate {
	if id != nil {
		scsc = scsc.SetLocationID(*id)
	}
	return scsc
}

// SetLocation sets the location edge to Location.
func (scsc *SurveyCellScanCreate) SetLocation(l *Location) *SurveyCellScanCreate {
	return scsc.SetLocationID(l.ID)
}

// Save creates the SurveyCellScan in the database.
func (scsc *SurveyCellScanCreate) Save(ctx context.Context) (*SurveyCellScan, error) {
	if _, ok := scsc.mutation.CreateTime(); !ok {
		v := surveycellscan.DefaultCreateTime()
		scsc.mutation.SetCreateTime(v)
	}
	if _, ok := scsc.mutation.UpdateTime(); !ok {
		v := surveycellscan.DefaultUpdateTime()
		scsc.mutation.SetUpdateTime(v)
	}
	if _, ok := scsc.mutation.NetworkType(); !ok {
		return nil, errors.New("ent: missing required field \"network_type\"")
	}
	if _, ok := scsc.mutation.SignalStrength(); !ok {
		return nil, errors.New("ent: missing required field \"signal_strength\"")
	}
	var (
		err  error
		node *SurveyCellScan
	)
	if len(scsc.hooks) == 0 {
		node, err = scsc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyCellScanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			scsc.mutation = mutation
			node, err = scsc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(scsc.hooks) - 1; i >= 0; i-- {
			mut = scsc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, scsc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (scsc *SurveyCellScanCreate) SaveX(ctx context.Context) *SurveyCellScan {
	v, err := scsc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (scsc *SurveyCellScanCreate) sqlSave(ctx context.Context) (*SurveyCellScan, error) {
	var (
		scs   = &SurveyCellScan{config: scsc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: surveycellscan.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: surveycellscan.FieldID,
			},
		}
	)
	if value, ok := scsc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveycellscan.FieldCreateTime,
		})
		scs.CreateTime = value
	}
	if value, ok := scsc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveycellscan.FieldUpdateTime,
		})
		scs.UpdateTime = value
	}
	if value, ok := scsc.mutation.NetworkType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldNetworkType,
		})
		scs.NetworkType = value
	}
	if value, ok := scsc.mutation.SignalStrength(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldSignalStrength,
		})
		scs.SignalStrength = value
	}
	if value, ok := scsc.mutation.Timestamp(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveycellscan.FieldTimestamp,
		})
		scs.Timestamp = value
	}
	if value, ok := scsc.mutation.BaseStationID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldBaseStationID,
		})
		scs.BaseStationID = value
	}
	if value, ok := scsc.mutation.NetworkID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldNetworkID,
		})
		scs.NetworkID = value
	}
	if value, ok := scsc.mutation.SystemID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldSystemID,
		})
		scs.SystemID = value
	}
	if value, ok := scsc.mutation.CellID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldCellID,
		})
		scs.CellID = value
	}
	if value, ok := scsc.mutation.LocationAreaCode(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldLocationAreaCode,
		})
		scs.LocationAreaCode = value
	}
	if value, ok := scsc.mutation.MobileCountryCode(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldMobileCountryCode,
		})
		scs.MobileCountryCode = value
	}
	if value, ok := scsc.mutation.MobileNetworkCode(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldMobileNetworkCode,
		})
		scs.MobileNetworkCode = value
	}
	if value, ok := scsc.mutation.PrimaryScramblingCode(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldPrimaryScramblingCode,
		})
		scs.PrimaryScramblingCode = value
	}
	if value, ok := scsc.mutation.Operator(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldOperator,
		})
		scs.Operator = value
	}
	if value, ok := scsc.mutation.Arfcn(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldArfcn,
		})
		scs.Arfcn = value
	}
	if value, ok := scsc.mutation.PhysicalCellID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldPhysicalCellID,
		})
		scs.PhysicalCellID = value
	}
	if value, ok := scsc.mutation.TrackingAreaCode(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveycellscan.FieldTrackingAreaCode,
		})
		scs.TrackingAreaCode = value
	}
	if value, ok := scsc.mutation.TimingAdvance(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldTimingAdvance,
		})
		scs.TimingAdvance = value
	}
	if value, ok := scsc.mutation.Earfcn(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldEarfcn,
		})
		scs.Earfcn = value
	}
	if value, ok := scsc.mutation.Uarfcn(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveycellscan.FieldUarfcn,
		})
		scs.Uarfcn = value
	}
	if value, ok := scsc.mutation.Latitude(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveycellscan.FieldLatitude,
		})
		scs.Latitude = value
	}
	if value, ok := scsc.mutation.Longitude(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveycellscan.FieldLongitude,
		})
		scs.Longitude = value
	}
	if nodes := scsc.mutation.ChecklistItemIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveycellscan.ChecklistItemTable,
			Columns: []string{surveycellscan.ChecklistItemColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistitem.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := scsc.mutation.SurveyQuestionIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := scsc.mutation.LocationIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, scsc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	scs.ID = int(id)
	return scs, nil
}
