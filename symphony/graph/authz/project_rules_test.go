// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
)

func TestProjectTypeWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := context.Background()
	projectType := c.ProjectType.Create().
		SetName("ProjectType").
		SaveX(ctx)
	createProjectType := func(ctx context.Context) error {
		_, err := c.ProjectType.Create().
			SetName("NewProjectType").
			Save(ctx)
		return err
	}
	updateProjectType := func(ctx context.Context) error {
		return c.ProjectType.UpdateOne(projectType).
			SetName("NewName").
			Exec(ctx)
	}
	deleteProjectType := func(ctx context.Context) error {
		return c.ProjectType.DeleteOne(projectType).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		getCud: func(p *models.PermissionSettings) *models.Cud {
			return p.WorkforcePolicy.Templates
		},
		create: createProjectType,
		update: updateProjectType,
		delete: deleteProjectType,
	})
}
