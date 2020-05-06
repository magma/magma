// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/log"
	"go.uber.org/zap"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/mempubsub"
)

const (
	// TenantHeader is the metadata key holding tenant name.
	TenantHeader = "event/tenant"
	// NameHeader is the metadata key holding event name.
	NameHeader = "event/name"
)

// Eventer generates events from mutations.
type Eventer struct {
	Logger  log.Logger
	Emitter Emitter
}

// HookTo hooks eventer to ent client.
func (e *Eventer) HookTo(client *ent.Client) {
	client.WorkOrder.Use(e.workOrderHook())
}

func (e *Eventer) emit(ctx context.Context, name string, value interface{}) {
	emit := func(err error) {
		if err != nil {
			return
		}
		logger := e.Logger.For(ctx).With(zap.String("name", name))
		body, err := Marshal(value)
		if err != nil {
			logger.Warn("cannot marshal event value", zap.Error(err))
			return
		}
		if err := e.Emitter.Emit(ctx, viewer.FromContext(ctx).Tenant(), name, body); err != nil {
			logger.Warn("cannot emit event", zap.Error(err))
		}
	}
	if tx := ent.TxFromContext(ctx); tx != nil {
		tx.OnCommit(emit)
	} else {
		emit(nil)
	}
}

// Marshal returns the event encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal decodes event data into v.
func Unmarshal(data []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(v)
}

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

// DefaultViews are predefined views for opencensus metrics.
var DefaultViews = pubsub.OpenCensusViews
