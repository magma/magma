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
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/floorplan"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// File is the model entity for the File schema.
type File struct {
	config `gqlgen:"-" json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Type holds the value of the "type" field.
	Type string `json:"type,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty" gqlgen:"fileName"`
	// Size holds the value of the "size" field.
	Size int `json:"size,omitempty" gqlgen:"sizeInBytes"`
	// ModifiedAt holds the value of the "modified_at" field.
	ModifiedAt time.Time `json:"modified_at,omitempty" gqlgen:"modified"`
	// UploadedAt holds the value of the "uploaded_at" field.
	UploadedAt time.Time `json:"uploaded_at,omitempty" gqlgen:"uploaded"`
	// ContentType holds the value of the "content_type" field.
	ContentType string `json:"content_type,omitempty" gqlgen:"mimeType"`
	// StoreKey holds the value of the "store_key" field.
	StoreKey string `json:"store_key,omitempty"`
	// Category holds the value of the "category" field.
	Category string `json:"category,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the FileQuery when eager-loading is set.
	Edges                      FileEdges `json:"edges"`
	check_list_item_files      *int
	equipment_files            *int
	floor_plan_image           *int
	location_files             *int
	survey_source_file         *int
	survey_question_photo_data *int
	survey_question_images     *int
	user_profile_photo         *int
	work_order_files           *int
}

// FileEdges holds the relations/edges for other nodes in the graph.
type FileEdges struct {
	// Location holds the value of the location edge.
	Location *Location
	// Equipment holds the value of the equipment edge.
	Equipment *Equipment
	// User holds the value of the user edge.
	User *User
	// WorkOrder holds the value of the work_order edge.
	WorkOrder *WorkOrder
	// ChecklistItem holds the value of the checklist_item edge.
	ChecklistItem *CheckListItem
	// Survey holds the value of the survey edge.
	Survey *Survey
	// FloorPlan holds the value of the floor_plan edge.
	FloorPlan *FloorPlan
	// PhotoSurveyQuestion holds the value of the photo_survey_question edge.
	PhotoSurveyQuestion *SurveyQuestion
	// SurveyQuestion holds the value of the survey_question edge.
	SurveyQuestion *SurveyQuestion
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [9]bool
}

// LocationOrErr returns the Location value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e FileEdges) LocationOrErr() (*Location, error) {
	if e.loadedTypes[0] {
		if e.Location == nil {
			// The edge location was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: location.Label}
		}
		return e.Location, nil
	}
	return nil, &NotLoadedError{edge: "location"}
}

// EquipmentOrErr returns the Equipment value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e FileEdges) EquipmentOrErr() (*Equipment, error) {
	if e.loadedTypes[1] {
		if e.Equipment == nil {
			// The edge equipment was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: equipment.Label}
		}
		return e.Equipment, nil
	}
	return nil, &NotLoadedError{edge: "equipment"}
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e FileEdges) UserOrErr() (*User, error) {
	if e.loadedTypes[2] {
		if e.User == nil {
			// The edge user was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: user.Label}
		}
		return e.User, nil
	}
	return nil, &NotLoadedError{edge: "user"}
}

// WorkOrderOrErr returns the WorkOrder value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e FileEdges) WorkOrderOrErr() (*WorkOrder, error) {
	if e.loadedTypes[3] {
		if e.WorkOrder == nil {
			// The edge work_order was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: workorder.Label}
		}
		return e.WorkOrder, nil
	}
	return nil, &NotLoadedError{edge: "work_order"}
}

// ChecklistItemOrErr returns the ChecklistItem value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e FileEdges) ChecklistItemOrErr() (*CheckListItem, error) {
	if e.loadedTypes[4] {
		if e.ChecklistItem == nil {
			// The edge checklist_item was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: checklistitem.Label}
		}
		return e.ChecklistItem, nil
	}
	return nil, &NotLoadedError{edge: "checklist_item"}
}

// SurveyOrErr returns the Survey value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e FileEdges) SurveyOrErr() (*Survey, error) {
	if e.loadedTypes[5] {
		if e.Survey == nil {
			// The edge survey was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: survey.Label}
		}
		return e.Survey, nil
	}
	return nil, &NotLoadedError{edge: "survey"}
}

