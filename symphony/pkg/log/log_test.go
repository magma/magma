// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestDefaultLogger(t *testing.T) {
	core, o := observer.New(zap.InfoLevel)
	logger := NewDefaultLogger(zap.New(core))
	assert.Implements(t, (*Logger)(nil), logger)

	msg := "context-less"
	logger.Background().Info(msg)
	assert.Len(t, o.TakeAll(), 1)

	ctx := context.Background()
	logger.For(ctx).Warn(msg)
	assert.Equal(t, 1, o.FilterMessage(msg).Len())

	exporter := &testExporter{}
	trace.RegisterExporter(exporter)

	ctx, span := trace.StartSpan(
		context.Background(), "test",
		trace.WithSampler(trace.AlwaysSample()),
	)
	field, msg := zap.Int("result", 42), "context-aware"
	logger.For(ctx).Info(msg, field)
	span.End()

	assert.Equal(t, 1, o.FilterField(field).FilterMessage(msg).Len())
	spans := exporter.spans
	require.Len(t, spans, 1)
	annotations := spans[0].Annotations
	require.Len(t, annotations, 1)
	assert.Equal(t, msg, annotations[0].Message)
	assert.EqualValues(t, 42, annotations[0].Attributes["result"])
}

func TestNopFactory(t *testing.T) {
	logger := NewNopLogger()
	assert.Implements(t, (*Logger)(nil), logger)
	assert.EqualValues(t, zap.NewNop(), logger.Background())
	assert.EqualValues(t, zap.NewNop(), logger.For(context.Background()))
}
