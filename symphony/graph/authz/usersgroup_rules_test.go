package authz_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/viewer"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
)

func TestUsersGroupWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	u := viewer.MustGetOrCreateUser(ctx, "AuthID", user.RoleUSER)
	g := c.UsersGroup.Create().
		SetName("Group").
		AddMembers(u).
		SaveX(ctx)
	createGroup := func(ctx context.Context) error {
		_, err := c.UsersGroup.Create().
			SetName("Group2").
			AddMembers(u).
			Save(ctx)
		return err
	}
	updateGroup := func(ctx context.Context) error {
		return c.UsersGroup.UpdateOne(g).
			SetName("NewName").
			Exec(ctx)
	}
	deleteGroup := func(ctx context.Context) error {
		return c.UsersGroup.DeleteOne(g).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.AdminPolicy.Access.IsAllowed = models2.PermissionValueYes
		},
		create: createGroup,
		update: updateGroup,
		delete: deleteGroup,
	})
}
