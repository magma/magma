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

	update_time            *time.Time
	form_name              *string
	clearform_name         bool
	form_description       *string
	clearform_description  bool
	form_index             *int
	addform_index          *int
	question_type          *string
	clearquestion_type     bool
	question_format        *string
	clearquestion_format   bool
	question_text          *string
	clearquestion_text     bool
	question_index         *int
	addquestion_index      *int
	bool_data              *bool
	clearbool_data         bool
	email_data             *string
	clearemail_data        bool
	latitude               *float64
	addlatitude            *float64
	clearlatitude          bool
	longitude              *float64
	addlongitude           *float64
	clearlongitude         bool
	location_accuracy      *float64
	addlocation_accuracy   *float64
	clearlocation_accuracy bool
	altitude               *float64
	addaltitude            *float64
	clearaltitude          bool
	phone_data             *string
	clearphone_data        bool
	text_data              *string
	cleartext_data         bool
	float_data             *float64
	addfloat_data          *float64
	clearfloat_data        bool
	int_data               *int
	addint_data            *int
	clearint_data          bool
	date_data              *time.Time
	cleardate_data         bool
	survey                 map[string]struct{}
	wifi_scan              map[string]struct{}
	cell_scan              map[string]struct{}
	photo_data             map[string]struct{}
	clearedSurvey          bool
	removedWifiScan        map[string]struct{}
	removedCellScan        map[string]struct{}
	removedPhotoData       map[string]struct{}
	predicates             []predicate.SurveyQuestion
}

// Where adds a new predicate for the builder.
func (squ *SurveyQuestionUpdate) Where(ps ...predicate.SurveyQuestion) *SurveyQuestionUpdate {
	squ.predicates = append(squ.predicates, ps...)
	return squ
}

// SetFormName sets the form_name field.
func (squ *SurveyQuestionUpdate) SetFormName(s string) *SurveyQuestionUpdate {
	squ.form_name = &s
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
	squ.form_name = nil
	squ.clearform_name = true
	return squ
}

// SetFormDescription sets the form_description field.
func (squ *SurveyQuestionUpdate) SetFormDescription(s string) *SurveyQuestionUpdate {
	squ.form_description = &s
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
	squ.form_description = nil
	squ.clearform_description = true
	return squ
}

// SetFormIndex sets the form_index field.
func (squ *SurveyQuestionUpdate) SetFormIndex(i int) *SurveyQuestionUpdate {
	squ.form_index = &i
	squ.addform_index = nil
	return squ
}

// AddFormIndex adds i to form_index.
func (squ *SurveyQuestionUpdate) AddFormIndex(i int) *SurveyQuestionUpdate {
	if squ.addform_index == nil {
		squ.addform_index = &i
	} else {
		*squ.addform_index += i
	}
	return squ
}

// SetQuestionType sets the question_type field.
func (squ *SurveyQuestionUpdate) SetQuestionType(s string) *SurveyQuestionUpdate {
	squ.question_type = &s
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
	squ.question_type = nil
	squ.clearquestion_type = true
	return squ
}

// SetQuestionFormat sets the question_format field.
func (squ *SurveyQuestionUpdate) SetQuestionFormat(s string) *SurveyQuestionUpdate {
	squ.question_format = &s
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
	squ.question_format = nil
	squ.clearquestion_format = true
	return squ
}

// SetQuestionText sets the question_text field.
func (squ *SurveyQuestionUpdate) SetQuestionText(s string) *SurveyQuestionUpdate {
	squ.question_text = &s
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
	squ.question_text = nil
	squ.clearquestion_text = true
	return squ
}

// SetQuestionIndex sets the question_index field.
func (squ *SurveyQuestionUpdate) SetQuestionIndex(i int) *SurveyQuestionUpdate {
	squ.question_index = &i
	squ.addquestion_index = nil
	return squ
}

// AddQuestionIndex adds i to question_index.
func (squ *SurveyQuestionUpdate) AddQuestionIndex(i int) *SurveyQuestionUpdate {
	if squ.addquestion_index == nil {
		squ.addquestion_index = &i
	} else {
		*squ.addquestion_index += i
	}
	return squ
}

// SetBoolData sets the bool_data field.
func (squ *SurveyQuestionUpdate) SetBoolData(b bool) *SurveyQuestionUpdate {
	squ.bool_data = &b
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
	squ.bool_data = nil
	squ.clearbool_data = true
	return squ
}

// SetEmailData sets the email_data field.
func (squ *SurveyQuestionUpdate) SetEmailData(s string) *SurveyQuestionUpdate {
	squ.email_data = &s
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
	squ.email_data = nil
	squ.clearemail_data = true
	return squ
}

