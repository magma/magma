// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jobs

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/pkg/ent/serviceendpoint"

	"github.com/facebookincubator/symphony/pkg/ent/serviceendpointdefinition"

	"github.com/facebookincubator/symphony/pkg/ent/service"
	"github.com/facebookincubator/symphony/pkg/ent/servicetype"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent/property"
	"github.com/facebookincubator/symphony/pkg/ent/propertytype"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

func TestGarbageCollectProperties(t *testing.T) {
	r := newJobsTestResolver(t)
	defer r.drv.Close()
	client := r.client
	ctx := viewertest.NewContext(context.Background(), client)
	locationType := client.LocationType.Create().
		SetName("LocationType").
		SaveX(ctx)
	location1 := client.Location.Create().
		SetName("Location1").
		SetType(locationType).
		SaveX(ctx)
	location2 := client.Location.Create().
		SetName("Location2").
		SetType(locationType).
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
		SetLocation(location1).
		SaveX(ctx)
	propToDelete2 := client.Property.Create().
		SetType(propTypeToDelete).
		SetStringVal("Prop2").
		SetLocation(location2).
		SaveX(ctx)
	propToDelete3 := client.Property.Create().
		SetType(propTypeToDelete2).
		SetBoolVal(true).
		SetLocation(location1).
		SaveX(ctx)
	prop := client.Property.Create().
		SetType(propType).
		SetIntVal(28).
		SetLocation(location1).
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

func TestGarbageCollectServices(t *testing.T) {
	r := newJobsTestResolver(t)
	defer r.drv.Close()
	client := r.client
	ctx := viewertest.NewContext(context.Background(), client)

	equipType := client.EquipmentType.Create().
		SetName("equipType").
		SaveX(ctx)
	pType := client.PropertyType.Create().
		SetName("p1").
		SetType("string").
		SetBoolVal(true).
		SetEquipmentType(equipType).
		SaveX(ctx)

	sType := client.ServiceType.Create().
		SetName("serviceType").
		AddPropertyTypes(pType).
		SetIsDeleted(true).
		SaveX(ctx)

	epType := client.ServiceEndpointDefinition.Create().
		SetName("ep1").
		SetEquipmentType(equipType).
		SetIndex(0).
		SetServiceType(sType).
		SaveX(ctx)

	eq := client.Equipment.Create().
		SetName("equip").
		SetType(equipType).
		SaveX(ctx)

	prop := client.Property.Create().
		SetType(pType).
		SetBoolVal(true).
		SetEquipment(eq).
		SaveX(ctx)

	s1 := client.Service.Create().
		SetName("s1").
		SetType(sType).
		SetStatus("PENDING").
		SaveX(ctx)
	ep := client.ServiceEndpoint.Create().
		SetDefinition(epType).
		SetService(s1).
		SetEquipment(eq).
		SaveX(ctx)

	s2 := client.Service.Create().
		SetName("s2").
		SetType(sType).
		AddProperties(prop).
		SetStatus("PENDING").
		SaveX(ctx)

	require.True(t, client.ServiceType.Query().Where(servicetype.ID(sType.ID)).ExistX(ctx))
	require.True(t, client.ServiceEndpointDefinition.Query().Where(serviceendpointdefinition.ID(epType.ID)).ExistX(ctx))
	require.True(t, client.PropertyType.Query().Where(propertytype.ID(pType.ID)).ExistX(ctx))

	require.True(t, client.Service.Query().Where(service.ID(s1.ID)).ExistX(ctx))
	require.True(t, client.Service.Query().Where(service.ID(s2.ID)).ExistX(ctx))
	require.True(t, client.ServiceEndpoint.Query().Where(serviceendpoint.ID(ep.ID)).ExistX(ctx))
	require.True(t, client.Property.Query().Where(property.ID(prop.ID)).ExistX(ctx))

	err := r.jobsRunner.collectServices(ctx)
	require.NoError(t, err)
	require.False(t, client.ServiceType.Query().Where(servicetype.ID(sType.ID)).ExistX(ctx))
	require.False(t, client.ServiceEndpointDefinition.Query().Where(serviceendpointdefinition.ID(epType.ID)).ExistX(ctx))
	require.False(t, client.PropertyType.Query().Where(propertytype.ID(pType.ID)).ExistX(ctx))

	require.False(t, client.Service.Query().Where(service.ID(s1.ID)).ExistX(ctx))
	require.False(t, client.Service.Query().Where(service.ID(s2.ID)).ExistX(ctx))
	require.False(t, client.ServiceEndpoint.Query().Where(serviceendpoint.ID(ep.ID)).ExistX(ctx))
	require.False(t, client.Property.Query().Where(property.ID(prop.ID)).ExistX(ctx))
}
