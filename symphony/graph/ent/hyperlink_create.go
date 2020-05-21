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
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// HyperlinkCreate is the builder for creating a Hyperlink entity.
type HyperlinkCreate struct {
	config
	mutation *HyperlinkMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (hc *HyperlinkCreate) SetCreateTime(t time.Time) *HyperlinkCreate {
	hc.mutation.SetCreateTime(t)
	return hc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (hc *HyperlinkCreate) SetNillableCreateTime(t *time.Time) *HyperlinkCreate {
	if t != nil {
		hc.SetCreateTime(*t)
	}
	return hc
}

// SetUpdateTime sets the update_time field.
func (hc *HyperlinkCreate) SetUpdateTime(t time.Time) *HyperlinkCreate {
	hc.mutation.SetUpdateTime(t)
	return hc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (hc *HyperlinkCreate) SetNillableUpdateTime(t *time.Time) *HyperlinkCreate {
	if t != nil {
		hc.SetUpdateTime(*t)
	}
	return hc
}

// SetURL sets the url field.
func (hc *HyperlinkCreate) SetURL(s string) *HyperlinkCreate {
	hc.mutation.SetURL(s)
	return hc
}

// SetName sets the name field.
func (hc *HyperlinkCreate) SetName(s string) *HyperlinkCreate {
	hc.mutation.SetName(s)
	return hc
}

// SetNillableName sets the name field if the given value is not nil.
func (hc *HyperlinkCreate) SetNillableName(s *string) *HyperlinkCreate {
	if s != nil {
		hc.SetName(*s)
	}
	return hc
}

// SetCategory sets the category field.
func (hc *HyperlinkCreate) SetCategory(s string) *HyperlinkCreate {
	hc.mutation.SetCategory(s)
	return hc
}

// SetNillableCategory sets the category field if the given value is not nil.
func (hc *HyperlinkCreate) SetNillableCategory(s *string) *HyperlinkCreate {
	if s != nil {
		hc.SetCategory(*s)
	}
	return hc
}

// SetEquipmentID sets the equipment edge to Equipment by id.
func (hc *HyperlinkCreate) SetEquipmentID(id int) *HyperlinkCreate {
	hc.mutation.SetEquipmentID(id)
	return hc
}

// SetNillableEquipmentID sets the equipment edge to Equipment by id if the given value is not nil.
func (hc *HyperlinkCreate) SetNillableEquipmentID(id *int) *HyperlinkCreate {
	if id != nil {
		hc = hc.SetEquipmentID(*id)
	}
	return hc
}

// SetEquipment sets the equipment edge to Equipment.
func (hc *HyperlinkCreate) SetEquipment(e *Equipment) *HyperlinkCreate {
	return hc.SetEquipmentID(e.ID)
}

// SetLocationID sets the location edge to Location by id.
func (hc *HyperlinkCreate) SetLocationID(id int) *HyperlinkCreate {
	hc.mutation.SetLocationID(id)
	return hc
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (hc *HyperlinkCreate) SetNillableLocationID(id *int) *HyperlinkCreate {
	if id != nil {
		hc = hc.SetLocationID(*id)
	}
	return hc
}

// SetLocation sets the location edge to Location.
func (hc *HyperlinkCreate) SetLocation(l *Location) *HyperlinkCreate {
	return hc.SetLocationID(l.ID)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (hc *HyperlinkCreate) SetWorkOrderID(id int) *HyperlinkCreate {
	hc.mutation.SetWorkOrderID(id)
	return hc
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (hc *HyperlinkCreate) SetNillableWorkOrderID(id *int) *HyperlinkCreate {
	if id != nil {
		hc = hc.SetWorkOrderID(*id)
	}
	return hc
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (hc *HyperlinkCreate) SetWorkOrder(w *WorkOrder) *HyperlinkCreate {
	return hc.SetWorkOrderID(w.ID)
}

// Save creates the Hyperlink in the database.
func (hc *HyperlinkCreate) Save(ctx context.Context) (*Hyperlink, error) {
	if _, ok := hc.mutation.CreateTime(); !ok {
		v := hyperlink.DefaultCreateTime()
		hc.mutation.SetCreateTime(v)
	}
	if _, ok := hc.mutation.UpdateTime(); !ok {
		v := hyperlink.DefaultUpdateTime()
		hc.mutation.SetUpdateTime(v)
	}
	if _, ok := hc.mutation.URL(); !ok {
		return nil, errors.New("ent: missing required field \"url\"")
	}
	var (
		err  error
		node *Hyperlink
	)
	if len(hc.hooks) == 0 {
		node, err = hc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*HyperlinkMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			hc.mutation = mutation
			node, err = hc.sqlSave(ctx)
			return node, err
		})
		for i := len(hc.hooks) - 1; i >= 0; i-- {
			mut = hc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, hc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (hc *HyperlinkCreate) SaveX(ctx context.Context) *Hyperlink {
	v, err := hc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (hc *HyperlinkCreate) sqlSave(ctx context.Context) (*Hyperlink, error) {
	var (
		h     = &Hyperlink{config: hc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: hyperlink.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: hyperlink.FieldID,
			},
		}
	)
	if value, ok := hc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: hyperlink.FieldCreateTime,
		})
		h.CreateTime = value
	}
	if value, ok := hc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: hyperlink.FieldUpdateTime,
		})
		h.UpdateTime = value
	}
	if value, ok := hc.mutation.URL(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: hyperlink.FieldURL,
		})
		h.URL = value
	}
	if value, ok := hc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: hyperlink.FieldName,
		})
		h.Name = value
	}
	if value, ok := hc.mutation.Category(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: hyperlink.FieldCategory,
		})
		h.Category = value
	}
	if nodes := hc.mutation.EquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   hyperlink.EquipmentTable,
			Columns: []string{hyperlink.EquipmentColumn},
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
	if nodes := hc.mutation.LocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   hyperlink.LocationTable,
			Columns: []string{hyperlink.LocationColumn},
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
	if nodes := hc.mutation.WorkOrderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   hyperlink.WorkOrderTable,
			Columns: []string{hyperlink.WorkOrderColumn},
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
	if err := sqlgraph.CreateNode(ctx, hc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	h.ID = int(id)
	return h, nil
}
