// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewertest_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	tenantName := "facebook"
	userName := "fbuser@fb.com"
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(
		context.Background(),
		c,
		viewertest.WithTenant(tenantName),
		viewertest.WithUser(userName),
	)
	got := viewer.FromContext(ctx)
	assert.Equal(t, tenantName, got.Tenant)
	u := got.User()
	assert.Equal(t, userName, u.AuthID)
	assert.Equal(t, userName, u.Email)
}
