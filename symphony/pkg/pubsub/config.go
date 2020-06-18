// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub

import (
	"context"
	"net/url"

	"github.com/google/wire"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Config configures this package.
type Config struct {
	PubURL *url.URL
	SubURL *url.URL
}

// AddFlagsVar adds the flags used by this package to the Kingpin application.
func AddFlagsVar(a *kingpin.Application, config *Config) {
	a.Flag("event.pub-url", "events pub url").
		Envar("EVENT_PUB_URL").
		Default("mem://events").
		URLVar(&config.PubURL)
	a.Flag("event.sub-url", "events sub url").
		Envar("EVENT_SUB_URL").
		Default("mem://events").
		URLVar(&config.SubURL)
}

// AddFlags adds the flags used by this package to the Kingpin application.
func AddFlags(a *kingpin.Application) *Config {
	config := &Config{}
	AddFlagsVar(a, config)
	return config
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
	emitter, err := NewTopicEmitter(ctx, cfg.PubURL.String())
	if err != nil {
		return nil, nil, err
	}
	return emitter, func() { _ = emitter.Shutdown(ctx) }, nil
}

// ProvideEmitter providers subscriber from config.
func ProvideSubscriber(cfg Config) URLSubscriber {
	return NewURLSubscriber(cfg.SubURL.String())
}
