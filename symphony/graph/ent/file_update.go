// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/floorplan"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// FileUpdate is the builder for updating File entities.
type FileUpdate struct {
	config
	hooks      []Hook
	mutation   *FileMutation
	predicates []predicate.File
}

// Where adds a new predicate for the builder.
func (fu *FileUpdate) Where(ps ...predicate.File) *FileUpdate {
	fu.predicates = append(fu.predicates, ps...)
	return fu
}

// SetType sets the type field.
func (fu *FileUpdate) SetType(s string) *FileUpdate {
	fu.mutation.SetType(s)
	return fu
}

// SetName sets the name field.
func (fu *FileUpdate) SetName(s string) *FileUpdate {
	fu.mutation.SetName(s)
	return fu
}

// SetSize sets the size field.
func (fu *FileUpdate) SetSize(i int) *FileUpdate {
	fu.mutation.ResetSize()
	fu.mutation.SetSize(i)
	return fu
}

// SetNillableSize sets the size field if the given value is not nil.
func (fu *FileUpdate) SetNillableSize(i *int) *FileUpdate {
	if i != nil {
		fu.SetSize(*i)
	}
	return fu
}

// AddSize adds i to size.
func (fu *FileUpdate) AddSize(i int) *FileUpdate {
	fu.mutation.AddSize(i)
	return fu
}

// ClearSize clears the value of size.
func (fu *FileUpdate) ClearSize() *FileUpdate {
	fu.mutation.ClearSize()
	return fu
}

// SetModifiedAt sets the modified_at field.
func (fu *FileUpdate) SetModifiedAt(t time.Time) *FileUpdate {
	fu.mutation.SetModifiedAt(t)
	return fu
}

// SetNillableModifiedAt sets the modified_at field if the given value is not nil.
func (fu *FileUpdate) SetNillableModifiedAt(t *time.Time) *FileUpdate {
	if t != nil {
		fu.SetModifiedAt(*t)
	}
	return fu
}

// ClearModifiedAt clears the value of modified_at.
func (fu *FileUpdate) ClearModifiedAt() *FileUpdate {
	fu.mutation.ClearModifiedAt()
	return fu
}

// SetUploadedAt sets the uploaded_at field.
func (fu *FileUpdate) SetUploadedAt(t time.Time) *FileUpdate {
	fu.mutation.SetUploadedAt(t)
	return fu
}

// SetNillableUploadedAt sets the uploaded_at field if the given value is not nil.
func (fu *FileUpdate) SetNillableUploadedAt(t *time.Time) *FileUpdate {
	if t != nil {
		fu.SetUploadedAt(*t)
	}
	return fu
}

// ClearUploadedAt clears the value of uploaded_at.
func (fu *FileUpdate) ClearUploadedAt() *FileUpdate {
	fu.mutation.ClearUploadedAt()
	return fu
}

// SetContentType sets the content_type field.
func (fu *FileUpdate) SetContentType(s string) *FileUpdate {
	fu.mutation.SetContentType(s)
	return fu
}

// SetStoreKey sets the store_key field.
func (fu *FileUpdate) SetStoreKey(s string) *FileUpdate {
	fu.mutation.SetStoreKey(s)
	return fu
}

// SetCategory sets the category field.
func (fu *FileUpdate) SetCategory(s string) *FileUpdate {
	fu.mutation.SetCategory(s)
	return fu
}

// SetNillableCategory sets the category field if the given value is not nil.
func (fu *FileUpdate) SetNillableCategory(s *string) *FileUpdate {
	if s != nil {
		fu.SetCategory(*s)
	}
	return fu
}

// ClearCategory clears the value of category.
func (fu *FileUpdate) ClearCategory() *FileUpdate {
	fu.mutation.ClearCategory()
	return fu
}

