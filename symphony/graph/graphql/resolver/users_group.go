// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/ent/usersgroup"
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
	mutation := client.UsersGroup.Create().
		SetName(input.Name).
		SetNillableDescription(input.Description).
		SetStatus(usersgroup.StatusACTIVE)
	if input.Policies != nil {
		mutation = mutation.AddPolicyIDs(input.Policies...)
	}
	if input.Members != nil {
		mutation = mutation.AddMemberIDs(input.Members...)
	}
	usersGroup, err := mutation.Save(ctx)
	if ent.IsConstraintError(err) {
		return nil, gqlerror.Errorf("A group with the given name already exists: %s", input.Name)
	}
	return usersGroup, err
}

func (r mutationResolver) EditUsersGroup(ctx context.Context, input models.EditUsersGroupInput) (*ent.UsersGroup, error) {
	client := r.ClientFrom(ctx)
	ug, err := client.UsersGroup.Get(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("querying usersGroup %q: %w", input.ID, err)
	}
	mutation := client.UsersGroup.UpdateOneID(input.ID).
		SetNillableDescription(input.Description).
		SetNillableStatus(input.Status)
	if input.Name != nil {
		mutation = mutation.SetName(*input.Name)
	}
	if input.Members != nil {
		currentMembers, err := client.User.Query().
			Where(user.HasGroupsWith(usersgroup.ID(input.ID))).
			IDs(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying members of group %q", input.ID)
		}
		addedMembers, removedMembers := resolverutil.GetDifferenceBetweenSlices(currentMembers, input.Members)
		mutation = mutation.
			AddMemberIDs(addedMembers...).
			RemoveMemberIDs(removedMembers...)
	}
	if input.Policies != nil {
		currentPolicies, err := client.UsersGroup.QueryPolicies(ug).IDs(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying policies of usersGroup %q", input.ID)
		}
		addPermissionsPolicyIds, removePermissionsPolicyIds := resolverutil.GetDifferenceBetweenSlices(currentPolicies, input.Policies)
		mutation = mutation.
			AddPolicyIDs(addPermissionsPolicyIds...).
			RemovePolicyIDs(removePermissionsPolicyIds...)
	}
	g, err := mutation.Save(ctx)
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
