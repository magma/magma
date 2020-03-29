// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/usersgroup"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

type usersGroupResolver struct{}

func (r queryResolver) UsersGroups(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int) (*ent.UsersGroupConnection, error) {
	return r.ClientFrom(ctx).UsersGroup.Query().
		Paginate(ctx, after, first, before, last)
}

func (usersGroupResolver) Members(ctx context.Context, obj *ent.UsersGroup) ([]*ent.User, error) {
	return obj.QueryMembers().All(ctx)
}

func (r mutationResolver) AddUsersGroup(ctx context.Context, input models.AddUsersGroupInput) (*ent.UsersGroup, error) {
	client := r.ClientFrom(ctx)

	return client.UsersGroup.Create().
		SetName(input.Name).
		SetNillableDescription(input.Description).
		SetStatus(usersgroup.StatusACTIVE).
		Save(ctx)
}

func (r mutationResolver) EditUsersGroup(ctx context.Context, input models.EditUsersGroupInput) (*ent.UsersGroup, error) {
	client := r.ClientFrom(ctx)

	m := client.UsersGroup.UpdateOneID(input.ID).
		SetNillableDescription(input.Description).
		SetNillableStatus(input.Status)

	if input.Name != nil {
		m = m.SetName(*input.Name)
	}

	return m.Save(ctx)
}