// FloorPlanOrErr returns the FloorPlan value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e FileEdges) FloorPlanOrErr() (*FloorPlan, error) {
	if e.loadedTypes[6] {
		if e.FloorPlan == nil {
			// The edge floor_plan was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: floorplan.Label}
		}
		return e.FloorPlan, nil
	}
	return nil, &NotLoadedError{edge: "floor_plan"}
}

// PhotoSurveyQuestionOrErr returns the PhotoSurveyQuestion value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e FileEdges) PhotoSurveyQuestionOrErr() (*SurveyQuestion, error) {
	if e.loadedTypes[7] {
		if e.PhotoSurveyQuestion == nil {
			// The edge photo_survey_question was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: surveyquestion.Label}
		}
		return e.PhotoSurveyQuestion, nil
	}
	return nil, &NotLoadedError{edge: "photo_survey_question"}
}

// SurveyQuestionOrErr returns the SurveyQuestion value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e FileEdges) SurveyQuestionOrErr() (*SurveyQuestion, error) {
	if e.loadedTypes[8] {
		if e.SurveyQuestion == nil {
			// The edge survey_question was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: surveyquestion.Label}
		}
		return e.SurveyQuestion, nil
	}
	return nil, &NotLoadedError{edge: "survey_question"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*File) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // type
		&sql.NullString{}, // name
		&sql.NullInt64{},  // size
		&sql.NullTime{},   // modified_at
		&sql.NullTime{},   // uploaded_at
		&sql.NullString{}, // content_type
		&sql.NullString{}, // store_key
		&sql.NullString{}, // category
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*File) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // check_list_item_files
		&sql.NullInt64{}, // equipment_files
		&sql.NullInt64{}, // floor_plan_image
		&sql.NullInt64{}, // location_files
		&sql.NullInt64{}, // survey_source_file
		&sql.NullInt64{}, // survey_question_photo_data
		&sql.NullInt64{}, // survey_question_images
		&sql.NullInt64{}, // user_profile_photo
		&sql.NullInt64{}, // work_order_files
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the File fields.
func (f *File) assignValues(values ...interface{}) error {
	if m, n := len(values), len(file.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	f.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		f.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		f.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field type", values[2])
	} else if value.Valid {
		f.Type = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[3])
	} else if value.Valid {
		f.Name = value.String
	}
	if value, ok := values[4].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field size", values[4])
	} else if value.Valid {
		f.Size = int(value.Int64)
	}
	if value, ok := values[5].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field modified_at", values[5])
	} else if value.Valid {
		f.ModifiedAt = value.Time
	}
	if value, ok := values[6].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field uploaded_at", values[6])
	} else if value.Valid {
		f.UploadedAt = value.Time
	}
	if value, ok := values[7].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field content_type", values[7])
	} else if value.Valid {
		f.ContentType = value.String
	}
	if value, ok := values[8].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field store_key", values[8])
	} else if value.Valid {
		f.StoreKey = value.String
	}
	if value, ok := values[9].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field category", values[9])
	} else if value.Valid {
		f.Category = value.String
	}
	values = values[10:]
	if len(values) == len(file.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field check_list_item_files", value)
		} else if value.Valid {
			f.check_list_item_files = new(int)
			*f.check_list_item_files = int(value.Int64)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field equipment_files", value)
		} else if value.Valid {
			f.equipment_files = new(int)
			*f.equipment_files = int(value.Int64)
		}
		if value, ok := values[2].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field floor_plan_image", value)
		} else if value.Valid {
			f.floor_plan_image = new(int)
			*f.floor_plan_image = int(value.Int64)
		}
		if value, ok := values[3].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field location_files", value)
		} else if value.Valid {
			f.location_files = new(int)
			*f.location_files = int(value.Int64)
		}
		if value, ok := values[4].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field survey_source_file", value)
		} else if value.Valid {
			f.survey_source_file = new(int)
			*f.survey_source_file = int(value.Int64)
		}
		if value, ok := values[5].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field survey_question_photo_data", value)
		} else if value.Valid {
			f.survey_question_photo_data = new(int)
			*f.survey_question_photo_data = int(value.Int64)
		}
		if value, ok := values[6].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field survey_question_images", value)
		} else if value.Valid {
			f.survey_question_images = new(int)
			*f.survey_question_images = int(value.Int64)
		}
		if value, ok := values[7].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field user_profile_photo", value)
		} else if value.Valid {
			f.user_profile_photo = new(int)
			*f.user_profile_photo = int(value.Int64)
		}
		if value, ok := values[8].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field work_order_files", value)
		} else if value.Valid {
			f.work_order_files = new(int)
			*f.work_order_files = int(value.Int64)
		}
	}
	return nil
}

