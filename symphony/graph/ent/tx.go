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

// Tx is a transactional client that is created by calling Client.Tx().
type Tx struct {
	config
	// ActionsRule is the client for interacting with the ActionsRule builders.
	ActionsRule *ActionsRuleClient
	// CheckListCategory is the client for interacting with the CheckListCategory builders.
	CheckListCategory *CheckListCategoryClient
	// CheckListItem is the client for interacting with the CheckListItem builders.
	CheckListItem *CheckListItemClient
	// CheckListItemDefinition is the client for interacting with the CheckListItemDefinition builders.
	CheckListItemDefinition *CheckListItemDefinitionClient
	// Comment is the client for interacting with the Comment builders.
	Comment *CommentClient
	// Customer is the client for interacting with the Customer builders.
	Customer *CustomerClient
	// Equipment is the client for interacting with the Equipment builders.
	Equipment *EquipmentClient
	// EquipmentCategory is the client for interacting with the EquipmentCategory builders.
	EquipmentCategory *EquipmentCategoryClient
	// EquipmentPort is the client for interacting with the EquipmentPort builders.
	EquipmentPort *EquipmentPortClient
	// EquipmentPortDefinition is the client for interacting with the EquipmentPortDefinition builders.
	EquipmentPortDefinition *EquipmentPortDefinitionClient
	// EquipmentPortType is the client for interacting with the EquipmentPortType builders.
	EquipmentPortType *EquipmentPortTypeClient
	// EquipmentPosition is the client for interacting with the EquipmentPosition builders.
	EquipmentPosition *EquipmentPositionClient
	// EquipmentPositionDefinition is the client for interacting with the EquipmentPositionDefinition builders.
	EquipmentPositionDefinition *EquipmentPositionDefinitionClient
	// EquipmentType is the client for interacting with the EquipmentType builders.
	EquipmentType *EquipmentTypeClient
	// File is the client for interacting with the File builders.
	File *FileClient
	// FloorPlan is the client for interacting with the FloorPlan builders.
	FloorPlan *FloorPlanClient
	// FloorPlanReferencePoint is the client for interacting with the FloorPlanReferencePoint builders.
	FloorPlanReferencePoint *FloorPlanReferencePointClient
	// FloorPlanScale is the client for interacting with the FloorPlanScale builders.
	FloorPlanScale *FloorPlanScaleClient
	// Hyperlink is the client for interacting with the Hyperlink builders.
	Hyperlink *HyperlinkClient
	// Link is the client for interacting with the Link builders.
	Link *LinkClient
	// Location is the client for interacting with the Location builders.
	Location *LocationClient
	// LocationType is the client for interacting with the LocationType builders.
	LocationType *LocationTypeClient
	// Project is the client for interacting with the Project builders.
	Project *ProjectClient
	// ProjectType is the client for interacting with the ProjectType builders.
	ProjectType *ProjectTypeClient
	// Property is the client for interacting with the Property builders.
	Property *PropertyClient
	// PropertyType is the client for interacting with the PropertyType builders.
	PropertyType *PropertyTypeClient
	// Service is the client for interacting with the Service builders.
	Service *ServiceClient
	// ServiceEndpoint is the client for interacting with the ServiceEndpoint builders.
	ServiceEndpoint *ServiceEndpointClient
	// ServiceType is the client for interacting with the ServiceType builders.
	ServiceType *ServiceTypeClient
	// Survey is the client for interacting with the Survey builders.
	Survey *SurveyClient
	// SurveyCellScan is the client for interacting with the SurveyCellScan builders.
	SurveyCellScan *SurveyCellScanClient
	// SurveyQuestion is the client for interacting with the SurveyQuestion builders.
	SurveyQuestion *SurveyQuestionClient
	// SurveyTemplateCategory is the client for interacting with the SurveyTemplateCategory builders.
	SurveyTemplateCategory *SurveyTemplateCategoryClient
	// SurveyTemplateQuestion is the client for interacting with the SurveyTemplateQuestion builders.
	SurveyTemplateQuestion *SurveyTemplateQuestionClient
	// SurveyWiFiScan is the client for interacting with the SurveyWiFiScan builders.
	SurveyWiFiScan *SurveyWiFiScanClient
	// Technician is the client for interacting with the Technician builders.
	Technician *TechnicianClient
	// User is the client for interacting with the User builders.
	User *UserClient
	// WorkOrder is the client for interacting with the WorkOrder builders.
	WorkOrder *WorkOrderClient
	// WorkOrderDefinition is the client for interacting with the WorkOrderDefinition builders.
	WorkOrderDefinition *WorkOrderDefinitionClient
	// WorkOrderType is the client for interacting with the WorkOrderType builders.
	WorkOrderType *WorkOrderTypeClient
}

// Commit commits the transaction.
func (tx *Tx) Commit() error {
	return tx.config.driver.(*txDriver).tx.Commit()
}

// Rollback rollbacks the transaction.
func (tx *Tx) Rollback() error {
	return tx.config.driver.(*txDriver).tx.Rollback()
}

