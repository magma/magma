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
	"github.com/facebookincubator/symphony/graph/ent/checklistcategory"
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/comment"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/technician"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// WorkOrderUpdate is the builder for updating WorkOrder entities.
type WorkOrderUpdate struct {
	config
	hooks      []Hook
	mutation   *WorkOrderMutation
	predicates []predicate.WorkOrder
}

// Where adds a new predicate for the builder.
func (wou *WorkOrderUpdate) Where(ps ...predicate.WorkOrder) *WorkOrderUpdate {
	wou.predicates = append(wou.predicates, ps...)
	return wou
}

// SetName sets the name field.
func (wou *WorkOrderUpdate) SetName(s string) *WorkOrderUpdate {
	wou.mutation.SetName(s)
	return wou
}

// SetStatus sets the status field.
func (wou *WorkOrderUpdate) SetStatus(s string) *WorkOrderUpdate {
	wou.mutation.SetStatus(s)
	return wou
}

// SetNillableStatus sets the status field if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillableStatus(s *string) *WorkOrderUpdate {
	if s != nil {
		wou.SetStatus(*s)
	}
	return wou
}

// SetPriority sets the priority field.
func (wou *WorkOrderUpdate) SetPriority(s string) *WorkOrderUpdate {
	wou.mutation.SetPriority(s)
	return wou
}

// SetNillablePriority sets the priority field if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillablePriority(s *string) *WorkOrderUpdate {
	if s != nil {
		wou.SetPriority(*s)
	}
	return wou
}

// SetDescription sets the description field.
func (wou *WorkOrderUpdate) SetDescription(s string) *WorkOrderUpdate {
	wou.mutation.SetDescription(s)
	return wou
}

// SetNillableDescription sets the description field if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillableDescription(s *string) *WorkOrderUpdate {
	if s != nil {
		wou.SetDescription(*s)
	}
	return wou
}

// ClearDescription clears the value of description.
func (wou *WorkOrderUpdate) ClearDescription() *WorkOrderUpdate {
	wou.mutation.ClearDescription()
	return wou
}

// SetOwnerName sets the owner_name field.
func (wou *WorkOrderUpdate) SetOwnerName(s string) *WorkOrderUpdate {
	wou.mutation.SetOwnerName(s)
	return wou
}

// SetInstallDate sets the install_date field.
func (wou *WorkOrderUpdate) SetInstallDate(t time.Time) *WorkOrderUpdate {
	wou.mutation.SetInstallDate(t)
	return wou
}

// SetNillableInstallDate sets the install_date field if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillableInstallDate(t *time.Time) *WorkOrderUpdate {
	if t != nil {
		wou.SetInstallDate(*t)
	}
	return wou
}

// ClearInstallDate clears the value of install_date.
func (wou *WorkOrderUpdate) ClearInstallDate() *WorkOrderUpdate {
	wou.mutation.ClearInstallDate()
	return wou
}

// SetCreationDate sets the creation_date field.
func (wou *WorkOrderUpdate) SetCreationDate(t time.Time) *WorkOrderUpdate {
	wou.mutation.SetCreationDate(t)
	return wou
}

// SetAssigneeName sets the assignee_name field.
func (wou *WorkOrderUpdate) SetAssigneeName(s string) *WorkOrderUpdate {
	wou.mutation.SetAssigneeName(s)
	return wou
}

// SetNillableAssigneeName sets the assignee_name field if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillableAssigneeName(s *string) *WorkOrderUpdate {
	if s != nil {
		wou.SetAssigneeName(*s)
	}
	return wou
}

// ClearAssigneeName clears the value of assignee_name.
func (wou *WorkOrderUpdate) ClearAssigneeName() *WorkOrderUpdate {
	wou.mutation.ClearAssigneeName()
	return wou
}

// SetIndex sets the index field.
func (wou *WorkOrderUpdate) SetIndex(i int) *WorkOrderUpdate {
	wou.mutation.ResetIndex()
	wou.mutation.SetIndex(i)
	return wou
}

// SetNillableIndex sets the index field if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillableIndex(i *int) *WorkOrderUpdate {
	if i != nil {
		wou.SetIndex(*i)
	}
	return wou
}

// AddIndex adds i to index.
func (wou *WorkOrderUpdate) AddIndex(i int) *WorkOrderUpdate {
	wou.mutation.AddIndex(i)
	return wou
}

// ClearIndex clears the value of index.
func (wou *WorkOrderUpdate) ClearIndex() *WorkOrderUpdate {
	wou.mutation.ClearIndex()
	return wou
}

// SetCloseDate sets the close_date field.
func (wou *WorkOrderUpdate) SetCloseDate(t time.Time) *WorkOrderUpdate {
	wou.mutation.SetCloseDate(t)
	return wou
}

// SetNillableCloseDate sets the close_date field if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillableCloseDate(t *time.Time) *WorkOrderUpdate {
	if t != nil {
		wou.SetCloseDate(*t)
	}
	return wou
}

// ClearCloseDate clears the value of close_date.
func (wou *WorkOrderUpdate) ClearCloseDate() *WorkOrderUpdate {
	wou.mutation.ClearCloseDate()
	return wou
}

// SetTypeID sets the type edge to WorkOrderType by id.
func (wou *WorkOrderUpdate) SetTypeID(id int) *WorkOrderUpdate {
	wou.mutation.SetTypeID(id)
	return wou
}

// SetNillableTypeID sets the type edge to WorkOrderType by id if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillableTypeID(id *int) *WorkOrderUpdate {
	if id != nil {
		wou = wou.SetTypeID(*id)
	}
	return wou
}

// SetType sets the type edge to WorkOrderType.
func (wou *WorkOrderUpdate) SetType(w *WorkOrderType) *WorkOrderUpdate {
	return wou.SetTypeID(w.ID)
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (wou *WorkOrderUpdate) AddEquipmentIDs(ids ...int) *WorkOrderUpdate {
	wou.mutation.AddEquipmentIDs(ids...)
	return wou
}

// AddEquipment adds the equipment edges to Equipment.
func (wou *WorkOrderUpdate) AddEquipment(e ...*Equipment) *WorkOrderUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return wou.AddEquipmentIDs(ids...)
}

