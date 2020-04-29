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

func TestWorkorderTypeWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := context.Background()
	workorderType := c.WorkOrderType.Create().
		SetName("WorkOrderType").
		SaveX(ctx)
	createWorkOrderType := func(ctx context.Context) error {
		_, err := c.WorkOrderType.Create().
			SetName("NewWorkOrderType").
			Save(ctx)
		return err
	}
	updateWorkOrderType := func(ctx context.Context) error {
		return c.WorkOrderType.UpdateOne(workorderType).
			SetName("NewName").
			Exec(ctx)
	}
	deleteWorkOrderType := func(ctx context.Context) error {
		return c.WorkOrderType.DeleteOne(workorderType).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		getCud: func(p *models.PermissionSettings) *models.Cud {
			return p.WorkforcePolicy.Templates
		},
		create: createWorkOrderType,
		update: updateWorkOrderType,
		delete: deleteWorkOrderType,
	})
}
