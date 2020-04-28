// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent/usersgroup"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/viewer"
)

// WritePermissionGroupName is the name of the group that its member has write permission for all symphony.
const WritePermissionGroupName = "Write Permission"

// mutationWithViewerRule returns a rule that checks for viewer and skip if it doesn't exist.
func mutationWithViewerRule(rule func(context.Context, ent.Mutation, *viewer.Viewer) error) privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			return privacy.Skip
		}
		return rule(ctx, m, v)
	})
}

// UserHasWritePermissions checks if user has write permissions based on role and group membership.
func UserHasWritePermissions(ctx context.Context, v *viewer.Viewer) (bool, error) {
	if !v.Features.Enabled(viewer.FeatureReadOnly) {
		return true, nil
	}
	u := v.User()
	if u.Role == user.RoleOWNER {
		return true, nil
	}
	return u.QueryGroups().
		Where(usersgroup.Name(WritePermissionGroupName)).
		Exist(ctx)
}

func userHasAdminPermissions(u *ent.User) bool {
	return u.Role == user.RoleADMIN || u.Role == user.RoleOWNER
}

// AllowViewerWritePermissionsRule grants write permission.
func AllowViewerWritePermissionsRule() privacy.MutationRule {
	return mutationWithViewerRule(func(ctx context.Context, _ ent.Mutation, v *viewer.Viewer) error {
		switch hasPerm, err := UserHasWritePermissions(ctx, v); {
		case err != nil:
			return privacy.Denyf("cannot get write permissions of user: %w", err)
		case hasPerm:
			return privacy.Allow
		default:
			return privacy.Skip
		}
	})
}

// AllowAdminRule grants access to admins.
func AllowAdminRule() privacy.MutationRule {
	return mutationWithViewerRule(func(_ context.Context, _ ent.Mutation, v *viewer.Viewer) error {
		if userHasAdminPermissions(v.User()) {
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// AlwaysAllowIfNoViewerRule grants access if no viewer is present on context.
func AlwaysAllowIfNoViewerRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, _ ent.Mutation) error {
		if viewer.FromContext(ctx) == nil {
			return privacy.Allow
		}
		return privacy.Skip
	})
}
