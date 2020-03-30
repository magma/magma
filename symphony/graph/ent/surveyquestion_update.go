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

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
)

// SurveyQuestionUpdate is the builder for updating SurveyQuestion entities.
type SurveyQuestionUpdate struct {
	config
	hooks      []Hook
	mutation   *SurveyQuestionMutation
	predicates []predicate.SurveyQuestion
}

// Where adds a new predicate for the builder.
func (squ *SurveyQuestionUpdate) Where(ps ...predicate.SurveyQuestion) *SurveyQuestionUpdate {
	squ.predicates = append(squ.predicates, ps...)
	return squ
}

// SetFormName sets the form_name field.
func (squ *SurveyQuestionUpdate) SetFormName(s string) *SurveyQuestionUpdate {
	squ.mutation.SetFormName(s)
	return squ
}

// SetNillableFormName sets the form_name field if the given value is not nil.
func (squ *SurveyQuestionUpdate) SetNillableFormName(s *string) *SurveyQuestionUpdate {
	if s != nil {
		squ.SetFormName(*s)
	}
	return squ
}

// ClearFormName clears the value of form_name.
func (squ *SurveyQuestionUpdate) ClearFormName() *SurveyQuestionUpdate {
	squ.mutation.ClearFormName()
	return squ
}

// SetFormDescription sets the form_description field.
func (squ *SurveyQuestionUpdate) SetFormDescription(s string) *SurveyQuestionUpdate {
	squ.mutation.SetFormDescription(s)
	return squ
}

// SetNillableFormDescription sets the form_description field if the given value is not nil.
func (squ *SurveyQuestionUpdate) SetNillableFormDescription(s *string) *SurveyQuestionUpdate {
	if s != nil {
		squ.SetFormDescription(*s)
	}
	return squ
}

// ClearFormDescription clears the value of form_description.
func (squ *SurveyQuestionUpdate) ClearFormDescription() *SurveyQuestionUpdate {
	squ.mutation.ClearFormDescription()
	return squ
}

// SetFormIndex sets the form_index field.
func (squ *SurveyQuestionUpdate) SetFormIndex(i int) *SurveyQuestionUpdate {
	squ.mutation.ResetFormIndex()
	squ.mutation.SetFormIndex(i)
	return squ
}

// AddFormIndex adds i to form_index.
func (squ *SurveyQuestionUpdate) AddFormIndex(i int) *SurveyQuestionUpdate {
	squ.mutation.AddFormIndex(i)
	return squ
}

// SetQuestionType sets the question_type field.
func (squ *SurveyQuestionUpdate) SetQuestionType(s string) *SurveyQuestionUpdate {
	squ.mutation.SetQuestionType(s)
	return squ
}

// SetNillableQuestionType sets the question_type field if the given value is not nil.
func (squ *SurveyQuestionUpdate) SetNillableQuestionType(s *string) *SurveyQuestionUpdate {
	if s != nil {
		squ.SetQuestionType(*s)
	}
	return squ
}

// ClearQuestionType clears the value of question_type.
func (squ *SurveyQuestionUpdate) ClearQuestionType() *SurveyQuestionUpdate {
	squ.mutation.ClearQuestionType()
	return squ
}

// SetQuestionFormat sets the question_format field.
func (squ *SurveyQuestionUpdate) SetQuestionFormat(s string) *SurveyQuestionUpdate {
	squ.mutation.SetQuestionFormat(s)
	return squ
}

// SetNillableQuestionFormat sets the question_format field if the given value is not nil.
func (squ *SurveyQuestionUpdate) SetNillableQuestionFormat(s *string) *SurveyQuestionUpdate {
	if s != nil {
		squ.SetQuestionFormat(*s)
	}
	return squ
}

// ClearQuestionFormat clears the value of question_format.
func (squ *SurveyQuestionUpdate) ClearQuestionFormat() *SurveyQuestionUpdate {
	squ.mutation.ClearQuestionFormat()
	return squ
}

// SetQuestionText sets the question_text field.
func (squ *SurveyQuestionUpdate) SetQuestionText(s string) *SurveyQuestionUpdate {
	squ.mutation.SetQuestionText(s)
	return squ
}

// SetNillableQuestionText sets the question_text field if the given value is not nil.
func (squ *SurveyQuestionUpdate) SetNillableQuestionText(s *string) *SurveyQuestionUpdate {
	if s != nil {
		squ.SetQuestionText(*s)
	}
	return squ
}

// ClearQuestionText clears the value of question_text.
func (squ *SurveyQuestionUpdate) ClearQuestionText() *SurveyQuestionUpdate {
	squ.mutation.ClearQuestionText()
	return squ
}

// SetQuestionIndex sets the question_index field.
func (squ *SurveyQuestionUpdate) SetQuestionIndex(i int) *SurveyQuestionUpdate {
	squ.mutation.ResetQuestionIndex()
	squ.mutation.SetQuestionIndex(i)
	return squ
}

// AddQuestionIndex adds i to question_index.
func (squ *SurveyQuestionUpdate) AddQuestionIndex(i int) *SurveyQuestionUpdate {
	squ.mutation.AddQuestionIndex(i)
	return squ
}

// SetBoolData sets the bool_data field.
func (squ *SurveyQuestionUpdate) SetBoolData(b bool) *SurveyQuestionUpdate {
	squ.mutation.SetBoolData(b)
	return squ
}

// SetNillableBoolData sets the bool_data field if the given value is not nil.
func (squ *SurveyQuestionUpdate) SetNillableBoolData(b *bool) *SurveyQuestionUpdate {
	if b != nil {
		squ.SetBoolData(*b)
	}
	return squ
}

// ClearBoolData clears the value of bool_data.
func (squ *SurveyQuestionUpdate) ClearBoolData() *SurveyQuestionUpdate {
	squ.mutation.ClearBoolData()
	return squ
}

// SetEmailData sets the email_data field.
func (squ *SurveyQuestionUpdate) SetEmailData(s string) *SurveyQuestionUpdate {
	squ.mutation.SetEmailData(s)
	return squ
}

