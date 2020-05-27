// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
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
	hooks      []Hook
	mutation   *PropertyTypeMutation
	predicates []predicate.PropertyType
}

// Where adds a new predicate for the builder.
func (ptu *PropertyTypeUpdate) Where(ps ...predicate.PropertyType) *PropertyTypeUpdate {
	ptu.predicates = append(ptu.predicates, ps...)
	return ptu
}

// SetType sets the type field.
func (ptu *PropertyTypeUpdate) SetType(s string) *PropertyTypeUpdate {
	ptu.mutation.SetType(s)
	return ptu
}

// SetName sets the name field.
func (ptu *PropertyTypeUpdate) SetName(s string) *PropertyTypeUpdate {
	ptu.mutation.SetName(s)
	return ptu
}

// SetExternalID sets the external_id field.
func (ptu *PropertyTypeUpdate) SetExternalID(s string) *PropertyTypeUpdate {
	ptu.mutation.SetExternalID(s)
	return ptu
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableExternalID(s *string) *PropertyTypeUpdate {
	if s != nil {
		ptu.SetExternalID(*s)
	}
	return ptu
}

// ClearExternalID clears the value of external_id.
func (ptu *PropertyTypeUpdate) ClearExternalID() *PropertyTypeUpdate {
	ptu.mutation.ClearExternalID()
	return ptu
}

// SetIndex sets the index field.
func (ptu *PropertyTypeUpdate) SetIndex(i int) *PropertyTypeUpdate {
	ptu.mutation.ResetIndex()
	ptu.mutation.SetIndex(i)
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
	ptu.mutation.AddIndex(i)
	return ptu
}

// ClearIndex clears the value of index.
func (ptu *PropertyTypeUpdate) ClearIndex() *PropertyTypeUpdate {
	ptu.mutation.ClearIndex()
	return ptu
}

// SetCategory sets the category field.
func (ptu *PropertyTypeUpdate) SetCategory(s string) *PropertyTypeUpdate {
	ptu.mutation.SetCategory(s)
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
	ptu.mutation.ClearCategory()
	return ptu
}

// SetIntVal sets the int_val field.
func (ptu *PropertyTypeUpdate) SetIntVal(i int) *PropertyTypeUpdate {
	ptu.mutation.ResetIntVal()
	ptu.mutation.SetIntVal(i)
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
	ptu.mutation.AddIntVal(i)
	return ptu
}

// ClearIntVal clears the value of int_val.
func (ptu *PropertyTypeUpdate) ClearIntVal() *PropertyTypeUpdate {
	ptu.mutation.ClearIntVal()
	return ptu
}

// SetBoolVal sets the bool_val field.
func (ptu *PropertyTypeUpdate) SetBoolVal(b bool) *PropertyTypeUpdate {
	ptu.mutation.SetBoolVal(b)
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
	ptu.mutation.ClearBoolVal()
	return ptu
}

// SetFloatVal sets the float_val field.
func (ptu *PropertyTypeUpdate) SetFloatVal(f float64) *PropertyTypeUpdate {
	ptu.mutation.ResetFloatVal()
	ptu.mutation.SetFloatVal(f)
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
	ptu.mutation.AddFloatVal(f)
	return ptu
}

// ClearFloatVal clears the value of float_val.
func (ptu *PropertyTypeUpdate) ClearFloatVal() *PropertyTypeUpdate {
	ptu.mutation.ClearFloatVal()
	return ptu
}

// SetLatitudeVal sets the latitude_val field.
func (ptu *PropertyTypeUpdate) SetLatitudeVal(f float64) *PropertyTypeUpdate {
	ptu.mutation.ResetLatitudeVal()
	ptu.mutation.SetLatitudeVal(f)
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
	ptu.mutation.AddLatitudeVal(f)
	return ptu
}

// ClearLatitudeVal clears the value of latitude_val.
func (ptu *PropertyTypeUpdate) ClearLatitudeVal() *PropertyTypeUpdate {
	ptu.mutation.ClearLatitudeVal()
	return ptu
}

// SetLongitudeVal sets the longitude_val field.
func (ptu *PropertyTypeUpdate) SetLongitudeVal(f float64) *PropertyTypeUpdate {
	ptu.mutation.ResetLongitudeVal()
	ptu.mutation.SetLongitudeVal(f)
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
	ptu.mutation.AddLongitudeVal(f)
	return ptu
}

// ClearLongitudeVal clears the value of longitude_val.
func (ptu *PropertyTypeUpdate) ClearLongitudeVal() *PropertyTypeUpdate {
	ptu.mutation.ClearLongitudeVal()
	return ptu
}

// SetStringVal sets the string_val field.
func (ptu *PropertyTypeUpdate) SetStringVal(s string) *PropertyTypeUpdate {
	ptu.mutation.SetStringVal(s)
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
	ptu.mutation.ClearStringVal()
	return ptu
}

// SetRangeFromVal sets the range_from_val field.
func (ptu *PropertyTypeUpdate) SetRangeFromVal(f float64) *PropertyTypeUpdate {
	ptu.mutation.ResetRangeFromVal()
	ptu.mutation.SetRangeFromVal(f)
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
	ptu.mutation.AddRangeFromVal(f)
	return ptu
}

// ClearRangeFromVal clears the value of range_from_val.
func (ptu *PropertyTypeUpdate) ClearRangeFromVal() *PropertyTypeUpdate {
	ptu.mutation.ClearRangeFromVal()
	return ptu
}

// SetRangeToVal sets the range_to_val field.
func (ptu *PropertyTypeUpdate) SetRangeToVal(f float64) *PropertyTypeUpdate {
	ptu.mutation.ResetRangeToVal()
	ptu.mutation.SetRangeToVal(f)
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
	ptu.mutation.AddRangeToVal(f)
	return ptu
}

// ClearRangeToVal clears the value of range_to_val.
func (ptu *PropertyTypeUpdate) ClearRangeToVal() *PropertyTypeUpdate {
	ptu.mutation.ClearRangeToVal()
	return ptu
}

// SetIsInstanceProperty sets the is_instance_property field.
func (ptu *PropertyTypeUpdate) SetIsInstanceProperty(b bool) *PropertyTypeUpdate {
	ptu.mutation.SetIsInstanceProperty(b)
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
	ptu.mutation.SetEditable(b)
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
	ptu.mutation.SetMandatory(b)
	return ptu
}

// SetNillableMandatory sets the mandatory field if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableMandatory(b *bool) *PropertyTypeUpdate {
	if b != nil {
		ptu.SetMandatory(*b)
	}
	return ptu
}

