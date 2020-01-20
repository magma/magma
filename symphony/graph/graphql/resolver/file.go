// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

type fileResolver struct{}

func (fileResolver) FileType(ctx context.Context, obj *ent.File) (*models.FileType, error) {
	ft := models.FileType(obj.Type)
	return &ft, nil
}
