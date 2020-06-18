// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/pubsub"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"go.uber.org/zap"
)

// Eventer generates events from mutations.
type Eventer struct {
	Logger  log.Logger
	Emitter pubsub.Emitter
}

// HookTo hooks eventer to ent client.
func (e *Eventer) HookTo(client *ent.Client) {
	client.Use(e.logHook())
	client.WorkOrder.Use(e.workOrderHook())
}

func (e *Eventer) emit(ctx context.Context, name string, value interface{}) {
	emit := func(err error) {
		if err != nil {
			return
		}
		logger := e.Logger.For(ctx).With(zap.String("name", name))
		body, err := pubsub.Marshal(value)
		if err != nil {
			logger.Warn("cannot marshal event value", zap.Error(err))
			return
		}
		if err := e.Emitter.Emit(ctx, viewer.FromContext(ctx).Tenant(), name, body); err != nil {
			logger.Warn("cannot emit event", zap.Error(err))
		}
		logger.Debug("emitting event")
	}
	if tx := ent.TxFromContext(ctx); tx != nil {
		tx.OnCommit(emit)
	} else {
		emit(nil)
	}
}
