// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/log"
	"go.uber.org/zap"
)

type subscriptionResolver struct{ resolver }

func (r subscriptionResolver) SubscribeAndListen(ctx context.Context, name string, handler event.Handler) {
	logger := r.logger.For(ctx)
	err := event.SubscribeAndListen(ctx, event.ListenerConfig{
		Subscriber: r.event.Subscriber,
		Logger:     logger,
		Tenant:     viewer.FromContext(ctx).Tenant,
		Events:     []string{name},
		Handler:    handler,
	})
	logger.Debug("subscription termination", zap.Error(err))
}

func (r subscriptionResolver) workOrderAddedDone(ctx context.Context, name string) (<-chan *ent.WorkOrder, error) {
	var (
		client = r.ClientFrom(ctx).WorkOrder
		events = make(chan *ent.WorkOrder, 1)
	)
	go func() {
		defer close(events)
		r.SubscribeAndListen(ctx, name,
			event.HandlerFunc(func(_ context.Context, _ string, body []byte) error {
				var wo *ent.WorkOrder
				if err := gob.NewDecoder(bytes.NewReader(body)).Decode(&wo); err != nil {
					return fmt.Errorf("cannot unmarshal work order: %w", err)
				}
				events <- client.Instantiate(wo)
				return nil
			}),
		)
	}()
	return events, nil
}

func (r subscriptionResolver) WorkOrderAdded(ctx context.Context) (<-chan *ent.WorkOrder, error) {
	return r.workOrderAddedDone(ctx, event.WorkOrderAdded)
}

func (r subscriptionResolver) WorkOrderDone(ctx context.Context) (<-chan *ent.WorkOrder, error) {
	return r.workOrderAddedDone(ctx, event.WorkOrderDone)
}

// eventResolver wraps a mutation resolver and emits events on mutations.
type eventResolver struct {
	generated.MutationResolver
	emitter event.Emitter
	logger  log.Logger
}

func (r eventResolver) event(ctx context.Context, name string, value interface{}) {
	logger := r.logger.For(ctx).With(zap.String("name", name))
	var body bytes.Buffer
	if err := gob.NewEncoder(&body).Encode(value); err != nil {
		logger.Warn("cannot marshal event value", zap.Error(err))
		return
	}
	if err := r.emitter.Emit(ctx, viewer.FromContext(ctx).Tenant, name, body.Bytes()); err != nil {
		logger.Warn("cannot emit event", zap.Error(err))
	}
}

func (r eventResolver) AddWorkOrder(
	ctx context.Context, input models.AddWorkOrderInput,
) (wo *ent.WorkOrder, err error) {
	defer func() {
		if err != nil {
			return
		}
		r.event(ctx, event.WorkOrderAdded, wo)
		if wo.Status == models.WorkOrderStatusDone.String() {
			r.event(ctx, event.WorkOrderDone, wo)
		}
	}()
	return r.MutationResolver.AddWorkOrder(ctx, input)
}

func (r eventResolver) EditWorkOrder(
	ctx context.Context, input models.EditWorkOrderInput,
) (wo *ent.WorkOrder, err error) {
	var exist bool
	if exist, err = ent.FromContext(ctx).
		WorkOrder.Query().
		Where(
			workorder.ID(input.ID),
			workorder.StatusNEQ(
				models.WorkOrderStatusDone.String(),
			),
		).
		Exist(ctx); err != nil {
		return nil, fmt.Errorf("querying work order existence: %w", err)
	}
	if exist {
		defer func() {
			if err == nil && wo.Status == models.WorkOrderStatusDone.String() {
				r.event(ctx, event.WorkOrderDone, wo)
			}
		}()
	}
	return r.MutationResolver.EditWorkOrder(ctx, input)
}