// AddLinkIDs adds the links edge to Link by ids.
func (wou *WorkOrderUpdate) AddLinkIDs(ids ...int) *WorkOrderUpdate {
	wou.mutation.AddLinkIDs(ids...)
	return wou
}

// AddLinks adds the links edges to Link.
func (wou *WorkOrderUpdate) AddLinks(l ...*Link) *WorkOrderUpdate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return wou.AddLinkIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (wou *WorkOrderUpdate) AddFileIDs(ids ...int) *WorkOrderUpdate {
	wou.mutation.AddFileIDs(ids...)
	return wou
}

// AddFiles adds the files edges to File.
func (wou *WorkOrderUpdate) AddFiles(f ...*File) *WorkOrderUpdate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return wou.AddFileIDs(ids...)
}

// AddHyperlinkIDs adds the hyperlinks edge to Hyperlink by ids.
func (wou *WorkOrderUpdate) AddHyperlinkIDs(ids ...int) *WorkOrderUpdate {
	wou.mutation.AddHyperlinkIDs(ids...)
	return wou
}

// AddHyperlinks adds the hyperlinks edges to Hyperlink.
func (wou *WorkOrderUpdate) AddHyperlinks(h ...*Hyperlink) *WorkOrderUpdate {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return wou.AddHyperlinkIDs(ids...)
}

// SetLocationID sets the location edge to Location by id.
func (wou *WorkOrderUpdate) SetLocationID(id int) *WorkOrderUpdate {
	wou.mutation.SetLocationID(id)
	return wou
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillableLocationID(id *int) *WorkOrderUpdate {
	if id != nil {
		wou = wou.SetLocationID(*id)
	}
	return wou
}

// SetLocation sets the location edge to Location.
func (wou *WorkOrderUpdate) SetLocation(l *Location) *WorkOrderUpdate {
	return wou.SetLocationID(l.ID)
}

// AddCommentIDs adds the comments edge to Comment by ids.
func (wou *WorkOrderUpdate) AddCommentIDs(ids ...int) *WorkOrderUpdate {
	wou.mutation.AddCommentIDs(ids...)
	return wou
}

// AddComments adds the comments edges to Comment.
func (wou *WorkOrderUpdate) AddComments(c ...*Comment) *WorkOrderUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wou.AddCommentIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (wou *WorkOrderUpdate) AddPropertyIDs(ids ...int) *WorkOrderUpdate {
	wou.mutation.AddPropertyIDs(ids...)
	return wou
}

// AddProperties adds the properties edges to Property.
func (wou *WorkOrderUpdate) AddProperties(p ...*Property) *WorkOrderUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wou.AddPropertyIDs(ids...)
}

// AddCheckListCategoryIDs adds the check_list_categories edge to CheckListCategory by ids.
func (wou *WorkOrderUpdate) AddCheckListCategoryIDs(ids ...int) *WorkOrderUpdate {
	wou.mutation.AddCheckListCategoryIDs(ids...)
	return wou
}

// AddCheckListCategories adds the check_list_categories edges to CheckListCategory.
func (wou *WorkOrderUpdate) AddCheckListCategories(c ...*CheckListCategory) *WorkOrderUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wou.AddCheckListCategoryIDs(ids...)
}

// AddCheckListItemIDs adds the check_list_items edge to CheckListItem by ids.
func (wou *WorkOrderUpdate) AddCheckListItemIDs(ids ...int) *WorkOrderUpdate {
	wou.mutation.AddCheckListItemIDs(ids...)
	return wou
}

// AddCheckListItems adds the check_list_items edges to CheckListItem.
func (wou *WorkOrderUpdate) AddCheckListItems(c ...*CheckListItem) *WorkOrderUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wou.AddCheckListItemIDs(ids...)
}

// SetTechnicianID sets the technician edge to Technician by id.
func (wou *WorkOrderUpdate) SetTechnicianID(id int) *WorkOrderUpdate {
	wou.mutation.SetTechnicianID(id)
	return wou
}

// SetNillableTechnicianID sets the technician edge to Technician by id if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillableTechnicianID(id *int) *WorkOrderUpdate {
	if id != nil {
		wou = wou.SetTechnicianID(*id)
	}
	return wou
}

// SetTechnician sets the technician edge to Technician.
func (wou *WorkOrderUpdate) SetTechnician(t *Technician) *WorkOrderUpdate {
	return wou.SetTechnicianID(t.ID)
}

// SetProjectID sets the project edge to Project by id.
func (wou *WorkOrderUpdate) SetProjectID(id int) *WorkOrderUpdate {
	wou.mutation.SetProjectID(id)
	return wou
}

// SetNillableProjectID sets the project edge to Project by id if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillableProjectID(id *int) *WorkOrderUpdate {
	if id != nil {
		wou = wou.SetProjectID(*id)
	}
	return wou
}

// SetProject sets the project edge to Project.
func (wou *WorkOrderUpdate) SetProject(p *Project) *WorkOrderUpdate {
	return wou.SetProjectID(p.ID)
}

// SetOwnerID sets the owner edge to User by id.
func (wou *WorkOrderUpdate) SetOwnerID(id int) *WorkOrderUpdate {
	wou.mutation.SetOwnerID(id)
	return wou
}

// SetNillableOwnerID sets the owner edge to User by id if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillableOwnerID(id *int) *WorkOrderUpdate {
	if id != nil {
		wou = wou.SetOwnerID(*id)
	}
	return wou
}

// SetOwner sets the owner edge to User.
func (wou *WorkOrderUpdate) SetOwner(u *User) *WorkOrderUpdate {
	return wou.SetOwnerID(u.ID)
}

// SetAssigneeID sets the assignee edge to User by id.
func (wou *WorkOrderUpdate) SetAssigneeID(id int) *WorkOrderUpdate {
	wou.mutation.SetAssigneeID(id)
	return wou
}

