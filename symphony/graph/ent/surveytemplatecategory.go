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
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatecategory"
)

// SurveyTemplateCategory is the model entity for the SurveyTemplateCategory schema.
type SurveyTemplateCategory struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// CategoryTitle holds the value of the "category_title" field.
	CategoryTitle string `json:"category_title,omitempty"`
	// CategoryDescription holds the value of the "category_description" field.
	CategoryDescription string `json:"category_description,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the SurveyTemplateCategoryQuery when eager-loading is set.
	Edges                                     SurveyTemplateCategoryEdges `json:"edges"`
	location_type_survey_template_category_id *string
}

// SurveyTemplateCategoryEdges holds the relations/edges for other nodes in the graph.
type SurveyTemplateCategoryEdges struct {
	// SurveyTemplateQuestions holds the value of the survey_template_questions edge.
	SurveyTemplateQuestions []*SurveyTemplateQuestion
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// SurveyTemplateQuestionsOrErr returns the SurveyTemplateQuestions value or an error if the edge
// was not loaded in eager-loading.
func (e SurveyTemplateCategoryEdges) SurveyTemplateQuestionsOrErr() ([]*SurveyTemplateQuestion, error) {
	if e.loadedTypes[0] {
		return e.SurveyTemplateQuestions, nil
	}
	return nil, &NotLoadedError{edge: "survey_template_questions"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*SurveyTemplateCategory) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // category_title
		&sql.NullString{}, // category_description
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*SurveyTemplateCategory) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // location_type_survey_template_category_id
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the SurveyTemplateCategory fields.
func (stc *SurveyTemplateCategory) assignValues(values ...interface{}) error {
	if m, n := len(values), len(surveytemplatecategory.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	stc.ID = strconv.FormatInt(value.Int64, 10)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		stc.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		stc.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field category_title", values[2])
	} else if value.Valid {
		stc.CategoryTitle = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field category_description", values[3])
	} else if value.Valid {
		stc.CategoryDescription = value.String
	}
	values = values[4:]
	if len(values) == len(surveytemplatecategory.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field location_type_survey_template_category_id", value)
		} else if value.Valid {
			stc.location_type_survey_template_category_id = new(string)
			*stc.location_type_survey_template_category_id = strconv.FormatInt(value.Int64, 10)
		}
	}
	return nil
}

// QuerySurveyTemplateQuestions queries the survey_template_questions edge of the SurveyTemplateCategory.
func (stc *SurveyTemplateCategory) QuerySurveyTemplateQuestions() *SurveyTemplateQuestionQuery {
	return (&SurveyTemplateCategoryClient{stc.config}).QuerySurveyTemplateQuestions(stc)
}

// Update returns a builder for updating this SurveyTemplateCategory.
// Note that, you need to call SurveyTemplateCategory.Unwrap() before calling this method, if this SurveyTemplateCategory
// was returned from a transaction, and the transaction was committed or rolled back.
func (stc *SurveyTemplateCategory) Update() *SurveyTemplateCategoryUpdateOne {
	return (&SurveyTemplateCategoryClient{stc.config}).UpdateOne(stc)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (stc *SurveyTemplateCategory) Unwrap() *SurveyTemplateCategory {
	tx, ok := stc.config.driver.(*txDriver)
	if !ok {
		panic("ent: SurveyTemplateCategory is not a transactional entity")
	}
	stc.config.driver = tx.drv
	return stc
}

// String implements the fmt.Stringer.
func (stc *SurveyTemplateCategory) String() string {
	var builder strings.Builder
	builder.WriteString("SurveyTemplateCategory(")
	builder.WriteString(fmt.Sprintf("id=%v", stc.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(stc.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(stc.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", category_title=")
	builder.WriteString(stc.CategoryTitle)
	builder.WriteString(", category_description=")
	builder.WriteString(stc.CategoryDescription)
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (stc *SurveyTemplateCategory) id() int {
	id, _ := strconv.Atoi(stc.ID)
	return id
}

// SurveyTemplateCategories is a parsable slice of SurveyTemplateCategory.
type SurveyTemplateCategories []*SurveyTemplateCategory

func (stc SurveyTemplateCategories) config(cfg config) {
	for _i := range stc {
		stc[_i].config = cfg
	}
}
