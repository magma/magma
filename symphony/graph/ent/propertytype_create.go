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
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// PropertyTypeCreate is the builder for creating a PropertyType entity.
type PropertyTypeCreate struct {
	config
	mutation *PropertyTypeMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (ptc *PropertyTypeCreate) SetCreateTime(t time.Time) *PropertyTypeCreate {
	ptc.mutation.SetCreateTime(t)
	return ptc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableCreateTime(t *time.Time) *PropertyTypeCreate {
	if t != nil {
		ptc.SetCreateTime(*t)
	}
	return ptc
}

// SetUpdateTime sets the update_time field.
func (ptc *PropertyTypeCreate) SetUpdateTime(t time.Time) *PropertyTypeCreate {
	ptc.mutation.SetUpdateTime(t)
	return ptc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableUpdateTime(t *time.Time) *PropertyTypeCreate {
	if t != nil {
		ptc.SetUpdateTime(*t)
	}
	return ptc
}

// SetType sets the type field.
func (ptc *PropertyTypeCreate) SetType(s string) *PropertyTypeCreate {
	ptc.mutation.SetType(s)
	return ptc
}

// SetName sets the name field.
func (ptc *PropertyTypeCreate) SetName(s string) *PropertyTypeCreate {
	ptc.mutation.SetName(s)
	return ptc
}

// SetIndex sets the index field.
func (ptc *PropertyTypeCreate) SetIndex(i int) *PropertyTypeCreate {
	ptc.mutation.SetIndex(i)
	return ptc
}

// SetNillableIndex sets the index field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableIndex(i *int) *PropertyTypeCreate {
	if i != nil {
		ptc.SetIndex(*i)
	}
	return ptc
}

// SetCategory sets the category field.
func (ptc *PropertyTypeCreate) SetCategory(s string) *PropertyTypeCreate {
	ptc.mutation.SetCategory(s)
	return ptc
}

// SetNillableCategory sets the category field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableCategory(s *string) *PropertyTypeCreate {
	if s != nil {
		ptc.SetCategory(*s)
	}
	return ptc
}

// SetIntVal sets the int_val field.
func (ptc *PropertyTypeCreate) SetIntVal(i int) *PropertyTypeCreate {
	ptc.mutation.SetIntVal(i)
	return ptc
}

// SetNillableIntVal sets the int_val field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableIntVal(i *int) *PropertyTypeCreate {
	if i != nil {
		ptc.SetIntVal(*i)
	}
	return ptc
}

// SetBoolVal sets the bool_val field.
func (ptc *PropertyTypeCreate) SetBoolVal(b bool) *PropertyTypeCreate {
	ptc.mutation.SetBoolVal(b)
	return ptc
}

// SetNillableBoolVal sets the bool_val field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableBoolVal(b *bool) *PropertyTypeCreate {
	if b != nil {
		ptc.SetBoolVal(*b)
	}
	return ptc
}

// SetFloatVal sets the float_val field.
func (ptc *PropertyTypeCreate) SetFloatVal(f float64) *PropertyTypeCreate {
	ptc.mutation.SetFloatVal(f)
	return ptc
}

// SetNillableFloatVal sets the float_val field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableFloatVal(f *float64) *PropertyTypeCreate {
	if f != nil {
		ptc.SetFloatVal(*f)
	}
	return ptc
}

// SetLatitudeVal sets the latitude_val field.
func (ptc *PropertyTypeCreate) SetLatitudeVal(f float64) *PropertyTypeCreate {
	ptc.mutation.SetLatitudeVal(f)
	return ptc
}

// SetNillableLatitudeVal sets the latitude_val field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableLatitudeVal(f *float64) *PropertyTypeCreate {
	if f != nil {
		ptc.SetLatitudeVal(*f)
	}
	return ptc
}

// SetLongitudeVal sets the longitude_val field.
func (ptc *PropertyTypeCreate) SetLongitudeVal(f float64) *PropertyTypeCreate {
	ptc.mutation.SetLongitudeVal(f)
	return ptc
}

// SetNillableLongitudeVal sets the longitude_val field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableLongitudeVal(f *float64) *PropertyTypeCreate {
	if f != nil {
		ptc.SetLongitudeVal(*f)
	}
	return ptc
}

