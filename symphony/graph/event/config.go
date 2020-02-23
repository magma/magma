// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/wire"
)

// Config configures this package.
type Config struct {
	url string
}

// String returns the textual representation of a config.
func (c Config) String() string {
	return c.url
}

// Set updates the value of the config.
func (c *Config) Set(v string) error {
	if _, err := url.Parse(v); err != nil {
		return fmt.Errorf("parsing url: %w", err)
	}
	c.url = v
	return nil
}

// UnmarshalFlag updates the value of the config.
func (c *Config) UnmarshalFlag(v string) error {
	return c.Set(v)
}

// Set is a wire provider that provides an emitter/subscriber
// given context and config.
var Set = wire.NewSet(
	Provide,
	wire.Bind(new(Emitter), new(TopicEmitter)),
	wire.Bind(new(Subscriber), new(URLSubscriber)),
)

// Provide event emitter / subscriber from config.
func Provide(ctx context.Context, cfg Config) (*TopicEmitter, *URLSubscriber, error) {
	emitter, err := NewTopicEmitter(ctx, cfg.url)
	if err != nil {
		return nil, nil, err
	}
	subscriber := URLSubscriber(cfg.url)
	return emitter, &subscriber, nil
}
