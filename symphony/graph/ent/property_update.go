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
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// PropertyUpdate is the builder for updating Property entities.
type PropertyUpdate struct {
	config

	update_time           *time.Time
	int_val               *int
	addint_val            *int
	clearint_val          bool
	bool_val              *bool
	clearbool_val         bool
	float_val             *float64
	addfloat_val          *float64
	clearfloat_val        bool
	latitude_val          *float64
	addlatitude_val       *float64
	clearlatitude_val     bool
	longitude_val         *float64
	addlongitude_val      *float64
	clearlongitude_val    bool
	range_from_val        *float64
	addrange_from_val     *float64
	clearrange_from_val   bool
	range_to_val          *float64
	addrange_to_val       *float64
	clearrange_to_val     bool
	string_val            *string
	clearstring_val       bool
	_type                 map[string]struct{}
	location              map[string]struct{}
	equipment             map[string]struct{}
	service               map[string]struct{}
	equipment_port        map[string]struct{}
	link                  map[string]struct{}
	work_order            map[string]struct{}
	project               map[string]struct{}
	equipment_value       map[string]struct{}
	location_value        map[string]struct{}
	service_value         map[string]struct{}
	clearedType           bool
	clearedLocation       bool
	clearedEquipment      bool
	clearedService        bool
	clearedEquipmentPort  bool
	clearedLink           bool
	clearedWorkOrder      bool
	clearedProject        bool
	clearedEquipmentValue bool
	clearedLocationValue  bool
	clearedServiceValue   bool
	predicates            []predicate.Property
}

// Where adds a new predicate for the builder.
func (pu *PropertyUpdate) Where(ps ...predicate.Property) *PropertyUpdate {
	pu.predicates = append(pu.predicates, ps...)
	return pu
}

// SetIntVal sets the int_val field.
func (pu *PropertyUpdate) SetIntVal(i int) *PropertyUpdate {
	pu.int_val = &i
	pu.addint_val = nil
	return pu
}

// SetNillableIntVal sets the int_val field if the given value is not nil.
func (pu *PropertyUpdate) SetNillableIntVal(i *int) *PropertyUpdate {
	if i != nil {
		pu.SetIntVal(*i)
	}
	return pu
}

// AddIntVal adds i to int_val.
func (pu *PropertyUpdate) AddIntVal(i int) *PropertyUpdate {
	if pu.addint_val == nil {
		pu.addint_val = &i
	} else {
		*pu.addint_val += i
	}
	return pu
}

// ClearIntVal clears the value of int_val.
func (pu *PropertyUpdate) ClearIntVal() *PropertyUpdate {
	pu.int_val = nil
	pu.clearint_val = true
	return pu
}

// SetBoolVal sets the bool_val field.
func (pu *PropertyUpdate) SetBoolVal(b bool) *PropertyUpdate {
	pu.bool_val = &b
	return pu
}

// SetNillableBoolVal sets the bool_val field if the given value is not nil.
func (pu *PropertyUpdate) SetNillableBoolVal(b *bool) *PropertyUpdate {
	if b != nil {
		pu.SetBoolVal(*b)
	}
	return pu
}

// ClearBoolVal clears the value of bool_val.
func (pu *PropertyUpdate) ClearBoolVal() *PropertyUpdate {
	pu.bool_val = nil
	pu.clearbool_val = true
	return pu
}

// SetFloatVal sets the float_val field.
func (pu *PropertyUpdate) SetFloatVal(f float64) *PropertyUpdate {
	pu.float_val = &f
	pu.addfloat_val = nil
	return pu
}

// SetNillableFloatVal sets the float_val field if the given value is not nil.
func (pu *PropertyUpdate) SetNillableFloatVal(f *float64) *PropertyUpdate {
	if f != nil {
		pu.SetFloatVal(*f)
	}
	return pu
}

// AddFloatVal adds f to float_val.
func (pu *PropertyUpdate) AddFloatVal(f float64) *PropertyUpdate {
	if pu.addfloat_val == nil {
		pu.addfloat_val = &f
	} else {
		*pu.addfloat_val += f
	}
	return pu
}

// ClearFloatVal clears the value of float_val.
func (pu *PropertyUpdate) ClearFloatVal() *PropertyUpdate {
	pu.float_val = nil
	pu.clearfloat_val = true
	return pu
}

// SetLatitudeVal sets the latitude_val field.
func (pu *PropertyUpdate) SetLatitudeVal(f float64) *PropertyUpdate {
	pu.latitude_val = &f
	pu.addlatitude_val = nil
	return pu
}

// SetNillableLatitudeVal sets the latitude_val field if the given value is not nil.
func (pu *PropertyUpdate) SetNillableLatitudeVal(f *float64) *PropertyUpdate {
	if f != nil {
		pu.SetLatitudeVal(*f)
	}
	return pu
}

// AddLatitudeVal adds f to latitude_val.
func (pu *PropertyUpdate) AddLatitudeVal(f float64) *PropertyUpdate {
	if pu.addlatitude_val == nil {
		pu.addlatitude_val = &f
	} else {
		*pu.addlatitude_val += f
	}
	return pu
}

// ClearLatitudeVal clears the value of latitude_val.
func (pu *PropertyUpdate) ClearLatitudeVal() *PropertyUpdate {
	pu.latitude_val = nil
	pu.clearlatitude_val = true
	return pu
}

// SetLongitudeVal sets the longitude_val field.
func (pu *PropertyUpdate) SetLongitudeVal(f float64) *PropertyUpdate {
	pu.longitude_val = &f
	pu.addlongitude_val = nil
	return pu
}

// SetNillableLongitudeVal sets the longitude_val field if the given value is not nil.
func (pu *PropertyUpdate) SetNillableLongitudeVal(f *float64) *PropertyUpdate {
	if f != nil {
		pu.SetLongitudeVal(*f)
	}
	return pu
}

// AddLongitudeVal adds f to longitude_val.
func (pu *PropertyUpdate) AddLongitudeVal(f float64) *PropertyUpdate {
	if pu.addlongitude_val == nil {
		pu.addlongitude_val = &f
	} else {
		*pu.addlongitude_val += f
	}
	return pu
}

