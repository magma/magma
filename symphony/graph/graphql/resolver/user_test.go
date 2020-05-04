// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func toStatusPointer(status user.Status) *user.Status {
	return &status
}

func TestEditUser(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	u := viewer.FromContext(ctx).(*viewer.UserViewer).User()
	require.Equal(t, user.StatusACTIVE, u.Status)
	require.Empty(t, u.FirstName)

	mr := r.Mutation()
	u, err := mr.EditUser(ctx, models.EditUserInput{ID: u.ID, Status: toStatusPointer(user.StatusDEACTIVATED), FirstName: pointer.ToString("John")})
	require.NoError(t, err)
	require.Equal(t, user.StatusDEACTIVATED, u.Status)
	require.Equal(t, "John", u.FirstName)
}

func TestAddAndDeleteProfileImage(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	u := viewer.FromContext(ctx).(*viewer.UserViewer).User()

	mr, ur := r.Mutation(), r.User()
	file1, err := mr.AddImage(ctx, models.AddImageInput{
		EntityType:  models.ImageEntityUser,
		EntityID:    u.ID,
		ImgKey:      uuid.New().String(),
		FileName:    "profile_photo.png",
		FileSize:    123,
		Modified:    time.Now(),
		ContentType: "image/png",
	})
	require.NoError(t, err)
	file, err := ur.ProfilePhoto(ctx, u)
	require.NoError(t, err)
	require.Equal(t, "profile_photo.png", file.Name)

	_, err = mr.DeleteImage(ctx, models.ImageEntityUser, u.ID, file1.ID)
	require.NoError(t, err)

	file, err = ur.ProfilePhoto(ctx, u)
	require.NoError(t, err)
	require.Nil(t, file)
}
