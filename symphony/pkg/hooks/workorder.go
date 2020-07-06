// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hooks

import (
	"context"
	"time"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/hook"
	"github.com/facebookincubator/symphony/pkg/ent/workorder"
)

// WorkOrderCloseDateHook modifies work order close date from status.
func WorkOrderCloseDateHook() ent.Hook {
	hk := func(next ent.Mutator) ent.Mutator {
		return hook.WorkOrderFunc(func(ctx context.Context, mutation *ent.WorkOrderMutation) (ent.Value, error) {
			status, exists := mutation.Status()
			if !exists {
				return next.Mutate(ctx, mutation)
			}
			if _, exists := mutation.CloseDate(); exists {
				return next.Mutate(ctx, mutation)
			}
			if mutation.Op().Is(ent.OpUpdateOne) {
				switch oldStatus, err := mutation.OldStatus(ctx); {
				case err != nil:
					return nil, err
				case status == oldStatus:
					return next.Mutate(ctx, mutation)
				}
			}
			switch status {
			case workorder.StatusDONE:
				mutation.SetCloseDate(time.Now())
			default:
				mutation.ClearCloseDate()
			}
			return next.Mutate(ctx, mutation)
		})
	}
	return hook.On(hk, ent.OpCreate|ent.OpUpdateOne|ent.OpUpdate)
}
