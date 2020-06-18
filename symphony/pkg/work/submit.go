// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package work

import (
	"context"
	"encoding/json"
	"fmt"

	"gocloud.dev/pubsub"
)

// Submitter is a work submitter over pubsub topic.
type Submitter struct {
	topic *pubsub.Topic
}

// NewSubmitter creates a new work submitter.
func NewSubmitter(topic *pubsub.Topic) *Submitter {
	return &Submitter{topic}
}

// NewSubmitterURL creates a new work submitter from topic address.
func NewSubmitterURL(ctx context.Context, url string) (*Submitter, error) {
	topic, err := pubsub.OpenTopic(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("opening topic: %w", err)
	}
	return NewSubmitter(topic), nil
}

// Submit implements work submitter interface.
func (s *Submitter) Submit(ctx context.Context, job Job) (err error) {
	var msg pubsub.Message
	if msg.Body, err = json.Marshal(job); err != nil {
		return fmt.Errorf("json encoding job: %w", err)
	}
	// TODO: handle trace context propagation
	return s.topic.Send(ctx, &msg)
}

// Shutdown flushes pending jobs and disconnects the Topic.
func (s *Submitter) Shutdown(ctx context.Context) error {
	return s.topic.Shutdown(ctx)
}

// Close invokes Shutdown with background context.
func (s *Submitter) Close() error {
	return s.Shutdown(context.Background())
}
