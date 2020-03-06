// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

func toStatusPointer(status user.Status) *user.Status {
	return &status
}

func TestEditUser(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	prepareUserData(t, ctx, r.client)

	u, err := viewer.UserFromContext(ctx)
	require.NoError(t, err)
	require.Equal(t, user.StatusActive, u.Status)
	require.Empty(t, u.FirstName)

	mr := r.Mutation()
	u, err = mr.EditUser(ctx, models.EditUserInput{ID: u.ID, Status: toStatusPointer(user.StatusDeactivated), FirstName: pointer.ToString("John")})
	require.NoError(t, err)
	require.Equal(t, user.StatusDeactivated, u.Status)
	require.Equal(t, "John", u.FirstName)
}
