package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent/usersgroup"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/viewer"
)

// MutationWithViewerRuleFunc type is an adapter to allow the use of
// ordinary functions with viewer parameter as mutation rules.
type MutationWithViewerRuleFunc func(context.Context, ent.Mutation, *viewer.Viewer) error

// WritePermissionGroupName is the name of the group that its member has write permission for all symphony.
const WritePermissionGroupName = "Write Permission"

// MutationWithViewerRule returns a rule that checks for viewer and skip if not exist
func MutationWithViewerRule(rule MutationWithViewerRuleFunc) privacy.MutationRule {
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

// AllowViewerWritePermissionsRule is a rule to grant write permission.
func AllowViewerWritePermissionsRule(ctx context.Context, _ ent.Mutation, v *viewer.Viewer) error {
	writePermissions, err := UserHasWritePermissions(ctx, v)
	if err != nil {
		return privacy.Denyf("cannot read write permissions of user .%w", err)
	}
	if writePermissions {
		return privacy.Allow
	}
	return privacy.Skip
}

// DenyRule is a rule that always deny
func DenyRule(_ context.Context, _ ent.Mutation, _ *viewer.Viewer) error {
	return privacy.Deny
}

// AllowAdminRule is a rule that allows permissions if user has at least admin permissions
func AllowAdminRule(_ context.Context, _ ent.Mutation, v *viewer.Viewer) error {
	u := v.User()
	if userHasAdminPermissions(u) {
		return privacy.Allow
	}
	return privacy.Skip
}
