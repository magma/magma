// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphgrpc

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/viewer"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CreateServiceContext(ctx context.Context, tenantName, serviceName string, role user.Role) (context.Context, error) {
	v := viewer.NewAutomation(tenantName, serviceName, role)
	ctx = viewer.NewContext(ctx, v)
	permissions, err := authz.Permissions(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "error getting permissions")
	}
	return authz.NewContext(ctx, permissions), nil
}
