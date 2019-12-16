// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package work

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
)

func TestArgsString(t *testing.T) {
	args := Args{"name": "foo", "age": 42}
	assert.Implements(t, (*fmt.Stringer)(nil), args)
	assert.JSONEq(t, `{"name":"foo","age":42}`, fmt.Sprint(args))
}

func TestJobString(t *testing.T) {
	job := Job{Handler: "bar", Args: Args{"name": "baz"}}
	assert.Implements(t, (*fmt.Stringer)(nil), job)
	assert.JSONEq(t, `{"handler":"bar","args":{"name": "baz"}}`, fmt.Sprint(job))
}

func TestJobObjectMarshal(t *testing.T) {
	job := Job{Handler: "h"}
	assert.Implements(t, (*zapcore.ObjectMarshaler)(nil), job)
	core, o := observer.New(zapcore.InfoLevel)
	logger := zaptest.NewLogger(t, zaptest.WrapOptions(
		zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(c, core)
		}),
	))
	logger.Info("logging job", zap.Object("job", job))
	entries := o.TakeAll()
	require.Len(t, entries, 1)
	field, ok := entries[0].ContextMap()["job"]
	require.True(t, ok)
	assert.Contains(t, field, "handler")
	assert.Contains(t, field, "args")
}
