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
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatecategory"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatequestion"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
	"github.com/facebookincubator/symphony/graph/ent/technician"
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
	c := config{log: log.Println}
	c.options(opts...)
	return &Client{
		config:                      c,
		Schema:                      migrate.NewSchema(c.driver),
		ActionsRule:                 NewActionsRuleClient(c),
		CheckListCategory:           NewCheckListCategoryClient(c),
		CheckListItem:               NewCheckListItemClient(c),
		CheckListItemDefinition:     NewCheckListItemDefinitionClient(c),
		Comment:                     NewCommentClient(c),
		Customer:                    NewCustomerClient(c),
		Equipment:                   NewEquipmentClient(c),
		EquipmentCategory:           NewEquipmentCategoryClient(c),
		EquipmentPort:               NewEquipmentPortClient(c),
		EquipmentPortDefinition:     NewEquipmentPortDefinitionClient(c),
		EquipmentPortType:           NewEquipmentPortTypeClient(c),
		EquipmentPosition:           NewEquipmentPositionClient(c),
		EquipmentPositionDefinition: NewEquipmentPositionDefinitionClient(c),
		EquipmentType:               NewEquipmentTypeClient(c),
		File:                        NewFileClient(c),
		FloorPlan:                   NewFloorPlanClient(c),
		FloorPlanReferencePoint:     NewFloorPlanReferencePointClient(c),
		FloorPlanScale:              NewFloorPlanScaleClient(c),
		Hyperlink:                   NewHyperlinkClient(c),
		Link:                        NewLinkClient(c),
		Location:                    NewLocationClient(c),
		LocationType:                NewLocationTypeClient(c),
		Project:                     NewProjectClient(c),
		ProjectType:                 NewProjectTypeClient(c),
		Property:                    NewPropertyClient(c),
		PropertyType:                NewPropertyTypeClient(c),
		Service:                     NewServiceClient(c),
		ServiceEndpoint:             NewServiceEndpointClient(c),
		ServiceType:                 NewServiceTypeClient(c),
		Survey:                      NewSurveyClient(c),
		SurveyCellScan:              NewSurveyCellScanClient(c),
		SurveyQuestion:              NewSurveyQuestionClient(c),
		SurveyTemplateCategory:      NewSurveyTemplateCategoryClient(c),
		SurveyTemplateQuestion:      NewSurveyTemplateQuestionClient(c),
		SurveyWiFiScan:              NewSurveyWiFiScanClient(c),
		Technician:                  NewTechnicianClient(c),
		WorkOrder:                   NewWorkOrderClient(c),
		WorkOrderDefinition:         NewWorkOrderDefinitionClient(c),
		WorkOrderType:               NewWorkOrderTypeClient(c),
	}
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
	cfg := config{driver: tx, log: c.log, debug: c.debug}
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
	cfg := config{driver: dialect.Debug(c.driver, c.log), log: c.log, debug: true}
	return &Client{
		config:                      cfg,
		Schema:                      migrate.NewSchema(cfg.driver),
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

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// ActionsRuleClient is a client for the ActionsRule schema.
type ActionsRuleClient struct {
	config
}

// NewActionsRuleClient returns a client for the ActionsRule from the given config.
func NewActionsRuleClient(c config) *ActionsRuleClient {
	return &ActionsRuleClient{config: c}
}

// Create returns a create builder for ActionsRule.
func (c *ActionsRuleClient) Create() *ActionsRuleCreate {
	return &ActionsRuleCreate{config: c.config}
}

// Update returns an update builder for ActionsRule.
func (c *ActionsRuleClient) Update() *ActionsRuleUpdate {
	return &ActionsRuleUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *ActionsRuleClient) UpdateOne(ar *ActionsRule) *ActionsRuleUpdateOne {
	return c.UpdateOneID(ar.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *ActionsRuleClient) UpdateOneID(id string) *ActionsRuleUpdateOne {
	return &ActionsRuleUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for ActionsRule.
func (c *ActionsRuleClient) Delete() *ActionsRuleDelete {
	return &ActionsRuleDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *ActionsRuleClient) DeleteOne(ar *ActionsRule) *ActionsRuleDeleteOne {
	return c.DeleteOneID(ar.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *ActionsRuleClient) DeleteOneID(id string) *ActionsRuleDeleteOne {
	return &ActionsRuleDeleteOne{c.Delete().Where(actionsrule.ID(id))}
}

// Create returns a query builder for ActionsRule.
func (c *ActionsRuleClient) Query() *ActionsRuleQuery {
	return &ActionsRuleQuery{config: c.config}
}

// Get returns a ActionsRule entity by its id.
func (c *ActionsRuleClient) Get(ctx context.Context, id string) (*ActionsRule, error) {
	return c.Query().Where(actionsrule.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ActionsRuleClient) GetX(ctx context.Context, id string) *ActionsRule {
	ar, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return ar
}

// CheckListCategoryClient is a client for the CheckListCategory schema.
type CheckListCategoryClient struct {
	config
}

// NewCheckListCategoryClient returns a client for the CheckListCategory from the given config.
func NewCheckListCategoryClient(c config) *CheckListCategoryClient {
	return &CheckListCategoryClient{config: c}
}

// Create returns a create builder for CheckListCategory.
func (c *CheckListCategoryClient) Create() *CheckListCategoryCreate {
	return &CheckListCategoryCreate{config: c.config}
}

// Update returns an update builder for CheckListCategory.
func (c *CheckListCategoryClient) Update() *CheckListCategoryUpdate {
	return &CheckListCategoryUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *CheckListCategoryClient) UpdateOne(clc *CheckListCategory) *CheckListCategoryUpdateOne {
	return c.UpdateOneID(clc.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *CheckListCategoryClient) UpdateOneID(id string) *CheckListCategoryUpdateOne {
	return &CheckListCategoryUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for CheckListCategory.
func (c *CheckListCategoryClient) Delete() *CheckListCategoryDelete {
	return &CheckListCategoryDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *CheckListCategoryClient) DeleteOne(clc *CheckListCategory) *CheckListCategoryDeleteOne {
	return c.DeleteOneID(clc.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *CheckListCategoryClient) DeleteOneID(id string) *CheckListCategoryDeleteOne {
	return &CheckListCategoryDeleteOne{c.Delete().Where(checklistcategory.ID(id))}
}

// Create returns a query builder for CheckListCategory.
func (c *CheckListCategoryClient) Query() *CheckListCategoryQuery {
	return &CheckListCategoryQuery{config: c.config}
}

// Get returns a CheckListCategory entity by its id.
func (c *CheckListCategoryClient) Get(ctx context.Context, id string) (*CheckListCategory, error) {
	return c.Query().Where(checklistcategory.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *CheckListCategoryClient) GetX(ctx context.Context, id string) *CheckListCategory {
	clc, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return clc
}

// QueryCheckListItems queries the check_list_items edge of a CheckListCategory.
func (c *CheckListCategoryClient) QueryCheckListItems(clc *CheckListCategory) *CheckListItemQuery {
	query := &CheckListItemQuery{config: c.config}
	id := clc.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(checklistcategory.Table, checklistcategory.FieldID, id),
		sqlgraph.To(checklistitem.Table, checklistitem.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, checklistcategory.CheckListItemsTable, checklistcategory.CheckListItemsColumn),
	)
	query.sql = sqlgraph.Neighbors(clc.driver.Dialect(), step)

	return query
}

// CheckListItemClient is a client for the CheckListItem schema.
type CheckListItemClient struct {
	config
}

// NewCheckListItemClient returns a client for the CheckListItem from the given config.
func NewCheckListItemClient(c config) *CheckListItemClient {
	return &CheckListItemClient{config: c}
}

// Create returns a create builder for CheckListItem.
func (c *CheckListItemClient) Create() *CheckListItemCreate {
	return &CheckListItemCreate{config: c.config}
}

// Update returns an update builder for CheckListItem.
func (c *CheckListItemClient) Update() *CheckListItemUpdate {
	return &CheckListItemUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *CheckListItemClient) UpdateOne(cli *CheckListItem) *CheckListItemUpdateOne {
	return c.UpdateOneID(cli.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *CheckListItemClient) UpdateOneID(id string) *CheckListItemUpdateOne {
	return &CheckListItemUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for CheckListItem.
func (c *CheckListItemClient) Delete() *CheckListItemDelete {
	return &CheckListItemDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *CheckListItemClient) DeleteOne(cli *CheckListItem) *CheckListItemDeleteOne {
	return c.DeleteOneID(cli.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *CheckListItemClient) DeleteOneID(id string) *CheckListItemDeleteOne {
	return &CheckListItemDeleteOne{c.Delete().Where(checklistitem.ID(id))}
}

// Create returns a query builder for CheckListItem.
func (c *CheckListItemClient) Query() *CheckListItemQuery {
	return &CheckListItemQuery{config: c.config}
}

// Get returns a CheckListItem entity by its id.
func (c *CheckListItemClient) Get(ctx context.Context, id string) (*CheckListItem, error) {
	return c.Query().Where(checklistitem.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *CheckListItemClient) GetX(ctx context.Context, id string) *CheckListItem {
	cli, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return cli
}

// QueryWorkOrder queries the work_order edge of a CheckListItem.
func (c *CheckListItemClient) QueryWorkOrder(cli *CheckListItem) *WorkOrderQuery {
	query := &WorkOrderQuery{config: c.config}
	id := cli.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(checklistitem.Table, checklistitem.FieldID, id),
		sqlgraph.To(workorder.Table, workorder.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, checklistitem.WorkOrderTable, checklistitem.WorkOrderColumn),
	)
	query.sql = sqlgraph.Neighbors(cli.driver.Dialect(), step)

	return query
}

// CheckListItemDefinitionClient is a client for the CheckListItemDefinition schema.
type CheckListItemDefinitionClient struct {
	config
}

// NewCheckListItemDefinitionClient returns a client for the CheckListItemDefinition from the given config.
func NewCheckListItemDefinitionClient(c config) *CheckListItemDefinitionClient {
	return &CheckListItemDefinitionClient{config: c}
}

// Create returns a create builder for CheckListItemDefinition.
func (c *CheckListItemDefinitionClient) Create() *CheckListItemDefinitionCreate {
	return &CheckListItemDefinitionCreate{config: c.config}
}

// Update returns an update builder for CheckListItemDefinition.
func (c *CheckListItemDefinitionClient) Update() *CheckListItemDefinitionUpdate {
	return &CheckListItemDefinitionUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *CheckListItemDefinitionClient) UpdateOne(clid *CheckListItemDefinition) *CheckListItemDefinitionUpdateOne {
	return c.UpdateOneID(clid.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *CheckListItemDefinitionClient) UpdateOneID(id string) *CheckListItemDefinitionUpdateOne {
	return &CheckListItemDefinitionUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for CheckListItemDefinition.
func (c *CheckListItemDefinitionClient) Delete() *CheckListItemDefinitionDelete {
	return &CheckListItemDefinitionDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *CheckListItemDefinitionClient) DeleteOne(clid *CheckListItemDefinition) *CheckListItemDefinitionDeleteOne {
	return c.DeleteOneID(clid.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *CheckListItemDefinitionClient) DeleteOneID(id string) *CheckListItemDefinitionDeleteOne {
	return &CheckListItemDefinitionDeleteOne{c.Delete().Where(checklistitemdefinition.ID(id))}
}

// Create returns a query builder for CheckListItemDefinition.
func (c *CheckListItemDefinitionClient) Query() *CheckListItemDefinitionQuery {
	return &CheckListItemDefinitionQuery{config: c.config}
}

// Get returns a CheckListItemDefinition entity by its id.
func (c *CheckListItemDefinitionClient) Get(ctx context.Context, id string) (*CheckListItemDefinition, error) {
	return c.Query().Where(checklistitemdefinition.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *CheckListItemDefinitionClient) GetX(ctx context.Context, id string) *CheckListItemDefinition {
	clid, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return clid
}

// QueryWorkOrderType queries the work_order_type edge of a CheckListItemDefinition.
func (c *CheckListItemDefinitionClient) QueryWorkOrderType(clid *CheckListItemDefinition) *WorkOrderTypeQuery {
	query := &WorkOrderTypeQuery{config: c.config}
	id := clid.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(checklistitemdefinition.Table, checklistitemdefinition.FieldID, id),
		sqlgraph.To(workordertype.Table, workordertype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, checklistitemdefinition.WorkOrderTypeTable, checklistitemdefinition.WorkOrderTypeColumn),
	)
	query.sql = sqlgraph.Neighbors(clid.driver.Dialect(), step)

	return query
}

// CommentClient is a client for the Comment schema.
type CommentClient struct {
	config
}

// NewCommentClient returns a client for the Comment from the given config.
func NewCommentClient(c config) *CommentClient {
	return &CommentClient{config: c}
}

// Create returns a create builder for Comment.
func (c *CommentClient) Create() *CommentCreate {
	return &CommentCreate{config: c.config}
}

// Update returns an update builder for Comment.
func (c *CommentClient) Update() *CommentUpdate {
	return &CommentUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *CommentClient) UpdateOne(co *Comment) *CommentUpdateOne {
	return c.UpdateOneID(co.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *CommentClient) UpdateOneID(id string) *CommentUpdateOne {
	return &CommentUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for Comment.
func (c *CommentClient) Delete() *CommentDelete {
	return &CommentDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *CommentClient) DeleteOne(co *Comment) *CommentDeleteOne {
	return c.DeleteOneID(co.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *CommentClient) DeleteOneID(id string) *CommentDeleteOne {
	return &CommentDeleteOne{c.Delete().Where(comment.ID(id))}
}

// Create returns a query builder for Comment.
func (c *CommentClient) Query() *CommentQuery {
	return &CommentQuery{config: c.config}
}

// Get returns a Comment entity by its id.
func (c *CommentClient) Get(ctx context.Context, id string) (*Comment, error) {
	return c.Query().Where(comment.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *CommentClient) GetX(ctx context.Context, id string) *Comment {
	co, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return co
}

// CustomerClient is a client for the Customer schema.
type CustomerClient struct {
	config
}

// NewCustomerClient returns a client for the Customer from the given config.
func NewCustomerClient(c config) *CustomerClient {
	return &CustomerClient{config: c}
}

// Create returns a create builder for Customer.
func (c *CustomerClient) Create() *CustomerCreate {
	return &CustomerCreate{config: c.config}
}

// Update returns an update builder for Customer.
func (c *CustomerClient) Update() *CustomerUpdate {
	return &CustomerUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *CustomerClient) UpdateOne(cu *Customer) *CustomerUpdateOne {
	return c.UpdateOneID(cu.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *CustomerClient) UpdateOneID(id string) *CustomerUpdateOne {
	return &CustomerUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for Customer.
func (c *CustomerClient) Delete() *CustomerDelete {
	return &CustomerDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *CustomerClient) DeleteOne(cu *Customer) *CustomerDeleteOne {
	return c.DeleteOneID(cu.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *CustomerClient) DeleteOneID(id string) *CustomerDeleteOne {
	return &CustomerDeleteOne{c.Delete().Where(customer.ID(id))}
}

// Create returns a query builder for Customer.
func (c *CustomerClient) Query() *CustomerQuery {
	return &CustomerQuery{config: c.config}
}

// Get returns a Customer entity by its id.
func (c *CustomerClient) Get(ctx context.Context, id string) (*Customer, error) {
	return c.Query().Where(customer.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *CustomerClient) GetX(ctx context.Context, id string) *Customer {
	cu, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return cu
}

// QueryServices queries the services edge of a Customer.
func (c *CustomerClient) QueryServices(cu *Customer) *ServiceQuery {
	query := &ServiceQuery{config: c.config}
	id := cu.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(customer.Table, customer.FieldID, id),
		sqlgraph.To(service.Table, service.FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, customer.ServicesTable, customer.ServicesPrimaryKey...),
	)
	query.sql = sqlgraph.Neighbors(cu.driver.Dialect(), step)

	return query
}

// EquipmentClient is a client for the Equipment schema.
type EquipmentClient struct {
	config
}

// NewEquipmentClient returns a client for the Equipment from the given config.
func NewEquipmentClient(c config) *EquipmentClient {
	return &EquipmentClient{config: c}
}

// Create returns a create builder for Equipment.
func (c *EquipmentClient) Create() *EquipmentCreate {
	return &EquipmentCreate{config: c.config}
}

// Update returns an update builder for Equipment.
func (c *EquipmentClient) Update() *EquipmentUpdate {
	return &EquipmentUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *EquipmentClient) UpdateOne(e *Equipment) *EquipmentUpdateOne {
	return c.UpdateOneID(e.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *EquipmentClient) UpdateOneID(id string) *EquipmentUpdateOne {
	return &EquipmentUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for Equipment.
func (c *EquipmentClient) Delete() *EquipmentDelete {
	return &EquipmentDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *EquipmentClient) DeleteOne(e *Equipment) *EquipmentDeleteOne {
	return c.DeleteOneID(e.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *EquipmentClient) DeleteOneID(id string) *EquipmentDeleteOne {
	return &EquipmentDeleteOne{c.Delete().Where(equipment.ID(id))}
}

// Create returns a query builder for Equipment.
func (c *EquipmentClient) Query() *EquipmentQuery {
	return &EquipmentQuery{config: c.config}
}

// Get returns a Equipment entity by its id.
func (c *EquipmentClient) Get(ctx context.Context, id string) (*Equipment, error) {
	return c.Query().Where(equipment.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EquipmentClient) GetX(ctx context.Context, id string) *Equipment {
	e, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return e
}

// QueryType queries the type edge of a Equipment.
func (c *EquipmentClient) QueryType(e *Equipment) *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: c.config}
	id := e.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, id),
		sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, equipment.TypeTable, equipment.TypeColumn),
	)
	query.sql = sqlgraph.Neighbors(e.driver.Dialect(), step)

	return query
}

// QueryLocation queries the location edge of a Equipment.
func (c *EquipmentClient) QueryLocation(e *Equipment) *LocationQuery {
	query := &LocationQuery{config: c.config}
	id := e.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, id),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, equipment.LocationTable, equipment.LocationColumn),
	)
	query.sql = sqlgraph.Neighbors(e.driver.Dialect(), step)

	return query
}

// QueryParentPosition queries the parent_position edge of a Equipment.
func (c *EquipmentClient) QueryParentPosition(e *Equipment) *EquipmentPositionQuery {
	query := &EquipmentPositionQuery{config: c.config}
	id := e.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, id),
		sqlgraph.To(equipmentposition.Table, equipmentposition.FieldID),
		sqlgraph.Edge(sqlgraph.O2O, true, equipment.ParentPositionTable, equipment.ParentPositionColumn),
	)
	query.sql = sqlgraph.Neighbors(e.driver.Dialect(), step)

	return query
}

// QueryPositions queries the positions edge of a Equipment.
func (c *EquipmentClient) QueryPositions(e *Equipment) *EquipmentPositionQuery {
	query := &EquipmentPositionQuery{config: c.config}
	id := e.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, id),
		sqlgraph.To(equipmentposition.Table, equipmentposition.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, equipment.PositionsTable, equipment.PositionsColumn),
	)
	query.sql = sqlgraph.Neighbors(e.driver.Dialect(), step)

	return query
}

// QueryPorts queries the ports edge of a Equipment.
func (c *EquipmentClient) QueryPorts(e *Equipment) *EquipmentPortQuery {
	query := &EquipmentPortQuery{config: c.config}
	id := e.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, id),
		sqlgraph.To(equipmentport.Table, equipmentport.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, equipment.PortsTable, equipment.PortsColumn),
	)
	query.sql = sqlgraph.Neighbors(e.driver.Dialect(), step)

	return query
}

// QueryWorkOrder queries the work_order edge of a Equipment.
func (c *EquipmentClient) QueryWorkOrder(e *Equipment) *WorkOrderQuery {
	query := &WorkOrderQuery{config: c.config}
	id := e.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, id),
		sqlgraph.To(workorder.Table, workorder.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, equipment.WorkOrderTable, equipment.WorkOrderColumn),
	)
	query.sql = sqlgraph.Neighbors(e.driver.Dialect(), step)

	return query
}

// QueryProperties queries the properties edge of a Equipment.
func (c *EquipmentClient) QueryProperties(e *Equipment) *PropertyQuery {
	query := &PropertyQuery{config: c.config}
	id := e.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, id),
		sqlgraph.To(property.Table, property.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, equipment.PropertiesTable, equipment.PropertiesColumn),
	)
	query.sql = sqlgraph.Neighbors(e.driver.Dialect(), step)

	return query
}

// QueryFiles queries the files edge of a Equipment.
func (c *EquipmentClient) QueryFiles(e *Equipment) *FileQuery {
	query := &FileQuery{config: c.config}
	id := e.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, id),
		sqlgraph.To(file.Table, file.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, equipment.FilesTable, equipment.FilesColumn),
	)
	query.sql = sqlgraph.Neighbors(e.driver.Dialect(), step)

	return query
}

// QueryHyperlinks queries the hyperlinks edge of a Equipment.
func (c *EquipmentClient) QueryHyperlinks(e *Equipment) *HyperlinkQuery {
	query := &HyperlinkQuery{config: c.config}
	id := e.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipment.Table, equipment.FieldID, id),
		sqlgraph.To(hyperlink.Table, hyperlink.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, equipment.HyperlinksTable, equipment.HyperlinksColumn),
	)
	query.sql = sqlgraph.Neighbors(e.driver.Dialect(), step)

	return query
}

// EquipmentCategoryClient is a client for the EquipmentCategory schema.
type EquipmentCategoryClient struct {
	config
}

// NewEquipmentCategoryClient returns a client for the EquipmentCategory from the given config.
func NewEquipmentCategoryClient(c config) *EquipmentCategoryClient {
	return &EquipmentCategoryClient{config: c}
}

// Create returns a create builder for EquipmentCategory.
func (c *EquipmentCategoryClient) Create() *EquipmentCategoryCreate {
	return &EquipmentCategoryCreate{config: c.config}
}

// Update returns an update builder for EquipmentCategory.
func (c *EquipmentCategoryClient) Update() *EquipmentCategoryUpdate {
	return &EquipmentCategoryUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *EquipmentCategoryClient) UpdateOne(ec *EquipmentCategory) *EquipmentCategoryUpdateOne {
	return c.UpdateOneID(ec.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *EquipmentCategoryClient) UpdateOneID(id string) *EquipmentCategoryUpdateOne {
	return &EquipmentCategoryUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for EquipmentCategory.
func (c *EquipmentCategoryClient) Delete() *EquipmentCategoryDelete {
	return &EquipmentCategoryDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *EquipmentCategoryClient) DeleteOne(ec *EquipmentCategory) *EquipmentCategoryDeleteOne {
	return c.DeleteOneID(ec.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *EquipmentCategoryClient) DeleteOneID(id string) *EquipmentCategoryDeleteOne {
	return &EquipmentCategoryDeleteOne{c.Delete().Where(equipmentcategory.ID(id))}
}

// Create returns a query builder for EquipmentCategory.
func (c *EquipmentCategoryClient) Query() *EquipmentCategoryQuery {
	return &EquipmentCategoryQuery{config: c.config}
}

// Get returns a EquipmentCategory entity by its id.
func (c *EquipmentCategoryClient) Get(ctx context.Context, id string) (*EquipmentCategory, error) {
	return c.Query().Where(equipmentcategory.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EquipmentCategoryClient) GetX(ctx context.Context, id string) *EquipmentCategory {
	ec, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return ec
}

// QueryTypes queries the types edge of a EquipmentCategory.
func (c *EquipmentCategoryClient) QueryTypes(ec *EquipmentCategory) *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: c.config}
	id := ec.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentcategory.Table, equipmentcategory.FieldID, id),
		sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, equipmentcategory.TypesTable, equipmentcategory.TypesColumn),
	)
	query.sql = sqlgraph.Neighbors(ec.driver.Dialect(), step)

	return query
}

// EquipmentPortClient is a client for the EquipmentPort schema.
type EquipmentPortClient struct {
	config
}

// NewEquipmentPortClient returns a client for the EquipmentPort from the given config.
func NewEquipmentPortClient(c config) *EquipmentPortClient {
	return &EquipmentPortClient{config: c}
}

// Create returns a create builder for EquipmentPort.
func (c *EquipmentPortClient) Create() *EquipmentPortCreate {
	return &EquipmentPortCreate{config: c.config}
}

// Update returns an update builder for EquipmentPort.
func (c *EquipmentPortClient) Update() *EquipmentPortUpdate {
	return &EquipmentPortUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *EquipmentPortClient) UpdateOne(ep *EquipmentPort) *EquipmentPortUpdateOne {
	return c.UpdateOneID(ep.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *EquipmentPortClient) UpdateOneID(id string) *EquipmentPortUpdateOne {
	return &EquipmentPortUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for EquipmentPort.
func (c *EquipmentPortClient) Delete() *EquipmentPortDelete {
	return &EquipmentPortDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *EquipmentPortClient) DeleteOne(ep *EquipmentPort) *EquipmentPortDeleteOne {
	return c.DeleteOneID(ep.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *EquipmentPortClient) DeleteOneID(id string) *EquipmentPortDeleteOne {
	return &EquipmentPortDeleteOne{c.Delete().Where(equipmentport.ID(id))}
}

// Create returns a query builder for EquipmentPort.
func (c *EquipmentPortClient) Query() *EquipmentPortQuery {
	return &EquipmentPortQuery{config: c.config}
}

// Get returns a EquipmentPort entity by its id.
func (c *EquipmentPortClient) Get(ctx context.Context, id string) (*EquipmentPort, error) {
	return c.Query().Where(equipmentport.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EquipmentPortClient) GetX(ctx context.Context, id string) *EquipmentPort {
	ep, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return ep
}

// QueryDefinition queries the definition edge of a EquipmentPort.
func (c *EquipmentPortClient) QueryDefinition(ep *EquipmentPort) *EquipmentPortDefinitionQuery {
	query := &EquipmentPortDefinitionQuery{config: c.config}
	id := ep.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentport.Table, equipmentport.FieldID, id),
		sqlgraph.To(equipmentportdefinition.Table, equipmentportdefinition.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, equipmentport.DefinitionTable, equipmentport.DefinitionColumn),
	)
	query.sql = sqlgraph.Neighbors(ep.driver.Dialect(), step)

	return query
}

// QueryParent queries the parent edge of a EquipmentPort.
func (c *EquipmentPortClient) QueryParent(ep *EquipmentPort) *EquipmentQuery {
	query := &EquipmentQuery{config: c.config}
	id := ep.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentport.Table, equipmentport.FieldID, id),
		sqlgraph.To(equipment.Table, equipment.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, equipmentport.ParentTable, equipmentport.ParentColumn),
	)
	query.sql = sqlgraph.Neighbors(ep.driver.Dialect(), step)

	return query
}

// QueryLink queries the link edge of a EquipmentPort.
func (c *EquipmentPortClient) QueryLink(ep *EquipmentPort) *LinkQuery {
	query := &LinkQuery{config: c.config}
	id := ep.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentport.Table, equipmentport.FieldID, id),
		sqlgraph.To(link.Table, link.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, equipmentport.LinkTable, equipmentport.LinkColumn),
	)
	query.sql = sqlgraph.Neighbors(ep.driver.Dialect(), step)

	return query
}

