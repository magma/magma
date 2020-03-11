// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/viewer"
)

type viewerResolver struct{}

func (viewerResolver) User(ctx context.Context, obj *viewer.Viewer) (*ent.User, error) {
	return viewer.UserFromContext(ctx)
}
