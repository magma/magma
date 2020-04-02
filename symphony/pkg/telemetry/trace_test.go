// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opencensus.io/trace"
)

type traceIniter struct {
	mock.Mock
}

func (ti *traceIniter) Init(opts TraceExporterOptions) (trace.Exporter, error) {
	args := ti.Called(opts)
	exporter, _ := args.Get(0).(trace.Exporter)
	return exporter, args.Error(1)
}

func TestGetTraceExporter(t *testing.T) {
	_, err := GetTraceExporter("noexist", TraceExporterOptions{})
	assert.EqualError(t, err, `trace exporter "noexist" not found`)
	var ti traceIniter
	ti.On("Init", mock.Anything).Return(nil, nil).Once()
	defer ti.AssertExpectations(t)
	assert.NotPanics(t, func() { MustRegisterTraceExporter(t.Name(), ti.Init) })
	defer traceExporters.Delete(t.Name())
	_, err = GetTraceExporter(t.Name(), TraceExporterOptions{})
	assert.NoError(t, err)
}

func TestAvailableTraceExporters(t *testing.T) {
	var ti traceIniter
	defer ti.AssertExpectations(t)
	suffixes := []string{"foo", "bar", "baz"}
	for _, suffix := range suffixes {
		err := RegisterTraceExporter(t.Name()+suffix, ti.Init)
		assert.NoError(t, err)
	}
	defer func() {
		for _, suffix := range suffixes {
			traceExporters.Delete(t.Name() + suffix)
		}
	}()
	assert.Panics(t, func() {
		MustRegisterTraceExporter(t.Name()+suffixes[0], ti.Init)
	})
	exporters := AvailableTraceExporters()
	assert.True(t, sort.IsSorted(sort.StringSlice(exporters)))
	for _, suffix := range suffixes {
		assert.Contains(t, exporters, t.Name()+suffix)
	}
}

func TestWithoutNameSampler(t *testing.T) {
	sampler := WithoutNameSampler("foo", "bar")
	decision := sampler(trace.SamplingParameters{Name: "foo"})
	assert.False(t, decision.Sample)
	decision = sampler(trace.SamplingParameters{Name: "bar"})
	assert.False(t, decision.Sample)
	decision = sampler(trace.SamplingParameters{Name: "baz"})
	assert.True(t, decision.Sample)
}
