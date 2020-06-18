// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent/permissionspolicy"
)

// PermissionsPolicy is the model entity for the PermissionsPolicy schema.
type PermissionsPolicy struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// IsGlobal holds the value of the "is_global" field.
	IsGlobal bool `json:"is_global,omitempty"`
	// InventoryPolicy holds the value of the "inventory_policy" field.
	InventoryPolicy *models.InventoryPolicyInput `json:"inventory_policy,omitempty"`
	// WorkforcePolicy holds the value of the "workforce_policy" field.
	WorkforcePolicy *models.WorkforcePolicyInput `json:"workforce_policy,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PermissionsPolicyQuery when eager-loading is set.
	Edges PermissionsPolicyEdges `json:"edges"`
}

// PermissionsPolicyEdges holds the relations/edges for other nodes in the graph.
type PermissionsPolicyEdges struct {
	// Groups holds the value of the groups edge.
	Groups []*UsersGroup
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// GroupsOrErr returns the Groups value or an error if the edge
// was not loaded in eager-loading.
func (e PermissionsPolicyEdges) GroupsOrErr() ([]*UsersGroup, error) {
	if e.loadedTypes[0] {
		return e.Groups, nil
	}
	return nil, &NotLoadedError{edge: "groups"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*PermissionsPolicy) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // name
		&sql.NullString{}, // description
		&sql.NullBool{},   // is_global
		&[]byte{},         // inventory_policy
		&[]byte{},         // workforce_policy
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the PermissionsPolicy fields.
func (pp *PermissionsPolicy) assignValues(values ...interface{}) error {
	if m, n := len(values), len(permissionspolicy.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	pp.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		pp.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		pp.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[2])
	} else if value.Valid {
		pp.Name = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field description", values[3])
	} else if value.Valid {
		pp.Description = value.String
	}
	if value, ok := values[4].(*sql.NullBool); !ok {
		return fmt.Errorf("unexpected type %T for field is_global", values[4])
	} else if value.Valid {
		pp.IsGlobal = value.Bool
	}

	if value, ok := values[5].(*[]byte); !ok {
		return fmt.Errorf("unexpected type %T for field inventory_policy", values[5])
	} else if value != nil && len(*value) > 0 {
		if err := json.Unmarshal(*value, &pp.InventoryPolicy); err != nil {
			return fmt.Errorf("unmarshal field inventory_policy: %v", err)
		}
	}

	if value, ok := values[6].(*[]byte); !ok {
		return fmt.Errorf("unexpected type %T for field workforce_policy", values[6])
	} else if value != nil && len(*value) > 0 {
		if err := json.Unmarshal(*value, &pp.WorkforcePolicy); err != nil {
			return fmt.Errorf("unmarshal field workforce_policy: %v", err)
		}
	}
	return nil
}

// QueryGroups queries the groups edge of the PermissionsPolicy.
func (pp *PermissionsPolicy) QueryGroups() *UsersGroupQuery {
	return (&PermissionsPolicyClient{config: pp.config}).QueryGroups(pp)
}

// Update returns a builder for updating this PermissionsPolicy.
// Note that, you need to call PermissionsPolicy.Unwrap() before calling this method, if this PermissionsPolicy
// was returned from a transaction, and the transaction was committed or rolled back.
func (pp *PermissionsPolicy) Update() *PermissionsPolicyUpdateOne {
	return (&PermissionsPolicyClient{config: pp.config}).UpdateOne(pp)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (pp *PermissionsPolicy) Unwrap() *PermissionsPolicy {
	tx, ok := pp.config.driver.(*txDriver)
	if !ok {
		panic("ent: PermissionsPolicy is not a transactional entity")
	}
	pp.config.driver = tx.drv
	return pp
}

// String implements the fmt.Stringer.
func (pp *PermissionsPolicy) String() string {
	var builder strings.Builder
	builder.WriteString("PermissionsPolicy(")
	builder.WriteString(fmt.Sprintf("id=%v", pp.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(pp.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(pp.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(pp.Name)
	builder.WriteString(", description=")
	builder.WriteString(pp.Description)
	builder.WriteString(", is_global=")
	builder.WriteString(fmt.Sprintf("%v", pp.IsGlobal))
	builder.WriteString(", inventory_policy=")
	builder.WriteString(fmt.Sprintf("%v", pp.InventoryPolicy))
	builder.WriteString(", workforce_policy=")
	builder.WriteString(fmt.Sprintf("%v", pp.WorkforcePolicy))
	builder.WriteByte(')')
	return builder.String()
}

// PermissionsPolicies is a parsable slice of PermissionsPolicy.
type PermissionsPolicies []*PermissionsPolicy

func (pp PermissionsPolicies) config(cfg config) {
	for _i := range pp {
		pp[_i].config = cfg
	}
}
