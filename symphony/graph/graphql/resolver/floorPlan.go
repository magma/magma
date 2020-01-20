// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
)

type floorPlanResolver struct{}

func (floorPlanResolver) LocationID(ctx context.Context, obj *ent.FloorPlan) (string, error) {
	return obj.QueryLocation().FirstID(ctx)
}

func (floorPlanResolver) Image(ctx context.Context, obj *ent.FloorPlan) (*ent.File, error) {
	return obj.QueryImage().First(ctx)
}

func (floorPlanResolver) ReferencePoint(ctx context.Context, obj *ent.FloorPlan) (*ent.FloorPlanReferencePoint, error) {
	return obj.QueryReferencePoint().First(ctx)
}

func (floorPlanResolver) Scale(ctx context.Context, obj *ent.FloorPlan) (*ent.FloorPlanScale, error) {
	return obj.QueryScale().First(ctx)
}
