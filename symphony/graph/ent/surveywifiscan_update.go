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
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
)

// SurveyWiFiScanUpdate is the builder for updating SurveyWiFiScan entities.
type SurveyWiFiScanUpdate struct {
	config
	hooks      []Hook
	mutation   *SurveyWiFiScanMutation
	predicates []predicate.SurveyWiFiScan
}

// Where adds a new predicate for the builder.
func (swfsu *SurveyWiFiScanUpdate) Where(ps ...predicate.SurveyWiFiScan) *SurveyWiFiScanUpdate {
	swfsu.predicates = append(swfsu.predicates, ps...)
	return swfsu
}

// SetSsid sets the ssid field.
func (swfsu *SurveyWiFiScanUpdate) SetSsid(s string) *SurveyWiFiScanUpdate {
	swfsu.mutation.SetSsid(s)
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
	swfsu.mutation.ClearSsid()
	return swfsu
}

// SetBssid sets the bssid field.
func (swfsu *SurveyWiFiScanUpdate) SetBssid(s string) *SurveyWiFiScanUpdate {
	swfsu.mutation.SetBssid(s)
	return swfsu
}

// SetTimestamp sets the timestamp field.
func (swfsu *SurveyWiFiScanUpdate) SetTimestamp(t time.Time) *SurveyWiFiScanUpdate {
	swfsu.mutation.SetTimestamp(t)
	return swfsu
}

// SetFrequency sets the frequency field.
func (swfsu *SurveyWiFiScanUpdate) SetFrequency(i int) *SurveyWiFiScanUpdate {
	swfsu.mutation.ResetFrequency()
	swfsu.mutation.SetFrequency(i)
	return swfsu
}

// AddFrequency adds i to frequency.
func (swfsu *SurveyWiFiScanUpdate) AddFrequency(i int) *SurveyWiFiScanUpdate {
	swfsu.mutation.AddFrequency(i)
	return swfsu
}

// SetChannel sets the channel field.
func (swfsu *SurveyWiFiScanUpdate) SetChannel(i int) *SurveyWiFiScanUpdate {
	swfsu.mutation.ResetChannel()
	swfsu.mutation.SetChannel(i)
	return swfsu
}

// AddChannel adds i to channel.
func (swfsu *SurveyWiFiScanUpdate) AddChannel(i int) *SurveyWiFiScanUpdate {
	swfsu.mutation.AddChannel(i)
	return swfsu
}

// SetBand sets the band field.
func (swfsu *SurveyWiFiScanUpdate) SetBand(s string) *SurveyWiFiScanUpdate {
	swfsu.mutation.SetBand(s)
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
	swfsu.mutation.ClearBand()
	return swfsu
}

// SetChannelWidth sets the channel_width field.
func (swfsu *SurveyWiFiScanUpdate) SetChannelWidth(i int) *SurveyWiFiScanUpdate {
	swfsu.mutation.ResetChannelWidth()
	swfsu.mutation.SetChannelWidth(i)
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
	swfsu.mutation.AddChannelWidth(i)
	return swfsu
}

// ClearChannelWidth clears the value of channel_width.
func (swfsu *SurveyWiFiScanUpdate) ClearChannelWidth() *SurveyWiFiScanUpdate {
	swfsu.mutation.ClearChannelWidth()
	return swfsu
}

// SetCapabilities sets the capabilities field.
func (swfsu *SurveyWiFiScanUpdate) SetCapabilities(s string) *SurveyWiFiScanUpdate {
	swfsu.mutation.SetCapabilities(s)
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
	swfsu.mutation.ClearCapabilities()
	return swfsu
}

// SetStrength sets the strength field.
func (swfsu *SurveyWiFiScanUpdate) SetStrength(i int) *SurveyWiFiScanUpdate {
	swfsu.mutation.ResetStrength()
	swfsu.mutation.SetStrength(i)
	return swfsu
}

// AddStrength adds i to strength.
func (swfsu *SurveyWiFiScanUpdate) AddStrength(i int) *SurveyWiFiScanUpdate {
	swfsu.mutation.AddStrength(i)
	return swfsu
}

// SetLatitude sets the latitude field.
func (swfsu *SurveyWiFiScanUpdate) SetLatitude(f float64) *SurveyWiFiScanUpdate {
	swfsu.mutation.ResetLatitude()
	swfsu.mutation.SetLatitude(f)
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
	swfsu.mutation.AddLatitude(f)
	return swfsu
}

// ClearLatitude clears the value of latitude.
func (swfsu *SurveyWiFiScanUpdate) ClearLatitude() *SurveyWiFiScanUpdate {
	swfsu.mutation.ClearLatitude()
	return swfsu
}

