// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"time"

	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/mempubsub"
)

const (
	// TenantHeader is the metadata key holding tenant name.
	TenantHeader = "event/tenant"
	// NameHeader is the metadata key holding event name.
	NameHeader = "event/name"
)

// Work order events.
const (
	WorkOrderAdded = "work_order/added"
	WorkOrderDone  = "work_order/done"
)

// Pipe creates a in memory emitter/subscriber pipe.
func Pipe() (Emitter, Subscriber) {
	topic := mempubsub.NewTopic()
	subscriber := SubscriberFunc(
		func(context.Context) (*pubsub.Subscription, error) {
			return mempubsub.NewSubscription(topic, time.Second), nil
		},
	)
	return TopicEmitter{topic: topic}, subscriber
}

// Views are predefined views for opencensus metrics.
var Views = pubsub.OpenCensusViews
