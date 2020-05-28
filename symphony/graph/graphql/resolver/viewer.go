// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/authz"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/viewer"
)

type viewerResolver struct{}

func (viewerResolver) User(_ context.Context, v viewer.Viewer) (*ent.User, error) {
	if v, ok := v.(*viewer.UserViewer); ok {
		return v.User(), nil
	}
	return nil, nil
}

func (viewerResolver) Permissions(ctx context.Context, _ viewer.Viewer) (*models.PermissionSettings, error) {
	return authz.Permissions(ctx)
}