// SetLocationID sets the location edge to Location by id.
func (fu *FileUpdate) SetLocationID(id int) *FileUpdate {
	fu.mutation.SetLocationID(id)
	return fu
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (fu *FileUpdate) SetNillableLocationID(id *int) *FileUpdate {
	if id != nil {
		fu = fu.SetLocationID(*id)
	}
	return fu
}

// SetLocation sets the location edge to Location.
func (fu *FileUpdate) SetLocation(l *Location) *FileUpdate {
	return fu.SetLocationID(l.ID)
}

// SetEquipmentID sets the equipment edge to Equipment by id.
func (fu *FileUpdate) SetEquipmentID(id int) *FileUpdate {
	fu.mutation.SetEquipmentID(id)
	return fu
}

// SetNillableEquipmentID sets the equipment edge to Equipment by id if the given value is not nil.
func (fu *FileUpdate) SetNillableEquipmentID(id *int) *FileUpdate {
	if id != nil {
		fu = fu.SetEquipmentID(*id)
	}
	return fu
}

// SetEquipment sets the equipment edge to Equipment.
func (fu *FileUpdate) SetEquipment(e *Equipment) *FileUpdate {
	return fu.SetEquipmentID(e.ID)
}

// SetUserID sets the user edge to User by id.
func (fu *FileUpdate) SetUserID(id int) *FileUpdate {
	fu.mutation.SetUserID(id)
	return fu
}

// SetNillableUserID sets the user edge to User by id if the given value is not nil.
func (fu *FileUpdate) SetNillableUserID(id *int) *FileUpdate {
	if id != nil {
		fu = fu.SetUserID(*id)
	}
	return fu
}

// SetUser sets the user edge to User.
func (fu *FileUpdate) SetUser(u *User) *FileUpdate {
	return fu.SetUserID(u.ID)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (fu *FileUpdate) SetWorkOrderID(id int) *FileUpdate {
	fu.mutation.SetWorkOrderID(id)
	return fu
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (fu *FileUpdate) SetNillableWorkOrderID(id *int) *FileUpdate {
	if id != nil {
		fu = fu.SetWorkOrderID(*id)
	}
	return fu
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (fu *FileUpdate) SetWorkOrder(w *WorkOrder) *FileUpdate {
	return fu.SetWorkOrderID(w.ID)
}

// SetChecklistItemID sets the checklist_item edge to CheckListItem by id.
func (fu *FileUpdate) SetChecklistItemID(id int) *FileUpdate {
	fu.mutation.SetChecklistItemID(id)
	return fu
}

// SetNillableChecklistItemID sets the checklist_item edge to CheckListItem by id if the given value is not nil.
func (fu *FileUpdate) SetNillableChecklistItemID(id *int) *FileUpdate {
	if id != nil {
		fu = fu.SetChecklistItemID(*id)
	}
	return fu
}

// SetChecklistItem sets the checklist_item edge to CheckListItem.
func (fu *FileUpdate) SetChecklistItem(c *CheckListItem) *FileUpdate {
	return fu.SetChecklistItemID(c.ID)
}

// SetSurveyID sets the survey edge to Survey by id.
func (fu *FileUpdate) SetSurveyID(id int) *FileUpdate {
	fu.mutation.SetSurveyID(id)
	return fu
}

// SetNillableSurveyID sets the survey edge to Survey by id if the given value is not nil.
func (fu *FileUpdate) SetNillableSurveyID(id *int) *FileUpdate {
	if id != nil {
		fu = fu.SetSurveyID(*id)
	}
	return fu
}

// SetSurvey sets the survey edge to Survey.
func (fu *FileUpdate) SetSurvey(s *Survey) *FileUpdate {
	return fu.SetSurveyID(s.ID)
}

// SetFloorPlanID sets the floor_plan edge to FloorPlan by id.
func (fu *FileUpdate) SetFloorPlanID(id int) *FileUpdate {
	fu.mutation.SetFloorPlanID(id)
	return fu
}

// SetNillableFloorPlanID sets the floor_plan edge to FloorPlan by id if the given value is not nil.
func (fu *FileUpdate) SetNillableFloorPlanID(id *int) *FileUpdate {
	if id != nil {
		fu = fu.SetFloorPlanID(*id)
	}
	return fu
}

// SetFloorPlan sets the floor_plan edge to FloorPlan.
func (fu *FileUpdate) SetFloorPlan(f *FloorPlan) *FileUpdate {
	return fu.SetFloorPlanID(f.ID)
}

// SetPhotoSurveyQuestionID sets the photo_survey_question edge to SurveyQuestion by id.
func (fu *FileUpdate) SetPhotoSurveyQuestionID(id int) *FileUpdate {
	fu.mutation.SetPhotoSurveyQuestionID(id)
	return fu
}

// SetNillablePhotoSurveyQuestionID sets the photo_survey_question edge to SurveyQuestion by id if the given value is not nil.
func (fu *FileUpdate) SetNillablePhotoSurveyQuestionID(id *int) *FileUpdate {
	if id != nil {
		fu = fu.SetPhotoSurveyQuestionID(*id)
	}
	return fu
}

// SetPhotoSurveyQuestion sets the photo_survey_question edge to SurveyQuestion.
func (fu *FileUpdate) SetPhotoSurveyQuestion(s *SurveyQuestion) *FileUpdate {
	return fu.SetPhotoSurveyQuestionID(s.ID)
}

// SetSurveyQuestionID sets the survey_question edge to SurveyQuestion by id.
func (fu *FileUpdate) SetSurveyQuestionID(id int) *FileUpdate {
	fu.mutation.SetSurveyQuestionID(id)
	return fu
}

// SetNillableSurveyQuestionID sets the survey_question edge to SurveyQuestion by id if the given value is not nil.
func (fu *FileUpdate) SetNillableSurveyQuestionID(id *int) *FileUpdate {
	if id != nil {
		fu = fu.SetSurveyQuestionID(*id)
	}
	return fu
}

// SetSurveyQuestion sets the survey_question edge to SurveyQuestion.
func (fu *FileUpdate) SetSurveyQuestion(s *SurveyQuestion) *FileUpdate {
	return fu.SetSurveyQuestionID(s.ID)
}

// ClearLocation clears the location edge to Location.
func (fu *FileUpdate) ClearLocation() *FileUpdate {
	fu.mutation.ClearLocation()
	return fu
}

// ClearEquipment clears the equipment edge to Equipment.
func (fu *FileUpdate) ClearEquipment() *FileUpdate {
	fu.mutation.ClearEquipment()
	return fu
}

// ClearUser clears the user edge to User.
func (fu *FileUpdate) ClearUser() *FileUpdate {
	fu.mutation.ClearUser()
	return fu
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (fu *FileUpdate) ClearWorkOrder() *FileUpdate {
	fu.mutation.ClearWorkOrder()
	return fu
}

// ClearChecklistItem clears the checklist_item edge to CheckListItem.
func (fu *FileUpdate) ClearChecklistItem() *FileUpdate {
	fu.mutation.ClearChecklistItem()
	return fu
}

// ClearSurvey clears the survey edge to Survey.
func (fu *FileUpdate) ClearSurvey() *FileUpdate {
	fu.mutation.ClearSurvey()
	return fu
}

// ClearFloorPlan clears the floor_plan edge to FloorPlan.
func (fu *FileUpdate) ClearFloorPlan() *FileUpdate {
	fu.mutation.ClearFloorPlan()
	return fu
}

// ClearPhotoSurveyQuestion clears the photo_survey_question edge to SurveyQuestion.
func (fu *FileUpdate) ClearPhotoSurveyQuestion() *FileUpdate {
	fu.mutation.ClearPhotoSurveyQuestion()
	return fu
}

// ClearSurveyQuestion clears the survey_question edge to SurveyQuestion.
func (fu *FileUpdate) ClearSurveyQuestion() *FileUpdate {
	fu.mutation.ClearSurveyQuestion()
	return fu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (fu *FileUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := fu.mutation.UpdateTime(); !ok {
		v := file.UpdateDefaultUpdateTime()
		fu.mutation.SetUpdateTime(v)
	}
	if v, ok := fu.mutation.Size(); ok {
		if err := file.SizeValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"size\": %v", err)
		}
	}

	var (
		err      error
		affected int
	)
	if len(fu.hooks) == 0 {
		affected, err = fu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FileMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fu.mutation = mutation
			affected, err = fu.sqlSave(ctx)
			return affected, err
		})
		for i := len(fu.hooks) - 1; i >= 0; i-- {
			mut = fu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (fu *FileUpdate) SaveX(ctx context.Context) int {
	affected, err := fu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fu *FileUpdate) Exec(ctx context.Context) error {
	_, err := fu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fu *FileUpdate) ExecX(ctx context.Context) {
	if err := fu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fu *FileUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   file.Table,
			Columns: file.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: file.FieldID,
			},
		},
	}
	if ps := fu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldUpdateTime,
		})
	}
	if value, ok := fu.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldType,
		})
	}
	if value, ok := fu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldName,
		})
	}
	if value, ok := fu.mutation.Size(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: file.FieldSize,
		})
	}
	if value, ok := fu.mutation.AddedSize(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: file.FieldSize,
		})
	}
	if fu.mutation.SizeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: file.FieldSize,
		})
	}
	if value, ok := fu.mutation.ModifiedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldModifiedAt,
		})
	}
	if fu.mutation.ModifiedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: file.FieldModifiedAt,
		})
	}
	if value, ok := fu.mutation.UploadedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldUploadedAt,
		})
	}
	if fu.mutation.UploadedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: file.FieldUploadedAt,
		})
	}
	if value, ok := fu.mutation.ContentType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldContentType,
		})
	}
	if value, ok := fu.mutation.StoreKey(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldStoreKey,
		})
	}
	if value, ok := fu.mutation.Category(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldCategory,
		})
	}
	if fu.mutation.CategoryCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: file.FieldCategory,
		})
	}
	if fu.mutation.LocationCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.LocationIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fu.mutation.EquipmentCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.EquipmentIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fu.mutation.UserCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.UserIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fu.mutation.WorkOrderCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.WorkOrderIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fu.mutation.ChecklistItemCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.ChecklistItemIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fu.mutation.SurveyCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.SurveyIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fu.mutation.FloorPlanCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.FloorPlanIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fu.mutation.PhotoSurveyQuestionCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.PhotoSurveyQuestionIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fu.mutation.SurveyQuestionCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.SurveyQuestionIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{file.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// FileUpdateOne is the builder for updating a single File entity.
type FileUpdateOne struct {
	config
	hooks    []Hook
	mutation *FileMutation
}

// SetType sets the type field.
func (fuo *FileUpdateOne) SetType(s string) *FileUpdateOne {
	fuo.mutation.SetType(s)
	return fuo
}

// SetName sets the name field.
func (fuo *FileUpdateOne) SetName(s string) *FileUpdateOne {
	fuo.mutation.SetName(s)
	return fuo
}

// SetSize sets the size field.
func (fuo *FileUpdateOne) SetSize(i int) *FileUpdateOne {
	fuo.mutation.ResetSize()
	fuo.mutation.SetSize(i)
	return fuo
}

// SetNillableSize sets the size field if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableSize(i *int) *FileUpdateOne {
	if i != nil {
		fuo.SetSize(*i)
	}
	return fuo
}

