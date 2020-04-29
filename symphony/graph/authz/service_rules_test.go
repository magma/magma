// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
)

func TestServiceTypeWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := context.Background()
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