// SetLongitude sets the longitude field.
func (swfsu *SurveyWiFiScanUpdate) SetLongitude(f float64) *SurveyWiFiScanUpdate {
	swfsu.mutation.ResetLongitude()
	swfsu.mutation.SetLongitude(f)
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
	swfsu.mutation.AddLongitude(f)
	return swfsu
}

// ClearLongitude clears the value of longitude.
func (swfsu *SurveyWiFiScanUpdate) ClearLongitude() *SurveyWiFiScanUpdate {
	swfsu.mutation.ClearLongitude()
	return swfsu
}

// SetSurveyQuestionID sets the survey_question edge to SurveyQuestion by id.
func (swfsu *SurveyWiFiScanUpdate) SetSurveyQuestionID(id int) *SurveyWiFiScanUpdate {
	swfsu.mutation.SetSurveyQuestionID(id)
	return swfsu
}

// SetNillableSurveyQuestionID sets the survey_question edge to SurveyQuestion by id if the given value is not nil.
func (swfsu *SurveyWiFiScanUpdate) SetNillableSurveyQuestionID(id *int) *SurveyWiFiScanUpdate {
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
func (swfsu *SurveyWiFiScanUpdate) SetLocationID(id int) *SurveyWiFiScanUpdate {
	swfsu.mutation.SetLocationID(id)
	return swfsu
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (swfsu *SurveyWiFiScanUpdate) SetNillableLocationID(id *int) *SurveyWiFiScanUpdate {
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
	swfsu.mutation.ClearSurveyQuestion()
	return swfsu
}

// ClearLocation clears the location edge to Location.
func (swfsu *SurveyWiFiScanUpdate) ClearLocation() *SurveyWiFiScanUpdate {
	swfsu.mutation.ClearLocation()
	return swfsu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (swfsu *SurveyWiFiScanUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := swfsu.mutation.UpdateTime(); !ok {
		v := surveywifiscan.UpdateDefaultUpdateTime()
		swfsu.mutation.SetUpdateTime(v)
	}

	var (
		err      error
		affected int
	)
	if len(swfsu.hooks) == 0 {
		affected, err = swfsu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyWiFiScanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			swfsu.mutation = mutation
			affected, err = swfsu.sqlSave(ctx)
			return affected, err
		})
		for i := len(swfsu.hooks) - 1; i >= 0; i-- {
			mut = swfsu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, swfsu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveywifiscan.Table,
			Columns: surveywifiscan.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: surveywifiscan.FieldID,
			},
		},
	}
	if ps := swfsu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := swfsu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveywifiscan.FieldUpdateTime,
		})
	}
	if value, ok := swfsu.mutation.Ssid(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveywifiscan.FieldSsid,
		})
	}
	if swfsu.mutation.SsidCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveywifiscan.FieldSsid,
		})
	}
	if value, ok := swfsu.mutation.Bssid(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveywifiscan.FieldBssid,
		})
	}
	if value, ok := swfsu.mutation.Timestamp(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveywifiscan.FieldTimestamp,
		})
	}
	if value, ok := swfsu.mutation.Frequency(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldFrequency,
		})
	}
	if value, ok := swfsu.mutation.AddedFrequency(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldFrequency,
		})
	}
	if value, ok := swfsu.mutation.Channel(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldChannel,
		})
	}
	if value, ok := swfsu.mutation.AddedChannel(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldChannel,
		})
	}
	if value, ok := swfsu.mutation.Band(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveywifiscan.FieldBand,
		})
	}
	if swfsu.mutation.BandCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveywifiscan.FieldBand,
		})
	}
	if value, ok := swfsu.mutation.ChannelWidth(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldChannelWidth,
		})
	}
	if value, ok := swfsu.mutation.AddedChannelWidth(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldChannelWidth,
		})
	}
	if swfsu.mutation.ChannelWidthCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveywifiscan.FieldChannelWidth,
		})
	}
	if value, ok := swfsu.mutation.Capabilities(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveywifiscan.FieldCapabilities,
		})
	}
	if swfsu.mutation.CapabilitiesCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveywifiscan.FieldCapabilities,
		})
	}
	if value, ok := swfsu.mutation.Strength(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldStrength,
		})
	}
	if value, ok := swfsu.mutation.AddedStrength(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldStrength,
		})
	}
	if value, ok := swfsu.mutation.Latitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveywifiscan.FieldLatitude,
		})
	}
	if value, ok := swfsu.mutation.AddedLatitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveywifiscan.FieldLatitude,
		})
	}
	if swfsu.mutation.LatitudeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveywifiscan.FieldLatitude,
		})
	}
	if value, ok := swfsu.mutation.Longitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveywifiscan.FieldLongitude,
		})
	}
	if value, ok := swfsu.mutation.AddedLongitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveywifiscan.FieldLongitude,
		})
	}
	if swfsu.mutation.LongitudeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveywifiscan.FieldLongitude,
		})
	}
	if swfsu.mutation.SurveyQuestionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveywifiscan.SurveyQuestionTable,
			Columns: []string{surveywifiscan.SurveyQuestionColumn},
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
	if nodes := swfsu.mutation.SurveyQuestionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveywifiscan.SurveyQuestionTable,
			Columns: []string{surveywifiscan.SurveyQuestionColumn},
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
	if swfsu.mutation.LocationCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveywifiscan.LocationTable,
			Columns: []string{surveywifiscan.LocationColumn},
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
	if nodes := swfsu.mutation.LocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveywifiscan.LocationTable,
			Columns: []string{surveywifiscan.LocationColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, swfsu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{surveywifiscan.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// SurveyWiFiScanUpdateOne is the builder for updating a single SurveyWiFiScan entity.
type SurveyWiFiScanUpdateOne struct {
	config
	hooks    []Hook
	mutation *SurveyWiFiScanMutation
}

// SetSsid sets the ssid field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetSsid(s string) *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.SetSsid(s)
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
	swfsuo.mutation.ClearSsid()
	return swfsuo
}

// SetBssid sets the bssid field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetBssid(s string) *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.SetBssid(s)
	return swfsuo
}

