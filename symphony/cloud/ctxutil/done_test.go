// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ctxutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoneCtx(t *testing.T) {
	t.Parallel()
	ctx := DoneCtx()
	assert.Implements(t, (*context.Context)(nil), ctx)
	select {
	case <-ctx.Done():
	default:
		assert.Fail(t, "unreadable done channel")
	}
	assert.EqualError(t, ctx.Err(), ErrDone.Error())
	deadline, ok := ctx.Deadline()
	assert.True(t, deadline.IsZero())
	assert.False(t, ok)
}
