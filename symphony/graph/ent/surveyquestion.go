// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
)

// SurveyQuestion is the model entity for the SurveyQuestion schema.
type SurveyQuestion struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// FormName holds the value of the "form_name" field.
	FormName string `json:"form_name,omitempty"`
	// FormDescription holds the value of the "form_description" field.
	FormDescription string `json:"form_description,omitempty"`
	// FormIndex holds the value of the "form_index" field.
	FormIndex int `json:"form_index,omitempty"`
	// QuestionType holds the value of the "question_type" field.
	QuestionType string `json:"question_type,omitempty"`
	// QuestionFormat holds the value of the "question_format" field.
	QuestionFormat string `json:"question_format,omitempty"`
	// QuestionText holds the value of the "question_text" field.
	QuestionText string `json:"question_text,omitempty"`
	// QuestionIndex holds the value of the "question_index" field.
	QuestionIndex int `json:"question_index,omitempty"`
	// BoolData holds the value of the "bool_data" field.
	BoolData bool `json:"bool_data,omitempty"`
	// EmailData holds the value of the "email_data" field.
	EmailData string `json:"email_data,omitempty"`
	// Latitude holds the value of the "latitude" field.
	Latitude float64 `json:"latitude,omitempty"`
	// Longitude holds the value of the "longitude" field.
	Longitude float64 `json:"longitude,omitempty"`
	// LocationAccuracy holds the value of the "location_accuracy" field.
	LocationAccuracy float64 `json:"location_accuracy,omitempty"`
	// Altitude holds the value of the "altitude" field.
	Altitude float64 `json:"altitude,omitempty"`
	// PhoneData holds the value of the "phone_data" field.
	PhoneData string `json:"phone_data,omitempty"`
	// TextData holds the value of the "text_data" field.
	TextData string `json:"text_data,omitempty"`
	// FloatData holds the value of the "float_data" field.
	FloatData float64 `json:"float_data,omitempty"`
	// IntData holds the value of the "int_data" field.
	IntData int `json:"int_data,omitempty"`
	// DateData holds the value of the "date_data" field.
	DateData time.Time `json:"date_data,omitempty"`
}

// FromRows scans the sql response data into SurveyQuestion.
func (sq *SurveyQuestion) FromRows(rows *sql.Rows) error {
	var scansq struct {
		ID               int
		CreateTime       sql.NullTime
		UpdateTime       sql.NullTime
		FormName         sql.NullString
		FormDescription  sql.NullString
		FormIndex        sql.NullInt64
		QuestionType     sql.NullString
		QuestionFormat   sql.NullString
		QuestionText     sql.NullString
		QuestionIndex    sql.NullInt64
		BoolData         sql.NullBool
		EmailData        sql.NullString
		Latitude         sql.NullFloat64
		Longitude        sql.NullFloat64
		LocationAccuracy sql.NullFloat64
		Altitude         sql.NullFloat64
		PhoneData        sql.NullString
		TextData         sql.NullString
		FloatData        sql.NullFloat64
		IntData          sql.NullInt64
		DateData         sql.NullTime
	}
	// the order here should be the same as in the `surveyquestion.Columns`.
	if err := rows.Scan(
		&scansq.ID,
		&scansq.CreateTime,
		&scansq.UpdateTime,
		&scansq.FormName,
		&scansq.FormDescription,
		&scansq.FormIndex,
		&scansq.QuestionType,
		&scansq.QuestionFormat,
		&scansq.QuestionText,
		&scansq.QuestionIndex,
		&scansq.BoolData,
		&scansq.EmailData,
		&scansq.Latitude,
		&scansq.Longitude,
		&scansq.LocationAccuracy,
		&scansq.Altitude,
		&scansq.PhoneData,
		&scansq.TextData,
		&scansq.FloatData,
		&scansq.IntData,
		&scansq.DateData,
	); err != nil {
		return err
	}
	sq.ID = strconv.Itoa(scansq.ID)
	sq.CreateTime = scansq.CreateTime.Time
	sq.UpdateTime = scansq.UpdateTime.Time
	sq.FormName = scansq.FormName.String
	sq.FormDescription = scansq.FormDescription.String
	sq.FormIndex = int(scansq.FormIndex.Int64)
	sq.QuestionType = scansq.QuestionType.String
	sq.QuestionFormat = scansq.QuestionFormat.String
	sq.QuestionText = scansq.QuestionText.String
	sq.QuestionIndex = int(scansq.QuestionIndex.Int64)
	sq.BoolData = scansq.BoolData.Bool
	sq.EmailData = scansq.EmailData.String
	sq.Latitude = scansq.Latitude.Float64
	sq.Longitude = scansq.Longitude.Float64
	sq.LocationAccuracy = scansq.LocationAccuracy.Float64
	sq.Altitude = scansq.Altitude.Float64
	sq.PhoneData = scansq.PhoneData.String
	sq.TextData = scansq.TextData.String
	sq.FloatData = scansq.FloatData.Float64
	sq.IntData = int(scansq.IntData.Int64)
	sq.DateData = scansq.DateData.Time
	return nil
}

