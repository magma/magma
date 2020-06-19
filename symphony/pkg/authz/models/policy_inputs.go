// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package models

type BasicPermissionRuleInput struct {
	IsAllowed PermissionValue `json:"isAllowed"`
}

type LocationPermissionRuleInput struct {
	IsAllowed       PermissionValue `json:"isAllowed"`
	LocationTypeIds []int           `json:"locationIds"`
}

type WorkforcePermissionRuleInput struct {
	IsAllowed        PermissionValue `json:"isAllowed"`
	ProjectTypeIds   []int           `json:"projectTypeIds"`
	WorkOrderTypeIds []int           `json:"workOrderTypeIds"`
}

type BasicCUDInput struct {
	Create *BasicPermissionRuleInput `json:"create"`
	Update *BasicPermissionRuleInput `json:"update"`
	Delete *BasicPermissionRuleInput `json:"delete"`
}

type LocationCUDInput struct {
	Create *BasicPermissionRuleInput    `json:"create"`
	Update *LocationPermissionRuleInput `json:"update"`
	Delete *BasicPermissionRuleInput    `json:"delete"`
}

type WorkforceCUDInput struct {
	Create            *BasicPermissionRuleInput `json:"create"`
	Update            *BasicPermissionRuleInput `json:"update"`
	Delete            *BasicPermissionRuleInput `json:"delete"`
	Assign            *BasicPermissionRuleInput `json:"assign"`
	TransferOwnership *BasicPermissionRuleInput `json:"transferOwnership"`
}

type InventoryPolicyInput struct {
	Read          *BasicPermissionRuleInput `json:"read"`
	Location      *LocationCUDInput         `json:"location"`
	Equipment     *BasicCUDInput            `json:"equipment"`
	EquipmentType *BasicCUDInput            `json:"equipmentType"`
	LocationType  *BasicCUDInput            `json:"locationType"`
	PortType      *BasicCUDInput            `json:"portType"`
	ServiceType   *BasicCUDInput            `json:"serviceType"`
}

type WorkforcePolicyInput struct {
	Read      *WorkforcePermissionRuleInput `json:"read"`
	Data      *WorkforceCUDInput            `json:"data"`
	Templates *BasicCUDInput                `json:"templates"`
}
