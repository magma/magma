// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"encoding/json"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent"
)

type reportFilterResolver struct{}

func (reportFilterResolver) Entity(_ context.Context, rf *ent.ReportFilter) (models.FilterEntity, error) {
	return models.FilterEntity(rf.Entity), nil
}

func (r reportFilterResolver) Filters(_ context.Context, rf *ent.ReportFilter) ([]*models.GeneralFilter, error) {
	var filters []*models.GeneralFilter
	err := json.Unmarshal([]byte(rf.Filters), &filters)
	return filters, err
}
