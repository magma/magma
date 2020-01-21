// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
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
	"github.com/facebookincubator/symphony/graph/ent/workorder"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// WorkOrderUpdate is the builder for updating WorkOrder entities.
type WorkOrderUpdate struct {
	config

	update_time           *time.Time
	name                  *string
	status                *string
	priority              *string
	description           *string
	cleardescription      bool
	owner_name            *string
	install_date          *time.Time
	clearinstall_date     bool
	creation_date         *time.Time
	assignee              *string
	clearassignee         bool
	index                 *int
	addindex              *int
	clearindex            bool
	_type                 map[string]struct{}
	equipment             map[string]struct{}
	links                 map[string]struct{}
	files                 map[string]struct{}
	hyperlinks            map[string]struct{}
	location              map[string]struct{}
	comments              map[string]struct{}
	properties            map[string]struct{}
	check_list_items      map[string]struct{}
	technician            map[string]struct{}
	project               map[string]struct{}
	clearedType           bool
	removedEquipment      map[string]struct{}
	removedLinks          map[string]struct{}
	removedFiles          map[string]struct{}
	removedHyperlinks     map[string]struct{}
	clearedLocation       bool
	removedComments       map[string]struct{}
	removedProperties     map[string]struct{}
	removedCheckListItems map[string]struct{}
	clearedTechnician     bool
	clearedProject        bool
	predicates            []predicate.WorkOrder
}

// Where adds a new predicate for the builder.
func (wou *WorkOrderUpdate) Where(ps ...predicate.WorkOrder) *WorkOrderUpdate {
	wou.predicates = append(wou.predicates, ps...)
	return wou
}

// SetName sets the name field.
func (wou *WorkOrderUpdate) SetName(s string) *WorkOrderUpdate {
	wou.name = &s
	return wou
}

// SetStatus sets the status field.
func (wou *WorkOrderUpdate) SetStatus(s string) *WorkOrderUpdate {
	wou.status = &s
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
	wou.priority = &s
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
	wou.description = &s
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
	wou.description = nil
	wou.cleardescription = true
	return wou
}

// SetOwnerName sets the owner_name field.
func (wou *WorkOrderUpdate) SetOwnerName(s string) *WorkOrderUpdate {
	wou.owner_name = &s
	return wou
}

// SetInstallDate sets the install_date field.
func (wou *WorkOrderUpdate) SetInstallDate(t time.Time) *WorkOrderUpdate {
	wou.install_date = &t
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
	wou.install_date = nil
	wou.clearinstall_date = true
	return wou
}

// SetCreationDate sets the creation_date field.
func (wou *WorkOrderUpdate) SetCreationDate(t time.Time) *WorkOrderUpdate {
	wou.creation_date = &t
	return wou
}

// SetAssignee sets the assignee field.
func (wou *WorkOrderUpdate) SetAssignee(s string) *WorkOrderUpdate {
	wou.assignee = &s
	return wou
}

// SetNillableAssignee sets the assignee field if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillableAssignee(s *string) *WorkOrderUpdate {
	if s != nil {
		wou.SetAssignee(*s)
	}
	return wou
}

// ClearAssignee clears the value of assignee.
func (wou *WorkOrderUpdate) ClearAssignee() *WorkOrderUpdate {
	wou.assignee = nil
	wou.clearassignee = true
	return wou
}

// SetIndex sets the index field.
func (wou *WorkOrderUpdate) SetIndex(i int) *WorkOrderUpdate {
	wou.index = &i
	wou.addindex = nil
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
	if wou.addindex == nil {
		wou.addindex = &i
	} else {
		*wou.addindex += i
	}
	return wou
}

// ClearIndex clears the value of index.
func (wou *WorkOrderUpdate) ClearIndex() *WorkOrderUpdate {
	wou.index = nil
	wou.clearindex = true
	return wou
}

// SetTypeID sets the type edge to WorkOrderType by id.
func (wou *WorkOrderUpdate) SetTypeID(id string) *WorkOrderUpdate {
	if wou._type == nil {
		wou._type = make(map[string]struct{})
	}
	wou._type[id] = struct{}{}
	return wou
}

// SetNillableTypeID sets the type edge to WorkOrderType by id if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillableTypeID(id *string) *WorkOrderUpdate {
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
func (wou *WorkOrderUpdate) AddEquipmentIDs(ids ...string) *WorkOrderUpdate {
	if wou.equipment == nil {
		wou.equipment = make(map[string]struct{})
	}
	for i := range ids {
		wou.equipment[ids[i]] = struct{}{}
	}
	return wou
}

// AddEquipment adds the equipment edges to Equipment.
func (wou *WorkOrderUpdate) AddEquipment(e ...*Equipment) *WorkOrderUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return wou.AddEquipmentIDs(ids...)
}

