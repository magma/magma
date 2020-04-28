// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

type permissionsPolicyResolver struct{}

func (r queryResolver) PermissionsPolicies(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int) (*ent.PermissionsPolicyConnection, error) {
	return r.ClientFrom(ctx).PermissionsPolicy.Query().
		Paginate(ctx, after, first, before, last)
}

func (r permissionsPolicyResolver) Policy(ctx context.Context, obj *ent.PermissionsPolicy) (models.SystemPolicy, error) {
	if obj.InventoryPolicy != nil {
		return authz.AppendInventoryPolicies(
			authz.NewInventoryPolicy(false, false),
			obj.InventoryPolicy), nil
	}
	return authz.AppendWorkforcePolicies(
		authz.NewWorkforcePolicy(false, false),
		obj.WorkforcePolicy), nil
}

func (permissionsPolicyResolver) Groups(ctx context.Context, obj *ent.PermissionsPolicy) ([]*ent.UsersGroup, error) {
	return obj.QueryGroups().All(ctx)
}

func (mutationResolver) AddPolicy(ctx context.Context, input models.AddPermissionsPolicyInput) (*ent.PermissionsPolicy, error) {
	client := ent.FromContext(ctx)
	mutation := client.PermissionsPolicy.Create().
		SetName(input.Name).
		SetNillableDescription(input.Description).
		SetNillableIsGlobal(input.IsGlobal)
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
		return nil, fmt.Errorf("policy with the given name already exists: %s", input.Name)
	}
	return policy, err
}
