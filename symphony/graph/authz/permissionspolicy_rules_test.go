// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
)

func TestPermissionsPolicyWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	policy := c.PermissionsPolicy.Create().
		SetName("Policy").
		SaveX(ctx)
	createPolicy := func(ctx context.Context) error {
		_, err := c.PermissionsPolicy.Create().
			SetName("Policy2").
			Save(ctx)
		return err
	}
	updatePolicy := func(ctx context.Context) error {
		return c.PermissionsPolicy.UpdateOne(policy).
			SetName("NewName").
			Exec(ctx)
	}
	deletePolicy := func(ctx context.Context) error {
		return c.PermissionsPolicy.DeleteOne(policy).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.AdminPolicy.Access.IsAllowed = models2.PermissionValueYes
		},
		create: createPolicy,
		update: updatePolicy,
		delete: deletePolicy,
	})
}
