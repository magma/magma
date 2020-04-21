// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"log"

	"github.com/facebookincubator/symphony/graph/ent/migrate"

	"github.com/facebookincubator/symphony/graph/ent/actionsrule"
	"github.com/facebookincubator/symphony/graph/ent/checklistcategory"
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/checklistitemdefinition"
	"github.com/facebookincubator/symphony/graph/ent/comment"
	"github.com/facebookincubator/symphony/graph/ent/customer"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentcategory"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/floorplan"
	"github.com/facebookincubator/symphony/graph/ent/floorplanreferencepoint"
	"github.com/facebookincubator/symphony/graph/ent/floorplanscale"
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/permissionspolicy"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/reportfilter"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpointdefinition"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatecategory"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatequestion"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
	"github.com/facebookincubator/symphony/graph/ent/technician"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/usersgroup"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
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
	// Technician is the client for interacting with the Technician builders.
	Technician *TechnicianClient
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

	// additional fields for node api
	tables tables
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	cfg := config{log: log.Println, hooks: &hooks{}}
	cfg.options(opts...)
	client := &Client{config: cfg}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.ActionsRule = NewActionsRuleClient(c.config)
	c.CheckListCategory = NewCheckListCategoryClient(c.config)
	c.CheckListItem = NewCheckListItemClient(c.config)
	c.CheckListItemDefinition = NewCheckListItemDefinitionClient(c.config)
	c.Comment = NewCommentClient(c.config)
	c.Customer = NewCustomerClient(c.config)
	c.Equipment = NewEquipmentClient(c.config)
	c.EquipmentCategory = NewEquipmentCategoryClient(c.config)
	c.EquipmentPort = NewEquipmentPortClient(c.config)
	c.EquipmentPortDefinition = NewEquipmentPortDefinitionClient(c.config)
	c.EquipmentPortType = NewEquipmentPortTypeClient(c.config)
	c.EquipmentPosition = NewEquipmentPositionClient(c.config)
	c.EquipmentPositionDefinition = NewEquipmentPositionDefinitionClient(c.config)
	c.EquipmentType = NewEquipmentTypeClient(c.config)
	c.File = NewFileClient(c.config)
	c.FloorPlan = NewFloorPlanClient(c.config)
	c.FloorPlanReferencePoint = NewFloorPlanReferencePointClient(c.config)
	c.FloorPlanScale = NewFloorPlanScaleClient(c.config)
	c.Hyperlink = NewHyperlinkClient(c.config)
	c.Link = NewLinkClient(c.config)
	c.Location = NewLocationClient(c.config)
	c.LocationType = NewLocationTypeClient(c.config)
	c.PermissionsPolicy = NewPermissionsPolicyClient(c.config)
	c.Project = NewProjectClient(c.config)
	c.ProjectType = NewProjectTypeClient(c.config)
	c.Property = NewPropertyClient(c.config)
	c.PropertyType = NewPropertyTypeClient(c.config)
	c.ReportFilter = NewReportFilterClient(c.config)
	c.Service = NewServiceClient(c.config)
	c.ServiceEndpoint = NewServiceEndpointClient(c.config)
	c.ServiceEndpointDefinition = NewServiceEndpointDefinitionClient(c.config)
	c.ServiceType = NewServiceTypeClient(c.config)
	c.Survey = NewSurveyClient(c.config)
	c.SurveyCellScan = NewSurveyCellScanClient(c.config)
	c.SurveyQuestion = NewSurveyQuestionClient(c.config)
	c.SurveyTemplateCategory = NewSurveyTemplateCategoryClient(c.config)
	c.SurveyTemplateQuestion = NewSurveyTemplateQuestionClient(c.config)
	c.SurveyWiFiScan = NewSurveyWiFiScanClient(c.config)
	c.Technician = NewTechnicianClient(c.config)
	c.User = NewUserClient(c.config)
	c.UsersGroup = NewUsersGroupClient(c.config)
	c.WorkOrder = NewWorkOrderClient(c.config)
	c.WorkOrderDefinition = NewWorkOrderDefinitionClient(c.config)
	c.WorkOrderType = NewWorkOrderTypeClient(c.config)
}

// Open opens a connection to the database specified by the driver name and a
// driver-specific data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// Tx returns a new transactional client.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %v", err)
	}
	cfg := config{driver: tx, log: c.log, debug: c.debug, hooks: c.hooks}
	return &Tx{
		config:                      cfg,
		ActionsRule:                 NewActionsRuleClient(cfg),
		CheckListCategory:           NewCheckListCategoryClient(cfg),
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
		PermissionsPolicy:           NewPermissionsPolicyClient(cfg),
		Project:                     NewProjectClient(cfg),
		ProjectType:                 NewProjectTypeClient(cfg),
		Property:                    NewPropertyClient(cfg),
		PropertyType:                NewPropertyTypeClient(cfg),
		ReportFilter:                NewReportFilterClient(cfg),
		Service:                     NewServiceClient(cfg),
		ServiceEndpoint:             NewServiceEndpointClient(cfg),
		ServiceEndpointDefinition:   NewServiceEndpointDefinitionClient(cfg),
		ServiceType:                 NewServiceTypeClient(cfg),
		Survey:                      NewSurveyClient(cfg),
		SurveyCellScan:              NewSurveyCellScanClient(cfg),
		SurveyQuestion:              NewSurveyQuestionClient(cfg),
		SurveyTemplateCategory:      NewSurveyTemplateCategoryClient(cfg),
		SurveyTemplateQuestion:      NewSurveyTemplateQuestionClient(cfg),
		SurveyWiFiScan:              NewSurveyWiFiScanClient(cfg),
		Technician:                  NewTechnicianClient(cfg),
		User:                        NewUserClient(cfg),
		UsersGroup:                  NewUsersGroupClient(cfg),
		WorkOrder:                   NewWorkOrderClient(cfg),
		WorkOrderDefinition:         NewWorkOrderDefinitionClient(cfg),
		WorkOrderType:               NewWorkOrderTypeClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(*sql.Driver).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %v", err)
	}
	cfg := config{driver: &txDriver{tx: tx, drv: c.driver}, log: c.log, debug: c.debug, hooks: c.hooks}
	return &Tx{
		config:                      cfg,
		ActionsRule:                 NewActionsRuleClient(cfg),
		CheckListCategory:           NewCheckListCategoryClient(cfg),
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
		PermissionsPolicy:           NewPermissionsPolicyClient(cfg),
		Project:                     NewProjectClient(cfg),
		ProjectType:                 NewProjectTypeClient(cfg),
		Property:                    NewPropertyClient(cfg),
		PropertyType:                NewPropertyTypeClient(cfg),
		ReportFilter:                NewReportFilterClient(cfg),
		Service:                     NewServiceClient(cfg),
		ServiceEndpoint:             NewServiceEndpointClient(cfg),
		ServiceEndpointDefinition:   NewServiceEndpointDefinitionClient(cfg),
		ServiceType:                 NewServiceTypeClient(cfg),
		Survey:                      NewSurveyClient(cfg),
		SurveyCellScan:              NewSurveyCellScanClient(cfg),
		SurveyQuestion:              NewSurveyQuestionClient(cfg),
		SurveyTemplateCategory:      NewSurveyTemplateCategoryClient(cfg),
		SurveyTemplateQuestion:      NewSurveyTemplateQuestionClient(cfg),
		SurveyWiFiScan:              NewSurveyWiFiScanClient(cfg),
		Technician:                  NewTechnicianClient(cfg),
		User:                        NewUserClient(cfg),
		UsersGroup:                  NewUsersGroupClient(cfg),
		WorkOrder:                   NewWorkOrderClient(cfg),
		WorkOrderDefinition:         NewWorkOrderDefinitionClient(cfg),
		WorkOrderType:               NewWorkOrderTypeClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		ActionsRule.
//		Query().
//		Count(ctx)
//
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := config{driver: dialect.Debug(c.driver, c.log), log: c.log, debug: true, hooks: c.hooks}
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.ActionsRule.Use(hooks...)
	c.CheckListCategory.Use(hooks...)
	c.CheckListItem.Use(hooks...)
	c.CheckListItemDefinition.Use(hooks...)
	c.Comment.Use(hooks...)
	c.Customer.Use(hooks...)
	c.Equipment.Use(hooks...)
	c.EquipmentCategory.Use(hooks...)
	c.EquipmentPort.Use(hooks...)
	c.EquipmentPortDefinition.Use(hooks...)
	c.EquipmentPortType.Use(hooks...)
	c.EquipmentPosition.Use(hooks...)
	c.EquipmentPositionDefinition.Use(hooks...)
	c.EquipmentType.Use(hooks...)
	c.File.Use(hooks...)
	c.FloorPlan.Use(hooks...)
	c.FloorPlanReferencePoint.Use(hooks...)
	c.FloorPlanScale.Use(hooks...)
	c.Hyperlink.Use(hooks...)
	c.Link.Use(hooks...)
	c.Location.Use(hooks...)
	c.LocationType.Use(hooks...)
	c.PermissionsPolicy.Use(hooks...)
	c.Project.Use(hooks...)
	c.ProjectType.Use(hooks...)
	c.Property.Use(hooks...)
	c.PropertyType.Use(hooks...)
	c.ReportFilter.Use(hooks...)
	c.Service.Use(hooks...)
	c.ServiceEndpoint.Use(hooks...)
	c.ServiceEndpointDefinition.Use(hooks...)
	c.ServiceType.Use(hooks...)
	c.Survey.Use(hooks...)
	c.SurveyCellScan.Use(hooks...)
	c.SurveyQuestion.Use(hooks...)
	c.SurveyTemplateCategory.Use(hooks...)
	c.SurveyTemplateQuestion.Use(hooks...)
	c.SurveyWiFiScan.Use(hooks...)
	c.Technician.Use(hooks...)
	c.User.Use(hooks...)
	c.UsersGroup.Use(hooks...)
	c.WorkOrder.Use(hooks...)
	c.WorkOrderDefinition.Use(hooks...)
	c.WorkOrderType.Use(hooks...)
}

// ActionsRuleClient is a client for the ActionsRule schema.
type ActionsRuleClient struct {
	config
}

