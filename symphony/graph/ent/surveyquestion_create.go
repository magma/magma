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
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/survey"
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
		sq   = &SurveyQuestion{config: sqc.config}
		spec = &sqlgraph.CreateSpec{
			Table: surveyquestion.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: surveyquestion.FieldID,
			},
		}
	)
	if value := sqc.create_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveyquestion.FieldCreateTime,
		})
		sq.CreateTime = *value
	}
	if value := sqc.update_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveyquestion.FieldUpdateTime,
		})
		sq.UpdateTime = *value
	}
	if value := sqc.form_name; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldFormName,
		})
		sq.FormName = *value
	}
	if value := sqc.form_description; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldFormDescription,
		})
		sq.FormDescription = *value
	}
	if value := sqc.form_index; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveyquestion.FieldFormIndex,
		})
		sq.FormIndex = *value
	}
	if value := sqc.question_type; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldQuestionType,
		})
		sq.QuestionType = *value
	}
	if value := sqc.question_format; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldQuestionFormat,
		})
		sq.QuestionFormat = *value
	}
	if value := sqc.question_text; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldQuestionText,
		})
		sq.QuestionText = *value
	}
	if value := sqc.question_index; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveyquestion.FieldQuestionIndex,
		})
		sq.QuestionIndex = *value
	}
	if value := sqc.bool_data; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: surveyquestion.FieldBoolData,
		})
		sq.BoolData = *value
	}
	if value := sqc.email_data; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldEmailData,
		})
		sq.EmailData = *value
	}
	if value := sqc.latitude; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldLatitude,
		})
		sq.Latitude = *value
	}
	if value := sqc.longitude; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldLongitude,
		})
		sq.Longitude = *value
	}
	if value := sqc.location_accuracy; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldLocationAccuracy,
		})
		sq.LocationAccuracy = *value
	}
	if value := sqc.altitude; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldAltitude,
		})
		sq.Altitude = *value
	}
	if value := sqc.phone_data; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldPhoneData,
		})
		sq.PhoneData = *value
	}
	if value := sqc.text_data; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveyquestion.FieldTextData,
		})
		sq.TextData = *value
	}
	if value := sqc.float_data; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: surveyquestion.FieldFloatData,
		})
		sq.FloatData = *value
	}
	if value := sqc.int_data; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveyquestion.FieldIntData,
		})
		sq.IntData = *value
	}
	if value := sqc.date_data; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveyquestion.FieldDateData,
		})
		sq.DateData = *value
	}
	if nodes := sqc.survey; len(nodes) > 0 {
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
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := sqc.wifi_scan; len(nodes) > 0 {
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
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := sqc.cell_scan; len(nodes) > 0 {
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
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := sqc.photo_data; len(nodes) > 0 {
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
		spec.Edges = append(spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, sqc.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := spec.ID.Value.(int64)
	sq.ID = strconv.FormatInt(id, 10)
	return sq, nil
}
