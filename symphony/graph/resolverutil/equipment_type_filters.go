// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolverutil

import (
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/equipment"
	"github.com/facebookincubator/symphony/pkg/ent/equipmenttype"

	"github.com/pkg/errors"
)

func handleEquipmentTypeFilter(q *ent.EquipmentQuery, filter *models.EquipmentFilterInput) (*ent.EquipmentQuery, error) {
	if filter.FilterType == models.EquipmentFilterTypeEquipmentType {
		return equipmentTypeFilter(q, filter)
	}
	return nil, errors.Errorf("filter type is not supported: %s", filter.FilterType)
}

func equipmentTypeFilter(q *ent.EquipmentQuery, filter *models.EquipmentFilterInput) (*ent.EquipmentQuery, error) {
	if filter.Operator == models.FilterOperatorIsOneOf {
		return q.Where(equipment.HasTypeWith(equipmenttype.IDIn(filter.IDSet...))), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}