// NewActionsRuleClient returns a client for the ActionsRule from the given config.
func NewActionsRuleClient(c config) *ActionsRuleClient {
	return &ActionsRuleClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `actionsrule.Hooks(f(g(h())))`.
func (c *ActionsRuleClient) Use(hooks ...Hook) {
	c.hooks.ActionsRule = append(c.hooks.ActionsRule, hooks...)
}

// Create returns a create builder for ActionsRule.
func (c *ActionsRuleClient) Create() *ActionsRuleCreate {
	mutation := newActionsRuleMutation(c.config, OpCreate)
	return &ActionsRuleCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for ActionsRule.
func (c *ActionsRuleClient) Update() *ActionsRuleUpdate {
	mutation := newActionsRuleMutation(c.config, OpUpdate)
	return &ActionsRuleUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ActionsRuleClient) UpdateOne(ar *ActionsRule) *ActionsRuleUpdateOne {
	return c.UpdateOneID(ar.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *ActionsRuleClient) UpdateOneID(id int) *ActionsRuleUpdateOne {
	mutation := newActionsRuleMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &ActionsRuleUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for ActionsRule.
func (c *ActionsRuleClient) Delete() *ActionsRuleDelete {
	mutation := newActionsRuleMutation(c.config, OpDelete)
	return &ActionsRuleDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *ActionsRuleClient) DeleteOne(ar *ActionsRule) *ActionsRuleDeleteOne {
	return c.DeleteOneID(ar.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *ActionsRuleClient) DeleteOneID(id int) *ActionsRuleDeleteOne {
	builder := c.Delete().Where(actionsrule.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ActionsRuleDeleteOne{builder}
}

// Create returns a query builder for ActionsRule.
func (c *ActionsRuleClient) Query() *ActionsRuleQuery {
	return &ActionsRuleQuery{config: c.config}
}

// Get returns a ActionsRule entity by its id.
func (c *ActionsRuleClient) Get(ctx context.Context, id int) (*ActionsRule, error) {
	return c.Query().Where(actionsrule.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ActionsRuleClient) GetX(ctx context.Context, id int) *ActionsRule {
	ar, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return ar
}

// Hooks returns the client hooks.
func (c *ActionsRuleClient) Hooks() []Hook {
	return c.hooks.ActionsRule
}

// CheckListCategoryClient is a client for the CheckListCategory schema.
type CheckListCategoryClient struct {
	config
}

// NewCheckListCategoryClient returns a client for the CheckListCategory from the given config.
func NewCheckListCategoryClient(c config) *CheckListCategoryClient {
	return &CheckListCategoryClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `checklistcategory.Hooks(f(g(h())))`.
func (c *CheckListCategoryClient) Use(hooks ...Hook) {
	c.hooks.CheckListCategory = append(c.hooks.CheckListCategory, hooks...)
}

// Create returns a create builder for CheckListCategory.
func (c *CheckListCategoryClient) Create() *CheckListCategoryCreate {
	mutation := newCheckListCategoryMutation(c.config, OpCreate)
	return &CheckListCategoryCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for CheckListCategory.
func (c *CheckListCategoryClient) Update() *CheckListCategoryUpdate {
	mutation := newCheckListCategoryMutation(c.config, OpUpdate)
	return &CheckListCategoryUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *CheckListCategoryClient) UpdateOne(clc *CheckListCategory) *CheckListCategoryUpdateOne {
	return c.UpdateOneID(clc.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *CheckListCategoryClient) UpdateOneID(id int) *CheckListCategoryUpdateOne {
	mutation := newCheckListCategoryMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &CheckListCategoryUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for CheckListCategory.
func (c *CheckListCategoryClient) Delete() *CheckListCategoryDelete {
	mutation := newCheckListCategoryMutation(c.config, OpDelete)
	return &CheckListCategoryDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *CheckListCategoryClient) DeleteOne(clc *CheckListCategory) *CheckListCategoryDeleteOne {
	return c.DeleteOneID(clc.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *CheckListCategoryClient) DeleteOneID(id int) *CheckListCategoryDeleteOne {
	builder := c.Delete().Where(checklistcategory.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &CheckListCategoryDeleteOne{builder}
}

// Create returns a query builder for CheckListCategory.
func (c *CheckListCategoryClient) Query() *CheckListCategoryQuery {
	return &CheckListCategoryQuery{config: c.config}
}

// Get returns a CheckListCategory entity by its id.
func (c *CheckListCategoryClient) Get(ctx context.Context, id int) (*CheckListCategory, error) {
	return c.Query().Where(checklistcategory.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *CheckListCategoryClient) GetX(ctx context.Context, id int) *CheckListCategory {
	clc, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return clc
}

// QueryCheckListItems queries the check_list_items edge of a CheckListCategory.
func (c *CheckListCategoryClient) QueryCheckListItems(clc *CheckListCategory) *CheckListItemQuery {
	query := &CheckListItemQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := clc.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(checklistcategory.Table, checklistcategory.FieldID, id),
			sqlgraph.To(checklistitem.Table, checklistitem.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, checklistcategory.CheckListItemsTable, checklistcategory.CheckListItemsColumn),
		)
		fromV = sqlgraph.Neighbors(clc.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *CheckListCategoryClient) Hooks() []Hook {
	return c.hooks.CheckListCategory
}

// CheckListItemClient is a client for the CheckListItem schema.
type CheckListItemClient struct {
	config
}

// NewCheckListItemClient returns a client for the CheckListItem from the given config.
func NewCheckListItemClient(c config) *CheckListItemClient {
	return &CheckListItemClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `checklistitem.Hooks(f(g(h())))`.
func (c *CheckListItemClient) Use(hooks ...Hook) {
	c.hooks.CheckListItem = append(c.hooks.CheckListItem, hooks...)
}

// Create returns a create builder for CheckListItem.
func (c *CheckListItemClient) Create() *CheckListItemCreate {
	mutation := newCheckListItemMutation(c.config, OpCreate)
	return &CheckListItemCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for CheckListItem.
func (c *CheckListItemClient) Update() *CheckListItemUpdate {
	mutation := newCheckListItemMutation(c.config, OpUpdate)
	return &CheckListItemUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *CheckListItemClient) UpdateOne(cli *CheckListItem) *CheckListItemUpdateOne {
	return c.UpdateOneID(cli.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *CheckListItemClient) UpdateOneID(id int) *CheckListItemUpdateOne {
	mutation := newCheckListItemMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &CheckListItemUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for CheckListItem.
func (c *CheckListItemClient) Delete() *CheckListItemDelete {
	mutation := newCheckListItemMutation(c.config, OpDelete)
	return &CheckListItemDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *CheckListItemClient) DeleteOne(cli *CheckListItem) *CheckListItemDeleteOne {
	return c.DeleteOneID(cli.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *CheckListItemClient) DeleteOneID(id int) *CheckListItemDeleteOne {
	builder := c.Delete().Where(checklistitem.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &CheckListItemDeleteOne{builder}
}

// Create returns a query builder for CheckListItem.
func (c *CheckListItemClient) Query() *CheckListItemQuery {
	return &CheckListItemQuery{config: c.config}
}

// Get returns a CheckListItem entity by its id.
func (c *CheckListItemClient) Get(ctx context.Context, id int) (*CheckListItem, error) {
	return c.Query().Where(checklistitem.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *CheckListItemClient) GetX(ctx context.Context, id int) *CheckListItem {
	cli, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return cli
}

// QueryFiles queries the files edge of a CheckListItem.
func (c *CheckListItemClient) QueryFiles(cli *CheckListItem) *FileQuery {
	query := &FileQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := cli.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(checklistitem.Table, checklistitem.FieldID, id),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, checklistitem.FilesTable, checklistitem.FilesColumn),
		)
		fromV = sqlgraph.Neighbors(cli.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryWifiScan queries the wifi_scan edge of a CheckListItem.
func (c *CheckListItemClient) QueryWifiScan(cli *CheckListItem) *SurveyWiFiScanQuery {
	query := &SurveyWiFiScanQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := cli.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(checklistitem.Table, checklistitem.FieldID, id),
			sqlgraph.To(surveywifiscan.Table, surveywifiscan.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, checklistitem.WifiScanTable, checklistitem.WifiScanColumn),
		)
		fromV = sqlgraph.Neighbors(cli.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryCellScan queries the cell_scan edge of a CheckListItem.
func (c *CheckListItemClient) QueryCellScan(cli *CheckListItem) *SurveyCellScanQuery {
	query := &SurveyCellScanQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := cli.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(checklistitem.Table, checklistitem.FieldID, id),
			sqlgraph.To(surveycellscan.Table, surveycellscan.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, checklistitem.CellScanTable, checklistitem.CellScanColumn),
		)
		fromV = sqlgraph.Neighbors(cli.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryWorkOrder queries the work_order edge of a CheckListItem.
func (c *CheckListItemClient) QueryWorkOrder(cli *CheckListItem) *WorkOrderQuery {
	query := &WorkOrderQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := cli.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(checklistitem.Table, checklistitem.FieldID, id),
			sqlgraph.To(workorder.Table, workorder.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, checklistitem.WorkOrderTable, checklistitem.WorkOrderColumn),
		)
		fromV = sqlgraph.Neighbors(cli.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *CheckListItemClient) Hooks() []Hook {
	return c.hooks.CheckListItem
}

// CheckListItemDefinitionClient is a client for the CheckListItemDefinition schema.
type CheckListItemDefinitionClient struct {
	config
}

// NewCheckListItemDefinitionClient returns a client for the CheckListItemDefinition from the given config.
func NewCheckListItemDefinitionClient(c config) *CheckListItemDefinitionClient {
	return &CheckListItemDefinitionClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `checklistitemdefinition.Hooks(f(g(h())))`.
func (c *CheckListItemDefinitionClient) Use(hooks ...Hook) {
	c.hooks.CheckListItemDefinition = append(c.hooks.CheckListItemDefinition, hooks...)
}

// Create returns a create builder for CheckListItemDefinition.
func (c *CheckListItemDefinitionClient) Create() *CheckListItemDefinitionCreate {
	mutation := newCheckListItemDefinitionMutation(c.config, OpCreate)
	return &CheckListItemDefinitionCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for CheckListItemDefinition.
func (c *CheckListItemDefinitionClient) Update() *CheckListItemDefinitionUpdate {
	mutation := newCheckListItemDefinitionMutation(c.config, OpUpdate)
	return &CheckListItemDefinitionUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *CheckListItemDefinitionClient) UpdateOne(clid *CheckListItemDefinition) *CheckListItemDefinitionUpdateOne {
	return c.UpdateOneID(clid.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *CheckListItemDefinitionClient) UpdateOneID(id int) *CheckListItemDefinitionUpdateOne {
	mutation := newCheckListItemDefinitionMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &CheckListItemDefinitionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for CheckListItemDefinition.
func (c *CheckListItemDefinitionClient) Delete() *CheckListItemDefinitionDelete {
	mutation := newCheckListItemDefinitionMutation(c.config, OpDelete)
	return &CheckListItemDefinitionDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *CheckListItemDefinitionClient) DeleteOne(clid *CheckListItemDefinition) *CheckListItemDefinitionDeleteOne {
	return c.DeleteOneID(clid.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *CheckListItemDefinitionClient) DeleteOneID(id int) *CheckListItemDefinitionDeleteOne {
	builder := c.Delete().Where(checklistitemdefinition.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &CheckListItemDefinitionDeleteOne{builder}
}

// Create returns a query builder for CheckListItemDefinition.
func (c *CheckListItemDefinitionClient) Query() *CheckListItemDefinitionQuery {
	return &CheckListItemDefinitionQuery{config: c.config}
}

// Get returns a CheckListItemDefinition entity by its id.
func (c *CheckListItemDefinitionClient) Get(ctx context.Context, id int) (*CheckListItemDefinition, error) {
	return c.Query().Where(checklistitemdefinition.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *CheckListItemDefinitionClient) GetX(ctx context.Context, id int) *CheckListItemDefinition {
	clid, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return clid
}

// QueryWorkOrderType queries the work_order_type edge of a CheckListItemDefinition.
func (c *CheckListItemDefinitionClient) QueryWorkOrderType(clid *CheckListItemDefinition) *WorkOrderTypeQuery {
	query := &WorkOrderTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := clid.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(checklistitemdefinition.Table, checklistitemdefinition.FieldID, id),
			sqlgraph.To(workordertype.Table, workordertype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, checklistitemdefinition.WorkOrderTypeTable, checklistitemdefinition.WorkOrderTypeColumn),
		)
		fromV = sqlgraph.Neighbors(clid.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *CheckListItemDefinitionClient) Hooks() []Hook {
	return c.hooks.CheckListItemDefinition
}

// CommentClient is a client for the Comment schema.
type CommentClient struct {
	config
}

// NewCommentClient returns a client for the Comment from the given config.
func NewCommentClient(c config) *CommentClient {
	return &CommentClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `comment.Hooks(f(g(h())))`.
func (c *CommentClient) Use(hooks ...Hook) {
	c.hooks.Comment = append(c.hooks.Comment, hooks...)
}

// Create returns a create builder for Comment.
func (c *CommentClient) Create() *CommentCreate {
	mutation := newCommentMutation(c.config, OpCreate)
	return &CommentCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for Comment.
func (c *CommentClient) Update() *CommentUpdate {
	mutation := newCommentMutation(c.config, OpUpdate)
	return &CommentUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *CommentClient) UpdateOne(co *Comment) *CommentUpdateOne {
	return c.UpdateOneID(co.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *CommentClient) UpdateOneID(id int) *CommentUpdateOne {
	mutation := newCommentMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &CommentUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Comment.
func (c *CommentClient) Delete() *CommentDelete {
	mutation := newCommentMutation(c.config, OpDelete)
	return &CommentDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *CommentClient) DeleteOne(co *Comment) *CommentDeleteOne {
	return c.DeleteOneID(co.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *CommentClient) DeleteOneID(id int) *CommentDeleteOne {
	builder := c.Delete().Where(comment.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &CommentDeleteOne{builder}
}

// Create returns a query builder for Comment.
func (c *CommentClient) Query() *CommentQuery {
	return &CommentQuery{config: c.config}
}

// Get returns a Comment entity by its id.
func (c *CommentClient) Get(ctx context.Context, id int) (*Comment, error) {
	return c.Query().Where(comment.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *CommentClient) GetX(ctx context.Context, id int) *Comment {
	co, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return co
}

// QueryAuthor queries the author edge of a Comment.
func (c *CommentClient) QueryAuthor(co *Comment) *UserQuery {
	query := &UserQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := co.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(comment.Table, comment.FieldID, id),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, comment.AuthorTable, comment.AuthorColumn),
		)
		fromV = sqlgraph.Neighbors(co.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *CommentClient) Hooks() []Hook {
	return c.hooks.Comment
}

// CustomerClient is a client for the Customer schema.
type CustomerClient struct {
	config
}

// NewCustomerClient returns a client for the Customer from the given config.
func NewCustomerClient(c config) *CustomerClient {
	return &CustomerClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `customer.Hooks(f(g(h())))`.
func (c *CustomerClient) Use(hooks ...Hook) {
	c.hooks.Customer = append(c.hooks.Customer, hooks...)
}

// Create returns a create builder for Customer.
func (c *CustomerClient) Create() *CustomerCreate {
	mutation := newCustomerMutation(c.config, OpCreate)
	return &CustomerCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for Customer.
func (c *CustomerClient) Update() *CustomerUpdate {
	mutation := newCustomerMutation(c.config, OpUpdate)
	return &CustomerUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *CustomerClient) UpdateOne(cu *Customer) *CustomerUpdateOne {
	return c.UpdateOneID(cu.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *CustomerClient) UpdateOneID(id int) *CustomerUpdateOne {
	mutation := newCustomerMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &CustomerUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Customer.
func (c *CustomerClient) Delete() *CustomerDelete {
	mutation := newCustomerMutation(c.config, OpDelete)
	return &CustomerDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *CustomerClient) DeleteOne(cu *Customer) *CustomerDeleteOne {
	return c.DeleteOneID(cu.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *CustomerClient) DeleteOneID(id int) *CustomerDeleteOne {
	builder := c.Delete().Where(customer.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &CustomerDeleteOne{builder}
}

// Create returns a query builder for Customer.
func (c *CustomerClient) Query() *CustomerQuery {
	return &CustomerQuery{config: c.config}
}

// Get returns a Customer entity by its id.
func (c *CustomerClient) Get(ctx context.Context, id int) (*Customer, error) {
	return c.Query().Where(customer.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *CustomerClient) GetX(ctx context.Context, id int) *Customer {
	cu, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return cu
}

// QueryServices queries the services edge of a Customer.
func (c *CustomerClient) QueryServices(cu *Customer) *ServiceQuery {
	query := &ServiceQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := cu.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(customer.Table, customer.FieldID, id),
			sqlgraph.To(service.Table, service.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, customer.ServicesTable, customer.ServicesPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(cu.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *CustomerClient) Hooks() []Hook {
	return c.hooks.Customer
}

// EquipmentClient is a client for the Equipment schema.
type EquipmentClient struct {
	config
}

// NewEquipmentClient returns a client for the Equipment from the given config.
func NewEquipmentClient(c config) *EquipmentClient {
	return &EquipmentClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `equipment.Hooks(f(g(h())))`.
func (c *EquipmentClient) Use(hooks ...Hook) {
	c.hooks.Equipment = append(c.hooks.Equipment, hooks...)
}

// Create returns a create builder for Equipment.
func (c *EquipmentClient) Create() *EquipmentCreate {
	mutation := newEquipmentMutation(c.config, OpCreate)
	return &EquipmentCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for Equipment.
func (c *EquipmentClient) Update() *EquipmentUpdate {
	mutation := newEquipmentMutation(c.config, OpUpdate)
	return &EquipmentUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *EquipmentClient) UpdateOne(e *Equipment) *EquipmentUpdateOne {
	return c.UpdateOneID(e.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *EquipmentClient) UpdateOneID(id int) *EquipmentUpdateOne {
	mutation := newEquipmentMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &EquipmentUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Equipment.
func (c *EquipmentClient) Delete() *EquipmentDelete {
	mutation := newEquipmentMutation(c.config, OpDelete)
	return &EquipmentDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *EquipmentClient) DeleteOne(e *Equipment) *EquipmentDeleteOne {
	return c.DeleteOneID(e.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *EquipmentClient) DeleteOneID(id int) *EquipmentDeleteOne {
	builder := c.Delete().Where(equipment.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &EquipmentDeleteOne{builder}
}

// Create returns a query builder for Equipment.
func (c *EquipmentClient) Query() *EquipmentQuery {
	return &EquipmentQuery{config: c.config}
}

// Get returns a Equipment entity by its id.
func (c *EquipmentClient) Get(ctx context.Context, id int) (*Equipment, error) {
	return c.Query().Where(equipment.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EquipmentClient) GetX(ctx context.Context, id int) *Equipment {
	e, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return e
}

// QueryType queries the type edge of a Equipment.
func (c *EquipmentClient) QueryType(e *Equipment) *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := e.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, id),
			sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, equipment.TypeTable, equipment.TypeColumn),
		)
		fromV = sqlgraph.Neighbors(e.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryLocation queries the location edge of a Equipment.
func (c *EquipmentClient) QueryLocation(e *Equipment) *LocationQuery {
	query := &LocationQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := e.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, id),
			sqlgraph.To(location.Table, location.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, equipment.LocationTable, equipment.LocationColumn),
		)
		fromV = sqlgraph.Neighbors(e.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryParentPosition queries the parent_position edge of a Equipment.
func (c *EquipmentClient) QueryParentPosition(e *Equipment) *EquipmentPositionQuery {
	query := &EquipmentPositionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := e.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, id),
			sqlgraph.To(equipmentposition.Table, equipmentposition.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, equipment.ParentPositionTable, equipment.ParentPositionColumn),
		)
		fromV = sqlgraph.Neighbors(e.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryPositions queries the positions edge of a Equipment.
func (c *EquipmentClient) QueryPositions(e *Equipment) *EquipmentPositionQuery {
	query := &EquipmentPositionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := e.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, id),
			sqlgraph.To(equipmentposition.Table, equipmentposition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipment.PositionsTable, equipment.PositionsColumn),
		)
		fromV = sqlgraph.Neighbors(e.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryPorts queries the ports edge of a Equipment.
func (c *EquipmentClient) QueryPorts(e *Equipment) *EquipmentPortQuery {
	query := &EquipmentPortQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := e.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, id),
			sqlgraph.To(equipmentport.Table, equipmentport.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipment.PortsTable, equipment.PortsColumn),
		)
		fromV = sqlgraph.Neighbors(e.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryWorkOrder queries the work_order edge of a Equipment.
func (c *EquipmentClient) QueryWorkOrder(e *Equipment) *WorkOrderQuery {
	query := &WorkOrderQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := e.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, id),
			sqlgraph.To(workorder.Table, workorder.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, equipment.WorkOrderTable, equipment.WorkOrderColumn),
		)
		fromV = sqlgraph.Neighbors(e.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryProperties queries the properties edge of a Equipment.
func (c *EquipmentClient) QueryProperties(e *Equipment) *PropertyQuery {
	query := &PropertyQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := e.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, id),
			sqlgraph.To(property.Table, property.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipment.PropertiesTable, equipment.PropertiesColumn),
		)
		fromV = sqlgraph.Neighbors(e.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryFiles queries the files edge of a Equipment.
func (c *EquipmentClient) QueryFiles(e *Equipment) *FileQuery {
	query := &FileQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := e.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, id),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipment.FilesTable, equipment.FilesColumn),
		)
		fromV = sqlgraph.Neighbors(e.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryHyperlinks queries the hyperlinks edge of a Equipment.
func (c *EquipmentClient) QueryHyperlinks(e *Equipment) *HyperlinkQuery {
	query := &HyperlinkQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := e.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, id),
			sqlgraph.To(hyperlink.Table, hyperlink.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipment.HyperlinksTable, equipment.HyperlinksColumn),
		)
		fromV = sqlgraph.Neighbors(e.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEndpoints queries the endpoints edge of a Equipment.
func (c *EquipmentClient) QueryEndpoints(e *Equipment) *ServiceEndpointQuery {
	query := &ServiceEndpointQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := e.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipment.Table, equipment.FieldID, id),
			sqlgraph.To(serviceendpoint.Table, serviceendpoint.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, equipment.EndpointsTable, equipment.EndpointsColumn),
		)
		fromV = sqlgraph.Neighbors(e.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *EquipmentClient) Hooks() []Hook {
	return c.hooks.Equipment
}

// EquipmentCategoryClient is a client for the EquipmentCategory schema.
type EquipmentCategoryClient struct {
	config
}

// NewEquipmentCategoryClient returns a client for the EquipmentCategory from the given config.
func NewEquipmentCategoryClient(c config) *EquipmentCategoryClient {
	return &EquipmentCategoryClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `equipmentcategory.Hooks(f(g(h())))`.
func (c *EquipmentCategoryClient) Use(hooks ...Hook) {
	c.hooks.EquipmentCategory = append(c.hooks.EquipmentCategory, hooks...)
}

// Create returns a create builder for EquipmentCategory.
func (c *EquipmentCategoryClient) Create() *EquipmentCategoryCreate {
	mutation := newEquipmentCategoryMutation(c.config, OpCreate)
	return &EquipmentCategoryCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for EquipmentCategory.
func (c *EquipmentCategoryClient) Update() *EquipmentCategoryUpdate {
	mutation := newEquipmentCategoryMutation(c.config, OpUpdate)
	return &EquipmentCategoryUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *EquipmentCategoryClient) UpdateOne(ec *EquipmentCategory) *EquipmentCategoryUpdateOne {
	return c.UpdateOneID(ec.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *EquipmentCategoryClient) UpdateOneID(id int) *EquipmentCategoryUpdateOne {
	mutation := newEquipmentCategoryMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &EquipmentCategoryUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for EquipmentCategory.
func (c *EquipmentCategoryClient) Delete() *EquipmentCategoryDelete {
	mutation := newEquipmentCategoryMutation(c.config, OpDelete)
	return &EquipmentCategoryDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *EquipmentCategoryClient) DeleteOne(ec *EquipmentCategory) *EquipmentCategoryDeleteOne {
	return c.DeleteOneID(ec.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *EquipmentCategoryClient) DeleteOneID(id int) *EquipmentCategoryDeleteOne {
	builder := c.Delete().Where(equipmentcategory.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &EquipmentCategoryDeleteOne{builder}
}

// Create returns a query builder for EquipmentCategory.
func (c *EquipmentCategoryClient) Query() *EquipmentCategoryQuery {
	return &EquipmentCategoryQuery{config: c.config}
}

// Get returns a EquipmentCategory entity by its id.
func (c *EquipmentCategoryClient) Get(ctx context.Context, id int) (*EquipmentCategory, error) {
	return c.Query().Where(equipmentcategory.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EquipmentCategoryClient) GetX(ctx context.Context, id int) *EquipmentCategory {
	ec, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return ec
}

// QueryTypes queries the types edge of a EquipmentCategory.
func (c *EquipmentCategoryClient) QueryTypes(ec *EquipmentCategory) *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := ec.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentcategory.Table, equipmentcategory.FieldID, id),
			sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, equipmentcategory.TypesTable, equipmentcategory.TypesColumn),
		)
		fromV = sqlgraph.Neighbors(ec.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *EquipmentCategoryClient) Hooks() []Hook {
	return c.hooks.EquipmentCategory
}

// EquipmentPortClient is a client for the EquipmentPort schema.
type EquipmentPortClient struct {
	config
}

// NewEquipmentPortClient returns a client for the EquipmentPort from the given config.
func NewEquipmentPortClient(c config) *EquipmentPortClient {
	return &EquipmentPortClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `equipmentport.Hooks(f(g(h())))`.
func (c *EquipmentPortClient) Use(hooks ...Hook) {
	c.hooks.EquipmentPort = append(c.hooks.EquipmentPort, hooks...)
}

// Create returns a create builder for EquipmentPort.
func (c *EquipmentPortClient) Create() *EquipmentPortCreate {
	mutation := newEquipmentPortMutation(c.config, OpCreate)
	return &EquipmentPortCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for EquipmentPort.
func (c *EquipmentPortClient) Update() *EquipmentPortUpdate {
	mutation := newEquipmentPortMutation(c.config, OpUpdate)
	return &EquipmentPortUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *EquipmentPortClient) UpdateOne(ep *EquipmentPort) *EquipmentPortUpdateOne {
	return c.UpdateOneID(ep.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *EquipmentPortClient) UpdateOneID(id int) *EquipmentPortUpdateOne {
	mutation := newEquipmentPortMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &EquipmentPortUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for EquipmentPort.
func (c *EquipmentPortClient) Delete() *EquipmentPortDelete {
	mutation := newEquipmentPortMutation(c.config, OpDelete)
	return &EquipmentPortDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *EquipmentPortClient) DeleteOne(ep *EquipmentPort) *EquipmentPortDeleteOne {
	return c.DeleteOneID(ep.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *EquipmentPortClient) DeleteOneID(id int) *EquipmentPortDeleteOne {
	builder := c.Delete().Where(equipmentport.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &EquipmentPortDeleteOne{builder}
}

// Create returns a query builder for EquipmentPort.
func (c *EquipmentPortClient) Query() *EquipmentPortQuery {
	return &EquipmentPortQuery{config: c.config}
}

// Get returns a EquipmentPort entity by its id.
func (c *EquipmentPortClient) Get(ctx context.Context, id int) (*EquipmentPort, error) {
	return c.Query().Where(equipmentport.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EquipmentPortClient) GetX(ctx context.Context, id int) *EquipmentPort {
	ep, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return ep
}

// QueryDefinition queries the definition edge of a EquipmentPort.
func (c *EquipmentPortClient) QueryDefinition(ep *EquipmentPort) *EquipmentPortDefinitionQuery {
	query := &EquipmentPortDefinitionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := ep.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentport.Table, equipmentport.FieldID, id),
			sqlgraph.To(equipmentportdefinition.Table, equipmentportdefinition.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, equipmentport.DefinitionTable, equipmentport.DefinitionColumn),
		)
		fromV = sqlgraph.Neighbors(ep.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryParent queries the parent edge of a EquipmentPort.
func (c *EquipmentPortClient) QueryParent(ep *EquipmentPort) *EquipmentQuery {
	query := &EquipmentQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := ep.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentport.Table, equipmentport.FieldID, id),
			sqlgraph.To(equipment.Table, equipment.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, equipmentport.ParentTable, equipmentport.ParentColumn),
		)
		fromV = sqlgraph.Neighbors(ep.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryLink queries the link edge of a EquipmentPort.
func (c *EquipmentPortClient) QueryLink(ep *EquipmentPort) *LinkQuery {
	query := &LinkQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := ep.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentport.Table, equipmentport.FieldID, id),
			sqlgraph.To(link.Table, link.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, equipmentport.LinkTable, equipmentport.LinkColumn),
		)
		fromV = sqlgraph.Neighbors(ep.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryProperties queries the properties edge of a EquipmentPort.
func (c *EquipmentPortClient) QueryProperties(ep *EquipmentPort) *PropertyQuery {
	query := &PropertyQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := ep.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentport.Table, equipmentport.FieldID, id),
			sqlgraph.To(property.Table, property.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipmentport.PropertiesTable, equipmentport.PropertiesColumn),
		)
		fromV = sqlgraph.Neighbors(ep.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEndpoints queries the endpoints edge of a EquipmentPort.
func (c *EquipmentPortClient) QueryEndpoints(ep *EquipmentPort) *ServiceEndpointQuery {
	query := &ServiceEndpointQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := ep.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentport.Table, equipmentport.FieldID, id),
			sqlgraph.To(serviceendpoint.Table, serviceendpoint.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, equipmentport.EndpointsTable, equipmentport.EndpointsColumn),
		)
		fromV = sqlgraph.Neighbors(ep.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *EquipmentPortClient) Hooks() []Hook {
	return c.hooks.EquipmentPort
}

// EquipmentPortDefinitionClient is a client for the EquipmentPortDefinition schema.
type EquipmentPortDefinitionClient struct {
	config
}

// NewEquipmentPortDefinitionClient returns a client for the EquipmentPortDefinition from the given config.
func NewEquipmentPortDefinitionClient(c config) *EquipmentPortDefinitionClient {
	return &EquipmentPortDefinitionClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `equipmentportdefinition.Hooks(f(g(h())))`.
func (c *EquipmentPortDefinitionClient) Use(hooks ...Hook) {
	c.hooks.EquipmentPortDefinition = append(c.hooks.EquipmentPortDefinition, hooks...)
}

// Create returns a create builder for EquipmentPortDefinition.
func (c *EquipmentPortDefinitionClient) Create() *EquipmentPortDefinitionCreate {
	mutation := newEquipmentPortDefinitionMutation(c.config, OpCreate)
	return &EquipmentPortDefinitionCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for EquipmentPortDefinition.
func (c *EquipmentPortDefinitionClient) Update() *EquipmentPortDefinitionUpdate {
	mutation := newEquipmentPortDefinitionMutation(c.config, OpUpdate)
	return &EquipmentPortDefinitionUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *EquipmentPortDefinitionClient) UpdateOne(epd *EquipmentPortDefinition) *EquipmentPortDefinitionUpdateOne {
	return c.UpdateOneID(epd.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *EquipmentPortDefinitionClient) UpdateOneID(id int) *EquipmentPortDefinitionUpdateOne {
	mutation := newEquipmentPortDefinitionMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &EquipmentPortDefinitionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for EquipmentPortDefinition.
func (c *EquipmentPortDefinitionClient) Delete() *EquipmentPortDefinitionDelete {
	mutation := newEquipmentPortDefinitionMutation(c.config, OpDelete)
	return &EquipmentPortDefinitionDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *EquipmentPortDefinitionClient) DeleteOne(epd *EquipmentPortDefinition) *EquipmentPortDefinitionDeleteOne {
	return c.DeleteOneID(epd.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *EquipmentPortDefinitionClient) DeleteOneID(id int) *EquipmentPortDefinitionDeleteOne {
	builder := c.Delete().Where(equipmentportdefinition.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &EquipmentPortDefinitionDeleteOne{builder}
}

// Create returns a query builder for EquipmentPortDefinition.
func (c *EquipmentPortDefinitionClient) Query() *EquipmentPortDefinitionQuery {
	return &EquipmentPortDefinitionQuery{config: c.config}
}

// Get returns a EquipmentPortDefinition entity by its id.
func (c *EquipmentPortDefinitionClient) Get(ctx context.Context, id int) (*EquipmentPortDefinition, error) {
	return c.Query().Where(equipmentportdefinition.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EquipmentPortDefinitionClient) GetX(ctx context.Context, id int) *EquipmentPortDefinition {
	epd, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return epd
}

// QueryEquipmentPortType queries the equipment_port_type edge of a EquipmentPortDefinition.
func (c *EquipmentPortDefinitionClient) QueryEquipmentPortType(epd *EquipmentPortDefinition) *EquipmentPortTypeQuery {
	query := &EquipmentPortTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := epd.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentportdefinition.Table, equipmentportdefinition.FieldID, id),
			sqlgraph.To(equipmentporttype.Table, equipmentporttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, equipmentportdefinition.EquipmentPortTypeTable, equipmentportdefinition.EquipmentPortTypeColumn),
		)
		fromV = sqlgraph.Neighbors(epd.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryPorts queries the ports edge of a EquipmentPortDefinition.
func (c *EquipmentPortDefinitionClient) QueryPorts(epd *EquipmentPortDefinition) *EquipmentPortQuery {
	query := &EquipmentPortQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := epd.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentportdefinition.Table, equipmentportdefinition.FieldID, id),
			sqlgraph.To(equipmentport.Table, equipmentport.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, equipmentportdefinition.PortsTable, equipmentportdefinition.PortsColumn),
		)
		fromV = sqlgraph.Neighbors(epd.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEquipmentType queries the equipment_type edge of a EquipmentPortDefinition.
func (c *EquipmentPortDefinitionClient) QueryEquipmentType(epd *EquipmentPortDefinition) *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := epd.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentportdefinition.Table, equipmentportdefinition.FieldID, id),
			sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, equipmentportdefinition.EquipmentTypeTable, equipmentportdefinition.EquipmentTypeColumn),
		)
		fromV = sqlgraph.Neighbors(epd.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *EquipmentPortDefinitionClient) Hooks() []Hook {
	return c.hooks.EquipmentPortDefinition
}

// EquipmentPortTypeClient is a client for the EquipmentPortType schema.
type EquipmentPortTypeClient struct {
	config
}

// NewEquipmentPortTypeClient returns a client for the EquipmentPortType from the given config.
func NewEquipmentPortTypeClient(c config) *EquipmentPortTypeClient {
	return &EquipmentPortTypeClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `equipmentporttype.Hooks(f(g(h())))`.
func (c *EquipmentPortTypeClient) Use(hooks ...Hook) {
	c.hooks.EquipmentPortType = append(c.hooks.EquipmentPortType, hooks...)
}

// Create returns a create builder for EquipmentPortType.
func (c *EquipmentPortTypeClient) Create() *EquipmentPortTypeCreate {
	mutation := newEquipmentPortTypeMutation(c.config, OpCreate)
	return &EquipmentPortTypeCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for EquipmentPortType.
func (c *EquipmentPortTypeClient) Update() *EquipmentPortTypeUpdate {
	mutation := newEquipmentPortTypeMutation(c.config, OpUpdate)
	return &EquipmentPortTypeUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *EquipmentPortTypeClient) UpdateOne(ept *EquipmentPortType) *EquipmentPortTypeUpdateOne {
	return c.UpdateOneID(ept.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *EquipmentPortTypeClient) UpdateOneID(id int) *EquipmentPortTypeUpdateOne {
	mutation := newEquipmentPortTypeMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &EquipmentPortTypeUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for EquipmentPortType.
func (c *EquipmentPortTypeClient) Delete() *EquipmentPortTypeDelete {
	mutation := newEquipmentPortTypeMutation(c.config, OpDelete)
	return &EquipmentPortTypeDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *EquipmentPortTypeClient) DeleteOne(ept *EquipmentPortType) *EquipmentPortTypeDeleteOne {
	return c.DeleteOneID(ept.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *EquipmentPortTypeClient) DeleteOneID(id int) *EquipmentPortTypeDeleteOne {
	builder := c.Delete().Where(equipmentporttype.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &EquipmentPortTypeDeleteOne{builder}
}

// Create returns a query builder for EquipmentPortType.
func (c *EquipmentPortTypeClient) Query() *EquipmentPortTypeQuery {
	return &EquipmentPortTypeQuery{config: c.config}
}

// Get returns a EquipmentPortType entity by its id.
func (c *EquipmentPortTypeClient) Get(ctx context.Context, id int) (*EquipmentPortType, error) {
	return c.Query().Where(equipmentporttype.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EquipmentPortTypeClient) GetX(ctx context.Context, id int) *EquipmentPortType {
	ept, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return ept
}

// QueryPropertyTypes queries the property_types edge of a EquipmentPortType.
func (c *EquipmentPortTypeClient) QueryPropertyTypes(ept *EquipmentPortType) *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := ept.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentporttype.Table, equipmentporttype.FieldID, id),
			sqlgraph.To(propertytype.Table, propertytype.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipmentporttype.PropertyTypesTable, equipmentporttype.PropertyTypesColumn),
		)
		fromV = sqlgraph.Neighbors(ept.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryLinkPropertyTypes queries the link_property_types edge of a EquipmentPortType.
func (c *EquipmentPortTypeClient) QueryLinkPropertyTypes(ept *EquipmentPortType) *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := ept.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentporttype.Table, equipmentporttype.FieldID, id),
			sqlgraph.To(propertytype.Table, propertytype.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipmentporttype.LinkPropertyTypesTable, equipmentporttype.LinkPropertyTypesColumn),
		)
		fromV = sqlgraph.Neighbors(ept.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryPortDefinitions queries the port_definitions edge of a EquipmentPortType.
func (c *EquipmentPortTypeClient) QueryPortDefinitions(ept *EquipmentPortType) *EquipmentPortDefinitionQuery {
	query := &EquipmentPortDefinitionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := ept.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentporttype.Table, equipmentporttype.FieldID, id),
			sqlgraph.To(equipmentportdefinition.Table, equipmentportdefinition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, equipmentporttype.PortDefinitionsTable, equipmentporttype.PortDefinitionsColumn),
		)
		fromV = sqlgraph.Neighbors(ept.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *EquipmentPortTypeClient) Hooks() []Hook {
	return c.hooks.EquipmentPortType
}

// EquipmentPositionClient is a client for the EquipmentPosition schema.
type EquipmentPositionClient struct {
	config
}

// NewEquipmentPositionClient returns a client for the EquipmentPosition from the given config.
func NewEquipmentPositionClient(c config) *EquipmentPositionClient {
	return &EquipmentPositionClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `equipmentposition.Hooks(f(g(h())))`.
func (c *EquipmentPositionClient) Use(hooks ...Hook) {
	c.hooks.EquipmentPosition = append(c.hooks.EquipmentPosition, hooks...)
}

// Create returns a create builder for EquipmentPosition.
func (c *EquipmentPositionClient) Create() *EquipmentPositionCreate {
	mutation := newEquipmentPositionMutation(c.config, OpCreate)
	return &EquipmentPositionCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for EquipmentPosition.
func (c *EquipmentPositionClient) Update() *EquipmentPositionUpdate {
	mutation := newEquipmentPositionMutation(c.config, OpUpdate)
	return &EquipmentPositionUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *EquipmentPositionClient) UpdateOne(ep *EquipmentPosition) *EquipmentPositionUpdateOne {
	return c.UpdateOneID(ep.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *EquipmentPositionClient) UpdateOneID(id int) *EquipmentPositionUpdateOne {
	mutation := newEquipmentPositionMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &EquipmentPositionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for EquipmentPosition.
func (c *EquipmentPositionClient) Delete() *EquipmentPositionDelete {
	mutation := newEquipmentPositionMutation(c.config, OpDelete)
	return &EquipmentPositionDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *EquipmentPositionClient) DeleteOne(ep *EquipmentPosition) *EquipmentPositionDeleteOne {
	return c.DeleteOneID(ep.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *EquipmentPositionClient) DeleteOneID(id int) *EquipmentPositionDeleteOne {
	builder := c.Delete().Where(equipmentposition.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &EquipmentPositionDeleteOne{builder}
}

// Create returns a query builder for EquipmentPosition.
func (c *EquipmentPositionClient) Query() *EquipmentPositionQuery {
	return &EquipmentPositionQuery{config: c.config}
}

// Get returns a EquipmentPosition entity by its id.
func (c *EquipmentPositionClient) Get(ctx context.Context, id int) (*EquipmentPosition, error) {
	return c.Query().Where(equipmentposition.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EquipmentPositionClient) GetX(ctx context.Context, id int) *EquipmentPosition {
	ep, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return ep
}

// QueryDefinition queries the definition edge of a EquipmentPosition.
func (c *EquipmentPositionClient) QueryDefinition(ep *EquipmentPosition) *EquipmentPositionDefinitionQuery {
	query := &EquipmentPositionDefinitionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := ep.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentposition.Table, equipmentposition.FieldID, id),
			sqlgraph.To(equipmentpositiondefinition.Table, equipmentpositiondefinition.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, equipmentposition.DefinitionTable, equipmentposition.DefinitionColumn),
		)
		fromV = sqlgraph.Neighbors(ep.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryParent queries the parent edge of a EquipmentPosition.
func (c *EquipmentPositionClient) QueryParent(ep *EquipmentPosition) *EquipmentQuery {
	query := &EquipmentQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := ep.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentposition.Table, equipmentposition.FieldID, id),
			sqlgraph.To(equipment.Table, equipment.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, equipmentposition.ParentTable, equipmentposition.ParentColumn),
		)
		fromV = sqlgraph.Neighbors(ep.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryAttachment queries the attachment edge of a EquipmentPosition.
func (c *EquipmentPositionClient) QueryAttachment(ep *EquipmentPosition) *EquipmentQuery {
	query := &EquipmentQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := ep.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentposition.Table, equipmentposition.FieldID, id),
			sqlgraph.To(equipment.Table, equipment.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, equipmentposition.AttachmentTable, equipmentposition.AttachmentColumn),
		)
		fromV = sqlgraph.Neighbors(ep.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *EquipmentPositionClient) Hooks() []Hook {
	return c.hooks.EquipmentPosition
}

// EquipmentPositionDefinitionClient is a client for the EquipmentPositionDefinition schema.
type EquipmentPositionDefinitionClient struct {
	config
}

// NewEquipmentPositionDefinitionClient returns a client for the EquipmentPositionDefinition from the given config.
func NewEquipmentPositionDefinitionClient(c config) *EquipmentPositionDefinitionClient {
	return &EquipmentPositionDefinitionClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `equipmentpositiondefinition.Hooks(f(g(h())))`.
func (c *EquipmentPositionDefinitionClient) Use(hooks ...Hook) {
	c.hooks.EquipmentPositionDefinition = append(c.hooks.EquipmentPositionDefinition, hooks...)
}

// Create returns a create builder for EquipmentPositionDefinition.
func (c *EquipmentPositionDefinitionClient) Create() *EquipmentPositionDefinitionCreate {
	mutation := newEquipmentPositionDefinitionMutation(c.config, OpCreate)
	return &EquipmentPositionDefinitionCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for EquipmentPositionDefinition.
func (c *EquipmentPositionDefinitionClient) Update() *EquipmentPositionDefinitionUpdate {
	mutation := newEquipmentPositionDefinitionMutation(c.config, OpUpdate)
	return &EquipmentPositionDefinitionUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *EquipmentPositionDefinitionClient) UpdateOne(epd *EquipmentPositionDefinition) *EquipmentPositionDefinitionUpdateOne {
	return c.UpdateOneID(epd.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *EquipmentPositionDefinitionClient) UpdateOneID(id int) *EquipmentPositionDefinitionUpdateOne {
	mutation := newEquipmentPositionDefinitionMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &EquipmentPositionDefinitionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for EquipmentPositionDefinition.
func (c *EquipmentPositionDefinitionClient) Delete() *EquipmentPositionDefinitionDelete {
	mutation := newEquipmentPositionDefinitionMutation(c.config, OpDelete)
	return &EquipmentPositionDefinitionDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *EquipmentPositionDefinitionClient) DeleteOne(epd *EquipmentPositionDefinition) *EquipmentPositionDefinitionDeleteOne {
	return c.DeleteOneID(epd.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *EquipmentPositionDefinitionClient) DeleteOneID(id int) *EquipmentPositionDefinitionDeleteOne {
	builder := c.Delete().Where(equipmentpositiondefinition.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &EquipmentPositionDefinitionDeleteOne{builder}
}

// Create returns a query builder for EquipmentPositionDefinition.
func (c *EquipmentPositionDefinitionClient) Query() *EquipmentPositionDefinitionQuery {
	return &EquipmentPositionDefinitionQuery{config: c.config}
}

// Get returns a EquipmentPositionDefinition entity by its id.
func (c *EquipmentPositionDefinitionClient) Get(ctx context.Context, id int) (*EquipmentPositionDefinition, error) {
	return c.Query().Where(equipmentpositiondefinition.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EquipmentPositionDefinitionClient) GetX(ctx context.Context, id int) *EquipmentPositionDefinition {
	epd, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return epd
}

// QueryPositions queries the positions edge of a EquipmentPositionDefinition.
func (c *EquipmentPositionDefinitionClient) QueryPositions(epd *EquipmentPositionDefinition) *EquipmentPositionQuery {
	query := &EquipmentPositionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := epd.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentpositiondefinition.Table, equipmentpositiondefinition.FieldID, id),
			sqlgraph.To(equipmentposition.Table, equipmentposition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, equipmentpositiondefinition.PositionsTable, equipmentpositiondefinition.PositionsColumn),
		)
		fromV = sqlgraph.Neighbors(epd.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEquipmentType queries the equipment_type edge of a EquipmentPositionDefinition.
func (c *EquipmentPositionDefinitionClient) QueryEquipmentType(epd *EquipmentPositionDefinition) *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := epd.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmentpositiondefinition.Table, equipmentpositiondefinition.FieldID, id),
			sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, equipmentpositiondefinition.EquipmentTypeTable, equipmentpositiondefinition.EquipmentTypeColumn),
		)
		fromV = sqlgraph.Neighbors(epd.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *EquipmentPositionDefinitionClient) Hooks() []Hook {
	return c.hooks.EquipmentPositionDefinition
}

// EquipmentTypeClient is a client for the EquipmentType schema.
type EquipmentTypeClient struct {
	config
}

// NewEquipmentTypeClient returns a client for the EquipmentType from the given config.
func NewEquipmentTypeClient(c config) *EquipmentTypeClient {
	return &EquipmentTypeClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `equipmenttype.Hooks(f(g(h())))`.
func (c *EquipmentTypeClient) Use(hooks ...Hook) {
	c.hooks.EquipmentType = append(c.hooks.EquipmentType, hooks...)
}

// Create returns a create builder for EquipmentType.
func (c *EquipmentTypeClient) Create() *EquipmentTypeCreate {
	mutation := newEquipmentTypeMutation(c.config, OpCreate)
	return &EquipmentTypeCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for EquipmentType.
func (c *EquipmentTypeClient) Update() *EquipmentTypeUpdate {
	mutation := newEquipmentTypeMutation(c.config, OpUpdate)
	return &EquipmentTypeUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *EquipmentTypeClient) UpdateOne(et *EquipmentType) *EquipmentTypeUpdateOne {
	return c.UpdateOneID(et.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *EquipmentTypeClient) UpdateOneID(id int) *EquipmentTypeUpdateOne {
	mutation := newEquipmentTypeMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &EquipmentTypeUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for EquipmentType.
func (c *EquipmentTypeClient) Delete() *EquipmentTypeDelete {
	mutation := newEquipmentTypeMutation(c.config, OpDelete)
	return &EquipmentTypeDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *EquipmentTypeClient) DeleteOne(et *EquipmentType) *EquipmentTypeDeleteOne {
	return c.DeleteOneID(et.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *EquipmentTypeClient) DeleteOneID(id int) *EquipmentTypeDeleteOne {
	builder := c.Delete().Where(equipmenttype.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &EquipmentTypeDeleteOne{builder}
}

// Create returns a query builder for EquipmentType.
func (c *EquipmentTypeClient) Query() *EquipmentTypeQuery {
	return &EquipmentTypeQuery{config: c.config}
}

// Get returns a EquipmentType entity by its id.
func (c *EquipmentTypeClient) Get(ctx context.Context, id int) (*EquipmentType, error) {
	return c.Query().Where(equipmenttype.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EquipmentTypeClient) GetX(ctx context.Context, id int) *EquipmentType {
	et, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return et
}

// QueryPortDefinitions queries the port_definitions edge of a EquipmentType.
func (c *EquipmentTypeClient) QueryPortDefinitions(et *EquipmentType) *EquipmentPortDefinitionQuery {
	query := &EquipmentPortDefinitionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := et.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmenttype.Table, equipmenttype.FieldID, id),
			sqlgraph.To(equipmentportdefinition.Table, equipmentportdefinition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipmenttype.PortDefinitionsTable, equipmenttype.PortDefinitionsColumn),
		)
		fromV = sqlgraph.Neighbors(et.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryPositionDefinitions queries the position_definitions edge of a EquipmentType.
func (c *EquipmentTypeClient) QueryPositionDefinitions(et *EquipmentType) *EquipmentPositionDefinitionQuery {
	query := &EquipmentPositionDefinitionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := et.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmenttype.Table, equipmenttype.FieldID, id),
			sqlgraph.To(equipmentpositiondefinition.Table, equipmentpositiondefinition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipmenttype.PositionDefinitionsTable, equipmenttype.PositionDefinitionsColumn),
		)
		fromV = sqlgraph.Neighbors(et.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryPropertyTypes queries the property_types edge of a EquipmentType.
func (c *EquipmentTypeClient) QueryPropertyTypes(et *EquipmentType) *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := et.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmenttype.Table, equipmenttype.FieldID, id),
			sqlgraph.To(propertytype.Table, propertytype.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipmenttype.PropertyTypesTable, equipmenttype.PropertyTypesColumn),
		)
		fromV = sqlgraph.Neighbors(et.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEquipment queries the equipment edge of a EquipmentType.
func (c *EquipmentTypeClient) QueryEquipment(et *EquipmentType) *EquipmentQuery {
	query := &EquipmentQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := et.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmenttype.Table, equipmenttype.FieldID, id),
			sqlgraph.To(equipment.Table, equipment.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, equipmenttype.EquipmentTable, equipmenttype.EquipmentColumn),
		)
		fromV = sqlgraph.Neighbors(et.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryCategory queries the category edge of a EquipmentType.
func (c *EquipmentTypeClient) QueryCategory(et *EquipmentType) *EquipmentCategoryQuery {
	query := &EquipmentCategoryQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := et.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmenttype.Table, equipmenttype.FieldID, id),
			sqlgraph.To(equipmentcategory.Table, equipmentcategory.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, equipmenttype.CategoryTable, equipmenttype.CategoryColumn),
		)
		fromV = sqlgraph.Neighbors(et.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryServiceEndpointDefinitions queries the service_endpoint_definitions edge of a EquipmentType.
func (c *EquipmentTypeClient) QueryServiceEndpointDefinitions(et *EquipmentType) *ServiceEndpointDefinitionQuery {
	query := &ServiceEndpointDefinitionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := et.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(equipmenttype.Table, equipmenttype.FieldID, id),
			sqlgraph.To(serviceendpointdefinition.Table, serviceendpointdefinition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, equipmenttype.ServiceEndpointDefinitionsTable, equipmenttype.ServiceEndpointDefinitionsColumn),
		)
		fromV = sqlgraph.Neighbors(et.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *EquipmentTypeClient) Hooks() []Hook {
	return c.hooks.EquipmentType
}

// FileClient is a client for the File schema.
type FileClient struct {
	config
}

// NewFileClient returns a client for the File from the given config.
func NewFileClient(c config) *FileClient {
	return &FileClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `file.Hooks(f(g(h())))`.
func (c *FileClient) Use(hooks ...Hook) {
	c.hooks.File = append(c.hooks.File, hooks...)
}

// Create returns a create builder for File.
func (c *FileClient) Create() *FileCreate {
	mutation := newFileMutation(c.config, OpCreate)
	return &FileCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for File.
func (c *FileClient) Update() *FileUpdate {
	mutation := newFileMutation(c.config, OpUpdate)
	return &FileUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *FileClient) UpdateOne(f *File) *FileUpdateOne {
	return c.UpdateOneID(f.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *FileClient) UpdateOneID(id int) *FileUpdateOne {
	mutation := newFileMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &FileUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for File.
func (c *FileClient) Delete() *FileDelete {
	mutation := newFileMutation(c.config, OpDelete)
	return &FileDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *FileClient) DeleteOne(f *File) *FileDeleteOne {
	return c.DeleteOneID(f.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *FileClient) DeleteOneID(id int) *FileDeleteOne {
	builder := c.Delete().Where(file.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &FileDeleteOne{builder}
}

// Create returns a query builder for File.
func (c *FileClient) Query() *FileQuery {
	return &FileQuery{config: c.config}
}

// Get returns a File entity by its id.
func (c *FileClient) Get(ctx context.Context, id int) (*File, error) {
	return c.Query().Where(file.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *FileClient) GetX(ctx context.Context, id int) *File {
	f, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return f
}

// Hooks returns the client hooks.
func (c *FileClient) Hooks() []Hook {
	return c.hooks.File
}

// FloorPlanClient is a client for the FloorPlan schema.
type FloorPlanClient struct {
	config
}

// NewFloorPlanClient returns a client for the FloorPlan from the given config.
func NewFloorPlanClient(c config) *FloorPlanClient {
	return &FloorPlanClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `floorplan.Hooks(f(g(h())))`.
func (c *FloorPlanClient) Use(hooks ...Hook) {
	c.hooks.FloorPlan = append(c.hooks.FloorPlan, hooks...)
}

// Create returns a create builder for FloorPlan.
func (c *FloorPlanClient) Create() *FloorPlanCreate {
	mutation := newFloorPlanMutation(c.config, OpCreate)
	return &FloorPlanCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for FloorPlan.
func (c *FloorPlanClient) Update() *FloorPlanUpdate {
	mutation := newFloorPlanMutation(c.config, OpUpdate)
	return &FloorPlanUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *FloorPlanClient) UpdateOne(fp *FloorPlan) *FloorPlanUpdateOne {
	return c.UpdateOneID(fp.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *FloorPlanClient) UpdateOneID(id int) *FloorPlanUpdateOne {
	mutation := newFloorPlanMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &FloorPlanUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for FloorPlan.
func (c *FloorPlanClient) Delete() *FloorPlanDelete {
	mutation := newFloorPlanMutation(c.config, OpDelete)
	return &FloorPlanDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *FloorPlanClient) DeleteOne(fp *FloorPlan) *FloorPlanDeleteOne {
	return c.DeleteOneID(fp.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *FloorPlanClient) DeleteOneID(id int) *FloorPlanDeleteOne {
	builder := c.Delete().Where(floorplan.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &FloorPlanDeleteOne{builder}
}

// Create returns a query builder for FloorPlan.
func (c *FloorPlanClient) Query() *FloorPlanQuery {
	return &FloorPlanQuery{config: c.config}
}

// Get returns a FloorPlan entity by its id.
func (c *FloorPlanClient) Get(ctx context.Context, id int) (*FloorPlan, error) {
	return c.Query().Where(floorplan.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *FloorPlanClient) GetX(ctx context.Context, id int) *FloorPlan {
	fp, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return fp
}

// QueryLocation queries the location edge of a FloorPlan.
func (c *FloorPlanClient) QueryLocation(fp *FloorPlan) *LocationQuery {
	query := &LocationQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := fp.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(floorplan.Table, floorplan.FieldID, id),
			sqlgraph.To(location.Table, location.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, floorplan.LocationTable, floorplan.LocationColumn),
		)
		fromV = sqlgraph.Neighbors(fp.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryReferencePoint queries the reference_point edge of a FloorPlan.
func (c *FloorPlanClient) QueryReferencePoint(fp *FloorPlan) *FloorPlanReferencePointQuery {
	query := &FloorPlanReferencePointQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := fp.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(floorplan.Table, floorplan.FieldID, id),
			sqlgraph.To(floorplanreferencepoint.Table, floorplanreferencepoint.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, floorplan.ReferencePointTable, floorplan.ReferencePointColumn),
		)
		fromV = sqlgraph.Neighbors(fp.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryScale queries the scale edge of a FloorPlan.
func (c *FloorPlanClient) QueryScale(fp *FloorPlan) *FloorPlanScaleQuery {
	query := &FloorPlanScaleQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := fp.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(floorplan.Table, floorplan.FieldID, id),
			sqlgraph.To(floorplanscale.Table, floorplanscale.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, floorplan.ScaleTable, floorplan.ScaleColumn),
		)
		fromV = sqlgraph.Neighbors(fp.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryImage queries the image edge of a FloorPlan.
func (c *FloorPlanClient) QueryImage(fp *FloorPlan) *FileQuery {
	query := &FileQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := fp.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(floorplan.Table, floorplan.FieldID, id),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, floorplan.ImageTable, floorplan.ImageColumn),
		)
		fromV = sqlgraph.Neighbors(fp.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *FloorPlanClient) Hooks() []Hook {
	return c.hooks.FloorPlan
}

// FloorPlanReferencePointClient is a client for the FloorPlanReferencePoint schema.
type FloorPlanReferencePointClient struct {
	config
}

// NewFloorPlanReferencePointClient returns a client for the FloorPlanReferencePoint from the given config.
func NewFloorPlanReferencePointClient(c config) *FloorPlanReferencePointClient {
	return &FloorPlanReferencePointClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `floorplanreferencepoint.Hooks(f(g(h())))`.
func (c *FloorPlanReferencePointClient) Use(hooks ...Hook) {
	c.hooks.FloorPlanReferencePoint = append(c.hooks.FloorPlanReferencePoint, hooks...)
}

// Create returns a create builder for FloorPlanReferencePoint.
func (c *FloorPlanReferencePointClient) Create() *FloorPlanReferencePointCreate {
	mutation := newFloorPlanReferencePointMutation(c.config, OpCreate)
	return &FloorPlanReferencePointCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for FloorPlanReferencePoint.
func (c *FloorPlanReferencePointClient) Update() *FloorPlanReferencePointUpdate {
	mutation := newFloorPlanReferencePointMutation(c.config, OpUpdate)
	return &FloorPlanReferencePointUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *FloorPlanReferencePointClient) UpdateOne(fprp *FloorPlanReferencePoint) *FloorPlanReferencePointUpdateOne {
	return c.UpdateOneID(fprp.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *FloorPlanReferencePointClient) UpdateOneID(id int) *FloorPlanReferencePointUpdateOne {
	mutation := newFloorPlanReferencePointMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &FloorPlanReferencePointUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for FloorPlanReferencePoint.
func (c *FloorPlanReferencePointClient) Delete() *FloorPlanReferencePointDelete {
	mutation := newFloorPlanReferencePointMutation(c.config, OpDelete)
	return &FloorPlanReferencePointDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *FloorPlanReferencePointClient) DeleteOne(fprp *FloorPlanReferencePoint) *FloorPlanReferencePointDeleteOne {
	return c.DeleteOneID(fprp.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *FloorPlanReferencePointClient) DeleteOneID(id int) *FloorPlanReferencePointDeleteOne {
	builder := c.Delete().Where(floorplanreferencepoint.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &FloorPlanReferencePointDeleteOne{builder}
}

// Create returns a query builder for FloorPlanReferencePoint.
func (c *FloorPlanReferencePointClient) Query() *FloorPlanReferencePointQuery {
	return &FloorPlanReferencePointQuery{config: c.config}
}

// Get returns a FloorPlanReferencePoint entity by its id.
func (c *FloorPlanReferencePointClient) Get(ctx context.Context, id int) (*FloorPlanReferencePoint, error) {
	return c.Query().Where(floorplanreferencepoint.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *FloorPlanReferencePointClient) GetX(ctx context.Context, id int) *FloorPlanReferencePoint {
	fprp, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return fprp
}

// Hooks returns the client hooks.
func (c *FloorPlanReferencePointClient) Hooks() []Hook {
	return c.hooks.FloorPlanReferencePoint
}

// FloorPlanScaleClient is a client for the FloorPlanScale schema.
type FloorPlanScaleClient struct {
	config
}

// NewFloorPlanScaleClient returns a client for the FloorPlanScale from the given config.
func NewFloorPlanScaleClient(c config) *FloorPlanScaleClient {
	return &FloorPlanScaleClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `floorplanscale.Hooks(f(g(h())))`.
func (c *FloorPlanScaleClient) Use(hooks ...Hook) {
	c.hooks.FloorPlanScale = append(c.hooks.FloorPlanScale, hooks...)
}

// Create returns a create builder for FloorPlanScale.
func (c *FloorPlanScaleClient) Create() *FloorPlanScaleCreate {
	mutation := newFloorPlanScaleMutation(c.config, OpCreate)
	return &FloorPlanScaleCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for FloorPlanScale.
func (c *FloorPlanScaleClient) Update() *FloorPlanScaleUpdate {
	mutation := newFloorPlanScaleMutation(c.config, OpUpdate)
	return &FloorPlanScaleUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *FloorPlanScaleClient) UpdateOne(fps *FloorPlanScale) *FloorPlanScaleUpdateOne {
	return c.UpdateOneID(fps.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *FloorPlanScaleClient) UpdateOneID(id int) *FloorPlanScaleUpdateOne {
	mutation := newFloorPlanScaleMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &FloorPlanScaleUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for FloorPlanScale.
func (c *FloorPlanScaleClient) Delete() *FloorPlanScaleDelete {
	mutation := newFloorPlanScaleMutation(c.config, OpDelete)
	return &FloorPlanScaleDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *FloorPlanScaleClient) DeleteOne(fps *FloorPlanScale) *FloorPlanScaleDeleteOne {
	return c.DeleteOneID(fps.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *FloorPlanScaleClient) DeleteOneID(id int) *FloorPlanScaleDeleteOne {
	builder := c.Delete().Where(floorplanscale.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &FloorPlanScaleDeleteOne{builder}
}

// Create returns a query builder for FloorPlanScale.
func (c *FloorPlanScaleClient) Query() *FloorPlanScaleQuery {
	return &FloorPlanScaleQuery{config: c.config}
}

// Get returns a FloorPlanScale entity by its id.
func (c *FloorPlanScaleClient) Get(ctx context.Context, id int) (*FloorPlanScale, error) {
	return c.Query().Where(floorplanscale.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *FloorPlanScaleClient) GetX(ctx context.Context, id int) *FloorPlanScale {
	fps, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return fps
}

// Hooks returns the client hooks.
func (c *FloorPlanScaleClient) Hooks() []Hook {
	return c.hooks.FloorPlanScale
}

// HyperlinkClient is a client for the Hyperlink schema.
type HyperlinkClient struct {
	config
}

// NewHyperlinkClient returns a client for the Hyperlink from the given config.
func NewHyperlinkClient(c config) *HyperlinkClient {
	return &HyperlinkClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `hyperlink.Hooks(f(g(h())))`.
func (c *HyperlinkClient) Use(hooks ...Hook) {
	c.hooks.Hyperlink = append(c.hooks.Hyperlink, hooks...)
}

// Create returns a create builder for Hyperlink.
func (c *HyperlinkClient) Create() *HyperlinkCreate {
	mutation := newHyperlinkMutation(c.config, OpCreate)
	return &HyperlinkCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for Hyperlink.
func (c *HyperlinkClient) Update() *HyperlinkUpdate {
	mutation := newHyperlinkMutation(c.config, OpUpdate)
	return &HyperlinkUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *HyperlinkClient) UpdateOne(h *Hyperlink) *HyperlinkUpdateOne {
	return c.UpdateOneID(h.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *HyperlinkClient) UpdateOneID(id int) *HyperlinkUpdateOne {
	mutation := newHyperlinkMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &HyperlinkUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Hyperlink.
func (c *HyperlinkClient) Delete() *HyperlinkDelete {
	mutation := newHyperlinkMutation(c.config, OpDelete)
	return &HyperlinkDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *HyperlinkClient) DeleteOne(h *Hyperlink) *HyperlinkDeleteOne {
	return c.DeleteOneID(h.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *HyperlinkClient) DeleteOneID(id int) *HyperlinkDeleteOne {
	builder := c.Delete().Where(hyperlink.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &HyperlinkDeleteOne{builder}
}

// Create returns a query builder for Hyperlink.
func (c *HyperlinkClient) Query() *HyperlinkQuery {
	return &HyperlinkQuery{config: c.config}
}

// Get returns a Hyperlink entity by its id.
func (c *HyperlinkClient) Get(ctx context.Context, id int) (*Hyperlink, error) {
	return c.Query().Where(hyperlink.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *HyperlinkClient) GetX(ctx context.Context, id int) *Hyperlink {
	h, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return h
}

// Hooks returns the client hooks.
func (c *HyperlinkClient) Hooks() []Hook {
	return c.hooks.Hyperlink
}

// LinkClient is a client for the Link schema.
type LinkClient struct {
	config
}

// NewLinkClient returns a client for the Link from the given config.
func NewLinkClient(c config) *LinkClient {
	return &LinkClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `link.Hooks(f(g(h())))`.
func (c *LinkClient) Use(hooks ...Hook) {
	c.hooks.Link = append(c.hooks.Link, hooks...)
}

// Create returns a create builder for Link.
func (c *LinkClient) Create() *LinkCreate {
	mutation := newLinkMutation(c.config, OpCreate)
	return &LinkCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for Link.
func (c *LinkClient) Update() *LinkUpdate {
	mutation := newLinkMutation(c.config, OpUpdate)
	return &LinkUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *LinkClient) UpdateOne(l *Link) *LinkUpdateOne {
	return c.UpdateOneID(l.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *LinkClient) UpdateOneID(id int) *LinkUpdateOne {
	mutation := newLinkMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &LinkUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Link.
func (c *LinkClient) Delete() *LinkDelete {
	mutation := newLinkMutation(c.config, OpDelete)
	return &LinkDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *LinkClient) DeleteOne(l *Link) *LinkDeleteOne {
	return c.DeleteOneID(l.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *LinkClient) DeleteOneID(id int) *LinkDeleteOne {
	builder := c.Delete().Where(link.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &LinkDeleteOne{builder}
}

// Create returns a query builder for Link.
func (c *LinkClient) Query() *LinkQuery {
	return &LinkQuery{config: c.config}
}

// Get returns a Link entity by its id.
func (c *LinkClient) Get(ctx context.Context, id int) (*Link, error) {
	return c.Query().Where(link.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *LinkClient) GetX(ctx context.Context, id int) *Link {
	l, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return l
}

// QueryPorts queries the ports edge of a Link.
func (c *LinkClient) QueryPorts(l *Link) *EquipmentPortQuery {
	query := &EquipmentPortQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := l.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(link.Table, link.FieldID, id),
			sqlgraph.To(equipmentport.Table, equipmentport.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, link.PortsTable, link.PortsColumn),
		)
		fromV = sqlgraph.Neighbors(l.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryWorkOrder queries the work_order edge of a Link.
func (c *LinkClient) QueryWorkOrder(l *Link) *WorkOrderQuery {
	query := &WorkOrderQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := l.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(link.Table, link.FieldID, id),
			sqlgraph.To(workorder.Table, workorder.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, link.WorkOrderTable, link.WorkOrderColumn),
		)
		fromV = sqlgraph.Neighbors(l.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryProperties queries the properties edge of a Link.
func (c *LinkClient) QueryProperties(l *Link) *PropertyQuery {
	query := &PropertyQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := l.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(link.Table, link.FieldID, id),
			sqlgraph.To(property.Table, property.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, link.PropertiesTable, link.PropertiesColumn),
		)
		fromV = sqlgraph.Neighbors(l.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryService queries the service edge of a Link.
func (c *LinkClient) QueryService(l *Link) *ServiceQuery {
	query := &ServiceQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := l.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(link.Table, link.FieldID, id),
			sqlgraph.To(service.Table, service.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, link.ServiceTable, link.ServicePrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(l.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *LinkClient) Hooks() []Hook {
	return c.hooks.Link
}

// LocationClient is a client for the Location schema.
type LocationClient struct {
	config
}

// NewLocationClient returns a client for the Location from the given config.
func NewLocationClient(c config) *LocationClient {
	return &LocationClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `location.Hooks(f(g(h())))`.
func (c *LocationClient) Use(hooks ...Hook) {
	c.hooks.Location = append(c.hooks.Location, hooks...)
}

// Create returns a create builder for Location.
func (c *LocationClient) Create() *LocationCreate {
	mutation := newLocationMutation(c.config, OpCreate)
	return &LocationCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for Location.
func (c *LocationClient) Update() *LocationUpdate {
	mutation := newLocationMutation(c.config, OpUpdate)
	return &LocationUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *LocationClient) UpdateOne(l *Location) *LocationUpdateOne {
	return c.UpdateOneID(l.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *LocationClient) UpdateOneID(id int) *LocationUpdateOne {
	mutation := newLocationMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &LocationUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Location.
func (c *LocationClient) Delete() *LocationDelete {
	mutation := newLocationMutation(c.config, OpDelete)
	return &LocationDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *LocationClient) DeleteOne(l *Location) *LocationDeleteOne {
	return c.DeleteOneID(l.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *LocationClient) DeleteOneID(id int) *LocationDeleteOne {
	builder := c.Delete().Where(location.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &LocationDeleteOne{builder}
}

// Create returns a query builder for Location.
func (c *LocationClient) Query() *LocationQuery {
	return &LocationQuery{config: c.config}
}

// Get returns a Location entity by its id.
func (c *LocationClient) Get(ctx context.Context, id int) (*Location, error) {
	return c.Query().Where(location.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *LocationClient) GetX(ctx context.Context, id int) *Location {
	l, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return l
}

// QueryType queries the type edge of a Location.
func (c *LocationClient) QueryType(l *Location) *LocationTypeQuery {
	query := &LocationTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := l.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(location.Table, location.FieldID, id),
			sqlgraph.To(locationtype.Table, locationtype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, location.TypeTable, location.TypeColumn),
		)
		fromV = sqlgraph.Neighbors(l.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryParent queries the parent edge of a Location.
func (c *LocationClient) QueryParent(l *Location) *LocationQuery {
	query := &LocationQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := l.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(location.Table, location.FieldID, id),
			sqlgraph.To(location.Table, location.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, location.ParentTable, location.ParentColumn),
		)
		fromV = sqlgraph.Neighbors(l.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryChildren queries the children edge of a Location.
func (c *LocationClient) QueryChildren(l *Location) *LocationQuery {
	query := &LocationQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := l.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(location.Table, location.FieldID, id),
			sqlgraph.To(location.Table, location.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, location.ChildrenTable, location.ChildrenColumn),
		)
		fromV = sqlgraph.Neighbors(l.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryFiles queries the files edge of a Location.
func (c *LocationClient) QueryFiles(l *Location) *FileQuery {
	query := &FileQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := l.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(location.Table, location.FieldID, id),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, location.FilesTable, location.FilesColumn),
		)
		fromV = sqlgraph.Neighbors(l.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryHyperlinks queries the hyperlinks edge of a Location.
func (c *LocationClient) QueryHyperlinks(l *Location) *HyperlinkQuery {
	query := &HyperlinkQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := l.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(location.Table, location.FieldID, id),
			sqlgraph.To(hyperlink.Table, hyperlink.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, location.HyperlinksTable, location.HyperlinksColumn),
		)
		fromV = sqlgraph.Neighbors(l.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEquipment queries the equipment edge of a Location.
func (c *LocationClient) QueryEquipment(l *Location) *EquipmentQuery {
	query := &EquipmentQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := l.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(location.Table, location.FieldID, id),
			sqlgraph.To(equipment.Table, equipment.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, location.EquipmentTable, location.EquipmentColumn),
		)
		fromV = sqlgraph.Neighbors(l.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryProperties queries the properties edge of a Location.
func (c *LocationClient) QueryProperties(l *Location) *PropertyQuery {
	query := &PropertyQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := l.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(location.Table, location.FieldID, id),
			sqlgraph.To(property.Table, property.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, location.PropertiesTable, location.PropertiesColumn),
		)
		fromV = sqlgraph.Neighbors(l.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QuerySurvey queries the survey edge of a Location.
func (c *LocationClient) QuerySurvey(l *Location) *SurveyQuery {
	query := &SurveyQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := l.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(location.Table, location.FieldID, id),
			sqlgraph.To(survey.Table, survey.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, location.SurveyTable, location.SurveyColumn),
		)
		fromV = sqlgraph.Neighbors(l.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryWifiScan queries the wifi_scan edge of a Location.
func (c *LocationClient) QueryWifiScan(l *Location) *SurveyWiFiScanQuery {
	query := &SurveyWiFiScanQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := l.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(location.Table, location.FieldID, id),
			sqlgraph.To(surveywifiscan.Table, surveywifiscan.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, location.WifiScanTable, location.WifiScanColumn),
		)
		fromV = sqlgraph.Neighbors(l.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryCellScan queries the cell_scan edge of a Location.
func (c *LocationClient) QueryCellScan(l *Location) *SurveyCellScanQuery {
	query := &SurveyCellScanQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := l.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(location.Table, location.FieldID, id),
			sqlgraph.To(surveycellscan.Table, surveycellscan.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, location.CellScanTable, location.CellScanColumn),
		)
		fromV = sqlgraph.Neighbors(l.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryWorkOrders queries the work_orders edge of a Location.
func (c *LocationClient) QueryWorkOrders(l *Location) *WorkOrderQuery {
	query := &WorkOrderQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := l.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(location.Table, location.FieldID, id),
			sqlgraph.To(workorder.Table, workorder.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, location.WorkOrdersTable, location.WorkOrdersColumn),
		)
		fromV = sqlgraph.Neighbors(l.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryFloorPlans queries the floor_plans edge of a Location.
func (c *LocationClient) QueryFloorPlans(l *Location) *FloorPlanQuery {
	query := &FloorPlanQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := l.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(location.Table, location.FieldID, id),
			sqlgraph.To(floorplan.Table, floorplan.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, location.FloorPlansTable, location.FloorPlansColumn),
		)
		fromV = sqlgraph.Neighbors(l.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *LocationClient) Hooks() []Hook {
	return c.hooks.Location
}

// LocationTypeClient is a client for the LocationType schema.
type LocationTypeClient struct {
	config
}

// NewLocationTypeClient returns a client for the LocationType from the given config.
func NewLocationTypeClient(c config) *LocationTypeClient {
	return &LocationTypeClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `locationtype.Hooks(f(g(h())))`.
func (c *LocationTypeClient) Use(hooks ...Hook) {
	c.hooks.LocationType = append(c.hooks.LocationType, hooks...)
}

// Create returns a create builder for LocationType.
func (c *LocationTypeClient) Create() *LocationTypeCreate {
	mutation := newLocationTypeMutation(c.config, OpCreate)
	return &LocationTypeCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for LocationType.
func (c *LocationTypeClient) Update() *LocationTypeUpdate {
	mutation := newLocationTypeMutation(c.config, OpUpdate)
	return &LocationTypeUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *LocationTypeClient) UpdateOne(lt *LocationType) *LocationTypeUpdateOne {
	return c.UpdateOneID(lt.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *LocationTypeClient) UpdateOneID(id int) *LocationTypeUpdateOne {
	mutation := newLocationTypeMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &LocationTypeUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for LocationType.
func (c *LocationTypeClient) Delete() *LocationTypeDelete {
	mutation := newLocationTypeMutation(c.config, OpDelete)
	return &LocationTypeDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *LocationTypeClient) DeleteOne(lt *LocationType) *LocationTypeDeleteOne {
	return c.DeleteOneID(lt.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *LocationTypeClient) DeleteOneID(id int) *LocationTypeDeleteOne {
	builder := c.Delete().Where(locationtype.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &LocationTypeDeleteOne{builder}
}

// Create returns a query builder for LocationType.
func (c *LocationTypeClient) Query() *LocationTypeQuery {
	return &LocationTypeQuery{config: c.config}
}

// Get returns a LocationType entity by its id.
func (c *LocationTypeClient) Get(ctx context.Context, id int) (*LocationType, error) {
	return c.Query().Where(locationtype.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *LocationTypeClient) GetX(ctx context.Context, id int) *LocationType {
	lt, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return lt
}

// QueryLocations queries the locations edge of a LocationType.
func (c *LocationTypeClient) QueryLocations(lt *LocationType) *LocationQuery {
	query := &LocationQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := lt.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(locationtype.Table, locationtype.FieldID, id),
			sqlgraph.To(location.Table, location.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, locationtype.LocationsTable, locationtype.LocationsColumn),
		)
		fromV = sqlgraph.Neighbors(lt.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryPropertyTypes queries the property_types edge of a LocationType.
func (c *LocationTypeClient) QueryPropertyTypes(lt *LocationType) *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := lt.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(locationtype.Table, locationtype.FieldID, id),
			sqlgraph.To(propertytype.Table, propertytype.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, locationtype.PropertyTypesTable, locationtype.PropertyTypesColumn),
		)
		fromV = sqlgraph.Neighbors(lt.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QuerySurveyTemplateCategories queries the survey_template_categories edge of a LocationType.
func (c *LocationTypeClient) QuerySurveyTemplateCategories(lt *LocationType) *SurveyTemplateCategoryQuery {
	query := &SurveyTemplateCategoryQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := lt.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(locationtype.Table, locationtype.FieldID, id),
			sqlgraph.To(surveytemplatecategory.Table, surveytemplatecategory.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, locationtype.SurveyTemplateCategoriesTable, locationtype.SurveyTemplateCategoriesColumn),
		)
		fromV = sqlgraph.Neighbors(lt.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *LocationTypeClient) Hooks() []Hook {
	return c.hooks.LocationType
}

// PermissionsPolicyClient is a client for the PermissionsPolicy schema.
type PermissionsPolicyClient struct {
	config
}

// NewPermissionsPolicyClient returns a client for the PermissionsPolicy from the given config.
func NewPermissionsPolicyClient(c config) *PermissionsPolicyClient {
	return &PermissionsPolicyClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `permissionspolicy.Hooks(f(g(h())))`.
func (c *PermissionsPolicyClient) Use(hooks ...Hook) {
	c.hooks.PermissionsPolicy = append(c.hooks.PermissionsPolicy, hooks...)
}

// Create returns a create builder for PermissionsPolicy.
func (c *PermissionsPolicyClient) Create() *PermissionsPolicyCreate {
	mutation := newPermissionsPolicyMutation(c.config, OpCreate)
	return &PermissionsPolicyCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for PermissionsPolicy.
func (c *PermissionsPolicyClient) Update() *PermissionsPolicyUpdate {
	mutation := newPermissionsPolicyMutation(c.config, OpUpdate)
	return &PermissionsPolicyUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *PermissionsPolicyClient) UpdateOne(pp *PermissionsPolicy) *PermissionsPolicyUpdateOne {
	return c.UpdateOneID(pp.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *PermissionsPolicyClient) UpdateOneID(id int) *PermissionsPolicyUpdateOne {
	mutation := newPermissionsPolicyMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &PermissionsPolicyUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for PermissionsPolicy.
func (c *PermissionsPolicyClient) Delete() *PermissionsPolicyDelete {
	mutation := newPermissionsPolicyMutation(c.config, OpDelete)
	return &PermissionsPolicyDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *PermissionsPolicyClient) DeleteOne(pp *PermissionsPolicy) *PermissionsPolicyDeleteOne {
	return c.DeleteOneID(pp.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *PermissionsPolicyClient) DeleteOneID(id int) *PermissionsPolicyDeleteOne {
	builder := c.Delete().Where(permissionspolicy.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &PermissionsPolicyDeleteOne{builder}
}

// Create returns a query builder for PermissionsPolicy.
func (c *PermissionsPolicyClient) Query() *PermissionsPolicyQuery {
	return &PermissionsPolicyQuery{config: c.config}
}

// Get returns a PermissionsPolicy entity by its id.
func (c *PermissionsPolicyClient) Get(ctx context.Context, id int) (*PermissionsPolicy, error) {
	return c.Query().Where(permissionspolicy.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *PermissionsPolicyClient) GetX(ctx context.Context, id int) *PermissionsPolicy {
	pp, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return pp
}

// QueryGroups queries the groups edge of a PermissionsPolicy.
func (c *PermissionsPolicyClient) QueryGroups(pp *PermissionsPolicy) *UsersGroupQuery {
	query := &UsersGroupQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pp.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(permissionspolicy.Table, permissionspolicy.FieldID, id),
			sqlgraph.To(usersgroup.Table, usersgroup.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, permissionspolicy.GroupsTable, permissionspolicy.GroupsPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(pp.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *PermissionsPolicyClient) Hooks() []Hook {
	return c.hooks.PermissionsPolicy
}

// ProjectClient is a client for the Project schema.
type ProjectClient struct {
	config
}

// NewProjectClient returns a client for the Project from the given config.
func NewProjectClient(c config) *ProjectClient {
	return &ProjectClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `project.Hooks(f(g(h())))`.
func (c *ProjectClient) Use(hooks ...Hook) {
	c.hooks.Project = append(c.hooks.Project, hooks...)
}

// Create returns a create builder for Project.
func (c *ProjectClient) Create() *ProjectCreate {
	mutation := newProjectMutation(c.config, OpCreate)
	return &ProjectCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for Project.
func (c *ProjectClient) Update() *ProjectUpdate {
	mutation := newProjectMutation(c.config, OpUpdate)
	return &ProjectUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ProjectClient) UpdateOne(pr *Project) *ProjectUpdateOne {
	return c.UpdateOneID(pr.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *ProjectClient) UpdateOneID(id int) *ProjectUpdateOne {
	mutation := newProjectMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &ProjectUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Project.
func (c *ProjectClient) Delete() *ProjectDelete {
	mutation := newProjectMutation(c.config, OpDelete)
	return &ProjectDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *ProjectClient) DeleteOne(pr *Project) *ProjectDeleteOne {
	return c.DeleteOneID(pr.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *ProjectClient) DeleteOneID(id int) *ProjectDeleteOne {
	builder := c.Delete().Where(project.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ProjectDeleteOne{builder}
}

// Create returns a query builder for Project.
func (c *ProjectClient) Query() *ProjectQuery {
	return &ProjectQuery{config: c.config}
}

// Get returns a Project entity by its id.
func (c *ProjectClient) Get(ctx context.Context, id int) (*Project, error) {
	return c.Query().Where(project.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ProjectClient) GetX(ctx context.Context, id int) *Project {
	pr, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return pr
}

// QueryType queries the type edge of a Project.
func (c *ProjectClient) QueryType(pr *Project) *ProjectTypeQuery {
	query := &ProjectTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(project.Table, project.FieldID, id),
			sqlgraph.To(projecttype.Table, projecttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, project.TypeTable, project.TypeColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryLocation queries the location edge of a Project.
func (c *ProjectClient) QueryLocation(pr *Project) *LocationQuery {
	query := &LocationQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(project.Table, project.FieldID, id),
			sqlgraph.To(location.Table, location.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, project.LocationTable, project.LocationColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryComments queries the comments edge of a Project.
func (c *ProjectClient) QueryComments(pr *Project) *CommentQuery {
	query := &CommentQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(project.Table, project.FieldID, id),
			sqlgraph.To(comment.Table, comment.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, project.CommentsTable, project.CommentsColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryWorkOrders queries the work_orders edge of a Project.
func (c *ProjectClient) QueryWorkOrders(pr *Project) *WorkOrderQuery {
	query := &WorkOrderQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(project.Table, project.FieldID, id),
			sqlgraph.To(workorder.Table, workorder.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, project.WorkOrdersTable, project.WorkOrdersColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryProperties queries the properties edge of a Project.
func (c *ProjectClient) QueryProperties(pr *Project) *PropertyQuery {
	query := &PropertyQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(project.Table, project.FieldID, id),
			sqlgraph.To(property.Table, property.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, project.PropertiesTable, project.PropertiesColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryCreator queries the creator edge of a Project.
func (c *ProjectClient) QueryCreator(pr *Project) *UserQuery {
	query := &UserQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(project.Table, project.FieldID, id),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, project.CreatorTable, project.CreatorColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ProjectClient) Hooks() []Hook {
	return c.hooks.Project
}

// ProjectTypeClient is a client for the ProjectType schema.
type ProjectTypeClient struct {
	config
}

// NewProjectTypeClient returns a client for the ProjectType from the given config.
func NewProjectTypeClient(c config) *ProjectTypeClient {
	return &ProjectTypeClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `projecttype.Hooks(f(g(h())))`.
func (c *ProjectTypeClient) Use(hooks ...Hook) {
	c.hooks.ProjectType = append(c.hooks.ProjectType, hooks...)
}

// Create returns a create builder for ProjectType.
func (c *ProjectTypeClient) Create() *ProjectTypeCreate {
	mutation := newProjectTypeMutation(c.config, OpCreate)
	return &ProjectTypeCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for ProjectType.
func (c *ProjectTypeClient) Update() *ProjectTypeUpdate {
	mutation := newProjectTypeMutation(c.config, OpUpdate)
	return &ProjectTypeUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ProjectTypeClient) UpdateOne(pt *ProjectType) *ProjectTypeUpdateOne {
	return c.UpdateOneID(pt.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *ProjectTypeClient) UpdateOneID(id int) *ProjectTypeUpdateOne {
	mutation := newProjectTypeMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &ProjectTypeUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for ProjectType.
func (c *ProjectTypeClient) Delete() *ProjectTypeDelete {
	mutation := newProjectTypeMutation(c.config, OpDelete)
	return &ProjectTypeDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *ProjectTypeClient) DeleteOne(pt *ProjectType) *ProjectTypeDeleteOne {
	return c.DeleteOneID(pt.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *ProjectTypeClient) DeleteOneID(id int) *ProjectTypeDeleteOne {
	builder := c.Delete().Where(projecttype.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ProjectTypeDeleteOne{builder}
}

// Create returns a query builder for ProjectType.
func (c *ProjectTypeClient) Query() *ProjectTypeQuery {
	return &ProjectTypeQuery{config: c.config}
}

// Get returns a ProjectType entity by its id.
func (c *ProjectTypeClient) Get(ctx context.Context, id int) (*ProjectType, error) {
	return c.Query().Where(projecttype.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ProjectTypeClient) GetX(ctx context.Context, id int) *ProjectType {
	pt, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return pt
}

// QueryProjects queries the projects edge of a ProjectType.
func (c *ProjectTypeClient) QueryProjects(pt *ProjectType) *ProjectQuery {
	query := &ProjectQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pt.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(projecttype.Table, projecttype.FieldID, id),
			sqlgraph.To(project.Table, project.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, projecttype.ProjectsTable, projecttype.ProjectsColumn),
		)
		fromV = sqlgraph.Neighbors(pt.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryProperties queries the properties edge of a ProjectType.
func (c *ProjectTypeClient) QueryProperties(pt *ProjectType) *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pt.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(projecttype.Table, projecttype.FieldID, id),
			sqlgraph.To(propertytype.Table, propertytype.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, projecttype.PropertiesTable, projecttype.PropertiesColumn),
		)
		fromV = sqlgraph.Neighbors(pt.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryWorkOrders queries the work_orders edge of a ProjectType.
func (c *ProjectTypeClient) QueryWorkOrders(pt *ProjectType) *WorkOrderDefinitionQuery {
	query := &WorkOrderDefinitionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pt.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(projecttype.Table, projecttype.FieldID, id),
			sqlgraph.To(workorderdefinition.Table, workorderdefinition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, projecttype.WorkOrdersTable, projecttype.WorkOrdersColumn),
		)
		fromV = sqlgraph.Neighbors(pt.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ProjectTypeClient) Hooks() []Hook {
	return c.hooks.ProjectType
}

// PropertyClient is a client for the Property schema.
type PropertyClient struct {
	config
}

// NewPropertyClient returns a client for the Property from the given config.
func NewPropertyClient(c config) *PropertyClient {
	return &PropertyClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `property.Hooks(f(g(h())))`.
func (c *PropertyClient) Use(hooks ...Hook) {
	c.hooks.Property = append(c.hooks.Property, hooks...)
}

// Create returns a create builder for Property.
func (c *PropertyClient) Create() *PropertyCreate {
	mutation := newPropertyMutation(c.config, OpCreate)
	return &PropertyCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for Property.
func (c *PropertyClient) Update() *PropertyUpdate {
	mutation := newPropertyMutation(c.config, OpUpdate)
	return &PropertyUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *PropertyClient) UpdateOne(pr *Property) *PropertyUpdateOne {
	return c.UpdateOneID(pr.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *PropertyClient) UpdateOneID(id int) *PropertyUpdateOne {
	mutation := newPropertyMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &PropertyUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Property.
func (c *PropertyClient) Delete() *PropertyDelete {
	mutation := newPropertyMutation(c.config, OpDelete)
	return &PropertyDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *PropertyClient) DeleteOne(pr *Property) *PropertyDeleteOne {
	return c.DeleteOneID(pr.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *PropertyClient) DeleteOneID(id int) *PropertyDeleteOne {
	builder := c.Delete().Where(property.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &PropertyDeleteOne{builder}
}

// Create returns a query builder for Property.
func (c *PropertyClient) Query() *PropertyQuery {
	return &PropertyQuery{config: c.config}
}

// Get returns a Property entity by its id.
func (c *PropertyClient) Get(ctx context.Context, id int) (*Property, error) {
	return c.Query().Where(property.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *PropertyClient) GetX(ctx context.Context, id int) *Property {
	pr, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return pr
}

// QueryType queries the type edge of a Property.
func (c *PropertyClient) QueryType(pr *Property) *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(property.Table, property.FieldID, id),
			sqlgraph.To(propertytype.Table, propertytype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, property.TypeTable, property.TypeColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryLocation queries the location edge of a Property.
func (c *PropertyClient) QueryLocation(pr *Property) *LocationQuery {
	query := &LocationQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(property.Table, property.FieldID, id),
			sqlgraph.To(location.Table, location.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, property.LocationTable, property.LocationColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEquipment queries the equipment edge of a Property.
func (c *PropertyClient) QueryEquipment(pr *Property) *EquipmentQuery {
	query := &EquipmentQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(property.Table, property.FieldID, id),
			sqlgraph.To(equipment.Table, equipment.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, property.EquipmentTable, property.EquipmentColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryService queries the service edge of a Property.
func (c *PropertyClient) QueryService(pr *Property) *ServiceQuery {
	query := &ServiceQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(property.Table, property.FieldID, id),
			sqlgraph.To(service.Table, service.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, property.ServiceTable, property.ServiceColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEquipmentPort queries the equipment_port edge of a Property.
func (c *PropertyClient) QueryEquipmentPort(pr *Property) *EquipmentPortQuery {
	query := &EquipmentPortQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(property.Table, property.FieldID, id),
			sqlgraph.To(equipmentport.Table, equipmentport.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, property.EquipmentPortTable, property.EquipmentPortColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryLink queries the link edge of a Property.
func (c *PropertyClient) QueryLink(pr *Property) *LinkQuery {
	query := &LinkQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(property.Table, property.FieldID, id),
			sqlgraph.To(link.Table, link.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, property.LinkTable, property.LinkColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryWorkOrder queries the work_order edge of a Property.
func (c *PropertyClient) QueryWorkOrder(pr *Property) *WorkOrderQuery {
	query := &WorkOrderQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(property.Table, property.FieldID, id),
			sqlgraph.To(workorder.Table, workorder.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, property.WorkOrderTable, property.WorkOrderColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryProject queries the project edge of a Property.
func (c *PropertyClient) QueryProject(pr *Property) *ProjectQuery {
	query := &ProjectQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(property.Table, property.FieldID, id),
			sqlgraph.To(project.Table, project.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, property.ProjectTable, property.ProjectColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEquipmentValue queries the equipment_value edge of a Property.
func (c *PropertyClient) QueryEquipmentValue(pr *Property) *EquipmentQuery {
	query := &EquipmentQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(property.Table, property.FieldID, id),
			sqlgraph.To(equipment.Table, equipment.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, property.EquipmentValueTable, property.EquipmentValueColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryLocationValue queries the location_value edge of a Property.
func (c *PropertyClient) QueryLocationValue(pr *Property) *LocationQuery {
	query := &LocationQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(property.Table, property.FieldID, id),
			sqlgraph.To(location.Table, location.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, property.LocationValueTable, property.LocationValueColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryServiceValue queries the service_value edge of a Property.
func (c *PropertyClient) QueryServiceValue(pr *Property) *ServiceQuery {
	query := &ServiceQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(property.Table, property.FieldID, id),
			sqlgraph.To(service.Table, service.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, property.ServiceValueTable, property.ServiceValueColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *PropertyClient) Hooks() []Hook {
	return c.hooks.Property
}

// PropertyTypeClient is a client for the PropertyType schema.
type PropertyTypeClient struct {
	config
}

// NewPropertyTypeClient returns a client for the PropertyType from the given config.
func NewPropertyTypeClient(c config) *PropertyTypeClient {
	return &PropertyTypeClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `propertytype.Hooks(f(g(h())))`.
func (c *PropertyTypeClient) Use(hooks ...Hook) {
	c.hooks.PropertyType = append(c.hooks.PropertyType, hooks...)
}

// Create returns a create builder for PropertyType.
func (c *PropertyTypeClient) Create() *PropertyTypeCreate {
	mutation := newPropertyTypeMutation(c.config, OpCreate)
	return &PropertyTypeCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for PropertyType.
func (c *PropertyTypeClient) Update() *PropertyTypeUpdate {
	mutation := newPropertyTypeMutation(c.config, OpUpdate)
	return &PropertyTypeUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *PropertyTypeClient) UpdateOne(pt *PropertyType) *PropertyTypeUpdateOne {
	return c.UpdateOneID(pt.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *PropertyTypeClient) UpdateOneID(id int) *PropertyTypeUpdateOne {
	mutation := newPropertyTypeMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &PropertyTypeUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for PropertyType.
func (c *PropertyTypeClient) Delete() *PropertyTypeDelete {
	mutation := newPropertyTypeMutation(c.config, OpDelete)
	return &PropertyTypeDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *PropertyTypeClient) DeleteOne(pt *PropertyType) *PropertyTypeDeleteOne {
	return c.DeleteOneID(pt.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *PropertyTypeClient) DeleteOneID(id int) *PropertyTypeDeleteOne {
	builder := c.Delete().Where(propertytype.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &PropertyTypeDeleteOne{builder}
}

// Create returns a query builder for PropertyType.
func (c *PropertyTypeClient) Query() *PropertyTypeQuery {
	return &PropertyTypeQuery{config: c.config}
}

// Get returns a PropertyType entity by its id.
func (c *PropertyTypeClient) Get(ctx context.Context, id int) (*PropertyType, error) {
	return c.Query().Where(propertytype.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *PropertyTypeClient) GetX(ctx context.Context, id int) *PropertyType {
	pt, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return pt
}

// QueryProperties queries the properties edge of a PropertyType.
func (c *PropertyTypeClient) QueryProperties(pt *PropertyType) *PropertyQuery {
	query := &PropertyQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pt.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, id),
			sqlgraph.To(property.Table, property.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, propertytype.PropertiesTable, propertytype.PropertiesColumn),
		)
		fromV = sqlgraph.Neighbors(pt.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryLocationType queries the location_type edge of a PropertyType.
func (c *PropertyTypeClient) QueryLocationType(pt *PropertyType) *LocationTypeQuery {
	query := &LocationTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pt.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, id),
			sqlgraph.To(locationtype.Table, locationtype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, propertytype.LocationTypeTable, propertytype.LocationTypeColumn),
		)
		fromV = sqlgraph.Neighbors(pt.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEquipmentPortType queries the equipment_port_type edge of a PropertyType.
func (c *PropertyTypeClient) QueryEquipmentPortType(pt *PropertyType) *EquipmentPortTypeQuery {
	query := &EquipmentPortTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pt.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, id),
			sqlgraph.To(equipmentporttype.Table, equipmentporttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, propertytype.EquipmentPortTypeTable, propertytype.EquipmentPortTypeColumn),
		)
		fromV = sqlgraph.Neighbors(pt.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryLinkEquipmentPortType queries the link_equipment_port_type edge of a PropertyType.
func (c *PropertyTypeClient) QueryLinkEquipmentPortType(pt *PropertyType) *EquipmentPortTypeQuery {
	query := &EquipmentPortTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pt.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, id),
			sqlgraph.To(equipmentporttype.Table, equipmentporttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, propertytype.LinkEquipmentPortTypeTable, propertytype.LinkEquipmentPortTypeColumn),
		)
		fromV = sqlgraph.Neighbors(pt.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEquipmentType queries the equipment_type edge of a PropertyType.
func (c *PropertyTypeClient) QueryEquipmentType(pt *PropertyType) *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pt.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, id),
			sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, propertytype.EquipmentTypeTable, propertytype.EquipmentTypeColumn),
		)
		fromV = sqlgraph.Neighbors(pt.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryServiceType queries the service_type edge of a PropertyType.
func (c *PropertyTypeClient) QueryServiceType(pt *PropertyType) *ServiceTypeQuery {
	query := &ServiceTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pt.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, id),
			sqlgraph.To(servicetype.Table, servicetype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, propertytype.ServiceTypeTable, propertytype.ServiceTypeColumn),
		)
		fromV = sqlgraph.Neighbors(pt.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryWorkOrderType queries the work_order_type edge of a PropertyType.
func (c *PropertyTypeClient) QueryWorkOrderType(pt *PropertyType) *WorkOrderTypeQuery {
	query := &WorkOrderTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pt.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, id),
			sqlgraph.To(workordertype.Table, workordertype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, propertytype.WorkOrderTypeTable, propertytype.WorkOrderTypeColumn),
		)
		fromV = sqlgraph.Neighbors(pt.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryProjectType queries the project_type edge of a PropertyType.
func (c *PropertyTypeClient) QueryProjectType(pt *PropertyType) *ProjectTypeQuery {
	query := &ProjectTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pt.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, id),
			sqlgraph.To(projecttype.Table, projecttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, propertytype.ProjectTypeTable, propertytype.ProjectTypeColumn),
		)
		fromV = sqlgraph.Neighbors(pt.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *PropertyTypeClient) Hooks() []Hook {
	return c.hooks.PropertyType
}

// ReportFilterClient is a client for the ReportFilter schema.
type ReportFilterClient struct {
	config
}

// NewReportFilterClient returns a client for the ReportFilter from the given config.
func NewReportFilterClient(c config) *ReportFilterClient {
	return &ReportFilterClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `reportfilter.Hooks(f(g(h())))`.
func (c *ReportFilterClient) Use(hooks ...Hook) {
	c.hooks.ReportFilter = append(c.hooks.ReportFilter, hooks...)
}

// Create returns a create builder for ReportFilter.
func (c *ReportFilterClient) Create() *ReportFilterCreate {
	mutation := newReportFilterMutation(c.config, OpCreate)
	return &ReportFilterCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for ReportFilter.
func (c *ReportFilterClient) Update() *ReportFilterUpdate {
	mutation := newReportFilterMutation(c.config, OpUpdate)
	return &ReportFilterUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ReportFilterClient) UpdateOne(rf *ReportFilter) *ReportFilterUpdateOne {
	return c.UpdateOneID(rf.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *ReportFilterClient) UpdateOneID(id int) *ReportFilterUpdateOne {
	mutation := newReportFilterMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &ReportFilterUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for ReportFilter.
func (c *ReportFilterClient) Delete() *ReportFilterDelete {
	mutation := newReportFilterMutation(c.config, OpDelete)
	return &ReportFilterDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *ReportFilterClient) DeleteOne(rf *ReportFilter) *ReportFilterDeleteOne {
	return c.DeleteOneID(rf.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *ReportFilterClient) DeleteOneID(id int) *ReportFilterDeleteOne {
	builder := c.Delete().Where(reportfilter.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ReportFilterDeleteOne{builder}
}

// Create returns a query builder for ReportFilter.
func (c *ReportFilterClient) Query() *ReportFilterQuery {
	return &ReportFilterQuery{config: c.config}
}

// Get returns a ReportFilter entity by its id.
func (c *ReportFilterClient) Get(ctx context.Context, id int) (*ReportFilter, error) {
	return c.Query().Where(reportfilter.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ReportFilterClient) GetX(ctx context.Context, id int) *ReportFilter {
	rf, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return rf
}

// Hooks returns the client hooks.
func (c *ReportFilterClient) Hooks() []Hook {
	return c.hooks.ReportFilter
}

// ServiceClient is a client for the Service schema.
type ServiceClient struct {
	config
}

// NewServiceClient returns a client for the Service from the given config.
func NewServiceClient(c config) *ServiceClient {
	return &ServiceClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `service.Hooks(f(g(h())))`.
func (c *ServiceClient) Use(hooks ...Hook) {
	c.hooks.Service = append(c.hooks.Service, hooks...)
}

// Create returns a create builder for Service.
func (c *ServiceClient) Create() *ServiceCreate {
	mutation := newServiceMutation(c.config, OpCreate)
	return &ServiceCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for Service.
func (c *ServiceClient) Update() *ServiceUpdate {
	mutation := newServiceMutation(c.config, OpUpdate)
	return &ServiceUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ServiceClient) UpdateOne(s *Service) *ServiceUpdateOne {
	return c.UpdateOneID(s.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *ServiceClient) UpdateOneID(id int) *ServiceUpdateOne {
	mutation := newServiceMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &ServiceUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Service.
func (c *ServiceClient) Delete() *ServiceDelete {
	mutation := newServiceMutation(c.config, OpDelete)
	return &ServiceDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *ServiceClient) DeleteOne(s *Service) *ServiceDeleteOne {
	return c.DeleteOneID(s.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *ServiceClient) DeleteOneID(id int) *ServiceDeleteOne {
	builder := c.Delete().Where(service.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ServiceDeleteOne{builder}
}

// Create returns a query builder for Service.
func (c *ServiceClient) Query() *ServiceQuery {
	return &ServiceQuery{config: c.config}
}

// Get returns a Service entity by its id.
func (c *ServiceClient) Get(ctx context.Context, id int) (*Service, error) {
	return c.Query().Where(service.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ServiceClient) GetX(ctx context.Context, id int) *Service {
	s, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return s
}

// QueryType queries the type edge of a Service.
func (c *ServiceClient) QueryType(s *Service) *ServiceTypeQuery {
	query := &ServiceTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(service.Table, service.FieldID, id),
			sqlgraph.To(servicetype.Table, servicetype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, service.TypeTable, service.TypeColumn),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryDownstream queries the downstream edge of a Service.
func (c *ServiceClient) QueryDownstream(s *Service) *ServiceQuery {
	query := &ServiceQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(service.Table, service.FieldID, id),
			sqlgraph.To(service.Table, service.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, service.DownstreamTable, service.DownstreamPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryUpstream queries the upstream edge of a Service.
func (c *ServiceClient) QueryUpstream(s *Service) *ServiceQuery {
	query := &ServiceQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(service.Table, service.FieldID, id),
			sqlgraph.To(service.Table, service.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, service.UpstreamTable, service.UpstreamPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryProperties queries the properties edge of a Service.
func (c *ServiceClient) QueryProperties(s *Service) *PropertyQuery {
	query := &PropertyQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(service.Table, service.FieldID, id),
			sqlgraph.To(property.Table, property.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, service.PropertiesTable, service.PropertiesColumn),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryLinks queries the links edge of a Service.
func (c *ServiceClient) QueryLinks(s *Service) *LinkQuery {
	query := &LinkQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(service.Table, service.FieldID, id),
			sqlgraph.To(link.Table, link.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, service.LinksTable, service.LinksPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryCustomer queries the customer edge of a Service.
func (c *ServiceClient) QueryCustomer(s *Service) *CustomerQuery {
	query := &CustomerQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(service.Table, service.FieldID, id),
			sqlgraph.To(customer.Table, customer.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, service.CustomerTable, service.CustomerPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEndpoints queries the endpoints edge of a Service.
func (c *ServiceClient) QueryEndpoints(s *Service) *ServiceEndpointQuery {
	query := &ServiceEndpointQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(service.Table, service.FieldID, id),
			sqlgraph.To(serviceendpoint.Table, serviceendpoint.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, service.EndpointsTable, service.EndpointsColumn),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ServiceClient) Hooks() []Hook {
	return c.hooks.Service
}

// ServiceEndpointClient is a client for the ServiceEndpoint schema.
type ServiceEndpointClient struct {
	config
}

// NewServiceEndpointClient returns a client for the ServiceEndpoint from the given config.
func NewServiceEndpointClient(c config) *ServiceEndpointClient {
	return &ServiceEndpointClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `serviceendpoint.Hooks(f(g(h())))`.
func (c *ServiceEndpointClient) Use(hooks ...Hook) {
	c.hooks.ServiceEndpoint = append(c.hooks.ServiceEndpoint, hooks...)
}

// Create returns a create builder for ServiceEndpoint.
func (c *ServiceEndpointClient) Create() *ServiceEndpointCreate {
	mutation := newServiceEndpointMutation(c.config, OpCreate)
	return &ServiceEndpointCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for ServiceEndpoint.
func (c *ServiceEndpointClient) Update() *ServiceEndpointUpdate {
	mutation := newServiceEndpointMutation(c.config, OpUpdate)
	return &ServiceEndpointUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ServiceEndpointClient) UpdateOne(se *ServiceEndpoint) *ServiceEndpointUpdateOne {
	return c.UpdateOneID(se.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *ServiceEndpointClient) UpdateOneID(id int) *ServiceEndpointUpdateOne {
	mutation := newServiceEndpointMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &ServiceEndpointUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for ServiceEndpoint.
func (c *ServiceEndpointClient) Delete() *ServiceEndpointDelete {
	mutation := newServiceEndpointMutation(c.config, OpDelete)
	return &ServiceEndpointDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *ServiceEndpointClient) DeleteOne(se *ServiceEndpoint) *ServiceEndpointDeleteOne {
	return c.DeleteOneID(se.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *ServiceEndpointClient) DeleteOneID(id int) *ServiceEndpointDeleteOne {
	builder := c.Delete().Where(serviceendpoint.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ServiceEndpointDeleteOne{builder}
}

// Create returns a query builder for ServiceEndpoint.
func (c *ServiceEndpointClient) Query() *ServiceEndpointQuery {
	return &ServiceEndpointQuery{config: c.config}
}

// Get returns a ServiceEndpoint entity by its id.
func (c *ServiceEndpointClient) Get(ctx context.Context, id int) (*ServiceEndpoint, error) {
	return c.Query().Where(serviceendpoint.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ServiceEndpointClient) GetX(ctx context.Context, id int) *ServiceEndpoint {
	se, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return se
}

// QueryPort queries the port edge of a ServiceEndpoint.
func (c *ServiceEndpointClient) QueryPort(se *ServiceEndpoint) *EquipmentPortQuery {
	query := &EquipmentPortQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := se.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(serviceendpoint.Table, serviceendpoint.FieldID, id),
			sqlgraph.To(equipmentport.Table, equipmentport.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, serviceendpoint.PortTable, serviceendpoint.PortColumn),
		)
		fromV = sqlgraph.Neighbors(se.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEquipment queries the equipment edge of a ServiceEndpoint.
func (c *ServiceEndpointClient) QueryEquipment(se *ServiceEndpoint) *EquipmentQuery {
	query := &EquipmentQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := se.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(serviceendpoint.Table, serviceendpoint.FieldID, id),
			sqlgraph.To(equipment.Table, equipment.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, serviceendpoint.EquipmentTable, serviceendpoint.EquipmentColumn),
		)
		fromV = sqlgraph.Neighbors(se.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryService queries the service edge of a ServiceEndpoint.
func (c *ServiceEndpointClient) QueryService(se *ServiceEndpoint) *ServiceQuery {
	query := &ServiceQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := se.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(serviceendpoint.Table, serviceendpoint.FieldID, id),
			sqlgraph.To(service.Table, service.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, serviceendpoint.ServiceTable, serviceendpoint.ServiceColumn),
		)
		fromV = sqlgraph.Neighbors(se.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryDefinition queries the definition edge of a ServiceEndpoint.
func (c *ServiceEndpointClient) QueryDefinition(se *ServiceEndpoint) *ServiceEndpointDefinitionQuery {
	query := &ServiceEndpointDefinitionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := se.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(serviceendpoint.Table, serviceendpoint.FieldID, id),
			sqlgraph.To(serviceendpointdefinition.Table, serviceendpointdefinition.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, serviceendpoint.DefinitionTable, serviceendpoint.DefinitionColumn),
		)
		fromV = sqlgraph.Neighbors(se.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ServiceEndpointClient) Hooks() []Hook {
	return c.hooks.ServiceEndpoint
}

// ServiceEndpointDefinitionClient is a client for the ServiceEndpointDefinition schema.
type ServiceEndpointDefinitionClient struct {
	config
}

// NewServiceEndpointDefinitionClient returns a client for the ServiceEndpointDefinition from the given config.
func NewServiceEndpointDefinitionClient(c config) *ServiceEndpointDefinitionClient {
	return &ServiceEndpointDefinitionClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `serviceendpointdefinition.Hooks(f(g(h())))`.
func (c *ServiceEndpointDefinitionClient) Use(hooks ...Hook) {
	c.hooks.ServiceEndpointDefinition = append(c.hooks.ServiceEndpointDefinition, hooks...)
}

// Create returns a create builder for ServiceEndpointDefinition.
func (c *ServiceEndpointDefinitionClient) Create() *ServiceEndpointDefinitionCreate {
	mutation := newServiceEndpointDefinitionMutation(c.config, OpCreate)
	return &ServiceEndpointDefinitionCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for ServiceEndpointDefinition.
func (c *ServiceEndpointDefinitionClient) Update() *ServiceEndpointDefinitionUpdate {
	mutation := newServiceEndpointDefinitionMutation(c.config, OpUpdate)
	return &ServiceEndpointDefinitionUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ServiceEndpointDefinitionClient) UpdateOne(sed *ServiceEndpointDefinition) *ServiceEndpointDefinitionUpdateOne {
	return c.UpdateOneID(sed.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *ServiceEndpointDefinitionClient) UpdateOneID(id int) *ServiceEndpointDefinitionUpdateOne {
	mutation := newServiceEndpointDefinitionMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &ServiceEndpointDefinitionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for ServiceEndpointDefinition.
func (c *ServiceEndpointDefinitionClient) Delete() *ServiceEndpointDefinitionDelete {
	mutation := newServiceEndpointDefinitionMutation(c.config, OpDelete)
	return &ServiceEndpointDefinitionDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *ServiceEndpointDefinitionClient) DeleteOne(sed *ServiceEndpointDefinition) *ServiceEndpointDefinitionDeleteOne {
	return c.DeleteOneID(sed.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *ServiceEndpointDefinitionClient) DeleteOneID(id int) *ServiceEndpointDefinitionDeleteOne {
	builder := c.Delete().Where(serviceendpointdefinition.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ServiceEndpointDefinitionDeleteOne{builder}
}

// Create returns a query builder for ServiceEndpointDefinition.
func (c *ServiceEndpointDefinitionClient) Query() *ServiceEndpointDefinitionQuery {
	return &ServiceEndpointDefinitionQuery{config: c.config}
}

// Get returns a ServiceEndpointDefinition entity by its id.
func (c *ServiceEndpointDefinitionClient) Get(ctx context.Context, id int) (*ServiceEndpointDefinition, error) {
	return c.Query().Where(serviceendpointdefinition.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ServiceEndpointDefinitionClient) GetX(ctx context.Context, id int) *ServiceEndpointDefinition {
	sed, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return sed
}

// QueryEndpoints queries the endpoints edge of a ServiceEndpointDefinition.
func (c *ServiceEndpointDefinitionClient) QueryEndpoints(sed *ServiceEndpointDefinition) *ServiceEndpointQuery {
	query := &ServiceEndpointQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := sed.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(serviceendpointdefinition.Table, serviceendpointdefinition.FieldID, id),
			sqlgraph.To(serviceendpoint.Table, serviceendpoint.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, serviceendpointdefinition.EndpointsTable, serviceendpointdefinition.EndpointsColumn),
		)
		fromV = sqlgraph.Neighbors(sed.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryServiceType queries the service_type edge of a ServiceEndpointDefinition.
func (c *ServiceEndpointDefinitionClient) QueryServiceType(sed *ServiceEndpointDefinition) *ServiceTypeQuery {
	query := &ServiceTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := sed.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(serviceendpointdefinition.Table, serviceendpointdefinition.FieldID, id),
			sqlgraph.To(servicetype.Table, servicetype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, serviceendpointdefinition.ServiceTypeTable, serviceendpointdefinition.ServiceTypeColumn),
		)
		fromV = sqlgraph.Neighbors(sed.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEquipmentType queries the equipment_type edge of a ServiceEndpointDefinition.
func (c *ServiceEndpointDefinitionClient) QueryEquipmentType(sed *ServiceEndpointDefinition) *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := sed.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(serviceendpointdefinition.Table, serviceendpointdefinition.FieldID, id),
			sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, serviceendpointdefinition.EquipmentTypeTable, serviceendpointdefinition.EquipmentTypeColumn),
		)
		fromV = sqlgraph.Neighbors(sed.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ServiceEndpointDefinitionClient) Hooks() []Hook {
	return c.hooks.ServiceEndpointDefinition
}

// ServiceTypeClient is a client for the ServiceType schema.
type ServiceTypeClient struct {
	config
}

// NewServiceTypeClient returns a client for the ServiceType from the given config.
func NewServiceTypeClient(c config) *ServiceTypeClient {
	return &ServiceTypeClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `servicetype.Hooks(f(g(h())))`.
func (c *ServiceTypeClient) Use(hooks ...Hook) {
	c.hooks.ServiceType = append(c.hooks.ServiceType, hooks...)
}

// Create returns a create builder for ServiceType.
func (c *ServiceTypeClient) Create() *ServiceTypeCreate {
	mutation := newServiceTypeMutation(c.config, OpCreate)
	return &ServiceTypeCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for ServiceType.
func (c *ServiceTypeClient) Update() *ServiceTypeUpdate {
	mutation := newServiceTypeMutation(c.config, OpUpdate)
	return &ServiceTypeUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ServiceTypeClient) UpdateOne(st *ServiceType) *ServiceTypeUpdateOne {
	return c.UpdateOneID(st.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *ServiceTypeClient) UpdateOneID(id int) *ServiceTypeUpdateOne {
	mutation := newServiceTypeMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &ServiceTypeUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for ServiceType.
func (c *ServiceTypeClient) Delete() *ServiceTypeDelete {
	mutation := newServiceTypeMutation(c.config, OpDelete)
	return &ServiceTypeDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *ServiceTypeClient) DeleteOne(st *ServiceType) *ServiceTypeDeleteOne {
	return c.DeleteOneID(st.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *ServiceTypeClient) DeleteOneID(id int) *ServiceTypeDeleteOne {
	builder := c.Delete().Where(servicetype.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ServiceTypeDeleteOne{builder}
}

// Create returns a query builder for ServiceType.
func (c *ServiceTypeClient) Query() *ServiceTypeQuery {
	return &ServiceTypeQuery{config: c.config}
}

// Get returns a ServiceType entity by its id.
func (c *ServiceTypeClient) Get(ctx context.Context, id int) (*ServiceType, error) {
	return c.Query().Where(servicetype.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ServiceTypeClient) GetX(ctx context.Context, id int) *ServiceType {
	st, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return st
}

// QueryServices queries the services edge of a ServiceType.
func (c *ServiceTypeClient) QueryServices(st *ServiceType) *ServiceQuery {
	query := &ServiceQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := st.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(servicetype.Table, servicetype.FieldID, id),
			sqlgraph.To(service.Table, service.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, servicetype.ServicesTable, servicetype.ServicesColumn),
		)
		fromV = sqlgraph.Neighbors(st.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryPropertyTypes queries the property_types edge of a ServiceType.
func (c *ServiceTypeClient) QueryPropertyTypes(st *ServiceType) *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := st.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(servicetype.Table, servicetype.FieldID, id),
			sqlgraph.To(propertytype.Table, propertytype.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, servicetype.PropertyTypesTable, servicetype.PropertyTypesColumn),
		)
		fromV = sqlgraph.Neighbors(st.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEndpointDefinitions queries the endpoint_definitions edge of a ServiceType.
func (c *ServiceTypeClient) QueryEndpointDefinitions(st *ServiceType) *ServiceEndpointDefinitionQuery {
	query := &ServiceEndpointDefinitionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := st.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(servicetype.Table, servicetype.FieldID, id),
			sqlgraph.To(serviceendpointdefinition.Table, serviceendpointdefinition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, servicetype.EndpointDefinitionsTable, servicetype.EndpointDefinitionsColumn),
		)
		fromV = sqlgraph.Neighbors(st.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ServiceTypeClient) Hooks() []Hook {
	return c.hooks.ServiceType
}

// SurveyClient is a client for the Survey schema.
type SurveyClient struct {
	config
}

// NewSurveyClient returns a client for the Survey from the given config.
func NewSurveyClient(c config) *SurveyClient {
	return &SurveyClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `survey.Hooks(f(g(h())))`.
func (c *SurveyClient) Use(hooks ...Hook) {
	c.hooks.Survey = append(c.hooks.Survey, hooks...)
}

// Create returns a create builder for Survey.
func (c *SurveyClient) Create() *SurveyCreate {
	mutation := newSurveyMutation(c.config, OpCreate)
	return &SurveyCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for Survey.
func (c *SurveyClient) Update() *SurveyUpdate {
	mutation := newSurveyMutation(c.config, OpUpdate)
	return &SurveyUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SurveyClient) UpdateOne(s *Survey) *SurveyUpdateOne {
	return c.UpdateOneID(s.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *SurveyClient) UpdateOneID(id int) *SurveyUpdateOne {
	mutation := newSurveyMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &SurveyUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Survey.
func (c *SurveyClient) Delete() *SurveyDelete {
	mutation := newSurveyMutation(c.config, OpDelete)
	return &SurveyDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *SurveyClient) DeleteOne(s *Survey) *SurveyDeleteOne {
	return c.DeleteOneID(s.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *SurveyClient) DeleteOneID(id int) *SurveyDeleteOne {
	builder := c.Delete().Where(survey.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SurveyDeleteOne{builder}
}

// Create returns a query builder for Survey.
func (c *SurveyClient) Query() *SurveyQuery {
	return &SurveyQuery{config: c.config}
}

// Get returns a Survey entity by its id.
func (c *SurveyClient) Get(ctx context.Context, id int) (*Survey, error) {
	return c.Query().Where(survey.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SurveyClient) GetX(ctx context.Context, id int) *Survey {
	s, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return s
}

// QueryLocation queries the location edge of a Survey.
func (c *SurveyClient) QueryLocation(s *Survey) *LocationQuery {
	query := &LocationQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(survey.Table, survey.FieldID, id),
			sqlgraph.To(location.Table, location.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, survey.LocationTable, survey.LocationColumn),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QuerySourceFile queries the source_file edge of a Survey.
func (c *SurveyClient) QuerySourceFile(s *Survey) *FileQuery {
	query := &FileQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(survey.Table, survey.FieldID, id),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, survey.SourceFileTable, survey.SourceFileColumn),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryQuestions queries the questions edge of a Survey.
func (c *SurveyClient) QueryQuestions(s *Survey) *SurveyQuestionQuery {
	query := &SurveyQuestionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(survey.Table, survey.FieldID, id),
			sqlgraph.To(surveyquestion.Table, surveyquestion.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, survey.QuestionsTable, survey.QuestionsColumn),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *SurveyClient) Hooks() []Hook {
	return c.hooks.Survey
}

// SurveyCellScanClient is a client for the SurveyCellScan schema.
type SurveyCellScanClient struct {
	config
}

// NewSurveyCellScanClient returns a client for the SurveyCellScan from the given config.
func NewSurveyCellScanClient(c config) *SurveyCellScanClient {
	return &SurveyCellScanClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `surveycellscan.Hooks(f(g(h())))`.
func (c *SurveyCellScanClient) Use(hooks ...Hook) {
	c.hooks.SurveyCellScan = append(c.hooks.SurveyCellScan, hooks...)
}

// Create returns a create builder for SurveyCellScan.
func (c *SurveyCellScanClient) Create() *SurveyCellScanCreate {
	mutation := newSurveyCellScanMutation(c.config, OpCreate)
	return &SurveyCellScanCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for SurveyCellScan.
func (c *SurveyCellScanClient) Update() *SurveyCellScanUpdate {
	mutation := newSurveyCellScanMutation(c.config, OpUpdate)
	return &SurveyCellScanUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SurveyCellScanClient) UpdateOne(scs *SurveyCellScan) *SurveyCellScanUpdateOne {
	return c.UpdateOneID(scs.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *SurveyCellScanClient) UpdateOneID(id int) *SurveyCellScanUpdateOne {
	mutation := newSurveyCellScanMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &SurveyCellScanUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for SurveyCellScan.
func (c *SurveyCellScanClient) Delete() *SurveyCellScanDelete {
	mutation := newSurveyCellScanMutation(c.config, OpDelete)
	return &SurveyCellScanDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *SurveyCellScanClient) DeleteOne(scs *SurveyCellScan) *SurveyCellScanDeleteOne {
	return c.DeleteOneID(scs.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *SurveyCellScanClient) DeleteOneID(id int) *SurveyCellScanDeleteOne {
	builder := c.Delete().Where(surveycellscan.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SurveyCellScanDeleteOne{builder}
}

// Create returns a query builder for SurveyCellScan.
func (c *SurveyCellScanClient) Query() *SurveyCellScanQuery {
	return &SurveyCellScanQuery{config: c.config}
}

// Get returns a SurveyCellScan entity by its id.
func (c *SurveyCellScanClient) Get(ctx context.Context, id int) (*SurveyCellScan, error) {
	return c.Query().Where(surveycellscan.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SurveyCellScanClient) GetX(ctx context.Context, id int) *SurveyCellScan {
	scs, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return scs
}

// QueryChecklistItem queries the checklist_item edge of a SurveyCellScan.
func (c *SurveyCellScanClient) QueryChecklistItem(scs *SurveyCellScan) *CheckListItemQuery {
	query := &CheckListItemQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := scs.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(surveycellscan.Table, surveycellscan.FieldID, id),
			sqlgraph.To(checklistitem.Table, checklistitem.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, surveycellscan.ChecklistItemTable, surveycellscan.ChecklistItemColumn),
		)
		fromV = sqlgraph.Neighbors(scs.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QuerySurveyQuestion queries the survey_question edge of a SurveyCellScan.
func (c *SurveyCellScanClient) QuerySurveyQuestion(scs *SurveyCellScan) *SurveyQuestionQuery {
	query := &SurveyQuestionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := scs.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(surveycellscan.Table, surveycellscan.FieldID, id),
			sqlgraph.To(surveyquestion.Table, surveyquestion.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, surveycellscan.SurveyQuestionTable, surveycellscan.SurveyQuestionColumn),
		)
		fromV = sqlgraph.Neighbors(scs.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryLocation queries the location edge of a SurveyCellScan.
func (c *SurveyCellScanClient) QueryLocation(scs *SurveyCellScan) *LocationQuery {
	query := &LocationQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := scs.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(surveycellscan.Table, surveycellscan.FieldID, id),
			sqlgraph.To(location.Table, location.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, surveycellscan.LocationTable, surveycellscan.LocationColumn),
		)
		fromV = sqlgraph.Neighbors(scs.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *SurveyCellScanClient) Hooks() []Hook {
	return c.hooks.SurveyCellScan
}

// SurveyQuestionClient is a client for the SurveyQuestion schema.
type SurveyQuestionClient struct {
	config
}

// NewSurveyQuestionClient returns a client for the SurveyQuestion from the given config.
func NewSurveyQuestionClient(c config) *SurveyQuestionClient {
	return &SurveyQuestionClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `surveyquestion.Hooks(f(g(h())))`.
func (c *SurveyQuestionClient) Use(hooks ...Hook) {
	c.hooks.SurveyQuestion = append(c.hooks.SurveyQuestion, hooks...)
}

// Create returns a create builder for SurveyQuestion.
func (c *SurveyQuestionClient) Create() *SurveyQuestionCreate {
	mutation := newSurveyQuestionMutation(c.config, OpCreate)
	return &SurveyQuestionCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for SurveyQuestion.
func (c *SurveyQuestionClient) Update() *SurveyQuestionUpdate {
	mutation := newSurveyQuestionMutation(c.config, OpUpdate)
	return &SurveyQuestionUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SurveyQuestionClient) UpdateOne(sq *SurveyQuestion) *SurveyQuestionUpdateOne {
	return c.UpdateOneID(sq.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *SurveyQuestionClient) UpdateOneID(id int) *SurveyQuestionUpdateOne {
	mutation := newSurveyQuestionMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &SurveyQuestionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for SurveyQuestion.
func (c *SurveyQuestionClient) Delete() *SurveyQuestionDelete {
	mutation := newSurveyQuestionMutation(c.config, OpDelete)
	return &SurveyQuestionDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *SurveyQuestionClient) DeleteOne(sq *SurveyQuestion) *SurveyQuestionDeleteOne {
	return c.DeleteOneID(sq.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *SurveyQuestionClient) DeleteOneID(id int) *SurveyQuestionDeleteOne {
	builder := c.Delete().Where(surveyquestion.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SurveyQuestionDeleteOne{builder}
}

// Create returns a query builder for SurveyQuestion.
func (c *SurveyQuestionClient) Query() *SurveyQuestionQuery {
	return &SurveyQuestionQuery{config: c.config}
}

// Get returns a SurveyQuestion entity by its id.
func (c *SurveyQuestionClient) Get(ctx context.Context, id int) (*SurveyQuestion, error) {
	return c.Query().Where(surveyquestion.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SurveyQuestionClient) GetX(ctx context.Context, id int) *SurveyQuestion {
	sq, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return sq
}

// QuerySurvey queries the survey edge of a SurveyQuestion.
func (c *SurveyQuestionClient) QuerySurvey(sq *SurveyQuestion) *SurveyQuery {
	query := &SurveyQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := sq.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, id),
			sqlgraph.To(survey.Table, survey.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, surveyquestion.SurveyTable, surveyquestion.SurveyColumn),
		)
		fromV = sqlgraph.Neighbors(sq.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryWifiScan queries the wifi_scan edge of a SurveyQuestion.
func (c *SurveyQuestionClient) QueryWifiScan(sq *SurveyQuestion) *SurveyWiFiScanQuery {
	query := &SurveyWiFiScanQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := sq.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, id),
			sqlgraph.To(surveywifiscan.Table, surveywifiscan.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, surveyquestion.WifiScanTable, surveyquestion.WifiScanColumn),
		)
		fromV = sqlgraph.Neighbors(sq.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryCellScan queries the cell_scan edge of a SurveyQuestion.
func (c *SurveyQuestionClient) QueryCellScan(sq *SurveyQuestion) *SurveyCellScanQuery {
	query := &SurveyCellScanQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := sq.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, id),
			sqlgraph.To(surveycellscan.Table, surveycellscan.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, surveyquestion.CellScanTable, surveyquestion.CellScanColumn),
		)
		fromV = sqlgraph.Neighbors(sq.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryPhotoData queries the photo_data edge of a SurveyQuestion.
func (c *SurveyQuestionClient) QueryPhotoData(sq *SurveyQuestion) *FileQuery {
	query := &FileQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := sq.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, id),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, surveyquestion.PhotoDataTable, surveyquestion.PhotoDataColumn),
		)
		fromV = sqlgraph.Neighbors(sq.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryImages queries the images edge of a SurveyQuestion.
func (c *SurveyQuestionClient) QueryImages(sq *SurveyQuestion) *FileQuery {
	query := &FileQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := sq.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, id),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, surveyquestion.ImagesTable, surveyquestion.ImagesColumn),
		)
		fromV = sqlgraph.Neighbors(sq.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *SurveyQuestionClient) Hooks() []Hook {
	return c.hooks.SurveyQuestion
}

// SurveyTemplateCategoryClient is a client for the SurveyTemplateCategory schema.
type SurveyTemplateCategoryClient struct {
	config
}

// NewSurveyTemplateCategoryClient returns a client for the SurveyTemplateCategory from the given config.
func NewSurveyTemplateCategoryClient(c config) *SurveyTemplateCategoryClient {
	return &SurveyTemplateCategoryClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `surveytemplatecategory.Hooks(f(g(h())))`.
func (c *SurveyTemplateCategoryClient) Use(hooks ...Hook) {
	c.hooks.SurveyTemplateCategory = append(c.hooks.SurveyTemplateCategory, hooks...)
}

// Create returns a create builder for SurveyTemplateCategory.
func (c *SurveyTemplateCategoryClient) Create() *SurveyTemplateCategoryCreate {
	mutation := newSurveyTemplateCategoryMutation(c.config, OpCreate)
	return &SurveyTemplateCategoryCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for SurveyTemplateCategory.
func (c *SurveyTemplateCategoryClient) Update() *SurveyTemplateCategoryUpdate {
	mutation := newSurveyTemplateCategoryMutation(c.config, OpUpdate)
	return &SurveyTemplateCategoryUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SurveyTemplateCategoryClient) UpdateOne(stc *SurveyTemplateCategory) *SurveyTemplateCategoryUpdateOne {
	return c.UpdateOneID(stc.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *SurveyTemplateCategoryClient) UpdateOneID(id int) *SurveyTemplateCategoryUpdateOne {
	mutation := newSurveyTemplateCategoryMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &SurveyTemplateCategoryUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for SurveyTemplateCategory.
func (c *SurveyTemplateCategoryClient) Delete() *SurveyTemplateCategoryDelete {
	mutation := newSurveyTemplateCategoryMutation(c.config, OpDelete)
	return &SurveyTemplateCategoryDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *SurveyTemplateCategoryClient) DeleteOne(stc *SurveyTemplateCategory) *SurveyTemplateCategoryDeleteOne {
	return c.DeleteOneID(stc.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *SurveyTemplateCategoryClient) DeleteOneID(id int) *SurveyTemplateCategoryDeleteOne {
	builder := c.Delete().Where(surveytemplatecategory.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SurveyTemplateCategoryDeleteOne{builder}
}

// Create returns a query builder for SurveyTemplateCategory.
func (c *SurveyTemplateCategoryClient) Query() *SurveyTemplateCategoryQuery {
	return &SurveyTemplateCategoryQuery{config: c.config}
}

// Get returns a SurveyTemplateCategory entity by its id.
func (c *SurveyTemplateCategoryClient) Get(ctx context.Context, id int) (*SurveyTemplateCategory, error) {
	return c.Query().Where(surveytemplatecategory.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SurveyTemplateCategoryClient) GetX(ctx context.Context, id int) *SurveyTemplateCategory {
	stc, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return stc
}

// QuerySurveyTemplateQuestions queries the survey_template_questions edge of a SurveyTemplateCategory.
func (c *SurveyTemplateCategoryClient) QuerySurveyTemplateQuestions(stc *SurveyTemplateCategory) *SurveyTemplateQuestionQuery {
	query := &SurveyTemplateQuestionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := stc.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(surveytemplatecategory.Table, surveytemplatecategory.FieldID, id),
			sqlgraph.To(surveytemplatequestion.Table, surveytemplatequestion.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, surveytemplatecategory.SurveyTemplateQuestionsTable, surveytemplatecategory.SurveyTemplateQuestionsColumn),
		)
		fromV = sqlgraph.Neighbors(stc.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *SurveyTemplateCategoryClient) Hooks() []Hook {
	return c.hooks.SurveyTemplateCategory
}

// SurveyTemplateQuestionClient is a client for the SurveyTemplateQuestion schema.
type SurveyTemplateQuestionClient struct {
	config
}

// NewSurveyTemplateQuestionClient returns a client for the SurveyTemplateQuestion from the given config.
func NewSurveyTemplateQuestionClient(c config) *SurveyTemplateQuestionClient {
	return &SurveyTemplateQuestionClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `surveytemplatequestion.Hooks(f(g(h())))`.
func (c *SurveyTemplateQuestionClient) Use(hooks ...Hook) {
	c.hooks.SurveyTemplateQuestion = append(c.hooks.SurveyTemplateQuestion, hooks...)
}

// Create returns a create builder for SurveyTemplateQuestion.
func (c *SurveyTemplateQuestionClient) Create() *SurveyTemplateQuestionCreate {
	mutation := newSurveyTemplateQuestionMutation(c.config, OpCreate)
	return &SurveyTemplateQuestionCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for SurveyTemplateQuestion.
func (c *SurveyTemplateQuestionClient) Update() *SurveyTemplateQuestionUpdate {
	mutation := newSurveyTemplateQuestionMutation(c.config, OpUpdate)
	return &SurveyTemplateQuestionUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SurveyTemplateQuestionClient) UpdateOne(stq *SurveyTemplateQuestion) *SurveyTemplateQuestionUpdateOne {
	return c.UpdateOneID(stq.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *SurveyTemplateQuestionClient) UpdateOneID(id int) *SurveyTemplateQuestionUpdateOne {
	mutation := newSurveyTemplateQuestionMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &SurveyTemplateQuestionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for SurveyTemplateQuestion.
func (c *SurveyTemplateQuestionClient) Delete() *SurveyTemplateQuestionDelete {
	mutation := newSurveyTemplateQuestionMutation(c.config, OpDelete)
	return &SurveyTemplateQuestionDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *SurveyTemplateQuestionClient) DeleteOne(stq *SurveyTemplateQuestion) *SurveyTemplateQuestionDeleteOne {
	return c.DeleteOneID(stq.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *SurveyTemplateQuestionClient) DeleteOneID(id int) *SurveyTemplateQuestionDeleteOne {
	builder := c.Delete().Where(surveytemplatequestion.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SurveyTemplateQuestionDeleteOne{builder}
}

// Create returns a query builder for SurveyTemplateQuestion.
func (c *SurveyTemplateQuestionClient) Query() *SurveyTemplateQuestionQuery {
	return &SurveyTemplateQuestionQuery{config: c.config}
}

// Get returns a SurveyTemplateQuestion entity by its id.
func (c *SurveyTemplateQuestionClient) Get(ctx context.Context, id int) (*SurveyTemplateQuestion, error) {
	return c.Query().Where(surveytemplatequestion.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SurveyTemplateQuestionClient) GetX(ctx context.Context, id int) *SurveyTemplateQuestion {
	stq, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return stq
}

// QueryCategory queries the category edge of a SurveyTemplateQuestion.
func (c *SurveyTemplateQuestionClient) QueryCategory(stq *SurveyTemplateQuestion) *SurveyTemplateCategoryQuery {
	query := &SurveyTemplateCategoryQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := stq.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(surveytemplatequestion.Table, surveytemplatequestion.FieldID, id),
			sqlgraph.To(surveytemplatecategory.Table, surveytemplatecategory.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, surveytemplatequestion.CategoryTable, surveytemplatequestion.CategoryColumn),
		)
		fromV = sqlgraph.Neighbors(stq.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *SurveyTemplateQuestionClient) Hooks() []Hook {
	return c.hooks.SurveyTemplateQuestion
}

// SurveyWiFiScanClient is a client for the SurveyWiFiScan schema.
type SurveyWiFiScanClient struct {
	config
}

// NewSurveyWiFiScanClient returns a client for the SurveyWiFiScan from the given config.
func NewSurveyWiFiScanClient(c config) *SurveyWiFiScanClient {
	return &SurveyWiFiScanClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `surveywifiscan.Hooks(f(g(h())))`.
func (c *SurveyWiFiScanClient) Use(hooks ...Hook) {
	c.hooks.SurveyWiFiScan = append(c.hooks.SurveyWiFiScan, hooks...)
}

// Create returns a create builder for SurveyWiFiScan.
func (c *SurveyWiFiScanClient) Create() *SurveyWiFiScanCreate {
	mutation := newSurveyWiFiScanMutation(c.config, OpCreate)
	return &SurveyWiFiScanCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for SurveyWiFiScan.
func (c *SurveyWiFiScanClient) Update() *SurveyWiFiScanUpdate {
	mutation := newSurveyWiFiScanMutation(c.config, OpUpdate)
	return &SurveyWiFiScanUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SurveyWiFiScanClient) UpdateOne(swfs *SurveyWiFiScan) *SurveyWiFiScanUpdateOne {
	return c.UpdateOneID(swfs.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *SurveyWiFiScanClient) UpdateOneID(id int) *SurveyWiFiScanUpdateOne {
	mutation := newSurveyWiFiScanMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &SurveyWiFiScanUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for SurveyWiFiScan.
func (c *SurveyWiFiScanClient) Delete() *SurveyWiFiScanDelete {
	mutation := newSurveyWiFiScanMutation(c.config, OpDelete)
	return &SurveyWiFiScanDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *SurveyWiFiScanClient) DeleteOne(swfs *SurveyWiFiScan) *SurveyWiFiScanDeleteOne {
	return c.DeleteOneID(swfs.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *SurveyWiFiScanClient) DeleteOneID(id int) *SurveyWiFiScanDeleteOne {
	builder := c.Delete().Where(surveywifiscan.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SurveyWiFiScanDeleteOne{builder}
}

// Create returns a query builder for SurveyWiFiScan.
func (c *SurveyWiFiScanClient) Query() *SurveyWiFiScanQuery {
	return &SurveyWiFiScanQuery{config: c.config}
}

// Get returns a SurveyWiFiScan entity by its id.
func (c *SurveyWiFiScanClient) Get(ctx context.Context, id int) (*SurveyWiFiScan, error) {
	return c.Query().Where(surveywifiscan.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SurveyWiFiScanClient) GetX(ctx context.Context, id int) *SurveyWiFiScan {
	swfs, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return swfs
}

// QueryChecklistItem queries the checklist_item edge of a SurveyWiFiScan.
func (c *SurveyWiFiScanClient) QueryChecklistItem(swfs *SurveyWiFiScan) *CheckListItemQuery {
	query := &CheckListItemQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := swfs.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(surveywifiscan.Table, surveywifiscan.FieldID, id),
			sqlgraph.To(checklistitem.Table, checklistitem.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, surveywifiscan.ChecklistItemTable, surveywifiscan.ChecklistItemColumn),
		)
		fromV = sqlgraph.Neighbors(swfs.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QuerySurveyQuestion queries the survey_question edge of a SurveyWiFiScan.
func (c *SurveyWiFiScanClient) QuerySurveyQuestion(swfs *SurveyWiFiScan) *SurveyQuestionQuery {
	query := &SurveyQuestionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := swfs.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(surveywifiscan.Table, surveywifiscan.FieldID, id),
			sqlgraph.To(surveyquestion.Table, surveyquestion.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, surveywifiscan.SurveyQuestionTable, surveywifiscan.SurveyQuestionColumn),
		)
		fromV = sqlgraph.Neighbors(swfs.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryLocation queries the location edge of a SurveyWiFiScan.
func (c *SurveyWiFiScanClient) QueryLocation(swfs *SurveyWiFiScan) *LocationQuery {
	query := &LocationQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := swfs.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(surveywifiscan.Table, surveywifiscan.FieldID, id),
			sqlgraph.To(location.Table, location.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, surveywifiscan.LocationTable, surveywifiscan.LocationColumn),
		)
		fromV = sqlgraph.Neighbors(swfs.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *SurveyWiFiScanClient) Hooks() []Hook {
	return c.hooks.SurveyWiFiScan
}

// TechnicianClient is a client for the Technician schema.
type TechnicianClient struct {
	config
}

// NewTechnicianClient returns a client for the Technician from the given config.
func NewTechnicianClient(c config) *TechnicianClient {
	return &TechnicianClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `technician.Hooks(f(g(h())))`.
func (c *TechnicianClient) Use(hooks ...Hook) {
	c.hooks.Technician = append(c.hooks.Technician, hooks...)
}

// Create returns a create builder for Technician.
func (c *TechnicianClient) Create() *TechnicianCreate {
	mutation := newTechnicianMutation(c.config, OpCreate)
	return &TechnicianCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for Technician.
func (c *TechnicianClient) Update() *TechnicianUpdate {
	mutation := newTechnicianMutation(c.config, OpUpdate)
	return &TechnicianUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *TechnicianClient) UpdateOne(t *Technician) *TechnicianUpdateOne {
	return c.UpdateOneID(t.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *TechnicianClient) UpdateOneID(id int) *TechnicianUpdateOne {
	mutation := newTechnicianMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &TechnicianUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Technician.
func (c *TechnicianClient) Delete() *TechnicianDelete {
	mutation := newTechnicianMutation(c.config, OpDelete)
	return &TechnicianDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *TechnicianClient) DeleteOne(t *Technician) *TechnicianDeleteOne {
	return c.DeleteOneID(t.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *TechnicianClient) DeleteOneID(id int) *TechnicianDeleteOne {
	builder := c.Delete().Where(technician.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &TechnicianDeleteOne{builder}
}

// Create returns a query builder for Technician.
func (c *TechnicianClient) Query() *TechnicianQuery {
	return &TechnicianQuery{config: c.config}
}

// Get returns a Technician entity by its id.
func (c *TechnicianClient) Get(ctx context.Context, id int) (*Technician, error) {
	return c.Query().Where(technician.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *TechnicianClient) GetX(ctx context.Context, id int) *Technician {
	t, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return t
}

// QueryWorkOrders queries the work_orders edge of a Technician.
func (c *TechnicianClient) QueryWorkOrders(t *Technician) *WorkOrderQuery {
	query := &WorkOrderQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := t.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(technician.Table, technician.FieldID, id),
			sqlgraph.To(workorder.Table, workorder.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, technician.WorkOrdersTable, technician.WorkOrdersColumn),
		)
		fromV = sqlgraph.Neighbors(t.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *TechnicianClient) Hooks() []Hook {
	return c.hooks.Technician
}

// UserClient is a client for the User schema.
type UserClient struct {
	config
}

// NewUserClient returns a client for the User from the given config.
func NewUserClient(c config) *UserClient {
	return &UserClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `user.Hooks(f(g(h())))`.
func (c *UserClient) Use(hooks ...Hook) {
	c.hooks.User = append(c.hooks.User, hooks...)
}

// Create returns a create builder for User.
func (c *UserClient) Create() *UserCreate {
	mutation := newUserMutation(c.config, OpCreate)
	return &UserCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for User.
func (c *UserClient) Update() *UserUpdate {
	mutation := newUserMutation(c.config, OpUpdate)
	return &UserUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *UserClient) UpdateOne(u *User) *UserUpdateOne {
	return c.UpdateOneID(u.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *UserClient) UpdateOneID(id int) *UserUpdateOne {
	mutation := newUserMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for User.
func (c *UserClient) Delete() *UserDelete {
	mutation := newUserMutation(c.config, OpDelete)
	return &UserDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *UserClient) DeleteOne(u *User) *UserDeleteOne {
	return c.DeleteOneID(u.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *UserClient) DeleteOneID(id int) *UserDeleteOne {
	builder := c.Delete().Where(user.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &UserDeleteOne{builder}
}

// Create returns a query builder for User.
func (c *UserClient) Query() *UserQuery {
	return &UserQuery{config: c.config}
}

// Get returns a User entity by its id.
func (c *UserClient) Get(ctx context.Context, id int) (*User, error) {
	return c.Query().Where(user.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *UserClient) GetX(ctx context.Context, id int) *User {
	u, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return u
}

// QueryProfilePhoto queries the profile_photo edge of a User.
func (c *UserClient) QueryProfilePhoto(u *User) *FileQuery {
	query := &FileQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := u.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(user.Table, user.FieldID, id),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, user.ProfilePhotoTable, user.ProfilePhotoColumn),
		)
		fromV = sqlgraph.Neighbors(u.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryGroups queries the groups edge of a User.
func (c *UserClient) QueryGroups(u *User) *UsersGroupQuery {
	query := &UsersGroupQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := u.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(user.Table, user.FieldID, id),
			sqlgraph.To(usersgroup.Table, usersgroup.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, user.GroupsTable, user.GroupsPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(u.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *UserClient) Hooks() []Hook {
	hooks := c.hooks.User
	return append(hooks[:len(hooks):len(hooks)], user.Hooks[:]...)
}

// UsersGroupClient is a client for the UsersGroup schema.
type UsersGroupClient struct {
	config
}

// NewUsersGroupClient returns a client for the UsersGroup from the given config.
func NewUsersGroupClient(c config) *UsersGroupClient {
	return &UsersGroupClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `usersgroup.Hooks(f(g(h())))`.
func (c *UsersGroupClient) Use(hooks ...Hook) {
	c.hooks.UsersGroup = append(c.hooks.UsersGroup, hooks...)
}

// Create returns a create builder for UsersGroup.
func (c *UsersGroupClient) Create() *UsersGroupCreate {
	mutation := newUsersGroupMutation(c.config, OpCreate)
	return &UsersGroupCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for UsersGroup.
func (c *UsersGroupClient) Update() *UsersGroupUpdate {
	mutation := newUsersGroupMutation(c.config, OpUpdate)
	return &UsersGroupUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *UsersGroupClient) UpdateOne(ug *UsersGroup) *UsersGroupUpdateOne {
	return c.UpdateOneID(ug.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *UsersGroupClient) UpdateOneID(id int) *UsersGroupUpdateOne {
	mutation := newUsersGroupMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &UsersGroupUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for UsersGroup.
func (c *UsersGroupClient) Delete() *UsersGroupDelete {
	mutation := newUsersGroupMutation(c.config, OpDelete)
	return &UsersGroupDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *UsersGroupClient) DeleteOne(ug *UsersGroup) *UsersGroupDeleteOne {
	return c.DeleteOneID(ug.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *UsersGroupClient) DeleteOneID(id int) *UsersGroupDeleteOne {
	builder := c.Delete().Where(usersgroup.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &UsersGroupDeleteOne{builder}
}

// Create returns a query builder for UsersGroup.
func (c *UsersGroupClient) Query() *UsersGroupQuery {
	return &UsersGroupQuery{config: c.config}
}

// Get returns a UsersGroup entity by its id.
func (c *UsersGroupClient) Get(ctx context.Context, id int) (*UsersGroup, error) {
	return c.Query().Where(usersgroup.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *UsersGroupClient) GetX(ctx context.Context, id int) *UsersGroup {
	ug, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return ug
}

// QueryMembers queries the members edge of a UsersGroup.
func (c *UsersGroupClient) QueryMembers(ug *UsersGroup) *UserQuery {
	query := &UserQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := ug.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(usersgroup.Table, usersgroup.FieldID, id),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, usersgroup.MembersTable, usersgroup.MembersPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(ug.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryPolicies queries the policies edge of a UsersGroup.
func (c *UsersGroupClient) QueryPolicies(ug *UsersGroup) *PermissionsPolicyQuery {
	query := &PermissionsPolicyQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := ug.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(usersgroup.Table, usersgroup.FieldID, id),
			sqlgraph.To(permissionspolicy.Table, permissionspolicy.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, usersgroup.PoliciesTable, usersgroup.PoliciesPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(ug.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *UsersGroupClient) Hooks() []Hook {
	return c.hooks.UsersGroup
}

// WorkOrderClient is a client for the WorkOrder schema.
type WorkOrderClient struct {
	config
}

// NewWorkOrderClient returns a client for the WorkOrder from the given config.
func NewWorkOrderClient(c config) *WorkOrderClient {
	return &WorkOrderClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `workorder.Hooks(f(g(h())))`.
func (c *WorkOrderClient) Use(hooks ...Hook) {
	c.hooks.WorkOrder = append(c.hooks.WorkOrder, hooks...)
}

// Create returns a create builder for WorkOrder.
func (c *WorkOrderClient) Create() *WorkOrderCreate {
	mutation := newWorkOrderMutation(c.config, OpCreate)
	return &WorkOrderCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for WorkOrder.
func (c *WorkOrderClient) Update() *WorkOrderUpdate {
	mutation := newWorkOrderMutation(c.config, OpUpdate)
	return &WorkOrderUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *WorkOrderClient) UpdateOne(wo *WorkOrder) *WorkOrderUpdateOne {
	return c.UpdateOneID(wo.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *WorkOrderClient) UpdateOneID(id int) *WorkOrderUpdateOne {
	mutation := newWorkOrderMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &WorkOrderUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for WorkOrder.
func (c *WorkOrderClient) Delete() *WorkOrderDelete {
	mutation := newWorkOrderMutation(c.config, OpDelete)
	return &WorkOrderDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *WorkOrderClient) DeleteOne(wo *WorkOrder) *WorkOrderDeleteOne {
	return c.DeleteOneID(wo.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *WorkOrderClient) DeleteOneID(id int) *WorkOrderDeleteOne {
	builder := c.Delete().Where(workorder.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &WorkOrderDeleteOne{builder}
}

// Create returns a query builder for WorkOrder.
func (c *WorkOrderClient) Query() *WorkOrderQuery {
	return &WorkOrderQuery{config: c.config}
}

// Get returns a WorkOrder entity by its id.
func (c *WorkOrderClient) Get(ctx context.Context, id int) (*WorkOrder, error) {
	return c.Query().Where(workorder.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *WorkOrderClient) GetX(ctx context.Context, id int) *WorkOrder {
	wo, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return wo
}

// QueryType queries the type edge of a WorkOrder.
func (c *WorkOrderClient) QueryType(wo *WorkOrder) *WorkOrderTypeQuery {
	query := &WorkOrderTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wo.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, id),
			sqlgraph.To(workordertype.Table, workordertype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, workorder.TypeTable, workorder.TypeColumn),
		)
		fromV = sqlgraph.Neighbors(wo.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEquipment queries the equipment edge of a WorkOrder.
func (c *WorkOrderClient) QueryEquipment(wo *WorkOrder) *EquipmentQuery {
	query := &EquipmentQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wo.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, id),
			sqlgraph.To(equipment.Table, equipment.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, workorder.EquipmentTable, workorder.EquipmentColumn),
		)
		fromV = sqlgraph.Neighbors(wo.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryLinks queries the links edge of a WorkOrder.
func (c *WorkOrderClient) QueryLinks(wo *WorkOrder) *LinkQuery {
	query := &LinkQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wo.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, id),
			sqlgraph.To(link.Table, link.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, workorder.LinksTable, workorder.LinksColumn),
		)
		fromV = sqlgraph.Neighbors(wo.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryFiles queries the files edge of a WorkOrder.
func (c *WorkOrderClient) QueryFiles(wo *WorkOrder) *FileQuery {
	query := &FileQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wo.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, id),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, workorder.FilesTable, workorder.FilesColumn),
		)
		fromV = sqlgraph.Neighbors(wo.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryHyperlinks queries the hyperlinks edge of a WorkOrder.
func (c *WorkOrderClient) QueryHyperlinks(wo *WorkOrder) *HyperlinkQuery {
	query := &HyperlinkQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wo.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, id),
			sqlgraph.To(hyperlink.Table, hyperlink.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, workorder.HyperlinksTable, workorder.HyperlinksColumn),
		)
		fromV = sqlgraph.Neighbors(wo.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryLocation queries the location edge of a WorkOrder.
func (c *WorkOrderClient) QueryLocation(wo *WorkOrder) *LocationQuery {
	query := &LocationQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wo.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, id),
			sqlgraph.To(location.Table, location.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, workorder.LocationTable, workorder.LocationColumn),
		)
		fromV = sqlgraph.Neighbors(wo.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryComments queries the comments edge of a WorkOrder.
func (c *WorkOrderClient) QueryComments(wo *WorkOrder) *CommentQuery {
	query := &CommentQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wo.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, id),
			sqlgraph.To(comment.Table, comment.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, workorder.CommentsTable, workorder.CommentsColumn),
		)
		fromV = sqlgraph.Neighbors(wo.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryProperties queries the properties edge of a WorkOrder.
func (c *WorkOrderClient) QueryProperties(wo *WorkOrder) *PropertyQuery {
	query := &PropertyQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wo.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, id),
			sqlgraph.To(property.Table, property.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, workorder.PropertiesTable, workorder.PropertiesColumn),
		)
		fromV = sqlgraph.Neighbors(wo.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryCheckListCategories queries the check_list_categories edge of a WorkOrder.
func (c *WorkOrderClient) QueryCheckListCategories(wo *WorkOrder) *CheckListCategoryQuery {
	query := &CheckListCategoryQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wo.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, id),
			sqlgraph.To(checklistcategory.Table, checklistcategory.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, workorder.CheckListCategoriesTable, workorder.CheckListCategoriesColumn),
		)
		fromV = sqlgraph.Neighbors(wo.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryCheckListItems queries the check_list_items edge of a WorkOrder.
func (c *WorkOrderClient) QueryCheckListItems(wo *WorkOrder) *CheckListItemQuery {
	query := &CheckListItemQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wo.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, id),
			sqlgraph.To(checklistitem.Table, checklistitem.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, workorder.CheckListItemsTable, workorder.CheckListItemsColumn),
		)
		fromV = sqlgraph.Neighbors(wo.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryTechnician queries the technician edge of a WorkOrder.
func (c *WorkOrderClient) QueryTechnician(wo *WorkOrder) *TechnicianQuery {
	query := &TechnicianQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wo.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, id),
			sqlgraph.To(technician.Table, technician.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, workorder.TechnicianTable, workorder.TechnicianColumn),
		)
		fromV = sqlgraph.Neighbors(wo.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryProject queries the project edge of a WorkOrder.
func (c *WorkOrderClient) QueryProject(wo *WorkOrder) *ProjectQuery {
	query := &ProjectQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wo.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, id),
			sqlgraph.To(project.Table, project.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, workorder.ProjectTable, workorder.ProjectColumn),
		)
		fromV = sqlgraph.Neighbors(wo.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryOwner queries the owner edge of a WorkOrder.
func (c *WorkOrderClient) QueryOwner(wo *WorkOrder) *UserQuery {
	query := &UserQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wo.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, id),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, workorder.OwnerTable, workorder.OwnerColumn),
		)
		fromV = sqlgraph.Neighbors(wo.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryAssignee queries the assignee edge of a WorkOrder.
func (c *WorkOrderClient) QueryAssignee(wo *WorkOrder) *UserQuery {
	query := &UserQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wo.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workorder.Table, workorder.FieldID, id),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, workorder.AssigneeTable, workorder.AssigneeColumn),
		)
		fromV = sqlgraph.Neighbors(wo.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *WorkOrderClient) Hooks() []Hook {
	return c.hooks.WorkOrder
}

// WorkOrderDefinitionClient is a client for the WorkOrderDefinition schema.
type WorkOrderDefinitionClient struct {
	config
}

// NewWorkOrderDefinitionClient returns a client for the WorkOrderDefinition from the given config.
func NewWorkOrderDefinitionClient(c config) *WorkOrderDefinitionClient {
	return &WorkOrderDefinitionClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `workorderdefinition.Hooks(f(g(h())))`.
func (c *WorkOrderDefinitionClient) Use(hooks ...Hook) {
	c.hooks.WorkOrderDefinition = append(c.hooks.WorkOrderDefinition, hooks...)
}

// Create returns a create builder for WorkOrderDefinition.
func (c *WorkOrderDefinitionClient) Create() *WorkOrderDefinitionCreate {
	mutation := newWorkOrderDefinitionMutation(c.config, OpCreate)
	return &WorkOrderDefinitionCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for WorkOrderDefinition.
func (c *WorkOrderDefinitionClient) Update() *WorkOrderDefinitionUpdate {
	mutation := newWorkOrderDefinitionMutation(c.config, OpUpdate)
	return &WorkOrderDefinitionUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *WorkOrderDefinitionClient) UpdateOne(wod *WorkOrderDefinition) *WorkOrderDefinitionUpdateOne {
	return c.UpdateOneID(wod.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *WorkOrderDefinitionClient) UpdateOneID(id int) *WorkOrderDefinitionUpdateOne {
	mutation := newWorkOrderDefinitionMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &WorkOrderDefinitionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for WorkOrderDefinition.
func (c *WorkOrderDefinitionClient) Delete() *WorkOrderDefinitionDelete {
	mutation := newWorkOrderDefinitionMutation(c.config, OpDelete)
	return &WorkOrderDefinitionDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *WorkOrderDefinitionClient) DeleteOne(wod *WorkOrderDefinition) *WorkOrderDefinitionDeleteOne {
	return c.DeleteOneID(wod.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *WorkOrderDefinitionClient) DeleteOneID(id int) *WorkOrderDefinitionDeleteOne {
	builder := c.Delete().Where(workorderdefinition.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &WorkOrderDefinitionDeleteOne{builder}
}

// Create returns a query builder for WorkOrderDefinition.
func (c *WorkOrderDefinitionClient) Query() *WorkOrderDefinitionQuery {
	return &WorkOrderDefinitionQuery{config: c.config}
}

// Get returns a WorkOrderDefinition entity by its id.
func (c *WorkOrderDefinitionClient) Get(ctx context.Context, id int) (*WorkOrderDefinition, error) {
	return c.Query().Where(workorderdefinition.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *WorkOrderDefinitionClient) GetX(ctx context.Context, id int) *WorkOrderDefinition {
	wod, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return wod
}

// QueryType queries the type edge of a WorkOrderDefinition.
func (c *WorkOrderDefinitionClient) QueryType(wod *WorkOrderDefinition) *WorkOrderTypeQuery {
	query := &WorkOrderTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wod.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workorderdefinition.Table, workorderdefinition.FieldID, id),
			sqlgraph.To(workordertype.Table, workordertype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, workorderdefinition.TypeTable, workorderdefinition.TypeColumn),
		)
		fromV = sqlgraph.Neighbors(wod.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryProjectType queries the project_type edge of a WorkOrderDefinition.
func (c *WorkOrderDefinitionClient) QueryProjectType(wod *WorkOrderDefinition) *ProjectTypeQuery {
	query := &ProjectTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wod.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workorderdefinition.Table, workorderdefinition.FieldID, id),
			sqlgraph.To(projecttype.Table, projecttype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, workorderdefinition.ProjectTypeTable, workorderdefinition.ProjectTypeColumn),
		)
		fromV = sqlgraph.Neighbors(wod.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *WorkOrderDefinitionClient) Hooks() []Hook {
	return c.hooks.WorkOrderDefinition
}

// WorkOrderTypeClient is a client for the WorkOrderType schema.
type WorkOrderTypeClient struct {
	config
}

// NewWorkOrderTypeClient returns a client for the WorkOrderType from the given config.
func NewWorkOrderTypeClient(c config) *WorkOrderTypeClient {
	return &WorkOrderTypeClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `workordertype.Hooks(f(g(h())))`.
func (c *WorkOrderTypeClient) Use(hooks ...Hook) {
	c.hooks.WorkOrderType = append(c.hooks.WorkOrderType, hooks...)
}

// Create returns a create builder for WorkOrderType.
func (c *WorkOrderTypeClient) Create() *WorkOrderTypeCreate {
	mutation := newWorkOrderTypeMutation(c.config, OpCreate)
	return &WorkOrderTypeCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for WorkOrderType.
func (c *WorkOrderTypeClient) Update() *WorkOrderTypeUpdate {
	mutation := newWorkOrderTypeMutation(c.config, OpUpdate)
	return &WorkOrderTypeUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *WorkOrderTypeClient) UpdateOne(wot *WorkOrderType) *WorkOrderTypeUpdateOne {
	return c.UpdateOneID(wot.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *WorkOrderTypeClient) UpdateOneID(id int) *WorkOrderTypeUpdateOne {
	mutation := newWorkOrderTypeMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &WorkOrderTypeUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for WorkOrderType.
func (c *WorkOrderTypeClient) Delete() *WorkOrderTypeDelete {
	mutation := newWorkOrderTypeMutation(c.config, OpDelete)
	return &WorkOrderTypeDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *WorkOrderTypeClient) DeleteOne(wot *WorkOrderType) *WorkOrderTypeDeleteOne {
	return c.DeleteOneID(wot.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *WorkOrderTypeClient) DeleteOneID(id int) *WorkOrderTypeDeleteOne {
	builder := c.Delete().Where(workordertype.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &WorkOrderTypeDeleteOne{builder}
}

// Create returns a query builder for WorkOrderType.
func (c *WorkOrderTypeClient) Query() *WorkOrderTypeQuery {
	return &WorkOrderTypeQuery{config: c.config}
}

// Get returns a WorkOrderType entity by its id.
func (c *WorkOrderTypeClient) Get(ctx context.Context, id int) (*WorkOrderType, error) {
	return c.Query().Where(workordertype.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *WorkOrderTypeClient) GetX(ctx context.Context, id int) *WorkOrderType {
	wot, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return wot
}

// QueryWorkOrders queries the work_orders edge of a WorkOrderType.
func (c *WorkOrderTypeClient) QueryWorkOrders(wot *WorkOrderType) *WorkOrderQuery {
	query := &WorkOrderQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wot.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workordertype.Table, workordertype.FieldID, id),
			sqlgraph.To(workorder.Table, workorder.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, workordertype.WorkOrdersTable, workordertype.WorkOrdersColumn),
		)
		fromV = sqlgraph.Neighbors(wot.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryPropertyTypes queries the property_types edge of a WorkOrderType.
func (c *WorkOrderTypeClient) QueryPropertyTypes(wot *WorkOrderType) *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wot.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workordertype.Table, workordertype.FieldID, id),
			sqlgraph.To(propertytype.Table, propertytype.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, workordertype.PropertyTypesTable, workordertype.PropertyTypesColumn),
		)
		fromV = sqlgraph.Neighbors(wot.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryDefinitions queries the definitions edge of a WorkOrderType.
func (c *WorkOrderTypeClient) QueryDefinitions(wot *WorkOrderType) *WorkOrderDefinitionQuery {
	query := &WorkOrderDefinitionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wot.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workordertype.Table, workordertype.FieldID, id),
			sqlgraph.To(workorderdefinition.Table, workorderdefinition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, workordertype.DefinitionsTable, workordertype.DefinitionsColumn),
		)
		fromV = sqlgraph.Neighbors(wot.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryCheckListCategories queries the check_list_categories edge of a WorkOrderType.
func (c *WorkOrderTypeClient) QueryCheckListCategories(wot *WorkOrderType) *CheckListCategoryQuery {
	query := &CheckListCategoryQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wot.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workordertype.Table, workordertype.FieldID, id),
			sqlgraph.To(checklistcategory.Table, checklistcategory.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, workordertype.CheckListCategoriesTable, workordertype.CheckListCategoriesColumn),
		)
		fromV = sqlgraph.Neighbors(wot.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryCheckListDefinitions queries the check_list_definitions edge of a WorkOrderType.
func (c *WorkOrderTypeClient) QueryCheckListDefinitions(wot *WorkOrderType) *CheckListItemDefinitionQuery {
	query := &CheckListItemDefinitionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wot.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(workordertype.Table, workordertype.FieldID, id),
			sqlgraph.To(checklistitemdefinition.Table, checklistitemdefinition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, workordertype.CheckListDefinitionsTable, workordertype.CheckListDefinitionsColumn),
		)
		fromV = sqlgraph.Neighbors(wot.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *WorkOrderTypeClient) Hooks() []Hook {
	return c.hooks.WorkOrderType
}
