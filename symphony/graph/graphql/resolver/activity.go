// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"strconv"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/activity"
)

type activityResolver struct{}

func (a activityResolver) Author(ctx context.Context, obj *ent.Activity) (*ent.User, error) {
	author, err := obj.QueryAuthor().Only(ctx)
	return author, ent.MaskNotFound(err)
}

func getNode(ctx context.Context, field activity.ChangedField, val string) (ent.Noder, error) {
	if val == "" {
		return nil, nil
	}
	switch field {
	case activity.ChangedFieldASSIGNEE:
		fallthrough
	case activity.ChangedFieldOWNER:
		client := ent.FromContext(ctx)
		intID, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		return client.Noder(ctx, intID)
	}
	return nil, nil
}

func (a activityResolver) NewRelatedNode(ctx context.Context, obj *ent.Activity) (ent.Noder, error) {
	return getNode(ctx, obj.ChangedField, obj.NewValue)
}

func (a activityResolver) OldRelatedNode(ctx context.Context, obj *ent.Activity) (ent.Noder, error) {
	return getNode(ctx, obj.ChangedField, obj.OldValue)
}

func (a activityResolver) WorkOrder(ctx context.Context, obj *ent.Activity) (*ent.WorkOrder, error) {
	return obj.QueryWorkOrder().Only(ctx)
}