// AddSize adds i to size.
func (fuo *FileUpdateOne) AddSize(i int) *FileUpdateOne {
	fuo.mutation.AddSize(i)
	return fuo
}

// ClearSize clears the value of size.
func (fuo *FileUpdateOne) ClearSize() *FileUpdateOne {
	fuo.mutation.ClearSize()
	return fuo
}

// SetModifiedAt sets the modified_at field.
func (fuo *FileUpdateOne) SetModifiedAt(t time.Time) *FileUpdateOne {
	fuo.mutation.SetModifiedAt(t)
	return fuo
}

// SetNillableModifiedAt sets the modified_at field if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableModifiedAt(t *time.Time) *FileUpdateOne {
	if t != nil {
		fuo.SetModifiedAt(*t)
	}
	return fuo
}

// ClearModifiedAt clears the value of modified_at.
func (fuo *FileUpdateOne) ClearModifiedAt() *FileUpdateOne {
	fuo.mutation.ClearModifiedAt()
	return fuo
}

// SetUploadedAt sets the uploaded_at field.
func (fuo *FileUpdateOne) SetUploadedAt(t time.Time) *FileUpdateOne {
	fuo.mutation.SetUploadedAt(t)
	return fuo
}

// SetNillableUploadedAt sets the uploaded_at field if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableUploadedAt(t *time.Time) *FileUpdateOne {
	if t != nil {
		fuo.SetUploadedAt(*t)
	}
	return fuo
}

