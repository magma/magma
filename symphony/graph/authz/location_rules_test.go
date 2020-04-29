package authz_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
)

func TestLocationWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := context.Background()
	locationType := c.LocationType.Create().
		SetName("LocationType").
		SaveX(ctx)
	location := c.Location.Create().
		SetName("Location").
		SetType(locationType).
		SaveX(ctx)
	createLocation := func(ctx context.Context) error {
		_, err := c.Location.Create().
			SetName("NewLocation").
			SetType(locationType).
			Save(ctx)
		return err
	}
	updateLocation := func(ctx context.Context) error {
		return c.Location.UpdateOne(location).
			SetName("NewName").
			Exec(ctx)
	}
	deleteLocation := func(ctx context.Context) error {
		return c.Location.DeleteOne(location).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		getCud: func(p *models.PermissionSettings) *models.Cud {
			return p.InventoryPolicy.Location
		},
		create: createLocation,
		update: updateLocation,
		delete: deleteLocation,
	})
}

func TestLocationTypeWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := context.Background()
	locationType := c.LocationType.Create().
		SetName("LocationType").
		SaveX(ctx)
	createLocationType := func(ctx context.Context) error {
		_, err := c.LocationType.Create().
			SetName("NewLocationType").
			Save(ctx)
		return err
	}
	updateLocationType := func(ctx context.Context) error {
		return c.LocationType.UpdateOne(locationType).
			SetName("NewName").
			Exec(ctx)
	}
	deleteLocationType := func(ctx context.Context) error {
		return c.LocationType.DeleteOne(locationType).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		getCud: func(p *models.PermissionSettings) *models.Cud {
			return p.InventoryPolicy.LocationType
		},
		create: createLocationType,
		update: updateLocationType,
		delete: deleteLocationType,
	})
}
