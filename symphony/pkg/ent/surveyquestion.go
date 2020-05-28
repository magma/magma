// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/pkg/ent/survey"
	"github.com/facebookincubator/symphony/pkg/ent/surveyquestion"
)

// SurveyQuestion is the model entity for the SurveyQuestion schema.
type SurveyQuestion struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
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
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the SurveyQuestionQuery when eager-loading is set.
	Edges                  SurveyQuestionEdges `json:"edges"`
	survey_question_survey *int
}

// SurveyQuestionEdges holds the relations/edges for other nodes in the graph.
type SurveyQuestionEdges struct {
	// Survey holds the value of the survey edge.
	Survey *Survey
	// WifiScan holds the value of the wifi_scan edge.
	WifiScan []*SurveyWiFiScan
	// CellScan holds the value of the cell_scan edge.
	CellScan []*SurveyCellScan
	// PhotoData holds the value of the photo_data edge.
	PhotoData []*File
	// Images holds the value of the images edge.
	Images []*File
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [5]bool
}

// SurveyOrErr returns the Survey value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e SurveyQuestionEdges) SurveyOrErr() (*Survey, error) {
	if e.loadedTypes[0] {
		if e.Survey == nil {
			// The edge survey was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: survey.Label}
		}
		return e.Survey, nil
	}
	return nil, &NotLoadedError{edge: "survey"}
}

// WifiScanOrErr returns the WifiScan value or an error if the edge
// was not loaded in eager-loading.
func (e SurveyQuestionEdges) WifiScanOrErr() ([]*SurveyWiFiScan, error) {
	if e.loadedTypes[1] {
		return e.WifiScan, nil
	}
	return nil, &NotLoadedError{edge: "wifi_scan"}
}

// CellScanOrErr returns the CellScan value or an error if the edge
// was not loaded in eager-loading.
func (e SurveyQuestionEdges) CellScanOrErr() ([]*SurveyCellScan, error) {
	if e.loadedTypes[2] {
		return e.CellScan, nil
	}
	return nil, &NotLoadedError{edge: "cell_scan"}
}

// PhotoDataOrErr returns the PhotoData value or an error if the edge
// was not loaded in eager-loading.
func (e SurveyQuestionEdges) PhotoDataOrErr() ([]*File, error) {
	if e.loadedTypes[3] {
		return e.PhotoData, nil
	}
	return nil, &NotLoadedError{edge: "photo_data"}
}