// SetLatitude sets the latitude field.
func (squ *SurveyQuestionUpdate) SetLatitude(f float64) *SurveyQuestionUpdate {
	squ.latitude = &f
	squ.addlatitude = nil
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
	if squ.addlatitude == nil {
		squ.addlatitude = &f
	} else {
		*squ.addlatitude += f
	}
	return squ
}

// ClearLatitude clears the value of latitude.
func (squ *SurveyQuestionUpdate) ClearLatitude() *SurveyQuestionUpdate {
	squ.latitude = nil
	squ.clearlatitude = true
	return squ
}

// SetLongitude sets the longitude field.
func (squ *SurveyQuestionUpdate) SetLongitude(f float64) *SurveyQuestionUpdate {
	squ.longitude = &f
	squ.addlongitude = nil
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
	if squ.addlongitude == nil {
		squ.addlongitude = &f
	} else {
		*squ.addlongitude += f
	}
	return squ
}

// ClearLongitude clears the value of longitude.
func (squ *SurveyQuestionUpdate) ClearLongitude() *SurveyQuestionUpdate {
	squ.longitude = nil
	squ.clearlongitude = true
	return squ
}

// SetLocationAccuracy sets the location_accuracy field.
func (squ *SurveyQuestionUpdate) SetLocationAccuracy(f float64) *SurveyQuestionUpdate {
	squ.location_accuracy = &f
	squ.addlocation_accuracy = nil
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
	if squ.addlocation_accuracy == nil {
		squ.addlocation_accuracy = &f
	} else {
		*squ.addlocation_accuracy += f
	}
	return squ
}

// ClearLocationAccuracy clears the value of location_accuracy.
func (squ *SurveyQuestionUpdate) ClearLocationAccuracy() *SurveyQuestionUpdate {
	squ.location_accuracy = nil
	squ.clearlocation_accuracy = true
	return squ
}

// SetAltitude sets the altitude field.
func (squ *SurveyQuestionUpdate) SetAltitude(f float64) *SurveyQuestionUpdate {
	squ.altitude = &f
	squ.addaltitude = nil
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
	if squ.addaltitude == nil {
		squ.addaltitude = &f
	} else {
		*squ.addaltitude += f
	}
	return squ
}

// ClearAltitude clears the value of altitude.
func (squ *SurveyQuestionUpdate) ClearAltitude() *SurveyQuestionUpdate {
	squ.altitude = nil
	squ.clearaltitude = true
	return squ
}

// SetPhoneData sets the phone_data field.
func (squ *SurveyQuestionUpdate) SetPhoneData(s string) *SurveyQuestionUpdate {
	squ.phone_data = &s
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
	squ.phone_data = nil
	squ.clearphone_data = true
	return squ
}

// SetTextData sets the text_data field.
func (squ *SurveyQuestionUpdate) SetTextData(s string) *SurveyQuestionUpdate {
	squ.text_data = &s
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
	squ.text_data = nil
	squ.cleartext_data = true
	return squ
}

// SetFloatData sets the float_data field.
func (squ *SurveyQuestionUpdate) SetFloatData(f float64) *SurveyQuestionUpdate {
	squ.float_data = &f
	squ.addfloat_data = nil
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
	if squ.addfloat_data == nil {
		squ.addfloat_data = &f
	} else {
		*squ.addfloat_data += f
	}
	return squ
}

// ClearFloatData clears the value of float_data.
func (squ *SurveyQuestionUpdate) ClearFloatData() *SurveyQuestionUpdate {
	squ.float_data = nil
	squ.clearfloat_data = true
	return squ
}

// SetIntData sets the int_data field.
func (squ *SurveyQuestionUpdate) SetIntData(i int) *SurveyQuestionUpdate {
	squ.int_data = &i
	squ.addint_data = nil
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
	if squ.addint_data == nil {
		squ.addint_data = &i
	} else {
		*squ.addint_data += i
	}
	return squ
}

// ClearIntData clears the value of int_data.
func (squ *SurveyQuestionUpdate) ClearIntData() *SurveyQuestionUpdate {
	squ.int_data = nil
	squ.clearint_data = true
	return squ
}

// SetDateData sets the date_data field.
func (squ *SurveyQuestionUpdate) SetDateData(t time.Time) *SurveyQuestionUpdate {
	squ.date_data = &t
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
	squ.date_data = nil
	squ.cleardate_data = true
	return squ
}

// SetSurveyID sets the survey edge to Survey by id.
func (squ *SurveyQuestionUpdate) SetSurveyID(id string) *SurveyQuestionUpdate {
	if squ.survey == nil {
		squ.survey = make(map[string]struct{})
	}
	squ.survey[id] = struct{}{}
	return squ
}

