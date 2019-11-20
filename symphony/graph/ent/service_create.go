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
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/service"
)

// ServiceCreate is the builder for creating a Service entity.
type ServiceCreate struct {
	config
	create_time        *time.Time
	update_time        *time.Time
	name               *string
	external_id        *string
	_type              map[string]struct{}
	downstream         map[string]struct{}
	upstream           map[string]struct{}
	properties         map[string]struct{}
	termination_points map[string]struct{}
	links              map[string]struct{}
	customer           map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (sc *ServiceCreate) SetCreateTime(t time.Time) *ServiceCreate {
	sc.create_time = &t
	return sc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (sc *ServiceCreate) SetNillableCreateTime(t *time.Time) *ServiceCreate {
	if t != nil {
		sc.SetCreateTime(*t)
	}
	return sc
}

// SetUpdateTime sets the update_time field.
func (sc *ServiceCreate) SetUpdateTime(t time.Time) *ServiceCreate {
	sc.update_time = &t
	return sc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (sc *ServiceCreate) SetNillableUpdateTime(t *time.Time) *ServiceCreate {
	if t != nil {
		sc.SetUpdateTime(*t)
	}
	return sc
}

// SetName sets the name field.
func (sc *ServiceCreate) SetName(s string) *ServiceCreate {
	sc.name = &s
	return sc
}

// SetExternalID sets the external_id field.
func (sc *ServiceCreate) SetExternalID(s string) *ServiceCreate {
	sc.external_id = &s
	return sc
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (sc *ServiceCreate) SetNillableExternalID(s *string) *ServiceCreate {
	if s != nil {
		sc.SetExternalID(*s)
	}
	return sc
}

// SetTypeID sets the type edge to ServiceType by id.
func (sc *ServiceCreate) SetTypeID(id string) *ServiceCreate {
	if sc._type == nil {
		sc._type = make(map[string]struct{})
	}
	sc._type[id] = struct{}{}
	return sc
}

// SetType sets the type edge to ServiceType.
func (sc *ServiceCreate) SetType(s *ServiceType) *ServiceCreate {
	return sc.SetTypeID(s.ID)
}

// AddDownstreamIDs adds the downstream edge to Service by ids.
func (sc *ServiceCreate) AddDownstreamIDs(ids ...string) *ServiceCreate {
	if sc.downstream == nil {
		sc.downstream = make(map[string]struct{})
	}
	for i := range ids {
		sc.downstream[ids[i]] = struct{}{}
	}
	return sc
}

// AddDownstream adds the downstream edges to Service.
func (sc *ServiceCreate) AddDownstream(s ...*Service) *ServiceCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sc.AddDownstreamIDs(ids...)
}

// AddUpstreamIDs adds the upstream edge to Service by ids.
func (sc *ServiceCreate) AddUpstreamIDs(ids ...string) *ServiceCreate {
	if sc.upstream == nil {
		sc.upstream = make(map[string]struct{})
	}
	for i := range ids {
		sc.upstream[ids[i]] = struct{}{}
	}
	return sc
}

// AddUpstream adds the upstream edges to Service.
func (sc *ServiceCreate) AddUpstream(s ...*Service) *ServiceCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sc.AddUpstreamIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (sc *ServiceCreate) AddPropertyIDs(ids ...string) *ServiceCreate {
	if sc.properties == nil {
		sc.properties = make(map[string]struct{})
	}
	for i := range ids {
		sc.properties[ids[i]] = struct{}{}
	}
	return sc
}

// AddProperties adds the properties edges to Property.
func (sc *ServiceCreate) AddProperties(p ...*Property) *ServiceCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return sc.AddPropertyIDs(ids...)
}

// AddTerminationPointIDs adds the termination_points edge to Equipment by ids.
func (sc *ServiceCreate) AddTerminationPointIDs(ids ...string) *ServiceCreate {
	if sc.termination_points == nil {
		sc.termination_points = make(map[string]struct{})
	}
	for i := range ids {
		sc.termination_points[ids[i]] = struct{}{}
	}
	return sc
}

// AddTerminationPoints adds the termination_points edges to Equipment.
func (sc *ServiceCreate) AddTerminationPoints(e ...*Equipment) *ServiceCreate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return sc.AddTerminationPointIDs(ids...)
}

// AddLinkIDs adds the links edge to Link by ids.
func (sc *ServiceCreate) AddLinkIDs(ids ...string) *ServiceCreate {
	if sc.links == nil {
		sc.links = make(map[string]struct{})
	}
	for i := range ids {
		sc.links[ids[i]] = struct{}{}
	}
	return sc
}

// AddLinks adds the links edges to Link.
func (sc *ServiceCreate) AddLinks(l ...*Link) *ServiceCreate {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return sc.AddLinkIDs(ids...)
}

// AddCustomerIDs adds the customer edge to Customer by ids.
func (sc *ServiceCreate) AddCustomerIDs(ids ...string) *ServiceCreate {
	if sc.customer == nil {
		sc.customer = make(map[string]struct{})
	}
	for i := range ids {
		sc.customer[ids[i]] = struct{}{}
	}
	return sc
}

// AddCustomer adds the customer edges to Customer.
func (sc *ServiceCreate) AddCustomer(c ...*Customer) *ServiceCreate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return sc.AddCustomerIDs(ids...)
}