// QueryProperties queries the properties edge of a EquipmentPort.
func (c *EquipmentPortClient) QueryProperties(ep *EquipmentPort) *PropertyQuery {
	query := &PropertyQuery{config: c.config}
	id := ep.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentport.Table, equipmentport.FieldID, id),
		sqlgraph.To(property.Table, property.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, equipmentport.PropertiesTable, equipmentport.PropertiesColumn),
	)
	query.sql = sqlgraph.Neighbors(ep.driver.Dialect(), step)

	return query
}

// QueryEndpoints queries the endpoints edge of a EquipmentPort.
func (c *EquipmentPortClient) QueryEndpoints(ep *EquipmentPort) *ServiceEndpointQuery {
	query := &ServiceEndpointQuery{config: c.config}
	id := ep.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentport.Table, equipmentport.FieldID, id),
		sqlgraph.To(serviceendpoint.Table, serviceendpoint.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, equipmentport.EndpointsTable, equipmentport.EndpointsColumn),
	)
	query.sql = sqlgraph.Neighbors(ep.driver.Dialect(), step)

	return query
}

// EquipmentPortDefinitionClient is a client for the EquipmentPortDefinition schema.
type EquipmentPortDefinitionClient struct {
	config
}

// NewEquipmentPortDefinitionClient returns a client for the EquipmentPortDefinition from the given config.
func NewEquipmentPortDefinitionClient(c config) *EquipmentPortDefinitionClient {
	return &EquipmentPortDefinitionClient{config: c}
}

// Create returns a create builder for EquipmentPortDefinition.
func (c *EquipmentPortDefinitionClient) Create() *EquipmentPortDefinitionCreate {
	return &EquipmentPortDefinitionCreate{config: c.config}
}

// Update returns an update builder for EquipmentPortDefinition.
func (c *EquipmentPortDefinitionClient) Update() *EquipmentPortDefinitionUpdate {
	return &EquipmentPortDefinitionUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *EquipmentPortDefinitionClient) UpdateOne(epd *EquipmentPortDefinition) *EquipmentPortDefinitionUpdateOne {
	return c.UpdateOneID(epd.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *EquipmentPortDefinitionClient) UpdateOneID(id string) *EquipmentPortDefinitionUpdateOne {
	return &EquipmentPortDefinitionUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for EquipmentPortDefinition.
func (c *EquipmentPortDefinitionClient) Delete() *EquipmentPortDefinitionDelete {
	return &EquipmentPortDefinitionDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *EquipmentPortDefinitionClient) DeleteOne(epd *EquipmentPortDefinition) *EquipmentPortDefinitionDeleteOne {
	return c.DeleteOneID(epd.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *EquipmentPortDefinitionClient) DeleteOneID(id string) *EquipmentPortDefinitionDeleteOne {
	return &EquipmentPortDefinitionDeleteOne{c.Delete().Where(equipmentportdefinition.ID(id))}
}

// Create returns a query builder for EquipmentPortDefinition.
func (c *EquipmentPortDefinitionClient) Query() *EquipmentPortDefinitionQuery {
	return &EquipmentPortDefinitionQuery{config: c.config}
}

// Get returns a EquipmentPortDefinition entity by its id.
func (c *EquipmentPortDefinitionClient) Get(ctx context.Context, id string) (*EquipmentPortDefinition, error) {
	return c.Query().Where(equipmentportdefinition.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EquipmentPortDefinitionClient) GetX(ctx context.Context, id string) *EquipmentPortDefinition {
	epd, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return epd
}

// QueryEquipmentPortType queries the equipment_port_type edge of a EquipmentPortDefinition.
func (c *EquipmentPortDefinitionClient) QueryEquipmentPortType(epd *EquipmentPortDefinition) *EquipmentPortTypeQuery {
	query := &EquipmentPortTypeQuery{config: c.config}
	id := epd.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentportdefinition.Table, equipmentportdefinition.FieldID, id),
		sqlgraph.To(equipmentporttype.Table, equipmentporttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, equipmentportdefinition.EquipmentPortTypeTable, equipmentportdefinition.EquipmentPortTypeColumn),
	)
	query.sql = sqlgraph.Neighbors(epd.driver.Dialect(), step)

	return query
}

// QueryPorts queries the ports edge of a EquipmentPortDefinition.
func (c *EquipmentPortDefinitionClient) QueryPorts(epd *EquipmentPortDefinition) *EquipmentPortQuery {
	query := &EquipmentPortQuery{config: c.config}
	id := epd.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentportdefinition.Table, equipmentportdefinition.FieldID, id),
		sqlgraph.To(equipmentport.Table, equipmentport.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, equipmentportdefinition.PortsTable, equipmentportdefinition.PortsColumn),
	)
	query.sql = sqlgraph.Neighbors(epd.driver.Dialect(), step)

	return query
}

// QueryEquipmentType queries the equipment_type edge of a EquipmentPortDefinition.
func (c *EquipmentPortDefinitionClient) QueryEquipmentType(epd *EquipmentPortDefinition) *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: c.config}
	id := epd.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentportdefinition.Table, equipmentportdefinition.FieldID, id),
		sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, equipmentportdefinition.EquipmentTypeTable, equipmentportdefinition.EquipmentTypeColumn),
	)
	query.sql = sqlgraph.Neighbors(epd.driver.Dialect(), step)

	return query
}

// EquipmentPortTypeClient is a client for the EquipmentPortType schema.
type EquipmentPortTypeClient struct {
	config
}

// NewEquipmentPortTypeClient returns a client for the EquipmentPortType from the given config.
func NewEquipmentPortTypeClient(c config) *EquipmentPortTypeClient {
	return &EquipmentPortTypeClient{config: c}
}

// Create returns a create builder for EquipmentPortType.
func (c *EquipmentPortTypeClient) Create() *EquipmentPortTypeCreate {
	return &EquipmentPortTypeCreate{config: c.config}
}

// Update returns an update builder for EquipmentPortType.
func (c *EquipmentPortTypeClient) Update() *EquipmentPortTypeUpdate {
	return &EquipmentPortTypeUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *EquipmentPortTypeClient) UpdateOne(ept *EquipmentPortType) *EquipmentPortTypeUpdateOne {
	return c.UpdateOneID(ept.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *EquipmentPortTypeClient) UpdateOneID(id string) *EquipmentPortTypeUpdateOne {
	return &EquipmentPortTypeUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for EquipmentPortType.
func (c *EquipmentPortTypeClient) Delete() *EquipmentPortTypeDelete {
	return &EquipmentPortTypeDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *EquipmentPortTypeClient) DeleteOne(ept *EquipmentPortType) *EquipmentPortTypeDeleteOne {
	return c.DeleteOneID(ept.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *EquipmentPortTypeClient) DeleteOneID(id string) *EquipmentPortTypeDeleteOne {
	return &EquipmentPortTypeDeleteOne{c.Delete().Where(equipmentporttype.ID(id))}
}

// Create returns a query builder for EquipmentPortType.
func (c *EquipmentPortTypeClient) Query() *EquipmentPortTypeQuery {
	return &EquipmentPortTypeQuery{config: c.config}
}

// Get returns a EquipmentPortType entity by its id.
func (c *EquipmentPortTypeClient) Get(ctx context.Context, id string) (*EquipmentPortType, error) {
	return c.Query().Where(equipmentporttype.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EquipmentPortTypeClient) GetX(ctx context.Context, id string) *EquipmentPortType {
	ept, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return ept
}

// QueryPropertyTypes queries the property_types edge of a EquipmentPortType.
func (c *EquipmentPortTypeClient) QueryPropertyTypes(ept *EquipmentPortType) *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: c.config}
	id := ept.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentporttype.Table, equipmentporttype.FieldID, id),
		sqlgraph.To(propertytype.Table, propertytype.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, equipmentporttype.PropertyTypesTable, equipmentporttype.PropertyTypesColumn),
	)
	query.sql = sqlgraph.Neighbors(ept.driver.Dialect(), step)

	return query
}

