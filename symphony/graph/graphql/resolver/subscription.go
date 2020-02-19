// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"time"

	"github.com/facebookincubator/symphony/graph/ent"
)

type subscriptionResolver struct{ resolver }

func (subscriptionResolver) WorkOrderAdded(ctx context.Context) (<-chan *ent.WorkOrder, error) {
	events := make(chan *ent.WorkOrder, 1)
	go func() {
		defer close(events)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				events <- &ent.WorkOrder{ID: "42"}
			case <-ctx.Done():
				return
			}
		}
	}()
	return events, nil
}
