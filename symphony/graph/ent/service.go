// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
)

// Service is the model entity for the Service schema.
type Service struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// ExternalID holds the value of the "external_id" field.
	ExternalID *string `json:"external_id,omitempty"`
	// Status holds the value of the "status" field.
	Status string `json:"status,omitempty"`
	// DiscoveryMethod holds the value of the "discovery_method" field.
	DiscoveryMethod service.DiscoveryMethod `json:"discovery_method,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ServiceQuery when eager-loading is set.
	Edges        ServiceEdges `json:"edges"`
	service_type *int
}

// ServiceEdges holds the relations/edges for other nodes in the graph.
type ServiceEdges struct {
	// Type holds the value of the type edge.
	Type *ServiceType
	// Downstream holds the value of the downstream edge.
	Downstream []*Service
	// Upstream holds the value of the upstream edge.
	Upstream []*Service
	// Properties holds the value of the properties edge.
	Properties []*Property
	// Links holds the value of the links edge.
	Links []*Link
	// Customer holds the value of the customer edge.
	Customer []*Customer
	// Endpoints holds the value of the endpoints edge.
	Endpoints []*ServiceEndpoint
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [7]bool
}

// TypeOrErr returns the Type value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceEdges) TypeOrErr() (*ServiceType, error) {
	if e.loadedTypes[0] {
		if e.Type == nil {
			// The edge type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: servicetype.Label}
		}
		return e.Type, nil
	}
	return nil, &NotLoadedError{edge: "type"}
}

// DownstreamOrErr returns the Downstream value or an error if the edge
// was not loaded in eager-loading.
func (e ServiceEdges) DownstreamOrErr() ([]*Service, error) {
	if e.loadedTypes[1] {
		return e.Downstream, nil
	}
	return nil, &NotLoadedError{edge: "downstream"}
}

// UpstreamOrErr returns the Upstream value or an error if the edge
// was not loaded in eager-loading.
func (e ServiceEdges) UpstreamOrErr() ([]*Service, error) {
	if e.loadedTypes[2] {
		return e.Upstream, nil
	}
	return nil, &NotLoadedError{edge: "upstream"}
}

// PropertiesOrErr returns the Properties value or an error if the edge
// was not loaded in eager-loading.
func (e ServiceEdges) PropertiesOrErr() ([]*Property, error) {
	if e.loadedTypes[3] {
		return e.Properties, nil
	}
	return nil, &NotLoadedError{edge: "properties"}
}

// LinksOrErr returns the Links value or an error if the edge
// was not loaded in eager-loading.
func (e ServiceEdges) LinksOrErr() ([]*Link, error) {
	if e.loadedTypes[4] {
		return e.Links, nil
	}
	return nil, &NotLoadedError{edge: "links"}
}

// CustomerOrErr returns the Customer value or an error if the edge
// was not loaded in eager-loading.
func (e ServiceEdges) CustomerOrErr() ([]*Customer, error) {
	if e.loadedTypes[5] {
		return e.Customer, nil
	}
	return nil, &NotLoadedError{edge: "customer"}
}

// EndpointsOrErr returns the Endpoints value or an error if the edge
// was not loaded in eager-loading.
func (e ServiceEdges) EndpointsOrErr() ([]*ServiceEndpoint, error) {
	if e.loadedTypes[6] {
		return e.Endpoints, nil
	}
	return nil, &NotLoadedError{edge: "endpoints"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Service) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // name
		&sql.NullString{}, // external_id
		&sql.NullString{}, // status
		&sql.NullString{}, // discovery_method
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*Service) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // service_type
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Service fields.
func (s *Service) assignValues(values ...interface{}) error {
	if m, n := len(values), len(service.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	s.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		s.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		s.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[2])
	} else if value.Valid {
		s.Name = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field external_id", values[3])
	} else if value.Valid {
		s.ExternalID = new(string)
		*s.ExternalID = value.String
	}
	if value, ok := values[4].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field status", values[4])
	} else if value.Valid {
		s.Status = value.String
	}
	if value, ok := values[5].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field discovery_method", values[5])
	} else if value.Valid {
		s.DiscoveryMethod = service.DiscoveryMethod(value.String)
	}
	values = values[6:]
	if len(values) == len(service.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field service_type", value)
		} else if value.Valid {
			s.service_type = new(int)
			*s.service_type = int(value.Int64)
		}
	}
	return nil
}

// QueryType queries the type edge of the Service.
func (s *Service) QueryType() *ServiceTypeQuery {
	return (&ServiceClient{config: s.config}).QueryType(s)
}

// QueryDownstream queries the downstream edge of the Service.
func (s *Service) QueryDownstream() *ServiceQuery {
	return (&ServiceClient{config: s.config}).QueryDownstream(s)
}

// QueryUpstream queries the upstream edge of the Service.
func (s *Service) QueryUpstream() *ServiceQuery {
	return (&ServiceClient{config: s.config}).QueryUpstream(s)
}

// QueryProperties queries the properties edge of the Service.
func (s *Service) QueryProperties() *PropertyQuery {
	return (&ServiceClient{config: s.config}).QueryProperties(s)
}

// QueryLinks queries the links edge of the Service.
func (s *Service) QueryLinks() *LinkQuery {
	return (&ServiceClient{config: s.config}).QueryLinks(s)
}

// QueryCustomer queries the customer edge of the Service.
func (s *Service) QueryCustomer() *CustomerQuery {
	return (&ServiceClient{config: s.config}).QueryCustomer(s)
}

// QueryEndpoints queries the endpoints edge of the Service.
func (s *Service) QueryEndpoints() *ServiceEndpointQuery {
	return (&ServiceClient{config: s.config}).QueryEndpoints(s)
}

// Update returns a builder for updating this Service.
// Note that, you need to call Service.Unwrap() before calling this method, if this Service
// was returned from a transaction, and the transaction was committed or rolled back.
func (s *Service) Update() *ServiceUpdateOne {
	return (&ServiceClient{config: s.config}).UpdateOne(s)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (s *Service) Unwrap() *Service {
	tx, ok := s.config.driver.(*txDriver)
	if !ok {
		panic("ent: Service is not a transactional entity")
	}
	s.config.driver = tx.drv
	return s
}

// String implements the fmt.Stringer.
func (s *Service) String() string {
	var builder strings.Builder
	builder.WriteString("Service(")
	builder.WriteString(fmt.Sprintf("id=%v", s.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(s.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(s.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(s.Name)
	if v := s.ExternalID; v != nil {
		builder.WriteString(", external_id=")
		builder.WriteString(*v)
	}
	builder.WriteString(", status=")
	builder.WriteString(s.Status)
	builder.WriteString(", discovery_method=")
	builder.WriteString(fmt.Sprintf("%v", s.DiscoveryMethod))
	builder.WriteByte(')')
	return builder.String()
}

// Services is a parsable slice of Service.
type Services []*Service

func (s Services) config(cfg config) {
	for _i := range s {
		s[_i].config = cfg
	}
}
