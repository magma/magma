// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub

import (
	"context"
	"fmt"

	"gocloud.dev/pubsub"
)

const (
	// TenantHeader is the metadata key holding tenant name.
	TenantHeader = "event/tenant"
	// NameHeader is the metadata key holding event name.
	NameHeader = "event/name"
)

// Emitter represents types than can emit events.
type Emitter interface {
	Emit(context.Context, string, string, []byte) error
}

// The EmitterFunc type is an adapter to allow the use of
// ordinary functions as event emitters.
type EmitterFunc func(context.Context, string, string, []byte) error

// Emit returns f(ctx, tenant, name, body).
func (f EmitterFunc) Emit(ctx context.Context, tenant, name string, body []byte) error {
	return f(ctx, tenant, name, body)
}

// TopicEmitter emits events to a pubsub topic.
type TopicEmitter struct {
	topic *pubsub.Topic
	url   string
}

// NewTopicEmitter creates an event emitter that writes to a pubsub topic.
func NewTopicEmitter(ctx context.Context, url string) (*TopicEmitter, error) {
	topic, err := pubsub.OpenTopic(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("opening topic: %w", err)
	}
	return &TopicEmitter{
		topic: topic,
		url:   url,
	}, nil
}

// Emit an event to a pubsub topic.
func (e TopicEmitter) Emit(ctx context.Context, tenant, name string, body []byte) (err error) {
	if err := e.topic.Send(ctx, &pubsub.Message{
		Metadata: map[string]string{
			TenantHeader: tenant,
			NameHeader:   name,
		},
		Body: body,
	}); err != nil {
		return fmt.Errorf("emitting event: %w", err)
	}
	return nil
}

// Shutdown shuts down event emitter.
func (e TopicEmitter) Shutdown(ctx context.Context) error {
	return e.topic.Shutdown(ctx)
}

// String returns the textual representation of topic emitter.
func (e TopicEmitter) String() string {
	return e.url
}

// Set updates the value of the topic emitter.
func (e *TopicEmitter) Set(url string) (err error) {
	if e.topic, err = pubsub.OpenTopic(context.Background(), url); err != nil {
		return fmt.Errorf("opening topic: %w", err)
	}
	e.url = url
	return nil
}

// NewNopEmitter returns an emitter that drops all events and never fails.
func NewNopEmitter() Emitter {
	return EmitterFunc(func(context.Context, string, string, []byte) error {
		return nil
	})
}
