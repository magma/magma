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
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
)

// SurveyWiFiScanUpdate is the builder for updating SurveyWiFiScan entities.
type SurveyWiFiScanUpdate struct {
	config

	update_time           *time.Time
	ssid                  *string
	clearssid             bool
	bssid                 *string
	timestamp             *time.Time
	frequency             *int
	addfrequency          *int
	channel               *int
	addchannel            *int
	band                  *string
	clearband             bool
	channel_width         *int
	addchannel_width      *int
	clearchannel_width    bool
	capabilities          *string
	clearcapabilities     bool
	strength              *int
	addstrength           *int
	latitude              *float64
	addlatitude           *float64
	clearlatitude         bool
	longitude             *float64
	addlongitude          *float64
	clearlongitude        bool
	survey_question       map[string]struct{}
	location              map[string]struct{}
	clearedSurveyQuestion bool
	clearedLocation       bool
	predicates            []predicate.SurveyWiFiScan
}

// Where adds a new predicate for the builder.
func (swfsu *SurveyWiFiScanUpdate) Where(ps ...predicate.SurveyWiFiScan) *SurveyWiFiScanUpdate {
	swfsu.predicates = append(swfsu.predicates, ps...)
	return swfsu
}

// SetSsid sets the ssid field.
func (swfsu *SurveyWiFiScanUpdate) SetSsid(s string) *SurveyWiFiScanUpdate {
	swfsu.ssid = &s
	return swfsu
}

// SetNillableSsid sets the ssid field if the given value is not nil.
func (swfsu *SurveyWiFiScanUpdate) SetNillableSsid(s *string) *SurveyWiFiScanUpdate {
	if s != nil {
		swfsu.SetSsid(*s)
	}
	return swfsu
}

// ClearSsid clears the value of ssid.
func (swfsu *SurveyWiFiScanUpdate) ClearSsid() *SurveyWiFiScanUpdate {
	swfsu.ssid = nil
	swfsu.clearssid = true
	return swfsu
}

// SetBssid sets the bssid field.
func (swfsu *SurveyWiFiScanUpdate) SetBssid(s string) *SurveyWiFiScanUpdate {
	swfsu.bssid = &s
	return swfsu
}

// SetTimestamp sets the timestamp field.
func (swfsu *SurveyWiFiScanUpdate) SetTimestamp(t time.Time) *SurveyWiFiScanUpdate {
	swfsu.timestamp = &t
	return swfsu
}

// SetFrequency sets the frequency field.
func (swfsu *SurveyWiFiScanUpdate) SetFrequency(i int) *SurveyWiFiScanUpdate {
	swfsu.frequency = &i
	swfsu.addfrequency = nil
	return swfsu
}

// AddFrequency adds i to frequency.
func (swfsu *SurveyWiFiScanUpdate) AddFrequency(i int) *SurveyWiFiScanUpdate {
	if swfsu.addfrequency == nil {
		swfsu.addfrequency = &i
	} else {
		*swfsu.addfrequency += i
	}
	return swfsu
}

// SetChannel sets the channel field.
func (swfsu *SurveyWiFiScanUpdate) SetChannel(i int) *SurveyWiFiScanUpdate {
	swfsu.channel = &i
	swfsu.addchannel = nil
	return swfsu
}

// AddChannel adds i to channel.
func (swfsu *SurveyWiFiScanUpdate) AddChannel(i int) *SurveyWiFiScanUpdate {
	if swfsu.addchannel == nil {
		swfsu.addchannel = &i
	} else {
		*swfsu.addchannel += i
	}
	return swfsu
}

// SetBand sets the band field.
func (swfsu *SurveyWiFiScanUpdate) SetBand(s string) *SurveyWiFiScanUpdate {
	swfsu.band = &s
	return swfsu
}

// SetNillableBand sets the band field if the given value is not nil.
func (swfsu *SurveyWiFiScanUpdate) SetNillableBand(s *string) *SurveyWiFiScanUpdate {
	if s != nil {
		swfsu.SetBand(*s)
	}
	return swfsu
}

// ClearBand clears the value of band.
func (swfsu *SurveyWiFiScanUpdate) ClearBand() *SurveyWiFiScanUpdate {
	swfsu.band = nil
	swfsu.clearband = true
	return swfsu
}

