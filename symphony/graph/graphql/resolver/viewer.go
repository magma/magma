// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/authz"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer"
)

type viewerResolver struct{}

func (viewerResolver) Email(_ context.Context, obj *viewer.Viewer) (string, error) {
	return obj.User().Email, nil
}

func (viewerResolver) User(_ context.Context, obj *viewer.Viewer) (*ent.User, error) {
	return obj.User(), nil
}

func (viewerResolver) Permissions(ctx context.Context, obj *viewer.Viewer) (*models.PermissionSettings, error) {
	writePermissions, err := authz.UserHasWritePermissions(ctx, obj)
	if err != nil {
		return nil, err
	}
	policiesEnabled := viewer.FromContext(ctx).Features.Enabled(viewer.FeatureUserManagementDev)
	inventoryPolicy := authz.NewInventoryPolicy(true, writePermissions)
	workforcePolicy := authz.NewWorkforcePolicy(true, writePermissions)
	if policiesEnabled {
		inventoryPolicy, workforcePolicy, err = authz.PermissionPolicies(ctx)
		if err != nil {
			return nil, err
		}
	}
	res := models.PermissionSettings{
		// TODO(T64743627): Deprecate CanWrite field
		CanWrite:            writePermissions,
		AdminPolicy:         authz.NewAdministrativePolicy(obj.User()),
		InventoryPolicy:     inventoryPolicy,
		WorkforcePermission: workforcePolicy,
	}
	return &res, nil
}
