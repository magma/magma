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
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/comment"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// WorkOrderCreate is the builder for creating a WorkOrder entity.
type WorkOrderCreate struct {
	config
	create_time      *time.Time
	update_time      *time.Time
	name             *string
	status           *string
	priority         *string
	description      *string
	owner_name       *string
	install_date     *time.Time
	creation_date    *time.Time
	assignee         *string
	index            *int
	_type            map[string]struct{}
	equipment        map[string]struct{}
	links            map[string]struct{}
	files            map[string]struct{}
	location         map[string]struct{}
	comments         map[string]struct{}
	properties       map[string]struct{}
	check_list_items map[string]struct{}
	technician       map[string]struct{}
	project          map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (woc *WorkOrderCreate) SetCreateTime(t time.Time) *WorkOrderCreate {
	woc.create_time = &t
	return woc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (woc *WorkOrderCreate) SetNillableCreateTime(t *time.Time) *WorkOrderCreate {
	if t != nil {
		woc.SetCreateTime(*t)
	}
	return woc
}

// SetUpdateTime sets the update_time field.
func (woc *WorkOrderCreate) SetUpdateTime(t time.Time) *WorkOrderCreate {
	woc.update_time = &t
	return woc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (woc *WorkOrderCreate) SetNillableUpdateTime(t *time.Time) *WorkOrderCreate {
	if t != nil {
		woc.SetUpdateTime(*t)
	}
	return woc
}

// SetName sets the name field.
func (woc *WorkOrderCreate) SetName(s string) *WorkOrderCreate {
	woc.name = &s
	return woc
}

// SetStatus sets the status field.
func (woc *WorkOrderCreate) SetStatus(s string) *WorkOrderCreate {
	woc.status = &s
	return woc
}

// SetNillableStatus sets the status field if the given value is not nil.
func (woc *WorkOrderCreate) SetNillableStatus(s *string) *WorkOrderCreate {
	if s != nil {
		woc.SetStatus(*s)
	}
	return woc
}

// SetPriority sets the priority field.
func (woc *WorkOrderCreate) SetPriority(s string) *WorkOrderCreate {
	woc.priority = &s
	return woc
}

// SetNillablePriority sets the priority field if the given value is not nil.
func (woc *WorkOrderCreate) SetNillablePriority(s *string) *WorkOrderCreate {
	if s != nil {
		woc.SetPriority(*s)
	}
	return woc
}

// SetDescription sets the description field.
func (woc *WorkOrderCreate) SetDescription(s string) *WorkOrderCreate {
	woc.description = &s
	return woc
}

// SetNillableDescription sets the description field if the given value is not nil.
func (woc *WorkOrderCreate) SetNillableDescription(s *string) *WorkOrderCreate {
	if s != nil {
		woc.SetDescription(*s)
	}
	return woc
}

// SetOwnerName sets the owner_name field.
func (woc *WorkOrderCreate) SetOwnerName(s string) *WorkOrderCreate {
	woc.owner_name = &s
	return woc
}

// SetInstallDate sets the install_date field.
func (woc *WorkOrderCreate) SetInstallDate(t time.Time) *WorkOrderCreate {
	woc.install_date = &t
	return woc
}

// SetNillableInstallDate sets the install_date field if the given value is not nil.
func (woc *WorkOrderCreate) SetNillableInstallDate(t *time.Time) *WorkOrderCreate {
	if t != nil {
		woc.SetInstallDate(*t)
	}
	return woc
}

// SetCreationDate sets the creation_date field.
func (woc *WorkOrderCreate) SetCreationDate(t time.Time) *WorkOrderCreate {
	woc.creation_date = &t
	return woc
}

// SetAssignee sets the assignee field.
func (woc *WorkOrderCreate) SetAssignee(s string) *WorkOrderCreate {
	woc.assignee = &s
	return woc
}

// SetNillableAssignee sets the assignee field if the given value is not nil.
func (woc *WorkOrderCreate) SetNillableAssignee(s *string) *WorkOrderCreate {
	if s != nil {
		woc.SetAssignee(*s)
	}
	return woc
}

// SetIndex sets the index field.
func (woc *WorkOrderCreate) SetIndex(i int) *WorkOrderCreate {
	woc.index = &i
	return woc
}

// SetNillableIndex sets the index field if the given value is not nil.
func (woc *WorkOrderCreate) SetNillableIndex(i *int) *WorkOrderCreate {
	if i != nil {
		woc.SetIndex(*i)
	}
	return woc
}

// SetTypeID sets the type edge to WorkOrderType by id.
func (woc *WorkOrderCreate) SetTypeID(id string) *WorkOrderCreate {
	if woc._type == nil {
		woc._type = make(map[string]struct{})
	}
	woc._type[id] = struct{}{}
	return woc
}

// SetNillableTypeID sets the type edge to WorkOrderType by id if the given value is not nil.
func (woc *WorkOrderCreate) SetNillableTypeID(id *string) *WorkOrderCreate {
	if id != nil {
		woc = woc.SetTypeID(*id)
	}
	return woc
}

// SetType sets the type edge to WorkOrderType.
func (woc *WorkOrderCreate) SetType(w *WorkOrderType) *WorkOrderCreate {
	return woc.SetTypeID(w.ID)
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (woc *WorkOrderCreate) AddEquipmentIDs(ids ...string) *WorkOrderCreate {
	if woc.equipment == nil {
		woc.equipment = make(map[string]struct{})
	}
	for i := range ids {
		woc.equipment[ids[i]] = struct{}{}
	}
	return woc
}

// AddEquipment adds the equipment edges to Equipment.
func (woc *WorkOrderCreate) AddEquipment(e ...*Equipment) *WorkOrderCreate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return woc.AddEquipmentIDs(ids...)
}

// AddLinkIDs adds the links edge to Link by ids.
func (woc *WorkOrderCreate) AddLinkIDs(ids ...string) *WorkOrderCreate {
	if woc.links == nil {
		woc.links = make(map[string]struct{})
	}
	for i := range ids {
		woc.links[ids[i]] = struct{}{}
	}
	return woc
}

// AddLinks adds the links edges to Link.
func (woc *WorkOrderCreate) AddLinks(l ...*Link) *WorkOrderCreate {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return woc.AddLinkIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (woc *WorkOrderCreate) AddFileIDs(ids ...string) *WorkOrderCreate {
	if woc.files == nil {
		woc.files = make(map[string]struct{})
	}
	for i := range ids {
		woc.files[ids[i]] = struct{}{}
	}
	return woc
}

// AddFiles adds the files edges to File.
func (woc *WorkOrderCreate) AddFiles(f ...*File) *WorkOrderCreate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return woc.AddFileIDs(ids...)
}

// SetLocationID sets the location edge to Location by id.
func (woc *WorkOrderCreate) SetLocationID(id string) *WorkOrderCreate {
	if woc.location == nil {
		woc.location = make(map[string]struct{})
	}
	woc.location[id] = struct{}{}
	return woc
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (woc *WorkOrderCreate) SetNillableLocationID(id *string) *WorkOrderCreate {
	if id != nil {
		woc = woc.SetLocationID(*id)
	}
	return woc
}

// SetLocation sets the location edge to Location.
func (woc *WorkOrderCreate) SetLocation(l *Location) *WorkOrderCreate {
	return woc.SetLocationID(l.ID)
}

// AddCommentIDs adds the comments edge to Comment by ids.
func (woc *WorkOrderCreate) AddCommentIDs(ids ...string) *WorkOrderCreate {
	if woc.comments == nil {
		woc.comments = make(map[string]struct{})
	}
	for i := range ids {
		woc.comments[ids[i]] = struct{}{}
	}
	return woc
}

// AddComments adds the comments edges to Comment.
func (woc *WorkOrderCreate) AddComments(c ...*Comment) *WorkOrderCreate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return woc.AddCommentIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (woc *WorkOrderCreate) AddPropertyIDs(ids ...string) *WorkOrderCreate {
	if woc.properties == nil {
		woc.properties = make(map[string]struct{})
	}
	for i := range ids {
		woc.properties[ids[i]] = struct{}{}
	}
	return woc
}

// AddProperties adds the properties edges to Property.
func (woc *WorkOrderCreate) AddProperties(p ...*Property) *WorkOrderCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return woc.AddPropertyIDs(ids...)
}

// AddCheckListItemIDs adds the check_list_items edge to CheckListItem by ids.
func (woc *WorkOrderCreate) AddCheckListItemIDs(ids ...string) *WorkOrderCreate {
	if woc.check_list_items == nil {
		woc.check_list_items = make(map[string]struct{})
	}
	for i := range ids {
		woc.check_list_items[ids[i]] = struct{}{}
	}
	return woc
}

// AddCheckListItems adds the check_list_items edges to CheckListItem.
func (woc *WorkOrderCreate) AddCheckListItems(c ...*CheckListItem) *WorkOrderCreate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return woc.AddCheckListItemIDs(ids...)
}

// SetTechnicianID sets the technician edge to Technician by id.
func (woc *WorkOrderCreate) SetTechnicianID(id string) *WorkOrderCreate {
	if woc.technician == nil {
		woc.technician = make(map[string]struct{})
	}
	woc.technician[id] = struct{}{}
	return woc
}

// SetNillableTechnicianID sets the technician edge to Technician by id if the given value is not nil.
func (woc *WorkOrderCreate) SetNillableTechnicianID(id *string) *WorkOrderCreate {
	if id != nil {
		woc = woc.SetTechnicianID(*id)
	}
	return woc
}

// SetTechnician sets the technician edge to Technician.
func (woc *WorkOrderCreate) SetTechnician(t *Technician) *WorkOrderCreate {
	return woc.SetTechnicianID(t.ID)
}

// SetProjectID sets the project edge to Project by id.
func (woc *WorkOrderCreate) SetProjectID(id string) *WorkOrderCreate {
	if woc.project == nil {
		woc.project = make(map[string]struct{})
	}
	woc.project[id] = struct{}{}
	return woc
}

// SetNillableProjectID sets the project edge to Project by id if the given value is not nil.
func (woc *WorkOrderCreate) SetNillableProjectID(id *string) *WorkOrderCreate {
	if id != nil {
		woc = woc.SetProjectID(*id)
	}
	return woc
}

// SetProject sets the project edge to Project.
func (woc *WorkOrderCreate) SetProject(p *Project) *WorkOrderCreate {
	return woc.SetProjectID(p.ID)
}

// Save creates the WorkOrder in the database.
func (woc *WorkOrderCreate) Save(ctx context.Context) (*WorkOrder, error) {
	if woc.create_time == nil {
		v := workorder.DefaultCreateTime()
		woc.create_time = &v
	}
	if woc.update_time == nil {
		v := workorder.DefaultUpdateTime()
		woc.update_time = &v
	}
	if woc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if err := workorder.NameValidator(*woc.name); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
	}
	if woc.status == nil {
		v := workorder.DefaultStatus
		woc.status = &v
	}
	if woc.priority == nil {
		v := workorder.DefaultPriority
		woc.priority = &v
	}
	if woc.owner_name == nil {
		return nil, errors.New("ent: missing required field \"owner_name\"")
	}
	if woc.creation_date == nil {
		return nil, errors.New("ent: missing required field \"creation_date\"")
	}
	if len(woc._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if len(woc.location) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	if len(woc.technician) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"technician\"")
	}
	if len(woc.project) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"project\"")
	}
	return woc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (woc *WorkOrderCreate) SaveX(ctx context.Context) *WorkOrder {
	v, err := woc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (woc *WorkOrderCreate) sqlSave(ctx context.Context) (*WorkOrder, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(woc.driver.Dialect())
		wo      = &WorkOrder{config: woc.config}
	)
	tx, err := woc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(workorder.Table).Default()
	if value := woc.create_time; value != nil {
		insert.Set(workorder.FieldCreateTime, *value)
		wo.CreateTime = *value
	}
	if value := woc.update_time; value != nil {
		insert.Set(workorder.FieldUpdateTime, *value)
		wo.UpdateTime = *value
	}
	if value := woc.name; value != nil {
		insert.Set(workorder.FieldName, *value)
		wo.Name = *value
	}
	if value := woc.status; value != nil {
		insert.Set(workorder.FieldStatus, *value)
		wo.Status = *value
	}
	if value := woc.priority; value != nil {
		insert.Set(workorder.FieldPriority, *value)
		wo.Priority = *value
	}
	if value := woc.description; value != nil {
		insert.Set(workorder.FieldDescription, *value)
		wo.Description = *value
	}
	if value := woc.owner_name; value != nil {
		insert.Set(workorder.FieldOwnerName, *value)
		wo.OwnerName = *value
	}
	if value := woc.install_date; value != nil {
		insert.Set(workorder.FieldInstallDate, *value)
		wo.InstallDate = *value
	}
	if value := woc.creation_date; value != nil {
		insert.Set(workorder.FieldCreationDate, *value)
		wo.CreationDate = *value
	}
	if value := woc.assignee; value != nil {
		insert.Set(workorder.FieldAssignee, *value)
		wo.Assignee = *value
	}
	if value := woc.index; value != nil {
		insert.Set(workorder.FieldIndex, *value)
		wo.Index = *value
	}
	id, err := insertLastID(ctx, tx, insert.Returning(workorder.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	wo.ID = strconv.FormatInt(id, 10)
	if len(woc._type) > 0 {
		for eid := range woc._type {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(workorder.TypeTable).
				Set(workorder.TypeColumn, eid).
				Where(sql.EQ(workorder.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(woc.equipment) > 0 {
		p := sql.P()
		for eid := range woc.equipment {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(equipment.FieldID, eid)
		}
		query, args := builder.Update(workorder.EquipmentTable).
			Set(workorder.EquipmentColumn, id).
			Where(sql.And(p, sql.IsNull(workorder.EquipmentColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(woc.equipment) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"equipment\" %v already connected to a different \"WorkOrder\"", keys(woc.equipment))})
		}
	}
	if len(woc.links) > 0 {
		p := sql.P()
		for eid := range woc.links {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(link.FieldID, eid)
		}
		query, args := builder.Update(workorder.LinksTable).
			Set(workorder.LinksColumn, id).
			Where(sql.And(p, sql.IsNull(workorder.LinksColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(woc.links) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"links\" %v already connected to a different \"WorkOrder\"", keys(woc.links))})
		}
	}
	if len(woc.files) > 0 {
		p := sql.P()
		for eid := range woc.files {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(file.FieldID, eid)
		}
		query, args := builder.Update(workorder.FilesTable).
			Set(workorder.FilesColumn, id).
			Where(sql.And(p, sql.IsNull(workorder.FilesColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(woc.files) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"files\" %v already connected to a different \"WorkOrder\"", keys(woc.files))})
		}
	}
	if len(woc.location) > 0 {
		for eid := range woc.location {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(workorder.LocationTable).
				Set(workorder.LocationColumn, eid).
				Where(sql.EQ(workorder.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(woc.comments) > 0 {
		p := sql.P()
		for eid := range woc.comments {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(comment.FieldID, eid)
		}
		query, args := builder.Update(workorder.CommentsTable).
			Set(workorder.CommentsColumn, id).
			Where(sql.And(p, sql.IsNull(workorder.CommentsColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(woc.comments) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"comments\" %v already connected to a different \"WorkOrder\"", keys(woc.comments))})
		}
	}
	if len(woc.properties) > 0 {
		p := sql.P()
		for eid := range woc.properties {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(property.FieldID, eid)
		}
		query, args := builder.Update(workorder.PropertiesTable).
			Set(workorder.PropertiesColumn, id).
			Where(sql.And(p, sql.IsNull(workorder.PropertiesColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(woc.properties) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"WorkOrder\"", keys(woc.properties))})
		}
	}
	if len(woc.check_list_items) > 0 {
		p := sql.P()
		for eid := range woc.check_list_items {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(checklistitem.FieldID, eid)
		}
		query, args := builder.Update(workorder.CheckListItemsTable).
			Set(workorder.CheckListItemsColumn, id).
			Where(sql.And(p, sql.IsNull(workorder.CheckListItemsColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(woc.check_list_items) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"check_list_items\" %v already connected to a different \"WorkOrder\"", keys(woc.check_list_items))})
		}
	}
	if len(woc.technician) > 0 {
		for eid := range woc.technician {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(workorder.TechnicianTable).
				Set(workorder.TechnicianColumn, eid).
				Where(sql.EQ(workorder.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(woc.project) > 0 {
		for eid := range woc.project {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(workorder.ProjectTable).
				Set(workorder.ProjectColumn, eid).
				Where(sql.EQ(workorder.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return wo, nil
}
