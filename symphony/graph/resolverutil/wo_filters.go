// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolverutil

import (
	"strconv"
	"time"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	"github.com/pkg/errors"
)

func handleWorkOrderFilter(q *ent.WorkOrderQuery, filter *models.WorkOrderFilterInput) (*ent.WorkOrderQuery, error) {
	if filter.FilterType == models.WorkOrderFilterTypeWorkOrderName {
		return nameFilter(q, filter)
	}
	if filter.FilterType == models.WorkOrderFilterTypeWorkOrderStatus {
		return statusFilter(q, filter)
	}
	if filter.FilterType == models.WorkOrderFilterTypeWorkOrderOwner {
		return ownerFilter(q, filter)
	}
	if filter.FilterType == models.WorkOrderFilterTypeWorkOrderType {
		return typeFilter(q, filter)
	}
	if filter.FilterType == models.WorkOrderFilterTypeWorkOrderAssignee {
		return assigneeFilter(q, filter)
	}
	if filter.FilterType == models.WorkOrderFilterTypeWorkOrderCreationDate {
		return creationDateFilter(q, filter)
	}
	if filter.FilterType == models.WorkOrderFilterTypeWorkOrderInstallDate {
		return installDateFilter(q, filter)
	}
	if filter.FilterType == models.WorkOrderFilterTypeWorkOrderLocationInst {
		return locationInstFilter(q, filter)
	}
	if filter.FilterType == models.WorkOrderFilterTypeWorkOrderPriority {
		return priorityFilter(q, filter)
	}
	return nil, errors.Errorf("filter type is not supported: %s", filter.FilterType)
}

func nameFilter(q *ent.WorkOrderQuery, filter *models.WorkOrderFilterInput) (*ent.WorkOrderQuery, error) {
	if filter.Operator == models.FilterOperatorContains && filter.StringValue != nil {
		return q.Where(workorder.NameContainsFold(*filter.StringValue)), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

func statusFilter(q *ent.WorkOrderQuery, filter *models.WorkOrderFilterInput) (*ent.WorkOrderQuery, error) {
	if filter.Operator == models.FilterOperatorIsOneOf {
		return q.Where(workorder.StatusIn(filter.StringSet...)), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

func ownerFilter(q *ent.WorkOrderQuery, filter *models.WorkOrderFilterInput) (*ent.WorkOrderQuery, error) {
	if filter.Operator == models.FilterOperatorIsOneOf {
		return q.Where(workorder.OwnerNameIn(filter.StringSet...)), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

func typeFilter(q *ent.WorkOrderQuery, filter *models.WorkOrderFilterInput) (*ent.WorkOrderQuery, error) {
	if filter.Operator == models.FilterOperatorIsOneOf {
		return q.Where(workorder.HasTypeWith(workordertype.IDIn(filter.IDSet...))), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

func assigneeFilter(q *ent.WorkOrderQuery, filter *models.WorkOrderFilterInput) (*ent.WorkOrderQuery, error) {
	if filter.Operator == models.FilterOperatorIsOneOf {
		return q.Where(workorder.AssigneeNameIn(filter.StringSet...)), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

func getStartAndEndOfDay(filterTime string) (*time.Time, *time.Time, error) {
	mtime, err := strconv.ParseInt(filterTime, 10, 64)
	if err != nil {
		return nil, nil, err
	}
	unix := time.Unix(mtime, 0)
	bod := unix.Truncate(time.Hour * 24).UTC()
	eod := bod.Add(time.Hour*24 - 1).UTC()
	return &bod, &eod, nil
}
func creationDateFilter(q *ent.WorkOrderQuery, filter *models.WorkOrderFilterInput) (*ent.WorkOrderQuery, error) {
	bod, eod, err := getStartAndEndOfDay(*filter.StringValue)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing modification time: mtime=%q", *filter.StringValue)
	}
	if filter.Operator == models.FilterOperatorIs {
		return q.Where(workorder.CreationDateGTE(*bod)).Where(workorder.CreationDateLTE(*eod)), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

func installDateFilter(q *ent.WorkOrderQuery, filter *models.WorkOrderFilterInput) (*ent.WorkOrderQuery, error) {
	bod, eod, err := getStartAndEndOfDay(*filter.StringValue)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing modification time: mtime=%q", *filter.StringValue)
	}
	if filter.Operator == models.FilterOperatorIs {
		return q.Where(workorder.InstallDateGTE(*bod)).Where(workorder.InstallDateLTE(*eod)), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

func locationInstFilter(q *ent.WorkOrderQuery, filter *models.WorkOrderFilterInput) (*ent.WorkOrderQuery, error) {
	if filter.Operator == models.FilterOperatorIsOneOf {
		return q.Where(workorder.HasLocationWith(location.IDIn(filter.IDSet...))), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

func priorityFilter(q *ent.WorkOrderQuery, filter *models.WorkOrderFilterInput) (*ent.WorkOrderQuery, error) {
	if filter.Operator == models.FilterOperatorIsOneOf {
		return q.Where(workorder.PriorityIn(filter.StringSet...)), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

func handleWOLocationFilter(q *ent.WorkOrderQuery, filter *models.WorkOrderFilterInput) (*ent.WorkOrderQuery, error) {
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
			ps = append(ps, workorder.HasLocationWith(BuildLocationAncestorFilter(lid, 1, *filter.MaxDepth)))
		}
		return q.Where(workorder.Or(ps...)), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}