// SetChannelWidth sets the channel_width field.
func (swfsu *SurveyWiFiScanUpdate) SetChannelWidth(i int) *SurveyWiFiScanUpdate {
	swfsu.channel_width = &i
	swfsu.addchannel_width = nil
	return swfsu
}

// SetNillableChannelWidth sets the channel_width field if the given value is not nil.
func (swfsu *SurveyWiFiScanUpdate) SetNillableChannelWidth(i *int) *SurveyWiFiScanUpdate {
	if i != nil {
		swfsu.SetChannelWidth(*i)
	}
	return swfsu
}

// AddChannelWidth adds i to channel_width.
func (swfsu *SurveyWiFiScanUpdate) AddChannelWidth(i int) *SurveyWiFiScanUpdate {
	if swfsu.addchannel_width == nil {
		swfsu.addchannel_width = &i
	} else {
		*swfsu.addchannel_width += i
	}
	return swfsu
}

// ClearChannelWidth clears the value of channel_width.
func (swfsu *SurveyWiFiScanUpdate) ClearChannelWidth() *SurveyWiFiScanUpdate {
	swfsu.channel_width = nil
	swfsu.clearchannel_width = true
	return swfsu
}

// SetCapabilities sets the capabilities field.
func (swfsu *SurveyWiFiScanUpdate) SetCapabilities(s string) *SurveyWiFiScanUpdate {
	swfsu.capabilities = &s
	return swfsu
}

// SetNillableCapabilities sets the capabilities field if the given value is not nil.
func (swfsu *SurveyWiFiScanUpdate) SetNillableCapabilities(s *string) *SurveyWiFiScanUpdate {
	if s != nil {
		swfsu.SetCapabilities(*s)
	}
	return swfsu
}

// ClearCapabilities clears the value of capabilities.
func (swfsu *SurveyWiFiScanUpdate) ClearCapabilities() *SurveyWiFiScanUpdate {
	swfsu.capabilities = nil
	swfsu.clearcapabilities = true
	return swfsu
}

// SetStrength sets the strength field.
func (swfsu *SurveyWiFiScanUpdate) SetStrength(i int) *SurveyWiFiScanUpdate {
	swfsu.strength = &i
	swfsu.addstrength = nil
	return swfsu
}

// AddStrength adds i to strength.
func (swfsu *SurveyWiFiScanUpdate) AddStrength(i int) *SurveyWiFiScanUpdate {
	if swfsu.addstrength == nil {
		swfsu.addstrength = &i
	} else {
		*swfsu.addstrength += i
	}
	return swfsu
}

// SetLatitude sets the latitude field.
func (swfsu *SurveyWiFiScanUpdate) SetLatitude(f float64) *SurveyWiFiScanUpdate {
	swfsu.latitude = &f
	swfsu.addlatitude = nil
	return swfsu
}

// SetNillableLatitude sets the latitude field if the given value is not nil.
func (swfsu *SurveyWiFiScanUpdate) SetNillableLatitude(f *float64) *SurveyWiFiScanUpdate {
	if f != nil {
		swfsu.SetLatitude(*f)
	}
	return swfsu
}

// AddLatitude adds f to latitude.
func (swfsu *SurveyWiFiScanUpdate) AddLatitude(f float64) *SurveyWiFiScanUpdate {
	if swfsu.addlatitude == nil {
		swfsu.addlatitude = &f
	} else {
		*swfsu.addlatitude += f
	}
	return swfsu
}

// ClearLatitude clears the value of latitude.
func (swfsu *SurveyWiFiScanUpdate) ClearLatitude() *SurveyWiFiScanUpdate {
	swfsu.latitude = nil
	swfsu.clearlatitude = true
	return swfsu
}

// SetLongitude sets the longitude field.
func (swfsu *SurveyWiFiScanUpdate) SetLongitude(f float64) *SurveyWiFiScanUpdate {
	swfsu.longitude = &f
	swfsu.addlongitude = nil
	return swfsu
}

// SetNillableLongitude sets the longitude field if the given value is not nil.
func (swfsu *SurveyWiFiScanUpdate) SetNillableLongitude(f *float64) *SurveyWiFiScanUpdate {
	if f != nil {
		swfsu.SetLongitude(*f)
	}
	return swfsu
}

// AddLongitude adds f to longitude.
func (swfsu *SurveyWiFiScanUpdate) AddLongitude(f float64) *SurveyWiFiScanUpdate {
	if swfsu.addlongitude == nil {
		swfsu.addlongitude = &f
	} else {
		*swfsu.addlongitude += f
	}
	return swfsu
}

