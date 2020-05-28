// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"time"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"go.uber.org/zap"
)

// LogEntry holds an information on a single ent mutation that happened
type LogEntry struct {
	UserName  string    `json:"user_name"`
	UserID    *int      `json:"user_id"`
	Time      time.Time `json:"time"`
	Operation ent.Op    `json:"operation"`
	PrevState *ent.Node `json:"prevState"`
	CurrState *ent.Node `json:"currState"`
}

func (e *Eventer) hookWithLog(handler func(context.Context, LogEntry) error, next ent.Mutator) ent.Mutator {
	return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
		v := viewer.FromContext(ctx)
		if m.Op().Is(ent.OpDelete) || m.Op().Is(ent.OpUpdate) || v == nil {
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
				node, err := prevNoder.Node(ctx)
				if err != nil {
					if !ent.IsNotFound(err) {
						e.Logger.For(ctx).Error("query mutation previous value", zap.Error(err))
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
				node, err := currNoder.Node(ctx)
				if err != nil {
					e.Logger.For(ctx).Error("query mutation current value", zap.Error(err))
					return value, err
				}
				entry.CurrState = node
			}
		}
		err = handler(ctx, entry)
		if err != nil {
			return value, err
		}
		return value, nil
	})
}
