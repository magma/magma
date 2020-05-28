// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolverutil

import (
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/permissionspolicy"

	"github.com/pkg/errors"
)

func handlePermissionsPolicyFilter(
	q *ent.PermissionsPolicyQuery,
	filter *models.PermissionsPolicyFilterInput,
) (*ent.PermissionsPolicyQuery, error) {
	if filter.FilterType == models.PermissionsPolicyFilterTypePermissionsPolicyName {
		return permissionsPolicyFilter(q, filter)
	}
	return nil, errors.Errorf("filter type is not supported: %s", filter.FilterType)
}

func permissionsPolicyFilter(
	q *ent.PermissionsPolicyQuery,
	filter *models.PermissionsPolicyFilterInput,
) (*ent.PermissionsPolicyQuery, error) {
	switch {
	case filter.Operator == models.FilterOperatorIs && filter.StringValue != nil:
		return q.Where(permissionspolicy.NameEqualFold(*filter.StringValue)), nil
	case filter.Operator == models.FilterOperatorContains && filter.StringValue != nil:
		return q.Where(permissionspolicy.NameContainsFold(*filter.StringValue)), nil
	}
	return nil, errors.Errorf("operation %s is not supported with value of %#v", filter.Operator, filter.StringValue)
}