// ClearLongitude clears the value of longitude.
func (swfsu *SurveyWiFiScanUpdate) ClearLongitude() *SurveyWiFiScanUpdate {
	swfsu.longitude = nil
	swfsu.clearlongitude = true
	return swfsu
}

// SetSurveyQuestionID sets the survey_question edge to SurveyQuestion by id.
func (swfsu *SurveyWiFiScanUpdate) SetSurveyQuestionID(id string) *SurveyWiFiScanUpdate {
	if swfsu.survey_question == nil {
		swfsu.survey_question = make(map[string]struct{})
	}
	swfsu.survey_question[id] = struct{}{}
	return swfsu
}

// SetNillableSurveyQuestionID sets the survey_question edge to SurveyQuestion by id if the given value is not nil.
func (swfsu *SurveyWiFiScanUpdate) SetNillableSurveyQuestionID(id *string) *SurveyWiFiScanUpdate {
	if id != nil {
		swfsu = swfsu.SetSurveyQuestionID(*id)
	}
	return swfsu
}

// SetSurveyQuestion sets the survey_question edge to SurveyQuestion.
func (swfsu *SurveyWiFiScanUpdate) SetSurveyQuestion(s *SurveyQuestion) *SurveyWiFiScanUpdate {
	return swfsu.SetSurveyQuestionID(s.ID)
}