// SetStringVal sets the string_val field.
func (ptc *PropertyTypeCreate) SetStringVal(s string) *PropertyTypeCreate {
	ptc.mutation.SetStringVal(s)
	return ptc
}

// SetNillableStringVal sets the string_val field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableStringVal(s *string) *PropertyTypeCreate {
	if s != nil {
		ptc.SetStringVal(*s)
	}
	return ptc
}

// SetRangeFromVal sets the range_from_val field.
func (ptc *PropertyTypeCreate) SetRangeFromVal(f float64) *PropertyTypeCreate {
	ptc.mutation.SetRangeFromVal(f)
	return ptc
}

// SetNillableRangeFromVal sets the range_from_val field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableRangeFromVal(f *float64) *PropertyTypeCreate {
	if f != nil {
		ptc.SetRangeFromVal(*f)
	}
	return ptc
}

// SetRangeToVal sets the range_to_val field.
func (ptc *PropertyTypeCreate) SetRangeToVal(f float64) *PropertyTypeCreate {
	ptc.mutation.SetRangeToVal(f)
	return ptc
}

// SetNillableRangeToVal sets the range_to_val field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableRangeToVal(f *float64) *PropertyTypeCreate {
	if f != nil {
		ptc.SetRangeToVal(*f)
	}
	return ptc
}

// SetIsInstanceProperty sets the is_instance_property field.
func (ptc *PropertyTypeCreate) SetIsInstanceProperty(b bool) *PropertyTypeCreate {
	ptc.mutation.SetIsInstanceProperty(b)
	return ptc
}

// SetNillableIsInstanceProperty sets the is_instance_property field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableIsInstanceProperty(b *bool) *PropertyTypeCreate {
	if b != nil {
		ptc.SetIsInstanceProperty(*b)
	}
	return ptc
}

// SetEditable sets the editable field.
func (ptc *PropertyTypeCreate) SetEditable(b bool) *PropertyTypeCreate {
	ptc.mutation.SetEditable(b)
	return ptc
}

// SetNillableEditable sets the editable field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableEditable(b *bool) *PropertyTypeCreate {
	if b != nil {
		ptc.SetEditable(*b)
	}
	return ptc
}

// SetMandatory sets the mandatory field.
func (ptc *PropertyTypeCreate) SetMandatory(b bool) *PropertyTypeCreate {
	ptc.mutation.SetMandatory(b)
	return ptc
}

// SetNillableMandatory sets the mandatory field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableMandatory(b *bool) *PropertyTypeCreate {
	if b != nil {
		ptc.SetMandatory(*b)
	}
	return ptc
}

// SetDeleted sets the deleted field.
func (ptc *PropertyTypeCreate) SetDeleted(b bool) *PropertyTypeCreate {
	ptc.mutation.SetDeleted(b)
	return ptc
}

// SetNillableDeleted sets the deleted field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableDeleted(b *bool) *PropertyTypeCreate {
	if b != nil {
		ptc.SetDeleted(*b)
	}
	return ptc
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (ptc *PropertyTypeCreate) AddPropertyIDs(ids ...int) *PropertyTypeCreate {
	ptc.mutation.AddPropertyIDs(ids...)
	return ptc
}

// AddProperties adds the properties edges to Property.
func (ptc *PropertyTypeCreate) AddProperties(p ...*Property) *PropertyTypeCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptc.AddPropertyIDs(ids...)
}

// SetLocationTypeID sets the location_type edge to LocationType by id.
func (ptc *PropertyTypeCreate) SetLocationTypeID(id int) *PropertyTypeCreate {
	ptc.mutation.SetLocationTypeID(id)
	return ptc
}

// SetNillableLocationTypeID sets the location_type edge to LocationType by id if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableLocationTypeID(id *int) *PropertyTypeCreate {
	if id != nil {
		ptc = ptc.SetLocationTypeID(*id)
	}
	return ptc
}

// SetLocationType sets the location_type edge to LocationType.
func (ptc *PropertyTypeCreate) SetLocationType(l *LocationType) *PropertyTypeCreate {
	return ptc.SetLocationTypeID(l.ID)
}

// SetEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id.
func (ptc *PropertyTypeCreate) SetEquipmentPortTypeID(id int) *PropertyTypeCreate {
	ptc.mutation.SetEquipmentPortTypeID(id)
	return ptc
}

// SetNillableEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableEquipmentPortTypeID(id *int) *PropertyTypeCreate {
	if id != nil {
		ptc = ptc.SetEquipmentPortTypeID(*id)
	}
	return ptc
}

// SetEquipmentPortType sets the equipment_port_type edge to EquipmentPortType.
func (ptc *PropertyTypeCreate) SetEquipmentPortType(e *EquipmentPortType) *PropertyTypeCreate {
	return ptc.SetEquipmentPortTypeID(e.ID)
}

// SetLinkEquipmentPortTypeID sets the link_equipment_port_type edge to EquipmentPortType by id.
func (ptc *PropertyTypeCreate) SetLinkEquipmentPortTypeID(id int) *PropertyTypeCreate {
	ptc.mutation.SetLinkEquipmentPortTypeID(id)
	return ptc
}

// SetNillableLinkEquipmentPortTypeID sets the link_equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableLinkEquipmentPortTypeID(id *int) *PropertyTypeCreate {
	if id != nil {
		ptc = ptc.SetLinkEquipmentPortTypeID(*id)
	}
	return ptc
}

// SetLinkEquipmentPortType sets the link_equipment_port_type edge to EquipmentPortType.
func (ptc *PropertyTypeCreate) SetLinkEquipmentPortType(e *EquipmentPortType) *PropertyTypeCreate {
	return ptc.SetLinkEquipmentPortTypeID(e.ID)
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (ptc *PropertyTypeCreate) SetEquipmentTypeID(id int) *PropertyTypeCreate {
	ptc.mutation.SetEquipmentTypeID(id)
	return ptc
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableEquipmentTypeID(id *int) *PropertyTypeCreate {
	if id != nil {
		ptc = ptc.SetEquipmentTypeID(*id)
	}
	return ptc
}

// SetEquipmentType sets the equipment_type edge to EquipmentType.
func (ptc *PropertyTypeCreate) SetEquipmentType(e *EquipmentType) *PropertyTypeCreate {
	return ptc.SetEquipmentTypeID(e.ID)
}

// SetServiceTypeID sets the service_type edge to ServiceType by id.
func (ptc *PropertyTypeCreate) SetServiceTypeID(id int) *PropertyTypeCreate {
	ptc.mutation.SetServiceTypeID(id)
	return ptc
}

// SetNillableServiceTypeID sets the service_type edge to ServiceType by id if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableServiceTypeID(id *int) *PropertyTypeCreate {
	if id != nil {
		ptc = ptc.SetServiceTypeID(*id)
	}
	return ptc
}

// SetServiceType sets the service_type edge to ServiceType.
func (ptc *PropertyTypeCreate) SetServiceType(s *ServiceType) *PropertyTypeCreate {
	return ptc.SetServiceTypeID(s.ID)
}

// SetWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id.
func (ptc *PropertyTypeCreate) SetWorkOrderTypeID(id int) *PropertyTypeCreate {
	ptc.mutation.SetWorkOrderTypeID(id)
	return ptc
}

// SetNillableWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableWorkOrderTypeID(id *int) *PropertyTypeCreate {
	if id != nil {
		ptc = ptc.SetWorkOrderTypeID(*id)
	}
	return ptc
}

// SetWorkOrderType sets the work_order_type edge to WorkOrderType.
func (ptc *PropertyTypeCreate) SetWorkOrderType(w *WorkOrderType) *PropertyTypeCreate {
	return ptc.SetWorkOrderTypeID(w.ID)
}

// SetProjectTypeID sets the project_type edge to ProjectType by id.
func (ptc *PropertyTypeCreate) SetProjectTypeID(id int) *PropertyTypeCreate {
	ptc.mutation.SetProjectTypeID(id)
	return ptc
}

// SetNillableProjectTypeID sets the project_type edge to ProjectType by id if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableProjectTypeID(id *int) *PropertyTypeCreate {
	if id != nil {
		ptc = ptc.SetProjectTypeID(*id)
	}
	return ptc
}

// SetProjectType sets the project_type edge to ProjectType.
func (ptc *PropertyTypeCreate) SetProjectType(p *ProjectType) *PropertyTypeCreate {
	return ptc.SetProjectTypeID(p.ID)
}