// ClearUploadedAt clears the value of uploaded_at.
func (fuo *FileUpdateOne) ClearUploadedAt() *FileUpdateOne {
	fuo.mutation.ClearUploadedAt()
	return fuo
}

// SetContentType sets the content_type field.
func (fuo *FileUpdateOne) SetContentType(s string) *FileUpdateOne {
	fuo.mutation.SetContentType(s)
	return fuo
}

// SetStoreKey sets the store_key field.
func (fuo *FileUpdateOne) SetStoreKey(s string) *FileUpdateOne {
	fuo.mutation.SetStoreKey(s)
	return fuo
}

// SetCategory sets the category field.
func (fuo *FileUpdateOne) SetCategory(s string) *FileUpdateOne {
	fuo.mutation.SetCategory(s)
	return fuo
}

// SetNillableCategory sets the category field if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableCategory(s *string) *FileUpdateOne {
	if s != nil {
		fuo.SetCategory(*s)
	}
	return fuo
}

// ClearCategory clears the value of category.
func (fuo *FileUpdateOne) ClearCategory() *FileUpdateOne {
	fuo.mutation.ClearCategory()
	return fuo
}

// SetLocationID sets the location edge to Location by id.
func (fuo *FileUpdateOne) SetLocationID(id int) *FileUpdateOne {
	fuo.mutation.SetLocationID(id)
	return fuo
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableLocationID(id *int) *FileUpdateOne {
	if id != nil {
		fuo = fuo.SetLocationID(*id)
	}
	return fuo
}

// SetLocation sets the location edge to Location.
func (fuo *FileUpdateOne) SetLocation(l *Location) *FileUpdateOne {
	return fuo.SetLocationID(l.ID)
}

// SetEquipmentID sets the equipment edge to Equipment by id.
func (fuo *FileUpdateOne) SetEquipmentID(id int) *FileUpdateOne {
	fuo.mutation.SetEquipmentID(id)
	return fuo
}

// SetNillableEquipmentID sets the equipment edge to Equipment by id if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableEquipmentID(id *int) *FileUpdateOne {
	if id != nil {
		fuo = fuo.SetEquipmentID(*id)
	}
	return fuo
}

