// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
)

// SurveyQuestionCreate is the builder for creating a SurveyQuestion entity.
type SurveyQuestionCreate struct {
	config
	create_time       *time.Time
	update_time       *time.Time
	form_name         *string
	form_description  *string
	form_index        *int
	question_type     *string
	question_format   *string
	question_text     *string
	question_index    *int
	bool_data         *bool
	email_data        *string
	latitude          *float64
	longitude         *float64
	location_accuracy *float64
	altitude          *float64
	phone_data        *string
	text_data         *string
	float_data        *float64
	int_data          *int
	date_data         *time.Time
	survey            map[string]struct{}
	wifi_scan         map[string]struct{}
	cell_scan         map[string]struct{}
	photo_data        map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (sqc *SurveyQuestionCreate) SetCreateTime(t time.Time) *SurveyQuestionCreate {
	sqc.create_time = &t
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
	sqc.update_time = &t
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
	sqc.form_name = &s
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
	sqc.form_description = &s
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
	sqc.form_index = &i
	return sqc
}

// SetQuestionType sets the question_type field.
func (sqc *SurveyQuestionCreate) SetQuestionType(s string) *SurveyQuestionCreate {
	sqc.question_type = &s
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
	sqc.question_format = &s
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
	sqc.question_text = &s
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
	sqc.question_index = &i
	return sqc
}

// SetBoolData sets the bool_data field.
func (sqc *SurveyQuestionCreate) SetBoolData(b bool) *SurveyQuestionCreate {
	sqc.bool_data = &b
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
	sqc.email_data = &s
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
	sqc.latitude = &f
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
	sqc.longitude = &f
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
	sqc.location_accuracy = &f
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
	sqc.altitude = &f
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
	sqc.phone_data = &s
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
	sqc.text_data = &s
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
	sqc.float_data = &f
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
	sqc.int_data = &i
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
	sqc.date_data = &t
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
func (sqc *SurveyQuestionCreate) SetSurveyID(id string) *SurveyQuestionCreate {
	if sqc.survey == nil {
		sqc.survey = make(map[string]struct{})
	}
	sqc.survey[id] = struct{}{}
	return sqc
}

// SetSurvey sets the survey edge to Survey.
func (sqc *SurveyQuestionCreate) SetSurvey(s *Survey) *SurveyQuestionCreate {
	return sqc.SetSurveyID(s.ID)
}

// AddWifiScanIDs adds the wifi_scan edge to SurveyWiFiScan by ids.
func (sqc *SurveyQuestionCreate) AddWifiScanIDs(ids ...string) *SurveyQuestionCreate {
	if sqc.wifi_scan == nil {
		sqc.wifi_scan = make(map[string]struct{})
	}
	for i := range ids {
		sqc.wifi_scan[ids[i]] = struct{}{}
	}
	return sqc
}

// AddWifiScan adds the wifi_scan edges to SurveyWiFiScan.
func (sqc *SurveyQuestionCreate) AddWifiScan(s ...*SurveyWiFiScan) *SurveyQuestionCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sqc.AddWifiScanIDs(ids...)
}

// AddCellScanIDs adds the cell_scan edge to SurveyCellScan by ids.
func (sqc *SurveyQuestionCreate) AddCellScanIDs(ids ...string) *SurveyQuestionCreate {
	if sqc.cell_scan == nil {
		sqc.cell_scan = make(map[string]struct{})
	}
	for i := range ids {
		sqc.cell_scan[ids[i]] = struct{}{}
	}
	return sqc
}

// AddCellScan adds the cell_scan edges to SurveyCellScan.
func (sqc *SurveyQuestionCreate) AddCellScan(s ...*SurveyCellScan) *SurveyQuestionCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sqc.AddCellScanIDs(ids...)
}

// AddPhotoDatumIDs adds the photo_data edge to File by ids.
func (sqc *SurveyQuestionCreate) AddPhotoDatumIDs(ids ...string) *SurveyQuestionCreate {
	if sqc.photo_data == nil {
		sqc.photo_data = make(map[string]struct{})
	}
	for i := range ids {
		sqc.photo_data[ids[i]] = struct{}{}
	}
	return sqc
}

// AddPhotoData adds the photo_data edges to File.
func (sqc *SurveyQuestionCreate) AddPhotoData(f ...*File) *SurveyQuestionCreate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return sqc.AddPhotoDatumIDs(ids...)
}