// SetNillableEmailData sets the email_data field if the given value is not nil.
func (squ *SurveyQuestionUpdate) SetNillableEmailData(s *string) *SurveyQuestionUpdate {
	if s != nil {
		squ.SetEmailData(*s)
	}
	return squ
}

// ClearEmailData clears the value of email_data.
func (squ *SurveyQuestionUpdate) ClearEmailData() *SurveyQuestionUpdate {
	squ.mutation.ClearEmailData()
	return squ
}

// SetLatitude sets the latitude field.
func (squ *SurveyQuestionUpdate) SetLatitude(f float64) *SurveyQuestionUpdate {
	squ.mutation.ResetLatitude()
	squ.mutation.SetLatitude(f)
	return squ
}

// SetNillableLatitude sets the latitude field if the given value is not nil.
func (squ *SurveyQuestionUpdate) SetNillableLatitude(f *float64) *SurveyQuestionUpdate {
	if f != nil {
		squ.SetLatitude(*f)
	}
	return squ
}

// AddLatitude adds f to latitude.
func (squ *SurveyQuestionUpdate) AddLatitude(f float64) *SurveyQuestionUpdate {
	squ.mutation.AddLatitude(f)
	return squ
}

// ClearLatitude clears the value of latitude.
func (squ *SurveyQuestionUpdate) ClearLatitude() *SurveyQuestionUpdate {
	squ.mutation.ClearLatitude()
	return squ
}

// SetLongitude sets the longitude field.
func (squ *SurveyQuestionUpdate) SetLongitude(f float64) *SurveyQuestionUpdate {
	squ.mutation.ResetLongitude()
	squ.mutation.SetLongitude(f)
	return squ
}

// SetNillableLongitude sets the longitude field if the given value is not nil.
func (squ *SurveyQuestionUpdate) SetNillableLongitude(f *float64) *SurveyQuestionUpdate {
	if f != nil {
		squ.SetLongitude(*f)
	}
	return squ
}

// AddLongitude adds f to longitude.
func (squ *SurveyQuestionUpdate) AddLongitude(f float64) *SurveyQuestionUpdate {
	squ.mutation.AddLongitude(f)
	return squ
}

// ClearLongitude clears the value of longitude.
func (squ *SurveyQuestionUpdate) ClearLongitude() *SurveyQuestionUpdate {
	squ.mutation.ClearLongitude()
	return squ
}

// SetLocationAccuracy sets the location_accuracy field.
func (squ *SurveyQuestionUpdate) SetLocationAccuracy(f float64) *SurveyQuestionUpdate {
	squ.mutation.ResetLocationAccuracy()
	squ.mutation.SetLocationAccuracy(f)
	return squ
}

// SetNillableLocationAccuracy sets the location_accuracy field if the given value is not nil.
func (squ *SurveyQuestionUpdate) SetNillableLocationAccuracy(f *float64) *SurveyQuestionUpdate {
	if f != nil {
		squ.SetLocationAccuracy(*f)
	}
	return squ
}

// AddLocationAccuracy adds f to location_accuracy.
func (squ *SurveyQuestionUpdate) AddLocationAccuracy(f float64) *SurveyQuestionUpdate {
	squ.mutation.AddLocationAccuracy(f)
	return squ
}

// ClearLocationAccuracy clears the value of location_accuracy.
func (squ *SurveyQuestionUpdate) ClearLocationAccuracy() *SurveyQuestionUpdate {
	squ.mutation.ClearLocationAccuracy()
	return squ
}

// SetAltitude sets the altitude field.
func (squ *SurveyQuestionUpdate) SetAltitude(f float64) *SurveyQuestionUpdate {
	squ.mutation.ResetAltitude()
	squ.mutation.SetAltitude(f)
	return squ
}

// SetNillableAltitude sets the altitude field if the given value is not nil.
func (squ *SurveyQuestionUpdate) SetNillableAltitude(f *float64) *SurveyQuestionUpdate {
	if f != nil {
		squ.SetAltitude(*f)
	}
	return squ
}

// AddAltitude adds f to altitude.
func (squ *SurveyQuestionUpdate) AddAltitude(f float64) *SurveyQuestionUpdate {
	squ.mutation.AddAltitude(f)
	return squ
}

// ClearAltitude clears the value of altitude.
func (squ *SurveyQuestionUpdate) ClearAltitude() *SurveyQuestionUpdate {
	squ.mutation.ClearAltitude()
	return squ
}

// SetPhoneData sets the phone_data field.
func (squ *SurveyQuestionUpdate) SetPhoneData(s string) *SurveyQuestionUpdate {
	squ.mutation.SetPhoneData(s)
	return squ
}

// SetNillablePhoneData sets the phone_data field if the given value is not nil.
func (squ *SurveyQuestionUpdate) SetNillablePhoneData(s *string) *SurveyQuestionUpdate {
	if s != nil {
		squ.SetPhoneData(*s)
	}
	return squ
}

// ClearPhoneData clears the value of phone_data.
func (squ *SurveyQuestionUpdate) ClearPhoneData() *SurveyQuestionUpdate {
	squ.mutation.ClearPhoneData()
	return squ
}

// SetTextData sets the text_data field.
func (squ *SurveyQuestionUpdate) SetTextData(s string) *SurveyQuestionUpdate {
	squ.mutation.SetTextData(s)
	return squ
}

// SetNillableTextData sets the text_data field if the given value is not nil.
func (squ *SurveyQuestionUpdate) SetNillableTextData(s *string) *SurveyQuestionUpdate {
	if s != nil {
		squ.SetTextData(*s)
	}
	return squ
}

// ClearTextData clears the value of text_data.
func (squ *SurveyQuestionUpdate) ClearTextData() *SurveyQuestionUpdate {
	squ.mutation.ClearTextData()
	return squ
}

// SetFloatData sets the float_data field.
func (squ *SurveyQuestionUpdate) SetFloatData(f float64) *SurveyQuestionUpdate {
	squ.mutation.ResetFloatData()
	squ.mutation.SetFloatData(f)
	return squ
}

