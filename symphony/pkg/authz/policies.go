// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/permissionspolicy"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/ent/usersgroup"
	"github.com/facebookincubator/symphony/pkg/viewer"
)

// WritePermissionGroupName is the name of the group that its member has write permission for all symphony.
const WritePermissionGroupName = "Write Permission"

var allowedEnums = map[models.PermissionValue]int{
	models.PermissionValueNo:          1,
	models.PermissionValueByCondition: 2,
	models.PermissionValueYes:         3,
}

func newBasicPermissionRule(allowed bool) *models.BasicPermissionRule {
	rule := models.PermissionValueNo
	if allowed {
		rule = models.PermissionValueYes
	}
	return &models.BasicPermissionRule{IsAllowed: rule}
}

func newLocationPermissionRule(allowed bool) *models.LocationPermissionRule {
	rule := models.PermissionValueNo
	if allowed {
		rule = models.PermissionValueYes
	}
	return &models.LocationPermissionRule{IsAllowed: rule}
}

func newWorkforcePermissionRule(allowed bool) *models.WorkforcePermissionRule {
	rule := models.PermissionValueNo
	if allowed {
		rule = models.PermissionValueYes
	}
	return &models.WorkforcePermissionRule{IsAllowed: rule}
}

func newCUD(allowed bool) *models.Cud {
	return &models.Cud{
		Create: newBasicPermissionRule(allowed),
		Update: newBasicPermissionRule(allowed),
		Delete: newBasicPermissionRule(allowed),
	}
}

func newLocationCUD(allowed bool) *models.LocationCud {
	return &models.LocationCud{
		Create: newLocationPermissionRule(allowed),
		Update: newLocationPermissionRule(allowed),
		Delete: newLocationPermissionRule(allowed),
	}
}

func newWorkforceCUD(allowed bool) *models.WorkforceCud {
	return &models.WorkforceCud{
		Create:            newWorkforcePermissionRule(allowed),
		Update:            newWorkforcePermissionRule(allowed),
		Delete:            newWorkforcePermissionRule(allowed),
		Assign:            newWorkforcePermissionRule(allowed),
		TransferOwnership: newWorkforcePermissionRule(allowed),
	}
}

// NewInventoryPolicy builds an inventory policy based on general restriction on read,write
func NewInventoryPolicy(writeAllowed bool) *models.InventoryPolicy {
	return &models.InventoryPolicy{
		Read:          newBasicPermissionRule(true),
		Location:      newLocationCUD(writeAllowed),
		Equipment:     newCUD(writeAllowed),
		EquipmentType: newCUD(writeAllowed),
		LocationType:  newCUD(writeAllowed),
		PortType:      newCUD(writeAllowed),
		ServiceType:   newCUD(writeAllowed),
	}
}

// NewWorkforcePolicy build a workforce policy based on general restriction on read,write
func NewWorkforcePolicy(readAllowed, writeAllowed bool) *models.WorkforcePolicy {
	return &models.WorkforcePolicy{
		Read:      newWorkforcePermissionRule(readAllowed),
		Data:      newWorkforceCUD(writeAllowed),
		Templates: newCUD(writeAllowed),
	}
}

// NewAdministrativePolicy builds administrative policy of given user
func NewAdministrativePolicy(isAdmin bool) *models.AdministrativePolicy {
	return &models.AdministrativePolicy{
		Access: newBasicPermissionRule(isAdmin),
	}
}

func appendBasicPermissionRule(rule *models.BasicPermissionRule, addRule *models.BasicPermissionRuleInput) *models.BasicPermissionRule {
	if addRule != nil && allowedEnums[addRule.IsAllowed] >= allowedEnums[rule.IsAllowed] {
		rule.IsAllowed = addRule.IsAllowed
	}
	return rule
}

func appendTopLevelLocationPermissionRuleInput(
	rule *models.LocationPermissionRuleInput, bottomRule *models.BasicPermissionRuleInput) *models.LocationPermissionRuleInput {
	if bottomRule == nil || bottomRule.IsAllowed == models.PermissionValueNo {
		return &models.LocationPermissionRuleInput{IsAllowed: models.PermissionValueNo}
	}
	return rule
}

func appendLocationPermissionRule(rule *models.LocationPermissionRule, addRule *models.LocationPermissionRuleInput) *models.LocationPermissionRule {
	if addRule == nil {
		return rule
	}
	if allowedEnums[addRule.IsAllowed] >= allowedEnums[rule.IsAllowed] {
		rule.IsAllowed = addRule.IsAllowed
	}
	switch rule.IsAllowed {
	case models.PermissionValueYes:
		rule.LocationTypeIds = nil
	case models.PermissionValueNo:
		rule.LocationTypeIds = nil
	case models.PermissionValueByCondition:
		rule.LocationTypeIds = append(rule.LocationTypeIds, addRule.LocationTypeIds...)
	}
	return rule
}

