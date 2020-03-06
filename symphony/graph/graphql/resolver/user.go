// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer"
)

type userResolver struct{}

func (r userResolver) ProfilePhoto(ctx context.Context, user *ent.User) (*ent.File, error) {
	profilePhoto, err := user.Edges.ProfilePhotoOrErr()
	if ent.IsNotLoaded(err) {
		profilePhoto, err = user.QueryProfilePhoto().Only(ctx)
	}
	return profilePhoto, ent.MaskNotFound(err)
}

func (r queryResolver) Users(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int) (*ent.UserConnection, error) {
	return r.ClientFrom(ctx).User.Query().
		Paginate(ctx, after, first, before, last)
}

func (r mutationResolver) EditUser(ctx context.Context, input models.EditUserInput) (*ent.User, error) {
	client := r.ClientFrom(ctx)
	u, err := viewer.UserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user from context: %w", err)
	}

	u, err = client.User.UpdateOne(u).
		SetNillableFirstName(input.FirstName).
		SetNillableLastName(input.LastName).
		SetNillableStatus(input.Status).
		SetNillableRole(input.Role).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("edit user: %w", err)
	}
	return u, nil
}
