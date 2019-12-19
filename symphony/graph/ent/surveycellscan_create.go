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

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
)

// SurveyCellScanCreate is the builder for creating a SurveyCellScan entity.
type SurveyCellScanCreate struct {
	config
	create_time             *time.Time
	update_time             *time.Time
	network_type            *string
	signal_strength         *int
	timestamp               *time.Time
	base_station_id         *string
	network_id              *string
	system_id               *string
	cell_id                 *string
	location_area_code      *string
	mobile_country_code     *string
	mobile_network_code     *string
	primary_scrambling_code *string
	operator                *string
	arfcn                   *int
	physical_cell_id        *string
	tracking_area_code      *string
	timing_advance          *int
	earfcn                  *int
	uarfcn                  *int
	latitude                *float64
	longitude               *float64
	survey_question         map[string]struct{}
	location                map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (scsc *SurveyCellScanCreate) SetCreateTime(t time.Time) *SurveyCellScanCreate {
	scsc.create_time = &t
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
	scsc.update_time = &t
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
	scsc.network_type = &s
	return scsc
}

// SetSignalStrength sets the signal_strength field.
func (scsc *SurveyCellScanCreate) SetSignalStrength(i int) *SurveyCellScanCreate {
	scsc.signal_strength = &i
	return scsc
}

// SetTimestamp sets the timestamp field.
func (scsc *SurveyCellScanCreate) SetTimestamp(t time.Time) *SurveyCellScanCreate {
	scsc.timestamp = &t
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
	scsc.base_station_id = &s
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
	scsc.network_id = &s
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
	scsc.system_id = &s
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
	scsc.cell_id = &s
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
	scsc.location_area_code = &s
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
	scsc.mobile_country_code = &s
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
	scsc.mobile_network_code = &s
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
	scsc.primary_scrambling_code = &s
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
	scsc.operator = &s
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
	scsc.arfcn = &i
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
	scsc.physical_cell_id = &s
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
	scsc.tracking_area_code = &s
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
	scsc.timing_advance = &i
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
	scsc.earfcn = &i
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
	scsc.uarfcn = &i
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
	scsc.latitude = &f
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
	scsc.longitude = &f
	return scsc
}

// SetNillableLongitude sets the longitude field if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableLongitude(f *float64) *SurveyCellScanCreate {
	if f != nil {
		scsc.SetLongitude(*f)
	}
	return scsc
}

// SetSurveyQuestionID sets the survey_question edge to SurveyQuestion by id.
func (scsc *SurveyCellScanCreate) SetSurveyQuestionID(id string) *SurveyCellScanCreate {
	if scsc.survey_question == nil {
		scsc.survey_question = make(map[string]struct{})
	}
	scsc.survey_question[id] = struct{}{}
	return scsc
}

// SetNillableSurveyQuestionID sets the survey_question edge to SurveyQuestion by id if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableSurveyQuestionID(id *string) *SurveyCellScanCreate {
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
func (scsc *SurveyCellScanCreate) SetLocationID(id string) *SurveyCellScanCreate {
	if scsc.location == nil {
		scsc.location = make(map[string]struct{})
	}
	scsc.location[id] = struct{}{}
	return scsc
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (scsc *SurveyCellScanCreate) SetNillableLocationID(id *string) *SurveyCellScanCreate {
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
	if scsc.create_time == nil {
		v := surveycellscan.DefaultCreateTime()
		scsc.create_time = &v
	}
	if scsc.update_time == nil {
		v := surveycellscan.DefaultUpdateTime()
		scsc.update_time = &v
	}
	if scsc.network_type == nil {
		return nil, errors.New("ent: missing required field \"network_type\"")
	}
	if scsc.signal_strength == nil {
		return nil, errors.New("ent: missing required field \"signal_strength\"")
	}
	if len(scsc.survey_question) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"survey_question\"")
	}
	if len(scsc.location) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	return scsc.sqlSave(ctx)
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
		scs  = &SurveyCellScan{config: scsc.config}
		spec = &sqlgraph.CreateSpec{
			Table: surveycellscan.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: surveycellscan.FieldID,
			},
		}
	)
	if value := scsc.create_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveycellscan.FieldCreateTime,
		})
		scs.CreateTime = *value
	}
	if value := scsc.update_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveycellscan.FieldUpdateTime,
		})
		scs.UpdateTime = *value
	}
	if value := scsc.network_type; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldNetworkType,
		})
		scs.NetworkType = *value
	}
	if value := scsc.signal_strength; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldSignalStrength,
		})
		scs.SignalStrength = *value
	}
	if value := scsc.timestamp; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveycellscan.FieldTimestamp,
		})
		scs.Timestamp = *value
	}
	if value := scsc.base_station_id; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldBaseStationID,
		})
		scs.BaseStationID = *value
	}
	if value := scsc.network_id; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldNetworkID,
		})
		scs.NetworkID = *value
	}
	if value := scsc.system_id; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldSystemID,
		})
		scs.SystemID = *value
	}
	if value := scsc.cell_id; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldCellID,
		})
		scs.CellID = *value
	}
	if value := scsc.location_area_code; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldLocationAreaCode,
		})
		scs.LocationAreaCode = *value
	}
	if value := scsc.mobile_country_code; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldMobileCountryCode,
		})
		scs.MobileCountryCode = *value
	}
	if value := scsc.mobile_network_code; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldMobileNetworkCode,
		})
		scs.MobileNetworkCode = *value
	}
	if value := scsc.primary_scrambling_code; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldPrimaryScramblingCode,
		})
		scs.PrimaryScramblingCode = *value
	}
	if value := scsc.operator; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldOperator,
		})
		scs.Operator = *value
	}
	if value := scsc.arfcn; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldArfcn,
		})
		scs.Arfcn = *value
	}
	if value := scsc.physical_cell_id; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldPhysicalCellID,
		})
		scs.PhysicalCellID = *value
	}
	if value := scsc.tracking_area_code; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveycellscan.FieldTrackingAreaCode,
		})
		scs.TrackingAreaCode = *value
	}
	if value := scsc.timing_advance; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldTimingAdvance,
		})
		scs.TimingAdvance = *value
	}
	if value := scsc.earfcn; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldEarfcn,
		})
		scs.Earfcn = *value
	}
	if value := scsc.uarfcn; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveycellscan.FieldUarfcn,
		})
		scs.Uarfcn = *value
	}
	if value := scsc.latitude; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveycellscan.FieldLatitude,
		})
		scs.Latitude = *value
	}
	if value := scsc.longitude; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveycellscan.FieldLongitude,
		})
		scs.Longitude = *value
	}
	if nodes := scsc.survey_question; len(nodes) > 0 {
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
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := scsc.location; len(nodes) > 0 {
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
		spec.Edges = append(spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, scsc.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := spec.ID.Value.(int64)
	scs.ID = strconv.FormatInt(id, 10)
	return scs, nil
}