// SetNillableFloatData sets the float_data field if the given value is not nil.
func (squ *SurveyQuestionUpdate) SetNillableFloatData(f *float64) *SurveyQuestionUpdate {
	if f != nil {
		squ.SetFloatData(*f)
	}
	return squ
}

// AddFloatData adds f to float_data.
func (squ *SurveyQuestionUpdate) AddFloatData(f float64) *SurveyQuestionUpdate {
	squ.mutation.AddFloatData(f)
	return squ
}

// ClearFloatData clears the value of float_data.
func (squ *SurveyQuestionUpdate) ClearFloatData() *SurveyQuestionUpdate {
	squ.mutation.ClearFloatData()
	return squ
}

// SetIntData sets the int_data field.
func (squ *SurveyQuestionUpdate) SetIntData(i int) *SurveyQuestionUpdate {
	squ.mutation.ResetIntData()
	squ.mutation.SetIntData(i)
	return squ
}

// SetNillableIntData sets the int_data field if the given value is not nil.
func (squ *SurveyQuestionUpdate) SetNillableIntData(i *int) *SurveyQuestionUpdate {
	if i != nil {
		squ.SetIntData(*i)
	}
	return squ
}

// AddIntData adds i to int_data.
func (squ *SurveyQuestionUpdate) AddIntData(i int) *SurveyQuestionUpdate {
	squ.mutation.AddIntData(i)
	return squ
}

// ClearIntData clears the value of int_data.
func (squ *SurveyQuestionUpdate) ClearIntData() *SurveyQuestionUpdate {
	squ.mutation.ClearIntData()
	return squ
}

// SetDateData sets the date_data field.
func (squ *SurveyQuestionUpdate) SetDateData(t time.Time) *SurveyQuestionUpdate {
	squ.mutation.SetDateData(t)
	return squ
}

// SetNillableDateData sets the date_data field if the given value is not nil.
func (squ *SurveyQuestionUpdate) SetNillableDateData(t *time.Time) *SurveyQuestionUpdate {
	if t != nil {
		squ.SetDateData(*t)
	}
	return squ
}

// ClearDateData clears the value of date_data.
func (squ *SurveyQuestionUpdate) ClearDateData() *SurveyQuestionUpdate {
	squ.mutation.ClearDateData()
	return squ
}

// SetSurveyID sets the survey edge to Survey by id.
func (squ *SurveyQuestionUpdate) SetSurveyID(id int) *SurveyQuestionUpdate {
	squ.mutation.SetSurveyID(id)
	return squ
}

// SetSurvey sets the survey edge to Survey.
func (squ *SurveyQuestionUpdate) SetSurvey(s *Survey) *SurveyQuestionUpdate {
	return squ.SetSurveyID(s.ID)
}

// AddWifiScanIDs adds the wifi_scan edge to SurveyWiFiScan by ids.
func (squ *SurveyQuestionUpdate) AddWifiScanIDs(ids ...int) *SurveyQuestionUpdate {
	squ.mutation.AddWifiScanIDs(ids...)
	return squ
}

// AddWifiScan adds the wifi_scan edges to SurveyWiFiScan.
func (squ *SurveyQuestionUpdate) AddWifiScan(s ...*SurveyWiFiScan) *SurveyQuestionUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squ.AddWifiScanIDs(ids...)
}

// AddCellScanIDs adds the cell_scan edge to SurveyCellScan by ids.
func (squ *SurveyQuestionUpdate) AddCellScanIDs(ids ...int) *SurveyQuestionUpdate {
	squ.mutation.AddCellScanIDs(ids...)
	return squ
}

// AddCellScan adds the cell_scan edges to SurveyCellScan.
func (squ *SurveyQuestionUpdate) AddCellScan(s ...*SurveyCellScan) *SurveyQuestionUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squ.AddCellScanIDs(ids...)
}

// AddPhotoDatumIDs adds the photo_data edge to File by ids.
func (squ *SurveyQuestionUpdate) AddPhotoDatumIDs(ids ...int) *SurveyQuestionUpdate {
	squ.mutation.AddPhotoDatumIDs(ids...)
	return squ
}

// AddPhotoData adds the photo_data edges to File.
func (squ *SurveyQuestionUpdate) AddPhotoData(f ...*File) *SurveyQuestionUpdate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return squ.AddPhotoDatumIDs(ids...)
}

// ClearSurvey clears the survey edge to Survey.
func (squ *SurveyQuestionUpdate) ClearSurvey() *SurveyQuestionUpdate {
	squ.mutation.ClearSurvey()
	return squ
}

// RemoveWifiScanIDs removes the wifi_scan edge to SurveyWiFiScan by ids.
func (squ *SurveyQuestionUpdate) RemoveWifiScanIDs(ids ...int) *SurveyQuestionUpdate {
	squ.mutation.RemoveWifiScanIDs(ids...)
	return squ
}

// RemoveWifiScan removes wifi_scan edges to SurveyWiFiScan.
func (squ *SurveyQuestionUpdate) RemoveWifiScan(s ...*SurveyWiFiScan) *SurveyQuestionUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squ.RemoveWifiScanIDs(ids...)
}

// RemoveCellScanIDs removes the cell_scan edge to SurveyCellScan by ids.
func (squ *SurveyQuestionUpdate) RemoveCellScanIDs(ids ...int) *SurveyQuestionUpdate {
	squ.mutation.RemoveCellScanIDs(ids...)
	return squ
}

// RemoveCellScan removes cell_scan edges to SurveyCellScan.
func (squ *SurveyQuestionUpdate) RemoveCellScan(s ...*SurveyCellScan) *SurveyQuestionUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squ.RemoveCellScanIDs(ids...)
}

// RemovePhotoDatumIDs removes the photo_data edge to File by ids.
func (squ *SurveyQuestionUpdate) RemovePhotoDatumIDs(ids ...int) *SurveyQuestionUpdate {
	squ.mutation.RemovePhotoDatumIDs(ids...)
	return squ
}

