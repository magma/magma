// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/viewer"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	models2 "github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
)

func getPropertyTypeCudOperations(ctx context.Context, c *ent.Client, setParent func(*ent.PropertyTypeCreate) *ent.PropertyTypeCreate) cudOperations {
	propertyTypeQuery := c.PropertyType.Create().
		SetName("OldPropertyType").
		SetType("string")
	propertyTypeQuery = setParent(propertyTypeQuery)
	propertyType := propertyTypeQuery.SaveX(ctx)
	createPropertyType := func(ctx context.Context) error {
		propertyTypeQuery := c.PropertyType.Create().
			SetName("PropertyType").
			SetType("string")
		propertyTypeQuery = setParent(propertyTypeQuery)
		_, err := propertyTypeQuery.Save(ctx)
		return err
	}
	updatePropertyType := func(ctx context.Context) error {
		return c.PropertyType.UpdateOne(propertyType).
			SetName("NewName").
			Exec(ctx)
	}
	deletePropertyType := func(ctx context.Context) error {
		return c.PropertyType.DeleteOne(propertyType).
			Exec(ctx)
	}
	return cudOperations{
		create: createPropertyType,
		update: updatePropertyType,
		delete: deletePropertyType,
	}
}

func getPropertyCudOperations(ctx context.Context, c *ent.Client, setParent func(*ent.PropertyCreate) *ent.PropertyCreate, setParent2 func(*ent.PropertyCreate) *ent.PropertyCreate) cudOperations {
	propertyQuery := c.Property.Create().
		SetStringVal("value")
	propertyQuery = setParent(propertyQuery)
	property := propertyQuery.SaveX(ctx)
	createProperty := func(ctx context.Context) error {
		propertyQuery := c.Property.Create()
		propertyQuery = setParent2(propertyQuery)
		_, err := propertyQuery.Save(ctx)
		return err
	}
	updateProperty := func(ctx context.Context) error {
		return c.Property.UpdateOne(property).
			SetStringVal("newValue").
			Exec(ctx)
	}
	deleteProperty := func(ctx context.Context) error {
		return c.Property.DeleteOne(property).
			Exec(ctx)
	}
	return cudOperations{
		create: createProperty,
		update: updateProperty,
		delete: deleteProperty,
	}
}

func TestLocationTypePropertyTypePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	locationType := c.LocationType.Create().
		SetName("LocationType").
		SaveX(ctx)
	cudOperations := getPropertyTypeCudOperations(ctx, c, func(ptc *ent.PropertyTypeCreate) *ent.PropertyTypeCreate {
		return ptc.SetLocationType(locationType)
	})
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.InventoryPolicy.LocationType.Update.IsAllowed = models2.PermissionValueYes
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
	})
}

func TestEquipmentTypePropertyTypePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	equipmentType := c.EquipmentType.Create().
		SetName("EquipmentType").
		SaveX(ctx)
	cudOperations := getPropertyTypeCudOperations(ctx, c, func(ptc *ent.PropertyTypeCreate) *ent.PropertyTypeCreate {
		return ptc.SetEquipmentType(equipmentType)
	})
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.InventoryPolicy.EquipmentType.Update.IsAllowed = models2.PermissionValueYes
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
	})
}

func TestPortTypePropertyTypePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	portType := c.EquipmentPortType.Create().
		SetName("EquipmentPortType").
		SaveX(ctx)
	cudOperations := getPropertyTypeCudOperations(ctx, c, func(ptc *ent.PropertyTypeCreate) *ent.PropertyTypeCreate {
		return ptc.SetEquipmentPortType(portType)
	})
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.InventoryPolicy.PortType.Update.IsAllowed = models2.PermissionValueYes
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
	})
}

func TestServiceTypePropertyTypePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	serviceType := c.ServiceType.Create().
		SetName("ServiceType").
		SaveX(ctx)
	cudOperations := getPropertyTypeCudOperations(ctx, c, func(ptc *ent.PropertyTypeCreate) *ent.PropertyTypeCreate {
		return ptc.SetServiceType(serviceType)
	})
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.InventoryPolicy.ServiceType.Update.IsAllowed = models2.PermissionValueYes
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
	})
}

func TestWorkOrderTypePropertyTypePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	workOrderType := c.WorkOrderType.Create().
		SetName("WorkOrderType").
		SaveX(ctx)
	cudOperations := getPropertyTypeCudOperations(ctx, c, func(ptc *ent.PropertyTypeCreate) *ent.PropertyTypeCreate {
		return ptc.SetWorkOrderType(workOrderType)
	})
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Templates.Update.IsAllowed = models2.PermissionValueYes
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
	})
}

func TestLocationPropertyPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	locationType := c.LocationType.Create().
		SetName("LocationType").
		SaveX(ctx)
	PropertyType := c.PropertyType.Create().
		SetName("PropertyType").
		SetType("string").
		SetLocationType(locationType).
		SaveX(ctx)

	PropertyType2 := c.PropertyType.Create().
		SetName("PropertyType2").
		SetType("string").
		SetLocationType(locationType).
		SaveX(ctx)

	location := c.Location.Create().
		SetName("Location").
		SetType(locationType).
		SaveX(ctx)

	cudOperations := getPropertyCudOperations(ctx, c, func(ptc *ent.PropertyCreate) *ent.PropertyCreate {
		return ptc.SetLocation(location).SetType(PropertyType)
	}, func(ptc *ent.PropertyCreate) *ent.PropertyCreate {
		return ptc.SetLocation(location).SetType(PropertyType2)
	})
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.InventoryPolicy.Location.Update.IsAllowed = models2.PermissionValueYes
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
	})
}

func TestLocationByConditionPropertyPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	locationType := c.LocationType.Create().
		SetName("LocationType").
		SaveX(ctx)
	locationType2 := c.LocationType.Create().
		SetName("LocationType2").
		SaveX(ctx)
	PropertyType := c.PropertyType.Create().
		SetName("PropertyType").
		SetType("string").
		SetLocationType(locationType).
		SaveX(ctx)
	PropertyType2 := c.PropertyType.Create().
		SetName("PropertyType2").
		SetType("string").
		SetLocationType(locationType).
		SaveX(ctx)
	location := c.Location.Create().
		SetName("Location").
		SetType(locationType).
		SaveX(ctx)

	cudOperations := getPropertyCudOperations(ctx, c, func(ptc *ent.PropertyCreate) *ent.PropertyCreate {
		return ptc.SetLocation(location).SetType(PropertyType)
	}, func(ptc *ent.PropertyCreate) *ent.PropertyCreate {
		return ptc.SetLocation(location).SetType(PropertyType2)
	})
	runCudPolicyTest(t, cudPolicyTest{
		initialPermissions: func(p *models.PermissionSettings) {
			p.InventoryPolicy.Location.Update.IsAllowed = models2.PermissionValueByCondition
			p.InventoryPolicy.Location.Update.LocationTypeIds = []int{locationType2.ID}
		},
		appendPermissions: func(p *models.PermissionSettings) {
			p.InventoryPolicy.Location.Update.LocationTypeIds = []int{locationType.ID, locationType2.ID}
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
	})
}

func TestEquipmentPropertyPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	equipmentType := c.EquipmentType.Create().
		SetName("EquipmentType").
		SaveX(ctx)
	PropertyType := c.PropertyType.Create().
		SetName("PropertyType").
		SetType("string").
		SetEquipmentType(equipmentType).
		SaveX(ctx)
	PropertyType2 := c.PropertyType.Create().
		SetName("PropertyType2").
		SetType("string").
		SetEquipmentType(equipmentType).
		SaveX(ctx)
	equipment := c.Equipment.Create().
		SetName("Equipment").
		SetType(equipmentType).
		SaveX(ctx)

	cudOperations := getPropertyCudOperations(ctx, c, func(ptc *ent.PropertyCreate) *ent.PropertyCreate {
		return ptc.SetEquipment(equipment).SetType(PropertyType)
	}, func(ptc *ent.PropertyCreate) *ent.PropertyCreate {
		return ptc.SetEquipment(equipment).SetType(PropertyType2)
	})
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.InventoryPolicy.Equipment.Update.IsAllowed = models2.PermissionValueYes
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
	})
}

func TestProjectPropertyPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	projectType := c.ProjectType.Create().
		SetName("ProjectType").
		SaveX(ctx)
	PropertyType := c.PropertyType.Create().
		SetName("PropertyType").
		SetType("string").
		SetProjectType(projectType).
		SaveX(ctx)
	PropertyType2 := c.PropertyType.Create().
		SetName("PropertyType2").
		SetType("string").
		SetProjectType(projectType).
		SaveX(ctx)
	project := c.Project.Create().
		SetName("Project").
		SetType(projectType).
		SaveX(ctx)

	cudOperations := getPropertyCudOperations(ctx, c, func(ptc *ent.PropertyCreate) *ent.PropertyCreate {
		return ptc.SetProject(project).SetType(PropertyType)
	}, func(ptc *ent.PropertyCreate) *ent.PropertyCreate {
		return ptc.SetProject(project).SetType(PropertyType2)
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

func TestProjectByConditionPropertyPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	projectType := c.ProjectType.Create().
		SetName("ProjectType").
		SaveX(ctx)
	projectType2 := c.ProjectType.Create().
		SetName("ProjectType2").
		SaveX(ctx)
	PropertyType := c.PropertyType.Create().
		SetName("PropertyType").
		SetType("string").
		SetProjectType(projectType).
		SaveX(ctx)
	PropertyType2 := c.PropertyType.Create().
		SetName("PropertyType2").
		SetType("string").
		SetProjectType(projectType).
		SaveX(ctx)
	project := c.Project.Create().
		SetName("Project").
		SetType(projectType).
		SaveX(ctx)

	cudOperations := getPropertyCudOperations(ctx, c, func(ptc *ent.PropertyCreate) *ent.PropertyCreate {
		return ptc.SetProject(project).SetType(PropertyType)
	}, func(ptc *ent.PropertyCreate) *ent.PropertyCreate {
		return ptc.SetProject(project).SetType(PropertyType2)
	})
	runCudPolicyTest(t, cudPolicyTest{
		initialPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Read.IsAllowed = models2.PermissionValueYes
			p.WorkforcePolicy.Data.Update.IsAllowed = models2.PermissionValueByCondition
			p.WorkforcePolicy.Data.Update.ProjectTypeIds = []int{projectType2.ID}
		},
		appendPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Data.Update.ProjectTypeIds = []int{projectType.ID, projectType2.ID}
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
	})
}

