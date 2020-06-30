// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
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
				p.WorkforcePolicy.Data.Create.IsAllowed = models.PermissionValueYes
			},
			operation: createProject,
		},
		{
			operationName: "CreateWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Create.IsAllowed = models.PermissionValueByCondition
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
				p.WorkforcePolicy.Data.Update.IsAllowed = models.PermissionValueYes
			},
			operation: updateProject,
		},
		{
			operationName: "UpdateWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Update.IsAllowed = models.PermissionValueByCondition
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
				p.WorkforcePolicy.Data.Delete.IsAllowed = models.PermissionValueYes
			},
			operation: deleteProject,
		},
		{
			operationName: "DeleteWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Delete.IsAllowed = models.PermissionValueByCondition
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
	u := viewer.MustGetOrCreateUser(ctx, "SomeUser", user.RoleUSER)
	appendTransferOwnership := func(p *models.PermissionSettings) {
		getCud(p).TransferOwnership.IsAllowed = models.PermissionValueYes
	}
	createProjectWithCreator := func(ctx context.Context) error {
		_, err := c.Project.Create().
			SetName("NewProject").
			SetType(projectType).
			SetCreator(u).
			Save(ctx)
		return err
	}
	updateProjectCreator := func(ctx context.Context) error {
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
			appendPermissions: func(p *models.PermissionSettings) {
				getCud(p).Create.IsAllowed = models.PermissionValueYes
			},
			operation: createProjectWithCreator,
		},
		{
			operationName: "UpdateWithCreator",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Update.IsAllowed = models.PermissionValueYes
			},
			appendPermissions: appendTransferOwnership,
			operation:         updateProjectCreator,
		},
		{
			operationName: "ClearProjectCreator",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Update.IsAllowed = models.PermissionValueYes
			},
			appendPermissions: appendTransferOwnership,
			operation:         clearProjectCreator,
		},
		{
			operationName: "UpdateWithCreatorWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Update.IsAllowed = models.PermissionValueByCondition
				getCud(p).Update.ProjectTypeIds = []int{projectType.ID}
			},
			appendPermissions: func(p *models.PermissionSettings) {
				getCud(p).TransferOwnership.IsAllowed = models.PermissionValueByCondition
				getCud(p).TransferOwnership.ProjectTypeIds = []int{projectType.ID}
			},
			operation: updateProjectCreator,
		},
		{
			operationName: "ClearWorkOrderAssigneeWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Update.IsAllowed = models.PermissionValueByCondition
				getCud(p).Update.ProjectTypeIds = []int{projectType.ID}
			},
			appendPermissions: func(p *models.PermissionSettings) {
				getCud(p).TransferOwnership.IsAllowed = models.PermissionValueByCondition
				getCud(p).TransferOwnership.ProjectTypeIds = []int{projectType.ID}
			},
			operation: clearProjectCreator,
		},
	}
	runPolicyTest(t, tests)
}

func TestProjectCreatorCanEditProjectButNoDelete(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	projectType, _ := prepareProjectData(ctx, c)
	u2 := viewer.MustGetOrCreateUser(ctx, "AnotherUser", user.RoleUSER)
	ctx = userContextWithPermissions(ctx, "SomeUser", func(p *models.PermissionSettings) {
		p.WorkforcePolicy.Data.Create.IsAllowed = models.PermissionValueYes
	})
	u := viewer.FromContext(ctx).(*viewer.UserViewer).User()
	proj, err := c.Project.Create().
		SetName("NewProject").
		SetType(projectType).
		SetCreator(u).
		Save(ctx)
	require.NoError(t, err)
	err = c.Project.UpdateOne(proj).
		SetName("NewName").
		Exec(ctx)
	require.NoError(t, err)
	err = c.Project.UpdateOne(proj).
		SetCreator(u2).
		Exec(ctx)
	require.NoError(t, err)
	err = c.Project.DeleteOne(proj).
		Exec(ctx)
	require.Error(t, err)
	proj2, err := c.Project.Create().
		SetName("NewProject2").
		SetType(projectType).
		SetCreator(u).
		Save(ctx)
	require.NoError(t, err)
	err = c.Project.DeleteOne(proj2).
		Exec(ctx)
	require.Error(t, err)
}

func TestProjectCreatorUnchangedWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	_, project := prepareProjectData(ctx, c)
	_, project2 := prepareProjectData(ctx, c)
	u := viewer.MustGetOrCreateUser(ctx, "SomeUser", user.RoleUSER)
	c.Project.UpdateOne(project).
		SetCreatorID(u.ID).
		ExecX(ctx)
	permissions := authz.EmptyPermissions()
	permissions.WorkforcePolicy.Data.Update.IsAllowed = models.PermissionValueYes
	ctx = viewertest.NewContext(
		context.Background(),
		c,
		viewertest.WithUser("user"),
		viewertest.WithRole(user.RoleUSER),
		viewertest.WithPermissions(permissions))
	err := c.Project.UpdateOne(project).
		SetCreatorID(u.ID).
		Exec(ctx)
	require.NoError(t, err)
	err = c.Project.UpdateOne(project2).
		ClearCreator().
		Exec(ctx)
	require.NoError(t, err)
	err = c.Project.UpdateOne(project).
		ClearCreator().
		Exec(ctx)
	require.Error(t, err)
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
		permissions.WorkforcePolicy.Read.IsAllowed = models.PermissionValueByCondition
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
		permissions.WorkforcePolicy.Read.IsAllowed = models.PermissionValueYes
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

func TestProjectOfWorkOrderReadPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	workOrderType, workOrder := prepareWorkOrderData(ctx, c)
	_, project := prepareProjectData(ctx, c)
	c.Project.UpdateOne(project).
		AddWorkOrders(workOrder).
		ExecX(ctx)
	permissions := authz.EmptyPermissions()
	permissions.WorkforcePolicy.Read.IsAllowed = models.PermissionValueByCondition
	permissions.WorkforcePolicy.Read.WorkOrderTypeIds = []int{workOrderType.ID}
	permissionsContext := viewertest.NewContext(
		context.Background(),
		c,
		viewertest.WithUser("theOwner"),
		viewertest.WithPermissions(permissions))
	count, err := c.Project.Query().Count(permissionsContext)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestProjectOfOwnedWorkOrderReadPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	_, workOrder := prepareWorkOrderData(ctx, c)
	_, project := prepareProjectData(ctx, c)
	c.Project.UpdateOne(project).
		AddWorkOrders(workOrder).
		ExecX(ctx)
	u := viewer.MustGetOrCreateUser(ctx, "theOwner", user.RoleUSER)
	c.WorkOrder.UpdateOne(workOrder).
		SetOwner(u).
		ExecX(ctx)
	permissions := authz.EmptyPermissions()
	permissionsContext := viewertest.NewContext(
		context.Background(),
		c,
		viewertest.WithUser("theOwner"),
		viewertest.WithPermissions(permissions))
	count, err := c.Project.Query().Count(permissionsContext)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestWorkOrderDefinitionWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)

	workOrderType := c.WorkOrderType.Create().
		SetName("WorkOrderType").
		SaveX(ctx)

	projectType := c.ProjectType.Create().
		SetName("ProjectType").
		SaveX(ctx)

	workOrderDef := c.WorkOrderDefinition.Create().
		SetProjectType(projectType).
		SetType(workOrderType).
		SaveX(ctx)

	createItem := func(ctx context.Context) error {
		_, err := c.WorkOrderDefinition.Create().
			SetProjectType(projectType).
			SetType(workOrderType).
			Save(ctx)
		return err
	}
	updateItem := func(ctx context.Context) error {
		return c.WorkOrderDefinition.UpdateOne(workOrderDef).
			SetIndex(1).
			Exec(ctx)
	}
	deleteItem := func(ctx context.Context) error {
		return c.WorkOrderDefinition.DeleteOne(workOrderDef).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Templates.Update.IsAllowed = models.PermissionValueYes
		},
		create: createItem,
		update: updateItem,
		delete: deleteItem,
	})
}
