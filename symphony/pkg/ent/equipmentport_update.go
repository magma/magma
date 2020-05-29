// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/ent/equipment"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentport"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/link"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/property"
	"github.com/facebookincubator/symphony/pkg/ent/serviceendpoint"
)

// EquipmentPortUpdate is the builder for updating EquipmentPort entities.
type EquipmentPortUpdate struct {
	config
	hooks      []Hook
	mutation   *EquipmentPortMutation
	predicates []predicate.EquipmentPort
}

// Where adds a new predicate for the builder.
func (epu *EquipmentPortUpdate) Where(ps ...predicate.EquipmentPort) *EquipmentPortUpdate {
	epu.predicates = append(epu.predicates, ps...)
	return epu
}

// SetDefinitionID sets the definition edge to EquipmentPortDefinition by id.
func (epu *EquipmentPortUpdate) SetDefinitionID(id int) *EquipmentPortUpdate {
	epu.mutation.SetDefinitionID(id)
	return epu
}

// SetDefinition sets the definition edge to EquipmentPortDefinition.
func (epu *EquipmentPortUpdate) SetDefinition(e *EquipmentPortDefinition) *EquipmentPortUpdate {
	return epu.SetDefinitionID(e.ID)
}

// SetParentID sets the parent edge to Equipment by id.
func (epu *EquipmentPortUpdate) SetParentID(id int) *EquipmentPortUpdate {
	epu.mutation.SetParentID(id)
	return epu
}

// SetNillableParentID sets the parent edge to Equipment by id if the given value is not nil.
func (epu *EquipmentPortUpdate) SetNillableParentID(id *int) *EquipmentPortUpdate {
	if id != nil {
		epu = epu.SetParentID(*id)
	}
	return epu
}

// SetParent sets the parent edge to Equipment.
func (epu *EquipmentPortUpdate) SetParent(e *Equipment) *EquipmentPortUpdate {
	return epu.SetParentID(e.ID)
}

// SetLinkID sets the link edge to Link by id.
func (epu *EquipmentPortUpdate) SetLinkID(id int) *EquipmentPortUpdate {
	epu.mutation.SetLinkID(id)
	return epu
}

// SetNillableLinkID sets the link edge to Link by id if the given value is not nil.
func (epu *EquipmentPortUpdate) SetNillableLinkID(id *int) *EquipmentPortUpdate {
	if id != nil {
		epu = epu.SetLinkID(*id)
	}
	return epu
}