// QueryLinkPropertyTypes queries the link_property_types edge of a EquipmentPortType.
func (c *EquipmentPortTypeClient) QueryLinkPropertyTypes(ept *EquipmentPortType) *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: c.config}
	id := ept.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentporttype.Table, equipmentporttype.FieldID, id),
		sqlgraph.To(propertytype.Table, propertytype.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, equipmentporttype.LinkPropertyTypesTable, equipmentporttype.LinkPropertyTypesColumn),
	)
	query.sql = sqlgraph.Neighbors(ept.driver.Dialect(), step)

	return query
}

// QueryPortDefinitions queries the port_definitions edge of a EquipmentPortType.
func (c *EquipmentPortTypeClient) QueryPortDefinitions(ept *EquipmentPortType) *EquipmentPortDefinitionQuery {
	query := &EquipmentPortDefinitionQuery{config: c.config}
	id := ept.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentporttype.Table, equipmentporttype.FieldID, id),
		sqlgraph.To(equipmentportdefinition.Table, equipmentportdefinition.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, equipmentporttype.PortDefinitionsTable, equipmentporttype.PortDefinitionsColumn),
	)
	query.sql = sqlgraph.Neighbors(ept.driver.Dialect(), step)

	return query
}

// EquipmentPositionClient is a client for the EquipmentPosition schema.
type EquipmentPositionClient struct {
	config
}

// NewEquipmentPositionClient returns a client for the EquipmentPosition from the given config.
func NewEquipmentPositionClient(c config) *EquipmentPositionClient {
	return &EquipmentPositionClient{config: c}
}

// Create returns a create builder for EquipmentPosition.
func (c *EquipmentPositionClient) Create() *EquipmentPositionCreate {
	return &EquipmentPositionCreate{config: c.config}
}

// Update returns an update builder for EquipmentPosition.
func (c *EquipmentPositionClient) Update() *EquipmentPositionUpdate {
	return &EquipmentPositionUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *EquipmentPositionClient) UpdateOne(ep *EquipmentPosition) *EquipmentPositionUpdateOne {
	return c.UpdateOneID(ep.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *EquipmentPositionClient) UpdateOneID(id string) *EquipmentPositionUpdateOne {
	return &EquipmentPositionUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for EquipmentPosition.
func (c *EquipmentPositionClient) Delete() *EquipmentPositionDelete {
	return &EquipmentPositionDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *EquipmentPositionClient) DeleteOne(ep *EquipmentPosition) *EquipmentPositionDeleteOne {
	return c.DeleteOneID(ep.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *EquipmentPositionClient) DeleteOneID(id string) *EquipmentPositionDeleteOne {
	return &EquipmentPositionDeleteOne{c.Delete().Where(equipmentposition.ID(id))}
}

// Create returns a query builder for EquipmentPosition.
func (c *EquipmentPositionClient) Query() *EquipmentPositionQuery {
	return &EquipmentPositionQuery{config: c.config}
}

// Get returns a EquipmentPosition entity by its id.
func (c *EquipmentPositionClient) Get(ctx context.Context, id string) (*EquipmentPosition, error) {
	return c.Query().Where(equipmentposition.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EquipmentPositionClient) GetX(ctx context.Context, id string) *EquipmentPosition {
	ep, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return ep
}

// QueryDefinition queries the definition edge of a EquipmentPosition.
func (c *EquipmentPositionClient) QueryDefinition(ep *EquipmentPosition) *EquipmentPositionDefinitionQuery {
	query := &EquipmentPositionDefinitionQuery{config: c.config}
	id := ep.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentposition.Table, equipmentposition.FieldID, id),
		sqlgraph.To(equipmentpositiondefinition.Table, equipmentpositiondefinition.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, equipmentposition.DefinitionTable, equipmentposition.DefinitionColumn),
	)
	query.sql = sqlgraph.Neighbors(ep.driver.Dialect(), step)

	return query
}

// QueryParent queries the parent edge of a EquipmentPosition.
func (c *EquipmentPositionClient) QueryParent(ep *EquipmentPosition) *EquipmentQuery {
	query := &EquipmentQuery{config: c.config}
	id := ep.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentposition.Table, equipmentposition.FieldID, id),
		sqlgraph.To(equipment.Table, equipment.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, equipmentposition.ParentTable, equipmentposition.ParentColumn),
	)
	query.sql = sqlgraph.Neighbors(ep.driver.Dialect(), step)

	return query
}

// QueryAttachment queries the attachment edge of a EquipmentPosition.
func (c *EquipmentPositionClient) QueryAttachment(ep *EquipmentPosition) *EquipmentQuery {
	query := &EquipmentQuery{config: c.config}
	id := ep.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentposition.Table, equipmentposition.FieldID, id),
		sqlgraph.To(equipment.Table, equipment.FieldID),
		sqlgraph.Edge(sqlgraph.O2O, false, equipmentposition.AttachmentTable, equipmentposition.AttachmentColumn),
	)
	query.sql = sqlgraph.Neighbors(ep.driver.Dialect(), step)

	return query
}

// EquipmentPositionDefinitionClient is a client for the EquipmentPositionDefinition schema.
type EquipmentPositionDefinitionClient struct {
	config
}

// NewEquipmentPositionDefinitionClient returns a client for the EquipmentPositionDefinition from the given config.
func NewEquipmentPositionDefinitionClient(c config) *EquipmentPositionDefinitionClient {
	return &EquipmentPositionDefinitionClient{config: c}
}

// Create returns a create builder for EquipmentPositionDefinition.
func (c *EquipmentPositionDefinitionClient) Create() *EquipmentPositionDefinitionCreate {
	return &EquipmentPositionDefinitionCreate{config: c.config}
}

// Update returns an update builder for EquipmentPositionDefinition.
func (c *EquipmentPositionDefinitionClient) Update() *EquipmentPositionDefinitionUpdate {
	return &EquipmentPositionDefinitionUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *EquipmentPositionDefinitionClient) UpdateOne(epd *EquipmentPositionDefinition) *EquipmentPositionDefinitionUpdateOne {
	return c.UpdateOneID(epd.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *EquipmentPositionDefinitionClient) UpdateOneID(id string) *EquipmentPositionDefinitionUpdateOne {
	return &EquipmentPositionDefinitionUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for EquipmentPositionDefinition.
func (c *EquipmentPositionDefinitionClient) Delete() *EquipmentPositionDefinitionDelete {
	return &EquipmentPositionDefinitionDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *EquipmentPositionDefinitionClient) DeleteOne(epd *EquipmentPositionDefinition) *EquipmentPositionDefinitionDeleteOne {
	return c.DeleteOneID(epd.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *EquipmentPositionDefinitionClient) DeleteOneID(id string) *EquipmentPositionDefinitionDeleteOne {
	return &EquipmentPositionDefinitionDeleteOne{c.Delete().Where(equipmentpositiondefinition.ID(id))}
}

// Create returns a query builder for EquipmentPositionDefinition.
func (c *EquipmentPositionDefinitionClient) Query() *EquipmentPositionDefinitionQuery {
	return &EquipmentPositionDefinitionQuery{config: c.config}
}

// Get returns a EquipmentPositionDefinition entity by its id.
func (c *EquipmentPositionDefinitionClient) Get(ctx context.Context, id string) (*EquipmentPositionDefinition, error) {
	return c.Query().Where(equipmentpositiondefinition.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EquipmentPositionDefinitionClient) GetX(ctx context.Context, id string) *EquipmentPositionDefinition {
	epd, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return epd
}

// QueryPositions queries the positions edge of a EquipmentPositionDefinition.
func (c *EquipmentPositionDefinitionClient) QueryPositions(epd *EquipmentPositionDefinition) *EquipmentPositionQuery {
	query := &EquipmentPositionQuery{config: c.config}
	id := epd.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentpositiondefinition.Table, equipmentpositiondefinition.FieldID, id),
		sqlgraph.To(equipmentposition.Table, equipmentposition.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, equipmentpositiondefinition.PositionsTable, equipmentpositiondefinition.PositionsColumn),
	)
	query.sql = sqlgraph.Neighbors(epd.driver.Dialect(), step)

	return query
}

// QueryEquipmentType queries the equipment_type edge of a EquipmentPositionDefinition.
func (c *EquipmentPositionDefinitionClient) QueryEquipmentType(epd *EquipmentPositionDefinition) *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: c.config}
	id := epd.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmentpositiondefinition.Table, equipmentpositiondefinition.FieldID, id),
		sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, equipmentpositiondefinition.EquipmentTypeTable, equipmentpositiondefinition.EquipmentTypeColumn),
	)
	query.sql = sqlgraph.Neighbors(epd.driver.Dialect(), step)

	return query
}

// EquipmentTypeClient is a client for the EquipmentType schema.
type EquipmentTypeClient struct {
	config
}

// NewEquipmentTypeClient returns a client for the EquipmentType from the given config.
func NewEquipmentTypeClient(c config) *EquipmentTypeClient {
	return &EquipmentTypeClient{config: c}
}

// Create returns a create builder for EquipmentType.
func (c *EquipmentTypeClient) Create() *EquipmentTypeCreate {
	return &EquipmentTypeCreate{config: c.config}
}

// Update returns an update builder for EquipmentType.
func (c *EquipmentTypeClient) Update() *EquipmentTypeUpdate {
	return &EquipmentTypeUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *EquipmentTypeClient) UpdateOne(et *EquipmentType) *EquipmentTypeUpdateOne {
	return c.UpdateOneID(et.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *EquipmentTypeClient) UpdateOneID(id string) *EquipmentTypeUpdateOne {
	return &EquipmentTypeUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for EquipmentType.
func (c *EquipmentTypeClient) Delete() *EquipmentTypeDelete {
	return &EquipmentTypeDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *EquipmentTypeClient) DeleteOne(et *EquipmentType) *EquipmentTypeDeleteOne {
	return c.DeleteOneID(et.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *EquipmentTypeClient) DeleteOneID(id string) *EquipmentTypeDeleteOne {
	return &EquipmentTypeDeleteOne{c.Delete().Where(equipmenttype.ID(id))}
}

// Create returns a query builder for EquipmentType.
func (c *EquipmentTypeClient) Query() *EquipmentTypeQuery {
	return &EquipmentTypeQuery{config: c.config}
}

// Get returns a EquipmentType entity by its id.
func (c *EquipmentTypeClient) Get(ctx context.Context, id string) (*EquipmentType, error) {
	return c.Query().Where(equipmenttype.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EquipmentTypeClient) GetX(ctx context.Context, id string) *EquipmentType {
	et, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return et
}

// QueryPortDefinitions queries the port_definitions edge of a EquipmentType.
func (c *EquipmentTypeClient) QueryPortDefinitions(et *EquipmentType) *EquipmentPortDefinitionQuery {
	query := &EquipmentPortDefinitionQuery{config: c.config}
	id := et.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmenttype.Table, equipmenttype.FieldID, id),
		sqlgraph.To(equipmentportdefinition.Table, equipmentportdefinition.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, equipmenttype.PortDefinitionsTable, equipmenttype.PortDefinitionsColumn),
	)
	query.sql = sqlgraph.Neighbors(et.driver.Dialect(), step)

	return query
}

// QueryPositionDefinitions queries the position_definitions edge of a EquipmentType.
func (c *EquipmentTypeClient) QueryPositionDefinitions(et *EquipmentType) *EquipmentPositionDefinitionQuery {
	query := &EquipmentPositionDefinitionQuery{config: c.config}
	id := et.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmenttype.Table, equipmenttype.FieldID, id),
		sqlgraph.To(equipmentpositiondefinition.Table, equipmentpositiondefinition.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, equipmenttype.PositionDefinitionsTable, equipmenttype.PositionDefinitionsColumn),
	)
	query.sql = sqlgraph.Neighbors(et.driver.Dialect(), step)

	return query
}

// QueryPropertyTypes queries the property_types edge of a EquipmentType.
func (c *EquipmentTypeClient) QueryPropertyTypes(et *EquipmentType) *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: c.config}
	id := et.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmenttype.Table, equipmenttype.FieldID, id),
		sqlgraph.To(propertytype.Table, propertytype.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, equipmenttype.PropertyTypesTable, equipmenttype.PropertyTypesColumn),
	)
	query.sql = sqlgraph.Neighbors(et.driver.Dialect(), step)

	return query
}

// QueryEquipment queries the equipment edge of a EquipmentType.
func (c *EquipmentTypeClient) QueryEquipment(et *EquipmentType) *EquipmentQuery {
	query := &EquipmentQuery{config: c.config}
	id := et.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmenttype.Table, equipmenttype.FieldID, id),
		sqlgraph.To(equipment.Table, equipment.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, equipmenttype.EquipmentTable, equipmenttype.EquipmentColumn),
	)
	query.sql = sqlgraph.Neighbors(et.driver.Dialect(), step)

	return query
}

// QueryCategory queries the category edge of a EquipmentType.
func (c *EquipmentTypeClient) QueryCategory(et *EquipmentType) *EquipmentCategoryQuery {
	query := &EquipmentCategoryQuery{config: c.config}
	id := et.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(equipmenttype.Table, equipmenttype.FieldID, id),
		sqlgraph.To(equipmentcategory.Table, equipmentcategory.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, equipmenttype.CategoryTable, equipmenttype.CategoryColumn),
	)
	query.sql = sqlgraph.Neighbors(et.driver.Dialect(), step)

	return query
}

// FileClient is a client for the File schema.
type FileClient struct {
	config
}

// NewFileClient returns a client for the File from the given config.
func NewFileClient(c config) *FileClient {
	return &FileClient{config: c}
}

// Create returns a create builder for File.
func (c *FileClient) Create() *FileCreate {
	return &FileCreate{config: c.config}
}

