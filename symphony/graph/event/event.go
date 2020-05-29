// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"sync"
	"time"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/hashicorp/go-multierror"
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
	client.Use(e.logHook())
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
		logger.Debug("emitting event")
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

type (
	// PipeEmitter is the emit half of the pipe.
	PipeEmitter = TopicEmitter

	// PipeSubscriber is the subscribe half of the pipe.
	PipeSubscriber struct {
		mu    sync.Mutex
		topic *pubsub.Topic
		subs  []*pubsub.Subscription
	}
)

// Subscribe implements Subscriber interface.
func (s *PipeSubscriber) Subscribe(context.Context) (*pubsub.Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.topic == nil {
		return nil, errors.New("event: subscriber pipe has been shutdown")
	}
	sub := mempubsub.NewSubscription(s.topic, time.Second)
	s.subs = append(s.subs, sub)
	return sub, nil
}

func (s *PipeSubscriber) shutdown(ctx context.Context, subs []*pubsub.Subscription) error {
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs = &multierror.Error{}
	)
	wg.Add(len(subs))
	for _, sub := range subs {
		go func(sub *pubsub.Subscription) {
			defer wg.Done()
			if err := sub.Shutdown(ctx); err != nil {
				mu.Lock()
				errs = multierror.Append(errs, err)
				mu.Unlock()
			}
		}(sub)
	}
	wg.Wait()
	return errs.ErrorOrNil()
}

// Shutdown shuts down the pipe subscriber.
func (s *PipeSubscriber) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	subs := s.subs
	s.subs = nil
	s.topic = nil
	s.mu.Unlock()
	return s.shutdown(ctx, subs)
}

// Pipe creates a in memory emitter/subscriber pipe.
func Pipe() (*PipeEmitter, *PipeSubscriber) {
	topic := mempubsub.NewTopic()
	return &PipeEmitter{topic: topic}, &PipeSubscriber{topic: topic}
}

// DefaultViews are predefined views for opencensus metrics.
var DefaultViews = pubsub.OpenCensusViews