// SetDeleted sets the deleted field.
func (ptu *PropertyTypeUpdate) SetDeleted(b bool) *PropertyTypeUpdate {
	ptu.mutation.SetDeleted(b)
	return ptu
}

// SetNillableDeleted sets the deleted field if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableDeleted(b *bool) *PropertyTypeUpdate {
	if b != nil {
		ptu.SetDeleted(*b)
	}
	return ptu
}

// SetNodeType sets the nodeType field.
func (ptu *PropertyTypeUpdate) SetNodeType(s string) *PropertyTypeUpdate {
	ptu.mutation.SetNodeType(s)
	return ptu
}

// SetNillableNodeType sets the nodeType field if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableNodeType(s *string) *PropertyTypeUpdate {
	if s != nil {
		ptu.SetNodeType(*s)
	}
	return ptu
}

// ClearNodeType clears the value of nodeType.
func (ptu *PropertyTypeUpdate) ClearNodeType() *PropertyTypeUpdate {
	ptu.mutation.ClearNodeType()
	return ptu
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (ptu *PropertyTypeUpdate) AddPropertyIDs(ids ...int) *PropertyTypeUpdate {
	ptu.mutation.AddPropertyIDs(ids...)
	return ptu
}

// AddProperties adds the properties edges to Property.
func (ptu *PropertyTypeUpdate) AddProperties(p ...*Property) *PropertyTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptu.AddPropertyIDs(ids...)
}

// SetLocationTypeID sets the location_type edge to LocationType by id.
func (ptu *PropertyTypeUpdate) SetLocationTypeID(id int) *PropertyTypeUpdate {
	ptu.mutation.SetLocationTypeID(id)
	return ptu
}

// SetNillableLocationTypeID sets the location_type edge to LocationType by id if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableLocationTypeID(id *int) *PropertyTypeUpdate {
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
func (ptu *PropertyTypeUpdate) SetEquipmentPortTypeID(id int) *PropertyTypeUpdate {
	ptu.mutation.SetEquipmentPortTypeID(id)
	return ptu
}

// SetNillableEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableEquipmentPortTypeID(id *int) *PropertyTypeUpdate {
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
func (ptu *PropertyTypeUpdate) SetLinkEquipmentPortTypeID(id int) *PropertyTypeUpdate {
	ptu.mutation.SetLinkEquipmentPortTypeID(id)
	return ptu
}