// Update returns an update builder for File.
func (c *FileClient) Update() *FileUpdate {
	return &FileUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *FileClient) UpdateOne(f *File) *FileUpdateOne {
	return c.UpdateOneID(f.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *FileClient) UpdateOneID(id string) *FileUpdateOne {
	return &FileUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for File.
func (c *FileClient) Delete() *FileDelete {
	return &FileDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *FileClient) DeleteOne(f *File) *FileDeleteOne {
	return c.DeleteOneID(f.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *FileClient) DeleteOneID(id string) *FileDeleteOne {
	return &FileDeleteOne{c.Delete().Where(file.ID(id))}
}

// Create returns a query builder for File.
func (c *FileClient) Query() *FileQuery {
	return &FileQuery{config: c.config}
}

// Get returns a File entity by its id.
func (c *FileClient) Get(ctx context.Context, id string) (*File, error) {
	return c.Query().Where(file.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *FileClient) GetX(ctx context.Context, id string) *File {
	f, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return f
}

// FloorPlanClient is a client for the FloorPlan schema.
type FloorPlanClient struct {
	config
}

// NewFloorPlanClient returns a client for the FloorPlan from the given config.
func NewFloorPlanClient(c config) *FloorPlanClient {
	return &FloorPlanClient{config: c}
}

// Create returns a create builder for FloorPlan.
func (c *FloorPlanClient) Create() *FloorPlanCreate {
	return &FloorPlanCreate{config: c.config}
}

// Update returns an update builder for FloorPlan.
func (c *FloorPlanClient) Update() *FloorPlanUpdate {
	return &FloorPlanUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *FloorPlanClient) UpdateOne(fp *FloorPlan) *FloorPlanUpdateOne {
	return c.UpdateOneID(fp.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *FloorPlanClient) UpdateOneID(id string) *FloorPlanUpdateOne {
	return &FloorPlanUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for FloorPlan.
func (c *FloorPlanClient) Delete() *FloorPlanDelete {
	return &FloorPlanDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *FloorPlanClient) DeleteOne(fp *FloorPlan) *FloorPlanDeleteOne {
	return c.DeleteOneID(fp.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *FloorPlanClient) DeleteOneID(id string) *FloorPlanDeleteOne {
	return &FloorPlanDeleteOne{c.Delete().Where(floorplan.ID(id))}
}

// Create returns a query builder for FloorPlan.
func (c *FloorPlanClient) Query() *FloorPlanQuery {
	return &FloorPlanQuery{config: c.config}
}

// Get returns a FloorPlan entity by its id.
func (c *FloorPlanClient) Get(ctx context.Context, id string) (*FloorPlan, error) {
	return c.Query().Where(floorplan.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *FloorPlanClient) GetX(ctx context.Context, id string) *FloorPlan {
	fp, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return fp
}

// QueryLocation queries the location edge of a FloorPlan.
func (c *FloorPlanClient) QueryLocation(fp *FloorPlan) *LocationQuery {
	query := &LocationQuery{config: c.config}
	id := fp.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(floorplan.Table, floorplan.FieldID, id),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, floorplan.LocationTable, floorplan.LocationColumn),
	)
	query.sql = sqlgraph.Neighbors(fp.driver.Dialect(), step)

	return query
}

// QueryReferencePoint queries the reference_point edge of a FloorPlan.
func (c *FloorPlanClient) QueryReferencePoint(fp *FloorPlan) *FloorPlanReferencePointQuery {
	query := &FloorPlanReferencePointQuery{config: c.config}
	id := fp.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(floorplan.Table, floorplan.FieldID, id),
		sqlgraph.To(floorplanreferencepoint.Table, floorplanreferencepoint.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, floorplan.ReferencePointTable, floorplan.ReferencePointColumn),
	)
	query.sql = sqlgraph.Neighbors(fp.driver.Dialect(), step)

	return query
}

// QueryScale queries the scale edge of a FloorPlan.
func (c *FloorPlanClient) QueryScale(fp *FloorPlan) *FloorPlanScaleQuery {
	query := &FloorPlanScaleQuery{config: c.config}
	id := fp.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(floorplan.Table, floorplan.FieldID, id),
		sqlgraph.To(floorplanscale.Table, floorplanscale.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, floorplan.ScaleTable, floorplan.ScaleColumn),
	)
	query.sql = sqlgraph.Neighbors(fp.driver.Dialect(), step)

	return query
}

// QueryImage queries the image edge of a FloorPlan.
func (c *FloorPlanClient) QueryImage(fp *FloorPlan) *FileQuery {
	query := &FileQuery{config: c.config}
	id := fp.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(floorplan.Table, floorplan.FieldID, id),
		sqlgraph.To(file.Table, file.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, floorplan.ImageTable, floorplan.ImageColumn),
	)
	query.sql = sqlgraph.Neighbors(fp.driver.Dialect(), step)

	return query
}

// FloorPlanReferencePointClient is a client for the FloorPlanReferencePoint schema.
type FloorPlanReferencePointClient struct {
	config
}

// NewFloorPlanReferencePointClient returns a client for the FloorPlanReferencePoint from the given config.
func NewFloorPlanReferencePointClient(c config) *FloorPlanReferencePointClient {
	return &FloorPlanReferencePointClient{config: c}
}

// Create returns a create builder for FloorPlanReferencePoint.
func (c *FloorPlanReferencePointClient) Create() *FloorPlanReferencePointCreate {
	return &FloorPlanReferencePointCreate{config: c.config}
}

// Update returns an update builder for FloorPlanReferencePoint.
func (c *FloorPlanReferencePointClient) Update() *FloorPlanReferencePointUpdate {
	return &FloorPlanReferencePointUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *FloorPlanReferencePointClient) UpdateOne(fprp *FloorPlanReferencePoint) *FloorPlanReferencePointUpdateOne {
	return c.UpdateOneID(fprp.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *FloorPlanReferencePointClient) UpdateOneID(id string) *FloorPlanReferencePointUpdateOne {
	return &FloorPlanReferencePointUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for FloorPlanReferencePoint.
func (c *FloorPlanReferencePointClient) Delete() *FloorPlanReferencePointDelete {
	return &FloorPlanReferencePointDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *FloorPlanReferencePointClient) DeleteOne(fprp *FloorPlanReferencePoint) *FloorPlanReferencePointDeleteOne {
	return c.DeleteOneID(fprp.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *FloorPlanReferencePointClient) DeleteOneID(id string) *FloorPlanReferencePointDeleteOne {
	return &FloorPlanReferencePointDeleteOne{c.Delete().Where(floorplanreferencepoint.ID(id))}
}

// Create returns a query builder for FloorPlanReferencePoint.
func (c *FloorPlanReferencePointClient) Query() *FloorPlanReferencePointQuery {
	return &FloorPlanReferencePointQuery{config: c.config}
}

// Get returns a FloorPlanReferencePoint entity by its id.
func (c *FloorPlanReferencePointClient) Get(ctx context.Context, id string) (*FloorPlanReferencePoint, error) {
	return c.Query().Where(floorplanreferencepoint.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *FloorPlanReferencePointClient) GetX(ctx context.Context, id string) *FloorPlanReferencePoint {
	fprp, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return fprp
}

// FloorPlanScaleClient is a client for the FloorPlanScale schema.
type FloorPlanScaleClient struct {
	config
}

// NewFloorPlanScaleClient returns a client for the FloorPlanScale from the given config.
func NewFloorPlanScaleClient(c config) *FloorPlanScaleClient {
	return &FloorPlanScaleClient{config: c}
}

// Create returns a create builder for FloorPlanScale.
func (c *FloorPlanScaleClient) Create() *FloorPlanScaleCreate {
	return &FloorPlanScaleCreate{config: c.config}
}

// Update returns an update builder for FloorPlanScale.
func (c *FloorPlanScaleClient) Update() *FloorPlanScaleUpdate {
	return &FloorPlanScaleUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *FloorPlanScaleClient) UpdateOne(fps *FloorPlanScale) *FloorPlanScaleUpdateOne {
	return c.UpdateOneID(fps.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *FloorPlanScaleClient) UpdateOneID(id string) *FloorPlanScaleUpdateOne {
	return &FloorPlanScaleUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for FloorPlanScale.
func (c *FloorPlanScaleClient) Delete() *FloorPlanScaleDelete {
	return &FloorPlanScaleDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *FloorPlanScaleClient) DeleteOne(fps *FloorPlanScale) *FloorPlanScaleDeleteOne {
	return c.DeleteOneID(fps.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *FloorPlanScaleClient) DeleteOneID(id string) *FloorPlanScaleDeleteOne {
	return &FloorPlanScaleDeleteOne{c.Delete().Where(floorplanscale.ID(id))}
}

// Create returns a query builder for FloorPlanScale.
func (c *FloorPlanScaleClient) Query() *FloorPlanScaleQuery {
	return &FloorPlanScaleQuery{config: c.config}
}

// Get returns a FloorPlanScale entity by its id.
func (c *FloorPlanScaleClient) Get(ctx context.Context, id string) (*FloorPlanScale, error) {
	return c.Query().Where(floorplanscale.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *FloorPlanScaleClient) GetX(ctx context.Context, id string) *FloorPlanScale {
	fps, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return fps
}

// HyperlinkClient is a client for the Hyperlink schema.
type HyperlinkClient struct {
	config
}

// NewHyperlinkClient returns a client for the Hyperlink from the given config.
func NewHyperlinkClient(c config) *HyperlinkClient {
	return &HyperlinkClient{config: c}
}

// Create returns a create builder for Hyperlink.
func (c *HyperlinkClient) Create() *HyperlinkCreate {
	return &HyperlinkCreate{config: c.config}
}

// Update returns an update builder for Hyperlink.
func (c *HyperlinkClient) Update() *HyperlinkUpdate {
	return &HyperlinkUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *HyperlinkClient) UpdateOne(h *Hyperlink) *HyperlinkUpdateOne {
	return c.UpdateOneID(h.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *HyperlinkClient) UpdateOneID(id string) *HyperlinkUpdateOne {
	return &HyperlinkUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for Hyperlink.
func (c *HyperlinkClient) Delete() *HyperlinkDelete {
	return &HyperlinkDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *HyperlinkClient) DeleteOne(h *Hyperlink) *HyperlinkDeleteOne {
	return c.DeleteOneID(h.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *HyperlinkClient) DeleteOneID(id string) *HyperlinkDeleteOne {
	return &HyperlinkDeleteOne{c.Delete().Where(hyperlink.ID(id))}
}

// Create returns a query builder for Hyperlink.
func (c *HyperlinkClient) Query() *HyperlinkQuery {
	return &HyperlinkQuery{config: c.config}
}

// Get returns a Hyperlink entity by its id.
func (c *HyperlinkClient) Get(ctx context.Context, id string) (*Hyperlink, error) {
	return c.Query().Where(hyperlink.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *HyperlinkClient) GetX(ctx context.Context, id string) *Hyperlink {
	h, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return h
}

// LinkClient is a client for the Link schema.
type LinkClient struct {
	config
}

// NewLinkClient returns a client for the Link from the given config.
func NewLinkClient(c config) *LinkClient {
	return &LinkClient{config: c}
}

// Create returns a create builder for Link.
func (c *LinkClient) Create() *LinkCreate {
	return &LinkCreate{config: c.config}
}

// Update returns an update builder for Link.
func (c *LinkClient) Update() *LinkUpdate {
	return &LinkUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *LinkClient) UpdateOne(l *Link) *LinkUpdateOne {
	return c.UpdateOneID(l.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *LinkClient) UpdateOneID(id string) *LinkUpdateOne {
	return &LinkUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for Link.
func (c *LinkClient) Delete() *LinkDelete {
	return &LinkDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *LinkClient) DeleteOne(l *Link) *LinkDeleteOne {
	return c.DeleteOneID(l.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *LinkClient) DeleteOneID(id string) *LinkDeleteOne {
	return &LinkDeleteOne{c.Delete().Where(link.ID(id))}
}

// Create returns a query builder for Link.
func (c *LinkClient) Query() *LinkQuery {
	return &LinkQuery{config: c.config}
}

// Get returns a Link entity by its id.
func (c *LinkClient) Get(ctx context.Context, id string) (*Link, error) {
	return c.Query().Where(link.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *LinkClient) GetX(ctx context.Context, id string) *Link {
	l, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return l
}

// QueryPorts queries the ports edge of a Link.
func (c *LinkClient) QueryPorts(l *Link) *EquipmentPortQuery {
	query := &EquipmentPortQuery{config: c.config}
	id := l.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(link.Table, link.FieldID, id),
		sqlgraph.To(equipmentport.Table, equipmentport.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, link.PortsTable, link.PortsColumn),
	)
	query.sql = sqlgraph.Neighbors(l.driver.Dialect(), step)

	return query
}

// QueryWorkOrder queries the work_order edge of a Link.
func (c *LinkClient) QueryWorkOrder(l *Link) *WorkOrderQuery {
	query := &WorkOrderQuery{config: c.config}
	id := l.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(link.Table, link.FieldID, id),
		sqlgraph.To(workorder.Table, workorder.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, link.WorkOrderTable, link.WorkOrderColumn),
	)
	query.sql = sqlgraph.Neighbors(l.driver.Dialect(), step)

	return query
}

// QueryProperties queries the properties edge of a Link.
func (c *LinkClient) QueryProperties(l *Link) *PropertyQuery {
	query := &PropertyQuery{config: c.config}
	id := l.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(link.Table, link.FieldID, id),
		sqlgraph.To(property.Table, property.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, link.PropertiesTable, link.PropertiesColumn),
	)
	query.sql = sqlgraph.Neighbors(l.driver.Dialect(), step)

	return query
}

// QueryService queries the service edge of a Link.
func (c *LinkClient) QueryService(l *Link) *ServiceQuery {
	query := &ServiceQuery{config: c.config}
	id := l.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(link.Table, link.FieldID, id),
		sqlgraph.To(service.Table, service.FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, link.ServiceTable, link.ServicePrimaryKey...),
	)
	query.sql = sqlgraph.Neighbors(l.driver.Dialect(), step)

	return query
}

// LocationClient is a client for the Location schema.
type LocationClient struct {
	config
}

// NewLocationClient returns a client for the Location from the given config.
func NewLocationClient(c config) *LocationClient {
	return &LocationClient{config: c}
}

// Create returns a create builder for Location.
func (c *LocationClient) Create() *LocationCreate {
	return &LocationCreate{config: c.config}
}

// Update returns an update builder for Location.
func (c *LocationClient) Update() *LocationUpdate {
	return &LocationUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *LocationClient) UpdateOne(l *Location) *LocationUpdateOne {
	return c.UpdateOneID(l.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *LocationClient) UpdateOneID(id string) *LocationUpdateOne {
	return &LocationUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for Location.
func (c *LocationClient) Delete() *LocationDelete {
	return &LocationDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *LocationClient) DeleteOne(l *Location) *LocationDeleteOne {
	return c.DeleteOneID(l.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *LocationClient) DeleteOneID(id string) *LocationDeleteOne {
	return &LocationDeleteOne{c.Delete().Where(location.ID(id))}
}

// Create returns a query builder for Location.
func (c *LocationClient) Query() *LocationQuery {
	return &LocationQuery{config: c.config}
}

// Get returns a Location entity by its id.
func (c *LocationClient) Get(ctx context.Context, id string) (*Location, error) {
	return c.Query().Where(location.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *LocationClient) GetX(ctx context.Context, id string) *Location {
	l, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return l
}

// QueryType queries the type edge of a Location.
func (c *LocationClient) QueryType(l *Location) *LocationTypeQuery {
	query := &LocationTypeQuery{config: c.config}
	id := l.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, id),
		sqlgraph.To(locationtype.Table, locationtype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, location.TypeTable, location.TypeColumn),
	)
	query.sql = sqlgraph.Neighbors(l.driver.Dialect(), step)

	return query
}

// QueryParent queries the parent edge of a Location.
func (c *LocationClient) QueryParent(l *Location) *LocationQuery {
	query := &LocationQuery{config: c.config}
	id := l.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, id),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, location.ParentTable, location.ParentColumn),
	)
	query.sql = sqlgraph.Neighbors(l.driver.Dialect(), step)

	return query
}

// QueryChildren queries the children edge of a Location.
func (c *LocationClient) QueryChildren(l *Location) *LocationQuery {
	query := &LocationQuery{config: c.config}
	id := l.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, id),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, location.ChildrenTable, location.ChildrenColumn),
	)
	query.sql = sqlgraph.Neighbors(l.driver.Dialect(), step)

	return query
}

// QueryFiles queries the files edge of a Location.
func (c *LocationClient) QueryFiles(l *Location) *FileQuery {
	query := &FileQuery{config: c.config}
	id := l.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, id),
		sqlgraph.To(file.Table, file.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, location.FilesTable, location.FilesColumn),
	)
	query.sql = sqlgraph.Neighbors(l.driver.Dialect(), step)

	return query
}

// QueryHyperlinks queries the hyperlinks edge of a Location.
func (c *LocationClient) QueryHyperlinks(l *Location) *HyperlinkQuery {
	query := &HyperlinkQuery{config: c.config}
	id := l.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, id),
		sqlgraph.To(hyperlink.Table, hyperlink.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, location.HyperlinksTable, location.HyperlinksColumn),
	)
	query.sql = sqlgraph.Neighbors(l.driver.Dialect(), step)

	return query
}

// QueryEquipment queries the equipment edge of a Location.
func (c *LocationClient) QueryEquipment(l *Location) *EquipmentQuery {
	query := &EquipmentQuery{config: c.config}
	id := l.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, id),
		sqlgraph.To(equipment.Table, equipment.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, location.EquipmentTable, location.EquipmentColumn),
	)
	query.sql = sqlgraph.Neighbors(l.driver.Dialect(), step)

	return query
}

// QueryProperties queries the properties edge of a Location.
func (c *LocationClient) QueryProperties(l *Location) *PropertyQuery {
	query := &PropertyQuery{config: c.config}
	id := l.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, id),
		sqlgraph.To(property.Table, property.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, location.PropertiesTable, location.PropertiesColumn),
	)
	query.sql = sqlgraph.Neighbors(l.driver.Dialect(), step)

	return query
}

// QuerySurvey queries the survey edge of a Location.
func (c *LocationClient) QuerySurvey(l *Location) *SurveyQuery {
	query := &SurveyQuery{config: c.config}
	id := l.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, id),
		sqlgraph.To(survey.Table, survey.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, location.SurveyTable, location.SurveyColumn),
	)
	query.sql = sqlgraph.Neighbors(l.driver.Dialect(), step)

	return query
}

// QueryWifiScan queries the wifi_scan edge of a Location.
func (c *LocationClient) QueryWifiScan(l *Location) *SurveyWiFiScanQuery {
	query := &SurveyWiFiScanQuery{config: c.config}
	id := l.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, id),
		sqlgraph.To(surveywifiscan.Table, surveywifiscan.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, location.WifiScanTable, location.WifiScanColumn),
	)
	query.sql = sqlgraph.Neighbors(l.driver.Dialect(), step)

	return query
}

// QueryCellScan queries the cell_scan edge of a Location.
func (c *LocationClient) QueryCellScan(l *Location) *SurveyCellScanQuery {
	query := &SurveyCellScanQuery{config: c.config}
	id := l.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, id),
		sqlgraph.To(surveycellscan.Table, surveycellscan.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, location.CellScanTable, location.CellScanColumn),
	)
	query.sql = sqlgraph.Neighbors(l.driver.Dialect(), step)

	return query
}

// QueryWorkOrders queries the work_orders edge of a Location.
func (c *LocationClient) QueryWorkOrders(l *Location) *WorkOrderQuery {
	query := &WorkOrderQuery{config: c.config}
	id := l.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, id),
		sqlgraph.To(workorder.Table, workorder.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, location.WorkOrdersTable, location.WorkOrdersColumn),
	)
	query.sql = sqlgraph.Neighbors(l.driver.Dialect(), step)

	return query
}

// QueryFloorPlans queries the floor_plans edge of a Location.
func (c *LocationClient) QueryFloorPlans(l *Location) *FloorPlanQuery {
	query := &FloorPlanQuery{config: c.config}
	id := l.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(location.Table, location.FieldID, id),
		sqlgraph.To(floorplan.Table, floorplan.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, location.FloorPlansTable, location.FloorPlansColumn),
	)
	query.sql = sqlgraph.Neighbors(l.driver.Dialect(), step)

	return query
}

// LocationTypeClient is a client for the LocationType schema.
type LocationTypeClient struct {
	config
}

// NewLocationTypeClient returns a client for the LocationType from the given config.
func NewLocationTypeClient(c config) *LocationTypeClient {
	return &LocationTypeClient{config: c}
}

// Create returns a create builder for LocationType.
func (c *LocationTypeClient) Create() *LocationTypeCreate {
	return &LocationTypeCreate{config: c.config}
}