func appendTopLevelWorkforcePermissionRuleInput(
	rule *models.WorkforcePermissionRuleInput, bottomRule *models.BasicPermissionRuleInput) *models.WorkforcePermissionRuleInput {
	if bottomRule == nil || bottomRule.IsAllowed == models.PermissionValueNo {
		return &models.WorkforcePermissionRuleInput{IsAllowed: models.PermissionValueNo}
	}
	return rule
}

func appendWorkforcePermissionRule(rule *models.WorkforcePermissionRule, addRule *models.WorkforcePermissionRuleInput) *models.WorkforcePermissionRule {
	if addRule == nil {
		return rule
	}
	if allowedEnums[addRule.IsAllowed] >= allowedEnums[rule.IsAllowed] {
		rule.IsAllowed = addRule.IsAllowed
	}
	switch rule.IsAllowed {
	case models.PermissionValueYes:
		rule.WorkOrderTypeIds = nil
		rule.ProjectTypeIds = nil
	case models.PermissionValueNo:
		rule.WorkOrderTypeIds = nil
		rule.ProjectTypeIds = nil
	case models.PermissionValueByCondition:
		rule.WorkOrderTypeIds = append(rule.WorkOrderTypeIds, addRule.WorkOrderTypeIds...)
		rule.ProjectTypeIds = append(rule.ProjectTypeIds, addRule.ProjectTypeIds...)
	}
	return rule
}

func appendCUD(cud *models.Cud, addCUD *models.BasicCUDInput) *models.Cud {
	if addCUD == nil {
		return cud
	}
	cud.Create = appendBasicPermissionRule(cud.Create, addCUD.Create)
	cud.Delete = appendBasicPermissionRule(cud.Delete, addCUD.Delete)
	cud.Update = appendBasicPermissionRule(cud.Update, addCUD.Update)
	return cud
}

func appendLocationCUD(cud *models.LocationCud, addCUD *models.LocationCUDInput) *models.LocationCud {
	if addCUD == nil {
		return cud
	}
	cud.Create = appendLocationPermissionRule(
		cud.Create, appendTopLevelLocationPermissionRuleInput(addCUD.Update, addCUD.Create))
	cud.Update = appendLocationPermissionRule(cud.Update, addCUD.Update)
	cud.Delete = appendLocationPermissionRule(
		cud.Delete, appendTopLevelLocationPermissionRuleInput(addCUD.Update, addCUD.Delete))
	return cud
}

func appendWorkforceCUD(cud *models.WorkforceCud, readRule *models.WorkforcePermissionRuleInput, addCUD *models.WorkforceCUDInput) *models.WorkforceCud {
	if addCUD == nil {
		return cud
	}
	cud.Create = appendWorkforcePermissionRule(cud.Create, appendTopLevelWorkforcePermissionRuleInput(readRule, addCUD.Create))
	cud.Delete = appendWorkforcePermissionRule(cud.Delete, appendTopLevelWorkforcePermissionRuleInput(readRule, addCUD.Delete))
	cud.Update = appendWorkforcePermissionRule(cud.Update, appendTopLevelWorkforcePermissionRuleInput(readRule, addCUD.Update))
	cud.Assign = appendWorkforcePermissionRule(cud.Assign, appendTopLevelWorkforcePermissionRuleInput(readRule, addCUD.Assign))
	cud.TransferOwnership = appendWorkforcePermissionRule(
		cud.TransferOwnership, appendTopLevelWorkforcePermissionRuleInput(readRule, addCUD.TransferOwnership))
	return cud
}

// AppendInventoryPolicies append a list of inventory policy inputs to a inventory policy
func AppendInventoryPolicies(policy *models.InventoryPolicy, inputs ...*models.InventoryPolicyInput) *models.InventoryPolicy {
	for _, input := range inputs {
		if input == nil {
			continue
		}
		policy.Read = appendBasicPermissionRule(policy.Read, input.Read)
		policy.Location = appendLocationCUD(policy.Location, input.Location)
		policy.Equipment = appendCUD(policy.Equipment, input.Equipment)
		policy.EquipmentType = appendCUD(policy.EquipmentType, input.EquipmentType)
		policy.LocationType = appendCUD(policy.LocationType, input.LocationType)
		policy.PortType = appendCUD(policy.PortType, input.PortType)
		policy.ServiceType = appendCUD(policy.ServiceType, input.ServiceType)
	}
	return policy
}

