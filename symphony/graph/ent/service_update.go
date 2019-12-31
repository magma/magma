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
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
)

// ServiceUpdate is the builder for updating Service entities.
type ServiceUpdate struct {
	config

	update_time              *time.Time
	name                     *string
	external_id              *string
	clearexternal_id         bool
	status                   *string
	_type                    map[string]struct{}
	downstream               map[string]struct{}
	upstream                 map[string]struct{}
	properties               map[string]struct{}
	termination_points       map[string]struct{}
	links                    map[string]struct{}
	customer                 map[string]struct{}
	endpoints                map[string]struct{}
	clearedType              bool
	removedDownstream        map[string]struct{}
	removedUpstream          map[string]struct{}
	removedProperties        map[string]struct{}
	removedTerminationPoints map[string]struct{}
	removedLinks             map[string]struct{}
	removedCustomer          map[string]struct{}
	removedEndpoints         map[string]struct{}
	predicates               []predicate.Service
}

// Where adds a new predicate for the builder.
func (su *ServiceUpdate) Where(ps ...predicate.Service) *ServiceUpdate {
	su.predicates = append(su.predicates, ps...)
	return su
}

// SetName sets the name field.
func (su *ServiceUpdate) SetName(s string) *ServiceUpdate {
	su.name = &s
	return su
}

// SetExternalID sets the external_id field.
func (su *ServiceUpdate) SetExternalID(s string) *ServiceUpdate {
	su.external_id = &s
	return su
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (su *ServiceUpdate) SetNillableExternalID(s *string) *ServiceUpdate {
	if s != nil {
		su.SetExternalID(*s)
	}
	return su
}

// ClearExternalID clears the value of external_id.
func (su *ServiceUpdate) ClearExternalID() *ServiceUpdate {
	su.external_id = nil
	su.clearexternal_id = true
	return su
}

// SetStatus sets the status field.
func (su *ServiceUpdate) SetStatus(s string) *ServiceUpdate {
	su.status = &s
	return su
}

// SetTypeID sets the type edge to ServiceType by id.
func (su *ServiceUpdate) SetTypeID(id string) *ServiceUpdate {
	if su._type == nil {
		su._type = make(map[string]struct{})
	}
	su._type[id] = struct{}{}
	return su
}

// SetType sets the type edge to ServiceType.
func (su *ServiceUpdate) SetType(s *ServiceType) *ServiceUpdate {
	return su.SetTypeID(s.ID)
}

// AddDownstreamIDs adds the downstream edge to Service by ids.
func (su *ServiceUpdate) AddDownstreamIDs(ids ...string) *ServiceUpdate {
	if su.downstream == nil {
		su.downstream = make(map[string]struct{})
	}
	for i := range ids {
		su.downstream[ids[i]] = struct{}{}
	}
	return su
}

// AddDownstream adds the downstream edges to Service.
func (su *ServiceUpdate) AddDownstream(s ...*Service) *ServiceUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddDownstreamIDs(ids...)
}

// AddUpstreamIDs adds the upstream edge to Service by ids.
func (su *ServiceUpdate) AddUpstreamIDs(ids ...string) *ServiceUpdate {
	if su.upstream == nil {
		su.upstream = make(map[string]struct{})
	}
	for i := range ids {
		su.upstream[ids[i]] = struct{}{}
	}
	return su
}

// AddUpstream adds the upstream edges to Service.
func (su *ServiceUpdate) AddUpstream(s ...*Service) *ServiceUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddUpstreamIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (su *ServiceUpdate) AddPropertyIDs(ids ...string) *ServiceUpdate {
	if su.properties == nil {
		su.properties = make(map[string]struct{})
	}
	for i := range ids {
		su.properties[ids[i]] = struct{}{}
	}
	return su
}

// AddProperties adds the properties edges to Property.
func (su *ServiceUpdate) AddProperties(p ...*Property) *ServiceUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return su.AddPropertyIDs(ids...)
}

// AddTerminationPointIDs adds the termination_points edge to Equipment by ids.
func (su *ServiceUpdate) AddTerminationPointIDs(ids ...string) *ServiceUpdate {
	if su.termination_points == nil {
		su.termination_points = make(map[string]struct{})
	}
	for i := range ids {
		su.termination_points[ids[i]] = struct{}{}
	}
	return su
}

