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

	if assignee, found := findEdge(wo.Edges, toCamelCase(workorder.EdgeAssignee)); found {
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

	if owner, found := findEdge(wo.Edges, toCamelCase(workorder.EdgeOwner)); found {
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

	if st, found := findField(wo.Fields, toCamelCase(workorder.FieldStatus)); found {
		status := st.Value
		_, err = client.Activity.Create().
			SetChangedField(activity.ChangedFieldSTATUS).
			SetIsCreate(true).
			SetNillableAuthorID(userID).
			SetWorkOrderID(wo.ID).
			SetNewValue(status).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	if prio, found := findField(wo.Fields, toCamelCase(workorder.FieldPriority)); found {
		pri := prio.Value
		_, err = client.Activity.Create().
			SetChangedField(activity.ChangedFieldPRIORITY).
			SetIsCreate(true).
			SetNillableAuthorID(userID).
			SetWorkOrderID(wo.ID).
			SetNewValue(pri).
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

	newVal, oldVal, shouldUpdate = getDiffValuesField(entry, workorder.FieldStatus)
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

	newVal, oldVal, shouldUpdate = getDiffValuesField(entry, workorder.FieldPriority)
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
