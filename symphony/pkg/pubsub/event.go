// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/mempubsub"
)

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