// AppendInventoryPolicies append a list of workforce policy inputs to a workforce policy
func AppendWorkforcePolicies(policy *models.WorkforcePolicy, inputs ...*models.WorkforcePolicyInput) *models.WorkforcePolicy {
	for _, input := range inputs {
		if input == nil {
			continue
		}
		policy.Read = appendWorkforcePermissionRule(policy.Read, input.Read)
		policy.Data = appendWorkforceCUD(policy.Data, input.Read, input.Data)
		policy.Templates = appendCUD(policy.Templates, input.Templates)
	}
	return policy
}

func permissionPolicies(ctx context.Context, v *viewer.UserViewer) (*models.InventoryPolicy, *models.WorkforcePolicy, error) {
	client := ent.FromContext(ctx)
	userID := v.User().ID
	inventoryPolicy := NewInventoryPolicy(false)
	workforcePolicy := NewWorkforcePolicy(false, false)
	policies, err := client.PermissionsPolicy.Query().
		Where(permissionspolicy.Or(
			permissionspolicy.IsGlobal(true),
			permissionspolicy.HasGroupsWith(
				usersgroup.HasMembersWith(user.ID(userID)),
				usersgroup.StatusEQ(usersgroup.StatusACTIVE),
			))).
		All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot query policies: %w", err)
	}
	for _, policy := range policies {
		switch {
		case policy.InventoryPolicy != nil:
			inventoryPolicy = AppendInventoryPolicies(inventoryPolicy, policy.InventoryPolicy)
		case policy.WorkforcePolicy != nil:
			workforcePolicy = AppendWorkforcePolicies(workforcePolicy, policy.WorkforcePolicy)
		default:
			return nil, nil, fmt.Errorf("empty policy found: %d", policy.ID)
		}
	}
	return inventoryPolicy, workforcePolicy, nil
}

func userHasWritePermissions(ctx context.Context) (bool, error) {
	v := viewer.FromContext(ctx)
	if userHasFullPermissions(v) {
		return true, nil
	}
	if v, ok := v.(*viewer.UserViewer); ok && !v.Features().Enabled(viewer.FeaturePermissionPolicies) {
		return v.User().QueryGroups().
			Where(usersgroup.Name(WritePermissionGroupName)).
			Exist(ctx)
	}
	return false, nil
}

// Permissions builds the aggregated permissions for the given viewer
func Permissions(ctx context.Context) (*models.PermissionSettings, error) {
	writePermissions, err := userHasWritePermissions(ctx)
	if err != nil {
		return nil, err
	}
	v := viewer.FromContext(ctx)
	policiesEnabled := v.Features().Enabled(viewer.FeaturePermissionPolicies)
	inventoryPolicy := NewInventoryPolicy(writePermissions)
	workforcePolicy := NewWorkforcePolicy(true, writePermissions)
	if policiesEnabled {
		if u, ok := v.(*viewer.UserViewer); ok && !writePermissions {
			inventoryPolicy, workforcePolicy, err = permissionPolicies(ctx, u)
			if err != nil {
				return nil, err
			}
		}
	}
	res := models.PermissionSettings{
		// TODO(T64743627): Deprecate CanWrite field
		CanWrite:        writePermissions,
		AdminPolicy:     NewAdministrativePolicy(userHasFullPermissions(v)),
		InventoryPolicy: inventoryPolicy,
		WorkforcePolicy: workforcePolicy,
	}
	return &res, nil
}

func FullPermissions() *models.PermissionSettings {
	return &models.PermissionSettings{
		CanWrite:        true,
		AdminPolicy:     NewAdministrativePolicy(true),
		InventoryPolicy: NewInventoryPolicy(true),
		WorkforcePolicy: NewWorkforcePolicy(true, true),
	}
}

func EmptyPermissions() *models.PermissionSettings {
	return &models.PermissionSettings{
		CanWrite:        false,
		AdminPolicy:     NewAdministrativePolicy(false),
		InventoryPolicy: NewInventoryPolicy(false),
		WorkforcePolicy: NewWorkforcePolicy(false, false),
	}
}

func AdminPermissions() *models.PermissionSettings {
	return &models.PermissionSettings{
		CanWrite:        false,
		AdminPolicy:     NewAdministrativePolicy(true),
		InventoryPolicy: NewInventoryPolicy(false),
		WorkforcePolicy: NewWorkforcePolicy(false, false),
	}
}