// SetSurvey sets the survey edge to Survey.
func (squ *SurveyQuestionUpdate) SetSurvey(s *Survey) *SurveyQuestionUpdate {
	return squ.SetSurveyID(s.ID)
}

// AddWifiScanIDs adds the wifi_scan edge to SurveyWiFiScan by ids.
func (squ *SurveyQuestionUpdate) AddWifiScanIDs(ids ...string) *SurveyQuestionUpdate {
	if squ.wifi_scan == nil {
		squ.wifi_scan = make(map[string]struct{})
	}
	for i := range ids {
		squ.wifi_scan[ids[i]] = struct{}{}
	}
	return squ
}

// AddWifiScan adds the wifi_scan edges to SurveyWiFiScan.
func (squ *SurveyQuestionUpdate) AddWifiScan(s ...*SurveyWiFiScan) *SurveyQuestionUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squ.AddWifiScanIDs(ids...)
}

// AddCellScanIDs adds the cell_scan edge to SurveyCellScan by ids.
func (squ *SurveyQuestionUpdate) AddCellScanIDs(ids ...string) *SurveyQuestionUpdate {
	if squ.cell_scan == nil {
		squ.cell_scan = make(map[string]struct{})
	}
	for i := range ids {
		squ.cell_scan[ids[i]] = struct{}{}
	}
	return squ
}

// AddCellScan adds the cell_scan edges to SurveyCellScan.
func (squ *SurveyQuestionUpdate) AddCellScan(s ...*SurveyCellScan) *SurveyQuestionUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squ.AddCellScanIDs(ids...)
}

// AddPhotoDatumIDs adds the photo_data edge to File by ids.
func (squ *SurveyQuestionUpdate) AddPhotoDatumIDs(ids ...string) *SurveyQuestionUpdate {
	if squ.photo_data == nil {
		squ.photo_data = make(map[string]struct{})
	}
	for i := range ids {
		squ.photo_data[ids[i]] = struct{}{}
	}
	return squ
}

// AddPhotoData adds the photo_data edges to File.
func (squ *SurveyQuestionUpdate) AddPhotoData(f ...*File) *SurveyQuestionUpdate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return squ.AddPhotoDatumIDs(ids...)
}

// ClearSurvey clears the survey edge to Survey.
func (squ *SurveyQuestionUpdate) ClearSurvey() *SurveyQuestionUpdate {
	squ.clearedSurvey = true
	return squ
}

// RemoveWifiScanIDs removes the wifi_scan edge to SurveyWiFiScan by ids.
func (squ *SurveyQuestionUpdate) RemoveWifiScanIDs(ids ...string) *SurveyQuestionUpdate {
	if squ.removedWifiScan == nil {
		squ.removedWifiScan = make(map[string]struct{})
	}
	for i := range ids {
		squ.removedWifiScan[ids[i]] = struct{}{}
	}
	return squ
}

// RemoveWifiScan removes wifi_scan edges to SurveyWiFiScan.
func (squ *SurveyQuestionUpdate) RemoveWifiScan(s ...*SurveyWiFiScan) *SurveyQuestionUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squ.RemoveWifiScanIDs(ids...)
}

// RemoveCellScanIDs removes the cell_scan edge to SurveyCellScan by ids.
func (squ *SurveyQuestionUpdate) RemoveCellScanIDs(ids ...string) *SurveyQuestionUpdate {
	if squ.removedCellScan == nil {
		squ.removedCellScan = make(map[string]struct{})
	}
	for i := range ids {
		squ.removedCellScan[ids[i]] = struct{}{}
	}
	return squ
}

// RemoveCellScan removes cell_scan edges to SurveyCellScan.
func (squ *SurveyQuestionUpdate) RemoveCellScan(s ...*SurveyCellScan) *SurveyQuestionUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squ.RemoveCellScanIDs(ids...)
}

// RemovePhotoDatumIDs removes the photo_data edge to File by ids.
func (squ *SurveyQuestionUpdate) RemovePhotoDatumIDs(ids ...string) *SurveyQuestionUpdate {
	if squ.removedPhotoData == nil {
		squ.removedPhotoData = make(map[string]struct{})
	}
	for i := range ids {
		squ.removedPhotoData[ids[i]] = struct{}{}
	}
	return squ
}

