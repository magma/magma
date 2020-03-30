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
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
)

// SurveyQuestionCreate is the builder for creating a SurveyQuestion entity.
type SurveyQuestionCreate struct {
	config
	mutation *SurveyQuestionMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (sqc *SurveyQuestionCreate) SetCreateTime(t time.Time) *SurveyQuestionCreate {
	sqc.mutation.SetCreateTime(t)
	return sqc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillableCreateTime(t *time.Time) *SurveyQuestionCreate {
	if t != nil {
		sqc.SetCreateTime(*t)
	}
	return sqc
}

// SetUpdateTime sets the update_time field.
func (sqc *SurveyQuestionCreate) SetUpdateTime(t time.Time) *SurveyQuestionCreate {
	sqc.mutation.SetUpdateTime(t)
	return sqc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillableUpdateTime(t *time.Time) *SurveyQuestionCreate {
	if t != nil {
		sqc.SetUpdateTime(*t)
	}
	return sqc
}

// SetFormName sets the form_name field.
func (sqc *SurveyQuestionCreate) SetFormName(s string) *SurveyQuestionCreate {
	sqc.mutation.SetFormName(s)
	return sqc
}

// SetNillableFormName sets the form_name field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillableFormName(s *string) *SurveyQuestionCreate {
	if s != nil {
		sqc.SetFormName(*s)
	}
	return sqc
}

// SetFormDescription sets the form_description field.
func (sqc *SurveyQuestionCreate) SetFormDescription(s string) *SurveyQuestionCreate {
	sqc.mutation.SetFormDescription(s)
	return sqc
}

// SetNillableFormDescription sets the form_description field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillableFormDescription(s *string) *SurveyQuestionCreate {
	if s != nil {
		sqc.SetFormDescription(*s)
	}
	return sqc
}

// SetFormIndex sets the form_index field.
func (sqc *SurveyQuestionCreate) SetFormIndex(i int) *SurveyQuestionCreate {
	sqc.mutation.SetFormIndex(i)
	return sqc
}

// SetQuestionType sets the question_type field.
func (sqc *SurveyQuestionCreate) SetQuestionType(s string) *SurveyQuestionCreate {
	sqc.mutation.SetQuestionType(s)
	return sqc
}

// SetNillableQuestionType sets the question_type field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillableQuestionType(s *string) *SurveyQuestionCreate {
	if s != nil {
		sqc.SetQuestionType(*s)
	}
	return sqc
}

// SetQuestionFormat sets the question_format field.
func (sqc *SurveyQuestionCreate) SetQuestionFormat(s string) *SurveyQuestionCreate {
	sqc.mutation.SetQuestionFormat(s)
	return sqc
}

// SetNillableQuestionFormat sets the question_format field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillableQuestionFormat(s *string) *SurveyQuestionCreate {
	if s != nil {
		sqc.SetQuestionFormat(*s)
	}
	return sqc
}

// SetQuestionText sets the question_text field.
func (sqc *SurveyQuestionCreate) SetQuestionText(s string) *SurveyQuestionCreate {
	sqc.mutation.SetQuestionText(s)
	return sqc
}

// SetNillableQuestionText sets the question_text field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillableQuestionText(s *string) *SurveyQuestionCreate {
	if s != nil {
		sqc.SetQuestionText(*s)
	}
	return sqc
}

// SetQuestionIndex sets the question_index field.
func (sqc *SurveyQuestionCreate) SetQuestionIndex(i int) *SurveyQuestionCreate {
	sqc.mutation.SetQuestionIndex(i)
	return sqc
}

// SetBoolData sets the bool_data field.
func (sqc *SurveyQuestionCreate) SetBoolData(b bool) *SurveyQuestionCreate {
	sqc.mutation.SetBoolData(b)
	return sqc
}

// SetNillableBoolData sets the bool_data field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillableBoolData(b *bool) *SurveyQuestionCreate {
	if b != nil {
		sqc.SetBoolData(*b)
	}
	return sqc
}

// SetEmailData sets the email_data field.
func (sqc *SurveyQuestionCreate) SetEmailData(s string) *SurveyQuestionCreate {
	sqc.mutation.SetEmailData(s)
	return sqc
}

// SetNillableEmailData sets the email_data field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillableEmailData(s *string) *SurveyQuestionCreate {
	if s != nil {
		sqc.SetEmailData(*s)
	}
	return sqc
}

// SetLatitude sets the latitude field.
func (sqc *SurveyQuestionCreate) SetLatitude(f float64) *SurveyQuestionCreate {
	sqc.mutation.SetLatitude(f)
	return sqc
}

// SetNillableLatitude sets the latitude field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillableLatitude(f *float64) *SurveyQuestionCreate {
	if f != nil {
		sqc.SetLatitude(*f)
	}
	return sqc
}

// SetLongitude sets the longitude field.
func (sqc *SurveyQuestionCreate) SetLongitude(f float64) *SurveyQuestionCreate {
	sqc.mutation.SetLongitude(f)
	return sqc
}

