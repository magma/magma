// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
)

type subscriptionResolver struct{ resolver }

func (subscriptionResolver) WorkOrderAdded(ctx context.Context) (<-chan *ent.WorkOrder, error) {
	events := make(chan *ent.WorkOrder, 1)
	go func() {
		defer close(events)
		<-ctx.Done()
	}()
	return events, nil
}
