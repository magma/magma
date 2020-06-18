// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ctxutil

import (
	"context"
	"errors"
	"time"
)

// ErrDone is the error returned by Context.Err for done context.
var ErrDone = errors.New("context done")

// DoneCtx returns a non-nil, done Context. It is always canceled, has no
// values, and has no deadline.
func DoneCtx() context.Context {
	return doneCtx(0)
}

// closedchan is a reusable closed channel.
var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

type doneCtx int

func (doneCtx) Deadline() (deadline time.Time, ok bool) { return }
func (doneCtx) Done() <-chan struct{}                   { return closedchan }
func (doneCtx) Err() error                              { return ErrDone }
func (doneCtx) Value(interface{}) interface{}           { return nil }