// ClearLongitudeVal clears the value of longitude_val.
func (pu *PropertyUpdate) ClearLongitudeVal() *PropertyUpdate {
	pu.longitude_val = nil
	pu.clearlongitude_val = true
	return pu
}

// SetRangeFromVal sets the range_from_val field.
func (pu *PropertyUpdate) SetRangeFromVal(f float64) *PropertyUpdate {
	pu.range_from_val = &f
	pu.addrange_from_val = nil
	return pu
}

// SetNillableRangeFromVal sets the range_from_val field if the given value is not nil.
func (pu *PropertyUpdate) SetNillableRangeFromVal(f *float64) *PropertyUpdate {
	if f != nil {
		pu.SetRangeFromVal(*f)
	}
	return pu
}

// AddRangeFromVal adds f to range_from_val.
func (pu *PropertyUpdate) AddRangeFromVal(f float64) *PropertyUpdate {
	if pu.addrange_from_val == nil {
		pu.addrange_from_val = &f
	} else {
		*pu.addrange_from_val += f
	}
	return pu
}

// ClearRangeFromVal clears the value of range_from_val.
func (pu *PropertyUpdate) ClearRangeFromVal() *PropertyUpdate {
	pu.range_from_val = nil
	pu.clearrange_from_val = true
	return pu
}

// SetRangeToVal sets the range_to_val field.
func (pu *PropertyUpdate) SetRangeToVal(f float64) *PropertyUpdate {
	pu.range_to_val = &f
	pu.addrange_to_val = nil
	return pu
}

// SetNillableRangeToVal sets the range_to_val field if the given value is not nil.
func (pu *PropertyUpdate) SetNillableRangeToVal(f *float64) *PropertyUpdate {
	if f != nil {
		pu.SetRangeToVal(*f)
	}
	return pu
}

// AddRangeToVal adds f to range_to_val.
func (pu *PropertyUpdate) AddRangeToVal(f float64) *PropertyUpdate {
	if pu.addrange_to_val == nil {
		pu.addrange_to_val = &f
	} else {
		*pu.addrange_to_val += f
	}
	return pu
}

// ClearRangeToVal clears the value of range_to_val.
func (pu *PropertyUpdate) ClearRangeToVal() *PropertyUpdate {
	pu.range_to_val = nil
	pu.clearrange_to_val = true
	return pu
}

// SetStringVal sets the string_val field.
func (pu *PropertyUpdate) SetStringVal(s string) *PropertyUpdate {
	pu.string_val = &s
	return pu
}

// SetNillableStringVal sets the string_val field if the given value is not nil.
func (pu *PropertyUpdate) SetNillableStringVal(s *string) *PropertyUpdate {
	if s != nil {
		pu.SetStringVal(*s)
	}
	return pu
}

// ClearStringVal clears the value of string_val.
func (pu *PropertyUpdate) ClearStringVal() *PropertyUpdate {
	pu.string_val = nil
	pu.clearstring_val = true
	return pu
}

// SetTypeID sets the type edge to PropertyType by id.
func (pu *PropertyUpdate) SetTypeID(id string) *PropertyUpdate {
	if pu._type == nil {
		pu._type = make(map[string]struct{})
	}
	pu._type[id] = struct{}{}
	return pu
}

// SetType sets the type edge to PropertyType.
func (pu *PropertyUpdate) SetType(p *PropertyType) *PropertyUpdate {
	return pu.SetTypeID(p.ID)
}

