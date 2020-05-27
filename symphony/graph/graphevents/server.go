// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphevents

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/user"

	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/log"
	"go.uber.org/zap"
)

const (
	serviceName = "EventLogService"
)

// A Handler handles incoming events.
type Handler interface {
	Handle(context.Context, event.LogEntry) error
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as handlers.
type HandlerFunc func(context.Context, event.LogEntry) error

// Handle returns f(ctx, entry).
func (f HandlerFunc) Handle(ctx context.Context, entry event.LogEntry) error {
	return f(ctx, entry)
}

type serverConfig struct {
	tenancy    viewer.Tenancy
	logger     log.Logger
	subscriber event.Subscriber
}

// NewServer is the events server.
type Server struct {
	tenancy    viewer.Tenancy
	subscriber event.Subscriber
	logger     log.Logger
	handlers   []Handler
}

func newServer(cfg serverConfig) (*Server, func(), error) {
	return &Server{
		tenancy:    cfg.tenancy,
		logger:     cfg.logger,
		subscriber: cfg.subscriber,
		handlers: []Handler{
			eventLog{
				logger: cfg.logger,
			},
		},
	}, nil, nil
}

// Subscribe returns listener to the relevant events.
func (s *Server) Subscribe(ctx context.Context) (*event.Listener, error) {
	return event.NewListener(ctx, event.ListenerConfig{
		Subscriber: s.subscriber,
		Logger:     s.logger.For(ctx),
		Events:     []string{event.EntMutation},
		Handler:    s.handleEventLog(s.handlers),
	})
}

func (s *Server) handleEventLog(handlers []Handler) event.Handler {
	return event.HandlerFunc(func(ctx context.Context, tenant string, name string, body []byte) error {
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

		var entry event.LogEntry
		err = event.Unmarshal(body, &entry)
		if err != nil {
			const msg = "cannot unmarshal log entry"
			s.logger.For(ctx).Error(msg,
				zap.Error(err),
			)
			return fmt.Errorf("%s: %w", msg, err)
		}
		for _, h := range handlers {
			if err := s.runHandlerWithTransaction(ctx, h, entry); err != nil {
				return fmt.Errorf("running handler: %w", err)
			}
		}
		return nil
	})
}
func (s *Server) runHandlerWithTransaction(ctx context.Context, h Handler, entry event.LogEntry) error {
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