// SetNillableAssigneeID sets the assignee edge to User by id if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillableAssigneeID(id *int) *WorkOrderUpdate {
	if id != nil {
		wou = wou.SetAssigneeID(*id)
	}
	return wou
}

// SetAssignee sets the assignee edge to User.
func (wou *WorkOrderUpdate) SetAssignee(u *User) *WorkOrderUpdate {
	return wou.SetAssigneeID(u.ID)
}

// ClearType clears the type edge to WorkOrderType.
func (wou *WorkOrderUpdate) ClearType() *WorkOrderUpdate {
	wou.mutation.ClearType()
	return wou
}

// RemoveEquipmentIDs removes the equipment edge to Equipment by ids.
func (wou *WorkOrderUpdate) RemoveEquipmentIDs(ids ...int) *WorkOrderUpdate {
	wou.mutation.RemoveEquipmentIDs(ids...)
	return wou
}

// RemoveEquipment removes equipment edges to Equipment.
func (wou *WorkOrderUpdate) RemoveEquipment(e ...*Equipment) *WorkOrderUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return wou.RemoveEquipmentIDs(ids...)
}

// RemoveLinkIDs removes the links edge to Link by ids.
func (wou *WorkOrderUpdate) RemoveLinkIDs(ids ...int) *WorkOrderUpdate {
	wou.mutation.RemoveLinkIDs(ids...)
	return wou
}

// RemoveLinks removes links edges to Link.
func (wou *WorkOrderUpdate) RemoveLinks(l ...*Link) *WorkOrderUpdate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return wou.RemoveLinkIDs(ids...)
}

// RemoveFileIDs removes the files edge to File by ids.
func (wou *WorkOrderUpdate) RemoveFileIDs(ids ...int) *WorkOrderUpdate {
	wou.mutation.RemoveFileIDs(ids...)
	return wou
}

// RemoveFiles removes files edges to File.
func (wou *WorkOrderUpdate) RemoveFiles(f ...*File) *WorkOrderUpdate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return wou.RemoveFileIDs(ids...)
}

// RemoveHyperlinkIDs removes the hyperlinks edge to Hyperlink by ids.
func (wou *WorkOrderUpdate) RemoveHyperlinkIDs(ids ...int) *WorkOrderUpdate {
	wou.mutation.RemoveHyperlinkIDs(ids...)
	return wou
}

// RemoveHyperlinks removes hyperlinks edges to Hyperlink.
func (wou *WorkOrderUpdate) RemoveHyperlinks(h ...*Hyperlink) *WorkOrderUpdate {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return wou.RemoveHyperlinkIDs(ids...)
}

// ClearLocation clears the location edge to Location.
func (wou *WorkOrderUpdate) ClearLocation() *WorkOrderUpdate {
	wou.mutation.ClearLocation()
	return wou
}

// RemoveCommentIDs removes the comments edge to Comment by ids.
func (wou *WorkOrderUpdate) RemoveCommentIDs(ids ...int) *WorkOrderUpdate {
	wou.mutation.RemoveCommentIDs(ids...)
	return wou
}

// RemoveComments removes comments edges to Comment.
func (wou *WorkOrderUpdate) RemoveComments(c ...*Comment) *WorkOrderUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wou.RemoveCommentIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (wou *WorkOrderUpdate) RemovePropertyIDs(ids ...int) *WorkOrderUpdate {
	wou.mutation.RemovePropertyIDs(ids...)
	return wou
}

// RemoveProperties removes properties edges to Property.
func (wou *WorkOrderUpdate) RemoveProperties(p ...*Property) *WorkOrderUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wou.RemovePropertyIDs(ids...)
}

// RemoveCheckListCategoryIDs removes the check_list_categories edge to CheckListCategory by ids.
func (wou *WorkOrderUpdate) RemoveCheckListCategoryIDs(ids ...int) *WorkOrderUpdate {
	wou.mutation.RemoveCheckListCategoryIDs(ids...)
	return wou
}

// RemoveCheckListCategories removes check_list_categories edges to CheckListCategory.
func (wou *WorkOrderUpdate) RemoveCheckListCategories(c ...*CheckListCategory) *WorkOrderUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wou.RemoveCheckListCategoryIDs(ids...)
}

// RemoveCheckListItemIDs removes the check_list_items edge to CheckListItem by ids.
func (wou *WorkOrderUpdate) RemoveCheckListItemIDs(ids ...int) *WorkOrderUpdate {
	wou.mutation.RemoveCheckListItemIDs(ids...)
	return wou
}

// RemoveCheckListItems removes check_list_items edges to CheckListItem.
func (wou *WorkOrderUpdate) RemoveCheckListItems(c ...*CheckListItem) *WorkOrderUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wou.RemoveCheckListItemIDs(ids...)
}

// ClearTechnician clears the technician edge to Technician.
func (wou *WorkOrderUpdate) ClearTechnician() *WorkOrderUpdate {
	wou.mutation.ClearTechnician()
	return wou
}

// ClearProject clears the project edge to Project.
func (wou *WorkOrderUpdate) ClearProject() *WorkOrderUpdate {
	wou.mutation.ClearProject()
	return wou
}

// ClearOwner clears the owner edge to User.
func (wou *WorkOrderUpdate) ClearOwner() *WorkOrderUpdate {
	wou.mutation.ClearOwner()
	return wou
}

