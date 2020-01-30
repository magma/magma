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
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
)

// ServiceEndpointCreate is the builder for creating a ServiceEndpoint entity.
type ServiceEndpointCreate struct {
	config
	create_time *time.Time
	update_time *time.Time
	role        *string
	port        map[string]struct{}
	service     map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (sec *ServiceEndpointCreate) SetCreateTime(t time.Time) *ServiceEndpointCreate {
	sec.create_time = &t
	return sec
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (sec *ServiceEndpointCreate) SetNillableCreateTime(t *time.Time) *ServiceEndpointCreate {
	if t != nil {
		sec.SetCreateTime(*t)
	}
	return sec
}

// SetUpdateTime sets the update_time field.
func (sec *ServiceEndpointCreate) SetUpdateTime(t time.Time) *ServiceEndpointCreate {
	sec.update_time = &t
	return sec
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (sec *ServiceEndpointCreate) SetNillableUpdateTime(t *time.Time) *ServiceEndpointCreate {
	if t != nil {
		sec.SetUpdateTime(*t)
	}
	return sec
}

// SetRole sets the role field.
func (sec *ServiceEndpointCreate) SetRole(s string) *ServiceEndpointCreate {
	sec.role = &s
	return sec
}

// SetPortID sets the port edge to EquipmentPort by id.
func (sec *ServiceEndpointCreate) SetPortID(id string) *ServiceEndpointCreate {
	if sec.port == nil {
		sec.port = make(map[string]struct{})
	}
	sec.port[id] = struct{}{}
	return sec
}

// SetNillablePortID sets the port edge to EquipmentPort by id if the given value is not nil.
func (sec *ServiceEndpointCreate) SetNillablePortID(id *string) *ServiceEndpointCreate {
	if id != nil {
		sec = sec.SetPortID(*id)
	}
	return sec
}

// SetPort sets the port edge to EquipmentPort.
func (sec *ServiceEndpointCreate) SetPort(e *EquipmentPort) *ServiceEndpointCreate {
	return sec.SetPortID(e.ID)
}

// SetServiceID sets the service edge to Service by id.
func (sec *ServiceEndpointCreate) SetServiceID(id string) *ServiceEndpointCreate {
	if sec.service == nil {
		sec.service = make(map[string]struct{})
	}
	sec.service[id] = struct{}{}
	return sec
}

// SetNillableServiceID sets the service edge to Service by id if the given value is not nil.
func (sec *ServiceEndpointCreate) SetNillableServiceID(id *string) *ServiceEndpointCreate {
	if id != nil {
		sec = sec.SetServiceID(*id)
	}
	return sec
}

// SetService sets the service edge to Service.
func (sec *ServiceEndpointCreate) SetService(s *Service) *ServiceEndpointCreate {
	return sec.SetServiceID(s.ID)
}

// Save creates the ServiceEndpoint in the database.
func (sec *ServiceEndpointCreate) Save(ctx context.Context) (*ServiceEndpoint, error) {
	if sec.create_time == nil {
		v := serviceendpoint.DefaultCreateTime()
		sec.create_time = &v
	}
	if sec.update_time == nil {
		v := serviceendpoint.DefaultUpdateTime()
		sec.update_time = &v
	}
	if sec.role == nil {
		return nil, errors.New("ent: missing required field \"role\"")
	}
	if len(sec.port) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"port\"")
	}
	if len(sec.service) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"service\"")
	}
	return sec.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (sec *ServiceEndpointCreate) SaveX(ctx context.Context) *ServiceEndpoint {
	v, err := sec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (sec *ServiceEndpointCreate) sqlSave(ctx context.Context) (*ServiceEndpoint, error) {
	var (
		se    = &ServiceEndpoint{config: sec.config}
		_spec = &sqlgraph.CreateSpec{
			Table: serviceendpoint.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: serviceendpoint.FieldID,
			},
		}
	)
	if value := sec.create_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: serviceendpoint.FieldCreateTime,
		})
		se.CreateTime = *value
	}
	if value := sec.update_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: serviceendpoint.FieldUpdateTime,
		})
		se.UpdateTime = *value
	}
	if value := sec.role; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: serviceendpoint.FieldRole,
		})
		se.Role = *value
	}
	if nodes := sec.port; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   serviceendpoint.PortTable,
			Columns: []string{serviceendpoint.PortColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := sec.service; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   serviceendpoint.ServiceTable,
			Columns: []string{serviceendpoint.ServiceColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, sec.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	se.ID = strconv.FormatInt(id, 10)
	return se, nil
}