// SetNillableLongitude sets the longitude field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillableLongitude(f *float64) *SurveyQuestionCreate {
	if f != nil {
		sqc.SetLongitude(*f)
	}
	return sqc
}

// SetLocationAccuracy sets the location_accuracy field.
func (sqc *SurveyQuestionCreate) SetLocationAccuracy(f float64) *SurveyQuestionCreate {
	sqc.mutation.SetLocationAccuracy(f)
	return sqc
}

// SetNillableLocationAccuracy sets the location_accuracy field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillableLocationAccuracy(f *float64) *SurveyQuestionCreate {
	if f != nil {
		sqc.SetLocationAccuracy(*f)
	}
	return sqc
}

// SetAltitude sets the altitude field.
func (sqc *SurveyQuestionCreate) SetAltitude(f float64) *SurveyQuestionCreate {
	sqc.mutation.SetAltitude(f)
	return sqc
}

// SetNillableAltitude sets the altitude field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillableAltitude(f *float64) *SurveyQuestionCreate {
	if f != nil {
		sqc.SetAltitude(*f)
	}
	return sqc
}

// SetPhoneData sets the phone_data field.
func (sqc *SurveyQuestionCreate) SetPhoneData(s string) *SurveyQuestionCreate {
	sqc.mutation.SetPhoneData(s)
	return sqc
}

// SetNillablePhoneData sets the phone_data field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillablePhoneData(s *string) *SurveyQuestionCreate {
	if s != nil {
		sqc.SetPhoneData(*s)
	}
	return sqc
}

// SetTextData sets the text_data field.
func (sqc *SurveyQuestionCreate) SetTextData(s string) *SurveyQuestionCreate {
	sqc.mutation.SetTextData(s)
	return sqc
}

// SetNillableTextData sets the text_data field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillableTextData(s *string) *SurveyQuestionCreate {
	if s != nil {
		sqc.SetTextData(*s)
	}
	return sqc
}

// SetFloatData sets the float_data field.
func (sqc *SurveyQuestionCreate) SetFloatData(f float64) *SurveyQuestionCreate {
	sqc.mutation.SetFloatData(f)
	return sqc
}

// SetNillableFloatData sets the float_data field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillableFloatData(f *float64) *SurveyQuestionCreate {
	if f != nil {
		sqc.SetFloatData(*f)
	}
	return sqc
}

// SetIntData sets the int_data field.
func (sqc *SurveyQuestionCreate) SetIntData(i int) *SurveyQuestionCreate {
	sqc.mutation.SetIntData(i)
	return sqc
}

// SetNillableIntData sets the int_data field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillableIntData(i *int) *SurveyQuestionCreate {
	if i != nil {
		sqc.SetIntData(*i)
	}
	return sqc
}

// SetDateData sets the date_data field.
func (sqc *SurveyQuestionCreate) SetDateData(t time.Time) *SurveyQuestionCreate {
	sqc.mutation.SetDateData(t)
	return sqc
}

// SetNillableDateData sets the date_data field if the given value is not nil.
func (sqc *SurveyQuestionCreate) SetNillableDateData(t *time.Time) *SurveyQuestionCreate {
	if t != nil {
		sqc.SetDateData(*t)
	}
	return sqc
}

// SetSurveyID sets the survey edge to Survey by id.
func (sqc *SurveyQuestionCreate) SetSurveyID(id int) *SurveyQuestionCreate {
	sqc.mutation.SetSurveyID(id)
	return sqc
}

// SetSurvey sets the survey edge to Survey.
func (sqc *SurveyQuestionCreate) SetSurvey(s *Survey) *SurveyQuestionCreate {
	return sqc.SetSurveyID(s.ID)
}

// AddWifiScanIDs adds the wifi_scan edge to SurveyWiFiScan by ids.
func (sqc *SurveyQuestionCreate) AddWifiScanIDs(ids ...int) *SurveyQuestionCreate {
	sqc.mutation.AddWifiScanIDs(ids...)
	return sqc
}

// AddWifiScan adds the wifi_scan edges to SurveyWiFiScan.
func (sqc *SurveyQuestionCreate) AddWifiScan(s ...*SurveyWiFiScan) *SurveyQuestionCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sqc.AddWifiScanIDs(ids...)
}

// AddCellScanIDs adds the cell_scan edge to SurveyCellScan by ids.
func (sqc *SurveyQuestionCreate) AddCellScanIDs(ids ...int) *SurveyQuestionCreate {
	sqc.mutation.AddCellScanIDs(ids...)
	return sqc
}

// AddCellScan adds the cell_scan edges to SurveyCellScan.
func (sqc *SurveyQuestionCreate) AddCellScan(s ...*SurveyCellScan) *SurveyQuestionCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sqc.AddCellScanIDs(ids...)
}

// AddPhotoDatumIDs adds the photo_data edge to File by ids.
func (sqc *SurveyQuestionCreate) AddPhotoDatumIDs(ids ...int) *SurveyQuestionCreate {
	sqc.mutation.AddPhotoDatumIDs(ids...)
	return sqc
}

