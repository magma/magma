// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

func preparePermissionsPolicyData(ctx context.Context) {
	client := ent.FromContext(ctx).PermissionsPolicy
	inventoryPolicyInput := getInventoryPolicyInput()
	workforceInventoryInput := getWorkforcePolicyInput()

	client.Create().
		SetName("policy1").
		SetDescription("empty_policy").
		SetInventoryPolicy(nil).
		SetWorkforcePolicy(nil).
		SaveX(ctx)

	client.Create().
		SetName("policy2").
		SetDescription("inventory_policy").
		SetInventoryPolicy(inventoryPolicyInput).
		SetWorkforcePolicy(nil).
		SaveX(ctx)

	client.Create().
		SetName("policy3").
		SetDescription("workforce_policy").
		SetInventoryPolicy(nil).
		SetWorkforcePolicy(workforceInventoryInput).
		SaveX(ctx)
}

func TestPermissionsPolicySearchByName(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	preparePermissionsPolicyData(ctx)
	ppName1 := "policy1"

	f1 := models.PermissionsPolicyFilterInput{
		FilterType:  models.PermissionsPolicyFilterTypePermissionsPolicyName,
		Operator:    models.FilterOperatorIs,
		StringValue: &ppName1,
	}
	resAll, err := r.Query().PermissionsPolicySearch(ctx, []*models.PermissionsPolicyFilterInput{&f1}, pointer.ToInt(100))
	require.NoError(t, err)
	require.Len(t, resAll.PermissionsPolicies, 1)
	require.Equal(t, resAll.Count, 1)
}