// Save creates the Service in the database.
func (sc *ServiceCreate) Save(ctx context.Context) (*Service, error) {
	if sc.create_time == nil {
		v := service.DefaultCreateTime()
		sc.create_time = &v
	}
	if sc.update_time == nil {
		v := service.DefaultUpdateTime()
		sc.update_time = &v
	}
	if sc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if err := service.NameValidator(*sc.name); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
	}
	if sc.external_id != nil {
		if err := service.ExternalIDValidator(*sc.external_id); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"external_id\": %v", err)
		}
	}
	if len(sc._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if sc._type == nil {
		return nil, errors.New("ent: missing required edge \"type\"")
	}
	return sc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (sc *ServiceCreate) SaveX(ctx context.Context) *Service {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (sc *ServiceCreate) sqlSave(ctx context.Context) (*Service, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(sc.driver.Dialect())
		s       = &Service{config: sc.config}
	)
	tx, err := sc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(service.Table).Default()
	if value := sc.create_time; value != nil {
		insert.Set(service.FieldCreateTime, *value)
		s.CreateTime = *value
	}
	if value := sc.update_time; value != nil {
		insert.Set(service.FieldUpdateTime, *value)
		s.UpdateTime = *value
	}
	if value := sc.name; value != nil {
		insert.Set(service.FieldName, *value)
		s.Name = *value
	}
	if value := sc.external_id; value != nil {
		insert.Set(service.FieldExternalID, *value)
		s.ExternalID = *value
	}
	id, err := insertLastID(ctx, tx, insert.Returning(service.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	s.ID = strconv.FormatInt(id, 10)
	if len(sc._type) > 0 {
		for eid := range sc._type {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(service.TypeTable).
				Set(service.TypeColumn, eid).
				Where(sql.EQ(service.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(sc.downstream) > 0 {
		for eid := range sc.downstream {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}

			query, args := builder.Insert(service.DownstreamTable).
				Columns(service.DownstreamPrimaryKey[1], service.DownstreamPrimaryKey[0]).
				Values(id, eid).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(sc.upstream) > 0 {
		for eid := range sc.upstream {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}

			query, args := builder.Insert(service.UpstreamTable).
				Columns(service.UpstreamPrimaryKey[0], service.UpstreamPrimaryKey[1]).
				Values(id, eid).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(sc.properties) > 0 {
		p := sql.P()
		for eid := range sc.properties {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(property.FieldID, eid)
		}
		query, args := builder.Update(service.PropertiesTable).
			Set(service.PropertiesColumn, id).
			Where(sql.And(p, sql.IsNull(service.PropertiesColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(sc.properties) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"Service\"", keys(sc.properties))})
		}
	}
	if len(sc.termination_points) > 0 {
		for eid := range sc.termination_points {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}

			query, args := builder.Insert(service.TerminationPointsTable).
				Columns(service.TerminationPointsPrimaryKey[0], service.TerminationPointsPrimaryKey[1]).
				Values(id, eid).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(sc.links) > 0 {
		for eid := range sc.links {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}

			query, args := builder.Insert(service.LinksTable).
				Columns(service.LinksPrimaryKey[0], service.LinksPrimaryKey[1]).
				Values(id, eid).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(sc.customer) > 0 {
		for eid := range sc.customer {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}

			query, args := builder.Insert(service.CustomerTable).
				Columns(service.CustomerPrimaryKey[0], service.CustomerPrimaryKey[1]).
				Values(id, eid).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s, nil
}