// AddTerminationPoints adds the termination_points edges to Equipment.
func (su *ServiceUpdate) AddTerminationPoints(e ...*Equipment) *ServiceUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return su.AddTerminationPointIDs(ids...)
}

// AddLinkIDs adds the links edge to Link by ids.
func (su *ServiceUpdate) AddLinkIDs(ids ...string) *ServiceUpdate {
	if su.links == nil {
		su.links = make(map[string]struct{})
	}
	for i := range ids {
		su.links[ids[i]] = struct{}{}
	}
	return su
}

// AddLinks adds the links edges to Link.
func (su *ServiceUpdate) AddLinks(l ...*Link) *ServiceUpdate {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return su.AddLinkIDs(ids...)
}

// AddCustomerIDs adds the customer edge to Customer by ids.
func (su *ServiceUpdate) AddCustomerIDs(ids ...string) *ServiceUpdate {
	if su.customer == nil {
		su.customer = make(map[string]struct{})
	}
	for i := range ids {
		su.customer[ids[i]] = struct{}{}
	}
	return su
}

// AddCustomer adds the customer edges to Customer.
func (su *ServiceUpdate) AddCustomer(c ...*Customer) *ServiceUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return su.AddCustomerIDs(ids...)
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (su *ServiceUpdate) AddEndpointIDs(ids ...string) *ServiceUpdate {
	if su.endpoints == nil {
		su.endpoints = make(map[string]struct{})
	}
	for i := range ids {
		su.endpoints[ids[i]] = struct{}{}
	}
	return su
}

// AddEndpoints adds the endpoints edges to ServiceEndpoint.
func (su *ServiceUpdate) AddEndpoints(s ...*ServiceEndpoint) *ServiceUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddEndpointIDs(ids...)
}

// ClearType clears the type edge to ServiceType.
func (su *ServiceUpdate) ClearType() *ServiceUpdate {
	su.clearedType = true
	return su
}

// RemoveDownstreamIDs removes the downstream edge to Service by ids.
func (su *ServiceUpdate) RemoveDownstreamIDs(ids ...string) *ServiceUpdate {
	if su.removedDownstream == nil {
		su.removedDownstream = make(map[string]struct{})
	}
	for i := range ids {
		su.removedDownstream[ids[i]] = struct{}{}
	}
	return su
}

// RemoveDownstream removes downstream edges to Service.
func (su *ServiceUpdate) RemoveDownstream(s ...*Service) *ServiceUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveDownstreamIDs(ids...)
}

// RemoveUpstreamIDs removes the upstream edge to Service by ids.
func (su *ServiceUpdate) RemoveUpstreamIDs(ids ...string) *ServiceUpdate {
	if su.removedUpstream == nil {
		su.removedUpstream = make(map[string]struct{})
	}
	for i := range ids {
		su.removedUpstream[ids[i]] = struct{}{}
	}
	return su
}

// RemoveUpstream removes upstream edges to Service.
func (su *ServiceUpdate) RemoveUpstream(s ...*Service) *ServiceUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveUpstreamIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (su *ServiceUpdate) RemovePropertyIDs(ids ...string) *ServiceUpdate {
	if su.removedProperties == nil {
		su.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		su.removedProperties[ids[i]] = struct{}{}
	}
	return su
}

// RemoveProperties removes properties edges to Property.
func (su *ServiceUpdate) RemoveProperties(p ...*Property) *ServiceUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return su.RemovePropertyIDs(ids...)
}

// RemoveTerminationPointIDs removes the termination_points edge to Equipment by ids.
func (su *ServiceUpdate) RemoveTerminationPointIDs(ids ...string) *ServiceUpdate {
	if su.removedTerminationPoints == nil {
		su.removedTerminationPoints = make(map[string]struct{})
	}
	for i := range ids {
		su.removedTerminationPoints[ids[i]] = struct{}{}
	}
	return su
}

// RemoveTerminationPoints removes termination_points edges to Equipment.
func (su *ServiceUpdate) RemoveTerminationPoints(e ...*Equipment) *ServiceUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return su.RemoveTerminationPointIDs(ids...)
}

// RemoveLinkIDs removes the links edge to Link by ids.
func (su *ServiceUpdate) RemoveLinkIDs(ids ...string) *ServiceUpdate {
	if su.removedLinks == nil {
		su.removedLinks = make(map[string]struct{})
	}
	for i := range ids {
		su.removedLinks[ids[i]] = struct{}{}
	}
	return su
}

