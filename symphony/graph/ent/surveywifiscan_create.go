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
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
)

// SurveyWiFiScanCreate is the builder for creating a SurveyWiFiScan entity.
type SurveyWiFiScanCreate struct {
	config
	mutation *SurveyWiFiScanMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (swfsc *SurveyWiFiScanCreate) SetCreateTime(t time.Time) *SurveyWiFiScanCreate {
	swfsc.mutation.SetCreateTime(t)
	return swfsc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableCreateTime(t *time.Time) *SurveyWiFiScanCreate {
	if t != nil {
		swfsc.SetCreateTime(*t)
	}
	return swfsc
}

// SetUpdateTime sets the update_time field.
func (swfsc *SurveyWiFiScanCreate) SetUpdateTime(t time.Time) *SurveyWiFiScanCreate {
	swfsc.mutation.SetUpdateTime(t)
	return swfsc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableUpdateTime(t *time.Time) *SurveyWiFiScanCreate {
	if t != nil {
		swfsc.SetUpdateTime(*t)
	}
	return swfsc
}

// SetSsid sets the ssid field.
func (swfsc *SurveyWiFiScanCreate) SetSsid(s string) *SurveyWiFiScanCreate {
	swfsc.mutation.SetSsid(s)
	return swfsc
}

// SetNillableSsid sets the ssid field if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableSsid(s *string) *SurveyWiFiScanCreate {
	if s != nil {
		swfsc.SetSsid(*s)
	}
	return swfsc
}

// SetBssid sets the bssid field.
func (swfsc *SurveyWiFiScanCreate) SetBssid(s string) *SurveyWiFiScanCreate {
	swfsc.mutation.SetBssid(s)
	return swfsc
}

// SetTimestamp sets the timestamp field.
func (swfsc *SurveyWiFiScanCreate) SetTimestamp(t time.Time) *SurveyWiFiScanCreate {
	swfsc.mutation.SetTimestamp(t)
	return swfsc
}

// SetFrequency sets the frequency field.
func (swfsc *SurveyWiFiScanCreate) SetFrequency(i int) *SurveyWiFiScanCreate {
	swfsc.mutation.SetFrequency(i)
	return swfsc
}

// SetChannel sets the channel field.
func (swfsc *SurveyWiFiScanCreate) SetChannel(i int) *SurveyWiFiScanCreate {
	swfsc.mutation.SetChannel(i)
	return swfsc
}

// SetBand sets the band field.
func (swfsc *SurveyWiFiScanCreate) SetBand(s string) *SurveyWiFiScanCreate {
	swfsc.mutation.SetBand(s)
	return swfsc
}

// SetNillableBand sets the band field if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableBand(s *string) *SurveyWiFiScanCreate {
	if s != nil {
		swfsc.SetBand(*s)
	}
	return swfsc
}

// SetChannelWidth sets the channel_width field.
func (swfsc *SurveyWiFiScanCreate) SetChannelWidth(i int) *SurveyWiFiScanCreate {
	swfsc.mutation.SetChannelWidth(i)
	return swfsc
}

// SetNillableChannelWidth sets the channel_width field if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableChannelWidth(i *int) *SurveyWiFiScanCreate {
	if i != nil {
		swfsc.SetChannelWidth(*i)
	}
	return swfsc
}

// SetCapabilities sets the capabilities field.
func (swfsc *SurveyWiFiScanCreate) SetCapabilities(s string) *SurveyWiFiScanCreate {
	swfsc.mutation.SetCapabilities(s)
	return swfsc
}

// SetNillableCapabilities sets the capabilities field if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableCapabilities(s *string) *SurveyWiFiScanCreate {
	if s != nil {
		swfsc.SetCapabilities(*s)
	}
	return swfsc
}

// SetStrength sets the strength field.
func (swfsc *SurveyWiFiScanCreate) SetStrength(i int) *SurveyWiFiScanCreate {
	swfsc.mutation.SetStrength(i)
	return swfsc
}

// SetLatitude sets the latitude field.
func (swfsc *SurveyWiFiScanCreate) SetLatitude(f float64) *SurveyWiFiScanCreate {
	swfsc.mutation.SetLatitude(f)
	return swfsc
}

// SetNillableLatitude sets the latitude field if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableLatitude(f *float64) *SurveyWiFiScanCreate {
	if f != nil {
		swfsc.SetLatitude(*f)
	}
	return swfsc
}

// SetLongitude sets the longitude field.
func (swfsc *SurveyWiFiScanCreate) SetLongitude(f float64) *SurveyWiFiScanCreate {
	swfsc.mutation.SetLongitude(f)
	return swfsc
}

// SetNillableLongitude sets the longitude field if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableLongitude(f *float64) *SurveyWiFiScanCreate {
	if f != nil {
		swfsc.SetLongitude(*f)
	}
	return swfsc
}

// SetChecklistItemID sets the checklist_item edge to CheckListItem by id.
func (swfsc *SurveyWiFiScanCreate) SetChecklistItemID(id int) *SurveyWiFiScanCreate {
	swfsc.mutation.SetChecklistItemID(id)
	return swfsc
}

// SetNillableChecklistItemID sets the checklist_item edge to CheckListItem by id if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableChecklistItemID(id *int) *SurveyWiFiScanCreate {
	if id != nil {
		swfsc = swfsc.SetChecklistItemID(*id)
	}
	return swfsc
}

// SetChecklistItem sets the checklist_item edge to CheckListItem.
func (swfsc *SurveyWiFiScanCreate) SetChecklistItem(c *CheckListItem) *SurveyWiFiScanCreate {
	return swfsc.SetChecklistItemID(c.ID)
}

