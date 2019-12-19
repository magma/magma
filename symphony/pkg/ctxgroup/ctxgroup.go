// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ctxgroup

import (
	"context"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type (
	// Group wraps errgroup.
	Group struct {
		wrapped *errgroup.Group
		ctx     context.Context
		sem     *semaphore.Weighted
	}

	// Option configures group.
	Option func(*Group)
)

// MaxConcurrency limits group concurrency.
func MaxConcurrency(n int64) Option {
	return func(g *Group) {
		if n <= 0 {
			panic("ctxgroup: concurrency must great than 0")
		}
		g.sem = semaphore.NewWeighted(n)
	}
}

// WithContext returns a new Group and an associated Context derived from ctx.
func WithContext(ctx context.Context, opts ...Option) *Group {
	var g Group
	g.wrapped, g.ctx = errgroup.WithContext(ctx)
	for _, opt := range opts {
		opt(&g)
	}
	return &g
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them. If Wait() is invoked
// after the context (originally supplied to WithContext) is canceled, Wait
// returns an error, even if no Go invocation did. In particular, calling
// Wait() after Done has been closed is guaranteed to return an error.
func (g *Group) Wait() error {
	ctxErr := g.ctx.Err()
	err := g.wrapped.Wait()
	if err != nil {
		return err
	}
	return ctxErr
}

// Go calls the given function in a new goroutine.
func (g *Group) Go(f func(context.Context) error) {
	g.wrapped.Go(func() error {
		if g.sem != nil {
			if err := g.sem.Acquire(g.ctx, 1); err != nil {
				return err
			}
			defer g.sem.Release(1)
		}
		return f(g.ctx)
	})
}