// Update returns an update builder for LocationType.
func (c *LocationTypeClient) Update() *LocationTypeUpdate {
	return &LocationTypeUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *LocationTypeClient) UpdateOne(lt *LocationType) *LocationTypeUpdateOne {
	return c.UpdateOneID(lt.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *LocationTypeClient) UpdateOneID(id string) *LocationTypeUpdateOne {
	return &LocationTypeUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for LocationType.
func (c *LocationTypeClient) Delete() *LocationTypeDelete {
	return &LocationTypeDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *LocationTypeClient) DeleteOne(lt *LocationType) *LocationTypeDeleteOne {
	return c.DeleteOneID(lt.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *LocationTypeClient) DeleteOneID(id string) *LocationTypeDeleteOne {
	return &LocationTypeDeleteOne{c.Delete().Where(locationtype.ID(id))}
}

// Create returns a query builder for LocationType.
func (c *LocationTypeClient) Query() *LocationTypeQuery {
	return &LocationTypeQuery{config: c.config}
}

// Get returns a LocationType entity by its id.
func (c *LocationTypeClient) Get(ctx context.Context, id string) (*LocationType, error) {
	return c.Query().Where(locationtype.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *LocationTypeClient) GetX(ctx context.Context, id string) *LocationType {
	lt, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return lt
}

// QueryLocations queries the locations edge of a LocationType.
func (c *LocationTypeClient) QueryLocations(lt *LocationType) *LocationQuery {
	query := &LocationQuery{config: c.config}
	id := lt.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(locationtype.Table, locationtype.FieldID, id),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, locationtype.LocationsTable, locationtype.LocationsColumn),
	)
	query.sql = sqlgraph.Neighbors(lt.driver.Dialect(), step)

	return query
}

// QueryPropertyTypes queries the property_types edge of a LocationType.
func (c *LocationTypeClient) QueryPropertyTypes(lt *LocationType) *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: c.config}
	id := lt.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(locationtype.Table, locationtype.FieldID, id),
		sqlgraph.To(propertytype.Table, propertytype.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, locationtype.PropertyTypesTable, locationtype.PropertyTypesColumn),
	)
	query.sql = sqlgraph.Neighbors(lt.driver.Dialect(), step)

	return query
}

// QuerySurveyTemplateCategories queries the survey_template_categories edge of a LocationType.
func (c *LocationTypeClient) QuerySurveyTemplateCategories(lt *LocationType) *SurveyTemplateCategoryQuery {
	query := &SurveyTemplateCategoryQuery{config: c.config}
	id := lt.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(locationtype.Table, locationtype.FieldID, id),
		sqlgraph.To(surveytemplatecategory.Table, surveytemplatecategory.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, locationtype.SurveyTemplateCategoriesTable, locationtype.SurveyTemplateCategoriesColumn),
	)
	query.sql = sqlgraph.Neighbors(lt.driver.Dialect(), step)

	return query
}

// ProjectClient is a client for the Project schema.
type ProjectClient struct {
	config
}

// NewProjectClient returns a client for the Project from the given config.
func NewProjectClient(c config) *ProjectClient {
	return &ProjectClient{config: c}
}

// Create returns a create builder for Project.
func (c *ProjectClient) Create() *ProjectCreate {
	return &ProjectCreate{config: c.config}
}

// Update returns an update builder for Project.
func (c *ProjectClient) Update() *ProjectUpdate {
	return &ProjectUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *ProjectClient) UpdateOne(pr *Project) *ProjectUpdateOne {
	return c.UpdateOneID(pr.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *ProjectClient) UpdateOneID(id string) *ProjectUpdateOne {
	return &ProjectUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for Project.
func (c *ProjectClient) Delete() *ProjectDelete {
	return &ProjectDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *ProjectClient) DeleteOne(pr *Project) *ProjectDeleteOne {
	return c.DeleteOneID(pr.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *ProjectClient) DeleteOneID(id string) *ProjectDeleteOne {
	return &ProjectDeleteOne{c.Delete().Where(project.ID(id))}
}

// Create returns a query builder for Project.
func (c *ProjectClient) Query() *ProjectQuery {
	return &ProjectQuery{config: c.config}
}

// Get returns a Project entity by its id.
func (c *ProjectClient) Get(ctx context.Context, id string) (*Project, error) {
	return c.Query().Where(project.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ProjectClient) GetX(ctx context.Context, id string) *Project {
	pr, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return pr
}

// QueryType queries the type edge of a Project.
func (c *ProjectClient) QueryType(pr *Project) *ProjectTypeQuery {
	query := &ProjectTypeQuery{config: c.config}
	id := pr.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(project.Table, project.FieldID, id),
		sqlgraph.To(projecttype.Table, projecttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, project.TypeTable, project.TypeColumn),
	)
	query.sql = sqlgraph.Neighbors(pr.driver.Dialect(), step)

	return query
}

// QueryLocation queries the location edge of a Project.
func (c *ProjectClient) QueryLocation(pr *Project) *LocationQuery {
	query := &LocationQuery{config: c.config}
	id := pr.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(project.Table, project.FieldID, id),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, project.LocationTable, project.LocationColumn),
	)
	query.sql = sqlgraph.Neighbors(pr.driver.Dialect(), step)

	return query
}

// QueryComments queries the comments edge of a Project.
func (c *ProjectClient) QueryComments(pr *Project) *CommentQuery {
	query := &CommentQuery{config: c.config}
	id := pr.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(project.Table, project.FieldID, id),
		sqlgraph.To(comment.Table, comment.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, project.CommentsTable, project.CommentsColumn),
	)
	query.sql = sqlgraph.Neighbors(pr.driver.Dialect(), step)

	return query
}

// QueryWorkOrders queries the work_orders edge of a Project.
func (c *ProjectClient) QueryWorkOrders(pr *Project) *WorkOrderQuery {
	query := &WorkOrderQuery{config: c.config}
	id := pr.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(project.Table, project.FieldID, id),
		sqlgraph.To(workorder.Table, workorder.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, project.WorkOrdersTable, project.WorkOrdersColumn),
	)
	query.sql = sqlgraph.Neighbors(pr.driver.Dialect(), step)

	return query
}

// QueryProperties queries the properties edge of a Project.
func (c *ProjectClient) QueryProperties(pr *Project) *PropertyQuery {
	query := &PropertyQuery{config: c.config}
	id := pr.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(project.Table, project.FieldID, id),
		sqlgraph.To(property.Table, property.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, project.PropertiesTable, project.PropertiesColumn),
	)
	query.sql = sqlgraph.Neighbors(pr.driver.Dialect(), step)

	return query
}

// ProjectTypeClient is a client for the ProjectType schema.
type ProjectTypeClient struct {
	config
}

// NewProjectTypeClient returns a client for the ProjectType from the given config.
func NewProjectTypeClient(c config) *ProjectTypeClient {
	return &ProjectTypeClient{config: c}
}

// Create returns a create builder for ProjectType.
func (c *ProjectTypeClient) Create() *ProjectTypeCreate {
	return &ProjectTypeCreate{config: c.config}
}

// Update returns an update builder for ProjectType.
func (c *ProjectTypeClient) Update() *ProjectTypeUpdate {
	return &ProjectTypeUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *ProjectTypeClient) UpdateOne(pt *ProjectType) *ProjectTypeUpdateOne {
	return c.UpdateOneID(pt.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *ProjectTypeClient) UpdateOneID(id string) *ProjectTypeUpdateOne {
	return &ProjectTypeUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for ProjectType.
func (c *ProjectTypeClient) Delete() *ProjectTypeDelete {
	return &ProjectTypeDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *ProjectTypeClient) DeleteOne(pt *ProjectType) *ProjectTypeDeleteOne {
	return c.DeleteOneID(pt.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *ProjectTypeClient) DeleteOneID(id string) *ProjectTypeDeleteOne {
	return &ProjectTypeDeleteOne{c.Delete().Where(projecttype.ID(id))}
}

// Create returns a query builder for ProjectType.
func (c *ProjectTypeClient) Query() *ProjectTypeQuery {
	return &ProjectTypeQuery{config: c.config}
}

// Get returns a ProjectType entity by its id.
func (c *ProjectTypeClient) Get(ctx context.Context, id string) (*ProjectType, error) {
	return c.Query().Where(projecttype.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ProjectTypeClient) GetX(ctx context.Context, id string) *ProjectType {
	pt, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return pt
}

// QueryProjects queries the projects edge of a ProjectType.
func (c *ProjectTypeClient) QueryProjects(pt *ProjectType) *ProjectQuery {
	query := &ProjectQuery{config: c.config}
	id := pt.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(projecttype.Table, projecttype.FieldID, id),
		sqlgraph.To(project.Table, project.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, projecttype.ProjectsTable, projecttype.ProjectsColumn),
	)
	query.sql = sqlgraph.Neighbors(pt.driver.Dialect(), step)

	return query
}

// QueryProperties queries the properties edge of a ProjectType.
func (c *ProjectTypeClient) QueryProperties(pt *ProjectType) *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: c.config}
	id := pt.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(projecttype.Table, projecttype.FieldID, id),
		sqlgraph.To(propertytype.Table, propertytype.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, projecttype.PropertiesTable, projecttype.PropertiesColumn),
	)
	query.sql = sqlgraph.Neighbors(pt.driver.Dialect(), step)

	return query
}

// QueryWorkOrders queries the work_orders edge of a ProjectType.
func (c *ProjectTypeClient) QueryWorkOrders(pt *ProjectType) *WorkOrderDefinitionQuery {
	query := &WorkOrderDefinitionQuery{config: c.config}
	id := pt.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(projecttype.Table, projecttype.FieldID, id),
		sqlgraph.To(workorderdefinition.Table, workorderdefinition.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, projecttype.WorkOrdersTable, projecttype.WorkOrdersColumn),
	)
	query.sql = sqlgraph.Neighbors(pt.driver.Dialect(), step)

	return query
}

// PropertyClient is a client for the Property schema.
type PropertyClient struct {
	config
}

// NewPropertyClient returns a client for the Property from the given config.
func NewPropertyClient(c config) *PropertyClient {
	return &PropertyClient{config: c}
}

// Create returns a create builder for Property.
func (c *PropertyClient) Create() *PropertyCreate {
	return &PropertyCreate{config: c.config}
}

// Update returns an update builder for Property.
func (c *PropertyClient) Update() *PropertyUpdate {
	return &PropertyUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *PropertyClient) UpdateOne(pr *Property) *PropertyUpdateOne {
	return c.UpdateOneID(pr.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *PropertyClient) UpdateOneID(id string) *PropertyUpdateOne {
	return &PropertyUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for Property.
func (c *PropertyClient) Delete() *PropertyDelete {
	return &PropertyDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *PropertyClient) DeleteOne(pr *Property) *PropertyDeleteOne {
	return c.DeleteOneID(pr.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *PropertyClient) DeleteOneID(id string) *PropertyDeleteOne {
	return &PropertyDeleteOne{c.Delete().Where(property.ID(id))}
}

// Create returns a query builder for Property.
func (c *PropertyClient) Query() *PropertyQuery {
	return &PropertyQuery{config: c.config}
}

// Get returns a Property entity by its id.
func (c *PropertyClient) Get(ctx context.Context, id string) (*Property, error) {
	return c.Query().Where(property.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *PropertyClient) GetX(ctx context.Context, id string) *Property {
	pr, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return pr
}

// QueryType queries the type edge of a Property.
func (c *PropertyClient) QueryType(pr *Property) *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: c.config}
	id := pr.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, id),
		sqlgraph.To(propertytype.Table, propertytype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, property.TypeTable, property.TypeColumn),
	)
	query.sql = sqlgraph.Neighbors(pr.driver.Dialect(), step)

	return query
}

// QueryLocation queries the location edge of a Property.
func (c *PropertyClient) QueryLocation(pr *Property) *LocationQuery {
	query := &LocationQuery{config: c.config}
	id := pr.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, id),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, property.LocationTable, property.LocationColumn),
	)
	query.sql = sqlgraph.Neighbors(pr.driver.Dialect(), step)

	return query
}

// QueryEquipment queries the equipment edge of a Property.
func (c *PropertyClient) QueryEquipment(pr *Property) *EquipmentQuery {
	query := &EquipmentQuery{config: c.config}
	id := pr.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, id),
		sqlgraph.To(equipment.Table, equipment.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, property.EquipmentTable, property.EquipmentColumn),
	)
	query.sql = sqlgraph.Neighbors(pr.driver.Dialect(), step)

	return query
}

// QueryService queries the service edge of a Property.
func (c *PropertyClient) QueryService(pr *Property) *ServiceQuery {
	query := &ServiceQuery{config: c.config}
	id := pr.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, id),
		sqlgraph.To(service.Table, service.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, property.ServiceTable, property.ServiceColumn),
	)
	query.sql = sqlgraph.Neighbors(pr.driver.Dialect(), step)

	return query
}

// QueryEquipmentPort queries the equipment_port edge of a Property.
func (c *PropertyClient) QueryEquipmentPort(pr *Property) *EquipmentPortQuery {
	query := &EquipmentPortQuery{config: c.config}
	id := pr.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, id),
		sqlgraph.To(equipmentport.Table, equipmentport.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, property.EquipmentPortTable, property.EquipmentPortColumn),
	)
	query.sql = sqlgraph.Neighbors(pr.driver.Dialect(), step)

	return query
}

// QueryLink queries the link edge of a Property.
func (c *PropertyClient) QueryLink(pr *Property) *LinkQuery {
	query := &LinkQuery{config: c.config}
	id := pr.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, id),
		sqlgraph.To(link.Table, link.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, property.LinkTable, property.LinkColumn),
	)
	query.sql = sqlgraph.Neighbors(pr.driver.Dialect(), step)

	return query
}

// QueryWorkOrder queries the work_order edge of a Property.
func (c *PropertyClient) QueryWorkOrder(pr *Property) *WorkOrderQuery {
	query := &WorkOrderQuery{config: c.config}
	id := pr.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, id),
		sqlgraph.To(workorder.Table, workorder.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, property.WorkOrderTable, property.WorkOrderColumn),
	)
	query.sql = sqlgraph.Neighbors(pr.driver.Dialect(), step)

	return query
}

// QueryProject queries the project edge of a Property.
func (c *PropertyClient) QueryProject(pr *Property) *ProjectQuery {
	query := &ProjectQuery{config: c.config}
	id := pr.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, id),
		sqlgraph.To(project.Table, project.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, property.ProjectTable, property.ProjectColumn),
	)
	query.sql = sqlgraph.Neighbors(pr.driver.Dialect(), step)

	return query
}

// QueryEquipmentValue queries the equipment_value edge of a Property.
func (c *PropertyClient) QueryEquipmentValue(pr *Property) *EquipmentQuery {
	query := &EquipmentQuery{config: c.config}
	id := pr.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, id),
		sqlgraph.To(equipment.Table, equipment.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, property.EquipmentValueTable, property.EquipmentValueColumn),
	)
	query.sql = sqlgraph.Neighbors(pr.driver.Dialect(), step)

	return query
}

// QueryLocationValue queries the location_value edge of a Property.
func (c *PropertyClient) QueryLocationValue(pr *Property) *LocationQuery {
	query := &LocationQuery{config: c.config}
	id := pr.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, id),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, property.LocationValueTable, property.LocationValueColumn),
	)
	query.sql = sqlgraph.Neighbors(pr.driver.Dialect(), step)

	return query
}

// QueryServiceValue queries the service_value edge of a Property.
func (c *PropertyClient) QueryServiceValue(pr *Property) *ServiceQuery {
	query := &ServiceQuery{config: c.config}
	id := pr.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(property.Table, property.FieldID, id),
		sqlgraph.To(service.Table, service.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, property.ServiceValueTable, property.ServiceValueColumn),
	)
	query.sql = sqlgraph.Neighbors(pr.driver.Dialect(), step)

	return query
}

// PropertyTypeClient is a client for the PropertyType schema.
type PropertyTypeClient struct {
	config
}

// NewPropertyTypeClient returns a client for the PropertyType from the given config.
func NewPropertyTypeClient(c config) *PropertyTypeClient {
	return &PropertyTypeClient{config: c}
}

// Create returns a create builder for PropertyType.
func (c *PropertyTypeClient) Create() *PropertyTypeCreate {
	return &PropertyTypeCreate{config: c.config}
}

