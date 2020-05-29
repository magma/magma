// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"sync"

	"github.com/facebookincubator/ent/dialect"
)

// Tx is a transactional client that is created by calling Client.Tx().
type Tx struct {
	config
	// ActionsRule is the client for interacting with the ActionsRule builders.
	ActionsRule *ActionsRuleClient
	// Activity is the client for interacting with the Activity builders.
	Activity *ActivityClient
	// CheckListCategory is the client for interacting with the CheckListCategory builders.
	CheckListCategory *CheckListCategoryClient
	// CheckListCategoryDefinition is the client for interacting with the CheckListCategoryDefinition builders.
	CheckListCategoryDefinition *CheckListCategoryDefinitionClient
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
	// PermissionsPolicy is the client for interacting with the PermissionsPolicy builders.
	PermissionsPolicy *PermissionsPolicyClient
	// Project is the client for interacting with the Project builders.
	Project *ProjectClient
	// ProjectType is the client for interacting with the ProjectType builders.
	ProjectType *ProjectTypeClient
	// Property is the client for interacting with the Property builders.
	Property *PropertyClient
	// PropertyType is the client for interacting with the PropertyType builders.
	PropertyType *PropertyTypeClient
	// ReportFilter is the client for interacting with the ReportFilter builders.
	ReportFilter *ReportFilterClient
	// Service is the client for interacting with the Service builders.
	Service *ServiceClient
	// ServiceEndpoint is the client for interacting with the ServiceEndpoint builders.
	ServiceEndpoint *ServiceEndpointClient
	// ServiceEndpointDefinition is the client for interacting with the ServiceEndpointDefinition builders.
	ServiceEndpointDefinition *ServiceEndpointDefinitionClient
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
	// User is the client for interacting with the User builders.
	User *UserClient
	// UsersGroup is the client for interacting with the UsersGroup builders.
	UsersGroup *UsersGroupClient
	// WorkOrder is the client for interacting with the WorkOrder builders.
	WorkOrder *WorkOrderClient
	// WorkOrderDefinition is the client for interacting with the WorkOrderDefinition builders.
	WorkOrderDefinition *WorkOrderDefinitionClient
	// WorkOrderType is the client for interacting with the WorkOrderType builders.
	WorkOrderType *WorkOrderTypeClient

	// lazily loaded.
	client     *Client
	clientOnce sync.Once

	// completion callbacks.
	mu         sync.Mutex
	onCommit   []func(error)
	onRollback []func(error)
}

// Commit commits the transaction.
func (tx *Tx) Commit() error {
	err := tx.config.driver.(*txDriver).tx.Commit()
	tx.mu.Lock()
	defer tx.mu.Unlock()
	for _, f := range tx.onCommit {
		f(err)
	}
	return err
}

// OnCommit adds a function to call on commit.
func (tx *Tx) OnCommit(f func(error)) {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	tx.onCommit = append(tx.onCommit, f)
}

// Rollback rollbacks the transaction.
func (tx *Tx) Rollback() error {
	err := tx.config.driver.(*txDriver).tx.Rollback()
	tx.mu.Lock()
	defer tx.mu.Unlock()
	for _, f := range tx.onRollback {
		f(err)
	}
	return err
}

// OnRollback adds a function to call on rollback.
func (tx *Tx) OnRollback(f func(error)) {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	tx.onRollback = append(tx.onRollback, f)
}

// Client returns a Client that binds to current transaction.
func (tx *Tx) Client() *Client {
	tx.clientOnce.Do(func() {
		tx.client = &Client{config: tx.config}
		tx.client.init()
	})
	return tx.client
}

func (tx *Tx) init() {
	tx.ActionsRule = NewActionsRuleClient(tx.config)
	tx.Activity = NewActivityClient(tx.config)
	tx.CheckListCategory = NewCheckListCategoryClient(tx.config)
	tx.CheckListCategoryDefinition = NewCheckListCategoryDefinitionClient(tx.config)
	tx.CheckListItem = NewCheckListItemClient(tx.config)
	tx.CheckListItemDefinition = NewCheckListItemDefinitionClient(tx.config)
	tx.Comment = NewCommentClient(tx.config)
	tx.Customer = NewCustomerClient(tx.config)
	tx.Equipment = NewEquipmentClient(tx.config)
	tx.EquipmentCategory = NewEquipmentCategoryClient(tx.config)
	tx.EquipmentPort = NewEquipmentPortClient(tx.config)
	tx.EquipmentPortDefinition = NewEquipmentPortDefinitionClient(tx.config)
	tx.EquipmentPortType = NewEquipmentPortTypeClient(tx.config)
	tx.EquipmentPosition = NewEquipmentPositionClient(tx.config)
	tx.EquipmentPositionDefinition = NewEquipmentPositionDefinitionClient(tx.config)
	tx.EquipmentType = NewEquipmentTypeClient(tx.config)
	tx.File = NewFileClient(tx.config)
	tx.FloorPlan = NewFloorPlanClient(tx.config)
	tx.FloorPlanReferencePoint = NewFloorPlanReferencePointClient(tx.config)
	tx.FloorPlanScale = NewFloorPlanScaleClient(tx.config)
	tx.Hyperlink = NewHyperlinkClient(tx.config)
	tx.Link = NewLinkClient(tx.config)
	tx.Location = NewLocationClient(tx.config)
	tx.LocationType = NewLocationTypeClient(tx.config)
	tx.PermissionsPolicy = NewPermissionsPolicyClient(tx.config)
	tx.Project = NewProjectClient(tx.config)
	tx.ProjectType = NewProjectTypeClient(tx.config)
	tx.Property = NewPropertyClient(tx.config)
	tx.PropertyType = NewPropertyTypeClient(tx.config)
	tx.ReportFilter = NewReportFilterClient(tx.config)
	tx.Service = NewServiceClient(tx.config)
	tx.ServiceEndpoint = NewServiceEndpointClient(tx.config)
	tx.ServiceEndpointDefinition = NewServiceEndpointDefinitionClient(tx.config)
	tx.ServiceType = NewServiceTypeClient(tx.config)
	tx.Survey = NewSurveyClient(tx.config)
	tx.SurveyCellScan = NewSurveyCellScanClient(tx.config)
	tx.SurveyQuestion = NewSurveyQuestionClient(tx.config)
	tx.SurveyTemplateCategory = NewSurveyTemplateCategoryClient(tx.config)
	tx.SurveyTemplateQuestion = NewSurveyTemplateQuestionClient(tx.config)
	tx.SurveyWiFiScan = NewSurveyWiFiScanClient(tx.config)
	tx.User = NewUserClient(tx.config)
	tx.UsersGroup = NewUsersGroupClient(tx.config)
	tx.WorkOrder = NewWorkOrderClient(tx.config)
	tx.WorkOrderDefinition = NewWorkOrderDefinitionClient(tx.config)
	tx.WorkOrderType = NewWorkOrderTypeClient(tx.config)
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