// AddPhotoData adds the photo_data edges to File.
func (sqc *SurveyQuestionCreate) AddPhotoData(f ...*File) *SurveyQuestionCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return sqc.AddPhotoDatumIDs(ids...)
}

// Save creates the SurveyQuestion in the database.
func (sqc *SurveyQuestionCreate) Save(ctx context.Context) (*SurveyQuestion, error) {
	if _, ok := sqc.mutation.CreateTime(); !ok {
		v := surveyquestion.DefaultCreateTime()
		sqc.mutation.SetCreateTime(v)
	}
	if _, ok := sqc.mutation.UpdateTime(); !ok {
		v := surveyquestion.DefaultUpdateTime()
		sqc.mutation.SetUpdateTime(v)
	}
	if _, ok := sqc.mutation.FormIndex(); !ok {
		return nil, errors.New("ent: missing required field \"form_index\"")
	}
	if _, ok := sqc.mutation.QuestionIndex(); !ok {
		return nil, errors.New("ent: missing required field \"question_index\"")
	}
	if _, ok := sqc.mutation.SurveyID(); !ok {
		return nil, errors.New("ent: missing required edge \"survey\"")
	}
	var (
		err  error
		node *SurveyQuestion
	)
	if len(sqc.hooks) == 0 {
		node, err = sqc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyQuestionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			sqc.mutation = mutation
			node, err = sqc.sqlSave(ctx)
			return node, err
		})
		for i := len(sqc.hooks) - 1; i >= 0; i-- {
			mut = sqc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, sqc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (sqc *SurveyQuestionCreate) SaveX(ctx context.Context) *SurveyQuestion {
	v, err := sqc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (sqc *SurveyQuestionCreate) sqlSave(ctx context.Context) (*SurveyQuestion, error) {
	var (
		sq    = &SurveyQuestion{config: sqc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: surveyquestion.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: surveyquestion.FieldID,
			},
		}
	)
	if value, ok := sqc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveyquestion.FieldCreateTime,
		})
		sq.CreateTime = value
	}
	if value, ok := sqc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveyquestion.FieldUpdateTime,
		})
		sq.UpdateTime = value
	}
	if value, ok := sqc.mutation.FormName(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldFormName,
		})
		sq.FormName = value
	}
	if value, ok := sqc.mutation.FormDescription(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldFormDescription,
		})
		sq.FormDescription = value
	}
	if value, ok := sqc.mutation.FormIndex(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveyquestion.FieldFormIndex,
		})
		sq.FormIndex = value
	}
	if value, ok := sqc.mutation.QuestionType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldQuestionType,
		})
		sq.QuestionType = value
	}
	if value, ok := sqc.mutation.QuestionFormat(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldQuestionFormat,
		})
		sq.QuestionFormat = value
	}
	if value, ok := sqc.mutation.QuestionText(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldQuestionText,
		})
		sq.QuestionText = value
	}
	if value, ok := sqc.mutation.QuestionIndex(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveyquestion.FieldQuestionIndex,
		})
		sq.QuestionIndex = value
	}
	if value, ok := sqc.mutation.BoolData(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: surveyquestion.FieldBoolData,
		})
		sq.BoolData = value
	}
	if value, ok := sqc.mutation.EmailData(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldEmailData,
		})
		sq.EmailData = value
	}
	if value, ok := sqc.mutation.Latitude(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldLatitude,
		})
		sq.Latitude = value
	}
	if value, ok := sqc.mutation.Longitude(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldLongitude,
		})
		sq.Longitude = value
	}
	if value, ok := sqc.mutation.LocationAccuracy(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldLocationAccuracy,
		})
		sq.LocationAccuracy = value
	}
	if value, ok := sqc.mutation.Altitude(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldAltitude,
		})
		sq.Altitude = value
	}
	if value, ok := sqc.mutation.PhoneData(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldPhoneData,
		})
		sq.PhoneData = value
	}
	if value, ok := sqc.mutation.TextData(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldTextData,
		})
		sq.TextData = value
	}
	if value, ok := sqc.mutation.FloatData(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldFloatData,
		})
		sq.FloatData = value
	}
	if value, ok := sqc.mutation.IntData(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveyquestion.FieldIntData,
		})
		sq.IntData = value
	}
	if value, ok := sqc.mutation.DateData(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveyquestion.FieldDateData,
		})
		sq.DateData = value
	}
	if nodes := sqc.mutation.SurveyIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveyquestion.SurveyTable,
			Columns: []string{surveyquestion.SurveyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: survey.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := sqc.mutation.WifiScanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   surveyquestion.WifiScanTable,
			Columns: []string{surveyquestion.WifiScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveywifiscan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := sqc.mutation.CellScanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   surveyquestion.CellScanTable,
			Columns: []string{surveyquestion.CellScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveycellscan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := sqc.mutation.PhotoDataIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   surveyquestion.PhotoDataTable,
			Columns: []string{surveyquestion.PhotoDataColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: file.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, sqc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	sq.ID = int(id)
	return sq, nil
}