// RemovePhotoData removes photo_data edges to File.
func (squ *SurveyQuestionUpdate) RemovePhotoData(f ...*File) *SurveyQuestionUpdate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return squ.RemovePhotoDatumIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (squ *SurveyQuestionUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := squ.mutation.UpdateTime(); !ok {
		v := surveyquestion.UpdateDefaultUpdateTime()
		squ.mutation.SetUpdateTime(v)
	}

	if _, ok := squ.mutation.SurveyID(); squ.mutation.SurveyCleared() && !ok {
		return 0, errors.New("ent: clearing a unique edge \"survey\"")
	}

	var (
		err      error
		affected int
	)
	if len(squ.hooks) == 0 {
		affected, err = squ.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyQuestionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			squ.mutation = mutation
			affected, err = squ.sqlSave(ctx)
			return affected, err
		})
		for i := len(squ.hooks) - 1; i >= 0; i-- {
			mut = squ.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, squ.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (squ *SurveyQuestionUpdate) SaveX(ctx context.Context) int {
	affected, err := squ.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (squ *SurveyQuestionUpdate) Exec(ctx context.Context) error {
	_, err := squ.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (squ *SurveyQuestionUpdate) ExecX(ctx context.Context) {
	if err := squ.Exec(ctx); err != nil {
		panic(err)
	}
}

func (squ *SurveyQuestionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveyquestion.Table,
			Columns: surveyquestion.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: surveyquestion.FieldID,
			},
		},
	}
	if ps := squ.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := squ.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveyquestion.FieldUpdateTime,
		})
	}
	if value, ok := squ.mutation.FormName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldFormName,
		})
	}
	if squ.mutation.FormNameCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldFormName,
		})
	}
	if value, ok := squ.mutation.FormDescription(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldFormDescription,
		})
	}
	if squ.mutation.FormDescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldFormDescription,
		})
	}
	if value, ok := squ.mutation.FormIndex(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveyquestion.FieldFormIndex,
		})
	}
	if value, ok := squ.mutation.AddedFormIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveyquestion.FieldFormIndex,
		})
	}
	if value, ok := squ.mutation.QuestionType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldQuestionType,
		})
	}
	if squ.mutation.QuestionTypeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldQuestionType,
		})
	}
	if value, ok := squ.mutation.QuestionFormat(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldQuestionFormat,
		})
	}
	if squ.mutation.QuestionFormatCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldQuestionFormat,
		})
	}
	if value, ok := squ.mutation.QuestionText(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldQuestionText,
		})
	}
	if squ.mutation.QuestionTextCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldQuestionText,
		})
	}
	if value, ok := squ.mutation.QuestionIndex(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveyquestion.FieldQuestionIndex,
		})
	}
	if value, ok := squ.mutation.AddedQuestionIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveyquestion.FieldQuestionIndex,
		})
	}
	if value, ok := squ.mutation.BoolData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: surveyquestion.FieldBoolData,
		})
	}
	if squ.mutation.BoolDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: surveyquestion.FieldBoolData,
		})
	}
	if value, ok := squ.mutation.EmailData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldEmailData,
		})
	}
	if squ.mutation.EmailDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldEmailData,
		})
	}
	if value, ok := squ.mutation.Latitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldLatitude,
		})
	}
	if value, ok := squ.mutation.AddedLatitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldLatitude,
		})
	}
	if squ.mutation.LatitudeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldLatitude,
		})
	}
	if value, ok := squ.mutation.Longitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldLongitude,
		})
	}
	if value, ok := squ.mutation.AddedLongitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldLongitude,
		})
	}
	if squ.mutation.LongitudeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldLongitude,
		})
	}
	if value, ok := squ.mutation.LocationAccuracy(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldLocationAccuracy,
		})
	}
	if value, ok := squ.mutation.AddedLocationAccuracy(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldLocationAccuracy,
		})
	}
	if squ.mutation.LocationAccuracyCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldLocationAccuracy,
		})
	}
	if value, ok := squ.mutation.Altitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldAltitude,
		})
	}
	if value, ok := squ.mutation.AddedAltitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldAltitude,
		})
	}
	if squ.mutation.AltitudeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldAltitude,
		})
	}
	if value, ok := squ.mutation.PhoneData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldPhoneData,
		})
	}
	if squ.mutation.PhoneDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldPhoneData,
		})
	}
	if value, ok := squ.mutation.TextData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldTextData,
		})
	}
	if squ.mutation.TextDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldTextData,
		})
	}
	if value, ok := squ.mutation.FloatData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldFloatData,
		})
	}
	if value, ok := squ.mutation.AddedFloatData(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldFloatData,
		})
	}
	if squ.mutation.FloatDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldFloatData,
		})
	}
	if value, ok := squ.mutation.IntData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveyquestion.FieldIntData,
		})
	}
	if value, ok := squ.mutation.AddedIntData(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveyquestion.FieldIntData,
		})
	}
	if squ.mutation.IntDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveyquestion.FieldIntData,
		})
	}
	if value, ok := squ.mutation.DateData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveyquestion.FieldDateData,
		})
	}
	if squ.mutation.DateDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: surveyquestion.FieldDateData,
		})
	}
	if squ.mutation.SurveyCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squ.mutation.SurveyIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := squ.mutation.RemovedWifiScanIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squ.mutation.WifiScanIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := squ.mutation.RemovedCellScanIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squ.mutation.CellScanIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := squ.mutation.RemovedPhotoDataIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squ.mutation.PhotoDataIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, squ.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{surveyquestion.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// SurveyQuestionUpdateOne is the builder for updating a single SurveyQuestion entity.
type SurveyQuestionUpdateOne struct {
	config
	hooks    []Hook
	mutation *SurveyQuestionMutation
}

// SetFormName sets the form_name field.
func (squo *SurveyQuestionUpdateOne) SetFormName(s string) *SurveyQuestionUpdateOne {
	squo.mutation.SetFormName(s)
	return squo
}

// SetNillableFormName sets the form_name field if the given value is not nil.
func (squo *SurveyQuestionUpdateOne) SetNillableFormName(s *string) *SurveyQuestionUpdateOne {
	if s != nil {
		squo.SetFormName(*s)
	}
	return squo
}

// ClearFormName clears the value of form_name.
func (squo *SurveyQuestionUpdateOne) ClearFormName() *SurveyQuestionUpdateOne {
	squo.mutation.ClearFormName()
	return squo
}

// SetFormDescription sets the form_description field.
func (squo *SurveyQuestionUpdateOne) SetFormDescription(s string) *SurveyQuestionUpdateOne {
	squo.mutation.SetFormDescription(s)
	return squo
}

// SetNillableFormDescription sets the form_description field if the given value is not nil.
func (squo *SurveyQuestionUpdateOne) SetNillableFormDescription(s *string) *SurveyQuestionUpdateOne {
	if s != nil {
		squo.SetFormDescription(*s)
	}
	return squo
}

// ClearFormDescription clears the value of form_description.
func (squo *SurveyQuestionUpdateOne) ClearFormDescription() *SurveyQuestionUpdateOne {
	squo.mutation.ClearFormDescription()
	return squo
}

// SetFormIndex sets the form_index field.
func (squo *SurveyQuestionUpdateOne) SetFormIndex(i int) *SurveyQuestionUpdateOne {
	squo.mutation.ResetFormIndex()
	squo.mutation.SetFormIndex(i)
	return squo
}

// AddFormIndex adds i to form_index.
func (squo *SurveyQuestionUpdateOne) AddFormIndex(i int) *SurveyQuestionUpdateOne {
	squo.mutation.AddFormIndex(i)
	return squo
}

// SetQuestionType sets the question_type field.
func (squo *SurveyQuestionUpdateOne) SetQuestionType(s string) *SurveyQuestionUpdateOne {
	squo.mutation.SetQuestionType(s)
	return squo
}

// SetNillableQuestionType sets the question_type field if the given value is not nil.
func (squo *SurveyQuestionUpdateOne) SetNillableQuestionType(s *string) *SurveyQuestionUpdateOne {
	if s != nil {
		squo.SetQuestionType(*s)
	}
	return squo
}

// ClearQuestionType clears the value of question_type.
func (squo *SurveyQuestionUpdateOne) ClearQuestionType() *SurveyQuestionUpdateOne {
	squo.mutation.ClearQuestionType()
	return squo
}

// SetQuestionFormat sets the question_format field.
func (squo *SurveyQuestionUpdateOne) SetQuestionFormat(s string) *SurveyQuestionUpdateOne {
	squo.mutation.SetQuestionFormat(s)
	return squo
}

// SetNillableQuestionFormat sets the question_format field if the given value is not nil.
func (squo *SurveyQuestionUpdateOne) SetNillableQuestionFormat(s *string) *SurveyQuestionUpdateOne {
	if s != nil {
		squo.SetQuestionFormat(*s)
	}
	return squo
}

// ClearQuestionFormat clears the value of question_format.
func (squo *SurveyQuestionUpdateOne) ClearQuestionFormat() *SurveyQuestionUpdateOne {
	squo.mutation.ClearQuestionFormat()
	return squo
}

// SetQuestionText sets the question_text field.
func (squo *SurveyQuestionUpdateOne) SetQuestionText(s string) *SurveyQuestionUpdateOne {
	squo.mutation.SetQuestionText(s)
	return squo
}

// SetNillableQuestionText sets the question_text field if the given value is not nil.
func (squo *SurveyQuestionUpdateOne) SetNillableQuestionText(s *string) *SurveyQuestionUpdateOne {
	if s != nil {
		squo.SetQuestionText(*s)
	}
	return squo
}

// ClearQuestionText clears the value of question_text.
func (squo *SurveyQuestionUpdateOne) ClearQuestionText() *SurveyQuestionUpdateOne {
	squo.mutation.ClearQuestionText()
	return squo
}

// SetQuestionIndex sets the question_index field.
func (squo *SurveyQuestionUpdateOne) SetQuestionIndex(i int) *SurveyQuestionUpdateOne {
	squo.mutation.ResetQuestionIndex()
	squo.mutation.SetQuestionIndex(i)
	return squo
}

// AddQuestionIndex adds i to question_index.
func (squo *SurveyQuestionUpdateOne) AddQuestionIndex(i int) *SurveyQuestionUpdateOne {
	squo.mutation.AddQuestionIndex(i)
	return squo
}

// SetBoolData sets the bool_data field.
func (squo *SurveyQuestionUpdateOne) SetBoolData(b bool) *SurveyQuestionUpdateOne {
	squo.mutation.SetBoolData(b)
	return squo
}

// SetNillableBoolData sets the bool_data field if the given value is not nil.
func (squo *SurveyQuestionUpdateOne) SetNillableBoolData(b *bool) *SurveyQuestionUpdateOne {
	if b != nil {
		squo.SetBoolData(*b)
	}
	return squo
}

// ClearBoolData clears the value of bool_data.
func (squo *SurveyQuestionUpdateOne) ClearBoolData() *SurveyQuestionUpdateOne {
	squo.mutation.ClearBoolData()
	return squo
}

// SetEmailData sets the email_data field.
func (squo *SurveyQuestionUpdateOne) SetEmailData(s string) *SurveyQuestionUpdateOne {
	squo.mutation.SetEmailData(s)
	return squo
}

// SetNillableEmailData sets the email_data field if the given value is not nil.
func (squo *SurveyQuestionUpdateOne) SetNillableEmailData(s *string) *SurveyQuestionUpdateOne {
	if s != nil {
		squo.SetEmailData(*s)
	}
	return squo
}

// ClearEmailData clears the value of email_data.
func (squo *SurveyQuestionUpdateOne) ClearEmailData() *SurveyQuestionUpdateOne {
	squo.mutation.ClearEmailData()
	return squo
}

// SetLatitude sets the latitude field.
func (squo *SurveyQuestionUpdateOne) SetLatitude(f float64) *SurveyQuestionUpdateOne {
	squo.mutation.ResetLatitude()
	squo.mutation.SetLatitude(f)
	return squo
}

// SetNillableLatitude sets the latitude field if the given value is not nil.
func (squo *SurveyQuestionUpdateOne) SetNillableLatitude(f *float64) *SurveyQuestionUpdateOne {
	if f != nil {
		squo.SetLatitude(*f)
	}
	return squo
}

// AddLatitude adds f to latitude.
func (squo *SurveyQuestionUpdateOne) AddLatitude(f float64) *SurveyQuestionUpdateOne {
	squo.mutation.AddLatitude(f)
	return squo
}

// ClearLatitude clears the value of latitude.
func (squo *SurveyQuestionUpdateOne) ClearLatitude() *SurveyQuestionUpdateOne {
	squo.mutation.ClearLatitude()
	return squo
}

// SetLongitude sets the longitude field.
func (squo *SurveyQuestionUpdateOne) SetLongitude(f float64) *SurveyQuestionUpdateOne {
	squo.mutation.ResetLongitude()
	squo.mutation.SetLongitude(f)
	return squo
}

// SetNillableLongitude sets the longitude field if the given value is not nil.
func (squo *SurveyQuestionUpdateOne) SetNillableLongitude(f *float64) *SurveyQuestionUpdateOne {
	if f != nil {
		squo.SetLongitude(*f)
	}
	return squo
}

// AddLongitude adds f to longitude.
func (squo *SurveyQuestionUpdateOne) AddLongitude(f float64) *SurveyQuestionUpdateOne {
	squo.mutation.AddLongitude(f)
	return squo
}

// ClearLongitude clears the value of longitude.
func (squo *SurveyQuestionUpdateOne) ClearLongitude() *SurveyQuestionUpdateOne {
	squo.mutation.ClearLongitude()
	return squo
}

// SetLocationAccuracy sets the location_accuracy field.
func (squo *SurveyQuestionUpdateOne) SetLocationAccuracy(f float64) *SurveyQuestionUpdateOne {
	squo.mutation.ResetLocationAccuracy()
	squo.mutation.SetLocationAccuracy(f)
	return squo
}

// SetNillableLocationAccuracy sets the location_accuracy field if the given value is not nil.
func (squo *SurveyQuestionUpdateOne) SetNillableLocationAccuracy(f *float64) *SurveyQuestionUpdateOne {
	if f != nil {
		squo.SetLocationAccuracy(*f)
	}
	return squo
}

// AddLocationAccuracy adds f to location_accuracy.
func (squo *SurveyQuestionUpdateOne) AddLocationAccuracy(f float64) *SurveyQuestionUpdateOne {
	squo.mutation.AddLocationAccuracy(f)
	return squo
}

// ClearLocationAccuracy clears the value of location_accuracy.
func (squo *SurveyQuestionUpdateOne) ClearLocationAccuracy() *SurveyQuestionUpdateOne {
	squo.mutation.ClearLocationAccuracy()
	return squo
}

// SetAltitude sets the altitude field.
func (squo *SurveyQuestionUpdateOne) SetAltitude(f float64) *SurveyQuestionUpdateOne {
	squo.mutation.ResetAltitude()
	squo.mutation.SetAltitude(f)
	return squo
}

// SetNillableAltitude sets the altitude field if the given value is not nil.
func (squo *SurveyQuestionUpdateOne) SetNillableAltitude(f *float64) *SurveyQuestionUpdateOne {
	if f != nil {
		squo.SetAltitude(*f)
	}
	return squo
}

// AddAltitude adds f to altitude.
func (squo *SurveyQuestionUpdateOne) AddAltitude(f float64) *SurveyQuestionUpdateOne {
	squo.mutation.AddAltitude(f)
	return squo
}

// ClearAltitude clears the value of altitude.
func (squo *SurveyQuestionUpdateOne) ClearAltitude() *SurveyQuestionUpdateOne {
	squo.mutation.ClearAltitude()
	return squo
}

// SetPhoneData sets the phone_data field.
func (squo *SurveyQuestionUpdateOne) SetPhoneData(s string) *SurveyQuestionUpdateOne {
	squo.mutation.SetPhoneData(s)
	return squo
}

// SetNillablePhoneData sets the phone_data field if the given value is not nil.
func (squo *SurveyQuestionUpdateOne) SetNillablePhoneData(s *string) *SurveyQuestionUpdateOne {
	if s != nil {
		squo.SetPhoneData(*s)
	}
	return squo
}

// ClearPhoneData clears the value of phone_data.
func (squo *SurveyQuestionUpdateOne) ClearPhoneData() *SurveyQuestionUpdateOne {
	squo.mutation.ClearPhoneData()
	return squo
}

// SetTextData sets the text_data field.
func (squo *SurveyQuestionUpdateOne) SetTextData(s string) *SurveyQuestionUpdateOne {
	squo.mutation.SetTextData(s)
	return squo
}

// SetNillableTextData sets the text_data field if the given value is not nil.
func (squo *SurveyQuestionUpdateOne) SetNillableTextData(s *string) *SurveyQuestionUpdateOne {
	if s != nil {
		squo.SetTextData(*s)
	}
	return squo
}

// ClearTextData clears the value of text_data.
func (squo *SurveyQuestionUpdateOne) ClearTextData() *SurveyQuestionUpdateOne {
	squo.mutation.ClearTextData()
	return squo
}

// SetFloatData sets the float_data field.
func (squo *SurveyQuestionUpdateOne) SetFloatData(f float64) *SurveyQuestionUpdateOne {
	squo.mutation.ResetFloatData()
	squo.mutation.SetFloatData(f)
	return squo
}

// SetNillableFloatData sets the float_data field if the given value is not nil.
func (squo *SurveyQuestionUpdateOne) SetNillableFloatData(f *float64) *SurveyQuestionUpdateOne {
	if f != nil {
		squo.SetFloatData(*f)
	}
	return squo
}

// AddFloatData adds f to float_data.
func (squo *SurveyQuestionUpdateOne) AddFloatData(f float64) *SurveyQuestionUpdateOne {
	squo.mutation.AddFloatData(f)
	return squo
}

// ClearFloatData clears the value of float_data.
func (squo *SurveyQuestionUpdateOne) ClearFloatData() *SurveyQuestionUpdateOne {
	squo.mutation.ClearFloatData()
	return squo
}

// SetIntData sets the int_data field.
func (squo *SurveyQuestionUpdateOne) SetIntData(i int) *SurveyQuestionUpdateOne {
	squo.mutation.ResetIntData()
	squo.mutation.SetIntData(i)
	return squo
}

// SetNillableIntData sets the int_data field if the given value is not nil.
func (squo *SurveyQuestionUpdateOne) SetNillableIntData(i *int) *SurveyQuestionUpdateOne {
	if i != nil {
		squo.SetIntData(*i)
	}
	return squo
}

// AddIntData adds i to int_data.
func (squo *SurveyQuestionUpdateOne) AddIntData(i int) *SurveyQuestionUpdateOne {
	squo.mutation.AddIntData(i)
	return squo
}

// ClearIntData clears the value of int_data.
func (squo *SurveyQuestionUpdateOne) ClearIntData() *SurveyQuestionUpdateOne {
	squo.mutation.ClearIntData()
	return squo
}

// SetDateData sets the date_data field.
func (squo *SurveyQuestionUpdateOne) SetDateData(t time.Time) *SurveyQuestionUpdateOne {
	squo.mutation.SetDateData(t)
	return squo
}

// SetNillableDateData sets the date_data field if the given value is not nil.
func (squo *SurveyQuestionUpdateOne) SetNillableDateData(t *time.Time) *SurveyQuestionUpdateOne {
	if t != nil {
		squo.SetDateData(*t)
	}
	return squo
}

// ClearDateData clears the value of date_data.
func (squo *SurveyQuestionUpdateOne) ClearDateData() *SurveyQuestionUpdateOne {
	squo.mutation.ClearDateData()
	return squo
}

// SetSurveyID sets the survey edge to Survey by id.
func (squo *SurveyQuestionUpdateOne) SetSurveyID(id int) *SurveyQuestionUpdateOne {
	squo.mutation.SetSurveyID(id)
	return squo
}

// SetSurvey sets the survey edge to Survey.
func (squo *SurveyQuestionUpdateOne) SetSurvey(s *Survey) *SurveyQuestionUpdateOne {
	return squo.SetSurveyID(s.ID)
}

// AddWifiScanIDs adds the wifi_scan edge to SurveyWiFiScan by ids.
func (squo *SurveyQuestionUpdateOne) AddWifiScanIDs(ids ...int) *SurveyQuestionUpdateOne {
	squo.mutation.AddWifiScanIDs(ids...)
	return squo
}

// AddWifiScan adds the wifi_scan edges to SurveyWiFiScan.
func (squo *SurveyQuestionUpdateOne) AddWifiScan(s ...*SurveyWiFiScan) *SurveyQuestionUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squo.AddWifiScanIDs(ids...)
}

