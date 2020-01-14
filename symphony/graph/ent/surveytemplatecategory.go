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
}

// scanValues returns the types for scanning values from sql.Rows.
func (*SurveyTemplateCategory) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},
		&sql.NullTime{},
		&sql.NullTime{},
		&sql.NullString{},
		&sql.NullString{},
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the SurveyTemplateCategory fields.
func (stc *SurveyTemplateCategory) assignValues(values ...interface{}) error {
	if m, n := len(values), len(surveytemplatecategory.Columns); m != n {
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
