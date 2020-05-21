// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent/user"

	"github.com/stretchr/testify/require"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent/privacy"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
)

type cudOperations struct {
	create func(ctx context.Context) error
	update func(ctx context.Context) error
	delete func(ctx context.Context) error
}

type policyTest struct {
	operationName      string
	initialPermissions func(p *models.PermissionSettings)
	appendPermissions  func(p *models.PermissionSettings)
	operation          func(ctx context.Context) error
}

type contextBasedPolicyTest struct {
	operationName          string
	noPermissionsContext   context.Context
	withPermissionsContext context.Context
	operation              func(ctx context.Context) error
}

type cudPolicyTest struct {
	getCud             func(p *models.PermissionSettings) *models.Cud
	initialPermissions func(p *models.PermissionSettings)
	appendPermissions  func(p *models.PermissionSettings)
	create             func(ctx context.Context) error
	update             func(ctx context.Context) error
	delete             func(ctx context.Context) error
}

func runPolicyTest(t *testing.T, tests []policyTest) {
	var contextBasedTests []contextBasedPolicyTest
	for _, test := range tests {
		c := viewertest.NewTestClient(t)
		noPermissions := authz.EmptyPermissions()
		if test.initialPermissions != nil {
			test.initialPermissions(noPermissions)
		}
		noPermissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(noPermissions))
		withPermissions := authz.EmptyPermissions()
		if test.initialPermissions != nil {
			test.initialPermissions(withPermissions)
		}
		test.appendPermissions(withPermissions)
		withPermissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(withPermissions))

		contextBasedTests = append(contextBasedTests, contextBasedPolicyTest{
			operationName:          test.operationName,
			noPermissionsContext:   noPermissionsContext,
			withPermissionsContext: withPermissionsContext,
			operation:              test.operation,
		})
	}
	runContextBasedPolicyTest(t, contextBasedTests)
}

func runContextBasedPolicyTest(t *testing.T, tests []contextBasedPolicyTest) {
	for _, test := range tests {
		t.Run(test.operationName, func(t *testing.T) {
			modes := map[string]bool{"Denied": false, "Allowed": true}
			for _, name := range []string{"Denied", "Allowed"} {
				t.Run(name, func(t *testing.T) {
					allowed := modes[name]
					if allowed {
						err := test.operation(test.withPermissionsContext)
						require.NoError(t, err)
					} else {
						err := test.operation(test.noPermissionsContext)
						require.True(t, errors.Is(err, privacy.Deny), fmt.Sprintf("Error is %v", err))
					}
				})
			}
		})
	}
}

func runCudPolicyTest(t *testing.T, test cudPolicyTest) {
	tests := []policyTest{
		{
			operationName:      "Create",
			initialPermissions: test.initialPermissions,
			appendPermissions: func(p *models.PermissionSettings) {
				if test.appendPermissions != nil {
					test.appendPermissions(p)
				} else {
					test.getCud(p).Create.IsAllowed = models2.PermissionValueYes
				}
			},
			operation: test.create,
		},
		{
			operationName:      "Update",
			initialPermissions: test.initialPermissions,
			appendPermissions: func(p *models.PermissionSettings) {
				if test.appendPermissions != nil {
					test.appendPermissions(p)
				} else {
					test.getCud(p).Update.IsAllowed = models2.PermissionValueYes
				}
			},
			operation: test.update,
		},
		{
			operationName:      "Delete",
			initialPermissions: test.initialPermissions,
			appendPermissions: func(p *models.PermissionSettings) {
				if test.appendPermissions != nil {
					test.appendPermissions(p)
				} else {
					test.getCud(p).Delete.IsAllowed = models2.PermissionValueYes
				}
			},
			operation: test.delete,
		},
	}
	runPolicyTest(t, tests)
}