// AddCellScanIDs adds the cell_scan edge to SurveyCellScan by ids.
func (squo *SurveyQuestionUpdateOne) AddCellScanIDs(ids ...int) *SurveyQuestionUpdateOne {
	squo.mutation.AddCellScanIDs(ids...)
	return squo
}

// AddCellScan adds the cell_scan edges to SurveyCellScan.
func (squo *SurveyQuestionUpdateOne) AddCellScan(s ...*SurveyCellScan) *SurveyQuestionUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squo.AddCellScanIDs(ids...)
}

// AddPhotoDatumIDs adds the photo_data edge to File by ids.
func (squo *SurveyQuestionUpdateOne) AddPhotoDatumIDs(ids ...int) *SurveyQuestionUpdateOne {
	squo.mutation.AddPhotoDatumIDs(ids...)
	return squo
}

// AddPhotoData adds the photo_data edges to File.
func (squo *SurveyQuestionUpdateOne) AddPhotoData(f ...*File) *SurveyQuestionUpdateOne {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return squo.AddPhotoDatumIDs(ids...)
}

// ClearSurvey clears the survey edge to Survey.
func (squo *SurveyQuestionUpdateOne) ClearSurvey() *SurveyQuestionUpdateOne {
	squo.mutation.ClearSurvey()
	return squo
}

