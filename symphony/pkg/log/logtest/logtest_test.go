// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logtest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestLogger(t *testing.T) {
	assert.Implements(t, (*TestingT)(nil), &testing.T{})
	assert.Implements(t, (*TestingT)(nil), &testing.B{})
	logger := NewTestLogger(t)
	assert.Equal(t, logger.Background(), logger.For(context.Background()))
}
