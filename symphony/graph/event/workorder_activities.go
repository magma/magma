// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"strconv"

	"github.com/AlekSi/pointer"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/activity"
	"github.com/facebookincubator/symphony/pkg/ent/workorder"
)

func updateActivitiesOnWOCreate(ctx context.Context, entry *LogEntry) error {
	userID := entry.UserID
	client := ent.FromContext(ctx)

	wo := entry.CurrState
	_, err := client.Activity.Create().
		SetChangedField(activity.ChangedFieldCREATIONDATE).
		SetIsCreate(true).
		SetNewValue(strconv.FormatInt(entry.Time.Unix(), 10)).
		SetNillableAuthorID(userID).
		SetWorkOrderID(wo.ID).
		Save(ctx)
	if err != nil {
		return err
	}

	if assignee, found := findEdge(wo.Edges, workorder.EdgeAssignee); found {
		assgnID := assignee.IDs[0]
		_, err = client.Activity.Create().
			SetChangedField(activity.ChangedFieldASSIGNEE).
			SetIsCreate(true).
			SetNillableAuthorID(userID).
			SetWorkOrderID(wo.ID).
			SetNewValue(strconv.Itoa(assgnID)).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	if owner, found := findEdge(wo.Edges, workorder.EdgeOwner); found {
		ownerID := owner.IDs[0]

		_, err = client.Activity.Create().
			SetChangedField(activity.ChangedFieldOWNER).
			SetIsCreate(true).
			SetNillableAuthorID(userID).
			SetWorkOrderID(wo.ID).
			SetNewValue(strconv.Itoa(ownerID)).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	if st, found := findField(wo.Fields, workorder.FieldStatus); found {
		_, err = client.Activity.Create().
			SetChangedField(activity.ChangedFieldSTATUS).
			SetIsCreate(true).
			SetNillableAuthorID(userID).
			SetWorkOrderID(wo.ID).
			SetNewValue(st.MustGetString()).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	if pri, found := findField(wo.Fields, workorder.FieldPriority); found {
		_, err = client.Activity.Create().
			SetChangedField(activity.ChangedFieldPRIORITY).
			SetIsCreate(true).
			SetNillableAuthorID(userID).
			SetWorkOrderID(wo.ID).
			SetNewValue(pri.MustGetString()).
			Save(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateActivitiesOnWOUpdate(ctx context.Context, entry *LogEntry) error {
	userID := entry.UserID
	client := ent.FromContext(ctx)

	newVal, oldVal, shouldUpdate := getDiffOfUniqueEdgeAsString(entry, workorder.EdgeAssignee)
	if shouldUpdate {
		_, err := client.Activity.Create().
			SetChangedField(activity.ChangedFieldASSIGNEE).
			SetIsCreate(false).
			SetNillableAuthorID(userID).
			SetWorkOrderID(entry.CurrState.ID).
			SetNillableOldValue(oldVal).
			SetNillableNewValue(newVal).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	newVal, oldVal, shouldUpdate = getDiffOfUniqueEdgeAsString(entry, workorder.EdgeOwner)
	if shouldUpdate {
		_, err := client.Activity.Create().
			SetChangedField(activity.ChangedFieldOWNER).
			SetIsCreate(false).
			SetNillableAuthorID(userID).
			SetWorkOrderID(entry.CurrState.ID).
			SetNillableOldValue(oldVal).
			SetNillableNewValue(newVal).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	newVal, oldVal, shouldUpdate = getStringDiffValuesField(entry, workorder.FieldStatus)
	if shouldUpdate {
		_, err := client.Activity.Create().
			SetChangedField(activity.ChangedFieldSTATUS).
			SetIsCreate(false).
			SetNillableAuthorID(userID).
			SetWorkOrderID(entry.CurrState.ID).
			SetNillableOldValue(oldVal).
			SetNillableNewValue(newVal).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	newVal, oldVal, shouldUpdate = getStringDiffValuesField(entry, workorder.FieldPriority)
	if shouldUpdate {
		_, err := client.Activity.Create().
			SetChangedField(activity.ChangedFieldPRIORITY).
			SetIsCreate(false).
			SetNillableAuthorID(userID).
			SetWorkOrderID(entry.CurrState.ID).
			SetNillableOldValue(oldVal).
			SetNillableNewValue(newVal).
			Save(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func getDiffOfUniqueEdgeAsString(entry *LogEntry, edge string) (*string, *string, bool) {
	newIntVal, oldIntVal, shouldUpdate := getDiffOfUniqueEdge(entry, edge)
	var newStrVal, oldStrVal *string
	if newIntVal != nil {
		newStrVal = pointer.ToString(strconv.Itoa(*newIntVal))
	}
	if oldIntVal != nil {
		oldStrVal = pointer.ToString(strconv.Itoa(*oldIntVal))
	}
	return newStrVal, oldStrVal, shouldUpdate
}