// AddLinkIDs adds the links edge to Link by ids.
func (wou *WorkOrderUpdate) AddLinkIDs(ids ...string) *WorkOrderUpdate {
	if wou.links == nil {
		wou.links = make(map[string]struct{})
	}
	for i := range ids {
		wou.links[ids[i]] = struct{}{}
	}
	return wou
}

// AddLinks adds the links edges to Link.
func (wou *WorkOrderUpdate) AddLinks(l ...*Link) *WorkOrderUpdate {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return wou.AddLinkIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (wou *WorkOrderUpdate) AddFileIDs(ids ...string) *WorkOrderUpdate {
	if wou.files == nil {
		wou.files = make(map[string]struct{})
	}
	for i := range ids {
		wou.files[ids[i]] = struct{}{}
	}
	return wou
}

// AddFiles adds the files edges to File.
func (wou *WorkOrderUpdate) AddFiles(f ...*File) *WorkOrderUpdate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return wou.AddFileIDs(ids...)
}

// AddHyperlinkIDs adds the hyperlinks edge to Hyperlink by ids.
func (wou *WorkOrderUpdate) AddHyperlinkIDs(ids ...string) *WorkOrderUpdate {
	if wou.hyperlinks == nil {
		wou.hyperlinks = make(map[string]struct{})
	}
	for i := range ids {
		wou.hyperlinks[ids[i]] = struct{}{}
	}
	return wou
}

// AddHyperlinks adds the hyperlinks edges to Hyperlink.
func (wou *WorkOrderUpdate) AddHyperlinks(h ...*Hyperlink) *WorkOrderUpdate {
	ids := make([]string, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return wou.AddHyperlinkIDs(ids...)
}

// SetLocationID sets the location edge to Location by id.
func (wou *WorkOrderUpdate) SetLocationID(id string) *WorkOrderUpdate {
	if wou.location == nil {
		wou.location = make(map[string]struct{})
	}
	wou.location[id] = struct{}{}
	return wou
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillableLocationID(id *string) *WorkOrderUpdate {
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
func (wou *WorkOrderUpdate) AddCommentIDs(ids ...string) *WorkOrderUpdate {
	if wou.comments == nil {
		wou.comments = make(map[string]struct{})
	}
	for i := range ids {
		wou.comments[ids[i]] = struct{}{}
	}
	return wou
}

// AddComments adds the comments edges to Comment.
func (wou *WorkOrderUpdate) AddComments(c ...*Comment) *WorkOrderUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wou.AddCommentIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (wou *WorkOrderUpdate) AddPropertyIDs(ids ...string) *WorkOrderUpdate {
	if wou.properties == nil {
		wou.properties = make(map[string]struct{})
	}
	for i := range ids {
		wou.properties[ids[i]] = struct{}{}
	}
	return wou
}

// AddProperties adds the properties edges to Property.
func (wou *WorkOrderUpdate) AddProperties(p ...*Property) *WorkOrderUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wou.AddPropertyIDs(ids...)
}

// AddCheckListItemIDs adds the check_list_items edge to CheckListItem by ids.
func (wou *WorkOrderUpdate) AddCheckListItemIDs(ids ...string) *WorkOrderUpdate {
	if wou.check_list_items == nil {
		wou.check_list_items = make(map[string]struct{})
	}
	for i := range ids {
		wou.check_list_items[ids[i]] = struct{}{}
	}
	return wou
}

// AddCheckListItems adds the check_list_items edges to CheckListItem.
func (wou *WorkOrderUpdate) AddCheckListItems(c ...*CheckListItem) *WorkOrderUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wou.AddCheckListItemIDs(ids...)
}

// SetTechnicianID sets the technician edge to Technician by id.
func (wou *WorkOrderUpdate) SetTechnicianID(id string) *WorkOrderUpdate {
	if wou.technician == nil {
		wou.technician = make(map[string]struct{})
	}
	wou.technician[id] = struct{}{}
	return wou
}

// SetNillableTechnicianID sets the technician edge to Technician by id if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillableTechnicianID(id *string) *WorkOrderUpdate {
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
func (wou *WorkOrderUpdate) SetProjectID(id string) *WorkOrderUpdate {
	if wou.project == nil {
		wou.project = make(map[string]struct{})
	}
	wou.project[id] = struct{}{}
	return wou
}

// SetNillableProjectID sets the project edge to Project by id if the given value is not nil.
func (wou *WorkOrderUpdate) SetNillableProjectID(id *string) *WorkOrderUpdate {
	if id != nil {
		wou = wou.SetProjectID(*id)
	}
	return wou
}

// SetProject sets the project edge to Project.
func (wou *WorkOrderUpdate) SetProject(p *Project) *WorkOrderUpdate {
	return wou.SetProjectID(p.ID)
}

// ClearType clears the type edge to WorkOrderType.
func (wou *WorkOrderUpdate) ClearType() *WorkOrderUpdate {
	wou.clearedType = true
	return wou
}

// RemoveEquipmentIDs removes the equipment edge to Equipment by ids.
func (wou *WorkOrderUpdate) RemoveEquipmentIDs(ids ...string) *WorkOrderUpdate {
	if wou.removedEquipment == nil {
		wou.removedEquipment = make(map[string]struct{})
	}
	for i := range ids {
		wou.removedEquipment[ids[i]] = struct{}{}
	}
	return wou
}

// RemoveEquipment removes equipment edges to Equipment.
func (wou *WorkOrderUpdate) RemoveEquipment(e ...*Equipment) *WorkOrderUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return wou.RemoveEquipmentIDs(ids...)
}

