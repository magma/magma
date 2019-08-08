/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggerFieldContext(t *testing.T) {
	core, o := observer.New(zap.InfoLevel)
	logger := NewFactory(zap.New(core))

	ctx := NewFieldsContext(context.Background(), zap.String("name", "test"))
	ctx = NewFieldsContext(ctx, zap.String("lang", "go"))
	logger.For(ctx).Info("test message", zap.Int("speed", 42))

	assert.Equal(t, 1, o.
		FilterMessage("test message").
		FilterField(zap.String("name", "test")).
		FilterField(zap.String("lang", "go")).
		Len(),
	)
}
