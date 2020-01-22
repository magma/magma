// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.
package ent

import (
	"context"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/symphony/graph/ent/migrate"
)

// ReadOnly returns a new readonly-client.
//
//	client := client.ReadOnly()
//
func (c *Client) ReadOnly() *Client {
	cfg := config{driver: &readonly{Driver: c.driver}, log: c.log}
	return &Client{
		config:                      cfg,
		Schema:                      migrate.NewSchema(cfg.driver),
		ActionsRule:                 NewActionsRuleClient(cfg),
		CheckListItem:               NewCheckListItemClient(cfg),
		CheckListItemDefinition:     NewCheckListItemDefinitionClient(cfg),
		Comment:                     NewCommentClient(cfg),
		Customer:                    NewCustomerClient(cfg),
		Equipment:                   NewEquipmentClient(cfg),
		EquipmentCategory:           NewEquipmentCategoryClient(cfg),
		EquipmentPort:               NewEquipmentPortClient(cfg),
		EquipmentPortDefinition:     NewEquipmentPortDefinitionClient(cfg),
		EquipmentPortType:           NewEquipmentPortTypeClient(cfg),
		EquipmentPosition:           NewEquipmentPositionClient(cfg),
		EquipmentPositionDefinition: NewEquipmentPositionDefinitionClient(cfg),
		EquipmentType:               NewEquipmentTypeClient(cfg),
		File:                        NewFileClient(cfg),
		FloorPlan:                   NewFloorPlanClient(cfg),
		FloorPlanReferencePoint:     NewFloorPlanReferencePointClient(cfg),
		FloorPlanScale:              NewFloorPlanScaleClient(cfg),
		Hyperlink:                   NewHyperlinkClient(cfg),
		Link:                        NewLinkClient(cfg),
		Location:                    NewLocationClient(cfg),
		LocationType:                NewLocationTypeClient(cfg),
		Project:                     NewProjectClient(cfg),
		ProjectType:                 NewProjectTypeClient(cfg),
		Property:                    NewPropertyClient(cfg),
		PropertyType:                NewPropertyTypeClient(cfg),
		Service:                     NewServiceClient(cfg),
		ServiceEndpoint:             NewServiceEndpointClient(cfg),
		ServiceType:                 NewServiceTypeClient(cfg),
		Survey:                      NewSurveyClient(cfg),
		SurveyCellScan:              NewSurveyCellScanClient(cfg),
		SurveyQuestion:              NewSurveyQuestionClient(cfg),
		SurveyTemplateCategory:      NewSurveyTemplateCategoryClient(cfg),
		SurveyTemplateQuestion:      NewSurveyTemplateQuestionClient(cfg),
		SurveyWiFiScan:              NewSurveyWiFiScanClient(cfg),
		Technician:                  NewTechnicianClient(cfg),
		WorkOrder:                   NewWorkOrderClient(cfg),
		WorkOrderDefinition:         NewWorkOrderDefinitionClient(cfg),
		WorkOrderType:               NewWorkOrderTypeClient(cfg),
	}
}

// ErrReadOnly returns when a readonly user tries to execute a write operation.
var ErrReadOnly = &PermissionError{cause: "permission denied: read-only user"}

// PermissionError represents a permission denied error.
type PermissionError struct {
	cause string
}

func (e PermissionError) Error() string { return e.cause }

type readonly struct {
	dialect.Driver
}

func (r *readonly) Exec(context.Context, string, interface{}, interface{}) error {
	return ErrReadOnly
}

func (r *readonly) Tx(context.Context) (dialect.Tx, error) {
	return nil, ErrReadOnly
}
