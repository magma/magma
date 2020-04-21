// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"
)

type permissionsPolicyResolver struct{}

func (permissionsPolicyResolver) covertToInventoryPolicy(input *models2.InventoryPolicyInput) *models.InventoryPolicy {
	res := authz.NewInventoryPolicy(false, false)
	if input == nil {
		return res
	}
	res.Read = authz.AppendBasicPermissionRule(res.Read, input.Read)
	res.Location = authz.AppendCUD(res.Location, input.Location)
	res.Equipment = authz.AppendCUD(res.Equipment, input.Equipment)
	res.EquipmentType = authz.AppendCUD(res.EquipmentType, input.EquipmentType)
	res.LocationType = authz.AppendCUD(res.LocationType, input.LocationType)
	res.PortType = authz.AppendCUD(res.PortType, input.PortType)
	res.ServiceType = authz.AppendCUD(res.ServiceType, input.ServiceType)
	return res
}

func (permissionsPolicyResolver) covertToWorkforcePolicy(input *models2.WorkforcePolicyInput) *models.WorkforcePolicy {
	res := authz.NewWorkforcePolicy(false, false)
	if input == nil {
		return res
	}
	res.Read = authz.AppendBasicPermissionRule(res.Read, input.Read)
	res.Data = authz.AppendWorkforceCUD(res.Data, input.Data)
	res.Templates = authz.AppendCUD(res.Templates, input.Templates)
	return res
}

func (r permissionsPolicyResolver) Policy(ctx context.Context, obj *ent.PermissionsPolicy) (models.SystemPolicy, error) {
	if obj.InventoryPolicy != nil {
		return r.covertToInventoryPolicy(obj.InventoryPolicy), nil
	}
	return r.covertToWorkforcePolicy(obj.WorkforcePolicy), nil
}

func (permissionsPolicyResolver) Groups(ctx context.Context, obj *ent.PermissionsPolicy) ([]*ent.UsersGroup, error) {
	return obj.QueryGroups().All(ctx)
}

func (mutationResolver) AddPolicy(ctx context.Context, input models.AddPermissionsPolicyInput) (*ent.PermissionsPolicy, error) {
	client := ent.FromContext(ctx)
	mutation := client.PermissionsPolicy.Create().
		SetName(input.Name).
		SetNillableDescription(input.Description).
		SetNillableIsGlobal(input.IsGlobal)
	if input.InventoryInput != nil && input.WorkforceInput != nil {
		return nil, fmt.Errorf("policy cannot be of both inventory and workforce types")
	}
	switch {
	case input.InventoryInput != nil:
		mutation.SetInventoryPolicy(input.InventoryInput)
	case input.WorkforceInput != nil:
		mutation.SetWorkforcePolicy(input.WorkforceInput)
	default:
		return nil, fmt.Errorf("no policy found in input")
	}
	policy, err := mutation.Save(ctx)
	if ent.IsConstraintError(err) {
		return nil, fmt.Errorf("policy with the given name already exists: %s", input.Name)
	}
	return policy, err
}
