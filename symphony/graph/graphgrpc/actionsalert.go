// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphgrpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/facebookincubator/symphony/graph/graphgrpc/schema"
	"github.com/facebookincubator/symphony/pkg/actions"
	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/facebookincubator/symphony/pkg/actions/executor"
	"github.com/facebookincubator/symphony/pkg/ent/user"
)

const ActionsAlertServiceName = "ActionsAlertService"

type (
	// ActionsAlertService is an ActionsAlertService
	ActionsAlertService struct {
		actionsProvider ActionsProvider
	}

	// ActionsProvider returns an actions client given a context and tenant
	ActionsProvider func(ctx context.Context, tenantID string) (*actions.Client, error)
)

//NewActionsAlertService returns a new ActionsAlertService
func NewActionsAlertService(provider ActionsProvider) ActionsAlertService {
	return ActionsAlertService{provider}
}

// Receive an alert payload and execute the triggered actions
func (s ActionsAlertService) Trigger(ctx context.Context, payload *schema.AlertPayload) (*schema.ExecutionResult, error) {
	triggerPayload := make(map[string]interface{})
	triggerPayload["networkID"] = payload.NetworkID
	for key, val := range payload.Labels {
		triggerPayload[key] = val
	}
	idToPayload := map[core.TriggerID]map[string]interface{}{core.MagmaAlertTriggerID: triggerPayload}

	client, err := s.actionsProvider(ctx, payload.TenantID)
	if err != nil {
		return &schema.ExecutionResult{}, status.Error(codes.Internal, "error getting tenant client")
	}
	ctx, err = CreateServiceContext(ctx, payload.TenantID, ActionsAlertServiceName, user.RoleOWNER)
	if err != nil {
		return &schema.ExecutionResult{}, status.Error(codes.Internal, "error getting service context")
	}
	res := client.Execute(ctx, "", idToPayload)

	return executorResultToMessage(res), nil
}

func executorResultToMessage(res executor.ExecutionResult) *schema.ExecutionResult {
	var successStrings []string
	for _, id := range res.Successes {
		successStrings = append(successStrings, string(id))
	}
	var errors []*schema.ExecutionError
	for _, err := range res.Errors {
		errors = append(errors, &schema.ExecutionError{Id: string(err.ID), Err: err.Error})
	}

	return &schema.ExecutionResult{
		Successes: successStrings,
		Errors:    errors,
	}
}
