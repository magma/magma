// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"
)

// AllowAdminRule grants access to admins.
func AllowAdminRule() privacy.MutationRule {
	return mutationWithPermissionRule(func(ctx context.Context, _ ent.Mutation, p *models.PermissionSettings) error {
		if p.AdminPolicy.Access.IsAllowed == models2.PermissionValueYes {
			return privacy.Allow
		}
		return privacy.Skip
	})
}
