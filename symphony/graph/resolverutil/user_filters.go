// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolverutil

import (
	"strings"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/pkg/errors"
)

func handleUserFilter(q *ent.UserQuery, filter *models.UserFilterInput) (*ent.UserQuery, error) {
	switch filter.FilterType {
	case models.UserFilterTypeUserName:
		return userFilter(q, filter)
	default:
		return nil, errors.Errorf("filter type is not supported: %s", filter.FilterType)
	}
}

func userFilter(q *ent.UserQuery, filter *models.UserFilterInput) (*ent.UserQuery, error) {
	if filter.Operator == models.FilterOperatorContains {
		terms := strings.Split(*filter.StringValue, " ")
		qp := user.And()
		for _, s := range terms {
			qp = user.And(qp, userStringPredicate(s))
		}
		return q.Where(qp), nil
	}
	return nil, errors.Errorf("operation %q not supported", filter.Operator)
}

func userStringPredicate(s string) predicate.User {
	return user.Or(user.EmailContainsFold(s),
		user.FirstNameContainsFold(s),
		user.LastNameContainsFold(s))
}
