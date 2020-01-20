// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleValidate(t *testing.T) {
	for _, role := range []Role{UserRole, SuperUser} {
		err := role.Validate()
		assert.NoErrorf(t, err, "role %d must be valid", role)
	}
	err := Role(-1).Validate()
	assert.EqualError(t, err, "invalid role value: -1")
}
