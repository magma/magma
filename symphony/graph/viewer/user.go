// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewer

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/usersgroup"
)

// WritePermissionGroupName is the name of the group that its member has write permission for all symphony.
const WritePermissionGroupName = "Write Permission"

func IsUserReadOnly(ctx context.Context, u *ent.User) (bool, error) {
	v := FromContext(ctx)
	if !v.Features.Enabled(FeatureReadOnly) {
		return false, nil
	}
	if u.Role == user.RoleOWNER {
		return false, nil
	}
	exist, err := u.QueryGroups().Where(usersgroup.Name(WritePermissionGroupName)).Exist(ctx)
	return !exist, err
}
