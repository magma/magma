// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/wire"
	"gopkg.in/alecthomas/kingpin.v2"
)

type URL string

// Config configures this package.
type Config struct {
	pubURL URL
	subURL URL
}

// AddFlagsVar adds the flags used by this package to the Kingpin application.
func AddFlagsVar(a *kingpin.Application, config *Config) {
	a.Flag("event.pub-url", "events pub url").
		Envar("EVENT_PUB_URL").
		Default("mem://events").
		SetValue(&config.pubURL)
	a.Flag("event.sub-url", "events sub url").
		Envar("EVENT_SUB_URL").
		Default("mem://events").
		SetValue(&config.subURL)
}

// String returns the textual representation of a config.
func (u URL) String() string {
	return string(u)
}

// Set updates the value of the config.
func (u *URL) Set(v string) error {
	if _, err := url.Parse(v); err != nil {
		return fmt.Errorf("parsing url: %w", err)
	}
	*u = URL(v)
	return nil
}

func newConfig(url string) Config {
	return Config{
		pubURL: URL(url),
		subURL: URL(url),
	}
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
	emitter, err := NewTopicEmitter(ctx, cfg.pubURL.String())
	if err != nil {
		return nil, nil, err
	}
	return emitter, func() { _ = emitter.Shutdown(ctx) }, nil
}

// ProvideEmitter providers subscriber from config.
func ProvideSubscriber(cfg Config) URLSubscriber {
	return NewURLSubscriber(cfg.subURL.String())
}
