// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"time"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/viewer"

	"go.uber.org/zap"
)

func LogHook(handler func(context.Context, LogEntry) error, logger log.Logger) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			v := viewer.FromContext(ctx)
			if m.Op().Is(ent.OpDelete|ent.OpUpdate) || v == nil {
				return next.Mutate(ctx, m)
			}
			entry := LogEntry{
				UserName:  v.Name(),
				Operation: m.Op(),
				Time:      time.Now(),
			}
			if v, ok := v.(*viewer.UserViewer); ok {
				entry.UserID = &v.User().ID
			}
			if !m.Op().Is(ent.OpCreate) {
				if prevNoder, ok := m.(ent.Noder); ok {
					node, err := prevNoder.Node(privacy.DecisionContext(ctx, privacy.Allow))
					if err != nil {
						if !ent.IsNotFound(err) {
							logger.For(ctx).Error("query mutation previous value", zap.Error(err))
						}
						return next.Mutate(ctx, m)
					}
					entry.PrevState = node
				}
			}
			value, err := next.Mutate(ctx, m)
			if err != nil {
				return value, err
			}
			if !m.Op().Is(ent.OpDeleteOne) {
				if currNoder, ok := value.(ent.Noder); ok {
					node, err := currNoder.Node(privacy.DecisionContext(ctx, privacy.Allow))
					if err != nil {
						logger.For(ctx).Error("query mutation current value", zap.Error(err))
						return value, err
					}
					entry.CurrState = node
				}
			}
			err = handler(ctx, entry)
			if err != nil {
				return nil, err
			}
			return value, nil
		})
	}
}
