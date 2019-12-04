// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/cloud/log"
	"github.com/facebookincubator/symphony/cloud/orc8r"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/viewer"
)

type resolver struct {
	log         log.Logger
	withTx      bool
	orc8rClient *orc8r.Client
}

// User information of the graphql request initiator
type User struct {
	email string
}

// New creates a graphql resolver.
func New(logger log.Logger, opts ...ResolveOption) (generated.ResolverRoot, error) {
	r := &resolver{log: logger, withTx: true}
	for _, opt := range opts {
		opt(r)
	}
	return r, nil
}

func (r resolver) ClientFrom(ctx context.Context) *ent.Client {
	client := ent.FromContext(ctx)
	if client == nil {
		panic("no ClientFrom attached to context")
	}
	return client
}

func (r resolver) User(ctx context.Context) User {
	return User{viewer.FromContext(ctx).User}
}

func (r resolver) Equipment() generated.EquipmentResolver {
	return equipmentResolver{r}
}

func (resolver) EquipmentPort() generated.EquipmentPortResolver {
	return equipmentPortResolver{}
}

func (resolver) EquipmentPosition() generated.EquipmentPositionResolver {
	return equipmentPositionResolver{}
}

func (resolver) EquipmentPortDefinition() generated.EquipmentPortDefinitionResolver {
	return equipmentPortDefinitionResolver{}
}

func (resolver) EquipmentPortType() generated.EquipmentPortTypeResolver {
	return equipmentPortTypeResolver{}
}

func (resolver) EquipmentType() generated.EquipmentTypeResolver {
	return equipmentTypeResolver{}
}

func (resolver) File() generated.FileResolver {
	return fileResolver{}
}

func (resolver) Link() generated.LinkResolver {
	return linkResolver{}
}

func (resolver) Location() generated.LocationResolver {
	return locationResolver{}
}

func (resolver) LocationType() generated.LocationTypeResolver {
	return locationTypeResolver{}
}

func (resolver) FloorPlan() generated.FloorPlanResolver {
	return floorPlanResolver{}
}

func (r resolver) Mutation() generated.MutationResolver {
	mr := mutationResolver{r}
	if r.withTx {
		return txResolver{mr}
	}
	return mr
}

func (r resolver) Query() generated.QueryResolver {
	return queryResolver{r}
}

func (resolver) WorkOrder() generated.WorkOrderResolver {
	return workOrderResolver{}
}

func (resolver) WorkOrderType() generated.WorkOrderTypeResolver {
	return workOrderTypeResolver{}
}

func (resolver) WorkOrderDefinition() generated.WorkOrderDefinitionResolver {
	return workOrderDefinitionResolver{}
}

func (r resolver) Survey() generated.SurveyResolver {
	return surveyResolver{}
}

func (r resolver) SurveyQuestion() generated.SurveyQuestionResolver {
	return surveyQuestionResolver{}
}

func (r resolver) SurveyTemplateQuestion() generated.SurveyTemplateQuestionResolver {
	return surveyTemplateQuestionResolver{}
}

func (r resolver) SurveyTemplateCategory() generated.SurveyTemplateCategoryResolver {
	return surveyTemplateCategoryResolver{}
}

func (r resolver) SurveyCellScan() generated.SurveyCellScanResolver {
	return surveyCellScanResolver{}
}

func (r resolver) SurveyWiFiScan() generated.SurveyWiFiScanResolver {
	return surveyWiFiScanResolver{}
}

func (r resolver) PropertyType() generated.PropertyTypeResolver {
	return propertyTypeResolver{}
}

func (r resolver) Property() generated.PropertyResolver {
	return propertyResolver{}
}

func (resolver) Service() generated.ServiceResolver {
	return serviceResolver{}
}

func (r resolver) ServiceType() generated.ServiceTypeResolver {
	return serviceTypeResolver{}
}

func (resolver) Project() generated.ProjectResolver {
	return projectResolver{}
}

func (resolver) ProjectType() generated.ProjectTypeResolver {
	return projectTypeResolver{}
}

func (resolver) ProjectTypeConnection() generated.ProjectTypeConnectionResolver {
	return projectTypeResolver{}
}

func (resolver) CustomerConnection() generated.CustomerConnectionResolver {
	return customerResolver{}
}

func (resolver) CheckListItem() generated.CheckListItemResolver {
	return checkListItemResolver{}
}

func (resolver) CheckListItemDefinition() generated.CheckListItemDefinitionResolver {
	return checkListItemDefinitionResolver{}
}

// ResolveOption allows for managing resolver configuration using functional options.
type ResolveOption func(*resolver)

// WithTransaction if set to true, will wraps the mutation with transaction.
func WithTransaction(b bool) ResolveOption {
	return func(r *resolver) {
		r.withTx = b
	}
}

// WithOrc8r is used to provide orchestrator configs to the resolver
func WithOrc8r(orc8r *orc8r.Client) ResolveOption {
	return func(r *resolver) {
		r.orc8rClient = orc8r
	}
}