// RemoveWifiScanIDs removes the wifi_scan edge to SurveyWiFiScan by ids.
func (squo *SurveyQuestionUpdateOne) RemoveWifiScanIDs(ids ...int) *SurveyQuestionUpdateOne {
	squo.mutation.RemoveWifiScanIDs(ids...)
	return squo
}

// RemoveWifiScan removes wifi_scan edges to SurveyWiFiScan.
func (squo *SurveyQuestionUpdateOne) RemoveWifiScan(s ...*SurveyWiFiScan) *SurveyQuestionUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squo.RemoveWifiScanIDs(ids...)
}

// RemoveCellScanIDs removes the cell_scan edge to SurveyCellScan by ids.
func (squo *SurveyQuestionUpdateOne) RemoveCellScanIDs(ids ...int) *SurveyQuestionUpdateOne {
	squo.mutation.RemoveCellScanIDs(ids...)
	return squo
}

// RemoveCellScan removes cell_scan edges to SurveyCellScan.
func (squo *SurveyQuestionUpdateOne) RemoveCellScan(s ...*SurveyCellScan) *SurveyQuestionUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squo.RemoveCellScanIDs(ids...)
}

// RemovePhotoDatumIDs removes the photo_data edge to File by ids.
func (squo *SurveyQuestionUpdateOne) RemovePhotoDatumIDs(ids ...int) *SurveyQuestionUpdateOne {
	squo.mutation.RemovePhotoDatumIDs(ids...)
	return squo
}