// SetEquipment sets the equipment edge to Equipment.
func (fuo *FileUpdateOne) SetEquipment(e *Equipment) *FileUpdateOne {
	return fuo.SetEquipmentID(e.ID)
}

// SetUserID sets the user edge to User by id.
func (fuo *FileUpdateOne) SetUserID(id int) *FileUpdateOne {
	fuo.mutation.SetUserID(id)
	return fuo
}

// SetNillableUserID sets the user edge to User by id if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableUserID(id *int) *FileUpdateOne {
	if id != nil {
		fuo = fuo.SetUserID(*id)
	}
	return fuo
}

// SetUser sets the user edge to User.
func (fuo *FileUpdateOne) SetUser(u *User) *FileUpdateOne {
	return fuo.SetUserID(u.ID)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (fuo *FileUpdateOne) SetWorkOrderID(id int) *FileUpdateOne {
	fuo.mutation.SetWorkOrderID(id)
	return fuo
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableWorkOrderID(id *int) *FileUpdateOne {
	if id != nil {
		fuo = fuo.SetWorkOrderID(*id)
	}
	return fuo
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (fuo *FileUpdateOne) SetWorkOrder(w *WorkOrder) *FileUpdateOne {
	return fuo.SetWorkOrderID(w.ID)
}

// SetChecklistItemID sets the checklist_item edge to CheckListItem by id.
func (fuo *FileUpdateOne) SetChecklistItemID(id int) *FileUpdateOne {
	fuo.mutation.SetChecklistItemID(id)
	return fuo
}

// SetNillableChecklistItemID sets the checklist_item edge to CheckListItem by id if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableChecklistItemID(id *int) *FileUpdateOne {
	if id != nil {
		fuo = fuo.SetChecklistItemID(*id)
	}
	return fuo
}

// SetChecklistItem sets the checklist_item edge to CheckListItem.
func (fuo *FileUpdateOne) SetChecklistItem(c *CheckListItem) *FileUpdateOne {
	return fuo.SetChecklistItemID(c.ID)
}

// SetSurveyID sets the survey edge to Survey by id.
func (fuo *FileUpdateOne) SetSurveyID(id int) *FileUpdateOne {
	fuo.mutation.SetSurveyID(id)
	return fuo
}

// SetNillableSurveyID sets the survey edge to Survey by id if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableSurveyID(id *int) *FileUpdateOne {
	if id != nil {
		fuo = fuo.SetSurveyID(*id)
	}
	return fuo
}

// SetSurvey sets the survey edge to Survey.
func (fuo *FileUpdateOne) SetSurvey(s *Survey) *FileUpdateOne {
	return fuo.SetSurveyID(s.ID)
}

// SetFloorPlanID sets the floor_plan edge to FloorPlan by id.
func (fuo *FileUpdateOne) SetFloorPlanID(id int) *FileUpdateOne {
	fuo.mutation.SetFloorPlanID(id)
	return fuo
}

// SetNillableFloorPlanID sets the floor_plan edge to FloorPlan by id if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableFloorPlanID(id *int) *FileUpdateOne {
	if id != nil {
		fuo = fuo.SetFloorPlanID(*id)
	}
	return fuo
}

