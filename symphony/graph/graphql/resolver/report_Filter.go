// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"encoding/json"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

type (
	reportFilterResolver struct{}
)

func (r reportFilterResolver) Entity(ctx context.Context, obj *ent.ReportFilter) (models.FilterEntity, error) {
	return models.FilterEntity(obj.Entity), nil
}

func (r reportFilterResolver) Filters(ctx context.Context, obj *ent.ReportFilter) ([]*models.GeneralFilter, error) {
	var f []*models.GeneralFilter
	err := json.Unmarshal([]byte(obj.Filters), &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}