// SetSurveyQuestionID sets the survey_question edge to SurveyQuestion by id.
func (swfsc *SurveyWiFiScanCreate) SetSurveyQuestionID(id int) *SurveyWiFiScanCreate {
	swfsc.mutation.SetSurveyQuestionID(id)
	return swfsc
}

// SetNillableSurveyQuestionID sets the survey_question edge to SurveyQuestion by id if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableSurveyQuestionID(id *int) *SurveyWiFiScanCreate {
	if id != nil {
		swfsc = swfsc.SetSurveyQuestionID(*id)
	}
	return swfsc
}

// SetSurveyQuestion sets the survey_question edge to SurveyQuestion.
func (swfsc *SurveyWiFiScanCreate) SetSurveyQuestion(s *SurveyQuestion) *SurveyWiFiScanCreate {
	return swfsc.SetSurveyQuestionID(s.ID)
}

// SetLocationID sets the location edge to Location by id.
func (swfsc *SurveyWiFiScanCreate) SetLocationID(id int) *SurveyWiFiScanCreate {
	swfsc.mutation.SetLocationID(id)
	return swfsc
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableLocationID(id *int) *SurveyWiFiScanCreate {
	if id != nil {
		swfsc = swfsc.SetLocationID(*id)
	}
	return swfsc
}

// SetLocation sets the location edge to Location.
func (swfsc *SurveyWiFiScanCreate) SetLocation(l *Location) *SurveyWiFiScanCreate {
	return swfsc.SetLocationID(l.ID)
}

// Save creates the SurveyWiFiScan in the database.
func (swfsc *SurveyWiFiScanCreate) Save(ctx context.Context) (*SurveyWiFiScan, error) {
	if _, ok := swfsc.mutation.CreateTime(); !ok {
		v := surveywifiscan.DefaultCreateTime()
		swfsc.mutation.SetCreateTime(v)
	}
	if _, ok := swfsc.mutation.UpdateTime(); !ok {
		v := surveywifiscan.DefaultUpdateTime()
		swfsc.mutation.SetUpdateTime(v)
	}
	if _, ok := swfsc.mutation.Bssid(); !ok {
		return nil, errors.New("ent: missing required field \"bssid\"")
	}
	if _, ok := swfsc.mutation.Timestamp(); !ok {
		return nil, errors.New("ent: missing required field \"timestamp\"")
	}
	if _, ok := swfsc.mutation.Frequency(); !ok {
		return nil, errors.New("ent: missing required field \"frequency\"")
	}
	if _, ok := swfsc.mutation.Channel(); !ok {
		return nil, errors.New("ent: missing required field \"channel\"")
	}
	if _, ok := swfsc.mutation.Strength(); !ok {
		return nil, errors.New("ent: missing required field \"strength\"")
	}
	var (
		err  error
		node *SurveyWiFiScan
	)
	if len(swfsc.hooks) == 0 {
		node, err = swfsc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyWiFiScanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			swfsc.mutation = mutation
			node, err = swfsc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(swfsc.hooks) - 1; i >= 0; i-- {
			mut = swfsc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, swfsc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (swfsc *SurveyWiFiScanCreate) SaveX(ctx context.Context) *SurveyWiFiScan {
	v, err := swfsc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (swfsc *SurveyWiFiScanCreate) sqlSave(ctx context.Context) (*SurveyWiFiScan, error) {
	var (
		swfs  = &SurveyWiFiScan{config: swfsc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: surveywifiscan.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: surveywifiscan.FieldID,
			},
		}
	)
	if value, ok := swfsc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveywifiscan.FieldCreateTime,
		})
		swfs.CreateTime = value
	}
	if value, ok := swfsc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveywifiscan.FieldUpdateTime,
		})
		swfs.UpdateTime = value
	}
	if value, ok := swfsc.mutation.Ssid(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveywifiscan.FieldSsid,
		})
		swfs.Ssid = value
	}
	if value, ok := swfsc.mutation.Bssid(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveywifiscan.FieldBssid,
		})
		swfs.Bssid = value
	}
	if value, ok := swfsc.mutation.Timestamp(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveywifiscan.FieldTimestamp,
		})
		swfs.Timestamp = value
	}
	if value, ok := swfsc.mutation.Frequency(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldFrequency,
		})
		swfs.Frequency = value
	}
	if value, ok := swfsc.mutation.Channel(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldChannel,
		})
		swfs.Channel = value
	}
	if value, ok := swfsc.mutation.Band(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveywifiscan.FieldBand,
		})
		swfs.Band = value
	}
	if value, ok := swfsc.mutation.ChannelWidth(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldChannelWidth,
		})
		swfs.ChannelWidth = value
	}
	if value, ok := swfsc.mutation.Capabilities(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveywifiscan.FieldCapabilities,
		})
		swfs.Capabilities = value
	}
	if value, ok := swfsc.mutation.Strength(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveywifiscan.FieldStrength,
		})
		swfs.Strength = value
	}
	if value, ok := swfsc.mutation.Latitude(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveywifiscan.FieldLatitude,
		})
		swfs.Latitude = value
	}
	if value, ok := swfsc.mutation.Longitude(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveywifiscan.FieldLongitude,
		})
		swfs.Longitude = value
	}
	if nodes := swfsc.mutation.ChecklistItemIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveywifiscan.ChecklistItemTable,
			Columns: []string{surveywifiscan.ChecklistItemColumn},
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
	if nodes := swfsc.mutation.SurveyQuestionIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := swfsc.mutation.LocationIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, swfsc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	swfs.ID = int(id)
	return swfs, nil
}
