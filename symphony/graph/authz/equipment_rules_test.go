package authz_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
)

func TestEquipmentWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := context.Background()
	equipmentType := c.EquipmentType.Create().
		SetName("EquipmentType").
		SaveX(ctx)
	equipment := c.Equipment.Create().
		SetName("Equipment").
		SetType(equipmentType).
		SaveX(ctx)
	createEquipment := func(ctx context.Context) error {
		_, err := c.Equipment.Create().
			SetName("NewEquipment").
			SetType(equipmentType).
			Save(ctx)
		return err
	}
	updateEquipment := func(ctx context.Context) error {
		return c.Equipment.UpdateOne(equipment).
			SetName("NewName").
			Exec(ctx)
	}
	deleteEquipment := func(ctx context.Context) error {
		return c.Equipment.DeleteOne(equipment).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		getCud: func(p *models.PermissionSettings) *models.Cud {
			return p.InventoryPolicy.Equipment
		},
		create: createEquipment,
		update: updateEquipment,
		delete: deleteEquipment,
	})
}

func TestEquipmentTypeWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := context.Background()
	equipmentType := c.EquipmentType.Create().
		SetName("EquipmentType").
		SaveX(ctx)
	createEquipmentType := func(ctx context.Context) error {
		_, err := c.EquipmentType.Create().
			SetName("NewEquipmentType").
			Save(ctx)
		return err
	}
	updateEquipmentType := func(ctx context.Context) error {
		return c.EquipmentType.UpdateOne(equipmentType).
			SetName("NewName").
			Exec(ctx)
	}
	deleteEquipmentType := func(ctx context.Context) error {
		return c.EquipmentType.DeleteOne(equipmentType).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		getCud: func(p *models.PermissionSettings) *models.Cud {
			return p.InventoryPolicy.EquipmentType
		},
		create: createEquipmentType,
		update: updateEquipmentType,
		delete: deleteEquipmentType,
	})
}

func TestEquipmentPortTypeWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := context.Background()
	equipmentType := c.EquipmentPortType.Create().
		SetName("EquipmentPortType").
		SaveX(ctx)
	createEquipmentPortType := func(ctx context.Context) error {
		_, err := c.EquipmentPortType.Create().
			SetName("NewEquipmentPortType").
			Save(ctx)
		return err
	}
	updateEquipmentPortType := func(ctx context.Context) error {
		return c.EquipmentPortType.UpdateOne(equipmentType).
			SetName("NewName").
			Exec(ctx)
	}
	deleteEquipmentPortType := func(ctx context.Context) error {
		return c.EquipmentPortType.DeleteOne(equipmentType).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		getCud: func(p *models.PermissionSettings) *models.Cud {
			return p.InventoryPolicy.PortType
		},
		create: createEquipmentPortType,
		update: updateEquipmentPortType,
		delete: deleteEquipmentPortType,
	})
}