func TestWorkOrderPropertyBasedOnOwnerPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)

	workOrderType, workOrder := prepareWorkOrderData(ctx, c)
	u := viewer.MustGetOrCreateUser(ctx, "MyOwner", user.RoleUSER)
	c.WorkOrder.UpdateOne(workOrder).
		SetOwner(u).
		ExecX(ctx)

	PropertyType := c.PropertyType.Create().
		SetName("PropertyType").
		SetType("string").
		SetWorkOrderType(workOrderType).
		SaveX(ctx)
	PropertyType2 := c.PropertyType.Create().
		SetName("PropertyType2").
		SetType("string").
		SetWorkOrderType(workOrderType).
		SaveX(ctx)
	withPermissionsContext := viewertest.NewContext(ctx, c,
		viewertest.WithUser("MyOwner"),
		viewertest.WithRole(user.RoleUSER),
		viewertest.WithPermissions(authz.EmptyPermissions()))
	noPermissionsContext := viewertest.NewContext(ctx, c,
		viewertest.WithUser("user"),
		viewertest.WithRole(user.RoleUSER),
		viewertest.WithPermissions(authz.EmptyPermissions()))

	cudOperations := getPropertyCudOperations(ctx, c, func(ptc *ent.PropertyCreate) *ent.PropertyCreate {
		return ptc.SetWorkOrder(workOrder).SetType(PropertyType)
	}, func(ptc *ent.PropertyCreate) *ent.PropertyCreate {
		return ptc.SetWorkOrder(workOrder).SetType(PropertyType2)
	})
	tests := []contextBasedPolicyTest{
		{
			operationName:          "Create",
			noPermissionsContext:   noPermissionsContext,
			withPermissionsContext: withPermissionsContext,
			operation:              cudOperations.create,
		},
		{
			operationName:          "Update",
			noPermissionsContext:   noPermissionsContext,
			withPermissionsContext: withPermissionsContext,
			operation:              cudOperations.update,
		},
		{
			operationName:          "Delete",
			noPermissionsContext:   noPermissionsContext,
			withPermissionsContext: withPermissionsContext,
			operation:              cudOperations.delete,
		},
	}
	runContextBasedPolicyTest(t, tests)
}

func TestPropertyOfWorkOrderReadPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	woType1, wo1 := prepareWorkOrderData(ctx, c)
	woType2, wo2 := prepareWorkOrderData(ctx, c)
	pType1 := c.PropertyType.Create().
		SetName("pType1").
		SetType("string").
		SetWorkOrderType(woType1).
		SaveX(ctx)
	pType2 := c.PropertyType.Create().
		SetName("pType2").
		SetType("string").
		SetWorkOrderType(woType2).
		SaveX(ctx)
	c.Property.Create().
		SetType(pType1).
		SetStringVal("Hi").
		SetWorkOrder(wo1).
		SaveX(ctx)
	c.Property.Create().
		SetType(pType2).
		SetStringVal("Hi").
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
		count, err := c.Property.Query().Count(permissionsContext)
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
		count, err := c.Property.Query().Count(permissionsContext)
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
		count, err := c.Property.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})
}

func TestPropertyOfProjectReadPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	projectType1, project1 := prepareProjectData(ctx, c)
	projectType2, project2 := prepareProjectData(ctx, c)
	pType1 := c.PropertyType.Create().
		SetName("pType1").
		SetType("string").
		SetProjectType(projectType1).
		SaveX(ctx)
	pType2 := c.PropertyType.Create().
		SetName("pType2").
		SetType("string").
		SetProjectType(projectType2).
		SaveX(ctx)
	c.Property.Create().
		SetType(pType1).
		SetStringVal("Hi").
		SetProject(project1).
		SaveX(ctx)
	c.Property.Create().
		SetType(pType2).
		SetStringVal("Hi").
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
		count, err := c.Property.Query().Count(permissionsContext)
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
		count, err := c.Property.Query().Count(permissionsContext)
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
		count, err := c.Property.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})
}
