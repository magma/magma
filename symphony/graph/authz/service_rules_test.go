// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
)

func TestServiceTypeWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	serviceType := c.ServiceType.Create().
		SetName("ServiceType").
		SaveX(ctx)
	createServiceType := func(ctx context.Context) error {
		_, err := c.ServiceType.Create().
			SetName("NewServiceType").
			Save(ctx)
		return err
	}
	updateServiceType := func(ctx context.Context) error {
		return c.ServiceType.UpdateOne(serviceType).
			SetName("NewName").
			Exec(ctx)
	}
	deleteServiceType := func(ctx context.Context) error {
		return c.ServiceType.DeleteOne(serviceType).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		getCud: func(p *models.PermissionSettings) *models.Cud {
			return p.InventoryPolicy.ServiceType
		},
		create: createServiceType,
		update: updateServiceType,
		delete: deleteServiceType,
	})
}

func TestServiceTypeUpdateWithIsDeleted(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	serviceType := c.ServiceType.Create().
		SetName("ServiceType").
		SaveX(ctx)
	updateServiceType := func(ctx context.Context) error {
		return c.ServiceType.UpdateOne(serviceType).
			SetIsDeleted(false).
			Exec(ctx)
	}
	deleteServiceType := func(ctx context.Context) error {
		return c.ServiceType.UpdateOne(serviceType).
			SetIsDeleted(true).
			Exec(ctx)
	}
	tests := []policyTest{
		{
			operationName: "Update",
			appendPermissions: func(p *models.PermissionSettings) {
				p.InventoryPolicy.ServiceType.Update.IsAllowed = models2.PermissionValueYes
			},
			operation: updateServiceType,
		},
		{
			operationName: "UpdateWithDelete",
			initialPermissions: func(p *models.PermissionSettings) {
				p.InventoryPolicy.ServiceType.Update.IsAllowed = models2.PermissionValueYes
			},
			appendPermissions: func(p *models.PermissionSettings) {
				p.InventoryPolicy.ServiceType.Delete.IsAllowed = models2.PermissionValueYes
			},
			operation: deleteServiceType,
		},
	}
	runPolicyTest(t, tests)
}

func TestServiceWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)

	srvType1 := c.ServiceType.Create().
		SetName("test service type1").
		SetHasCustomer(false).
		SaveX(ctx)
	service := c.Service.Create().
		SetName("test service 1").
		SetTypeID(srvType1.ID).
		SetStatus("PLANNED").
		SaveX(ctx)
	createService := func(ctx context.Context) error {
		_, err := c.Service.Create().
			SetName("new service").
			SetTypeID(srvType1.ID).
			SetStatus("PLANNED").
			Save(ctx)
		return err
	}
	updateService := func(ctx context.Context) error {
		return c.Service.UpdateOne(service).
			SetName("NewName").
			SetExternalID("123").
			Exec(ctx)
	}
	deleteService := func(ctx context.Context) error {
		return c.Service.DeleteOne(service).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		getCud: func(p *models.PermissionSettings) *models.Cud {
			return p.InventoryPolicy.Equipment
		},
		create: createService,
		update: updateService,
		delete: deleteService,
	})
}

func TestServiceEndpointsWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)

	equipmentType := c.EquipmentType.Create().
		SetName("EquipmentType").
		SaveX(ctx)

	endpointDef1 := c.ServiceEndpointDefinition.Create().
		SetName("ep1").SetIndex(0).SetEquipmentType(equipmentType).SaveX(ctx)
	endpointDef2 := c.ServiceEndpointDefinition.Create().
		SetName("ep2").SetIndex(1).SetEquipmentType(equipmentType).SaveX(ctx)

	srvType1 := c.ServiceType.Create().
		SetName("test service type1").
		SetHasCustomer(false).
		AddEndpointDefinitions(endpointDef1, endpointDef2).
		SaveX(ctx)

	equipment1 := c.Equipment.Create().
		SetName("Equipment").
		SetType(equipmentType).
		SaveX(ctx)

	equipment2 := c.Equipment.Create().
		SetName("Equipment2").
		SetType(equipmentType).
		SaveX(ctx)
	equipment3 := c.Equipment.Create().
		SetName("Equipment3").
		SetType(equipmentType).
		SaveX(ctx)
	service := c.Service.Create().
		SetName("test service 1").
		SetTypeID(srvType1.ID).
		SetStatus("PLANNED").
		SaveX(ctx)

	endpoint1 := c.ServiceEndpoint.Create().SetDefinition(endpointDef1).SetEquipment(equipment1).
		SetService(service).SaveX(ctx)

	createServiceEP := func(ctx context.Context) error {
		_, err := c.ServiceEndpoint.Create().SetDefinition(endpointDef2).SetEquipment(equipment2).
			SetService(service).Save(ctx)
		return err
	}
	updateServiceEP := func(ctx context.Context) error {
		return c.ServiceEndpoint.UpdateOne(endpoint1).
			SetEquipment(equipment3).
			Exec(ctx)
	}
	deleteServiceEP := func(ctx context.Context) error {
		return c.ServiceEndpoint.DeleteOne(endpoint1).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.InventoryPolicy.Equipment.Update.IsAllowed = models2.PermissionValueYes
		},
		create: createServiceEP,
		update: updateServiceEP,
		delete: deleteServiceEP,
	})
}

func TestServiceEndpointDefinitionWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)

	equipmentType := c.EquipmentType.Create().
		SetName("EquipmentType").
		SaveX(ctx)

	serviceType := c.ServiceType.Create().
		SetName("ServiceType").
		SaveX(ctx)

	endpointDef1 := c.ServiceEndpointDefinition.Create().
		SetServiceType(serviceType).
		SetName("ep1").
		SetIndex(0).
		SetEquipmentType(equipmentType).
		SaveX(ctx)

	createServiceEndpointDefinition := func(ctx context.Context) error {
		_, err := c.ServiceEndpointDefinition.Create().
			SetServiceType(serviceType).
			SetName("ep2").
			SetIndex(1).
			SetEquipmentType(equipmentType).
			Save(ctx)
		return err
	}
	updateServiceEndpointDefinition := func(ctx context.Context) error {
		return c.ServiceEndpointDefinition.UpdateOne(endpointDef1).
			SetName("NewName").
			Exec(ctx)
	}
	deleteServiceEndpointDefinition := func(ctx context.Context) error {
		return c.ServiceEndpointDefinition.DeleteOne(endpointDef1).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.InventoryPolicy.ServiceType.Update.IsAllowed = models2.PermissionValueYes
		},
		create: createServiceEndpointDefinition,
		update: updateServiceEndpointDefinition,
		delete: deleteServiceEndpointDefinition,
	})
}