// SetFloorPlan sets the floor_plan edge to FloorPlan.
func (fuo *FileUpdateOne) SetFloorPlan(f *FloorPlan) *FileUpdateOne {
	return fuo.SetFloorPlanID(f.ID)
}

// SetPhotoSurveyQuestionID sets the photo_survey_question edge to SurveyQuestion by id.
func (fuo *FileUpdateOne) SetPhotoSurveyQuestionID(id int) *FileUpdateOne {
	fuo.mutation.SetPhotoSurveyQuestionID(id)
	return fuo
}

// SetNillablePhotoSurveyQuestionID sets the photo_survey_question edge to SurveyQuestion by id if the given value is not nil.
func (fuo *FileUpdateOne) SetNillablePhotoSurveyQuestionID(id *int) *FileUpdateOne {
	if id != nil {
		fuo = fuo.SetPhotoSurveyQuestionID(*id)
	}
	return fuo
}

// SetPhotoSurveyQuestion sets the photo_survey_question edge to SurveyQuestion.
func (fuo *FileUpdateOne) SetPhotoSurveyQuestion(s *SurveyQuestion) *FileUpdateOne {
	return fuo.SetPhotoSurveyQuestionID(s.ID)
}

// SetSurveyQuestionID sets the survey_question edge to SurveyQuestion by id.
func (fuo *FileUpdateOne) SetSurveyQuestionID(id int) *FileUpdateOne {
	fuo.mutation.SetSurveyQuestionID(id)
	return fuo
}

// SetNillableSurveyQuestionID sets the survey_question edge to SurveyQuestion by id if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableSurveyQuestionID(id *int) *FileUpdateOne {
	if id != nil {
		fuo = fuo.SetSurveyQuestionID(*id)
	}
	return fuo
}

// SetSurveyQuestion sets the survey_question edge to SurveyQuestion.
func (fuo *FileUpdateOne) SetSurveyQuestion(s *SurveyQuestion) *FileUpdateOne {
	return fuo.SetSurveyQuestionID(s.ID)
}

// ClearLocation clears the location edge to Location.
func (fuo *FileUpdateOne) ClearLocation() *FileUpdateOne {
	fuo.mutation.ClearLocation()
	return fuo
}

// ClearEquipment clears the equipment edge to Equipment.
func (fuo *FileUpdateOne) ClearEquipment() *FileUpdateOne {
	fuo.mutation.ClearEquipment()
	return fuo
}

