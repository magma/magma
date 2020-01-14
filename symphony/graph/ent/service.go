// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/service"
)

// Service is the model entity for the Service schema.
type Service struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
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
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Service) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},
		&sql.NullTime{},
		&sql.NullTime{},
		&sql.NullString{},
		&sql.NullString{},
		&sql.NullString{},
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Service fields.
func (s *Service) assignValues(values ...interface{}) error {
	if m, n := len(values), len(service.Columns); m != n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	s.ID = strconv.FormatInt(value.Int64, 10)
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
	return nil
}

// QueryType queries the type edge of the Service.
func (s *Service) QueryType() *ServiceTypeQuery {
	return (&ServiceClient{s.config}).QueryType(s)
}

// QueryDownstream queries the downstream edge of the Service.
func (s *Service) QueryDownstream() *ServiceQuery {
	return (&ServiceClient{s.config}).QueryDownstream(s)
}

// QueryUpstream queries the upstream edge of the Service.
func (s *Service) QueryUpstream() *ServiceQuery {
	return (&ServiceClient{s.config}).QueryUpstream(s)
}

// QueryProperties queries the properties edge of the Service.
func (s *Service) QueryProperties() *PropertyQuery {
	return (&ServiceClient{s.config}).QueryProperties(s)
}

// QueryLinks queries the links edge of the Service.
func (s *Service) QueryLinks() *LinkQuery {
	return (&ServiceClient{s.config}).QueryLinks(s)
}

// QueryCustomer queries the customer edge of the Service.
func (s *Service) QueryCustomer() *CustomerQuery {
	return (&ServiceClient{s.config}).QueryCustomer(s)
}

// QueryEndpoints queries the endpoints edge of the Service.
func (s *Service) QueryEndpoints() *ServiceEndpointQuery {
	return (&ServiceClient{s.config}).QueryEndpoints(s)
}

// Update returns a builder for updating this Service.
// Note that, you need to call Service.Unwrap() before calling this method, if this Service
// was returned from a transaction, and the transaction was committed or rolled back.
func (s *Service) Update() *ServiceUpdateOne {
	return (&ServiceClient{s.config}).UpdateOne(s)
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
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (s *Service) id() int {
	id, _ := strconv.Atoi(s.ID)
	return id
}

// Services is a parsable slice of Service.
type Services []*Service

func (s Services) config(cfg config) {
	for _i := range s {
		s[_i].config = cfg
	}
}
