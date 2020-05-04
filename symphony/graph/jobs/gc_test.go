package jobs

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	"github.com/stretchr/testify/require"
)

func TestGarbageCollectProperties(t *testing.T) {
	r := newJobsTestResolver(t)
	defer r.drv.Close()
	client := r.client
	ctx := ent.NewContext(context.Background(), client)
	locationType := client.LocationType.Create().
		SetName("LocationType").
		SaveX(ctx)
	propTypeToDelete := client.PropertyType.Create().
		SetName("PropToDelete").
		SetLocationType(locationType).
		SetType(models.PropertyKindString.String()).
		SetDeleted(true).
		SaveX(ctx)
	propTypeToDelete2 := client.PropertyType.Create().
		SetName("PropToDelete2").
		SetLocationType(locationType).
		SetType(models.PropertyKindBool.String()).
		SetDeleted(true).
		SaveX(ctx)
	propType := client.PropertyType.Create().
		SetName("Prop").
		SetLocationType(locationType).
		SetType(models.PropertyKindInt.String()).
		SaveX(ctx)
	_ = client.Location.Create().
		SetName("Location").
		SetType(locationType).
		SaveX(ctx)
	propToDelete1 := client.Property.Create().
		SetType(propTypeToDelete).
		SetStringVal("Prop1").
		SaveX(ctx)
	propToDelete2 := client.Property.Create().
		SetType(propTypeToDelete).
		SetStringVal("Prop2").
		SaveX(ctx)
	propToDelete3 := client.Property.Create().
		SetType(propTypeToDelete2).
		SetBoolVal(true).
		SaveX(ctx)
	prop := client.Property.Create().
		SetType(propType).
		SetIntVal(28).
		SaveX(ctx)
	err := r.jobsRunner.collectProperties(ctx)
	require.NoError(t, err)
	require.False(t, client.PropertyType.Query().Where(propertytype.ID(propTypeToDelete.ID)).ExistX(ctx))
	require.False(t, client.Property.Query().Where(property.ID(propToDelete1.ID)).ExistX(ctx))
	require.False(t, client.Property.Query().Where(property.ID(propToDelete2.ID)).ExistX(ctx))
	require.False(t, client.Property.Query().Where(property.ID(propToDelete3.ID)).ExistX(ctx))
	require.True(t, client.PropertyType.Query().Where(propertytype.ID(propType.ID)).ExistX(ctx))
	require.True(t, client.Property.Query().Where(property.ID(prop.ID)).ExistX(ctx))
}