// RemoveLinks removes links edges to Link.
func (su *ServiceUpdate) RemoveLinks(l ...*Link) *ServiceUpdate {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return su.RemoveLinkIDs(ids...)
}

// RemoveCustomerIDs removes the customer edge to Customer by ids.
func (su *ServiceUpdate) RemoveCustomerIDs(ids ...string) *ServiceUpdate {
	if su.removedCustomer == nil {
		su.removedCustomer = make(map[string]struct{})
	}
	for i := range ids {
		su.removedCustomer[ids[i]] = struct{}{}
	}
	return su
}

// RemoveCustomer removes customer edges to Customer.
func (su *ServiceUpdate) RemoveCustomer(c ...*Customer) *ServiceUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return su.RemoveCustomerIDs(ids...)
}

// RemoveEndpointIDs removes the endpoints edge to ServiceEndpoint by ids.
func (su *ServiceUpdate) RemoveEndpointIDs(ids ...string) *ServiceUpdate {
	if su.removedEndpoints == nil {
		su.removedEndpoints = make(map[string]struct{})
	}
	for i := range ids {
		su.removedEndpoints[ids[i]] = struct{}{}
	}
	return su
}

// RemoveEndpoints removes endpoints edges to ServiceEndpoint.
func (su *ServiceUpdate) RemoveEndpoints(s ...*ServiceEndpoint) *ServiceUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveEndpointIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (su *ServiceUpdate) Save(ctx context.Context) (int, error) {
	if su.update_time == nil {
		v := service.UpdateDefaultUpdateTime()
		su.update_time = &v
	}
	if su.name != nil {
		if err := service.NameValidator(*su.name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if su.external_id != nil {
		if err := service.ExternalIDValidator(*su.external_id); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"external_id\": %v", err)
		}
	}
	if len(su._type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if su.clearedType && su._type == nil {
		return 0, errors.New("ent: clearing a unique edge \"type\"")
	}
	return su.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (su *ServiceUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *ServiceUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *ServiceUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

func (su *ServiceUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(su.driver.Dialect())
		selector = builder.Select(service.FieldID).From(builder.Table(service.Table))
	)
	for _, p := range su.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = su.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := su.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(service.Table)
	)
	updater = updater.Where(sql.InInts(service.FieldID, ids...))
	if value := su.update_time; value != nil {
		updater.Set(service.FieldUpdateTime, *value)
	}
	if value := su.name; value != nil {
		updater.Set(service.FieldName, *value)
	}
	if value := su.external_id; value != nil {
		updater.Set(service.FieldExternalID, *value)
	}
	if su.clearexternal_id {
		updater.SetNull(service.FieldExternalID)
	}
	if value := su.status; value != nil {
		updater.Set(service.FieldStatus, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if su.clearedType {
		query, args := builder.Update(service.TypeTable).
			SetNull(service.TypeColumn).
			Where(sql.InInts(servicetype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(su._type) > 0 {
		for eid := range su._type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(service.TypeTable).
				Set(service.TypeColumn, eid).
				Where(sql.InInts(service.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if len(su.removedDownstream) > 0 {
		eids := make([]int, len(su.removedDownstream))
		for eid := range su.removedDownstream {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Delete(service.DownstreamTable).
			Where(sql.InInts(service.DownstreamPrimaryKey[1], ids...)).
			Where(sql.InInts(service.DownstreamPrimaryKey[0], eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(su.downstream) > 0 {
		values := make([][]int, 0, len(ids))
		for _, id := range ids {
			for eid := range su.downstream {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				values = append(values, []int{id, eid})
			}
		}
		builder := builder.Insert(service.DownstreamTable).
			Columns(service.DownstreamPrimaryKey[1], service.DownstreamPrimaryKey[0])
		for _, v := range values {
			builder.Values(v[0], v[1])
		}
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(su.removedUpstream) > 0 {
		eids := make([]int, len(su.removedUpstream))
		for eid := range su.removedUpstream {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Delete(service.UpstreamTable).
			Where(sql.InInts(service.UpstreamPrimaryKey[0], ids...)).
			Where(sql.InInts(service.UpstreamPrimaryKey[1], eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(su.upstream) > 0 {
		values := make([][]int, 0, len(ids))
		for _, id := range ids {
			for eid := range su.upstream {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				values = append(values, []int{id, eid})
			}
		}
		builder := builder.Insert(service.UpstreamTable).
			Columns(service.UpstreamPrimaryKey[0], service.UpstreamPrimaryKey[1])
		for _, v := range values {
			builder.Values(v[0], v[1])
		}
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(su.removedProperties) > 0 {
		eids := make([]int, len(su.removedProperties))
		for eid := range su.removedProperties {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(service.PropertiesTable).
			SetNull(service.PropertiesColumn).
			Where(sql.InInts(service.PropertiesColumn, ids...)).
			Where(sql.InInts(property.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(su.properties) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range su.properties {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(property.FieldID, eid)
			}
			query, args := builder.Update(service.PropertiesTable).
				Set(service.PropertiesColumn, id).
				Where(sql.And(p, sql.IsNull(service.PropertiesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(su.properties) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"Service\"", keys(su.properties))})
			}
		}
	}
	if len(su.removedTerminationPoints) > 0 {
		eids := make([]int, len(su.removedTerminationPoints))
		for eid := range su.removedTerminationPoints {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Delete(service.TerminationPointsTable).
			Where(sql.InInts(service.TerminationPointsPrimaryKey[0], ids...)).
			Where(sql.InInts(service.TerminationPointsPrimaryKey[1], eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(su.termination_points) > 0 {
		values := make([][]int, 0, len(ids))
		for _, id := range ids {
			for eid := range su.termination_points {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				values = append(values, []int{id, eid})
			}
		}
		builder := builder.Insert(service.TerminationPointsTable).
			Columns(service.TerminationPointsPrimaryKey[0], service.TerminationPointsPrimaryKey[1])
		for _, v := range values {
			builder.Values(v[0], v[1])
		}
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(su.removedLinks) > 0 {
		eids := make([]int, len(su.removedLinks))
		for eid := range su.removedLinks {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Delete(service.LinksTable).
			Where(sql.InInts(service.LinksPrimaryKey[0], ids...)).
			Where(sql.InInts(service.LinksPrimaryKey[1], eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(su.links) > 0 {
		values := make([][]int, 0, len(ids))
		for _, id := range ids {
			for eid := range su.links {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				values = append(values, []int{id, eid})
			}
		}
		builder := builder.Insert(service.LinksTable).
			Columns(service.LinksPrimaryKey[0], service.LinksPrimaryKey[1])
		for _, v := range values {
			builder.Values(v[0], v[1])
		}
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(su.removedCustomer) > 0 {
		eids := make([]int, len(su.removedCustomer))
		for eid := range su.removedCustomer {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Delete(service.CustomerTable).
			Where(sql.InInts(service.CustomerPrimaryKey[0], ids...)).
			Where(sql.InInts(service.CustomerPrimaryKey[1], eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(su.customer) > 0 {
		values := make([][]int, 0, len(ids))
		for _, id := range ids {
			for eid := range su.customer {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				values = append(values, []int{id, eid})
			}
		}
		builder := builder.Insert(service.CustomerTable).
			Columns(service.CustomerPrimaryKey[0], service.CustomerPrimaryKey[1])
		for _, v := range values {
			builder.Values(v[0], v[1])
		}
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(su.removedEndpoints) > 0 {
		eids := make([]int, len(su.removedEndpoints))
		for eid := range su.removedEndpoints {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(service.EndpointsTable).
			SetNull(service.EndpointsColumn).
			Where(sql.InInts(service.EndpointsColumn, ids...)).
			Where(sql.InInts(serviceendpoint.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(su.endpoints) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range su.endpoints {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(serviceendpoint.FieldID, eid)
			}
			query, args := builder.Update(service.EndpointsTable).
				Set(service.EndpointsColumn, id).
				Where(sql.And(p, sql.IsNull(service.EndpointsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(su.endpoints) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"endpoints\" %v already connected to a different \"Service\"", keys(su.endpoints))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// ServiceUpdateOne is the builder for updating a single Service entity.
type ServiceUpdateOne struct {
	config
	id string

	update_time              *time.Time
	name                     *string
	external_id              *string
	clearexternal_id         bool
	status                   *string
	_type                    map[string]struct{}
	downstream               map[string]struct{}
	upstream                 map[string]struct{}
	properties               map[string]struct{}
	termination_points       map[string]struct{}
	links                    map[string]struct{}
	customer                 map[string]struct{}
	endpoints                map[string]struct{}
	clearedType              bool
	removedDownstream        map[string]struct{}
	removedUpstream          map[string]struct{}
	removedProperties        map[string]struct{}
	removedTerminationPoints map[string]struct{}
	removedLinks             map[string]struct{}
	removedCustomer          map[string]struct{}
	removedEndpoints         map[string]struct{}
}

// SetName sets the name field.
func (suo *ServiceUpdateOne) SetName(s string) *ServiceUpdateOne {
	suo.name = &s
	return suo
}

// SetExternalID sets the external_id field.
func (suo *ServiceUpdateOne) SetExternalID(s string) *ServiceUpdateOne {
	suo.external_id = &s
	return suo
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableExternalID(s *string) *ServiceUpdateOne {
	if s != nil {
		suo.SetExternalID(*s)
	}
	return suo
}

// ClearExternalID clears the value of external_id.
func (suo *ServiceUpdateOne) ClearExternalID() *ServiceUpdateOne {
	suo.external_id = nil
	suo.clearexternal_id = true
	return suo
}

// SetStatus sets the status field.
func (suo *ServiceUpdateOne) SetStatus(s string) *ServiceUpdateOne {
	suo.status = &s
	return suo
}

// SetTypeID sets the type edge to ServiceType by id.
func (suo *ServiceUpdateOne) SetTypeID(id string) *ServiceUpdateOne {
	if suo._type == nil {
		suo._type = make(map[string]struct{})
	}
	suo._type[id] = struct{}{}
	return suo
}

// SetType sets the type edge to ServiceType.
func (suo *ServiceUpdateOne) SetType(s *ServiceType) *ServiceUpdateOne {
	return suo.SetTypeID(s.ID)
}

// AddDownstreamIDs adds the downstream edge to Service by ids.
func (suo *ServiceUpdateOne) AddDownstreamIDs(ids ...string) *ServiceUpdateOne {
	if suo.downstream == nil {
		suo.downstream = make(map[string]struct{})
	}
	for i := range ids {
		suo.downstream[ids[i]] = struct{}{}
	}
	return suo
}

// AddDownstream adds the downstream edges to Service.
func (suo *ServiceUpdateOne) AddDownstream(s ...*Service) *ServiceUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddDownstreamIDs(ids...)
}

// AddUpstreamIDs adds the upstream edge to Service by ids.
func (suo *ServiceUpdateOne) AddUpstreamIDs(ids ...string) *ServiceUpdateOne {
	if suo.upstream == nil {
		suo.upstream = make(map[string]struct{})
	}
	for i := range ids {
		suo.upstream[ids[i]] = struct{}{}
	}
	return suo
}

// AddUpstream adds the upstream edges to Service.
func (suo *ServiceUpdateOne) AddUpstream(s ...*Service) *ServiceUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddUpstreamIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (suo *ServiceUpdateOne) AddPropertyIDs(ids ...string) *ServiceUpdateOne {
	if suo.properties == nil {
		suo.properties = make(map[string]struct{})
	}
	for i := range ids {
		suo.properties[ids[i]] = struct{}{}
	}
	return suo
}

// AddProperties adds the properties edges to Property.
func (suo *ServiceUpdateOne) AddProperties(p ...*Property) *ServiceUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return suo.AddPropertyIDs(ids...)
}

// AddTerminationPointIDs adds the termination_points edge to Equipment by ids.
func (suo *ServiceUpdateOne) AddTerminationPointIDs(ids ...string) *ServiceUpdateOne {
	if suo.termination_points == nil {
		suo.termination_points = make(map[string]struct{})
	}
	for i := range ids {
		suo.termination_points[ids[i]] = struct{}{}
	}
	return suo
}

// AddTerminationPoints adds the termination_points edges to Equipment.
func (suo *ServiceUpdateOne) AddTerminationPoints(e ...*Equipment) *ServiceUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return suo.AddTerminationPointIDs(ids...)
}

// AddLinkIDs adds the links edge to Link by ids.
func (suo *ServiceUpdateOne) AddLinkIDs(ids ...string) *ServiceUpdateOne {
	if suo.links == nil {
		suo.links = make(map[string]struct{})
	}
	for i := range ids {
		suo.links[ids[i]] = struct{}{}
	}
	return suo
}

// AddLinks adds the links edges to Link.
func (suo *ServiceUpdateOne) AddLinks(l ...*Link) *ServiceUpdateOne {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return suo.AddLinkIDs(ids...)
}

// AddCustomerIDs adds the customer edge to Customer by ids.
func (suo *ServiceUpdateOne) AddCustomerIDs(ids ...string) *ServiceUpdateOne {
	if suo.customer == nil {
		suo.customer = make(map[string]struct{})
	}
	for i := range ids {
		suo.customer[ids[i]] = struct{}{}
	}
	return suo
}

// AddCustomer adds the customer edges to Customer.
func (suo *ServiceUpdateOne) AddCustomer(c ...*Customer) *ServiceUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return suo.AddCustomerIDs(ids...)
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (suo *ServiceUpdateOne) AddEndpointIDs(ids ...string) *ServiceUpdateOne {
	if suo.endpoints == nil {
		suo.endpoints = make(map[string]struct{})
	}
	for i := range ids {
		suo.endpoints[ids[i]] = struct{}{}
	}
	return suo
}

// AddEndpoints adds the endpoints edges to ServiceEndpoint.
func (suo *ServiceUpdateOne) AddEndpoints(s ...*ServiceEndpoint) *ServiceUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddEndpointIDs(ids...)
}

// ClearType clears the type edge to ServiceType.
func (suo *ServiceUpdateOne) ClearType() *ServiceUpdateOne {
	suo.clearedType = true
	return suo
}

// RemoveDownstreamIDs removes the downstream edge to Service by ids.
func (suo *ServiceUpdateOne) RemoveDownstreamIDs(ids ...string) *ServiceUpdateOne {
	if suo.removedDownstream == nil {
		suo.removedDownstream = make(map[string]struct{})
	}
	for i := range ids {
		suo.removedDownstream[ids[i]] = struct{}{}
	}
	return suo
}

// RemoveDownstream removes downstream edges to Service.
func (suo *ServiceUpdateOne) RemoveDownstream(s ...*Service) *ServiceUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveDownstreamIDs(ids...)
}

// RemoveUpstreamIDs removes the upstream edge to Service by ids.
func (suo *ServiceUpdateOne) RemoveUpstreamIDs(ids ...string) *ServiceUpdateOne {
	if suo.removedUpstream == nil {
		suo.removedUpstream = make(map[string]struct{})
	}
	for i := range ids {
		suo.removedUpstream[ids[i]] = struct{}{}
	}
	return suo
}

// RemoveUpstream removes upstream edges to Service.
func (suo *ServiceUpdateOne) RemoveUpstream(s ...*Service) *ServiceUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveUpstreamIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (suo *ServiceUpdateOne) RemovePropertyIDs(ids ...string) *ServiceUpdateOne {
	if suo.removedProperties == nil {
		suo.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		suo.removedProperties[ids[i]] = struct{}{}
	}
	return suo
}

// RemoveProperties removes properties edges to Property.
func (suo *ServiceUpdateOne) RemoveProperties(p ...*Property) *ServiceUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return suo.RemovePropertyIDs(ids...)
}

// RemoveTerminationPointIDs removes the termination_points edge to Equipment by ids.
func (suo *ServiceUpdateOne) RemoveTerminationPointIDs(ids ...string) *ServiceUpdateOne {
	if suo.removedTerminationPoints == nil {
		suo.removedTerminationPoints = make(map[string]struct{})
	}
	for i := range ids {
		suo.removedTerminationPoints[ids[i]] = struct{}{}
	}
	return suo
}

// RemoveTerminationPoints removes termination_points edges to Equipment.
func (suo *ServiceUpdateOne) RemoveTerminationPoints(e ...*Equipment) *ServiceUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return suo.RemoveTerminationPointIDs(ids...)
}

// RemoveLinkIDs removes the links edge to Link by ids.
func (suo *ServiceUpdateOne) RemoveLinkIDs(ids ...string) *ServiceUpdateOne {
	if suo.removedLinks == nil {
		suo.removedLinks = make(map[string]struct{})
	}
	for i := range ids {
		suo.removedLinks[ids[i]] = struct{}{}
	}
	return suo
}

// RemoveLinks removes links edges to Link.
func (suo *ServiceUpdateOne) RemoveLinks(l ...*Link) *ServiceUpdateOne {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return suo.RemoveLinkIDs(ids...)
}

// RemoveCustomerIDs removes the customer edge to Customer by ids.
func (suo *ServiceUpdateOne) RemoveCustomerIDs(ids ...string) *ServiceUpdateOne {
	if suo.removedCustomer == nil {
		suo.removedCustomer = make(map[string]struct{})
	}
	for i := range ids {
		suo.removedCustomer[ids[i]] = struct{}{}
	}
	return suo
}

// RemoveCustomer removes customer edges to Customer.
func (suo *ServiceUpdateOne) RemoveCustomer(c ...*Customer) *ServiceUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return suo.RemoveCustomerIDs(ids...)
}

// RemoveEndpointIDs removes the endpoints edge to ServiceEndpoint by ids.
func (suo *ServiceUpdateOne) RemoveEndpointIDs(ids ...string) *ServiceUpdateOne {
	if suo.removedEndpoints == nil {
		suo.removedEndpoints = make(map[string]struct{})
	}
	for i := range ids {
		suo.removedEndpoints[ids[i]] = struct{}{}
	}
	return suo
}

// RemoveEndpoints removes endpoints edges to ServiceEndpoint.
func (suo *ServiceUpdateOne) RemoveEndpoints(s ...*ServiceEndpoint) *ServiceUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveEndpointIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (suo *ServiceUpdateOne) Save(ctx context.Context) (*Service, error) {
	if suo.update_time == nil {
		v := service.UpdateDefaultUpdateTime()
		suo.update_time = &v
	}
	if suo.name != nil {
		if err := service.NameValidator(*suo.name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if suo.external_id != nil {
		if err := service.ExternalIDValidator(*suo.external_id); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"external_id\": %v", err)
		}
	}
	if len(suo._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if suo.clearedType && suo._type == nil {
		return nil, errors.New("ent: clearing a unique edge \"type\"")
	}
	return suo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (suo *ServiceUpdateOne) SaveX(ctx context.Context) *Service {
	s, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return s
}

// Exec executes the query on the entity.
func (suo *ServiceUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *ServiceUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (suo *ServiceUpdateOne) sqlSave(ctx context.Context) (s *Service, err error) {
	var (
		builder  = sql.Dialect(suo.driver.Dialect())
		selector = builder.Select(service.Columns...).From(builder.Table(service.Table))
	)
	service.ID(suo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = suo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		s = &Service{config: suo.config}
		if err := s.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into Service: %v", err)
		}
		id = s.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("Service with id: %v", suo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one Service with the same id: %v", suo.id)
	}

	tx, err := suo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(service.Table)
	)
	updater = updater.Where(sql.InInts(service.FieldID, ids...))
	if value := suo.update_time; value != nil {
		updater.Set(service.FieldUpdateTime, *value)
		s.UpdateTime = *value
	}
	if value := suo.name; value != nil {
		updater.Set(service.FieldName, *value)
		s.Name = *value
	}
	if value := suo.external_id; value != nil {
		updater.Set(service.FieldExternalID, *value)
		s.ExternalID = value
	}
	if suo.clearexternal_id {
		s.ExternalID = nil
		updater.SetNull(service.FieldExternalID)
	}
	if value := suo.status; value != nil {
		updater.Set(service.FieldStatus, *value)
		s.Status = *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if suo.clearedType {
		query, args := builder.Update(service.TypeTable).
			SetNull(service.TypeColumn).
			Where(sql.InInts(servicetype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(suo._type) > 0 {
		for eid := range suo._type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(service.TypeTable).
				Set(service.TypeColumn, eid).
				Where(sql.InInts(service.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(suo.removedDownstream) > 0 {
		eids := make([]int, len(suo.removedDownstream))
		for eid := range suo.removedDownstream {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Delete(service.DownstreamTable).
			Where(sql.InInts(service.DownstreamPrimaryKey[1], ids...)).
			Where(sql.InInts(service.DownstreamPrimaryKey[0], eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(suo.downstream) > 0 {
		values := make([][]int, 0, len(ids))
		for _, id := range ids {
			for eid := range suo.downstream {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				values = append(values, []int{id, eid})
			}
		}
		builder := builder.Insert(service.DownstreamTable).
			Columns(service.DownstreamPrimaryKey[1], service.DownstreamPrimaryKey[0])
		for _, v := range values {
			builder.Values(v[0], v[1])
		}
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(suo.removedUpstream) > 0 {
		eids := make([]int, len(suo.removedUpstream))
		for eid := range suo.removedUpstream {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Delete(service.UpstreamTable).
			Where(sql.InInts(service.UpstreamPrimaryKey[0], ids...)).
			Where(sql.InInts(service.UpstreamPrimaryKey[1], eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(suo.upstream) > 0 {
		values := make([][]int, 0, len(ids))
		for _, id := range ids {
			for eid := range suo.upstream {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				values = append(values, []int{id, eid})
			}
		}
		builder := builder.Insert(service.UpstreamTable).
			Columns(service.UpstreamPrimaryKey[0], service.UpstreamPrimaryKey[1])
		for _, v := range values {
			builder.Values(v[0], v[1])
		}
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(suo.removedProperties) > 0 {
		eids := make([]int, len(suo.removedProperties))
		for eid := range suo.removedProperties {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(service.PropertiesTable).
			SetNull(service.PropertiesColumn).
			Where(sql.InInts(service.PropertiesColumn, ids...)).
			Where(sql.InInts(property.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(suo.properties) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range suo.properties {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(suo.properties) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"Service\"", keys(suo.properties))})
			}
		}
	}
	if len(suo.removedTerminationPoints) > 0 {
		eids := make([]int, len(suo.removedTerminationPoints))
		for eid := range suo.removedTerminationPoints {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Delete(service.TerminationPointsTable).
			Where(sql.InInts(service.TerminationPointsPrimaryKey[0], ids...)).
			Where(sql.InInts(service.TerminationPointsPrimaryKey[1], eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(suo.termination_points) > 0 {
		values := make([][]int, 0, len(ids))
		for _, id := range ids {
			for eid := range suo.termination_points {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				values = append(values, []int{id, eid})
			}
		}
		builder := builder.Insert(service.TerminationPointsTable).
			Columns(service.TerminationPointsPrimaryKey[0], service.TerminationPointsPrimaryKey[1])
		for _, v := range values {
			builder.Values(v[0], v[1])
		}
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(suo.removedLinks) > 0 {
		eids := make([]int, len(suo.removedLinks))
		for eid := range suo.removedLinks {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Delete(service.LinksTable).
			Where(sql.InInts(service.LinksPrimaryKey[0], ids...)).
			Where(sql.InInts(service.LinksPrimaryKey[1], eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(suo.links) > 0 {
		values := make([][]int, 0, len(ids))
		for _, id := range ids {
			for eid := range suo.links {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				values = append(values, []int{id, eid})
			}
		}
		builder := builder.Insert(service.LinksTable).
			Columns(service.LinksPrimaryKey[0], service.LinksPrimaryKey[1])
		for _, v := range values {
			builder.Values(v[0], v[1])
		}
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(suo.removedCustomer) > 0 {
		eids := make([]int, len(suo.removedCustomer))
		for eid := range suo.removedCustomer {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Delete(service.CustomerTable).
			Where(sql.InInts(service.CustomerPrimaryKey[0], ids...)).
			Where(sql.InInts(service.CustomerPrimaryKey[1], eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(suo.customer) > 0 {
		values := make([][]int, 0, len(ids))
		for _, id := range ids {
			for eid := range suo.customer {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				values = append(values, []int{id, eid})
			}
		}
		builder := builder.Insert(service.CustomerTable).
			Columns(service.CustomerPrimaryKey[0], service.CustomerPrimaryKey[1])
		for _, v := range values {
			builder.Values(v[0], v[1])
		}
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(suo.removedEndpoints) > 0 {
		eids := make([]int, len(suo.removedEndpoints))
		for eid := range suo.removedEndpoints {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(service.EndpointsTable).
			SetNull(service.EndpointsColumn).
			Where(sql.InInts(service.EndpointsColumn, ids...)).
			Where(sql.InInts(serviceendpoint.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(suo.endpoints) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range suo.endpoints {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(serviceendpoint.FieldID, eid)
			}
			query, args := builder.Update(service.EndpointsTable).
				Set(service.EndpointsColumn, id).
				Where(sql.And(p, sql.IsNull(service.EndpointsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(suo.endpoints) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"endpoints\" %v already connected to a different \"Service\"", keys(suo.endpoints))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return s, nil
}