// SetLocationID sets the location edge to Location by id.
func (swfsu *SurveyWiFiScanUpdate) SetLocationID(id string) *SurveyWiFiScanUpdate {
	if swfsu.location == nil {
		swfsu.location = make(map[string]struct{})
	}
	swfsu.location[id] = struct{}{}
	return swfsu
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (swfsu *SurveyWiFiScanUpdate) SetNillableLocationID(id *string) *SurveyWiFiScanUpdate {
	if id != nil {
		swfsu = swfsu.SetLocationID(*id)
	}
	return swfsu
}

// SetLocation sets the location edge to Location.
func (swfsu *SurveyWiFiScanUpdate) SetLocation(l *Location) *SurveyWiFiScanUpdate {
	return swfsu.SetLocationID(l.ID)
}

// ClearSurveyQuestion clears the survey_question edge to SurveyQuestion.
func (swfsu *SurveyWiFiScanUpdate) ClearSurveyQuestion() *SurveyWiFiScanUpdate {
	swfsu.clearedSurveyQuestion = true
	return swfsu
}

// ClearLocation clears the location edge to Location.
func (swfsu *SurveyWiFiScanUpdate) ClearLocation() *SurveyWiFiScanUpdate {
	swfsu.clearedLocation = true
	return swfsu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (swfsu *SurveyWiFiScanUpdate) Save(ctx context.Context) (int, error) {
	if swfsu.update_time == nil {
		v := surveywifiscan.UpdateDefaultUpdateTime()
		swfsu.update_time = &v
	}
	if len(swfsu.survey_question) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"survey_question\"")
	}
	if len(swfsu.location) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	return swfsu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (swfsu *SurveyWiFiScanUpdate) SaveX(ctx context.Context) int {
	affected, err := swfsu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (swfsu *SurveyWiFiScanUpdate) Exec(ctx context.Context) error {
	_, err := swfsu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (swfsu *SurveyWiFiScanUpdate) ExecX(ctx context.Context) {
	if err := swfsu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (swfsu *SurveyWiFiScanUpdate) sqlSave(ctx context.Context) (n int, err error) {
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveywifiscan.Table,
			Columns: surveywifiscan.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: surveywifiscan.FieldID,
			},
		},
	}
	if ps := swfsu.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := swfsu.update_time; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveywifiscan.FieldUpdateTime,
		})
	}
	if value := swfsu.ssid; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveywifiscan.FieldSsid,
		})
	}
	if swfsu.clearssid {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveywifiscan.FieldSsid,
		})
	}
	if value := swfsu.bssid; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveywifiscan.FieldBssid,
		})
	}
	if value := swfsu.timestamp; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveywifiscan.FieldTimestamp,
		})
	}
	if value := swfsu.frequency; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveywifiscan.FieldFrequency,
		})
	}
	if value := swfsu.addfrequency; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveywifiscan.FieldFrequency,
		})
	}
	if value := swfsu.channel; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveywifiscan.FieldChannel,
		})
	}
	if value := swfsu.addchannel; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveywifiscan.FieldChannel,
		})
	}
	if value := swfsu.band; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveywifiscan.FieldBand,
		})
	}
	if swfsu.clearband {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveywifiscan.FieldBand,
		})
	}
	if value := swfsu.channel_width; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveywifiscan.FieldChannelWidth,
		})
	}
	if value := swfsu.addchannel_width; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveywifiscan.FieldChannelWidth,
		})
	}
	if swfsu.clearchannel_width {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveywifiscan.FieldChannelWidth,
		})
	}
	if value := swfsu.capabilities; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveywifiscan.FieldCapabilities,
		})
	}
	if swfsu.clearcapabilities {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveywifiscan.FieldCapabilities,
		})
	}
	if value := swfsu.strength; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveywifiscan.FieldStrength,
		})
	}
	if value := swfsu.addstrength; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveywifiscan.FieldStrength,
		})
	}
	if value := swfsu.latitude; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveywifiscan.FieldLatitude,
		})
	}
	if value := swfsu.addlatitude; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveywifiscan.FieldLatitude,
		})
	}
	if swfsu.clearlatitude {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveywifiscan.FieldLatitude,
		})
	}
	if value := swfsu.longitude; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveywifiscan.FieldLongitude,
		})
	}
	if value := swfsu.addlongitude; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveywifiscan.FieldLongitude,
		})
	}
	if swfsu.clearlongitude {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveywifiscan.FieldLongitude,
		})
	}
	if swfsu.clearedSurveyQuestion {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveywifiscan.SurveyQuestionTable,
			Columns: []string{surveywifiscan.SurveyQuestionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveyquestion.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := swfsu.survey_question; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveywifiscan.SurveyQuestionTable,
			Columns: []string{surveywifiscan.SurveyQuestionColumn},
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
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if swfsu.clearedLocation {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveywifiscan.LocationTable,
			Columns: []string{surveywifiscan.LocationColumn},
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
	if nodes := swfsu.location; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveywifiscan.LocationTable,
			Columns: []string{surveywifiscan.LocationColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, swfsu.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// SurveyWiFiScanUpdateOne is the builder for updating a single SurveyWiFiScan entity.
type SurveyWiFiScanUpdateOne struct {
	config
	id string

	update_time           *time.Time
	ssid                  *string
	clearssid             bool
	bssid                 *string
	timestamp             *time.Time
	frequency             *int
	addfrequency          *int
	channel               *int
	addchannel            *int
	band                  *string
	clearband             bool
	channel_width         *int
	addchannel_width      *int
	clearchannel_width    bool
	capabilities          *string
	clearcapabilities     bool
	strength              *int
	addstrength           *int
	latitude              *float64
	addlatitude           *float64
	clearlatitude         bool
	longitude             *float64
	addlongitude          *float64
	clearlongitude        bool
	survey_question       map[string]struct{}
	location              map[string]struct{}
	clearedSurveyQuestion bool
	clearedLocation       bool
}

// SetSsid sets the ssid field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetSsid(s string) *SurveyWiFiScanUpdateOne {
	swfsuo.ssid = &s
	return swfsuo
}

// SetNillableSsid sets the ssid field if the given value is not nil.
func (swfsuo *SurveyWiFiScanUpdateOne) SetNillableSsid(s *string) *SurveyWiFiScanUpdateOne {
	if s != nil {
		swfsuo.SetSsid(*s)
	}
	return swfsuo
}

// ClearSsid clears the value of ssid.
func (swfsuo *SurveyWiFiScanUpdateOne) ClearSsid() *SurveyWiFiScanUpdateOne {
	swfsuo.ssid = nil
	swfsuo.clearssid = true
	return swfsuo
}

// SetBssid sets the bssid field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetBssid(s string) *SurveyWiFiScanUpdateOne {
	swfsuo.bssid = &s
	return swfsuo
}

// SetTimestamp sets the timestamp field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetTimestamp(t time.Time) *SurveyWiFiScanUpdateOne {
	swfsuo.timestamp = &t
	return swfsuo
}

// SetFrequency sets the frequency field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetFrequency(i int) *SurveyWiFiScanUpdateOne {
	swfsuo.frequency = &i
	swfsuo.addfrequency = nil
	return swfsuo
}

