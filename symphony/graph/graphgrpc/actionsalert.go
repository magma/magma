// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphgrpc

import (
	"context"
	"github.com/facebookincubator/symphony/pkg/actions"
	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/golang/protobuf/ptypes/empty"
)

type (
	// ActionsAlertService is an ActionsAlertService
	ActionsAlertService struct {
		actionsProvider ActionsProvider
	}

	// ActionsProvider returns an actions client given a context and tenant
	ActionsProvider func(ctx context.Context, tenantID string) (*actions.Client, error)
)

// NewActionsAlertService returns a new ActionsAlertService
func NewActionsAlertService(provider ActionsProvider) ActionsAlertService {
	return ActionsAlertService{provider}
}

// Receive an alert payload and execute the triggered actions
func (s ActionsAlertService) Trigger(ctx context.Context, payload *AlertPayload) (*empty.Empty, error) {
	triggerPayload := make(map[string]interface{})
	triggerPayload["networkID"] = payload.NetworkID
	for key, val := range payload.Labels {
		triggerPayload[key] = val
	}
	idToPayload := map[core.TriggerID]map[string]interface{}{core.MagmaAlertTriggerID: triggerPayload}

	client, err := s.actionsProvider(ctx, payload.TenantID)
	if err != nil {
		return &empty.Empty{}, err
	}
	client.Execute(context.Background(), "", idToPayload)

	return &empty.Empty{}, nil
}