// Client returns a Client that binds to current transaction.
func (tx *Tx) Client() *Client {
	return &Client{
		config:                      tx.config,
		Schema:                      migrate.NewSchema(tx.driver),
		ActionsRule:                 NewActionsRuleClient(tx.config),
		CheckListCategory:           NewCheckListCategoryClient(tx.config),
		CheckListItem:               NewCheckListItemClient(tx.config),
		CheckListItemDefinition:     NewCheckListItemDefinitionClient(tx.config),
		Comment:                     NewCommentClient(tx.config),
		Customer:                    NewCustomerClient(tx.config),
		Equipment:                   NewEquipmentClient(tx.config),
		EquipmentCategory:           NewEquipmentCategoryClient(tx.config),
		EquipmentPort:               NewEquipmentPortClient(tx.config),
		EquipmentPortDefinition:     NewEquipmentPortDefinitionClient(tx.config),
		EquipmentPortType:           NewEquipmentPortTypeClient(tx.config),
		EquipmentPosition:           NewEquipmentPositionClient(tx.config),
		EquipmentPositionDefinition: NewEquipmentPositionDefinitionClient(tx.config),
		EquipmentType:               NewEquipmentTypeClient(tx.config),
		File:                        NewFileClient(tx.config),
		FloorPlan:                   NewFloorPlanClient(tx.config),
		FloorPlanReferencePoint:     NewFloorPlanReferencePointClient(tx.config),
		FloorPlanScale:              NewFloorPlanScaleClient(tx.config),
		Hyperlink:                   NewHyperlinkClient(tx.config),
		Link:                        NewLinkClient(tx.config),
		Location:                    NewLocationClient(tx.config),
		LocationType:                NewLocationTypeClient(tx.config),
		Project:                     NewProjectClient(tx.config),
		ProjectType:                 NewProjectTypeClient(tx.config),
		Property:                    NewPropertyClient(tx.config),
		PropertyType:                NewPropertyTypeClient(tx.config),
		Service:                     NewServiceClient(tx.config),
		ServiceEndpoint:             NewServiceEndpointClient(tx.config),
		ServiceType:                 NewServiceTypeClient(tx.config),
		Survey:                      NewSurveyClient(tx.config),
		SurveyCellScan:              NewSurveyCellScanClient(tx.config),
		SurveyQuestion:              NewSurveyQuestionClient(tx.config),
		SurveyTemplateCategory:      NewSurveyTemplateCategoryClient(tx.config),
		SurveyTemplateQuestion:      NewSurveyTemplateQuestionClient(tx.config),
		SurveyWiFiScan:              NewSurveyWiFiScanClient(tx.config),
		Technician:                  NewTechnicianClient(tx.config),
		User:                        NewUserClient(tx.config),
		WorkOrder:                   NewWorkOrderClient(tx.config),
		WorkOrderDefinition:         NewWorkOrderDefinitionClient(tx.config),
		WorkOrderType:               NewWorkOrderTypeClient(tx.config),
	}
}

// txDriver wraps the given dialect.Tx with a nop dialect.Driver implementation.
// The idea is to support transactions without adding any extra code to the builders.
// When a builder calls to driver.Tx(), it gets the same dialect.Tx instance.
// Commit and Rollback are nop for the internal builders and the user must call one
// of them in order to commit or rollback the transaction.
//
// If a closed transaction is embedded in one of the generated entities, and the entity
// applies a query, for example: ActionsRule.QueryXXX(), the query will be executed
// through the driver which created this transaction.
//
// Note that txDriver is not goroutine safe.
type txDriver struct {
	// the driver we started the transaction from.
	drv dialect.Driver
	// tx is the underlying transaction.
	tx dialect.Tx
}

// newTx creates a new transactional driver.
func newTx(ctx context.Context, drv dialect.Driver) (*txDriver, error) {
	tx, err := drv.Tx(ctx)
	if err != nil {
		return nil, err
	}
	return &txDriver{tx: tx, drv: drv}, nil
}

// Tx returns the transaction wrapper (txDriver) to avoid Commit or Rollback calls
// from the internal builders. Should be called only by the internal builders.
func (tx *txDriver) Tx(context.Context) (dialect.Tx, error) { return tx, nil }

// Dialect returns the dialect of the driver we started the transaction from.
func (tx *txDriver) Dialect() string { return tx.drv.Dialect() }

// Close is a nop close.
func (*txDriver) Close() error { return nil }

// Commit is a nop commit for the internal builders.
// User must call `Tx.Commit` in order to commit the transaction.
func (*txDriver) Commit() error { return nil }

// Rollback is a nop rollback for the internal builders.
// User must call `Tx.Rollback` in order to rollback the transaction.
func (*txDriver) Rollback() error { return nil }

// Exec calls tx.Exec.
func (tx *txDriver) Exec(ctx context.Context, query string, args, v interface{}) error {
	return tx.tx.Exec(ctx, query, args, v)
}

// Query calls tx.Query.
func (tx *txDriver) Query(ctx context.Context, query string, args, v interface{}) error {
	return tx.tx.Query(ctx, query, args, v)
}

var _ dialect.Driver = (*txDriver)(nil)