// Save creates the SurveyQuestion in the database.
func (sqc *SurveyQuestionCreate) Save(ctx context.Context) (*SurveyQuestion, error) {
	if sqc.create_time == nil {
		v := surveyquestion.DefaultCreateTime()
		sqc.create_time = &v
	}
	if sqc.update_time == nil {
		v := surveyquestion.DefaultUpdateTime()
		sqc.update_time = &v
	}
	if sqc.form_index == nil {
		return nil, errors.New("ent: missing required field \"form_index\"")
	}
	if sqc.question_index == nil {
		return nil, errors.New("ent: missing required field \"question_index\"")
	}
	if len(sqc.survey) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"survey\"")
	}
	if sqc.survey == nil {
		return nil, errors.New("ent: missing required edge \"survey\"")
	}
	return sqc.sqlSave(ctx)
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
		res     sql.Result
		builder = sql.Dialect(sqc.driver.Dialect())
		sq      = &SurveyQuestion{config: sqc.config}
	)
	tx, err := sqc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(surveyquestion.Table).Default()
	if value := sqc.create_time; value != nil {
		insert.Set(surveyquestion.FieldCreateTime, *value)
		sq.CreateTime = *value
	}
	if value := sqc.update_time; value != nil {
		insert.Set(surveyquestion.FieldUpdateTime, *value)
		sq.UpdateTime = *value
	}
	if value := sqc.form_name; value != nil {
		insert.Set(surveyquestion.FieldFormName, *value)
		sq.FormName = *value
	}
	if value := sqc.form_description; value != nil {
		insert.Set(surveyquestion.FieldFormDescription, *value)
		sq.FormDescription = *value
	}
	if value := sqc.form_index; value != nil {
		insert.Set(surveyquestion.FieldFormIndex, *value)
		sq.FormIndex = *value
	}
	if value := sqc.question_type; value != nil {
		insert.Set(surveyquestion.FieldQuestionType, *value)
		sq.QuestionType = *value
	}
	if value := sqc.question_format; value != nil {
		insert.Set(surveyquestion.FieldQuestionFormat, *value)
		sq.QuestionFormat = *value
	}
	if value := sqc.question_text; value != nil {
		insert.Set(surveyquestion.FieldQuestionText, *value)
		sq.QuestionText = *value
	}
	if value := sqc.question_index; value != nil {
		insert.Set(surveyquestion.FieldQuestionIndex, *value)
		sq.QuestionIndex = *value
	}
	if value := sqc.bool_data; value != nil {
		insert.Set(surveyquestion.FieldBoolData, *value)
		sq.BoolData = *value
	}
	if value := sqc.email_data; value != nil {
		insert.Set(surveyquestion.FieldEmailData, *value)
		sq.EmailData = *value
	}
	if value := sqc.latitude; value != nil {
		insert.Set(surveyquestion.FieldLatitude, *value)
		sq.Latitude = *value
	}
	if value := sqc.longitude; value != nil {
		insert.Set(surveyquestion.FieldLongitude, *value)
		sq.Longitude = *value
	}
	if value := sqc.location_accuracy; value != nil {
		insert.Set(surveyquestion.FieldLocationAccuracy, *value)
		sq.LocationAccuracy = *value
	}
	if value := sqc.altitude; value != nil {
		insert.Set(surveyquestion.FieldAltitude, *value)
		sq.Altitude = *value
	}
	if value := sqc.phone_data; value != nil {
		insert.Set(surveyquestion.FieldPhoneData, *value)
		sq.PhoneData = *value
	}
	if value := sqc.text_data; value != nil {
		insert.Set(surveyquestion.FieldTextData, *value)
		sq.TextData = *value
	}
	if value := sqc.float_data; value != nil {
		insert.Set(surveyquestion.FieldFloatData, *value)
		sq.FloatData = *value
	}
	if value := sqc.int_data; value != nil {
		insert.Set(surveyquestion.FieldIntData, *value)
		sq.IntData = *value
	}
	if value := sqc.date_data; value != nil {
		insert.Set(surveyquestion.FieldDateData, *value)
		sq.DateData = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(surveyquestion.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	sq.ID = strconv.FormatInt(id, 10)
	if len(sqc.survey) > 0 {
		for eid := range sqc.survey {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(surveyquestion.SurveyTable).
				Set(surveyquestion.SurveyColumn, eid).
				Where(sql.EQ(surveyquestion.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(sqc.wifi_scan) > 0 {
		p := sql.P()
		for eid := range sqc.wifi_scan {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(surveywifiscan.FieldID, eid)
		}
		query, args := builder.Update(surveyquestion.WifiScanTable).
			Set(surveyquestion.WifiScanColumn, id).
			Where(sql.And(p, sql.IsNull(surveyquestion.WifiScanColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(sqc.wifi_scan) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"wifi_scan\" %v already connected to a different \"SurveyQuestion\"", keys(sqc.wifi_scan))})
		}
	}
	if len(sqc.cell_scan) > 0 {
		p := sql.P()
		for eid := range sqc.cell_scan {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(surveycellscan.FieldID, eid)
		}
		query, args := builder.Update(surveyquestion.CellScanTable).
			Set(surveyquestion.CellScanColumn, id).
			Where(sql.And(p, sql.IsNull(surveyquestion.CellScanColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(sqc.cell_scan) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"cell_scan\" %v already connected to a different \"SurveyQuestion\"", keys(sqc.cell_scan))})
		}
	}
	if len(sqc.photo_data) > 0 {
		p := sql.P()
		for eid := range sqc.photo_data {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(file.FieldID, eid)
		}
		query, args := builder.Update(surveyquestion.PhotoDataTable).
			Set(surveyquestion.PhotoDataColumn, id).
			Where(sql.And(p, sql.IsNull(surveyquestion.PhotoDataColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(sqc.photo_data) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"photo_data\" %v already connected to a different \"SurveyQuestion\"", keys(sqc.photo_data))})
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return sq, nil
}
