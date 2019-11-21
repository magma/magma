// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/property"
)

// PropertyCreate is the builder for creating a Property entity.
type PropertyCreate struct {
	config
	create_time     *time.Time
	update_time     *time.Time
	int_val         *int
	bool_val        *bool
	float_val       *float64
	latitude_val    *float64
	longitude_val   *float64
	range_from_val  *float64
	range_to_val    *float64
	string_val      *string
	_type           map[string]struct{}
	location        map[string]struct{}
	equipment       map[string]struct{}
	service         map[string]struct{}
	equipment_port  map[string]struct{}
	link            map[string]struct{}
	work_order      map[string]struct{}
	project         map[string]struct{}
	equipment_value map[string]struct{}
	location_value  map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (pc *PropertyCreate) SetCreateTime(t time.Time) *PropertyCreate {
	pc.create_time = &t
	return pc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (pc *PropertyCreate) SetNillableCreateTime(t *time.Time) *PropertyCreate {
	if t != nil {
		pc.SetCreateTime(*t)
	}
	return pc
}

// SetUpdateTime sets the update_time field.
func (pc *PropertyCreate) SetUpdateTime(t time.Time) *PropertyCreate {
	pc.update_time = &t
	return pc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (pc *PropertyCreate) SetNillableUpdateTime(t *time.Time) *PropertyCreate {
	if t != nil {
		pc.SetUpdateTime(*t)
	}
	return pc
}

// SetIntVal sets the int_val field.
func (pc *PropertyCreate) SetIntVal(i int) *PropertyCreate {
	pc.int_val = &i
	return pc
}

// SetNillableIntVal sets the int_val field if the given value is not nil.
func (pc *PropertyCreate) SetNillableIntVal(i *int) *PropertyCreate {
	if i != nil {
		pc.SetIntVal(*i)
	}
	return pc
}

// SetBoolVal sets the bool_val field.
func (pc *PropertyCreate) SetBoolVal(b bool) *PropertyCreate {
	pc.bool_val = &b
	return pc
}

// SetNillableBoolVal sets the bool_val field if the given value is not nil.
func (pc *PropertyCreate) SetNillableBoolVal(b *bool) *PropertyCreate {
	if b != nil {
		pc.SetBoolVal(*b)
	}
	return pc
}

// SetFloatVal sets the float_val field.
func (pc *PropertyCreate) SetFloatVal(f float64) *PropertyCreate {
	pc.float_val = &f
	return pc
}

// SetNillableFloatVal sets the float_val field if the given value is not nil.
func (pc *PropertyCreate) SetNillableFloatVal(f *float64) *PropertyCreate {
	if f != nil {
		pc.SetFloatVal(*f)
	}
	return pc
}

// SetLatitudeVal sets the latitude_val field.
func (pc *PropertyCreate) SetLatitudeVal(f float64) *PropertyCreate {
	pc.latitude_val = &f
	return pc
}

// SetNillableLatitudeVal sets the latitude_val field if the given value is not nil.
func (pc *PropertyCreate) SetNillableLatitudeVal(f *float64) *PropertyCreate {
	if f != nil {
		pc.SetLatitudeVal(*f)
	}
	return pc
}

// SetLongitudeVal sets the longitude_val field.
func (pc *PropertyCreate) SetLongitudeVal(f float64) *PropertyCreate {
	pc.longitude_val = &f
	return pc
}

// SetNillableLongitudeVal sets the longitude_val field if the given value is not nil.
func (pc *PropertyCreate) SetNillableLongitudeVal(f *float64) *PropertyCreate {
	if f != nil {
		pc.SetLongitudeVal(*f)
	}
	return pc
}

// SetRangeFromVal sets the range_from_val field.
func (pc *PropertyCreate) SetRangeFromVal(f float64) *PropertyCreate {
	pc.range_from_val = &f
	return pc
}

// SetNillableRangeFromVal sets the range_from_val field if the given value is not nil.
func (pc *PropertyCreate) SetNillableRangeFromVal(f *float64) *PropertyCreate {
	if f != nil {
		pc.SetRangeFromVal(*f)
	}
	return pc
}

// SetRangeToVal sets the range_to_val field.
func (pc *PropertyCreate) SetRangeToVal(f float64) *PropertyCreate {
	pc.range_to_val = &f
	return pc
}

// SetNillableRangeToVal sets the range_to_val field if the given value is not nil.
func (pc *PropertyCreate) SetNillableRangeToVal(f *float64) *PropertyCreate {
	if f != nil {
		pc.SetRangeToVal(*f)
	}
	return pc
}

// SetStringVal sets the string_val field.
func (pc *PropertyCreate) SetStringVal(s string) *PropertyCreate {
	pc.string_val = &s
	return pc
}

// SetNillableStringVal sets the string_val field if the given value is not nil.
func (pc *PropertyCreate) SetNillableStringVal(s *string) *PropertyCreate {
	if s != nil {
		pc.SetStringVal(*s)
	}
	return pc
}

// SetTypeID sets the type edge to PropertyType by id.
func (pc *PropertyCreate) SetTypeID(id string) *PropertyCreate {
	if pc._type == nil {
		pc._type = make(map[string]struct{})
	}
	pc._type[id] = struct{}{}
	return pc
}

// SetType sets the type edge to PropertyType.
func (pc *PropertyCreate) SetType(p *PropertyType) *PropertyCreate {
	return pc.SetTypeID(p.ID)
}

// SetLocationID sets the location edge to Location by id.
func (pc *PropertyCreate) SetLocationID(id string) *PropertyCreate {
	if pc.location == nil {
		pc.location = make(map[string]struct{})
	}
	pc.location[id] = struct{}{}
	return pc
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableLocationID(id *string) *PropertyCreate {
	if id != nil {
		pc = pc.SetLocationID(*id)
	}
	return pc
}

// SetLocation sets the location edge to Location.
func (pc *PropertyCreate) SetLocation(l *Location) *PropertyCreate {
	return pc.SetLocationID(l.ID)
}

// SetEquipmentID sets the equipment edge to Equipment by id.
func (pc *PropertyCreate) SetEquipmentID(id string) *PropertyCreate {
	if pc.equipment == nil {
		pc.equipment = make(map[string]struct{})
	}
	pc.equipment[id] = struct{}{}
	return pc
}

// SetNillableEquipmentID sets the equipment edge to Equipment by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableEquipmentID(id *string) *PropertyCreate {
	if id != nil {
		pc = pc.SetEquipmentID(*id)
	}
	return pc
}

// SetEquipment sets the equipment edge to Equipment.
func (pc *PropertyCreate) SetEquipment(e *Equipment) *PropertyCreate {
	return pc.SetEquipmentID(e.ID)
}

// SetServiceID sets the service edge to Service by id.
func (pc *PropertyCreate) SetServiceID(id string) *PropertyCreate {
	if pc.service == nil {
		pc.service = make(map[string]struct{})
	}
	pc.service[id] = struct{}{}
	return pc
}

// SetNillableServiceID sets the service edge to Service by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableServiceID(id *string) *PropertyCreate {
	if id != nil {
		pc = pc.SetServiceID(*id)
	}
	return pc
}

// SetService sets the service edge to Service.
func (pc *PropertyCreate) SetService(s *Service) *PropertyCreate {
	return pc.SetServiceID(s.ID)
}

// SetEquipmentPortID sets the equipment_port edge to EquipmentPort by id.
func (pc *PropertyCreate) SetEquipmentPortID(id string) *PropertyCreate {
	if pc.equipment_port == nil {
		pc.equipment_port = make(map[string]struct{})
	}
	pc.equipment_port[id] = struct{}{}
	return pc
}

// SetNillableEquipmentPortID sets the equipment_port edge to EquipmentPort by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableEquipmentPortID(id *string) *PropertyCreate {
	if id != nil {
		pc = pc.SetEquipmentPortID(*id)
	}
	return pc
}

// SetEquipmentPort sets the equipment_port edge to EquipmentPort.
func (pc *PropertyCreate) SetEquipmentPort(e *EquipmentPort) *PropertyCreate {
	return pc.SetEquipmentPortID(e.ID)
}

// SetLinkID sets the link edge to Link by id.
func (pc *PropertyCreate) SetLinkID(id string) *PropertyCreate {
	if pc.link == nil {
		pc.link = make(map[string]struct{})
	}
	pc.link[id] = struct{}{}
	return pc
}

// SetNillableLinkID sets the link edge to Link by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableLinkID(id *string) *PropertyCreate {
	if id != nil {
		pc = pc.SetLinkID(*id)
	}
	return pc
}