// RemovePhotoData removes photo_data edges to File.
func (squo *SurveyQuestionUpdateOne) RemovePhotoData(f ...*File) *SurveyQuestionUpdateOne {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return squo.RemovePhotoDatumIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (squo *SurveyQuestionUpdateOne) Save(ctx context.Context) (*SurveyQuestion, error) {
	if _, ok := squo.mutation.UpdateTime(); !ok {
		v := surveyquestion.UpdateDefaultUpdateTime()
		squo.mutation.SetUpdateTime(v)
	}

	if _, ok := squo.mutation.SurveyID(); squo.mutation.SurveyCleared() && !ok {
		return nil, errors.New("ent: clearing a unique edge \"survey\"")
	}

	var (
		err  error
		node *SurveyQuestion
	)
	if len(squo.hooks) == 0 {
		node, err = squo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyQuestionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			squo.mutation = mutation
			node, err = squo.sqlSave(ctx)
			return node, err
		})
		for i := len(squo.hooks) - 1; i >= 0; i-- {
			mut = squo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, squo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (squo *SurveyQuestionUpdateOne) SaveX(ctx context.Context) *SurveyQuestion {
	sq, err := squo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return sq
}

// Exec executes the query on the entity.
func (squo *SurveyQuestionUpdateOne) Exec(ctx context.Context) error {
	_, err := squo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (squo *SurveyQuestionUpdateOne) ExecX(ctx context.Context) {
	if err := squo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (squo *SurveyQuestionUpdateOne) sqlSave(ctx context.Context) (sq *SurveyQuestion, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveyquestion.Table,
			Columns: surveyquestion.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: surveyquestion.FieldID,
			},
		},
	}
	id, ok := squo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing SurveyQuestion.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := squo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveyquestion.FieldUpdateTime,
		})
	}
	if value, ok := squo.mutation.FormName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldFormName,
		})
	}
	if squo.mutation.FormNameCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldFormName,
		})
	}
	if value, ok := squo.mutation.FormDescription(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldFormDescription,
		})
	}
	if squo.mutation.FormDescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldFormDescription,
		})
	}
	if value, ok := squo.mutation.FormIndex(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveyquestion.FieldFormIndex,
		})
	}
	if value, ok := squo.mutation.AddedFormIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveyquestion.FieldFormIndex,
		})
	}
	if value, ok := squo.mutation.QuestionType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldQuestionType,
		})
	}
	if squo.mutation.QuestionTypeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldQuestionType,
		})
	}
	if value, ok := squo.mutation.QuestionFormat(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldQuestionFormat,
		})
	}
	if squo.mutation.QuestionFormatCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldQuestionFormat,
		})
	}
	if value, ok := squo.mutation.QuestionText(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldQuestionText,
		})
	}
	if squo.mutation.QuestionTextCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldQuestionText,
		})
	}
	if value, ok := squo.mutation.QuestionIndex(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveyquestion.FieldQuestionIndex,
		})
	}
	if value, ok := squo.mutation.AddedQuestionIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveyquestion.FieldQuestionIndex,
		})
	}
	if value, ok := squo.mutation.BoolData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: surveyquestion.FieldBoolData,
		})
	}
	if squo.mutation.BoolDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: surveyquestion.FieldBoolData,
		})
	}
	if value, ok := squo.mutation.EmailData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldEmailData,
		})
	}
	if squo.mutation.EmailDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldEmailData,
		})
	}
	if value, ok := squo.mutation.Latitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldLatitude,
		})
	}
	if value, ok := squo.mutation.AddedLatitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldLatitude,
		})
	}
	if squo.mutation.LatitudeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldLatitude,
		})
	}
	if value, ok := squo.mutation.Longitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldLongitude,
		})
	}
	if value, ok := squo.mutation.AddedLongitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldLongitude,
		})
	}
	if squo.mutation.LongitudeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldLongitude,
		})
	}
	if value, ok := squo.mutation.LocationAccuracy(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldLocationAccuracy,
		})
	}
	if value, ok := squo.mutation.AddedLocationAccuracy(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldLocationAccuracy,
		})
	}
	if squo.mutation.LocationAccuracyCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldLocationAccuracy,
		})
	}
	if value, ok := squo.mutation.Altitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldAltitude,
		})
	}
	if value, ok := squo.mutation.AddedAltitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldAltitude,
		})
	}
	if squo.mutation.AltitudeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldAltitude,
		})
	}
	if value, ok := squo.mutation.PhoneData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldPhoneData,
		})
	}
	if squo.mutation.PhoneDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldPhoneData,
		})
	}
	if value, ok := squo.mutation.TextData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: surveyquestion.FieldTextData,
		})
	}
	if squo.mutation.TextDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldTextData,
		})
	}
	if value, ok := squo.mutation.FloatData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldFloatData,
		})
	}
	if value, ok := squo.mutation.AddedFloatData(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: surveyquestion.FieldFloatData,
		})
	}
	if squo.mutation.FloatDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldFloatData,
		})
	}
	if value, ok := squo.mutation.IntData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveyquestion.FieldIntData,
		})
	}
	if value, ok := squo.mutation.AddedIntData(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: surveyquestion.FieldIntData,
		})
	}
	if squo.mutation.IntDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveyquestion.FieldIntData,
		})
	}
	if value, ok := squo.mutation.DateData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: surveyquestion.FieldDateData,
		})
	}
	if squo.mutation.DateDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: surveyquestion.FieldDateData,
		})
	}
	if squo.mutation.SurveyCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squo.mutation.SurveyIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := squo.mutation.RemovedWifiScanIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squo.mutation.WifiScanIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := squo.mutation.RemovedCellScanIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squo.mutation.CellScanIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := squo.mutation.RemovedPhotoDataIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squo.mutation.PhotoDataIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	sq = &SurveyQuestion{config: squo.config}
	_spec.Assign = sq.assignValues
	_spec.ScanValues = sq.scanValues()
	if err = sqlgraph.UpdateNode(ctx, squo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{surveyquestion.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return sq, nil
}