// Update returns an update builder for PropertyType.
func (c *PropertyTypeClient) Update() *PropertyTypeUpdate {
	return &PropertyTypeUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *PropertyTypeClient) UpdateOne(pt *PropertyType) *PropertyTypeUpdateOne {
	return c.UpdateOneID(pt.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *PropertyTypeClient) UpdateOneID(id string) *PropertyTypeUpdateOne {
	return &PropertyTypeUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for PropertyType.
func (c *PropertyTypeClient) Delete() *PropertyTypeDelete {
	return &PropertyTypeDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *PropertyTypeClient) DeleteOne(pt *PropertyType) *PropertyTypeDeleteOne {
	return c.DeleteOneID(pt.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *PropertyTypeClient) DeleteOneID(id string) *PropertyTypeDeleteOne {
	return &PropertyTypeDeleteOne{c.Delete().Where(propertytype.ID(id))}
}

// Create returns a query builder for PropertyType.
func (c *PropertyTypeClient) Query() *PropertyTypeQuery {
	return &PropertyTypeQuery{config: c.config}
}

// Get returns a PropertyType entity by its id.
func (c *PropertyTypeClient) Get(ctx context.Context, id string) (*PropertyType, error) {
	return c.Query().Where(propertytype.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *PropertyTypeClient) GetX(ctx context.Context, id string) *PropertyType {
	pt, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return pt
}

// QueryProperties queries the properties edge of a PropertyType.
func (c *PropertyTypeClient) QueryProperties(pt *PropertyType) *PropertyQuery {
	query := &PropertyQuery{config: c.config}
	id := pt.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(propertytype.Table, propertytype.FieldID, id),
		sqlgraph.To(property.Table, property.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, propertytype.PropertiesTable, propertytype.PropertiesColumn),
	)
	query.sql = sqlgraph.Neighbors(pt.driver.Dialect(), step)

	return query
}

// QueryLocationType queries the location_type edge of a PropertyType.
func (c *PropertyTypeClient) QueryLocationType(pt *PropertyType) *LocationTypeQuery {
	query := &LocationTypeQuery{config: c.config}
	id := pt.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(propertytype.Table, propertytype.FieldID, id),
		sqlgraph.To(locationtype.Table, locationtype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, propertytype.LocationTypeTable, propertytype.LocationTypeColumn),
	)
	query.sql = sqlgraph.Neighbors(pt.driver.Dialect(), step)

	return query
}

// QueryEquipmentPortType queries the equipment_port_type edge of a PropertyType.
func (c *PropertyTypeClient) QueryEquipmentPortType(pt *PropertyType) *EquipmentPortTypeQuery {
	query := &EquipmentPortTypeQuery{config: c.config}
	id := pt.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(propertytype.Table, propertytype.FieldID, id),
		sqlgraph.To(equipmentporttype.Table, equipmentporttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, propertytype.EquipmentPortTypeTable, propertytype.EquipmentPortTypeColumn),
	)
	query.sql = sqlgraph.Neighbors(pt.driver.Dialect(), step)

	return query
}

// QueryLinkEquipmentPortType queries the link_equipment_port_type edge of a PropertyType.
func (c *PropertyTypeClient) QueryLinkEquipmentPortType(pt *PropertyType) *EquipmentPortTypeQuery {
	query := &EquipmentPortTypeQuery{config: c.config}
	id := pt.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(propertytype.Table, propertytype.FieldID, id),
		sqlgraph.To(equipmentporttype.Table, equipmentporttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, propertytype.LinkEquipmentPortTypeTable, propertytype.LinkEquipmentPortTypeColumn),
	)
	query.sql = sqlgraph.Neighbors(pt.driver.Dialect(), step)

	return query
}

// QueryEquipmentType queries the equipment_type edge of a PropertyType.
func (c *PropertyTypeClient) QueryEquipmentType(pt *PropertyType) *EquipmentTypeQuery {
	query := &EquipmentTypeQuery{config: c.config}
	id := pt.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(propertytype.Table, propertytype.FieldID, id),
		sqlgraph.To(equipmenttype.Table, equipmenttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, propertytype.EquipmentTypeTable, propertytype.EquipmentTypeColumn),
	)
	query.sql = sqlgraph.Neighbors(pt.driver.Dialect(), step)

	return query
}

// QueryServiceType queries the service_type edge of a PropertyType.
func (c *PropertyTypeClient) QueryServiceType(pt *PropertyType) *ServiceTypeQuery {
	query := &ServiceTypeQuery{config: c.config}
	id := pt.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(propertytype.Table, propertytype.FieldID, id),
		sqlgraph.To(servicetype.Table, servicetype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, propertytype.ServiceTypeTable, propertytype.ServiceTypeColumn),
	)
	query.sql = sqlgraph.Neighbors(pt.driver.Dialect(), step)

	return query
}

// QueryWorkOrderType queries the work_order_type edge of a PropertyType.
func (c *PropertyTypeClient) QueryWorkOrderType(pt *PropertyType) *WorkOrderTypeQuery {
	query := &WorkOrderTypeQuery{config: c.config}
	id := pt.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(propertytype.Table, propertytype.FieldID, id),
		sqlgraph.To(workordertype.Table, workordertype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, propertytype.WorkOrderTypeTable, propertytype.WorkOrderTypeColumn),
	)
	query.sql = sqlgraph.Neighbors(pt.driver.Dialect(), step)

	return query
}

// QueryProjectType queries the project_type edge of a PropertyType.
func (c *PropertyTypeClient) QueryProjectType(pt *PropertyType) *ProjectTypeQuery {
	query := &ProjectTypeQuery{config: c.config}
	id := pt.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(propertytype.Table, propertytype.FieldID, id),
		sqlgraph.To(projecttype.Table, projecttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, propertytype.ProjectTypeTable, propertytype.ProjectTypeColumn),
	)
	query.sql = sqlgraph.Neighbors(pt.driver.Dialect(), step)

	return query
}

// ServiceClient is a client for the Service schema.
type ServiceClient struct {
	config
}

// NewServiceClient returns a client for the Service from the given config.
func NewServiceClient(c config) *ServiceClient {
	return &ServiceClient{config: c}
}

// Create returns a create builder for Service.
func (c *ServiceClient) Create() *ServiceCreate {
	return &ServiceCreate{config: c.config}
}

// Update returns an update builder for Service.
func (c *ServiceClient) Update() *ServiceUpdate {
	return &ServiceUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *ServiceClient) UpdateOne(s *Service) *ServiceUpdateOne {
	return c.UpdateOneID(s.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *ServiceClient) UpdateOneID(id string) *ServiceUpdateOne {
	return &ServiceUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for Service.
func (c *ServiceClient) Delete() *ServiceDelete {
	return &ServiceDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *ServiceClient) DeleteOne(s *Service) *ServiceDeleteOne {
	return c.DeleteOneID(s.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *ServiceClient) DeleteOneID(id string) *ServiceDeleteOne {
	return &ServiceDeleteOne{c.Delete().Where(service.ID(id))}
}

// Create returns a query builder for Service.
func (c *ServiceClient) Query() *ServiceQuery {
	return &ServiceQuery{config: c.config}
}

// Get returns a Service entity by its id.
func (c *ServiceClient) Get(ctx context.Context, id string) (*Service, error) {
	return c.Query().Where(service.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ServiceClient) GetX(ctx context.Context, id string) *Service {
	s, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return s
}

// QueryType queries the type edge of a Service.
func (c *ServiceClient) QueryType(s *Service) *ServiceTypeQuery {
	query := &ServiceTypeQuery{config: c.config}
	id := s.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(service.Table, service.FieldID, id),
		sqlgraph.To(servicetype.Table, servicetype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, service.TypeTable, service.TypeColumn),
	)
	query.sql = sqlgraph.Neighbors(s.driver.Dialect(), step)

	return query
}

// QueryDownstream queries the downstream edge of a Service.
func (c *ServiceClient) QueryDownstream(s *Service) *ServiceQuery {
	query := &ServiceQuery{config: c.config}
	id := s.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(service.Table, service.FieldID, id),
		sqlgraph.To(service.Table, service.FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, service.DownstreamTable, service.DownstreamPrimaryKey...),
	)
	query.sql = sqlgraph.Neighbors(s.driver.Dialect(), step)

	return query
}

// QueryUpstream queries the upstream edge of a Service.
func (c *ServiceClient) QueryUpstream(s *Service) *ServiceQuery {
	query := &ServiceQuery{config: c.config}
	id := s.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(service.Table, service.FieldID, id),
		sqlgraph.To(service.Table, service.FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, service.UpstreamTable, service.UpstreamPrimaryKey...),
	)
	query.sql = sqlgraph.Neighbors(s.driver.Dialect(), step)

	return query
}

// QueryProperties queries the properties edge of a Service.
func (c *ServiceClient) QueryProperties(s *Service) *PropertyQuery {
	query := &PropertyQuery{config: c.config}
	id := s.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(service.Table, service.FieldID, id),
		sqlgraph.To(property.Table, property.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, service.PropertiesTable, service.PropertiesColumn),
	)
	query.sql = sqlgraph.Neighbors(s.driver.Dialect(), step)

	return query
}

// QueryLinks queries the links edge of a Service.
func (c *ServiceClient) QueryLinks(s *Service) *LinkQuery {
	query := &LinkQuery{config: c.config}
	id := s.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(service.Table, service.FieldID, id),
		sqlgraph.To(link.Table, link.FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, service.LinksTable, service.LinksPrimaryKey...),
	)
	query.sql = sqlgraph.Neighbors(s.driver.Dialect(), step)

	return query
}

// QueryCustomer queries the customer edge of a Service.
func (c *ServiceClient) QueryCustomer(s *Service) *CustomerQuery {
	query := &CustomerQuery{config: c.config}
	id := s.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(service.Table, service.FieldID, id),
		sqlgraph.To(customer.Table, customer.FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, service.CustomerTable, service.CustomerPrimaryKey...),
	)
	query.sql = sqlgraph.Neighbors(s.driver.Dialect(), step)

	return query
}

// QueryEndpoints queries the endpoints edge of a Service.
func (c *ServiceClient) QueryEndpoints(s *Service) *ServiceEndpointQuery {
	query := &ServiceEndpointQuery{config: c.config}
	id := s.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(service.Table, service.FieldID, id),
		sqlgraph.To(serviceendpoint.Table, serviceendpoint.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, service.EndpointsTable, service.EndpointsColumn),
	)
	query.sql = sqlgraph.Neighbors(s.driver.Dialect(), step)

	return query
}

// ServiceEndpointClient is a client for the ServiceEndpoint schema.
type ServiceEndpointClient struct {
	config
}

// NewServiceEndpointClient returns a client for the ServiceEndpoint from the given config.
func NewServiceEndpointClient(c config) *ServiceEndpointClient {
	return &ServiceEndpointClient{config: c}
}

// Create returns a create builder for ServiceEndpoint.
func (c *ServiceEndpointClient) Create() *ServiceEndpointCreate {
	return &ServiceEndpointCreate{config: c.config}
}

// Update returns an update builder for ServiceEndpoint.
func (c *ServiceEndpointClient) Update() *ServiceEndpointUpdate {
	return &ServiceEndpointUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *ServiceEndpointClient) UpdateOne(se *ServiceEndpoint) *ServiceEndpointUpdateOne {
	return c.UpdateOneID(se.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *ServiceEndpointClient) UpdateOneID(id string) *ServiceEndpointUpdateOne {
	return &ServiceEndpointUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for ServiceEndpoint.
func (c *ServiceEndpointClient) Delete() *ServiceEndpointDelete {
	return &ServiceEndpointDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *ServiceEndpointClient) DeleteOne(se *ServiceEndpoint) *ServiceEndpointDeleteOne {
	return c.DeleteOneID(se.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *ServiceEndpointClient) DeleteOneID(id string) *ServiceEndpointDeleteOne {
	return &ServiceEndpointDeleteOne{c.Delete().Where(serviceendpoint.ID(id))}
}

// Create returns a query builder for ServiceEndpoint.
func (c *ServiceEndpointClient) Query() *ServiceEndpointQuery {
	return &ServiceEndpointQuery{config: c.config}
}

// Get returns a ServiceEndpoint entity by its id.
func (c *ServiceEndpointClient) Get(ctx context.Context, id string) (*ServiceEndpoint, error) {
	return c.Query().Where(serviceendpoint.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ServiceEndpointClient) GetX(ctx context.Context, id string) *ServiceEndpoint {
	se, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return se
}

// QueryPort queries the port edge of a ServiceEndpoint.
func (c *ServiceEndpointClient) QueryPort(se *ServiceEndpoint) *EquipmentPortQuery {
	query := &EquipmentPortQuery{config: c.config}
	id := se.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(serviceendpoint.Table, serviceendpoint.FieldID, id),
		sqlgraph.To(equipmentport.Table, equipmentport.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, serviceendpoint.PortTable, serviceendpoint.PortColumn),
	)
	query.sql = sqlgraph.Neighbors(se.driver.Dialect(), step)

	return query
}

// QueryService queries the service edge of a ServiceEndpoint.
func (c *ServiceEndpointClient) QueryService(se *ServiceEndpoint) *ServiceQuery {
	query := &ServiceQuery{config: c.config}
	id := se.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(serviceendpoint.Table, serviceendpoint.FieldID, id),
		sqlgraph.To(service.Table, service.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, serviceendpoint.ServiceTable, serviceendpoint.ServiceColumn),
	)
	query.sql = sqlgraph.Neighbors(se.driver.Dialect(), step)

	return query
}

// ServiceTypeClient is a client for the ServiceType schema.
type ServiceTypeClient struct {
	config
}

// NewServiceTypeClient returns a client for the ServiceType from the given config.
func NewServiceTypeClient(c config) *ServiceTypeClient {
	return &ServiceTypeClient{config: c}
}

// Create returns a create builder for ServiceType.
func (c *ServiceTypeClient) Create() *ServiceTypeCreate {
	return &ServiceTypeCreate{config: c.config}
}

// Update returns an update builder for ServiceType.
func (c *ServiceTypeClient) Update() *ServiceTypeUpdate {
	return &ServiceTypeUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *ServiceTypeClient) UpdateOne(st *ServiceType) *ServiceTypeUpdateOne {
	return c.UpdateOneID(st.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *ServiceTypeClient) UpdateOneID(id string) *ServiceTypeUpdateOne {
	return &ServiceTypeUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for ServiceType.
func (c *ServiceTypeClient) Delete() *ServiceTypeDelete {
	return &ServiceTypeDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *ServiceTypeClient) DeleteOne(st *ServiceType) *ServiceTypeDeleteOne {
	return c.DeleteOneID(st.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *ServiceTypeClient) DeleteOneID(id string) *ServiceTypeDeleteOne {
	return &ServiceTypeDeleteOne{c.Delete().Where(servicetype.ID(id))}
}

// Create returns a query builder for ServiceType.
func (c *ServiceTypeClient) Query() *ServiceTypeQuery {
	return &ServiceTypeQuery{config: c.config}
}

// Get returns a ServiceType entity by its id.
func (c *ServiceTypeClient) Get(ctx context.Context, id string) (*ServiceType, error) {
	return c.Query().Where(servicetype.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ServiceTypeClient) GetX(ctx context.Context, id string) *ServiceType {
	st, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return st
}

// QueryServices queries the services edge of a ServiceType.
func (c *ServiceTypeClient) QueryServices(st *ServiceType) *ServiceQuery {
	query := &ServiceQuery{config: c.config}
	id := st.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(servicetype.Table, servicetype.FieldID, id),
		sqlgraph.To(service.Table, service.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, servicetype.ServicesTable, servicetype.ServicesColumn),
	)
	query.sql = sqlgraph.Neighbors(st.driver.Dialect(), step)

	return query
}

// QueryPropertyTypes queries the property_types edge of a ServiceType.
func (c *ServiceTypeClient) QueryPropertyTypes(st *ServiceType) *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: c.config}
	id := st.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(servicetype.Table, servicetype.FieldID, id),
		sqlgraph.To(propertytype.Table, propertytype.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, servicetype.PropertyTypesTable, servicetype.PropertyTypesColumn),
	)
	query.sql = sqlgraph.Neighbors(st.driver.Dialect(), step)

	return query
}

// SurveyClient is a client for the Survey schema.
type SurveyClient struct {
	config
}

// NewSurveyClient returns a client for the Survey from the given config.
func NewSurveyClient(c config) *SurveyClient {
	return &SurveyClient{config: c}
}

// Create returns a create builder for Survey.
func (c *SurveyClient) Create() *SurveyCreate {
	return &SurveyCreate{config: c.config}
}

// Update returns an update builder for Survey.
func (c *SurveyClient) Update() *SurveyUpdate {
	return &SurveyUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *SurveyClient) UpdateOne(s *Survey) *SurveyUpdateOne {
	return c.UpdateOneID(s.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *SurveyClient) UpdateOneID(id string) *SurveyUpdateOne {
	return &SurveyUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for Survey.
func (c *SurveyClient) Delete() *SurveyDelete {
	return &SurveyDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *SurveyClient) DeleteOne(s *Survey) *SurveyDeleteOne {
	return c.DeleteOneID(s.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *SurveyClient) DeleteOneID(id string) *SurveyDeleteOne {
	return &SurveyDeleteOne{c.Delete().Where(survey.ID(id))}
}

// Create returns a query builder for Survey.
func (c *SurveyClient) Query() *SurveyQuery {
	return &SurveyQuery{config: c.config}
}

// Get returns a Survey entity by its id.
func (c *SurveyClient) Get(ctx context.Context, id string) (*Survey, error) {
	return c.Query().Where(survey.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SurveyClient) GetX(ctx context.Context, id string) *Survey {
	s, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return s
}

// QueryLocation queries the location edge of a Survey.
func (c *SurveyClient) QueryLocation(s *Survey) *LocationQuery {
	query := &LocationQuery{config: c.config}
	id := s.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(survey.Table, survey.FieldID, id),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, survey.LocationTable, survey.LocationColumn),
	)
	query.sql = sqlgraph.Neighbors(s.driver.Dialect(), step)

	return query
}

// QuerySourceFile queries the source_file edge of a Survey.
func (c *SurveyClient) QuerySourceFile(s *Survey) *FileQuery {
	query := &FileQuery{config: c.config}
	id := s.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(survey.Table, survey.FieldID, id),
		sqlgraph.To(file.Table, file.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, survey.SourceFileTable, survey.SourceFileColumn),
	)
	query.sql = sqlgraph.Neighbors(s.driver.Dialect(), step)

	return query
}

// QueryQuestions queries the questions edge of a Survey.
func (c *SurveyClient) QueryQuestions(s *Survey) *SurveyQuestionQuery {
	query := &SurveyQuestionQuery{config: c.config}
	id := s.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(survey.Table, survey.FieldID, id),
		sqlgraph.To(surveyquestion.Table, surveyquestion.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, survey.QuestionsTable, survey.QuestionsColumn),
	)
	query.sql = sqlgraph.Neighbors(s.driver.Dialect(), step)

	return query
}

// SurveyCellScanClient is a client for the SurveyCellScan schema.
type SurveyCellScanClient struct {
	config
}

// NewSurveyCellScanClient returns a client for the SurveyCellScan from the given config.
func NewSurveyCellScanClient(c config) *SurveyCellScanClient {
	return &SurveyCellScanClient{config: c}
}

// Create returns a create builder for SurveyCellScan.
func (c *SurveyCellScanClient) Create() *SurveyCellScanCreate {
	return &SurveyCellScanCreate{config: c.config}
}

// Update returns an update builder for SurveyCellScan.
func (c *SurveyCellScanClient) Update() *SurveyCellScanUpdate {
	return &SurveyCellScanUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *SurveyCellScanClient) UpdateOne(scs *SurveyCellScan) *SurveyCellScanUpdateOne {
	return c.UpdateOneID(scs.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *SurveyCellScanClient) UpdateOneID(id string) *SurveyCellScanUpdateOne {
	return &SurveyCellScanUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for SurveyCellScan.
func (c *SurveyCellScanClient) Delete() *SurveyCellScanDelete {
	return &SurveyCellScanDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *SurveyCellScanClient) DeleteOne(scs *SurveyCellScan) *SurveyCellScanDeleteOne {
	return c.DeleteOneID(scs.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *SurveyCellScanClient) DeleteOneID(id string) *SurveyCellScanDeleteOne {
	return &SurveyCellScanDeleteOne{c.Delete().Where(surveycellscan.ID(id))}
}

// Create returns a query builder for SurveyCellScan.
func (c *SurveyCellScanClient) Query() *SurveyCellScanQuery {
	return &SurveyCellScanQuery{config: c.config}
}

// Get returns a SurveyCellScan entity by its id.
func (c *SurveyCellScanClient) Get(ctx context.Context, id string) (*SurveyCellScan, error) {
	return c.Query().Where(surveycellscan.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SurveyCellScanClient) GetX(ctx context.Context, id string) *SurveyCellScan {
	scs, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return scs
}

// QuerySurveyQuestion queries the survey_question edge of a SurveyCellScan.
func (c *SurveyCellScanClient) QuerySurveyQuestion(scs *SurveyCellScan) *SurveyQuestionQuery {
	query := &SurveyQuestionQuery{config: c.config}
	id := scs.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(surveycellscan.Table, surveycellscan.FieldID, id),
		sqlgraph.To(surveyquestion.Table, surveyquestion.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, surveycellscan.SurveyQuestionTable, surveycellscan.SurveyQuestionColumn),
	)
	query.sql = sqlgraph.Neighbors(scs.driver.Dialect(), step)

	return query
}

// QueryLocation queries the location edge of a SurveyCellScan.
func (c *SurveyCellScanClient) QueryLocation(scs *SurveyCellScan) *LocationQuery {
	query := &LocationQuery{config: c.config}
	id := scs.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(surveycellscan.Table, surveycellscan.FieldID, id),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, surveycellscan.LocationTable, surveycellscan.LocationColumn),
	)
	query.sql = sqlgraph.Neighbors(scs.driver.Dialect(), step)

	return query
}

// SurveyQuestionClient is a client for the SurveyQuestion schema.
type SurveyQuestionClient struct {
	config
}

// NewSurveyQuestionClient returns a client for the SurveyQuestion from the given config.
func NewSurveyQuestionClient(c config) *SurveyQuestionClient {
	return &SurveyQuestionClient{config: c}
}

// Create returns a create builder for SurveyQuestion.
func (c *SurveyQuestionClient) Create() *SurveyQuestionCreate {
	return &SurveyQuestionCreate{config: c.config}
}

// Update returns an update builder for SurveyQuestion.
func (c *SurveyQuestionClient) Update() *SurveyQuestionUpdate {
	return &SurveyQuestionUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *SurveyQuestionClient) UpdateOne(sq *SurveyQuestion) *SurveyQuestionUpdateOne {
	return c.UpdateOneID(sq.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *SurveyQuestionClient) UpdateOneID(id string) *SurveyQuestionUpdateOne {
	return &SurveyQuestionUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for SurveyQuestion.
func (c *SurveyQuestionClient) Delete() *SurveyQuestionDelete {
	return &SurveyQuestionDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *SurveyQuestionClient) DeleteOne(sq *SurveyQuestion) *SurveyQuestionDeleteOne {
	return c.DeleteOneID(sq.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *SurveyQuestionClient) DeleteOneID(id string) *SurveyQuestionDeleteOne {
	return &SurveyQuestionDeleteOne{c.Delete().Where(surveyquestion.ID(id))}
}

// Create returns a query builder for SurveyQuestion.
func (c *SurveyQuestionClient) Query() *SurveyQuestionQuery {
	return &SurveyQuestionQuery{config: c.config}
}

// Get returns a SurveyQuestion entity by its id.
func (c *SurveyQuestionClient) Get(ctx context.Context, id string) (*SurveyQuestion, error) {
	return c.Query().Where(surveyquestion.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SurveyQuestionClient) GetX(ctx context.Context, id string) *SurveyQuestion {
	sq, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return sq
}

// QuerySurvey queries the survey edge of a SurveyQuestion.
func (c *SurveyQuestionClient) QuerySurvey(sq *SurveyQuestion) *SurveyQuery {
	query := &SurveyQuery{config: c.config}
	id := sq.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, id),
		sqlgraph.To(survey.Table, survey.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, surveyquestion.SurveyTable, surveyquestion.SurveyColumn),
	)
	query.sql = sqlgraph.Neighbors(sq.driver.Dialect(), step)

	return query
}

// QueryWifiScan queries the wifi_scan edge of a SurveyQuestion.
func (c *SurveyQuestionClient) QueryWifiScan(sq *SurveyQuestion) *SurveyWiFiScanQuery {
	query := &SurveyWiFiScanQuery{config: c.config}
	id := sq.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, id),
		sqlgraph.To(surveywifiscan.Table, surveywifiscan.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, surveyquestion.WifiScanTable, surveyquestion.WifiScanColumn),
	)
	query.sql = sqlgraph.Neighbors(sq.driver.Dialect(), step)

	return query
}

// QueryCellScan queries the cell_scan edge of a SurveyQuestion.
func (c *SurveyQuestionClient) QueryCellScan(sq *SurveyQuestion) *SurveyCellScanQuery {
	query := &SurveyCellScanQuery{config: c.config}
	id := sq.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, id),
		sqlgraph.To(surveycellscan.Table, surveycellscan.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, surveyquestion.CellScanTable, surveyquestion.CellScanColumn),
	)
	query.sql = sqlgraph.Neighbors(sq.driver.Dialect(), step)

	return query
}

// QueryPhotoData queries the photo_data edge of a SurveyQuestion.
func (c *SurveyQuestionClient) QueryPhotoData(sq *SurveyQuestion) *FileQuery {
	query := &FileQuery{config: c.config}
	id := sq.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(surveyquestion.Table, surveyquestion.FieldID, id),
		sqlgraph.To(file.Table, file.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, surveyquestion.PhotoDataTable, surveyquestion.PhotoDataColumn),
	)
	query.sql = sqlgraph.Neighbors(sq.driver.Dialect(), step)

	return query
}

