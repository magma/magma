// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/graph/resolverutil"
	"github.com/pkg/errors"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type permissionsPolicyResolver struct{}

func (r queryResolver) PermissionsPolicies(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int) (*ent.PermissionsPolicyConnection, error) {
	return r.ClientFrom(ctx).PermissionsPolicy.Query().
		Paginate(ctx, after, first, before, last)
}

func (r permissionsPolicyResolver) Policy(ctx context.Context, obj *ent.PermissionsPolicy) (models.SystemPolicy, error) {
	if obj.InventoryPolicy != nil {
		return authz.AppendInventoryPolicies(
			authz.NewInventoryPolicy(false),
			obj.InventoryPolicy), nil
	}
	return authz.AppendWorkforcePolicies(
		authz.NewWorkforcePolicy(false, false),
		obj.WorkforcePolicy), nil
}

func (permissionsPolicyResolver) Groups(ctx context.Context, obj *ent.PermissionsPolicy) ([]*ent.UsersGroup, error) {
	return obj.QueryGroups().All(ctx)
}

func (mutationResolver) AddPermissionsPolicy(
	ctx context.Context,
	input models.AddPermissionsPolicyInput,
) (*ent.PermissionsPolicy, error) {
	client := ent.FromContext(ctx)
	mutation := client.PermissionsPolicy.Create().
		SetName(input.Name).
		SetNillableDescription(input.Description).
		SetNillableIsGlobal(input.IsGlobal)
	if input.Groups != nil {
		mutation = mutation.AddGroupIDs(input.Groups...)
	}
	if input.InventoryInput != nil && input.WorkforceInput != nil {
		return nil, fmt.Errorf("policy cannot be of both inventory and workforce types")
	}
	switch {
	case input.InventoryInput != nil:
		mutation.SetInventoryPolicy(input.InventoryInput)
	case input.WorkforceInput != nil:
		mutation.SetWorkforcePolicy(input.WorkforceInput)
	default:
		return nil, fmt.Errorf("no policy found in input")
	}
	policy, err := mutation.Save(ctx)
	if ent.IsConstraintError(err) {
		return nil, gqlerror.Errorf("A policy with the given name already exists: %s", input.Name)
	}
	return policy, err
}

func (mutationResolver) EditPermissionsPolicy(
	ctx context.Context,
	input models.EditPermissionsPolicyInput,
) (*ent.PermissionsPolicy, error) {
	client := ent.FromContext(ctx)
	p, err := client.PermissionsPolicy.Get(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("querying permissionsPolicy %q: %w", input.ID, err)
	}
	upd := client.PermissionsPolicy.
		UpdateOne(p).
		SetNillableDescription(input.Description).
		SetNillableIsGlobal(input.IsGlobal)
	if input.Name != nil {
		upd = upd.SetName(*input.Name)
	}
	if input.Groups != nil {
		currentGroups, err := client.PermissionsPolicy.QueryGroups(p).IDs(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying groups of permissionPolicy %q", input.ID)
		}
		addGroupIds, removeGroupIds := resolverutil.GetDifferenceBetweenSlices(currentGroups, input.Groups)
		upd = upd.
			AddGroupIDs(addGroupIds...).
			RemoveGroupIDs(removeGroupIds...)
	}
	switch {
	case input.InventoryInput != nil && input.WorkforceInput != nil:
		return nil, fmt.Errorf("policy cannot be of both inventory and workforce types")
	case input.InventoryInput != nil && p.WorkforcePolicy != nil:
		return nil, fmt.Errorf("only workforce policy is legal to edit")
	case input.WorkforceInput != nil && p.InventoryPolicy != nil:
		return nil, fmt.Errorf("only inventory policy is legal to edit")
	case input.InventoryInput != nil && p.InventoryPolicy != nil:
		upd.SetInventoryPolicy(input.InventoryInput)
	case input.WorkforceInput != nil && p.WorkforcePolicy != nil:
		upd.SetWorkforcePolicy(input.WorkforceInput)
	}
	p, err = upd.Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, gqlerror.Errorf("A policy with the name %v already exists", *input.Name)
		}
		return nil, fmt.Errorf("updating permissionsPolicy %q: %w", input.ID, err)
	}
	return p, nil
}

func (r mutationResolver) DeletePermissionsPolicy(ctx context.Context, id int) (bool, error) {
	client := r.ClientFrom(ctx)
	if err := client.PermissionsPolicy.DeleteOneID(id).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return false, gqlerror.Errorf("permissionsPolicy doesn't exist")
		}
		return false, fmt.Errorf("deleting permissionsPolicy: %w", err)
	}
	return true, nil
}
