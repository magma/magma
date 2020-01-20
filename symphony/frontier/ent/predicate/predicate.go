// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package predicate

import (
	"github.com/facebookincubator/ent/dialect/sql"
)

// AuditLog is the predicate function for auditlog builders.
type AuditLog func(*sql.Selector)

// Tenant is the predicate function for tenant builders.
type Tenant func(*sql.Selector)

// Token is the predicate function for token builders.
type Token func(*sql.Selector)

// User is the predicate function for user builders.
type User func(*sql.Selector)