// SurveyTemplateCategoryClient is a client for the SurveyTemplateCategory schema.
type SurveyTemplateCategoryClient struct {
	config
}

// NewSurveyTemplateCategoryClient returns a client for the SurveyTemplateCategory from the given config.
func NewSurveyTemplateCategoryClient(c config) *SurveyTemplateCategoryClient {
	return &SurveyTemplateCategoryClient{config: c}
}

// Create returns a create builder for SurveyTemplateCategory.
func (c *SurveyTemplateCategoryClient) Create() *SurveyTemplateCategoryCreate {
	return &SurveyTemplateCategoryCreate{config: c.config}
}

// Update returns an update builder for SurveyTemplateCategory.
func (c *SurveyTemplateCategoryClient) Update() *SurveyTemplateCategoryUpdate {
	return &SurveyTemplateCategoryUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *SurveyTemplateCategoryClient) UpdateOne(stc *SurveyTemplateCategory) *SurveyTemplateCategoryUpdateOne {
	return c.UpdateOneID(stc.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *SurveyTemplateCategoryClient) UpdateOneID(id string) *SurveyTemplateCategoryUpdateOne {
	return &SurveyTemplateCategoryUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for SurveyTemplateCategory.
func (c *SurveyTemplateCategoryClient) Delete() *SurveyTemplateCategoryDelete {
	return &SurveyTemplateCategoryDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *SurveyTemplateCategoryClient) DeleteOne(stc *SurveyTemplateCategory) *SurveyTemplateCategoryDeleteOne {
	return c.DeleteOneID(stc.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *SurveyTemplateCategoryClient) DeleteOneID(id string) *SurveyTemplateCategoryDeleteOne {
	return &SurveyTemplateCategoryDeleteOne{c.Delete().Where(surveytemplatecategory.ID(id))}
}

// Create returns a query builder for SurveyTemplateCategory.
func (c *SurveyTemplateCategoryClient) Query() *SurveyTemplateCategoryQuery {
	return &SurveyTemplateCategoryQuery{config: c.config}
}

// Get returns a SurveyTemplateCategory entity by its id.
func (c *SurveyTemplateCategoryClient) Get(ctx context.Context, id string) (*SurveyTemplateCategory, error) {
	return c.Query().Where(surveytemplatecategory.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SurveyTemplateCategoryClient) GetX(ctx context.Context, id string) *SurveyTemplateCategory {
	stc, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return stc
}

// QuerySurveyTemplateQuestions queries the survey_template_questions edge of a SurveyTemplateCategory.
func (c *SurveyTemplateCategoryClient) QuerySurveyTemplateQuestions(stc *SurveyTemplateCategory) *SurveyTemplateQuestionQuery {
	query := &SurveyTemplateQuestionQuery{config: c.config}
	id := stc.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(surveytemplatecategory.Table, surveytemplatecategory.FieldID, id),
		sqlgraph.To(surveytemplatequestion.Table, surveytemplatequestion.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, surveytemplatecategory.SurveyTemplateQuestionsTable, surveytemplatecategory.SurveyTemplateQuestionsColumn),
	)
	query.sql = sqlgraph.Neighbors(stc.driver.Dialect(), step)

	return query
}

// SurveyTemplateQuestionClient is a client for the SurveyTemplateQuestion schema.
type SurveyTemplateQuestionClient struct {
	config
}

// NewSurveyTemplateQuestionClient returns a client for the SurveyTemplateQuestion from the given config.
func NewSurveyTemplateQuestionClient(c config) *SurveyTemplateQuestionClient {
	return &SurveyTemplateQuestionClient{config: c}
}

// Create returns a create builder for SurveyTemplateQuestion.
func (c *SurveyTemplateQuestionClient) Create() *SurveyTemplateQuestionCreate {
	return &SurveyTemplateQuestionCreate{config: c.config}
}

// Update returns an update builder for SurveyTemplateQuestion.
func (c *SurveyTemplateQuestionClient) Update() *SurveyTemplateQuestionUpdate {
	return &SurveyTemplateQuestionUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *SurveyTemplateQuestionClient) UpdateOne(stq *SurveyTemplateQuestion) *SurveyTemplateQuestionUpdateOne {
	return c.UpdateOneID(stq.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *SurveyTemplateQuestionClient) UpdateOneID(id string) *SurveyTemplateQuestionUpdateOne {
	return &SurveyTemplateQuestionUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for SurveyTemplateQuestion.
func (c *SurveyTemplateQuestionClient) Delete() *SurveyTemplateQuestionDelete {
	return &SurveyTemplateQuestionDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *SurveyTemplateQuestionClient) DeleteOne(stq *SurveyTemplateQuestion) *SurveyTemplateQuestionDeleteOne {
	return c.DeleteOneID(stq.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *SurveyTemplateQuestionClient) DeleteOneID(id string) *SurveyTemplateQuestionDeleteOne {
	return &SurveyTemplateQuestionDeleteOne{c.Delete().Where(surveytemplatequestion.ID(id))}
}

// Create returns a query builder for SurveyTemplateQuestion.
func (c *SurveyTemplateQuestionClient) Query() *SurveyTemplateQuestionQuery {
	return &SurveyTemplateQuestionQuery{config: c.config}
}

// Get returns a SurveyTemplateQuestion entity by its id.
func (c *SurveyTemplateQuestionClient) Get(ctx context.Context, id string) (*SurveyTemplateQuestion, error) {
	return c.Query().Where(surveytemplatequestion.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SurveyTemplateQuestionClient) GetX(ctx context.Context, id string) *SurveyTemplateQuestion {
	stq, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return stq
}

// QueryCategory queries the category edge of a SurveyTemplateQuestion.
func (c *SurveyTemplateQuestionClient) QueryCategory(stq *SurveyTemplateQuestion) *SurveyTemplateCategoryQuery {
	query := &SurveyTemplateCategoryQuery{config: c.config}
	id := stq.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(surveytemplatequestion.Table, surveytemplatequestion.FieldID, id),
		sqlgraph.To(surveytemplatecategory.Table, surveytemplatecategory.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, surveytemplatequestion.CategoryTable, surveytemplatequestion.CategoryColumn),
	)
	query.sql = sqlgraph.Neighbors(stq.driver.Dialect(), step)

	return query
}

// SurveyWiFiScanClient is a client for the SurveyWiFiScan schema.
type SurveyWiFiScanClient struct {
	config
}

// NewSurveyWiFiScanClient returns a client for the SurveyWiFiScan from the given config.
func NewSurveyWiFiScanClient(c config) *SurveyWiFiScanClient {
	return &SurveyWiFiScanClient{config: c}
}

// Create returns a create builder for SurveyWiFiScan.
func (c *SurveyWiFiScanClient) Create() *SurveyWiFiScanCreate {
	return &SurveyWiFiScanCreate{config: c.config}
}

// Update returns an update builder for SurveyWiFiScan.
func (c *SurveyWiFiScanClient) Update() *SurveyWiFiScanUpdate {
	return &SurveyWiFiScanUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *SurveyWiFiScanClient) UpdateOne(swfs *SurveyWiFiScan) *SurveyWiFiScanUpdateOne {
	return c.UpdateOneID(swfs.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *SurveyWiFiScanClient) UpdateOneID(id string) *SurveyWiFiScanUpdateOne {
	return &SurveyWiFiScanUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for SurveyWiFiScan.
func (c *SurveyWiFiScanClient) Delete() *SurveyWiFiScanDelete {
	return &SurveyWiFiScanDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *SurveyWiFiScanClient) DeleteOne(swfs *SurveyWiFiScan) *SurveyWiFiScanDeleteOne {
	return c.DeleteOneID(swfs.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *SurveyWiFiScanClient) DeleteOneID(id string) *SurveyWiFiScanDeleteOne {
	return &SurveyWiFiScanDeleteOne{c.Delete().Where(surveywifiscan.ID(id))}
}

// Create returns a query builder for SurveyWiFiScan.
func (c *SurveyWiFiScanClient) Query() *SurveyWiFiScanQuery {
	return &SurveyWiFiScanQuery{config: c.config}
}

// Get returns a SurveyWiFiScan entity by its id.
func (c *SurveyWiFiScanClient) Get(ctx context.Context, id string) (*SurveyWiFiScan, error) {
	return c.Query().Where(surveywifiscan.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SurveyWiFiScanClient) GetX(ctx context.Context, id string) *SurveyWiFiScan {
	swfs, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return swfs
}

// QuerySurveyQuestion queries the survey_question edge of a SurveyWiFiScan.
func (c *SurveyWiFiScanClient) QuerySurveyQuestion(swfs *SurveyWiFiScan) *SurveyQuestionQuery {
	query := &SurveyQuestionQuery{config: c.config}
	id := swfs.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(surveywifiscan.Table, surveywifiscan.FieldID, id),
		sqlgraph.To(surveyquestion.Table, surveyquestion.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, surveywifiscan.SurveyQuestionTable, surveywifiscan.SurveyQuestionColumn),
	)
	query.sql = sqlgraph.Neighbors(swfs.driver.Dialect(), step)

	return query
}

// QueryLocation queries the location edge of a SurveyWiFiScan.
func (c *SurveyWiFiScanClient) QueryLocation(swfs *SurveyWiFiScan) *LocationQuery {
	query := &LocationQuery{config: c.config}
	id := swfs.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(surveywifiscan.Table, surveywifiscan.FieldID, id),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, surveywifiscan.LocationTable, surveywifiscan.LocationColumn),
	)
	query.sql = sqlgraph.Neighbors(swfs.driver.Dialect(), step)

	return query
}

// TechnicianClient is a client for the Technician schema.
type TechnicianClient struct {
	config
}

// NewTechnicianClient returns a client for the Technician from the given config.
func NewTechnicianClient(c config) *TechnicianClient {
	return &TechnicianClient{config: c}
}

// Create returns a create builder for Technician.
func (c *TechnicianClient) Create() *TechnicianCreate {
	return &TechnicianCreate{config: c.config}
}

// Update returns an update builder for Technician.
func (c *TechnicianClient) Update() *TechnicianUpdate {
	return &TechnicianUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *TechnicianClient) UpdateOne(t *Technician) *TechnicianUpdateOne {
	return c.UpdateOneID(t.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *TechnicianClient) UpdateOneID(id string) *TechnicianUpdateOne {
	return &TechnicianUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for Technician.
func (c *TechnicianClient) Delete() *TechnicianDelete {
	return &TechnicianDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *TechnicianClient) DeleteOne(t *Technician) *TechnicianDeleteOne {
	return c.DeleteOneID(t.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *TechnicianClient) DeleteOneID(id string) *TechnicianDeleteOne {
	return &TechnicianDeleteOne{c.Delete().Where(technician.ID(id))}
}

// Create returns a query builder for Technician.
func (c *TechnicianClient) Query() *TechnicianQuery {
	return &TechnicianQuery{config: c.config}
}

// Get returns a Technician entity by its id.
func (c *TechnicianClient) Get(ctx context.Context, id string) (*Technician, error) {
	return c.Query().Where(technician.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *TechnicianClient) GetX(ctx context.Context, id string) *Technician {
	t, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return t
}

// QueryWorkOrders queries the work_orders edge of a Technician.
func (c *TechnicianClient) QueryWorkOrders(t *Technician) *WorkOrderQuery {
	query := &WorkOrderQuery{config: c.config}
	id := t.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(technician.Table, technician.FieldID, id),
		sqlgraph.To(workorder.Table, workorder.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, technician.WorkOrdersTable, technician.WorkOrdersColumn),
	)
	query.sql = sqlgraph.Neighbors(t.driver.Dialect(), step)

	return query
}

// WorkOrderClient is a client for the WorkOrder schema.
type WorkOrderClient struct {
	config
}

// NewWorkOrderClient returns a client for the WorkOrder from the given config.
func NewWorkOrderClient(c config) *WorkOrderClient {
	return &WorkOrderClient{config: c}
}

// Create returns a create builder for WorkOrder.
func (c *WorkOrderClient) Create() *WorkOrderCreate {
	return &WorkOrderCreate{config: c.config}
}

// Update returns an update builder for WorkOrder.
func (c *WorkOrderClient) Update() *WorkOrderUpdate {
	return &WorkOrderUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *WorkOrderClient) UpdateOne(wo *WorkOrder) *WorkOrderUpdateOne {
	return c.UpdateOneID(wo.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *WorkOrderClient) UpdateOneID(id string) *WorkOrderUpdateOne {
	return &WorkOrderUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for WorkOrder.
func (c *WorkOrderClient) Delete() *WorkOrderDelete {
	return &WorkOrderDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *WorkOrderClient) DeleteOne(wo *WorkOrder) *WorkOrderDeleteOne {
	return c.DeleteOneID(wo.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *WorkOrderClient) DeleteOneID(id string) *WorkOrderDeleteOne {
	return &WorkOrderDeleteOne{c.Delete().Where(workorder.ID(id))}
}

// Create returns a query builder for WorkOrder.
func (c *WorkOrderClient) Query() *WorkOrderQuery {
	return &WorkOrderQuery{config: c.config}
}

// Get returns a WorkOrder entity by its id.
func (c *WorkOrderClient) Get(ctx context.Context, id string) (*WorkOrder, error) {
	return c.Query().Where(workorder.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *WorkOrderClient) GetX(ctx context.Context, id string) *WorkOrder {
	wo, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return wo
}

// QueryType queries the type edge of a WorkOrder.
func (c *WorkOrderClient) QueryType(wo *WorkOrder) *WorkOrderTypeQuery {
	query := &WorkOrderTypeQuery{config: c.config}
	id := wo.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workorder.Table, workorder.FieldID, id),
		sqlgraph.To(workordertype.Table, workordertype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, workorder.TypeTable, workorder.TypeColumn),
	)
	query.sql = sqlgraph.Neighbors(wo.driver.Dialect(), step)

	return query
}

// QueryEquipment queries the equipment edge of a WorkOrder.
func (c *WorkOrderClient) QueryEquipment(wo *WorkOrder) *EquipmentQuery {
	query := &EquipmentQuery{config: c.config}
	id := wo.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workorder.Table, workorder.FieldID, id),
		sqlgraph.To(equipment.Table, equipment.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, workorder.EquipmentTable, workorder.EquipmentColumn),
	)
	query.sql = sqlgraph.Neighbors(wo.driver.Dialect(), step)

	return query
}

// QueryLinks queries the links edge of a WorkOrder.
func (c *WorkOrderClient) QueryLinks(wo *WorkOrder) *LinkQuery {
	query := &LinkQuery{config: c.config}
	id := wo.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workorder.Table, workorder.FieldID, id),
		sqlgraph.To(link.Table, link.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, workorder.LinksTable, workorder.LinksColumn),
	)
	query.sql = sqlgraph.Neighbors(wo.driver.Dialect(), step)

	return query
}

// QueryFiles queries the files edge of a WorkOrder.
func (c *WorkOrderClient) QueryFiles(wo *WorkOrder) *FileQuery {
	query := &FileQuery{config: c.config}
	id := wo.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workorder.Table, workorder.FieldID, id),
		sqlgraph.To(file.Table, file.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, workorder.FilesTable, workorder.FilesColumn),
	)
	query.sql = sqlgraph.Neighbors(wo.driver.Dialect(), step)

	return query
}

// QueryHyperlinks queries the hyperlinks edge of a WorkOrder.
func (c *WorkOrderClient) QueryHyperlinks(wo *WorkOrder) *HyperlinkQuery {
	query := &HyperlinkQuery{config: c.config}
	id := wo.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workorder.Table, workorder.FieldID, id),
		sqlgraph.To(hyperlink.Table, hyperlink.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, workorder.HyperlinksTable, workorder.HyperlinksColumn),
	)
	query.sql = sqlgraph.Neighbors(wo.driver.Dialect(), step)

	return query
}

// QueryLocation queries the location edge of a WorkOrder.
func (c *WorkOrderClient) QueryLocation(wo *WorkOrder) *LocationQuery {
	query := &LocationQuery{config: c.config}
	id := wo.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workorder.Table, workorder.FieldID, id),
		sqlgraph.To(location.Table, location.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, workorder.LocationTable, workorder.LocationColumn),
	)
	query.sql = sqlgraph.Neighbors(wo.driver.Dialect(), step)

	return query
}

// QueryComments queries the comments edge of a WorkOrder.
func (c *WorkOrderClient) QueryComments(wo *WorkOrder) *CommentQuery {
	query := &CommentQuery{config: c.config}
	id := wo.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workorder.Table, workorder.FieldID, id),
		sqlgraph.To(comment.Table, comment.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, workorder.CommentsTable, workorder.CommentsColumn),
	)
	query.sql = sqlgraph.Neighbors(wo.driver.Dialect(), step)

	return query
}