// AddFrequency adds i to frequency.
func (swfsuo *SurveyWiFiScanUpdateOne) AddFrequency(i int) *SurveyWiFiScanUpdateOne {
	if swfsuo.addfrequency == nil {
		swfsuo.addfrequency = &i
	} else {
		*swfsuo.addfrequency += i
	}
	return swfsuo
}

// SetChannel sets the channel field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetChannel(i int) *SurveyWiFiScanUpdateOne {
	swfsuo.channel = &i
	swfsuo.addchannel = nil
	return swfsuo
}

// AddChannel adds i to channel.
func (swfsuo *SurveyWiFiScanUpdateOne) AddChannel(i int) *SurveyWiFiScanUpdateOne {
	if swfsuo.addchannel == nil {
		swfsuo.addchannel = &i
	} else {
		*swfsuo.addchannel += i
	}
	return swfsuo
}

// SetBand sets the band field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetBand(s string) *SurveyWiFiScanUpdateOne {
	swfsuo.band = &s
	return swfsuo
}

// SetNillableBand sets the band field if the given value is not nil.
func (swfsuo *SurveyWiFiScanUpdateOne) SetNillableBand(s *string) *SurveyWiFiScanUpdateOne {
	if s != nil {
		swfsuo.SetBand(*s)
	}
	return swfsuo
}

// ClearBand clears the value of band.
func (swfsuo *SurveyWiFiScanUpdateOne) ClearBand() *SurveyWiFiScanUpdateOne {
	swfsuo.band = nil
	swfsuo.clearband = true
	return swfsuo
}

// SetChannelWidth sets the channel_width field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetChannelWidth(i int) *SurveyWiFiScanUpdateOne {
	swfsuo.channel_width = &i
	swfsuo.addchannel_width = nil
	return swfsuo
}

// SetNillableChannelWidth sets the channel_width field if the given value is not nil.
func (swfsuo *SurveyWiFiScanUpdateOne) SetNillableChannelWidth(i *int) *SurveyWiFiScanUpdateOne {
	if i != nil {
		swfsuo.SetChannelWidth(*i)
	}
	return swfsuo
}

// AddChannelWidth adds i to channel_width.
func (swfsuo *SurveyWiFiScanUpdateOne) AddChannelWidth(i int) *SurveyWiFiScanUpdateOne {
	if swfsuo.addchannel_width == nil {
		swfsuo.addchannel_width = &i
	} else {
		*swfsuo.addchannel_width += i
	}
	return swfsuo
}

// ClearChannelWidth clears the value of channel_width.
func (swfsuo *SurveyWiFiScanUpdateOne) ClearChannelWidth() *SurveyWiFiScanUpdateOne {
	swfsuo.channel_width = nil
	swfsuo.clearchannel_width = true
	return swfsuo
}

// SetCapabilities sets the capabilities field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetCapabilities(s string) *SurveyWiFiScanUpdateOne {
	swfsuo.capabilities = &s
	return swfsuo
}

// SetNillableCapabilities sets the capabilities field if the given value is not nil.
func (swfsuo *SurveyWiFiScanUpdateOne) SetNillableCapabilities(s *string) *SurveyWiFiScanUpdateOne {
	if s != nil {
		swfsuo.SetCapabilities(*s)
	}
	return swfsuo
}

// ClearCapabilities clears the value of capabilities.
func (swfsuo *SurveyWiFiScanUpdateOne) ClearCapabilities() *SurveyWiFiScanUpdateOne {
	swfsuo.capabilities = nil
	swfsuo.clearcapabilities = true
	return swfsuo
}

// SetStrength sets the strength field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetStrength(i int) *SurveyWiFiScanUpdateOne {
	swfsuo.strength = &i
	swfsuo.addstrength = nil
	return swfsuo
}

// AddStrength adds i to strength.
func (swfsuo *SurveyWiFiScanUpdateOne) AddStrength(i int) *SurveyWiFiScanUpdateOne {
	if swfsuo.addstrength == nil {
		swfsuo.addstrength = &i
	} else {
		*swfsuo.addstrength += i
	}
	return swfsuo
}

// SetLatitude sets the latitude field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetLatitude(f float64) *SurveyWiFiScanUpdateOne {
	swfsuo.latitude = &f
	swfsuo.addlatitude = nil
	return swfsuo
}

// SetNillableLatitude sets the latitude field if the given value is not nil.
func (swfsuo *SurveyWiFiScanUpdateOne) SetNillableLatitude(f *float64) *SurveyWiFiScanUpdateOne {
	if f != nil {
		swfsuo.SetLatitude(*f)
	}
	return swfsuo
}

