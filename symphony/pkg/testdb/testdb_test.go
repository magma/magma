// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	db, name, err := Open()
	assert.NotNil(t, db)
	assert.NotEmpty(t, name)
	assert.NoError(t, err)
}