// Save creates the PropertyType in the database.
func (ptc *PropertyTypeCreate) Save(ctx context.Context) (*PropertyType, error) {
	if _, ok := ptc.mutation.CreateTime(); !ok {
		v := propertytype.DefaultCreateTime()
		ptc.mutation.SetCreateTime(v)
	}
	if _, ok := ptc.mutation.UpdateTime(); !ok {
		v := propertytype.DefaultUpdateTime()
		ptc.mutation.SetUpdateTime(v)
	}
	if _, ok := ptc.mutation.GetType(); !ok {
		return nil, errors.New("ent: missing required field \"type\"")
	}
	if _, ok := ptc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if _, ok := ptc.mutation.IsInstanceProperty(); !ok {
		v := propertytype.DefaultIsInstanceProperty
		ptc.mutation.SetIsInstanceProperty(v)
	}
	if _, ok := ptc.mutation.Editable(); !ok {
		v := propertytype.DefaultEditable
		ptc.mutation.SetEditable(v)
	}
	if _, ok := ptc.mutation.Mandatory(); !ok {
		v := propertytype.DefaultMandatory
		ptc.mutation.SetMandatory(v)
	}
	if _, ok := ptc.mutation.Deleted(); !ok {
		v := propertytype.DefaultDeleted
		ptc.mutation.SetDeleted(v)
	}
	var (
		err  error
		node *PropertyType
	)
	if len(ptc.hooks) == 0 {
		node, err = ptc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PropertyTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ptc.mutation = mutation
			node, err = ptc.sqlSave(ctx)
			return node, err
		})
		for i := len(ptc.hooks); i > 0; i-- {
			mut = ptc.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, ptc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (ptc *PropertyTypeCreate) SaveX(ctx context.Context) *PropertyType {
	v, err := ptc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ptc *PropertyTypeCreate) sqlSave(ctx context.Context) (*PropertyType, error) {
	var (
		pt    = &PropertyType{config: ptc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: propertytype.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: propertytype.FieldID,
			},
		}
	)
	if value, ok := ptc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: propertytype.FieldCreateTime,
		})
		pt.CreateTime = value
	}
	if value, ok := ptc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: propertytype.FieldUpdateTime,
		})
		pt.UpdateTime = value
	}
	if value, ok := ptc.mutation.GetType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldType,
		})
		pt.Type = value
	}
	if value, ok := ptc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldName,
		})
		pt.Name = value
	}
	if value, ok := ptc.mutation.Index(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: propertytype.FieldIndex,
		})
		pt.Index = value
	}
	if value, ok := ptc.mutation.Category(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldCategory,
		})
		pt.Category = value
	}
	if value, ok := ptc.mutation.IntVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: propertytype.FieldIntVal,
		})
		pt.IntVal = value
	}
	if value, ok := ptc.mutation.BoolVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldBoolVal,
		})
		pt.BoolVal = value
	}
	if value, ok := ptc.mutation.FloatVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldFloatVal,
		})
		pt.FloatVal = value
	}
	if value, ok := ptc.mutation.LatitudeVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldLatitudeVal,
		})
		pt.LatitudeVal = value
	}
	if value, ok := ptc.mutation.LongitudeVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldLongitudeVal,
		})
		pt.LongitudeVal = value
	}
	if value, ok := ptc.mutation.StringVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldStringVal,
		})
		pt.StringVal = value
	}
	if value, ok := ptc.mutation.RangeFromVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldRangeFromVal,
		})
		pt.RangeFromVal = value
	}
	if value, ok := ptc.mutation.RangeToVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldRangeToVal,
		})
		pt.RangeToVal = value
	}
	if value, ok := ptc.mutation.IsInstanceProperty(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldIsInstanceProperty,
		})
		pt.IsInstanceProperty = value
	}
	if value, ok := ptc.mutation.Editable(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldEditable,
		})
		pt.Editable = value
	}
	if value, ok := ptc.mutation.Mandatory(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldMandatory,
		})
		pt.Mandatory = value
	}
	if value, ok := ptc.mutation.Deleted(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldDeleted,
		})
		pt.Deleted = value
	}
	if nodes := ptc.mutation.PropertiesIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ptc.mutation.LocationTypeIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ptc.mutation.EquipmentPortTypeIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ptc.mutation.LinkEquipmentPortTypeIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ptc.mutation.EquipmentTypeIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ptc.mutation.ServiceTypeIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ptc.mutation.WorkOrderTypeIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ptc.mutation.ProjectTypeIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, ptc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	pt.ID = int(id)
	return pt, nil
}
