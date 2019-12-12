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
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// PropertyTypeUpdate is the builder for updating PropertyType entities.
type PropertyTypeUpdate struct {
	config

	update_time                  *time.Time
	_type                        *string
	name                         *string
	index                        *int
	addindex                     *int
	clearindex                   bool
	category                     *string
	clearcategory                bool
	int_val                      *int
	addint_val                   *int
	clearint_val                 bool
	bool_val                     *bool
	clearbool_val                bool
	float_val                    *float64
	addfloat_val                 *float64
	clearfloat_val               bool
	latitude_val                 *float64
	addlatitude_val              *float64
	clearlatitude_val            bool
	longitude_val                *float64
	addlongitude_val             *float64
	clearlongitude_val           bool
	string_val                   *string
	clearstring_val              bool
	range_from_val               *float64
	addrange_from_val            *float64
	clearrange_from_val          bool
	range_to_val                 *float64
	addrange_to_val              *float64
	clearrange_to_val            bool
	is_instance_property         *bool
	editable                     *bool
	mandatory                    *bool
	properties                   map[string]struct{}
	location_type                map[string]struct{}
	equipment_port_type          map[string]struct{}
	link_equipment_port_type     map[string]struct{}
	equipment_type               map[string]struct{}
	service_type                 map[string]struct{}
	work_order_type              map[string]struct{}
	project_type                 map[string]struct{}
	removedProperties            map[string]struct{}
	clearedLocationType          bool
	clearedEquipmentPortType     bool
	clearedLinkEquipmentPortType bool
	clearedEquipmentType         bool
	clearedServiceType           bool
	clearedWorkOrderType         bool
	clearedProjectType           bool
	predicates                   []predicate.PropertyType
}

// Where adds a new predicate for the builder.
func (ptu *PropertyTypeUpdate) Where(ps ...predicate.PropertyType) *PropertyTypeUpdate {
	ptu.predicates = append(ptu.predicates, ps...)
	return ptu
}

// SetType sets the type field.
func (ptu *PropertyTypeUpdate) SetType(s string) *PropertyTypeUpdate {
	ptu._type = &s
	return ptu
}

// SetName sets the name field.
func (ptu *PropertyTypeUpdate) SetName(s string) *PropertyTypeUpdate {
	ptu.name = &s
	return ptu
}

// SetIndex sets the index field.
func (ptu *PropertyTypeUpdate) SetIndex(i int) *PropertyTypeUpdate {
	ptu.index = &i
	ptu.addindex = nil
	return ptu
}

// SetNillableIndex sets the index field if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableIndex(i *int) *PropertyTypeUpdate {
	if i != nil {
		ptu.SetIndex(*i)
	}
	return ptu
}

// AddIndex adds i to index.
func (ptu *PropertyTypeUpdate) AddIndex(i int) *PropertyTypeUpdate {
	if ptu.addindex == nil {
		ptu.addindex = &i
	} else {
		*ptu.addindex += i
	}
	return ptu
}

// ClearIndex clears the value of index.
func (ptu *PropertyTypeUpdate) ClearIndex() *PropertyTypeUpdate {
	ptu.index = nil
	ptu.clearindex = true
	return ptu
}

// SetCategory sets the category field.
func (ptu *PropertyTypeUpdate) SetCategory(s string) *PropertyTypeUpdate {
	ptu.category = &s
	return ptu
}

// SetNillableCategory sets the category field if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableCategory(s *string) *PropertyTypeUpdate {
	if s != nil {
		ptu.SetCategory(*s)
	}
	return ptu
}

// ClearCategory clears the value of category.
func (ptu *PropertyTypeUpdate) ClearCategory() *PropertyTypeUpdate {
	ptu.category = nil
	ptu.clearcategory = true
	return ptu
}

// SetIntVal sets the int_val field.
func (ptu *PropertyTypeUpdate) SetIntVal(i int) *PropertyTypeUpdate {
	ptu.int_val = &i
	ptu.addint_val = nil
	return ptu
}

// SetNillableIntVal sets the int_val field if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableIntVal(i *int) *PropertyTypeUpdate {
	if i != nil {
		ptu.SetIntVal(*i)
	}
	return ptu
}

// AddIntVal adds i to int_val.
func (ptu *PropertyTypeUpdate) AddIntVal(i int) *PropertyTypeUpdate {
	if ptu.addint_val == nil {
		ptu.addint_val = &i
	} else {
		*ptu.addint_val += i
	}
	return ptu
}

// ClearIntVal clears the value of int_val.
func (ptu *PropertyTypeUpdate) ClearIntVal() *PropertyTypeUpdate {
	ptu.int_val = nil
	ptu.clearint_val = true
	return ptu
}

// SetBoolVal sets the bool_val field.
func (ptu *PropertyTypeUpdate) SetBoolVal(b bool) *PropertyTypeUpdate {
	ptu.bool_val = &b
	return ptu
}

// SetNillableBoolVal sets the bool_val field if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableBoolVal(b *bool) *PropertyTypeUpdate {
	if b != nil {
		ptu.SetBoolVal(*b)
	}
	return ptu
}

// ClearBoolVal clears the value of bool_val.
func (ptu *PropertyTypeUpdate) ClearBoolVal() *PropertyTypeUpdate {
	ptu.bool_val = nil
	ptu.clearbool_val = true
	return ptu
}

// SetFloatVal sets the float_val field.
func (ptu *PropertyTypeUpdate) SetFloatVal(f float64) *PropertyTypeUpdate {
	ptu.float_val = &f
	ptu.addfloat_val = nil
	return ptu
}

// SetNillableFloatVal sets the float_val field if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableFloatVal(f *float64) *PropertyTypeUpdate {
	if f != nil {
		ptu.SetFloatVal(*f)
	}
	return ptu
}

// AddFloatVal adds f to float_val.
func (ptu *PropertyTypeUpdate) AddFloatVal(f float64) *PropertyTypeUpdate {
	if ptu.addfloat_val == nil {
		ptu.addfloat_val = &f
	} else {
		*ptu.addfloat_val += f
	}
	return ptu
}

// ClearFloatVal clears the value of float_val.
func (ptu *PropertyTypeUpdate) ClearFloatVal() *PropertyTypeUpdate {
	ptu.float_val = nil
	ptu.clearfloat_val = true
	return ptu
}

// SetLatitudeVal sets the latitude_val field.
func (ptu *PropertyTypeUpdate) SetLatitudeVal(f float64) *PropertyTypeUpdate {
	ptu.latitude_val = &f
	ptu.addlatitude_val = nil
	return ptu
}

// SetNillableLatitudeVal sets the latitude_val field if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableLatitudeVal(f *float64) *PropertyTypeUpdate {
	if f != nil {
		ptu.SetLatitudeVal(*f)
	}
	return ptu
}

// AddLatitudeVal adds f to latitude_val.
func (ptu *PropertyTypeUpdate) AddLatitudeVal(f float64) *PropertyTypeUpdate {
	if ptu.addlatitude_val == nil {
		ptu.addlatitude_val = &f
	} else {
		*ptu.addlatitude_val += f
	}
	return ptu
}

// ClearLatitudeVal clears the value of latitude_val.
func (ptu *PropertyTypeUpdate) ClearLatitudeVal() *PropertyTypeUpdate {
	ptu.latitude_val = nil
	ptu.clearlatitude_val = true
	return ptu
}

// SetLongitudeVal sets the longitude_val field.
func (ptu *PropertyTypeUpdate) SetLongitudeVal(f float64) *PropertyTypeUpdate {
	ptu.longitude_val = &f
	ptu.addlongitude_val = nil
	return ptu
}

// SetNillableLongitudeVal sets the longitude_val field if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableLongitudeVal(f *float64) *PropertyTypeUpdate {
	if f != nil {
		ptu.SetLongitudeVal(*f)
	}
	return ptu
}

