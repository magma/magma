// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package role

import (
	"fmt"
)

// Role defines a user role.
type Role int

// Allowed user roles.
const (
	UserRole  Role = 0
	SuperUser Role = 3
)

// Validate role value is a valid one.
func (r Role) Validate() error {
	switch r {
	case UserRole, SuperUser:
		return nil
	default:
		return fmt.Errorf("invalid role value: %d", r)
	}
}

// ValidateValue can be used as an int field validator for roles.
func ValidateValue(v int) error {
	return Role(v).Validate()
}
