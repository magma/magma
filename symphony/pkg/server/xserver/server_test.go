// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xserver

import (
	"regexp"
	"testing"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestViews(t *testing.T) {
	r := regexp.MustCompile(
		`http_(request|response)_.*(total|bytes|seconds)`,
	)
	for _, v := range DefaultViews() {
		assert.Regexp(t, r, v.Name)
	}
}

func TestDefaultPrometheusRegisterer(t *testing.T) {
	assert.IsType(t, (*prometheus.Registry)(nil), prometheus.DefaultRegisterer)
}

func TestNewJaegerExporter(t *testing.T) {
	logger := log.NewNopLogger()
	t.Run("Simple", func(t *testing.T) {
		exporter, cleaner, err := NewJaegerExporter(logger, jaeger.Options{
			AgentEndpoint: "localhost:6831",
			Process: jaeger.Process{
				ServiceName: "test",
			},
		})
		assert.NotNil(t, exporter)
		assert.NotNil(t, cleaner)
		assert.NoError(t, err)
	})
	t.Run("Empty", func(t *testing.T) {
		exporter, cleaner, err := NewJaegerExporter(logger, jaeger.Options{})
		assert.Nil(t, exporter)
		assert.NotNil(t, cleaner)
		assert.NoError(t, err)
	})
}
