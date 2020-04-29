// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"errors"
	"testing"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

type policyTest struct {
	operationName     string
	appendPermissions func(p *models.PermissionSettings)
	operation         func(ctx context.Context) error
}

type cudPolicyTest struct {
	getCud func(p *models.PermissionSettings) *models.Cud
	create func(ctx context.Context) error
	update func(ctx context.Context) error
	delete func(ctx context.Context) error
}

func runPolicyTest(t *testing.T, tests []policyTest) {
	for _, test := range tests {
		t.Run(test.operationName, func(t *testing.T) {
			for name, allowed := range map[string]bool{"Denied": true, "Allowed": false} {
				t.Run(name, func(t *testing.T) {
					c := viewertest.NewTestClient(t)
					permissions := authz.EmptyPermissions()
					if allowed {
						test.appendPermissions(permissions)
					}
					ctx := viewertest.NewContext(context.Background(), c, viewertest.WithPermissions(permissions))
					err := test.operation(ctx)
					if allowed {
						require.NoError(t, err)
					} else {
						require.True(t, errors.Is(err, privacy.Deny))
					}
				})
			}
		})
	}
}

func runCudPolicyTest(t *testing.T, test cudPolicyTest) {
	tests := []policyTest{
		{
			operationName: "Create",
			appendPermissions: func(p *models.PermissionSettings) {
				test.getCud(p).Create.IsAllowed = models2.PermissionValueYes
			},
			operation: test.create,
		},
		{
			operationName: "Update",
			appendPermissions: func(p *models.PermissionSettings) {
				test.getCud(p).Update.IsAllowed = models2.PermissionValueYes
			},
			operation: test.update,
		},
		{
			operationName: "Delete",
			appendPermissions: func(p *models.PermissionSettings) {
				test.getCud(p).Delete.IsAllowed = models2.PermissionValueYes
			},
			operation: test.delete,
		},
	}
	runPolicyTest(t, tests)
}