// AddLatitude adds f to latitude.
func (swfsuo *SurveyWiFiScanUpdateOne) AddLatitude(f float64) *SurveyWiFiScanUpdateOne {
	if swfsuo.addlatitude == nil {
		swfsuo.addlatitude = &f
	} else {
		*swfsuo.addlatitude += f
	}
	return swfsuo
}

// ClearLatitude clears the value of latitude.
func (swfsuo *SurveyWiFiScanUpdateOne) ClearLatitude() *SurveyWiFiScanUpdateOne {
	swfsuo.latitude = nil
	swfsuo.clearlatitude = true
	return swfsuo
}

// SetLongitude sets the longitude field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetLongitude(f float64) *SurveyWiFiScanUpdateOne {
	swfsuo.longitude = &f
	swfsuo.addlongitude = nil
	return swfsuo
}

// SetNillableLongitude sets the longitude field if the given value is not nil.
func (swfsuo *SurveyWiFiScanUpdateOne) SetNillableLongitude(f *float64) *SurveyWiFiScanUpdateOne {
	if f != nil {
		swfsuo.SetLongitude(*f)
	}
	return swfsuo
}

// AddLongitude adds f to longitude.
func (swfsuo *SurveyWiFiScanUpdateOne) AddLongitude(f float64) *SurveyWiFiScanUpdateOne {
	if swfsuo.addlongitude == nil {
		swfsuo.addlongitude = &f
	} else {
		*swfsuo.addlongitude += f
	}
	return swfsuo
}

// ClearLongitude clears the value of longitude.
func (swfsuo *SurveyWiFiScanUpdateOne) ClearLongitude() *SurveyWiFiScanUpdateOne {
	swfsuo.longitude = nil
	swfsuo.clearlongitude = true
	return swfsuo
}

// SetSurveyQuestionID sets the survey_question edge to SurveyQuestion by id.
func (swfsuo *SurveyWiFiScanUpdateOne) SetSurveyQuestionID(id string) *SurveyWiFiScanUpdateOne {
	if swfsuo.survey_question == nil {
		swfsuo.survey_question = make(map[string]struct{})
	}
	swfsuo.survey_question[id] = struct{}{}
	return swfsuo
}

// SetNillableSurveyQuestionID sets the survey_question edge to SurveyQuestion by id if the given value is not nil.
func (swfsuo *SurveyWiFiScanUpdateOne) SetNillableSurveyQuestionID(id *string) *SurveyWiFiScanUpdateOne {
	if id != nil {
		swfsuo = swfsuo.SetSurveyQuestionID(*id)
	}
	return swfsuo
}

// SetSurveyQuestion sets the survey_question edge to SurveyQuestion.
func (swfsuo *SurveyWiFiScanUpdateOne) SetSurveyQuestion(s *SurveyQuestion) *SurveyWiFiScanUpdateOne {
	return swfsuo.SetSurveyQuestionID(s.ID)
}

// SetLocationID sets the location edge to Location by id.
func (swfsuo *SurveyWiFiScanUpdateOne) SetLocationID(id string) *SurveyWiFiScanUpdateOne {
	if swfsuo.location == nil {
		swfsuo.location = make(map[string]struct{})
	}
	swfsuo.location[id] = struct{}{}
	return swfsuo
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (swfsuo *SurveyWiFiScanUpdateOne) SetNillableLocationID(id *string) *SurveyWiFiScanUpdateOne {
	if id != nil {
		swfsuo = swfsuo.SetLocationID(*id)
	}
	return swfsuo
}

// SetLocation sets the location edge to Location.
func (swfsuo *SurveyWiFiScanUpdateOne) SetLocation(l *Location) *SurveyWiFiScanUpdateOne {
	return swfsuo.SetLocationID(l.ID)
}

// ClearSurveyQuestion clears the survey_question edge to SurveyQuestion.
func (swfsuo *SurveyWiFiScanUpdateOne) ClearSurveyQuestion() *SurveyWiFiScanUpdateOne {
	swfsuo.clearedSurveyQuestion = true
	return swfsuo
}

// ClearLocation clears the location edge to Location.
func (swfsuo *SurveyWiFiScanUpdateOne) ClearLocation() *SurveyWiFiScanUpdateOne {
	swfsuo.clearedLocation = true
	return swfsuo
}

// Save executes the query and returns the updated entity.
func (swfsuo *SurveyWiFiScanUpdateOne) Save(ctx context.Context) (*SurveyWiFiScan, error) {
	if swfsuo.update_time == nil {
		v := surveywifiscan.UpdateDefaultUpdateTime()
		swfsuo.update_time = &v
	}
	if len(swfsuo.survey_question) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"survey_question\"")
	}
	if len(swfsuo.location) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	return swfsuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (swfsuo *SurveyWiFiScanUpdateOne) SaveX(ctx context.Context) *SurveyWiFiScan {
	swfs, err := swfsuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return swfs
}

