// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gocloud.dev/pubsub"
)

// A Handler handles incoming events.
type Handler interface {
	Handle(context.Context, string, string, []byte) error
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as event handlers.
type HandlerFunc func(context.Context, string, string, []byte) error

// Handle returns f(ctx, name, body).
func (f HandlerFunc) Handle(ctx context.Context, tenant string, name string, body []byte) error {
	return f(ctx, tenant, name, body)
}

type (
	// ListenerConfig configures event listener.
	ListenerConfig struct {
		Subscriber Subscriber
		Logger     *zap.Logger
		Tenant     *string
		Events     []string
		Handler    Handler
	}

	// Listener handles incoming events.
	Listener struct {
		subscription *pubsub.Subscription
		logger       *zap.Logger
		tenant       *string
		events       map[string]struct{}
		handler      Handler
	}
)

// NewListener creates an events listener from config.
func NewListener(ctx context.Context, cfg ListenerConfig) (*Listener, error) {
	if len(cfg.Events) == 0 {
		return nil, errors.New("events cannot be empty")
	}

	subscription, err := cfg.Subscriber.Subscribe(ctx)
	if err != nil {
		return nil, fmt.Errorf("opening subscription: %w", err)
	}

	if cfg.Logger == nil {
		cfg.Logger = zap.L()
	}
	events := make(map[string]struct{})
	for _, name := range cfg.Events {
		events[name] = struct{}{}
	}

	return &Listener{
		subscription: subscription,
		logger:       cfg.Logger,
		tenant:       cfg.Tenant,
		events:       events,
		handler:      cfg.Handler,
	}, nil
}

// Shutdown shuts down listener.
func (l *Listener) Shutdown(ctx context.Context) error {
	return l.subscription.Shutdown(ctx)
}

// Listen listens for incoming events.
func (l *Listener) Listen(ctx context.Context) error {
	for {
		if err := l.receive(ctx); err != nil {
			return err
		}
	}
}

func (l *Listener) receive(ctx context.Context) error {
	msg, err := l.subscription.Receive(ctx)
	if err != nil {
		return fmt.Errorf("receiving from subscription: %w", err)
	}
	defer msg.Ack()

	tenant, ok := msg.Metadata[TenantHeader]
	if !ok {
		l.logger.Warn("received event without tenant header")
		return nil
	}
	if l.tenant != nil && tenant != *l.tenant {
		return nil
	}

	name, ok := msg.Metadata[NameHeader]
	if !ok {
		l.logger.Warn("received event without name header")
		return nil
	}
	if _, ok := l.events[name]; !ok {
		return nil
	}

	if err := l.handler.Handle(ctx, tenant, name, msg.Body); err != nil {
		return fmt.Errorf("handling event %q: %w", name, err)
	}
	l.logger.Debug("handled event", zap.String("name", name))
	return nil
}

// SubscribeAndListen creates a listener and listens for incoming events.
func SubscribeAndListen(ctx context.Context, cfg ListenerConfig) error {
	listener, err := NewListener(ctx, cfg)
	if err != nil {
		return nil
	}
	defer listener.Shutdown(ctx)
	return listener.Listen(ctx)
}