// QuerySurvey queries the survey edge of the SurveyQuestion.
func (sq *SurveyQuestion) QuerySurvey() *SurveyQuery {
	return (&SurveyQuestionClient{sq.config}).QuerySurvey(sq)
}

// QueryWifiScan queries the wifi_scan edge of the SurveyQuestion.
func (sq *SurveyQuestion) QueryWifiScan() *SurveyWiFiScanQuery {
	return (&SurveyQuestionClient{sq.config}).QueryWifiScan(sq)
}

// QueryCellScan queries the cell_scan edge of the SurveyQuestion.
func (sq *SurveyQuestion) QueryCellScan() *SurveyCellScanQuery {
	return (&SurveyQuestionClient{sq.config}).QueryCellScan(sq)
}

// QueryPhotoData queries the photo_data edge of the SurveyQuestion.
func (sq *SurveyQuestion) QueryPhotoData() *FileQuery {
	return (&SurveyQuestionClient{sq.config}).QueryPhotoData(sq)
}

// Update returns a builder for updating this SurveyQuestion.
// Note that, you need to call SurveyQuestion.Unwrap() before calling this method, if this SurveyQuestion
// was returned from a transaction, and the transaction was committed or rolled back.
func (sq *SurveyQuestion) Update() *SurveyQuestionUpdateOne {
	return (&SurveyQuestionClient{sq.config}).UpdateOne(sq)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (sq *SurveyQuestion) Unwrap() *SurveyQuestion {
	tx, ok := sq.config.driver.(*txDriver)
	if !ok {
		panic("ent: SurveyQuestion is not a transactional entity")
	}
	sq.config.driver = tx.drv
	return sq
}

// String implements the fmt.Stringer.
func (sq *SurveyQuestion) String() string {
	var builder strings.Builder
	builder.WriteString("SurveyQuestion(")
	builder.WriteString(fmt.Sprintf("id=%v", sq.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(sq.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(sq.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", form_name=")
	builder.WriteString(sq.FormName)
	builder.WriteString(", form_description=")
	builder.WriteString(sq.FormDescription)
	builder.WriteString(", form_index=")
	builder.WriteString(fmt.Sprintf("%v", sq.FormIndex))
	builder.WriteString(", question_type=")
	builder.WriteString(sq.QuestionType)
	builder.WriteString(", question_format=")
	builder.WriteString(sq.QuestionFormat)
	builder.WriteString(", question_text=")
	builder.WriteString(sq.QuestionText)
	builder.WriteString(", question_index=")
	builder.WriteString(fmt.Sprintf("%v", sq.QuestionIndex))
	builder.WriteString(", bool_data=")
	builder.WriteString(fmt.Sprintf("%v", sq.BoolData))
	builder.WriteString(", email_data=")
	builder.WriteString(sq.EmailData)
	builder.WriteString(", latitude=")
	builder.WriteString(fmt.Sprintf("%v", sq.Latitude))
	builder.WriteString(", longitude=")
	builder.WriteString(fmt.Sprintf("%v", sq.Longitude))
	builder.WriteString(", location_accuracy=")
	builder.WriteString(fmt.Sprintf("%v", sq.LocationAccuracy))
	builder.WriteString(", altitude=")
	builder.WriteString(fmt.Sprintf("%v", sq.Altitude))
	builder.WriteString(", phone_data=")
	builder.WriteString(sq.PhoneData)
	builder.WriteString(", text_data=")
	builder.WriteString(sq.TextData)
	builder.WriteString(", float_data=")
	builder.WriteString(fmt.Sprintf("%v", sq.FloatData))
	builder.WriteString(", int_data=")
	builder.WriteString(fmt.Sprintf("%v", sq.IntData))
	builder.WriteString(", date_data=")
	builder.WriteString(sq.DateData.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (sq *SurveyQuestion) id() int {
	id, _ := strconv.Atoi(sq.ID)
	return id
}

// SurveyQuestions is a parsable slice of SurveyQuestion.
type SurveyQuestions []*SurveyQuestion

// FromRows scans the sql response data into SurveyQuestions.
func (sq *SurveyQuestions) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scansq := &SurveyQuestion{}
		if err := scansq.FromRows(rows); err != nil {
			return err
		}
		*sq = append(*sq, scansq)
	}
	return nil
}

func (sq SurveyQuestions) config(cfg config) {
	for _i := range sq {
		sq[_i].config = cfg
	}
}