// RemoveLinkIDs removes the links edge to Link by ids.
func (wou *WorkOrderUpdate) RemoveLinkIDs(ids ...string) *WorkOrderUpdate {
	if wou.removedLinks == nil {
		wou.removedLinks = make(map[string]struct{})
	}
	for i := range ids {
		wou.removedLinks[ids[i]] = struct{}{}
	}
	return wou
}

// RemoveLinks removes links edges to Link.
func (wou *WorkOrderUpdate) RemoveLinks(l ...*Link) *WorkOrderUpdate {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return wou.RemoveLinkIDs(ids...)
}

// RemoveFileIDs removes the files edge to File by ids.
func (wou *WorkOrderUpdate) RemoveFileIDs(ids ...string) *WorkOrderUpdate {
	if wou.removedFiles == nil {
		wou.removedFiles = make(map[string]struct{})
	}
	for i := range ids {
		wou.removedFiles[ids[i]] = struct{}{}
	}
	return wou
}

// RemoveFiles removes files edges to File.
func (wou *WorkOrderUpdate) RemoveFiles(f ...*File) *WorkOrderUpdate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return wou.RemoveFileIDs(ids...)
}

// RemoveHyperlinkIDs removes the hyperlinks edge to Hyperlink by ids.
func (wou *WorkOrderUpdate) RemoveHyperlinkIDs(ids ...string) *WorkOrderUpdate {
	if wou.removedHyperlinks == nil {
		wou.removedHyperlinks = make(map[string]struct{})
	}
	for i := range ids {
		wou.removedHyperlinks[ids[i]] = struct{}{}
	}
	return wou
}

// RemoveHyperlinks removes hyperlinks edges to Hyperlink.
func (wou *WorkOrderUpdate) RemoveHyperlinks(h ...*Hyperlink) *WorkOrderUpdate {
	ids := make([]string, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return wou.RemoveHyperlinkIDs(ids...)
}

// ClearLocation clears the location edge to Location.
func (wou *WorkOrderUpdate) ClearLocation() *WorkOrderUpdate {
	wou.clearedLocation = true
	return wou
}

// RemoveCommentIDs removes the comments edge to Comment by ids.
func (wou *WorkOrderUpdate) RemoveCommentIDs(ids ...string) *WorkOrderUpdate {
	if wou.removedComments == nil {
		wou.removedComments = make(map[string]struct{})
	}
	for i := range ids {
		wou.removedComments[ids[i]] = struct{}{}
	}
	return wou
}

// RemoveComments removes comments edges to Comment.
func (wou *WorkOrderUpdate) RemoveComments(c ...*Comment) *WorkOrderUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wou.RemoveCommentIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (wou *WorkOrderUpdate) RemovePropertyIDs(ids ...string) *WorkOrderUpdate {
	if wou.removedProperties == nil {
		wou.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		wou.removedProperties[ids[i]] = struct{}{}
	}
	return wou
}

// RemoveProperties removes properties edges to Property.
func (wou *WorkOrderUpdate) RemoveProperties(p ...*Property) *WorkOrderUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wou.RemovePropertyIDs(ids...)
}

// RemoveCheckListItemIDs removes the check_list_items edge to CheckListItem by ids.
func (wou *WorkOrderUpdate) RemoveCheckListItemIDs(ids ...string) *WorkOrderUpdate {
	if wou.removedCheckListItems == nil {
		wou.removedCheckListItems = make(map[string]struct{})
	}
	for i := range ids {
		wou.removedCheckListItems[ids[i]] = struct{}{}
	}
	return wou
}

// RemoveCheckListItems removes check_list_items edges to CheckListItem.
func (wou *WorkOrderUpdate) RemoveCheckListItems(c ...*CheckListItem) *WorkOrderUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wou.RemoveCheckListItemIDs(ids...)
}

// ClearTechnician clears the technician edge to Technician.
func (wou *WorkOrderUpdate) ClearTechnician() *WorkOrderUpdate {
	wou.clearedTechnician = true
	return wou
}

