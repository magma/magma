package resolverutil

import (
	"github.com/facebookincubator/symphony/graph/ent"
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
		return q.Where(user.Or(user.EmailContainsFold(*filter.StringValue),
			user.FirstNameContainsFold(*filter.StringValue),
			user.LastNameContainsFold(*filter.StringValue))), nil
	}
	return nil, errors.Errorf("operation %q not supported", filter.Operator)
}
