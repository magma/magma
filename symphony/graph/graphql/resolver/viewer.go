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

func (viewerResolver) Permissions(ctx context.Context, _ *viewer.Viewer) (*models.PermissionSettings, error) {
	return authz.Permissions(ctx)
}
