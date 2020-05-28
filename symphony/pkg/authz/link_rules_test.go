// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/pkg/ent"

	models2 "github.com/facebookincubator/symphony/pkg/authz/models"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
)

func TestLinkWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	equipmentPortType := c.EquipmentPortType.Create().
		SetName("EquipmentPortType").
		SaveX(ctx)
	equipmentPortDefinition1 := c.EquipmentPortDefinition.Create().
		SetName("EquipmentPortDefinition").
		SetEquipmentPortType(equipmentPortType).
		SaveX(ctx)
	equipmentPort1 := c.EquipmentPort.Create().
		SetDefinition(equipmentPortDefinition1).
		SaveX(ctx)
	equipmentPortDefinition2 := c.EquipmentPortDefinition.Create().
		SetName("EquipmentPortDefinition2").
		SetEquipmentPortType(equipmentPortType).
		SaveX(ctx)
	equipmentPort2 := c.EquipmentPort.Create().
		SetDefinition(equipmentPortDefinition2).
		SaveX(ctx)
	equipmentPortDefinition3 := c.EquipmentPortDefinition.Create().
		SetName("EquipmentPortDefinition3").
		SetEquipmentPortType(equipmentPortType).
		SaveX(ctx)
	equipmentPort3 := c.EquipmentPort.Create().
		SetDefinition(equipmentPortDefinition3).
		SaveX(ctx)
	var (
		link *ent.Link
		err  error
	)
	createLink := func(ctx context.Context) error {
		link, err = c.Link.Create().
			AddPortIDs(equipmentPort1.ID, equipmentPort2.ID).
			Save(ctx)
		return err
	}
	updateLink := func(ctx context.Context) error {
		_, err := c.Link.UpdateOne(link).
			AddPortIDs(equipmentPort3.ID).
			Save(ctx)
		return err
	}
	deleteLink := func(ctx context.Context) error {
		return c.Link.DeleteOne(link).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.InventoryPolicy.Equipment.Update.IsAllowed = models2.PermissionValueYes
		},
		create: createLink,
		update: updateLink,
		delete: deleteLink,
	})
}