// SetNillableLinkEquipmentPortTypeID sets the link_equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableLinkEquipmentPortTypeID(id *int) *PropertyTypeUpdate {
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
func (ptu *PropertyTypeUpdate) SetEquipmentTypeID(id int) *PropertyTypeUpdate {
	ptu.mutation.SetEquipmentTypeID(id)
	return ptu
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableEquipmentTypeID(id *int) *PropertyTypeUpdate {
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
func (ptu *PropertyTypeUpdate) SetServiceTypeID(id int) *PropertyTypeUpdate {
	ptu.mutation.SetServiceTypeID(id)
	return ptu
}

// SetNillableServiceTypeID sets the service_type edge to ServiceType by id if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableServiceTypeID(id *int) *PropertyTypeUpdate {
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
func (ptu *PropertyTypeUpdate) SetWorkOrderTypeID(id int) *PropertyTypeUpdate {
	ptu.mutation.SetWorkOrderTypeID(id)
	return ptu
}

// SetNillableWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableWorkOrderTypeID(id *int) *PropertyTypeUpdate {
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
func (ptu *PropertyTypeUpdate) SetProjectTypeID(id int) *PropertyTypeUpdate {
	ptu.mutation.SetProjectTypeID(id)
	return ptu
}

// SetNillableProjectTypeID sets the project_type edge to ProjectType by id if the given value is not nil.
func (ptu *PropertyTypeUpdate) SetNillableProjectTypeID(id *int) *PropertyTypeUpdate {
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
func (ptu *PropertyTypeUpdate) RemovePropertyIDs(ids ...int) *PropertyTypeUpdate {
	ptu.mutation.RemovePropertyIDs(ids...)
	return ptu
}

// RemoveProperties removes properties edges to Property.
func (ptu *PropertyTypeUpdate) RemoveProperties(p ...*Property) *PropertyTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptu.RemovePropertyIDs(ids...)
}

// ClearLocationType clears the location_type edge to LocationType.
func (ptu *PropertyTypeUpdate) ClearLocationType() *PropertyTypeUpdate {
	ptu.mutation.ClearLocationType()
	return ptu
}

// ClearEquipmentPortType clears the equipment_port_type edge to EquipmentPortType.
func (ptu *PropertyTypeUpdate) ClearEquipmentPortType() *PropertyTypeUpdate {
	ptu.mutation.ClearEquipmentPortType()
	return ptu
}

// ClearLinkEquipmentPortType clears the link_equipment_port_type edge to EquipmentPortType.
func (ptu *PropertyTypeUpdate) ClearLinkEquipmentPortType() *PropertyTypeUpdate {
	ptu.mutation.ClearLinkEquipmentPortType()
	return ptu
}

// ClearEquipmentType clears the equipment_type edge to EquipmentType.
func (ptu *PropertyTypeUpdate) ClearEquipmentType() *PropertyTypeUpdate {
	ptu.mutation.ClearEquipmentType()
	return ptu
}

// ClearServiceType clears the service_type edge to ServiceType.
func (ptu *PropertyTypeUpdate) ClearServiceType() *PropertyTypeUpdate {
	ptu.mutation.ClearServiceType()
	return ptu
}

// ClearWorkOrderType clears the work_order_type edge to WorkOrderType.
func (ptu *PropertyTypeUpdate) ClearWorkOrderType() *PropertyTypeUpdate {
	ptu.mutation.ClearWorkOrderType()
	return ptu
}

// ClearProjectType clears the project_type edge to ProjectType.
func (ptu *PropertyTypeUpdate) ClearProjectType() *PropertyTypeUpdate {
	ptu.mutation.ClearProjectType()
	return ptu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (ptu *PropertyTypeUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := ptu.mutation.UpdateTime(); !ok {
		v := propertytype.UpdateDefaultUpdateTime()
		ptu.mutation.SetUpdateTime(v)
	}

	var (
		err      error
		affected int
	)
	if len(ptu.hooks) == 0 {
		affected, err = ptu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PropertyTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ptu.mutation = mutation
			affected, err = ptu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(ptu.hooks) - 1; i >= 0; i-- {
			mut = ptu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ptu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   propertytype.Table,
			Columns: propertytype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: propertytype.FieldID,
			},
		},
	}
	if ps := ptu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ptu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: propertytype.FieldUpdateTime,
		})
	}
	if value, ok := ptu.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldType,
		})
	}
	if value, ok := ptu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldName,
		})
	}
	if value, ok := ptu.mutation.ExternalID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldExternalID,
		})
	}
	if ptu.mutation.ExternalIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: propertytype.FieldExternalID,
		})
	}
	if value, ok := ptu.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: propertytype.FieldIndex,
		})
	}
	if value, ok := ptu.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: propertytype.FieldIndex,
		})
	}
	if ptu.mutation.IndexCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: propertytype.FieldIndex,
		})
	}
	if value, ok := ptu.mutation.Category(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldCategory,
		})
	}
	if ptu.mutation.CategoryCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: propertytype.FieldCategory,
		})
	}
	if value, ok := ptu.mutation.IntVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: propertytype.FieldIntVal,
		})
	}
	if value, ok := ptu.mutation.AddedIntVal(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: propertytype.FieldIntVal,
		})
	}
	if ptu.mutation.IntValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: propertytype.FieldIntVal,
		})
	}
	if value, ok := ptu.mutation.BoolVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldBoolVal,
		})
	}
	if ptu.mutation.BoolValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: propertytype.FieldBoolVal,
		})
	}
	if value, ok := ptu.mutation.FloatVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldFloatVal,
		})
	}
	if value, ok := ptu.mutation.AddedFloatVal(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldFloatVal,
		})
	}
	if ptu.mutation.FloatValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: propertytype.FieldFloatVal,
		})
	}
	if value, ok := ptu.mutation.LatitudeVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldLatitudeVal,
		})
	}
	if value, ok := ptu.mutation.AddedLatitudeVal(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldLatitudeVal,
		})
	}
	if ptu.mutation.LatitudeValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: propertytype.FieldLatitudeVal,
		})
	}
	if value, ok := ptu.mutation.LongitudeVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldLongitudeVal,
		})
	}
	if value, ok := ptu.mutation.AddedLongitudeVal(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldLongitudeVal,
		})
	}
	if ptu.mutation.LongitudeValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: propertytype.FieldLongitudeVal,
		})
	}
	if value, ok := ptu.mutation.StringVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldStringVal,
		})
	}
	if ptu.mutation.StringValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: propertytype.FieldStringVal,
		})
	}
	if value, ok := ptu.mutation.RangeFromVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldRangeFromVal,
		})
	}
	if value, ok := ptu.mutation.AddedRangeFromVal(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldRangeFromVal,
		})
	}
	if ptu.mutation.RangeFromValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: propertytype.FieldRangeFromVal,
		})
	}
	if value, ok := ptu.mutation.RangeToVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldRangeToVal,
		})
	}
	if value, ok := ptu.mutation.AddedRangeToVal(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldRangeToVal,
		})
	}
	if ptu.mutation.RangeToValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: propertytype.FieldRangeToVal,
		})
	}
	if value, ok := ptu.mutation.IsInstanceProperty(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldIsInstanceProperty,
		})
	}
	if value, ok := ptu.mutation.Editable(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldEditable,
		})
	}
	if value, ok := ptu.mutation.Mandatory(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldMandatory,
		})
	}
	if value, ok := ptu.mutation.Deleted(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldDeleted,
		})
	}
	if value, ok := ptu.mutation.NodeType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldNodeType,
		})
	}
	if ptu.mutation.NodeTypeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: propertytype.FieldNodeType,
		})
	}
	if nodes := ptu.mutation.RemovedPropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   propertytype.PropertiesTable,
			Columns: []string{propertytype.PropertiesColumn},
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
	if nodes := ptu.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   propertytype.PropertiesTable,
			Columns: []string{propertytype.PropertiesColumn},
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
	if ptu.mutation.LocationTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.LocationTypeTable,
			Columns: []string{propertytype.LocationTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: locationtype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptu.mutation.LocationTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.LocationTypeTable,
			Columns: []string{propertytype.LocationTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: locationtype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ptu.mutation.EquipmentPortTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.EquipmentPortTypeTable,
			Columns: []string{propertytype.EquipmentPortTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentporttype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptu.mutation.EquipmentPortTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.EquipmentPortTypeTable,
			Columns: []string{propertytype.EquipmentPortTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentporttype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ptu.mutation.LinkEquipmentPortTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.LinkEquipmentPortTypeTable,
			Columns: []string{propertytype.LinkEquipmentPortTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentporttype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptu.mutation.LinkEquipmentPortTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.LinkEquipmentPortTypeTable,
			Columns: []string{propertytype.LinkEquipmentPortTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentporttype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ptu.mutation.EquipmentTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.EquipmentTypeTable,
			Columns: []string{propertytype.EquipmentTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmenttype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptu.mutation.EquipmentTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.EquipmentTypeTable,
			Columns: []string{propertytype.EquipmentTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmenttype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ptu.mutation.ServiceTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.ServiceTypeTable,
			Columns: []string{propertytype.ServiceTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: servicetype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptu.mutation.ServiceTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.ServiceTypeTable,
			Columns: []string{propertytype.ServiceTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: servicetype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ptu.mutation.WorkOrderTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.WorkOrderTypeTable,
			Columns: []string{propertytype.WorkOrderTypeColumn},
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
	if nodes := ptu.mutation.WorkOrderTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.WorkOrderTypeTable,
			Columns: []string{propertytype.WorkOrderTypeColumn},
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
	if ptu.mutation.ProjectTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.ProjectTypeTable,
			Columns: []string{propertytype.ProjectTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: projecttype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptu.mutation.ProjectTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.ProjectTypeTable,
			Columns: []string{propertytype.ProjectTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: projecttype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ptu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{propertytype.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// PropertyTypeUpdateOne is the builder for updating a single PropertyType entity.
type PropertyTypeUpdateOne struct {
	config
	hooks    []Hook
	mutation *PropertyTypeMutation
}

// SetType sets the type field.
func (ptuo *PropertyTypeUpdateOne) SetType(s string) *PropertyTypeUpdateOne {
	ptuo.mutation.SetType(s)
	return ptuo
}

// SetName sets the name field.
func (ptuo *PropertyTypeUpdateOne) SetName(s string) *PropertyTypeUpdateOne {
	ptuo.mutation.SetName(s)
	return ptuo
}

// SetExternalID sets the external_id field.
func (ptuo *PropertyTypeUpdateOne) SetExternalID(s string) *PropertyTypeUpdateOne {
	ptuo.mutation.SetExternalID(s)
	return ptuo
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableExternalID(s *string) *PropertyTypeUpdateOne {
	if s != nil {
		ptuo.SetExternalID(*s)
	}
	return ptuo
}

// ClearExternalID clears the value of external_id.
func (ptuo *PropertyTypeUpdateOne) ClearExternalID() *PropertyTypeUpdateOne {
	ptuo.mutation.ClearExternalID()
	return ptuo
}

// SetIndex sets the index field.
func (ptuo *PropertyTypeUpdateOne) SetIndex(i int) *PropertyTypeUpdateOne {
	ptuo.mutation.ResetIndex()
	ptuo.mutation.SetIndex(i)
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
	ptuo.mutation.AddIndex(i)
	return ptuo
}

// ClearIndex clears the value of index.
func (ptuo *PropertyTypeUpdateOne) ClearIndex() *PropertyTypeUpdateOne {
	ptuo.mutation.ClearIndex()
	return ptuo
}

// SetCategory sets the category field.
func (ptuo *PropertyTypeUpdateOne) SetCategory(s string) *PropertyTypeUpdateOne {
	ptuo.mutation.SetCategory(s)
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
	ptuo.mutation.ClearCategory()
	return ptuo
}

// SetIntVal sets the int_val field.
func (ptuo *PropertyTypeUpdateOne) SetIntVal(i int) *PropertyTypeUpdateOne {
	ptuo.mutation.ResetIntVal()
	ptuo.mutation.SetIntVal(i)
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
	ptuo.mutation.AddIntVal(i)
	return ptuo
}

// ClearIntVal clears the value of int_val.
func (ptuo *PropertyTypeUpdateOne) ClearIntVal() *PropertyTypeUpdateOne {
	ptuo.mutation.ClearIntVal()
	return ptuo
}

// SetBoolVal sets the bool_val field.
func (ptuo *PropertyTypeUpdateOne) SetBoolVal(b bool) *PropertyTypeUpdateOne {
	ptuo.mutation.SetBoolVal(b)
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
	ptuo.mutation.ClearBoolVal()
	return ptuo
}

// SetFloatVal sets the float_val field.
func (ptuo *PropertyTypeUpdateOne) SetFloatVal(f float64) *PropertyTypeUpdateOne {
	ptuo.mutation.ResetFloatVal()
	ptuo.mutation.SetFloatVal(f)
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
	ptuo.mutation.AddFloatVal(f)
	return ptuo
}

// ClearFloatVal clears the value of float_val.
func (ptuo *PropertyTypeUpdateOne) ClearFloatVal() *PropertyTypeUpdateOne {
	ptuo.mutation.ClearFloatVal()
	return ptuo
}

// SetLatitudeVal sets the latitude_val field.
func (ptuo *PropertyTypeUpdateOne) SetLatitudeVal(f float64) *PropertyTypeUpdateOne {
	ptuo.mutation.ResetLatitudeVal()
	ptuo.mutation.SetLatitudeVal(f)
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
	ptuo.mutation.AddLatitudeVal(f)
	return ptuo
}

// ClearLatitudeVal clears the value of latitude_val.
func (ptuo *PropertyTypeUpdateOne) ClearLatitudeVal() *PropertyTypeUpdateOne {
	ptuo.mutation.ClearLatitudeVal()
	return ptuo
}

// SetLongitudeVal sets the longitude_val field.
func (ptuo *PropertyTypeUpdateOne) SetLongitudeVal(f float64) *PropertyTypeUpdateOne {
	ptuo.mutation.ResetLongitudeVal()
	ptuo.mutation.SetLongitudeVal(f)
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
	ptuo.mutation.AddLongitudeVal(f)
	return ptuo
}

// ClearLongitudeVal clears the value of longitude_val.
func (ptuo *PropertyTypeUpdateOne) ClearLongitudeVal() *PropertyTypeUpdateOne {
	ptuo.mutation.ClearLongitudeVal()
	return ptuo
}

// SetStringVal sets the string_val field.
func (ptuo *PropertyTypeUpdateOne) SetStringVal(s string) *PropertyTypeUpdateOne {
	ptuo.mutation.SetStringVal(s)
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
	ptuo.mutation.ClearStringVal()
	return ptuo
}

// SetRangeFromVal sets the range_from_val field.
func (ptuo *PropertyTypeUpdateOne) SetRangeFromVal(f float64) *PropertyTypeUpdateOne {
	ptuo.mutation.ResetRangeFromVal()
	ptuo.mutation.SetRangeFromVal(f)
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
	ptuo.mutation.AddRangeFromVal(f)
	return ptuo
}

// ClearRangeFromVal clears the value of range_from_val.
func (ptuo *PropertyTypeUpdateOne) ClearRangeFromVal() *PropertyTypeUpdateOne {
	ptuo.mutation.ClearRangeFromVal()
	return ptuo
}

// SetRangeToVal sets the range_to_val field.
func (ptuo *PropertyTypeUpdateOne) SetRangeToVal(f float64) *PropertyTypeUpdateOne {
	ptuo.mutation.ResetRangeToVal()
	ptuo.mutation.SetRangeToVal(f)
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
	ptuo.mutation.AddRangeToVal(f)
	return ptuo
}

// ClearRangeToVal clears the value of range_to_val.
func (ptuo *PropertyTypeUpdateOne) ClearRangeToVal() *PropertyTypeUpdateOne {
	ptuo.mutation.ClearRangeToVal()
	return ptuo
}

// SetIsInstanceProperty sets the is_instance_property field.
func (ptuo *PropertyTypeUpdateOne) SetIsInstanceProperty(b bool) *PropertyTypeUpdateOne {
	ptuo.mutation.SetIsInstanceProperty(b)
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
	ptuo.mutation.SetEditable(b)
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
	ptuo.mutation.SetMandatory(b)
	return ptuo
}

// SetNillableMandatory sets the mandatory field if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableMandatory(b *bool) *PropertyTypeUpdateOne {
	if b != nil {
		ptuo.SetMandatory(*b)
	}
	return ptuo
}

// SetDeleted sets the deleted field.
func (ptuo *PropertyTypeUpdateOne) SetDeleted(b bool) *PropertyTypeUpdateOne {
	ptuo.mutation.SetDeleted(b)
	return ptuo
}

// SetNillableDeleted sets the deleted field if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableDeleted(b *bool) *PropertyTypeUpdateOne {
	if b != nil {
		ptuo.SetDeleted(*b)
	}
	return ptuo
}

// SetNodeType sets the nodeType field.
func (ptuo *PropertyTypeUpdateOne) SetNodeType(s string) *PropertyTypeUpdateOne {
	ptuo.mutation.SetNodeType(s)
	return ptuo
}

// SetNillableNodeType sets the nodeType field if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableNodeType(s *string) *PropertyTypeUpdateOne {
	if s != nil {
		ptuo.SetNodeType(*s)
	}
	return ptuo
}

// ClearNodeType clears the value of nodeType.
func (ptuo *PropertyTypeUpdateOne) ClearNodeType() *PropertyTypeUpdateOne {
	ptuo.mutation.ClearNodeType()
	return ptuo
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (ptuo *PropertyTypeUpdateOne) AddPropertyIDs(ids ...int) *PropertyTypeUpdateOne {
	ptuo.mutation.AddPropertyIDs(ids...)
	return ptuo
}

// AddProperties adds the properties edges to Property.
func (ptuo *PropertyTypeUpdateOne) AddProperties(p ...*Property) *PropertyTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptuo.AddPropertyIDs(ids...)
}

// SetLocationTypeID sets the location_type edge to LocationType by id.
func (ptuo *PropertyTypeUpdateOne) SetLocationTypeID(id int) *PropertyTypeUpdateOne {
	ptuo.mutation.SetLocationTypeID(id)
	return ptuo
}

// SetNillableLocationTypeID sets the location_type edge to LocationType by id if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableLocationTypeID(id *int) *PropertyTypeUpdateOne {
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
func (ptuo *PropertyTypeUpdateOne) SetEquipmentPortTypeID(id int) *PropertyTypeUpdateOne {
	ptuo.mutation.SetEquipmentPortTypeID(id)
	return ptuo
}

// SetNillableEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableEquipmentPortTypeID(id *int) *PropertyTypeUpdateOne {
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
func (ptuo *PropertyTypeUpdateOne) SetLinkEquipmentPortTypeID(id int) *PropertyTypeUpdateOne {
	ptuo.mutation.SetLinkEquipmentPortTypeID(id)
	return ptuo
}

// SetNillableLinkEquipmentPortTypeID sets the link_equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableLinkEquipmentPortTypeID(id *int) *PropertyTypeUpdateOne {
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
func (ptuo *PropertyTypeUpdateOne) SetEquipmentTypeID(id int) *PropertyTypeUpdateOne {
	ptuo.mutation.SetEquipmentTypeID(id)
	return ptuo
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableEquipmentTypeID(id *int) *PropertyTypeUpdateOne {
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
func (ptuo *PropertyTypeUpdateOne) SetServiceTypeID(id int) *PropertyTypeUpdateOne {
	ptuo.mutation.SetServiceTypeID(id)
	return ptuo
}

// SetNillableServiceTypeID sets the service_type edge to ServiceType by id if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableServiceTypeID(id *int) *PropertyTypeUpdateOne {
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
func (ptuo *PropertyTypeUpdateOne) SetWorkOrderTypeID(id int) *PropertyTypeUpdateOne {
	ptuo.mutation.SetWorkOrderTypeID(id)
	return ptuo
}

// SetNillableWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableWorkOrderTypeID(id *int) *PropertyTypeUpdateOne {
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
func (ptuo *PropertyTypeUpdateOne) SetProjectTypeID(id int) *PropertyTypeUpdateOne {
	ptuo.mutation.SetProjectTypeID(id)
	return ptuo
}

// SetNillableProjectTypeID sets the project_type edge to ProjectType by id if the given value is not nil.
func (ptuo *PropertyTypeUpdateOne) SetNillableProjectTypeID(id *int) *PropertyTypeUpdateOne {
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
func (ptuo *PropertyTypeUpdateOne) RemovePropertyIDs(ids ...int) *PropertyTypeUpdateOne {
	ptuo.mutation.RemovePropertyIDs(ids...)
	return ptuo
}

// RemoveProperties removes properties edges to Property.
func (ptuo *PropertyTypeUpdateOne) RemoveProperties(p ...*Property) *PropertyTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptuo.RemovePropertyIDs(ids...)
}

// ClearLocationType clears the location_type edge to LocationType.
func (ptuo *PropertyTypeUpdateOne) ClearLocationType() *PropertyTypeUpdateOne {
	ptuo.mutation.ClearLocationType()
	return ptuo
}

// ClearEquipmentPortType clears the equipment_port_type edge to EquipmentPortType.
func (ptuo *PropertyTypeUpdateOne) ClearEquipmentPortType() *PropertyTypeUpdateOne {
	ptuo.mutation.ClearEquipmentPortType()
	return ptuo
}

// ClearLinkEquipmentPortType clears the link_equipment_port_type edge to EquipmentPortType.
func (ptuo *PropertyTypeUpdateOne) ClearLinkEquipmentPortType() *PropertyTypeUpdateOne {
	ptuo.mutation.ClearLinkEquipmentPortType()
	return ptuo
}

// ClearEquipmentType clears the equipment_type edge to EquipmentType.
func (ptuo *PropertyTypeUpdateOne) ClearEquipmentType() *PropertyTypeUpdateOne {
	ptuo.mutation.ClearEquipmentType()
	return ptuo
}

// ClearServiceType clears the service_type edge to ServiceType.
func (ptuo *PropertyTypeUpdateOne) ClearServiceType() *PropertyTypeUpdateOne {
	ptuo.mutation.ClearServiceType()
	return ptuo
}

// ClearWorkOrderType clears the work_order_type edge to WorkOrderType.
func (ptuo *PropertyTypeUpdateOne) ClearWorkOrderType() *PropertyTypeUpdateOne {
	ptuo.mutation.ClearWorkOrderType()
	return ptuo
}

// ClearProjectType clears the project_type edge to ProjectType.
func (ptuo *PropertyTypeUpdateOne) ClearProjectType() *PropertyTypeUpdateOne {
	ptuo.mutation.ClearProjectType()
	return ptuo
}

// Save executes the query and returns the updated entity.
func (ptuo *PropertyTypeUpdateOne) Save(ctx context.Context) (*PropertyType, error) {
	if _, ok := ptuo.mutation.UpdateTime(); !ok {
		v := propertytype.UpdateDefaultUpdateTime()
		ptuo.mutation.SetUpdateTime(v)
	}

	var (
		err  error
		node *PropertyType
	)
	if len(ptuo.hooks) == 0 {
		node, err = ptuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PropertyTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ptuo.mutation = mutation
			node, err = ptuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(ptuo.hooks) - 1; i >= 0; i-- {
			mut = ptuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ptuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   propertytype.Table,
			Columns: propertytype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: propertytype.FieldID,
			},
		},
	}
	id, ok := ptuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing PropertyType.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := ptuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: propertytype.FieldUpdateTime,
		})
	}
	if value, ok := ptuo.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldType,
		})
	}
	if value, ok := ptuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldName,
		})
	}
	if value, ok := ptuo.mutation.ExternalID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldExternalID,
		})
	}
	if ptuo.mutation.ExternalIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: propertytype.FieldExternalID,
		})
	}
	if value, ok := ptuo.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: propertytype.FieldIndex,
		})
	}
	if value, ok := ptuo.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: propertytype.FieldIndex,
		})
	}
	if ptuo.mutation.IndexCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: propertytype.FieldIndex,
		})
	}
	if value, ok := ptuo.mutation.Category(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldCategory,
		})
	}
	if ptuo.mutation.CategoryCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: propertytype.FieldCategory,
		})
	}
	if value, ok := ptuo.mutation.IntVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: propertytype.FieldIntVal,
		})
	}
	if value, ok := ptuo.mutation.AddedIntVal(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: propertytype.FieldIntVal,
		})
	}
	if ptuo.mutation.IntValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: propertytype.FieldIntVal,
		})
	}
	if value, ok := ptuo.mutation.BoolVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldBoolVal,
		})
	}
	if ptuo.mutation.BoolValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: propertytype.FieldBoolVal,
		})
	}
	if value, ok := ptuo.mutation.FloatVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldFloatVal,
		})
	}
	if value, ok := ptuo.mutation.AddedFloatVal(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldFloatVal,
		})
	}
	if ptuo.mutation.FloatValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: propertytype.FieldFloatVal,
		})
	}
	if value, ok := ptuo.mutation.LatitudeVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldLatitudeVal,
		})
	}
	if value, ok := ptuo.mutation.AddedLatitudeVal(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldLatitudeVal,
		})
	}
	if ptuo.mutation.LatitudeValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: propertytype.FieldLatitudeVal,
		})
	}
	if value, ok := ptuo.mutation.LongitudeVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldLongitudeVal,
		})
	}
	if value, ok := ptuo.mutation.AddedLongitudeVal(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldLongitudeVal,
		})
	}
	if ptuo.mutation.LongitudeValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: propertytype.FieldLongitudeVal,
		})
	}
	if value, ok := ptuo.mutation.StringVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldStringVal,
		})
	}
	if ptuo.mutation.StringValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: propertytype.FieldStringVal,
		})
	}
	if value, ok := ptuo.mutation.RangeFromVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldRangeFromVal,
		})
	}
	if value, ok := ptuo.mutation.AddedRangeFromVal(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldRangeFromVal,
		})
	}
	if ptuo.mutation.RangeFromValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: propertytype.FieldRangeFromVal,
		})
	}
	if value, ok := ptuo.mutation.RangeToVal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldRangeToVal,
		})
	}
	if value, ok := ptuo.mutation.AddedRangeToVal(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldRangeToVal,
		})
	}
	if ptuo.mutation.RangeToValCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: propertytype.FieldRangeToVal,
		})
	}
	if value, ok := ptuo.mutation.IsInstanceProperty(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldIsInstanceProperty,
		})
	}
	if value, ok := ptuo.mutation.Editable(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldEditable,
		})
	}
	if value, ok := ptuo.mutation.Mandatory(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldMandatory,
		})
	}
	if value, ok := ptuo.mutation.Deleted(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldDeleted,
		})
	}
	if value, ok := ptuo.mutation.NodeType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldNodeType,
		})
	}
	if ptuo.mutation.NodeTypeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: propertytype.FieldNodeType,
		})
	}
	if nodes := ptuo.mutation.RemovedPropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   propertytype.PropertiesTable,
			Columns: []string{propertytype.PropertiesColumn},
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
	if nodes := ptuo.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   propertytype.PropertiesTable,
			Columns: []string{propertytype.PropertiesColumn},
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
	if ptuo.mutation.LocationTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.LocationTypeTable,
			Columns: []string{propertytype.LocationTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: locationtype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptuo.mutation.LocationTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.LocationTypeTable,
			Columns: []string{propertytype.LocationTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: locationtype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ptuo.mutation.EquipmentPortTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.EquipmentPortTypeTable,
			Columns: []string{propertytype.EquipmentPortTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentporttype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptuo.mutation.EquipmentPortTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.EquipmentPortTypeTable,
			Columns: []string{propertytype.EquipmentPortTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentporttype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ptuo.mutation.LinkEquipmentPortTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.LinkEquipmentPortTypeTable,
			Columns: []string{propertytype.LinkEquipmentPortTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentporttype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptuo.mutation.LinkEquipmentPortTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.LinkEquipmentPortTypeTable,
			Columns: []string{propertytype.LinkEquipmentPortTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentporttype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ptuo.mutation.EquipmentTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.EquipmentTypeTable,
			Columns: []string{propertytype.EquipmentTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmenttype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptuo.mutation.EquipmentTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.EquipmentTypeTable,
			Columns: []string{propertytype.EquipmentTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmenttype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ptuo.mutation.ServiceTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.ServiceTypeTable,
			Columns: []string{propertytype.ServiceTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: servicetype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptuo.mutation.ServiceTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.ServiceTypeTable,
			Columns: []string{propertytype.ServiceTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: servicetype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ptuo.mutation.WorkOrderTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.WorkOrderTypeTable,
			Columns: []string{propertytype.WorkOrderTypeColumn},
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
	if nodes := ptuo.mutation.WorkOrderTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.WorkOrderTypeTable,
			Columns: []string{propertytype.WorkOrderTypeColumn},
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
	if ptuo.mutation.ProjectTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.ProjectTypeTable,
			Columns: []string{propertytype.ProjectTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: projecttype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ptuo.mutation.ProjectTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.ProjectTypeTable,
			Columns: []string{propertytype.ProjectTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: projecttype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	pt = &PropertyType{config: ptuo.config}
	_spec.Assign = pt.assignValues
	_spec.ScanValues = pt.scanValues()
	if err = sqlgraph.UpdateNode(ctx, ptuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{propertytype.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return pt, nil
}