// AddLongitudeVal adds f to longitude_val.
func (ptu *PropertyTypeUpdate) AddLongitudeVal(f float64) *PropertyTypeUpdate {
	if ptu.addlongitude_val == nil {
		ptu.addlongitude_val = &f
	} else {
		*ptu.addlongitude_val += f
	}
	return ptu
}

// ClearLongitudeVal clears the value of longitude_val.
func (ptu *PropertyTypeUpdate) ClearLongitudeVal() *PropertyTypeUpdate {
	ptu.longitude_val = nil
	ptu.clearlongitude_val = true
	return ptu
}

// SetStringVal sets the string_val field.
func (ptu *PropertyTypeUpdate) SetStringVal(s string) *PropertyTypeUpdate {
	ptu.string_val = &s
	return ptu
}

// SetNillableStringVal sets the string_val field if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableStringVal(s *string) *PropertyTypeUpdate {
	if s != nil {
		ptu.SetStringVal(*s)
	}
	return ptu
}

// ClearStringVal clears the value of string_val.
func (ptu *PropertyTypeUpdate) ClearStringVal() *PropertyTypeUpdate {
	ptu.string_val = nil
	ptu.clearstring_val = true
	return ptu
}

// SetRangeFromVal sets the range_from_val field.
func (ptu *PropertyTypeUpdate) SetRangeFromVal(f float64) *PropertyTypeUpdate {
	ptu.range_from_val = &f
	ptu.addrange_from_val = nil
	return ptu
}

// SetNillableRangeFromVal sets the range_from_val field if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableRangeFromVal(f *float64) *PropertyTypeUpdate {
	if f != nil {
		ptu.SetRangeFromVal(*f)
	}
	return ptu
}

// AddRangeFromVal adds f to range_from_val.
func (ptu *PropertyTypeUpdate) AddRangeFromVal(f float64) *PropertyTypeUpdate {
	if ptu.addrange_from_val == nil {
		ptu.addrange_from_val = &f
	} else {
		*ptu.addrange_from_val += f
	}
	return ptu
}

// ClearRangeFromVal clears the value of range_from_val.
func (ptu *PropertyTypeUpdate) ClearRangeFromVal() *PropertyTypeUpdate {
	ptu.range_from_val = nil
	ptu.clearrange_from_val = true
	return ptu
}

// SetRangeToVal sets the range_to_val field.
func (ptu *PropertyTypeUpdate) SetRangeToVal(f float64) *PropertyTypeUpdate {
	ptu.range_to_val = &f
	ptu.addrange_to_val = nil
	return ptu
}

// SetNillableRangeToVal sets the range_to_val field if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableRangeToVal(f *float64) *PropertyTypeUpdate {
	if f != nil {
		ptu.SetRangeToVal(*f)
	}
	return ptu
}

// AddRangeToVal adds f to range_to_val.
func (ptu *PropertyTypeUpdate) AddRangeToVal(f float64) *PropertyTypeUpdate {
	if ptu.addrange_to_val == nil {
		ptu.addrange_to_val = &f
	} else {
		*ptu.addrange_to_val += f
	}
	return ptu
}

// ClearRangeToVal clears the value of range_to_val.
func (ptu *PropertyTypeUpdate) ClearRangeToVal() *PropertyTypeUpdate {
	ptu.range_to_val = nil
	ptu.clearrange_to_val = true
	return ptu
}

// SetIsInstanceProperty sets the is_instance_property field.
func (ptu *PropertyTypeUpdate) SetIsInstanceProperty(b bool) *PropertyTypeUpdate {
	ptu.is_instance_property = &b
	return ptu
}

// SetNillableIsInstanceProperty sets the is_instance_property field if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableIsInstanceProperty(b *bool) *PropertyTypeUpdate {
	if b != nil {
		ptu.SetIsInstanceProperty(*b)
	}
	return ptu
}

// SetEditable sets the editable field.
func (ptu *PropertyTypeUpdate) SetEditable(b bool) *PropertyTypeUpdate {
	ptu.editable = &b
	return ptu
}

// SetNillableEditable sets the editable field if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableEditable(b *bool) *PropertyTypeUpdate {
	if b != nil {
		ptu.SetEditable(*b)
	}
	return ptu
}

// SetMandatory sets the mandatory field.
func (ptu *PropertyTypeUpdate) SetMandatory(b bool) *PropertyTypeUpdate {
	ptu.mandatory = &b
	return ptu
}

// SetNillableMandatory sets the mandatory field if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableMandatory(b *bool) *PropertyTypeUpdate {
	if b != nil {
		ptu.SetMandatory(*b)
	}
	return ptu
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (ptu *PropertyTypeUpdate) AddPropertyIDs(ids ...string) *PropertyTypeUpdate {
	if ptu.properties == nil {
		ptu.properties = make(map[string]struct{})
	}
	for i := range ids {
		ptu.properties[ids[i]] = struct{}{}
	}
	return ptu
}

// AddProperties adds the properties edges to Property.
func (ptu *PropertyTypeUpdate) AddProperties(p ...*Property) *PropertyTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptu.AddPropertyIDs(ids...)
}

// SetLocationTypeID sets the location_type edge to LocationType by id.
func (ptu *PropertyTypeUpdate) SetLocationTypeID(id string) *PropertyTypeUpdate {
	if ptu.location_type == nil {
		ptu.location_type = make(map[string]struct{})
	}
	ptu.location_type[id] = struct{}{}
	return ptu
}

// SetNillableLocationTypeID sets the location_type edge to LocationType by id if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableLocationTypeID(id *string) *PropertyTypeUpdate {
	if id != nil {
		ptu = ptu.SetLocationTypeID(*id)
	}
	return ptu
}

// SetLocationType sets the location_type edge to LocationType.
func (ptu *PropertyTypeUpdate) SetLocationType(l *LocationType) *PropertyTypeUpdate {
	return ptu.SetLocationTypeID(l.ID)
}

// SetEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id.
func (ptu *PropertyTypeUpdate) SetEquipmentPortTypeID(id string) *PropertyTypeUpdate {
	if ptu.equipment_port_type == nil {
		ptu.equipment_port_type = make(map[string]struct{})
	}
	ptu.equipment_port_type[id] = struct{}{}
	return ptu
}

// SetNillableEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableEquipmentPortTypeID(id *string) *PropertyTypeUpdate {
	if id != nil {
		ptu = ptu.SetEquipmentPortTypeID(*id)
	}
	return ptu
}

// SetEquipmentPortType sets the equipment_port_type edge to EquipmentPortType.
func (ptu *PropertyTypeUpdate) SetEquipmentPortType(e *EquipmentPortType) *PropertyTypeUpdate {
	return ptu.SetEquipmentPortTypeID(e.ID)
}

// SetLinkEquipmentPortTypeID sets the link_equipment_port_type edge to EquipmentPortType by id.
func (ptu *PropertyTypeUpdate) SetLinkEquipmentPortTypeID(id string) *PropertyTypeUpdate {
	if ptu.link_equipment_port_type == nil {
		ptu.link_equipment_port_type = make(map[string]struct{})
	}
	ptu.link_equipment_port_type[id] = struct{}{}
	return ptu
}

// SetNillableLinkEquipmentPortTypeID sets the link_equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableLinkEquipmentPortTypeID(id *string) *PropertyTypeUpdate {
	if id != nil {
		ptu = ptu.SetLinkEquipmentPortTypeID(*id)
	}
	return ptu
}