// ImagesOrErr returns the Images value or an error if the edge
// was not loaded in eager-loading.
func (e SurveyQuestionEdges) ImagesOrErr() ([]*File, error) {
	if e.loadedTypes[4] {
		return e.Images, nil
	}
	return nil, &NotLoadedError{edge: "images"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*SurveyQuestion) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},   // id
		&sql.NullTime{},    // create_time
		&sql.NullTime{},    // update_time
		&sql.NullString{},  // form_name
		&sql.NullString{},  // form_description
		&sql.NullInt64{},   // form_index
		&sql.NullString{},  // question_type
		&sql.NullString{},  // question_format
		&sql.NullString{},  // question_text
		&sql.NullInt64{},   // question_index
		&sql.NullBool{},    // bool_data
		&sql.NullString{},  // email_data
		&sql.NullFloat64{}, // latitude
		&sql.NullFloat64{}, // longitude
		&sql.NullFloat64{}, // location_accuracy
		&sql.NullFloat64{}, // altitude
		&sql.NullString{},  // phone_data
		&sql.NullString{},  // text_data
		&sql.NullFloat64{}, // float_data
		&sql.NullInt64{},   // int_data
		&sql.NullTime{},    // date_data
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*SurveyQuestion) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // survey_question_survey
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the SurveyQuestion fields.
func (sq *SurveyQuestion) assignValues(values ...interface{}) error {
	if m, n := len(values), len(surveyquestion.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	sq.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		sq.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		sq.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field form_name", values[2])
	} else if value.Valid {
		sq.FormName = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field form_description", values[3])
	} else if value.Valid {
		sq.FormDescription = value.String
	}
	if value, ok := values[4].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field form_index", values[4])
	} else if value.Valid {
		sq.FormIndex = int(value.Int64)
	}
	if value, ok := values[5].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field question_type", values[5])
	} else if value.Valid {
		sq.QuestionType = value.String
	}
	if value, ok := values[6].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field question_format", values[6])
	} else if value.Valid {
		sq.QuestionFormat = value.String
	}
	if value, ok := values[7].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field question_text", values[7])
	} else if value.Valid {
		sq.QuestionText = value.String
	}
	if value, ok := values[8].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field question_index", values[8])
	} else if value.Valid {
		sq.QuestionIndex = int(value.Int64)
	}
	if value, ok := values[9].(*sql.NullBool); !ok {
		return fmt.Errorf("unexpected type %T for field bool_data", values[9])
	} else if value.Valid {
		sq.BoolData = value.Bool
	}
	if value, ok := values[10].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field email_data", values[10])
	} else if value.Valid {
		sq.EmailData = value.String
	}
	if value, ok := values[11].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field latitude", values[11])
	} else if value.Valid {
		sq.Latitude = value.Float64
	}
	if value, ok := values[12].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field longitude", values[12])
	} else if value.Valid {
		sq.Longitude = value.Float64
	}
	if value, ok := values[13].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field location_accuracy", values[13])
	} else if value.Valid {
		sq.LocationAccuracy = value.Float64
	}
	if value, ok := values[14].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field altitude", values[14])
	} else if value.Valid {
		sq.Altitude = value.Float64
	}
	if value, ok := values[15].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field phone_data", values[15])
	} else if value.Valid {
		sq.PhoneData = value.String
	}
	if value, ok := values[16].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field text_data", values[16])
	} else if value.Valid {
		sq.TextData = value.String
	}
	if value, ok := values[17].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field float_data", values[17])
	} else if value.Valid {
		sq.FloatData = value.Float64
	}
	if value, ok := values[18].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field int_data", values[18])
	} else if value.Valid {
		sq.IntData = int(value.Int64)
	}
	if value, ok := values[19].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field date_data", values[19])
	} else if value.Valid {
		sq.DateData = value.Time
	}
	values = values[20:]
	if len(values) == len(surveyquestion.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field survey_question_survey", value)
		} else if value.Valid {
			sq.survey_question_survey = new(int)
			*sq.survey_question_survey = int(value.Int64)
		}
	}
	return nil
}

// QuerySurvey queries the survey edge of the SurveyQuestion.
func (sq *SurveyQuestion) QuerySurvey() *SurveyQuery {
	return (&SurveyQuestionClient{config: sq.config}).QuerySurvey(sq)
}

// QueryWifiScan queries the wifi_scan edge of the SurveyQuestion.
func (sq *SurveyQuestion) QueryWifiScan() *SurveyWiFiScanQuery {
	return (&SurveyQuestionClient{config: sq.config}).QueryWifiScan(sq)
}

// QueryCellScan queries the cell_scan edge of the SurveyQuestion.
func (sq *SurveyQuestion) QueryCellScan() *SurveyCellScanQuery {
	return (&SurveyQuestionClient{config: sq.config}).QueryCellScan(sq)
}

// QueryPhotoData queries the photo_data edge of the SurveyQuestion.
func (sq *SurveyQuestion) QueryPhotoData() *FileQuery {
	return (&SurveyQuestionClient{config: sq.config}).QueryPhotoData(sq)
}

// QueryImages queries the images edge of the SurveyQuestion.
func (sq *SurveyQuestion) QueryImages() *FileQuery {
	return (&SurveyQuestionClient{config: sq.config}).QueryImages(sq)
}

// Update returns a builder for updating this SurveyQuestion.
// Note that, you need to call SurveyQuestion.Unwrap() before calling this method, if this SurveyQuestion
// was returned from a transaction, and the transaction was committed or rolled back.
func (sq *SurveyQuestion) Update() *SurveyQuestionUpdateOne {
	return (&SurveyQuestionClient{config: sq.config}).UpdateOne(sq)
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

// SurveyQuestions is a parsable slice of SurveyQuestion.
type SurveyQuestions []*SurveyQuestion

func (sq SurveyQuestions) config(cfg config) {
	for _i := range sq {
		sq[_i].config = cfg
	}
}
