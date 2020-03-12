// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewer

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/user"
)

// UserFromContext returns the ent user using authID in context.
func UserFromContext(ctx context.Context) (*ent.User, error) {
	client := ent.FromContext(ctx)
	v := FromContext(ctx)
	return client.User.Query().Where(user.AuthID(v.User)).Only(ctx)
}