// SetLinkEquipmentPortType sets the link_equipment_port_type edge to EquipmentPortType.
func (ptu *PropertyTypeUpdate) SetLinkEquipmentPortType(e *EquipmentPortType) *PropertyTypeUpdate {
	return ptu.SetLinkEquipmentPortTypeID(e.ID)
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (ptu *PropertyTypeUpdate) SetEquipmentTypeID(id string) *PropertyTypeUpdate {
	if ptu.equipment_type == nil {
		ptu.equipment_type = make(map[string]struct{})
	}
	ptu.equipment_type[id] = struct{}{}
	return ptu
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableEquipmentTypeID(id *string) *PropertyTypeUpdate {
	if id != nil {
		ptu = ptu.SetEquipmentTypeID(*id)
	}
	return ptu
}

// SetEquipmentType sets the equipment_type edge to EquipmentType.
func (ptu *PropertyTypeUpdate) SetEquipmentType(e *EquipmentType) *PropertyTypeUpdate {
	return ptu.SetEquipmentTypeID(e.ID)
}

// SetServiceTypeID sets the service_type edge to ServiceType by id.
func (ptu *PropertyTypeUpdate) SetServiceTypeID(id string) *PropertyTypeUpdate {
	if ptu.service_type == nil {
		ptu.service_type = make(map[string]struct{})
	}
	ptu.service_type[id] = struct{}{}
	return ptu
}

// SetNillableServiceTypeID sets the service_type edge to ServiceType by id if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableServiceTypeID(id *string) *PropertyTypeUpdate {
	if id != nil {
		ptu = ptu.SetServiceTypeID(*id)
	}
	return ptu
}

// SetServiceType sets the service_type edge to ServiceType.
func (ptu *PropertyTypeUpdate) SetServiceType(s *ServiceType) *PropertyTypeUpdate {
	return ptu.SetServiceTypeID(s.ID)
}

// SetWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id.
func (ptu *PropertyTypeUpdate) SetWorkOrderTypeID(id string) *PropertyTypeUpdate {
	if ptu.work_order_type == nil {
		ptu.work_order_type = make(map[string]struct{})
	}
	ptu.work_order_type[id] = struct{}{}
	return ptu
}

// SetNillableWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableWorkOrderTypeID(id *string) *PropertyTypeUpdate {
	if id != nil {
		ptu = ptu.SetWorkOrderTypeID(*id)
	}
	return ptu
}

// SetWorkOrderType sets the work_order_type edge to WorkOrderType.
func (ptu *PropertyTypeUpdate) SetWorkOrderType(w *WorkOrderType) *PropertyTypeUpdate {
	return ptu.SetWorkOrderTypeID(w.ID)
}

// SetProjectTypeID sets the project_type edge to ProjectType by id.
func (ptu *PropertyTypeUpdate) SetProjectTypeID(id string) *PropertyTypeUpdate {
	if ptu.project_type == nil {
		ptu.project_type = make(map[string]struct{})
	}
	ptu.project_type[id] = struct{}{}
	return ptu
}

// SetNillableProjectTypeID sets the project_type edge to ProjectType by id if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableProjectTypeID(id *string) *PropertyTypeUpdate {
	if id != nil {
		ptu = ptu.SetProjectTypeID(*id)
	}
	return ptu
}

// SetProjectType sets the project_type edge to ProjectType.
func (ptu *PropertyTypeUpdate) SetProjectType(p *ProjectType) *PropertyTypeUpdate {
	return ptu.SetProjectTypeID(p.ID)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (ptu *PropertyTypeUpdate) RemovePropertyIDs(ids ...string) *PropertyTypeUpdate {
	if ptu.removedProperties == nil {
		ptu.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		ptu.removedProperties[ids[i]] = struct{}{}
	}
	return ptu
}

// RemoveProperties removes properties edges to Property.
func (ptu *PropertyTypeUpdate) RemoveProperties(p ...*Property) *PropertyTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptu.RemovePropertyIDs(ids...)
}

// ClearLocationType clears the location_type edge to LocationType.
func (ptu *PropertyTypeUpdate) ClearLocationType() *PropertyTypeUpdate {
	ptu.clearedLocationType = true
	return ptu
}

// ClearEquipmentPortType clears the equipment_port_type edge to EquipmentPortType.
func (ptu *PropertyTypeUpdate) ClearEquipmentPortType() *PropertyTypeUpdate {
	ptu.clearedEquipmentPortType = true
	return ptu
}

// ClearLinkEquipmentPortType clears the link_equipment_port_type edge to EquipmentPortType.
func (ptu *PropertyTypeUpdate) ClearLinkEquipmentPortType() *PropertyTypeUpdate {
	ptu.clearedLinkEquipmentPortType = true
	return ptu
}

// ClearEquipmentType clears the equipment_type edge to EquipmentType.
func (ptu *PropertyTypeUpdate) ClearEquipmentType() *PropertyTypeUpdate {
	ptu.clearedEquipmentType = true
	return ptu
}

// ClearServiceType clears the service_type edge to ServiceType.
func (ptu *PropertyTypeUpdate) ClearServiceType() *PropertyTypeUpdate {
	ptu.clearedServiceType = true
	return ptu
}

// ClearWorkOrderType clears the work_order_type edge to WorkOrderType.
func (ptu *PropertyTypeUpdate) ClearWorkOrderType() *PropertyTypeUpdate {
	ptu.clearedWorkOrderType = true
	return ptu
}

