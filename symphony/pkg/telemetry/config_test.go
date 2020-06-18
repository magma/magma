// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry_test

import (
	"os"
	"testing"

	"github.com/facebookincubator/symphony/pkg/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/alecthomas/kingpin.v2"
)

func TestFlags(t *testing.T) {
	a := kingpin.New(t.Name(), "")
	c := telemetry.AddFlags(a)
	_, err := a.Parse([]string{
		"--" + telemetry.TraceExporterFlag, "nop",
		"--" + telemetry.TraceSamplingProbabilityFlag, "0.5",
		"--" + telemetry.TraceServiceFlag, t.Name(),
		"--" + telemetry.TraceTagsFlag, "one:1",
		"--" + telemetry.TraceTagsFlag, "two=2",
		"--" + telemetry.ViewExporterFlag, "nop",
		"--" + telemetry.ViewLabelsFlag, "three:3",
	})
	require.NoError(t, err)
	assert.Equal(t, "nop", c.Trace.ExporterName)
	assert.Equal(t, 0.5, c.Trace.SamplingProbability)
	assert.Equal(t, t.Name(), c.Trace.ServiceName)
	assert.Equal(t, map[string]string{"one": "1", "two": "2"}, c.Trace.Tags)
	assert.Equal(t, "nop", c.View.ExporterName)
	assert.Equal(t, map[string]string{"three": "3"}, c.View.Labels)
}

func TestEnvarFlags(t *testing.T) {
	vars := map[string]string{
		telemetry.TraceExporterEnvar:            "nop",
		telemetry.TraceSamplingProbabilityEnvar: "0.2",
		telemetry.TraceServiceEnvar:             t.Name(),
		telemetry.ViewExporterEnvar:             "nop",
	}
	for key, value := range vars {
		err := os.Setenv(key, value)
		require.NoError(t, err)
	}
	defer func() {
		for key := range vars {
			os.Unsetenv(key)
		}
	}()
	a := kingpin.New(t.Name(), "")
	c := telemetry.AddFlags(a)
	_, err := a.Parse(nil)
	require.NoError(t, err)
	assert.Equal(t, "nop", c.Trace.ExporterName)
	assert.Equal(t, 0.2, c.Trace.SamplingProbability)
	assert.Equal(t, t.Name(), c.Trace.ServiceName)
	assert.Equal(t, "nop", c.View.ExporterName)
}

func TestProvider(t *testing.T) {
	err := os.Setenv("JAEGER_AGENT_ENDPOINT", "localhost:6831")
	require.NoError(t, err)
	defer os.Unsetenv("JAEGER_AGENT_ENDPOINT")
	a := kingpin.New(t.Name(), "")
	c := telemetry.AddFlags(a)
	_, err = a.Parse([]string{
		"--" + telemetry.TraceExporterFlag, "jaeger",
		"--" + telemetry.ViewExporterFlag, "prometheus",
	})
	require.NoError(t, err)
	te, flusher, err := telemetry.ProvideTraceExporter(c)
	require.NoError(t, err)
	assert.NotNil(t, te)
	assert.NotNil(t, flusher)
	sampler := telemetry.ProvideTraceSampler(c)
	assert.NotNil(t, sampler)
	ve, err := telemetry.ProvideViewExporter(c)
	require.NoError(t, err)
	assert.NotNil(t, ve)
}
