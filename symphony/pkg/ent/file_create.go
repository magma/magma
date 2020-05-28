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

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/ent/checklistitem"
	"github.com/facebookincubator/symphony/pkg/ent/equipment"
	"github.com/facebookincubator/symphony/pkg/ent/file"
	"github.com/facebookincubator/symphony/pkg/ent/floorplan"
	"github.com/facebookincubator/symphony/pkg/ent/location"
	"github.com/facebookincubator/symphony/pkg/ent/survey"
	"github.com/facebookincubator/symphony/pkg/ent/surveyquestion"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/ent/workorder"
)

// FileCreate is the builder for creating a File entity.
type FileCreate struct {
	config
	mutation *FileMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (fc *FileCreate) SetCreateTime(t time.Time) *FileCreate {
	fc.mutation.SetCreateTime(t)
	return fc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (fc *FileCreate) SetNillableCreateTime(t *time.Time) *FileCreate {
	if t != nil {
		fc.SetCreateTime(*t)
	}
	return fc
}

// SetUpdateTime sets the update_time field.
func (fc *FileCreate) SetUpdateTime(t time.Time) *FileCreate {
	fc.mutation.SetUpdateTime(t)
	return fc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (fc *FileCreate) SetNillableUpdateTime(t *time.Time) *FileCreate {
	if t != nil {
		fc.SetUpdateTime(*t)
	}
	return fc
}

// SetType sets the type field.
func (fc *FileCreate) SetType(s string) *FileCreate {
	fc.mutation.SetType(s)
	return fc
}

// SetName sets the name field.
func (fc *FileCreate) SetName(s string) *FileCreate {
	fc.mutation.SetName(s)
	return fc
}

// SetSize sets the size field.
func (fc *FileCreate) SetSize(i int) *FileCreate {
	fc.mutation.SetSize(i)
	return fc
}

// SetNillableSize sets the size field if the given value is not nil.
func (fc *FileCreate) SetNillableSize(i *int) *FileCreate {
	if i != nil {
		fc.SetSize(*i)
	}
	return fc
}

// SetModifiedAt sets the modified_at field.
func (fc *FileCreate) SetModifiedAt(t time.Time) *FileCreate {
	fc.mutation.SetModifiedAt(t)
	return fc
}

// SetNillableModifiedAt sets the modified_at field if the given value is not nil.
func (fc *FileCreate) SetNillableModifiedAt(t *time.Time) *FileCreate {
	if t != nil {
		fc.SetModifiedAt(*t)
	}
	return fc
}

// SetUploadedAt sets the uploaded_at field.
func (fc *FileCreate) SetUploadedAt(t time.Time) *FileCreate {
	fc.mutation.SetUploadedAt(t)
	return fc
}

// SetNillableUploadedAt sets the uploaded_at field if the given value is not nil.
func (fc *FileCreate) SetNillableUploadedAt(t *time.Time) *FileCreate {
	if t != nil {
		fc.SetUploadedAt(*t)
	}
	return fc
}

// SetContentType sets the content_type field.
func (fc *FileCreate) SetContentType(s string) *FileCreate {
	fc.mutation.SetContentType(s)
	return fc
}

// SetStoreKey sets the store_key field.
func (fc *FileCreate) SetStoreKey(s string) *FileCreate {
	fc.mutation.SetStoreKey(s)
	return fc
}

// SetCategory sets the category field.
func (fc *FileCreate) SetCategory(s string) *FileCreate {
	fc.mutation.SetCategory(s)
	return fc
}

// SetNillableCategory sets the category field if the given value is not nil.
func (fc *FileCreate) SetNillableCategory(s *string) *FileCreate {
	if s != nil {
		fc.SetCategory(*s)
	}
	return fc
}

// SetLocationID sets the location edge to Location by id.
func (fc *FileCreate) SetLocationID(id int) *FileCreate {
	fc.mutation.SetLocationID(id)
	return fc
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (fc *FileCreate) SetNillableLocationID(id *int) *FileCreate {
	if id != nil {
		fc = fc.SetLocationID(*id)
	}
	return fc
}

// SetLocation sets the location edge to Location.
func (fc *FileCreate) SetLocation(l *Location) *FileCreate {
	return fc.SetLocationID(l.ID)
}

// SetEquipmentID sets the equipment edge to Equipment by id.
func (fc *FileCreate) SetEquipmentID(id int) *FileCreate {
	fc.mutation.SetEquipmentID(id)
	return fc
}

// SetNillableEquipmentID sets the equipment edge to Equipment by id if the given value is not nil.
func (fc *FileCreate) SetNillableEquipmentID(id *int) *FileCreate {
	if id != nil {
		fc = fc.SetEquipmentID(*id)
	}
	return fc
}

// SetEquipment sets the equipment edge to Equipment.
func (fc *FileCreate) SetEquipment(e *Equipment) *FileCreate {
	return fc.SetEquipmentID(e.ID)
}

// SetUserID sets the user edge to User by id.
func (fc *FileCreate) SetUserID(id int) *FileCreate {
	fc.mutation.SetUserID(id)
	return fc
}

// SetNillableUserID sets the user edge to User by id if the given value is not nil.
func (fc *FileCreate) SetNillableUserID(id *int) *FileCreate {
	if id != nil {
		fc = fc.SetUserID(*id)
	}
	return fc
}

// SetUser sets the user edge to User.
func (fc *FileCreate) SetUser(u *User) *FileCreate {
	return fc.SetUserID(u.ID)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (fc *FileCreate) SetWorkOrderID(id int) *FileCreate {
	fc.mutation.SetWorkOrderID(id)
	return fc
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (fc *FileCreate) SetNillableWorkOrderID(id *int) *FileCreate {
	if id != nil {
		fc = fc.SetWorkOrderID(*id)
	}
	return fc
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (fc *FileCreate) SetWorkOrder(w *WorkOrder) *FileCreate {
	return fc.SetWorkOrderID(w.ID)
}

// SetChecklistItemID sets the checklist_item edge to CheckListItem by id.
func (fc *FileCreate) SetChecklistItemID(id int) *FileCreate {
	fc.mutation.SetChecklistItemID(id)
	return fc
}

// SetNillableChecklistItemID sets the checklist_item edge to CheckListItem by id if the given value is not nil.
func (fc *FileCreate) SetNillableChecklistItemID(id *int) *FileCreate {
	if id != nil {
		fc = fc.SetChecklistItemID(*id)
	}
	return fc
}

// SetChecklistItem sets the checklist_item edge to CheckListItem.
func (fc *FileCreate) SetChecklistItem(c *CheckListItem) *FileCreate {
	return fc.SetChecklistItemID(c.ID)
}

// SetSurveyID sets the survey edge to Survey by id.
func (fc *FileCreate) SetSurveyID(id int) *FileCreate {
	fc.mutation.SetSurveyID(id)
	return fc
}

// SetNillableSurveyID sets the survey edge to Survey by id if the given value is not nil.
func (fc *FileCreate) SetNillableSurveyID(id *int) *FileCreate {
	if id != nil {
		fc = fc.SetSurveyID(*id)
	}
	return fc
}

// SetSurvey sets the survey edge to Survey.
func (fc *FileCreate) SetSurvey(s *Survey) *FileCreate {
	return fc.SetSurveyID(s.ID)
}

// SetFloorPlanID sets the floor_plan edge to FloorPlan by id.
func (fc *FileCreate) SetFloorPlanID(id int) *FileCreate {
	fc.mutation.SetFloorPlanID(id)
	return fc
}

// SetNillableFloorPlanID sets the floor_plan edge to FloorPlan by id if the given value is not nil.
func (fc *FileCreate) SetNillableFloorPlanID(id *int) *FileCreate {
	if id != nil {
		fc = fc.SetFloorPlanID(*id)
	}
	return fc
}

// SetFloorPlan sets the floor_plan edge to FloorPlan.
func (fc *FileCreate) SetFloorPlan(f *FloorPlan) *FileCreate {
	return fc.SetFloorPlanID(f.ID)
}

// SetPhotoSurveyQuestionID sets the photo_survey_question edge to SurveyQuestion by id.
func (fc *FileCreate) SetPhotoSurveyQuestionID(id int) *FileCreate {
	fc.mutation.SetPhotoSurveyQuestionID(id)
	return fc
}

// SetNillablePhotoSurveyQuestionID sets the photo_survey_question edge to SurveyQuestion by id if the given value is not nil.
func (fc *FileCreate) SetNillablePhotoSurveyQuestionID(id *int) *FileCreate {
	if id != nil {
		fc = fc.SetPhotoSurveyQuestionID(*id)
	}
	return fc
}

// SetPhotoSurveyQuestion sets the photo_survey_question edge to SurveyQuestion.
func (fc *FileCreate) SetPhotoSurveyQuestion(s *SurveyQuestion) *FileCreate {
	return fc.SetPhotoSurveyQuestionID(s.ID)
}

// SetSurveyQuestionID sets the survey_question edge to SurveyQuestion by id.
func (fc *FileCreate) SetSurveyQuestionID(id int) *FileCreate {
	fc.mutation.SetSurveyQuestionID(id)
	return fc
}

// SetNillableSurveyQuestionID sets the survey_question edge to SurveyQuestion by id if the given value is not nil.
func (fc *FileCreate) SetNillableSurveyQuestionID(id *int) *FileCreate {
	if id != nil {
		fc = fc.SetSurveyQuestionID(*id)
	}
	return fc
}

// SetSurveyQuestion sets the survey_question edge to SurveyQuestion.
func (fc *FileCreate) SetSurveyQuestion(s *SurveyQuestion) *FileCreate {
	return fc.SetSurveyQuestionID(s.ID)
}

// Save creates the File in the database.
func (fc *FileCreate) Save(ctx context.Context) (*File, error) {
	if _, ok := fc.mutation.CreateTime(); !ok {
		v := file.DefaultCreateTime()
		fc.mutation.SetCreateTime(v)
	}
	if _, ok := fc.mutation.UpdateTime(); !ok {
		v := file.DefaultUpdateTime()
		fc.mutation.SetUpdateTime(v)
	}
	if _, ok := fc.mutation.GetType(); !ok {
		return nil, errors.New("ent: missing required field \"type\"")
	}
	if _, ok := fc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if v, ok := fc.mutation.Size(); ok {
		if err := file.SizeValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"size\": %v", err)
		}
	}
	if _, ok := fc.mutation.ContentType(); !ok {
		return nil, errors.New("ent: missing required field \"content_type\"")
	}
	if _, ok := fc.mutation.StoreKey(); !ok {
		return nil, errors.New("ent: missing required field \"store_key\"")
	}
	var (
		err  error
		node *File
	)
	if len(fc.hooks) == 0 {
		node, err = fc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FileMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fc.mutation = mutation
			node, err = fc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(fc.hooks) - 1; i >= 0; i-- {
			mut = fc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (fc *FileCreate) SaveX(ctx context.Context) *File {
	v, err := fc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (fc *FileCreate) sqlSave(ctx context.Context) (*File, error) {
	var (
		f     = &File{config: fc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: file.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: file.FieldID,
			},
		}
	)
	if value, ok := fc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldCreateTime,
		})
		f.CreateTime = value
	}
	if value, ok := fc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldUpdateTime,
		})
		f.UpdateTime = value
	}
	if value, ok := fc.mutation.GetType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldType,
		})
		f.Type = value
	}
	if value, ok := fc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldName,
		})
		f.Name = value
	}
	if value, ok := fc.mutation.Size(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: file.FieldSize,
		})
		f.Size = value
	}
	if value, ok := fc.mutation.ModifiedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldModifiedAt,
		})
		f.ModifiedAt = value
	}
	if value, ok := fc.mutation.UploadedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldUploadedAt,
		})
		f.UploadedAt = value
	}
	if value, ok := fc.mutation.ContentType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldContentType,
		})
		f.ContentType = value
	}
	if value, ok := fc.mutation.StoreKey(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldStoreKey,
		})
		f.StoreKey = value
	}
	if value, ok := fc.mutation.Category(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldCategory,
		})
		f.Category = value
	}
	if nodes := fc.mutation.LocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   file.LocationTable,
			Columns: []string{file.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fc.mutation.EquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   file.EquipmentTable,
			Columns: []string{file.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fc.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   file.UserTable,
			Columns: []string{file.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fc.mutation.WorkOrderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   file.WorkOrderTable,
			Columns: []string{file.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fc.mutation.ChecklistItemIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   file.ChecklistItemTable,
			Columns: []string{file.ChecklistItemColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistitem.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fc.mutation.SurveyIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   file.SurveyTable,
			Columns: []string{file.SurveyColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fc.mutation.FloorPlanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   file.FloorPlanTable,
			Columns: []string{file.FloorPlanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: floorplan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fc.mutation.PhotoSurveyQuestionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   file.PhotoSurveyQuestionTable,
			Columns: []string{file.PhotoSurveyQuestionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveyquestion.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fc.mutation.SurveyQuestionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   file.SurveyQuestionTable,
			Columns: []string{file.SurveyQuestionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveyquestion.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, fc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	f.ID = int(id)
	return f, nil
}