// ClearUser clears the user edge to User.
func (fuo *FileUpdateOne) ClearUser() *FileUpdateOne {
	fuo.mutation.ClearUser()
	return fuo
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (fuo *FileUpdateOne) ClearWorkOrder() *FileUpdateOne {
	fuo.mutation.ClearWorkOrder()
	return fuo
}

// ClearChecklistItem clears the checklist_item edge to CheckListItem.
func (fuo *FileUpdateOne) ClearChecklistItem() *FileUpdateOne {
	fuo.mutation.ClearChecklistItem()
	return fuo
}

// ClearSurvey clears the survey edge to Survey.
func (fuo *FileUpdateOne) ClearSurvey() *FileUpdateOne {
	fuo.mutation.ClearSurvey()
	return fuo
}

// ClearFloorPlan clears the floor_plan edge to FloorPlan.
func (fuo *FileUpdateOne) ClearFloorPlan() *FileUpdateOne {
	fuo.mutation.ClearFloorPlan()
	return fuo
}

// ClearPhotoSurveyQuestion clears the photo_survey_question edge to SurveyQuestion.
func (fuo *FileUpdateOne) ClearPhotoSurveyQuestion() *FileUpdateOne {
	fuo.mutation.ClearPhotoSurveyQuestion()
	return fuo
}

// ClearSurveyQuestion clears the survey_question edge to SurveyQuestion.
func (fuo *FileUpdateOne) ClearSurveyQuestion() *FileUpdateOne {
	fuo.mutation.ClearSurveyQuestion()
	return fuo
}

// Save executes the query and returns the updated entity.
func (fuo *FileUpdateOne) Save(ctx context.Context) (*File, error) {
	if _, ok := fuo.mutation.UpdateTime(); !ok {
		v := file.UpdateDefaultUpdateTime()
		fuo.mutation.SetUpdateTime(v)
	}
	if v, ok := fuo.mutation.Size(); ok {
		if err := file.SizeValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"size\": %v", err)
		}
	}

	var (
		err  error
		node *File
	)
	if len(fuo.hooks) == 0 {
		node, err = fuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FileMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fuo.mutation = mutation
			node, err = fuo.sqlSave(ctx)
			return node, err
		})
		for i := len(fuo.hooks) - 1; i >= 0; i-- {
			mut = fuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (fuo *FileUpdateOne) SaveX(ctx context.Context) *File {
	f, err := fuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return f
}

// Exec executes the query on the entity.
func (fuo *FileUpdateOne) Exec(ctx context.Context) error {
	_, err := fuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fuo *FileUpdateOne) ExecX(ctx context.Context) {
	if err := fuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fuo *FileUpdateOne) sqlSave(ctx context.Context) (f *File, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   file.Table,
			Columns: file.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: file.FieldID,
			},
		},
	}
	id, ok := fuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing File.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := fuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldUpdateTime,
		})
	}
	if value, ok := fuo.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldType,
		})
	}
	if value, ok := fuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldName,
		})
	}
	if value, ok := fuo.mutation.Size(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: file.FieldSize,
		})
	}
	if value, ok := fuo.mutation.AddedSize(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: file.FieldSize,
		})
	}
	if fuo.mutation.SizeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: file.FieldSize,
		})
	}
	if value, ok := fuo.mutation.ModifiedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldModifiedAt,
		})
	}
	if fuo.mutation.ModifiedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: file.FieldModifiedAt,
		})
	}
	if value, ok := fuo.mutation.UploadedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldUploadedAt,
		})
	}
	if fuo.mutation.UploadedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: file.FieldUploadedAt,
		})
	}
	if value, ok := fuo.mutation.ContentType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldContentType,
		})
	}
	if value, ok := fuo.mutation.StoreKey(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldStoreKey,
		})
	}
	if value, ok := fuo.mutation.Category(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldCategory,
		})
	}
	if fuo.mutation.CategoryCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: file.FieldCategory,
		})
	}
	if fuo.mutation.LocationCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.LocationIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fuo.mutation.EquipmentCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.EquipmentIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fuo.mutation.UserCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.UserIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fuo.mutation.WorkOrderCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.WorkOrderIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fuo.mutation.ChecklistItemCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.ChecklistItemIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fuo.mutation.SurveyCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.SurveyIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fuo.mutation.FloorPlanCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.FloorPlanIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fuo.mutation.PhotoSurveyQuestionCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.PhotoSurveyQuestionIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fuo.mutation.SurveyQuestionCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.SurveyQuestionIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	f = &File{config: fuo.config}
	_spec.Assign = f.assignValues
	_spec.ScanValues = f.scanValues()
	if err = sqlgraph.UpdateNode(ctx, fuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{file.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return f, nil
}
