// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ctxutil

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignal(t *testing.T) {
	ctxdone := func(ctx context.Context) bool {
		timer := time.NewTimer(100 * time.Millisecond)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return true
		case <-timer.C:
			return false
		}
	}
	tests := []struct {
		name   string
		sig    []os.Signal
		notify func(func(os.Signal), context.CancelFunc)
		expect func(*testing.T, context.Context)
	}{
		{
			name: "Simple",
			notify: func(signal func(os.Signal), _ context.CancelFunc) {
				signal(os.Interrupt)
			},
			expect: func(t *testing.T, ctx context.Context) {
				assert.True(t, ctxdone(ctx))
				assert.EqualError(t, ctx.Err(), ErrSignal{os.Interrupt}.Error())
			},
		},
		{
			name: "ParentCanceled",
			notify: func(_ func(os.Signal), cancel context.CancelFunc) {
				cancel()
			},
			expect: func(t *testing.T, ctx context.Context) {
				assert.True(t, ctxdone(ctx))
				assert.EqualError(t, ctx.Err(), context.Canceled.Error())
			},
		},
		{
			name:   "NoError",
			notify: func(func(os.Signal), context.CancelFunc) {},
			expect: func(t *testing.T, ctx context.Context) {
				assert.False(t, ctxdone(ctx))
				assert.NoError(t, ctx.Err())
			},
		},
		{
			name: "MaskedSig",
			sig:  []os.Signal{syscall.SIGTERM},
			notify: func(signal func(os.Signal), _ context.CancelFunc) {
				signal(syscall.SIGINT)
			},
			expect: func(t *testing.T, ctx context.Context) {
				assert.False(t, ctxdone(ctx))
				assert.NoError(t, ctx.Err())
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			notifier := &testNotifier{}
			ctx, cancel := context.WithCancel(context.Background())
			ctx = withSignal(ctx, notifier, tc.sig...)
			tc.notify(notifier.signal, cancel)
			tc.expect(t, ctx)
		})
	}
}

type testNotifier struct {
	c   chan<- os.Signal
	sig []os.Signal
}

func (n *testNotifier) Notify(c chan<- os.Signal, sig ...os.Signal) {
	n.c = c
	n.sig = sig
}

func (testNotifier) Stop(chan<- os.Signal) {}

func (n *testNotifier) signal(sig os.Signal) {
	if len(n.sig) == 0 {
		n.c <- sig
		return
	}
	for _, s := range n.sig {
		if s == sig {
			n.c <- sig
			return
		}
	}
}
