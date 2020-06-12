// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/pubsub"
	"github.com/facebookincubator/symphony/pkg/viewer"
)

func (e *Eventer) logHook() ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return e.hookWithLog(func(ctx context.Context, entry pubsub.LogEntry) error {
			v := viewer.FromContext(ctx)
			if v == nil ||
				!v.Features().Enabled(viewer.FeatureGraphEventLogging) {
				return nil
			}
			e.emit(ctx, pubsub.EntMutation, entry)
			return nil
		}, next)
	}
}
