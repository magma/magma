// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package actions

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/facebookincubator/symphony/pkg/actions/executor"
)

// Client provides a client interface to the actions framework.
// It essentially provides a sub-interface of Executor/Registry
type Client struct {
	executor *executor.Executor
}

// NewClient returns a new executor client
func NewClient(exec *executor.Executor) *Client {
	return &Client{exec}
}

// ActionForID delegates to executor.Registry.ActionForID
func (c *Client) ActionForID(str core.ActionID) (core.Action, error) {
	return c.executor.Registry.ActionForID(str)
}

// TriggerForID delegates to executor.Registry.TriggerForID
func (c *Client) TriggerForID(str core.TriggerID) (core.Trigger, error) {
	return c.executor.Registry.TriggerForID(str)
}

// Triggers delegates to executor.Registry.Triggers
func (c *Client) Triggers() []core.Trigger {
	return c.executor.Registry.Triggers()
}

// Execute delegates to executor.Execute
func (c *Client) Execute(ctx context.Context, objectID string, triggerToPayload map[core.TriggerID]map[string]interface{}) executor.ExecutionResult {
	return c.executor.Execute(ctx, objectID, triggerToPayload)
}

type contextKey struct{}

// FromContext returns an executor stored in a context
func FromContext(ctx context.Context) *Client {
	e, ok := ctx.Value(contextKey{}).(*executor.Executor)
	if !ok {
		return nil
	}
	return &Client{e}
}

// NewContext returns a new context with the given Executor attached.
func NewContext(parent context.Context, executor *executor.Executor) context.Context {
	return context.WithValue(parent, contextKey{}, executor)
}
