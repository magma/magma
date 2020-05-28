// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"

	"github.com/AlekSi/pointer"

	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"go.uber.org/zap"
)

type subscriptionResolver struct{ resolver }

func (r subscriptionResolver) subscribeAndListen(ctx context.Context, name string, handler event.Handler) {
	logger := r.logger.For(ctx)
	err := event.SubscribeAndListen(ctx, event.ListenerConfig{
		Subscriber: r.event.Subscriber,
		Logger:     logger,
		Tenant:     pointer.ToString(viewer.FromContext(ctx).Tenant()),
		Events:     []string{name},
		Handler:    handler,
	})
	logger.Info("subscription termination", zap.Error(err))
}

func (r subscriptionResolver) workOrderAddedDone(ctx context.Context, name string) (<-chan *ent.WorkOrder, error) {
	var (
		client = r.ClientFrom(ctx).WorkOrder
		events = make(chan *ent.WorkOrder, 1)
	)
	go func() {
		defer close(events)
		r.subscribeAndListen(ctx, name,
			event.HandlerFunc(func(_ context.Context, _, _ string, body []byte) error {
				var wo *ent.WorkOrder
				if err := event.Unmarshal(body, &wo); err != nil {
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
