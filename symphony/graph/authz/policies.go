// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	models2 "github.com/facebookincubator/symphony/graph/authz/models"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

var allowedEnums = map[models2.PermissionValue]int{
	models2.PermissionValueNo:          1,
	models2.PermissionValueByCondition: 2,
	models2.PermissionValueYes:         3,
}

func NewBasicPermissionRule(allowed bool) *models.BasicPermissionRule {
	rule := models2.PermissionValueNo
	if allowed {
		rule = models2.PermissionValueYes
	}
	return &models.BasicPermissionRule{IsAllowed: rule}
}

func NewCUD(allowed bool) *models.Cud {
	return &models.Cud{
		Create: NewBasicPermissionRule(allowed),
		Update: NewBasicPermissionRule(allowed),
		Delete: NewBasicPermissionRule(allowed),
	}
}

func NewWorkforceCUD(allowed bool) *models.WorkforceCud {
	return &models.WorkforceCud{
		Create:            NewBasicPermissionRule(allowed),
		Update:            NewBasicPermissionRule(allowed),
		Delete:            NewBasicPermissionRule(allowed),
		Assign:            NewBasicPermissionRule(allowed),
		TransferOwnership: NewBasicPermissionRule(allowed),
	}
}

func NewInventoryPolicy(readAllowed, writeAllowed bool) *models.InventoryPolicy {
	return &models.InventoryPolicy{
		Read:          NewBasicPermissionRule(readAllowed),
		Location:      NewCUD(writeAllowed),
		Equipment:     NewCUD(writeAllowed),
		EquipmentType: NewCUD(writeAllowed),
		LocationType:  NewCUD(writeAllowed),
		PortType:      NewCUD(writeAllowed),
		ServiceType:   NewCUD(writeAllowed),
	}
}

func NewWorkforcePolicy(readAllowed, writeAllowed bool) *models.WorkforcePolicy {
	return &models.WorkforcePolicy{
		Read:      NewBasicPermissionRule(readAllowed),
		Data:      NewWorkforceCUD(writeAllowed),
		Templates: NewCUD(writeAllowed),
	}
}

func NewAdministrativePolicy(u *ent.User) *models.AdministrativePolicy {
	allowed := u.Role == user.RoleADMIN || u.Role == user.RoleOWNER
	return &models.AdministrativePolicy{Access: NewBasicPermissionRule(allowed)}
}

func AppendBasicPermissionRule(rule *models.BasicPermissionRule, addRule *models2.BasicPermissionRuleInput) *models.BasicPermissionRule {
	if addRule != nil && allowedEnums[addRule.IsAllowed] >= allowedEnums[rule.IsAllowed] {
		rule.IsAllowed = addRule.IsAllowed
	}
	return rule
}

func AppendCUD(cud *models.Cud, addCUD *models2.BasicCUDInput) *models.Cud {
	if addCUD == nil {
		return cud
	}
	cud.Create = AppendBasicPermissionRule(cud.Create, addCUD.Create)
	cud.Delete = AppendBasicPermissionRule(cud.Delete, addCUD.Delete)
	cud.Update = AppendBasicPermissionRule(cud.Update, addCUD.Update)
	return cud
}

func AppendWorkforceCUD(cud *models.WorkforceCud, addCUD *models2.BasicWorkforceCUDInput) *models.WorkforceCud {
	if addCUD == nil {
		return cud
	}
	cud.Create = AppendBasicPermissionRule(cud.Create, addCUD.Create)
	cud.Delete = AppendBasicPermissionRule(cud.Delete, addCUD.Delete)
	cud.Update = AppendBasicPermissionRule(cud.Update, addCUD.Update)
	cud.Assign = AppendBasicPermissionRule(cud.Assign, addCUD.Assign)
	cud.TransferOwnership = AppendBasicPermissionRule(cud.TransferOwnership, addCUD.TransferOwnership)
	return cud
}