// Exec executes the query on the entity.
func (swfsuo *SurveyWiFiScanUpdateOne) Exec(ctx context.Context) error {
	_, err := swfsuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (swfsuo *SurveyWiFiScanUpdateOne) ExecX(ctx context.Context) {
	if err := swfsuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (swfsuo *SurveyWiFiScanUpdateOne) sqlSave(ctx context.Context) (swfs *SurveyWiFiScan, err error) {
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveywifiscan.Table,
			Columns: surveywifiscan.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  swfsuo.id,
				Type:   field.TypeString,
				Column: surveywifiscan.FieldID,
			},
		},
	}
	if value := swfsuo.update_time; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveywifiscan.FieldUpdateTime,
		})
	}
	if value := swfsuo.ssid; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveywifiscan.FieldSsid,
		})
	}
	if swfsuo.clearssid {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveywifiscan.FieldSsid,
		})
	}
	if value := swfsuo.bssid; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveywifiscan.FieldBssid,
		})
	}
	if value := swfsuo.timestamp; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveywifiscan.FieldTimestamp,
		})
	}
	if value := swfsuo.frequency; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveywifiscan.FieldFrequency,
		})
	}
	if value := swfsuo.addfrequency; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveywifiscan.FieldFrequency,
		})
	}
	if value := swfsuo.channel; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveywifiscan.FieldChannel,
		})
	}
	if value := swfsuo.addchannel; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveywifiscan.FieldChannel,
		})
	}
	if value := swfsuo.band; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveywifiscan.FieldBand,
		})
	}
	if swfsuo.clearband {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveywifiscan.FieldBand,
		})
	}
	if value := swfsuo.channel_width; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveywifiscan.FieldChannelWidth,
		})
	}
	if value := swfsuo.addchannel_width; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveywifiscan.FieldChannelWidth,
		})
	}
	if swfsuo.clearchannel_width {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveywifiscan.FieldChannelWidth,
		})
	}
	if value := swfsuo.capabilities; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveywifiscan.FieldCapabilities,
		})
	}
	if swfsuo.clearcapabilities {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveywifiscan.FieldCapabilities,
		})
	}
	if value := swfsuo.strength; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveywifiscan.FieldStrength,
		})
	}
	if value := swfsuo.addstrength; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveywifiscan.FieldStrength,
		})
	}
	if value := swfsuo.latitude; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveywifiscan.FieldLatitude,
		})
	}
	if value := swfsuo.addlatitude; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveywifiscan.FieldLatitude,
		})
	}
	if swfsuo.clearlatitude {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveywifiscan.FieldLatitude,
		})
	}
	if value := swfsuo.longitude; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveywifiscan.FieldLongitude,
		})
	}
	if value := swfsuo.addlongitude; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveywifiscan.FieldLongitude,
		})
	}
	if swfsuo.clearlongitude {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveywifiscan.FieldLongitude,
		})
	}
	if swfsuo.clearedSurveyQuestion {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveywifiscan.SurveyQuestionTable,
			Columns: []string{surveywifiscan.SurveyQuestionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveyquestion.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := swfsuo.survey_question; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveywifiscan.SurveyQuestionTable,
			Columns: []string{surveywifiscan.SurveyQuestionColumn},
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
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if swfsuo.clearedLocation {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveywifiscan.LocationTable,
			Columns: []string{surveywifiscan.LocationColumn},
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
	if nodes := swfsuo.location; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveywifiscan.LocationTable,
			Columns: []string{surveywifiscan.LocationColumn},
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
	swfs = &SurveyWiFiScan{config: swfsuo.config}
	spec.Assign = swfs.assignValues
	spec.ScanValues = swfs.scanValues()
	if err = sqlgraph.UpdateNode(ctx, swfsuo.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return swfs, nil
}
