// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"
	"time"

	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/stretchr/testify/require"

	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/viewer"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	models2 "github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
)

func getCommentCudOperations(
	ctx context.Context,
	c *ent.Client,
	setParent func(*ent.CommentCreate) *ent.CommentCreate,
) cudOperations {
	author := viewer.MustGetOrCreateUser(ctx, "AuthID", user.RoleOWNER)
	commentQuery := c.Comment.Create().
		SetAuthor(author).
		SetText("comment")
	commentQuery = setParent(commentQuery)
	commentEntity := commentQuery.SaveX(ctx)
	createComment := func(ctx context.Context) error {
		commentQuery := c.Comment.Create().
			SetText("comment").
			SetAuthor(author)
		commentQuery = setParent(commentQuery)
		_, err := commentQuery.Save(ctx)
		return err
	}
	updateComment := func(ctx context.Context) error {
		return c.Comment.UpdateOne(commentEntity).
			SetText("newComment").
			Exec(ctx)
	}
	deleteComment := func(ctx context.Context) error {
		return c.Comment.DeleteOne(commentEntity).
			Exec(ctx)
	}
	return cudOperations{
		create: createComment,
		update: updateComment,
		delete: deleteComment,
	}
}

func TestCommentOfWorkOrderReadPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	u := viewer.MustGetOrCreateUser(ctx, "AuthID", user.RoleUSER)
	woType1, wo1 := prepareWorkOrderData(ctx, c)
	_, wo2 := prepareWorkOrderData(ctx, c)
	c.Comment.Create().
		SetText("Hi!").
		SetAuthor(u).
		SetWorkOrder(wo1).
		SaveX(ctx)
	c.Comment.Create().
		SetText("Hi!").
		SetAuthor(u).
		SetWorkOrder(wo2).
		SaveX(ctx)
	t.Run("EmptyPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.Comment.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Zero(t, count)
	})
	t.Run("PartialPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissions.WorkforcePolicy.Read.IsAllowed = models2.PermissionValueByCondition
		permissions.WorkforcePolicy.Read.WorkOrderTypeIds = []int{woType1.ID}
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.Comment.Query().Count(permissionsContext)
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
		count, err := c.Comment.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})
}
func TestCommentOfProjectReadPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	u := viewer.MustGetOrCreateUser(ctx, "AuthID", user.RoleUSER)
	projectType1, project1 := prepareProjectData(ctx, c)
	_, project2 := prepareProjectData(ctx, c)
	c.Comment.Create().
		SetText("Hi!").
		SetAuthor(u).
		SetProject(project1).
		SaveX(ctx)
	c.Comment.Create().
		SetText("Hi!").
		SetAuthor(u).
		SetProject(project2).
		SaveX(ctx)
	t.Run("EmptyPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.Comment.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Zero(t, count)
	})
	t.Run("PartialPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissions.WorkforcePolicy.Read.IsAllowed = models2.PermissionValueByCondition
		permissions.WorkforcePolicy.Read.ProjectTypeIds = []int{projectType1.ID}
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.Comment.Query().Count(permissionsContext)
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
		count, err := c.Comment.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})
}

func TestProjectCommentPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	projectType := c.ProjectType.Create().
		SetName("ProjectType").
		SaveX(ctx)

	project := c.Project.Create().
		SetName("Project").
		SetType(projectType).
		SaveX(ctx)

	cudOperations := getCommentCudOperations(ctx, c, func(ptc *ent.CommentCreate) *ent.CommentCreate {
		return ptc.SetProject(project)
	})
	runCudPolicyTest(t, cudPolicyTest{
		initialPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Read.IsAllowed = models2.PermissionValueYes
		},
		appendPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Data.Update.IsAllowed = models2.PermissionValueYes
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
	})
}

func TestWorkOrderCommentPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	u := viewer.MustGetOrCreateUser(ctx, "AuthID", user.RoleOWNER)
	workOrderType := c.WorkOrderType.Create().
		SetName("WorkOrderType").
		SaveX(ctx)
	workOrder := c.WorkOrder.Create().
		SetName("workOrder").
		SetType(workOrderType).
		SetOwner(u).
		SetCreationDate(time.Now()).
		SaveX(ctx)

	cudOperations := getCommentCudOperations(ctx, c, func(ptc *ent.CommentCreate) *ent.CommentCreate {
		return ptc.SetWorkOrder(workOrder)
	})
	runCudPolicyTest(t, cudPolicyTest{
		initialPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Read.IsAllowed = models2.PermissionValueYes
		},
		appendPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Data.Update.IsAllowed = models2.PermissionValueYes
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
	})
}
