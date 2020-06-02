// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub

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

// Set is a wire provider that provides an emitter/subscriber
// given context and config.
var Set = wire.NewSet(
	ProvideEmitter,
	ProvideSubscriber,
	wire.Bind(new(Emitter), new(*TopicEmitter)),
	wire.Bind(new(Subscriber), new(URLSubscriber)),
)

// ProvideEmitter providers emitter from config.
func ProvideEmitter(ctx context.Context, cfg Config) (*TopicEmitter, func(), error) {
	emitter, err := NewTopicEmitter(ctx, cfg.url)
	if err != nil {
		return nil, nil, err
	}
	return emitter, func() { _ = emitter.Shutdown(ctx) }, nil
}

// ProvideEmitter providers subscriber from config.
func ProvideSubscriber(cfg Config) URLSubscriber {
	return NewURLSubscriber(cfg.url)
}
