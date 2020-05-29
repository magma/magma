// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewer

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/hook"
)

// UpdateCurrentUser updates user stored in viewer.
func UpdateCurrentUser() ent.Hook {
	hk := func(next ent.Mutator) ent.Mutator {
		return hook.UserFunc(func(ctx context.Context, m *ent.UserMutation) (ent.Value, error) {
			value, err := next.Mutate(ctx, m)
			if err != nil {
				return value, err
			}
			if v, ok := FromContext(ctx).(*UserViewer); ok && v.User().ID == value.(*ent.User).ID {
				v.user.Store(value)
			}
			return value, nil
		})
	}
	return hook.NewChain(
		hook.On(hk, ent.OpUpdateOne),
		hook.Reject(ent.OpUpdate),
	).Hook()
}
