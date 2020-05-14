// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent/user"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"

	"github.com/facebookincubator/symphony/graph/viewer"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"
	"github.com/facebookincubator/symphony/graph/ent"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
)

func prepareProjectData(ctx context.Context, c *ent.Client) (*ent.ProjectType, *ent.Project) {
	projectTypeName := uuid.New().String()
	projectName := uuid.New().String()
	projectType := c.ProjectType.Create().
		SetName(projectTypeName).
		SaveX(ctx)
	project := c.Project.Create().
		SetName(projectName).
		SetType(projectType).
		SaveX(ctx)
	return projectType, project
}
func TestProjectWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	projectType, project := prepareProjectData(ctx, c)
	projectType2, project2 := prepareProjectData(ctx, c)
	createProject := func(ctx context.Context) error {
		_, err := c.Project.Create().
			SetName("NewProject").
			SetType(projectType).
			Save(ctx)
		return err
	}
	createProject2 := func(ctx context.Context) error {
		_, err := c.Project.Create().
			SetName("NewProject2").
			SetType(projectType2).
			Save(ctx)
		return err
	}
	updateProject := func(ctx context.Context) error {
		return c.Project.UpdateOne(project).
			SetName("NewName").
			Exec(ctx)
	}
	updateProject2 := func(ctx context.Context) error {
		return c.Project.UpdateOne(project2).
			SetName("NewName2").
			Exec(ctx)
	}
	deleteProject := func(ctx context.Context) error {
		return c.Project.DeleteOne(project).
			Exec(ctx)
	}
	deleteProject2 := func(ctx context.Context) error {
		return c.Project.DeleteOne(project2).
			Exec(ctx)
	}
	tests := []policyTest{
		{
			operationName: "Create",
			appendPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Create.IsAllowed = models2.PermissionValueYes
			},
			operation: createProject,
		},
		{
			operationName: "CreateWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Create.IsAllowed = models2.PermissionValueByCondition
				p.WorkforcePolicy.Data.Create.ProjectTypeIds = []int{projectType.ID}
			},
			appendPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Create.ProjectTypeIds = []int{projectType.ID, projectType2.ID}
			},
			operation: createProject2,
		},
		{
			operationName: "Update",
			appendPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Update.IsAllowed = models2.PermissionValueYes
			},
			operation: updateProject,
		},
		{
			operationName: "UpdateWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Update.IsAllowed = models2.PermissionValueByCondition
				p.WorkforcePolicy.Data.Update.ProjectTypeIds = []int{projectType.ID}
			},
			appendPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Update.ProjectTypeIds = []int{projectType.ID, projectType2.ID}
			},
			operation: updateProject2,
		},
		{
			operationName: "Delete",
			appendPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Delete.IsAllowed = models2.PermissionValueYes
			},
			operation: deleteProject,
		},
		{
			operationName: "DeleteWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Delete.IsAllowed = models2.PermissionValueByCondition
				p.WorkforcePolicy.Data.Delete.ProjectTypeIds = []int{projectType.ID}
			},
			appendPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Delete.ProjectTypeIds = []int{projectType.ID, projectType2.ID}
			},
			operation: deleteProject2,
		},
	}
	runPolicyTest(t, tests)
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
	ctx := viewertest.NewContext(context.Background(), c)
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

func TestProjectReadPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	projectType, _ := prepareProjectData(ctx, c)
	_, _ = prepareProjectData(ctx, c)
	t.Run("EmptyPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.Project.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Zero(t, count)
	})
	t.Run("PartialPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissions.WorkforcePolicy.Read.IsAllowed = models2.PermissionValueByCondition
		permissions.WorkforcePolicy.Read.ProjectTypeIds = []int{projectType.ID}
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.Project.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})
	t.Run("FullPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissions.WorkforcePolicy.Read.IsAllowed = models2.PermissionValueYes
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.Project.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})
}