// RemovePhotoData removes photo_data edges to File.
func (squ *SurveyQuestionUpdate) RemovePhotoData(f ...*File) *SurveyQuestionUpdate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return squ.RemovePhotoDatumIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (squ *SurveyQuestionUpdate) Save(ctx context.Context) (int, error) {
	if squ.update_time == nil {
		v := surveyquestion.UpdateDefaultUpdateTime()
		squ.update_time = &v
	}
	if len(squ.survey) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"survey\"")
	}
	if squ.clearedSurvey && squ.survey == nil {
		return 0, errors.New("ent: clearing a unique edge \"survey\"")
	}
	return squ.sqlSave(ctx)
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
				Type:   field.TypeString,
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
	if value := squ.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveyquestion.FieldUpdateTime,
		})
	}
	if value := squ.form_name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldFormName,
		})
	}
	if squ.clearform_name {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldFormName,
		})
	}
	if value := squ.form_description; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldFormDescription,
		})
	}
	if squ.clearform_description {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldFormDescription,
		})
	}
	if value := squ.form_index; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveyquestion.FieldFormIndex,
		})
	}
	if value := squ.addform_index; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveyquestion.FieldFormIndex,
		})
	}
	if value := squ.question_type; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldQuestionType,
		})
	}
	if squ.clearquestion_type {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldQuestionType,
		})
	}
	if value := squ.question_format; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldQuestionFormat,
		})
	}
	if squ.clearquestion_format {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldQuestionFormat,
		})
	}
	if value := squ.question_text; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldQuestionText,
		})
	}
	if squ.clearquestion_text {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldQuestionText,
		})
	}
	if value := squ.question_index; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveyquestion.FieldQuestionIndex,
		})
	}
	if value := squ.addquestion_index; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveyquestion.FieldQuestionIndex,
		})
	}
	if value := squ.bool_data; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: surveyquestion.FieldBoolData,
		})
	}
	if squ.clearbool_data {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: surveyquestion.FieldBoolData,
		})
	}
	if value := squ.email_data; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldEmailData,
		})
	}
	if squ.clearemail_data {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldEmailData,
		})
	}
	if value := squ.latitude; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldLatitude,
		})
	}
	if value := squ.addlatitude; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldLatitude,
		})
	}
	if squ.clearlatitude {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldLatitude,
		})
	}
	if value := squ.longitude; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldLongitude,
		})
	}
	if value := squ.addlongitude; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldLongitude,
		})
	}
	if squ.clearlongitude {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldLongitude,
		})
	}
	if value := squ.location_accuracy; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldLocationAccuracy,
		})
	}
	if value := squ.addlocation_accuracy; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldLocationAccuracy,
		})
	}
	if squ.clearlocation_accuracy {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldLocationAccuracy,
		})
	}
	if value := squ.altitude; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldAltitude,
		})
	}
	if value := squ.addaltitude; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldAltitude,
		})
	}
	if squ.clearaltitude {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldAltitude,
		})
	}
	if value := squ.phone_data; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldPhoneData,
		})
	}
	if squ.clearphone_data {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldPhoneData,
		})
	}
	if value := squ.text_data; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldTextData,
		})
	}
	if squ.cleartext_data {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldTextData,
		})
	}
	if value := squ.float_data; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldFloatData,
		})
	}
	if value := squ.addfloat_data; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldFloatData,
		})
	}
	if squ.clearfloat_data {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldFloatData,
		})
	}
	if value := squ.int_data; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveyquestion.FieldIntData,
		})
	}
	if value := squ.addint_data; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveyquestion.FieldIntData,
		})
	}
	if squ.clearint_data {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveyquestion.FieldIntData,
		})
	}
	if value := squ.date_data; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveyquestion.FieldDateData,
		})
	}
	if squ.cleardate_data {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: surveyquestion.FieldDateData,
		})
	}
	if squ.clearedSurvey {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveyquestion.SurveyTable,
			Columns: []string{surveyquestion.SurveyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: survey.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squ.survey; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveyquestion.SurveyTable,
			Columns: []string{surveyquestion.SurveyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: survey.FieldID,
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
	if nodes := squ.removedWifiScan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   surveyquestion.WifiScanTable,
			Columns: []string{surveyquestion.WifiScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveywifiscan.FieldID,
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squ.wifi_scan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   surveyquestion.WifiScanTable,
			Columns: []string{surveyquestion.WifiScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveywifiscan.FieldID,
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
	if nodes := squ.removedCellScan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   surveyquestion.CellScanTable,
			Columns: []string{surveyquestion.CellScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveycellscan.FieldID,
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squ.cell_scan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   surveyquestion.CellScanTable,
			Columns: []string{surveyquestion.CellScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveycellscan.FieldID,
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
	if nodes := squ.removedPhotoData; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   surveyquestion.PhotoDataTable,
			Columns: []string{surveyquestion.PhotoDataColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: file.FieldID,
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squ.photo_data; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   surveyquestion.PhotoDataTable,
			Columns: []string{surveyquestion.PhotoDataColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: file.FieldID,
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
	id string

	update_time            *time.Time
	form_name              *string
	clearform_name         bool
	form_description       *string
	clearform_description  bool
	form_index             *int
	addform_index          *int
	question_type          *string
	clearquestion_type     bool
	question_format        *string
	clearquestion_format   bool
	question_text          *string
	clearquestion_text     bool
	question_index         *int
	addquestion_index      *int
	bool_data              *bool
	clearbool_data         bool
	email_data             *string
	clearemail_data        bool
	latitude               *float64
	addlatitude            *float64
	clearlatitude          bool
	longitude              *float64
	addlongitude           *float64
	clearlongitude         bool
	location_accuracy      *float64
	addlocation_accuracy   *float64
	clearlocation_accuracy bool
	altitude               *float64
	addaltitude            *float64
	clearaltitude          bool
	phone_data             *string
	clearphone_data        bool
	text_data              *string
	cleartext_data         bool
	float_data             *float64
	addfloat_data          *float64
	clearfloat_data        bool
	int_data               *int
	addint_data            *int
	clearint_data          bool
	date_data              *time.Time
	cleardate_data         bool
	survey                 map[string]struct{}
	wifi_scan              map[string]struct{}
	cell_scan              map[string]struct{}
	photo_data             map[string]struct{}
	clearedSurvey          bool
	removedWifiScan        map[string]struct{}
	removedCellScan        map[string]struct{}
	removedPhotoData       map[string]struct{}
}

// SetFormName sets the form_name field.
func (squo *SurveyQuestionUpdateOne) SetFormName(s string) *SurveyQuestionUpdateOne {
	squo.form_name = &s
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
	squo.form_name = nil
	squo.clearform_name = true
	return squo
}

// SetFormDescription sets the form_description field.
func (squo *SurveyQuestionUpdateOne) SetFormDescription(s string) *SurveyQuestionUpdateOne {
	squo.form_description = &s
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
	squo.form_description = nil
	squo.clearform_description = true
	return squo
}

// SetFormIndex sets the form_index field.
func (squo *SurveyQuestionUpdateOne) SetFormIndex(i int) *SurveyQuestionUpdateOne {
	squo.form_index = &i
	squo.addform_index = nil
	return squo
}

// AddFormIndex adds i to form_index.
func (squo *SurveyQuestionUpdateOne) AddFormIndex(i int) *SurveyQuestionUpdateOne {
	if squo.addform_index == nil {
		squo.addform_index = &i
	} else {
		*squo.addform_index += i
	}
	return squo
}

// SetQuestionType sets the question_type field.
func (squo *SurveyQuestionUpdateOne) SetQuestionType(s string) *SurveyQuestionUpdateOne {
	squo.question_type = &s
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
	squo.question_type = nil
	squo.clearquestion_type = true
	return squo
}

// SetQuestionFormat sets the question_format field.
func (squo *SurveyQuestionUpdateOne) SetQuestionFormat(s string) *SurveyQuestionUpdateOne {
	squo.question_format = &s
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
	squo.question_format = nil
	squo.clearquestion_format = true
	return squo
}

// SetQuestionText sets the question_text field.
func (squo *SurveyQuestionUpdateOne) SetQuestionText(s string) *SurveyQuestionUpdateOne {
	squo.question_text = &s
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
	squo.question_text = nil
	squo.clearquestion_text = true
	return squo
}

// SetQuestionIndex sets the question_index field.
func (squo *SurveyQuestionUpdateOne) SetQuestionIndex(i int) *SurveyQuestionUpdateOne {
	squo.question_index = &i
	squo.addquestion_index = nil
	return squo
}

// AddQuestionIndex adds i to question_index.
func (squo *SurveyQuestionUpdateOne) AddQuestionIndex(i int) *SurveyQuestionUpdateOne {
	if squo.addquestion_index == nil {
		squo.addquestion_index = &i
	} else {
		*squo.addquestion_index += i
	}
	return squo
}

// SetBoolData sets the bool_data field.
func (squo *SurveyQuestionUpdateOne) SetBoolData(b bool) *SurveyQuestionUpdateOne {
	squo.bool_data = &b
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
	squo.bool_data = nil
	squo.clearbool_data = true
	return squo
}

// SetEmailData sets the email_data field.
func (squo *SurveyQuestionUpdateOne) SetEmailData(s string) *SurveyQuestionUpdateOne {
	squo.email_data = &s
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
	squo.email_data = nil
	squo.clearemail_data = true
	return squo
}

// SetLatitude sets the latitude field.
func (squo *SurveyQuestionUpdateOne) SetLatitude(f float64) *SurveyQuestionUpdateOne {
	squo.latitude = &f
	squo.addlatitude = nil
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
	if squo.addlatitude == nil {
		squo.addlatitude = &f
	} else {
		*squo.addlatitude += f
	}
	return squo
}

// ClearLatitude clears the value of latitude.
func (squo *SurveyQuestionUpdateOne) ClearLatitude() *SurveyQuestionUpdateOne {
	squo.latitude = nil
	squo.clearlatitude = true
	return squo
}

// SetLongitude sets the longitude field.
func (squo *SurveyQuestionUpdateOne) SetLongitude(f float64) *SurveyQuestionUpdateOne {
	squo.longitude = &f
	squo.addlongitude = nil
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
	if squo.addlongitude == nil {
		squo.addlongitude = &f
	} else {
		*squo.addlongitude += f
	}
	return squo
}

// ClearLongitude clears the value of longitude.
func (squo *SurveyQuestionUpdateOne) ClearLongitude() *SurveyQuestionUpdateOne {
	squo.longitude = nil
	squo.clearlongitude = true
	return squo
}

// SetLocationAccuracy sets the location_accuracy field.
func (squo *SurveyQuestionUpdateOne) SetLocationAccuracy(f float64) *SurveyQuestionUpdateOne {
	squo.location_accuracy = &f
	squo.addlocation_accuracy = nil
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
	if squo.addlocation_accuracy == nil {
		squo.addlocation_accuracy = &f
	} else {
		*squo.addlocation_accuracy += f
	}
	return squo
}

// ClearLocationAccuracy clears the value of location_accuracy.
func (squo *SurveyQuestionUpdateOne) ClearLocationAccuracy() *SurveyQuestionUpdateOne {
	squo.location_accuracy = nil
	squo.clearlocation_accuracy = true
	return squo
}

// SetAltitude sets the altitude field.
func (squo *SurveyQuestionUpdateOne) SetAltitude(f float64) *SurveyQuestionUpdateOne {
	squo.altitude = &f
	squo.addaltitude = nil
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
	if squo.addaltitude == nil {
		squo.addaltitude = &f
	} else {
		*squo.addaltitude += f
	}
	return squo
}

// ClearAltitude clears the value of altitude.
func (squo *SurveyQuestionUpdateOne) ClearAltitude() *SurveyQuestionUpdateOne {
	squo.altitude = nil
	squo.clearaltitude = true
	return squo
}

// SetPhoneData sets the phone_data field.
func (squo *SurveyQuestionUpdateOne) SetPhoneData(s string) *SurveyQuestionUpdateOne {
	squo.phone_data = &s
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
	squo.phone_data = nil
	squo.clearphone_data = true
	return squo
}

// SetTextData sets the text_data field.
func (squo *SurveyQuestionUpdateOne) SetTextData(s string) *SurveyQuestionUpdateOne {
	squo.text_data = &s
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
	squo.text_data = nil
	squo.cleartext_data = true
	return squo
}

// SetFloatData sets the float_data field.
func (squo *SurveyQuestionUpdateOne) SetFloatData(f float64) *SurveyQuestionUpdateOne {
	squo.float_data = &f
	squo.addfloat_data = nil
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
	if squo.addfloat_data == nil {
		squo.addfloat_data = &f
	} else {
		*squo.addfloat_data += f
	}
	return squo
}

// ClearFloatData clears the value of float_data.
func (squo *SurveyQuestionUpdateOne) ClearFloatData() *SurveyQuestionUpdateOne {
	squo.float_data = nil
	squo.clearfloat_data = true
	return squo
}

// SetIntData sets the int_data field.
func (squo *SurveyQuestionUpdateOne) SetIntData(i int) *SurveyQuestionUpdateOne {
	squo.int_data = &i
	squo.addint_data = nil
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
	if squo.addint_data == nil {
		squo.addint_data = &i
	} else {
		*squo.addint_data += i
	}
	return squo
}

// ClearIntData clears the value of int_data.
func (squo *SurveyQuestionUpdateOne) ClearIntData() *SurveyQuestionUpdateOne {
	squo.int_data = nil
	squo.clearint_data = true
	return squo
}

// SetDateData sets the date_data field.
func (squo *SurveyQuestionUpdateOne) SetDateData(t time.Time) *SurveyQuestionUpdateOne {
	squo.date_data = &t
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
	squo.date_data = nil
	squo.cleardate_data = true
	return squo
}

// SetSurveyID sets the survey edge to Survey by id.
func (squo *SurveyQuestionUpdateOne) SetSurveyID(id string) *SurveyQuestionUpdateOne {
	if squo.survey == nil {
		squo.survey = make(map[string]struct{})
	}
	squo.survey[id] = struct{}{}
	return squo
}

// SetSurvey sets the survey edge to Survey.
func (squo *SurveyQuestionUpdateOne) SetSurvey(s *Survey) *SurveyQuestionUpdateOne {
	return squo.SetSurveyID(s.ID)
}

// AddWifiScanIDs adds the wifi_scan edge to SurveyWiFiScan by ids.
func (squo *SurveyQuestionUpdateOne) AddWifiScanIDs(ids ...string) *SurveyQuestionUpdateOne {
	if squo.wifi_scan == nil {
		squo.wifi_scan = make(map[string]struct{})
	}
	for i := range ids {
		squo.wifi_scan[ids[i]] = struct{}{}
	}
	return squo
}

// AddWifiScan adds the wifi_scan edges to SurveyWiFiScan.
func (squo *SurveyQuestionUpdateOne) AddWifiScan(s ...*SurveyWiFiScan) *SurveyQuestionUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squo.AddWifiScanIDs(ids...)
}

// AddCellScanIDs adds the cell_scan edge to SurveyCellScan by ids.
func (squo *SurveyQuestionUpdateOne) AddCellScanIDs(ids ...string) *SurveyQuestionUpdateOne {
	if squo.cell_scan == nil {
		squo.cell_scan = make(map[string]struct{})
	}
	for i := range ids {
		squo.cell_scan[ids[i]] = struct{}{}
	}
	return squo
}

// AddCellScan adds the cell_scan edges to SurveyCellScan.
func (squo *SurveyQuestionUpdateOne) AddCellScan(s ...*SurveyCellScan) *SurveyQuestionUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squo.AddCellScanIDs(ids...)
}

// AddPhotoDatumIDs adds the photo_data edge to File by ids.
func (squo *SurveyQuestionUpdateOne) AddPhotoDatumIDs(ids ...string) *SurveyQuestionUpdateOne {
	if squo.photo_data == nil {
		squo.photo_data = make(map[string]struct{})
	}
	for i := range ids {
		squo.photo_data[ids[i]] = struct{}{}
	}
	return squo
}

// AddPhotoData adds the photo_data edges to File.
func (squo *SurveyQuestionUpdateOne) AddPhotoData(f ...*File) *SurveyQuestionUpdateOne {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return squo.AddPhotoDatumIDs(ids...)
}

// ClearSurvey clears the survey edge to Survey.
func (squo *SurveyQuestionUpdateOne) ClearSurvey() *SurveyQuestionUpdateOne {
	squo.clearedSurvey = true
	return squo
}

// RemoveWifiScanIDs removes the wifi_scan edge to SurveyWiFiScan by ids.
func (squo *SurveyQuestionUpdateOne) RemoveWifiScanIDs(ids ...string) *SurveyQuestionUpdateOne {
	if squo.removedWifiScan == nil {
		squo.removedWifiScan = make(map[string]struct{})
	}
	for i := range ids {
		squo.removedWifiScan[ids[i]] = struct{}{}
	}
	return squo
}

// RemoveWifiScan removes wifi_scan edges to SurveyWiFiScan.
func (squo *SurveyQuestionUpdateOne) RemoveWifiScan(s ...*SurveyWiFiScan) *SurveyQuestionUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squo.RemoveWifiScanIDs(ids...)
}

// RemoveCellScanIDs removes the cell_scan edge to SurveyCellScan by ids.
func (squo *SurveyQuestionUpdateOne) RemoveCellScanIDs(ids ...string) *SurveyQuestionUpdateOne {
	if squo.removedCellScan == nil {
		squo.removedCellScan = make(map[string]struct{})
	}
	for i := range ids {
		squo.removedCellScan[ids[i]] = struct{}{}
	}
	return squo
}

// RemoveCellScan removes cell_scan edges to SurveyCellScan.
func (squo *SurveyQuestionUpdateOne) RemoveCellScan(s ...*SurveyCellScan) *SurveyQuestionUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squo.RemoveCellScanIDs(ids...)
}

// RemovePhotoDatumIDs removes the photo_data edge to File by ids.
func (squo *SurveyQuestionUpdateOne) RemovePhotoDatumIDs(ids ...string) *SurveyQuestionUpdateOne {
	if squo.removedPhotoData == nil {
		squo.removedPhotoData = make(map[string]struct{})
	}
	for i := range ids {
		squo.removedPhotoData[ids[i]] = struct{}{}
	}
	return squo
}

// RemovePhotoData removes photo_data edges to File.
func (squo *SurveyQuestionUpdateOne) RemovePhotoData(f ...*File) *SurveyQuestionUpdateOne {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return squo.RemovePhotoDatumIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (squo *SurveyQuestionUpdateOne) Save(ctx context.Context) (*SurveyQuestion, error) {
	if squo.update_time == nil {
		v := surveyquestion.UpdateDefaultUpdateTime()
		squo.update_time = &v
	}
	if len(squo.survey) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"survey\"")
	}
	if squo.clearedSurvey && squo.survey == nil {
		return nil, errors.New("ent: clearing a unique edge \"survey\"")
	}
	return squo.sqlSave(ctx)
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
				Value:  squo.id,
				Type:   field.TypeString,
				Column: surveyquestion.FieldID,
			},
		},
	}
	if value := squo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveyquestion.FieldUpdateTime,
		})
	}
	if value := squo.form_name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldFormName,
		})
	}
	if squo.clearform_name {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldFormName,
		})
	}
	if value := squo.form_description; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldFormDescription,
		})
	}
	if squo.clearform_description {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldFormDescription,
		})
	}
	if value := squo.form_index; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveyquestion.FieldFormIndex,
		})
	}
	if value := squo.addform_index; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveyquestion.FieldFormIndex,
		})
	}
	if value := squo.question_type; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldQuestionType,
		})
	}
	if squo.clearquestion_type {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldQuestionType,
		})
	}
	if value := squo.question_format; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldQuestionFormat,
		})
	}
	if squo.clearquestion_format {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldQuestionFormat,
		})
	}
	if value := squo.question_text; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldQuestionText,
		})
	}
	if squo.clearquestion_text {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldQuestionText,
		})
	}
	if value := squo.question_index; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveyquestion.FieldQuestionIndex,
		})
	}
	if value := squo.addquestion_index; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveyquestion.FieldQuestionIndex,
		})
	}
	if value := squo.bool_data; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: surveyquestion.FieldBoolData,
		})
	}
	if squo.clearbool_data {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: surveyquestion.FieldBoolData,
		})
	}
	if value := squo.email_data; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldEmailData,
		})
	}
	if squo.clearemail_data {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldEmailData,
		})
	}
	if value := squo.latitude; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldLatitude,
		})
	}
	if value := squo.addlatitude; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldLatitude,
		})
	}
	if squo.clearlatitude {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldLatitude,
		})
	}
	if value := squo.longitude; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldLongitude,
		})
	}
	if value := squo.addlongitude; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldLongitude,
		})
	}
	if squo.clearlongitude {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldLongitude,
		})
	}
	if value := squo.location_accuracy; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldLocationAccuracy,
		})
	}
	if value := squo.addlocation_accuracy; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldLocationAccuracy,
		})
	}
	if squo.clearlocation_accuracy {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldLocationAccuracy,
		})
	}
	if value := squo.altitude; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldAltitude,
		})
	}
	if value := squo.addaltitude; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldAltitude,
		})
	}
	if squo.clearaltitude {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldAltitude,
		})
	}
	if value := squo.phone_data; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldPhoneData,
		})
	}
	if squo.clearphone_data {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldPhoneData,
		})
	}
	if value := squo.text_data; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldTextData,
		})
	}
	if squo.cleartext_data {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: surveyquestion.FieldTextData,
		})
	}
	if value := squo.float_data; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldFloatData,
		})
	}
	if value := squo.addfloat_data; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldFloatData,
		})
	}
	if squo.clearfloat_data {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: surveyquestion.FieldFloatData,
		})
	}
	if value := squo.int_data; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveyquestion.FieldIntData,
		})
	}
	if value := squo.addint_data; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveyquestion.FieldIntData,
		})
	}
	if squo.clearint_data {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: surveyquestion.FieldIntData,
		})
	}
	if value := squo.date_data; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveyquestion.FieldDateData,
		})
	}
	if squo.cleardate_data {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: surveyquestion.FieldDateData,
		})
	}
	if squo.clearedSurvey {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveyquestion.SurveyTable,
			Columns: []string{surveyquestion.SurveyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: survey.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squo.survey; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   surveyquestion.SurveyTable,
			Columns: []string{surveyquestion.SurveyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: survey.FieldID,
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
	if nodes := squo.removedWifiScan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   surveyquestion.WifiScanTable,
			Columns: []string{surveyquestion.WifiScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveywifiscan.FieldID,
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squo.wifi_scan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   surveyquestion.WifiScanTable,
			Columns: []string{surveyquestion.WifiScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveywifiscan.FieldID,
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
	if nodes := squo.removedCellScan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   surveyquestion.CellScanTable,
			Columns: []string{surveyquestion.CellScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveycellscan.FieldID,
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squo.cell_scan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   surveyquestion.CellScanTable,
			Columns: []string{surveyquestion.CellScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveycellscan.FieldID,
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
	if nodes := squo.removedPhotoData; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   surveyquestion.PhotoDataTable,
			Columns: []string{surveyquestion.PhotoDataColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: file.FieldID,
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squo.photo_data; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   surveyquestion.PhotoDataTable,
			Columns: []string{surveyquestion.PhotoDataColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: file.FieldID,
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
