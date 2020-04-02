// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer"
)

type viewerResolver struct{}

func (viewerResolver) User(ctx context.Context, obj *viewer.Viewer) (*ent.User, error) {
	return viewer.UserFromContext(ctx)
}

func (viewerResolver) Permissions(ctx context.Context, obj *viewer.Viewer) (*models.PermissionSettings, error) {
	u, err := viewer.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	adminPolicy := models.AdministrativePolicy{
		CanRead: u.Role == user.RoleADMIN || u.Role == user.RoleOWNER,
	}
	res := models.PermissionSettings{
		CanWrite:    obj.Role != "readonly",
		AdminPolicy: &adminPolicy,
	}
	return &res, nil
}
