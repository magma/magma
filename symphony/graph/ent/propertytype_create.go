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
	create_time              *time.Time
	update_time              *time.Time
	_type                    *string
	name                     *string
	index                    *int
	category                 *string
	int_val                  *int
	bool_val                 *bool
	float_val                *float64
	latitude_val             *float64
	longitude_val            *float64
	string_val               *string
	range_from_val           *float64
	range_to_val             *float64
	is_instance_property     *bool
	editable                 *bool
	mandatory                *bool
	deleted                  *bool
	properties               map[string]struct{}
	location_type            map[string]struct{}
	equipment_port_type      map[string]struct{}
	link_equipment_port_type map[string]struct{}
	equipment_type           map[string]struct{}
	service_type             map[string]struct{}
	work_order_type          map[string]struct{}
	project_type             map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (ptc *PropertyTypeCreate) SetCreateTime(t time.Time) *PropertyTypeCreate {
	ptc.create_time = &t
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
	ptc.update_time = &t
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
	ptc._type = &s
	return ptc
}

// SetName sets the name field.
func (ptc *PropertyTypeCreate) SetName(s string) *PropertyTypeCreate {
	ptc.name = &s
	return ptc
}

// SetIndex sets the index field.
func (ptc *PropertyTypeCreate) SetIndex(i int) *PropertyTypeCreate {
	ptc.index = &i
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
	ptc.category = &s
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
	ptc.int_val = &i
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
	ptc.bool_val = &b
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
	ptc.float_val = &f
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
	ptc.latitude_val = &f
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
	ptc.longitude_val = &f
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
	ptc.string_val = &s
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
	ptc.range_from_val = &f
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
	ptc.range_to_val = &f
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
	ptc.is_instance_property = &b
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
	ptc.editable = &b
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
	ptc.mandatory = &b
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
	ptc.deleted = &b
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
func (ptc *PropertyTypeCreate) AddPropertyIDs(ids ...string) *PropertyTypeCreate {
	if ptc.properties == nil {
		ptc.properties = make(map[string]struct{})
	}
	for i := range ids {
		ptc.properties[ids[i]] = struct{}{}
	}
	return ptc
}

// AddProperties adds the properties edges to Property.
func (ptc *PropertyTypeCreate) AddProperties(p ...*Property) *PropertyTypeCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptc.AddPropertyIDs(ids...)
}

// SetLocationTypeID sets the location_type edge to LocationType by id.
func (ptc *PropertyTypeCreate) SetLocationTypeID(id string) *PropertyTypeCreate {
	if ptc.location_type == nil {
		ptc.location_type = make(map[string]struct{})
	}
	ptc.location_type[id] = struct{}{}
	return ptc
}

// SetNillableLocationTypeID sets the location_type edge to LocationType by id if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableLocationTypeID(id *string) *PropertyTypeCreate {
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
func (ptc *PropertyTypeCreate) SetEquipmentPortTypeID(id string) *PropertyTypeCreate {
	if ptc.equipment_port_type == nil {
		ptc.equipment_port_type = make(map[string]struct{})
	}
	ptc.equipment_port_type[id] = struct{}{}
	return ptc
}

// SetNillableEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableEquipmentPortTypeID(id *string) *PropertyTypeCreate {
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
func (ptc *PropertyTypeCreate) SetLinkEquipmentPortTypeID(id string) *PropertyTypeCreate {
	if ptc.link_equipment_port_type == nil {
		ptc.link_equipment_port_type = make(map[string]struct{})
	}
	ptc.link_equipment_port_type[id] = struct{}{}
	return ptc
}