// QueryLocation queries the location edge of the File.
func (f *File) QueryLocation() *LocationQuery {
	return (&FileClient{config: f.config}).QueryLocation(f)
}

// QueryEquipment queries the equipment edge of the File.
func (f *File) QueryEquipment() *EquipmentQuery {
	return (&FileClient{config: f.config}).QueryEquipment(f)
}

// QueryUser queries the user edge of the File.
func (f *File) QueryUser() *UserQuery {
	return (&FileClient{config: f.config}).QueryUser(f)
}

// QueryWorkOrder queries the work_order edge of the File.
func (f *File) QueryWorkOrder() *WorkOrderQuery {
	return (&FileClient{config: f.config}).QueryWorkOrder(f)
}

// QueryChecklistItem queries the checklist_item edge of the File.
func (f *File) QueryChecklistItem() *CheckListItemQuery {
	return (&FileClient{config: f.config}).QueryChecklistItem(f)
}

// QuerySurvey queries the survey edge of the File.
func (f *File) QuerySurvey() *SurveyQuery {
	return (&FileClient{config: f.config}).QuerySurvey(f)
}

// QueryFloorPlan queries the floor_plan edge of the File.
func (f *File) QueryFloorPlan() *FloorPlanQuery {
	return (&FileClient{config: f.config}).QueryFloorPlan(f)
}

// QueryPhotoSurveyQuestion queries the photo_survey_question edge of the File.
func (f *File) QueryPhotoSurveyQuestion() *SurveyQuestionQuery {
	return (&FileClient{config: f.config}).QueryPhotoSurveyQuestion(f)
}

// QuerySurveyQuestion queries the survey_question edge of the File.
func (f *File) QuerySurveyQuestion() *SurveyQuestionQuery {
	return (&FileClient{config: f.config}).QuerySurveyQuestion(f)
}

// Update returns a builder for updating this File.
// Note that, you need to call File.Unwrap() before calling this method, if this File
// was returned from a transaction, and the transaction was committed or rolled back.
func (f *File) Update() *FileUpdateOne {
	return (&FileClient{config: f.config}).UpdateOne(f)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (f *File) Unwrap() *File {
	tx, ok := f.config.driver.(*txDriver)
	if !ok {
		panic("ent: File is not a transactional entity")
	}
	f.config.driver = tx.drv
	return f
}

// String implements the fmt.Stringer.
func (f *File) String() string {
	var builder strings.Builder
	builder.WriteString("File(")
	builder.WriteString(fmt.Sprintf("id=%v", f.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(f.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(f.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", type=")
	builder.WriteString(f.Type)
	builder.WriteString(", name=")
	builder.WriteString(f.Name)
	builder.WriteString(", size=")
	builder.WriteString(fmt.Sprintf("%v", f.Size))
	builder.WriteString(", modified_at=")
	builder.WriteString(f.ModifiedAt.Format(time.ANSIC))
	builder.WriteString(", uploaded_at=")
	builder.WriteString(f.UploadedAt.Format(time.ANSIC))
	builder.WriteString(", content_type=")
	builder.WriteString(f.ContentType)
	builder.WriteString(", store_key=")
	builder.WriteString(f.StoreKey)
	builder.WriteString(", category=")
	builder.WriteString(f.Category)
	builder.WriteByte(')')
	return builder.String()
}

// Files is a parsable slice of File.
type Files []*File

func (f Files) config(cfg config) {
	for _i := range f {
		f[_i].config = cfg
	}
}