// QueryProperties queries the properties edge of a WorkOrder.
func (c *WorkOrderClient) QueryProperties(wo *WorkOrder) *PropertyQuery {
	query := &PropertyQuery{config: c.config}
	id := wo.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workorder.Table, workorder.FieldID, id),
		sqlgraph.To(property.Table, property.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, workorder.PropertiesTable, workorder.PropertiesColumn),
	)
	query.sql = sqlgraph.Neighbors(wo.driver.Dialect(), step)

	return query
}

// QueryCheckListCategories queries the check_list_categories edge of a WorkOrder.
func (c *WorkOrderClient) QueryCheckListCategories(wo *WorkOrder) *CheckListCategoryQuery {
	query := &CheckListCategoryQuery{config: c.config}
	id := wo.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workorder.Table, workorder.FieldID, id),
		sqlgraph.To(checklistcategory.Table, checklistcategory.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, workorder.CheckListCategoriesTable, workorder.CheckListCategoriesColumn),
	)
	query.sql = sqlgraph.Neighbors(wo.driver.Dialect(), step)

	return query
}

// QueryCheckListItems queries the check_list_items edge of a WorkOrder.
func (c *WorkOrderClient) QueryCheckListItems(wo *WorkOrder) *CheckListItemQuery {
	query := &CheckListItemQuery{config: c.config}
	id := wo.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workorder.Table, workorder.FieldID, id),
		sqlgraph.To(checklistitem.Table, checklistitem.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, workorder.CheckListItemsTable, workorder.CheckListItemsColumn),
	)
	query.sql = sqlgraph.Neighbors(wo.driver.Dialect(), step)

	return query
}

// QueryTechnician queries the technician edge of a WorkOrder.
func (c *WorkOrderClient) QueryTechnician(wo *WorkOrder) *TechnicianQuery {
	query := &TechnicianQuery{config: c.config}
	id := wo.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workorder.Table, workorder.FieldID, id),
		sqlgraph.To(technician.Table, technician.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, workorder.TechnicianTable, workorder.TechnicianColumn),
	)
	query.sql = sqlgraph.Neighbors(wo.driver.Dialect(), step)

	return query
}

// QueryProject queries the project edge of a WorkOrder.
func (c *WorkOrderClient) QueryProject(wo *WorkOrder) *ProjectQuery {
	query := &ProjectQuery{config: c.config}
	id := wo.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workorder.Table, workorder.FieldID, id),
		sqlgraph.To(project.Table, project.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, workorder.ProjectTable, workorder.ProjectColumn),
	)
	query.sql = sqlgraph.Neighbors(wo.driver.Dialect(), step)

	return query
}

// WorkOrderDefinitionClient is a client for the WorkOrderDefinition schema.
type WorkOrderDefinitionClient struct {
	config
}

// NewWorkOrderDefinitionClient returns a client for the WorkOrderDefinition from the given config.
func NewWorkOrderDefinitionClient(c config) *WorkOrderDefinitionClient {
	return &WorkOrderDefinitionClient{config: c}
}

// Create returns a create builder for WorkOrderDefinition.
func (c *WorkOrderDefinitionClient) Create() *WorkOrderDefinitionCreate {
	return &WorkOrderDefinitionCreate{config: c.config}
}

// Update returns an update builder for WorkOrderDefinition.
func (c *WorkOrderDefinitionClient) Update() *WorkOrderDefinitionUpdate {
	return &WorkOrderDefinitionUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *WorkOrderDefinitionClient) UpdateOne(wod *WorkOrderDefinition) *WorkOrderDefinitionUpdateOne {
	return c.UpdateOneID(wod.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *WorkOrderDefinitionClient) UpdateOneID(id string) *WorkOrderDefinitionUpdateOne {
	return &WorkOrderDefinitionUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for WorkOrderDefinition.
func (c *WorkOrderDefinitionClient) Delete() *WorkOrderDefinitionDelete {
	return &WorkOrderDefinitionDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *WorkOrderDefinitionClient) DeleteOne(wod *WorkOrderDefinition) *WorkOrderDefinitionDeleteOne {
	return c.DeleteOneID(wod.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *WorkOrderDefinitionClient) DeleteOneID(id string) *WorkOrderDefinitionDeleteOne {
	return &WorkOrderDefinitionDeleteOne{c.Delete().Where(workorderdefinition.ID(id))}
}

// Create returns a query builder for WorkOrderDefinition.
func (c *WorkOrderDefinitionClient) Query() *WorkOrderDefinitionQuery {
	return &WorkOrderDefinitionQuery{config: c.config}
}

// Get returns a WorkOrderDefinition entity by its id.
func (c *WorkOrderDefinitionClient) Get(ctx context.Context, id string) (*WorkOrderDefinition, error) {
	return c.Query().Where(workorderdefinition.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *WorkOrderDefinitionClient) GetX(ctx context.Context, id string) *WorkOrderDefinition {
	wod, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return wod
}

// QueryType queries the type edge of a WorkOrderDefinition.
func (c *WorkOrderDefinitionClient) QueryType(wod *WorkOrderDefinition) *WorkOrderTypeQuery {
	query := &WorkOrderTypeQuery{config: c.config}
	id := wod.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workorderdefinition.Table, workorderdefinition.FieldID, id),
		sqlgraph.To(workordertype.Table, workordertype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, workorderdefinition.TypeTable, workorderdefinition.TypeColumn),
	)
	query.sql = sqlgraph.Neighbors(wod.driver.Dialect(), step)

	return query
}

// QueryProjectType queries the project_type edge of a WorkOrderDefinition.
func (c *WorkOrderDefinitionClient) QueryProjectType(wod *WorkOrderDefinition) *ProjectTypeQuery {
	query := &ProjectTypeQuery{config: c.config}
	id := wod.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workorderdefinition.Table, workorderdefinition.FieldID, id),
		sqlgraph.To(projecttype.Table, projecttype.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, workorderdefinition.ProjectTypeTable, workorderdefinition.ProjectTypeColumn),
	)
	query.sql = sqlgraph.Neighbors(wod.driver.Dialect(), step)

	return query
}

// WorkOrderTypeClient is a client for the WorkOrderType schema.
type WorkOrderTypeClient struct {
	config
}

// NewWorkOrderTypeClient returns a client for the WorkOrderType from the given config.
func NewWorkOrderTypeClient(c config) *WorkOrderTypeClient {
	return &WorkOrderTypeClient{config: c}
}

// Create returns a create builder for WorkOrderType.
func (c *WorkOrderTypeClient) Create() *WorkOrderTypeCreate {
	return &WorkOrderTypeCreate{config: c.config}
}

// Update returns an update builder for WorkOrderType.
func (c *WorkOrderTypeClient) Update() *WorkOrderTypeUpdate {
	return &WorkOrderTypeUpdate{config: c.config}
}

// UpdateOne returns an update builder for the given entity.
func (c *WorkOrderTypeClient) UpdateOne(wot *WorkOrderType) *WorkOrderTypeUpdateOne {
	return c.UpdateOneID(wot.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *WorkOrderTypeClient) UpdateOneID(id string) *WorkOrderTypeUpdateOne {
	return &WorkOrderTypeUpdateOne{config: c.config, id: id}
}

// Delete returns a delete builder for WorkOrderType.
func (c *WorkOrderTypeClient) Delete() *WorkOrderTypeDelete {
	return &WorkOrderTypeDelete{config: c.config}
}

// DeleteOne returns a delete builder for the given entity.
func (c *WorkOrderTypeClient) DeleteOne(wot *WorkOrderType) *WorkOrderTypeDeleteOne {
	return c.DeleteOneID(wot.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *WorkOrderTypeClient) DeleteOneID(id string) *WorkOrderTypeDeleteOne {
	return &WorkOrderTypeDeleteOne{c.Delete().Where(workordertype.ID(id))}
}

// Create returns a query builder for WorkOrderType.
func (c *WorkOrderTypeClient) Query() *WorkOrderTypeQuery {
	return &WorkOrderTypeQuery{config: c.config}
}

// Get returns a WorkOrderType entity by its id.
func (c *WorkOrderTypeClient) Get(ctx context.Context, id string) (*WorkOrderType, error) {
	return c.Query().Where(workordertype.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *WorkOrderTypeClient) GetX(ctx context.Context, id string) *WorkOrderType {
	wot, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return wot
}

// QueryWorkOrders queries the work_orders edge of a WorkOrderType.
func (c *WorkOrderTypeClient) QueryWorkOrders(wot *WorkOrderType) *WorkOrderQuery {
	query := &WorkOrderQuery{config: c.config}
	id := wot.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workordertype.Table, workordertype.FieldID, id),
		sqlgraph.To(workorder.Table, workorder.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, workordertype.WorkOrdersTable, workordertype.WorkOrdersColumn),
	)
	query.sql = sqlgraph.Neighbors(wot.driver.Dialect(), step)

	return query
}

// QueryPropertyTypes queries the property_types edge of a WorkOrderType.
func (c *WorkOrderTypeClient) QueryPropertyTypes(wot *WorkOrderType) *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: c.config}
	id := wot.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workordertype.Table, workordertype.FieldID, id),
		sqlgraph.To(propertytype.Table, propertytype.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, workordertype.PropertyTypesTable, workordertype.PropertyTypesColumn),
	)
	query.sql = sqlgraph.Neighbors(wot.driver.Dialect(), step)

	return query
}

// QueryDefinitions queries the definitions edge of a WorkOrderType.
func (c *WorkOrderTypeClient) QueryDefinitions(wot *WorkOrderType) *WorkOrderDefinitionQuery {
	query := &WorkOrderDefinitionQuery{config: c.config}
	id := wot.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workordertype.Table, workordertype.FieldID, id),
		sqlgraph.To(workorderdefinition.Table, workorderdefinition.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, workordertype.DefinitionsTable, workordertype.DefinitionsColumn),
	)
	query.sql = sqlgraph.Neighbors(wot.driver.Dialect(), step)

	return query
}

// QueryCheckListCategories queries the check_list_categories edge of a WorkOrderType.
func (c *WorkOrderTypeClient) QueryCheckListCategories(wot *WorkOrderType) *CheckListCategoryQuery {
	query := &CheckListCategoryQuery{config: c.config}
	id := wot.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workordertype.Table, workordertype.FieldID, id),
		sqlgraph.To(checklistcategory.Table, checklistcategory.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, workordertype.CheckListCategoriesTable, workordertype.CheckListCategoriesColumn),
	)
	query.sql = sqlgraph.Neighbors(wot.driver.Dialect(), step)

	return query
}

// QueryCheckListDefinitions queries the check_list_definitions edge of a WorkOrderType.
func (c *WorkOrderTypeClient) QueryCheckListDefinitions(wot *WorkOrderType) *CheckListItemDefinitionQuery {
	query := &CheckListItemDefinitionQuery{config: c.config}
	id := wot.id()
	step := sqlgraph.NewStep(
		sqlgraph.From(workordertype.Table, workordertype.FieldID, id),
		sqlgraph.To(checklistitemdefinition.Table, checklistitemdefinition.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, workordertype.CheckListDefinitionsTable, workordertype.CheckListDefinitionsColumn),
	)
	query.sql = sqlgraph.Neighbors(wot.driver.Dialect(), step)

	return query
}