// SetLink sets the link edge to Link.
func (pc *PropertyCreate) SetLink(l *Link) *PropertyCreate {
	return pc.SetLinkID(l.ID)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (pc *PropertyCreate) SetWorkOrderID(id string) *PropertyCreate {
	if pc.work_order == nil {
		pc.work_order = make(map[string]struct{})
	}
	pc.work_order[id] = struct{}{}
	return pc
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableWorkOrderID(id *string) *PropertyCreate {
	if id != nil {
		pc = pc.SetWorkOrderID(*id)
	}
	return pc
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (pc *PropertyCreate) SetWorkOrder(w *WorkOrder) *PropertyCreate {
	return pc.SetWorkOrderID(w.ID)
}

// SetProjectID sets the project edge to Project by id.
func (pc *PropertyCreate) SetProjectID(id string) *PropertyCreate {
	if pc.project == nil {
		pc.project = make(map[string]struct{})
	}
	pc.project[id] = struct{}{}
	return pc
}

// SetNillableProjectID sets the project edge to Project by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableProjectID(id *string) *PropertyCreate {
	if id != nil {
		pc = pc.SetProjectID(*id)
	}
	return pc
}

// SetProject sets the project edge to Project.
func (pc *PropertyCreate) SetProject(p *Project) *PropertyCreate {
	return pc.SetProjectID(p.ID)
}

// SetEquipmentValueID sets the equipment_value edge to Equipment by id.
func (pc *PropertyCreate) SetEquipmentValueID(id string) *PropertyCreate {
	if pc.equipment_value == nil {
		pc.equipment_value = make(map[string]struct{})
	}
	pc.equipment_value[id] = struct{}{}
	return pc
}

// SetNillableEquipmentValueID sets the equipment_value edge to Equipment by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableEquipmentValueID(id *string) *PropertyCreate {
	if id != nil {
		pc = pc.SetEquipmentValueID(*id)
	}
	return pc
}

// SetEquipmentValue sets the equipment_value edge to Equipment.
func (pc *PropertyCreate) SetEquipmentValue(e *Equipment) *PropertyCreate {
	return pc.SetEquipmentValueID(e.ID)
}

// SetLocationValueID sets the location_value edge to Location by id.
func (pc *PropertyCreate) SetLocationValueID(id string) *PropertyCreate {
	if pc.location_value == nil {
		pc.location_value = make(map[string]struct{})
	}
	pc.location_value[id] = struct{}{}
	return pc
}

// SetNillableLocationValueID sets the location_value edge to Location by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableLocationValueID(id *string) *PropertyCreate {
	if id != nil {
		pc = pc.SetLocationValueID(*id)
	}
	return pc
}