// ClearProjectType clears the project_type edge to ProjectType.
func (ptu *PropertyTypeUpdate) ClearProjectType() *PropertyTypeUpdate {
	ptu.clearedProjectType = true
	return ptu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (ptu *PropertyTypeUpdate) Save(ctx context.Context) (int, error) {
	if ptu.update_time == nil {
		v := propertytype.UpdateDefaultUpdateTime()
		ptu.update_time = &v
	}
	if len(ptu.location_type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"location_type\"")
	}
	if len(ptu.equipment_port_type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"equipment_port_type\"")
	}
	if len(ptu.link_equipment_port_type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"link_equipment_port_type\"")
	}
	if len(ptu.equipment_type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"equipment_type\"")
	}
	if len(ptu.service_type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"service_type\"")
	}
	if len(ptu.work_order_type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"work_order_type\"")
	}
	if len(ptu.project_type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"project_type\"")
	}
	return ptu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (ptu *PropertyTypeUpdate) SaveX(ctx context.Context) int {
	affected, err := ptu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ptu *PropertyTypeUpdate) Exec(ctx context.Context) error {
	_, err := ptu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ptu *PropertyTypeUpdate) ExecX(ctx context.Context) {
	if err := ptu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ptu *PropertyTypeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(ptu.driver.Dialect())
		selector = builder.Select(propertytype.FieldID).From(builder.Table(propertytype.Table))
	)
	for _, p := range ptu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = ptu.driver.Query(ctx, query, args, rows); err != nil {
		return 0, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ent: failed reading id: %v", err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return 0, nil
	}

	tx, err := ptu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(propertytype.Table)
	)
	updater = updater.Where(sql.InInts(propertytype.FieldID, ids...))
	if value := ptu.update_time; value != nil {
		updater.Set(propertytype.FieldUpdateTime, *value)
	}
	if value := ptu._type; value != nil {
		updater.Set(propertytype.FieldType, *value)
	}
	if value := ptu.name; value != nil {
		updater.Set(propertytype.FieldName, *value)
	}
	if value := ptu.index; value != nil {
		updater.Set(propertytype.FieldIndex, *value)
	}
	if value := ptu.addindex; value != nil {
		updater.Add(propertytype.FieldIndex, *value)
	}
	if ptu.clearindex {
		updater.SetNull(propertytype.FieldIndex)
	}
	if value := ptu.category; value != nil {
		updater.Set(propertytype.FieldCategory, *value)
	}
	if ptu.clearcategory {
		updater.SetNull(propertytype.FieldCategory)
	}
	if value := ptu.int_val; value != nil {
		updater.Set(propertytype.FieldIntVal, *value)
	}
	if value := ptu.addint_val; value != nil {
		updater.Add(propertytype.FieldIntVal, *value)
	}
	if ptu.clearint_val {
		updater.SetNull(propertytype.FieldIntVal)
	}
	if value := ptu.bool_val; value != nil {
		updater.Set(propertytype.FieldBoolVal, *value)
	}
	if ptu.clearbool_val {
		updater.SetNull(propertytype.FieldBoolVal)
	}
	if value := ptu.float_val; value != nil {
		updater.Set(propertytype.FieldFloatVal, *value)
	}
	if value := ptu.addfloat_val; value != nil {
		updater.Add(propertytype.FieldFloatVal, *value)
	}
	if ptu.clearfloat_val {
		updater.SetNull(propertytype.FieldFloatVal)
	}
	if value := ptu.latitude_val; value != nil {
		updater.Set(propertytype.FieldLatitudeVal, *value)
	}
	if value := ptu.addlatitude_val; value != nil {
		updater.Add(propertytype.FieldLatitudeVal, *value)
	}
	if ptu.clearlatitude_val {
		updater.SetNull(propertytype.FieldLatitudeVal)
	}
	if value := ptu.longitude_val; value != nil {
		updater.Set(propertytype.FieldLongitudeVal, *value)
	}
	if value := ptu.addlongitude_val; value != nil {
		updater.Add(propertytype.FieldLongitudeVal, *value)
	}
	if ptu.clearlongitude_val {
		updater.SetNull(propertytype.FieldLongitudeVal)
	}
	if value := ptu.string_val; value != nil {
		updater.Set(propertytype.FieldStringVal, *value)
	}
	if ptu.clearstring_val {
		updater.SetNull(propertytype.FieldStringVal)
	}
	if value := ptu.range_from_val; value != nil {
		updater.Set(propertytype.FieldRangeFromVal, *value)
	}
	if value := ptu.addrange_from_val; value != nil {
		updater.Add(propertytype.FieldRangeFromVal, *value)
	}
	if ptu.clearrange_from_val {
		updater.SetNull(propertytype.FieldRangeFromVal)
	}
	if value := ptu.range_to_val; value != nil {
		updater.Set(propertytype.FieldRangeToVal, *value)
	}
	if value := ptu.addrange_to_val; value != nil {
		updater.Add(propertytype.FieldRangeToVal, *value)
	}
	if ptu.clearrange_to_val {
		updater.SetNull(propertytype.FieldRangeToVal)
	}
	if value := ptu.is_instance_property; value != nil {
		updater.Set(propertytype.FieldIsInstanceProperty, *value)
	}
	if value := ptu.editable; value != nil {
		updater.Set(propertytype.FieldEditable, *value)
	}
	if value := ptu.mandatory; value != nil {
		updater.Set(propertytype.FieldMandatory, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ptu.removedProperties) > 0 {
		eids := make([]int, len(ptu.removedProperties))
		for eid := range ptu.removedProperties {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(propertytype.PropertiesTable).
			SetNull(propertytype.PropertiesColumn).
			Where(sql.InInts(propertytype.PropertiesColumn, ids...)).
			Where(sql.InInts(property.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ptu.properties) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range ptu.properties {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(property.FieldID, eid)
			}
			query, args := builder.Update(propertytype.PropertiesTable).
				Set(propertytype.PropertiesColumn, id).
				Where(sql.And(p, sql.IsNull(propertytype.PropertiesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(ptu.properties) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"PropertyType\"", keys(ptu.properties))})
			}
		}
	}
	if ptu.clearedLocationType {
		query, args := builder.Update(propertytype.LocationTypeTable).
			SetNull(propertytype.LocationTypeColumn).
			Where(sql.InInts(locationtype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ptu.location_type) > 0 {
		for eid := range ptu.location_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(propertytype.LocationTypeTable).
				Set(propertytype.LocationTypeColumn, eid).
				Where(sql.InInts(propertytype.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if ptu.clearedEquipmentPortType {
		query, args := builder.Update(propertytype.EquipmentPortTypeTable).
			SetNull(propertytype.EquipmentPortTypeColumn).
			Where(sql.InInts(equipmentporttype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ptu.equipment_port_type) > 0 {
		for eid := range ptu.equipment_port_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(propertytype.EquipmentPortTypeTable).
				Set(propertytype.EquipmentPortTypeColumn, eid).
				Where(sql.InInts(propertytype.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if ptu.clearedLinkEquipmentPortType {
		query, args := builder.Update(propertytype.LinkEquipmentPortTypeTable).
			SetNull(propertytype.LinkEquipmentPortTypeColumn).
			Where(sql.InInts(equipmentporttype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ptu.link_equipment_port_type) > 0 {
		for eid := range ptu.link_equipment_port_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(propertytype.LinkEquipmentPortTypeTable).
				Set(propertytype.LinkEquipmentPortTypeColumn, eid).
				Where(sql.InInts(propertytype.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if ptu.clearedEquipmentType {
		query, args := builder.Update(propertytype.EquipmentTypeTable).
			SetNull(propertytype.EquipmentTypeColumn).
			Where(sql.InInts(equipmenttype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ptu.equipment_type) > 0 {
		for eid := range ptu.equipment_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(propertytype.EquipmentTypeTable).
				Set(propertytype.EquipmentTypeColumn, eid).
				Where(sql.InInts(propertytype.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if ptu.clearedServiceType {
		query, args := builder.Update(propertytype.ServiceTypeTable).
			SetNull(propertytype.ServiceTypeColumn).
			Where(sql.InInts(servicetype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ptu.service_type) > 0 {
		for eid := range ptu.service_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(propertytype.ServiceTypeTable).
				Set(propertytype.ServiceTypeColumn, eid).
				Where(sql.InInts(propertytype.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if ptu.clearedWorkOrderType {
		query, args := builder.Update(propertytype.WorkOrderTypeTable).
			SetNull(propertytype.WorkOrderTypeColumn).
			Where(sql.InInts(workordertype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ptu.work_order_type) > 0 {
		for eid := range ptu.work_order_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(propertytype.WorkOrderTypeTable).
				Set(propertytype.WorkOrderTypeColumn, eid).
				Where(sql.InInts(propertytype.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if ptu.clearedProjectType {
		query, args := builder.Update(propertytype.ProjectTypeTable).
			SetNull(propertytype.ProjectTypeColumn).
			Where(sql.InInts(projecttype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ptu.project_type) > 0 {
		for eid := range ptu.project_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(propertytype.ProjectTypeTable).
				Set(propertytype.ProjectTypeColumn, eid).
				Where(sql.InInts(propertytype.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// PropertyTypeUpdateOne is the builder for updating a single PropertyType entity.
type PropertyTypeUpdateOne struct {
	config
	id string

	update_time                  *time.Time
	_type                        *string
	name                         *string
	index                        *int
	addindex                     *int
	clearindex                   bool
	category                     *string
	clearcategory                bool
	int_val                      *int
	addint_val                   *int
	clearint_val                 bool
	bool_val                     *bool
	clearbool_val                bool
	float_val                    *float64
	addfloat_val                 *float64
	clearfloat_val               bool
	latitude_val                 *float64
	addlatitude_val              *float64
	clearlatitude_val            bool
	longitude_val                *float64
	addlongitude_val             *float64
	clearlongitude_val           bool
	string_val                   *string
	clearstring_val              bool
	range_from_val               *float64
	addrange_from_val            *float64
	clearrange_from_val          bool
	range_to_val                 *float64
	addrange_to_val              *float64
	clearrange_to_val            bool
	is_instance_property         *bool
	editable                     *bool
	mandatory                    *bool
	properties                   map[string]struct{}
	location_type                map[string]struct{}
	equipment_port_type          map[string]struct{}
	link_equipment_port_type     map[string]struct{}
	equipment_type               map[string]struct{}
	service_type                 map[string]struct{}
	work_order_type              map[string]struct{}
	project_type                 map[string]struct{}
	removedProperties            map[string]struct{}
	clearedLocationType          bool
	clearedEquipmentPortType     bool
	clearedLinkEquipmentPortType bool
	clearedEquipmentType         bool
	clearedServiceType           bool
	clearedWorkOrderType         bool
	clearedProjectType           bool
}

// SetType sets the type field.
func (ptuo *PropertyTypeUpdateOne) SetType(s string) *PropertyTypeUpdateOne {
	ptuo._type = &s
	return ptuo
}

// SetName sets the name field.
func (ptuo *PropertyTypeUpdateOne) SetName(s string) *PropertyTypeUpdateOne {
	ptuo.name = &s
	return ptuo
}

// SetIndex sets the index field.
func (ptuo *PropertyTypeUpdateOne) SetIndex(i int) *PropertyTypeUpdateOne {
	ptuo.index = &i
	ptuo.addindex = nil
	return ptuo
}

// SetNillableIndex sets the index field if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableIndex(i *int) *PropertyTypeUpdateOne {
	if i != nil {
		ptuo.SetIndex(*i)
	}
	return ptuo
}

// AddIndex adds i to index.
func (ptuo *PropertyTypeUpdateOne) AddIndex(i int) *PropertyTypeUpdateOne {
	if ptuo.addindex == nil {
		ptuo.addindex = &i
	} else {
		*ptuo.addindex += i
	}
	return ptuo
}

// ClearIndex clears the value of index.
func (ptuo *PropertyTypeUpdateOne) ClearIndex() *PropertyTypeUpdateOne {
	ptuo.index = nil
	ptuo.clearindex = true
	return ptuo
}

// SetCategory sets the category field.
func (ptuo *PropertyTypeUpdateOne) SetCategory(s string) *PropertyTypeUpdateOne {
	ptuo.category = &s
	return ptuo
}

// SetNillableCategory sets the category field if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableCategory(s *string) *PropertyTypeUpdateOne {
	if s != nil {
		ptuo.SetCategory(*s)
	}
	return ptuo
}

// ClearCategory clears the value of category.
func (ptuo *PropertyTypeUpdateOne) ClearCategory() *PropertyTypeUpdateOne {
	ptuo.category = nil
	ptuo.clearcategory = true
	return ptuo
}

// SetIntVal sets the int_val field.
func (ptuo *PropertyTypeUpdateOne) SetIntVal(i int) *PropertyTypeUpdateOne {
	ptuo.int_val = &i
	ptuo.addint_val = nil
	return ptuo
}

// SetNillableIntVal sets the int_val field if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableIntVal(i *int) *PropertyTypeUpdateOne {
	if i != nil {
		ptuo.SetIntVal(*i)
	}
	return ptuo
}

// AddIntVal adds i to int_val.
func (ptuo *PropertyTypeUpdateOne) AddIntVal(i int) *PropertyTypeUpdateOne {
	if ptuo.addint_val == nil {
		ptuo.addint_val = &i
	} else {
		*ptuo.addint_val += i
	}
	return ptuo
}

// ClearIntVal clears the value of int_val.
func (ptuo *PropertyTypeUpdateOne) ClearIntVal() *PropertyTypeUpdateOne {
	ptuo.int_val = nil
	ptuo.clearint_val = true
	return ptuo
}

// SetBoolVal sets the bool_val field.
func (ptuo *PropertyTypeUpdateOne) SetBoolVal(b bool) *PropertyTypeUpdateOne {
	ptuo.bool_val = &b
	return ptuo
}

// SetNillableBoolVal sets the bool_val field if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableBoolVal(b *bool) *PropertyTypeUpdateOne {
	if b != nil {
		ptuo.SetBoolVal(*b)
	}
	return ptuo
}

// ClearBoolVal clears the value of bool_val.
func (ptuo *PropertyTypeUpdateOne) ClearBoolVal() *PropertyTypeUpdateOne {
	ptuo.bool_val = nil
	ptuo.clearbool_val = true
	return ptuo
}

// SetFloatVal sets the float_val field.
func (ptuo *PropertyTypeUpdateOne) SetFloatVal(f float64) *PropertyTypeUpdateOne {
	ptuo.float_val = &f
	ptuo.addfloat_val = nil
	return ptuo
}

// SetNillableFloatVal sets the float_val field if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableFloatVal(f *float64) *PropertyTypeUpdateOne {
	if f != nil {
		ptuo.SetFloatVal(*f)
	}
	return ptuo
}

// AddFloatVal adds f to float_val.
func (ptuo *PropertyTypeUpdateOne) AddFloatVal(f float64) *PropertyTypeUpdateOne {
	if ptuo.addfloat_val == nil {
		ptuo.addfloat_val = &f
	} else {
		*ptuo.addfloat_val += f
	}
	return ptuo
}

// ClearFloatVal clears the value of float_val.
func (ptuo *PropertyTypeUpdateOne) ClearFloatVal() *PropertyTypeUpdateOne {
	ptuo.float_val = nil
	ptuo.clearfloat_val = true
	return ptuo
}

// SetLatitudeVal sets the latitude_val field.
func (ptuo *PropertyTypeUpdateOne) SetLatitudeVal(f float64) *PropertyTypeUpdateOne {
	ptuo.latitude_val = &f
	ptuo.addlatitude_val = nil
	return ptuo
}

// SetNillableLatitudeVal sets the latitude_val field if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableLatitudeVal(f *float64) *PropertyTypeUpdateOne {
	if f != nil {
		ptuo.SetLatitudeVal(*f)
	}
	return ptuo
}

// AddLatitudeVal adds f to latitude_val.
func (ptuo *PropertyTypeUpdateOne) AddLatitudeVal(f float64) *PropertyTypeUpdateOne {
	if ptuo.addlatitude_val == nil {
		ptuo.addlatitude_val = &f
	} else {
		*ptuo.addlatitude_val += f
	}
	return ptuo
}

// ClearLatitudeVal clears the value of latitude_val.
func (ptuo *PropertyTypeUpdateOne) ClearLatitudeVal() *PropertyTypeUpdateOne {
	ptuo.latitude_val = nil
	ptuo.clearlatitude_val = true
	return ptuo
}

// SetLongitudeVal sets the longitude_val field.
func (ptuo *PropertyTypeUpdateOne) SetLongitudeVal(f float64) *PropertyTypeUpdateOne {
	ptuo.longitude_val = &f
	ptuo.addlongitude_val = nil
	return ptuo
}

// SetNillableLongitudeVal sets the longitude_val field if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableLongitudeVal(f *float64) *PropertyTypeUpdateOne {
	if f != nil {
		ptuo.SetLongitudeVal(*f)
	}
	return ptuo
}

// AddLongitudeVal adds f to longitude_val.
func (ptuo *PropertyTypeUpdateOne) AddLongitudeVal(f float64) *PropertyTypeUpdateOne {
	if ptuo.addlongitude_val == nil {
		ptuo.addlongitude_val = &f
	} else {
		*ptuo.addlongitude_val += f
	}
	return ptuo
}

// ClearLongitudeVal clears the value of longitude_val.
func (ptuo *PropertyTypeUpdateOne) ClearLongitudeVal() *PropertyTypeUpdateOne {
	ptuo.longitude_val = nil
	ptuo.clearlongitude_val = true
	return ptuo
}

// SetStringVal sets the string_val field.
func (ptuo *PropertyTypeUpdateOne) SetStringVal(s string) *PropertyTypeUpdateOne {
	ptuo.string_val = &s
	return ptuo
}

// SetNillableStringVal sets the string_val field if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableStringVal(s *string) *PropertyTypeUpdateOne {
	if s != nil {
		ptuo.SetStringVal(*s)
	}
	return ptuo
}

// ClearStringVal clears the value of string_val.
func (ptuo *PropertyTypeUpdateOne) ClearStringVal() *PropertyTypeUpdateOne {
	ptuo.string_val = nil
	ptuo.clearstring_val = true
	return ptuo
}

// SetRangeFromVal sets the range_from_val field.
func (ptuo *PropertyTypeUpdateOne) SetRangeFromVal(f float64) *PropertyTypeUpdateOne {
	ptuo.range_from_val = &f
	ptuo.addrange_from_val = nil
	return ptuo
}

// SetNillableRangeFromVal sets the range_from_val field if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableRangeFromVal(f *float64) *PropertyTypeUpdateOne {
	if f != nil {
		ptuo.SetRangeFromVal(*f)
	}
	return ptuo
}

// AddRangeFromVal adds f to range_from_val.
func (ptuo *PropertyTypeUpdateOne) AddRangeFromVal(f float64) *PropertyTypeUpdateOne {
	if ptuo.addrange_from_val == nil {
		ptuo.addrange_from_val = &f
	} else {
		*ptuo.addrange_from_val += f
	}
	return ptuo
}

// ClearRangeFromVal clears the value of range_from_val.
func (ptuo *PropertyTypeUpdateOne) ClearRangeFromVal() *PropertyTypeUpdateOne {
	ptuo.range_from_val = nil
	ptuo.clearrange_from_val = true
	return ptuo
}

// SetRangeToVal sets the range_to_val field.
func (ptuo *PropertyTypeUpdateOne) SetRangeToVal(f float64) *PropertyTypeUpdateOne {
	ptuo.range_to_val = &f
	ptuo.addrange_to_val = nil
	return ptuo
}

// SetNillableRangeToVal sets the range_to_val field if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableRangeToVal(f *float64) *PropertyTypeUpdateOne {
	if f != nil {
		ptuo.SetRangeToVal(*f)
	}
	return ptuo
}

// AddRangeToVal adds f to range_to_val.
func (ptuo *PropertyTypeUpdateOne) AddRangeToVal(f float64) *PropertyTypeUpdateOne {
	if ptuo.addrange_to_val == nil {
		ptuo.addrange_to_val = &f
	} else {
		*ptuo.addrange_to_val += f
	}
	return ptuo
}

// ClearRangeToVal clears the value of range_to_val.
func (ptuo *PropertyTypeUpdateOne) ClearRangeToVal() *PropertyTypeUpdateOne {
	ptuo.range_to_val = nil
	ptuo.clearrange_to_val = true
	return ptuo
}

// SetIsInstanceProperty sets the is_instance_property field.
func (ptuo *PropertyTypeUpdateOne) SetIsInstanceProperty(b bool) *PropertyTypeUpdateOne {
	ptuo.is_instance_property = &b
	return ptuo
}

// SetNillableIsInstanceProperty sets the is_instance_property field if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableIsInstanceProperty(b *bool) *PropertyTypeUpdateOne {
	if b != nil {
		ptuo.SetIsInstanceProperty(*b)
	}
	return ptuo
}

// SetEditable sets the editable field.
func (ptuo *PropertyTypeUpdateOne) SetEditable(b bool) *PropertyTypeUpdateOne {
	ptuo.editable = &b
	return ptuo
}

// SetNillableEditable sets the editable field if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableEditable(b *bool) *PropertyTypeUpdateOne {
	if b != nil {
		ptuo.SetEditable(*b)
	}
	return ptuo
}

// SetMandatory sets the mandatory field.
func (ptuo *PropertyTypeUpdateOne) SetMandatory(b bool) *PropertyTypeUpdateOne {
	ptuo.mandatory = &b
	return ptuo
}

// SetNillableMandatory sets the mandatory field if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableMandatory(b *bool) *PropertyTypeUpdateOne {
	if b != nil {
		ptuo.SetMandatory(*b)
	}
	return ptuo
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (ptuo *PropertyTypeUpdateOne) AddPropertyIDs(ids ...string) *PropertyTypeUpdateOne {
	if ptuo.properties == nil {
		ptuo.properties = make(map[string]struct{})
	}
	for i := range ids {
		ptuo.properties[ids[i]] = struct{}{}
	}
	return ptuo
}

// AddProperties adds the properties edges to Property.
func (ptuo *PropertyTypeUpdateOne) AddProperties(p ...*Property) *PropertyTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptuo.AddPropertyIDs(ids...)
}

// SetLocationTypeID sets the location_type edge to LocationType by id.
func (ptuo *PropertyTypeUpdateOne) SetLocationTypeID(id string) *PropertyTypeUpdateOne {
	if ptuo.location_type == nil {
		ptuo.location_type = make(map[string]struct{})
	}
	ptuo.location_type[id] = struct{}{}
	return ptuo
}

// SetNillableLocationTypeID sets the location_type edge to LocationType by id if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableLocationTypeID(id *string) *PropertyTypeUpdateOne {
	if id != nil {
		ptuo = ptuo.SetLocationTypeID(*id)
	}
	return ptuo
}

// SetLocationType sets the location_type edge to LocationType.
func (ptuo *PropertyTypeUpdateOne) SetLocationType(l *LocationType) *PropertyTypeUpdateOne {
	return ptuo.SetLocationTypeID(l.ID)
}

// SetEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id.
func (ptuo *PropertyTypeUpdateOne) SetEquipmentPortTypeID(id string) *PropertyTypeUpdateOne {
	if ptuo.equipment_port_type == nil {
		ptuo.equipment_port_type = make(map[string]struct{})
	}
	ptuo.equipment_port_type[id] = struct{}{}
	return ptuo
}

// SetNillableEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableEquipmentPortTypeID(id *string) *PropertyTypeUpdateOne {
	if id != nil {
		ptuo = ptuo.SetEquipmentPortTypeID(*id)
	}
	return ptuo
}

// SetEquipmentPortType sets the equipment_port_type edge to EquipmentPortType.
func (ptuo *PropertyTypeUpdateOne) SetEquipmentPortType(e *EquipmentPortType) *PropertyTypeUpdateOne {
	return ptuo.SetEquipmentPortTypeID(e.ID)
}

// SetLinkEquipmentPortTypeID sets the link_equipment_port_type edge to EquipmentPortType by id.
func (ptuo *PropertyTypeUpdateOne) SetLinkEquipmentPortTypeID(id string) *PropertyTypeUpdateOne {
	if ptuo.link_equipment_port_type == nil {
		ptuo.link_equipment_port_type = make(map[string]struct{})
	}
	ptuo.link_equipment_port_type[id] = struct{}{}
	return ptuo
}

// SetNillableLinkEquipmentPortTypeID sets the link_equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableLinkEquipmentPortTypeID(id *string) *PropertyTypeUpdateOne {
	if id != nil {
		ptuo = ptuo.SetLinkEquipmentPortTypeID(*id)
	}
	return ptuo
}

// SetLinkEquipmentPortType sets the link_equipment_port_type edge to EquipmentPortType.
func (ptuo *PropertyTypeUpdateOne) SetLinkEquipmentPortType(e *EquipmentPortType) *PropertyTypeUpdateOne {
	return ptuo.SetLinkEquipmentPortTypeID(e.ID)
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (ptuo *PropertyTypeUpdateOne) SetEquipmentTypeID(id string) *PropertyTypeUpdateOne {
	if ptuo.equipment_type == nil {
		ptuo.equipment_type = make(map[string]struct{})
	}
	ptuo.equipment_type[id] = struct{}{}
	return ptuo
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableEquipmentTypeID(id *string) *PropertyTypeUpdateOne {
	if id != nil {
		ptuo = ptuo.SetEquipmentTypeID(*id)
	}
	return ptuo
}

// SetEquipmentType sets the equipment_type edge to EquipmentType.
func (ptuo *PropertyTypeUpdateOne) SetEquipmentType(e *EquipmentType) *PropertyTypeUpdateOne {
	return ptuo.SetEquipmentTypeID(e.ID)
}

// SetServiceTypeID sets the service_type edge to ServiceType by id.
func (ptuo *PropertyTypeUpdateOne) SetServiceTypeID(id string) *PropertyTypeUpdateOne {
	if ptuo.service_type == nil {
		ptuo.service_type = make(map[string]struct{})
	}
	ptuo.service_type[id] = struct{}{}
	return ptuo
}

// SetNillableServiceTypeID sets the service_type edge to ServiceType by id if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableServiceTypeID(id *string) *PropertyTypeUpdateOne {
	if id != nil {
		ptuo = ptuo.SetServiceTypeID(*id)
	}
	return ptuo
}

// SetServiceType sets the service_type edge to ServiceType.
func (ptuo *PropertyTypeUpdateOne) SetServiceType(s *ServiceType) *PropertyTypeUpdateOne {
	return ptuo.SetServiceTypeID(s.ID)
}

// SetWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id.
func (ptuo *PropertyTypeUpdateOne) SetWorkOrderTypeID(id string) *PropertyTypeUpdateOne {
	if ptuo.work_order_type == nil {
		ptuo.work_order_type = make(map[string]struct{})
	}
	ptuo.work_order_type[id] = struct{}{}
	return ptuo
}

// SetNillableWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableWorkOrderTypeID(id *string) *PropertyTypeUpdateOne {
	if id != nil {
		ptuo = ptuo.SetWorkOrderTypeID(*id)
	}
	return ptuo
}

// SetWorkOrderType sets the work_order_type edge to WorkOrderType.
func (ptuo *PropertyTypeUpdateOne) SetWorkOrderType(w *WorkOrderType) *PropertyTypeUpdateOne {
	return ptuo.SetWorkOrderTypeID(w.ID)
}

// SetProjectTypeID sets the project_type edge to ProjectType by id.
func (ptuo *PropertyTypeUpdateOne) SetProjectTypeID(id string) *PropertyTypeUpdateOne {
	if ptuo.project_type == nil {
		ptuo.project_type = make(map[string]struct{})
	}
	ptuo.project_type[id] = struct{}{}
	return ptuo
}

// SetNillableProjectTypeID sets the project_type edge to ProjectType by id if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableProjectTypeID(id *string) *PropertyTypeUpdateOne {
	if id != nil {
		ptuo = ptuo.SetProjectTypeID(*id)
	}
	return ptuo
}

// SetProjectType sets the project_type edge to ProjectType.
func (ptuo *PropertyTypeUpdateOne) SetProjectType(p *ProjectType) *PropertyTypeUpdateOne {
	return ptuo.SetProjectTypeID(p.ID)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (ptuo *PropertyTypeUpdateOne) RemovePropertyIDs(ids ...string) *PropertyTypeUpdateOne {
	if ptuo.removedProperties == nil {
		ptuo.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		ptuo.removedProperties[ids[i]] = struct{}{}
	}
	return ptuo
}

// RemoveProperties removes properties edges to Property.
func (ptuo *PropertyTypeUpdateOne) RemoveProperties(p ...*Property) *PropertyTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptuo.RemovePropertyIDs(ids...)
}

// ClearLocationType clears the location_type edge to LocationType.
func (ptuo *PropertyTypeUpdateOne) ClearLocationType() *PropertyTypeUpdateOne {
	ptuo.clearedLocationType = true
	return ptuo
}

// ClearEquipmentPortType clears the equipment_port_type edge to EquipmentPortType.
func (ptuo *PropertyTypeUpdateOne) ClearEquipmentPortType() *PropertyTypeUpdateOne {
	ptuo.clearedEquipmentPortType = true
	return ptuo
}

// ClearLinkEquipmentPortType clears the link_equipment_port_type edge to EquipmentPortType.
func (ptuo *PropertyTypeUpdateOne) ClearLinkEquipmentPortType() *PropertyTypeUpdateOne {
	ptuo.clearedLinkEquipmentPortType = true
	return ptuo
}

// ClearEquipmentType clears the equipment_type edge to EquipmentType.
func (ptuo *PropertyTypeUpdateOne) ClearEquipmentType() *PropertyTypeUpdateOne {
	ptuo.clearedEquipmentType = true
	return ptuo
}

// ClearServiceType clears the service_type edge to ServiceType.
func (ptuo *PropertyTypeUpdateOne) ClearServiceType() *PropertyTypeUpdateOne {
	ptuo.clearedServiceType = true
	return ptuo
}

// ClearWorkOrderType clears the work_order_type edge to WorkOrderType.
func (ptuo *PropertyTypeUpdateOne) ClearWorkOrderType() *PropertyTypeUpdateOne {
	ptuo.clearedWorkOrderType = true
	return ptuo
}

// ClearProjectType clears the project_type edge to ProjectType.
func (ptuo *PropertyTypeUpdateOne) ClearProjectType() *PropertyTypeUpdateOne {
	ptuo.clearedProjectType = true
	return ptuo
}

// Save executes the query and returns the updated entity.
func (ptuo *PropertyTypeUpdateOne) Save(ctx context.Context) (*PropertyType, error) {
	if ptuo.update_time == nil {
		v := propertytype.UpdateDefaultUpdateTime()
		ptuo.update_time = &v
	}
	if len(ptuo.location_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location_type\"")
	}
	if len(ptuo.equipment_port_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"equipment_port_type\"")
	}
	if len(ptuo.link_equipment_port_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"link_equipment_port_type\"")
	}
	if len(ptuo.equipment_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"equipment_type\"")
	}
	if len(ptuo.service_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"service_type\"")
	}
	if len(ptuo.work_order_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"work_order_type\"")
	}
	if len(ptuo.project_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"project_type\"")
	}
	return ptuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (ptuo *PropertyTypeUpdateOne) SaveX(ctx context.Context) *PropertyType {
	pt, err := ptuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return pt
}

// Exec executes the query on the entity.
func (ptuo *PropertyTypeUpdateOne) Exec(ctx context.Context) error {
	_, err := ptuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ptuo *PropertyTypeUpdateOne) ExecX(ctx context.Context) {
	if err := ptuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ptuo *PropertyTypeUpdateOne) sqlSave(ctx context.Context) (pt *PropertyType, err error) {
	var (
		builder  = sql.Dialect(ptuo.driver.Dialect())
		selector = builder.Select(propertytype.Columns...).From(builder.Table(propertytype.Table))
	)
	propertytype.ID(ptuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = ptuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		pt = &PropertyType{config: ptuo.config}
		if err := pt.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into PropertyType: %v", err)
		}
		id = pt.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("PropertyType with id: %v", ptuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one PropertyType with the same id: %v", ptuo.id)
	}

	tx, err := ptuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(propertytype.Table)
	)
	updater = updater.Where(sql.InInts(propertytype.FieldID, ids...))
	if value := ptuo.update_time; value != nil {
		updater.Set(propertytype.FieldUpdateTime, *value)
		pt.UpdateTime = *value
	}
	if value := ptuo._type; value != nil {
		updater.Set(propertytype.FieldType, *value)
		pt.Type = *value
	}
	if value := ptuo.name; value != nil {
		updater.Set(propertytype.FieldName, *value)
		pt.Name = *value
	}
	if value := ptuo.index; value != nil {
		updater.Set(propertytype.FieldIndex, *value)
		pt.Index = *value
	}
	if value := ptuo.addindex; value != nil {
		updater.Add(propertytype.FieldIndex, *value)
		pt.Index += *value
	}
	if ptuo.clearindex {
		var value int
		pt.Index = value
		updater.SetNull(propertytype.FieldIndex)
	}
	if value := ptuo.category; value != nil {
		updater.Set(propertytype.FieldCategory, *value)
		pt.Category = *value
	}
	if ptuo.clearcategory {
		var value string
		pt.Category = value
		updater.SetNull(propertytype.FieldCategory)
	}
	if value := ptuo.int_val; value != nil {
		updater.Set(propertytype.FieldIntVal, *value)
		pt.IntVal = *value
	}
	if value := ptuo.addint_val; value != nil {
		updater.Add(propertytype.FieldIntVal, *value)
		pt.IntVal += *value
	}
	if ptuo.clearint_val {
		var value int
		pt.IntVal = value
		updater.SetNull(propertytype.FieldIntVal)
	}
	if value := ptuo.bool_val; value != nil {
		updater.Set(propertytype.FieldBoolVal, *value)
		pt.BoolVal = *value
	}
	if ptuo.clearbool_val {
		var value bool
		pt.BoolVal = value
		updater.SetNull(propertytype.FieldBoolVal)
	}
	if value := ptuo.float_val; value != nil {
		updater.Set(propertytype.FieldFloatVal, *value)
		pt.FloatVal = *value
	}
	if value := ptuo.addfloat_val; value != nil {
		updater.Add(propertytype.FieldFloatVal, *value)
		pt.FloatVal += *value
	}
	if ptuo.clearfloat_val {
		var value float64
		pt.FloatVal = value
		updater.SetNull(propertytype.FieldFloatVal)
	}
	if value := ptuo.latitude_val; value != nil {
		updater.Set(propertytype.FieldLatitudeVal, *value)
		pt.LatitudeVal = *value
	}
	if value := ptuo.addlatitude_val; value != nil {
		updater.Add(propertytype.FieldLatitudeVal, *value)
		pt.LatitudeVal += *value
	}
	if ptuo.clearlatitude_val {
		var value float64
		pt.LatitudeVal = value
		updater.SetNull(propertytype.FieldLatitudeVal)
	}
	if value := ptuo.longitude_val; value != nil {
		updater.Set(propertytype.FieldLongitudeVal, *value)
		pt.LongitudeVal = *value
	}
	if value := ptuo.addlongitude_val; value != nil {
		updater.Add(propertytype.FieldLongitudeVal, *value)
		pt.LongitudeVal += *value
	}
	if ptuo.clearlongitude_val {
		var value float64
		pt.LongitudeVal = value
		updater.SetNull(propertytype.FieldLongitudeVal)
	}
	if value := ptuo.string_val; value != nil {
		updater.Set(propertytype.FieldStringVal, *value)
		pt.StringVal = *value
	}
	if ptuo.clearstring_val {
		var value string
		pt.StringVal = value
		updater.SetNull(propertytype.FieldStringVal)
	}
	if value := ptuo.range_from_val; value != nil {
		updater.Set(propertytype.FieldRangeFromVal, *value)
		pt.RangeFromVal = *value
	}
	if value := ptuo.addrange_from_val; value != nil {
		updater.Add(propertytype.FieldRangeFromVal, *value)
		pt.RangeFromVal += *value
	}
	if ptuo.clearrange_from_val {
		var value float64
		pt.RangeFromVal = value
		updater.SetNull(propertytype.FieldRangeFromVal)
	}
	if value := ptuo.range_to_val; value != nil {
		updater.Set(propertytype.FieldRangeToVal, *value)
		pt.RangeToVal = *value
	}
	if value := ptuo.addrange_to_val; value != nil {
		updater.Add(propertytype.FieldRangeToVal, *value)
		pt.RangeToVal += *value
	}
	if ptuo.clearrange_to_val {
		var value float64
		pt.RangeToVal = value
		updater.SetNull(propertytype.FieldRangeToVal)
	}
	if value := ptuo.is_instance_property; value != nil {
		updater.Set(propertytype.FieldIsInstanceProperty, *value)
		pt.IsInstanceProperty = *value
	}
	if value := ptuo.editable; value != nil {
		updater.Set(propertytype.FieldEditable, *value)
		pt.Editable = *value
	}
	if value := ptuo.mandatory; value != nil {
		updater.Set(propertytype.FieldMandatory, *value)
		pt.Mandatory = *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ptuo.removedProperties) > 0 {
		eids := make([]int, len(ptuo.removedProperties))
		for eid := range ptuo.removedProperties {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(propertytype.PropertiesTable).
			SetNull(propertytype.PropertiesColumn).
			Where(sql.InInts(propertytype.PropertiesColumn, ids...)).
			Where(sql.InInts(property.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ptuo.properties) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range ptuo.properties {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(property.FieldID, eid)
			}
			query, args := builder.Update(propertytype.PropertiesTable).
				Set(propertytype.PropertiesColumn, id).
				Where(sql.And(p, sql.IsNull(propertytype.PropertiesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(ptuo.properties) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"PropertyType\"", keys(ptuo.properties))})
			}
		}
	}
	if ptuo.clearedLocationType {
		query, args := builder.Update(propertytype.LocationTypeTable).
			SetNull(propertytype.LocationTypeColumn).
			Where(sql.InInts(locationtype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ptuo.location_type) > 0 {
		for eid := range ptuo.location_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(propertytype.LocationTypeTable).
				Set(propertytype.LocationTypeColumn, eid).
				Where(sql.InInts(propertytype.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if ptuo.clearedEquipmentPortType {
		query, args := builder.Update(propertytype.EquipmentPortTypeTable).
			SetNull(propertytype.EquipmentPortTypeColumn).
			Where(sql.InInts(equipmentporttype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ptuo.equipment_port_type) > 0 {
		for eid := range ptuo.equipment_port_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(propertytype.EquipmentPortTypeTable).
				Set(propertytype.EquipmentPortTypeColumn, eid).
				Where(sql.InInts(propertytype.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if ptuo.clearedLinkEquipmentPortType {
		query, args := builder.Update(propertytype.LinkEquipmentPortTypeTable).
			SetNull(propertytype.LinkEquipmentPortTypeColumn).
			Where(sql.InInts(equipmentporttype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ptuo.link_equipment_port_type) > 0 {
		for eid := range ptuo.link_equipment_port_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(propertytype.LinkEquipmentPortTypeTable).
				Set(propertytype.LinkEquipmentPortTypeColumn, eid).
				Where(sql.InInts(propertytype.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if ptuo.clearedEquipmentType {
		query, args := builder.Update(propertytype.EquipmentTypeTable).
			SetNull(propertytype.EquipmentTypeColumn).
			Where(sql.InInts(equipmenttype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ptuo.equipment_type) > 0 {
		for eid := range ptuo.equipment_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(propertytype.EquipmentTypeTable).
				Set(propertytype.EquipmentTypeColumn, eid).
				Where(sql.InInts(propertytype.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if ptuo.clearedServiceType {
		query, args := builder.Update(propertytype.ServiceTypeTable).
			SetNull(propertytype.ServiceTypeColumn).
			Where(sql.InInts(servicetype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ptuo.service_type) > 0 {
		for eid := range ptuo.service_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(propertytype.ServiceTypeTable).
				Set(propertytype.ServiceTypeColumn, eid).
				Where(sql.InInts(propertytype.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if ptuo.clearedWorkOrderType {
		query, args := builder.Update(propertytype.WorkOrderTypeTable).
			SetNull(propertytype.WorkOrderTypeColumn).
			Where(sql.InInts(workordertype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ptuo.work_order_type) > 0 {
		for eid := range ptuo.work_order_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(propertytype.WorkOrderTypeTable).
				Set(propertytype.WorkOrderTypeColumn, eid).
				Where(sql.InInts(propertytype.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if ptuo.clearedProjectType {
		query, args := builder.Update(propertytype.ProjectTypeTable).
			SetNull(propertytype.ProjectTypeColumn).
			Where(sql.InInts(projecttype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ptuo.project_type) > 0 {
		for eid := range ptuo.project_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(propertytype.ProjectTypeTable).
				Set(propertytype.ProjectTypeColumn, eid).
				Where(sql.InInts(propertytype.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return pt, nil
}
