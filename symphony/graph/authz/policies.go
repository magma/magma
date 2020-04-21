// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

func NewBasicPermissionRule(allowed bool) *models.BasicPermissionRule {
	rule := models.PermissionValueNo
	if allowed {
		rule = models.PermissionValueYes
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
