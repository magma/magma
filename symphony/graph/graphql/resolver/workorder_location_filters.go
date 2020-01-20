// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"

	"github.com/pkg/errors"
)

func (r *queryResolver) handleWOLocationFilter(q *ent.WorkOrderQuery, filter *models.WorkOrderFilterInput) (*ent.WorkOrderQuery, error) {
	if filter.FilterType == models.WorkOrderFilterTypeLocationInst {
		return woLocationFilter(q, filter)
	}
	return nil, errors.Errorf("filter type is not supported: %s", filter.FilterType)
}

func woLocationFilter(q *ent.WorkOrderQuery, filter *models.WorkOrderFilterInput) (*ent.WorkOrderQuery, error) {
	if filter.Operator == models.FilterOperatorIsOneOf {
		if filter.MaxDepth == nil {
			return nil, errors.New("max depth not supplied to location filter")
		}
		var ps []predicate.WorkOrder
		for _, lid := range filter.IDSet {
			ps = append(ps, workorder.HasLocationWith(resolverutil.BuildLocationAncestorFilter(lid, 1, *filter.MaxDepth)))
		}
		return q.Where(workorder.Or(ps...)), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}
