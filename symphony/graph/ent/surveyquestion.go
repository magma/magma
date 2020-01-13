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
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
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

// scanValues returns the types for scanning values from sql.Rows.
func (*SurveyQuestion) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},
		&sql.NullTime{},
		&sql.NullTime{},
		&sql.NullString{},
		&sql.NullString{},
		&sql.NullInt64{},
		&sql.NullString{},
		&sql.NullString{},
		&sql.NullString{},
		&sql.NullInt64{},
		&sql.NullBool{},
		&sql.NullString{},
		&sql.NullFloat64{},
		&sql.NullFloat64{},
		&sql.NullFloat64{},
		&sql.NullFloat64{},
		&sql.NullString{},
		&sql.NullString{},
		&sql.NullFloat64{},
		&sql.NullInt64{},
		&sql.NullTime{},
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the SurveyQuestion fields.
func (sq *SurveyQuestion) assignValues(values ...interface{}) error {
	if m, n := len(values), len(surveyquestion.Columns); m != n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	sq.ID = strconv.FormatInt(value.Int64, 10)
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

func (sq SurveyQuestions) config(cfg config) {
	for _i := range sq {
		sq[_i].config = cfg
	}
}