// SetLocationID sets the location edge to Location by id.
func (pu *PropertyUpdate) SetLocationID(id string) *PropertyUpdate {
	if pu.location == nil {
		pu.location = make(map[string]struct{})
	}
	pu.location[id] = struct{}{}
	return pu
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (pu *PropertyUpdate) SetNillableLocationID(id *string) *PropertyUpdate {
	if id != nil {
		pu = pu.SetLocationID(*id)
	}
	return pu
}

// SetLocation sets the location edge to Location.
func (pu *PropertyUpdate) SetLocation(l *Location) *PropertyUpdate {
	return pu.SetLocationID(l.ID)
}

// SetEquipmentID sets the equipment edge to Equipment by id.
func (pu *PropertyUpdate) SetEquipmentID(id string) *PropertyUpdate {
	if pu.equipment == nil {
		pu.equipment = make(map[string]struct{})
	}
	pu.equipment[id] = struct{}{}
	return pu
}

// SetNillableEquipmentID sets the equipment edge to Equipment by id if the given value is not nil.
func (pu *PropertyUpdate) SetNillableEquipmentID(id *string) *PropertyUpdate {
	if id != nil {
		pu = pu.SetEquipmentID(*id)
	}
	return pu
}

// SetEquipment sets the equipment edge to Equipment.
func (pu *PropertyUpdate) SetEquipment(e *Equipment) *PropertyUpdate {
	return pu.SetEquipmentID(e.ID)
}

// SetServiceID sets the service edge to Service by id.
func (pu *PropertyUpdate) SetServiceID(id string) *PropertyUpdate {
	if pu.service == nil {
		pu.service = make(map[string]struct{})
	}
	pu.service[id] = struct{}{}
	return pu
}

// SetNillableServiceID sets the service edge to Service by id if the given value is not nil.
func (pu *PropertyUpdate) SetNillableServiceID(id *string) *PropertyUpdate {
	if id != nil {
		pu = pu.SetServiceID(*id)
	}
	return pu
}

// SetService sets the service edge to Service.
func (pu *PropertyUpdate) SetService(s *Service) *PropertyUpdate {
	return pu.SetServiceID(s.ID)
}

// SetEquipmentPortID sets the equipment_port edge to EquipmentPort by id.
func (pu *PropertyUpdate) SetEquipmentPortID(id string) *PropertyUpdate {
	if pu.equipment_port == nil {
		pu.equipment_port = make(map[string]struct{})
	}
	pu.equipment_port[id] = struct{}{}
	return pu
}

// SetNillableEquipmentPortID sets the equipment_port edge to EquipmentPort by id if the given value is not nil.
func (pu *PropertyUpdate) SetNillableEquipmentPortID(id *string) *PropertyUpdate {
	if id != nil {
		pu = pu.SetEquipmentPortID(*id)
	}
	return pu
}

// SetEquipmentPort sets the equipment_port edge to EquipmentPort.
func (pu *PropertyUpdate) SetEquipmentPort(e *EquipmentPort) *PropertyUpdate {
	return pu.SetEquipmentPortID(e.ID)
}

// SetLinkID sets the link edge to Link by id.
func (pu *PropertyUpdate) SetLinkID(id string) *PropertyUpdate {
	if pu.link == nil {
		pu.link = make(map[string]struct{})
	}
	pu.link[id] = struct{}{}
	return pu
}

// SetNillableLinkID sets the link edge to Link by id if the given value is not nil.
func (pu *PropertyUpdate) SetNillableLinkID(id *string) *PropertyUpdate {
	if id != nil {
		pu = pu.SetLinkID(*id)
	}
	return pu
}

// SetLink sets the link edge to Link.
func (pu *PropertyUpdate) SetLink(l *Link) *PropertyUpdate {
	return pu.SetLinkID(l.ID)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (pu *PropertyUpdate) SetWorkOrderID(id string) *PropertyUpdate {
	if pu.work_order == nil {
		pu.work_order = make(map[string]struct{})
	}
	pu.work_order[id] = struct{}{}
	return pu
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (pu *PropertyUpdate) SetNillableWorkOrderID(id *string) *PropertyUpdate {
	if id != nil {
		pu = pu.SetWorkOrderID(*id)
	}
	return pu
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (pu *PropertyUpdate) SetWorkOrder(w *WorkOrder) *PropertyUpdate {
	return pu.SetWorkOrderID(w.ID)
}

// SetProjectID sets the project edge to Project by id.
func (pu *PropertyUpdate) SetProjectID(id string) *PropertyUpdate {
	if pu.project == nil {
		pu.project = make(map[string]struct{})
	}
	pu.project[id] = struct{}{}
	return pu
}

// SetNillableProjectID sets the project edge to Project by id if the given value is not nil.
func (pu *PropertyUpdate) SetNillableProjectID(id *string) *PropertyUpdate {
	if id != nil {
		pu = pu.SetProjectID(*id)
	}
	return pu
}

// SetProject sets the project edge to Project.
func (pu *PropertyUpdate) SetProject(p *Project) *PropertyUpdate {
	return pu.SetProjectID(p.ID)
}

// SetEquipmentValueID sets the equipment_value edge to Equipment by id.
func (pu *PropertyUpdate) SetEquipmentValueID(id string) *PropertyUpdate {
	if pu.equipment_value == nil {
		pu.equipment_value = make(map[string]struct{})
	}
	pu.equipment_value[id] = struct{}{}
	return pu
}

// SetNillableEquipmentValueID sets the equipment_value edge to Equipment by id if the given value is not nil.
func (pu *PropertyUpdate) SetNillableEquipmentValueID(id *string) *PropertyUpdate {
	if id != nil {
		pu = pu.SetEquipmentValueID(*id)
	}
	return pu
}

// SetEquipmentValue sets the equipment_value edge to Equipment.
func (pu *PropertyUpdate) SetEquipmentValue(e *Equipment) *PropertyUpdate {
	return pu.SetEquipmentValueID(e.ID)
}

// SetLocationValueID sets the location_value edge to Location by id.
func (pu *PropertyUpdate) SetLocationValueID(id string) *PropertyUpdate {
	if pu.location_value == nil {
		pu.location_value = make(map[string]struct{})
	}
	pu.location_value[id] = struct{}{}
	return pu
}

// SetNillableLocationValueID sets the location_value edge to Location by id if the given value is not nil.
func (pu *PropertyUpdate) SetNillableLocationValueID(id *string) *PropertyUpdate {
	if id != nil {
		pu = pu.SetLocationValueID(*id)
	}
	return pu
}

// SetLocationValue sets the location_value edge to Location.
func (pu *PropertyUpdate) SetLocationValue(l *Location) *PropertyUpdate {
	return pu.SetLocationValueID(l.ID)
}

// SetServiceValueID sets the service_value edge to Service by id.
func (pu *PropertyUpdate) SetServiceValueID(id string) *PropertyUpdate {
	if pu.service_value == nil {
		pu.service_value = make(map[string]struct{})
	}
	pu.service_value[id] = struct{}{}
	return pu
}

// SetNillableServiceValueID sets the service_value edge to Service by id if the given value is not nil.
func (pu *PropertyUpdate) SetNillableServiceValueID(id *string) *PropertyUpdate {
	if id != nil {
		pu = pu.SetServiceValueID(*id)
	}
	return pu
}

// SetServiceValue sets the service_value edge to Service.
func (pu *PropertyUpdate) SetServiceValue(s *Service) *PropertyUpdate {
	return pu.SetServiceValueID(s.ID)
}

// ClearType clears the type edge to PropertyType.
func (pu *PropertyUpdate) ClearType() *PropertyUpdate {
	pu.clearedType = true
	return pu
}

// ClearLocation clears the location edge to Location.
func (pu *PropertyUpdate) ClearLocation() *PropertyUpdate {
	pu.clearedLocation = true
	return pu
}

// ClearEquipment clears the equipment edge to Equipment.
func (pu *PropertyUpdate) ClearEquipment() *PropertyUpdate {
	pu.clearedEquipment = true
	return pu
}

// ClearService clears the service edge to Service.
func (pu *PropertyUpdate) ClearService() *PropertyUpdate {
	pu.clearedService = true
	return pu
}

// ClearEquipmentPort clears the equipment_port edge to EquipmentPort.
func (pu *PropertyUpdate) ClearEquipmentPort() *PropertyUpdate {
	pu.clearedEquipmentPort = true
	return pu
}

// ClearLink clears the link edge to Link.
func (pu *PropertyUpdate) ClearLink() *PropertyUpdate {
	pu.clearedLink = true
	return pu
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (pu *PropertyUpdate) ClearWorkOrder() *PropertyUpdate {
	pu.clearedWorkOrder = true
	return pu
}

// ClearProject clears the project edge to Project.
func (pu *PropertyUpdate) ClearProject() *PropertyUpdate {
	pu.clearedProject = true
	return pu
}

// ClearEquipmentValue clears the equipment_value edge to Equipment.
func (pu *PropertyUpdate) ClearEquipmentValue() *PropertyUpdate {
	pu.clearedEquipmentValue = true
	return pu
}

// ClearLocationValue clears the location_value edge to Location.
func (pu *PropertyUpdate) ClearLocationValue() *PropertyUpdate {
	pu.clearedLocationValue = true
	return pu
}

// ClearServiceValue clears the service_value edge to Service.
func (pu *PropertyUpdate) ClearServiceValue() *PropertyUpdate {
	pu.clearedServiceValue = true
	return pu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (pu *PropertyUpdate) Save(ctx context.Context) (int, error) {
	if pu.update_time == nil {
		v := property.UpdateDefaultUpdateTime()
		pu.update_time = &v
	}
	if len(pu._type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if pu.clearedType && pu._type == nil {
		return 0, errors.New("ent: clearing a unique edge \"type\"")
	}
	if len(pu.location) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	if len(pu.equipment) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"equipment\"")
	}
	if len(pu.service) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"service\"")
	}
	if len(pu.equipment_port) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"equipment_port\"")
	}
	if len(pu.link) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"link\"")
	}
	if len(pu.work_order) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"work_order\"")
	}
	if len(pu.project) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"project\"")
	}
	if len(pu.equipment_value) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"equipment_value\"")
	}
	if len(pu.location_value) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"location_value\"")
	}
	if len(pu.service_value) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"service_value\"")
	}
	return pu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (pu *PropertyUpdate) SaveX(ctx context.Context) int {
	affected, err := pu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pu *PropertyUpdate) Exec(ctx context.Context) error {
	_, err := pu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pu *PropertyUpdate) ExecX(ctx context.Context) {
	if err := pu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (pu *PropertyUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   property.Table,
			Columns: property.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: property.FieldID,
			},
		},
	}
	if ps := pu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := pu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: property.FieldUpdateTime,
		})
	}
	if value := pu.int_val; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: property.FieldIntVal,
		})
	}
	if value := pu.addint_val; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: property.FieldIntVal,
		})
	}
	if pu.clearint_val {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: property.FieldIntVal,
		})
	}
	if value := pu.bool_val; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: property.FieldBoolVal,
		})
	}
	if pu.clearbool_val {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: property.FieldBoolVal,
		})
	}
	if value := pu.float_val; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldFloatVal,
		})
	}
	if value := pu.addfloat_val; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldFloatVal,
		})
	}
	if pu.clearfloat_val {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: property.FieldFloatVal,
		})
	}
	if value := pu.latitude_val; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldLatitudeVal,
		})
	}
	if value := pu.addlatitude_val; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldLatitudeVal,
		})
	}
	if pu.clearlatitude_val {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: property.FieldLatitudeVal,
		})
	}
	if value := pu.longitude_val; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldLongitudeVal,
		})
	}
	if value := pu.addlongitude_val; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldLongitudeVal,
		})
	}
	if pu.clearlongitude_val {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: property.FieldLongitudeVal,
		})
	}
	if value := pu.range_from_val; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldRangeFromVal,
		})
	}
	if value := pu.addrange_from_val; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldRangeFromVal,
		})
	}
	if pu.clearrange_from_val {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: property.FieldRangeFromVal,
		})
	}
	if value := pu.range_to_val; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldRangeToVal,
		})
	}
	if value := pu.addrange_to_val; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldRangeToVal,
		})
	}
	if pu.clearrange_to_val {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: property.FieldRangeToVal,
		})
	}
	if value := pu.string_val; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: property.FieldStringVal,
		})
	}
	if pu.clearstring_val {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: property.FieldStringVal,
		})
	}
	if pu.clearedType {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.TypeTable,
			Columns: []string{property.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: propertytype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu._type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.TypeTable,
			Columns: []string{property.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: propertytype.FieldID,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pu.clearedLocation {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.LocationTable,
			Columns: []string{property.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.location; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.LocationTable,
			Columns: []string{property.LocationColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pu.clearedEquipment {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.EquipmentTable,
			Columns: []string{property.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipment.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.equipment; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.EquipmentTable,
			Columns: []string{property.EquipmentColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pu.clearedService {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.ServiceTable,
			Columns: []string{property.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.service; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.ServiceTable,
			Columns: []string{property.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pu.clearedEquipmentPort {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.EquipmentPortTable,
			Columns: []string{property.EquipmentPortColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentport.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.equipment_port; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.EquipmentPortTable,
			Columns: []string{property.EquipmentPortColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentport.FieldID,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pu.clearedLink {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.LinkTable,
			Columns: []string{property.LinkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: link.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.link; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.LinkTable,
			Columns: []string{property.LinkColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pu.clearedWorkOrder {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.WorkOrderTable,
			Columns: []string{property.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.work_order; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.WorkOrderTable,
			Columns: []string{property.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pu.clearedProject {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.ProjectTable,
			Columns: []string{property.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: project.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.project; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.ProjectTable,
			Columns: []string{property.ProjectColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pu.clearedEquipmentValue {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.EquipmentValueTable,
			Columns: []string{property.EquipmentValueColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipment.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.equipment_value; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.EquipmentValueTable,
			Columns: []string{property.EquipmentValueColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pu.clearedLocationValue {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.LocationValueTable,
			Columns: []string{property.LocationValueColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.location_value; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.LocationValueTable,
			Columns: []string{property.LocationValueColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pu.clearedServiceValue {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.ServiceValueTable,
			Columns: []string{property.ServiceValueColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.service_value; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.ServiceValueTable,
			Columns: []string{property.ServiceValueColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, pu.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// PropertyUpdateOne is the builder for updating a single Property entity.
type PropertyUpdateOne struct {
	config
	id string

	update_time           *time.Time
	int_val               *int
	addint_val            *int
	clearint_val          bool
	bool_val              *bool
	clearbool_val         bool
	float_val             *float64
	addfloat_val          *float64
	clearfloat_val        bool
	latitude_val          *float64
	addlatitude_val       *float64
	clearlatitude_val     bool
	longitude_val         *float64
	addlongitude_val      *float64
	clearlongitude_val    bool
	range_from_val        *float64
	addrange_from_val     *float64
	clearrange_from_val   bool
	range_to_val          *float64
	addrange_to_val       *float64
	clearrange_to_val     bool
	string_val            *string
	clearstring_val       bool
	_type                 map[string]struct{}
	location              map[string]struct{}
	equipment             map[string]struct{}
	service               map[string]struct{}
	equipment_port        map[string]struct{}
	link                  map[string]struct{}
	work_order            map[string]struct{}
	project               map[string]struct{}
	equipment_value       map[string]struct{}
	location_value        map[string]struct{}
	service_value         map[string]struct{}
	clearedType           bool
	clearedLocation       bool
	clearedEquipment      bool
	clearedService        bool
	clearedEquipmentPort  bool
	clearedLink           bool
	clearedWorkOrder      bool
	clearedProject        bool
	clearedEquipmentValue bool
	clearedLocationValue  bool
	clearedServiceValue   bool
}

// SetIntVal sets the int_val field.
func (puo *PropertyUpdateOne) SetIntVal(i int) *PropertyUpdateOne {
	puo.int_val = &i
	puo.addint_val = nil
	return puo
}

// SetNillableIntVal sets the int_val field if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableIntVal(i *int) *PropertyUpdateOne {
	if i != nil {
		puo.SetIntVal(*i)
	}
	return puo
}

// AddIntVal adds i to int_val.
func (puo *PropertyUpdateOne) AddIntVal(i int) *PropertyUpdateOne {
	if puo.addint_val == nil {
		puo.addint_val = &i
	} else {
		*puo.addint_val += i
	}
	return puo
}

// ClearIntVal clears the value of int_val.
func (puo *PropertyUpdateOne) ClearIntVal() *PropertyUpdateOne {
	puo.int_val = nil
	puo.clearint_val = true
	return puo
}

// SetBoolVal sets the bool_val field.
func (puo *PropertyUpdateOne) SetBoolVal(b bool) *PropertyUpdateOne {
	puo.bool_val = &b
	return puo
}

// SetNillableBoolVal sets the bool_val field if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableBoolVal(b *bool) *PropertyUpdateOne {
	if b != nil {
		puo.SetBoolVal(*b)
	}
	return puo
}

// ClearBoolVal clears the value of bool_val.
func (puo *PropertyUpdateOne) ClearBoolVal() *PropertyUpdateOne {
	puo.bool_val = nil
	puo.clearbool_val = true
	return puo
}

// SetFloatVal sets the float_val field.
func (puo *PropertyUpdateOne) SetFloatVal(f float64) *PropertyUpdateOne {
	puo.float_val = &f
	puo.addfloat_val = nil
	return puo
}

// SetNillableFloatVal sets the float_val field if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableFloatVal(f *float64) *PropertyUpdateOne {
	if f != nil {
		puo.SetFloatVal(*f)
	}
	return puo
}

// AddFloatVal adds f to float_val.
func (puo *PropertyUpdateOne) AddFloatVal(f float64) *PropertyUpdateOne {
	if puo.addfloat_val == nil {
		puo.addfloat_val = &f
	} else {
		*puo.addfloat_val += f
	}
	return puo
}

// ClearFloatVal clears the value of float_val.
func (puo *PropertyUpdateOne) ClearFloatVal() *PropertyUpdateOne {
	puo.float_val = nil
	puo.clearfloat_val = true
	return puo
}

// SetLatitudeVal sets the latitude_val field.
func (puo *PropertyUpdateOne) SetLatitudeVal(f float64) *PropertyUpdateOne {
	puo.latitude_val = &f
	puo.addlatitude_val = nil
	return puo
}

// SetNillableLatitudeVal sets the latitude_val field if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableLatitudeVal(f *float64) *PropertyUpdateOne {
	if f != nil {
		puo.SetLatitudeVal(*f)
	}
	return puo
}

// AddLatitudeVal adds f to latitude_val.
func (puo *PropertyUpdateOne) AddLatitudeVal(f float64) *PropertyUpdateOne {
	if puo.addlatitude_val == nil {
		puo.addlatitude_val = &f
	} else {
		*puo.addlatitude_val += f
	}
	return puo
}

// ClearLatitudeVal clears the value of latitude_val.
func (puo *PropertyUpdateOne) ClearLatitudeVal() *PropertyUpdateOne {
	puo.latitude_val = nil
	puo.clearlatitude_val = true
	return puo
}

// SetLongitudeVal sets the longitude_val field.
func (puo *PropertyUpdateOne) SetLongitudeVal(f float64) *PropertyUpdateOne {
	puo.longitude_val = &f
	puo.addlongitude_val = nil
	return puo
}

// SetNillableLongitudeVal sets the longitude_val field if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableLongitudeVal(f *float64) *PropertyUpdateOne {
	if f != nil {
		puo.SetLongitudeVal(*f)
	}
	return puo
}

// AddLongitudeVal adds f to longitude_val.
func (puo *PropertyUpdateOne) AddLongitudeVal(f float64) *PropertyUpdateOne {
	if puo.addlongitude_val == nil {
		puo.addlongitude_val = &f
	} else {
		*puo.addlongitude_val += f
	}
	return puo
}

// ClearLongitudeVal clears the value of longitude_val.
func (puo *PropertyUpdateOne) ClearLongitudeVal() *PropertyUpdateOne {
	puo.longitude_val = nil
	puo.clearlongitude_val = true
	return puo
}

// SetRangeFromVal sets the range_from_val field.
func (puo *PropertyUpdateOne) SetRangeFromVal(f float64) *PropertyUpdateOne {
	puo.range_from_val = &f
	puo.addrange_from_val = nil
	return puo
}

// SetNillableRangeFromVal sets the range_from_val field if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableRangeFromVal(f *float64) *PropertyUpdateOne {
	if f != nil {
		puo.SetRangeFromVal(*f)
	}
	return puo
}

// AddRangeFromVal adds f to range_from_val.
func (puo *PropertyUpdateOne) AddRangeFromVal(f float64) *PropertyUpdateOne {
	if puo.addrange_from_val == nil {
		puo.addrange_from_val = &f
	} else {
		*puo.addrange_from_val += f
	}
	return puo
}

// ClearRangeFromVal clears the value of range_from_val.
func (puo *PropertyUpdateOne) ClearRangeFromVal() *PropertyUpdateOne {
	puo.range_from_val = nil
	puo.clearrange_from_val = true
	return puo
}

// SetRangeToVal sets the range_to_val field.
func (puo *PropertyUpdateOne) SetRangeToVal(f float64) *PropertyUpdateOne {
	puo.range_to_val = &f
	puo.addrange_to_val = nil
	return puo
}

// SetNillableRangeToVal sets the range_to_val field if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableRangeToVal(f *float64) *PropertyUpdateOne {
	if f != nil {
		puo.SetRangeToVal(*f)
	}
	return puo
}

// AddRangeToVal adds f to range_to_val.
func (puo *PropertyUpdateOne) AddRangeToVal(f float64) *PropertyUpdateOne {
	if puo.addrange_to_val == nil {
		puo.addrange_to_val = &f
	} else {
		*puo.addrange_to_val += f
	}
	return puo
}

// ClearRangeToVal clears the value of range_to_val.
func (puo *PropertyUpdateOne) ClearRangeToVal() *PropertyUpdateOne {
	puo.range_to_val = nil
	puo.clearrange_to_val = true
	return puo
}

// SetStringVal sets the string_val field.
func (puo *PropertyUpdateOne) SetStringVal(s string) *PropertyUpdateOne {
	puo.string_val = &s
	return puo
}

// SetNillableStringVal sets the string_val field if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableStringVal(s *string) *PropertyUpdateOne {
	if s != nil {
		puo.SetStringVal(*s)
	}
	return puo
}

// ClearStringVal clears the value of string_val.
func (puo *PropertyUpdateOne) ClearStringVal() *PropertyUpdateOne {
	puo.string_val = nil
	puo.clearstring_val = true
	return puo
}

// SetTypeID sets the type edge to PropertyType by id.
func (puo *PropertyUpdateOne) SetTypeID(id string) *PropertyUpdateOne {
	if puo._type == nil {
		puo._type = make(map[string]struct{})
	}
	puo._type[id] = struct{}{}
	return puo
}

// SetType sets the type edge to PropertyType.
func (puo *PropertyUpdateOne) SetType(p *PropertyType) *PropertyUpdateOne {
	return puo.SetTypeID(p.ID)
}

// SetLocationID sets the location edge to Location by id.
func (puo *PropertyUpdateOne) SetLocationID(id string) *PropertyUpdateOne {
	if puo.location == nil {
		puo.location = make(map[string]struct{})
	}
	puo.location[id] = struct{}{}
	return puo
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableLocationID(id *string) *PropertyUpdateOne {
	if id != nil {
		puo = puo.SetLocationID(*id)
	}
	return puo
}

// SetLocation sets the location edge to Location.
func (puo *PropertyUpdateOne) SetLocation(l *Location) *PropertyUpdateOne {
	return puo.SetLocationID(l.ID)
}

// SetEquipmentID sets the equipment edge to Equipment by id.
func (puo *PropertyUpdateOne) SetEquipmentID(id string) *PropertyUpdateOne {
	if puo.equipment == nil {
		puo.equipment = make(map[string]struct{})
	}
	puo.equipment[id] = struct{}{}
	return puo
}

// SetNillableEquipmentID sets the equipment edge to Equipment by id if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableEquipmentID(id *string) *PropertyUpdateOne {
	if id != nil {
		puo = puo.SetEquipmentID(*id)
	}
	return puo
}

// SetEquipment sets the equipment edge to Equipment.
func (puo *PropertyUpdateOne) SetEquipment(e *Equipment) *PropertyUpdateOne {
	return puo.SetEquipmentID(e.ID)
}

// SetServiceID sets the service edge to Service by id.
func (puo *PropertyUpdateOne) SetServiceID(id string) *PropertyUpdateOne {
	if puo.service == nil {
		puo.service = make(map[string]struct{})
	}
	puo.service[id] = struct{}{}
	return puo
}

// SetNillableServiceID sets the service edge to Service by id if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableServiceID(id *string) *PropertyUpdateOne {
	if id != nil {
		puo = puo.SetServiceID(*id)
	}
	return puo
}

// SetService sets the service edge to Service.
func (puo *PropertyUpdateOne) SetService(s *Service) *PropertyUpdateOne {
	return puo.SetServiceID(s.ID)
}

// SetEquipmentPortID sets the equipment_port edge to EquipmentPort by id.
func (puo *PropertyUpdateOne) SetEquipmentPortID(id string) *PropertyUpdateOne {
	if puo.equipment_port == nil {
		puo.equipment_port = make(map[string]struct{})
	}
	puo.equipment_port[id] = struct{}{}
	return puo
}

// SetNillableEquipmentPortID sets the equipment_port edge to EquipmentPort by id if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableEquipmentPortID(id *string) *PropertyUpdateOne {
	if id != nil {
		puo = puo.SetEquipmentPortID(*id)
	}
	return puo
}

// SetEquipmentPort sets the equipment_port edge to EquipmentPort.
func (puo *PropertyUpdateOne) SetEquipmentPort(e *EquipmentPort) *PropertyUpdateOne {
	return puo.SetEquipmentPortID(e.ID)
}

// SetLinkID sets the link edge to Link by id.
func (puo *PropertyUpdateOne) SetLinkID(id string) *PropertyUpdateOne {
	if puo.link == nil {
		puo.link = make(map[string]struct{})
	}
	puo.link[id] = struct{}{}
	return puo
}

// SetNillableLinkID sets the link edge to Link by id if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableLinkID(id *string) *PropertyUpdateOne {
	if id != nil {
		puo = puo.SetLinkID(*id)
	}
	return puo
}

// SetLink sets the link edge to Link.
func (puo *PropertyUpdateOne) SetLink(l *Link) *PropertyUpdateOne {
	return puo.SetLinkID(l.ID)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (puo *PropertyUpdateOne) SetWorkOrderID(id string) *PropertyUpdateOne {
	if puo.work_order == nil {
		puo.work_order = make(map[string]struct{})
	}
	puo.work_order[id] = struct{}{}
	return puo
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableWorkOrderID(id *string) *PropertyUpdateOne {
	if id != nil {
		puo = puo.SetWorkOrderID(*id)
	}
	return puo
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (puo *PropertyUpdateOne) SetWorkOrder(w *WorkOrder) *PropertyUpdateOne {
	return puo.SetWorkOrderID(w.ID)
}

// SetProjectID sets the project edge to Project by id.
func (puo *PropertyUpdateOne) SetProjectID(id string) *PropertyUpdateOne {
	if puo.project == nil {
		puo.project = make(map[string]struct{})
	}
	puo.project[id] = struct{}{}
	return puo
}

// SetNillableProjectID sets the project edge to Project by id if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableProjectID(id *string) *PropertyUpdateOne {
	if id != nil {
		puo = puo.SetProjectID(*id)
	}
	return puo
}

// SetProject sets the project edge to Project.
func (puo *PropertyUpdateOne) SetProject(p *Project) *PropertyUpdateOne {
	return puo.SetProjectID(p.ID)
}

// SetEquipmentValueID sets the equipment_value edge to Equipment by id.
func (puo *PropertyUpdateOne) SetEquipmentValueID(id string) *PropertyUpdateOne {
	if puo.equipment_value == nil {
		puo.equipment_value = make(map[string]struct{})
	}
	puo.equipment_value[id] = struct{}{}
	return puo
}

// SetNillableEquipmentValueID sets the equipment_value edge to Equipment by id if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableEquipmentValueID(id *string) *PropertyUpdateOne {
	if id != nil {
		puo = puo.SetEquipmentValueID(*id)
	}
	return puo
}

// SetEquipmentValue sets the equipment_value edge to Equipment.
func (puo *PropertyUpdateOne) SetEquipmentValue(e *Equipment) *PropertyUpdateOne {
	return puo.SetEquipmentValueID(e.ID)
}

// SetLocationValueID sets the location_value edge to Location by id.
func (puo *PropertyUpdateOne) SetLocationValueID(id string) *PropertyUpdateOne {
	if puo.location_value == nil {
		puo.location_value = make(map[string]struct{})
	}
	puo.location_value[id] = struct{}{}
	return puo
}

// SetNillableLocationValueID sets the location_value edge to Location by id if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableLocationValueID(id *string) *PropertyUpdateOne {
	if id != nil {
		puo = puo.SetLocationValueID(*id)
	}
	return puo
}

// SetLocationValue sets the location_value edge to Location.
func (puo *PropertyUpdateOne) SetLocationValue(l *Location) *PropertyUpdateOne {
	return puo.SetLocationValueID(l.ID)
}

// SetServiceValueID sets the service_value edge to Service by id.
func (puo *PropertyUpdateOne) SetServiceValueID(id string) *PropertyUpdateOne {
	if puo.service_value == nil {
		puo.service_value = make(map[string]struct{})
	}
	puo.service_value[id] = struct{}{}
	return puo
}

// SetNillableServiceValueID sets the service_value edge to Service by id if the given value is not nil.
func (puo *PropertyUpdateOne) SetNillableServiceValueID(id *string) *PropertyUpdateOne {
	if id != nil {
		puo = puo.SetServiceValueID(*id)
	}
	return puo
}

// SetServiceValue sets the service_value edge to Service.
func (puo *PropertyUpdateOne) SetServiceValue(s *Service) *PropertyUpdateOne {
	return puo.SetServiceValueID(s.ID)
}

// ClearType clears the type edge to PropertyType.
func (puo *PropertyUpdateOne) ClearType() *PropertyUpdateOne {
	puo.clearedType = true
	return puo
}

// ClearLocation clears the location edge to Location.
func (puo *PropertyUpdateOne) ClearLocation() *PropertyUpdateOne {
	puo.clearedLocation = true
	return puo
}

// ClearEquipment clears the equipment edge to Equipment.
func (puo *PropertyUpdateOne) ClearEquipment() *PropertyUpdateOne {
	puo.clearedEquipment = true
	return puo
}

// ClearService clears the service edge to Service.
func (puo *PropertyUpdateOne) ClearService() *PropertyUpdateOne {
	puo.clearedService = true
	return puo
}

// ClearEquipmentPort clears the equipment_port edge to EquipmentPort.
func (puo *PropertyUpdateOne) ClearEquipmentPort() *PropertyUpdateOne {
	puo.clearedEquipmentPort = true
	return puo
}

// ClearLink clears the link edge to Link.
func (puo *PropertyUpdateOne) ClearLink() *PropertyUpdateOne {
	puo.clearedLink = true
	return puo
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (puo *PropertyUpdateOne) ClearWorkOrder() *PropertyUpdateOne {
	puo.clearedWorkOrder = true
	return puo
}

// ClearProject clears the project edge to Project.
func (puo *PropertyUpdateOne) ClearProject() *PropertyUpdateOne {
	puo.clearedProject = true
	return puo
}

// ClearEquipmentValue clears the equipment_value edge to Equipment.
func (puo *PropertyUpdateOne) ClearEquipmentValue() *PropertyUpdateOne {
	puo.clearedEquipmentValue = true
	return puo
}

// ClearLocationValue clears the location_value edge to Location.
func (puo *PropertyUpdateOne) ClearLocationValue() *PropertyUpdateOne {
	puo.clearedLocationValue = true
	return puo
}

// ClearServiceValue clears the service_value edge to Service.
func (puo *PropertyUpdateOne) ClearServiceValue() *PropertyUpdateOne {
	puo.clearedServiceValue = true
	return puo
}

// Save executes the query and returns the updated entity.
func (puo *PropertyUpdateOne) Save(ctx context.Context) (*Property, error) {
	if puo.update_time == nil {
		v := property.UpdateDefaultUpdateTime()
		puo.update_time = &v
	}
	if len(puo._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if puo.clearedType && puo._type == nil {
		return nil, errors.New("ent: clearing a unique edge \"type\"")
	}
	if len(puo.location) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	if len(puo.equipment) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"equipment\"")
	}
	if len(puo.service) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"service\"")
	}
	if len(puo.equipment_port) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"equipment_port\"")
	}
	if len(puo.link) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"link\"")
	}
	if len(puo.work_order) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"work_order\"")
	}
	if len(puo.project) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"project\"")
	}
	if len(puo.equipment_value) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"equipment_value\"")
	}
	if len(puo.location_value) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location_value\"")
	}
	if len(puo.service_value) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"service_value\"")
	}
	return puo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (puo *PropertyUpdateOne) SaveX(ctx context.Context) *Property {
	pr, err := puo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return pr
}

// Exec executes the query on the entity.
func (puo *PropertyUpdateOne) Exec(ctx context.Context) error {
	_, err := puo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (puo *PropertyUpdateOne) ExecX(ctx context.Context) {
	if err := puo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (puo *PropertyUpdateOne) sqlSave(ctx context.Context) (pr *Property, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   property.Table,
			Columns: property.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  puo.id,
				Type:   field.TypeString,
				Column: property.FieldID,
			},
		},
	}
	if value := puo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: property.FieldUpdateTime,
		})
	}
	if value := puo.int_val; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: property.FieldIntVal,
		})
	}
	if value := puo.addint_val; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: property.FieldIntVal,
		})
	}
	if puo.clearint_val {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: property.FieldIntVal,
		})
	}
	if value := puo.bool_val; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: property.FieldBoolVal,
		})
	}
	if puo.clearbool_val {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: property.FieldBoolVal,
		})
	}
	if value := puo.float_val; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldFloatVal,
		})
	}
	if value := puo.addfloat_val; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldFloatVal,
		})
	}
	if puo.clearfloat_val {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: property.FieldFloatVal,
		})
	}
	if value := puo.latitude_val; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldLatitudeVal,
		})
	}
	if value := puo.addlatitude_val; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldLatitudeVal,
		})
	}
	if puo.clearlatitude_val {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: property.FieldLatitudeVal,
		})
	}
	if value := puo.longitude_val; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldLongitudeVal,
		})
	}
	if value := puo.addlongitude_val; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldLongitudeVal,
		})
	}
	if puo.clearlongitude_val {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: property.FieldLongitudeVal,
		})
	}
	if value := puo.range_from_val; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldRangeFromVal,
		})
	}
	if value := puo.addrange_from_val; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldRangeFromVal,
		})
	}
	if puo.clearrange_from_val {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: property.FieldRangeFromVal,
		})
	}
	if value := puo.range_to_val; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldRangeToVal,
		})
	}
	if value := puo.addrange_to_val; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: property.FieldRangeToVal,
		})
	}
	if puo.clearrange_to_val {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: property.FieldRangeToVal,
		})
	}
	if value := puo.string_val; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: property.FieldStringVal,
		})
	}
	if puo.clearstring_val {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: property.FieldStringVal,
		})
	}
	if puo.clearedType {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.TypeTable,
			Columns: []string{property.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: propertytype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo._type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.TypeTable,
			Columns: []string{property.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: propertytype.FieldID,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if puo.clearedLocation {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.LocationTable,
			Columns: []string{property.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.location; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.LocationTable,
			Columns: []string{property.LocationColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if puo.clearedEquipment {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.EquipmentTable,
			Columns: []string{property.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipment.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.equipment; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.EquipmentTable,
			Columns: []string{property.EquipmentColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if puo.clearedService {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.ServiceTable,
			Columns: []string{property.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.service; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.ServiceTable,
			Columns: []string{property.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if puo.clearedEquipmentPort {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.EquipmentPortTable,
			Columns: []string{property.EquipmentPortColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentport.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.equipment_port; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.EquipmentPortTable,
			Columns: []string{property.EquipmentPortColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentport.FieldID,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if puo.clearedLink {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.LinkTable,
			Columns: []string{property.LinkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: link.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.link; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.LinkTable,
			Columns: []string{property.LinkColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if puo.clearedWorkOrder {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.WorkOrderTable,
			Columns: []string{property.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.work_order; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.WorkOrderTable,
			Columns: []string{property.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if puo.clearedProject {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.ProjectTable,
			Columns: []string{property.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: project.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.project; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.ProjectTable,
			Columns: []string{property.ProjectColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if puo.clearedEquipmentValue {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.EquipmentValueTable,
			Columns: []string{property.EquipmentValueColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipment.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.equipment_value; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.EquipmentValueTable,
			Columns: []string{property.EquipmentValueColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if puo.clearedLocationValue {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.LocationValueTable,
			Columns: []string{property.LocationValueColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.location_value; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.LocationValueTable,
			Columns: []string{property.LocationValueColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if puo.clearedServiceValue {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.ServiceValueTable,
			Columns: []string{property.ServiceValueColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.service_value; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.ServiceValueTable,
			Columns: []string{property.ServiceValueColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	pr = &Property{config: puo.config}
	_spec.Assign = pr.assignValues
	_spec.ScanValues = pr.scanValues()
	if err = sqlgraph.UpdateNode(ctx, puo.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return pr, nil
}
