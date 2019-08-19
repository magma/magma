/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package logtest

import (
	"context"
	"testing"

	"fbc/lib/go/log"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type testFactory struct {
	*zap.Logger
}

// NewFactory creates a new testing logger factory
func NewFactory(t *testing.T) log.Factory {
	return testFactory{zaptest.NewLogger(t)}
}

func (f testFactory) Bg() *zap.Logger                 { return f.Logger }
func (f testFactory) For(context.Context) *zap.Logger { return f.Logger }

func (f testFactory) With(fields ...zap.Field) log.Factory {
	return testFactory{f.Logger.With(fields...)}
}
