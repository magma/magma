// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewer

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent/hook"

	"github.com/facebookincubator/symphony/graph/ent"
)

// UpdateCurrentUser is a hook to update user ent cache in Viewer object
func UpdateCurrentUser() ent.Hook {
	chain := hook.NewChain(
		hook.Reject(ent.OpUpdate),
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
					value, err := next.Mutate(ctx, m)
					if err != nil {
						return value, err
					}
					v := FromContext(ctx)
					if v == nil {
						return value, err
					}
					u := value.(*ent.User)
					if u.ID == v.User().ID {
						v.mu.Lock()
						v.user = u
						v.mu.Unlock()
					}
					return value, nil
				})
			},
			ent.OpUpdateOne,
		),
	)
	return chain.Hook()
}