// SetLink sets the link edge to Link.
func (epu *EquipmentPortUpdate) SetLink(l *Link) *EquipmentPortUpdate {
	return epu.SetLinkID(l.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (epu *EquipmentPortUpdate) AddPropertyIDs(ids ...int) *EquipmentPortUpdate {
	epu.mutation.AddPropertyIDs(ids...)
	return epu
}

// AddProperties adds the properties edges to Property.
func (epu *EquipmentPortUpdate) AddProperties(p ...*Property) *EquipmentPortUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return epu.AddPropertyIDs(ids...)
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (epu *EquipmentPortUpdate) AddEndpointIDs(ids ...int) *EquipmentPortUpdate {
	epu.mutation.AddEndpointIDs(ids...)
	return epu
}

// AddEndpoints adds the endpoints edges to ServiceEndpoint.
func (epu *EquipmentPortUpdate) AddEndpoints(s ...*ServiceEndpoint) *EquipmentPortUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return epu.AddEndpointIDs(ids...)
}

// ClearDefinition clears the definition edge to EquipmentPortDefinition.
func (epu *EquipmentPortUpdate) ClearDefinition() *EquipmentPortUpdate {
	epu.mutation.ClearDefinition()
	return epu
}

// ClearParent clears the parent edge to Equipment.
func (epu *EquipmentPortUpdate) ClearParent() *EquipmentPortUpdate {
	epu.mutation.ClearParent()
	return epu
}

// ClearLink clears the link edge to Link.
func (epu *EquipmentPortUpdate) ClearLink() *EquipmentPortUpdate {
	epu.mutation.ClearLink()
	return epu
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (epu *EquipmentPortUpdate) RemovePropertyIDs(ids ...int) *EquipmentPortUpdate {
	epu.mutation.RemovePropertyIDs(ids...)
	return epu
}

// RemoveProperties removes properties edges to Property.
func (epu *EquipmentPortUpdate) RemoveProperties(p ...*Property) *EquipmentPortUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return epu.RemovePropertyIDs(ids...)
}

// RemoveEndpointIDs removes the endpoints edge to ServiceEndpoint by ids.
func (epu *EquipmentPortUpdate) RemoveEndpointIDs(ids ...int) *EquipmentPortUpdate {
	epu.mutation.RemoveEndpointIDs(ids...)
	return epu
}

// RemoveEndpoints removes endpoints edges to ServiceEndpoint.
func (epu *EquipmentPortUpdate) RemoveEndpoints(s ...*ServiceEndpoint) *EquipmentPortUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return epu.RemoveEndpointIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (epu *EquipmentPortUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := epu.mutation.UpdateTime(); !ok {
		v := equipmentport.UpdateDefaultUpdateTime()
		epu.mutation.SetUpdateTime(v)
	}

	if _, ok := epu.mutation.DefinitionID(); epu.mutation.DefinitionCleared() && !ok {
		return 0, errors.New("ent: clearing a unique edge \"definition\"")
	}

	var (
		err      error
		affected int
	)
	if len(epu.hooks) == 0 {
		affected, err = epu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentPortMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			epu.mutation = mutation
			affected, err = epu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(epu.hooks) - 1; i >= 0; i-- {
			mut = epu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, epu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (epu *EquipmentPortUpdate) SaveX(ctx context.Context) int {
	affected, err := epu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (epu *EquipmentPortUpdate) Exec(ctx context.Context) error {
	_, err := epu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (epu *EquipmentPortUpdate) ExecX(ctx context.Context) {
	if err := epu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (epu *EquipmentPortUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentport.Table,
			Columns: equipmentport.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentport.FieldID,
			},
		},
	}
	if ps := epu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := epu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentport.FieldUpdateTime,
		})
	}
	if epu.mutation.DefinitionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentport.DefinitionTable,
			Columns: []string{equipmentport.DefinitionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentportdefinition.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epu.mutation.DefinitionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentport.DefinitionTable,
			Columns: []string{equipmentport.DefinitionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentportdefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if epu.mutation.ParentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipmentport.ParentTable,
			Columns: []string{equipmentport.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epu.mutation.ParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipmentport.ParentTable,
			Columns: []string{equipmentport.ParentColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if epu.mutation.LinkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentport.LinkTable,
			Columns: []string{equipmentport.LinkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: link.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epu.mutation.LinkIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentport.LinkTable,
			Columns: []string{equipmentport.LinkColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := epu.mutation.RemovedPropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentport.PropertiesTable,
			Columns: []string{equipmentport.PropertiesColumn},
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
	if nodes := epu.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentport.PropertiesTable,
			Columns: []string{equipmentport.PropertiesColumn},
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
	if nodes := epu.mutation.RemovedEndpointsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentport.EndpointsTable,
			Columns: []string{equipmentport.EndpointsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: serviceendpoint.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epu.mutation.EndpointsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentport.EndpointsTable,
			Columns: []string{equipmentport.EndpointsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: serviceendpoint.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, epu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{equipmentport.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// EquipmentPortUpdateOne is the builder for updating a single EquipmentPort entity.
type EquipmentPortUpdateOne struct {
	config
	hooks    []Hook
	mutation *EquipmentPortMutation
}

// SetDefinitionID sets the definition edge to EquipmentPortDefinition by id.
func (epuo *EquipmentPortUpdateOne) SetDefinitionID(id int) *EquipmentPortUpdateOne {
	epuo.mutation.SetDefinitionID(id)
	return epuo
}

// SetDefinition sets the definition edge to EquipmentPortDefinition.
func (epuo *EquipmentPortUpdateOne) SetDefinition(e *EquipmentPortDefinition) *EquipmentPortUpdateOne {
	return epuo.SetDefinitionID(e.ID)
}

// SetParentID sets the parent edge to Equipment by id.
func (epuo *EquipmentPortUpdateOne) SetParentID(id int) *EquipmentPortUpdateOne {
	epuo.mutation.SetParentID(id)
	return epuo
}

// SetNillableParentID sets the parent edge to Equipment by id if the given value is not nil.
func (epuo *EquipmentPortUpdateOne) SetNillableParentID(id *int) *EquipmentPortUpdateOne {
	if id != nil {
		epuo = epuo.SetParentID(*id)
	}
	return epuo
}

// SetParent sets the parent edge to Equipment.
func (epuo *EquipmentPortUpdateOne) SetParent(e *Equipment) *EquipmentPortUpdateOne {
	return epuo.SetParentID(e.ID)
}

// SetLinkID sets the link edge to Link by id.
func (epuo *EquipmentPortUpdateOne) SetLinkID(id int) *EquipmentPortUpdateOne {
	epuo.mutation.SetLinkID(id)
	return epuo
}

// SetNillableLinkID sets the link edge to Link by id if the given value is not nil.
func (epuo *EquipmentPortUpdateOne) SetNillableLinkID(id *int) *EquipmentPortUpdateOne {
	if id != nil {
		epuo = epuo.SetLinkID(*id)
	}
	return epuo
}

// SetLink sets the link edge to Link.
func (epuo *EquipmentPortUpdateOne) SetLink(l *Link) *EquipmentPortUpdateOne {
	return epuo.SetLinkID(l.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (epuo *EquipmentPortUpdateOne) AddPropertyIDs(ids ...int) *EquipmentPortUpdateOne {
	epuo.mutation.AddPropertyIDs(ids...)
	return epuo
}

// AddProperties adds the properties edges to Property.
func (epuo *EquipmentPortUpdateOne) AddProperties(p ...*Property) *EquipmentPortUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return epuo.AddPropertyIDs(ids...)
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (epuo *EquipmentPortUpdateOne) AddEndpointIDs(ids ...int) *EquipmentPortUpdateOne {
	epuo.mutation.AddEndpointIDs(ids...)
	return epuo
}

// AddEndpoints adds the endpoints edges to ServiceEndpoint.
func (epuo *EquipmentPortUpdateOne) AddEndpoints(s ...*ServiceEndpoint) *EquipmentPortUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return epuo.AddEndpointIDs(ids...)
}

// ClearDefinition clears the definition edge to EquipmentPortDefinition.
func (epuo *EquipmentPortUpdateOne) ClearDefinition() *EquipmentPortUpdateOne {
	epuo.mutation.ClearDefinition()
	return epuo
}

// ClearParent clears the parent edge to Equipment.
func (epuo *EquipmentPortUpdateOne) ClearParent() *EquipmentPortUpdateOne {
	epuo.mutation.ClearParent()
	return epuo
}

// ClearLink clears the link edge to Link.
func (epuo *EquipmentPortUpdateOne) ClearLink() *EquipmentPortUpdateOne {
	epuo.mutation.ClearLink()
	return epuo
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (epuo *EquipmentPortUpdateOne) RemovePropertyIDs(ids ...int) *EquipmentPortUpdateOne {
	epuo.mutation.RemovePropertyIDs(ids...)
	return epuo
}

// RemoveProperties removes properties edges to Property.
func (epuo *EquipmentPortUpdateOne) RemoveProperties(p ...*Property) *EquipmentPortUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return epuo.RemovePropertyIDs(ids...)
}

// RemoveEndpointIDs removes the endpoints edge to ServiceEndpoint by ids.
func (epuo *EquipmentPortUpdateOne) RemoveEndpointIDs(ids ...int) *EquipmentPortUpdateOne {
	epuo.mutation.RemoveEndpointIDs(ids...)
	return epuo
}

// RemoveEndpoints removes endpoints edges to ServiceEndpoint.
func (epuo *EquipmentPortUpdateOne) RemoveEndpoints(s ...*ServiceEndpoint) *EquipmentPortUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return epuo.RemoveEndpointIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (epuo *EquipmentPortUpdateOne) Save(ctx context.Context) (*EquipmentPort, error) {
	if _, ok := epuo.mutation.UpdateTime(); !ok {
		v := equipmentport.UpdateDefaultUpdateTime()
		epuo.mutation.SetUpdateTime(v)
	}

	if _, ok := epuo.mutation.DefinitionID(); epuo.mutation.DefinitionCleared() && !ok {
		return nil, errors.New("ent: clearing a unique edge \"definition\"")
	}

	var (
		err  error
		node *EquipmentPort
	)
	if len(epuo.hooks) == 0 {
		node, err = epuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentPortMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			epuo.mutation = mutation
			node, err = epuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(epuo.hooks) - 1; i >= 0; i-- {
			mut = epuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, epuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (epuo *EquipmentPortUpdateOne) SaveX(ctx context.Context) *EquipmentPort {
	ep, err := epuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return ep
}

// Exec executes the query on the entity.
func (epuo *EquipmentPortUpdateOne) Exec(ctx context.Context) error {
	_, err := epuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (epuo *EquipmentPortUpdateOne) ExecX(ctx context.Context) {
	if err := epuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (epuo *EquipmentPortUpdateOne) sqlSave(ctx context.Context) (ep *EquipmentPort, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentport.Table,
			Columns: equipmentport.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentport.FieldID,
			},
		},
	}
	id, ok := epuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing EquipmentPort.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := epuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentport.FieldUpdateTime,
		})
	}
	if epuo.mutation.DefinitionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentport.DefinitionTable,
			Columns: []string{equipmentport.DefinitionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentportdefinition.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epuo.mutation.DefinitionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentport.DefinitionTable,
			Columns: []string{equipmentport.DefinitionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentportdefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if epuo.mutation.ParentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipmentport.ParentTable,
			Columns: []string{equipmentport.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epuo.mutation.ParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipmentport.ParentTable,
			Columns: []string{equipmentport.ParentColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if epuo.mutation.LinkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentport.LinkTable,
			Columns: []string{equipmentport.LinkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: link.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epuo.mutation.LinkIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentport.LinkTable,
			Columns: []string{equipmentport.LinkColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := epuo.mutation.RemovedPropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentport.PropertiesTable,
			Columns: []string{equipmentport.PropertiesColumn},
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
	if nodes := epuo.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentport.PropertiesTable,
			Columns: []string{equipmentport.PropertiesColumn},
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
	if nodes := epuo.mutation.RemovedEndpointsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentport.EndpointsTable,
			Columns: []string{equipmentport.EndpointsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: serviceendpoint.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epuo.mutation.EndpointsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentport.EndpointsTable,
			Columns: []string{equipmentport.EndpointsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: serviceendpoint.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	ep = &EquipmentPort{config: epuo.config}
	_spec.Assign = ep.assignValues
	_spec.ScanValues = ep.scanValues()
	if err = sqlgraph.UpdateNode(ctx, epuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{equipmentport.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return ep, nil
}