// SetLocationValue sets the location_value edge to Location.
func (pc *PropertyCreate) SetLocationValue(l *Location) *PropertyCreate {
	return pc.SetLocationValueID(l.ID)
}

// Save creates the Property in the database.
func (pc *PropertyCreate) Save(ctx context.Context) (*Property, error) {
	if pc.create_time == nil {
		v := property.DefaultCreateTime()
		pc.create_time = &v
	}
	if pc.update_time == nil {
		v := property.DefaultUpdateTime()
		pc.update_time = &v
	}
	if len(pc._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if pc._type == nil {
		return nil, errors.New("ent: missing required edge \"type\"")
	}
	if len(pc.location) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	if len(pc.equipment) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"equipment\"")
	}
	if len(pc.service) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"service\"")
	}
	if len(pc.equipment_port) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"equipment_port\"")
	}
	if len(pc.link) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"link\"")
	}
	if len(pc.work_order) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"work_order\"")
	}
	if len(pc.project) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"project\"")
	}
	if len(pc.equipment_value) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"equipment_value\"")
	}
	if len(pc.location_value) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location_value\"")
	}
	return pc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (pc *PropertyCreate) SaveX(ctx context.Context) *Property {
	v, err := pc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (pc *PropertyCreate) sqlSave(ctx context.Context) (*Property, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(pc.driver.Dialect())
		pr      = &Property{config: pc.config}
	)
	tx, err := pc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(property.Table).Default()
	if value := pc.create_time; value != nil {
		insert.Set(property.FieldCreateTime, *value)
		pr.CreateTime = *value
	}
	if value := pc.update_time; value != nil {
		insert.Set(property.FieldUpdateTime, *value)
		pr.UpdateTime = *value
	}
	if value := pc.int_val; value != nil {
		insert.Set(property.FieldIntVal, *value)
		pr.IntVal = *value
	}
	if value := pc.bool_val; value != nil {
		insert.Set(property.FieldBoolVal, *value)
		pr.BoolVal = *value
	}
	if value := pc.float_val; value != nil {
		insert.Set(property.FieldFloatVal, *value)
		pr.FloatVal = *value
	}
	if value := pc.latitude_val; value != nil {
		insert.Set(property.FieldLatitudeVal, *value)
		pr.LatitudeVal = *value
	}
	if value := pc.longitude_val; value != nil {
		insert.Set(property.FieldLongitudeVal, *value)
		pr.LongitudeVal = *value
	}
	if value := pc.range_from_val; value != nil {
		insert.Set(property.FieldRangeFromVal, *value)
		pr.RangeFromVal = *value
	}
	if value := pc.range_to_val; value != nil {
		insert.Set(property.FieldRangeToVal, *value)
		pr.RangeToVal = *value
	}
	if value := pc.string_val; value != nil {
		insert.Set(property.FieldStringVal, *value)
		pr.StringVal = *value
	}
	id, err := insertLastID(ctx, tx, insert.Returning(property.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	pr.ID = strconv.FormatInt(id, 10)
	if len(pc._type) > 0 {
		for eid := range pc._type {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(property.TypeTable).
				Set(property.TypeColumn, eid).
				Where(sql.EQ(property.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(pc.location) > 0 {
		for eid := range pc.location {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(property.LocationTable).
				Set(property.LocationColumn, eid).
				Where(sql.EQ(property.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(pc.equipment) > 0 {
		for eid := range pc.equipment {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(property.EquipmentTable).
				Set(property.EquipmentColumn, eid).
				Where(sql.EQ(property.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(pc.service) > 0 {
		for eid := range pc.service {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(property.ServiceTable).
				Set(property.ServiceColumn, eid).
				Where(sql.EQ(property.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(pc.equipment_port) > 0 {
		for eid := range pc.equipment_port {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(property.EquipmentPortTable).
				Set(property.EquipmentPortColumn, eid).
				Where(sql.EQ(property.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(pc.link) > 0 {
		for eid := range pc.link {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(property.LinkTable).
				Set(property.LinkColumn, eid).
				Where(sql.EQ(property.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(pc.work_order) > 0 {
		for eid := range pc.work_order {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(property.WorkOrderTable).
				Set(property.WorkOrderColumn, eid).
				Where(sql.EQ(property.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(pc.project) > 0 {
		for eid := range pc.project {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(property.ProjectTable).
				Set(property.ProjectColumn, eid).
				Where(sql.EQ(property.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(pc.equipment_value) > 0 {
		for eid := range pc.equipment_value {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(property.EquipmentValueTable).
				Set(property.EquipmentValueColumn, eid).
				Where(sql.EQ(property.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(pc.location_value) > 0 {
		for eid := range pc.location_value {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(property.LocationValueTable).
				Set(property.LocationValueColumn, eid).
				Where(sql.EQ(property.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return pr, nil
}