// ClearAssignee clears the assignee edge to User.
func (wou *WorkOrderUpdate) ClearAssignee() *WorkOrderUpdate {
	wou.mutation.ClearAssignee()
	return wou
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (wou *WorkOrderUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := wou.mutation.UpdateTime(); !ok {
		v := workorder.UpdateDefaultUpdateTime()
		wou.mutation.SetUpdateTime(v)
	}
	if v, ok := wou.mutation.Name(); ok {
		if err := workorder.NameValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}

	var (
		err      error
		affected int
	)
	if len(wou.hooks) == 0 {
		affected, err = wou.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*WorkOrderMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			wou.mutation = mutation
			affected, err = wou.sqlSave(ctx)
			return affected, err
		})
		for i := len(wou.hooks); i > 0; i-- {
			mut = wou.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, wou.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (wou *WorkOrderUpdate) SaveX(ctx context.Context) int {
	affected, err := wou.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (wou *WorkOrderUpdate) Exec(ctx context.Context) error {
	_, err := wou.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wou *WorkOrderUpdate) ExecX(ctx context.Context) {
	if err := wou.Exec(ctx); err != nil {
		panic(err)
	}
}

func (wou *WorkOrderUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   workorder.Table,
			Columns: workorder.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: workorder.FieldID,
			},
		},
	}
	if ps := wou.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := wou.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: workorder.FieldUpdateTime,
		})
	}
	if value, ok := wou.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workorder.FieldName,
		})
	}
	if value, ok := wou.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workorder.FieldStatus,
		})
	}
	if value, ok := wou.mutation.Priority(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workorder.FieldPriority,
		})
	}
	if value, ok := wou.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workorder.FieldDescription,
		})
	}
	if wou.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: workorder.FieldDescription,
		})
	}
	if value, ok := wou.mutation.OwnerName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workorder.FieldOwnerName,
		})
	}
	if value, ok := wou.mutation.InstallDate(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: workorder.FieldInstallDate,
		})
	}
	if wou.mutation.InstallDateCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: workorder.FieldInstallDate,
		})
	}
	if value, ok := wou.mutation.CreationDate(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: workorder.FieldCreationDate,
		})
	}
	if value, ok := wou.mutation.AssigneeName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workorder.FieldAssigneeName,
		})
	}
	if wou.mutation.AssigneeNameCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: workorder.FieldAssigneeName,
		})
	}
	if value, ok := wou.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: workorder.FieldIndex,
		})
	}
	if value, ok := wou.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: workorder.FieldIndex,
		})
	}
	if wou.mutation.IndexCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: workorder.FieldIndex,
		})
	}
	if value, ok := wou.mutation.CloseDate(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: workorder.FieldCloseDate,
		})
	}
	if wou.mutation.CloseDateCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: workorder.FieldCloseDate,
		})
	}
	if wou.mutation.TypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.TypeTable,
			Columns: []string{workorder.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workordertype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wou.mutation.TypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.TypeTable,
			Columns: []string{workorder.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workordertype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wou.mutation.RemovedEquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workorder.EquipmentTable,
			Columns: []string{workorder.EquipmentColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wou.mutation.EquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workorder.EquipmentTable,
			Columns: []string{workorder.EquipmentColumn},
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
	if nodes := wou.mutation.RemovedLinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workorder.LinksTable,
			Columns: []string{workorder.LinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: link.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wou.mutation.LinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workorder.LinksTable,
			Columns: []string{workorder.LinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: link.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wou.mutation.RemovedFilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.FilesTable,
			Columns: []string{workorder.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: file.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wou.mutation.FilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.FilesTable,
			Columns: []string{workorder.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: file.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wou.mutation.RemovedHyperlinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.HyperlinksTable,
			Columns: []string{workorder.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: hyperlink.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wou.mutation.HyperlinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.HyperlinksTable,
			Columns: []string{workorder.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: hyperlink.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if wou.mutation.LocationCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.LocationTable,
			Columns: []string{workorder.LocationColumn},
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
	if nodes := wou.mutation.LocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.LocationTable,
			Columns: []string{workorder.LocationColumn},
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
	if nodes := wou.mutation.RemovedCommentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CommentsTable,
			Columns: []string{workorder.CommentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: comment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wou.mutation.CommentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CommentsTable,
			Columns: []string{workorder.CommentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: comment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wou.mutation.RemovedPropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.PropertiesTable,
			Columns: []string{workorder.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wou.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.PropertiesTable,
			Columns: []string{workorder.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wou.mutation.RemovedCheckListCategoriesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CheckListCategoriesTable,
			Columns: []string{workorder.CheckListCategoriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistcategory.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wou.mutation.CheckListCategoriesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CheckListCategoriesTable,
			Columns: []string{workorder.CheckListCategoriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistcategory.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wou.mutation.RemovedCheckListItemsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CheckListItemsTable,
			Columns: []string{workorder.CheckListItemsColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wou.mutation.CheckListItemsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CheckListItemsTable,
			Columns: []string{workorder.CheckListItemsColumn},
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
	if wou.mutation.TechnicianCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.TechnicianTable,
			Columns: []string{workorder.TechnicianColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: technician.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wou.mutation.TechnicianIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.TechnicianTable,
			Columns: []string{workorder.TechnicianColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: technician.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if wou.mutation.ProjectCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   workorder.ProjectTable,
			Columns: []string{workorder.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: project.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wou.mutation.ProjectIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   workorder.ProjectTable,
			Columns: []string{workorder.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: project.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if wou.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.OwnerTable,
			Columns: []string{workorder.OwnerColumn},
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
	if nodes := wou.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.OwnerTable,
			Columns: []string{workorder.OwnerColumn},
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
	if wou.mutation.AssigneeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.AssigneeTable,
			Columns: []string{workorder.AssigneeColumn},
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
	if nodes := wou.mutation.AssigneeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.AssigneeTable,
			Columns: []string{workorder.AssigneeColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, wou.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{workorder.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// WorkOrderUpdateOne is the builder for updating a single WorkOrder entity.
type WorkOrderUpdateOne struct {
	config
	hooks    []Hook
	mutation *WorkOrderMutation
}

// SetName sets the name field.
func (wouo *WorkOrderUpdateOne) SetName(s string) *WorkOrderUpdateOne {
	wouo.mutation.SetName(s)
	return wouo
}

// SetStatus sets the status field.
func (wouo *WorkOrderUpdateOne) SetStatus(s string) *WorkOrderUpdateOne {
	wouo.mutation.SetStatus(s)
	return wouo
}

// SetNillableStatus sets the status field if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillableStatus(s *string) *WorkOrderUpdateOne {
	if s != nil {
		wouo.SetStatus(*s)
	}
	return wouo
}

// SetPriority sets the priority field.
func (wouo *WorkOrderUpdateOne) SetPriority(s string) *WorkOrderUpdateOne {
	wouo.mutation.SetPriority(s)
	return wouo
}

// SetNillablePriority sets the priority field if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillablePriority(s *string) *WorkOrderUpdateOne {
	if s != nil {
		wouo.SetPriority(*s)
	}
	return wouo
}

// SetDescription sets the description field.
func (wouo *WorkOrderUpdateOne) SetDescription(s string) *WorkOrderUpdateOne {
	wouo.mutation.SetDescription(s)
	return wouo
}

// SetNillableDescription sets the description field if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillableDescription(s *string) *WorkOrderUpdateOne {
	if s != nil {
		wouo.SetDescription(*s)
	}
	return wouo
}

// ClearDescription clears the value of description.
func (wouo *WorkOrderUpdateOne) ClearDescription() *WorkOrderUpdateOne {
	wouo.mutation.ClearDescription()
	return wouo
}

// SetOwnerName sets the owner_name field.
func (wouo *WorkOrderUpdateOne) SetOwnerName(s string) *WorkOrderUpdateOne {
	wouo.mutation.SetOwnerName(s)
	return wouo
}

// SetInstallDate sets the install_date field.
func (wouo *WorkOrderUpdateOne) SetInstallDate(t time.Time) *WorkOrderUpdateOne {
	wouo.mutation.SetInstallDate(t)
	return wouo
}

// SetNillableInstallDate sets the install_date field if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillableInstallDate(t *time.Time) *WorkOrderUpdateOne {
	if t != nil {
		wouo.SetInstallDate(*t)
	}
	return wouo
}

// ClearInstallDate clears the value of install_date.
func (wouo *WorkOrderUpdateOne) ClearInstallDate() *WorkOrderUpdateOne {
	wouo.mutation.ClearInstallDate()
	return wouo
}

// SetCreationDate sets the creation_date field.
func (wouo *WorkOrderUpdateOne) SetCreationDate(t time.Time) *WorkOrderUpdateOne {
	wouo.mutation.SetCreationDate(t)
	return wouo
}

// SetAssigneeName sets the assignee_name field.
func (wouo *WorkOrderUpdateOne) SetAssigneeName(s string) *WorkOrderUpdateOne {
	wouo.mutation.SetAssigneeName(s)
	return wouo
}

// SetNillableAssigneeName sets the assignee_name field if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillableAssigneeName(s *string) *WorkOrderUpdateOne {
	if s != nil {
		wouo.SetAssigneeName(*s)
	}
	return wouo
}

// ClearAssigneeName clears the value of assignee_name.
func (wouo *WorkOrderUpdateOne) ClearAssigneeName() *WorkOrderUpdateOne {
	wouo.mutation.ClearAssigneeName()
	return wouo
}

// SetIndex sets the index field.
func (wouo *WorkOrderUpdateOne) SetIndex(i int) *WorkOrderUpdateOne {
	wouo.mutation.ResetIndex()
	wouo.mutation.SetIndex(i)
	return wouo
}

// SetNillableIndex sets the index field if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillableIndex(i *int) *WorkOrderUpdateOne {
	if i != nil {
		wouo.SetIndex(*i)
	}
	return wouo
}

// AddIndex adds i to index.
func (wouo *WorkOrderUpdateOne) AddIndex(i int) *WorkOrderUpdateOne {
	wouo.mutation.AddIndex(i)
	return wouo
}

// ClearIndex clears the value of index.
func (wouo *WorkOrderUpdateOne) ClearIndex() *WorkOrderUpdateOne {
	wouo.mutation.ClearIndex()
	return wouo
}

// SetCloseDate sets the close_date field.
func (wouo *WorkOrderUpdateOne) SetCloseDate(t time.Time) *WorkOrderUpdateOne {
	wouo.mutation.SetCloseDate(t)
	return wouo
}

// SetNillableCloseDate sets the close_date field if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillableCloseDate(t *time.Time) *WorkOrderUpdateOne {
	if t != nil {
		wouo.SetCloseDate(*t)
	}
	return wouo
}

// ClearCloseDate clears the value of close_date.
func (wouo *WorkOrderUpdateOne) ClearCloseDate() *WorkOrderUpdateOne {
	wouo.mutation.ClearCloseDate()
	return wouo
}

// SetTypeID sets the type edge to WorkOrderType by id.
func (wouo *WorkOrderUpdateOne) SetTypeID(id int) *WorkOrderUpdateOne {
	wouo.mutation.SetTypeID(id)
	return wouo
}

// SetNillableTypeID sets the type edge to WorkOrderType by id if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillableTypeID(id *int) *WorkOrderUpdateOne {
	if id != nil {
		wouo = wouo.SetTypeID(*id)
	}
	return wouo
}

// SetType sets the type edge to WorkOrderType.
func (wouo *WorkOrderUpdateOne) SetType(w *WorkOrderType) *WorkOrderUpdateOne {
	return wouo.SetTypeID(w.ID)
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (wouo *WorkOrderUpdateOne) AddEquipmentIDs(ids ...int) *WorkOrderUpdateOne {
	wouo.mutation.AddEquipmentIDs(ids...)
	return wouo
}

// AddEquipment adds the equipment edges to Equipment.
func (wouo *WorkOrderUpdateOne) AddEquipment(e ...*Equipment) *WorkOrderUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return wouo.AddEquipmentIDs(ids...)
}

// AddLinkIDs adds the links edge to Link by ids.
func (wouo *WorkOrderUpdateOne) AddLinkIDs(ids ...int) *WorkOrderUpdateOne {
	wouo.mutation.AddLinkIDs(ids...)
	return wouo
}

// AddLinks adds the links edges to Link.
func (wouo *WorkOrderUpdateOne) AddLinks(l ...*Link) *WorkOrderUpdateOne {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return wouo.AddLinkIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (wouo *WorkOrderUpdateOne) AddFileIDs(ids ...int) *WorkOrderUpdateOne {
	wouo.mutation.AddFileIDs(ids...)
	return wouo
}

// AddFiles adds the files edges to File.
func (wouo *WorkOrderUpdateOne) AddFiles(f ...*File) *WorkOrderUpdateOne {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return wouo.AddFileIDs(ids...)
}

// AddHyperlinkIDs adds the hyperlinks edge to Hyperlink by ids.
func (wouo *WorkOrderUpdateOne) AddHyperlinkIDs(ids ...int) *WorkOrderUpdateOne {
	wouo.mutation.AddHyperlinkIDs(ids...)
	return wouo
}

// AddHyperlinks adds the hyperlinks edges to Hyperlink.
func (wouo *WorkOrderUpdateOne) AddHyperlinks(h ...*Hyperlink) *WorkOrderUpdateOne {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return wouo.AddHyperlinkIDs(ids...)
}

// SetLocationID sets the location edge to Location by id.
func (wouo *WorkOrderUpdateOne) SetLocationID(id int) *WorkOrderUpdateOne {
	wouo.mutation.SetLocationID(id)
	return wouo
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillableLocationID(id *int) *WorkOrderUpdateOne {
	if id != nil {
		wouo = wouo.SetLocationID(*id)
	}
	return wouo
}

// SetLocation sets the location edge to Location.
func (wouo *WorkOrderUpdateOne) SetLocation(l *Location) *WorkOrderUpdateOne {
	return wouo.SetLocationID(l.ID)
}

// AddCommentIDs adds the comments edge to Comment by ids.
func (wouo *WorkOrderUpdateOne) AddCommentIDs(ids ...int) *WorkOrderUpdateOne {
	wouo.mutation.AddCommentIDs(ids...)
	return wouo
}

// AddComments adds the comments edges to Comment.
func (wouo *WorkOrderUpdateOne) AddComments(c ...*Comment) *WorkOrderUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wouo.AddCommentIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (wouo *WorkOrderUpdateOne) AddPropertyIDs(ids ...int) *WorkOrderUpdateOne {
	wouo.mutation.AddPropertyIDs(ids...)
	return wouo
}

// AddProperties adds the properties edges to Property.
func (wouo *WorkOrderUpdateOne) AddProperties(p ...*Property) *WorkOrderUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wouo.AddPropertyIDs(ids...)
}

// AddCheckListCategoryIDs adds the check_list_categories edge to CheckListCategory by ids.
func (wouo *WorkOrderUpdateOne) AddCheckListCategoryIDs(ids ...int) *WorkOrderUpdateOne {
	wouo.mutation.AddCheckListCategoryIDs(ids...)
	return wouo
}

// AddCheckListCategories adds the check_list_categories edges to CheckListCategory.
func (wouo *WorkOrderUpdateOne) AddCheckListCategories(c ...*CheckListCategory) *WorkOrderUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wouo.AddCheckListCategoryIDs(ids...)
}

// AddCheckListItemIDs adds the check_list_items edge to CheckListItem by ids.
func (wouo *WorkOrderUpdateOne) AddCheckListItemIDs(ids ...int) *WorkOrderUpdateOne {
	wouo.mutation.AddCheckListItemIDs(ids...)
	return wouo
}

// AddCheckListItems adds the check_list_items edges to CheckListItem.
func (wouo *WorkOrderUpdateOne) AddCheckListItems(c ...*CheckListItem) *WorkOrderUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wouo.AddCheckListItemIDs(ids...)
}

// SetTechnicianID sets the technician edge to Technician by id.
func (wouo *WorkOrderUpdateOne) SetTechnicianID(id int) *WorkOrderUpdateOne {
	wouo.mutation.SetTechnicianID(id)
	return wouo
}

// SetNillableTechnicianID sets the technician edge to Technician by id if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillableTechnicianID(id *int) *WorkOrderUpdateOne {
	if id != nil {
		wouo = wouo.SetTechnicianID(*id)
	}
	return wouo
}

// SetTechnician sets the technician edge to Technician.
func (wouo *WorkOrderUpdateOne) SetTechnician(t *Technician) *WorkOrderUpdateOne {
	return wouo.SetTechnicianID(t.ID)
}

// SetProjectID sets the project edge to Project by id.
func (wouo *WorkOrderUpdateOne) SetProjectID(id int) *WorkOrderUpdateOne {
	wouo.mutation.SetProjectID(id)
	return wouo
}

// SetNillableProjectID sets the project edge to Project by id if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillableProjectID(id *int) *WorkOrderUpdateOne {
	if id != nil {
		wouo = wouo.SetProjectID(*id)
	}
	return wouo
}

// SetProject sets the project edge to Project.
func (wouo *WorkOrderUpdateOne) SetProject(p *Project) *WorkOrderUpdateOne {
	return wouo.SetProjectID(p.ID)
}

// SetOwnerID sets the owner edge to User by id.
func (wouo *WorkOrderUpdateOne) SetOwnerID(id int) *WorkOrderUpdateOne {
	wouo.mutation.SetOwnerID(id)
	return wouo
}

// SetNillableOwnerID sets the owner edge to User by id if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillableOwnerID(id *int) *WorkOrderUpdateOne {
	if id != nil {
		wouo = wouo.SetOwnerID(*id)
	}
	return wouo
}

// SetOwner sets the owner edge to User.
func (wouo *WorkOrderUpdateOne) SetOwner(u *User) *WorkOrderUpdateOne {
	return wouo.SetOwnerID(u.ID)
}

// SetAssigneeID sets the assignee edge to User by id.
func (wouo *WorkOrderUpdateOne) SetAssigneeID(id int) *WorkOrderUpdateOne {
	wouo.mutation.SetAssigneeID(id)
	return wouo
}

// SetNillableAssigneeID sets the assignee edge to User by id if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillableAssigneeID(id *int) *WorkOrderUpdateOne {
	if id != nil {
		wouo = wouo.SetAssigneeID(*id)
	}
	return wouo
}

// SetAssignee sets the assignee edge to User.
func (wouo *WorkOrderUpdateOne) SetAssignee(u *User) *WorkOrderUpdateOne {
	return wouo.SetAssigneeID(u.ID)
}

// ClearType clears the type edge to WorkOrderType.
func (wouo *WorkOrderUpdateOne) ClearType() *WorkOrderUpdateOne {
	wouo.mutation.ClearType()
	return wouo
}

// RemoveEquipmentIDs removes the equipment edge to Equipment by ids.
func (wouo *WorkOrderUpdateOne) RemoveEquipmentIDs(ids ...int) *WorkOrderUpdateOne {
	wouo.mutation.RemoveEquipmentIDs(ids...)
	return wouo
}

// RemoveEquipment removes equipment edges to Equipment.
func (wouo *WorkOrderUpdateOne) RemoveEquipment(e ...*Equipment) *WorkOrderUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return wouo.RemoveEquipmentIDs(ids...)
}

// RemoveLinkIDs removes the links edge to Link by ids.
func (wouo *WorkOrderUpdateOne) RemoveLinkIDs(ids ...int) *WorkOrderUpdateOne {
	wouo.mutation.RemoveLinkIDs(ids...)
	return wouo
}

// RemoveLinks removes links edges to Link.
func (wouo *WorkOrderUpdateOne) RemoveLinks(l ...*Link) *WorkOrderUpdateOne {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return wouo.RemoveLinkIDs(ids...)
}

// RemoveFileIDs removes the files edge to File by ids.
func (wouo *WorkOrderUpdateOne) RemoveFileIDs(ids ...int) *WorkOrderUpdateOne {
	wouo.mutation.RemoveFileIDs(ids...)
	return wouo
}

// RemoveFiles removes files edges to File.
func (wouo *WorkOrderUpdateOne) RemoveFiles(f ...*File) *WorkOrderUpdateOne {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return wouo.RemoveFileIDs(ids...)
}

// RemoveHyperlinkIDs removes the hyperlinks edge to Hyperlink by ids.
func (wouo *WorkOrderUpdateOne) RemoveHyperlinkIDs(ids ...int) *WorkOrderUpdateOne {
	wouo.mutation.RemoveHyperlinkIDs(ids...)
	return wouo
}

// RemoveHyperlinks removes hyperlinks edges to Hyperlink.
func (wouo *WorkOrderUpdateOne) RemoveHyperlinks(h ...*Hyperlink) *WorkOrderUpdateOne {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return wouo.RemoveHyperlinkIDs(ids...)
}

// ClearLocation clears the location edge to Location.
func (wouo *WorkOrderUpdateOne) ClearLocation() *WorkOrderUpdateOne {
	wouo.mutation.ClearLocation()
	return wouo
}

// RemoveCommentIDs removes the comments edge to Comment by ids.
func (wouo *WorkOrderUpdateOne) RemoveCommentIDs(ids ...int) *WorkOrderUpdateOne {
	wouo.mutation.RemoveCommentIDs(ids...)
	return wouo
}

// RemoveComments removes comments edges to Comment.
func (wouo *WorkOrderUpdateOne) RemoveComments(c ...*Comment) *WorkOrderUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wouo.RemoveCommentIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (wouo *WorkOrderUpdateOne) RemovePropertyIDs(ids ...int) *WorkOrderUpdateOne {
	wouo.mutation.RemovePropertyIDs(ids...)
	return wouo
}

// RemoveProperties removes properties edges to Property.
func (wouo *WorkOrderUpdateOne) RemoveProperties(p ...*Property) *WorkOrderUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wouo.RemovePropertyIDs(ids...)
}

// RemoveCheckListCategoryIDs removes the check_list_categories edge to CheckListCategory by ids.
func (wouo *WorkOrderUpdateOne) RemoveCheckListCategoryIDs(ids ...int) *WorkOrderUpdateOne {
	wouo.mutation.RemoveCheckListCategoryIDs(ids...)
	return wouo
}

// RemoveCheckListCategories removes check_list_categories edges to CheckListCategory.
func (wouo *WorkOrderUpdateOne) RemoveCheckListCategories(c ...*CheckListCategory) *WorkOrderUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wouo.RemoveCheckListCategoryIDs(ids...)
}

// RemoveCheckListItemIDs removes the check_list_items edge to CheckListItem by ids.
func (wouo *WorkOrderUpdateOne) RemoveCheckListItemIDs(ids ...int) *WorkOrderUpdateOne {
	wouo.mutation.RemoveCheckListItemIDs(ids...)
	return wouo
}

// RemoveCheckListItems removes check_list_items edges to CheckListItem.
func (wouo *WorkOrderUpdateOne) RemoveCheckListItems(c ...*CheckListItem) *WorkOrderUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wouo.RemoveCheckListItemIDs(ids...)
}

// ClearTechnician clears the technician edge to Technician.
func (wouo *WorkOrderUpdateOne) ClearTechnician() *WorkOrderUpdateOne {
	wouo.mutation.ClearTechnician()
	return wouo
}

// ClearProject clears the project edge to Project.
func (wouo *WorkOrderUpdateOne) ClearProject() *WorkOrderUpdateOne {
	wouo.mutation.ClearProject()
	return wouo
}

// ClearOwner clears the owner edge to User.
func (wouo *WorkOrderUpdateOne) ClearOwner() *WorkOrderUpdateOne {
	wouo.mutation.ClearOwner()
	return wouo
}

// ClearAssignee clears the assignee edge to User.
func (wouo *WorkOrderUpdateOne) ClearAssignee() *WorkOrderUpdateOne {
	wouo.mutation.ClearAssignee()
	return wouo
}

// Save executes the query and returns the updated entity.
func (wouo *WorkOrderUpdateOne) Save(ctx context.Context) (*WorkOrder, error) {
	if _, ok := wouo.mutation.UpdateTime(); !ok {
		v := workorder.UpdateDefaultUpdateTime()
		wouo.mutation.SetUpdateTime(v)
	}
	if v, ok := wouo.mutation.Name(); ok {
		if err := workorder.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}

	var (
		err  error
		node *WorkOrder
	)
	if len(wouo.hooks) == 0 {
		node, err = wouo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*WorkOrderMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			wouo.mutation = mutation
			node, err = wouo.sqlSave(ctx)
			return node, err
		})
		for i := len(wouo.hooks); i > 0; i-- {
			mut = wouo.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, wouo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (wouo *WorkOrderUpdateOne) SaveX(ctx context.Context) *WorkOrder {
	wo, err := wouo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return wo
}

// Exec executes the query on the entity.
func (wouo *WorkOrderUpdateOne) Exec(ctx context.Context) error {
	_, err := wouo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wouo *WorkOrderUpdateOne) ExecX(ctx context.Context) {
	if err := wouo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (wouo *WorkOrderUpdateOne) sqlSave(ctx context.Context) (wo *WorkOrder, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   workorder.Table,
			Columns: workorder.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: workorder.FieldID,
			},
		},
	}
	id, ok := wouo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing WorkOrder.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := wouo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: workorder.FieldUpdateTime,
		})
	}
	if value, ok := wouo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workorder.FieldName,
		})
	}
	if value, ok := wouo.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workorder.FieldStatus,
		})
	}
	if value, ok := wouo.mutation.Priority(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workorder.FieldPriority,
		})
	}
	if value, ok := wouo.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workorder.FieldDescription,
		})
	}
	if wouo.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: workorder.FieldDescription,
		})
	}
	if value, ok := wouo.mutation.OwnerName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workorder.FieldOwnerName,
		})
	}
	if value, ok := wouo.mutation.InstallDate(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: workorder.FieldInstallDate,
		})
	}
	if wouo.mutation.InstallDateCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: workorder.FieldInstallDate,
		})
	}
	if value, ok := wouo.mutation.CreationDate(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: workorder.FieldCreationDate,
		})
	}
	if value, ok := wouo.mutation.AssigneeName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: workorder.FieldAssigneeName,
		})
	}
	if wouo.mutation.AssigneeNameCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: workorder.FieldAssigneeName,
		})
	}
	if value, ok := wouo.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: workorder.FieldIndex,
		})
	}
	if value, ok := wouo.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: workorder.FieldIndex,
		})
	}
	if wouo.mutation.IndexCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: workorder.FieldIndex,
		})
	}
	if value, ok := wouo.mutation.CloseDate(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: workorder.FieldCloseDate,
		})
	}
	if wouo.mutation.CloseDateCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: workorder.FieldCloseDate,
		})
	}
	if wouo.mutation.TypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.TypeTable,
			Columns: []string{workorder.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workordertype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wouo.mutation.TypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.TypeTable,
			Columns: []string{workorder.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workordertype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wouo.mutation.RemovedEquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workorder.EquipmentTable,
			Columns: []string{workorder.EquipmentColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wouo.mutation.EquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workorder.EquipmentTable,
			Columns: []string{workorder.EquipmentColumn},
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
	if nodes := wouo.mutation.RemovedLinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workorder.LinksTable,
			Columns: []string{workorder.LinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: link.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wouo.mutation.LinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workorder.LinksTable,
			Columns: []string{workorder.LinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: link.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wouo.mutation.RemovedFilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.FilesTable,
			Columns: []string{workorder.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: file.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wouo.mutation.FilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.FilesTable,
			Columns: []string{workorder.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: file.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wouo.mutation.RemovedHyperlinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.HyperlinksTable,
			Columns: []string{workorder.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: hyperlink.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wouo.mutation.HyperlinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.HyperlinksTable,
			Columns: []string{workorder.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: hyperlink.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if wouo.mutation.LocationCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.LocationTable,
			Columns: []string{workorder.LocationColumn},
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
	if nodes := wouo.mutation.LocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.LocationTable,
			Columns: []string{workorder.LocationColumn},
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
	if nodes := wouo.mutation.RemovedCommentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CommentsTable,
			Columns: []string{workorder.CommentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: comment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wouo.mutation.CommentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CommentsTable,
			Columns: []string{workorder.CommentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: comment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wouo.mutation.RemovedPropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.PropertiesTable,
			Columns: []string{workorder.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wouo.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.PropertiesTable,
			Columns: []string{workorder.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wouo.mutation.RemovedCheckListCategoriesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CheckListCategoriesTable,
			Columns: []string{workorder.CheckListCategoriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistcategory.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wouo.mutation.CheckListCategoriesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CheckListCategoriesTable,
			Columns: []string{workorder.CheckListCategoriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistcategory.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := wouo.mutation.RemovedCheckListItemsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CheckListItemsTable,
			Columns: []string{workorder.CheckListItemsColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wouo.mutation.CheckListItemsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CheckListItemsTable,
			Columns: []string{workorder.CheckListItemsColumn},
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
	if wouo.mutation.TechnicianCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.TechnicianTable,
			Columns: []string{workorder.TechnicianColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: technician.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wouo.mutation.TechnicianIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.TechnicianTable,
			Columns: []string{workorder.TechnicianColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: technician.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if wouo.mutation.ProjectCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   workorder.ProjectTable,
			Columns: []string{workorder.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: project.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wouo.mutation.ProjectIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   workorder.ProjectTable,
			Columns: []string{workorder.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: project.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if wouo.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.OwnerTable,
			Columns: []string{workorder.OwnerColumn},
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
	if nodes := wouo.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.OwnerTable,
			Columns: []string{workorder.OwnerColumn},
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
	if wouo.mutation.AssigneeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.AssigneeTable,
			Columns: []string{workorder.AssigneeColumn},
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
	if nodes := wouo.mutation.AssigneeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.AssigneeTable,
			Columns: []string{workorder.AssigneeColumn},
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
	wo = &WorkOrder{config: wouo.config}
	_spec.Assign = wo.assignValues
	_spec.ScanValues = wo.scanValues()
	if err = sqlgraph.UpdateNode(ctx, wouo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{workorder.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return wo, nil
}