// SetTimestamp sets the timestamp field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetTimestamp(t time.Time) *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.SetTimestamp(t)
	return swfsuo
}

// SetFrequency sets the frequency field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetFrequency(i int) *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.ResetFrequency()
	swfsuo.mutation.SetFrequency(i)
	return swfsuo
}

// AddFrequency adds i to frequency.
func (swfsuo *SurveyWiFiScanUpdateOne) AddFrequency(i int) *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.AddFrequency(i)
	return swfsuo
}

// SetChannel sets the channel field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetChannel(i int) *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.ResetChannel()
	swfsuo.mutation.SetChannel(i)
	return swfsuo
}

// AddChannel adds i to channel.
func (swfsuo *SurveyWiFiScanUpdateOne) AddChannel(i int) *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.AddChannel(i)
	return swfsuo
}

// SetBand sets the band field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetBand(s string) *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.SetBand(s)
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
	swfsuo.mutation.ClearBand()
	return swfsuo
}

// SetChannelWidth sets the channel_width field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetChannelWidth(i int) *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.ResetChannelWidth()
	swfsuo.mutation.SetChannelWidth(i)
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
	swfsuo.mutation.AddChannelWidth(i)
	return swfsuo
}

// ClearChannelWidth clears the value of channel_width.
func (swfsuo *SurveyWiFiScanUpdateOne) ClearChannelWidth() *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.ClearChannelWidth()
	return swfsuo
}

// SetCapabilities sets the capabilities field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetCapabilities(s string) *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.SetCapabilities(s)
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
	swfsuo.mutation.ClearCapabilities()
	return swfsuo
}

// SetStrength sets the strength field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetStrength(i int) *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.ResetStrength()
	swfsuo.mutation.SetStrength(i)
	return swfsuo
}

// AddStrength adds i to strength.
func (swfsuo *SurveyWiFiScanUpdateOne) AddStrength(i int) *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.AddStrength(i)
	return swfsuo
}

// SetLatitude sets the latitude field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetLatitude(f float64) *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.ResetLatitude()
	swfsuo.mutation.SetLatitude(f)
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
	swfsuo.mutation.AddLatitude(f)
	return swfsuo
}

// ClearLatitude clears the value of latitude.
func (swfsuo *SurveyWiFiScanUpdateOne) ClearLatitude() *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.ClearLatitude()
	return swfsuo
}

// SetLongitude sets the longitude field.
func (swfsuo *SurveyWiFiScanUpdateOne) SetLongitude(f float64) *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.ResetLongitude()
	swfsuo.mutation.SetLongitude(f)
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
	swfsuo.mutation.AddLongitude(f)
	return swfsuo
}

// ClearLongitude clears the value of longitude.
func (swfsuo *SurveyWiFiScanUpdateOne) ClearLongitude() *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.ClearLongitude()
	return swfsuo
}

// SetSurveyQuestionID sets the survey_question edge to SurveyQuestion by id.
func (swfsuo *SurveyWiFiScanUpdateOne) SetSurveyQuestionID(id int) *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.SetSurveyQuestionID(id)
	return swfsuo
}

