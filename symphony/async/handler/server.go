// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"context"
	"fmt"
	"sync"

	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/pubsub"
	"github.com/facebookincubator/symphony/pkg/telemetry"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"go.uber.org/zap"
)

const (
	serviceName = "EventLogService"
)

// A Handler handles incoming events.
type Handler interface {
	Handle(context.Context, pubsub.LogEntry) error
}

// The Func type is an adapter to allow the use of
// ordinary functions as handlers.
type Func func(context.Context, pubsub.LogEntry) error

// Handle returns f(ctx, entry).
func (f Func) Handle(ctx context.Context, entry pubsub.LogEntry) error {
	return f(ctx, entry)
}

// NewServer is the events server.
type Server struct {
	tenancy    viewer.Tenancy
	subscriber pubsub.Subscriber
	logger     log.Logger
	handlers   []Handler
}

// Config defines the async server config.
type Config struct {
	Tenancy    viewer.Tenancy
	Logger     log.Logger
	Subscriber pubsub.Subscriber
	Telemetry  *telemetry.Config
}

func NewServer(cfg Config) *Server {
	return &Server{
		tenancy:    cfg.Tenancy,
		logger:     cfg.Logger,
		subscriber: cfg.Subscriber,
		handlers: []Handler{
			eventLog{
				logger: cfg.Logger,
			},
		},
	}
}

// Subscribe returns listener to the relevant events.
func (s *Server) Subscribe(ctx context.Context, wg *sync.WaitGroup) (*pubsub.Listener, error) {
	return pubsub.NewListener(ctx, pubsub.ListenerConfig{
		Subscriber: s.subscriber,
		Logger:     s.logger.For(ctx),
		Events:     []string{pubsub.EntMutation},
		Handler: pubsub.HandlerFunc(func(ctx context.Context, tenant string, _ string, body []byte) error {
			wg.Add(1)
			var entry pubsub.LogEntry
			err := pubsub.Unmarshal(body, &entry)
			if err != nil {
				wg.Done()
				return fmt.Errorf("cannot unmarshal log entry: %w", err)
			}
			go func() {
				defer wg.Done()
				err := s.handleEventLog(s.handlers)(context.Background(), tenant, entry)
				if err != nil {
					s.logger.For(ctx).Error("failed to handle event", zap.Error(err))
				}
			}()
			return nil
		}),
	})
}

func (s *Server) handleEventLog(handlers []Handler) func(context.Context, string, pubsub.LogEntry) error {
	return func(ctx context.Context, tenant string, entry pubsub.LogEntry) error {
		client, err := s.tenancy.ClientFor(ctx, tenant)
		if err != nil {
			const msg = "cannot get tenancy client"
			s.logger.For(ctx).Error(msg, zap.Error(err))
			return fmt.Errorf("%s. tenant: %s", msg, tenant)
		}
		ctx = ent.NewContext(ctx, client)
		v := viewer.NewAutomation(tenant, serviceName, user.RoleOWNER)
		ctx = log.NewFieldsContext(ctx, zap.Object("viewer", v))
		ctx = viewer.NewContext(ctx, v)
		permissions, err := authz.Permissions(ctx)
		if err != nil {
			const msg = "cannot get permissions"
			s.logger.For(ctx).Error(msg,
				zap.Error(err),
			)
			return fmt.Errorf("%s. tenant: %s, name: %s", msg, tenant, serviceName)
		}
		ctx = authz.NewContext(ctx, permissions)

		for _, h := range handlers {
			if err := s.runHandlerWithTransaction(ctx, h, entry); err != nil {
				s.logger.For(ctx).Error("running handler", zap.Error(err))
			}
		}
		return nil
	}
}
func (s *Server) runHandlerWithTransaction(ctx context.Context, h Handler, entry pubsub.LogEntry) error {
	tx, err := ent.FromContext(ctx).Tx(ctx)
	if err != nil {
		return fmt.Errorf("creating transaction: %w", err)
	}
	ctx = ent.NewTxContext(ctx, tx)
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
	}()
	ctx = ent.NewContext(ctx, tx.Client())
	if err := h.Handle(ctx, entry); err != nil {
		if r := tx.Rollback(); r != nil {
			err = fmt.Errorf("rolling back transaction: %v", r)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
