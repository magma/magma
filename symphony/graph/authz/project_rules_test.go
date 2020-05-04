// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/viewer"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"
	"github.com/facebookincubator/symphony/graph/ent"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
)

func prepareProjectData(ctx context.Context, c *ent.Client) (*ent.ProjectType, *ent.Project) {
	projectType := c.ProjectType.Create().
		SetName("ProjectType").
		SaveX(ctx)
	project := c.Project.Create().
		SetName("Project").
		SetType(projectType).
		SaveX(ctx)
	return projectType, project
}
func TestProjectWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := context.Background()
	projectType, project := prepareProjectData(ctx, c)
	createProject := func(ctx context.Context) error {
		_, err := c.Project.Create().
			SetName("NewProject").
			SetType(projectType).
			Save(ctx)
		return err
	}
	updateProject := func(ctx context.Context) error {
		return c.Project.UpdateOne(project).
			SetName("NewName").
			Exec(ctx)
	}
	deleteProject := func(ctx context.Context) error {
		return c.Project.DeleteOne(project).
			Exec(ctx)
	}
	getCud := func(p *models.PermissionSettings) *models.WorkforceCud {
		return p.WorkforcePolicy.Data
	}
	runCudPolicyTest(t, cudPolicyTest{
		getCud: func(p *models.PermissionSettings) *models.Cud {
			return &models.Cud{
				Create: getCud(p).Create,
				Update: getCud(p).Update,
				Delete: getCud(p).Delete,
			}
		},
		create: createProject,
		update: updateProject,
		delete: deleteProject,
	})
}

func TestProjectTransferOwnershipWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	projectType, project := prepareProjectData(ctx, c)
	getCud := func(p *models.PermissionSettings) *models.WorkforceCud {
		return p.WorkforcePolicy.Data
	}
	appendTransferOwnership := func(p *models.PermissionSettings) {
		getCud(p).TransferOwnership.IsAllowed = models2.PermissionValueYes
	}
	createProjectWithCreator := func(ctx context.Context) error {
		u := viewer.FromContext(ctx).(*viewer.UserViewer).User()
		_, err := c.Project.Create().
			SetName("NewProject").
			SetType(projectType).
			SetCreator(u).
			Save(ctx)
		return err
	}
	updateProjectCreator := func(ctx context.Context) error {
		u := viewer.FromContext(ctx).(*viewer.UserViewer).User()
		_, err := c.Project.UpdateOne(project).
			SetCreator(u).
			Save(ctx)
		return err
	}
	clearProjectCreator := func(ctx context.Context) error {
		_, err := c.Project.UpdateOne(project).
			ClearCreator().
			Save(ctx)
		return err
	}
	tests := []policyTest{
		{
			operationName: "CreateWithCreator",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Create.IsAllowed = models2.PermissionValueYes
			},
			appendPermissions: appendTransferOwnership,
			operation:         createProjectWithCreator,
		},
		{
			operationName: "UpdateWithCreator",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Update.IsAllowed = models2.PermissionValueYes
			},
			appendPermissions: appendTransferOwnership,
			operation:         updateProjectCreator,
		},
		{
			operationName: "ClearProjectCreator",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Update.IsAllowed = models2.PermissionValueYes
			},
			appendPermissions: appendTransferOwnership,
			operation:         clearProjectCreator,
		},
	}
	runPolicyTest(t, tests)
}

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
