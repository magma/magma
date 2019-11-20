// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	"golang.org/x/xerrors"
)

type (
	customerResolver struct{}
)

func (customerResolver) TotalCount(ctx context.Context, obj *models.CustomerConnection) (int, error) {
	count, err := ent.FromContext(ctx).Customer.Query().Count(ctx)
	if err != nil {
		return 0, xerrors.Errorf("querying customer count: %w", err)
	}
	return count, nil
}