// SetNillableLinkEquipmentPortTypeID sets the link_equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableLinkEquipmentPortTypeID(id *string) *PropertyTypeCreate {
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
func (ptc *PropertyTypeCreate) SetEquipmentTypeID(id string) *PropertyTypeCreate {
	if ptc.equipment_type == nil {
		ptc.equipment_type = make(map[string]struct{})
	}
	ptc.equipment_type[id] = struct{}{}
	return ptc
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableEquipmentTypeID(id *string) *PropertyTypeCreate {
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
func (ptc *PropertyTypeCreate) SetServiceTypeID(id string) *PropertyTypeCreate {
	if ptc.service_type == nil {
		ptc.service_type = make(map[string]struct{})
	}
	ptc.service_type[id] = struct{}{}
	return ptc
}

// SetNillableServiceTypeID sets the service_type edge to ServiceType by id if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableServiceTypeID(id *string) *PropertyTypeCreate {
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
func (ptc *PropertyTypeCreate) SetWorkOrderTypeID(id string) *PropertyTypeCreate {
	if ptc.work_order_type == nil {
		ptc.work_order_type = make(map[string]struct{})
	}
	ptc.work_order_type[id] = struct{}{}
	return ptc
}

// SetNillableWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableWorkOrderTypeID(id *string) *PropertyTypeCreate {
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
func (ptc *PropertyTypeCreate) SetProjectTypeID(id string) *PropertyTypeCreate {
	if ptc.project_type == nil {
		ptc.project_type = make(map[string]struct{})
	}
	ptc.project_type[id] = struct{}{}
	return ptc
}

// SetNillableProjectTypeID sets the project_type edge to ProjectType by id if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableProjectTypeID(id *string) *PropertyTypeCreate {
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
	if ptc.create_time == nil {
		v := propertytype.DefaultCreateTime()
		ptc.create_time = &v
	}
	if ptc.update_time == nil {
		v := propertytype.DefaultUpdateTime()
		ptc.update_time = &v
	}
	if ptc._type == nil {
		return nil, errors.New("ent: missing required field \"type\"")
	}
	if ptc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if ptc.is_instance_property == nil {
		v := propertytype.DefaultIsInstanceProperty
		ptc.is_instance_property = &v
	}
	if ptc.editable == nil {
		v := propertytype.DefaultEditable
		ptc.editable = &v
	}
	if ptc.mandatory == nil {
		v := propertytype.DefaultMandatory
		ptc.mandatory = &v
	}
	if ptc.deleted == nil {
		v := propertytype.DefaultDeleted
		ptc.deleted = &v
	}
	if len(ptc.location_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location_type\"")
	}
	if len(ptc.equipment_port_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"equipment_port_type\"")
	}
	if len(ptc.link_equipment_port_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"link_equipment_port_type\"")
	}
	if len(ptc.equipment_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"equipment_type\"")
	}
	if len(ptc.service_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"service_type\"")
	}
	if len(ptc.work_order_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"work_order_type\"")
	}
	if len(ptc.project_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"project_type\"")
	}
	return ptc.sqlSave(ctx)
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
				Type:   field.TypeString,
				Column: propertytype.FieldID,
			},
		}
	)
	if value := ptc.create_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: propertytype.FieldCreateTime,
		})
		pt.CreateTime = *value
	}
	if value := ptc.update_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: propertytype.FieldUpdateTime,
		})
		pt.UpdateTime = *value
	}
	if value := ptc._type; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: propertytype.FieldType,
		})
		pt.Type = *value
	}
	if value := ptc.name; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: propertytype.FieldName,
		})
		pt.Name = *value
	}
	if value := ptc.index; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: propertytype.FieldIndex,
		})
		pt.Index = *value
	}
	if value := ptc.category; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: propertytype.FieldCategory,
		})
		pt.Category = *value
	}
	if value := ptc.int_val; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: propertytype.FieldIntVal,
		})
		pt.IntVal = *value
	}
	if value := ptc.bool_val; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: propertytype.FieldBoolVal,
		})
		pt.BoolVal = *value
	}
	if value := ptc.float_val; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: propertytype.FieldFloatVal,
		})
		pt.FloatVal = *value
	}
	if value := ptc.latitude_val; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: propertytype.FieldLatitudeVal,
		})
		pt.LatitudeVal = *value
	}
	if value := ptc.longitude_val; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: propertytype.FieldLongitudeVal,
		})
		pt.LongitudeVal = *value
	}
	if value := ptc.string_val; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: propertytype.FieldStringVal,
		})
		pt.StringVal = *value
	}
	if value := ptc.range_from_val; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: propertytype.FieldRangeFromVal,
		})
		pt.RangeFromVal = *value
	}
	if value := ptc.range_to_val; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: propertytype.FieldRangeToVal,
		})
		pt.RangeToVal = *value
	}
	if value := ptc.is_instance_property; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: propertytype.FieldIsInstanceProperty,
		})
		pt.IsInstanceProperty = *value
	}
	if value := ptc.editable; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: propertytype.FieldEditable,
		})
		pt.Editable = *value
	}
	if value := ptc.mandatory; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: propertytype.FieldMandatory,
		})
		pt.Mandatory = *value
	}
	if value := ptc.deleted; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: propertytype.FieldDeleted,
		})
		pt.Deleted = *value
	}
	if nodes := ptc.properties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   propertytype.PropertiesTable,
			Columns: []string{propertytype.PropertiesColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ptc.location_type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.LocationTypeTable,
			Columns: []string{propertytype.LocationTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: locationtype.FieldID,
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ptc.equipment_port_type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.EquipmentPortTypeTable,
			Columns: []string{propertytype.EquipmentPortTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentporttype.FieldID,
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ptc.link_equipment_port_type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.LinkEquipmentPortTypeTable,
			Columns: []string{propertytype.LinkEquipmentPortTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentporttype.FieldID,
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ptc.equipment_type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.EquipmentTypeTable,
			Columns: []string{propertytype.EquipmentTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmenttype.FieldID,
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ptc.service_type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.ServiceTypeTable,
			Columns: []string{propertytype.ServiceTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: servicetype.FieldID,
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ptc.work_order_type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.WorkOrderTypeTable,
			Columns: []string{propertytype.WorkOrderTypeColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ptc.project_type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.ProjectTypeTable,
			Columns: []string{propertytype.ProjectTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: projecttype.FieldID,
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, ptc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	pt.ID = strconv.FormatInt(id, 10)
	return pt, nil
}