// SetNillableSurveyQuestionID sets the survey_question edge to SurveyQuestion by id if the given value is not nil.
func (swfsuo *SurveyWiFiScanUpdateOne) SetNillableSurveyQuestionID(id *int) *SurveyWiFiScanUpdateOne {
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
func (swfsuo *SurveyWiFiScanUpdateOne) SetLocationID(id int) *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.SetLocationID(id)
	return swfsuo
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (swfsuo *SurveyWiFiScanUpdateOne) SetNillableLocationID(id *int) *SurveyWiFiScanUpdateOne {
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
	swfsuo.mutation.ClearSurveyQuestion()
	return swfsuo
}

// ClearLocation clears the location edge to Location.
func (swfsuo *SurveyWiFiScanUpdateOne) ClearLocation() *SurveyWiFiScanUpdateOne {
	swfsuo.mutation.ClearLocation()
	return swfsuo
}

// Save executes the query and returns the updated entity.
func (swfsuo *SurveyWiFiScanUpdateOne) Save(ctx context.Context) (*SurveyWiFiScan, error) {
	if _, ok := swfsuo.mutation.UpdateTime(); !ok {
		v := surveywifiscan.UpdateDefaultUpdateTime()
		swfsuo.mutation.SetUpdateTime(v)
	}

	var (
		err  error
		node *SurveyWiFiScan
	)
	if len(swfsuo.hooks) == 0 {
		node, err = swfsuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyWiFiScanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			swfsuo.mutation = mutation
			node, err = swfsuo.sqlSave(ctx)
			return node, err
		})
		for i := len(swfsuo.hooks) - 1; i >= 0; i-- {
			mut = swfsuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, swfsuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveywifiscan.Table,
			Columns: surveywifiscan.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: surveywifiscan.FieldID,
			},
		},
	}
	id, ok := swfsuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing SurveyWiFiScan.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := swfsuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveywifiscan.FieldUpdateTime,
		})
	}
	if value, ok := swfsuo.mutation.Ssid(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveywifiscan.FieldSsid,
		})
	}
	if swfsuo.mutation.SsidCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveywifiscan.FieldSsid,
		})
	}
	if value, ok := swfsuo.mutation.Bssid(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveywifiscan.FieldBssid,
		})
	}
	if value, ok := swfsuo.mutation.Timestamp(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveywifiscan.FieldTimestamp,
		})
	}
	if value, ok := swfsuo.mutation.Frequency(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldFrequency,
		})
	}
	if value, ok := swfsuo.mutation.AddedFrequency(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldFrequency,
		})
	}
	if value, ok := swfsuo.mutation.Channel(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldChannel,
		})
	}
	if value, ok := swfsuo.mutation.AddedChannel(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldChannel,
		})
	}
	if value, ok := swfsuo.mutation.Band(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveywifiscan.FieldBand,
		})
	}
	if swfsuo.mutation.BandCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveywifiscan.FieldBand,
		})
	}
	if value, ok := swfsuo.mutation.ChannelWidth(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldChannelWidth,
		})
	}
	if value, ok := swfsuo.mutation.AddedChannelWidth(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldChannelWidth,
		})
	}
	if swfsuo.mutation.ChannelWidthCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveywifiscan.FieldChannelWidth,
		})
	}
	if value, ok := swfsuo.mutation.Capabilities(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveywifiscan.FieldCapabilities,
		})
	}
	if swfsuo.mutation.CapabilitiesCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveywifiscan.FieldCapabilities,
		})
	}
	if value, ok := swfsuo.mutation.Strength(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldStrength,
		})
	}
	if value, ok := swfsuo.mutation.AddedStrength(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldStrength,
		})
	}
	if value, ok := swfsuo.mutation.Latitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveywifiscan.FieldLatitude,
		})
	}
	if value, ok := swfsuo.mutation.AddedLatitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveywifiscan.FieldLatitude,
		})
	}
	if swfsuo.mutation.LatitudeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveywifiscan.FieldLatitude,
		})
	}
	if value, ok := swfsuo.mutation.Longitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveywifiscan.FieldLongitude,
		})
	}
	if value, ok := swfsuo.mutation.AddedLongitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveywifiscan.FieldLongitude,
		})
	}
	if swfsuo.mutation.LongitudeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveywifiscan.FieldLongitude,
		})
	}
	if swfsuo.mutation.SurveyQuestionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveywifiscan.SurveyQuestionTable,
			Columns: []string{surveywifiscan.SurveyQuestionColumn},
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
	if nodes := swfsuo.mutation.SurveyQuestionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveywifiscan.SurveyQuestionTable,
			Columns: []string{surveywifiscan.SurveyQuestionColumn},
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
	if swfsuo.mutation.LocationCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveywifiscan.LocationTable,
			Columns: []string{surveywifiscan.LocationColumn},
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
	if nodes := swfsuo.mutation.LocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveywifiscan.LocationTable,
			Columns: []string{surveywifiscan.LocationColumn},
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
	swfs = &SurveyWiFiScan{config: swfsuo.config}
	_spec.Assign = swfs.assignValues
	_spec.ScanValues = swfs.scanValues()
	if err = sqlgraph.UpdateNode(ctx, swfsuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{surveywifiscan.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return swfs, nil
}
