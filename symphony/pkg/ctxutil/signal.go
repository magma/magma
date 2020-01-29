// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ctxutil

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
)

// ErrSignal is returned when a signal is received.
type ErrSignal struct {
	Signal os.Signal
}

// Error implements error interface.
func (e ErrSignal) Error() string {
	return "received signal: " + e.Signal.String()
}

// WithSignal returns a context which is done when an OS signal is received.
func WithSignal(parent context.Context, sig ...os.Signal) context.Context {
	return withSignal(parent, signalNotifier{}, sig...)
}

func withSignal(parent context.Context, notifier notifier, sig ...os.Signal) context.Context {
	s := &signalCtx{
		Context: parent,
		done:    make(chan struct{}),
	}
	sigCh := make(chan os.Signal, 1)
	notifier.Notify(sigCh, sig...)
	runtime.SetFinalizer(s, func(*signalCtx) { notifier.Stop(sigCh) })
	go s.watch(sigCh)
	return s
}

type signalCtx struct {
	context.Context
	done chan struct{}
	err  atomic.Value
}

func (s *signalCtx) Done() <-chan struct{} {
	return s.done
}

func (s *signalCtx) Err() error {
	if err := s.err.Load(); err != nil {
		return err.(ErrSignal)
	}
	return s.Context.Err()
}

func (s *signalCtx) watch(sigCh chan os.Signal) {
	select {
	case <-s.Context.Done():
	case sig := <-sigCh:
		s.err.Store(ErrSignal{sig})
	}
	close(s.done)
}

type notifier interface {
	Notify(chan<- os.Signal, ...os.Signal)
	Stop(chan<- os.Signal)
}

type signalNotifier struct{}

func (signalNotifier) Notify(c chan<- os.Signal, sig ...os.Signal) {
	signal.Notify(c, sig...)
}

func (signalNotifier) Stop(c chan<- os.Signal) {
	signal.Stop(c)
}
