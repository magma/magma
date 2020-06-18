// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent"
)

type fileResolver struct{}

func (fileResolver) FileType(_ context.Context, file *ent.File) (*models.FileType, error) {
	ft := models.FileType(file.Type)
	return &ft, nil
}
