// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/usersgroup"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type usersGroupResolver struct{}

func (r queryResolver) UsersGroups(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int) (*ent.UsersGroupConnection, error) {
	return r.ClientFrom(ctx).UsersGroup.Query().
		Paginate(ctx, after, first, before, last)
}

func (usersGroupResolver) Members(ctx context.Context, obj *ent.UsersGroup) ([]*ent.User, error) {
	return obj.QueryMembers().All(ctx)
}

func (usersGroupResolver) Policies(ctx context.Context, obj *ent.UsersGroup) ([]*ent.PermissionsPolicy, error) {
	return obj.QueryPolicies().All(ctx)
}

func (r mutationResolver) AddUsersGroup(ctx context.Context, input models.AddUsersGroupInput) (*ent.UsersGroup, error) {
	client := r.ClientFrom(ctx)
	g, err := client.UsersGroup.Create().
		SetName(input.Name).
		SetNillableDescription(input.Description).
		SetStatus(usersgroup.StatusACTIVE).
		Save(ctx)
	if ent.IsConstraintError(err) {
		return nil, gqlerror.Errorf("A group with the name %s already exists", input.Name)
	}
	return client.UsersGroup.UpdateOneID(g.ID).
		AddMemberIDs(input.Members...).
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
	if input.Members != nil {
		currentMembers, err := client.User.Query().
			Where(user.HasGroupsWith(usersgroup.ID(input.ID))).
			IDs(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying members of group %q", input.ID)
		}
		addedMembers, removedMembers := resolverutil.GetDifferenceBetweenSlices(currentMembers, input.Members)
		m = m.
			AddMemberIDs(addedMembers...).
			RemoveMemberIDs(removedMembers...)
	}
	g, err := m.Save(ctx)
	if ent.IsConstraintError(err) {
		return nil, gqlerror.Errorf("A group with the name %s already exists", *input.Name)
	}
	return g, nil
}

func (r mutationResolver) DeleteUsersGroup(ctx context.Context, id int) (bool, error) {
	client := r.ClientFrom(ctx)
	if err := client.UsersGroup.DeleteOneID(id).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return false, gqlerror.Errorf("users group doesn't exist")
		}
		return false, fmt.Errorf("deleting users group: %w", err)
	}
	return true, nil
}

func (r mutationResolver) UpdatePermissionsPoliciesInUsersGroup(
	ctx context.Context,
	input models.UpdatePermissionsPoliciesInUsersGroupInput,
) (*ent.UsersGroup, error) {

	return r.ClientFrom(ctx).UsersGroup.UpdateOneID(input.ID).
		AddPolicyIDs(input.AddPermissionsPolicyIds...).
		RemovePolicyIDs(input.RemovePermissionsPolicyIds...).
		Save(ctx)
}
