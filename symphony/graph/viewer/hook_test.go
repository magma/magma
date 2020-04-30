// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewer_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserInViewer(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	ctx2 := viewertest.NewContext(context.Background(), c, viewertest.WithUser("tester2@example.com"))

	u := viewer.FromContext(ctx).(*viewer.UserViewer).User()
	require.Equal(t, viewertest.DefaultUser, u.Email)
	u2 := viewer.FromContext(ctx2).(*viewer.UserViewer).User()
	require.Equal(t, "tester2@example.com", u2.Email)

	err := c.User.UpdateOneID(u.ID).SetEmail("new_tester@example.com").Exec(ctx)
	require.NoError(t, err)

	u = viewer.FromContext(ctx).(*viewer.UserViewer).User()
	require.Equal(t, "new_tester@example.com", u.Email)
	u2 = viewer.FromContext(ctx2).(*viewer.UserViewer).User()
	require.Equal(t, "tester2@example.com", u2.Email)
}
