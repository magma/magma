// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/authz/models"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
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

// AllowViewerWritePermissionsRule grants write permission.
func AllowViewerWritePermissionsRule() privacy.MutationRule {
	return mutationWithViewerRule(func(ctx context.Context, _ ent.Mutation, v *viewer.Viewer) error {
		if FromContext(ctx).CanWrite {
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// AllowAdminRule grants access to admins.
func AllowAdminRule() privacy.MutationRule {
	return mutationWithViewerRule(func(ctx context.Context, _ ent.Mutation, v *viewer.Viewer) error {
		if FromContext(ctx).AdminPolicy.Access.IsAllowed == models.PermissionValueYes {
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