// ClearProject clears the project edge to Project.
func (wou *WorkOrderUpdate) ClearProject() *WorkOrderUpdate {
	wou.clearedProject = true
	return wou
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (wou *WorkOrderUpdate) Save(ctx context.Context) (int, error) {
	if wou.update_time == nil {
		v := workorder.UpdateDefaultUpdateTime()
		wou.update_time = &v
	}
	if wou.name != nil {
		if err := workorder.NameValidator(*wou.name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if len(wou._type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if len(wou.location) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	if len(wou.technician) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"technician\"")
	}
	if len(wou.project) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"project\"")
	}
	return wou.sqlSave(ctx)
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
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   workorder.Table,
			Columns: workorder.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: workorder.FieldID,
			},
		},
	}
	if ps := wou.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := wou.update_time; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: workorder.FieldUpdateTime,
		})
	}
	if value := wou.name; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: workorder.FieldName,
		})
	}
	if value := wou.status; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: workorder.FieldStatus,
		})
	}
	if value := wou.priority; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: workorder.FieldPriority,
		})
	}
	if value := wou.description; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: workorder.FieldDescription,
		})
	}
	if wou.cleardescription {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: workorder.FieldDescription,
		})
	}
	if value := wou.owner_name; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: workorder.FieldOwnerName,
		})
	}
	if value := wou.install_date; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: workorder.FieldInstallDate,
		})
	}
	if wou.clearinstall_date {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: workorder.FieldInstallDate,
		})
	}
	if value := wou.creation_date; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: workorder.FieldCreationDate,
		})
	}
	if value := wou.assignee; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: workorder.FieldAssignee,
		})
	}
	if wou.clearassignee {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: workorder.FieldAssignee,
		})
	}
	if value := wou.index; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: workorder.FieldIndex,
		})
	}
	if value := wou.addindex; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: workorder.FieldIndex,
		})
	}
	if wou.clearindex {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: workorder.FieldIndex,
		})
	}
	if wou.clearedType {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.TypeTable,
			Columns: []string{workorder.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workordertype.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wou._type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.TypeTable,
			Columns: []string{workorder.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workordertype.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := wou.removedEquipment; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workorder.EquipmentTable,
			Columns: []string{workorder.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipment.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wou.equipment; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workorder.EquipmentTable,
			Columns: []string{workorder.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipment.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := wou.removedLinks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workorder.LinksTable,
			Columns: []string{workorder.LinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: link.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wou.links; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workorder.LinksTable,
			Columns: []string{workorder.LinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: link.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := wou.removedFiles; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.FilesTable,
			Columns: []string{workorder.FilesColumn},
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wou.files; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.FilesTable,
			Columns: []string{workorder.FilesColumn},
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := wou.removedHyperlinks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.HyperlinksTable,
			Columns: []string{workorder.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: hyperlink.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wou.hyperlinks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.HyperlinksTable,
			Columns: []string{workorder.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: hyperlink.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if wou.clearedLocation {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.LocationTable,
			Columns: []string{workorder.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wou.location; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.LocationTable,
			Columns: []string{workorder.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := wou.removedComments; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CommentsTable,
			Columns: []string{workorder.CommentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: comment.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wou.comments; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CommentsTable,
			Columns: []string{workorder.CommentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: comment.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := wou.removedProperties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.PropertiesTable,
			Columns: []string{workorder.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wou.properties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.PropertiesTable,
			Columns: []string{workorder.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := wou.removedCheckListItems; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CheckListItemsTable,
			Columns: []string{workorder.CheckListItemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: checklistitem.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wou.check_list_items; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CheckListItemsTable,
			Columns: []string{workorder.CheckListItemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: checklistitem.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if wou.clearedTechnician {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.TechnicianTable,
			Columns: []string{workorder.TechnicianColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: technician.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wou.technician; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.TechnicianTable,
			Columns: []string{workorder.TechnicianColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: technician.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if wou.clearedProject {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   workorder.ProjectTable,
			Columns: []string{workorder.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: project.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wou.project; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   workorder.ProjectTable,
			Columns: []string{workorder.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: project.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, wou.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// WorkOrderUpdateOne is the builder for updating a single WorkOrder entity.
type WorkOrderUpdateOne struct {
	config
	id string

	update_time           *time.Time
	name                  *string
	status                *string
	priority              *string
	description           *string
	cleardescription      bool
	owner_name            *string
	install_date          *time.Time
	clearinstall_date     bool
	creation_date         *time.Time
	assignee              *string
	clearassignee         bool
	index                 *int
	addindex              *int
	clearindex            bool
	_type                 map[string]struct{}
	equipment             map[string]struct{}
	links                 map[string]struct{}
	files                 map[string]struct{}
	hyperlinks            map[string]struct{}
	location              map[string]struct{}
	comments              map[string]struct{}
	properties            map[string]struct{}
	check_list_items      map[string]struct{}
	technician            map[string]struct{}
	project               map[string]struct{}
	clearedType           bool
	removedEquipment      map[string]struct{}
	removedLinks          map[string]struct{}
	removedFiles          map[string]struct{}
	removedHyperlinks     map[string]struct{}
	clearedLocation       bool
	removedComments       map[string]struct{}
	removedProperties     map[string]struct{}
	removedCheckListItems map[string]struct{}
	clearedTechnician     bool
	clearedProject        bool
}

// SetName sets the name field.
func (wouo *WorkOrderUpdateOne) SetName(s string) *WorkOrderUpdateOne {
	wouo.name = &s
	return wouo
}

// SetStatus sets the status field.
func (wouo *WorkOrderUpdateOne) SetStatus(s string) *WorkOrderUpdateOne {
	wouo.status = &s
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
	wouo.priority = &s
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
	wouo.description = &s
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
	wouo.description = nil
	wouo.cleardescription = true
	return wouo
}

// SetOwnerName sets the owner_name field.
func (wouo *WorkOrderUpdateOne) SetOwnerName(s string) *WorkOrderUpdateOne {
	wouo.owner_name = &s
	return wouo
}

// SetInstallDate sets the install_date field.
func (wouo *WorkOrderUpdateOne) SetInstallDate(t time.Time) *WorkOrderUpdateOne {
	wouo.install_date = &t
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
	wouo.install_date = nil
	wouo.clearinstall_date = true
	return wouo
}

// SetCreationDate sets the creation_date field.
func (wouo *WorkOrderUpdateOne) SetCreationDate(t time.Time) *WorkOrderUpdateOne {
	wouo.creation_date = &t
	return wouo
}

// SetAssignee sets the assignee field.
func (wouo *WorkOrderUpdateOne) SetAssignee(s string) *WorkOrderUpdateOne {
	wouo.assignee = &s
	return wouo
}

// SetNillableAssignee sets the assignee field if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillableAssignee(s *string) *WorkOrderUpdateOne {
	if s != nil {
		wouo.SetAssignee(*s)
	}
	return wouo
}

// ClearAssignee clears the value of assignee.
func (wouo *WorkOrderUpdateOne) ClearAssignee() *WorkOrderUpdateOne {
	wouo.assignee = nil
	wouo.clearassignee = true
	return wouo
}

// SetIndex sets the index field.
func (wouo *WorkOrderUpdateOne) SetIndex(i int) *WorkOrderUpdateOne {
	wouo.index = &i
	wouo.addindex = nil
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
	if wouo.addindex == nil {
		wouo.addindex = &i
	} else {
		*wouo.addindex += i
	}
	return wouo
}

// ClearIndex clears the value of index.
func (wouo *WorkOrderUpdateOne) ClearIndex() *WorkOrderUpdateOne {
	wouo.index = nil
	wouo.clearindex = true
	return wouo
}

// SetTypeID sets the type edge to WorkOrderType by id.
func (wouo *WorkOrderUpdateOne) SetTypeID(id string) *WorkOrderUpdateOne {
	if wouo._type == nil {
		wouo._type = make(map[string]struct{})
	}
	wouo._type[id] = struct{}{}
	return wouo
}

// SetNillableTypeID sets the type edge to WorkOrderType by id if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillableTypeID(id *string) *WorkOrderUpdateOne {
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
func (wouo *WorkOrderUpdateOne) AddEquipmentIDs(ids ...string) *WorkOrderUpdateOne {
	if wouo.equipment == nil {
		wouo.equipment = make(map[string]struct{})
	}
	for i := range ids {
		wouo.equipment[ids[i]] = struct{}{}
	}
	return wouo
}

// AddEquipment adds the equipment edges to Equipment.
func (wouo *WorkOrderUpdateOne) AddEquipment(e ...*Equipment) *WorkOrderUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return wouo.AddEquipmentIDs(ids...)
}

// AddLinkIDs adds the links edge to Link by ids.
func (wouo *WorkOrderUpdateOne) AddLinkIDs(ids ...string) *WorkOrderUpdateOne {
	if wouo.links == nil {
		wouo.links = make(map[string]struct{})
	}
	for i := range ids {
		wouo.links[ids[i]] = struct{}{}
	}
	return wouo
}

// AddLinks adds the links edges to Link.
func (wouo *WorkOrderUpdateOne) AddLinks(l ...*Link) *WorkOrderUpdateOne {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return wouo.AddLinkIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (wouo *WorkOrderUpdateOne) AddFileIDs(ids ...string) *WorkOrderUpdateOne {
	if wouo.files == nil {
		wouo.files = make(map[string]struct{})
	}
	for i := range ids {
		wouo.files[ids[i]] = struct{}{}
	}
	return wouo
}

// AddFiles adds the files edges to File.
func (wouo *WorkOrderUpdateOne) AddFiles(f ...*File) *WorkOrderUpdateOne {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return wouo.AddFileIDs(ids...)
}

// AddHyperlinkIDs adds the hyperlinks edge to Hyperlink by ids.
func (wouo *WorkOrderUpdateOne) AddHyperlinkIDs(ids ...string) *WorkOrderUpdateOne {
	if wouo.hyperlinks == nil {
		wouo.hyperlinks = make(map[string]struct{})
	}
	for i := range ids {
		wouo.hyperlinks[ids[i]] = struct{}{}
	}
	return wouo
}

// AddHyperlinks adds the hyperlinks edges to Hyperlink.
func (wouo *WorkOrderUpdateOne) AddHyperlinks(h ...*Hyperlink) *WorkOrderUpdateOne {
	ids := make([]string, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return wouo.AddHyperlinkIDs(ids...)
}

// SetLocationID sets the location edge to Location by id.
func (wouo *WorkOrderUpdateOne) SetLocationID(id string) *WorkOrderUpdateOne {
	if wouo.location == nil {
		wouo.location = make(map[string]struct{})
	}
	wouo.location[id] = struct{}{}
	return wouo
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillableLocationID(id *string) *WorkOrderUpdateOne {
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
func (wouo *WorkOrderUpdateOne) AddCommentIDs(ids ...string) *WorkOrderUpdateOne {
	if wouo.comments == nil {
		wouo.comments = make(map[string]struct{})
	}
	for i := range ids {
		wouo.comments[ids[i]] = struct{}{}
	}
	return wouo
}

// AddComments adds the comments edges to Comment.
func (wouo *WorkOrderUpdateOne) AddComments(c ...*Comment) *WorkOrderUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wouo.AddCommentIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (wouo *WorkOrderUpdateOne) AddPropertyIDs(ids ...string) *WorkOrderUpdateOne {
	if wouo.properties == nil {
		wouo.properties = make(map[string]struct{})
	}
	for i := range ids {
		wouo.properties[ids[i]] = struct{}{}
	}
	return wouo
}

// AddProperties adds the properties edges to Property.
func (wouo *WorkOrderUpdateOne) AddProperties(p ...*Property) *WorkOrderUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wouo.AddPropertyIDs(ids...)
}

// AddCheckListItemIDs adds the check_list_items edge to CheckListItem by ids.
func (wouo *WorkOrderUpdateOne) AddCheckListItemIDs(ids ...string) *WorkOrderUpdateOne {
	if wouo.check_list_items == nil {
		wouo.check_list_items = make(map[string]struct{})
	}
	for i := range ids {
		wouo.check_list_items[ids[i]] = struct{}{}
	}
	return wouo
}

// AddCheckListItems adds the check_list_items edges to CheckListItem.
func (wouo *WorkOrderUpdateOne) AddCheckListItems(c ...*CheckListItem) *WorkOrderUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wouo.AddCheckListItemIDs(ids...)
}

// SetTechnicianID sets the technician edge to Technician by id.
func (wouo *WorkOrderUpdateOne) SetTechnicianID(id string) *WorkOrderUpdateOne {
	if wouo.technician == nil {
		wouo.technician = make(map[string]struct{})
	}
	wouo.technician[id] = struct{}{}
	return wouo
}

// SetNillableTechnicianID sets the technician edge to Technician by id if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillableTechnicianID(id *string) *WorkOrderUpdateOne {
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
func (wouo *WorkOrderUpdateOne) SetProjectID(id string) *WorkOrderUpdateOne {
	if wouo.project == nil {
		wouo.project = make(map[string]struct{})
	}
	wouo.project[id] = struct{}{}
	return wouo
}

// SetNillableProjectID sets the project edge to Project by id if the given value is not nil.
func (wouo *WorkOrderUpdateOne) SetNillableProjectID(id *string) *WorkOrderUpdateOne {
	if id != nil {
		wouo = wouo.SetProjectID(*id)
	}
	return wouo
}

// SetProject sets the project edge to Project.
func (wouo *WorkOrderUpdateOne) SetProject(p *Project) *WorkOrderUpdateOne {
	return wouo.SetProjectID(p.ID)
}

// ClearType clears the type edge to WorkOrderType.
func (wouo *WorkOrderUpdateOne) ClearType() *WorkOrderUpdateOne {
	wouo.clearedType = true
	return wouo
}

// RemoveEquipmentIDs removes the equipment edge to Equipment by ids.
func (wouo *WorkOrderUpdateOne) RemoveEquipmentIDs(ids ...string) *WorkOrderUpdateOne {
	if wouo.removedEquipment == nil {
		wouo.removedEquipment = make(map[string]struct{})
	}
	for i := range ids {
		wouo.removedEquipment[ids[i]] = struct{}{}
	}
	return wouo
}

// RemoveEquipment removes equipment edges to Equipment.
func (wouo *WorkOrderUpdateOne) RemoveEquipment(e ...*Equipment) *WorkOrderUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return wouo.RemoveEquipmentIDs(ids...)
}

// RemoveLinkIDs removes the links edge to Link by ids.
func (wouo *WorkOrderUpdateOne) RemoveLinkIDs(ids ...string) *WorkOrderUpdateOne {
	if wouo.removedLinks == nil {
		wouo.removedLinks = make(map[string]struct{})
	}
	for i := range ids {
		wouo.removedLinks[ids[i]] = struct{}{}
	}
	return wouo
}

// RemoveLinks removes links edges to Link.
func (wouo *WorkOrderUpdateOne) RemoveLinks(l ...*Link) *WorkOrderUpdateOne {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return wouo.RemoveLinkIDs(ids...)
}

// RemoveFileIDs removes the files edge to File by ids.
func (wouo *WorkOrderUpdateOne) RemoveFileIDs(ids ...string) *WorkOrderUpdateOne {
	if wouo.removedFiles == nil {
		wouo.removedFiles = make(map[string]struct{})
	}
	for i := range ids {
		wouo.removedFiles[ids[i]] = struct{}{}
	}
	return wouo
}

// RemoveFiles removes files edges to File.
func (wouo *WorkOrderUpdateOne) RemoveFiles(f ...*File) *WorkOrderUpdateOne {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return wouo.RemoveFileIDs(ids...)
}

// RemoveHyperlinkIDs removes the hyperlinks edge to Hyperlink by ids.
func (wouo *WorkOrderUpdateOne) RemoveHyperlinkIDs(ids ...string) *WorkOrderUpdateOne {
	if wouo.removedHyperlinks == nil {
		wouo.removedHyperlinks = make(map[string]struct{})
	}
	for i := range ids {
		wouo.removedHyperlinks[ids[i]] = struct{}{}
	}
	return wouo
}

// RemoveHyperlinks removes hyperlinks edges to Hyperlink.
func (wouo *WorkOrderUpdateOne) RemoveHyperlinks(h ...*Hyperlink) *WorkOrderUpdateOne {
	ids := make([]string, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return wouo.RemoveHyperlinkIDs(ids...)
}

// ClearLocation clears the location edge to Location.
func (wouo *WorkOrderUpdateOne) ClearLocation() *WorkOrderUpdateOne {
	wouo.clearedLocation = true
	return wouo
}

// RemoveCommentIDs removes the comments edge to Comment by ids.
func (wouo *WorkOrderUpdateOne) RemoveCommentIDs(ids ...string) *WorkOrderUpdateOne {
	if wouo.removedComments == nil {
		wouo.removedComments = make(map[string]struct{})
	}
	for i := range ids {
		wouo.removedComments[ids[i]] = struct{}{}
	}
	return wouo
}

// RemoveComments removes comments edges to Comment.
func (wouo *WorkOrderUpdateOne) RemoveComments(c ...*Comment) *WorkOrderUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wouo.RemoveCommentIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (wouo *WorkOrderUpdateOne) RemovePropertyIDs(ids ...string) *WorkOrderUpdateOne {
	if wouo.removedProperties == nil {
		wouo.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		wouo.removedProperties[ids[i]] = struct{}{}
	}
	return wouo
}

// RemoveProperties removes properties edges to Property.
func (wouo *WorkOrderUpdateOne) RemoveProperties(p ...*Property) *WorkOrderUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wouo.RemovePropertyIDs(ids...)
}

// RemoveCheckListItemIDs removes the check_list_items edge to CheckListItem by ids.
func (wouo *WorkOrderUpdateOne) RemoveCheckListItemIDs(ids ...string) *WorkOrderUpdateOne {
	if wouo.removedCheckListItems == nil {
		wouo.removedCheckListItems = make(map[string]struct{})
	}
	for i := range ids {
		wouo.removedCheckListItems[ids[i]] = struct{}{}
	}
	return wouo
}

// RemoveCheckListItems removes check_list_items edges to CheckListItem.
func (wouo *WorkOrderUpdateOne) RemoveCheckListItems(c ...*CheckListItem) *WorkOrderUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wouo.RemoveCheckListItemIDs(ids...)
}

// ClearTechnician clears the technician edge to Technician.
func (wouo *WorkOrderUpdateOne) ClearTechnician() *WorkOrderUpdateOne {
	wouo.clearedTechnician = true
	return wouo
}

// ClearProject clears the project edge to Project.
func (wouo *WorkOrderUpdateOne) ClearProject() *WorkOrderUpdateOne {
	wouo.clearedProject = true
	return wouo
}

// Save executes the query and returns the updated entity.
func (wouo *WorkOrderUpdateOne) Save(ctx context.Context) (*WorkOrder, error) {
	if wouo.update_time == nil {
		v := workorder.UpdateDefaultUpdateTime()
		wouo.update_time = &v
	}
	if wouo.name != nil {
		if err := workorder.NameValidator(*wouo.name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if len(wouo._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if len(wouo.location) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	if len(wouo.technician) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"technician\"")
	}
	if len(wouo.project) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"project\"")
	}
	return wouo.sqlSave(ctx)
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
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   workorder.Table,
			Columns: workorder.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  wouo.id,
				Type:   field.TypeString,
				Column: workorder.FieldID,
			},
		},
	}
	if value := wouo.update_time; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: workorder.FieldUpdateTime,
		})
	}
	if value := wouo.name; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: workorder.FieldName,
		})
	}
	if value := wouo.status; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: workorder.FieldStatus,
		})
	}
	if value := wouo.priority; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: workorder.FieldPriority,
		})
	}
	if value := wouo.description; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: workorder.FieldDescription,
		})
	}
	if wouo.cleardescription {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: workorder.FieldDescription,
		})
	}
	if value := wouo.owner_name; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: workorder.FieldOwnerName,
		})
	}
	if value := wouo.install_date; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: workorder.FieldInstallDate,
		})
	}
	if wouo.clearinstall_date {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: workorder.FieldInstallDate,
		})
	}
	if value := wouo.creation_date; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: workorder.FieldCreationDate,
		})
	}
	if value := wouo.assignee; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: workorder.FieldAssignee,
		})
	}
	if wouo.clearassignee {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: workorder.FieldAssignee,
		})
	}
	if value := wouo.index; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: workorder.FieldIndex,
		})
	}
	if value := wouo.addindex; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: workorder.FieldIndex,
		})
	}
	if wouo.clearindex {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: workorder.FieldIndex,
		})
	}
	if wouo.clearedType {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.TypeTable,
			Columns: []string{workorder.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workordertype.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wouo._type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.TypeTable,
			Columns: []string{workorder.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workordertype.FieldID,
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
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := wouo.removedEquipment; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workorder.EquipmentTable,
			Columns: []string{workorder.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipment.FieldID,
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
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wouo.equipment; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workorder.EquipmentTable,
			Columns: []string{workorder.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipment.FieldID,
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
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := wouo.removedLinks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workorder.LinksTable,
			Columns: []string{workorder.LinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: link.FieldID,
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
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wouo.links; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   workorder.LinksTable,
			Columns: []string{workorder.LinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: link.FieldID,
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
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := wouo.removedFiles; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.FilesTable,
			Columns: []string{workorder.FilesColumn},
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
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wouo.files; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.FilesTable,
			Columns: []string{workorder.FilesColumn},
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
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := wouo.removedHyperlinks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.HyperlinksTable,
			Columns: []string{workorder.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: hyperlink.FieldID,
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
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wouo.hyperlinks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.HyperlinksTable,
			Columns: []string{workorder.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: hyperlink.FieldID,
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
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if wouo.clearedLocation {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.LocationTable,
			Columns: []string{workorder.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wouo.location; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.LocationTable,
			Columns: []string{workorder.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
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
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := wouo.removedComments; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CommentsTable,
			Columns: []string{workorder.CommentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: comment.FieldID,
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
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wouo.comments; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CommentsTable,
			Columns: []string{workorder.CommentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: comment.FieldID,
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
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := wouo.removedProperties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.PropertiesTable,
			Columns: []string{workorder.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
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
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wouo.properties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.PropertiesTable,
			Columns: []string{workorder.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
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
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := wouo.removedCheckListItems; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CheckListItemsTable,
			Columns: []string{workorder.CheckListItemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: checklistitem.FieldID,
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
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wouo.check_list_items; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   workorder.CheckListItemsTable,
			Columns: []string{workorder.CheckListItemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: checklistitem.FieldID,
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
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if wouo.clearedTechnician {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.TechnicianTable,
			Columns: []string{workorder.TechnicianColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: technician.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wouo.technician; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorder.TechnicianTable,
			Columns: []string{workorder.TechnicianColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: technician.FieldID,
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
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if wouo.clearedProject {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   workorder.ProjectTable,
			Columns: []string{workorder.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: project.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := wouo.project; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   workorder.ProjectTable,
			Columns: []string{workorder.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: project.FieldID,
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
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	wo = &WorkOrder{config: wouo.config}
	spec.Assign = wo.assignValues
	spec.ScanValues = wo.scanValues()
	if err = sqlgraph.UpdateNode(ctx, wouo.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return wo, nil
}
