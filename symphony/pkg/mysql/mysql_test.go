// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mysql

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
)

type testExporter struct {
	spans []*trace.SpanData
}

func (e *testExporter) ExportSpan(s *trace.SpanData) {
	e.spans = append(e.spans, s)
}

func TestOpen(t *testing.T) {
	dsn, ok := os.LookupEnv("MYSQL_DSN")
	if !ok {
		t.Skip("MYSQL_DSN not provided")
	}
	e := &testExporter{}
	trace.RegisterExporter(e)
	defer trace.UnregisterExporter(e)

	db := Open(dsn)
	require.NotNil(t, db)
	ctx, span := trace.StartSpan(context.Background(), "test",
		trace.WithSampler(trace.AlwaysSample()),
	)
	err := db.PingContext(ctx)
	assert.NoError(t, err)
	span.End()
	assert.Len(t, e.spans, 2)
	err = db.Close()
	assert.NoError(t, err)
}
