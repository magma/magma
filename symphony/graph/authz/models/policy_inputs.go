// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package models

type PermissionValue string

const (
	PermissionValueYes         PermissionValue = "YES"
	PermissionValueNo          PermissionValue = "NO"
	PermissionValueByCondition PermissionValue = "BY_CONDITION"
)

type BasicPermissionRuleInput struct {
	IsAllowed PermissionValue `json:"isAllowed"`
}

type BasicCUDInput struct {
	Create *BasicPermissionRuleInput `json:"create"`
	Update *BasicPermissionRuleInput `json:"update"`
	Delete *BasicPermissionRuleInput `json:"delete"`
}

type BasicWorkforceCUDInput struct {
	Create            *BasicPermissionRuleInput `json:"create"`
	Update            *BasicPermissionRuleInput `json:"update"`
	Delete            *BasicPermissionRuleInput `json:"delete"`
	Assign            *BasicPermissionRuleInput `json:"assign"`
	TransferOwnership *BasicPermissionRuleInput `json:"transferOwnership"`
}

type InventoryPolicyInput struct {
	Read          *BasicPermissionRuleInput `json:"read"`
	Location      *BasicCUDInput            `json:"location"`
	Equipment     *BasicCUDInput            `json:"equipment"`
	EquipmentType *BasicCUDInput            `json:"equipmentType"`
	LocationType  *BasicCUDInput            `json:"locationType"`
	PortType      *BasicCUDInput            `json:"portType"`
	ServiceType   *BasicCUDInput            `json:"serviceType"`
}

type WorkforcePolicyInput struct {
	Read      *BasicPermissionRuleInput `json:"read"`
	Data      *BasicWorkforceCUDInput   `json:"data"`
	Templates *BasicCUDInput            `json:"templates"`
}
