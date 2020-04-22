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
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// PropertyCreate is the builder for creating a Property entity.
type PropertyCreate struct {
	config
	mutation *PropertyMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (pc *PropertyCreate) SetCreateTime(t time.Time) *PropertyCreate {
	pc.mutation.SetCreateTime(t)
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
	pc.mutation.SetUpdateTime(t)
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
	pc.mutation.SetIntVal(i)
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
	pc.mutation.SetBoolVal(b)
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
	pc.mutation.SetFloatVal(f)
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
	pc.mutation.SetLatitudeVal(f)
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
	pc.mutation.SetLongitudeVal(f)
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
	pc.mutation.SetRangeFromVal(f)
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
	pc.mutation.SetRangeToVal(f)
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
	pc.mutation.SetStringVal(s)
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
func (pc *PropertyCreate) SetTypeID(id int) *PropertyCreate {
	pc.mutation.SetTypeID(id)
	return pc
}

// SetType sets the type edge to PropertyType.
func (pc *PropertyCreate) SetType(p *PropertyType) *PropertyCreate {
	return pc.SetTypeID(p.ID)
}

// SetLocationID sets the location edge to Location by id.
func (pc *PropertyCreate) SetLocationID(id int) *PropertyCreate {
	pc.mutation.SetLocationID(id)
	return pc
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableLocationID(id *int) *PropertyCreate {
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
func (pc *PropertyCreate) SetEquipmentID(id int) *PropertyCreate {
	pc.mutation.SetEquipmentID(id)
	return pc
}

// SetNillableEquipmentID sets the equipment edge to Equipment by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableEquipmentID(id *int) *PropertyCreate {
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
func (pc *PropertyCreate) SetServiceID(id int) *PropertyCreate {
	pc.mutation.SetServiceID(id)
	return pc
}

// SetNillableServiceID sets the service edge to Service by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableServiceID(id *int) *PropertyCreate {
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
func (pc *PropertyCreate) SetEquipmentPortID(id int) *PropertyCreate {
	pc.mutation.SetEquipmentPortID(id)
	return pc
}

// SetNillableEquipmentPortID sets the equipment_port edge to EquipmentPort by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableEquipmentPortID(id *int) *PropertyCreate {
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
func (pc *PropertyCreate) SetLinkID(id int) *PropertyCreate {
	pc.mutation.SetLinkID(id)
	return pc
}

// SetNillableLinkID sets the link edge to Link by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableLinkID(id *int) *PropertyCreate {
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
func (pc *PropertyCreate) SetWorkOrderID(id int) *PropertyCreate {
	pc.mutation.SetWorkOrderID(id)
	return pc
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableWorkOrderID(id *int) *PropertyCreate {
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
func (pc *PropertyCreate) SetProjectID(id int) *PropertyCreate {
	pc.mutation.SetProjectID(id)
	return pc
}

// SetNillableProjectID sets the project edge to Project by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableProjectID(id *int) *PropertyCreate {
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
func (pc *PropertyCreate) SetEquipmentValueID(id int) *PropertyCreate {
	pc.mutation.SetEquipmentValueID(id)
	return pc
}

// SetNillableEquipmentValueID sets the equipment_value edge to Equipment by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableEquipmentValueID(id *int) *PropertyCreate {
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
func (pc *PropertyCreate) SetLocationValueID(id int) *PropertyCreate {
	pc.mutation.SetLocationValueID(id)
	return pc
}

// SetNillableLocationValueID sets the location_value edge to Location by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableLocationValueID(id *int) *PropertyCreate {
	if id != nil {
		pc = pc.SetLocationValueID(*id)
	}
	return pc
}

// SetLocationValue sets the location_value edge to Location.
func (pc *PropertyCreate) SetLocationValue(l *Location) *PropertyCreate {
	return pc.SetLocationValueID(l.ID)
}

// SetServiceValueID sets the service_value edge to Service by id.
func (pc *PropertyCreate) SetServiceValueID(id int) *PropertyCreate {
	pc.mutation.SetServiceValueID(id)
	return pc
}

// SetNillableServiceValueID sets the service_value edge to Service by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableServiceValueID(id *int) *PropertyCreate {
	if id != nil {
		pc = pc.SetServiceValueID(*id)
	}
	return pc
}

// SetServiceValue sets the service_value edge to Service.
func (pc *PropertyCreate) SetServiceValue(s *Service) *PropertyCreate {
	return pc.SetServiceValueID(s.ID)
}

// SetWorkOrderValueID sets the work_order_value edge to WorkOrder by id.
func (pc *PropertyCreate) SetWorkOrderValueID(id int) *PropertyCreate {
	pc.mutation.SetWorkOrderValueID(id)
	return pc
}

// SetNillableWorkOrderValueID sets the work_order_value edge to WorkOrder by id if the given value is not nil.
func (pc *PropertyCreate) SetNillableWorkOrderValueID(id *int) *PropertyCreate {
	if id != nil {
		pc = pc.SetWorkOrderValueID(*id)
	}
	return pc
}

// SetWorkOrderValue sets the work_order_value edge to WorkOrder.
func (pc *PropertyCreate) SetWorkOrderValue(w *WorkOrder) *PropertyCreate {
	return pc.SetWorkOrderValueID(w.ID)
}

// Save creates the Property in the database.
func (pc *PropertyCreate) Save(ctx context.Context) (*Property, error) {
	if _, ok := pc.mutation.CreateTime(); !ok {
		v := property.DefaultCreateTime()
		pc.mutation.SetCreateTime(v)
	}
	if _, ok := pc.mutation.UpdateTime(); !ok {
		v := property.DefaultUpdateTime()
		pc.mutation.SetUpdateTime(v)
	}
	if _, ok := pc.mutation.TypeID(); !ok {
		return nil, errors.New("ent: missing required edge \"type\"")
	}
	var (
		err  error
		node *Property
	)
	if len(pc.hooks) == 0 {
		node, err = pc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PropertyMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			pc.mutation = mutation
			node, err = pc.sqlSave(ctx)
			return node, err
		})
		for i := len(pc.hooks) - 1; i >= 0; i-- {
			mut = pc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, pc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
		pr    = &Property{config: pc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: property.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: property.FieldID,
			},
		}
	)
	if value, ok := pc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: property.FieldCreateTime,
		})
		pr.CreateTime = value
	}
	if value, ok := pc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: property.FieldUpdateTime,
		})
		pr.UpdateTime = value
	}
	if value, ok := pc.mutation.IntVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: property.FieldIntVal,
		})
		pr.IntVal = value
	}
	if value, ok := pc.mutation.BoolVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: property.FieldBoolVal,
		})
		pr.BoolVal = value
	}
	if value, ok := pc.mutation.FloatVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: property.FieldFloatVal,
		})
		pr.FloatVal = value
	}
	if value, ok := pc.mutation.LatitudeVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: property.FieldLatitudeVal,
		})
		pr.LatitudeVal = value
	}
	if value, ok := pc.mutation.LongitudeVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: property.FieldLongitudeVal,
		})
		pr.LongitudeVal = value
	}
	if value, ok := pc.mutation.RangeFromVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: property.FieldRangeFromVal,
		})
		pr.RangeFromVal = value
	}
	if value, ok := pc.mutation.RangeToVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: property.FieldRangeToVal,
		})
		pr.RangeToVal = value
	}
	if value, ok := pc.mutation.StringVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: property.FieldStringVal,
		})
		pr.StringVal = value
	}
	if nodes := pc.mutation.TypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.TypeTable,
			Columns: []string{property.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: propertytype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.LocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.LocationTable,
			Columns: []string{property.LocationColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.EquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.EquipmentTable,
			Columns: []string{property.EquipmentColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.ServiceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.ServiceTable,
			Columns: []string{property.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: service.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.EquipmentPortIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.EquipmentPortTable,
			Columns: []string{property.EquipmentPortColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentport.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.LinkIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.LinkTable,
			Columns: []string{property.LinkColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.WorkOrderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.WorkOrderTable,
			Columns: []string{property.WorkOrderColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.ProjectIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   property.ProjectTable,
			Columns: []string{property.ProjectColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.EquipmentValueIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.EquipmentValueTable,
			Columns: []string{property.EquipmentValueColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.LocationValueIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.LocationValueTable,
			Columns: []string{property.LocationValueColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.ServiceValueIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.ServiceValueTable,
			Columns: []string{property.ServiceValueColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: service.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.WorkOrderValueIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.WorkOrderValueTable,
			Columns: []string{property.WorkOrderValueColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, pc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	pr.ID = int(id)
	return pr, nil
}
